## 1. OpenSpec

- [x] 1.1 Add proposal, tasks, and `specs/account-provisioning/spec.md`
- [x] 1.2 Run `openspec validate add-account-provisioning --strict`

## 2. Backend

- [x] 2.1 Add `provisioning_targets`, `provisioning_runs`, and user–external mapping models + AutoMigrate
- [x] 2.2 Implement Provider interface, Registry, HTTP stubs (SCIM, GitLab, Jenkins, Grafana, Harbor, ArgoCD, Jumpserver, Nightingale)
- [x] 2.3 Worker pool + enqueue from user create/update/delete + integration decrypt for outbound HTTP
- [x] 2.4 HTTP handlers: list targets, manual sync; seed `provisioning:*` permission and route map

## 3. Frontend

- [x] 3.1 Page `ProvisioningTargets.tsx` with PageContainer / PageHeader / EmptyState
- [x] 3.2 Menu「账号同步」under 系统管理 + route + API client

## 4. Verification

- [x] 4.1 `go build ./...` and `npm run build`
