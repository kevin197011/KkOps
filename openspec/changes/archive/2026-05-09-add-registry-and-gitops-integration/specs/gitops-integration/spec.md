## ADDED Requirements

### Requirement: GitOps APIs

The system SHALL expose HTTP endpoints to list Argo CD applications for an argocd integration and to request sync for a named application. Endpoints SHALL require `gitops:*`.

#### Scenario: Sync failure

- **WHEN** Argo CD rejects or fails a sync operation
- **THEN** the server SHALL return a non-2xx response with an explanatory error body.

### Requirement: GitOps UI

The web application SHALL provide `/gitops` to list applications with health and sync status and to trigger sync for a selected application.

#### Scenario: Sync from UI

- **WHEN** the user clicks sync for an application row and the API succeeds
- **THEN** the UI SHALL confirm success and refresh or update the application list state.
