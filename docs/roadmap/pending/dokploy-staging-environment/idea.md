# Dokploy staging environment

- **Slug:** dokploy-staging-environment
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Run a second khhub stack on the same Hetzner cx23 and the same Dokploy panel (`admin.wijan.dev`), in a `staging` environment next to `production`. Testers use `staging.khhub.app` / `apistaging.khhub.app` against their own Postgres. Production data and volumes stay untouched.

Build images once on `dev` (immutable SHA tags) and deploy that pair to staging. A release PR `dev` → `main` does **not** rebuild: retag the same digests and webhook production so it points at the image already exercised on staging.

## Hosts

| Role | Hostname |
| --- | --- |
| Web | `staging.khhub.app` |
| API | `apistaging.khhub.app` |

## Out of scope

- A second Dokploy install or a second VPS
- `APP_ENV=development` or `POST /dev/reset-seed` on a public hostname
- Sharing the production `khhub_pg` volume
- Preview deploys for every pull request
- Real congregation PII in the staging database (fictional seed only)
- Deploying production from `dev` (production moves only via the `main` promote webhook)
