# Cost analysis

- **Score:** 3
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.00

## Justification

Steps 1–5 are a localized CI change (path filters, conditional build/retag, four secrets). Steps 6–8 are the real cost: two Dokploy apps, four webhooks, and moving Postgres off the Compose volume without mixing staging and production data. Step 9 is docs. Step 10 needs live GHCR and Dokploy, not unit tests. Several days and more than one layer (CI + Dokploy + deploy docs). No application schema. Risk sits in the database move, not in the path-filter YAML.
