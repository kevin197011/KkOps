# Tasks: Remove Unused `code` Fields

## Phase 1: Database Migrations

- [ ] 1.1 Create migration script to drop `code` column and unique index from `assets` table
- [ ] 1.2 Create migration script to drop `code` column and unique index from `projects` table
- [ ] 1.3 Add unique constraint on `projects.name` if not already present
- [ ] 1.4 Create migration script to drop `code` column and unique index from `asset_categories` table
- [ ] 1.5 Create migration script to drop `code` column and unique index from `environments` table
- [ ] 1.6 Add unique constraint on `environments.name` if not already present
- [ ] 1.7 Create migration script to drop `code` column and unique index from `roles` table
- [ ] 1.8 Add unique constraint on `roles.name` if not already present
- [ ] 1.9 Create migration script to drop `code` column and unique index from `permissions` table
- [ ] 1.10 Add unique constraint on `permissions(resource, action)` combination if not already present
- [ ] 1.11 Verify migration scripts handle indexes correctly
- [ ] 1.12 Document rollback procedure

## Phase 2: Backend Model Changes

### Assets
- [ ] 2.1 Remove `Code` field from Asset model struct
- [ ] 2.2 Verify no other references to `asset.Code` in models

### Projects
- [ ] 2.3 Remove `Code` field from Project model struct
- [ ] 2.4 Verify no other references to `project.Code` in models

### Asset Categories
- [ ] 2.5 Remove `Code` field from AssetCategory model struct
- [ ] 2.6 Verify no other references to `category.Code` in models

### Environments
- [ ] 2.7 Remove `Code` field from Environment model struct
- [ ] 2.8 Verify no other references to `environment.Code` in models

### Roles
- [ ] 2.9 Remove `Code` field from Role model struct
- [ ] 2.10 Verify no other references to `role.Code` in models

### Permissions
- [ ] 2.11 Remove `Code` field from Permission model struct
- [ ] 2.12 Verify no other references to `permission.Code` in models

## Phase 3: Backend Service Layer - Assets

- [ ] 3.1 Remove `Code` field from `CreateAssetRequest` struct
- [ ] 3.2 Remove `Code` field from `UpdateAssetRequest` struct
- [ ] 3.3 Remove `Code` field from `AssetResponse` struct
- [ ] 3.4 Update `CreateAsset` method to remove code handling
- [ ] 3.5 Update `UpdateAsset` method to remove code handling
- [ ] 3.6 Update `assetToResponse` method to remove code field
- [ ] 3.7 Update `ListAssets` search query to remove `code` from WHERE clause
- [ ] 3.8 Remove code-related validation logic

## Phase 4: Backend Service Layer - Projects

- [ ] 4.1 Remove `Code` field from `CreateProjectRequest` struct
- [ ] 4.2 Remove `Code` field from `ProjectResponse` struct
- [ ] 4.3 Update `CreateProject` method to remove code handling
- [ ] 4.4 Update uniqueness check to use `name` instead of `code`
- [ ] 4.5 Update `projectToResponse` method to remove code field
- [ ] 4.6 Remove code-related validation logic

## Phase 5: Backend Service Layer - Asset Categories

- [ ] 5.1 Remove `Code` field from category request/response structs
- [ ] 5.2 Update category service methods to remove code handling
- [ ] 5.3 Remove code uniqueness checks

## Phase 5a: Backend Service Layer - Environments

- [ ] 5a.1 Remove `Code` field from `CreateEnvironmentRequest` struct
- [ ] 5a.2 Remove `Code` field from `EnvironmentResponse` struct
- [ ] 5a.3 Update `CreateEnvironment` method to remove code handling
- [ ] 5a.4 Update uniqueness check to use `name` instead of `code`
- [ ] 5a.5 Update `environmentToResponse` method to remove code field
- [ ] 5a.6 Remove code-related validation logic

## Phase 5b: Backend Service Layer - Roles

- [ ] 5b.1 Remove `Code` field from role request/response structs
- [ ] 5b.2 Update role service methods to remove code handling
- [ ] 5b.3 Update uniqueness check to use `name` instead of `code`
- [ ] 5b.4 Remove code-related validation logic

## Phase 5c: Backend Service Layer - Permissions & RBAC

- [ ] 5c.1 Remove `Code` field from permission request/response structs (if any)
- [ ] 5c.2 Update permission service methods to remove code handling
- [ ] 5c.3 Remove `HasPermissionByCode` method (or mark as deprecated if used elsewhere)
- [ ] 5c.4 Remove `RequirePermissionByCode` middleware (or mark as deprecated if used elsewhere)
- [ ] 5c.5 Find and update all permission code references in handlers/middleware
- [ ] 5c.6 Update permission response structures to remove code field
- [ ] 5c.7 Verify all permission checks use `HasPermission(resource, action)` instead

## Phase 5d: Backend Service Layer - Auth Service

- [ ] 5d.1 Update auth service JWT token generation to use role names instead of role codes
- [ ] 5d.2 Update `Login` method to use role names in JWT claims
- [ ] 5d.3 Update `GetCurrentUser` method to use role names instead of role codes
- [ ] 5d.4 Update role-related response structures to remove code field

## Phase 6: Backend Import/Export - Assets

- [ ] 6.1 Remove "Code" from CSV export header
- [ ] 6.2 Remove code from CSV export data rows
- [ ] 6.3 Remove "Code" from CSV import expected headers
- [ ] 6.4 Remove code parsing from CSV import logic
- [ ] 6.5 Update import validation messages

## Phase 7: Frontend API Client

### Assets
- [ ] 7.1 Remove `code` field from `Asset` interface
- [ ] 7.2 Remove `code` field from `CreateAssetRequest` interface
- [ ] 7.3 Remove `code` field from `UpdateAssetRequest` interface

### Projects
- [ ] 7.4 Remove `code` field from `Project` interface (if exists)
- [ ] 7.5 Remove `code` field from `CreateProjectRequest` interface (if exists)

### Categories
- [ ] 7.6 Remove `code` field from category interfaces (if exists)

### Environments
- [ ] 7.7 Remove `code` field from `Environment` interface
- [ ] 7.8 Remove `code` field from `CreateEnvironmentRequest` interface

### Roles
- [ ] 7.9 Remove `code` field from `Role` interface (if exists)
- [ ] 7.10 Remove `code` field from role request interfaces (if exists)

### Permissions
- [ ] 7.11 Remove `code` field from `Permission` interface (if exists)

## Phase 8: Frontend Components - Assets

- [ ] 8.1 Remove code input field from AssetList create/edit form
- [ ] 8.2 Remove code column from AssetList table
- [ ] 8.3 Remove code display from AssetDetail page/modal
- [ ] 8.4 Update search placeholder text (remove "代码")
- [ ] 8.5 Update form field references in handleEdit method

## Phase 9: Frontend Components - Projects

- [ ] 9.1 Remove code input field from ProjectList create/edit form
- [ ] 9.2 Remove code column from ProjectList table
- [ ] 9.3 Update form field references

## Phase 10: Frontend Components - Categories

- [ ] 10.1 Remove code input field from CategoryList form (if exists)
- [ ] 10.2 Remove code column from CategoryList table (if exists)

## Phase 10a: Frontend Components - Environments

- [ ] 10a.1 Remove code input field from EnvironmentList create/edit form
- [ ] 10a.2 Remove code column from EnvironmentList table
- [ ] 10a.3 Update form field references in handleEdit method

## Phase 10b: Frontend Components - Roles

- [ ] 10b.1 Remove code input field from RoleList create/edit form
- [ ] 10b.2 Remove code column from RoleList table
- [ ] 10b.3 Update form field references in handleEdit method

## Phase 11: Testing & Validation

- [ ] 11.1 Run database migrations on test environment
- [ ] 11.2 Test asset creation without code field
- [ ] 11.3 Test asset update without code field
- [ ] 11.4 Test asset search (verify code is not in search)
- [ ] 11.5 Test project creation without code field
- [ ] 11.6 Test project name uniqueness validation
- [ ] 11.7 Test environment creation without code field
- [ ] 11.8 Test environment name uniqueness validation
- [ ] 11.9 Test role creation without code field
- [ ] 11.10 Test role name uniqueness validation
- [ ] 11.11 Test permission resource/action uniqueness validation
- [ ] 11.12 Test RBAC permission checks with resource/action instead of code
- [ ] 11.13 Test CSV import without code column
- [ ] 11.14 Test CSV export without code column
- [ ] 11.15 Verify all UI forms work correctly
- [ ] 11.16 Run all backend tests
- [ ] 11.17 Run all frontend tests

## Phase 12: Documentation

- [ ] 12.1 Update API documentation to remove code fields
- [ ] 12.2 Update CSV import/export documentation
- [ ] 12.3 Document breaking change in changelog
