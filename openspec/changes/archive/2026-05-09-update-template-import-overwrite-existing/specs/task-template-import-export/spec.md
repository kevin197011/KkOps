# Task Template Import/Export – Overwrite Existing on Import

## MODIFIED Requirements

### Requirement: Task Template Import
The system SHALL provide functionality to import task templates from a JSON file. When an imported template has the same name as an existing template, the system SHALL **overwrite** (update) the existing record with the imported content instead of skipping it.

#### Scenario: Import new templates
- **WHEN** a user clicks "导入配置" (Import Config) button
- **AND** selects a valid JSON file containing task templates
- **AND** clicks "确认导入"
- **THEN** the system SHALL validate the JSON format
- **AND** SHALL validate required fields (name, content) for each template
- **AND** SHALL create templates that do not already exist (by name)
- **AND** SHALL update existing templates when the name matches (overwrite with imported description, content, type)
- **AND** SHALL return import result with statistics including:
  - `total`: number of templates in the import file
  - `success`: number of newly created templates
  - `updated`: number of existing templates overwritten
  - `failed`: number of templates that failed validation or save
  - `errors`: list of error messages for failed items
  - `updated_items`: optional list of template names that were updated

#### Scenario: Import with duplicate names (overwrite)
- **WHEN** importing templates that have the same name as existing templates
- **THEN** the system SHALL update those existing templates with the imported content (description, content, type)
- **AND** SHALL count them as updated (not skipped)
- **AND** SHALL include their names in the import result (e.g. `updated_items`)
- **AND** SHALL continue processing remaining templates

#### Scenario: Import with validation errors
- **WHEN** importing a template with missing required fields (e.g. name or content)
- **OR** invalid JSON format
- **THEN** the system SHALL mark that template as failed
- **AND** SHALL add error message to `errors` array
- **AND** SHALL continue processing remaining templates
- **AND** SHALL return the import result with failure count

#### Scenario: Import preview (optional)
- **WHEN** a user selects a JSON file for import
- **AND** clicks "预览" (Preview) button (if implemented)
- **THEN** the system SHALL validate the JSON format and required fields
- **AND** SHALL indicate which templates will be created and which existing ones will be updated (by name)
- **AND** SHALL NOT actually create or update any templates
