# Change: Alert aggregator

## Why

Operators need a single place to view alerts from Prometheus Alertmanager, Nightingale, and webhook ingress, with consistent severity and lifecycle (acknowledge/dismiss).

## What Changes

- Normalize external alert payloads into a common `Alert` model persisted in Postgres.
- Backend APIs: list alerts, sync from integrations, acknowledge/dismiss, optional `POST /api/v1/alerts/webhook` protected by a configurable secret header.
- Frontend **告警** page with table (severity, source integration, labels) and silence stub placeholder.
- RBAC resource `alerts:*` and route registration.

## Impact

- Affected specs: `alert-aggregator`
- Affected code: `backend/internal/{model,handler/alert,service/alert,integration/monitoring}`, `frontend/src/pages/alerts`, `MainLayout`, `App.tsx`, permissions.
