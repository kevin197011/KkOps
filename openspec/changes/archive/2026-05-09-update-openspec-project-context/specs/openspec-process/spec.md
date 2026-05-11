## ADDED Requirements

### Requirement: Proposal-First Development

The repository MUST require an OpenSpec change proposal under
`openspec/changes/<change-id>/` for all substantive code changes (new
features, breaking changes, architecture shifts, security or performance
work) before implementation begins.

#### Scenario: Feature work without a proposal

- **WHEN** a contributor begins implementing a new feature
- **THEN** a corresponding `openspec/changes/<id>/proposal.md` MUST exist
- **AND** at least one spec delta file under `specs/<capability>/spec.md`
  MUST be present
- **AND** `openspec validate <id> --strict` MUST pass before merge

#### Scenario: Bug fix that restores spec behavior

- **WHEN** a contributor fixes a regression
- **THEN** the proposal step MAY be skipped, provided the fix only restores
  the existing spec behavior

### Requirement: Project Context Document

The repository MUST maintain `openspec/project.md` populated with the project
purpose, tech stack, conventions, domain model, constraints and external
dependencies. New change proposals SHOULD reference `project.md` rather than
restating shared context.

#### Scenario: New AI assistant joins

- **WHEN** an AI assistant or new contributor opens the repo
- **THEN** reading `openspec/project.md` is sufficient to understand the
  layered backend, the React + AntD frontend, the RBAC model and the
  integrations framework expectations
