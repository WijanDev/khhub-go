---
name: khhub-docs
description: Write and update khhub documentation in English while keeping the UI in Spanish. Use when editing README, ROADMAP, AGENTS.md, deploy docs, plans, or when the user asks to document a feature.
---

# khhub documentation

Load `developer-docs-planning` / `developer-docs-drafting` when scoping new pages.

## Language

- **All markdown and comments:** English.
- **UI:** Spanish. Do not translate `frontend/src/lib/labels.ts` or feature copy into English unless the user asks.

## Where things live

| Audience | File |
| --- | --- |
| Humans, quick start | `README.md` |
| Agents, always-on rules | `AGENTS.md` |
| Living checklist | `ROADMAP.md` |
| Production deploy | `docs/deploy-dokploy.md` |
| Feature implementation plans | `docs/plans/YYYY-MM-DD-<feature>.md` |
| Skill catalog | `.agents/README.md` |

Do not add Cursor-only paths (`.cursor/rules`, `.cursor/skills`) for project guidance. Repo skills stay under `.agents/skills/`.

## Style

- State what is true and how to run it. No marketing.
- Repeat the non-official disclaimer in human-facing docs (`README.md`, deploy doc).
- Never paste real `.env` values. Point at `.env.example`.
- When behavior changes, update the matching doc in the same change. Then tick or add a line in `ROADMAP.md` (`maintain-roadmap`).

## Feature docs

A new user-facing capability needs:

1. A `ROADMAP.md` item (done or planned).
2. Enough README or deploy detail for a human to use or ship it.
3. An implementation plan in `docs/plans/` only when the work is multi-step.

Do not invent a docs site or Docusaurus unless the roadmap asks for it.
