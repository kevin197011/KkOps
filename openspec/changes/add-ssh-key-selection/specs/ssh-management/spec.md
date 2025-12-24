## MODIFIED Requirements

### Requirement: SSH Connection Management
Users SHALL be able to create, view, edit, and delete SSH connections. When creating or editing an SSH connection, users MUST specify:
- Project association
- Connection name
- Host information (hostname/IP, port)
- Username
- Authentication method (password or key)
- **When authentication method is "key", users MUST select a configured SSH key from the available keys list**
- Connection status

The system SHALL display SSH connection information including:
- Connection details (name, hostname, port, username)
- Authentication type
- **Associated SSH key name (when key authentication is used)**
- Connection status
- Last connection time

#### Scenario: Create SSH connection with key authentication
- **WHEN** user selects "key" as authentication type in SSH connection form
- **THEN** a dropdown list of available SSH keys is displayed
- **AND** user must select a key from the list
- **WHEN** user submits the form without selecting a key
- **THEN** validation error is shown requiring key selection
- **WHEN** user selects a key and submits the form
- **THEN** SSH connection is created with the selected key associated

#### Scenario: Edit SSH connection to change key
- **WHEN** user edits an existing SSH connection with key authentication
- **THEN** the currently associated key is pre-selected in the key dropdown
- **WHEN** user changes the selected key
- **THEN** the connection is updated with the new key

#### Scenario: View SSH connection with key information
- **WHEN** user views the SSH connection list
- **THEN** connections using key authentication display the associated key name
- **AND** connections using password authentication show "密码" authentication type

#### Scenario: Switch authentication type
- **WHEN** user changes authentication type from "password" to "key"
- **THEN** key selection dropdown appears
- **AND** key selection becomes required
- **WHEN** user changes authentication type from "key" to "password"
- **THEN** key selection dropdown is hidden
- **AND** key selection requirement is removed

