## MODIFIED Requirements

### Requirement: SSH Key Management
The system SHALL support managing SSH keys for secure authentication to remote hosts.

#### Scenario: Create SSH key
- **WHEN** a user creates or imports an SSH key with name, private key, and optional SSH username
- **THEN** the system automatically extracts the public key from the private key if not provided
- **AND** the private key is encrypted (AES-256) and stored securely
- **AND** the public key fingerprint is calculated and stored
- **AND** the key is associated with the user
- **AND** the public key extraction supports RSA, Ed25519, ECDSA, and DSA key types
- **AND** if the private key is encrypted, the system requires a passphrase to extract the public key

#### Scenario: Create SSH key with provided public key (backward compatibility)
- **WHEN** a user creates an SSH key and provides both private and public keys
- **THEN** the system uses the provided public key
- **AND** validates that the provided public key matches the private key (optional validation)
- **AND** the private key is encrypted and stored securely
- **AND** the key is associated with the user

#### Scenario: List SSH keys
- **WHEN** a user requests their SSH key list
- **THEN** the system returns all keys belonging to that user
- **AND** private keys are never exposed
- **AND** key information includes name, type, fingerprint, username, and last used date

#### Scenario: Test SSH key
- **WHEN** a user tests an SSH key connection to a host
- **THEN** the system attempts to authenticate using the key
- **AND** returns connection success or failure status

#### Scenario: Delete SSH key
- **WHEN** a user deletes an SSH key
- **THEN** the system verifies the key is not in use
- **AND** the key record is removed from the database
- **AND** encrypted private key data is securely deleted

#### Scenario: SSH key creation form layout
- **WHEN** a user opens the SSH key creation form
- **THEN** the form displays fields in the following order: name, SSH user, private key, passphrase (if needed), description
- **AND** the public key field is not displayed (public key is auto-extracted)
- **AND** only name and private key are required fields