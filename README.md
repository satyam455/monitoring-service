# SmartX Token Monitor

Priority 1 is the Go service foundation for a chain-neutral, real-time monitoring pipeline.

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

Run the local infrastructure with `docker compose up --build`. PostgreSQL is
initialized from `migrations/001_monitor_events.sql`; Redis is available at
`localhost:6379`.

Configuration is supplied through `HTTP_ADDR`, `WORKERS`, and `SHUTDOWN_TIMEOUT`.
The default Solana target is devnet; set `SOLANA_WS_URL` and `SOLANA_NETWORK`
explicitly before using another environment.

## Planned build order

1. Go service foundation and normalized event pipeline (this stage)
2. Solana WebSocket/RPC adapter for token launches and account updates (complete)
3. EVM WebSocket adapter for logs, swaps, and token launches (complete)
4. PostgreSQL persistence and Redis hot cache
5. Large-wallet detection, signal calculation, and SmartX API
