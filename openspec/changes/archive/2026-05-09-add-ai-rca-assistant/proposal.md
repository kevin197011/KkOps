# Change: AI root cause analysis

## Why

For incidents, correlated alerts, logs, metrics, K8s events, and pipelines should be summarized into a structured RCA report operators can share.

## What Changes

- `POST /api/v1/ai/rca` builds Markdown RCA from incident id; persists `AIRcaReport`.
- Incident detail UI: button to trigger RCA; listing page for reports per incident.

## Impact

- Affected specs: `ai-rca`
- Affected code: `backend/internal/service/aisvc/rca.go`, incident UI, `pages/ai/RcaReports.tsx`.
