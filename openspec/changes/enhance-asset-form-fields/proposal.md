# Change: Enhance Asset Form to Display All Fields

## Why

Currently, the asset creation and edit form only displays a subset of available asset fields. Many important fields such as SSH connection details (SSH port, SSH key, SSH user) and hardware specifications (CPU, memory, disk) are not visible in the form, even though they exist in the data model and are displayed in the asset detail view.

This limitation causes:
- **Incomplete data entry**: Users cannot easily enter all asset information during creation or editing
- **Poor user experience**: Users must rely on other methods (import, API) or edit after creation to add missing information
- **Inconsistent workflow**: Detail view shows all fields, but form view doesn't allow editing them
- **Reduced functionality**: Fields that could be manually entered must be left empty or updated separately

Enhancing the form to display all fields will:
- Provide a complete and consistent user experience
- Allow users to enter all asset information in one place
- Support both manual entry and future automatic collection
- Improve data completeness and accuracy
- Align form capabilities with data model capabilities

## What Changes

- **MODIFIED**: Asset creation/edit form to display all available asset fields
- **ADDED**: Field grouping and organization (basic info, SSH connection, hardware specs, etc.)
- **ADDED**: Required field indicators and validation
- **ADDED**: Optional field indicators (fields that can be filled later or auto-collected)
- **MODIFIED**: Form layout to accommodate all fields (multi-column, sections, etc.)
- **MAINTAINED**: Backward compatibility - all fields remain optional except hostName

## Impact

- **UI Changes**: Form will be more comprehensive and potentially longer
- **User Experience**: Users can enter complete asset information in one place
- **Data Quality**: Improved data completeness as all fields are accessible
- **Future Compatibility**: Foundation for automatic asset information collection features

## Field Categories

### Required Fields
- `hostName` (already required)

### Optional but Important Fields (can be filled later or auto-collected)
- `project_id` - Project association
- `cloud_platform_id` - Cloud platform
- `environment_id` - Environment
- `ip` - IP address
- `ssh_port` - SSH port (default: 22)
- `ssh_key_id` - SSH key for authentication
- `ssh_user` - SSH username
- `cpu` - CPU information
- `memory` - Memory information
- `disk` - Disk information
- `status` - Asset status (default: active)
- `description` - Description
- `tag_ids` - Tags for categorization

## Considerations

- **Form Layout**: Use sections or columns to organize fields logically
- **Field Visibility**: All fields should be visible, but some may be in collapsible sections
- **Validation**: Only hostName should be required; all other fields remain optional
- **Future Auto-Collection**: Fields like CPU, memory, disk can be auto-collected later via SSH/agent
- **User Guidance**: Provide hints/tooltips for fields that can be auto-collected
- **Default Values**: Maintain sensible defaults (e.g., ssh_port: 22, status: active)
