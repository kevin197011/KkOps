## REMOVED Requirements

### Requirement: Direct User Asset Authorization
**Reason**: Simplifying to role-based authorization only. Users now get asset access through their assigned roles.
**Migration**: Existing user_assets data should be migrated to role_assets before removal.

## MODIFIED Requirements

### Requirement: Combined Authorization Check
The system SHALL check user authorization through roles only when determining asset access.

#### Scenario: Admin role access
- **WHEN** a user has a role with is_admin=true
- **THEN** the user SHALL have access to all assets

#### Scenario: Role-based authorization
- **WHEN** a user has a role that is authorized for an asset (via role_assets)
- **THEN** the user SHALL have access to that asset

#### Scenario: Multiple roles authorization
- **WHEN** a user has multiple roles with different asset authorizations
- **THEN** the user SHALL have access to the union of all authorized assets from all roles

#### Scenario: No authorization
- **WHEN** a user has no admin role and no role with authorization for an asset
- **THEN** the user SHALL be denied access to that asset

### Requirement: User Management Simplified
The system SHALL only provide role assignment functionality in user management, without direct asset authorization.

#### Scenario: View user roles
- **WHEN** administrator views the user list
- **THEN** each user's assigned roles SHALL be displayed as tags

#### Scenario: Assign roles to user
- **WHEN** administrator clicks "角色" button for a user
- **THEN** a modal SHALL display all available roles with checkboxes
- **AND** currently assigned roles SHALL be pre-selected
- **WHEN** administrator selects/deselects roles and clicks save
- **THEN** the user's roles SHALL be updated accordingly

#### Scenario: No direct asset authorization
- **WHEN** administrator views user management page
- **THEN** there SHALL be no "授权" button for direct asset assignment
- **AND** there SHALL be no asset authorization modal
