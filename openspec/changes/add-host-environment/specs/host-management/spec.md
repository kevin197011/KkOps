## MODIFIED Requirements

### Requirement: Host Management
The system SHALL provide the ability to manage hosts, including host configuration, status monitoring, and organization. Host information SHALL include:
- Host ID
- **Project association (displayed in table after ID column)**
- Hostname
- IP address
- **Environment classification (dev, uat, staging, prod, test, etc.)**
- Salt Minion ID
- Operating system information
- Hardware specifications (CPU, memory, disk)
- Status (online, offline, unknown)
- SSH port
- Host groups and tags
- Last seen timestamp

#### Scenario: View host list with project and environment
- **WHEN** user views the host management page
- **THEN** the table displays columns in order: ID, Project, Hostname, IP Address, Environment, and other fields
- **AND** the Project column shows the project name for each host
- **AND** the Environment column shows the environment tag (dev, uat, prod, etc.) for each host

#### Scenario: Create host with environment
- **WHEN** user creates a new host
- **THEN** the form includes an "环境" (Environment) field with dropdown options (dev, uat, staging, prod, test)
- **WHEN** user selects an environment and submits the form
- **THEN** the host is created with the selected environment value
- **AND** the environment is displayed in the host list table

#### Scenario: Edit host environment
- **WHEN** user edits an existing host
- **THEN** the environment field is pre-populated with the current environment value
- **WHEN** user changes the environment and saves
- **THEN** the host environment is updated
- **AND** the updated environment is displayed in the host list table

#### Scenario: Filter hosts by environment
- **WHEN** user wants to filter hosts by environment
- **THEN** an environment filter option is available (optional enhancement)
- **WHEN** user selects an environment filter
- **THEN** only hosts matching that environment are displayed

