# Implementation plan

1. In Dokploy project `khhub`, create environment `beta` (do not clone production volumes).
2. Duplicate the compose service into `beta`, same git source (`WijanDev/khhub-go`, `docker-compose.yml`). Override images if needed (`KHHUB_WEB_IMAGE` / `KHHUB_API_IMAGE`).
3. Give the compose its own env: new `POSTGRES_*`, `SESSION_SECRET`, `ADMIN_*`. Set `CORS_ORIGINS=https://beta.khhub.app`, `COOKIE_SECURE=true`. Set `APP_ENV` so `/dev/reset-seed` is **not** registered (not `development`).
4. DNS (Cloudflare, grey cloud): `beta.khhub.app` and `api.beta.khhub.app` A → `138.201.156.89`.
5. Dokploy domains: `beta.khhub.app` → `web:80`, `api.beta.khhub.app` → `api:8080`, HTTPS Let’s Encrypt. Pass path `/` without Git Bash rewriting it.
6. Bake a **web** image whose `VITE_API_URL` is `https://api.beta.khhub.app` (today that value is compile-time). Add a CI job or workflow input that pushes e.g. `ghcr.io/wijandev/khhub-web:beta` without moving production `:latest`.
7. Deploy `beta` only. Confirm `GET https://api.beta.khhub.app/health` and login on `https://beta.khhub.app`. Confirm production hosts still hit the production compose.
8. Document hosts and “no prod volume” in `docs/deploy-dokploy.md`.
