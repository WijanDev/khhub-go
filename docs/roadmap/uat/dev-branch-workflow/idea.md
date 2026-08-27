# PRs target `dev`, not `main`

- **Slug:** dev-branch-workflow
- **Status:** uat
- **Merge date:**
- **App version:**

## Summary

`main` is production (images + Dokploy). Feature work opens PRs against **`dev`**. Promote `dev` → `main` with a release PR when we intend to ship. Agents and humans follow the same default.

## Out of scope

- Git Flow with `release/*` and `hotfix/*` unless we later need it.
- Deploying production from `dev`.
- Forcing a second VPS (staging stays `dokploy-staging-environment`).
