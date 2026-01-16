# Change: Rename Asset "名称" to "主机名"

## Why

The current UI label "名称" (Name) for assets is too generic and doesn't clearly indicate what this field represents in the context of IT asset management. Changing it to "主机名" (Hostname) provides better semantic clarity and aligns with common IT infrastructure terminology. Hostname is more specific and immediately conveys that this field represents the hostname of the server or machine.

## What Changes

- **MODIFIED**: All UI labels displaying "名称" for assets changed to "主机名"
- **MODIFIED**: Form labels, table column headers, and validation messages updated
- **MODIFIED**: Search placeholder text updated to reflect the new terminology
- **MODIFIED**: Asset detail view labels updated
- **MODIFIED**: Asset management specification updated to use "hostname" terminology
- **NOTE**: Database column name (`name`) and API field name (`name`) remain unchanged - only display labels are updated

## Impact

- **UI/UX**: Improved semantic clarity for users
- **Frontend**: All asset-related UI components updated with new labels
- **Backend**: No breaking changes - field names remain the same
- **Database**: No migration required - column name unchanged
- **API**: No breaking changes - JSON field names remain `name`
- **Documentation**: Specification updated to reflect terminology change
- **Affected Specs**: `asset-management` specification updated to use hostname terminology

## Considerations

- Import/export CSV headers may need to be updated if they display column names in Chinese
- Error messages in backend may need localization updates (currently using English "name is required")
- Other parts of the system that display asset names (e.g., in task execution, SSH connection) continue to work as the underlying field name remains `name`
