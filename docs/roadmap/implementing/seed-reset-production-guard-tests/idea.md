# Tests that seed reset is blocked in production

- **Slug:** seed-reset-production-guard-tests
- **Status:** in-progress
- **Merge date:**
- **App version:**

## Summary

`POST /dev/reset-seed` is an explicit allowlist (`development` and `staging`). Production never registers the route, even if a reset callback is passed in. Reset wipes demo congregation rows and keeps the admin user. Staging may use the button to restore the fictional seed; do not put real congregation PII there.

## Out of scope

- Removing the reset button from the UI.
- Enabling reset when `APP_ENV=production`.
- Postgres integration tests or testcontainers.
