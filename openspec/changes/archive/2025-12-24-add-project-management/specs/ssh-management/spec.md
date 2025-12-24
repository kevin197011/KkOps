## MODIFIED Requirements

### Requirement: SSH Connection Management
The system SHALL provide the ability to manage SSH connections to hosts, including connection configuration and session management. All SSH connections MUST be associated with a project.

#### Scenario: Create SSH connection configuration
- **WHEN** a user creates an SSH connection configuration with host, port, username, authentication method, and project
- **THEN** the system SHALL validate that the project exists and the user has access to it
- **AND** the system SHALL store the connection configuration associated with the project
- **AND** the system SHALL encrypt sensitive information (passwords, keys) before storage

#### Scenario: Test SSH connection
- **WHEN** a user tests an SSH connection
- **THEN** the system SHALL attempt to establish an SSH connection
- **AND** the system SHALL return connection status (success/failure) and error message if failed
- **AND** the system SHALL validate that the user has access to the connection's project

#### Scenario: List SSH connections
- **WHEN** a user requests the list of SSH connections
- **THEN** the system SHALL return all configured SSH connections
- **AND** the system SHALL support filtering by project
- **AND** the system SHALL only return connections from projects the user has access to
- **AND** the system SHALL not return sensitive information (passwords, private keys)

### Requirement: SSH Terminal Session
The system SHALL provide the ability to establish interactive SSH terminal sessions through the web interface. Sessions are associated with the connection's project.

#### Scenario: Open SSH terminal session
- **WHEN** a user opens an SSH terminal session for a connection
- **THEN** the system SHALL validate that the user has access to the connection's project
- **AND** the system SHALL establish an SSH connection to the host
- **AND** the system SHALL create a WebSocket connection for terminal interaction
- **AND** the system SHALL stream terminal input/output in real-time

#### Scenario: Query SSH session history
- **WHEN** a user queries SSH session history
- **THEN** the system SHALL return a paginated list of past SSH sessions
- **AND** the system SHALL support filtering by project
- **AND** the system SHALL only return sessions from projects the user has access to
- **AND** the system SHALL include session start/end time, host, user, and duration

