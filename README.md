# Curate

Curate fetches cryptocurrency exchange rates for a target quote currency, stores history in InfluxDB (or in memory), and
exposes them
over an HTTP API.

## Environment

All flags can be provided via CLI or env vars with prefix `CURATE_`.

| Env Var                    | CLI Flag              | Purpose                                        | Default                             |
|----------------------------|-----------------------|------------------------------------------------|-------------------------------------|
| `CURATE_REST_HOST`         | `--rest-host`         | HTTP server host                               | `127.0.0.1`                         |
| `CURATE_REST_PORT`         | `--rest-port`         | HTTP server port                               | `8080`                              |
| `CURATE_POLLING_INTERVAL`  | `--polling-interval`  | Frequency of rate providers polling            | `10s`                               |
| `CURATE_CURRENCIES`        | `--currencies`        | Comma‑separated list of currencies             | `BTC,ETH,TRX`                       |
| `CURATE_QUOTE`             | `--quote`             | Quote currency                                 | `USD`                               |
| `CURATE_HISTORY_LIMIT`     | `--history-limit`     | Max history points that the service can return | `10`                                |
| `CURATE_COINDESK_URL`      | `--coindesk-url`      | Upstream price API base URL                    | `https://min-api.cryptocompare.com` |
| `CURATE_COINDESK_TOKEN`    | `--coindesk-token`    | Upstream API token (optional)                  | ``                                  |
| `CURATE_IN_MEMORY_STORAGE` | `--in-memory-storage` | Use in-memory storage instead of InfluxDB      | `false`                             |
| `CURATE_INFLUXDB_URL`      | `--influxdb-uri`      | InfluxDB URL                                   | `http://127.0.0.1:8181`             |
| `CURATE_INFLUXDB_TOKEN`    | `--influxdb-token`    | InfluxDB token (not required if auth disabled) | `dev-token`                         |
| `CURATE_INFLUXDB_BUCKET`   | `--influxdb-bucket`   | InfluxDB bucket/database                       | `curate`                            |

#### Run locally:

```
make docker-build
make run
# or if you don't need InfluxDB:
go run ./cmd/service --rest-host=127.0.0.1 --rest-port=8080 --in-memory-storage=true
```

---

## API

- OpenAPI spec: `api/openapi.yaml`.
- To view Swagger locally: `npx swagger-ui-watcher api/openapi.yaml` (or import into Postman/Insomnia).

Endpoints:

- `GET /api/v1/rates?limit={n}` — all rates (limit not more than `HISTORY_LIMIT`)
- `GET /api/v1/rates/{currency}?limit={n}` — rate history for a currency
- `GET /api/v1/currencies` — supported currencies
- `GET /api/v1/health` — liveness (204)

**Note:**

- If `limit` exceeds configured `HISTORY_LIMIT`, only the latest `HISTORY_LIMIT` points are returned to avoid overload.
- If `limit=0`, an empty list is returned.

---

## Components

### `cmd/service/main.go`

- entrypoint; reads CLI/env config, builds dependencies, starts HTTP server, scheduler, and storage.

### `internal/service`

- core orchestrator with gocron. Schedule provider fetches, writes to storage, serves data to
  HTTP.

### `internal/provider`

- Upstream price provider (CryptoCompare). Additional providers can be added behind the same
  interface.

### `internal/storage`

- Two implementations available:
-
    - **influxdb v3**:
      Default external database. Persists data to `/var/lib/influxdb3`.
      Authentication disabled by default. (see `docker-compose.yml`)

-
    - **in-memory storage**: A simple in-memory storage implemented via go map for local development and testing.

### `internal/transport/http`

- Fiber server exposing the OpenAPI contract, handlers, validation, request/response
  mapping.

### `internal/clients/rest`

- shared HTTP client wrapper (resty) with logging for upstream calls.

### `internal/log`

slog wrapper to standardize structured logging across components.

---

