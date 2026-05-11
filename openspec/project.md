# Project Context

## Purpose

KkOps is an Intelligent Operations Platform that unifies identity, asset management,
remote execution (WebSSH + script tasks), deployment, audit, observability and
AI-assisted operations across an organization's open-source ops toolchain.

The platform plays three roles:

1. **Identity Hub** — KkOps is the central account / role / permission source of
   truth. It is both an OIDC IdP for downstream tools (GitLab, Grafana, Jenkins,
   Harbor, ArgoCD, Jumpserver, Nightingale, ...) and a provisioner that pushes
   accounts and groups outward.
2. **Operations Cockpit** — A single console that aggregates monitoring, logging,
   CI/CD, registry, GitOps, Kubernetes and on-call data from multiple integrated
   tools, with consistent navigation and audited actions.
3. **AI Ops Brain** — A function-calling LLM assistant that can query metrics
   and logs, run approved Runbooks, generate scripts/SQL/YAML and analyze
   alert root causes.

## Tech Stack

### Backend
- Go 1.24+
- Gin (HTTP), GORM (PostgreSQL), Redis
- JWT (auth), Cookie sessions (IdP), RS256 OIDC
- gorilla/websocket (WebSSH, log streaming)
- viper (configuration), zap (logging)
- swaggo (Swagger docs)

### Frontend
- React 18 + TypeScript 5
- Vite 5, Ant Design 5, `@ant-design/icons`
- React Router v6, Zustand (state)
- Axios (HTTP), xterm.js (terminals)

### Infrastructure
- Docker Compose for development
- PostgreSQL 16, Redis 7

## Project Conventions

### Code Style

- Comments and logs MUST be written in English. User-visible strings go through
  i18n (currently embedded zh-CN in pages until we externalize).
- Backend: standard Go formatting (`gofmt`, `golangci-lint`). Package layout
  follows `internal/{handler,service,model,middleware,config,integration,...}`.
- Frontend: ESLint + TypeScript strict. Prefer functional components, hooks,
  Zustand stores; no Redux.
- License headers: every new source file starts with the SPDX-style block used
  across the repo (`Copyright (c) 2025 kk` / MIT).
- Avoid noisy comments; comments explain *why*, not *what*.

### Architecture Patterns

- **Layered backend**: Handler -> Service -> Model (GORM). Services contain
  business rules; handlers only translate HTTP <-> service calls.
- **RBAC + menu permissions**: see `backend/internal/model/permissions.go` and
  `backend/internal/middleware/permission.go`. New menu modules MUST register
  a permission (resource + action) and a route mapping.
- **Integrations**: every external system uses the
  `backend/internal/integration` framework (Provider interface + registry +
  encrypted credential store + http factory + webhook hub).
- **Identity**: KkOps owns the user/role tables. Outbound provisioning, OIDC
  IdP (`backend/internal/idp`) and SSO inbound (`backend/internal/service/auth/sso.go`)
  all reference the same canonical user.
- **Frontend shell**: `frontend/src/layouts/MainLayout.tsx` is the global
  shell (sider + topbar + theme + permission filter). New top-level pages must
  add a route in `App.tsx`, an item in `MainLayout.tsx` menu, and a permission
  entry both in backend `permissions.go` and frontend `stores/permission.ts`.

### Testing Strategy

- Test scripts are written in **Ruby 3.1+** and live under `scripts/`. They
  must be non-interactive and CI-friendly (per `.cursor/rules/`).
- Backend unit tests: `go test ./...` against an isolated database via
  Docker Compose.
- Frontend: TypeScript `tsc` build is the smoke check today; component tests
  may be added per feature.

### Git Workflow

- Conventional Commits: `feat(scope): ...`, `fix(...): ...`, `docs(...)`,
  `refactor(...)`, etc.
- Each substantive change MUST be backed by an OpenSpec change proposal
  (`openspec/changes/<id>/`) before code lands. Bug fixes that restore spec
  behavior may skip the proposal.

## Domain Context

- **Resource hierarchy**: Project -> Environment -> Asset (host / container /
  service). Tags add cross-cutting grouping.
- **Execution**: Templates (reusable scripts/playbooks) + Executions (a run on
  selected assets) + Scheduled Tasks (cron-driven Executions). Logs stream over
  WebSocket.
- **Identity**: User -> Role -> Permission (menu + asset). The `Source` field
  on `User` distinguishes `local` from `sso` accounts.
- **Audit**: Every protected mutation is captured by `AuditMiddleware`. WebSSH
  sessions are recorded into `connection_audit` with replayable transcripts.

## Important Constraints

- Must run end-to-end inside `docker-compose` without external services for
  development. Real integrations (Prometheus, GitLab, ...) are configured at
  runtime by an operator.
- All credentials (SSH private keys, IdP client secrets, integration tokens,
  kubeconfigs) are encrypted at rest using `Encryption.Key`. Plaintext secrets
  are never returned in API responses.
- AI features are opt-in: no LLM calls happen unless a provider is configured
  in the AI settings page; high-risk tool actions called by AI must be gated
  behind explicit human approval.

## External Dependencies

- PostgreSQL (state of record).
- Redis (cache, future task queues).
- Optional integrations (added phase by phase): Prometheus, Nightingale,
  Loki / Elasticsearch, Grafana, Jenkins, GitLab, Harbor, ArgoCD,
  Kubernetes API, OpenAI / Anthropic / Ollama.
