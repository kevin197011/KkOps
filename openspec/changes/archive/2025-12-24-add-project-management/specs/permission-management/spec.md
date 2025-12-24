## MODIFIED Requirements

### Requirement: Project-Level Permission Control
The system SHALL provide the ability to control user access to projects and resources within projects.

#### Scenario: Assign project role to user
- **WHEN** a user assigns a role to another user for a specific project
- **THEN** the system SHALL validate that the user has permission to manage project members
- **AND** the system SHALL associate the user with the project and role
- **AND** the system SHALL grant the user access to resources within that project based on the role

#### Scenario: Query user's project access
- **WHEN** a user queries their project access
- **THEN** the system SHALL return a list of all projects the user has access to
- **AND** the system SHALL include the user's role in each project

#### Scenario: Check project resource access
- **WHEN** a user attempts to access a resource (host, deployment, etc.)
- **THEN** the system SHALL check if the user has access to the resource's project
- **AND** the system SHALL check if the user has the required permission for the operation
- **AND** the system SHALL allow the operation only if both checks pass

