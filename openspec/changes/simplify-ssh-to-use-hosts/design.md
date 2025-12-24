# Design: Simplify SSH Management to Use Host Data

## Context
The current implementation maintains separate `SSHConnection` records that duplicate information from the `Host` table. This creates redundancy and requires users to manually create connections for each host.

## Goals / Non-Goals

### Goals
- Eliminate SSHConnection table and use Host data directly
- Provide WebSSH access for all hosts automatically
- Simplify user workflow (no separate connection creation)
- Maintain SSH key and session management functionality

### Non-Goals
- Changing SSH key management (remains as is)
- Changing session management (only the reference changes)
- Changing host management (only adding optional SSH fields)

## Decisions

### Decision: Remove SSHConnection Table
- **Rationale**: Host table already contains all necessary information (hostname, IP, SSH port)
- **Trade-offs**: 
  - ✅ Eliminates data duplication
  - ✅ Automatic SSH access for all hosts
  - ❌ Loses ability to have multiple SSH configurations per host (rarely needed)
  - ❌ Requires data migration for existing deployments

### Decision: Add Optional SSH Fields to Host Model
- **Rationale**: Some hosts may need specific SSH username or key configuration
- **Fields to add**:
  - `SSHUsername` (string, optional) - Default SSH username for this host
  - `SSHKeyID` (uint64, optional, foreign key to ssh_keys) - Default SSH key for this host
- **Trade-offs**:
  - ✅ Allows per-host SSH configuration
  - ✅ Maintains flexibility
  - ❌ Adds optional fields to Host model

### Decision: Support Password Authentication via Prompt
- **Rationale**: Not all hosts use key-based authentication
- **Implementation**: 
  - When connecting to a host without configured SSH key, prompt user for password
  - Store password in session (encrypted, temporary)
  - Do not store passwords in database
- **Trade-offs**:
  - ✅ Supports password-based authentication
  - ❌ Requires user to enter password each time (or use SSH key)

### Decision: Keep SSH Key Management Separate
- **Rationale**: SSH keys are reusable across multiple hosts
- **Implementation**: SSH keys remain in separate table, can be associated with hosts
- **Trade-offs**:
  - ✅ Keys can be reused
  - ✅ Centralized key management
  - ✅ No change to existing key management functionality

## Architecture

### Data Model Changes
```
Before:
Host (id, hostname, ip, ssh_port, ...)
SSHConnection (id, host_id, hostname, port, username, auth_type, key_id, ...)
SSHSession (id, connection_id, ...)

After:
Host (id, hostname, ip, ssh_port, ssh_username, ssh_key_id, ...)
SSHSession (id, host_id, ...)
SSHKey (id, ...) - unchanged
```

### WebSSH Connection Flow
```
1. User selects host from SSH management page
2. Frontend calls WebSocket: /api/v1/ssh/terminal/:host_id
3. Backend loads host data (hostname, IP, SSH port, SSH username, SSH key)
4. If SSH key configured: use key authentication
5. If no SSH key: prompt user for password (via WebSocket message)
6. Establish SSH connection using host data
7. Create SSHSession with host_id
```

### SSH Management Page
```
Before:
- Tab 1: SSH Connections (list of SSHConnection records)
- Tab 2: SSH Keys
- Tab 3: SSH Sessions

After:
- Tab 1: Hosts (list of Host records, filtered by project)
- Tab 2: SSH Keys
- Tab 3: SSH Sessions
```

## Risks / Trade-offs

### Data Loss Risk
- **Risk**: Existing SSH connections will be lost
- **Mitigation**: Provide migration script to convert connections to host SSH configurations

### Authentication Complexity
- **Risk**: Password authentication requires user interaction
- **Mitigation**: Encourage SSH key usage, provide clear UI for password entry

### Multiple SSH Configurations per Host
- **Risk**: Some users may need multiple SSH configs per host (different users, keys)
- **Mitigation**: 
  - Primary: Use host's default SSH configuration
  - Future: Add "SSH Profiles" feature if needed

### Performance
- **Risk**: Loading all hosts for SSH management may be slow with many hosts
- **Mitigation**: 
  - Add pagination to host list
  - Add filtering by project, status, etc.
  - Lazy load host details

## Migration Plan

### Phase 1: Database Migration
1. Add optional SSH fields to Host table
2. Migrate SSHConnection data to Host table (if preserving data)
3. Update SSHSession table (connection_id → host_id)
4. Drop SSHConnection table

### Phase 2: Backend Changes
1. Remove SSHConnection model
2. Update SSHSession model (HostID instead of ConnectionID)
3. Update SSH terminal handler to use host data
4. Remove SSH connection service and repository
5. Update session service to use host ID

### Phase 3: Frontend Changes
1. Update SSH management page to display hosts
2. Update Terminal component to accept host ID
3. Remove SSH connection management UI
4. Update session display to show host information

### Phase 4: Testing
1. Test WebSSH connection with key authentication
2. Test WebSSH connection with password authentication
3. Test session management
4. Test with multiple hosts

## Open Questions
- Should we support multiple SSH keys per host? (Answer: No, use one default key)
- Should we store SSH passwords? (Answer: No, prompt each time)
- How to handle hosts without SSH port configured? (Answer: Use default port 22)
- Should we validate SSH connectivity before allowing terminal access? (Answer: No, let connection fail naturally)

