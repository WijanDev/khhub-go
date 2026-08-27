# Implementation plan

1. In `.github/workflows/ci.yml`, add a path-filter step (for example `dorny/paths-filter`) on the `images` and `promote` jobs. Compare against the previous commit on that branch (`github.event.before`). Treat an empty / zero `before` SHA as “both images changed.”
2. Define two filters:
   - **API:** `backend/` except `*_test.go` and testdata.
   - **Web:** `frontend/` except test/spec files.
   Everything else (`docs/`, `docs/roadmap/`, `.agents/`, README, and similar) is ignored.
3. On push to `dev` (after `secrets` + `test`): if the API filter matches, build and push `ghcr.io/wijandev/khhub-api:$sha` and `:staging`, then POST `DOKPLOY_STAGING_API_HOOK`. If the web filter matches, the same for `khhub-web` and `DOKPLOY_STAGING_WEB_HOOK`. If neither matches, skip both publishes and hooks.
4. On push to `main` (no rebuild): if the API filter matches, `docker buildx imagetools create` `khhub-api:staging` → `:latest` and POST `DOKPLOY_PROD_API_HOOK`. Same for web with `DOKPLOY_PROD_WEB_HOOK`. If neither matches, skip `promote`.
5. Add four GitHub Actions secrets (do not commit URLs): `DOKPLOY_STAGING_API_HOOK`, `DOKPLOY_STAGING_WEB_HOOK`, `DOKPLOY_PROD_API_HOOK`, `DOKPLOY_PROD_WEB_HOOK`. Remove use of `DOKPLOY_STAGING_DEPLOY_HOOK` and `DOKPLOY_DEPLOY_HOOK` from the workflow once the new hooks work.
6. In Dokploy project `khhub`, keep environments `staging` and `production`. Create applications `khhub-api` and `khhub-web` (Docker image from GHCR, not a git build). Attach each app to both environments. Do not clone production volumes into staging.
7. Create one Dokploy Postgres per environment. Dump/restore (or attach) the existing `khhub_pg` data so production congregation data and the staging fictional seed stay on their own disks. Point `khhub-api` at that instance with `DATABASE_URL`. `khhub-web` has no database. Keep current env vars (`APP_ENV`, `CORS_ORIGINS`, `COOKIE_SECURE`, `SESSION_SECRET`, `ADMIN_*`, `KHHUB_API_URL`).
8. Move domains onto the new apps: `api.khhub.app` / `apistaging.khhub.app` → API `:8080`; `khhub.app` / `staging.khhub.app` → web `:80`. Leave the old Compose stack stopped only after health checks pass.
9. Update `docs/deploy-dokploy.md` (two apps, four hooks, Dokploy Postgres). Update `README.md` and `AGENTS.md` if they still describe one Compose webhook. Keep `docker-compose.yml` as an optional full-stack reference; Dokploy no longer runs it.
10. Verify: a docs-only push to `dev` publishes nothing and fires no hook; a frontend-only push publishes and deploys only web on staging; a backend-only merge to `main` retags and deploys only API on production. `GET /health` and login still work on both environments.
