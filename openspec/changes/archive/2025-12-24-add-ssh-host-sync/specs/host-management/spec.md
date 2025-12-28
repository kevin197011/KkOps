## MODIFIED Requirements

### Requirement: Host SSH Port Configuration
The system SHALL provide the ability to configure SSH port for each host, which can be used when creating SSH connections.

#### Scenario: Configure SSH port for host
- **WHEN** a user creates or updates a host with an SSH port
- **THEN** the system SHALL store the SSH port value (default: 22)
- **AND** the system SHALL validate that the port is in the range 1-65535
- **AND** the system SHALL use this port when creating SSH connections for this host

#### Scenario: View host SSH port
- **WHEN** a user views host details
- **THEN** the system SHALL display the SSH port if configured
- **AND** the system SHALL display the default port (22) if not explicitly configured

