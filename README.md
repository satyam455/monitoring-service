# SmartX Token Monitor

SmartX Token Monitor is a Go service for chain-neutral, real-time token monitoring.

It includes Solana and EVM WebSocket ingestion, normalized events, concurrent
processing, PostgreSQL persistence, Redis caching, wallet signals, API-key
protection, Prometheus metrics, reconnect backoff, CI, and Docker deployment.

## Run

```bash
go test ./...
go run ./cmd/monitor
curl http://localhost:8080/healthz
```

The same workflows are available through `make test`, `make run`, and
`make compose-up`.

Pushes and pull requests are validated by GitHub Actions with Go 1.25,
race-enabled tests, `go vet`, and a container build.

The SmartX-facing API is documented in [openapi.yaml](openapi.yaml).
Set `API_KEY` to protect `/v1/*` endpoints with the `X-API-Key` header.

For Docker Compose, PostgreSQL persistence is enabled automatically. Set a
strong API key before deployment:

```bash
API_KEY='replace-with-a-random-secret' docker compose up --build
```

Run the local infrastructure with `docker compose up --build`. PostgreSQL is
initialized from `migrations/001_monitor_events.sql`; Redis is available at
`localhost:6379`.

Configuration is supplied through `HTTP_ADDR`, `WORKERS`, and `SHUTDOWN_TIMEOUT`.
The default Solana target is devnet; set `SOLANA_WS_URL` and `SOLANA_NETWORK`
explicitly before using another environment.

## Build history

1. Go service foundation and normalized event pipeline
2. Solana WebSocket/RPC adapter
3. EVM WebSocket adapter
4. PostgreSQL persistence and Redis hot cache
5. Large-wallet detection and SmartX signals API
6. Reconnect, metrics, validation, security, CI, and deployment tooling
