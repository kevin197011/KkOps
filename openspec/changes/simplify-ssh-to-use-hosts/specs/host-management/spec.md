## MODIFIED Requirements

### Requirement: Host SSH Configuration
The system SHALL support optional SSH configuration fields for hosts, allowing per-host SSH username and key settings.

#### Scenario: Configure SSH username for host
- **WHEN** a user creates or updates a host
- **THEN** the system SHALL allow setting an optional SSH username field
- **AND** the system SHALL store this in the host record
- **AND** the system SHALL use this username for SSH connections to this host

#### Scenario: Configure SSH key for host
- **WHEN** a user creates or updates a host
- **THEN** the system SHALL allow selecting an optional SSH key
- **AND** the system SHALL store the key ID in the host record
- **AND** the system SHALL use this key for SSH authentication to this host

#### Scenario: Display SSH configuration in host list
- **WHEN** a user views the host list
- **THEN** the system SHALL optionally display SSH configuration status
- **AND** the system SHALL indicate if host has SSH username or key configured
- **AND** the system SHALL allow quick access to SSH terminal

