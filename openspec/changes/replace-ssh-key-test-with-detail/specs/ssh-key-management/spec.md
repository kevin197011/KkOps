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

#### Scenario: View SSH key details
- **WHEN** a user clicks the "详情" (Detail) button for an SSH key in the list
- **THEN** a modal displays complete key information including:
  - ID, name, type, SSH user
  - Fingerprint and public key (full content)
  - Description
  - Last used time, created time, updated time
- **AND** the public key can be copied to clipboard with a single click
- **AND** all information is displayed in a structured, readable format

#### Scenario: Test SSH key
- **WHEN** a user tests an SSH key connection to a host
- **THEN** the system attempts to authenticate using the key
- **AND** returns connection success or failure status
- **Note**: Test functionality remains available via API but is not exposed in the main UI list operations

#### Scenario: Delete SSH key
- **WHEN** a user deletes an SSH key
- **THEN** the system verifies the key is not in use
- **AND** the key record is removed from the database
- **AND** encrypted private key data is securely deleted