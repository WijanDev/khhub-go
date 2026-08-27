# Agent harness

This folder is the **vendor-neutral** Agent Skills location. The [Agent Skills specification](https://agentskills.io/specification) does not mandate a path; [`.agents/skills/`](https://agentskills.io/client-implementation/adding-skills-support) is the cross-client convention used by Codex, Amp, Cursor (project scope), Gemini, and others.

Always-on project rules live in [`../AGENTS.md`](../AGENTS.md), not in skills.

## Layout

```
.agents/
  README.md          # this file
  skills/
    <name>/SKILL.md  # one directory per skill
```

Install or refresh third-party skills with the [skills CLI](https://skills.sh) into Universal (writes here, no Cursor-only copy):

```bash
npx skills add <owner/repo> -a universal --copy -y -s <skill-name>
```

`../skills-lock.json` pins source + hash. Project-owned skills (`khhub-stack`, `khhub-docs`, `maintain-roadmap`) are not in that lockfile.

## Installed skills

### Planning and documentation

| Skill | Source | Use |
| --- | --- | --- |
| `brainstorming` | [obra/superpowers](https://github.com/obra/superpowers) | Design before coding |
| `writing-plans` | obra/superpowers | Bite-sized implementation plans |
| `executing-plans` | obra/superpowers | Execute a written plan with checkpoints |
| `docs-ai-prd` | [vasilyu1983/ai-agents-public](https://github.com/vasilyu1983/ai-agents-public) | PRDs, specs, acceptance criteria |
| `developer-docs-planning` | [lvtd-llc/skills](https://github.com/lvtd-llc/skills) | What docs to write before drafting |
| `developer-docs-drafting` | lvtd-llc/skills | Draft developer-facing docs |
| `maintain-roadmap` | this repo | Keep `ROADMAP.md` honest |
| `khhub-docs` | this repo | English docs + Spanish UI |

### Stack

| Skill | Source | Use |
| --- | --- | --- |
| `khhub-stack` | this repo | sqlc, cookie auth, Base UI — **overrides** generic Gin/React advice |
| `golang-gin-api` | [henriqueatila/golang-gin-best-practices](https://github.com/henriqueatila/golang-gin-best-practices) | Gin routing and handlers |
| `golang-gin-testing` | same | httptest / table-driven tests |
| `golang-gin-psql-dba` | same | schema, indexes, migration safety |
| `tanstack-query` | [tanstack-skills/tanstack-skills](https://github.com/tanstack-skills/tanstack-skills) | Query v5 overview |
| `tanstack-query-best-practices` | [deckardger/tanstack-agent-skills](https://github.com/deckardger/tanstack-agent-skills) | Query caching and mutations |
| `tanstack-router` | tanstack-skills/tanstack-skills | Router overview |
| `tanstack-router-best-practices` | deckardger/tanstack-agent-skills | Search params, loaders, navigation |
| `tanstack-table` | tanstack-skills/tanstack-skills | Table **v8** (`useReactTable`) |
| `hetzner-deploy` | [fcakyon/claude-codex-settings](https://github.com/fcakyon/claude-codex-settings) | `hcloud` CLI: servers, networks, and live type/pricing lookup |
| `dokploy-deploy` | same pack | Dokploy Compose projects, domains, and databases |

### Quality

| Skill | Source | Use |
| --- | --- | --- |
| `systematic-debugging` | obra/superpowers | Bugs and unexpected behavior |
| `verification-before-completion` | obra/superpowers | Evidence before “done” |
| `requesting-code-review` | obra/superpowers | Review after sizable work |

## Deliberately not installed

- Full `obra/superpowers` pack (`using-superpowers` forces a skill on every reply).
- Gin JWT / GORM skills — they fight this repo’s cookie + sqlc design.
- Vercel React / composition / writing skills — this SPA uses TanStack, not Next.js.
- TanStack Start skills — khhub is a Vite SPA; the API is Go.
- Official Table v9 Intent skills — not installed; the app is `@tanstack/react-table` v8.
- Hetzner MCP servers — the skill uses the official `hcloud` CLI instead. Same API token as the CLI (`HCLOUD_TOKEN` or `hcloud context`). Do not put the token in git.

Official Router skills also live in `frontend/node_modules/@tanstack/router-core/skills/` (TanStack Intent). Load with `npx @tanstack/intent@latest load @tanstack/router-core#router-core` from `frontend/`. Query and Table packages do not ship Intent skills at the versions pinned here.
