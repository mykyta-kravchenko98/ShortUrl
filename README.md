# ShortUrl

A high-performance URL shortening service built with Go. Application repo of
a two-repo GitOps portfolio project - infra, Argo CD, Terraform and Helm live
in the companion `shorturl-gitops` repo.

## Features
- REST API for creating and resolving short URLs
- LRU cache for fast redirect lookups
- Snowflake ID generation for unique URL keys, scrambled through a keyed
  bijective mix (`pkg/obfuscate`) before encoding - short_url is not a
  direct encoding of the timestamp-ordered ID, so it can't be enumerated
- Base62 encoding for compact short URLs
- PostgreSQL for persistent storage (parameterized queries)
- `/healthz` (liveness) and `/readyz` (readiness, checks DB) probes
- Structured JSON logging (`log/slog`)
- OpenTelemetry traces + metrics, exported over OTLP/gRPC to
  `OTEL_EXPORTER_OTLP_ENDPOINT` (defaults to `localhost:4317` - the address
  a sidecar collector injected by the `otel-sidecar-injector` controller
  listens on)
- Graceful shutdown on SIGINT/SIGTERM
- Fully containerized: CI builds/tests, then builds and pushes a Docker
  image to GHCR on `master`. No bare-metal/EC2 deploy path anymore.

## Tech Stack
- Go 1.24
- PostgreSQL
- Echo v4
- OpenTelemetry SDK (OTLP/gRPC)
- Docker (distroless runtime image)
- GitHub Actions

## Architecture
Clean layered architecture:
- `internal/handler` - HTTP handlers and routing
- `internal/service` - business logic
- `internal/db` - PostgreSQL repository
- `internal/cache` - LRU cache layer
- `internal/observability` - slog + OTel setup
- `pkg/generator` - Snowflake ID generation
- `pkg/hash_function` - Base62 encoding
- `pkg/obfuscate` - keyed bit-mix that stands between the Snowflake ID and
  the short_url, so short_url doesn't leak generation order

## Run locally

With Docker Compose (Postgres + OTel collector + app) - no config file needed,
Postgres connection details are passed as env vars in `docker-compose.yml`:
```bash
docker compose up --build
```

Without Docker:
```bash
cp config/dev.json.example config/dev.json   # then edit the password if you want one
make migrateup
go run .
```

## Endpoints
| Method | Path                    | Description            |
|--------|--------------------------|-------------------------|
| POST   | `/api/v1/data/shorten`   | Create a short URL     |
| GET    | `/api/v1/{hash}`         | Redirect to long URL   |
| GET    | `/api/v1/status`         | Trivial status string  |
| GET    | `/healthz`               | Liveness probe         |
| GET    | `/readyz`                | Readiness probe (DB)   |

## Deployment
No more zip + S3 + AWS CodeDeploy. CI (`.github/workflows/go.yml`) builds
and pushes a container image to `ghcr.io/<owner>/shorturl` on every push to
`master`. Rollout is handled by ArgoCD from the `shorturl-gitops` repo,
which watches for new image tags - this repo doesn't push to a server
directly anymore.

## Configuration & secrets
`config/dev.json` is gitignored - it's your local file, never committed.
`config/dev.json.example` is the committed template with placeholder
values. Any key in it can also be overridden by an env var without
touching the file, e.g. `POSTGRESDB_PASSWORD`, `POSTGRESDB_HOST`,
`SERVER_RESTPORT` (see `internal/config/config.go`) - that's how
`docker-compose.yml` points the app at the `postgres` service instead of
`localhost` without a second copy of the config file.

`config/config.yml` (the prod-path config, loaded via `LoadConfigYAML`)
only ever contains placeholder values in this repo - real values are
supplied at deploy time by the `shorturl-gitops` Helm chart, which mounts
its own ConfigMap over `/app/config/config.yml`.

`ID_OBFUSCATION_KEY` (optional, 16 hex chars / 64 bits) seeds
`pkg/obfuscate`. If unset, a random key is generated per process at
startup - that's fine correctness-wise (short_url is stored once at
creation and looked up by string afterwards, never re-derived from the
key), it just means the same Snowflake ID wouldn't map to the same
short_url across a restart, which nothing here relies on. Set it
explicitly only if you want that reproducibility for some other reason.
