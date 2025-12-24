# Implementation Tasks

## 1. Backend Model Changes

- [x] 1.1 Remove SSH models
  - Delete `backend/internal/models/ssh.go` (entire file)
  - Remove SSH model imports from `backend/internal/models/database.go`
  - Remove SSH models from AutoMigrate list

- [x] 1.2 Update Host model
  - Remove `ssh_username` field (if exists)
  - Remove `ssh_key_id` field (if exists)
  - Remove `SSHKey` association (if exists)
  - Keep `ssh_port` field (used for WebSSH)

- [x] 1.3 Create database migration
  - Create migration script to drop SSH tables:
    - `ssh_connections`
    - `ssh_keys`
    - `ssh_sessions`
  - Create migration script to remove SSH columns from `hosts` table:
    - `ssh_username` (if exists)
    - `ssh_key_id` (if exists)

## 2. Backend Repository Removal

- [x] 2.1 Remove SSH repositories
  - Delete `backend/internal/repository/ssh_repository.go` (entire file)
  - Remove SSH repository initialization from `backend/cmd/api/main.go`

## 3. Backend Service Removal

- [x] 3.1 Remove SSH services
  - Delete `backend/internal/service/ssh_service.go` (entire file)
  - Remove SSH service initialization from `backend/cmd/api/main.go`

## 4. Backend Handler Changes

- [x] 4.1 Remove SSH handlers
  - Delete `backend/internal/handler/ssh_handler.go` (entire file)
  - Delete `backend/internal/handler/ssh_terminal_handler.go` (entire file)
  - Remove SSH handler initialization from `backend/cmd/api/main.go`

- [x] 4.2 Create WebSSH handler
  - Create `backend/internal/handler/webssh_handler.go`
  - Implement WebSocket handler for terminal connections
  - Use host data for SSH connection configuration
  - Support password authentication (prompt via WebSocket)
  - Handle terminal I/O (stdin/stdout/stderr)
  - Support terminal resize

- [x] 4.3 Update routes
  - Remove all `/api/v1/ssh/*` routes from `backend/cmd/api/main.go`
  - Add `/api/v1/webssh/terminal/:host_id` WebSocket route
  - Add authentication middleware for WebSSH route

## 5. Backend Utils Cleanup

- [x] 5.1 Review SSH-related utilities
  - Keep `backend/internal/utils/encrypt.go` (may be used elsewhere)
  - Remove `backend/internal/utils/ssh_fingerprint.go` (if only used for SSH key management)
  - Check if encryption utilities are used by other features

## 6. Frontend Service Changes

- [x] 6.1 Remove SSH service
  - Delete `frontend/src/services/ssh.ts` (entire file)

- [x] 6.2 Create WebSSH service
  - Create `frontend/src/services/webssh.ts`
  - Implement WebSocket client for terminal connections
  - Handle terminal message protocol (input, resize, ping/pong)

## 7. Frontend Component Changes

- [x] 7.1 Update Terminal component
  - Modify `frontend/src/components/Terminal.tsx`
  - Change props from `connectionId` to `hostId`
  - Update WebSocket URL to use `/api/v1/webssh/terminal/:host_id`
  - Update connection logic to use host data

- [x] 7.2 Remove SSH page
  - Delete `frontend/src/pages/SSH.tsx` (entire file)

- [x] 7.3 Update Host management page
  - Modify `frontend/src/pages/Hosts.tsx`
  - Add "打开终端" button for each host
  - Implement terminal tab management
  - Support opening multiple terminal sessions

## 8. Frontend Routing and Menu

- [x] 8.1 Update routes
  - Remove SSH route from `frontend/src/App.tsx`
  - Remove SSH route configuration

- [x] 8.2 Update menu
  - Remove SSH menu item from `frontend/src/components/MainLayout.tsx`
  - Update menu structure

## 9. Testing

- [ ] 9.1 Backend tests
  - Remove SSH-related unit tests
  - Remove SSH-related integration tests
  - Add WebSSH handler unit tests
  - Add WebSSH WebSocket integration tests

- [ ] 9.2 Frontend tests
  - Remove SSH page tests
  - Update Terminal component tests
  - Update Host management page tests
  - Add WebSSH service tests

- [ ] 9.3 E2E tests
  - Remove SSH management E2E tests
  - Add WebSSH terminal E2E tests
  - Test terminal connection from host management page

## 10. Documentation

- [ ] 10.1 Update API documentation
  - Remove SSH API documentation
  - Add WebSSH API documentation
  - Document WebSocket message protocol

- [ ] 10.2 Update user manual
  - Remove SSH management section
  - Add WebSSH usage section
  - Document how to open terminal from host management

- [ ] 10.3 Update README
  - Remove SSH management features from feature list
  - Add WebSSH management to feature list
  - Update architecture diagram (if exists)

## 11. Database Migration Execution

- [x] 11.1 Create migration scripts
  - Write SQL migration to drop SSH tables
  - Write SQL migration to remove SSH columns from hosts
  - Test migration scripts on development database

- [ ] 11.2 Execute migration
  - Run migration on development environment
  - Verify tables and columns are removed
  - Test application functionality after migration

## 12. Validation

- [x] 12.1 Code validation
  - Run `go build` to verify no compilation errors
  - Run `npm run build` to verify frontend compiles
  - Fix any linting errors

- [ ] 12.2 Functional validation
  - Verify WebSSH terminal opens from host management page
  - Verify terminal connection works with password authentication
  - Verify terminal I/O works correctly
  - Verify terminal resize works
  - Verify multiple terminal tabs work
  - Verify terminal disconnection works

- [x] 12.3 OpenSpec validation
  - Run `openspec validate remove-ssh-add-webssh --strict`
  - Fix any validation errors

