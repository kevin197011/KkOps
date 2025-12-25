## MODIFIED Requirements

### Requirement: SSH Key Management UI Location
The system SHALL provide SSH key management functionality on a dedicated page, separate from WebSSH terminal access.

#### Scenario: Access SSH key management page
- **WHEN** a user navigates to the SSH Keys management page
- **THEN** the system SHALL display a dedicated page for SSH key management
- **AND** the system SHALL provide access via menu item "SSH密钥管理"
- **AND** the system SHALL provide access via route `/ssh-keys`
- **AND** the system SHALL NOT require navigation through WebSSH page

#### Scenario: SSH key management not in WebSSH page
- **WHEN** a user navigates to the WebSSH management page
- **THEN** the system SHALL NOT display SSH key management tab or UI
- **AND** the system SHALL display only host list and terminal access functionality
- **AND** the system SHALL provide a way to navigate to SSH key management (via menu)

### Requirement: SSH Key Management (Unchanged Core Functionality)
The system SHALL provide the ability to manage SSH keys for WebSSH authentication, including uploading, viewing, editing, and deleting SSH keys.

#### Scenario: Upload SSH private key
- **WHEN** a user uploads an SSH private key in the SSH Keys management page
- **THEN** the system SHALL validate the key format (RSA, ED25519, etc.)
- **AND** the system SHALL extract the public key and calculate fingerprint
- **AND** the system SHALL encrypt and store the private key securely
- **AND** the system SHALL associate the key with the current user
- **AND** the system SHALL display the key in the SSH key list

#### Scenario: List SSH keys
- **WHEN** a user views the SSH Keys management page
- **THEN** the system SHALL display all SSH keys belonging to the current user
- **AND** the system SHALL display key metadata: name, type, fingerprint, username, created date
- **AND** the system SHALL NOT display the private key content
- **AND** the system SHALL allow editing key name and username
- **AND** the system SHALL allow deleting keys

#### Scenario: Delete SSH key
- **WHEN** a user deletes an SSH key from the SSH Keys management page
- **THEN** the system SHALL remove the key from storage
- **AND** the system SHALL prevent deletion if the key is configured as default for any host (warn user)
- **AND** the system SHALL update the SSH key list

#### Scenario: SSH key security (Unchanged)
- **WHEN** an SSH key is stored in the system
- **THEN** the system SHALL encrypt the private key before storage
- **AND** the system SHALL only allow the key owner to access their keys
- **AND** the system SHALL NOT expose private keys in API responses

### Requirement: SSH Key Selection in Terminal Connection (Unchanged)
The system SHALL allow users to select an SSH key when opening a WebSSH terminal connection.

#### Scenario: Select SSH key for terminal connection
- **WHEN** a user opens a terminal for a host in WebSSH page
- **THEN** the system SHALL show authentication method selection (key or password)
- **AND** if key method is selected, the system SHALL display a list of available SSH keys
- **AND** the system SHALL allow the user to select a key from the list
- **AND** the system SHALL use the selected key for SSH authentication
- **AND** if no key is configured for the host, the system SHALL show all user's keys

#### Scenario: Auto-select default key (Unchanged)
- **WHEN** a host has a default SSH key configured
- **AND** a user opens a terminal for that host
- **THEN** the system SHALL pre-select the default key in the key selection
- **AND** the system SHALL allow the user to override and select a different key or use password

#### Scenario: Key authentication fallback (Unchanged)
- **WHEN** key-based authentication fails during terminal connection
- **THEN** the system SHALL display an error message
- **AND** the system SHALL allow the user to retry with a different key or switch to password authentication

