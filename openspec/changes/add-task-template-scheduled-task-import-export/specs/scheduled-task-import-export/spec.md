# Scheduled Task Import/Export Specification

## ADDED Requirements

### Requirement: Scheduled Task Export
The system SHALL provide functionality to export all scheduled tasks to a JSON file.

#### Scenario: Export all scheduled tasks
- **WHEN** a user clicks "导出配置" (Export Config) button in the scheduled task list page
- **THEN** the system SHALL call the export API endpoint
- **AND** SHALL return a JSON file with the following structure:
  ```json
  {
    "version": "1.0",
    "export_at": "2025-01-27T10:00:00Z",
    "tasks": [
      {
        "name": "task_name",
        "description": "task_description",
        "cron_expression": "0 0 * * * *",
        "template_name": "template_name",  // 可选，如果有关联模板
        "content": "script_content",  // 如果没有关联模板，使用此内容
        "type": "shell",
        "timeout": 600,
        "enabled": true,
        "update_assets": false,
        "target_hosts": ["host01", "192.168.1.1"]  // 主机名或 IP 列表
      }
    ]
  }
  ```
- **AND** SHALL convert `AssetIDs` (comma-separated string) to hostname/IP array
- **AND** SHALL convert template ID to template name (if template is associated)
- **AND** SHALL set `Content-Disposition` header to trigger file download
- **AND** SHALL use filename format: `scheduled-tasks-YYYYMMDD-HHMMSS.json`

#### Scenario: Export task without template
- **WHEN** a scheduled task does not have an associated template
- **THEN** the exported JSON SHALL include `content` and `type` fields
- **AND** SHALL set `template_name` to `null` or omit it

#### Scenario: Export task with template
- **WHEN** a scheduled task has an associated template
- **THEN** the exported JSON SHALL include `template_name` field
- **AND** SHALL include `content` and `type` fields (current values, not template values)

#### Scenario: Export empty task list
- **WHEN** there are no scheduled tasks in the system
- **THEN** the system SHALL return a JSON file with empty `tasks` array

### Requirement: Scheduled Task Import
The system SHALL provide functionality to import scheduled tasks from a JSON file.

#### Scenario: Import new tasks
- **WHEN** a user clicks "导入配置" (Import Config) button
- **AND** selects a valid JSON file containing scheduled tasks
- **AND** clicks "确认导入"
- **THEN** the system SHALL validate the JSON format
- **AND** SHALL validate required fields (name, cron_expression, content/type)
- **AND** SHALL validate Cron expression format (6-field format)
- **AND** SHALL match hosts by hostname or IP from `target_hosts` array
- **AND** SHALL match template by name from `template_name` (if provided)
- **AND** SHALL create tasks that do not already exist
- **AND** SHALL skip tasks with duplicate names
- **AND** SHALL return import result with statistics

#### Scenario: Import with host matching
- **WHEN** importing a task with `target_hosts` array
- **THEN** the system SHALL match each hostname/IP to an asset in the database
- **AND** SHALL collect matched asset IDs
- **AND** SHALL add unmatched hosts to warnings (not errors)
- **AND** SHALL continue importing the task even if some hosts are unmatched
- **AND** SHALL include unmatched hosts in `skipped_items` in import result

#### Scenario: Import with template matching
- **WHEN** importing a task with `template_name` field
- **THEN** the system SHALL look up the template by name
- **IF** template exists:
  - SHALL set `template_id` to the matched template ID
  - SHALL use provided `content` and `type` if specified
  - SHALL use template's content and type if `content` is empty
- **IF** template does not exist:
  - SHALL set `template_id` to `null`
  - SHALL add warning to `skipped_items`
  - SHALL require `content` and `type` to be provided
  - SHALL fail import if `content` is empty and template not found

#### Scenario: Import with duplicate names
- **WHEN** importing tasks that have the same name as existing tasks
- **THEN** the system SHALL skip those tasks
- **AND** SHALL add them to `skipped_items` in the import result
- **AND** SHALL continue processing remaining tasks

#### Scenario: Import with validation errors
- **WHEN** importing a task with:
  - Missing required fields (name, cron_expression)
  - Invalid Cron expression format
  - Missing content and template not found
- **THEN** the system SHALL mark that task as failed
- **AND** SHALL add error message to `errors` array
- **AND** SHALL continue processing remaining tasks

#### Scenario: Import preview (optional)
- **WHEN** a user selects a JSON file for import
- **AND** clicks "预览" (Preview) button
- **THEN** the system SHALL validate the JSON format, required fields, and Cron expressions
- **AND** SHALL check for duplicate names
- **AND** SHALL attempt to match hosts and templates (without creating tasks)
- **AND** SHALL return preview result showing:
  - Which tasks will be created
  - Which tasks will be skipped (duplicates)
  - Which tasks have validation errors
  - Which hosts were matched/unmatched
  - Which templates were matched/unmatched
- **BUT** SHALL NOT actually create any tasks

### Requirement: Host Matching Helper
The system SHALL provide functionality to match hosts by hostname or IP.

#### Scenario: Match hosts by hostname
- **WHEN** importing a task with `target_hosts: ["host01", "host02"]`
- **THEN** the system SHALL query assets where `host_name IN ('host01', 'host02')`
- **AND** SHALL return matched asset IDs

#### Scenario: Match hosts by IP
- **WHEN** importing a task with `target_hosts: ["192.168.1.1", "192.168.1.2"]`
- **THEN** the system SHALL query assets where `ip IN ('192.168.1.1', '192.168.1.2')`
- **AND** SHALL return matched asset IDs

#### Scenario: Match hosts by mixed hostname and IP
- **WHEN** importing a task with `target_hosts: ["host01", "192.168.1.1"]`
- **THEN** the system SHALL match by hostname first, then by IP
- **AND** SHALL return all matched asset IDs
- **AND** SHALL not create duplicate matches

## MODIFIED Requirements

### Requirement: Scheduled Task Management API
The scheduled task API SHALL support import/export operations.

#### Scenario: Export endpoint
- **WHEN** making a `GET` request to `/api/v1/tasks/export`
- **THEN** the system SHALL return all scheduled tasks in JSON format
- **AND** SHALL include export metadata (version, export_at)
- **AND** SHALL convert internal IDs to human-readable names (hosts, templates)

#### Scenario: Import endpoint
- **WHEN** making a `POST` request to `/api/v1/tasks/import`
- **WITH** a JSON body containing import configuration
- **THEN** the system SHALL validate and import tasks
- **AND** SHALL match hosts and templates by name/IP
- **AND** SHALL return import result with statistics
