# Tasks: Enhance Asset Form to Display All Fields

## Phase 1: Frontend Form Enhancement

- [x] 1.1 Update `frontend/src/pages/assets/AssetList.tsx` asset creation/edit form:
  - Add SSH connection fields section: `ssh_port`, `ssh_key_id`, `ssh_user`
  - Add hardware specifications section: `cpu`, `memory`, `disk`
  - Add tags field (`tag_ids`) with multi-select
  - Organize fields into logical sections
  - Ensure all fields from Asset model are displayed

- [x] 1.2 Fetch SSH keys list for `ssh_key_id` dropdown:
  - Import SSH key API
  - Add state for SSH keys
  - Fetch SSH keys on component mount
  - Populate dropdown options

- [x] 1.3 Fetch tags list for `tag_ids` multi-select:
  - Import tag API (if not already imported)
  - Add state for tags (if not already present)
  - Fetch tags on component mount
  - Implement multi-select for tags

- [x] 1.4 Update form layout and styling:
  - Organize fields into sections with clear grouping
  - Consider using Ant Design Form sections or dividers
  - Ensure form is scrollable if needed
  - Maintain responsive design

- [x] 1.5 Update form validation:
  - Ensure only `hostName` is required
  - All other fields remain optional
  - Add appropriate field types (number for ssh_port, etc.)

- [x] 1.6 Update `handleEdit` function:
  - Include all new fields in form initialization
  - Set `ssh_port`, `ssh_key_id`, `ssh_user`, `cpu`, `memory`, `disk`, `tag_ids` values

- [x] 1.7 Update `handleSubmit` function:
  - Ensure all field values are submitted
  - Handle `tag_ids` array properly
  - Handle optional fields (empty strings vs undefined)

## Phase 2: Testing & Validation

- [ ] 2.1 Test asset creation with all fields:
  - Create asset with only required field (hostName)
  - Create asset with all fields filled
  - Create asset with partial fields
  - Verify all fields are saved correctly

- [ ] 2.2 Test asset editing with all fields:
  - Edit asset and update various fields
  - Verify field values are pre-populated correctly
  - Verify updates are saved correctly

- [ ] 2.3 Verify form validation:
  - Confirm hostName is required
  - Confirm other fields are optional
  - Test form submission with various field combinations

- [ ] 2.4 Verify UI/UX:
  - Check form layout and organization
  - Verify all fields are accessible
  - Check responsive design on different screen sizes
  - Verify dropdown options load correctly

## Phase 3: Documentation & Cleanup

- [x] 3.1 Verify all fields from Asset model are included
- [x] 3.2 Review form organization and grouping
- [x] 3.3 Update OpenSpec tasks status
- [ ] 3.4 Run OpenSpec validation: `openspec validate enhance-asset-form-fields --strict`
