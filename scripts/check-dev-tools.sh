#!/usr/bin/env bash
# Report khhub toolchain status. Safe to run on Git Bash (Windows), macOS, and Linux.
# Does not print secret values. Exit 1 if a required coding tool is missing or too old.

set -u

required_fail=0

row() {
  printf '%-12s %-14s %-22s %s\n' "$1" "$2" "$3" "$4"
}

have() {
  command -v "$1" >/dev/null 2>&1
}

ver_or_dash() {
  if have "$1"; then
    shift
    "$@" 2>/dev/null | head -n 1 | tr -d '\r'
  else
    echo "-"
  fi
}

go_major_minor() {
  go version 2>/dev/null | sed -n 's/.*go\([0-9][0-9]*\)\.\([0-9][0-9]*\).*/\1 \2/p'
}

node_major() {
  node -v 2>/dev/null | sed -n 's/^v\([0-9][0-9]*\).*/\1/p'
}

npm_major_minor() {
  npm --version 2>/dev/null | sed -n 's/^\([0-9][0-9]*\)\.\([0-9][0-9]*\).*/\1 \2/p'
}

check_required() {
  local name="$1"
  local bin="$2"
  shift 2
  if ! have "$bin"; then
    row "missing" "$name" "-" "install; see install-repo/install.md"
    required_fail=1
    return
  fi
  local v
  v="$(ver_or_dash "$bin" "$@")"
  row "ok" "$name" "${v:0:22}" ""
}

echo "khhub tool check"
echo "repo: $(pwd)"
echo
row "STATUS" "TOOL" "VERSION" "NOTE"
row "------" "----" "-------" "----"

# --- required ---
check_required git git git --version

if ! have npm; then
  row "missing" "npm" "-" "need 11.11+ (Node 24 LTS)"
  required_fail=1
else
  read -r npmmaj npmmin <<EOF
$(npm_major_minor)
EOF
  npmv="$(npm --version 2>/dev/null | tr -d '\r')"
  if [ -z "${npmmaj:-}" ] || [ "$npmmaj" -lt 11 ] || { [ "$npmmaj" -eq 11 ] && [ "${npmmin:-0}" -lt 11 ]; }; then
    row "old" "npm" "$npmv" "need 11.11+ (Node 24 LTS)"
    required_fail=1
  else
    row "ok" "npm" "$npmv" ""
  fi
fi

if ! have go; then
  row "missing" "go" "-" "need 1.24+"
  required_fail=1
else
  read -r gmaj gmin <<EOF
$(go_major_minor)
EOF
  gv="$(go version 2>/dev/null | tr -d '\r')"
  if [ -z "${gmaj:-}" ] || [ "$gmaj" -lt 1 ] || { [ "$gmaj" -eq 1 ] && [ "${gmin:-0}" -lt 24 ]; }; then
    row "old" "go" "${gv:0:22}" "need 1.24+"
    required_fail=1
  else
    row "ok" "go" "${gv:0:22}" ""
  fi
fi

if ! have node; then
  row "missing" "node" "-" "need 24+"
  required_fail=1
else
  nmaj="$(node_major)"
  nv="$(node -v 2>/dev/null | tr -d '\r')"
  if [ -z "${nmaj:-}" ] || [ "$nmaj" -lt 24 ]; then
    row "old" "node" "$nv" "need 24+"
    required_fail=1
  else
    row "ok" "node" "$nv" ""
  fi
fi

if ! have docker; then
  row "missing" "docker" "-" "Docker Desktop or Engine"
  required_fail=1
else
  dv="$(docker --version 2>/dev/null | tr -d '\r')"
  if docker info >/dev/null 2>&1; then
    row "ok" "docker" "${dv:0:22}" "daemon up"
  else
    row "down" "docker" "${dv:0:22}" "start Docker Desktop"
    required_fail=1
  fi
fi

if docker compose version >/dev/null 2>&1; then
  cv="$(docker compose version 2>/dev/null | tr -d '\r')"
  row "ok" "compose" "${cv:0:22}" ""
else
  row "missing" "compose" "-" "Docker Compose v2 plugin"
  required_fail=1
fi

if have sqlc; then
  row "ok" "sqlc" "$(sqlc version 2>/dev/null | tr -d '\r' | head -n 1)" ""
else
  row "missing" "sqlc" "-" "go install .../sqlc@latest + GOPATH/bin"
  required_fail=1
fi

# --- helpful ---
if have make; then
  row "ok" "make" "$(ver_or_dash make make --version)" ""
else
  row "skip" "make" "-" "optional; use README commands"
fi

if have air; then
  row "ok" "air" "installed" "API reload"
else
  row "skip" "air" "-" "optional go install air"
fi

if have gitleaks; then
  row "ok" "gitleaks" "$(ver_or_dash gitleaks gitleaks version)" ""
else
  row "skip" "gitleaks" "-" "optional; CI runs it"
fi

# --- cloud: present? authenticated? ---
if have gh; then
  if gh auth status >/dev/null 2>&1; then
    row "ok" "gh" "$(gh --version 2>/dev/null | head -n 1 | tr -d '\r')" "logged in"
  else
    row "needs-auth" "gh" "$(gh --version 2>/dev/null | head -n 1 | tr -d '\r')" "gh auth login"
  fi
else
  row "missing" "gh" "-" "PRs and GHCR"
fi

if have hcloud; then
  if hcloud context active >/dev/null 2>&1 || [ -n "${HCLOUD_TOKEN:-}" ]; then
    row "ok" "hcloud" "$(hcloud version 2>/dev/null | tr -d '\r')" "context or HCLOUD_TOKEN"
  else
    row "needs-auth" "hcloud" "$(hcloud version 2>/dev/null | tr -d '\r')" "hcloud context create"
  fi
else
  row "missing" "hcloud" "-" "Hetzner VPS"
fi

if have dokploy; then
  if dokploy project all >/dev/null 2>&1 || { [ -n "${DOKPLOY_URL:-}" ] && [ -n "${DOKPLOY_API_KEY:-}" ]; }; then
    row "ok" "dokploy" "installed" "authenticated"
  else
    row "needs-auth" "dokploy" "installed" "dokploy auth -u https://admin.wijan.dev -t"
  fi
else
  row "missing" "dokploy" "-" "npm i -g @dokploy/cli"
fi

if have wrangler; then
  if wrangler whoami >/dev/null 2>&1; then
    row "ok" "wrangler" "$(wrangler --version 2>/dev/null | tr -d '\r')" "logged in"
  else
    row "needs-auth" "wrangler" "$(wrangler --version 2>/dev/null | tr -d '\r')" "wrangler login"
  fi
else
  row "missing" "wrangler" "-" "Cloudflare / R2 identity"
fi

# --- repo files (no secret values) ---
echo
if [ -f .env ]; then
  row "ok" ".env" "-" "present (not shown)"
else
  row "missing" ".env" "-" "cp .env.example .env"
fi

if [ -d frontend/node_modules ]; then
  row "ok" "npm-modules" "-" "frontend/node_modules"
else
  row "missing" "npm-modules" "-" "cd frontend && npm ci"
fi

if [ -d .git ]; then
  remote="$(git remote get-url origin 2>/dev/null || true)"
  branch="$(git branch --show-current 2>/dev/null || true)"
  row "ok" "git-repo" "$branch" "${remote:-no origin}"
else
  row "missing" "git-repo" "-" "clone WijanDev/khhub-go"
fi

echo
if [ "$required_fail" -ne 0 ]; then
  echo "Required coding tools are incomplete."
  exit 1
fi
echo "Required coding tools look ready. Fix any needs-auth rows before cloud work."
exit 0
