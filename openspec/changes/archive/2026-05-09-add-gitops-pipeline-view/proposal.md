# Proposal: Unified GitOps pipeline timeline

## Why

Release visibility improves when CI runs and Argo CD syncs appear on one timeline per application.

## What

- Service aggregates recent CI/CD pipeline runs, Argo CD application revision history, and optional metadata into `PipelineEvent` records.
- `GET /api/v1/gitops/pipeline-view` with `app` filter; RBAC aligned with existing gitops routes.

## Impact

- Spec capability: `gitops-pipeline-view`
- Backend: `service/gitopsview`, Argo client history helper; frontend `PipelineView.tsx`.
