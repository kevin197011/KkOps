# Change: Auto-Generate SSH Connections for Linux/Unix Hosts

## Why
Currently, users must manually create SSH connections for each host in the system. This is tedious and error-prone, especially when managing many Linux/Unix hosts. Since the system already tracks host operating system information (OSType) and SSH port details, we should automatically generate SSH connections for all Linux/Unix hosts to improve operational efficiency.

## What Changes
- **ADDED**: Automatic SSH connection generation for Linux/Unix hosts
  - System SHALL automatically create SSH connections for hosts with OSType matching Linux/Unix patterns
  - System SHALL provide a batch generation API endpoint
  - System SHALL provide a frontend UI for triggering auto-generation
  - System SHALL skip hosts that already have SSH connections
  - System SHALL use host's SSH port and IP address/hostname from host management
  - System SHALL allow configuration of default username and authentication method

- **MODIFIED**: SSH Connection Management
  - SSH connection list SHALL display auto-generated connections with appropriate indicators
  - System SHALL allow manual override of auto-generated connections

## Impact
- **Affected specs**: `ssh-management`
- **Affected code**:
  - `backend/internal/service/ssh_service.go` - Add auto-generation logic
  - `backend/internal/handler/ssh_handler.go` - Add batch generation endpoint
  - `backend/internal/service/host_service.go` - Query hosts by OS type
  - `frontend/src/pages/SSH.tsx` - Add auto-generation UI
  - `frontend/src/services/ssh.ts` - Add auto-generation API call

- **Database**: No schema changes required (uses existing SSHConnection and Host tables)

- **Dependencies**: 
  - Requires host OS type information to be populated (from Salt Grains)
  - Requires host SSH port to be configured

