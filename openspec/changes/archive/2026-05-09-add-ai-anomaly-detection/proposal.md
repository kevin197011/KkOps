# Change: AI anomaly detection

## Why

Scheduled metric evaluation plus LLM reasoning can surface anomalies early and optionally open critical incidents.

## What Changes

- Models `AIAnomalyRule`, `AIAnomalyFinding`; CRUD rules APIs and list findings API.
- In-process cron (`robfig/cron/v3`) executing rules against monitoring integrations; persists findings; optional critical incident creation.

## Impact

- Affected specs: `ai-anomaly-detection`
- Affected code: `backend/internal/service/aisvc/anomaly.go`, `main.go` worker startup, frontend rule and finding pages.
