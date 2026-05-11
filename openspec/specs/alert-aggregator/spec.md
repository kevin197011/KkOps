# alert-aggregator Specification

## Purpose
TBD - created by archiving change add-alert-aggregator. Update Purpose after archive.
## Requirements
### Requirement: Normalized alert storage

The system SHALL persist alerts with fields including severity, title, status, source kind, optional integration reference, label map (JSON), fingerprint for deduplication, and start/end timestamps.

#### Scenario: Webhook ingest

- **WHEN** a client sends a valid Alertmanager-style JSON payload to `POST /api/v1/alerts/webhook` with the configured secret header
- **THEN** the server SHALL normalize each alert, upsert by fingerprint, and respond with success without panicking.

### Requirement: Alert APIs

The system SHALL expose authenticated APIs under `/api/v1/alerts` to list alerts, optionally sync from configured monitoring integrations, acknowledge an alert, and dismiss an alert. The server SHALL require `alerts:*`.

#### Scenario: Upstream integration failure

- **WHEN** sync calls an external Alertmanager or Nightingale endpoint and it returns an error
- **THEN** the API SHALL return a non-2xx or partial success message without panicking.

### Requirement: Alert Center UI

The web application SHALL provide a `/alerts` page listing alerts with severity, source integration name, and labels, plus actions for acknowledge and dismiss and a placeholder area for future silence rules.

#### Scenario: Empty state

- **WHEN** no alerts exist and sync returns none
- **THEN** the UI SHALL show an empty state without breaking navigation.

