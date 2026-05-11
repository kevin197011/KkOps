## ADDED Requirements

### Requirement: RCA generation API

The system SHALL expose `POST /api/v1/ai/rca` with body containing `incident_id` and optional `integration_id_overrides`. The response SHALL include `report_md` in Markdown and `cited_tool_calls` describing tool invocations used for evidence.

#### Scenario: Persist report

- **WHEN** RCA generation succeeds
- **THEN** the server SHALL store an `AIRcaReport` row linked to the incident with raw tool log for audit.

### Requirement: Incident UI integration

The incident detail experience SHALL offer a control labeled for AI root-cause analysis that triggers RCA generation and SHALL link to a reports list filtered by incident where applicable.

#### Scenario: Unauthorized

- **WHEN** a user lacks `ai:*`
- **THEN** RCA APIs SHALL return forbidden and UI controls SHALL be hidden per RBAC rules.
