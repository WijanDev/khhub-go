# Implementation plan

Two slices. Do **only slice 1** until the human asks for planned versions.

## Slice 1 — counts and percentages

1. In `docs/roadmap/summary.md`, add a line under the intro: **Active total** = pending + implementing + UAT + implemented. Dropped folders stay listed in the intro and stay out of the total.
2. Change each section heading to include the count and the share of that total, one decimal if needed so they add to 100% (or whole percents that sum to 100 by rounding the largest remainder). Example: `## Pending (39 · 81%)`.
3. Update `.agents/skills/roadmap-plan/SKILL.md` **Rebuild `summary.md`**: every rebuild must recount folders (not the previous heading numbers) and rewrite the heading suffixes. Same rule in `AGENTS.md` only if it already describes the summary tables.
4. Verify: open `summary.md`; each heading count matches the table rows (or `—` when empty); the four shares refer to the same active total; a dropped slug is not in the total.

## Slice 2 — planned + target version (later)

5. Give `planned` its own folder `docs/roadmap/planned/<slug>/` (move out of `pending/`). Status `planned` no longer shares the pending table.
6. Add **Target version** on `idea.md` (empty until planned). Rebuild `summary.md` with a **Planned** table between pending and implementing: idea, cost, utility, ratio, target version. Sort by target version, then ratio, then slug.
7. Update the skill layout, status table, templates, and `khhub-docs`. Do not invent a version; copy `VERSION` or the next bump the human names.
8. Verify: a planned idea has a folder under `planned/`, appears only in that table, and pending/implementing/UAT/implemented counts still exclude dropped.
