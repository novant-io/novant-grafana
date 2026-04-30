# Changelog

## Version 1.2.0 (working)
* Add `Point Types` filter for `points` and `values` queries — comma-separated
  point-type allowlist (e.g. `zone_air_temp_sensor,discharge_air_temp_sensor`),
  passed through to the `point_types` API parameter. Supports template variables.
* Fix decoding of 4xx API errors. The Novant API gzip-compresses error
  responses; the client now decompresses on the error path so users see the
  actual error message instead of a grpc UTF-8 marshaling failure.

## Version 1.1.0 (29-Apr-2026)
* Display human-readable point names in Live Values and Trends panels
  (point name in the legend / `name` column; raw `point_id` preserved as a
  field label). Point metadata is cached per source for 24 hours to avoid
  hitting `/v1/points` on every query. A "Clear cache" button on the data
  source settings page forces a refresh on the next query.

## Version 1.0.0 (29-Apr-2026)
* Initial MVP
