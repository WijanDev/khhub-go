# AGENTS.md

Operating contract for coding agents in this repository. Humans start at [README.md](README.md). Living work lives in [docs/roadmap/](docs/roadmap/summary.md). Skills live in [`.agents/skills/`](.agents/skills/) (vendor-neutral Agent Skills path).

## Product

khhub is a **self-hosted congregation ops tool**. It is **not official Jehovah’s Witnesses software**, does not talk to jw.org, and does not replace S-21 cards or JW Hub. Dashboard totals are copied to the branch by hand.

- **UI copy:** Spanish.
- **Code, comments, API fields, cookies, routes, variables, and all documentation:** English.

Do not commit congregation PII or real secrets. Seed data stays fictional.

## Stack (do not drift)

| Layer | Choice |
| --- | --- |
| API | Go 1.24, Gin, `log` today (prefer `log/slog` on new work) |
| DB | PostgreSQL 16, pgx, **sqlc**, golang-migrate |
| Auth | httpOnly cookie `khhub_session`, SHA-256 token hash, bcrypt. **No JWT. No public signup.** |
| Frontend | Vite SPA (no SSR), React 19, TanStack Router + Query + Table |
| UI | Official shadcn **Base UI** (`base-nova`, `@base-ui/react`). **Not Radix.** Forms use NativeSelect. |
| Deploy | Docker Compose on Hetzner + Dokploy. Dev servers stay on the host. |

When a vendored skill mentions GORM, sqlx, JWT, Next.js, TanStack Start, or Radix, **this repo wins**. Use sqlc + pgx, cookie sessions, Vite SPA (no SSR), Table **v8** (`useReactTable`), and Base UI.

## Layout

```
backend/cmd/api/            # process entry
backend/internal/http/      # Gin routes and handlers
backend/internal/service/   # domain rules (activity, report hours)
backend/internal/store/     # sqlc output + migrate
backend/internal/domain/    # pure types and service-year helpers
backend/internal/auth/      # password + session helpers
backend/internal/seed/      # fictional demo congregation
backend/db/migrations/      # golang-migrate SQL
backend/db/queries/         # sqlc SQL
frontend/src/features/     # pages
frontend/src/components/ui # shadcn Base UI
```

## Commands

```bash
# Postgres only (preferred local setup)
make dev-db

# API (loads .env from backend/ or repo root)
make api
# or: cd backend && go run ./cmd/api

# Vite (uses VITE_API_URL, default http://127.0.0.1:8080)
make web

# Backend unit tests
make test
# or: cd backend && go test ./...

# Go format and vet (same checks as CI)
make lint

# After changing SQL
cd backend && sqlc generate

# Frontend
cd frontend && npm run lint && npm run build
```

Do **not** dockerize the API or Vite for day-to-day development. `docker-compose.yml` is for production-like / Dokploy runs.

## Hard rules

1. English documentation only. Spanish stays in the UI (`frontend/src/lib/labels.ts` and feature copy).
2. Never commit `.env`. Ship placeholders in `.env.example` only.
3. Schema changes: new migration + `db/queries` + `sqlc generate`. Do not hand-edit `internal/store/*.sql.go`.
4. Hours are recorded only for regular/special pioneers or auxiliary pioneers that month (`domain.ReportsHours`).
5. Service year is 1 September – 31 August (`domain.ServiceYear`).
6. Activity status is derived (regular / irregular / inactive) over a 6-month window. Do not persist it as a source of truth.
7. `POST /dev/reset-seed` exists only when `APP_ENV=development`. Staging (`APP_ENV=staging`) loads the fictional demo seed on an empty directory but does not expose the reset route.
8. Card padding: official shadcn `Card` pads header/content, not arbitrary children. Put page body in `CardContent` or keep the existing root padding convention.
9. Verify UI changes in the browser (or the closest substitute) before claiming done.
10. Hetzner: before creating a billed resource, and after listing live servers, load `hetzner-deploy` and estimate cost from live CLI prices (`hcloud server-type describe`, load-balancer types, volumes). Report hourly + monthly EUR. This is catalog rate × resources, not the Hetzner invoice (traffic overage and snapshots can add more). Unexpected spend → `hcloud server delete` after confirming the name.

## Skills

Load a skill when the task matches. Do not load all of them.

| Task | Skill |
| --- | --- |
| New feature / unclear design | `brainstorming` (then `roadmap-plan` after the design is approved; `docs-ai-prd` if a spec is still needed) |
| Multi-step implementation | `writing-plans` → `executing-plans` |
| Capture or rank a roadmap idea | `roadmap-plan` |
| Write or refresh English docs | `khhub-docs`, `developer-docs-planning`, `developer-docs-drafting` |
| Go/Gin handlers or tests | `khhub-stack`, `golang-gin-api`, `golang-gin-testing` |
| Schema / indexes / migrations | `khhub-stack`, `golang-gin-psql-dba` |
| TanStack Router | `tanstack-router`, `tanstack-router-best-practices`, then official `router-core` via Intent (below) |
| TanStack Query | `tanstack-query`, `tanstack-query-best-practices` |
| TanStack Table (v8) | `tanstack-table` |
| Hetzner Cloud / `hcloud` | `hetzner-deploy` |
| Dokploy deploy | `dokploy-deploy` |
| Bug | `systematic-debugging` |
| Before claiming done | `verification-before-completion` |
| After a sizable change | `requesting-code-review` |

Implementation plans go in `docs/plans/YYYY-MM-DD-<feature>.md` (not `docs/superpowers/plans/`).

Official TanStack Router skills ship inside `@tanstack/router-core` (versioned with npm). From `frontend/`:

```bash
npx @tanstack/intent@latest list
npx @tanstack/intent@latest load @tanstack/router-core#router-core
```

Query and Table do not ship Intent skills at the versions we use; use the vendored `.agents/skills/tanstack-*` copies. Do not follow Table v9 (`useTable`) docs — this app is Table v8.

Catalog and sources: [`.agents/README.md`](.agents/README.md).

## Git

- Default integration branch is **`dev`**. Open feature PRs against `dev`.
- Do **not** merge a feature branch into `main`.
- `main` is production. Release by opening a PR **`dev` → `main`**. That is the only path that should publish GHCR images and trigger Dokploy.
- Roadmap status **`uat`**: the idea is on `dev` (and staging when it is user-visible) but not yet on `main`. Folder: `docs/roadmap/uat/`.
- Dependabot PRs target `dev`.

## Security

- Congregation data is personal. Default to least privilege, no public listing of publishers, no debug dumps.
- Production: `APP_ENV=production`, `COOKIE_SECURE=true`, `CORS_ORIGINS=https://khhub.app` (SPA and API are different hosts), strong `ADMIN_PASSWORD`, Postgres backups on day one.
- Do not expose 5432 or 8080 to the internet.
