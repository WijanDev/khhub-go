# Tests that seed reset is blocked in production

- **Slug:** seed-reset-production-guard-tests
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Prove `POST /dev/reset-seed` is absent or rejected when `APP_ENV=production` and that reset keeps the admin user in development.

## Out of scope

- Removing the reset button from the UI.
