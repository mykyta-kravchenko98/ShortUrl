# ShortUrl

A high-performance URL shortening service built with Go.

## Features
- REST API for creating and resolving short URLs
- LRU cache for fast redirect lookups
- Snowflake ID generation for unique URL keys
- Base62 encoding for compact short URLs
- PostgreSQL for persistent storage
- CI/CD via GitHub Actions
- AWS CodeDeploy deployment support

## Tech Stack
- Go
- PostgreSQL
- LRU Cache (in-memory)
- GitHub Actions
- AWS CodeDeploy

## Architecture
Clean layered architecture:
- `internal/handler` – HTTP handlers and routing
- `internal/service` – business logic
- `internal/db` – PostgreSQL repository
- `internal/cache` – LRU cache layer
- `pkg/generator` – Snowflake ID generation
- `pkg/hash_function` – Base62 encoding

## Run locally
```bash
cp config/dev.json config/local.json
make run
```
