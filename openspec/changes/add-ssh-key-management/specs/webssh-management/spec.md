## MODIFIED Requirements

### Requirement: Terminal Authentication
The system SHALL support both SSH key-based and password-based authentication for WebSSH terminal connections.

#### Scenario: Key-based terminal authentication
- **WHEN** a user opens a WebSSH terminal and selects key-based authentication
- **AND** the user selects an SSH key from their available keys
- **THEN** the system SHALL use the selected SSH key for authentication
- **AND** the system SHALL NOT prompt for password
- **AND** the system SHALL establish the SSH connection using the key

#### Scenario: Password-based terminal authentication
- **WHEN** a user opens a WebSSH terminal and selects password-based authentication
- **THEN** the system SHALL prompt for username and password
- **AND** the system SHALL establish the SSH connection using the provided credentials
- **AND** the system SHALL work as before (existing password authentication flow)

#### Scenario: Authentication method selection
- **WHEN** a user clicks "打开终端" for a host
- **THEN** the system SHALL show a dialog to select authentication method (key or password)
- **AND** if key method is selected, the system SHALL show SSH key selection dropdown
- **AND** if password method is selected, the system SHALL proceed with password prompt
- **AND** the system SHALL remember the user's choice for the current session (optional)

### Requirement: Terminal Integration with Host Management
The system SHALL provide a dedicated WebSSH management page for terminal access, separate from host CRUD operations.

#### Scenario: Terminal access from WebSSH management page
- **WHEN** a user navigates to the WebSSH management page
- **THEN** the system SHALL display a list of all hosts
- **AND** the system SHALL display an "打开终端" button for each host in the operations column
- **AND** the system SHALL allow opening a terminal for any host in the list
- **AND** the system SHALL support opening multiple terminals for different hosts
- **AND** the system SHALL provide SSH key management in a separate tab

#### Scenario: SSH key management in WebSSH page
- **WHEN** a user navigates to the WebSSH management page
- **THEN** the system SHALL display a "SSH密钥" tab alongside the "主机列表" tab
- **AND** the system SHALL allow managing SSH keys (upload, view, delete) in this tab
- **AND** the system SHALL display all SSH keys belonging to the current user

