# Chat Micro-services (Auth & Chat)

A minimal gRPC backend split into two independent Go services:

1. **auth** – user registration / authentication
2. **chat_server** – chat rooms & messaging

Each service owns its own PostgreSQL database and is exposed via gRPC with protobuf contracts located under `*/api/*/*.proto`.

## Tech Stack

- Go 1.22
- gRPC + Protocol Buffers
- PostgreSQL
- `goose` for SQL migrations (invoked via Makefiles)
- Docker & Docker Compose (local databases)
- Makefile driven developer workflow

## Quick Start

### 0. Prerequisites

* Docker + Docker Compose
* Make (installed by default on macOS/Linux)
* Go 1.22+ (only if you want to run binaries outside Docker)

### 1. Start the databases

```bash
make db-up        # spins up two Postgres containers (auth & chat)
```

### 2. Run database migrations (both services)

```bash
make install-deps       # installs goose once
make local-migration-up # applies all pending migrations
```

> Need only one service? Use `make local-migration-auth` or `make local-migration-chat`.

### 3. Launch the Go services

Open **two** terminals:

```bash
# Terminal 1 – Auth service
cd auth
go run ./cmd/server
```

```bash
# Terminal 2 – Chat service
cd chat_server
go run ./cmd/server
```

Both servers read connection settings from their local `.env` files (already checked in).

### 4. Tear everything down

```bash
make db-down   # stops and removes Postgres containers
```

## Project Structure (simplified)

```
|-- auth/
|   |-- cmd/server          # main entry-point
|   |-- api/user_v1/*.proto # gRPC contract
|   |-- internal/...        # handlers, services, repo, etc.
|   |-- migrations/*.sql    # goose migrations
|
|-- chat_server/            # same layout for chat
|
|-- docker-compose.yaml     # two Postgres instances
|-- Makefile                # root-level helper targets
```

## Trying the APIs

Generate gRPC stubs in your preferred language from the `.proto` files or use `grpcurl`, e.g.:

```bash
grpcurl -d '{"login":"bob","password":"secret"}' \
  -plaintext localhost:50051 user_v1.UserService/CreateUser
```

---

That’s it – spin up Postgres, run migrations, and start each Go service.

## Monitoring

- Run `docker compose up prometheus grafana` from the repo root to start Prometheus and Grafana
- Prometheus scrapes the auth service on `host.docker.internal:2112` using `auth/prometheus.yml`.
- Grafana auto-loads a datasource pointing at Prometheus and ships with an `Auth Service Overview` dashboard (`http://localhost:3000`, login `admin/admin`).
- Dashboard JSON and provisioning files live under `auth/grafana/` if you want to tweak or extend the shipped panels.
