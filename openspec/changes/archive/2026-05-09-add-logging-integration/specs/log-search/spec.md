## ADDED Requirements

### Requirement: Log search API

The system SHALL expose `POST /api/v1/logging/search` with `integration_id`, `query`, optional `start`, `end`, and `limit`. The server SHALL require `logging:*` and SHALL return a JSON array of normalized log lines.

#### Scenario: Backend unavailable

- **WHEN** Loki or Elasticsearch returns an error or cannot be reached
- **THEN** the server SHALL return HTTP 502 (or equivalent) with an error message and SHALL NOT panic.

### Requirement: Log search UI

The web application SHALL provide `/logging` to select a loki or elasticsearch integration, enter a query string, choose a time range, and view results in a scrollable monospace list.

#### Scenario: Search displays lines

- **WHEN** the user selects an integration, enters a query, chooses a time window, and the API returns log lines
- **THEN** the UI SHALL render each line with timestamp and message in a scrollable region.
