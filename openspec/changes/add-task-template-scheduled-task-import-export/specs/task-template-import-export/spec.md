# Task Template Import/Export Specification

## ADDED Requirements

### Requirement: Task Template Export
The system SHALL provide functionality to export all task templates to a JSON file.

#### Scenario: Export all task templates
- **WHEN** a user clicks "导出配置" (Export Config) button in the task template list page
- **THEN** the system SHALL call the export API endpoint
- **AND** SHALL return a JSON file with the following structure:
  ```json
  {
    "version": "1.0",
    "export_at": "2025-01-27T10:00:00Z",
    "templates": [
      {
        "name": "template_name",
        "description": "template_description",
        "content": "script_content",
        "type": "shell"
      }
    ]
  }
  ```
- **AND** SHALL set `Content-Disposition` header to trigger file download
- **AND** SHALL set `Content-Type` to `application/json`
- **AND** SHALL use filename format: `task-templates-YYYYMMDD-HHMMSS.json`

#### Scenario: Export empty template list
- **WHEN** there are no task templates in the system
- **AND** user clicks "导出配置"
- **THEN** the system SHALL return a JSON file with empty `templates` array:
  ```json
  {
    "version": "1.0",
    "export_at": "...",
    "templates": []
  }
  ```

### Requirement: Task Template Import
The system SHALL provide functionality to import task templates from a JSON file.

#### Scenario: Import new templates
- **WHEN** a user clicks "导入配置" (Import Config) button
- **AND** selects a valid JSON file containing task templates
- **AND** clicks "确认导入"
- **THEN** the system SHALL validate the JSON format
- **AND** SHALL validate required fields (name, content, type) for each template
- **AND** SHALL create templates that do not already exist
- **AND** SHALL skip templates with duplicate names (if already exists)
- **AND** SHALL return import result with statistics:
  ```json
  {
    "total": 10,
    "success": 8,
    "failed": 1,
    "skipped": 1,
    "errors": ["模板 'test': 缺少必填字段 'content'"],
    "skipped_items": ["模板 'daily': 已存在，已跳过"]
  }
  ```

#### Scenario: Import with duplicate names
- **WHEN** importing templates that have the same name as existing templates
- **THEN** the system SHALL skip those templates
- **AND** SHALL add them to `skipped_items` in the import result
- **AND** SHALL continue processing remaining templates

#### Scenario: Import with validation errors
- **WHEN** importing a template with missing required fields
- **OR** invalid JSON format
- **THEN** the system SHALL mark that template as failed
- **AND** SHALL add error message to `errors` array
- **AND** SHALL continue processing remaining templates
- **AND** SHALL return the import result with failure count

#### Scenario: Import preview (optional)
- **WHEN** a user selects a JSON file for import
- **AND** clicks "预览" (Preview) button
- **THEN** the system SHALL validate the JSON format and required fields
- **AND** SHALL check for duplicate names
- **AND** SHALL return preview result showing:
  - Which templates will be created
  - Which templates will be skipped (duplicates)
  - Which templates have validation errors
- **BUT** SHALL NOT actually create any templates

## MODIFIED Requirements

### Requirement: Task Template Management API
The task template API SHALL support import/export operations.

#### Scenario: Export endpoint
- **WHEN** making a `GET` request to `/api/v1/templates/export`
- **THEN** the system SHALL return all task templates in JSON format
- **AND** SHALL include export metadata (version, export_at)

#### Scenario: Import endpoint
- **WHEN** making a `POST` request to `/api/v1/templates/import`
- **WITH** a JSON body containing import configuration
- **THEN** the system SHALL validate and import templates
- **AND** SHALL return import result with statistics
