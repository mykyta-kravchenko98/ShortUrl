# ShortUrl

A high-performance URL shortening service built with Go. Application repo of
a two-repo GitOps portfolio project - infra, Argo CD, Terraform and Helm live
in the companion `shorturl-gitops` repo.

## Features
- REST API for creating and resolving short URLs
- LRU cache for fast redirect lookups
- Snowflake ID generation for unique URL keys
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
