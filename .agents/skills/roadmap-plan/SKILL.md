---
name: roadmap-plan
description: Capture and rank khhub ideas under docs/roadmap. Use after brainstorming approves a design, when adding or scoring an idea, when marking an idea implemented, or when the user mentions the roadmap.
---

# Roadmap plan

Living work lives only in [docs/roadmap/](../../../docs/roadmap/). There is no root `ROADMAP.md`. Do not invent a second backlog.

`brainstorming` must load this skill after the human approves a design (bounded or architectural). A spike does **not** create an idea.

English only. Spanish stays in the UI.

## Layout

```
docs/roadmap/
  summary.md
  pending/<slug>/
  implementing/<slug>/
  implemented/<slug>/
  dropped/<slug>/
```

Each idea folder still has `idea.md`, `implementation-plan.md`, `cost-analysis.md`, `utility-analysis.md`.

`summary.md` is the index. Every pending, implementing, or implemented row must have a matching folder. Never list an idea that has no folder.

Create new work under `pending/<slug>/`. **Move** the folder when status changes:

| Status | Folder |
| --- | --- |
| `proposed` or `planned` | `pending/<slug>/` |
| `in-progress` | `implementing/<slug>/` |
| `implemented` | `implemented/<slug>/` |
| `dropped` | `dropped/<slug>/` |

## Scores (1–5)

| Score | Cost (1 = cheap) | Utility (1 = low) |
| --- | --- | --- |
| 1 | Hours, few files, no schema | Nice-to-have |
| 2 | About a day, localized change | Helps one occasional flow |
| 3 | Several days, more than one layer | Helps a weekly congregation task |
| 4 | A week-class change, tests + docs | Unblocks a core weekly job or production safety |
| 5 | Multi-week or new subsystem | Required for safe daily use or a hard constraint |

**Ratio** = utility ÷ cost (higher is better). Example: utility 5 / cost 2 = `2.50`.

Write `implementation-plan.md` **before** `cost-analysis.md`. Cost must cite the plan steps. Utility does not depend on the plan.

## Create an idea (from brainstorming)

After the human approves the designed idea:

1. Choose a kebab-case slug (`pwa-installable-shell`, not `PWA`).
2. Create `docs/roadmap/pending/<slug>/` with all four files. Use [references/templates.md](references/templates.md).
3. Status in `idea.md` starts as `proposed`.
4. Fill `implementation-plan.md` from the approved design (steps an agent can execute).
5. Fill `cost-analysis.md` from that plan (score + justification).
6. Fill `utility-analysis.md` (score + justification). Who benefits, how often, what breaks if we skip it.
7. Rebuild `summary.md` (see below).
8. Stop. Do not implement the product work until the human asks. Capturing the idea is not approval to code the feature.

If a folder already exists for the same outcome, update it. Do not duplicate slugs.

## Mark in-progress

When the human asks to build the idea:

1. Set `idea.md` status to `in-progress`.
2. Move `docs/roadmap/pending/<slug>/` to `docs/roadmap/implementing/<slug>/`.
3. Rebuild `summary.md`.

## Mark implemented

Only after the work is merged to `main` (release PR from `dev`) and deployed:

1. Set `idea.md` status to `implemented`.
2. Record **merge date** (`YYYY-MM-DD` of the merge commit on `main`).
3. Record **app version** from the repo root `VERSION` file at that deploy (create `VERSION` on the first tagged release if it is missing; do not invent a version).
4. Move `docs/roadmap/implementing/<slug>/` (or `pending/<slug>/` if it never moved) to `docs/roadmap/implemented/<slug>/`.
5. Rebuild `summary.md`.

## Rebuild `summary.md`

Three lists. No fourth “unscored” section.

**Pending** — `proposed` or `planned`. Sort by ratio descending, then utility descending, then slug. Columns: idea (link to `pending/<slug>/idea.md`), cost, utility, ratio.

**Implementing** — `in-progress`. Same sort and columns; links use `implementing/`.

**Implemented** — `implemented`. Sort by **merge date, newest first**. Columns: idea (link to `implemented/<slug>/idea.md`), cost, utility, merge date, version.

Keep the non-official product disclaimer at the top. Repeat out-of-scope constraints in the intro (no jw.org, no S-21 replacement, no public signup, no real PII in git). Do not turn those into idea folders.

## Status values (`idea.md`)

`proposed` · `planned` · `in-progress` · `implemented` · `dropped`

`dropped` stays on disk under `dropped/` and is omitted from the `summary.md` tables.

## Do not

- Score cost without an implementation plan.
- Cache or invent merge dates or versions.
- Keep a parallel checklist in `README.md` or chat.
- Start coding from brainstorming until this folder exists and the human says to implement.
