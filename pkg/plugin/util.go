//
// Copyright (c) 2021, Novant LLC
// Licensed under the MIT License
//
// History:
//   13 Dec 2021  Andy Frank  Creation
//

package plugin

import (
  "time"
)

type Map = map[string]interface{}
type List = []interface{}

// Convert given Time to midnight on the same day.
func toMidnight(t time.Time) time.Time {
  y,m,d := t.Date()
  return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// Flatten and filter /points response for just the given
// id list and return lookup for point meta data.
func filterPointMeta(res Map, pointIds []string) Map {
  acc := Map{}
  sources := res["sources"].(List)
  for i := range sources {
    src := sources[i].(Map)
    pts := src["points"].(List)
    for j := range pts {
      p := pts[j].(Map)
      acc[p["id"].(string)] = p
    }
  }
  return acc
}