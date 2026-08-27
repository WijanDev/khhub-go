# Deploy on Hetzner + Dokploy

khhub is **not official Jehovah’s Witnesses software**. It does not talk to jw.org and does not replace S-21 cards or JW Hub. Dashboard totals are copied to the branch by hand.

Images are built on GitHub Actions and stored in GHCR. The VPS pulls; it does not `docker build` the app.

- SPA: `https://khhub.app`
- API: `https://api.khhub.app`
- Staging SPA: `https://staging.khhub.app`
- Staging API: `https://apistaging.khhub.app`
- Dokploy panel: `https://admin.wijan.dev` (personal; not a khhub hostname)

## Images

| Image | Source |
| --- | --- |
| `ghcr.io/wijandev/khhub-api` | `backend/` |
| `ghcr.io/wijandev/khhub-web` | `frontend/` (static files; API origin from `KHHUB_API_URL` at container start) |
| `postgres:16-alpine` | Dokploy database per environment (not in the app images) |

CI publishes an image only when that image’s source tree changed (`backend/` or `frontend/`, excluding tests). Docs, roadmap, and other non-image paths do not push to GHCR. A frontend-only change never publishes the API image, and the reverse. Tags: git SHA and `staging` on a matching push to `dev`; a matching merge to `main` retags that `:staging` digest as `:latest` (no rebuild). After the first push, set both GHCR packages to **public** if GitHub created them as private.

## Dokploy

Project `khhub` keeps environments `staging` and `production`. Deploy as two **applications**, not one Compose stack:

| Application | Source | Staging image | Production image |
| --- | --- | --- | --- |
| `khhub-api` | Docker (GHCR), not a git build | `ghcr.io/wijandev/khhub-api:staging` | `ghcr.io/wijandev/khhub-api:latest` |
| `khhub-web` | Docker (GHCR), not a git build | `ghcr.io/wijandev/khhub-web:staging` | `ghcr.io/wijandev/khhub-web:latest` |

Create one **Dokploy Postgres** per environment. Point `khhub-api` at it with `DATABASE_URL`. `khhub-web` has no database. **Do not clone production volumes** into staging.

`docker-compose.yml` in this repo is an optional full-stack reference for a local prod-like run. Dokploy does not deploy it.

1. Environments and apps as in the table. `APP_ENV=staging` loads the fictional demo seed on an empty directory and does **not** register `POST /dev/reset-seed`.

   | Environment | `KHHUB_API_URL` | `CORS_ORIGINS` | `APP_ENV` |
   | --- | --- | --- | --- |
   | `production` | `https://api.khhub.app` | `https://khhub.app` | `production` |
   | `staging` | `https://apistaging.khhub.app` | `https://staging.khhub.app` | `staging` |

2. Environment variables (never commit them). Each environment has its own Postgres, `SESSION_SECRET`, and `ADMIN_*`.

   On `khhub-api`:

   - `DATABASE_URL` — internal Dokploy Postgres URL (`sslmode=disable` on the VPS network). Do not expose 5432.
   - `SESSION_SECRET` — 32+ random characters
   - `ADMIN_EMAIL`
   - `ADMIN_PASSWORD` — at least 10 characters in production and staging
   - `APP_ENV` — `production` or `staging` (never `development` on a public hostname)
   - `COOKIE_SECURE=true`
   - `CORS_ORIGINS` — see table

   On `khhub-web`:

   - `KHHUB_API_URL` — see table

3. Domains (HTTPS, Let’s Encrypt), DNS-only A records to the VPS IPv4 first:

   - `khhub.app` → production `khhub-web` port 80
   - `api.khhub.app` → production `khhub-api` port 8080
   - `staging.khhub.app` → staging `khhub-web` port 80
   - `apistaging.khhub.app` → staging `khhub-api` port 8080

4. GitHub Actions secrets (do not commit them). Each app/environment has its own Dokploy webhook:

   - `DOKPLOY_STAGING_API_HOOK` — push to `dev` that touches API paths builds `:sha` + `:staging` and fires this hook
   - `DOKPLOY_STAGING_WEB_HOOK` — same for the web image
   - `DOKPLOY_PROD_API_HOOK` — a release PR `dev` → `main` that touches API paths retags `:staging` → `:latest` and fires this hook
   - `DOKPLOY_PROD_WEB_HOOK` — same for the web image

   A docs-only or tests-only merge publishes nothing and fires no hook. Remove the old `DOKPLOY_STAGING_DEPLOY_HOOK` and `DOKPLOY_DEPLOY_HOOK` secrets after the four new hooks work.

5. Postgres dumps go to the Cloudflare R2 bucket `khhub-backups` (EU). Add an S3 destination in Dokploy and a **weekly Sunday midnight** backup of each environment’s Dokploy Postgres, database `khhub`. Production stores congregation personal data.
6. Firewall: 22/80/443 only. Do not expose 5432 or 8080.
7. After the first login, change the admin password under **Congregación**.

### Cut over from the old Compose stack

Do this once, staging first, then production. Keep the old Compose running until the new apps pass health checks.

1. In Dokploy, create Postgres in that environment. Do not attach the production volume to staging.
2. Dump the existing Compose Postgres (`pg_dump` of database `khhub`) and restore into the new Dokploy Postgres. Confirm row counts before switching traffic.
3. Create `khhub-api` and `khhub-web` (Docker provider, image tags from the table). Paste env vars. Set `DATABASE_URL` to the new internal URL.
4. Attach the domains to the new apps (or move them from the Compose service). Check `GET /health` and login.
5. Copy each app’s deploy webhook into the matching GitHub secret.
6. Stop the old Compose stack. Do not delete the old volume until a successful R2 dump exists from the new Postgres.
7. Point the weekly backup at the Dokploy database, not the old Compose `postgres` service.

Check `GET https://api.khhub.app/health` → `{"ok":true}`. Staging: `GET https://apistaging.khhub.app/health`. A release to `main` only works after `dev` has published a `:staging` tag for that image at least once.

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

Database backup (each environment’s Dokploy Postgres → Backups):

- Type: Postgres (Dokploy database, not the retired Compose service)
- Database: `khhub`
- Schedule: **Sunday midnight** as set in the Dokploy UI (typical cron `0 0 * * 0`). Dokploy cron is usually **UTC** — Sunday 00:00 UTC is 02:00 CEST / 01:00 CET. Confirm the timezone shown in the panel.
- Prefix: Dokploy chooses a database prefix. After the first dump, record the key pattern here; do not assume `khhub/postgres/`.
- Keep latest: `14` (about three months of weekly dumps if unchanged)

R2 free tier covers 10 GB-month and well above our dump size. Overages are Cloudflare R2 catalog rates, not Hetzner.

Optional extra (Hetzner, billed): whole-VPS automatic backups are **20% of the cx23** — about **1.10 € net / 1.33 € gross per month** at the current `fsn1` catalog price (`hcloud server-type describe cx23`). That is a disk snapshot, not `pg_dump`. Enable only if you want a second copy on Hetzner: `hcloud server enable-backup khhub`.

## Restore

1. In Dokploy, open the environment’s Postgres → **Backups** → **Restore**.
2. Pick destination `r2-khhub-backups` and the dump.
3. Database name: `khhub`.
4. Stop `khhub-api` (and `khhub-web` if you want the UI offline) before restore if Dokploy does not do it.
5. Run restore. Dokploy uses `pg_dump -Fc` + gzip on the way out; its Restore button uses the matching restore command.
6. Start `khhub-api` again.
7. Check `GET https://api.khhub.app/health` → `{"ok":true}` and log in at `https://khhub.app`.

If you only have a Hetzner server backup: restore that snapshot in the Hetzner console (rolls the whole VPS back), then check `/health`. Do not treat that as a point-in-time Postgres restore.
