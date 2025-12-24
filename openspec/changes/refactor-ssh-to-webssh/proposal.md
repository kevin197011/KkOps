# Change: Refactor SSH Management to WebSSH Terminal

## Why
The current SSH management implementation only provides connection configuration and session metadata management. Users need a true WebSSH terminal interface that allows them to:
- Open interactive SSH terminals directly in the browser
- Execute commands and see real-time output
- Manage multiple terminal sessions simultaneously
- Have a terminal emulator experience similar to native SSH clients

This refactoring transforms the SSH management from a configuration-only interface into a fully functional WebSSH terminal management system.

## What Changes
- **MODIFIED**: SSH Management becomes WebSSH Terminal Management
  - Replace connection list view with terminal session management
  - Add WebSocket-based terminal communication
  - Implement terminal emulator UI component (xterm.js or similar)
  - Support real-time command execution and output display
  - Support multiple concurrent terminal sessions
  - Add terminal session controls (resize, clear, disconnect)

- **ADDED**: WebSocket Server for Terminal Communication
  - Backend WebSocket handler for SSH terminal connections
  - Real-time bidirectional communication (stdin/stdout/stderr)
  - Terminal size negotiation (rows/columns)
  - Session authentication and authorization

- **ADDED**: Terminal Emulator Frontend
  - xterm.js-based terminal component
  - Multiple tab support for concurrent sessions
  - Terminal controls (connect, disconnect, clear, resize)
  - Connection status indicators

- **MODIFIED**: SSH Connection Model
  - Keep connection configuration (hostname, port, credentials)
  - Add terminal-specific fields (terminal type, initial size)
  - Support quick connect from connection list

## Impact
- **Affected specs**: `ssh-management`
- **Affected code**:
  - `backend/internal/handler/ssh_handler.go` - Add WebSocket handler
  - `backend/internal/service/ssh_service.go` - Add terminal session management
  - `frontend/src/pages/SSH.tsx` - Complete rewrite as terminal interface
  - `frontend/src/components/Terminal.tsx` - New terminal component
  - `frontend/package.json` - Add xterm.js dependency

- **Database**: No schema changes required (existing SSHConnection and SSHSession models are sufficient)

- **Dependencies**: 
  - Backend: `golang.org/x/crypto/ssh` (already available)
  - Backend: WebSocket support (gorilla/websocket or similar)
  - Frontend: `xterm.js` for terminal emulation
  - Frontend: `xterm-addon-fit` for terminal resizing

- **Breaking Changes**: 
  - SSH management page UI completely changes from list-based to terminal-based
  - Existing session management API may need updates for terminal-specific features

