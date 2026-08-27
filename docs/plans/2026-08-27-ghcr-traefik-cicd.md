# GHCR + Traefik CI/CD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `.agents/skills/executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build khhub images in GitHub Actions, store them in GHCR, serve the SPA at `khhub.app` and the API at `api.khhub.app` with no `/api` path prefix, and stop building on the VPS.

**Architecture:** Actions push `ghcr.io/wijandev/khhub-api` and `khhub-web`. Dokploy/Traefik routes two hostnames. The browser calls `VITE_API_URL` with credentialed fetch. Postgres stays a Compose service on the VPS.

**Tech Stack:** Go 1.24, Gin, Vite, GitHub Actions, GHCR, Dokploy, Traefik, PostgreSQL 16.

**Spec:** [docs/plans/2026-08-27-ghcr-traefik-cicd-design.md](2026-08-27-ghcr-traefik-cicd-design.md)

## Global Constraints

- UI copy stays Spanish; docs, comments, API field names stay English.
- Auth is httpOnly cookie `khhub_session`. No JWT. No public signup.
- Hours rules (`domain.ReportsHours`), service year, derived activity: unchanged.
- Never commit `.env` or Dokploy/GitHub secrets.
- Do not expose 5432 or 8080. Production: `APP_ENV=production`, `COOKIE_SECURE=true`, `CORS_ORIGINS=https://khhub.app`.
- GHCR names are lowercase: `ghcr.io/wijandev/khhub-api`, `ghcr.io/wijandev/khhub-web`.
- `backend/cmd/api` is the process entry; do not rename it.

---

### Task 1: Gin routes without `/api`

**Files:**
- Modify: `backend/internal/http/router.go`
- Create: `backend/internal/http/router_test.go`

**Interfaces:**
- Consumes: existing handlers in `router.go`
- Produces: `GET /health`, `POST /auth/login`, `POST /auth/logout`, authed group at `/` (`/auth/me`, `/congregation`, …)

- [ ] **Step 1: Write failing handler tests**

```go
package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"khhub/internal/config"
)

func TestHealthHasNoAPIPrefix(t *testing.T) {
	h := NewRouter(config.Config{AppEnv: "development"}, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health: got %d", rec.Code)
	}
	old := httptest.NewRecorder()
	h.ServeHTTP(old, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if old.Code == http.StatusOK {
		t.Fatal("GET /api/health must not remain mounted")
	}
}
```

`NewRouter` currently requires a live `*store.Queries` for session middleware. If `nil` panics, construct the smallest test helper the file already needs (or skip session by only hitting `/health` before that middleware — today session middleware is registered before `/health` wait: `/health` is registered BEFORE `sessionMiddleware`. Good. But `NewRouter(..., nil, nil)` will still compile; session middleware is not used on `/health`.

- [ ] **Step 2: Run the test**

```bash
cd backend && go test ./internal/http -run TestHealthHasNoAPIPrefix -v
```

Expected: FAIL (404 on `/health` or compile if you need to adjust the test).

- [ ] **Step 3: Change routes**

In `router.go`:

- `r.GET("/health", ...)`
- `r.POST("/auth/login", ...)` and `r.POST("/auth/logout", ...)`
- `authed := r.Group("/")` with `requireAuth()`, same child paths (`/auth/me`, `/congregation`, `/dev/reset-seed`, …).

Do not leave a compatibility `/api` group.

- [ ] **Step 4: Re-run tests**

```bash
cd backend && go test ./internal/http -run TestHealthHasNoAPIPrefix -v && go test ./...
```

Expected: PASS.

- [ ] **Step 5: Commit** (only if the user asked to commit in that session)

```bash
git add backend/internal/http/router.go backend/internal/http/router_test.go
git commit -m "fix: serve API routes at the host root, not under /api"
```

---

### Task 2: Frontend client uses `VITE_API_URL`

**Files:**
- Modify: `frontend/src/lib/api.ts`
- Modify: every `api("/api/...")` call site under `frontend/src/`
- Modify: `frontend/vite.config.ts`
- Modify: `frontend/.env.example` if present, else repo `.env.example`

**Interfaces:**
- Consumes: `import.meta.env.VITE_API_URL` (string, no trailing slash)
- Produces: `api("/congregation")` → `fetch(`${base}/congregation`, { credentials: "include" })`

- [ ] **Step 1: Change `api.ts`**

```ts
function apiURL(path: string): string {
  const base = (import.meta.env.VITE_API_URL ?? "").replace(/\/$/, "");
  if (!path.startsWith("/")) {
    throw new Error("api() path must start with /");
  }
  return `${base}${path}`;
}

export async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (init?.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  const res = await fetch(apiURL(path), {
    ...init,
    credentials: "include",
    headers,
  });
  // ... existing status / json handling unchanged
}
```

- [ ] **Step 2: Replace call sites**

All `"/api/auth/login"` → `"/auth/login"`, `"/api/congregation"` → `"/congregation"`, `"/api/dev/reset-seed"` → `"/dev/reset-seed"`, and the same for households, publishers, reports, attendance, dashboard. Search `"/api/` in `frontend/src` until zero matches.

- [ ] **Step 3: Drop the Vite proxy**

`frontend/vite.config.ts`: remove `server.proxy`. Local browser talks to `http://127.0.0.1:8080` via `VITE_API_URL`.

- [ ] **Step 4: Document env**

In `.env.example` add:

```env
# Vite (frontend). No trailing slash. Local: API on the host. Production image bake: https://api.khhub.app
VITE_API_URL=http://127.0.0.1:8080
```

Change the `CORS_ORIGINS` comment: production is **not** empty; it must be `https://khhub.app`.

- [ ] **Step 5: Verify locally**

```bash
# API already allows CORS_ORIGINS=http://localhost:5173
cd frontend && echo VITE_API_URL=http://127.0.0.1:8080 > .env.development.local
npm run lint && npm run build
```

`.env.development.local` is gitignored. Do not commit it.

Manual: login on `http://localhost:5173` against `http://127.0.0.1:8080`.

---

### Task 3: Static web image (no nginx API proxy)

**Files:**
- Modify: `frontend/Dockerfile`
- Delete: `frontend/nginx.conf` (only if nothing else references it)
- Create: `frontend/static-web-server.toml` (or equivalent config next to the Dockerfile)

**Interfaces:**
- Produces: container listens on 80, serves `dist/`, SPA fallback to `index.html`, security headers comparable to the old nginx file (nosniff, DENY frame, same-origin referrer).

Use `ghcr.io/static-web-server/static-web-server:2` (or pin a digest). Example Dockerfile:

```dockerfile
FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM ghcr.io/static-web-server/static-web-server:2
COPY --from=build /app/dist /public
COPY static-web-server.toml /config.toml
ENV SERVER_CONFIG_FILE=/config.toml
EXPOSE 80
```

Configure page fallback to `/index.html`. Do not add a reverse-proxy location.

- [ ] **Step 1: Write config + Dockerfile**
- [ ] **Step 2: Build locally (optional if Docker is up)**

```bash
docker build -t khhub-web:test --build-arg VITE_API_URL=https://api.khhub.app ./frontend
```

Expected: image well under the old ~50 MB nginx image.

---

### Task 4: Production Compose uses GHCR images

**Files:**
- Modify: `docker-compose.yml`

Replace `api` and `web` `build:` with:

```yaml
  api:
    image: ${KHHUB_API_IMAGE:-ghcr.io/wijandev/khhub-api:latest}
    # existing environment, depends_on, expose 8080, no ports
  web:
    image: ${KHHUB_WEB_IMAGE:-ghcr.io/wijandev/khhub-web:latest}
    # no depends_on api required for HTTP; keep restart/networks
    expose:
      - "80"
```

Keep `postgres` as `postgres:16-alpine` and `khhub_pg`. Do not publish 5432 or 8080.

`CORS_ORIGINS` stays an env var (Dokploy sets `https://khhub.app`).

- [ ] **Step 1: Edit compose**
- [ ] **Step 2: `docker compose config`** to validate YAML

---

### Task 5: GitHub Actions — test, build, push, deploy hook

**Files:**
- Create: `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read
  packages: write

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
          cache-dependency-path: backend/go.sum
      - uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: npm
          cache-dependency-path: frontend/package-lock.json
      - name: Go tests
        working-directory: backend
        run: go test ./...
      - name: Frontend lint and build
        working-directory: frontend
        env:
          VITE_API_URL: https://api.khhub.app
        run: npm ci && npm run lint && npm run build

  images:
    needs: test
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/setup-buildx-action@v3
      - name: Build and push API
        uses: docker/build-push-action@v6
        with:
          context: ./backend
          push: true
          tags: |
            ghcr.io/wijandev/khhub-api:${{ github.sha }}
            ghcr.io/wijandev/khhub-api:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
      - name: Build and push web
        uses: docker/build-push-action@v6
        with:
          context: ./frontend
          push: true
          build-args: |
            VITE_API_URL=https://api.khhub.app
          tags: |
            ghcr.io/wijandev/khhub-web:${{ github.sha }}
            ghcr.io/wijandev/khhub-web:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
      - name: Trigger Dokploy
        env:
          HOOK: ${{ secrets.DOKPLOY_DEPLOY_HOOK }}
        run: |
          if [ -z "$HOOK" ]; then
            echo "DOKPLOY_DEPLOY_HOOK not set; skip deploy trigger"
            exit 0
          fi
          curl -fsS -X POST "$HOOK"
```

GitHub secret `DOKPLOY_DEPLOY_HOOK` is the compose webhook URL from Dokploy (or a small script that POSTs `compose.deploy` with `x-api-key`). Do not put the hook in the repo.

- [ ] **Step 1: Add the workflow**
- [ ] **Step 2: After first green push, set package visibility to Public** on both GHCR packages (GitHub → Packages) if they default to private.

---

### Task 6: Docs and roadmap

**Files:**
- Modify: `AGENTS.md` (CORS empty → `https://khhub.app` in production; health is `/health`; no `/api` prefix)
- Modify: `docs/deploy-dokploy.md` (GHCR images, `khhub.app` + `api.khhub.app`, env list, verify `/health`)
- Modify: `README.md` (Vite needs `VITE_API_URL`; API is `:8080`, not proxied)
- Modify: `ROADMAP.md` — tick only what this change ships: GitHub Actions test+build, production CORS/cookie note. Leave golangci-lint, backups runbook, handler auth tests (beyond `/health`) unchecked unless those tests were added.

- [ ] **Step 1: Edit the four files in English**
- [ ] **Step 2: Repeat the non-official disclaimer in README and deploy doc**

---

### Task 7: Dokploy and DNS (ops, after code is on `main`)

Not a repo commit. Do this from Cloudflare + Dokploy CLI/UI.

- [ ] **Step 1:** DNS-only A records: `khhub.app` and `api.khhub.app` → `138.201.156.89`
- [ ] **Step 2:** Dokploy compose env: production secrets listed in the spec
- [ ] **Step 3:** Domains on the compose services: `khhub.app` → web:80 LE; `api.khhub.app` → api:8080 LE
- [ ] **Step 4:** Confirm webhook/hook secret in GitHub Actions
- [ ] **Step 5:** First deploy, then `curl -sS https://api.khhub.app/health` and login on `https://khhub.app`

Do not close or reopen port 3000. Do not dockerize the API/Vite for day-to-day local work.

---

## Self-review

- Spec coverage: routes, CORS/cookie, GHCR names, Actions, Compose, Dokploy/DNS, docs — each has a task.
- No `/api` compatibility shim.
- Image names lowercase.
- Deploy hook is a secret, not source.
