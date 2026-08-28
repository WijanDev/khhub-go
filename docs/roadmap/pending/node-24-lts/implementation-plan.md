# Implementation plan

1. Discard any dirty `frontend/package-lock.json` produced by npm 11.6.x (`git checkout -- frontend/package-lock.json` if it is only lockfile-format churn).
2. In `.github/workflows/ci.yml`, set `actions/setup-node` `node-version` to `"24"` (latest 24.x; 24.20+ bundles npm 11.19).
3. In `frontend/Dockerfile`, change the build stage from `node:22-alpine` to `node:24-alpine`.
4. In `frontend/package.json`, add `engines`: `node` `>=24.14.0` and `npm` `>=11.11.0` (warning only; do not add `engine-strict`).
5. On Node **≥ 24.14** (prefer current 24 LTS), from `frontend/`: run `npm install` once so the lockfile is rewritten by npm ≥ 11.11, then confirm `npm ci && npm run lint && npm run build`. Commit lockfile changes only if they are the stable 11.11+ format (keep `libc` on optional natives; no bogus `"peer": true` on `react` / `typescript` / similar).
6. Update living docs and the tool check: `README.md` (Node 24+), `.agents/skills/install-repo/SKILL.md` and `install.md` (Node 24 LTS; `winget` LTS / `brew install node`, not `node@22`). In `scripts/check-dev-tools.sh`, fail required if Node major is below 24 or npm is below 11.11.
7. Verify: `bash scripts/check-dev-tools.sh` is green on a 24.14+ machine. CI `test` job uses Node 24. A second `npm install` in `frontend/` does not rewrite the lockfile. Do not edit historical plan markdown.
