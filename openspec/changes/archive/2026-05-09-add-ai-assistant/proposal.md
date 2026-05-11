# Change: AI assistant (chat ops)

## Why

Operators need a conversational assistant that streams responses and can invoke read-only platform tools against alerts, incidents, logs, metrics, Kubernetes, pipelines, and integrations.

## What Changes

- `POST /api/v1/ai/chat` with SSE token streaming; tool bridge via `<<TOOL: name args>>` markers and optional native tool calls when supported.
- Session persistence (`AIChatSession`, `AIChatMessage`) and CRUD session APIs.
- Frontend `/ai/chat` with session rail and markdown rendering; Cmd/Ctrl+L shortcut; nav under AI 运维.

## Impact

- Affected specs: `ai-assistant`
- Affected code: `backend/internal/service/aisvc`, `backend/internal/handler/ai`, `frontend/src/pages/ai/Chat.tsx`, `App.tsx`, `MainLayout.tsx`.
