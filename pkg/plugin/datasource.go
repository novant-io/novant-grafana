//
// Copyright (c) 2023, Novant LLC
// Licensed under the MIT License
//
// History:
//   4 Dec 2023  Andy Frank  Creation
//

package plugin

import (
	"context"
	"encoding/json"
  "errors"
  "math"
	"math/rand"
  "net/url"
  "strings"
  "time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, _ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &Datasource{}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct{}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct{}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

  // unmarshal the JSON params
  params := map[string]interface{}{}
  response.Error = json.Unmarshal(query.JSON, &params)
  if response.Error != nil {
    return response
  }
  start := toMidnight(query.TimeRange.From)
  end   := toMidnight(query.TimeRange.To).Add(24 * time.Hour)
  sourceId := strings.TrimSpace(params["sourceId"].(string))
  pointIds := strings.TrimSpace(params["pointIds"].(string))

  // validate
  if sourceId == "" {
    response.Error = errors.New("Missing sourceId value")
    return response
  }

  // validate
  if pointIds == "" {
    response.Error = errors.New("Missing pointIds value")
    return response
  }

  // stub out working data structures
  tss  := []time.Time{}
  vals := [][]float64{}
  pids := strings.Split(pointIds, ",")
  for _ = range pids {
    vals = append(vals, []float64{})
  }

  // request /points to get point meta
  pointsReq, err := novantReq(pCtx, "points", url.Values{ "source_id": {sourceId}})
  if err != nil {
    log.DefaultLogger.Error("/points request failed", "error", err)
    response.Error = err
    return response
  }
  pnames := []string{}
  points := filterPointMeta(pointsReq, pids)
  for i := range pids {
    id := pids[i]
    p  := points[id].(Map)
    pnames = append(pnames, p["name"].(string))
  }

  // iterate from time range
  cur := start
  for cur.Before(end) {
    // query /trends for data
    args := url.Values{
      "source_id": {sourceId},
      "point_ids": {pointIds},
      "date":      {cur.Format("2006-01-02")},
      "interval":  {"15min"}, // TODO
    }
    trends, err := novantReq(pCtx, "trends", args)
    if err != nil {
      log.DefaultLogger.Error("/trends request failed", "error", err)
      response.Error = err
      return response
    }

    // map values to frame format
    list := trends["trends"].([]interface{})
    for i := range list {
      row := list[i].(map[string]interface{})

      // decode ts
      ts, err := time.Parse(time.RFC3339, row["ts"].(string))
      if err != nil {
        response.Error = err
        return response
      }
      tss = append(tss, ts)

      // map trends to data frame
      for j := range pids {
        pid := pids[j]
        val := row[pid]
        switch t := val.(type) {
          case float64:
            vals[j] = append(vals[j], val.(float64))

          // TODO FIXIT: using nan for 'null' is not the same thing; but
          // does this have the right effect for now?
          default:
            _ = t
            vals[j] = append(vals[j], math.NaN())
        }
      }
    }

    // advance
    cur = cur.Add(24 * time.Hour)
  }

  // create data frame response
  frame := data.NewFrame("response")
  frame.Fields = append(frame.Fields,
    data.NewField("time", nil, tss),
  )
  for i := range vals {
    frame.Fields = append(frame.Fields, data.NewField(pnames[i], nil, vals[i]))
  }

  // If query called with streaming on then return a channel
  // to subscribe on a client-side and consume updates from a plugin.
  // Feel free to remove this if you don't need streaming for your datasource.
  /*
  if qm.WithStreaming {
    channel := live.Channel{
      Scope:     live.ScopeDatasource,
      Namespace: pCtx.DataSourceInstanceSettings.UID,
      Path:      "stream",
    }
    frame.SetMeta(&data.FrameMeta{Channel: channel.String()})
  }
  */

  // add the frames to the response
  response.Frames = append(response.Frames, frame)
  return response

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	if rand.Int()%2 == 0 {
		status = backend.HealthStatusError
		message = "randomized error"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
