# Tasks: Add Cloud Platform Management

## Phase 1: Backend Model & Database

- [x] 1.1 Create `CloudPlatform` model in `backend/internal/model/cloud_platform.go`
- [x] 1.2 Update `Asset` model to replace `CloudPlatform string` with `CloudPlatformID *uint` and relationship
- [x] 1.3 Create database migration script to:
  - Create `cloud_platforms` table
  - Extract unique values from `assets.cloud_platform`
  - Create cloud platform entries
  - Add `cloud_platform_id` column to `assets` table
  - Migrate data from string to foreign key
  - Drop `cloud_platform` column

## Phase 2: Backend Service Layer

- [x] 2.1 Create service package `backend/internal/service/cloudplatform/`
- [x] 2.2 Implement `CreateCloudPlatformRequest`, `UpdateCloudPlatformRequest`, `CloudPlatformResponse` structs
- [x] 2.3 Implement CRUD methods:
  - `CreateCloudPlatform`
  - `GetCloudPlatform`
  - `ListCloudPlatforms`
  - `UpdateCloudPlatform`
  - `DeleteCloudPlatform`

## Phase 3: Backend Handler & Router

- [ ] 3.1 Create handler package `backend/internal/handler/cloudplatform/`
- [ ] 3.2 Implement HTTP handlers:
  - `CreateCloudPlatform`
  - `GetCloudPlatform`
  - `ListCloudPlatforms`
  - `UpdateCloudPlatform`
  - `DeleteCloudPlatform`
- [ ] 3.3 Register routes in `backend/cmd/server/main.go`

## Phase 4: Backend Asset Service Updates

- [x] 4.1 Update `CreateAssetRequest` to use `CloudPlatformID *uint` instead of `CloudPlatform string`
- [x] 4.2 Update `UpdateAssetRequest` to use `CloudPlatformID *uint`
- [x] 4.3 Update `AssetResponse` to include cloud platform relationship
- [x] 4.4 Update `ListAssets` filter to use `CloudPlatformID` instead of string matching
- [x] 4.5 Update asset creation/update logic to handle cloud platform ID

## Phase 5: Frontend API Client

- [x] 5.1 Create `frontend/src/api/cloudPlatform.ts` with:
  - `CloudPlatform` interface
  - `CreateCloudPlatformRequest`, `UpdateCloudPlatformRequest` interfaces
  - API methods: `list`, `get`, `create`, `update`, `delete`

## Phase 6: Frontend Page Component

- [x] 6.1 Create `frontend/src/pages/cloudPlatforms/CloudPlatformList.tsx`
- [x] 6.2 Implement table view with columns: ID, Name, Description, Actions
- [x] 6.3 Implement create/edit modal form
- [x] 6.4 Implement delete confirmation dialog
- [x] 6.5 Follow same pattern as `EnvironmentList.tsx`

## Phase 7: Frontend Asset Form Updates

- [ ] 7.1 Update `frontend/src/api/asset.ts`:
  - Remove `cloud_platform: string` from `Asset` interface
  - Add `cloud_platform_id?: number` and `cloud_platform?: CloudPlatformInfo`
  - Update request interfaces
- [ ] 7.2 Update `frontend/src/pages/assets/AssetList.tsx`:
  - Fetch cloud platforms list
  - Update table to display cloud platform name instead of string
- [ ] 7.3 Update asset create/edit form:
  - Replace text input with dropdown select
  - Populate options from cloud platforms API

## Phase 8: Frontend Routing & Menu

- [x] 8.1 Add route in `frontend/src/App.tsx` for `/cloud-platforms`
- [x] 8.2 Add menu item in `frontend/src/layouts/MainLayout.tsx`:
  - Import `CloudOutlined` icon from `@ant-design/icons`
  - Add menu item after "环境管理"
  - Set route to `/cloud-platforms`

## Phase 9: Database Migration & Testing

- [ ] 9.1 Run database migration on test environment
- [ ] 9.2 Verify existing asset data is correctly migrated
- [ ] 9.3 Test cloud platform CRUD operations
- [ ] 9.4 Test asset creation with cloud platform selection
- [ ] 9.5 Test asset update with cloud platform change
- [ ] 9.6 Verify asset list displays cloud platform names correctly
- [ ] 9.7 Test asset filtering by cloud platform

## Phase 10: Documentation & Validation

- [ ] 10.1 Update API documentation
- [ ] 10.2 Document migration process
- [ ] 10.3 Run OpenSpec validation: `openspec validate add-cloud-platform-management --strict`
