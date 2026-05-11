# Proposal: Backend integrations framework

## Why

Third-party connectors (alerting, ticketing, GitOps, CMDB) need a shared pattern: persisted integration records with encrypted credential payloads and RBAC-protected APIs.

## What

- `integrations` table and GORM model.
- AES-GCM encryption at rest for JSON credential blobs using the existing `encryption.key`.
- REST CRUD under `/api/v1/integrations` with `integrations:*` RBAC.
- API responses MUST NOT return decrypted secrets (metadata only).

## Scope

Initial delivery is framework-only (CRUD + crypto); connector-specific workers are follow-ups.
