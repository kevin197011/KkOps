## ADDED Requirements

### Requirement: Automatic SSH Connection Generation
The system SHALL automatically generate SSH connections for Linux/Unix hosts based on host management data, eliminating the need for manual connection creation.

#### Scenario: Auto-generate connections for all eligible hosts
- **WHEN** a user triggers auto-generation for all projects
- **THEN** the system SHALL identify all hosts with Linux/Unix OS types (e.g., "Linux", "Ubuntu", "CentOS", "RedHat", "Debian", "FreeBSD", "OpenBSD")
- **AND** the system SHALL skip hosts that already have SSH connections
- **AND** the system SHALL create SSH connections for eligible hosts using:
  - hostname from host.IPAddress or host.Hostname
  - port from host.SSHPort (default: 22)
  - project_id from host.ProjectID
  - host_id from host.ID
  - configurable default username (e.g., "root", "ubuntu", "admin")
  - configurable default authentication type ("password" or "key")
  - auto-generated name (e.g., "Auto: {hostname}:{port}")
- **AND** the system SHALL return a summary of created and skipped connections

#### Scenario: Auto-generate connections for a specific project
- **WHEN** a user triggers auto-generation for a specific project
- **THEN** the system SHALL only process hosts belonging to that project
- **AND** the system SHALL apply the same generation logic as for all projects
- **AND** the system SHALL return project-scoped results

#### Scenario: Auto-generate connection for a single host
- **WHEN** a user triggers auto-generation for a specific host
- **THEN** the system SHALL check if the host is Linux/Unix
- **AND** the system SHALL check if the host already has an SSH connection
- **AND** if eligible, the system SHALL create an SSH connection for that host
- **AND** if not eligible (wrong OS or connection exists), the system SHALL return an appropriate error

#### Scenario: Skip hosts with existing connections
- **WHEN** auto-generation is triggered
- **THEN** the system SHALL check for existing SSH connections by host_id
- **AND** the system SHALL skip hosts that already have connections
- **AND** the system SHALL include skipped hosts in the summary response

#### Scenario: Validation of host requirements
- **WHEN** auto-generation processes a host
- **THEN** the system SHALL validate that the host has:
  - Valid IP address or hostname
  - Valid SSH port (1-65535)
  - Linux/Unix OS type
- **AND** if validation fails, the system SHALL skip that host and log the reason

#### Scenario: Configure default username and authentication
- **WHEN** a user triggers auto-generation
- **THEN** the system SHALL allow configuration of:
  - Default username (optional, default: "root")
  - Default authentication type (optional, default: "password")
- **AND** the system SHALL apply these defaults to all generated connections

## MODIFIED Requirements

### Requirement: SSH Connection Management
The system SHALL provide comprehensive SSH connection management capabilities, including manual creation, automatic generation, and connection lifecycle management.

#### Scenario: Display auto-generated connections
- **WHEN** a user views the SSH connection list
- **THEN** the system SHALL display all connections (manual and auto-generated)
- **AND** the system SHALL indicate which connections were auto-generated (e.g., with a badge or icon)
- **AND** the system SHALL allow editing and deletion of auto-generated connections

#### Scenario: Manual override of auto-generated connections
- **WHEN** a user edits an auto-generated SSH connection
- **THEN** the system SHALL allow modification of all connection fields
- **AND** the system SHALL preserve the connection's auto-generated status indicator
- **AND** the system SHALL allow the user to change authentication credentials

#### Scenario: Trigger auto-generation from UI
- **WHEN** a user navigates to the SSH management page
- **THEN** the system SHALL provide an "自动生成 SSH 连接" button
- **AND** the system SHALL show a modal with configuration options:
  - Project filter (optional, default: all projects)
  - Default username (optional, default: "root")
  - Default authentication type (optional, default: "password")
  - Preview of eligible hosts count
- **AND** upon confirmation, the system SHALL trigger auto-generation
- **AND** the system SHALL display results (created, skipped) and refresh the connection list

#### Scenario: Quick generate for single host
- **WHEN** a user views a Linux/Unix host in host management
- **THEN** the system SHALL provide a "生成 SSH 连接" action
- **AND** when clicked, the system SHALL generate an SSH connection for that host
- **AND** if a connection already exists, the system SHALL show an appropriate message

