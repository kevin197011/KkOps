## 1. Backend Model and Database
- [x] 1.1 Create SSH Key model
  - Create `backend/internal/models/ssh_key.go`
  - Define `SSHKey` struct with fields: id, user_id, name, key_type, private_key (encrypted), public_key, fingerprint, created_at, updated_at
  - Add to AutoMigrate in `database.go`
- [x] 1.2 Add optional ssh_key_id to Host model
  - Add `SSHKeyID` field to `backend/internal/models/host.go`
  - Add foreign key relationship to SSHKey (optional)
  - Create database migration to add `ssh_key_id` column to `hosts` table
- [x] 1.3 Create database migration
  - Create migration script to create `ssh_keys` table
  - Add foreign key constraint to `users` table
  - Add optional foreign key from `hosts.ssh_key_id` to `ssh_keys.id`

## 2. Backend Repository
- [x] 2.1 Create SSH Key repository
  - Create `backend/internal/repository/ssh_key_repository.go`
  - Implement CRUD operations: Create, GetByID, GetByUserID, List, Update, Delete
  - Implement encryption/decryption of private keys

## 3. Backend Service
- [x] 3.1 Create SSH Key service
  - Create `backend/internal/service/ssh_key_service.go`
  - Implement business logic: CreateKey, GetKey, ListKeys, UpdateKey, DeleteKey
  - Implement key validation (format, type detection)
  - Implement fingerprint calculation
  - Implement encryption/decryption using existing utilities

## 4. Backend Handler
- [x] 4.1 Create SSH Key handler
  - Create `backend/internal/handler/ssh_key_handler.go`
  - Implement REST API endpoints:
    - POST `/api/v1/ssh-keys` - Upload/create SSH key
    - GET `/api/v1/ssh-keys` - List user's SSH keys
    - GET `/api/v1/ssh-keys/:id` - Get SSH key details (without private key)
    - PUT `/api/v1/ssh-keys/:id` - Update SSH key (name only)
    - DELETE `/api/v1/ssh-keys/:id` - Delete SSH key
  - Add authentication middleware
  - Add user-scoped access control (users can only access their own keys)
- [x] 4.2 Update WebSSH handler
  - Modify `backend/internal/handler/webssh_handler.go`
  - Update `establishSSHConnection` to support key-based authentication
  - Add logic to retrieve SSH key from database if key_id is provided
  - Support both key-based and password-based authentication
  - Update `requestCredentials` to request key selection or password

## 5. Frontend Service
- [x] 5.1 Create SSH Key service
  - Create `frontend/src/services/sshKey.ts`
  - Implement API calls: create, list, get, update, delete
  - Define TypeScript interfaces for SSH key data

## 6. Frontend WebSSH Page
- [x] 6.1 Add SSH Key management tab
  - Modify `frontend/src/pages/WebSSH.tsx`
  - Add "SSH密钥" tab to Tabs component
  - Implement SSH key list table with columns: name, type, fingerprint, created date, actions
  - Add "上传密钥" button
  - Implement upload modal/form for SSH key
  - Implement edit/delete actions for each key
- [x] 6.2 Add key selection to terminal connection
  - Add authentication method selection (key/password) when opening terminal
  - Show key selection dropdown if key method is selected
  - Pass selected key_id to Terminal component
  - Support auto-selecting default key if host has ssh_key_id configured

## 7. Frontend Terminal Component
- [x] 7.1 Update Terminal component for key authentication
  - Modify `frontend/src/components/Terminal.tsx`
  - Add props for `sshKeyId` (optional)
  - Update connection flow to send key_id if provided
  - Update authentication UI to show key selection if needed
  - Handle key authentication errors gracefully

## 8. Frontend Host Management
- [x] 8.1 Add SSH key selection to host form
  - Modify `frontend/src/pages/Hosts.tsx`
  - Add SSH key selection dropdown in host create/edit form
  - Load user's SSH keys for selection
  - Save ssh_key_id when creating/updating host

## 9. Testing
- [ ] 9.1 Backend tests
  - Add SSH key repository unit tests
  - Add SSH key service unit tests
  - Add SSH key handler integration tests
  - Add WebSSH handler key authentication tests
- [ ] 9.2 Frontend tests
  - Add SSH key management UI tests
  - Add terminal key selection tests
  - Add host SSH key configuration tests

## 10. Validation
- [x] 10.1 Code validation
  - Run `go build` to verify backend compiles
  - Run `npm run build` to verify frontend compiles
  - Fix any linting errors
- [x] 10.2 Functional validation
  - Verify SSH key upload works
  - Verify SSH key list displays correctly
  - Verify key-based terminal connection works
  - Verify password fallback works if key auth fails
  - Verify per-host key configuration works
  - Verify key selection in terminal connection flow
- [x] 10.3 OpenSpec validation
  - Run `openspec validate add-ssh-key-management --strict`
  - Fix any validation errors

