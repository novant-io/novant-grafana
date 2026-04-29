package plugin

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const apiBaseURL = "https://api.novant.io"

// Client is an HTTP client for the Novant API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Novant API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// post performs a POST request with form-encoded body and returns the decoded response.
func (c *Client) post(path string, params url.Values, result interface{}) error {
	var body io.Reader
	if params != nil {
		body = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequest("POST", apiBaseURL+path, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Accept-Encoding", "gzip")
	if params != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(b))
	}

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("creating gzip reader: %w", err)
		}
		defer gr.Close()
		reader = gr
	}

	if err := json.NewDecoder(reader).Decode(result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

func (c *Client) GetProject() (*ProjectResp, error) {
	var resp ProjectResp
	if err := c.post("/v1/project", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetZones(zoneIDs string) (*ZonesResp, error) {
	params := url.Values{}
	if zoneIDs != "" {
		params.Set("zone_ids", zoneIDs)
	}
	var resp ZonesResp
	if err := c.post("/v1/zones", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetSpaces(spaceIDs string) (*SpacesResp, error) {
	params := url.Values{}
	if spaceIDs != "" {
		params.Set("space_ids", spaceIDs)
	}
	var resp SpacesResp
	if err := c.post("/v1/spaces", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetAssets(assetIDs string) (*AssetsResp, error) {
	params := url.Values{}
	if assetIDs != "" {
		params.Set("asset_ids", assetIDs)
	}
	var resp AssetsResp
	if err := c.post("/v1/assets", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetSources(sourceIDs string, boundOnly bool) (*SourcesResp, error) {
	params := url.Values{}
	if sourceIDs != "" {
		params.Set("source_ids", sourceIDs)
	}
	if boundOnly {
		params.Set("bound_only", "true")
	}
	var resp SourcesResp
	if err := c.post("/v1/sources", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetPoints(sourceID, assetID, spaceID, pointIDs string) (*PointsResp, error) {
	params := url.Values{}
	if sourceID != "" {
		params.Set("source_id", sourceID)
	} else if assetID != "" {
		params.Set("asset_id", assetID)
	} else if spaceID != "" {
		params.Set("space_id", spaceID)
	}
	if pointIDs != "" {
		params.Set("point_ids", pointIDs)
	}
	var resp PointsResp
	if err := c.post("/v1/points", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetValues(sourceID, assetID, spaceID, pointIDs string) (*ValuesResp, error) {
	params := url.Values{}
	if sourceID != "" {
		params.Set("source_id", sourceID)
	} else if assetID != "" {
		params.Set("asset_id", assetID)
	} else if spaceID != "" {
		params.Set("space_id", spaceID)
	}
	if pointIDs != "" {
		params.Set("point_ids", pointIDs)
	}
	var resp ValuesResp
	if err := c.post("/v1/values", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetTrends(pointIDs, startDate, endDate, interval, aggregate string) (*TrendsResp, error) {
	params := url.Values{}
	params.Set("point_ids", pointIDs)
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	if interval != "" {
		params.Set("interval", interval)
	}
	if aggregate != "" {
		params.Set("aggregate", aggregate)
	}

	// The trends response has dynamic keys per point ID in each trend row,
	// so we decode to a raw structure first.
	var raw json.RawMessage
	if err := c.post("/v1/trends", params, &raw); err != nil {
		return nil, err
	}

	// Decode known fields
	var resp TrendsResp
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("decoding trends metadata: %w", err)
	}

	// Decode the trends array with dynamic point value keys
	var rawOuter map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rawOuter); err != nil {
		return nil, fmt.Errorf("decoding trends raw: %w", err)
	}

	trendsRaw, ok := rawOuter["trends"]
	if !ok {
		return &resp, nil
	}

	var rawTrends []map[string]interface{}
	if err := json.Unmarshal(trendsRaw, &rawTrends); err != nil {
		return nil, fmt.Errorf("decoding trend rows: %w", err)
	}

	resp.Trends = make([]TrendRow, len(rawTrends))
	for i, row := range rawTrends {
		ts, _ := row["ts"].(string)
		resp.Trends[i] = TrendRow{
			Ts:     ts,
			Values: make(map[string]interface{}),
		}
		for k, v := range row {
			if k != "ts" {
				resp.Trends[i].Values[k] = v
			}
		}
	}

	return &resp, nil
}
