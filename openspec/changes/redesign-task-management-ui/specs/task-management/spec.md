# Task Management Capability

## ADDED Requirements

### Requirement: Dual-Column Task Management Layout
The system SHALL display a dual-column layout with workflow list on the left and execution details on the right.

#### Scenario: User views task management page
- **WHEN** user navigates to the task management page
- **THEN** the left panel displays a list of workflow cards
- **AND** the right panel displays details of the selected workflow or a prompt to select one

#### Scenario: User selects a workflow
- **WHEN** user clicks on a workflow card in the left panel
- **THEN** the right panel updates to show the workflow details and execution history
- **AND** the selected card is visually highlighted

### Requirement: Workflow Status Indicators
The system SHALL display visual status indicators on workflow cards showing the current or last execution state.

#### Scenario: Workflow with successful execution
- **WHEN** the last execution of a workflow completed successfully
- **THEN** the workflow card displays a green success indicator
- **AND** shows the time since last execution

#### Scenario: Workflow currently running
- **WHEN** a workflow has an execution in progress
- **THEN** the workflow card displays a blue running indicator with pulse animation
- **AND** shows "Running..." status text

#### Scenario: Workflow with failed execution
- **WHEN** the last execution of a workflow failed
- **THEN** the workflow card displays a red failed indicator
- **AND** shows the time since last execution

### Requirement: Inline Execution Log Viewer
The system SHALL display execution logs inline within the execution details panel without page navigation.

#### Scenario: User expands execution logs
- **WHEN** user clicks "View Logs" on an execution record
- **THEN** the logs section expands below the execution record
- **AND** logs stream in real-time via WebSocket if execution is running
- **AND** logs display with per-host separation

#### Scenario: User collapses execution logs
- **WHEN** user clicks on an expanded log section header
- **THEN** the log section collapses
- **AND** the WebSocket connection for that execution is closed

### Requirement: Workflow Action Controls
The system SHALL provide action buttons to edit, execute, and cancel workflows from the execution details panel.

#### Scenario: User executes a workflow
- **WHEN** user clicks the "Execute" button in the details panel
- **THEN** a confirmation dialog appears with execution type selection (sync/async)
- **AND** upon confirmation, the workflow execution starts
- **AND** the workflow card updates to show running status

#### Scenario: User cancels a running execution
- **WHEN** user clicks "Cancel" on a running execution
- **THEN** a confirmation dialog appears
- **AND** upon confirmation, the execution is cancelled
- **AND** the status updates to "cancelled"

#### Scenario: User edits a workflow
- **WHEN** user clicks the "Edit" button
- **THEN** a modal opens with the workflow configuration
- **AND** user can modify name, content, and target hosts

### Requirement: Project/Environment-based Host Selection
The system SHALL allow users to filter and select execution target hosts by project or environment.

#### Scenario: User selects hosts by project
- **WHEN** user opens the host selection in task creation/edit modal
- **AND** selects a project from the filter dropdown
- **THEN** only assets belonging to that project are displayed
- **AND** user can multi-select from the filtered list

#### Scenario: User selects hosts by environment
- **WHEN** user opens the host selection in task creation/edit modal
- **AND** selects an environment from the filter dropdown
- **THEN** only assets in that environment are displayed
- **AND** user can multi-select from the filtered list

### Requirement: Execution History Display
The system SHALL display a list of recent executions for the selected workflow with status and timing information.

#### Scenario: User views execution history
- **WHEN** user selects a workflow
- **THEN** the execution history shows the 5 most recent executions
- **AND** each entry displays run number, status, and relative time
- **AND** a "Load more" option is available for older executions

#### Scenario: Execution history with per-host details
- **WHEN** user views an execution that ran on multiple hosts
- **THEN** the execution shows aggregate status
- **AND** expanding the execution reveals per-host status and logs

### Requirement: Real-time Status Updates
The system SHALL update workflow and execution statuses in real-time without page refresh.

#### Scenario: Execution completes while user is viewing
- **WHEN** an execution completes (success or failure)
- **AND** user is viewing the task management page
- **THEN** the workflow card status updates automatically
- **AND** the execution history updates with final status
- **AND** a notification appears indicating completion
