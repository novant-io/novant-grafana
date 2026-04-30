package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// Datasource is the Novant data source plugin.
type Datasource struct {
	client     *Client
	pointCache *pointCache
}

// NewDatasource creates a new Novant data source instance.
func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	apiKey, ok := settings.DecryptedSecureJSONData["apiKey"]
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	return &Datasource{
		client:     NewClient(apiKey),
		pointCache: newPointCache(),
	}, nil
}

// Dispose cleans up resources.
func (d *Datasource) Dispose() {}

// CheckHealth validates the data source configuration.
func (d *Datasource) CheckHealth(_ context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	proj, err := d.client.GetProject()
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to connect: %v", err),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("Connected to %s (%s)", proj.ProjName, proj.City),
	}, nil
}

// CallResource handles HTTP calls to /api/datasources/uid/<uid>/resources/<path>.
// Used by the data source config UI to clear the in-memory point name cache.
func (d *Datasource) CallResource(_ context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	switch req.Path {
	case "clear-cache":
		if req.Method != http.MethodPost {
			return sender.Send(&backend.CallResourceResponse{
				Status: http.StatusMethodNotAllowed,
				Body:   []byte(`{"error":"method not allowed"}`),
			})
		}
		d.pointCache.clear()
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusOK,
			Body:   []byte(`{"status":"ok"}`),
		})
	default:
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
			Body:   []byte(`{"error":"not found"}`),
		})
	}
}

// QueryData handles multiple queries.
func (d *Datasource) QueryData(_ context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		response.Responses[q.RefID] = d.query(q)
	}
	return response, nil
}

func (d *Datasource) query(q backend.DataQuery) backend.DataResponse {
	var qm QueryModel
	if err := json.Unmarshal(q.JSON, &qm); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err))
	}

	switch q.QueryType {
	case "zones":
		return d.queryZones(qm)
	case "spaces":
		return d.querySpaces(qm)
	case "assets":
		return d.queryAssets(qm)
	case "sources":
		return d.querySources(qm)
	case "points":
		return d.queryPoints(qm)
	case "values":
		return d.queryValues(qm)
	case "trends":
		return d.queryTrends(q, qm)
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unknown query type: %s", q.QueryType))
	}
}

func (d *Datasource) queryZones(qm QueryModel) backend.DataResponse {
	resp, err := d.client.GetZones(qm.ZoneIDs)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	return backend.DataResponse{Frames: data.Frames{buildZonesFrame(resp)}}
}

func (d *Datasource) querySpaces(qm QueryModel) backend.DataResponse {
	resp, err := d.client.GetSpaces(qm.SpaceIDs)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	return backend.DataResponse{Frames: data.Frames{buildSpacesFrame(resp)}}
}

func (d *Datasource) queryAssets(qm QueryModel) backend.DataResponse {
	resp, err := d.client.GetAssets(qm.AssetIDs)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	return backend.DataResponse{Frames: data.Frames{buildAssetsFrame(resp)}}
}

func (d *Datasource) querySources(qm QueryModel) backend.DataResponse {
	resp, err := d.client.GetSources(qm.SourceIDs, qm.BoundOnly)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	return backend.DataResponse{Frames: data.Frames{buildSourcesFrame(resp)}}
}

func (d *Datasource) queryPoints(qm QueryModel) backend.DataResponse {
	resp, err := d.client.GetPoints(qm.SourceID, qm.AssetID, qm.SpaceID, qm.PointIDs, qm.PointTypes)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	return backend.DataResponse{Frames: data.Frames{buildPointsFrame(resp)}}
}

func (d *Datasource) queryValues(qm QueryModel) backend.DataResponse {
	resp, err := d.client.GetValues(qm.SourceID, qm.AssetID, qm.SpaceID, qm.PointIDs, qm.PointTypes)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	pointIDs := make([]string, len(resp.Values))
	for i, v := range resp.Values {
		pointIDs[i] = v.ID
	}
	names := d.pointCache.resolveNames(d.client, pointIDs)
	return backend.DataResponse{Frames: data.Frames{buildValuesFrame(resp, names)}}
}

func (d *Datasource) queryTrends(q backend.DataQuery, qm QueryModel) backend.DataResponse {
	if qm.PointIDs == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "point_ids is required for trends")
	}

	startDate := q.TimeRange.From.Format("2006-01-02")
	endDate := q.TimeRange.To.Format("2006-01-02")

	resp, err := d.client.GetTrends(qm.PointIDs, startDate, endDate, qm.Interval, qm.Aggregate)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}

	names := d.pointCache.resolveNames(d.client, resp.PointIDs)

	frames, err := buildTrendsFrames(resp, names)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}

	return backend.DataResponse{Frames: frames}
}
