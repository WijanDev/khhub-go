# Utility analysis

- **Score:** 3
- **Scale:** 1 = low … 5 = high

## Justification

The operator (not the secretary) benefits on every merge to `dev` and every release. Docs, roadmap, and test-only changes stop creating unused GHCR packages and stop bouncing both containers. A frontend-only release no longer restarts the API or Postgres. Skipping it keeps the current “always build both, always deploy the whole stack” loop: extra CI minutes, extra tags, and unnecessary staging/production restarts when the image did not change.
