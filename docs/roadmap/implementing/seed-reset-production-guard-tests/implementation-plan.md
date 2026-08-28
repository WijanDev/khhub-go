# Implementation plan

1. Change `config.AllowsSeedReset` to `development` or `staging` (never “not production”). Keep `config_test.go` in lockstep; reject unknown `APP_ENV` values.
2. In `NewRouter`, force `resetSeed = nil` when `cfg.Production()` so a miswired callback cannot mount the route on production.
3. Add httptest coverage in `backend/internal/http/congregation_test.go`: production + callback → 404; nil callback → 404; development/staging + callback without cookie → 401; `postResetSeed(nil)` → 403; authenticated stub → 200.
4. Extract the `TRUNCATE` table list in `internal/seed` and test that `users` and `sessions` are not wiped.
5. Update living docs (`AGENTS.md`, `docs/deploy-dokploy.md`) so staging is allowlisted and production stays blocked.
6. Verify: `cd backend && go test ./...`
