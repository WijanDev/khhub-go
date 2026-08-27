# Utility analysis

- **Score:** 4
- **Scale:** 1 = low … 5 = high

## Justification

Staging testers must not sit on the production congregation database. A separate environment lets 2–3 people try builds on `staging.khhub.app` before they hit `khhub.app`. Promoting the **same** image digest (instead of rebuilding on `main`) means production runs what staging already exercised. That is production safety, not a new ministry feature. Skipping it means testing on live data, not testing with others, or shipping a newly built `:latest` that staging never ran. One extra VPS would isolate more but is unnecessary for this size.
