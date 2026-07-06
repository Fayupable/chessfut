# Chessfut — Backend

Go API that turns a Chess.com username into a FIFA Ultimate Team-style player card: OVR, six attributes (PAC/SHO/PAS/DRI/DEF/PHY), position, work rate, and badges — derived from a FIDE-anchored statistical scoring engine, not a simple average.

Hexagonal architecture (`domain` → `application` → `adapter`), PostgreSQL as the source of truth, Redis as a cache layer, chess.com's public API as the only external data source.

## Architecture

```
domain/              # pure Go structs — Player, Card, Game, Attributes, Position
application/
port/input/         # use case interfaces
port/output/        # driven ports (repository, cache, chess.com client)
service/             # use cases + the scoring engine (card_engine.go)
adapter/
chesscom/            # chess.com API client (rate-limited, deduped, backoff on 429)
postgres/            # source-of-truth repository + migrations runner
redis/               # cache adapter (titled players cached immediately, others promoted after 5+ views/hour)
http/                # handlers, middleware (CORS, admin auth, logging)
health/               # background health checkers for Postgres/Redis
```

## Data Flow

1. **Redis** (fast path) — titled players cached immediately on first fetch; non-titled players promoted to cache after 5+ views within an hour
2. **Postgres** (source of truth) — always served even if stale; a 10-minute cron refreshes cards older than their TTL
3. **chess.com API** — only hit for genuinely new usernames or the stale-refresh cron; never blocks a user request behind a live chess.com call if cached/stored data exists

Concurrent requests for the same never-before-seen username are deduplicated via `singleflight`, so a burst of traffic for one player triggers exactly one chess.com fetch.

## API

### Public

| Method | Path                              | Description                          |
|--------|------------------------------------|----------------------------------------|
| GET    | `/api/v1/player/{username}`        | Fast card (profile + stats only)       |
| GET    | `/api/v1/player/{username}/detailed` | Detailed card (+ games, openings, play style) |
| GET    | `/api/v1/leaderboard?limit=`       | Top players by OVR                    |
| GET    | `/api/v1/stats`                    | Total cards rated                     |
| GET    | `/api/v1/search?q=&limit=&offset=` | Username autocomplete (Postgres, not cached) |

### Admin (requires `X-Admin-Key` header, not CORS-enabled)

| Method | Path                                  | Description                    |
|--------|----------------------------------------|----------------------------------|
| POST   | `/api/admin/refresh/{username}`         | Force-rebuild one player's card |
| POST   | `/api/admin/sync-titled`                | Sync all titled players from chess.com |
| POST   | `/api/admin/refresh-stale`              | Manually trigger the stale-card refresh |

## Getting Started

```bash
cp .env.example .env   # fill in DATABASE_URL, REDIS_ADDR, CHESSCOM_USER_AGENT, CHESSCOM_BASE_URL, ADMIN_API_KEY
go run .
```

Migrations run automatically on startup via `golang-migrate`.

## Environment Variables

| Variable               | Required | Default          | Description                                  |
|-------------------------|----------|-------------------|-------------------------------------------------|
| `PORT`                  | No       | `8080`            | HTTP listen port                                 |
| `DATABASE_URL`          | Yes      | —                 | Postgres connection string                       |
| `REDIS_ADDR`            | No       | `localhost:6379`  | Redis address                                    |
| `REDIS_PASSWORD`        | No       | —                 | Redis password                                   |
| `ADMIN_API_KEY`         | No       | —                 | Required header value for `/api/admin/*` routes  |
| `CHESSCOM_USER_AGENT`   | Yes      | —                 | User-Agent sent to chess.com (contact info required by their API policy) |
| `CHESSCOM_BASE_URL`     | Yes      | —                 | Base URL for chess.com's public API              |
| `MIGRATIONS_PATH`       | No       | `db/migration`    | Path to Flyway-style SQL migration files         |

## Testing

```bash
go test -race ./domain/... ./application/service/... ./adapter/chesscom/...
```

Coverage thresholds are enforced in CI (`.testcoverage.yml`): 80% total, 90% for `application/service`.

## Rate Limiting & chess.com Etiquette

The chess.com client enforces a 1 req/s limiter (burst 2) and backs off to 0.2 req/s for 2 minutes on any 429, to stay well within chess.com's published API guidelines and avoid IP bans.

## Docker

```bash
docker build -t chessfut-be .
docker run -p 8099:8099 --env-file .env chessfut-be
```

## Linting & Formatting

```bash
gofmt -l .
go vet ./...
golangci-lint run --config=.golangci.yml
```
