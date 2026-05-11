## 1. Implementation

- [x] Add `AlertRecord` model, GORM migrate, normalize Alertmanager + webhook JSON
- [x] Alert service: list, ack/dismiss, sync from Prometheus Alertmanager API + Nightingale stub
- [x] Handler + routes: GET/POST `/api/v1/alerts`, ack/dismiss, POST webhook with secret header
- [x] Config `alerts.webhook_secret` + middleware for webhook
- [x] Frontend Alert Center page + API client + RBAC `alerts:*`
- [x] `openspec validate add-alert-aggregator --strict`
