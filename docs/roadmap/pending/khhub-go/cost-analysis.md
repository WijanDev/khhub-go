# Cost analysis

- **Score:** 2
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.50

## Justification

Steps 1–3 are markdown in `.agents/` and `AGENTS.md` (no schema, no UI). Steps 4–5 are a localized Go move plus two error mappings and one helper call, with service-package tests. About a day, one layer (API + agent docs). Risk is low: no route or cookie change; the only user-visible shift is Spanish copy that already should have been Spanish.
