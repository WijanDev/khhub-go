# Cost analysis

- **Score:** 1
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 2.00

## Justification

Steps 2–4 and 6 are pin and copy edits (CI, Dockerfile, `engines`, README, install-repo, check script). Step 5 is one lockfile regenerate plus lint/build. No API, schema, or UI. Risk is lockfile noise if someone regenerates with npm 11.6.x (step 1 and the check in step 7). Hours, few files, CI-only toolchain.
