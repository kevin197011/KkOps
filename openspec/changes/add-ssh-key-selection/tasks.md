## 1. Frontend Implementation
- [x] 1.1 Load SSH keys list when opening SSH connection form
- [x] 1.2 Add SSH key selection dropdown in connection form (show when `auth_type === 'key'`)
- [x] 1.3 Add validation: require `key_id` when `auth_type` is "key"
- [x] 1.4 Update connection form to handle key selection state changes
- [x] 1.5 Display associated key name/info in SSH connection list table
- [x] 1.6 Update connection edit form to pre-populate selected key

## 2. Backend Validation (if needed)
- [x] 2.1 Verify `key_id` validation in SSH connection service
- [x] 2.2 Ensure `key_id` is properly saved when creating/updating connections
- [x] 2.3 Verify key existence check when `key_id` is provided

## 3. Testing
- [x] 3.1 Test creating SSH connection with key authentication
- [x] 3.2 Test editing SSH connection to change/select key
- [x] 3.3 Test validation: key required when auth_type is "key"
- [x] 3.4 Test switching between password and key authentication
- [x] 3.5 Verify key information displays correctly in connection list

