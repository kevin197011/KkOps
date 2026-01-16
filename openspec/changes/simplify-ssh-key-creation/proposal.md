# Change: Simplify SSH Key Creation - Private Key Only

## Why

Currently, when creating an SSH key, users must provide both the private key and public key separately. This creates an unnecessary burden because:

- **Redundant input**: Public keys can be automatically extracted from private keys
- **User confusion**: Users may not understand why both keys are needed
- **Error-prone**: Users might provide mismatched key pairs
- **Poor UX**: The form requires more fields than necessary

Additionally, the SSH user field is positioned after other fields, which is not intuitive since SSH user is a fundamental property that should be configured early in the form.

Simplifying the SSH key creation process to require only the private key (and automatically extract the public key) will:
- Reduce user effort and potential errors
- Provide a more streamlined and intuitive user experience
- Align with common SSH key management practices
- Make the form more focused on essential information

## What Changes

- **MODIFIED**: SSH key creation form to require only private key input (public key auto-extracted)
- **MODIFIED**: Remove public key input field from creation form
- **MODIFIED**: Move SSH user field to appear immediately after name field
- **MODIFIED**: Backend service to automatically extract public key from private key during creation
- **MAINTAINED**: Backend API continues to accept public_key (for backward compatibility), but it becomes optional
- **MAINTAINED**: Edit form behavior (only name, SSH user, description can be edited)

## Impact

- **UI Changes**: Creation form becomes simpler with fewer required fields
- **Backend Changes**: Add logic to extract public key from private key if not provided
- **User Experience**: Faster and easier SSH key creation
- **Data Quality**: Reduced risk of mismatched key pairs

## Field Order

### New Form Layout (Create)
1. **Name** (required)
2. **SSH User** (optional) - moved here for better UX
3. **Private Key** (required) - public key auto-extracted
4. **Private Key Passphrase** (optional) - if key is encrypted
5. **Description** (optional)

### Edit Form (unchanged)
- Name
- SSH User
- Description

## Technical Considerations

- **Public Key Extraction**: Use Go's `golang.org/x/crypto/ssh` package to parse private key and extract public key
- **Error Handling**: Provide clear error messages if private key format is invalid
- **Backward Compatibility**: API still accepts `public_key` parameter but treats it as optional
- **Validation**: Validate private key format and extract public key before storing