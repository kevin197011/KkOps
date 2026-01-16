## MODIFIED Requirements

### Requirement: Asset Management
The system SHALL provide comprehensive IT asset management including creation, retrieval, update, and deletion of assets.

#### Scenario: Create asset
- **WHEN** a user with appropriate permissions creates a new asset with required information (hostname, code, status)
- **THEN** a new asset record is created
- **AND** the asset can include optional information (project, cloud_platform, IP, SSH credentials, hardware specs, tags)
- **AND** the asset does NOT require a category to be specified
- **NOTE**: The hostname field is labeled as "主机名" in the UI and stored in the `name` field in the database

#### Scenario: List assets
- **WHEN** a user requests the asset list
- **THEN** the system returns a paginated list of assets
- **AND** the list displays hostname in the "主机名" column
- **AND** the list supports filtering by project, cloud_platform, environment, status, tags, and IP
- **AND** the list supports searching by hostname, code, IP, or other fields

#### Scenario: View asset details
- **WHEN** a user views an asset detail page or modal
- **THEN** the system displays complete asset information including hostname (labeled as "主机名"), SSH connection info, hardware specs, tags, and relationships
- **AND** the hostname is clearly labeled and displayed

#### Scenario: Update asset
- **WHEN** a user with appropriate permissions updates asset information
- **THEN** the asset record is updated
- **AND** the updated_at timestamp is automatically updated
- **AND** the hostname field can be updated (displayed as "主机名" in the UI)

#### Scenario: Search assets
- **WHEN** a user searches for assets by keyword
- **THEN** the system searches across multiple fields including hostname, code, IP, description
- **AND** returns matching assets
- **AND** the search placeholder indicates searching by "资产主机名、代码、IP"
