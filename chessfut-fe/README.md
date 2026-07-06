# Chessfut — Frontend

Turns your Chess.com stats into a FIFA Ultimate Team-style player card, rated out of 99.

Built with Next.js 16 (App Router) and Tailwind CSS 4. Talks to the [`chessfut-be`](../chessfut-be) Go API for all data.

## Features

- **Player card** (`/[username]`) — OVR, six attributes (PAC/SHO/PAS/DRI/DEF/PHY), position, work rate, and badges rendered as a scalable card graphic
- **Compare** (`/compare`) — side-by-side stat comparison between two players, with a live example pair shown by default
- **Leaderboard** (`/leaderboard`) — top-rated players ranked by OVR
- **Username autocomplete** — debounced, paginated search against the backend, database-only (not cached)
- **Scouting metrics panel** — per-time-control ratings (bullet/blitz/rapid/daily), top openings, work rate

## Tech Stack

- Next.js 16 (App Router, Server Components for data fetching)
- React 19
- Tailwind CSS 4
- Vitest + Testing Library (pure-function unit tests)
- Docker (multi-stage build, standalone output, non-root runtime user)

## Getting Started

```bash
npm install
cp .env.example .env.local   # set NEXT_PUBLIC_API_URL to your backend URL
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Environment Variables

| Variable               | Required | Description                                      |
|-------------------------|----------|---------------------------------------------------|
| `NEXT_PUBLIC_API_URL`   | Yes      | Base URL of the `chessfut-be` API (e.g. `http://localhost:8099/api/v1`) |

## Scripts

| Command          | Description                          |
|-------------------|---------------------------------------|
| `npm run dev`     | Start the dev server                  |
| `npm run build`   | Production build (`output: standalone`) |
| `npm run start`   | Run the production build              |
| `npm run lint`    | ESLint (Next.js core-web-vitals rules)|
| `npm run test`    | Run the Vitest unit test suite        |

## Project Structure

```
src/
├── app/                  # routes (/, /[username], /compare, /leaderboard)
├── components/
│   ├── card/             # PlayerCard rendering, tier/frame resolution, name/stat formatting
│   ├── compare/           # comparison table and search
│   ├── home/              # landing page pieces (search form, card fan, mascot)
│   ├── layout/             # TopBar, Footer
│   ├── panels/             # scouting metrics panel
│   └── search/             # username autocomplete
├── lib/api.ts             # typed fetch client for the backend
└── types/card.types.ts    # mirrors the backend's card response shape
```

## Docker

```bash
docker build --build-arg NEXT_PUBLIC_API_URL=https://api.example.com -t chessfut-fe .
docker run -p 3139:3139 chessfut-fe
```

The image runs as a non-root user and exposes port `3139`.

## Attribution

Card frame and flag assets are adapted from [GitFut](https://github.com/Younesfdj/gitfut) (MIT License). See [THIRD_PARTY_LICENSES.md](./THIRD_PARTY_LICENSES.md) for details. No code was copied — only visual assets.
