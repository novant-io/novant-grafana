package plugin

import (
	"encoding/json"
	"strings"
)

// FlexString unmarshals a JSON string, number, bool, or null into a string.
// Used for fields like device_id where the API may return either a quoted
// string or a bare number depending on source type (e.g. BACnet integer IDs).
type FlexString string

func (s *FlexString) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*s = ""
		return nil
	}
	if trimmed[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = FlexString(str)
		return nil
	}
	*s = FlexString(trimmed)
	return nil
}

// QueryModel is the frontend query deserialized from JSON.
// The queryType is read from backend.DataQuery.QueryType (the top-level SDK field), not from here.
type QueryModel struct {
	// Entity filters
	ZoneIDs   string `json:"zoneIds"`
	SpaceIDs  string `json:"spaceIds"`
	AssetIDs  string `json:"assetIds"`
	SourceIDs string `json:"sourceIds"`
	// Points/values context
	SourceID string `json:"sourceId"`
	AssetID  string `json:"assetId"`
	SpaceID  string `json:"spaceId"`
	PointIDs string `json:"pointIds"`
	// Sources filter
	BoundOnly bool `json:"boundOnly"`
	// Trend options
	Interval  string `json:"interval"`
	Aggregate string `json:"aggregate"`
}

// Novant API response types

type ProjectResp struct {
	ProjID   int    `json:"proj_id"`
	ProjName string `json:"proj_name"`
	City     string `json:"city"`
	Tz       string `json:"tz"`
	Usage    int    `json:"usage"`
	Capacity int    `json:"capacity"`
}

type Zone struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	FedByAssetIDs []string `json:"fed_by_asset_ids"`
	FeedsSpaceIDs []string `json:"feeds_space_ids"`
}

type ZonesResp struct {
	Zones []Zone `json:"zones"`
}

type Space struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Type             string                `json:"type"`
	ParentSpaceID    string                `json:"parent_space_id"`
	ParentZoneID     string                `json:"parent_zone_id"`
	ContainsAssetIDs []string              `json:"contains_asset_ids"`
	Props            map[string]FlexString `json:"props"`
}

type SpacesResp struct {
	Spaces []Space `json:"spaces"`
}

type Asset struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Type      string                `json:"type"`
	Props     map[string]FlexString `json:"props"`
	SourceIDs []string              `json:"source_ids"`
}

type AssetsResp struct {
	Currency string  `json:"currency"`
	Assets   []Asset `json:"assets"`
}

type Source struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Addr          FlexString `json:"addr"`
	DeviceID      FlexString `json:"device_id"`
	Vendor        string     `json:"vendor"`
	Model         string     `json:"model"`
	Enabled       bool       `json:"enabled"`
	Bound         bool       `json:"bound"`
	ParentAssetID string     `json:"parent_asset_id"`
}

type SourcesResp struct {
	Sources []Source `json:"sources"`
}

type Point struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	Kind     string `json:"kind"`
	Unit     string `json:"unit"`
	Writable bool   `json:"writable"`
}

type PointsResp struct {
	SourceID    string  `json:"source_id"`
	SourceName  string  `json:"source_name"`
	SourceBound bool    `json:"source_bound"`
	Points      []Point `json:"points"`
}

type PointValue struct {
	ID     string      `json:"id"`
	Val    interface{} `json:"val"`
	Status string      `json:"status"`
}

type ValuesResp struct {
	SourceID string       `json:"source_id"`
	Values   []PointValue `json:"values"`
}

type TrendRow struct {
	Ts     string                 `json:"ts"`
	Values map[string]interface{} `json:"-"`
}

type TrendsResp struct {
	Start     string     `json:"start"`
	End       string     `json:"end"`
	Tz        string     `json:"tz"`
	Interval  string     `json:"interval"`
	Aggregate string     `json:"aggregate"`
	PointIDs  []string   `json:"point_ids"`
	Trends    []TrendRow `json:"trends"`
}
