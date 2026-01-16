## ADDED Requirements

### Requirement: SSH Key Management
The system SHALL support managing SSH keys for secure authentication to remote hosts.

#### Scenario: Create SSH key
- **WHEN** a user creates or imports an SSH key with name, type, public key, private key, and optional username
- **THEN** the private key is encrypted (AES-256) and stored securely
- **AND** the public key fingerprint is calculated and stored
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

### Requirement: WebSSH Terminal Connection
The system SHALL provide browser-based SSH terminal access to remote hosts.

#### Scenario: Establish SSH connection
- **WHEN** a user initiates an SSH connection with host, username, and authentication (password or SSH key)
- **THEN** the system establishes a WebSocket connection
- **AND** the system creates an SSH session to the target host
- **AND** the terminal session is ready for input/output

#### Scenario: Terminal input/output
- **WHEN** a user types commands in the terminal
- **THEN** input is sent to the remote host via SSH
- **AND** output from the remote host is displayed in the terminal
- **AND** the communication is bidirectional and real-time

#### Scenario: Terminal resize
- **WHEN** a user resizes the browser terminal window
- **THEN** the terminal PTY size is updated on the remote host
- **AND** applications on the remote host adjust to the new size

#### Scenario: Disconnect SSH session
- **WHEN** a user disconnects an SSH session
- **THEN** the SSH connection to the remote host is closed
- **AND** the WebSocket connection is closed
- **AND** session duration is recorded

### Requirement: SSH Session Management
The system SHALL support managing multiple SSH sessions with reconnection capability.

#### Scenario: Multiple sessions
- **WHEN** a user opens multiple terminal tabs
- **THEN** each tab maintains an independent SSH session
- **AND** sessions are tracked separately

#### Scenario: Session reconnection
- **WHEN** a user reconnects to a previously established session
- **THEN** the system attempts to resume the session if possible
- **OR** establishes a new session if the previous one is no longer available

### Requirement: File Transfer Support
The system SHALL support file transfer via SFTP and lrzsz (ZMODEM) protocols.

#### Scenario: SFTP file transfer
- **WHEN** a user uploads or downloads files via SFTP
- **THEN** the system uses the SSH connection to transfer files
- **AND** file transfer progress is displayed
- **AND** file transfer uses the same authentication as the SSH session

### Requirement: SSH Session Permission Control
The system SHALL enforce permission-based access control for SSH sessions.

#### Scenario: Permission check for SSH access
- **WHEN** a user attempts to connect to a host via SSH
- **THEN** the system verifies the user has permission to access that host
- **AND** if permission is denied, the connection is rejected
- **AND** permission is checked based on RBAC rules
