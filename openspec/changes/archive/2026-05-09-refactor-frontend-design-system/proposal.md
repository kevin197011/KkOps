# Change: Introduce a unified frontend design system shell

## Why

Pages were built ad-hoc and now diverge in spacing, page headers, empty
states and density. v2 brings many new modules (integrations, alerts,
incidents, K8s, AI ops). Without a shared shell and primitives, every new
page reintroduces inconsistency. This change establishes the base shell
primitives once and uses them everywhere going forward.

## What Changes

- Add `frontend/src/components/shell/` primitives:
  `PageContainer`, `PageHeader`, `EmptyState`, `Section`, and a global
  `CommandPalette` (CMD/CTRL + K) that mirrors the navigation menu and
  exposes "Go to" / "Search" actions.
- Wire the command palette into `MainLayout` so it works on every page.
- Centralize page paddings, max-width and section spacing through
  `PageContainer`.
- Keep current themes (`light`/`dark`) but expose density tokens (`compact`)
  used by tables and forms.

## Impact

- Affected specs: new `frontend-design-system` capability.
- Affected code: `frontend/src/components/shell/*`,
  `frontend/src/layouts/MainLayout.tsx`.
- No backend or API changes.
- Existing pages keep working; future pages MUST adopt the shell primitives.
