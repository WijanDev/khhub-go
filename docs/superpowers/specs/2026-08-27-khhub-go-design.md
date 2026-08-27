# Design: khhub-go skill and a small backend alignment pass

**Date:** 2026-08-27  
**Status:** draft — awaiting review  
**Product:** khhub (not official JW software; no jw.org; UI in Spanish, docs in English)

## Problem

Agents and humans coming from class-oriented languages (and from the vendored `golang-gin-api` skill) treat one type per file and a service-per-resource stack as “correct Go.” In this repo that fights idiom and YAGNI. The HTTP package already groups work by concern (`auth.go`, `publishers.go`, `reports.go`). Domain rules already live in `domain` and `service` (hours eligibility, activity status, service year). CRUD handlers call sqlc directly.

Without a project skill that **wins** over generic Gin advice, new backend work will keep being “fixed” in the wrong direction: split `auth.go`, invent `PublisherService`, or copy `internal/handler` + `internal/repository`.

## Goals

1. A project skill `khhub-go` that agents load **when writing Go** and **when finishing backend work**, so the same bar applies before and after a change.
2. The skill encodes this repo’s layering and idiomatic Go package layout. It explicitly overrides `golang-gin-api` where they conflict.
3. A small alignment pass on the current backend: move the one real orchestration rule still sitting in HTTP, stop duplicating `IsHourReporter`, and stop leaking `err.Error()` to the SPA.
4. Catalog wiring so the skill is discoverable (`AGENTS.md`, `.agents/README.md`) and hooked from `khhub-stack`, `requesting-code-review`, and `verification-before-completion` when the diff touches product Go.

## Non-goals

- One type (or one handler) per file.
- A service or repository type per resource (`PublisherService`, `HouseholdRepository`).
- Renaming `internal/http` → `internal/handler` or `internal/store` → `internal/repository`.
- JWT, public signup, or any auth redesign.
- A full httptest suite for reports/publishers in this cycle (those stay their own roadmap ideas).
- Splitting `auth.go` because it holds several request types and handlers.
- Changing generated `internal/store/*.sql.go`.
- Frontend work.

## Decisions (locked)

| Topic | Choice |
| --- | --- |
| Path | Architectural process, approach 2: standard + extract real rules only |
| Skill name | `khhub-go` |
| Location | `.agents/skills/khhub-go/` (project skill, English) |
| Invocation | Auto-discoverable (no `disable-model-invocation`) |
| Write-time | Load when adding or changing product Go under `backend/` |
| Review-time | Load after a backend task, before claiming done; also when the user asks for a Go review |
| Review grades | Critical / Important / Minor. Fix Critical and Important in the same task. Note Minor. |
| Package vs file | Package is the unit. One file per **concern**, not per type. |
| When to split a file | Mixed concerns or hard to follow — not “this file has three structs.” |
| Handler → store | Allowed for CRUD (list/get/create/update/delete with no congregation rule). |
| Handler → service | Required when the work is a congregation rule, is used in more than one place, or should be tested without Gin. |
| Error to client | Fixed Spanish messages. Never `err.Error()` on the wire. |
| Interfaces | Small, defined by the consumer (existing `sessionQuerier` pattern). |
| First backend pass | Only the items in [Alignment pass](#alignment-pass). No resource-service rewrite. |

## Architecture

```
Agent writing or reviewing backend Go
  → load khhub-stack (sqlc, cookie, Spanish UI)
  → load khhub-go (package layout, layers, idiom, review checklist)
  → golang-gin-api / golang-gin-testing only where they do not conflict

http          bind, cookies, status codes, JSON; may call store for CRUD
service       congregation rules and orchestration (activity, later similar)
domain        pure types and predicates (ReportsHours, MustReport, ServiceYear)
store         sqlc + migrate only
auth          password hash and session tokens
```

`khhub-go` does not replace `khhub-stack`. Stack stays the source for sqlc, cookies, and product constraints. `khhub-go` is the source for **how to structure and review Go**.

## Skill shape

```
.agents/skills/khhub-go/
  SKILL.md       when to load, layering, anti-patterns, review workflow
  checklist.md   full review list (read only on a review pass)
```

`SKILL.md` stays short. The description (third person, WHAT + WHEN) must mention: Go backend, handlers, package layout, layers, review after backend changes, and that it overrides `golang-gin-api` on structure.

### Catalog

In `AGENTS.md` Skills table:

- Row “Go/Gin handlers or tests” becomes `khhub-stack`, `khhub-go`, `golang-gin-api`, `golang-gin-testing`.
- Add a row: after backend Go changes (or before claiming done on backend work) → `khhub-go`.

In `.agents/README.md`, list `khhub-go` under Quality (project-owned, not in `skills-lock.json`).

### Hooks in existing skills

- `khhub-stack`: one line — for Go package layout and layering, load `khhub-go`.
- `requesting-code-review` and `verification-before-completion`: if the change set includes `backend/**/*.go` excluding generated `internal/store/*.sql.go`, load `khhub-go` and run its review pass.

Do not vendor a second copy of the checklist into those skills.

## Checklist (normative)

### Organization

- Group by concern in the same package (`auth.go` = login, logout, me, password, cookies).
- Do not split a file because it declares more than one type or handler.
- Tests live in `*_test.go` beside the code.
- Do not hand-edit sqlc output.

### Layers

- `http` owns HTTP. It must not grow new congregation rules.
- `domain` owns pure predicates and calendar helpers already there (`ReportsHours`, `MustReport`, `IsHourReporter`, `ServiceYear`, `LastNMonths`).
- `service` owns derived status and any orchestration that loads store data to apply those rules (today: activity recompute).
- `store` is persistence only.

### Idiom agents often break

- `context.Context` is the first argument of blocking functions.
- Wrap internal errors with `%w`.
- Client JSON `error` is a fixed Spanish string, never the raw Go error.
- Exported identifiers get a comment when they leave the package.
- Prefer table-driven tests and `httptest` for handlers; do not invent mocks for simple CRUD.

### Anti-patterns (reject in review)

- “Move this struct to its own file.”
- New `*Service` that only forwards to `store.Queries`.
- JWT, `gin.Default()`, or `internal/handler` / `internal/repository` layout.
- Persisting activity status as a source of truth (the column may be a cache; `service.ActivityStatus` is the rule).

## Alignment pass

Apply in the same implementation cycle as the skill. Do not expand it.

1. Move `recomputeActivity` from `internal/http/reports.go` to `internal/service`. The function still uses `store.Queries` (or a small consumer-defined interface) to load shares and write `activity_status`. `putReports` calls it and maps failure to a Spanish 500.
2. In `listReports`, replace the inline `isRegularPioneer || isSpecialPioneer || auxiliary` with `domain.IsHourReporter`.
3. Map `domain.ReportsHours` failures in `putReports` by comparing the English `error` (domain stays English for tests). Do not send `err.Error()` to the client:

   | `ReportsHours` error | JSON `error` (400) |
   | --- | --- |
   | hours are only recorded for pioneers | las horas solo se registran para precursores |
   | hours cannot be negative | las horas no pueden ser negativas |
   | hours are required when a pioneer shared in the ministry | las horas son obligatorias cuando un precursor participó en el ministerio |
   | any other | datos de informes no válidos |

4. Map `auth.HashPassword` failure in `postChangePassword` to HTTP 400 and `no se pudo cambiar la contraseña`. Do not send `err.Error()`.

Leave `auth.go` as one file. Leave household, publisher, congregation, attendance, and dashboard handlers calling `store` directly.

## Testing

- New or moved `recomputeActivity` tests live under `internal/service` (no Gin).
- `make lint` and `make test` must pass.
- Existing `http` tests (`auth_test.go`, `router_test.go`) stay green; only `postChangePassword` response text may change if a test asserted the leaked error (today they do not cover that path).
- Do not add a reports httptest suite here.

## Error handling

- Review findings are graded; they are not a linter binary.
- If `golang-gin-api` says “handlers must not call the DB” and `khhub-go` says CRUD may call sqlc, **`khhub-go` wins**.
- Unknown `ReportsHours` errors use `datos de informes no válidos` (see the table in the alignment pass).

## Success criteria

- An agent loading skills for a new handler reads `khhub-go` and does not split types or add a pass-through service.
- After a backend change, claiming done includes a `khhub-go` review pass when product Go changed.
- `recomputeActivity` is callable and tested without spinning Gin.
- Password-change and report-hour validation never put a raw Go error in JSON.

## Out of scope for follow-ups (not this spec)

Handler tests for report upsert, `golangci-lint` CI (already a UAT idea), consistent JSON error **shape** (existing `consistent-json-errors` idea), and any further extraction from HTTP. Those stay separate roadmap items.
