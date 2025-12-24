# Change: Remove SSH Management and Add WebSSH Management

## Why
The current SSH management system is overly complex, maintaining separate models for SSH connections, SSH keys, and SSH sessions. This creates:
- Data duplication (SSH connection info duplicates host data)
- Maintenance overhead (users must manually create SSH connections for each host)
- Unnecessary complexity (SSH key management, connection management, session management as separate concerns)

By removing SSH management and replacing it with a simplified WebSSH management system that directly uses host data, we can:
- Eliminate data redundancy
- Simplify the user experience (direct terminal access from host management)
- Reduce maintenance burden (host updates automatically reflect in WebSSH)
- Focus on the core value: interactive terminal access in the browser

## What Changes

### REMOVED: SSH Management System
- **SSHConnection model and table**
  - Remove `SSHConnection` model from backend
  - Remove `ssh_connections` table from database
  - Remove SSH connection CRUD APIs (`/api/v1/ssh/connections/*`)
  - Remove SSH connection management UI

- **SSHKey model and table**
  - Remove `SSHKey` model from backend
  - Remove `ssh_keys` table from database
  - Remove SSH key CRUD APIs (`/api/v1/ssh/keys/*`)
  - Remove SSH key management UI

- **SSHSession model and table**
  - Remove `SSHSession` model from backend
  - Remove `ssh_sessions` table from database
  - Remove SSH session CRUD APIs (`/api/v1/ssh/sessions/*`)
  - Remove SSH session management UI

- **SSH Management Page**
  - Remove `frontend/src/pages/SSH.tsx` (entire page)
  - Remove SSH-related menu items
  - Remove SSH-related routes

### ADDED: WebSSH Management
- **WebSSH Terminal Access**
  - Direct terminal access from host management page
  - WebSocket-based terminal communication
  - Real-time interactive terminal (xterm.js)
  - Support multiple concurrent terminal sessions (tabs)
  - Terminal controls (connect, disconnect, clear, resize)

- **Host-based SSH Configuration**
  - SSH connection uses host data directly:
    - hostname: from `host.hostname` or `host.ip_address`
    - port: from `host.ssh_port` (default: 22)
    - username: prompt user when connecting (no pre-configuration)
    - authentication: prompt for password when connecting (no key management)

- **WebSSH Terminal Handler**
  - Backend WebSocket handler at `/api/v1/webssh/terminal/:host_id`
  - Real-time bidirectional communication (stdin/stdout/stderr)
  - Terminal size negotiation (rows/columns)
  - Session authentication and authorization via JWT

- **WebSSH Terminal Component**
  - Reusable terminal component (`frontend/src/components/Terminal.tsx`)
  - Accepts `hostId` and `hostName` as props
  - Integrates with host management page
  - Supports terminal resizing and multiple tabs

### MODIFIED: Host Management
- **Host Management Page Integration**
  - Add "打开终端" button for each host in host list
  - Support opening terminal in new tab
  - Display terminal connection status

- **Host Model**
  - Keep existing `ssh_port` field (used for WebSSH connection)
  - Remove SSH-related fields if they exist (`ssh_username`, `ssh_key_id`)
  - No additional SSH configuration fields needed

## Impact

### Affected Specs
- `ssh-management` - REMOVED (entire capability)
- `webssh-management` - ADDED (new capability)
- `host-management` - MODIFIED (add terminal access)

### Affected Code

**Backend:**
- `backend/internal/models/ssh.go` - REMOVED (entire file)
- `backend/internal/models/host.go` - MODIFIED (remove SSH-related fields)
- `backend/internal/models/database.go` - MODIFIED (remove SSH model migrations)
- `backend/internal/repository/ssh_repository.go` - REMOVED (entire file)
- `backend/internal/service/ssh_service.go` - REMOVED (entire file)
- `backend/internal/handler/ssh_handler.go` - REMOVED (entire file)
- `backend/internal/handler/ssh_terminal_handler.go` - REMOVED (replaced with webssh handler)
- `backend/internal/handler/webssh_handler.go` - ADDED (new WebSSH handler)
- `backend/cmd/api/main.go` - MODIFIED (remove SSH routes, add WebSSH routes)

**Frontend:**
- `frontend/src/pages/SSH.tsx` - REMOVED (entire file)
- `frontend/src/components/Terminal.tsx` - MODIFIED (use hostId instead of connectionId)
- `frontend/src/pages/Hosts.tsx` - MODIFIED (add terminal access)
- `frontend/src/services/ssh.ts` - REMOVED (entire file)
- `frontend/src/services/webssh.ts` - ADDED (new WebSSH service)
- `frontend/src/components/MainLayout.tsx` - MODIFIED (remove SSH menu item)
- `frontend/src/App.tsx` - MODIFIED (remove SSH route)

### Database
- **REMOVED tables:**
  - `ssh_connections`
  - `ssh_keys`
  - `ssh_sessions`
- **MODIFIED tables:**
  - `hosts` - Remove `ssh_username` and `ssh_key_id` fields (if they exist)
- **Migration required:**
  - Drop SSH-related tables
  - Remove SSH-related columns from hosts table
  - No data migration needed (SSH data will be lost)

### Dependencies
- **Backend:**
  - Keep `golang.org/x/crypto/ssh` (for SSH client)
  - Keep `github.com/gorilla/websocket` (for WebSocket)
  - Remove any SSH key parsing libraries (if only used for SSH key management)
- **Frontend:**
  - Keep `xterm` and `xterm-addon-fit` (for terminal emulation)
  - Remove any SSH key management UI libraries (if any)

### Breaking Changes
- **API Breaking Changes:**
  - All `/api/v1/ssh/*` endpoints removed
  - New `/api/v1/webssh/terminal/:host_id` endpoint added
- **Frontend Breaking Changes:**
  - SSH management page removed
  - SSH menu item removed
  - Terminal access moved to host management page
- **Database Breaking Changes:**
  - All SSH-related tables removed
  - SSH-related columns removed from hosts table
- **Data Loss:**
  - All SSH connection configurations will be lost
  - All SSH keys will be lost
  - All SSH session history will be lost

### Migration Strategy
- **For existing deployments:**
  - No data migration (SSH data intentionally discarded)
  - Users will need to reconnect terminals (password prompt on first connection)
- **For new deployments:**
  - No migration needed
  - Clean start with WebSSH only

