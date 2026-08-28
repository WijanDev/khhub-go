# Cost analysis

- **Score:** 4
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.00

## Justification

Steps 1–2 are two new tables (versioned timetable + dated exceptions). Step 3 changes `kind` and the unique key. Steps 4–6 are handler rules, settings UI, attendance form, and dashboard averages that must skip cancelled dates and keep Memorial out of the weekly mean. That is schema + API + two UI surfaces + tests — a week-class change, score 4. Still not a meeting-program product (assignments, talks). Depends on `congregation_id` or one implicit congregation.
