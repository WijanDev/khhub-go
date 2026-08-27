# Cost analysis

- **Score:** 3
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.33

## Justification

Duplicating a compose in Dokploy (steps 1–3, 5) is a day of panel work, but the plan is not only that. The SPA bakes `VITE_API_URL` (step 6), so beta needs a second image tag and a CI path that must not overwrite production `:latest`. DNS and a second Let’s Encrypt pair (steps 4–5) add another layer. RAM on the cx23 is enough for one extra stack (~100 MiB), so cost is operational complexity, not a new server. That is several days across Dokploy, Cloudflare, and GitHub Actions — score 3, not 2.
