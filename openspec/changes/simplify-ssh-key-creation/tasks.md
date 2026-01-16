# Tasks: Simplify SSH Key Creation

## 1. Backend Implementation

- [x] 1.1 Update `CreateSSHKeyRequest` struct to make `public_key` optional
  - Modify `backend/internal/service/sshkey/service.go`
  - Remove `binding:"required"` from `PublicKey` field
  - Update struct comments

- [x] 1.2 Implement public key extraction from private key
  - Create helper function `extractPublicKeyFromPrivateKey(privateKey, passphrase string) (string, error)`
  - Support RSA, Ed25519, ECDSA, DSA key types
  - Handle encrypted private keys using passphrase
  - Use `golang.org/x/crypto/ssh` package for parsing
  - Format public key in OpenSSH format

- [x] 1.3 Update `CreateSSHKey` service method
  - If `req.PublicKey` is provided and not empty, use it (backward compatibility)
  - If `req.PublicKey` is empty, extract from `req.PrivateKey`
  - Validate extracted public key format
  - If both provided, validate they match (optional enhancement)
  - Add error handling for invalid private key formats

- [ ] 1.4 Add unit tests for public key extraction
  - Test RSA key extraction
  - Test Ed25519 key extraction
  - Test ECDSA key extraction
  - Test encrypted private key with passphrase
  - Test invalid key formats
  - Test mismatched key pairs (if both provided)

## 2. Frontend Implementation

- [x] 2.1 Update TypeScript interface for `CreateSSHKeyRequest`
  - Modify `frontend/src/api/sshkey.ts`
  - Make `public_key` field optional

- [x] 2.2 Remove public key input field from creation form
  - Modify `frontend/src/pages/ssh/SSHKeyList.tsx`
  - Remove `public_key` Form.Item from creation modal
  - Keep field in edit modal disabled (if displayed) or remove completely

- [x] 2.3 Move SSH user field after name field
  - Reorder Form.Item components in creation modal
  - Update field order: name → ssh_user → private_key → passphrase → description

- [x] 2.4 Update form validation
  - Ensure `name` and `private_key` are required
  - Remove validation for `public_key`
  - Update error messages if needed

- [x] 2.5 Update form labels and help text
  - Update private key label/placeholder to indicate public key is auto-extracted
  - Add tooltip or help text explaining public key auto-extraction
  - Ensure SSH user field has clear label

## 3. Testing

- [ ] 3.1 Test SSH key creation with private key only
  - Create key with RSA private key
  - Create key with Ed25519 private key
  - Create key with ECDSA private key
  - Verify public key is correctly extracted and stored

- [ ] 3.2 Test SSH key creation with encrypted private key
  - Create key with encrypted RSA private key + passphrase
  - Verify public key extraction works with passphrase

- [ ] 3.3 Test backward compatibility
  - Create key with both private and public key (should use provided public key)
  - Verify existing API clients still work

- [ ] 3.4 Test error handling
  - Invalid private key format → appropriate error message
  - Missing passphrase for encrypted key → clear error
  - Unsupported key type → informative error

- [ ] 3.5 Test form UI
  - Verify SSH user field appears after name field
  - Verify public key field is not shown
  - Verify form submission works correctly
  - Verify form validation works as expected

## 4. Documentation

- [ ] 4.1 Update API documentation (if exists)
  - Document that `public_key` is now optional
  - Document automatic public key extraction behavior

- [ ] 4.2 Update user documentation (if exists)
  - Document simplified SSH key creation process
  - Explain that only private key is required

## 5. Validation

- [ ] 5.1 Run backend tests
  - `go test ./internal/service/sshkey/...`

- [ ] 5.2 Run frontend build
  - `npm run build` (or equivalent)

- [ ] 5.3 Manual testing
  - Create SSH key via UI with private key only
  - Verify key appears in list with correct public key
  - Test SSH key in asset connection (if applicable)

- [ ] 5.4 Validate OpenSpec proposal
  - `openspec validate simplify-ssh-key-creation --strict`