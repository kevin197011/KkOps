## MODIFIED Requirements

### Requirement: Environment Management
The system SHALL provide environment management including creation, retrieval, update, and deletion of environments.

#### Scenario: Create environment
- **WHEN** a user with appropriate permissions creates a new environment with required information (name, description optional)
- **THEN** a new environment record is created
- **AND** the environment name must be unique
- **AND** the environment does NOT require a code field
- **AND** the API response includes the name but NOT a code field

#### Scenario: List environments
- **WHEN** a user requests the environment list
- **THEN** the system returns a list of environments
- **AND** each environment includes the `name` field in the JSON response
- **AND** the list does NOT include a code field

#### Scenario: Update environment
- **WHEN** a user with appropriate permissions updates environment information
- **THEN** the environment record is updated
- **AND** the updated environment name must remain unique
- **AND** the update does NOT involve a code field
