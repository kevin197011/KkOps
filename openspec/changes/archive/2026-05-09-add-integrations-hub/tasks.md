## 1. Backend

- [x] 1.1 `internal/integration/provider`: `Provider`, `Registry`, stub providers per kind, `VerifyHTTPGET`.
- [x] 1.2 `POST /api/v1/integrations/:id/test` wired to registry `Verify`.
- [x] 1.3 RBAC: reuse `integrations:*` in `AllMenuPermissions` and `RoutePermissionMap`.

## 2. Frontend

- [x] 2.1 Page `IntegrationsHub.tsx` with catalog grid, table, wizard, edit/delete, test connection.
- [x] 2.2 Register route and sidebar entry under menu group 集成中心.

## 3. Spec

- [x] 3.1 Capability `integrations-hub` spec with MUST requirements.
- [x] 3.2 `openspec validate add-integrations-hub --strict` passes.
