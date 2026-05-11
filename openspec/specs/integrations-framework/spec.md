# integrations-framework Specification

## Purpose
TBD - created by archiving change add-backend-integrations-framework. Update Purpose after archive.
## Requirements
### Requirement: Integration records

The system SHALL persist integration definitions in a relational store with at least: unique identifier, human-readable name, integration kind (string), enabled flag, encrypted credential/configuration payload, and timestamps.

#### Scenario: Create integration

- **WHEN** an authorized client submits a valid create request with name, kind, optional plaintext JSON credentials, and enabled flag
- **THEN** the server SHALL store credentials only in encrypted form and SHALL NOT persist plaintext secrets.

### Requirement: Encrypted credentials at rest

The system SHALL encrypt integration credential/configuration payloads using the application encryption key before persistence.

#### Scenario: Read integration via API

- **WHEN** a client reads an integration via the public HTTP API
- **THEN** the response SHALL NOT include decrypted credential fields.

### Requirement: Authorization

The system SHALL require `integrations:*` permission for all `/api/v1/integrations` HTTP endpoints used for mutating or listing integration records.

#### Scenario: Unauthorized access

- **WHEN** a non-admin user without `integrations:*` calls an integrations endpoint that requires that permission
- **THEN** the server SHALL respond with HTTP 403.

