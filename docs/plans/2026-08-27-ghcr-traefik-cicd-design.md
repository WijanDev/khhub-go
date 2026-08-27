# Design: GHCR + Traefik subdomains (no `/api` prefix)

**Date:** 2026-08-27  
**Status:** draft — review before implementation  
**Product:** khhub (not official JW software; no jw.org; UI in Spanish, docs in English)

## Problem

Production Compose currently `build:`s the API and SPA on the Hetzner `cx23`. That burns RAM and can stall Dokploy. The SPA nginx container also reverse-proxies `/api` so the browser sees one origin. That coupling makes a fat web image and hides the API behind a path prefix.

## Goals

1. Build images in GitHub Actions. Store them in GHCR. The VPS only pulls.
2. Expose the SPA at `https://khhub.app` and the API at `https://api.khhub.app` via Dokploy/Traefik.
3. Remove the `/api` URL prefix from Gin and from the frontend client. The host is the API.
4. Drop the nginx `/api` proxy. The web image is static files plus SPA fallback only.
5. Do not bake congregation secrets into images. Do not expose 5432 or 8080 on the public internet.

## Non-goals

- Self-hosted registry on the VPS.
- Building images on the `cx23`.
- Changing auth (still httpOnly cookie `khhub_session`, no JWT, no public signup).
- First production data migration or Postgres backups runbook (still a roadmap item).
- Multi-arch images (amd64 only; matches the VPS).

## Architecture

```
Browser
  ├─ https://khhub.app          → Traefik → web (static)
  └─ https://api.khhub.app      → Traefik → api :8080
https://admin.khhub.app         → Dokploy panel (unchanged)

GitHub main
  → Actions build/push
  → ghcr.io/wijandev/khhub-api:<sha>|latest
  → ghcr.io/wijandev/khhub-web:<sha>|latest
  → Dokploy compose pull + up (webhook or API)
```

Local development stays on the host: `make dev-db`, `go run ./cmd/api`, Vite on `:5173`. The API origin in the browser is `http://127.0.0.1:8080` with CORS. No Vite `/api` proxy.

## HTTP contract (after)

Cookie: `khhub_session`, Path `/`, host-only on the API host, `Secure` in production, `SameSite=Lax`. Sent on credentialed fetches to `api.khhub.app` (same eTLD+1).

CORS production: `CORS_ORIGINS=https://khhub.app` (exact origin, credentials allowed).  
CORS local: `CORS_ORIGINS=http://localhost:5173` (already in `.env.example`).

| Method | Path |
| --- | --- |
| GET | `/health` |
| POST | `/auth/login` |
| POST | `/auth/logout` |
| GET | `/auth/me` |
| POST | `/auth/change-password` |
| GET, PUT | `/congregation` |
| POST | `/dev/reset-seed` (non-production only) |
| GET, POST | `/households` |
| PUT, DELETE | `/households/:id` |
| GET, POST | `/publishers` |
| GET, PUT, DELETE | `/publishers/:id` |
| GET, PUT | `/reports` |
| GET, PUT | `/attendance` |
| DELETE | `/attendance/:id` |
| GET | `/dashboard` |

No redirect from `/api/...`. Old paths are gone.

Frontend: `VITE_API_URL` (no trailing slash) + paths like `api("/congregation")`. Production image build must set `VITE_API_URL=https://api.khhub.app`.

## Images

| Image | Context | Notes |
| --- | --- | --- |
| `ghcr.io/wijandev/khhub-api` | `backend/` | Existing multi-stage Go Dockerfile. |
| `ghcr.io/wijandev/khhub-web` | `frontend/` | Node build stage, then a small static server with SPA fallback. No `/api` proxy. |

Tags per build: git SHA (immutable) and `latest` (Dokploy default pull). Packages are public (repo is public). GHCR container storage is currently unbilled; treat that as a policy, not a contract.

`docker-compose.yml` (Dokploy / prod-like): `image:` for `api` and `web`, `postgres:16-alpine` unchanged, named volume for Postgres. No `build:`.

`docker-compose.dev.yml`: unchanged (Postgres only).

## CI

Workflow on `push`/`pull_request` to `main`:

1. `go test ./...` from `backend/`.
2. `npm ci && npm run lint && npm run build` from `frontend/` (PR can use a dummy `VITE_API_URL`).
3. On push to `main` only: build/push both images to GHCR with `GITHUB_TOKEN`.
4. Trigger Dokploy compose deploy (webhook URL or `x-api-key` + compose deploy). Secret lives in GitHub Actions, not in the repo.

Do not start the production deploy until env vars exist in Dokploy (see below).

## Dokploy / DNS

Dokploy already has project `khhub`, environment `production`, compose service pointed at `WijanDev/khhub-go` `main`. After this work it must:

- Use the repo Compose file that references GHCR images (still git pull of the compose file, or paste — git pull of compose is enough; containers come from GHCR).
- Register GHCR in Settings → Registry only if packages are private. Public packages need no login.
- Domains: `khhub.app` → `web` port 80, HTTPS LE; `api.khhub.app` → `api` port 8080, HTTPS LE. Grey-cloud DNS A records to `138.201.156.89` first.
- Env (Dokploy compose `.env`, never committed): `POSTGRES_*`, `SESSION_SECRET`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`, `APP_ENV=production`, `COOKIE_SECURE=true`, `CORS_ORIGINS=https://khhub.app`.

Firewall stays 22/80/443.

## Security

- Images contain no passwords, session secrets, or congregation PII.
- Seed reset stays off when `APP_ENV=production`.
- Do not publish 8080 or 5432.
- Cookie remains httpOnly; CORS is an allow-list, not `*`.

## Docs to update in the same change

- `AGENTS.md` (health path, CORS, no `/api` prefix).
- `docs/deploy-dokploy.md` (GHCR, domains, env).
- `README.md` (local API origin, no Vite proxy).
- `.env.example` (document `VITE_API_URL` for Vite).
- `ROADMAP.md` (tick CI / cookie-CORS items that this change actually ships).

## Success

- `curl -sS https://api.khhub.app/health` → `{"ok":true}`.
- Browser on `https://khhub.app` logs in; session cookie is set for `api.khhub.app`.
- Push to `main` publishes new images and Dokploy pulls them. The VPS does not run `docker build` for khhub.
- Local `make api` + `make web` still works against `localhost:5173` → `127.0.0.1:8080`.
