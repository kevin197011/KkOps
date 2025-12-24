# Change: Add SSH Key Selection in SSH Connection Form

## Why
Currently, SSH connections can be configured with key-based authentication (`auth_type: 'key'`), but the frontend form does not provide a UI to select which SSH key to use. Users need to be able to:
1. Configure SSH keys in the SSH Keys tab
2. Select a specific SSH key when creating/editing SSH connections with key authentication
3. View which key is associated with each SSH connection

## What Changes
- **Frontend**: Add SSH key selection dropdown in SSH connection form when `auth_type` is "key"
- **Frontend**: Display associated key information in SSH connection list
- **Frontend**: Load SSH keys list when opening connection form
- **Backend**: Ensure `key_id` field is properly validated and saved
- **UI/UX**: Show/hide key selection based on authentication type selection

## Impact
- **Affected specs**: `ssh-management`
- **Affected code**: 
  - `frontend/src/pages/SSH.tsx` - Connection form and list display
  - `frontend/src/services/ssh.ts` - API interface (already has `key_id` field)
  - `backend/internal/handler/ssh_handler.go` - Validation (if needed)
  - `backend/internal/service/ssh_service.go` - Business logic (if needed)

