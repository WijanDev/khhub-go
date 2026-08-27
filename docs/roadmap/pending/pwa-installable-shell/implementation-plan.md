# Implementation plan

1. Add `vite-plugin-pwa` in `frontend/vite.config.ts` with `manifest` name `khhub`, `start_url` `/`, `display` `standalone`, and a theme color from the existing CSS tokens.
2. Generate 192 and 512 PNG icons (plus `apple-touch-icon`) from `frontend/public/favicon.svg`. Put them in `frontend/public/`.
3. Configure the service worker to precache built JS/CSS/HTML only. Set runtime caching for `https://api.khhub.app` to network-only (no Cache API for JSON or cookies).
4. Add Apple meta tags and the web manifest link in `frontend/index.html`.
5. Add a `ROADMAP` line is not required; this folder is the source. Mention install in `docs/deploy-dokploy.md` in one sentence if the production host is listed there.
6. Verify: `cd frontend && npm run lint && npm run build` emits `manifest.webmanifest` and a service worker. Confirm `GET https://khhub.app/publishers` still returns the SPA HTML (200). On a phone: Android Chrome → Install; iOS Safari → Share → Add to Home Screen. Do not cache `/auth/me` or report payloads.
