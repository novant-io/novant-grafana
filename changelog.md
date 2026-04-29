# Changelog

## Version 1.1.0 (working)
* Display human-readable point names in Live Values and Trends panels
  (point name in the legend / `name` column; raw `point_id` preserved as a
  field label). Point metadata is cached per source for 24 hours to avoid
  hitting `/v1/points` on every query. A "Clear cache" button on the data
  source settings page forces a refresh on the next query.

## Version 1.0.0 (29-Apr-2026)
* Initial MVP
