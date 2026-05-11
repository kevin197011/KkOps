# Change: Unified operations dashboard summary

## Why

The home dashboard should surface cross-cutting operational signals (alerts, incidents, integrations, provisioning, AI anomalies) in one glance.

## What Changes

- Backend **`GET /api/v1/dashboard/summary`** aggregating counts from existing tables (open alerts, open incidents, integrations, recent provisioning runs, AI anomaly findings in the last 24 hours).
- Frontend enhances **`Dashboard.tsx`** with summary cards and shortcuts (Alerts, Incidents, Integrations, AI chat). Existing **`GET /dashboard/stats`** remains unchanged.
- RBAC: reuse **`dashboard:read`** (same as dashboard home).

## Impact

- Affected specs: `unified-dashboard`
- Affected code: `backend/internal/service/dashboard`, `backend/internal/handler/dashboard`, `frontend/src/pages/Dashboard.tsx`, `frontend/src/api/dashboard.ts`.
