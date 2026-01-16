## ADDED Requirements

### Requirement: Asset Management
The system SHALL provide comprehensive IT asset management including creation, retrieval, update, and deletion of assets.

#### Scenario: Create asset
- **WHEN** a user with appropriate permissions creates a new asset with required information (name, code, category, status)
- **THEN** a new asset record is created
- **AND** the asset can include optional information (project, cloud_platform, IP, SSH credentials, hardware specs, tags)

#### Scenario: List assets
- **WHEN** a user requests the asset list
- **THEN** the system returns a paginated list of assets
- **AND** the list supports filtering by project, cloud_platform, environment, status, tags, and IP
- **AND** the list supports searching by name, code, IP, or other fields

#### Scenario: View asset details
- **WHEN** a user views an asset detail page
- **THEN** the system displays complete asset information including basic info, SSH connection info, hardware specs, tags, and relationships
- **AND** the user can see related assets if any

#### Scenario: Update asset
- **WHEN** a user with appropriate permissions updates asset information
- **THEN** the asset record is updated
- **AND** the updated_at timestamp is automatically updated

#### Scenario: Delete asset
- **WHEN** a user with appropriate permissions deletes an asset
- **THEN** the asset is removed from the database
- **AND** associated relationships are handled appropriately

### Requirement: Asset Status Management
The system SHALL support asset status management with active and disabled states.

#### Scenario: Disable asset
- **WHEN** an asset status is set to disabled
- **THEN** the asset is marked as disabled
- **AND** the asset may be excluded from certain operations (e.g., SSH connections, task execution)

#### Scenario: Enable asset
- **WHEN** an asset status is set to active
- **THEN** the asset is marked as active
- **AND** the asset is available for all operations

### Requirement: Project Association
The system SHALL support associating assets with projects for organizational purposes.

#### Scenario: Assign asset to project
- **WHEN** an asset is assigned to a project
- **THEN** the asset's project_id is updated
- **AND** the asset can be filtered and grouped by project

### Requirement: Tag Management
The system SHALL support tagging assets for categorization and filtering.

#### Scenario: Create tag
- **WHEN** an administrator creates a tag with name, color, and description
- **THEN** a new tag is created
- **AND** the tag can be assigned to assets

#### Scenario: Assign tags to asset
- **WHEN** tags are assigned to an asset
- **THEN** the asset-tag relationships are established
- **AND** the asset can be filtered by any of its tags
- **AND** an asset can have multiple tags

#### Scenario: Filter assets by tags
- **WHEN** a user filters assets by one or more tags
- **THEN** the system returns only assets that have all specified tags

### Requirement: Asset Import/Export
The system SHALL support bulk import and export of asset data.

#### Scenario: Export assets
- **WHEN** a user exports assets
- **THEN** the system generates a file (CSV or Excel) containing asset data
- **AND** the exported file includes all asset fields

#### Scenario: Import assets
- **WHEN** a user uploads a file containing asset data
- **THEN** the system validates and imports the assets
- **AND** validation errors are reported for invalid records
- **AND** successfully imported assets are created in the database

### Requirement: Asset Search and Filtering
The system SHALL provide comprehensive search and filtering capabilities for assets.

#### Scenario: Search assets
- **WHEN** a user searches for assets by keyword
- **THEN** the system searches across multiple fields (name, code, IP, description)
- **AND** returns matching assets

#### Scenario: Filter assets
- **WHEN** a user applies filters (project, cloud_platform, environment, status, tags, IP)
- **THEN** the system returns only assets matching all filter criteria
- **AND** multiple filter values within the same category use OR logic (e.g., multiple tags)
