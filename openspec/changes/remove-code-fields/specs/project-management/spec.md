## MODIFIED Requirements

### Requirement: Project Management
The system SHALL provide project management including creation, retrieval, update, and deletion of projects.

#### Scenario: Create project
- **WHEN** a user with appropriate permissions creates a new project with required information (name, description optional)
- **THEN** a new project record is created
- **AND** the project name must be unique
- **AND** the project does NOT require a code field
- **AND** the API response includes the name but NOT a code field

#### Scenario: List projects
- **WHEN** a user requests the project list
- **THEN** the system returns a list of projects
- **AND** each project includes the `name` field in the JSON response
- **AND** the list does NOT include a code field

#### Scenario: Update project
- **WHEN** a user with appropriate permissions updates project information
- **THEN** the project record is updated
- **AND** the updated project name must remain unique
- **AND** the update does NOT involve a code field
