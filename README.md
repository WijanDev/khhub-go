# khhub

Personal tool for congregation operations: publishers, monthly field service reports, and meeting attendance. **Not official software.** It does not talk to jw.org and does not replace S-21 cards or JW Hub.

The UI is in Spanish. Code, API, and documentation are in English.

Dashboard totals are copied to the branch by hand.

## Requirements

- Go 1.24+
- Node 22+
- Docker (Postgres only in development)

## Local setup

```bash
cp .env.example .env
# set ADMIN_PASSWORD (and ADMIN_EMAIL if you want)
docker compose -f docker-compose.dev.yml --env-file .env up -d
cd backend && go run ./cmd/api   # loads ../.env automatically
```

In another terminal:

```bash
cd frontend && npm install && npm run dev
```

Copy `.env.example` so Vite has `VITE_API_URL=http://127.0.0.1:8080`. Open http://localhost:5173 and sign in with `ADMIN_EMAIL` / `ADMIN_PASSWORD`. The browser calls the API on port 8080 (CORS).

API reload: `go install github.com/air-verse/air@latest` and run `air` inside `backend/`.

Makefile shortcuts: `make dev-db`, `make api`, `make web`, `make test`.

## Service year

1 September – 31 August. Reports follow current practice: everyone records whether they shared in the ministry and their Bible studies; only pioneers (and auxiliary pioneers that month) record hours.

## Git

The default branch is **`dev`**. Open pull requests against `dev`. Promote a release with a PR from `dev` to `main`. Pushing `main` builds GHCR images and deploys khhub.app.

## Deploy

See [docs/deploy-dokploy.md](docs/deploy-dokploy.md).

## Agents

Project rules: [AGENTS.md](AGENTS.md). Living roadmap: [docs/roadmap/summary.md](docs/roadmap/summary.md). Skills (Cursor-agnostic): [.agents/README.md](.agents/README.md).
