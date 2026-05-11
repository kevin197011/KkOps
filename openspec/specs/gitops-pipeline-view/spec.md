# gitops-pipeline-view Specification

## Purpose
TBD - created by archiving change add-gitops-pipeline-view. Update Purpose after archive.
## Requirements
### Requirement: Pipeline view API

The system SHALL provide `GET /api/v1/gitops/pipeline-view` returning a time-ordered list of normalized events combining recent CI/CD activity and Argo CD deployment history for a requested application name.

#### Scenario: Valid application

- **WHEN** a client passes query `app` matching an Argo CD application name discoverable via configured integrations
- **THEN** the response includes `PipelineEvent` entries with timestamp, kind, source, ref, status, and link fields

#### Scenario: Missing application name

- **WHEN** `app` is empty
- **THEN** the API responds with a client error describing the missing parameter

### Requirement: Authorization

The pipeline view endpoint SHALL use the same permission mapping as other `/api/v1/gitops` routes (`gitops:*`).

#### Scenario: Unauthorized

- **WHEN** the user lacks `gitops:*`
- **THEN** access is denied by the standard RBAC middleware

