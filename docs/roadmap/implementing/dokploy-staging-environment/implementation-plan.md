# Implementation plan

1. In Dokploy project `khhub`, create environment `staging` (do not clone production volumes).
2. Duplicate the compose service into `staging`, same git source (`WijanDev/khhub-go`, `docker-compose.yml`). Pin images by **SHA digest or immutable tag**, not a floating `:latest` that production also follows.
3. Give the compose its own env: new `POSTGRES_*`, `SESSION_SECRET`, `ADMIN_*`. Set `CORS_ORIGINS=https://staging.khhub.app`, `COOKIE_SECURE=true`. Set `APP_ENV` so `/dev/reset-seed` is **not** registered (not `development`).
4. DNS (Cloudflare, grey cloud): `staging.khhub.app` and `apistaging.khhub.app` A → `138.201.156.89`.
5. Dokploy domains: `staging.khhub.app` → `web:80`, `apistaging.khhub.app` → `api:8080`, HTTPS Let’s Encrypt. Pass path `/` without Git Bash rewriting it.
6. **API origin on the SPA.** Today `VITE_API_URL` is baked at image build. Prefer a **runtime** config (for example a small `config.js` or static-server env) so the same web digest can serve `https://apistaging.khhub.app` on staging and `https://api.khhub.app` on production. If runtime is too large for the first slice, CI must build **two** web images from the same commit (staging bake vs production bake) and promote that pair — not “the last `:latest`”.
7. **CI on push to `dev`:** after `secrets` + `test`, build and push `khhub-api` and `khhub-web` tagged with `github.sha` (and a `staging` pointer if useful). Do **not** move production `:latest` here. Fire the **staging** Dokploy deploy hook so staging pulls those SHA tags.
8. **CI on push to `main` (release):** do **not** rebuild. Retag the SHA already on GHCR as the production pointer (`:latest` or a `prod` tag) and fire the **production** Dokploy hook so production pulls that same digest. Production and staging must not share one compose webhook.
9. Verify: `GET https://apistaging.khhub.app/health` and login on `https://staging.khhub.app`. Confirm `khhub.app` / `api.khhub.app` still hit the production compose. Confirm a merge to `dev` updates staging only; a merge to `main` updates production to the **same** image digest (or the matching prod-baked web if dual-build).
10. Document hosts, hooks, pin-by-SHA, and “no prod volume” in `docs/deploy-dokploy.md`. Store `DOKPLOY_STAGING_DEPLOY_HOOK` (and keep `DOKPLOY_DEPLOY_HOOK` for production) as GitHub Actions secrets; do not commit them.
