## ADDED Requirements

### Requirement: User Role Management
The system SHALL allow administrators to manage user-role assignments through the UI.

#### Scenario: View user roles
- **WHEN** administrator views the user list
- **THEN** each user's assigned roles SHALL be displayed as tags

#### Scenario: Assign roles to user
- **WHEN** administrator clicks "角色" button for a user
- **THEN** a modal SHALL display all available roles with checkboxes
- **AND** currently assigned roles SHALL be pre-selected
- **WHEN** administrator selects/deselects roles and clicks save
- **THEN** the user's roles SHALL be updated accordingly

#### Scenario: Remove role from user
- **WHEN** administrator deselects a role in the role assignment modal
- **AND** clicks save
- **THEN** the role SHALL be removed from the user

### Requirement: Role Asset Authorization
The system SHALL allow administrators to authorize assets for non-admin roles.

#### Scenario: View role authorized assets
- **WHEN** administrator views the role list
- **THEN** each role SHALL display the count of authorized assets
- **AND** admin roles SHALL display "全部" (all) instead of count

#### Scenario: Authorize assets for role
- **WHEN** administrator clicks "授权" button for a non-admin role
- **THEN** a modal SHALL display a Transfer component for asset selection
- **AND** currently authorized assets SHALL be in the right panel
- **WHEN** administrator moves assets and clicks save
- **THEN** the role's authorized assets SHALL be updated

#### Scenario: Batch authorize by project/environment
- **WHEN** administrator selects a project and/or environment in the authorization modal
- **AND** clicks "添加到授权"
- **THEN** all matching assets SHALL be added to the authorized list

### Requirement: Combined Authorization Check
The system SHALL check both direct user authorization and role-based authorization when determining asset access.

#### Scenario: Admin role access
- **WHEN** a user has a role with is_admin=true
- **THEN** the user SHALL have access to all assets

#### Scenario: Direct user authorization
- **WHEN** a user has a direct authorization record in user_assets
- **THEN** the user SHALL have access to that asset

#### Scenario: Role-based authorization
- **WHEN** a user has a role that is authorized for an asset (via role_assets)
- **THEN** the user SHALL have access to that asset

#### Scenario: Combined authorization
- **WHEN** a user has both direct and role-based authorization for different assets
- **THEN** the user SHALL have access to all authorized assets from both sources

#### Scenario: No authorization
- **WHEN** a user has no admin role, no direct authorization, and no role-based authorization for an asset
- **THEN** the user SHALL be denied access to that asset

### Requirement: Role Assets API
The system SHALL provide REST API endpoints for managing role-asset authorizations.

#### Scenario: Get role assets
- **WHEN** GET /api/v1/roles/:id/assets is called
- **THEN** the response SHALL contain the list of assets authorized for the role

#### Scenario: Grant role assets
- **WHEN** POST /api/v1/roles/:id/assets is called with asset_ids
- **THEN** the specified assets SHALL be authorized for the role

#### Scenario: Revoke role asset
- **WHEN** DELETE /api/v1/roles/:id/assets/:asset_id is called
- **THEN** the authorization for that asset SHALL be removed from the role

### Requirement: User Roles API
The system SHALL provide REST API endpoints for managing user-role assignments.

#### Scenario: Get user roles
- **WHEN** GET /api/v1/users/:id/roles is called
- **THEN** the response SHALL contain the list of roles assigned to the user

#### Scenario: Assign user roles
- **WHEN** POST /api/v1/users/:id/roles is called with role_ids
- **THEN** the specified roles SHALL be assigned to the user

#### Scenario: Remove user role
- **WHEN** DELETE /api/v1/users/:id/roles/:role_id is called
- **THEN** the role SHALL be removed from the user
