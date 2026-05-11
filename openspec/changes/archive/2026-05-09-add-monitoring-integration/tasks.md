## 1. Backend

- [x] 1.1 `internal/integration/monitoring`: Prometheus + Nightingale clients, `MetricSeries` model.
- [x] 1.2 `POST /api/v1/monitoring/query` with `monitoring:*` permission.

## 2. Frontend

- [x] 2.1 `MonitoringQuery.tsx` with integration picker, instant/range mode, table + SVG sparkline.

## 3. Spec

- [x] 3.1 Capability `monitoring-query` spec.
- [x] 3.2 `openspec validate add-monitoring-integration --strict` passes.
