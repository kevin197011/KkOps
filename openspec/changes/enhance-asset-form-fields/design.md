# Design: Enhance Asset Form Fields

## Form Structure

The enhanced asset form should display all available fields organized into logical sections:

### Section 1: Basic Information (Required)
- **hostName** (Required) - Asset hostname

### Section 2: Organization (Optional)
- **project_id** - Project association (dropdown)
- **cloud_platform_id** - Cloud platform (dropdown)
- **environment_id** - Environment (dropdown)

### Section 3: Network Information (Optional)
- **ip** - IP address
- **status** - Asset status (dropdown: active/disabled, default: active)

### Section 4: SSH Connection (Optional, can be auto-collected)
- **ssh_port** - SSH port (number input, default: 22)
- **ssh_key_id** - SSH key (dropdown)
- **ssh_user** - SSH username

### Section 5: Hardware Specifications (Optional, can be auto-collected)
- **cpu** - CPU information
- **memory** - Memory information
- **disk** - Disk information

### Section 6: Additional Information (Optional)
- **description** - Description (textarea)
- **tag_ids** - Tags (multi-select)

## UI Layout Options

### Option 1: Single Column with Sections
- All fields in a single scrollable column
- Sections with headers/separators
- Simple and straightforward
- Works well on all screen sizes

### Option 2: Two-Column Layout
- Basic info and organization on left
- Network, SSH, hardware on right
- More compact
- Better use of horizontal space
- May need responsive handling for mobile

### Recommended: Option 1 with Collapsible Sections
- Single column layout for consistency
- Collapsible sections for optional fields
- "Basic Information" always expanded
- Other sections can be collapsed by default
- Better UX for users who want to enter minimal info

## Field Requirements

### Required Fields
- `hostName` - Must be provided, cannot be empty

### Optional Fields (All others)
- Can be left empty during creation
- Can be filled later through editing
- Some fields (CPU, memory, disk, SSH info) can be auto-collected via agents/scripts
- Validation should only enforce required fields

## Visual Indicators

- **Required fields**: Mark with asterisk (*) and "required" validation
- **Optional fields**: No special marking, but can add tooltip/hint
- **Auto-collectable fields**: Add hint/tooltip indicating they can be auto-collected
- **Default values**: Show in placeholder or as pre-filled values

## Backend Considerations

- No changes to API structure needed
- All fields already supported in CreateAssetRequest and UpdateAssetRequest
- Validation logic remains the same (only hostName required)
- Default values handled on backend (ssh_port: 22, status: active)

## Frontend Implementation

### Form Component Structure
```tsx
<Form>
  {/* Section 1: Basic Information */}
  <Form.Item name="hostName" required>
    <Input />
  </Form.Item>

  {/* Section 2: Organization */}
  <Form.Item name="project_id">
    <Select />
  </Form.Item>
  {/* ... other organization fields */}

  {/* Section 3: Network Information */}
  {/* ... network fields */}

  {/* Section 4: SSH Connection */}
  {/* ... SSH fields with hints about auto-collection */}

  {/* Section 5: Hardware Specifications */}
  {/* ... hardware fields with hints about auto-collection */}

  {/* Section 6: Additional Information */}
  {/* ... description, tags */}
</Form>
```

### State Management
- Form state managed by Ant Design Form
- Fetch related data (projects, cloud platforms, environments, SSH keys, tags) on component mount
- Handle form submission with all field values

### Validation
- Only hostName is required
- All other fields are optional
- Field-level validation for format (e.g., IP address format, port range)

## Future Auto-Collection Support

The form design should accommodate future automatic collection features:
- Fields that can be auto-collected should be clearly marked
- UI should indicate when fields are auto-collected vs manually entered
- Consider adding a "Refresh from Agent" button in edit mode
- Fields may have read-only indicators when auto-managed
