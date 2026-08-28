# Pin Node 24 LTS and npm 11.11+

- **Slug:** node-24-lts
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Operators and CI build the Vite SPA on **Node 24 LTS** (Active LTS) with the npm that ships in a recent 24.x (**≥ 11.11.0**, which keeps `libc` in the lockfile and does not stamp false `"peer": true`). The congregation secretary sees no UI change. Local machines still on Node 24.11 / npm 11.6 must update Node so `npm install` stops rewriting `package-lock.json`.

## Out of scope

- Node 26 (Current), or a separate global `npm i -g npm@…` as the supported path.
- Rewriting historical plans (for example `docs/plans/2026-08-27-ghcr-traefik-cicd.md`).
- Changing frontend dependencies or Vite behavior beyond lockfile metadata from the newer npm.
- Runtime Node in production: the published image stays static files; only the build stage changes.
