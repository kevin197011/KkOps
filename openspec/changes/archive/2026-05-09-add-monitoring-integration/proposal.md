# Proposal: Monitoring integration (Prometheus / Nightingale)

## Why

Unified PromQL execution against configured Prometheus or Nightingale endpoints reduces context switching and enables asset-centric workflows later.

## What

- Go clients with `Query` and `QueryRange` returning normalized `MetricSeries`.
- `POST /api/v1/monitoring/query` with instant or range body; RBAC `monitoring:*`.
- Frontend `/monitoring` page: integration picker, PromQL, table + sparkline preview.

## Scope

Read-only query proxy; alerting and Grafana embed are out of scope for this change.
