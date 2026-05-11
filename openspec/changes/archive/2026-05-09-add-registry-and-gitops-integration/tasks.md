## 1. Backend — Harbor

- [x] 1.1 `internal/integration/registry`: Harbor client (repositories, tags, vulnerabilities).
- [x] 1.2 REST routes under `/api/v1/registry` with `registry:*`.

## 2. Backend — Argo CD

- [x] 2.1 `internal/integration/gitops`: Argo CD client (list, sync).
- [x] 2.2 REST routes under `/api/v1/gitops` with `gitops:*`.

## 3. Frontend

- [x] 3.1 `Repositories.tsx` for Harbor browsing.
- [x] 3.2 `Applications.tsx` for Argo CD applications and sync.

## 4. Spec

- [x] 4.1 Capabilities `registry-integration` and `gitops-integration`.
- [x] 4.2 `openspec validate add-registry-and-gitops-integration --strict` passes.
