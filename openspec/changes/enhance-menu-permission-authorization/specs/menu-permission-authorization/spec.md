# Menu Permission Authorization Specification

## ADDED Requirements

### Requirement: Menu Permission Definition
The system SHALL define permissions for all menu functions using resource/action format.

#### Scenario: Define menu permissions
- **WHEN** the system initializes
- **THEN** it SHALL create permissions for all menu functions:
  - `dashboard:read` - 仪表板查看
  - `projects:*` - 项目管理所有操作
  - `environments:*` - 环境管理所有操作
  - `cloud-platforms:*` - 云平台管理所有操作
  - `assets:*` - 资产管理所有操作（菜单权限）
  - `executions:*` - 运维执行所有操作
  - `templates:*` - 任务模板所有操作
  - `tasks:*` - 任务执行所有操作
  - `deployments:*` - 部署管理所有操作
  - `ssh-keys:*` - SSH 密钥管理所有操作
  - `users:*` - 用户管理所有操作
  - `roles:*` - 角色权限管理所有操作
  - `audit-logs:read` - 审计日志查看
- **AND** SHALL use format `<resource>:<action>` where `*` represents all actions

### Requirement: Permission Initialization
The system SHALL automatically initialize menu permissions during database setup.

#### Scenario: Initialize permissions on first startup
- **WHEN** the system starts for the first time
- **OR** when permissions are missing
- **THEN** it SHALL create all required menu permissions in the database
- **AND** SHALL ensure idempotent execution (no error on repeated execution)

#### Scenario: Permission initialization is idempotent
- **WHEN** permission initialization is executed multiple times
- **THEN** it SHALL not create duplicate permissions
- **AND** SHALL not return errors

### Requirement: Admin Permission Override
Administrator users SHALL automatically have access to all menu functions without explicit permission assignment.

#### Scenario: Admin bypasses permission check
- **WHEN** a user has a role with `is_admin=true`
- **THEN** permission checks SHALL automatically return `true` for all permissions
- **AND** the user SHALL be able to access all menu functions
- **AND** the user SHALL see all menu items in the UI

### Requirement: User Permission Retrieval API
The system SHALL provide an API endpoint to retrieve the current user's permissions.

#### Scenario: Get current user permissions
- **WHEN** a user makes a `GET` request to `/api/v1/user/permissions`
- **AND** the user is authenticated
- **THEN** the system SHALL return a list of permissions in format:
  ```json
  {
    "permissions": [
      { "resource": "dashboard", "action": "read" },
      { "resource": "projects", "action": "*" },
      ...
    ]
  }
  ```
- **AND** for admin users, SHALL return all menu permissions

#### Scenario: Get permissions for admin user
- **WHEN** an admin user requests their permissions
- **THEN** the system SHALL return all defined menu permissions
- **AND** SHALL include all resource/action combinations

### Requirement: Route Permission Middleware
All API routes SHALL check user permissions before allowing access.

#### Scenario: Route with permission check
- **WHEN** a user makes a request to a protected route
- **THEN** the system SHALL check if the user has the required permission
- **AND** if the user has permission OR is admin, SHALL allow access
- **AND** if the user lacks permission, SHALL return 403 Forbidden

#### Scenario: Admin bypasses route permission check
- **WHEN** an admin user makes a request to any route
- **THEN** the permission check SHALL automatically pass
- **AND** the request SHALL proceed normally

### Requirement: Frontend Menu Permission Control
The frontend SHALL display menu items based on user permissions.

#### Scenario: Menu items filtered by permissions
- **WHEN** a user logs in
- **AND** the frontend retrieves user permissions
- **THEN** the menu SHALL only display items for which the user has permission
- **AND** menu items without permission SHALL be hidden

#### Scenario: Admin sees all menu items
- **WHEN** an admin user logs in
- **THEN** the menu SHALL display all available menu items
- **AND** all menu items SHALL be accessible

#### Scenario: Menu permission mapping
- **WHEN** rendering the menu
- **THEN** each menu item SHALL be checked against its required permission:
  - `/dashboard` → `dashboard:read`
  - `/projects` → `projects:*`
  - `/environments` → `environments:*`
  - `/cloud-platforms` → `cloud-platforms:*`
  - `/assets` → `assets:*`
  - `/executions` → `executions:*`
  - `/templates` → `templates:*`
  - `/tasks` → `tasks:*`
  - `/deployments` → `deployments:*`
  - `/ssh/keys` → `ssh-keys:*`
  - `/users` → `users:*`
  - `/roles` → `roles:*`
  - `/audit-logs` → `audit-logs:read`

### Requirement: Frontend Route Guard
The frontend SHALL prevent access to routes for which the user lacks permission.

#### Scenario: Access denied for unauthorized route
- **WHEN** a user attempts to access a route without permission
- **OR** directly enters a URL without permission
- **THEN** the frontend SHALL redirect to an error page
- **OR** display a 403 Forbidden message
- **AND** SHALL not render the page content

### Requirement: Role Permission Assignment UI
The role management interface SHALL allow assigning permissions to roles.

#### Scenario: Assign permissions to role
- **WHEN** a user edits a role in the role management interface
- **THEN** the interface SHALL display a list of available permissions
- **AND** SHALL allow selecting/deselecting permissions
- **AND** SHALL group permissions by resource (e.g., "Projects", "Environments")
- **AND** SHALL support bulk selection (select all, select by resource)
- **AND** SHALL display currently assigned permissions

#### Scenario: Save permission assignments
- **WHEN** a user assigns permissions to a role and saves
- **THEN** the system SHALL update the `role_permissions` table
- **AND** SHALL reflect changes immediately for users with that role
- **AND** SHALL show success/error feedback

## MODIFIED Requirements

### Requirement: RBAC Permission Check
The RBAC service SHALL check admin status before checking specific permissions.

#### Scenario: Permission check with admin override
- **WHEN** `HasPermission(userID, resource, action)` is called
- **THEN** it SHALL first check if the user has admin role (`is_admin=true`)
- **IF** admin, SHALL return `true` immediately
- **IF** not admin, SHALL check user's assigned permissions
- **AND** SHALL return `true` if permission exists, `false` otherwise

### Requirement: Protected Routes
All routes SHALL be protected by permission middleware where applicable.

#### Scenario: Route protection
- **WHEN** defining routes in `backend/cmd/server/main.go`
- **THEN** protected routes SHALL use `RequirePermission` middleware
- **AND** SHALL specify the required resource and action
- **AND** system management routes SHALL require admin role

### Requirement: User Login Flow
The login flow SHALL retrieve and store user permissions.

#### Scenario: Post-login permission retrieval
- **WHEN** a user successfully logs in
- **THEN** the frontend SHALL immediately fetch user permissions
- **AND** SHALL store permissions in the permission store
- **AND** SHALL use permissions for menu rendering and route guards
