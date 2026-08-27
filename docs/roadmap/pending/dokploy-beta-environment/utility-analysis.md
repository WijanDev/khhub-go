# Utility analysis

- **Score:** 4
- **Scale:** 1 = low … 5 = high

## Justification

Beta testers must not sit on the production congregation database. A separate environment lets 2–3 people try builds before they hit `khhub.app`. That is production safety, not a new ministry feature. Skipping it means testing on live data or not testing with others at all. One extra VPS would isolate more but is unnecessary for this size.
