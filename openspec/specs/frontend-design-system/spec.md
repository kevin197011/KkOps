# frontend-design-system Specification

## Purpose
TBD - created by archiving change refactor-frontend-design-system. Update Purpose after archive.
## Requirements
### Requirement: Page Shell Primitives

The frontend MUST expose a small library of shell primitives (`PageContainer`,
`PageHeader`, `EmptyState`, `Section`) that all top-level pages use to
guarantee consistent padding, page header layout and empty-state styling.

#### Scenario: New top-level page

- **WHEN** a developer creates a new top-level page under
  `frontend/src/pages/`
- **THEN** the page MUST render content inside `PageContainer`
- **AND** SHOULD render a `PageHeader` providing the title, description and
  primary action

#### Scenario: Empty list

- **WHEN** a list page has no rows and is not loading
- **THEN** an `EmptyState` MUST be rendered with a description and an
  optional call-to-action

### Requirement: Global Command Palette

The application MUST expose a global command palette opened with
`Ctrl+K` / `Cmd+K` listing the navigation entries the current user has
permission to access.

#### Scenario: Authenticated user presses Cmd+K

- **WHEN** the user is authenticated and presses `Cmd+K` (or `Ctrl+K` on
  non-mac)
- **THEN** the command palette MUST open
- **AND** the listed entries MUST be filtered by the same permission rules
  that filter the side menu

#### Scenario: Selecting an entry

- **WHEN** the user selects an entry in the palette
- **THEN** the application MUST navigate to that route and close the palette

