# Design: Simplify SSH Key Creation

## Context

The current SSH key creation process requires users to provide both private and public keys manually. This is unnecessarily complex since public keys can be programmatically extracted from private keys.

## Goals

1. Simplify SSH key creation to require only private key input
2. Automatically extract public key from private key
3. Improve form layout by moving SSH user field to a more intuitive position
4. Maintain backward compatibility with existing API

## Non-Goals

- Changing the encryption/security model for private keys
- Modifying SSH key editing capabilities (still limited to name, SSH user, description)
- Changing SSH key deletion or testing functionality

## Decisions

### Decision: Auto-extract Public Key from Private Key

**What**: When a user provides only a private key, automatically extract and use the corresponding public key.

**Why**: 
- Eliminates redundant user input
- Reduces errors from mismatched key pairs
- Aligns with standard SSH key management practices

**How**:
- Use Go's `golang.org/x/crypto/ssh` package to parse the private key
- Extract the public key using `ssh.PublicKeys()` or similar method
- Support common private key formats: RSA, Ed25519, ECDSA, DSA
- Support both encrypted and unencrypted private keys

**Alternatives considered**:
- Continue requiring both keys: Rejected - adds unnecessary friction
- Generate key pair automatically: Rejected - users may want to use existing keys

### Decision: Make Public Key Optional in API

**What**: Update the API to treat `public_key` as optional in `CreateSSHKeyRequest`.

**Why**: 
- Maintains backward compatibility if clients still send public_key
- Allows frontend to omit public_key entirely

**Implementation**:
- If `public_key` is provided, use it (for backward compatibility)
- If `public_key` is empty but `private_key` is provided, extract public key from private key
- If both are empty, return validation error

### Decision: Move SSH User Field

**What**: Position SSH user field immediately after the name field in the creation form.

**Why**: 
- SSH user is a fundamental property that should be configured early
- Logical grouping: name → user → key material → description
- Better visual flow and user experience

**UI Layout**:
```
┌─────────────────────────────────────┐
│ Name: [________________]            │
│ SSH User: [_______________]         │ ← Moved here
│ Private Key: [text area]            │
│ Passphrase: [____________]          │
│ Description: [text area]            │
└─────────────────────────────────────┘
```

## Implementation Details

### Backend Changes

#### Service Layer (`backend/internal/service/sshkey/service.go`)

**Function**: `CreateSSHKey`

1. If `req.PublicKey` is provided and valid, use it (backward compatibility)
2. If `req.PublicKey` is empty but `req.PrivateKey` is provided:
   - Parse private key using `ssh.ParsePrivateKey()` or format-specific parsers
   - Extract public key using `ssh.PublicKeys()`
   - Derive public key in OpenSSH format (e.g., `ssh-rsa AAAAB3Nza... user@host`)
3. Validate extracted public key format
4. Continue with existing encryption and storage logic

**Error Handling**:
- Invalid private key format → clear error message
- Unsupported key type → error with supported types
- Encrypted key without passphrase → error prompting for passphrase

**Dependencies**:
- `golang.org/x/crypto/ssh` - for parsing private keys
- Existing encryption utilities (`utils.Encrypt`)

### Frontend Changes

#### Component: `frontend/src/pages/ssh/SSHKeyList.tsx`

1. Remove `public_key` input field from creation form
2. Move `ssh_user` field to appear after `name` field
3. Update form validation to only require `name` and `private_key`
4. Update TypeScript interface to make `public_key` optional in `CreateSSHKeyRequest`

**Form Structure**:
```tsx
<Form.Item name="name" label="名称" required>
<Form.Item name="ssh_user" label="SSH 用户">  {/* Moved here */}
<Form.Item name="private_key" label="私钥" required>
<Form.Item name="passphrase" label="私钥密码">
<Form.Item name="description" label="描述">
```

### API Interface Changes

#### `CreateSSHKeyRequest`

**Before**:
```go
type CreateSSHKeyRequest struct {
    Name        string `json:"name" binding:"required"`
    Type        string `json:"type"`
    PublicKey   string `json:"public_key" binding:"required"`  // Required
    PrivateKey  string `json:"private_key" binding:"required"`
    SSHUser     string `json:"ssh_user"`
    Passphrase  string `json:"passphrase"`
    Description string `json:"description"`
}
```

**After**:
```go
type CreateSSHKeyRequest struct {
    Name        string `json:"name" binding:"required"`
    Type        string `json:"type"`
    PublicKey   string `json:"public_key"`  // Optional - auto-extracted if empty
    PrivateKey  string `json:"private_key" binding:"required"`
    SSHUser     string `json:"ssh_user"`
    Passphrase  string `json:"passphrase"`
    Description string `json:"description"`
}
```

## Risks / Trade-offs

### Risk: Private Key Parsing Errors

**Risk**: Users may provide malformed private keys that fail to parse.

**Mitigation**: 
- Provide clear, specific error messages
- Support common formats (OpenSSH, PEM)
- Document supported key types

### Risk: Encrypted Key Handling

**Risk**: Encrypted private keys require passphrase to extract public key.

**Mitigation**:
- If passphrase is provided, use it to decrypt before extraction
- If passphrase is missing but key is encrypted, return clear error
- Update UI to indicate passphrase is required for encrypted keys

### Trade-off: Backward Compatibility

**Trade-off**: Making `public_key` optional could break clients that don't send it.

**Decision**: Keep API accepting `public_key` but make it optional. If provided, use it; if not, extract from private key. This maintains backward compatibility.

## Migration Plan

1. Update backend service to support optional `public_key` and auto-extraction
2. Update frontend form to remove public key field and reorder fields
3. Test with various private key formats (RSA, Ed25519, ECDSA)
4. Test backward compatibility (API still accepts `public_key`)
5. Update API documentation if needed

## Open Questions

- Should we validate that extracted public key matches provided public key (if both are given)?
  - **Answer**: Yes, if both are provided, validate they match and return error if mismatch

- Should we support uploading public key files separately (future enhancement)?
  - **Answer**: Out of scope for this change