# Change: Add Cloud Platform Management

## Why

Currently, cloud platforms are stored as simple string values in the `assets.cloud_platform` field. This approach has limitations:

- **No standardization**: Users can enter any string value, leading to inconsistencies (e.g., "aliyun", "Alibaba Cloud", "阿里云")
- **No metadata**: Cannot store additional information about cloud platforms (icons, regions, API endpoints, etc.)
- **No validation**: Typo errors can cause data quality issues
- **No relationship tracking**: Cannot easily track which assets belong to which cloud platforms

Creating a dedicated Cloud Platform management feature (similar to Project and Environment management) will:
- Standardize cloud platform names and identifiers
- Enable data consistency across assets
- Support future features like cloud-specific integrations
- Provide a foundation for cloud platform analytics and reporting
- Improve user experience with dropdown selection instead of free-text input

## What Changes

- **ADDED**: `CloudPlatform` model with fields: `id`, `name`, `code` (optional identifier), `icon` (optional), `description`, `created_at`, `updated_at`, `deleted_at`
- **ADDED**: Backend CRUD API endpoints for cloud platform management (`/cloud-platforms`)
- **ADDED**: Frontend cloud platform management page with list, create, edit, delete functionality
- **ADDED**: Menu item for "云平台管理" (Cloud Platform Management) in the sidebar
- **MODIFIED**: `Asset` model to change `cloud_platform` from string to foreign key relationship (`CloudPlatformID *uint`)
- **MODIFIED**: Asset creation/edit forms to use dropdown selection instead of text input
- **MODIFIED**: Asset list and detail views to display cloud platform name instead of string value
- **MODIFIED**: Asset search/filter to use cloud platform ID instead of string matching

## Impact

- **Database Migration**: 
  - Create `cloud_platforms` table
  - Migrate existing `assets.cloud_platform` string values to `cloud_platforms` table
  - Add `cloud_platform_id` foreign key column to `assets` table
  - Update asset records to reference cloud platform entries
  
- **Breaking Change**: 
  - Asset API field changes from `cloud_platform: string` to `cloud_platform_id: number`
  - Existing API clients need to be updated
  
- **Migration Strategy**:
  1. Create `cloud_platforms` table
  2. Extract unique values from `assets.cloud_platform`
  3. Create cloud platform entries for each unique value
  4. Add `cloud_platform_id` column to `assets` table
  5. Update asset records to set `cloud_platform_id` based on `cloud_platform` string
  6. Drop `cloud_platform` string column
  7. Deploy backend changes
  8. Deploy frontend changes

## Considerations

- **Data Migration**: Need to handle existing data gracefully
  - Extract all unique `cloud_platform` values from existing assets
  - Create corresponding cloud platform entries
  - Map asset records to cloud platform IDs
  
- **Default Cloud Platforms**: Consider pre-populating common platforms (Aliyun, Tencent Cloud, AWS, Azure, etc.)

- **Backward Compatibility**: During migration, maintain support for both string and ID formats temporarily

- **Validation**: Ensure cloud platform names are unique

- **UI Placement**: Add "云平台管理" menu item near "环境管理" for logical grouping

- **Asset Forms**: Update asset creation/edit forms to show cloud platform dropdown instead of text input
