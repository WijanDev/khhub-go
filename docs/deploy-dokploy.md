# Deploy on Hetzner + Dokploy

khhub is **not official Jehovah’s Witnesses software**. It does not talk to jw.org and does not replace S-21 cards or JW Hub. Dashboard totals are copied to the branch by hand.

Images are built on GitHub Actions and stored in GHCR. The VPS pulls; it does not `docker build` the app.

- SPA: `https://khhub.app`
- API: `https://api.khhub.app`
- Dokploy panel: `https://admin.khhub.app`

## Images

| Image | Source |
| --- | --- |
| `ghcr.io/wijandev/khhub-api` | `backend/` |
| `ghcr.io/wijandev/khhub-web` | `frontend/` (static files; `VITE_API_URL` baked at build) |
| `postgres:16-alpine` | official |

Tags: git SHA and `latest`. After the first push, set both GHCR packages to **public** if GitHub created them as private.

## Dokploy

1. Compose project pointed at [WijanDev/khhub-go](https://github.com/WijanDev/khhub-go) `main` (this repo’s `docker-compose.yml`).
2. Environment variables (never commit them):

   - `POSTGRES_USER` (for example `khhub`)
   - `POSTGRES_PASSWORD` — long and random
   - `POSTGRES_DB` (`khhub`)
   - `SESSION_SECRET` — 32+ random characters
   - `ADMIN_EMAIL`
   - `ADMIN_PASSWORD` — at least 10 characters in production
   - `APP_ENV=production`
   - `COOKIE_SECURE=true`
   - `CORS_ORIGINS=https://khhub.app`

3. Domains (HTTPS, Let’s Encrypt), DNS-only A records to the VPS IPv4 first:

   - `khhub.app` → `web` port 80
   - `api.khhub.app` → `api` port 8080

4. GitHub Actions secret `DOKPLOY_DEPLOY_HOOK` = the compose webhook URL (Deployments tab). Push to `main` then deploys.
5. Enable **Postgres backups** to an external destination the same day. This database stores congregation personal data.
6. Firewall: 22/80/443 only. Do not expose 5432 or 8080.
7. After the first login, change the admin password under **Congregación**.

Check `GET https://api.khhub.app/health` → `{"ok":true}`.

## Restore (outline)

Restoration steps should be written down when backups are first enabled: which Dokploy backup, how to stop `api`/`web`, how to restore the volume or `pg_restore`, and how to verify `/health`. Track the detailed runbook under ROADMAP → “Document restore-from-backup”.
