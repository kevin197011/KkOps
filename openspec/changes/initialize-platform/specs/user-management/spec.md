## ADDED Requirements

### Requirement: User Account Management
The system SHALL provide user account management functionality including creation, retrieval, update, and deletion of user accounts.

#### Scenario: Create user
- **WHEN** an administrator creates a new user with valid information (username, email, password, real_name, department)
- **THEN** a new user account is created and stored in the database with encrypted password
- **AND** the user account status is set to active by default

#### Scenario: List users
- **WHEN** a user with appropriate permissions requests the user list
- **THEN** the system returns a paginated list of users with their basic information
- **AND** the list supports filtering and searching

#### Scenario: Update user
- **WHEN** an administrator updates user information
- **THEN** the user record is updated in the database
- **AND** the updated_at timestamp is automatically updated

#### Scenario: Delete user
- **WHEN** an administrator deletes a user account
- **THEN** the user account is removed from the database
- **AND** associated data (roles, tokens) are handled according to data integrity rules

### Requirement: User Status Management
The system SHALL support enabling and disabling user accounts.

#### Scenario: Disable user
- **WHEN** an administrator disables a user account
- **THEN** the user status is set to disabled
- **AND** the user cannot log in or access the system

#### Scenario: Enable user
- **WHEN** an administrator enables a disabled user account
- **THEN** the user status is set to active
- **AND** the user can log in and access the system

### Requirement: Password Management
The system SHALL provide password management functionality including password reset and password change.

#### Scenario: Reset password
- **WHEN** an administrator resets a user's password
- **THEN** a new password is set for the user
- **AND** the password is stored using bcrypt hashing

#### Scenario: Change password
- **WHEN** a user changes their own password with current password verification
- **THEN** the new password is stored using bcrypt hashing
- **AND** the user must provide the current password for verification

### Requirement: Department Management
The system SHALL support organizational structure through departments with hierarchical relationships.

#### Scenario: Create department
- **WHEN** an administrator creates a department with name and optional parent department
- **THEN** a new department is created
- **AND** the department hierarchy is maintained

#### Scenario: Assign user to department
- **WHEN** a user is assigned to a department
- **THEN** the user's department_id is updated
- **AND** the user belongs to the organizational structure

### Requirement: API Token Management
The system SHALL allow users to generate and manage API tokens for programmatic API access.

#### Scenario: Generate API token
- **WHEN** a user generates a new API token with a name and optional expiration date
- **THEN** a unique token is generated and hashed for storage
- **AND** the full token is returned once for the user to save
- **AND** subsequent views show only a token prefix

#### Scenario: List API tokens
- **WHEN** a user requests their API token list
- **THEN** the system returns all tokens belonging to that user
- **AND** tokens show name, creation date, last used date, and status (not the full token)

#### Scenario: Revoke API token
- **WHEN** a user revokes or deletes an API token
- **THEN** the token status is set to disabled or the token is deleted
- **AND** the token can no longer be used for API authentication

#### Scenario: API token expiration
- **WHEN** an API token has an expiration date and that date is reached
- **THEN** the token is automatically invalidated
- **AND** API calls using that token are rejected
