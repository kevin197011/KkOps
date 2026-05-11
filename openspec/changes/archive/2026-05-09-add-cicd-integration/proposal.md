# Proposal: CI/CD integration (Jenkins / GitLab CI)

## Why

Pipeline visibility and trigger from KkOps shortens the path from change to deployment verification.

## What

- Jenkins and GitLab CI clients: list pipelines, get detail, trigger run, fetch logs.
- REST: `GET /api/v1/cicd/pipelines`, `POST /api/v1/cicd/pipelines/:id/run`, `GET /api/v1/cicd/pipelines/:id/logs` with `cicd:*`.
- Frontend `/cicd` listing pipelines with run modal and logs drawer.

## Scope

Deployment-module linkage is deferred.
