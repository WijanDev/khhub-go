# Cost analysis

- **Score:** 2
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 2.00

## Justification

Steps 1–2 are a branch, protection, and a few lines in `ci.yml`. Steps 3–4 are docs. About a day plus a human click on GitHub settings. No schema. Risk is building `:latest` from `dev` by mistake; the plan keeps the image job on `main` only.
