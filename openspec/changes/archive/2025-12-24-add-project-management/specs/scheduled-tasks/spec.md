## MODIFIED Requirements

### Requirement: Scheduled Task Creation
The system SHALL provide the ability to create scheduled tasks that execute Salt commands or states at specified times. All scheduled tasks MUST be associated with a project.

#### Scenario: Create a scheduled task
- **WHEN** a user creates a scheduled task with name, schedule (cron expression), Salt command/state, target hosts, and project
- **THEN** the system SHALL validate that the project exists and the user has access to it
- **AND** the system SHALL validate that all target hosts belong to the same project
- **AND** the system SHALL validate the cron expression
- **AND** the system SHALL store the task configuration associated with the project
- **AND** the system SHALL return the created task with a unique ID

#### Scenario: List scheduled tasks
- **WHEN** a user requests the list of scheduled tasks
- **THEN** the system SHALL return all tasks
- **AND** the system SHALL support filtering by name, status, and project
- **AND** the system SHALL only return tasks from projects the user has access to

### Requirement: Scheduled Task Execution
The system SHALL provide the ability to execute scheduled tasks and track their execution status. Task executions inherit the project from the task.

#### Scenario: Execute scheduled task automatically
- **WHEN** a scheduled task's execution time arrives
- **THEN** the system SHALL create an execution record with status "running" and the task's project
- **AND** the system SHALL execute the Salt command/state on target hosts
- **AND** the system SHALL track the execution progress
- **AND** the system SHALL update the execution record with results

#### Scenario: Query task execution history
- **WHEN** a user queries execution history for tasks
- **THEN** the system SHALL return a paginated list of execution records
- **AND** the system SHALL support filtering by status, date range, execution type, and project
- **AND** the system SHALL only return executions from projects the user has access to
- **AND** the system SHALL include execution time, status, duration, and results summary

