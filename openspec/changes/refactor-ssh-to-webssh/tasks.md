## 1. Backend WebSocket Infrastructure

- [x] 1.1 Add WebSocket dependency
  - Add `github.com/gorilla/websocket` to `go.mod`
  - Update dependencies

- [x] 1.2 Create WebSocket handler for SSH terminal
  - Create `TerminalHandler` in `backend/internal/handler/ssh_terminal_handler.go`
  - Implement WebSocket upgrade logic
  - Add authentication middleware for WebSocket connections
  - Handle WebSocket connection lifecycle

- [x] 1.3 Implement SSH connection establishment in handler
  - Load SSH connection configuration from database
  - Decrypt credentials (password or private key)
  - Establish SSH client connection using `golang.org/x/crypto/ssh`
  - Handle connection errors and authentication failures

- [x] 1.4 Implement terminal I/O forwarding
  - Create SSH session and request PTY (pseudo-terminal)
  - Set terminal type (e.g., "xterm-256color")
  - Forward stdin from WebSocket to SSH session
  - Forward stdout/stderr from SSH session to WebSocket
  - Handle terminal resize events

- [x] 1.5 Add terminal size negotiation
  - Accept terminal dimensions (rows/columns) from WebSocket messages
  - Update SSH session window size
  - Handle window resize events

- [x] 1.6 Add WebSocket route
  - Add route `/api/v1/ssh/terminal/:connection_id` in `main.go`
  - Register WebSocket handler
  - Add authentication middleware

- [x] 1.7 Implement session lifecycle management
  - Create SSHSession record when terminal connects
  - Update session status (active/closed)
  - Clean up SSH connection on disconnect
  - Handle connection timeouts

## 2. Frontend Terminal Component

- [x] 2.1 Install terminal emulator dependencies
  - Add `xterm` and `xterm-addon-fit` to `package.json`
  - Run `npm install`

- [x] 2.2 Create Terminal component
  - Create `frontend/src/components/Terminal.tsx`
  - Initialize xterm.js terminal instance
  - Configure terminal options (theme, font, cursor)
  - Implement terminal rendering

- [x] 2.3 Implement WebSocket client
  - Create WebSocket connection to backend
  - Send terminal input to WebSocket
  - Receive terminal output from WebSocket
  - Handle WebSocket connection lifecycle (connect, disconnect, error)

- [x] 2.4 Add terminal controls
  - Connect/disconnect button
  - Clear terminal button
  - Terminal resize handling
  - Connection status indicator

- [x] 2.5 Handle terminal events
  - Keyboard input handling
  - Paste support
  - Terminal resize events
  - Window resize events (using xterm-addon-fit)

## 3. SSH Management Page Refactoring

- [x] 3.1 Redesign SSH management page layout
  - Replace connection list with terminal interface
  - Add connection selector (sidebar or dropdown)
  - Add terminal tabs for multiple sessions
  - Keep connection management in separate modal or sidebar

- [x] 3.2 Implement connection selection
  - Display list of SSH connections
  - Allow selecting connection to open terminal
  - Show connection status (online/offline)
  - Support quick connect from connection list

- [x] 3.3 Implement multi-tab terminal support
  - Support multiple terminal tabs
  - Each tab maintains its own WebSocket connection
  - Tab management (new, close, switch)
  - Tab titles show connection name

- [x] 3.4 Add terminal session management UI
  - Display active sessions
  - Show session status
  - Allow closing sessions
  - Show session duration

- [x] 3.5 Update connection management
  - Keep connection CRUD in modal or separate page
  - Add "Open Terminal" action to connection list
  - Show terminal status for each connection

## 4. Integration & Testing

- [ ] 4.1 Test SSH connection establishment
  - Test password authentication
  - Test key-based authentication
  - Test connection failures
  - Test authentication failures

- [ ] 4.2 Test terminal functionality
  - Test command execution
  - Test output display
  - Test terminal resize
  - Test special characters and ANSI codes

- [ ] 4.3 Test concurrent sessions
  - Test multiple terminal tabs
  - Test session isolation
  - Test resource cleanup

- [ ] 4.4 Test error handling
  - Test network disconnections
  - Test SSH connection timeouts
  - Test invalid credentials
  - Test connection refused errors

- [ ] 4.5 Performance testing
  - Test with large output
  - Test with many concurrent sessions
  - Test terminal scrolling performance
  - Test memory usage

## 5. Security & Audit

- [ ] 5.1 Implement WebSocket authentication
  - Require JWT token in WebSocket handshake
  - Validate token before establishing SSH connection
  - Check user permissions for SSH connection

- [ ] 5.2 Add audit logging for terminal sessions
  - Log terminal session start/end
  - Log connection attempts
  - Log authentication failures
  - Consider command logging (optional, for compliance)

- [ ] 5.3 Implement connection limits
  - Limit concurrent sessions per user
  - Limit total concurrent sessions
  - Add rate limiting for connection attempts

## 6. Documentation

- [ ] 6.1 Update API documentation
  - Document WebSocket endpoint
  - Document message protocol
  - Document terminal size negotiation

- [ ] 6.2 Update user manual
  - Document terminal usage
  - Document connection management
  - Document keyboard shortcuts
  - Document troubleshooting

