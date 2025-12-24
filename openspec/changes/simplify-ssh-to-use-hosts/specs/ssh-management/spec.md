## REMOVED Requirements

### Requirement: SSH Connection Management
The system SHALL NOT maintain separate SSH connection records. SSH connections SHALL be managed directly through host data.

#### Scenario: SSH connection CRUD operations removed
- **WHEN** a user attempts to create, list, update, or delete SSH connections
- **THEN** the system SHALL NOT provide these operations
- **AND** the system SHALL use host data directly for SSH access

## MODIFIED Requirements

### Requirement: SSH Terminal Access
The system SHALL provide WebSSH terminal access directly through host management, using host data for connection configuration.

#### Scenario: Open terminal for a host
- **WHEN** a user selects a host from the SSH management page and clicks "打开终端"
- **THEN** the system SHALL open a WebSSH terminal for that host
- **AND** the system SHALL use the host's data for connection:
  - hostname: from host.hostname or host.ip_address
  - port: from host.ssh_port (default: 22)
  - username: from host.ssh_username (if configured) or prompt user
  - authentication: from host.ssh_key_id (if configured) or prompt for password
- **AND** the system SHALL establish the SSH connection using this configuration

#### Scenario: Configure SSH for a host
- **WHEN** a user configures SSH settings for a host
- **THEN** the system SHALL allow setting:
  - SSH username (optional, default can be prompted)
  - SSH key (optional, for key-based authentication)
- **AND** the system SHALL store this configuration in the Host record
- **AND** the system SHALL use this configuration when opening terminals

#### Scenario: SSH authentication with key
- **WHEN** a host has an SSH key configured (host.ssh_key_id)
- **THEN** the system SHALL use that SSH key for authentication
- **AND** the system SHALL NOT prompt for password
- **AND** the system SHALL establish the SSH connection using key authentication

#### Scenario: SSH authentication with password
- **WHEN** a host does not have an SSH key configured
- **THEN** the system SHALL prompt the user for password
- **AND** the system SHALL use the provided password for authentication
- **AND** the system SHALL NOT store the password (only use it for the session)

### Requirement: SSH Management Page
The system SHALL display hosts from host management instead of SSH connections, allowing direct terminal access for each host.

#### Scenario: View hosts in SSH management
- **WHEN** a user navigates to the SSH management page
- **THEN** the system SHALL display a list of hosts (from host management)
- **AND** the system SHALL show host information:
  - Host name
  - IP address
  - SSH port
  - SSH configuration status (has key/username configured)
- **AND** the system SHALL allow filtering by project
- **AND** the system SHALL allow opening terminal for any host

#### Scenario: Open terminal from host list
- **WHEN** a user clicks "打开终端" for a host
- **THEN** the system SHALL open a new terminal tab
- **AND** the system SHALL establish WebSSH connection using host data
- **AND** the system SHALL display the terminal with host name as tab title

### Requirement: SSH Session Management
The system SHALL manage SSH sessions with reference to hosts instead of connections.

#### Scenario: Create SSH session
- **WHEN** a user opens a WebSSH terminal for a host
- **THEN** the system SHALL create an SSHSession record with:
  - host_id: reference to the host
  - user_id: current user
  - session_token: unique session identifier
  - status: "active"
- **AND** the system SHALL NOT reference a connection

#### Scenario: List SSH sessions
- **WHEN** a user views SSH sessions
- **THEN** the system SHALL display sessions with host information
- **AND** the system SHALL show:
  - Host name (from host.hostname)
  - Host IP address (from host.ip_address)
  - Session start time
  - Session status
- **AND** the system SHALL allow filtering by host

#### Scenario: Close SSH session
- **WHEN** a user closes a terminal or session
- **THEN** the system SHALL update the SSHSession record:
  - status: "closed"
  - end_time: current time
  - duration: calculated
- **AND** the system SHALL close the SSH connection
- **AND** the system SHALL clean up resources

## ADDED Requirements

### Requirement: Host SSH Configuration
The system SHALL allow configuring SSH settings per host, including default username and SSH key.

#### Scenario: Configure SSH username for host
- **WHEN** a user configures SSH settings for a host
- **THEN** the system SHALL allow setting a default SSH username
- **AND** the system SHALL store this in host.ssh_username
- **AND** the system SHALL use this username when opening terminals (unless overridden)

#### Scenario: Configure SSH key for host
- **WHEN** a user configures SSH settings for a host
- **THEN** the system SHALL allow selecting an SSH key from available keys
- **AND** the system SHALL store the key ID in host.ssh_key_id
- **AND** the system SHALL use this key for authentication when opening terminals

#### Scenario: Use default SSH port
- **WHEN** a host does not have ssh_port configured
- **THEN** the system SHALL use port 22 as default
- **AND** the system SHALL allow connection using default port

