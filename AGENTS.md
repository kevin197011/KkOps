<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

When implementing code (backend, frontend, API, database), follow the development rules in `rules/` (alwaysApply Cursor rules). See `docs/README.md` for the rules index.

All development changes that add features, introduce breaking changes, or change architecture/behavior MUST follow the OpenSpec workflow: create a change proposal first, get approval, then implement; see `rules/07-openspec-workflow-rules.mdc` and `openspec/AGENTS.md`.

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->