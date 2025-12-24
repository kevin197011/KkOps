# Design: Remove SSH Management and Add WebSSH Management

## Context
The current SSH management system maintains separate models for SSH connections, SSH keys, and SSH sessions, creating unnecessary complexity and data duplication. This design removes all SSH management functionality and replaces it with a simplified WebSSH management system that directly uses host data.

## Goals / Non-Goals

### Goals
- Remove all SSH management complexity (connections, keys, sessions)
- Provide direct terminal access from host management
- Simplify user experience (no pre-configuration needed)
- Maintain core WebSSH terminal functionality
- Use host data directly for SSH connections

### Non-Goals
- SSH key management (users will use password authentication)
- SSH connection pre-configuration (connections are established on-demand)
- SSH session history (sessions are ephemeral)
- SSH connection testing (connections are established when terminal opens)
- SSH key generation or management

## Decisions

### Decision: Remove SSH Key Management
- **Rationale**: SSH key management adds complexity without clear value. Users can use password authentication, which is simpler for a web-based terminal.
- **Alternatives considered**:
  - Keep SSH key management but simplify it
  - Support both password and key authentication
- **Trade-offs**: 
  - Simpler implementation
  - Users must enter password each time (no key-based auth)
  - Less secure than key-based auth (but acceptable for internal tools)

### Decision: Remove SSH Connection Pre-configuration
- **Rationale**: Connection information already exists in the Host table. Pre-configuring connections creates redundancy.
- **Alternatives considered**:
  - Keep connection model but auto-generate from hosts
  - Use connection model but simplify it
- **Trade-offs**:
  - No connection management overhead
  - Users must enter credentials each time (no saved credentials)
  - Simpler data model

### Decision: Remove SSH Session History
- **Rationale**: Session history is not essential for terminal access. Users can reconnect as needed.
- **Alternatives considered**:
  - Keep session model but simplify it
  - Keep session history for audit purposes
- **Trade-offs**:
  - Simpler implementation
  - No session history or audit trail
  - Sessions are ephemeral (no persistence)

### Decision: Use Host Data Directly
- **Rationale**: Host table already contains all necessary information (hostname, IP, SSH port).
- **Implementation**: 
  - Terminal connection uses `host.hostname` or `host.ip_address`
  - Terminal connection uses `host.ssh_port` (default: 22)
  - Username and password prompted via WebSocket when connecting
- **Trade-offs**:
  - No separate connection configuration
  - Users must enter credentials each time
  - Host updates automatically reflect in terminal access

### Decision: Password Authentication Only
- **Rationale**: Simplifies implementation and user experience. No key management needed.
- **Implementation**:
  - WebSocket handler prompts for password when connecting
  - Password sent securely via WebSocket (encrypted in transit)
  - Password not stored (entered each time)
- **Trade-offs**:
  - Less secure than key-based auth
  - Users must enter password each time
  - Simpler implementation

## Architecture

### Backend Flow
```
Browser (xterm.js)
  → WebSocket Connection (/api/v1/webssh/terminal/:host_id)
    → WebSSH Handler (authenticate, get host data)
      → SSH Client (golang.org/x/crypto/ssh)
        → Remote SSH Server
```

### Frontend Flow
```
Host Management Page
  → User clicks "打开终端"
    → Terminal Component (xterm.js)
      → WebSocket Client
        → Backend WebSocket Handler
```

### Data Flow
1. User clicks "打开终端" on a host in host management page
2. Frontend opens WebSocket connection to `/api/v1/webssh/terminal/:host_id`
3. Backend authenticates user (JWT token)
4. Backend loads host data from database
5. Backend prompts for username and password via WebSocket
6. Frontend sends username and password
7. Backend establishes SSH connection using host data and credentials
8. Terminal I/O flows bidirectionally:
   - User input → WebSocket → SSH stdin
   - SSH stdout/stderr → WebSocket → Terminal display

## Security Considerations

### Authentication
- **User Authentication**: JWT token required in WebSocket handshake
- **SSH Authentication**: Password authentication (prompted via WebSocket)
- **No Credential Storage**: Passwords not stored, entered each time

### Data Security
- **WebSocket Encryption**: Use WSS (WebSocket Secure) in production
- **Password Transmission**: Password sent via WebSocket (encrypted in transit if using WSS)
- **No Password Storage**: Passwords never stored in database

### Access Control
- **Host Access**: Users can only access hosts they have permission to view
- **Project-based Access**: Terminal access respects project-level permissions

## Risks / Trade-offs

### Security Risks
- **Risk**: Password authentication less secure than key-based auth
  - **Mitigation**: Acceptable for internal tools, can add key support later if needed
- **Risk**: Password sent via WebSocket
  - **Mitigation**: Use WSS in production, password encrypted in transit
- **Risk**: No session audit trail
  - **Mitigation**: Can add audit logging later if needed

### UX Trade-offs
- **Trade-off**: Users must enter password each time
  - **Decision**: Acceptable for simplicity, can add credential caching later
- **Trade-off**: No connection pre-configuration
  - **Decision**: Acceptable, host data is sufficient
- **Trade-off**: No session history
  - **Decision**: Acceptable, sessions are ephemeral

### Performance Risks
- **Risk**: Many concurrent WebSocket connections
  - **Mitigation**: Connection limits, resource monitoring
- **Risk**: Large terminal output buffers
  - **Mitigation**: Limit buffer size, implement scrolling

## Migration Plan

### Phase 1: Backend Removal
1. Remove SSH models, repositories, services, handlers
2. Create WebSSH handler
3. Update routes
4. Test backend functionality

### Phase 2: Database Migration
1. Create migration scripts
2. Test migration on development database
3. Execute migration
4. Verify tables and columns removed

### Phase 3: Frontend Removal and Update
1. Remove SSH page and service
2. Create WebSSH service
3. Update Terminal component
4. Update Host management page
5. Update routes and menu
6. Test frontend functionality

### Phase 4: Integration Testing
1. Test terminal connection from host management
2. Test password authentication
3. Test terminal I/O
4. Test terminal resize
5. Test multiple terminal tabs
6. Test terminal disconnection

### Phase 5: Documentation and Cleanup
1. Update API documentation
2. Update user manual
3. Update README
4. Remove unused code and dependencies
5. Final validation

## Open Questions
- Should we support SSH key authentication in the future? (Can be added later if needed)
- Should we add session audit logging? (Can be added later if needed)
- Should we support credential caching? (Can be added later if needed)
- What is the maximum number of concurrent terminal sessions per user?

