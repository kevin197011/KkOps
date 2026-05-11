# Proposal: Logging integration (Loki / Elasticsearch)

## Why

Operators need a single search experience across log backends tied to encrypted integration credentials.

## What

- Loki and Elasticsearch clients with `Search(ctx, query, start, end, limit)` returning `LogLine` records.
- `POST /api/v1/logging/search`; RBAC `logging:*`.
- Frontend `/logging` with integration picker, time range, scrollable log list.

## Scope

Read-only search; saved queries and deep links are future work.
