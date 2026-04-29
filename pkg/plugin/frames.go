package plugin

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func buildZonesFrame(resp *ZonesResp) *data.Frame {
	count := len(resp.Zones)
	ids := make([]string, count)
	names := make([]string, count)
	types := make([]string, count)
	fedBy := make([]string, count)
	feeds := make([]string, count)

	for i, z := range resp.Zones {
		ids[i] = z.ID
		names[i] = z.Name
		types[i] = z.Type
		fedBy[i] = strings.Join(z.FedByAssetIDs, ", ")
		feeds[i] = strings.Join(z.FeedsSpaceIDs, ", ")
	}

	return data.NewFrame("zones",
		data.NewField("id", nil, ids),
		data.NewField("name", nil, names),
		data.NewField("type", nil, types),
		data.NewField("fed_by_asset_ids", nil, fedBy),
		data.NewField("feeds_space_ids", nil, feeds),
	)
}

func buildSpacesFrame(resp *SpacesResp) *data.Frame {
	count := len(resp.Spaces)
	ids := make([]string, count)
	names := make([]string, count)
	types := make([]string, count)
	parentSpace := make([]string, count)
	parentZone := make([]string, count)
	containsAssets := make([]string, count)

	for i, s := range resp.Spaces {
		ids[i] = s.ID
		names[i] = s.Name
		types[i] = s.Type
		parentSpace[i] = s.ParentSpaceID
		parentZone[i] = s.ParentZoneID
		containsAssets[i] = strings.Join(s.ContainsAssetIDs, ", ")
	}

	return data.NewFrame("spaces",
		data.NewField("id", nil, ids),
		data.NewField("name", nil, names),
		data.NewField("type", nil, types),
		data.NewField("parent_space_id", nil, parentSpace),
		data.NewField("parent_zone_id", nil, parentZone),
		data.NewField("contains_asset_ids", nil, containsAssets),
	)
}

func buildAssetsFrame(resp *AssetsResp) *data.Frame {
	count := len(resp.Assets)
	ids := make([]string, count)
	names := make([]string, count)
	types := make([]string, count)
	sourceIDs := make([]string, count)

	for i, a := range resp.Assets {
		ids[i] = a.ID
		names[i] = a.Name
		types[i] = a.Type
		sourceIDs[i] = strings.Join(a.SourceIDs, ", ")
	}

	return data.NewFrame("assets",
		data.NewField("id", nil, ids),
		data.NewField("name", nil, names),
		data.NewField("type", nil, types),
		data.NewField("source_ids", nil, sourceIDs),
	)
}

func buildSourcesFrame(resp *SourcesResp) *data.Frame {
	count := len(resp.Sources)
	ids := make([]string, count)
	names := make([]string, count)
	types := make([]string, count)
	addrs := make([]string, count)
	vendors := make([]string, count)
	models := make([]string, count)
	enabled := make([]bool, count)
	bound := make([]bool, count)
	parentAsset := make([]string, count)

	for i, s := range resp.Sources {
		ids[i] = s.ID
		names[i] = s.Name
		types[i] = s.Type
		addrs[i] = string(s.Addr)
		vendors[i] = s.Vendor
		models[i] = s.Model
		enabled[i] = s.Enabled
		bound[i] = s.Bound
		parentAsset[i] = s.ParentAssetID
	}

	return data.NewFrame("sources",
		data.NewField("id", nil, ids),
		data.NewField("name", nil, names),
		data.NewField("type", nil, types),
		data.NewField("addr", nil, addrs),
		data.NewField("vendor", nil, vendors),
		data.NewField("model", nil, models),
		data.NewField("enabled", nil, enabled),
		data.NewField("bound", nil, bound),
		data.NewField("parent_asset_id", nil, parentAsset),
	)
}

func buildPointsFrame(resp *PointsResp) *data.Frame {
	count := len(resp.Points)
	ids := make([]string, count)
	names := make([]string, count)
	types := make([]string, count)
	addrs := make([]string, count)
	kinds := make([]string, count)
	units := make([]string, count)
	writable := make([]bool, count)

	for i, p := range resp.Points {
		ids[i] = p.ID
		names[i] = p.Name
		types[i] = p.Type
		addrs[i] = p.Addr
		kinds[i] = p.Kind
		units[i] = p.Unit
		writable[i] = p.Writable
	}

	frame := data.NewFrame("points",
		data.NewField("id", nil, ids),
		data.NewField("name", nil, names),
		data.NewField("type", nil, types),
		data.NewField("addr", nil, addrs),
		data.NewField("kind", nil, kinds),
		data.NewField("unit", nil, units),
		data.NewField("writable", nil, writable),
	)
	frame.Meta = &data.FrameMeta{
		Custom: map[string]interface{}{
			"source_id":   resp.SourceID,
			"source_name": resp.SourceName,
		},
	}
	return frame
}

func buildValuesFrame(resp *ValuesResp) *data.Frame {
	count := len(resp.Values)
	ids := make([]string, count)
	vals := make([]*float64, count)
	statuses := make([]string, count)

	for i, v := range resp.Values {
		ids[i] = v.ID
		statuses[i] = v.Status
		switch val := v.Val.(type) {
		case float64:
			f := val
			vals[i] = &f
		case nil:
			vals[i] = nil
		default:
			// Try to handle other numeric types
			vals[i] = nil
		}
	}

	return data.NewFrame("values",
		data.NewField("id", nil, ids),
		data.NewField("value", nil, vals),
		data.NewField("status", nil, statuses),
	)
}

func buildTrendsFrames(resp *TrendsResp) (data.Frames, error) {
	if len(resp.Trends) == 0 || len(resp.PointIDs) == 0 {
		return data.Frames{}, nil
	}

	count := len(resp.Trends)
	timestamps := make([]time.Time, count)

	// Parse timestamps
	for i, row := range resp.Trends {
		t, err := time.Parse(time.RFC3339, row.Ts)
		if err != nil {
			// Try alternate formats
			t, err = time.Parse("2006-01-02T15:04:05", row.Ts)
			if err != nil {
				return nil, fmt.Errorf("parsing timestamp %q: %w", row.Ts, err)
			}
		}
		timestamps[i] = t
	}

	// Build one field per point ID
	fields := make([]*data.Field, 0, len(resp.PointIDs)+1)
	fields = append(fields, data.NewField("time", nil, timestamps))

	for _, pid := range resp.PointIDs {
		values := make([]*float64, count)
		for i, row := range resp.Trends {
			if v, ok := row.Values[pid]; ok {
				if f, ok := v.(float64); ok {
					values[i] = &f
				}
			}
		}
		fields = append(fields, data.NewField(pid, nil, values))
	}

	frame := data.NewFrame("trends", fields...)
	frame.Meta = &data.FrameMeta{
		PreferredVisualization: data.VisTypeGraph,
	}

	return data.Frames{frame}, nil
}
