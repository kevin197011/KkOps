## 1. Backend

- [x] 1.1 `internal/integration/k8s` registry and kubeconfig parsing (`client-go`).
- [x] 1.2 K8s HTTP handler and `/api/v1/k8s/clusters/:id/...` routes.
- [x] 1.3 Provider kind `kubernetes` stub + registry registration.

## 2. Frontend

- [x] 2.1 `pages/k8s/Clusters.tsx` and `ClusterDetail.tsx` with shell components.
- [x] 2.2 API client and navigation under **Kubernetes**.

## 3. RBAC & spec

- [x] 3.1 `kubernetes:*` in permissions, route map, menu, App routes, permission store.
- [x] 3.2 Capability spec and `openspec validate add-kubernetes-multi-cluster --strict`.
