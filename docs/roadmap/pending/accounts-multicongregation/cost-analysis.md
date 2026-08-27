# Cost analysis

- **Score:** 5
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.00

## Justification

Steps 1–2 are a schema rewrite and every list/query grows a tenant filter. Steps 3–6 add a new auth surface (passkeys, enrollment tokens, `/auth/me` claims) and a second SPA shell. Step 7 is platform admin. Tests in step 8 must cover isolation and the permission matrix. That is a multi-week subsystem (score 5), not a week-class feature. Risk is leaking another congregation’s rows or caching the wrong active congregation on the session.
