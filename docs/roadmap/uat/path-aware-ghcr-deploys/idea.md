# Path-aware GHCR publishes and split Dokploy apps

- **Slug:** path-aware-ghcr-deploys
- **Status:** uat
- **Merge date:**
- **App version:**

## Summary

Publish a GHCR image and fire a Dokploy hook only when that image’s contents changed. Docs, roadmap, tests, and other non-image paths do not create packages or deploys. A frontend-only change never builds or deploys the API, and the reverse.

Dokploy project `khhub` keeps the current `staging` and `production` environments. Replace the single Compose stack with two apps, `khhub-api` and `khhub-web`. One Dokploy Postgres per environment; the API points at it with `DATABASE_URL`. Four independent deploy hooks (API/web × staging/production). A release to `main` retags and deploys only the side that changed.

## Out of scope

- Preview deploys for every pull request
- Rebuilding images on `main` (promote stays retag-from-`:staging`)
- Changing hosts (`khhub.app`, `api.khhub.app`, `staging.khhub.app`, `apistaging.khhub.app`)
- A second VPS or a second Dokploy install
- Splitting Postgres across more than one database per environment
- Baking congregation secrets into images
