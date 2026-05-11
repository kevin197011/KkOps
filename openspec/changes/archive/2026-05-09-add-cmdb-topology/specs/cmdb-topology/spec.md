## ADDED Requirements

### Requirement: CMDB configuration items

The system SHALL persist configuration items in table `cmdb_assets` with fields: name, kind (`service`, `host`, `database`, `cluster`, `other`), environment label, optional owner user ID, labels (JSONB), optional integration ID, optional external reference, notes, and timestamps. This SHALL be separate from infrastructure host records in table `assets`.

#### Scenario: List with filters

- **WHEN** an authorized user calls `GET /api/v1/cmdb/assets` with optional query parameters `kind`, `env`, and `q` (search on name and notes)
- **THEN** the server SHALL return a paginated list of CMDB items.

### Requirement: CMDB REST APIs

The system SHALL expose CRUD on `/api/v1/cmdb/assets` and `/api/v1/cmdb/assets/:id`, and list/create/delete on `/api/v1/cmdb/asset-relations` and `/api/v1/cmdb/asset-relations/:id`. All SHALL require permission `cmdb:*`.

#### Scenario: Create relation

- **WHEN** an authorized user posts valid `from_asset_id`, `to_asset_id`, and `relation_type` (`depends_on`, `runs_on`, `calls`)
- **THEN** the server SHALL store the relation and SHALL reject self-references.

### Requirement: Topology graph API

The system SHALL expose `GET /api/v1/topology/graph` returning JSON with `nodes` (CMDB assets plus optional derived placeholder metadata), `edges` from `asset_relations`, and integration-derived counts (e.g. Kubernetes connectors, Harbor connectors) without requiring a graph database.

#### Scenario: Read topology

- **WHEN** an authorized user calls `GET /api/v1/topology/graph`
- **THEN** the server SHALL return nodes and edges suitable for a simple UI visualization.

### Requirement: Topology UI choice

The web application SHALL implement **`/topology`** without adding heavy graph libraries (for example `@antv/g6`). The UI SHALL use Ant Design **Tree** for hierarchy and a table or list for edges.

#### Scenario: Navigate topology

- **WHEN** a user with `cmdb:*` opens `/topology`
- **THEN** they SHALL see CMDB-backed structure and dependency edges without loading a force-directed canvas library.

### Requirement: CMDB UI

The web application SHALL provide **`/cmdb`** with a table and forms to manage CMDB assets.

#### Scenario: Unauthorized access

- **WHEN** a user lacks `cmdb:*`
- **THEN** CMDB and topology routes SHALL be hidden from the shell menu and APIs SHALL return forbidden.
