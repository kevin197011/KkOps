# Change: Simplify SSH Management to Use Host Data Directly

## Why
Currently, SSH management maintains a separate `SSHConnection` table that duplicates information already stored in the `Host` table (hostname, IP address, SSH port, etc.). This creates data redundancy and maintenance overhead. Users must manually create SSH connections for each host, even though the host information already exists.

By reusing the `Host` table directly, we can:
- Eliminate data duplication
- Automatically provide WebSSH access for all hosts in the system
- Simplify the user experience (no need to create SSH connections separately)
- Reduce maintenance burden (host updates automatically reflect in SSH management)

## What Changes
- **REMOVED**: SSHConnection model and table
  - Remove `SSHConnection` model from backend
  - Remove `ssh_connections` table from database
  - Remove SSH connection CRUD APIs
  - Remove SSH connection management UI

- **MODIFIED**: SSH Management to use Host data directly
  - SSH management page displays hosts from host management
  - Each host can directly open WebSSH terminal
  - WebSSH terminal connection uses host data (hostname, IP, SSH port)
  - SSH authentication uses SSH keys or prompts for password

- **MODIFIED**: SSHSession model
  - Change `ConnectionID` to `HostID` (reference to Host table)
  - Update session creation to use host ID instead of connection ID

- **KEPT**: SSH Key management
  - SSH key CRUD operations remain unchanged
  - SSH keys can be associated with hosts for authentication

- **KEPT**: SSH Session management
  - Session listing, viewing, and closing remain
  - Sessions now reference hosts instead of connections

- **ADDED**: Host-based SSH authentication configuration
  - Allow configuring default SSH username per host (optional field in Host model)
  - Allow associating SSH key with host (optional field in Host model)
  - Support password authentication (prompt user when connecting)

## Impact
- **Affected specs**: `ssh-management`, `host-management`
- **Affected code**:
  - `backend/internal/models/ssh.go` - Remove SSHConnection, modify SSHSession
  - `backend/internal/models/host.go` - Add optional SSH-related fields
  - `backend/internal/repository/ssh_repository.go` - Remove SSHConnection repository
  - `backend/internal/service/ssh_service.go` - Remove SSHConnection service, modify session service
  - `backend/internal/handler/ssh_handler.go` - Remove connection handlers, modify terminal handler
  - `backend/internal/handler/ssh_terminal_handler.go` - Use host data instead of connection data
  - `frontend/src/services/ssh.ts` - Remove connection service methods
  - `frontend/src/pages/SSH.tsx` - Display hosts instead of connections
  - `frontend/src/components/Terminal.tsx` - Accept host ID instead of connection ID

- **Database**: 
  - Remove `ssh_connections` table
  - Modify `ssh_sessions` table: change `connection_id` to `host_id`
  - Optionally add `ssh_username` and `ssh_key_id` fields to `hosts` table

- **Breaking Changes**: 
  - All existing SSH connections will be lost (migration needed if data preservation is required)
  - API endpoints for SSH connections (`/api/v1/ssh/connections/*`) will be removed
  - Frontend SSH management page will completely change (hosts instead of connections)

- **Migration Strategy**:
  - For existing deployments: Create migration script to convert SSH connections to host SSH configurations
  - For new deployments: No migration needed

