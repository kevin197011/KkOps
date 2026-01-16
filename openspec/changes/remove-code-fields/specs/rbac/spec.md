## MODIFIED Requirements

### Requirement: Role Management
The system SHALL provide role management including creation, retrieval, update, and deletion of roles.

#### Scenario: Create role
- **WHEN** a user with appropriate permissions creates a new role with required information (name, description optional)
- **THEN** a new role record is created
- **AND** the role name must be unique
- **AND** the role does NOT require a code field
- **AND** the API response includes the name but NOT a code field

#### Scenario: Update role
- **WHEN** a user with appropriate permissions updates role information
- **THEN** the role record is updated
- **AND** the updated role name must remain unique
- **AND** the update does NOT involve a code field

### Requirement: Permission Management
The system SHALL provide permission management and checking based on resource and action.

#### Scenario: Check permission
- **WHEN** the system checks if a user has permission to perform an action on a resource
- **THEN** the system uses the combination of `resource` and `action` to identify the permission
- **AND** the system does NOT use permission codes
- **AND** permissions are identified uniquely by `resource` + `action` combination

#### Scenario: Require permission middleware
- **WHEN** a middleware requires a specific permission
- **THEN** the middleware accepts `resource` and `action` parameters instead of permission code
- **AND** the middleware checks user permissions based on resource/action combination
