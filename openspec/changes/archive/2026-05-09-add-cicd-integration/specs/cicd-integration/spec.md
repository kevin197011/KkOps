## ADDED Requirements

### Requirement: CI/CD proxy APIs

The system SHALL expose Git-backed CI operations under `/api/v1/cicd` including listing pipelines for an integration, triggering a pipeline run with optional variables, and retrieving textual logs. All endpoints SHALL require `cicd:*`.

#### Scenario: Trigger rejected

- **WHEN** the upstream CI system rejects a trigger or logs request
- **THEN** the server SHALL surface the failure as an HTTP error with a clear message.

### Requirement: CI/CD UI

The web application SHALL provide `/cicd` to choose a jenkins or gitlab integration, list pipelines, trigger runs, and view logs in a drawer.

#### Scenario: View pipeline logs

- **WHEN** the user opens the logs action for a pipeline row and the logs API returns text
- **THEN** the UI SHALL show the log body inside a drawer or modal without navigating away from the list.
