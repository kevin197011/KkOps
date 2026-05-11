# integrations-hub Specification

## Purpose
TBD - created by archiving change add-integrations-hub. Update Purpose after archive.
## Requirements
### Requirement: Provider registry

The system SHALL expose a registry of integration providers keyed by normalized kind string. Each provider SHALL implement `Kind`, `Verify(context)`, and `Metadata` for discovery.

#### Scenario: Verify integration

- **WHEN** an authorized client calls `POST /api/v1/integrations/:id/test` for an existing integration
- **THEN** the server SHALL decrypt credentials for that integration, resolve the provider by kind, call `Verify`, and SHALL NOT panic on upstream failure (SHALL return an HTTP error with a message instead).

### Requirement: Integrations hub UI

The web application SHALL provide a page at path `/integrations` listing configured integrations and supported kinds, with actions to add, edit, delete, and test connection where RBAC allows.

#### Scenario: Unauthorized hub access

- **WHEN** a user without `integrations:*` opens `/integrations`
- **THEN** the client SHALL block navigation per existing protected-route rules.

