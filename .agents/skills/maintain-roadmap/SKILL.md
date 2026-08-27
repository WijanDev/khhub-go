---
name: maintain-roadmap
description: Keep ROADMAP.md current. Use when shipping a feature, finishing a stabilization item, planning new work, or when the user mentions the roadmap, checklist, or next priorities.
---

# Maintain the roadmap

`ROADMAP.md` is the living checklist. Agents do not invent a second backlog.

## When work finishes

1. Open `ROADMAP.md`.
2. Mark the matching item `[x]` and keep the original wording if it still describes the outcome.
3. If the work was not listed, add it under the right section as `[x]` with one line of outcome.
4. If the change creates follow-up work, add unchecked items in the same section. Do not hide debt.
5. Do not delete completed items. The checked history is the changelog of intent.

## When planning

1. Load `brainstorming` and/or `docs-ai-prd` for the design.
2. Add unchecked items to `ROADMAP.md` before coding, unless the user asked only for a plan file.
3. Split functional product work from engineering stabilization. Do not mix them in one bullet.
4. Multi-step builds also get `docs/plans/YYYY-MM-DD-<feature>.md` via `writing-plans`.

## Item style

```markdown
- [ ] Short outcome in English, not a file list
```

Good: `- [ ] Add GitHub Actions for go test and frontend build`
Bad: `- [ ] Fix stuff` / `- [ ] Edit main.go`

## Priority

Unless the user reorders: finish **Development stabilization** before large later features (territories, LMM, cart, literature, shepherding).
