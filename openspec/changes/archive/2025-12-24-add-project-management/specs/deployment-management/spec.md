## MODIFIED Requirements

### Requirement: Deployment Configuration
The system SHALL provide the ability to create and manage deployment configurations that define how applications are deployed to hosts. All deployment configurations MUST be associated with a project.

#### Scenario: Create deployment configuration
- **WHEN** a user creates a deployment configuration with application name, version, target hosts, Salt state files, and project
- **THEN** the system SHALL validate that the project exists and the user has access to it
- **AND** the system SHALL validate that all target hosts belong to the same project
- **AND** the system SHALL store the configuration associated with the project
- **AND** the system SHALL validate that the Salt state files exist
- **AND** the system SHALL return the created configuration with a unique ID

#### Scenario: List deployment configurations
- **WHEN** a user requests the list of deployment configurations
- **THEN** the system SHALL return all configurations
- **AND** the system SHALL support filtering by application name, status, and project
- **AND** the system SHALL only return configurations from projects the user has access to

### Requirement: Deployment Execution
The system SHALL provide the ability to execute deployments using Salt to target hosts. Deployments inherit the project from their configuration.

#### Scenario: Start a deployment
- **WHEN** a user starts a deployment with a configuration and target hosts
- **THEN** the system SHALL validate that the user has access to the configuration's project
- **AND** the system SHALL validate that all target hosts belong to the same project
- **AND** the system SHALL create a deployment record with status "in_progress"
- **AND** the system SHALL execute Salt state.apply on target hosts
- **AND** the system SHALL track the deployment progress
- **AND** the system SHALL record the deployment in audit logs

#### Scenario: Query deployment history
- **WHEN** a user queries deployment history
- **THEN** the system SHALL return a list of all deployments ordered by deployment time
- **AND** the system SHALL support filtering by version, status, date range, and project
- **AND** the system SHALL only return deployments from projects the user has access to
- **AND** the system SHALL include version, status, target hosts, and execution time

