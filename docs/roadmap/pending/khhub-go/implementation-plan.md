# Implementation plan

Follow [docs/superpowers/specs/2026-08-27-khhub-go-design.md](../../../superpowers/specs/2026-08-27-khhub-go-design.md). Do not expand the alignment pass.

1. Add `.agents/skills/khhub-go/SKILL.md` (when to load, layering, anti-patterns, review workflow: Critical / Important / Minor) and `checklist.md` (full review list). Description must mention Go backend, handlers, package layout, layers, review after backend changes, and that this skill overrides `golang-gin-api` on structure. No `disable-model-invocation`.
2. Catalog: in `AGENTS.md`, add `khhub-go` to the “Go/Gin handlers or tests” row and add a row for after backend Go changes. List it under Quality in `.agents/README.md` (project-owned, not in `skills-lock.json`).
3. Hooks (one line each, do not copy the checklist): `khhub-stack` points here for package layout and layering; `requesting-code-review` and `verification-before-completion` load `khhub-go` when the change set includes `backend/**/*.go` except generated `internal/store/*.sql.go`.
4. Move `recomputeActivity` from `backend/internal/http/reports.go` to `backend/internal/service`. Keep a small consumer-defined store interface if that keeps tests free of Gin. `putReports` calls it and maps failure to a Spanish 500. Add table-driven tests in `internal/service` (no Gin).
5. In `listReports`, use `domain.IsHourReporter` instead of the inline `||`. In `putReports`, map `domain.ReportsHours` errors to the Spanish table in the spec (domain errors stay English). In `postChangePassword`, map `auth.HashPassword` failure to 400 `no se pudo cambiar la contraseña`. Do not send `err.Error()`. Leave `auth.go` and other CRUD handlers as they are.
6. Verify: `make lint` and `make test` (`cd backend && go test ./...`). Existing `auth_test.go` and `router_test.go` stay green.
