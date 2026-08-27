# Cost analysis

- **Score:** 4
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.00

## Justification

Steps 1–5 are still a Dokploy + Cloudflare day (second compose, DNS, Let’s Encrypt). The rewrite adds a real CD change: images on `dev`, two deploy hooks, retag-not-rebuild on `main` (steps 7–8), and a decision on SPA API origin (step 6). Runtime config touches the Vite build and the static image; dual-build is less frontend work but more CI and a pair of web tags to promote. Either path is tests + docs + Actions + Dokploy — a week-class change, score 4. RAM on the cx23 still fits one extra stack; no second VPS.
