# Design: WebSSH Terminal Implementation

## Context
The current SSH management provides only connection configuration and session metadata. Users need interactive terminal access through the browser, similar to tools like Bastillion or WebSSH.

## Goals / Non-Goals

### Goals
- Provide interactive SSH terminal in browser
- Support multiple concurrent terminal sessions
- Real-time command execution and output
- Secure WebSocket-based communication
- Terminal emulation with proper ANSI support
- Session management (connect, disconnect, reconnect)

### Non-Goals
- File transfer (SFTP) - can be added later
- Port forwarding - can be added later
- Terminal recording/replay - can be added later
- Multi-user terminal sharing - can be added later

## Decisions

### Decision: Use xterm.js for Terminal Emulation
- **Rationale**: xterm.js is the most popular and mature browser-based terminal emulator
- **Alternatives considered**: 
  - hterm (deprecated)
  - Terminal.js (less features)
- **Trade-offs**: Adds ~200KB to bundle size, but provides excellent terminal emulation

### Decision: WebSocket for Real-time Communication
- **Rationale**: WebSocket provides full-duplex communication needed for terminal I/O
- **Alternatives considered**:
  - Server-Sent Events (SSE) - only one-way
  - Long polling - inefficient
- **Trade-offs**: Requires WebSocket server, but standard and well-supported

### Decision: Backend SSH Client per Session
- **Rationale**: Each WebSocket connection maintains its own SSH client connection
- **Alternatives considered**:
  - Shared SSH connections - complex connection pooling
  - SSH multiplexing - adds complexity
- **Trade-offs**: One SSH connection per terminal session, but simpler and more reliable

### Decision: Terminal Size Negotiation
- **Rationale**: Terminal size affects command output formatting
- **Implementation**: Send terminal dimensions (rows/columns) via WebSocket messages
- **Trade-offs**: Requires handling window resize events

## Architecture

### Backend Flow
```
Browser (xterm.js) 
  → WebSocket Connection 
    → SSH Handler (gorilla/websocket)
      → SSH Client (golang.org/x/crypto/ssh)
        → Remote SSH Server
```

### Frontend Flow
```
SSH Management Page
  → Terminal Component (xterm.js)
    → WebSocket Client
      → Backend WebSocket Handler
```

### Data Flow
1. User clicks "Connect" on a connection
2. Frontend opens WebSocket connection to `/api/v1/ssh/terminal/:connection_id`
3. Backend authenticates and establishes SSH connection
4. Backend creates SSHSession record
5. Terminal I/O flows bidirectionally:
   - User input → WebSocket → SSH stdin
   - SSH stdout/stderr → WebSocket → Terminal display

## Risks / Trade-offs

### Security Risks
- **Risk**: SSH credentials exposed in browser
  - **Mitigation**: Credentials stored encrypted on backend, never sent to frontend
- **Risk**: WebSocket connections not authenticated
  - **Mitigation**: Require JWT token in WebSocket handshake
- **Risk**: Terminal output may contain sensitive data
  - **Mitigation**: Audit logging, session recording (future)

### Performance Risks
- **Risk**: Many concurrent WebSocket connections
  - **Mitigation**: Connection limits, resource monitoring
- **Risk**: Large terminal output buffers
  - **Mitigation**: Limit buffer size, implement scrolling

### UX Trade-offs
- **Trade-off**: Terminal vs. connection list view
  - **Decision**: Make terminal primary, keep connection list as sidebar or modal
- **Trade-off**: Multiple tabs vs. single terminal
  - **Decision**: Support multiple terminal tabs for concurrent sessions

## Migration Plan

### Phase 1: Backend WebSocket Infrastructure
1. Add WebSocket handler
2. Implement SSH connection establishment
3. Implement terminal I/O forwarding
4. Add session management

### Phase 2: Frontend Terminal Component
1. Install xterm.js
2. Create Terminal component
3. Implement WebSocket client
4. Add terminal controls

### Phase 3: Integration
1. Update SSH management page
2. Add connection selection UI
3. Implement multi-tab support
4. Add session management UI

### Phase 4: Testing & Polish
1. Test various terminal types
2. Test connection failures
3. Test concurrent sessions
4. Performance optimization

## Open Questions
- Should we support terminal themes/color schemes?
- Should we implement terminal history/search?
- Should we add terminal recording for audit purposes?
- What is the maximum number of concurrent sessions per user?

