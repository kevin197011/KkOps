# Proposal: Kubernetes multi-cluster browsing

## Why

Operators need namespace/workload/pod visibility across clusters from KkOps without exposing kubeconfig to browsers.

## What

- Encrypted `kubeconfig` in integration kind `kubernetes`; server-side `client-go` with cached clientsets.
- REST under `/api/v1/k8s/clusters/:id/*` for namespaces, workloads, pods, logs, events, nodes.
- RBAC `kubernetes:*`; frontend Kubernetes cluster list and detail tabs.

## Impact

- Spec capability: `kubernetes-multi-cluster`
- Backend: `internal/integration/k8s`, `handler/k8s`; `permissions`, routes.
