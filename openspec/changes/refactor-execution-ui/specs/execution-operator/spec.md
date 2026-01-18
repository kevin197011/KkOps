# Execution Operator Specification

## ADDED Requirements

### Requirement: Quick Execution Interface
The system SHALL provide a quick execution interface that allows users to execute scripts on selected hosts without creating a task first.

#### Scenario: Execute from template
- **WHEN** a user selects "模板模式" (Template Mode)
- **AND** selects a task template from the dropdown
- **THEN** the script content and type SHALL be automatically populated from the template
- **AND** the script editor SHALL be read-only
- **AND** the template description SHALL be displayed

#### Scenario: Execute custom script
- **WHEN** a user selects "自定义模式" (Custom Mode)
- **THEN** the script editor SHALL be editable
- **AND** the user SHALL be able to input script content
- **AND** the user SHALL be able to select script type (shell/python)

#### Scenario: Select execution hosts
- **WHEN** a user is on the execution operator page
- **THEN** the host selector SHALL display all available hosts
- **AND** the user SHALL be able to filter hosts by project
- **AND** the user SHALL be able to filter hosts by environment
- **AND** the user SHALL be able to search hosts by hostname or IP
- **AND** the user SHALL be able to select multiple hosts using checkboxes
- **AND** the user SHALL be able to use batch selection buttons:
  - "全选" (Select All) - selects all filtered hosts
  - "按项目选" (Select by Project) - selects all hosts in the selected project
  - "按环境选" (Select by Environment) - selects all hosts in the selected environment
  - "清空" (Clear) - clears all selections
- **AND** the selected host count SHALL be displayed

#### Scenario: Configure execution options
- **WHEN** a user is on the execution operator page
- **THEN** the execution options SHALL include:
  - Execution mode (synchronous/asynchronous)
  - Timeout setting (in seconds, default: 600)
  - "保存为任务" (Save as Task) checkbox (optional)
- **AND** if "保存为任务" is checked, the user SHALL be prompted to enter task name and description after execution

#### Scenario: Execute script immediately
- **WHEN** a user has selected:
  - Template or entered custom script
  - At least one execution host
  - Execution options
- **AND** clicks the "执行" (Execute) button
- **THEN** the system SHALL immediately start execution on all selected hosts
- **AND** SHALL display execution status for each host
- **AND** SHALL stream execution logs in real-time for each host
- **AND** if "保存为任务" is checked, SHALL create a task record with the execution details

#### Scenario: View execution results
- **WHEN** execution starts
- **THEN** the results area SHALL display:
  - Overall execution status (running, success, failed, partial)
  - Individual host execution status cards
  - Expandable log viewer for each host
  - Execution statistics (total hosts, success count, failed count)
- **AND** each host card SHALL show:
  - Hostname and IP
  - Execution status badge
  - Execution time
  - Exit code (if available)
  - Expandable log viewer button

#### Scenario: Real-time log streaming
- **WHEN** execution is running
- **THEN** the system SHALL connect to WebSocket endpoint for log streaming
- **AND** SHALL receive log messages for each host
- **AND** SHALL update the corresponding host's log viewer in real-time
- **AND** SHALL auto-scroll to latest log line when viewing logs
- **AND** SHALL handle WebSocket disconnection gracefully

#### Scenario: Save execution as task (optional)
- **WHEN** a user has executed a script
- **AND** "保存为任务" was checked before execution
- **THEN** after execution starts, the system SHALL prompt the user to enter:
  - Task name
  - Task description (optional)
- **AND** SHALL create a task record with:
  - Script content and type
  - Selected host IDs
  - Execution options (timeout, etc.)
  - Reference to execution records
- **AND** SHALL provide a link to view the task in `/tasks` page

### Requirement: Temporary Execution
The system SHALL support temporary execution that does not create a task record.

#### Scenario: Temporary execution without saving
- **WHEN** a user executes a script without checking "保存为任务"
- **THEN** the system SHALL execute the script on selected hosts
- **AND** SHALL create execution records in `task_executions` table
- **AND** SHALL create a temporary task record (marked as `is_temporary=true`)
- **OR** SHALL create execution records with NULL `task_id`
- **AND** SHALL allow viewing execution results and logs
- **BUT** SHALL NOT appear in the task management list

### Requirement: Host Selection Enhancement
The host selector SHALL provide efficient filtering and selection capabilities.

#### Scenario: Filter hosts by project
- **WHEN** a user selects a project from the project filter
- **THEN** the host list SHALL only display hosts belonging to that project
- **AND** the selection SHALL update immediately

#### Scenario: Filter hosts by environment
- **WHEN** a user selects an environment from the environment filter
- **THEN** the host list SHALL only display hosts in that environment
- **AND** the selection SHALL update immediately
- **AND** if a project is also selected, SHALL show hosts matching both project AND environment

#### Scenario: Search hosts
- **WHEN** a user enters text in the host search box
- **THEN** the host list SHALL filter hosts by:
  - Hostname (case-insensitive partial match)
  - IP address (partial match)
- **AND** the filter SHALL apply in real-time as the user types

#### Scenario: Batch host selection
- **WHEN** a user clicks "全选" (Select All)
- **THEN** all currently filtered hosts SHALL be selected
- **WHEN** a user clicks "按项目选" (Select by Project)
- **AND** a project is selected
- **THEN** all hosts in that project (matching current filters) SHALL be selected
- **WHEN** a user clicks "按环境选" (Select by Environment)
- **AND** an environment is selected
- **THEN** all hosts in that environment (matching current filters) SHALL be selected
- **WHEN** a user clicks "清空" (Clear)
- **THEN** all host selections SHALL be cleared

### Requirement: Execution Form Validation
The system SHALL validate the execution form before allowing execution.

#### Scenario: Validate required fields
- **WHEN** a user clicks "执行" without selecting template or entering script
- **THEN** the system SHALL display an error: "请选择模板或输入脚本内容"
- **WHEN** a user clicks "执行" without selecting any hosts
- **THEN** the system SHALL display an error: "请至少选择一个执行主机"
- **WHEN** a user clicks "执行" with valid inputs
- **THEN** the system SHALL proceed with execution

## MODIFIED Requirements

### Requirement: Execution Page Layout (from task-management)
The execution page SHALL use a single-page form layout instead of dual-column workflow list.

#### Scenario: Single-page execution form
- **WHEN** a user navigates to `/executions`
- **THEN** the page SHALL display a single-page form layout with:
  - Execution mode selector at the top
  - Script editor in the middle
  - Host selector below script editor
  - Execution options below host selector
  - Execute button at the bottom
  - Results area (hidden until execution starts)
- **AND** SHALL NOT display a workflow list on the left
- **AND** SHALL provide a link to view saved tasks in `/tasks` page
