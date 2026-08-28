# Install and authenticate

Prefer the package manager already on the machine (`winget`, `scoop`, `choco`, `brew`, `apt`). Put `$(go env GOPATH)/bin` on `PATH` after any `go install`.

Dokploy panel for this project: `https://admin.wijan.dev`.

## Required

| Tool | Windows (Git Bash) | macOS | Linux |
| --- | --- | --- | --- |
| Git | `winget install Git.Git` | `brew install git` | `sudo apt install git` |
| Go 1.24+ | `winget install GoLang.Go` | `brew install go` | [go.dev/dl](https://go.dev/dl/) |
| Node 22+ | `winget install OpenJS.NodeJS.LTS` | `brew install node@22` | [nodejs.org](https://nodejs.org/) or NodeSource |
| Docker | `winget install Docker.DockerDesktop` then start Docker Desktop | `brew install --cask docker` | Docker Engine + Compose plugin |
| sqlc | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` | same | same |

Verify Docker with `docker info` (daemon must be running).

## Helpful

| Tool | Install |
| --- | --- |
| GNU make | Windows: `scoop install make` or `choco install make`. macOS/Linux: usually present; else `brew install make` / `sudo apt install make`. |
| air | `go install github.com/air-verse/air@latest` then `air` inside `backend/`. |

## Cloud CLIs

| Tool | Install | Auth (human, interactive) | Ready when |
| --- | --- | --- | --- |
| GitHub `gh` | `winget install GitHub.cli` / `brew install gh` | `gh auth login` (HTTPS, browser) | `gh auth status` and `gh repo view WijanDev/khhub-go` |
| Hetzner `hcloud` | `go install github.com/hetznercloud/cli/cmd/hcloud@latest` or `brew install hcloud` | `hcloud context create khhub` (paste a **read** token unless they need writes) | `hcloud context active` |
| Dokploy | `npm install -g @dokploy/cli` | API token from the panel (Settings → profile / API). Then `dokploy auth -u https://admin.wijan.dev -t` and let the human paste the token. | A read command succeeds, e.g. project list. Do **not** put the token in repo `.env`. |
| Cloudflare `wrangler` | `npm install -g wrangler` | `wrangler login` (browser) | `wrangler whoami` |

GHCR: after `gh auth login`, `gh auth token` can feed `docker login ghcr.io` if they need to pull images. Day-to-day coding does not need that.

R2 backups stay in the Dokploy UI. `wrangler` is for account identity and occasional R2 checks, not for pasting R2 access keys into git.

## Cursor MCP

Enable the GitHub, Cloudflare, and Snyk MCP servers in Cursor on the new PC. Authenticate there. Do not copy MCP tokens into this repo.

## Local app after CLIs

```bash
cp .env.example .env          # if missing; edit ADMIN_PASSWORD and SESSION_SECRET
cd frontend && npm ci
cd ..
make dev-db                   # or the docker compose line in README.md
cd backend && go test ./...
```
