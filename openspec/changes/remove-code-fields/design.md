# Design: Remove Unused `code` Fields

## Scope

This change removes the `code` field from all models that currently have it:
1. **Asset** - Replaced by `hostName` for identification
2. **Project** - Replaced by `name` for identification (with unique constraint)
3. **AssetCategory** - Obsolete model (categories removed from asset workflow)
4. **Environment** - Replaced by `name` for identification (with unique constraint)
5. **Role** - Replaced by `name` for identification (with unique constraint)
6. **Permission** - Replaced by `resource` + `action` combination for identification

## Uniqueness Strategy

### Current State
- Assets: `code` has unique index, `hostName` may not have unique constraint
- Projects: `code` has unique index, `name` may not have unique constraint
- AssetCategory: `code` has unique index

### After Removal

1. **Assets**:
   - Remove unique index on `code` column
   - Verify/add unique constraint on `host_name` if not already present
   - Use `hostName` as the business identifier

2. **Projects**:
   - Remove unique index on `code` column
   - Add unique constraint on `name` column if not already present
   - Use `name` as the business identifier
   - Update validation logic to check for duplicate names

3. **AssetCategory**:
   - Remove unique index on `code` column
   - No uniqueness constraints needed (feature is deprecated)

4. **Environment**:
   - Remove unique index on `code` column
   - Add unique constraint on `name` column if not already present
   - Use `name` as the business identifier
   - Update validation logic to check for duplicate names

5. **Role**:
   - Remove unique index on `code` column
   - Add unique constraint on `name` column if not already present
   - Use `name` as the business identifier
   - Update validation logic to check for duplicate names
   - Update RBAC permission checks to use resource/action instead of role code

6. **Permission**:
   - Remove unique index on `code` column
   - Add unique constraint on `resource` + `action` combination if not already present
   - Use `resource` + `action` as the business identifier
   - Update RBAC permission checks to use resource/action instead of permission code

## Database Migration

Three separate migrations (or one combined):

```sql
-- Migration 1: Assets
DROP INDEX IF EXISTS idx_assets_code;
ALTER TABLE assets DROP COLUMN IF EXISTS code;

-- Migration 2: Projects
DROP INDEX IF EXISTS idx_projects_code;
ALTER TABLE projects DROP COLUMN IF EXISTS code;
-- Add unique constraint on name if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_name ON projects(name) WHERE deleted_at IS NULL;

-- Migration 3: Asset Categories
DROP INDEX IF EXISTS idx_asset_categories_code;
ALTER TABLE asset_categories DROP COLUMN IF EXISTS code;

-- Migration 4: Environments
DROP INDEX IF EXISTS idx_environments_code;
ALTER TABLE environments DROP COLUMN IF EXISTS code;
-- Add unique constraint on name if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_environments_name ON environments(name) WHERE deleted_at IS NULL;

-- Migration 5: Roles
DROP INDEX IF EXISTS idx_roles_code;
ALTER TABLE roles DROP COLUMN IF EXISTS code;
-- Add unique constraint on name if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON roles(name) WHERE deleted_at IS NULL;

-- Migration 6: Permissions
DROP INDEX IF EXISTS idx_permissions_code;
ALTER TABLE permissions DROP COLUMN IF EXISTS code;
-- Add unique constraint on resource + action if needed
CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action) WHERE deleted_at IS NULL;
```

## Search Functionality Impact

### Assets
Current search includes code:
```sql
WHERE host_name LIKE ? OR code LIKE ? OR ip LIKE ? OR description LIKE ?
```

After removal:
```sql
WHERE host_name LIKE ? OR ip LIKE ? OR description LIKE ?
```

### Projects
Search should use `name` field (if search exists).

## API Changes

### Asset API
- **REMOVED**: `code` field from `CreateAssetRequest`
- **REMOVED**: `code` field from `UpdateAssetRequest`
- **REMOVED**: `code` field from `AssetResponse`

### Project API
- **REMOVED**: `code` field from `CreateProjectRequest`
- **REMOVED**: `code` field from `ProjectResponse`
- **MODIFIED**: Uniqueness validation to check `name` instead of `code`

### AssetCategory API
- **REMOVED**: `code` field from all request/response structs

### Environment API
- **REMOVED**: `code` field from `CreateEnvironmentRequest`
- **REMOVED**: `code` field from `EnvironmentResponse`
- **MODIFIED**: Uniqueness validation to check `name` instead of `code`

### Role API
- **REMOVED**: `code` field from role request/response structs
- **MODIFIED**: Uniqueness validation to check `name` instead of `code`

### Permission API
- **REMOVED**: `code` field from permission request/response structs
- **MODIFIED**: Uniqueness validation to check `resource` + `action` combination instead of `code`

## Frontend Changes

### Asset Management
- Remove code input field from create/edit forms
- Remove code column from asset list table
- Remove code display from asset detail view
- Update search placeholder text

### Project Management
- Remove code input field from create/edit forms
- Remove code column from project list table (if displayed)
- Update uniqueness validation messaging

### Category Management
- Remove code input field from forms
- Remove code column from category list (if displayed)

## CSV Import/Export

### Assets
- **REMOVED**: "Code" column from CSV export
- **REMOVED**: "Code" column requirement from CSV import
- Update import validation to require only "Hostname" and other essential fields

### Projects
- **REMOVED**: "Code" column from CSV export (if exists)
- **REMOVED**: "Code" column requirement from CSV import (if exists)

### Environments
- No CSV import/export currently implemented (no changes needed)

### RBAC Permission Checking Changes

**Current Implementation**:
- `HasPermissionByCode(userID, "user:create")` - checks permission by code
- `RequirePermissionByCode(rbacService, "user:create")` - middleware uses code

**New Implementation**:
- `HasPermission(userID, "user", "create")` - already exists, checks by resource/action
- `RequirePermission(rbacService, "user", "create")` - already exists middleware

This change requires:
- Remove or deprecate `HasPermissionByCode` method
- Remove or deprecate `RequirePermissionByCode` middleware  
- Update all code references to use resource/action-based methods
- Remove permission code from role/permission response structures

### JWT Token Changes

**Current Implementation**:
JWT tokens include role codes in the `roles` claim:
```go
roleCodes := []string{role.Code, ...}
GenerateJWT(userID, username, roleCodes, ...)
```

**New Implementation**:
JWT tokens include role names in the `roles` claim:
```go
roleNames := []string{role.Name, ...}
GenerateJWT(userID, username, roleNames, ...)
```

This change requires:
- Update `Login` method in auth service to use `role.Name` instead of `role.Code`
- Update `GetCurrentUser` method to return role names instead of role codes
- JWT token structure remains the same (string array), only the content changes from codes to names
- Frontend receives role names in JWT token (no breaking change for frontend, just semantic change)

## Rollback Strategy

If rollback is needed:
```sql
-- Restore columns (data will be lost)
ALTER TABLE assets ADD COLUMN code VARCHAR(50);
ALTER TABLE projects ADD COLUMN code VARCHAR(50);
ALTER TABLE asset_categories ADD COLUMN code VARCHAR(50);
ALTER TABLE environments ADD COLUMN code VARCHAR(50);
ALTER TABLE roles ADD COLUMN code VARCHAR(50);
ALTER TABLE permissions ADD COLUMN code VARCHAR(100);

-- Restore indexes
CREATE UNIQUE INDEX idx_assets_code ON assets(code);
CREATE UNIQUE INDEX idx_projects_code ON projects(code);
CREATE UNIQUE INDEX idx_asset_categories_code ON asset_categories(code);
CREATE UNIQUE INDEX idx_environments_code ON environments(code);
CREATE UNIQUE INDEX idx_roles_code ON roles(code);
CREATE UNIQUE INDEX idx_permissions_code ON permissions(code);

-- Drop unique constraints added during migration
DROP INDEX IF EXISTS idx_projects_name;
DROP INDEX IF EXISTS idx_environments_name;
DROP INDEX IF EXISTS idx_roles_name;
DROP INDEX IF EXISTS idx_permissions_resource_action;
```

Then revert code changes to previous version.
