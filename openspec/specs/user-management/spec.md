# User Management Specification

## Overview
用户管理系统提供完整的用户认证、角色权限管理和用户生命周期管理功能。

## Requirements

### Requirement: User Authentication
The system SHALL provide secure user authentication with JWT tokens.

#### Scenario: Successful Login
- **WHEN** user provides valid username and password
- **THEN** system returns JWT access token
- **AND** token contains user ID, username, and role information
- **AND** token expires after 24 hours

#### Scenario: Invalid Credentials
- **WHEN** user provides invalid username or password
- **THEN** system returns 401 Unauthorized
- **AND** includes error message "Invalid credentials"

#### Scenario: Account Locked
- **WHEN** user account is locked due to security policy
- **THEN** system returns 423 Locked status
- **AND** includes message about account status

### Requirement: User Registration
The system SHALL allow administrators to create new user accounts.

#### Scenario: Admin Creates User
- **WHEN** administrator provides username, email, and initial password
- **THEN** system creates new user account
- **AND** assigns default role (viewer)
- **AND** sends welcome notification

### Requirement: Password Management
The system SHALL provide secure password management with encryption.

#### Scenario: Password Hashing
- **WHEN** user sets or changes password
- **THEN** system hashes password using bcrypt
- **AND** stores only hashed version in database

#### Scenario: Password Change
- **WHEN** user requests password change
- **THEN** system validates current password
- **AND** requires strong password policy
- **AND** updates password hash in database

### Requirement: Role-Based Access Control (RBAC)
The system SHALL implement RBAC with predefined roles and granular permissions.

#### Scenario: Role Assignment
- **WHEN** administrator assigns role to user
- **THEN** user inherits all permissions from that role
- **AND** permissions are cached for performance

#### Scenario: Permission Check
- **WHEN** user attempts to access protected resource
- **THEN** system validates user has required permission
- **AND** returns 403 Forbidden if permission denied

### Requirement: User Profile Management
The system SHALL allow users to manage their profile information.

#### Scenario: Profile Update
- **WHEN** user updates display name or email
- **THEN** system validates input data
- **AND** updates user record in database
- **AND** reflects changes in JWT token

## Predefined Roles

### Administrator (admin)
- Full system access
- User management
- System configuration
- All permissions

### Operator (operator)
- Host management
- Deployment operations
- SSH access
- Batch operations
- Formula management

### Viewer (viewer)
- Read-only access to most resources
- Cannot modify system state
- Can view logs and reports

## Permissions List

### User Management
- `user:create` - Create new users
- `user:read` - View user information
- `user:update` - Update user profiles
- `user:delete` - Delete users
- `user:manage` - Full user administration

### Host Management
- `host:create` - Create hosts
- `host:read` - View host information
- `host:update` - Update host details
- `host:delete` - Delete hosts

### SSH Management
- `ssh:execute` - Execute SSH connections
- `ssh:key:create` - Create SSH keys
- `ssh:key:read` - View SSH keys
- `ssh:key:update` - Update SSH keys
- `ssh:key:delete` - Delete SSH keys

### WebSSH Management
- `webssh:execute` - Access WebSSH terminals

### Batch Operations
- `batch_operation:read` - View batch operations
- `batch_operation:create` - Create batch operations
- `batch_operation:execute` - Execute batch operations
- `batch_operation:cancel` - Cancel running operations

### Formula Management
- `formula:read` - View formulas
- `formula:create` - Create formulas
- `formula:update` - Update formulas
- `formula:delete` - Delete formulas
- `formula:execute` - Execute formula deployments

### Deployment Management
- `deployment:read` - View deployments
- `deployment:create` - Create deployments
- `deployment:execute` - Execute deployments

### System Settings
- `system:read` - View system settings
- `system:update` - Update system settings

### Audit Management
- `audit:read` - View audit logs

## Security Considerations

### Password Policy
- Minimum 8 characters
- Must contain uppercase, lowercase, and numeric characters
- Cannot reuse recent passwords

### Session Management
- JWT tokens expire after 24 hours
- Refresh tokens for extended sessions
- Automatic logout on suspicious activity

### Audit Logging
- All authentication attempts logged
- Password changes tracked
- Role assignments audited
- Permission checks recorded
