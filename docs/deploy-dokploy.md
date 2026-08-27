# Deploy on Hetzner + Dokploy

khhub is **not official Jehovah’s Witnesses software**. It does not talk to jw.org and does not replace S-21 cards or JW Hub. Dashboard totals are copied to the branch by hand.

Images are built on GitHub Actions and stored in GHCR. The VPS pulls; it does not `docker build` the app.

- SPA: `https://khhub.app`
- API: `https://api.khhub.app`
- Dokploy panel: `https://admin.wijan.dev` (personal; not a khhub hostname)

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

4. GitHub Actions secret `DOKPLOY_DEPLOY_HOOK` = the compose webhook URL from **https://admin.wijan.dev** → khhub compose → Deployments. Copy it there (do not commit it). A merge to `main` (release PR from `dev`) then deploys.
5. Postgres dumps go to the Cloudflare R2 bucket `khhub-backups` (EU). Add an S3 destination in Dokploy and a **weekly Sunday midnight** compose backup of service `postgres`, database `khhub`. This database stores congregation personal data.
6. Firewall: 22/80/443 only. Do not expose 5432 or 8080.
7. After the first login, change the admin password under **Congregación**.

Check `GET https://api.khhub.app/health` → `{"ok":true}`.

## Backups (R2)

Destination (Dokploy → Settings → Destinations):

| Field | Value |
| --- | --- |
| Name | `r2-khhub-backups` |
| Provider | S3 |
| Bucket | `khhub-backups` |
| Region | `auto` |
| Endpoint | `https://6e9f13c95076e80e4350c4d4fcbc4b85.r2.cloudflarestorage.com` |
| Access key / secret | Cloudflare R2 API token (Object Read and Write, scoped to `khhub-backups`). Create it in R2 → Manage API Tokens. Paste into Dokploy, not into git or chat. |

Compose backup (khhub → Backups):

- Type: compose / Postgres
- Service: `postgres`
- Database: `khhub`
- Schedule: **Sunday midnight** as set in the Dokploy UI (typical cron `0 0 * * 0`). Dokploy cron is usually **UTC** — Sunday 00:00 UTC is 02:00 CEST / 01:00 CET. Confirm the timezone shown in the panel.
- Prefix: Dokploy chooses a compose-service prefix (first verified dump: `compose-…_postgres/YYYY-MM-DD….sql.gz`). Do not assume `khhub/postgres/`.
- Keep latest: `14` (about three months of weekly dumps if unchanged)

R2 free tier covers 10 GB-month and well above our dump size. Overages are Cloudflare R2 catalog rates, not Hetzner.

Optional extra (Hetzner, billed): whole-VPS automatic backups are **20% of the cx23** — about **1.10 € net / 1.33 € gross per month** at the current `fsn1` catalog price (`hcloud server-type describe cx23`). That is a disk snapshot, not `pg_dump`. Enable only if you want a second copy on Hetzner: `hcloud server enable-backup khhub`.

## Restore

1. In Dokploy, open the khhub compose → **Backups** → **Restore**.
2. Pick destination `r2-khhub-backups` and the dump (keys look like `compose-…_postgres/*.sql.gz`).
3. Database name: `khhub`.
4. Stop `api` (and `web` if you want the UI offline) before restore if Dokploy does not do it: compose → Stop, or `docker stop` the `api` container on the VPS.
5. Run restore. Dokploy uses `pg_dump -Fc` + gzip on the way out; its Restore button uses the matching restore command.
6. Start the compose again.
7. Check `GET https://api.khhub.app/health` → `{"ok":true}` and log in at `https://khhub.app`.

If you only have a Hetzner server backup: restore that snapshot in the Hetzner console (rolls the whole VPS back), then check `/health`. Do not treat that as a point-in-time Postgres restore.
