# Tasks: Rename Asset `name` Field to `hostName`

## Phase 1: Database Migration

- [x] 1.1 Create database migration script to rename `name` column to `host_name`
- [x] 1.2 Verify migration script handles indexes and constraints correctly
- [x] 1.3 Test migration script on development database (manual testing required)
- [x] 1.4 Document rollback procedure for migration (included in migration script)

## Phase 2: Backend Model Changes

- [x] 2.1 Update Asset model struct field from `Name` to `HostName`
- [x] 2.2 Update GORM tag (column name explicitly set to `host_name`)
- [x] 2.3 Update JSON tag from `json:"name"` to `json:"hostName"`

## Phase 3: Backend Service Layer

- [x] 3.1 Update `CreateAssetRequest` struct field `Name` to `HostName` with `json:"hostName"` tag
- [x] 3.2 Update `UpdateAssetRequest` struct field `Name` to `HostName` with `json:"hostName"` tag
- [x] 3.3 Update `AssetResponse` struct field `Name` to `HostName` with `json:"hostName"` tag
- [x] 3.4 Update `TagInfo` struct if needed (no changes needed - tag name field is separate)
- [x] 3.5 Update `CreateAsset` method to use `req.HostName` instead of `req.Name`
- [x] 3.6 Update `UpdateAsset` method to use `req.HostName` instead of `req.Name`
- [x] 3.7 Update `assetToResponse` method to use `asset.HostName` instead of `asset.Name`
- [x] 3.8 Update `ListAssets` search query to use `host_name` instead of `name` in WHERE clause
- [x] 3.9 Update all references to `asset.Name` in service layer

## Phase 4: Backend Import/Export

- [x] 4.1 Update CSV export header from `"Name"` to `"Hostname"`
- [x] 4.2 Update CSV export to write `asset.HostName` instead of `asset.Name`
- [x] 4.3 Update CSV import expected header from `"Name"` to `"Hostname"`
- [x] 4.4 Update CSV import parsing to read `"Hostname"` column
- [x] 4.5 Update CSV import validation error message (changed to "hostname is required")

## Phase 5: Backend Handler Layer

- [x] 5.1 Review handler code for any direct field access (no direct field access found)
- [x] 5.2 Update handler tests if they reference `name` field directly (no tests found in handler)
- [x] 5.3 Verify handler binding works correctly with `hostName` field (handler uses service layer)

## Phase 6: Frontend API Client

- [x] 6.1 Update `Asset` interface field from `name: string` to `hostName: string`
- [x] 6.2 Update `CreateAssetRequest` interface field from `name: string` to `hostName: string`
- [x] 6.3 Update `UpdateAssetRequest` interface field from `name?: string` to `hostName?: string`
- [x] 6.4 Update `ListAssetsParams` interface if it has search fields (no changes needed - search uses query parameter)

## Phase 7: Frontend Components

- [x] 7.1 Update AssetList component: `asset.name` → `asset.hostName` in all references
- [x] 7.2 Update AssetList form field: `name="name"` → `name="hostName"`
- [x] 7.3 Update AssetDetail component: `asset.name` → `asset.hostName`
- [x] 7.4 Update WebSSHTerminal component: `asset.name` → `asset.hostName`
- [x] 7.5 Update TaskList component: `a.name` → `a.hostName` (in asset selection)
- [x] 7.6 Verify all display components show correct field

## Phase 8: Testing & Validation

- [ ] 8.1 Run database migration on test environment
- [ ] 8.2 Test asset creation with new `hostName` field
- [ ] 8.3 Test asset update with new `hostName` field
- [ ] 8.4 Test asset list/search functionality
- [ ] 8.5 Test CSV import with new `"Hostname"` header
- [ ] 8.6 Test CSV export with new `"Hostname"` header
- [ ] 8.7 Test asset detail display
- [ ] 8.8 Verify asset display in task execution context
- [ ] 8.9 Verify asset display in SSH connection context
- [ ] 8.10 Run all backend tests
- [ ] 8.11 Run all frontend tests
- [ ] 8.12 Test API endpoints with API client tools (Postman/curl)

## Phase 9: Documentation

- [ ] 9.1 Update API documentation with new field name
- [ ] 9.2 Update CSV import/export documentation
- [ ] 9.3 Update migration guide for existing data
- [ ] 9.4 Document breaking change in changelog
