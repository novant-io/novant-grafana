# Novant plugin for Grafana

The [Novant](https://novant.io) plugin for [Grafana](https://grafana.com)
provides data source integration for the Novant API.

## Limitations

This plugin is currently unsigned, which creates some restrictions:

  * This plugin cannot be used in a paid Grafana hosted account
  * To enable plugin in self-hosted you must whitelist in `grafana.ini`

## Installing Plugin

To install the `novant-datasource` plugin:

  * Download the latest `novant-datasource.zip` from
    [Releases](https://github.com/novant-io/novant-grafana/releases/tag/0.2)

  * Unzip `novant-datasource.zip` into your Grafana `plugins/` directory
    (should look like `plugins/novant-datasource/`)

  * Update this line in your`grafana.ini` to:

        allow_loading_unsigned_plugins = novant-datasource

  * Restart Grafana for changes to take effect


## Adding Novant Datasource

To add a new Novant datasource:

  * In Grafana navigate to `Add data source`
  * Search or find `Novant` in the source list
  * Supply a `name` and enter your `API key`
  * Click `Save and Test` to validate your key and connectivity

## Using your Novant Datasource

In your dashboard:

  * Create a new panel
  * Select your datasource
  * Configure the `device_id` and `point_ids`

The visualization should automatically update when changes are made.  The
`point_ids` is a comma-separated list of points to display, ie: `"p1,p2,p3"`.