## MODIFIED Requirements

### Requirement: SSH Connection Host Selection
The system SHALL provide the ability to select a host from host management when creating SSH connections, automatically syncing host information.

#### Scenario: Create SSH connection from host selection
- **WHEN** a user creates an SSH connection and selects a host from the host management list
- **THEN** the system SHALL automatically populate the hostname, IP address, and SSH port from the selected host
- **AND** the system SHALL associate the SSH connection with the selected host (via host_id)
- **AND** the system SHALL allow the user to modify the auto-filled information if needed

#### Scenario: Create SSH connection manually
- **WHEN** a user creates an SSH connection without selecting a host
- **THEN** the system SHALL allow manual input of hostname, IP address, and port
- **AND** the system SHALL not require host selection (host_id can be null)

#### Scenario: Update SSH connection with host selection
- **WHEN** a user updates an SSH connection and changes the selected host
- **THEN** the system SHALL update the hostname, IP address, and SSH port from the newly selected host
- **AND** the system SHALL preserve other connection settings (username, auth type, etc.)

