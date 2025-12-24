## REMOVED Requirements

### Requirement: SSH Connection Management
The system SHALL NOT provide SSH connection management functionality. SSH connections SHALL NOT be pre-configured or stored separately from host data.

#### Scenario: SSH connection CRUD operations removed
- **WHEN** a user attempts to create, list, update, or delete SSH connections
- **THEN** the system SHALL NOT provide these operations
- **AND** the system SHALL NOT maintain a separate SSH connection model or table
- **AND** the system SHALL use host data directly for SSH access

### Requirement: SSH Key Management
The system SHALL NOT provide SSH key management functionality. SSH keys SHALL NOT be stored or managed by the system.

#### Scenario: SSH key CRUD operations removed
- **WHEN** a user attempts to create, list, update, or delete SSH keys
- **THEN** the system SHALL NOT provide these operations
- **AND** the system SHALL NOT maintain a separate SSH key model or table
- **AND** the system SHALL NOT support key-based SSH authentication

#### Scenario: SSH key generation removed
- **WHEN** a user attempts to generate an SSH key pair
- **THEN** the system SHALL NOT provide this functionality
- **AND** the system SHALL NOT generate or store SSH keys

### Requirement: SSH Session Management
The system SHALL NOT provide SSH session management functionality. SSH sessions SHALL NOT be tracked or stored.

#### Scenario: SSH session CRUD operations removed
- **WHEN** a user attempts to create, list, view, or close SSH sessions
- **THEN** the system SHALL NOT provide these operations
- **AND** the system SHALL NOT maintain a separate SSH session model or table
- **AND** the system SHALL NOT track session history or metadata

#### Scenario: SSH session history removed
- **WHEN** a user attempts to view SSH session history
- **THEN** the system SHALL NOT provide this functionality
- **AND** the system SHALL NOT store session start/end times, duration, or other metadata

### Requirement: SSH Management Page
The system SHALL NOT provide a dedicated SSH management page. SSH-related functionality SHALL be integrated into host management.

#### Scenario: SSH management page removed
- **WHEN** a user attempts to access the SSH management page
- **THEN** the system SHALL NOT provide this page
- **AND** the system SHALL NOT display SSH-related menu items
- **AND** the system SHALL integrate terminal access into the host management page

