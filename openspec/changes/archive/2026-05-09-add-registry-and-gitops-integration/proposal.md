# Proposal: Registry (Harbor) and GitOps (Argo CD)

## Why

Image inventory and GitOps application state are core operational views that should reuse the integrations framework.

## What

- Harbor client: repositories, tags, vulnerability summary fetch.
- Argo CD client: list applications, sync application.
- REST under `/api/v1/registry` and `/api/v1/gitops` with `registry:*` and `gitops:*`.
- Frontend `/registry` and `/gitops` pages.

## Scope

Deep vulnerability analytics and multi-cluster Argo RBAC are out of scope.
