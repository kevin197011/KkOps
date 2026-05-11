## 1. Data model and crypto

- [x] 1.1 `Integration` model with encrypted `config_encrypted` column.
- [x] 1.2 Use `utils.Encrypt` / `utils.Decrypt` with `encryption.key` from config.

## 2. API and RBAC

- [x] 2.1 REST handlers: list, get, create, update, delete under `/api/v1/integrations`.
- [x] 2.2 Register `integrations:*` in `AllMenuPermissions` and `RoutePermissionMap`.

## 3. Spec validation

- [x] 3.1 Capability spec `integrations-framework` with MUST requirements.
- [x] 3.2 `openspec validate add-backend-integrations-framework --strict` passes.
