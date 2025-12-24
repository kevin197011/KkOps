## Context
Currently, WebSSH only supports password-based authentication. Users must enter username and password each time they connect to a host. This is inconvenient and less secure than key-based authentication. We need to add SSH key management that integrates seamlessly with the existing WebSSH workflow.

## Goals / Non-Goals

### Goals
- Provide SSH key management (upload, view, delete) within WebSSH context
- Support key-based authentication in WebSSH terminal connections
- Allow users to choose between key and password authentication when connecting
- Support optional per-host default SSH key configuration
- Encrypt and securely store SSH private keys
- Maintain simplified architecture (no separate SSH management page)

### Non-Goals
- SSH key generation (users upload their own keys)
- SSH key sharing between users (keys are user-specific)
- SSH key rotation automation
- SSH key expiration management
- Multiple keys per host (one default key per host)

## Decisions

### Decision: Integrate SSH Key Management into WebSSH Page
**What**: Add SSH key management as a tab in the WebSSH management page, alongside the host list tab.

**Why**: 
- Keeps SSH-related functionality together (keys are used for WebSSH)
- Maintains simplified navigation (no separate menu item needed)
- Users managing terminals will naturally find key management nearby
- Consistent with the design principle of keeping related features together

**Alternatives considered**:
- Separate SSH Key management page - **Rejected**: Adds navigation complexity
- SSH Key management in host management page - **Rejected**: Keys are not host-specific, they're user-specific
- SSH Key management in user settings - **Rejected**: Keys are used for WebSSH, should be accessible from WebSSH page

### Decision: User-Scoped SSH Keys
**What**: SSH keys are associated with the user who uploads them, not with hosts.

**Why**: 
- Users may have different keys for different purposes
- Keys are personal credentials, should be user-scoped
- Allows users to manage their own keys independently

**Alternatives considered**:
- System-wide shared keys - **Rejected**: Security risk, keys should be user-specific
- Project-scoped keys - **Rejected**: Adds complexity, user-scoped is simpler

### Decision: Optional Per-Host Key Configuration
**What**: Allow (but not require) configuring a default SSH key for each host in the host management page.

**Why**: 
- Provides convenience (auto-select key when opening terminal)
- Maintains flexibility (users can still choose different keys or use password)
- Optional configuration doesn't force users to configure keys

**Alternatives considered**:
- No per-host key configuration - **Rejected**: Less convenient, users must select key each time
- Required per-host key configuration - **Rejected**: Too restrictive, password auth should remain available

### Decision: Key Selection in Terminal Connection Flow
**What**: When opening a terminal, show a dialog/modal to select authentication method (key or password) and choose a key if key method is selected.

**Why**: 
- Clear user choice at connection time
- Allows users to override default key selection
- Maintains password authentication as fallback option

**Alternatives considered**:
- Auto-use default key without prompt - **Rejected**: Users should have control over authentication
- Key selection in WebSSH page before opening terminal - **Rejected**: Less intuitive, connection-time selection is clearer

### Decision: Encrypt SSH Private Keys
**What**: Store SSH private keys encrypted in the database using the existing encryption utilities.

**Why**: 
- Security best practice
- Protects sensitive key material
- Uses existing encryption infrastructure

**Alternatives considered**:
- Store keys in plaintext - **Rejected**: Security risk
- Use external key management service - **Rejected**: Adds complexity, current encryption is sufficient

## Risks / Trade-offs

### Risk: Key Management Complexity
**Mitigation**: 
- Keep UI simple (upload, view, delete)
- Clear documentation on key format requirements
- Validation of key format on upload

### Risk: Key Security
**Mitigation**: 
- Encrypt keys at rest
- Use existing encryption utilities
- Keys are user-scoped (users can only see their own keys)
- Consider key access logging for audit

### Risk: Authentication Flow Complexity
**Mitigation**: 
- Clear UI for selecting authentication method
- Fallback to password if key auth fails
- Good error messages

## Migration Plan
1. Create SSH key model and database table
2. Add SSH key management UI to WebSSH page
3. Update WebSSH handler to support key authentication
4. Add optional `ssh_key_id` to Host model
5. Update Terminal component for key selection
6. Test key-based authentication flow
7. Deploy

No data migration needed - this is a new feature.

## Open Questions
- Should we support passphrase-protected keys? (Initial: No, simplify by requiring keys without passphrase)
- Should we validate key format on upload? (Initial: Yes, validate it's a valid SSH private key)
- Should we show key fingerprint for verification? (Initial: Yes, display fingerprint for user verification)

