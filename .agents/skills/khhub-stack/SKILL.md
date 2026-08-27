---
name: khhub-stack
description: khhub-specific stack conventions (Go/Gin/sqlc/pgx cookie auth, Vite React, shadcn Base UI). Use when adding API endpoints, migrations, queries, UI pages, or when a generic Gin/React skill conflicts with this repo.
---

# khhub stack

Read `AGENTS.md` first. This skill wins over vendored Gin/React skills.

## Backend

- Handlers in `backend/internal/http`. Domain rules in `internal/service` or `internal/domain`. Persistence only through sqlc `store.Queries`.
- New SQL: `backend/db/queries/*.sql` + a new golang-migrate file in `backend/db/migrations/`. Then `sqlc generate` from `backend/`. Never edit generated `internal/store/*.sql.go` by hand.
- Auth is an httpOnly cookie (`khhub_session`), not JWT. Protect routes with `requireAuth()`. Login is rate-limited. No public registration.
- Validate with go-playground/validator where the handler already does; return JSON errors, not HTML.
- Hours: `domain.ReportsHours`. Service year: `domain.ServiceYear`. Activity: `service.ActivityStatus` (derived, 6 months).
- Demo seed (`internal/seed`) is fictional and development-only.

## Frontend

- Vite SPA (not TanStack Start, no SSR). Routes in `frontend/src/main.tsx` (TanStack Router). Data via `frontend/src/lib/api.ts` + TanStack Query (`credentials: "include"`).
- Tables: TanStack Table **v8** (`useReactTable`). Keep `data` and `columns` referentially stable. Do not migrate to v9 `useTable()` unless the roadmap says so.
- UI: shadcn **base-nova** / `@base-ui/react`. Do not add Radix. Prefer NativeSelect for forms. Alert Dialog for destructive confirms.
- User-visible strings in Spanish. Component names, props, and API keys in English.
- Keep `Card` content inside padded regions.

## Out of scope for a typical change

Roles, territories, LMM, cart, literature, and shepherding notes are later roadmap items. Do not scaffold them inside an unrelated fix.
