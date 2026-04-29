# Novant Data Source for Grafana

A Grafana data source plugin for [Novant](https://novant.io) — visualize zones,
spaces, assets, sources, and point data from your building automation systems
directly inside Grafana dashboards.

- **Plugin ID:** `novant-datasource`
- **Type:** Data source (with backend)
- **License:** MIT

## Features

- Query historical time series ("trends") for any Novant point with selectable
  intervals and aggregations
- Query live current values for points
- Browse Novant entity metadata as Grafana tables: zones, spaces, assets,
  sources, and points
- Grafana template variable support across all entity/point ID fields
- Alerting compatible (`alerting: true`)

## Query Types

| Type        | Returns                          | Required fields           |
| ----------- | -------------------------------- | ------------------------- |
| `trends`    | Time series of point values      | `Point IDs`               |
| `values`    | Current point values (table)     | `Source ID`, `Asset ID`, or `Space ID` |
| `points`    | Point metadata (table)           | `Source ID`, `Asset ID`, or `Space ID` |
| `sources`   | Source devices (table)           | — (optional `Source IDs`) |
| `assets`    | Equipment / assets (table)       | — (optional `Asset IDs`)  |
| `spaces`    | Building spaces (table)          | — (optional `Space IDs`)  |
| `zones`     | Building zones (table)           | — (optional `Zone IDs`)   |

For `trends`, the dashboard time range is sent as `start_date` / `end_date` to
the Novant API. `interval` and `aggregate` default to `auto`.

## Prerequisites

- [Grafana](https://grafana.com/) ≥ 10.0
- A Novant account and API key (starts with `ak_`) — sign up at
  [novant.io](https://novant.io)

## Installing

The plugin is not yet listed in the Grafana plugin catalog. For now, install
the latest release `.zip` from this repo's
[Releases page](https://github.com/novant-io/novant-grafana/releases):

```bash
# Download the .zip from a release, then:
unzip novant-datasource-<version>.zip -d /var/lib/grafana/plugins/
```

Or via `grafana-cli` with the release URL:

```bash
grafana-cli --pluginUrl https://github.com/novant-io/novant-grafana/releases/download/v<version>/novant-datasource-<version>.zip \
  plugins install novant-datasource
```

Because the plugin is unsigned, allow it to load by adding the following to
your `grafana.ini` under `[plugins]`:

```ini
allow_loading_unsigned_plugins = novant-datasource
```

Then restart Grafana.

## Configuring the Data Source

1. In Grafana, go to **Connections → Data sources → Add data source** and pick
   **Novant**.
2. Enter your Novant **API key** (`ak_...`). The key is stored as
   `secureJsonData` and decrypted only on the backend.
3. Click **Save & test** — a successful health check shows the connected
   project name and city.

## Building Queries

Open any dashboard panel, choose the Novant data source, then in the query
editor:

1. Pick a **Query Type**.
2. Fill in the relevant fields. ID fields accept comma-separated lists and
   support Grafana [template variables](https://grafana.com/docs/grafana/latest/dashboards/variables/)
   (e.g. `$point_ids`). Variable names are conventionally `snake_case` to
   match the Novant API parameter names (`point_ids`, `source_ids`,
   `zone_ids`, etc.).
3. For `trends`, select an **Interval** (`auto`, `5min`, `15min`, `30min`,
   `1hr`, `1day`, `1mo`, `raw`) and **Aggregate** (`auto`, `mean`, `sum`,
   `min`, `max`, `diff`).

### Example: trend chart

- Query Type: `Trends`
- Point IDs: `s.1.1,s.1.2`
- Interval: `15min`
- Aggregate: `mean`

### Example: live values table

- Query Type: `Live Values`
- Source ID: `s.1`  *(or Asset ID: `a.1`, or Space ID: `sp.1`)*
- Point IDs: *(optional — leave empty for all)*

## Contributing

To build, modify, or contribute to the plugin, see
[DEVELOPERS.md](DEVELOPERS.md).
