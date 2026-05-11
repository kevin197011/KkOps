# Change: Populate openspec/project.md as the v2 baseline

## Why

`openspec/project.md` is still the empty scaffold. Every Phase of the v2
roadmap depends on shared context (tech stack, layered architecture, RBAC
contract, integrations framework) being explicit so subsequent change
proposals can reference it instead of re-stating it.

## What Changes

- Fill in `openspec/project.md` with: purpose, tech stack, code style,
  architecture patterns, testing strategy, git workflow, domain context,
  constraints, and external dependencies.
- Add an "OpenSpec Process" capability that codifies the rules already
  applied in this repo (proposal-first, validate strict, archive after ship).

## Impact

- Affected specs: new `openspec-process` capability.
- Affected code: `openspec/project.md`.
- No runtime/code impact.
