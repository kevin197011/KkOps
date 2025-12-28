## MODIFIED Requirements

### Requirement: Role Management UI
The system SHALL provide a user interface for managing roles including creation, update, deletion, and permission assignment.

#### Scenario: View role list
- **WHEN** a user accesses the role management page
- **THEN** the system SHALL display a paginated list of all roles
- **AND** the system SHALL show role name, description, and permission count for each role
- **AND** the system SHALL provide actions to edit, delete, and manage permissions for each role

#### Scenario: Create a new role
- **WHEN** a user creates a new role through the UI
- **THEN** the system SHALL display a form for entering role name and description
- **AND** the system SHALL validate that the role name is unique
- **AND** the system SHALL create the role and return to the role list

#### Scenario: Edit a role
- **WHEN** a user edits a role through the UI
- **THEN** the system SHALL display a form pre-filled with the role's current information
- **AND** the system SHALL allow updating the role name and description
- **AND** the system SHALL validate that the updated role name is unique (if changed)

#### Scenario: Assign permissions to a role
- **WHEN** a user manages permissions for a role
- **THEN** the system SHALL display all available permissions grouped by resource type
- **AND** the system SHALL show which permissions are currently assigned to the role
- **AND** the system SHALL allow the user to check/uncheck permissions
- **AND** the system SHALL immediately update the role's permissions when changes are made

#### Scenario: Delete a role
- **WHEN** a user deletes a role through the UI
- **THEN** the system SHALL display a confirmation dialog
- **AND** the system SHALL check if any users are assigned to the role
- **AND** the system SHALL prevent deletion if users are assigned
- **AND** the system SHALL remove the role if deletion is confirmed and no users are assigned

### Requirement: Permission Management UI
The system SHALL provide a user interface for viewing and managing permissions.

#### Scenario: View permission list
- **WHEN** a user accesses the permission management page
- **THEN** the system SHALL display a paginated list of all permissions
- **AND** the system SHALL show permission code, name, resource type, action, and description
- **AND** the system SHALL support filtering by resource type
- **AND** the system SHALL support searching by permission name, code, or description

#### Scenario: Filter permissions by resource type
- **WHEN** a user filters permissions by resource type
- **THEN** the system SHALL display only permissions matching the selected resource type
- **AND** the system SHALL update the list dynamically

#### Scenario: Search permissions
- **WHEN** a user searches for permissions
- **THEN** the system SHALL search across permission name, code, and description
- **AND** the system SHALL display matching permissions
- **AND** the system SHALL support case-insensitive search

