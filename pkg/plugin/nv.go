//
// Copyright (c) 2021, Novant LLC
// Licensed under the MIT License
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

package plugin

import (
  "encoding/json"
  "errors"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
  "github.com/grafana/grafana-plugin-sdk-go/backend"
)

func novantReq(cx backend.PluginContext, op string, args url.Values) (Map, error) {
  apiKey := cx.DataSourceInstanceSettings.DecryptedSecureJSONData["apiKey"]
  client := &http.Client{}
  reader := strings.NewReader(args.Encode())

  req, err := http.NewRequest("POST", "https://api.novant.io/v1/" + op, reader)
  if err != nil {
    return nil, err
  }

  // setup auth and send request
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.SetBasicAuth(apiKey, "")
  res, err := client.Do(req)
  if err != nil {
    return nil, err
  }

  // read raw response
  resBytes, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  // decode json
  resMap := Map{}
  err = json.Unmarshal([]byte(resBytes), &resMap)
  if err != nil {
    return nil, err
  }

  // check status code
  if res.StatusCode != 200 {
    return nil, errors.New(resMap["err_msg"].(string))
  }

  return resMap, nil
}