# Roadmap

Living checklist for khhub. Product work and engineering work stay in separate sections. Check items in place; do not delete them.

Agents: load `.agents/skills/maintain-roadmap` before editing this file.

## How to use

1. Pick the next unchecked item in **Development stabilization** unless a product need is urgent.
2. For multi-step work, write `docs/plans/YYYY-MM-DD-<feature>.md` first.
3. Mark the item `[x]` in the same change that lands the work.

## Shipped (MVP)

- [x] Single-admin login (email/password, httpOnly session cookie)
- [x] Congregation settings
- [x] Households and publishers CRUD
- [x] Monthly field service reports (shared + bible studies; hours only for RP/SP or auxiliary that month)
- [x] Derived activity status (regular / irregular / inactive over 6 months)
- [x] Weekly attendance (midweek / weekend)
- [x] Dashboard totals for hand-copy to the branch
- [x] Fictional demo seed and congregation reset (development only)
- [x] Vite SPA (Spanish UI) + Go API + Compose for Dokploy
- [x] Public repo snapshot at [WijanDev/khhub-go](https://github.com/WijanDev/khhub-go)

## Development stabilization

Engineering work so the MVP is safe to change.

### Tooling and CI

- [x] GitHub Actions: `go test ./...`, `npm run lint`, `npm run build`, GHCR image push (no `sqlc diff` yet)
- [ ] golangci-lint (or `go vet` + `gofmt` check) in CI and locally
- [ ] Dependabot or equivalent for Go modules and npm
- [ ] LICENSE file for the public repo

### Tests

- [x] Unit tests for derived activity status
- [ ] Table-driven tests for `domain.ReportsHours` and service-year helpers
- [ ] Handler tests with `httptest` for auth (login, session, unauthorized)
- [ ] Handler tests for report upsert rules (hours rejected for non-pioneers)
- [ ] Seed reset keeps the admin user and is blocked when `APP_ENV=production`

### Observability and API hygiene

- [ ] Structured logging (`log/slog`) with request IDs
- [ ] Consistent JSON error shape on all handlers
- [ ] Pagination or hard caps on list endpoints
- [ ] OpenAPI sketch or documented route list in `docs/`

### Security and operations

- [ ] Session listing / revoke-all on password change
- [x] Confirm production cookie flags, CORS, and secure headers in Compose (`COOKIE_SECURE`, `CORS_ORIGINS=https://khhub.app`)
- [ ] Document restore-from-backup in `docs/deploy-dokploy.md`
- [ ] Secret scan in CI (do not rely on `.gitignore` alone)

### Frontend quality

- [ ] Empty, loading, and error states on every feature page
- [ ] Keyboard and screen-reader pass on login, publishers table, and report form
- [ ] Mobile layout pass (publishers table and monthly reports)
- [ ] Drop leftover Vite starter assets that are not part of the product

### Docs and harness

- [x] English README, deploy doc, and `AGENTS.md`
- [x] Agent skills under `.agents/skills/`
- [ ] Short `docs/architecture.md` (request path: UI → `api.khhub.app` → handler → sqlc)
- [ ] `CONTRIBUTING.md` with local commands and sqlc workflow

## Product hardening (still MVP)

Improvements to what is already shipped, not new ministries.

- [ ] Change admin email (not only password)
- [ ] Publisher filters: activity status, pioneer, missing report
- [ ] Missing-report list as a first-class view (if not already obvious on Reports)
- [ ] Attendance: prevent duplicate meeting rows for the same date + type
- [ ] Dashboard: copy-to-clipboard for branch totals
- [ ] Soft-delete or archive publishers (keep historical reports)
- [ ] Export monthly reports / dashboard as CSV

## Later features (out of current MVP)

Do not start these until stabilization is in good shape, unless a real congregation need jumps the queue.

- [ ] Roles beyond single admin (e.g. secretary vs view-only)
- [ ] Territories
- [ ] Midweek meeting (LMM) assignments
- [ ] Cart / public witnessing schedule
- [ ] Literature inventory
- [ ] Shepherding notes (access control required)

## Explicitly out of scope

- Official JW systems, jw.org APIs, or S-21 replacement
- Automatic submission to the branch
- Multi-congregation SaaS / public signup
- Storing real congregation PII in the public git history
