# Roadmap

khhub is **not official Jehovah’s Witnesses software**. It does not talk to jw.org and does not replace S-21 cards or JW Hub. Dashboard totals are copied to the branch by hand.

Living work lives in this folder only. Agents follow `.agents/skills/roadmap-plan`. Idea folders sit in `pending/`, `implementing/`, or `implemented/` (move the folder when status changes). Dropped ideas sit in `dropped/`.

**Scores:** cost 1 = cheap … 5 = expensive. Utility 1 = low … 5 = high. **Ratio** = utility ÷ cost.

**Out of scope (not ideas):** official JW systems or jw.org APIs; automatic submission to the branch; public signup / self-service “create a congregation”; stored shepherding note text; real congregation PII in git. One install may host many congregations; that is `accounts-multicongregation`, not a public SaaS signup product.

Imported unchecked items from the retired root `ROADMAP.md` start as `proposed` with first-pass scores. Re-score after a real `implementation-plan.md` rewrite if the work grows.

Dropped folders are omitted from the tables: [`roles-beyond-admin`](dropped/roles-beyond-admin/idea.md), [`shepherding-notes`](dropped/shepherding-notes/idea.md).

## Pending

Sorted by ratio (highest first), then utility, then slug.

| Idea | Cost | Utility | Ratio |
| --- | --- | --- | --- |
| [httptest coverage for auth](pending/auth-handler-tests/idea.md) | 2 | 5 | 2.50 |
| [Secret scan in CI](pending/secret-scan-ci/idea.md) | 2 | 5 | 2.50 |
| [Prevent duplicate attendance rows](pending/attendance-duplicate-guard/idea.md) | 2 | 4 | 2.00 |
| [Copy dashboard totals to the clipboard](pending/dashboard-copy-totals/idea.md) | 2 | 4 | 2.00 |
| [PRs target `dev`, not `main`](pending/dev-branch-workflow/idea.md) | 2 | 4 | 2.00 |
| [golangci-lint or go vet + gofmt in CI](pending/golangci-lint-ci/idea.md) | 2 | 4 | 2.00 |
| [Login rate limit](pending/login-rate-limit/idea.md) | 2 | 4 | 2.00 |
| [Handler tests for report hour rules](pending/report-upsert-handler-tests/idea.md) | 2 | 4 | 2.00 |
| [Tests for ReportsHours and service year](pending/reports-hours-service-year-tests/idea.md) | 2 | 4 | 2.00 |
| [Tests that seed reset is blocked in production](pending/seed-reset-production-guard-tests/idea.md) | 2 | 4 | 2.00 |
| [Drop leftover Vite starter assets](pending/drop-vite-starter-assets/idea.md) | 1 | 2 | 2.00 |
| [LICENSE for the public repo](pending/license-file/idea.md) | 1 | 2 | 2.00 |
| [docs/architecture.md](pending/architecture-doc/idea.md) | 2 | 3 | 1.50 |
| [Change admin email](pending/change-admin-email/idea.md) | 2 | 3 | 1.50 |
| [CONTRIBUTING.md](pending/contributing-md/idea.md) | 2 | 3 | 1.50 |
| [Installable PWA shell](pending/pwa-installable-shell/idea.md) | 2 | 3 | 1.50 |
| [Attention list (shepherding planning)](pending/attention-list/idea.md) | 3 | 4 | 1.33 |
| [Dokploy beta environment](pending/dokploy-beta-environment/idea.md) | 3 | 4 | 1.33 |
| [Empty, loading, and error states on every page](pending/empty-loading-error-states/idea.md) | 3 | 4 | 1.33 |
| [Missing-report list as a first-class view](pending/missing-report-list/idea.md) | 3 | 4 | 1.33 |
| [Mobile layout pass](pending/mobile-layout-pass/idea.md) | 3 | 4 | 1.33 |
| [Month lock](pending/month-lock/idea.md) | 3 | 4 | 1.33 |
| [Publisher filters](pending/publisher-filters/idea.md) | 3 | 4 | 1.33 |
| [Service-year card per publisher](pending/service-year-publisher-card/idea.md) | 3 | 4 | 1.33 |
| [Revoke sessions on password change](pending/session-revoke-on-password-change/idea.md) | 3 | 4 | 1.33 |
| [SonarCloud and minimum quality bar](pending/sonarcloud-quality-gate/idea.md) | 3 | 4 | 1.33 |
| [Accounts, publisher portal, multi-congregation](pending/accounts-multicongregation/idea.md) | 5 | 5 | 1.00 |
| [Archive publishers](pending/archive-publishers/idea.md) | 4 | 4 | 1.00 |
| [Field service groups](pending/field-service-groups/idea.md) | 4 | 4 | 1.00 |
| [Keyboard and screen-reader pass](pending/a11y-pass/idea.md) | 3 | 3 | 1.00 |
| [Consistent JSON error shape](pending/consistent-json-errors/idea.md) | 3 | 3 | 1.00 |
| [CSV export of reports and dashboard](pending/export-csv/idea.md) | 3 | 3 | 1.00 |
| [CSV import](pending/import-csv/idea.md) | 3 | 3 | 1.00 |
| [Pagination or hard caps on list endpoints](pending/list-pagination/idea.md) | 3 | 3 | 1.00 |
| [Memorial attendance](pending/memorial-attendance/idea.md) | 3 | 3 | 1.00 |
| [Missing-report email reminder](pending/missing-report-email/idea.md) | 3 | 3 | 1.00 |
| [Pioneer hour goals](pending/pioneer-hour-goals/idea.md) | 3 | 3 | 1.00 |
| [Privilege history](pending/privilege-history/idea.md) | 3 | 3 | 1.00 |
| [Structured logging with request IDs](pending/structured-logging/idea.md) | 3 | 3 | 1.00 |
| [Dark mode](pending/dark-mode/idea.md) | 2 | 2 | 1.00 |
| [Audit log](pending/audit-log/idea.md) | 4 | 3 | 0.75 |
| [Public talks and visiting speakers](pending/public-talks/idea.md) | 4 | 3 | 0.75 |
| [OpenAPI sketch or documented route list](pending/openapi-or-route-list/idea.md) | 3 | 2 | 0.67 |
| [Cart / public witnessing schedule](pending/cart-schedule/idea.md) | 5 | 3 | 0.60 |
| [Midweek meeting (LMM) assignments](pending/lmm-assignments/idea.md) | 5 | 3 | 0.60 |
| [Territories](pending/territories/idea.md) | 5 | 3 | 0.60 |
| [Literature inventory](pending/literature-inventory/idea.md) | 5 | 2 | 0.40 |

## Implementing

Status `in-progress`. Folders live in `implementing/`. Same sort as pending.

| Idea | Cost | Utility | Ratio |
| --- | --- | --- | --- |
| — | — | — | — |

## Implemented

Sorted by merge date, newest first. Each row needs the merge date on `main` and the app version from `VERSION`.

| Idea | Cost | Utility | Merged | Version |
| --- | --- | --- | --- | --- |
| [Dependabot for Go and npm](implemented/dependabot/idea.md) | 1 | 3 | 2026-08-27 | 0.1.0 |
