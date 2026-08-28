---
name: install-repo
description: Installs and verifies khhub local and cloud CLIs (Git, Go, Node, Docker, sqlc, gh, hcloud, Dokploy, Wrangler). Use when setting up a new machine, a personal PC, a weekend checkout, missing tools, or GitHub / Dokploy / Cloudflare / Hetzner auth.
---

# Install repo tools

Bring a fresh clone (or a second PC) to a working khhub checkout. Load this skill instead of improvising install steps.

## Rules

- Never commit `.env` or paste secrets into chat or git.
- Never put `HCLOUD_TOKEN`, `DOKPLOY_API_KEY`, Cloudflare tokens, or GHCR passwords in the repo `.env`. That file is the **app** env only.
- Never create billed Hetzner resources from this skill. Listing servers is read-only; creating or deleting is `hetzner-deploy`.
- Never copy a production `.env` or congregation PII onto the new machine.
- On Windows, use **Git Bash** (not PowerShell) for the check script and `make`.
- Install only what the check marks missing. Re-run the script after each group.

## Workflow

1. Confirm the working directory is the khhub repo root. If the clone is missing: `git clone https://github.com/WijanDev/khhub-go.git` and `git checkout dev`.
2. Run the check (agent executes this):

   ```bash
   bash scripts/check-dev-tools.sh
   ```

3. Install each **missing** required tool. Commands: [install.md](install.md).
4. Prompt the human for each **cloud login** that is `needs-auth`. Do not invent tokens.
5. Local app files:
   - If `.env` is missing: `cp .env.example .env` and ask the human to set `ADMIN_PASSWORD` and `SESSION_SECRET` (32+ random characters). Leave `APP_ENV=development`.
   - `cd frontend && npm ci` (or `npm install` if `npm ci` fails).
   - `make dev-db` then `cd backend && go test ./...`.
6. Run `bash scripts/check-dev-tools.sh` again.
7. Report the table from the script. Say what is still missing and the next human action (browser login, Docker Desktop running, PATH).
8. Mention Cursor MCP: GitHub, Cloudflare, and Snyk stay in Cursor settings. This skill does not write MCP tokens.

## Tool groups

| Group | Tools | Needed to |
| --- | --- | --- |
| Required | Git, Go 1.24+, Node 22+, npm, Docker + Compose, `sqlc` | Code, tests, `sqlc generate`, local Postgres |
| Helpful | GNU `make`, `air` | Makefile shortcuts, API reload |
| Cloud | `gh`, `hcloud`, `dokploy` (`@dokploy/cli`), `wrangler` | PRs, VPS, Dokploy panel, R2 |
| Optional | `gitleaks` | Same secret scan as CI |

`make` is helpful, not required. Without it, use the commands in `README.md`.

## Out of scope

Installing Dokploy on the VPS, changing DNS, creating GHCR packages, cloning production volumes, or enabling Hetzner server backups.

## Additional resources

- Per-OS install and auth: [install.md](install.md)
- Local commands: [README.md](../../../README.md)
- Deploy hosts: [docs/deploy-dokploy.md](../../../docs/deploy-dokploy.md)

If the human asked to set up the machine, stop after the second check report. Do not open a PR unless they ask.
