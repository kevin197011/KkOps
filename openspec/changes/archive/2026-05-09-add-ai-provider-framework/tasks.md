## 1. Implementation

- [x] Add `internal/ai` Provider interface and OpenAI-compatible, Anthropic, webhook providers
- [x] Registry loading AI integrations from DB; `Default()` picks first enabled `kind=ai`
- [x] REST `GET /api/v1/ai/providers`, `POST /api/v1/ai/test`
- [x] RBAC `ai:*`, `RoutePermissionMap`, seed menu permission
- [x] Spec delta and validation
