# Change: CMDB and topology (MVP)

## Why

Operators need a lightweight configuration inventory and a simple dependency view without adopting a full enterprise CMDB product.

## What Changes

- **CMDB models** (`cmdb_assets`, `asset_relations`) distinct from existing host-centric `assets` table.
- REST APIs under **`/api/v1/cmdb/assets`** and **`/api/v1/cmdb/asset-relations`** (not `/api/v1/assets`, which remains infrastructure hosts).
- **`GET /api/v1/topology/graph`** returns CMDB nodes, edges from relations, and minimal derived counts from integrations (Kubernetes/Harbor placeholders).
- Frontend: **`/cmdb`** (CRUD table) and **`/topology`** (Ant Design Tree hierarchy + edge list; no `@antv/g6` to keep bundle small).
- RBAC **`cmdb:*`**.

## Impact

- Affected specs: `cmdb-topology`
- Affected code: `backend/internal/{model,service/cmdb,handler/cmdb,handler/topology}`, `frontend/src/pages/cmdb`, `frontend/src/pages/topology`, routing, permissions.
