## 1. Database Schema Changes

- [ ] 1.1 Add optional SSH fields to Host table
  - Add `ssh_username` field (VARCHAR(100), nullable)
  - Add `ssh_key_id` field (BIGINT, nullable, foreign key to ssh_keys)
  - Create migration script

- [ ] 1.2 Update SSHSession table
  - Rename `connection_id` to `host_id`
  - Update foreign key constraint to reference hosts table
  - Create migration script

- [ ] 1.3 Data migration (if preserving existing connections)
  - Create script to migrate SSHConnection data to Host table
  - Map connection hostname/IP to host records
  - Preserve SSH username and key associations

- [ ] 1.4 Remove SSHConnection table
  - Drop `ssh_connections` table
  - Remove foreign key constraints
  - Create migration script

## 2. Backend Model Changes

- [ ] 2.1 Update Host model
  - Add `SSHUsername` field (string, optional)
  - Add `SSHKeyID` field (*uint64, optional)
  - Add `SSHKey` association (gorm foreign key)

- [ ] 2.2 Update SSHSession model
  - Change `ConnectionID` to `HostID`
  - Update foreign key to reference Host
  - Update `Connection` association to `Host`

- [ ] 2.3 Remove SSHConnection model
  - Delete `backend/internal/models/ssh.go` SSHConnection struct
  - Update imports and references

## 3. Backend Repository Changes

- [ ] 3.1 Remove SSHConnection repository
  - Delete `SSHConnectionRepository` interface
  - Delete `sshConnectionRepository` implementation
  - Remove from repository initialization

- [ ] 3.2 Update SSHSession repository
  - Change `connection_id` filter to `host_id`
  - Update `Preload` to load Host instead of Connection
  - Update query methods

- [ ] 3.3 Update Host repository (if needed)
  - Add methods to query hosts with SSH configuration
  - Add filtering by SSH key

## 4. Backend Service Changes

- [ ] 4.1 Remove SSHConnection service
  - Delete `SSHConnectionService` interface
  - Delete `sshConnectionService` implementation
  - Remove auto-generation methods (no longer needed)

- [ ] 4.2 Update SSHSession service
  - Change session creation to use host ID
  - Update session queries to filter by host ID
  - Update session retrieval methods

- [ ] 4.3 Update SSH terminal handler
  - Change to accept `host_id` instead of `connection_id`
  - Load host data instead of connection data
  - Use host's SSH configuration (username, key, port)
  - Handle password authentication (prompt via WebSocket)

## 5. Backend Handler Changes

- [ ] 5.1 Remove SSH connection handlers
  - Delete `CreateConnection`, `ListConnections`, `GetConnection`, `UpdateConnection`, `DeleteConnection` handlers
  - Delete `AutoGenerateConnections` handlers
  - Remove connection routes from main.go

- [ ] 5.2 Update SSH terminal handler route
  - Change route from `/api/v1/ssh/terminal/:connection_id` to `/api/v1/ssh/terminal/:host_id`
  - Update handler to use host ID

- [ ] 5.3 Update SSH session handlers
  - Update session creation to use host ID
  - Update session listing to filter by host ID
  - Update session display to show host information

## 6. Frontend Service Changes

- [ ] 6.1 Remove SSH connection service methods
  - Delete `createConnection`, `listConnections`, `getConnection`, `updateConnection`, `deleteConnection` from `ssh.ts`
  - Delete `autoGenerateConnections` methods
  - Remove `SSHConnection` interface

- [ ] 6.2 Update SSH session service
  - Change session creation to use host ID
  - Update session interfaces to reference host instead of connection

- [ ] 6.3 Update Terminal component
  - Change prop from `connectionId` to `hostId`
  - Update WebSocket URL to use host ID
  - Update connection name display

## 7. Frontend Page Changes

- [ ] 7.1 Redesign SSH management page
  - Replace "SSH连接" tab with "主机" tab
  - Display hosts from host management (reuse host service)
  - Add project filter for hosts
  - Add "打开终端" button for each host

- [ ] 7.2 Update host list display
  - Show host name, IP address, SSH port
  - Show SSH configuration status (has key/username configured)
  - Add quick actions: "打开终端", "配置SSH"

- [ ] 7.3 Add SSH configuration for hosts
  - Add modal/form to configure SSH username and key for a host
  - Allow setting default SSH credentials per host
  - Update host via host service

- [ ] 7.4 Update terminal tab management
  - Change terminal tabs to use host ID and host name
  - Update tab creation and closing logic

- [ ] 7.5 Update session management tab
  - Display sessions with host information instead of connection
  - Filter sessions by host
  - Show host name in session list

## 8. Testing

- [ ] 8.1 Test database migration
  - Test adding SSH fields to Host table
  - Test updating SSHSession table
  - Test data migration (if applicable)
  - Test dropping SSHConnection table

- [ ] 8.2 Test WebSSH with host data
  - Test connection using host's SSH key
  - Test connection using host's SSH username
  - Test connection with password prompt
  - Test connection with default port (22)

- [ ] 8.3 Test SSH management page
  - Test host list display
  - Test opening terminal for host
  - Test SSH configuration for host
  - Test session management

- [ ] 8.4 Test edge cases
  - Test host without SSH port configured
  - Test host without SSH key or username
  - Test multiple terminal sessions for same host
  - Test session cleanup

## 9. Documentation

- [ ] 9.1 Update API documentation
  - Remove SSH connection endpoints
  - Update SSH terminal endpoint (host_id instead of connection_id)
  - Update session endpoints

- [ ] 9.2 Update user manual
  - Document new SSH management workflow
  - Document host-based SSH configuration
  - Document password authentication

- [ ] 9.3 Update migration guide
  - Document database migration steps
  - Document data migration (if applicable)
  - Document breaking changes

