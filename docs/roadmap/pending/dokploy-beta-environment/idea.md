# Dokploy beta environment

- **Slug:** dokploy-beta-environment
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Run a second khhub stack on the same Hetzner cx23 and the same Dokploy panel (`admin.wijan.dev`), in a new environment next to `production`. Testers use `beta.khhub.app` / `api.beta.khhub.app` against their own Postgres. Production data and volumes stay untouched.

## Out of scope

- A second Dokploy install or a second VPS
- `APP_ENV=development` or `POST /dev/reset-seed` on a public hostname
- Sharing the production `khhub_pg` volume
- Preview deploys for every pull request
- Real congregation PII in the beta database (fictional seed only)
