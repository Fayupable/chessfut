# Chessfut

Turns your Chess.com stats into a FIFA Ultimate Team-style player card — OVR, six attributes (PAC/SHO/PAS/DRI/DEF/PHY), position, work rate, and badges — rated out of 99.

Inspired by [GitFut](https://github.com/Younesfdj/gitfut), which does the same for GitHub profiles. Card frame and flag assets are adapted from GitFut under the MIT License (see each package's `THIRD_PARTY_LICENSES.md`); no code was copied.

## Monorepo Structure

| Package                        | Description                                       |
|----------------------------------|------------------------------------------------------|
| [`chessfut-be`](./chessfut-be)   | Go API — hexagonal architecture, chess.com integration, FIDE-anchored scoring engine |
| [`chessfut-fe`](./chessfut-fe)   | Next.js frontend — player card, compare, leaderboard |

Each package has its own README with setup instructions, environment variables, and scripts.

## How Rating Works

Chessfut doesn't just average your chess.com stats. The scoring engine anchors each player's attribute spread around a skill baseline derived from their FIDE rating (or a title-based default for verified titled players missing one), then distributes chess.com's own signals (blitz/bullet/rapid ratings, tactics, puzzle rush, opening diversity) around that anchor. This keeps two players with similar chess.com ratings but very different real-world strength meaningfully apart on the card.

## Quick Start (Docker Compose)

```bash
cp .env.example .env   # fill in POSTGRES_*, REDIS_PASSWORD, ADMIN_API_KEY, CHESSCOM_USER_AGENT, NEXT_PUBLIC_API_URL
docker compose up --build
```

- Frontend: `http://localhost:3139`
- Backend API: `http://localhost:8099/api/v1`

For local development without Docker, see the per-package READMEs.

## Architecture

- **Backend**: hexagonal (ports & adapters), Postgres as source of truth, Redis as a cache layer, chess.com's public API as the only external data source, background cron for stale-card refresh and titled-player sync.
- **Frontend**: Next.js App Router, Server Components for data fetching, Tailwind CSS 4.

## CI

- `.github/workflows/backend-ci.yml` — gofmt, vet, build, golangci-lint, tests with coverage gate
- `.github/workflows/frontend-ci.yml` — lint, type-check, unit tests, production build

## License

MIT — see [LICENSE](./LICENSE). Card frame and flag graphics are third-party (GitFut, MIT) — see attribution notices in each package.

## Contributing

Issues and PRs are welcome. Please run the relevant package's lint/test/build commands locally before opening a PR — CI will run them again automatically.
