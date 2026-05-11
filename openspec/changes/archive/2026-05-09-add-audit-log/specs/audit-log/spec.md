# Audit Log Specification

## ADDED Requirements

### Requirement: Audit Log Recording
The system SHALL automatically record audit logs for all critical user operations.

#### Scenario: User login audit
- **WHEN** a user successfully logs in
- **THEN** an audit log entry SHALL be created with:
  - `module`: "auth"
  - `action`: "login"
  - `username`: the logged-in user's username
  - `status`: "success"
  - `ip_address`: the client's IP address

#### Scenario: User login failure audit
- **WHEN** a user fails to log in
- **THEN** an audit log entry SHALL be created with:
  - `module`: "auth"
  - `action`: "login"
  - `username`: the attempted username
  - `status`: "failed"
  - `error_msg`: the failure reason

#### Scenario: Resource creation audit
- **WHEN** a user creates a resource (user, role, asset, task, deployment module)
- **THEN** an audit log entry SHALL be created with:
  - `module`: the resource type
  - `action`: "create"
  - `resource_id`: the created resource's ID
  - `resource_name`: the created resource's name
  - `detail`: JSON containing the creation parameters (excluding sensitive fields)

#### Scenario: Resource update audit
- **WHEN** a user updates a resource
- **THEN** an audit log entry SHALL be created with:
  - `module`: the resource type
  - `action`: "update"
  - `resource_id`: the updated resource's ID
  - `detail`: JSON containing the changed fields (before and after values)

#### Scenario: Resource deletion audit
- **WHEN** a user deletes a resource
- **THEN** an audit log entry SHALL be created with:
  - `module`: the resource type
  - `action`: "delete"
  - `resource_id`: the deleted resource's ID
  - `resource_name`: the deleted resource's name

#### Scenario: Task execution audit
- **WHEN** a user executes a task
- **THEN** an audit log entry SHALL be created with:
  - `module`: "task"
  - `action`: "execute"
  - `resource_id`: the task ID
  - `detail`: JSON containing target hosts and execution mode

### Requirement: Audit Log Query
The system SHALL provide an API to query audit logs with filtering and pagination.

#### Scenario: Query audit logs with filters
- **WHEN** an administrator queries audit logs with filters
- **THEN** the system SHALL return audit logs matching:
  - `user_id` or `username` (optional)
  - `module` (optional)
  - `action` (optional)
  - `status` (optional)
  - `start_time` and `end_time` (optional)
  - `keyword` for searching in resource_name or detail (optional)

#### Scenario: Paginated audit log results
- **WHEN** querying audit logs
- **THEN** results SHALL be paginated with:
  - `page` (default: 1)
  - `page_size` (default: 20, max: 100)
  - Total count in response

### Requirement: Audit Log Export
The system SHALL support exporting audit logs in CSV and JSON formats.

#### Scenario: Export audit logs as CSV
- **WHEN** an administrator exports audit logs with format=csv
- **THEN** the system SHALL return a downloadable CSV file containing:
  - All columns: ID, Time, Username, Module, Action, Resource, Status, IP Address
  - Filtered by the same parameters as the query API

#### Scenario: Export audit logs as JSON
- **WHEN** an administrator exports audit logs with format=json
- **THEN** the system SHALL return a downloadable JSON file containing:
  - All audit log fields
  - Filtered by the same parameters as the query API

### Requirement: Audit Log Immutability
Audit logs SHALL be immutable and cannot be modified or deleted through the API.

#### Scenario: Attempt to delete audit log
- **WHEN** any user attempts to delete an audit log
- **THEN** the system SHALL reject the request with a 403 Forbidden error

#### Scenario: Attempt to modify audit log
- **WHEN** any user attempts to modify an audit log
- **THEN** the system SHALL reject the request with a 403 Forbidden error

### Requirement: Audit Log Retention
The system SHALL support configurable audit log retention with automatic cleanup.

#### Scenario: Automatic cleanup of old logs
- **WHEN** the configured retention period (default: 90 days) expires
- **THEN** the system SHALL automatically delete audit logs older than the retention period
- **AND** this cleanup operation itself SHALL be logged

### Requirement: Sensitive Data Protection
Audit logs SHALL NOT contain sensitive information.

#### Scenario: Password field filtering
- **WHEN** recording an audit log for user creation or password change
- **THEN** the `detail` field SHALL NOT contain any password values
- **AND** password fields SHALL be replaced with "[FILTERED]"

#### Scenario: SSH key content filtering
- **WHEN** recording an audit log for SSH key operations
- **THEN** the `detail` field SHALL NOT contain private key content
- **AND** private key fields SHALL be replaced with "[FILTERED]"

### Requirement: Audit Log UI
The system SHALL provide a user interface for viewing and exporting audit logs.

#### Scenario: View audit log list
- **WHEN** an administrator navigates to the audit log page
- **THEN** the system SHALL display a table with:
  - Time, Username, Module, Action, Resource, Status, IP Address columns
  - Sortable by time (default: newest first)
  - Clickable rows to view details

#### Scenario: Filter audit logs
- **WHEN** an administrator applies filters
- **THEN** the table SHALL update to show only matching logs
- **AND** filters SHALL include: user, module, action, status, time range, keyword

#### Scenario: View audit log detail
- **WHEN** an administrator clicks on an audit log entry
- **THEN** a modal SHALL display the full audit log details including:
  - All basic fields
  - Formatted JSON detail
  - User agent information
