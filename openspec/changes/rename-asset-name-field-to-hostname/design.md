# Design: Rename Asset `name` Field to `hostName`

## Architectural Considerations

### Field Naming Conventions

1. **Database Level**: 
   - Column name: `name` → `host_name` (snake_case, PostgreSQL convention)
   - GORM automatically converts Go struct field `HostName` to `host_name` column

2. **Go Struct Level**:
   - Field name: `Name` → `HostName` (PascalCase, Go convention)
   - JSON tag: `json:"name"` → `json:"hostName"` (camelCase, JSON convention)

3. **TypeScript/Frontend**:
   - Interface field: `name: string` → `hostName: string` (camelCase, TypeScript convention)
   - Usage: `asset.name` → `asset.hostName`

4. **CSV Import/Export**:
   - Column header: `"Name"` → `"Hostname"` (PascalCase for clarity in CSV)

### Database Migration Strategy

The migration is straightforward as it's a simple column rename:

```sql
ALTER TABLE assets RENAME COLUMN name TO host_name;
```

**Considerations**:
- GORM indexes and constraints are preserved automatically
- No data transformation needed (values remain the same)
- Column type and constraints remain unchanged
- Foreign key relationships are unaffected

### API Compatibility

This is a **breaking change** with no backward compatibility. All API clients must be updated simultaneously with the backend deployment.

**Options considered**:
1. **No versioning** (chosen): Simple, clean, requires coordinated deployment
2. **API versioning**: Adds complexity but allows gradual migration
3. **Dual support**: Temporarily support both `name` and `hostName` - too complex

**Rationale**: The system is likely still in active development with controlled client access, making a clean break acceptable.

### Search Functionality

The search functionality in `ListAssets` uses SQL LIKE queries on the `name` column:

```go
query = query.Where("name LIKE ? OR code LIKE ? OR ip LIKE ? OR description LIKE ?", ...)
```

This must be updated to:

```go
query = query.Where("host_name LIKE ? OR code LIKE ? OR ip LIKE ? OR description LIKE ?", ...)
```

### CSV Import/Export Compatibility

**Breaking Change**: Existing CSV files with `"Name"` header will fail.

**Migration Path**:
1. Users must update CSV files to use `"Hostname"` header
2. Alternatively, provide migration script/tool (optional, not in scope)
3. Document the change clearly in user guide

### Deployment Sequence

**Critical**: Backend must be deployed before frontend.

1. **Step 1**: Deploy backend with new field names (API now returns `hostName`)
2. **Step 2**: Run database migration
3. **Step 3**: Deploy frontend (now expects `hostName` in API responses)

**Rollback Strategy**:
- Database: `ALTER TABLE assets RENAME COLUMN host_name TO name;`
- Code: Revert to previous version that uses `name` field

### Testing Strategy

1. **Unit Tests**: Update all tests referencing `Name` field
2. **Integration Tests**: Verify API endpoints return `hostName`
3. **E2E Tests**: Verify UI displays and edits hostname correctly
4. **Migration Tests**: Test migration script on copy of production data

### Impact on Other Systems

**Internal Dependencies**:
- Task execution: Uses asset hostname - must update field access
- SSH connections: Uses asset hostname - must update field access
- Asset filtering/search: All queries updated

**External Dependencies**:
- None identified (but document API change for future integrations)

### Performance Considerations

- Column rename is a metadata operation in PostgreSQL - very fast
- No impact on query performance
- Indexes are automatically preserved
