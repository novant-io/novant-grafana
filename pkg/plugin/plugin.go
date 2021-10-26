//
// Copyright (c) 2021, Novant LLC
// Licensed under the MIT License
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

package plugin

import (
  "context"
  "encoding/json"
  "errors"
  "net/url"
  "strings"
  "time"

  "github.com/grafana/grafana-plugin-sdk-go/backend"
  "github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
  "github.com/grafana/grafana-plugin-sdk-go/backend/log"
  "github.com/grafana/grafana-plugin-sdk-go/data"
  // "github.com/grafana/grafana-plugin-sdk-go/live"
)

// Make sure NvDatasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler, backend.StreamHandler interfaces. Plugin should not
// implement all these interfaces - only those which are required for a particular task.
// For example if plugin does not need streaming functionality then you are free to remove
// methods that implement backend.StreamHandler. Implementing instancemgmt.InstanceDisposer
// is useful to clean up resources used by previous datasource instance when a new datasource
// instance created upon datasource settings changed.
var (
  _ backend.QueryDataHandler      = (*NvDatasource)(nil)
  _ backend.CheckHealthHandler    = (*NvDatasource)(nil)
  _ backend.StreamHandler         = (*NvDatasource)(nil)
  _ instancemgmt.InstanceDisposer = (*NvDatasource)(nil)
)

// NewNvDatasource creates a new datasource instance.
func NewNvDatasource(_ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
  return &NvDatasource{}, nil
}

// NvDatasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type NvDatasource struct{}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewNvDatasource factory function.
func (d *NvDatasource) Dispose() {
  // Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *NvDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
  log.DefaultLogger.Info("QueryData called", "request", req)

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

func (d *NvDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
  response := backend.DataResponse{}

  // unmarshal the JSON params
  params := map[string]interface{}{}
  response.Error = json.Unmarshal(query.JSON, &params)
  if response.Error != nil {
    return response
  }
  ts := query.TimeRange.From
  deviceId := strings.TrimSpace(params["deviceId"].(string))
  pointIds := strings.TrimSpace(params["pointIds"].(string))

  // validate
  if deviceId == "" {
    response.Error = errors.New("Missing deviceId value")
    return response
  }

  // validate
  if pointIds == "" {
    response.Error = errors.New("Missing pointIds value")
    return response
  }

  // query /trends for data
  args := url.Values{
    "device_id": {deviceId},
    "point_ids": {pointIds},
    "date":      {ts.Format("2006-01-02")},
    "interval":  {"15min"}, // TODO
  }
  trends, err := novantReq(pCtx, "trends", args)
  if err != nil {
    response.Error = err
    return response
  }

  // stub out working data structures
  tss  := []time.Time{}
  vals := [][]float64{}
  pids := strings.Split(pointIds, ",")
  for _ = range pids {
    vals = append(vals, []float64{})
  }

  // map values to frame format
  list := trends["data"].([]interface{})
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

        // TODO FIXIT: for now skip if nil/nan
        default:
          _ = t
      }
    }
  }

  // create data frame response
  frame := data.NewFrame("response")
  frame.Fields = append(frame.Fields,
    data.NewField("time", nil, tss), //[]time.Time{query.TimeRange.From, query.TimeRange.To}),
  )
  for i := range vals {
    frame.Fields = append(frame.Fields, data.NewField("values", nil, vals[i]))
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
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *NvDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
  log.DefaultLogger.Info("CheckHealth called", "request", req)

  var status = backend.HealthStatusOk
  var message = "Data source is working"

  // use /ping to test API key
  _, err := novantReq(req.PluginContext, "ping", url.Values{})
  if err != nil {
    status = backend.HealthStatusError
    message = err.Error()
  }

  return &backend.CheckHealthResult{
    Status:  status,
    Message: message,
  }, nil
}

// SubscribeStream is called when a client wants to connect to a stream. This callback
// allows sending the first message.
func (d *NvDatasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
  log.DefaultLogger.Info("SubscribeStream called", "request", req)

  status := backend.SubscribeStreamStatusPermissionDenied
  if req.Path == "stream" {
    // Allow subscribing only on expected path.
    status = backend.SubscribeStreamStatusOK
  }
  return &backend.SubscribeStreamResponse{
    Status: status,
  }, nil
}

// RunStream is called once for any open channel.  Results are shared with everyone
// subscribed to the same channel.
func (d *NvDatasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
  log.DefaultLogger.Info("RunStream called", "request", req)

  // Create the same data frame as for query data.
  frame := data.NewFrame("response")

  // Add fields (matching the same schema used in QueryData).
  frame.Fields = append(frame.Fields,
    data.NewField("time", nil, make([]time.Time, 1)),
    data.NewField("values", nil, make([]int64, 1)),
  )

  counter := 0

  // Stream data frames periodically till stream closed by Grafana.
  for {
    select {
    case <-ctx.Done():
      log.DefaultLogger.Info("Context done, finish streaming", "path", req.Path)
      return nil
    case <-time.After(time.Second):
      // Send new data periodically.
      frame.Fields[0].Set(0, time.Now())
      frame.Fields[1].Set(0, int64(10*(counter%2+1)))

      counter++

      err := sender.SendFrame(frame, data.IncludeAll)
      if err != nil {
        log.DefaultLogger.Error("Error sending frame", "error", err)
        continue
      }
    }
  }
}

// PublishStream is called when a client sends a message to the stream.
func (d *NvDatasource) PublishStream(_ context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
  log.DefaultLogger.Info("PublishStream called", "request", req)

  // Do not allow publishing at all.
  return &backend.PublishStreamResponse{
    Status: backend.PublishStreamStatusPermissionDenied,
  }, nil
}
