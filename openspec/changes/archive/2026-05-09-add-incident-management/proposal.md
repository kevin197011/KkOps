# Change: Incident management (minimal)

## Why

Teams need to track operational incidents with severity and lifecycle, optionally linked to aggregated alerts.

## What Changes

- Minimal incident CRUD: title, severity, status (open/acknowledged/resolved), optional JSON array of linked alert IDs, optional assignee user ID.
- GORM AutoMigrate models; REST APIs under `/api/v1/incidents`.
- Frontend **事件** page: list, create modal, detail drawer.
- RBAC `incidents:*`.

## Impact

- Affected specs: `incident-management`
- Affected code: `backend/internal/{model,handler/incident,service/incident}`, `frontend/src/pages/incidents`, routing and permissions.
