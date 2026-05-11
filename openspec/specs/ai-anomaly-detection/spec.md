# ai-anomaly-detection Specification

## Purpose
TBD - created by archiving change add-ai-anomaly-detection. Update Purpose after archive.
## Requirements
### Requirement: Anomaly rules

The system SHALL store `AIAnomalyRule` records with integration reference, PromQL query, cron schedule, enabled flag, and prompt template. Operators SHALL manage rules via `CRUD /api/v1/ai/anomaly/rules` with permission `ai:*`.

#### Scenario: Create rule

- **WHEN** an authorized user creates a rule with valid cron and integration id
- **THEN** the server SHALL persist the rule and the worker SHALL pick it up on schedule.

### Requirement: Findings and incidents

The system SHALL persist `AIAnomalyFinding` rows with severity and summary. When severity is `critical`, the worker MAY create an incident using the existing incident service.

#### Scenario: List findings

- **WHEN** `GET /api/v1/ai/anomaly/findings` is called with optional `rule_id` and `since`
- **THEN** the server SHALL return matching findings ordered by timestamp descending.

