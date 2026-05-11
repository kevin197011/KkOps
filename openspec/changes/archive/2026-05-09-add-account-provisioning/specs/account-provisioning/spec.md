## ADDED Requirements

### Requirement: Provisioning targets

The system SHALL persist provisioning targets linking each integration record to a provider kind (for example `scim`, `gitlab`, `jenkins`) with sync status and last successful sync timestamp.

#### Scenario: List targets

- **WHEN** an authorized client calls the provisioning targets list API
- **THEN** the server SHALL return targets including integration reference, provider kind, status, and `last_sync_at` when present.

### Requirement: Provisioning run history

The system SHALL record each provisioning attempt in a run history with outcome status and timestamps for auditing and troubleshooting.

#### Scenario: Runs visible after sync

- **WHEN** a sync completes for a target
- **THEN** a provisioning run row SHALL exist associating the target with the outcome.

### Requirement: User lifecycle hooks

The system SHALL enqueue asynchronous provisioning work when users are created, updated, or deleted so external systems can be notified without blocking the HTTP request.

#### Scenario: Create user enqueues work

- **WHEN** a user is successfully created through the user API
- **THEN** the system SHALL enqueue provisioning sync work for enabled targets.

### Requirement: Manual sync API

The system SHALL expose an endpoint to trigger synchronization for a specific provisioning target on demand.

#### Scenario: Authorized sync

- **WHEN** a caller with `provisioning:*` invokes manual sync for an existing target
- **THEN** the server SHALL accept the request and initiate a sync operation for that target.

### Requirement: Authorization

The system SHALL require `provisioning:*` permission for provisioning management HTTP endpoints under `/api/v1/provisioning`.

#### Scenario: Forbidden without permission

- **WHEN** a user without `provisioning:*` calls a mutating or listing provisioning endpoint
- **THEN** the server SHALL respond with HTTP 403.
