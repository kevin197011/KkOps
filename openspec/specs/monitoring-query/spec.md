# monitoring-query Specification

## Purpose
TBD - created by archiving change add-monitoring-integration. Update Purpose after archive.
## Requirements
### Requirement: Monitoring query API

The system SHALL expose `POST /api/v1/monitoring/query` accepting `integration_id`, `query`, optional instant `time`, or optional `range` with `start`, `end`, and `step`. The server SHALL require `monitoring:*` and SHALL return normalized series suitable for charting.

#### Scenario: Upstream failure

- **WHEN** the upstream Prometheus-compatible endpoint is unreachable or returns an error
- **THEN** the server SHALL respond with a non-2xx HTTP status and a typed error message without panicking.

### Requirement: Monitoring query UI

The web application SHALL provide a `/monitoring` page to pick a prometheus or nightingale integration, enter PromQL, run instant or range queries, and display results in a table with a simple line visualization.

#### Scenario: Successful query display

- **WHEN** the user selects an integration, submits a valid PromQL query, and the API returns series data
- **THEN** the UI SHALL render a results table and show at least one polyline preview for the first series.

