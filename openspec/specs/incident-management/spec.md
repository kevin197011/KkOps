# incident-management Specification

## Purpose
TBD - created by archiving change add-incident-management. Update Purpose after archive.
## Requirements
### Requirement: Incident persistence

The system SHALL store incidents with title, severity, status (`open`, `acknowledged`, `resolved`), optional linked alert IDs as JSON array, optional assignee user ID, and timestamps.

#### Scenario: Create incident

- **WHEN** an authorized user submits `POST /api/v1/incidents` with a title and severity
- **THEN** the server SHALL create an incident with status `open` and return its identifier.

### Requirement: Incident APIs

The system SHALL expose `GET` and `POST` on `/api/v1/incidents`, and `GET` and `PATCH` on `/api/v1/incidents/:id`. All SHALL require `incidents:*`.

#### Scenario: Update status

- **WHEN** an authorized user sends `PATCH /api/v1/incidents/:id` with a valid status transition
- **THEN** the server SHALL update the incident and return the updated record.

### Requirement: Incidents UI

The web application SHALL provide a `/incidents` page with a table of incidents, a create modal, and a detail drawer showing fields and linked alert IDs.

#### Scenario: Unauthorized access

- **WHEN** a user lacks `incidents:*`
- **THEN** the route SHALL not be reachable in the shell menu and API calls SHALL return forbidden.

