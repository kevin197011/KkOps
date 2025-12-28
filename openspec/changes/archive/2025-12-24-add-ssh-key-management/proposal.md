# Change: Add SSH Key Management Integrated with WebSSH

## Why
Currently, WebSSH only supports password-based authentication, requiring users to enter credentials each time they connect. This is inconvenient and less secure than key-based authentication. Users need the ability to:
- Upload and manage SSH keys for secure authentication
- Use SSH keys when connecting via WebSSH terminals
- Optionally configure default SSH keys per host for automatic authentication

By adding SSH key management integrated with WebSSH, we can:
- Provide more secure authentication (key-based auth)
- Improve user experience (no need to enter password each time)
- Allow per-host key configuration for automation
- Maintain the simplified architecture (keys managed within WebSSH context)

## What Changes
- **ADDED**: SSH Key Management in WebSSH page
  - Add "SSH密钥" tab to WebSSH management page
  - Support uploading SSH private keys (RSA, ED25519, etc.)
  - Support viewing, editing, and deleting SSH keys
  - Encrypt and store SSH keys securely in database
  - Display key metadata (name, type, fingerprint, created date)

- **ADDED**: SSH Key selection in WebSSH terminal connection
  - When opening a terminal, allow user to choose authentication method (key or password)
  - If key is selected, show list of available SSH keys
  - Support key-based SSH authentication in WebSSH handler
  - Fallback to password authentication if key authentication fails

- **ADDED**: Optional per-host SSH key configuration
  - Add optional `ssh_key_id` field to Host model
  - Allow configuring default SSH key for a host in host management
  - Auto-select configured key when opening terminal for that host
  - Support overriding default key selection

- **MODIFIED**: WebSSH authentication flow
  - Support both key-based and password-based authentication
  - Update WebSSH handler to handle SSH key authentication
  - Update Terminal component to support key selection UI

## Impact
- **Affected specs**: `ssh-key-management` (new), `webssh-management` (modified)
- **Affected code**:
  - `backend/internal/models/ssh_key.go` - New SSH key model
  - `backend/internal/repository/ssh_key_repository.go` - New repository
  - `backend/internal/service/ssh_key_service.go` - New service
  - `backend/internal/handler/ssh_key_handler.go` - New handler
  - `backend/internal/handler/webssh_handler.go` - Modify to support key auth
  - `backend/internal/models/host.go` - Add optional `ssh_key_id` field
  - `frontend/src/pages/WebSSH.tsx` - Add SSH key management tab
  - `frontend/src/components/Terminal.tsx` - Add key selection UI
  - `frontend/src/services/sshKey.ts` - New service for SSH key API
- **Database**: New `ssh_keys` table, optional `ssh_key_id` column in `hosts` table
- **User experience**: Users can manage SSH keys and use them for secure terminal access

