# Proposal: Integrations hub UI and provider verification

## Why

Operators need one place to configure third-party connectors, verify connectivity, and see status across Prometheus, Nightingale, Loki, Elasticsearch, Grafana, Jenkins, GitLab, Harbor, and Argo CD.

## What

- Shared `Provider` interface (`Kind`, `Verify`, `Metadata`) with a kind-keyed registry and HTTP-based verification stubs.
- REST `POST /api/v1/integrations/:id/test` invoking `Provider.Verify` with decrypted credentials.
- Frontend route `/integrations` with catalog cards, CRUD wizard, test connection, RBAC `integrations:*`.

## Scope

Connector-specific business logic beyond HTTP probe remains minimal (stubs); Phase 2 follow-ups add domain APIs per subsystem.
