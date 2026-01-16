# Change: Remove Unused `code` Fields from All Models

## Why

The `code` fields across multiple models were intended as unique business identifiers but serve no actual purpose in the current system:

- **Asset**: Has `hostName` which serves as the unique identifier. The `code` field is redundant and adds unnecessary complexity.
- **Project**: Has `name` which can serve as the unique identifier. The `code` field adds no value.
- **AssetCategory**: The asset category feature has already been removed from the asset management workflow, making this field obsolete.
- **Environment**: Has `name` which can serve as the unique identifier. The `code` field adds no value.
- **Role**: Has `name` which can serve as the unique identifier. The `code` field is only used internally for permission checks, but can be replaced with name-based checks.
- **Permission**: Has `resource` and `action` combination which serves as the unique identifier. The `code` field is only used internally, but can be replaced with resource/action-based checks.

Removing these fields will:
- Simplify the data model and reduce maintenance overhead
- Reduce API payload size
- Eliminate unnecessary validation and uniqueness checks
- Simplify user interfaces by removing redundant fields
- Reduce confusion for users who must provide a "code" without clear purpose
- Standardize identification approach across all models

## What Changes

- **REMOVED**: `code` field from Asset model, API requests/responses, and database schema
- **REMOVED**: `code` field from Project model, API requests/responses, and database schema
- **REMOVED**: `code` field from AssetCategory model, API requests/responses, and database schema
- **REMOVED**: `code` field from Environment model, API requests/responses, and database schema
- **REMOVED**: `code` field from Role model, API requests/responses, and database schema
- **REMOVED**: `code` field from Permission model, API requests/responses, and database schema
- **MODIFIED**: Asset search functionality to remove code from searchable fields
- **MODIFIED**: CSV import/export to remove code column
- **MODIFIED**: All frontend forms and displays to remove code fields
- **MODIFIED**: Uniqueness validation logic to use name or hostName instead of code
- **MODIFIED**: RBAC permission checks to use resource/action instead of permission code

## Impact

- **Breaking Change**: This is a breaking change affecting:
  - Database schema (requires migration to drop columns)
  - All API endpoints for assets, projects, and categories
  - Frontend forms and displays
  - CSV import/export format

- **Database**: Migration required to drop `code` columns and their unique indexes from:
  - `assets` table
  - `projects` table
  - `asset_categories` table
  - `environments` table
  - `roles` table
  - `permissions` table

- **Backend**: 
  - All service, handler, and model code updated
  - Uniqueness checks moved from code to name/hostName
  - Search queries updated

- **Frontend**: 
  - All forms updated to remove code fields
  - All list/detail views updated
  - Search functionality updated

- **Import/Export**: CSV format changed - existing CSV files will need to be updated

## Migration Strategy

1. **Database Migration**:
   - Drop unique index on `code` columns
   - Drop `code` columns from `assets`, `projects`, and `asset_categories` tables

2. **Backend Deployment**:
   - Deploy backend changes (code updates)
   - Run database migration
   - Backend will no longer accept or return `code` fields

3. **Frontend Deployment**:
   - Deploy frontend changes after backend is deployed
   - Frontend will no longer display or require code fields

## Considerations

- Existing data with code values will be lost (but this is acceptable as the field has no meaning)
- Users will need to update CSV import files to remove code columns
- API documentation must be updated
- Uniqueness validation will need to be updated:
  - Assets: Use `hostName` for uniqueness (already unique)
  - Projects: Use `name` for uniqueness (may need to add unique constraint)
  - AssetCategory: No uniqueness needed (category feature is deprecated)
  - Environments: Use `name` for uniqueness (may need to add unique constraint)
  - Roles: Use `name` for uniqueness (may need to add unique constraint)
  - Permissions: Use `resource` + `action` combination for uniqueness (may need to add unique constraint)

- RBAC permission checking logic needs to be updated:
  - Replace permission code-based checks with resource/action-based checks
  - Update any hardcoded permission code references