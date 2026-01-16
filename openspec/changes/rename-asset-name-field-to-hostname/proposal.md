# Change: Rename Asset `name` Field to `hostName` in Database and API

## Why

Previously, we changed the UI label from "名称" to "主机名" for better semantic clarity, but the underlying field name remained `name`. To achieve complete consistency and better align with domain terminology, we need to rename the database column and API field name from `name` to `hostName`. This change:
- Aligns field names with the UI terminology (主机名 = hostname)
- Improves API clarity - `hostName` is more descriptive than `name`
- Follows common IT infrastructure naming conventions
- Makes the codebase more self-documenting

## What Changes

- **BREAKING CHANGE**: Database column `name` renamed to `host_name` in `assets` table
- **BREAKING CHANGE**: Backend Go struct field `Name` renamed to `HostName` in Asset model
- **BREAKING CHANGE**: JSON API field name changed from `name` to `hostName` in all asset-related endpoints
- **BREAKING CHANGE**: Frontend TypeScript interface `name` field renamed to `hostName`
- **BREAKING CHANGE**: CSV import/export column header changed from `"Name"` to `"Hostname"`
- **MODIFIED**: All backend code referencing `asset.Name` updated to `asset.HostName`
- **MODIFIED**: All frontend code referencing `asset.name` updated to `asset.hostName`
- **MODIFIED**: Database migration script to rename column
- **MODIFIED**: Search queries updated to use `host_name` instead of `name`

## Impact

- **Breaking Change**: This is a major breaking change affecting:
  - Database schema (requires migration)
  - All API endpoints returning or accepting asset data
  - All frontend components displaying or editing asset information
  - CSV import/export format
  - Any external integrations using the asset API

- **Database**: Migration required to rename `name` column to `host_name`
- **Backend**: All service, handler, and model code updated
- **Frontend**: All TypeScript interfaces and component code updated
- **API Contract**: JSON field names changed - existing API clients will break
- **Import/Export**: CSV format changed - existing CSV files will need to be updated

## Migration Strategy

1. **Database Migration**: 
   - Create migration script: `ALTER TABLE assets RENAME COLUMN name TO host_name;`
   - Verify all constraints and indexes are preserved

2. **Backend Deployment**:
   - Deploy backend changes (code updates)
   - Run database migration
   - Backend will now use `hostName` in JSON responses

3. **Frontend Deployment**:
   - Deploy frontend changes after backend is deployed
   - Frontend will now expect `hostName` in API responses

4. **Backward Compatibility**:
   - Consider API versioning if supporting old clients is required
   - Document the breaking change clearly

## Considerations

- All existing data must be migrated - the `name` column contains actual hostname values
- Existing CSV import files using "Name" header will fail - users need to update them
- API documentation must be updated
- Integration tests must be updated
- Consider deprecation period if external integrations exist
