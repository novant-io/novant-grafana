# Developing the Novant Grafana Plugin

This guide is for contributors working on the plugin itself. End users don't
need any of this — they install the plugin into their own Grafana per the
[README](README.md).

## Prerequisites

- Node.js ≥ 20
- Go 1.21 (`GOTOOLCHAIN=local` is set during builds)
- Docker (for the local Grafana dev server)
- [Mage](https://magefile.org/) for backend builds
  (`go install github.com/magefile/mage@latest`, then add `~/go/bin` to PATH)

## Building from Source

Clone the repo, install frontend dependencies once, then build both halves of
the plugin with a single command:

```bash
npm install
npm run build:all       # frontend + backend
```

Or build each half on its own:

```bash
npm run build           # frontend (production)
npm run dev             # frontend (watch mode)
npm run build:backend   # backend (mage -v)
```

The build output lands in `dist/`, which is what Grafana loads as the plugin.

## Running Locally

A `docker-compose.yml` runs Grafana with this plugin mounted from `dist/` and
unsigned plugin loading enabled. Before the first run, set up your **`.env`**
file with your Novant API key:

```bash
cp .env.example .env
# edit .env, set NOVANT_API_KEY=ak_your_real_key
```

The API key is auto-provisioned into the dev Grafana data source from this env
var on every startup, so it survives container resets.

Then run Grafana:

```bash
docker compose up           # foreground, Ctrl+C to stop
docker compose up -d        # detached
docker compose logs -f      # tail logs
docker compose down         # stop and remove the container (data persists)
```

Open <http://localhost:3000> (default Grafana login: `admin` / `admin`). The
**Novant** data source is auto-provisioned and the **Novant Overview**
dashboard (Dashboards → Novant folder) is too — entity tables plus Live
Values and Trends panels driven by dashboard variables.

### State persistence

Grafana state (dashboards you've edited in the UI, etc.) is stored in a named
Docker volume called `grafana-data`. It survives `docker compose down`,
container rebuilds, and Grafana version bumps. The data source's API key is
re-applied from `.env` on every startup, so it never goes stale.

### Resetting state

When you want a fresh Grafana — for example to test first-run provisioning:

```bash
docker compose down -v          # removes containers AND named volumes
docker compose up               # fresh DB, re-provisions everything from yaml
```

The API key still gets provisioned from `.env`, so no manual reconfig needed.

## Project Layout

```
.
├── pkg/                  Go backend (plugin binary)
│   ├── main.go           Entry point
│   └── plugin/
│       ├── datasource.go QueryData / CheckHealth handlers
│       ├── client.go     Novant API HTTP client
│       ├── models.go     Request / response types
│       └── frames.go     Grafana data frame builders
├── src/                  TypeScript frontend
│   ├── module.ts         Plugin registration
│   ├── datasource.ts     DataSourceWithBackend implementation
│   ├── types.ts          Query and config types
│   ├── plugin.json       Plugin manifest
│   └── components/
│       ├── ConfigEditor.tsx  API key configuration UI
│       └── QueryEditor.tsx   Query builder UI
├── provisioning/         Grafana provisioning for the dev container
├── docker-compose.yml    Local dev Grafana + plugin mount
├── Magefile.go           Go build entry point
├── go.mod                Go module: github.com/novant-io/novant-grafana
├── package.json          Frontend build / lint / test scripts
└── .env.example          Template for local dev API key
```

## Development Scripts

```bash
npm run build         # Production frontend build
npm run build:backend # Backend plugin binary (mage -v)
npm run build:all     # Frontend + backend
npm run package       # Build + zip a release artifact into rel/
npm run dev           # Watch frontend
npm run typecheck     # tsc --noEmit
npm run lint          # ESLint
npm run lint:fix      # ESLint with --fix
npm run test          # Jest in watch mode
npm run test:ci       # Jest one-shot
```

To run Grafana, use Docker directly: `docker compose up`. See
[Running Locally](#running-locally).

For ad-hoc Go work:

```bash
GOTOOLCHAIN=local go build ./pkg/...    # quick compile check, no plugin binary
```

## Cutting a Release

End users install the plugin from a `.zip` attached to a GitHub Release. To
cut one:

1. Bump `version` in `package.json` (and run `npm install` to update the lockfile).
2. Build and package:
   ```bash
   npm run package
   ```
   This runs `build:all`, stages `dist/` into a properly-named `novant-datasource/`
   directory, and zips it as `rel/novant-datasource-<version>.zip`. The `rel/`
   directory is gitignored.
3. Tag and push:
   ```bash
   git tag v$(node -p "require('./package.json').version")
   git push --tags
   ```
4. Create the GitHub Release and attach the zip:
   ```bash
   gh release create v1.0.0 rel/novant-datasource-1.0.0.zip \
     --title "v1.0.0" --notes "Initial release"
   ```

The plugin is unsigned, so users must add `allow_loading_unsigned_plugins =
novant-datasource` to their `grafana.ini` until/unless we sign it through
Grafana's catalog.

## How It Works

The frontend (`DataSourceWithBackend`) forwards every query to the Go backend,
which calls the Novant REST API at `https://api.novant.io`:

- All requests are `POST` with `application/x-www-form-urlencoded` bodies
- HTTP Basic Auth — API key as username, empty password
- Responses are gzip-decoded and converted to Grafana
  [data frames](https://grafana.com/developers/plugin-tools/key-concepts/data-frames)

`CheckHealth` calls `/v1/project` and reports the project name and city on
success, making misconfigured API keys obvious from the data source page.
