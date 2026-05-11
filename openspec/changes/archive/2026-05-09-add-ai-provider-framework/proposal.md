# Change: AI provider framework

## Why

KkOps needs a pluggable LLM layer so operators can attach OpenAI-compatible, Anthropic, or custom webhook backends using the existing integrations store.

## What Changes

- Backend `internal/ai` package with `Provider` interface and HTTP-based implementations (OpenAI-compatible, Anthropic Messages API, generic webhook).
- AI integrations use `integrations.kind = ai` with encrypted JSON config.
- APIs `GET /api/v1/ai/providers` and `POST /api/v1/ai/test`.
- RBAC permission `ai:*` and route mapping for `/api/v1/ai`.

## Impact

- Affected specs: `ai-provider-framework`
- Affected code: `backend/internal/ai`, `backend/internal/model/permissions.go`, `backend/cmd/server/main.go`, frontend permission store (admin menu sync).
