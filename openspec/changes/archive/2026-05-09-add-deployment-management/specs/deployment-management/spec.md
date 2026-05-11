## ADDED Requirements

### Requirement: Deployment Module Management

The system SHALL provide deployment module management functionality, allowing users to configure deployment modules under projects with version data sources and deployment scripts.

#### Scenario: Create deployment module

- **WHEN** user creates a deployment module with name, project, version source URL, and deploy script
- **THEN** the system SHALL save the module configuration
- **AND** the module SHALL be associated with the specified project

#### Scenario: List deployment modules

- **WHEN** user requests deployment module list
- **THEN** the system SHALL return all modules the user has access to
- **AND** support filtering by project ID

#### Scenario: Update deployment module

- **WHEN** user updates a deployment module's configuration
- **THEN** the system SHALL persist the changes

#### Scenario: Delete deployment module

- **WHEN** user deletes a deployment module
- **THEN** the system SHALL soft-delete the module
- **AND** preserve associated deployment history records

---

### Requirement: Version Source Integration

The system SHALL fetch available versions from configured HTTP JSON data sources.

#### Scenario: Fetch versions from HTTP source

- **WHEN** user requests versions for a deployment module
- **THEN** the system SHALL make HTTP GET request to the configured version_source_url
- **AND** parse the JSON response with format `{"versions": [...], "latest": "..."}`
- **AND** return the version list to the user

#### Scenario: Version source unavailable

- **WHEN** the version source URL is unreachable or returns invalid JSON
- **THEN** the system SHALL return an appropriate error message
- **AND** NOT block the user from manually specifying a version

---

### Requirement: Deployment Execution

The system SHALL execute deployments with selected version on target assets.

#### Scenario: Execute deployment with version

- **WHEN** user triggers deployment with a selected version and target assets
- **THEN** the system SHALL create a deployment record with status "pending"
- **AND** replace `${VERSION}`, `${MODULE_NAME}`, `${PROJECT_NAME}` variables in the deploy script
- **AND** execute the script on each target asset via SSH
- **AND** update deployment status to "running", then "success" or "failed"

#### Scenario: Deployment timeout

- **WHEN** deployment execution exceeds the configured timeout
- **THEN** the system SHALL terminate the execution
- **AND** mark the deployment as "failed" with timeout error

#### Scenario: Cancel deployment

- **WHEN** user cancels a running deployment
- **THEN** the system SHALL attempt to terminate the execution
- **AND** mark the deployment as "cancelled"

---

### Requirement: Deployment History

The system SHALL maintain a complete history of all deployment executions.

#### Scenario: View deployment history

- **WHEN** user requests deployment history
- **THEN** the system SHALL return deployment records sorted by creation time (newest first)
- **AND** support filtering by module ID

#### Scenario: View deployment details

- **WHEN** user requests deployment details
- **THEN** the system SHALL return the deployment record
- **AND** include execution output and error logs

---

### Requirement: Deployment UI

The system SHALL provide a user-friendly interface for deployment management.

#### Scenario: Deployment module list page

- **WHEN** user navigates to deployment management
- **THEN** the system SHALL display a list of deployment modules
- **AND** provide options to create, edit, delete modules
- **AND** provide a "Deploy" action button for each module

#### Scenario: Deployment execution modal

- **WHEN** user clicks "Deploy" on a module
- **THEN** the system SHALL display a modal with version dropdown
- **AND** fetch and display available versions from the configured source
- **AND** allow selection of target assets
- **AND** provide "Execute Deployment" button
