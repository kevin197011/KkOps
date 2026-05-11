## ADDED Requirements

### Requirement: Dashboard summary API

The system SHALL expose `GET /api/v1/dashboard/summary` requiring **`dashboard:read`**. The response SHALL include aggregated counters such as: open alerts, open incidents, integration totals and healthy counts (enabled integrations treated as healthy for MVP), recent provisioning run count (last 24 hours), and AI anomaly findings count (last 24 hours), plus any stable keys documented in the handler.

#### Scenario: Authenticated summary

- **WHEN** an authenticated user with `dashboard:read` requests `GET /api/v1/dashboard/summary`
- **THEN** the server SHALL return a JSON object with numeric summary fields.

### Requirement: Unified dashboard UI

The web application SHALL enhance the main dashboard page to display summary metrics (using Ant Design **Statistic** or equivalent) and SHALL provide navigation shortcuts to Alerts, Incidents, Integrations hub, and AI assistant where permissions allow.

#### Scenario: Operators see signals

- **WHEN** a user opens the dashboard home
- **THEN** they SHALL see operational summary cards populated from `/api/v1/dashboard/summary` alongside existing `/dashboard/stats` content.
