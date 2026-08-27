# Idea file templates

Replace `SLUG`, scores, and copy. Keep English.

## idea.md

```markdown
# Title

- **Slug:** SLUG
- **Status:** proposed
  (`proposed` · `planned` · `in-progress` · `uat` · `implemented` · `dropped`)
- **Merge date:** (empty until implemented)
- **App version:** (empty until implemented)

## Summary

One short paragraph: what changes for the congregation secretary or the operator.

## Out of scope

Bullets for what this idea will not do.
```

## implementation-plan.md

```markdown
# Implementation plan

Numbered steps an agent can execute. Name the files and checks.
Cost analysis is written from this list.

1. …
2. …
3. Verify: …
```

## cost-analysis.md

```markdown
# Cost analysis

- **Score:** N
- **Scale:** 1 = cheap … 5 = expensive

## Justification

Tie the score to the steps in `implementation-plan.md`. Mention layers (API, UI, schema, CI) and risk.
```

## utility-analysis.md

```markdown
# Utility analysis

- **Score:** N
- **Scale:** 1 = low … 5 = high

## Justification

Who benefits, how often, and what we lose if we skip it.
```
