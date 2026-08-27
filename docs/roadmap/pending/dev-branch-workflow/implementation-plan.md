# Implementation plan

1. Create branch `dev` from current `main` and push it. In GitHub: default branch stays `main` (clone/docs) **or** switch default to `dev` so new PRs aim there — prefer **default = `dev`** so “New pull request” is correct. Protect `main`: no direct pushes; PRs from `dev` only. Protect `dev`: PRs required, CI must pass.
2. Update `.github/workflows/ci.yml`: `pull_request` branches `[dev, main]`; `push` tests on `dev` and `main`. Keep the **images + Dokploy hook** job on `push` to `main` only (`github.ref == refs/heads/main`).
3. Update `docs/plans/2026-08-27-ghcr-traefik-cicd.md` (and the design note if it still says PR-to-main only), `README.md`, `AGENTS.md` Commands/workflow: “open PRs against `dev`”. `CONTRIBUTING.md` when that idea exists, same line.
4. Update agent-facing text: do not merge to `main` from a feature branch. Release = PR `dev` → `main`.
5. Verify: a feature PR into `dev` runs `test` and does not push GHCR. Merging `dev` → `main` runs `test` then `images`. Confirm branch protection on GitHub (human click).
