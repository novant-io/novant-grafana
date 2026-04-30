package plugin

import (
	"strings"
	"sync"
	"time"
)

// pointCacheTTL is how long cached point metadata is considered fresh.
// Point names rarely change after commissioning, so a long TTL is fine; on
// miss we re-fetch /v1/points. Users can also clear the cache manually via
// the data source config UI.
const pointCacheTTL = 24 * time.Hour

type sourceEntry struct {
	fetched time.Time
	points  map[string]Point // pointID -> Point
}

// pointCache caches /v1/points responses per source so we can decorate live values
// and trend series with human-readable point names without hitting the API on every query.
type pointCache struct {
	mu      sync.RWMutex
	sources map[string]*sourceEntry
}

func newPointCache() *pointCache {
	return &pointCache{sources: make(map[string]*sourceEntry)}
}

// extractSourceID returns the "s.<n>" prefix of a point ID. Point IDs follow
// the format "s.<sourceId>.<pointId>" (e.g. "s.1.1" → "s.1"). Returns "" if
// the format doesn't match.
func extractSourceID(pointID string) string {
	parts := strings.SplitN(pointID, ".", 3)
	if len(parts) < 2 || parts[0] != "s" {
		return ""
	}
	return parts[0] + "." + parts[1]
}

// ensureSource fetches and caches the points for a source if not yet cached or stale.
// On API errors, an existing (stale) cache entry is left in place; otherwise the
// miss is silent and lookups will fall back to the raw point ID.
func (c *pointCache) ensureSource(client *Client, sourceID string) {
	c.mu.RLock()
	entry, ok := c.sources[sourceID]
	c.mu.RUnlock()

	if ok && time.Since(entry.fetched) < pointCacheTTL {
		return
	}

	resp, err := client.GetPoints(sourceID, "", "", "", "")
	if err != nil {
		return
	}

	points := make(map[string]Point, len(resp.Points))
	for _, p := range resp.Points {
		points[p.ID] = p
	}

	c.mu.Lock()
	c.sources[sourceID] = &sourceEntry{fetched: time.Now(), points: points}
	c.mu.Unlock()
}

// clear removes all cached source entries. The next query will re-fetch
// /v1/points for every source it touches.
func (c *pointCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sources = make(map[string]*sourceEntry)
}

// resolveNames returns a map of pointID → display name for the given point IDs.
// Falls back to the point ID itself if no cached name is available.
func (c *pointCache) resolveNames(client *Client, pointIDs []string) map[string]string {
	// Collect unique source IDs so we fetch each source's points only once.
	sources := make(map[string]struct{})
	for _, pid := range pointIDs {
		if sid := extractSourceID(pid); sid != "" {
			sources[sid] = struct{}{}
		}
	}
	for sid := range sources {
		c.ensureSource(client, sid)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	names := make(map[string]string, len(pointIDs))
	for _, pid := range pointIDs {
		names[pid] = pid // default fallback
		if sid := extractSourceID(pid); sid != "" {
			if entry, ok := c.sources[sid]; ok {
				if p, found := entry.points[pid]; found && p.Name != "" {
					names[pid] = p.Name
				}
			}
		}
	}
	return names
}
