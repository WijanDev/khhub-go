# Cost analysis

- **Score:** 2
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.50

## Justification

The plan is frontend-only: one Vite plugin, icon export from the existing favicon, manifest fields, and a service-worker cache policy that excludes the API (steps 1–4). No schema, no Gin routes, no cookie changes. A day-class change with a lint/build check and a phone install smoke test (steps 6). Risk is misconfigured caching of congregation JSON; the plan forbids that, which keeps the work localized rather than a new subsystem.
