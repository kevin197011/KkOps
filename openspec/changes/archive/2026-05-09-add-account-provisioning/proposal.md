# Change: Account provisioning to external systems

## Why

KkOps should act as an identity hub by syncing user lifecycle events to external tools (SCIM, GitLab, CI/CD, monitoring, etc.) using the integrations framework for encrypted credentials.

## What Changes

- OpenSpec capability `account-provisioning` with requirements for targets, runs, async jobs, and APIs.
- Backend: `provisioning` package with `Provider` interface, registry, HTTP-based provider stubs, persistence for targets and runs, in-process worker queue, hooks from user CRUD, REST endpoints with `provisioning:*` permission.
- Frontend: provisioning targets page under 系统管理 with list, sync, and run history.

## Impact

- Affected specs: `account-provisioning` (new delta).
- Affected code: `backend/internal/provisioning/`, `backend/internal/handler/`, `backend/internal/model/`, `backend/cmd/server/main.go`, `frontend/src/pages/provisioning/`, `MainLayout.tsx`.
