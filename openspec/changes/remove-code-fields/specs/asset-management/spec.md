## MODIFIED Requirements

### Requirement: Asset Management
The system SHALL provide comprehensive IT asset management including creation, retrieval, update, and deletion of assets.

#### Scenario: Create asset
- **WHEN** a user with appropriate permissions creates a new asset with required information (hostname, status)
- **THEN** a new asset record is created
- **AND** the asset hostname is stored in the `host_name` database column
- **AND** the API response includes the hostname in the `hostName` JSON field
- **AND** the asset does NOT require a code field
- **AND** the asset can include optional information (project, cloud_platform, IP, SSH credentials, hardware specs, tags)
- **AND** the asset does NOT require a category to be specified

#### Scenario: List assets
- **WHEN** a user requests the asset list
- **THEN** the system returns a paginated list of assets
- **AND** each asset includes the `hostName` field in the JSON response
- **AND** the list does NOT include a code field
- **AND** the list supports filtering by project, cloud_platform, environment, status, tags, and IP
- **AND** the list supports searching by hostname, IP, or description (code is NOT searchable)

#### Scenario: Search assets
- **WHEN** a user searches for assets by keyword
- **THEN** the system searches across multiple fields including hostname (stored in `host_name` column), IP, description
- **AND** the system does NOT search the code field (field does not exist)
- **AND** returns matching assets with `hostName` field in response

#### Scenario: Export assets
- **WHEN** a user exports assets
- **THEN** the system generates a CSV file containing asset data
- **AND** the CSV file includes a `Hostname` column header
- **AND** the CSV file does NOT include a `Code` column
- **AND** the hostname values are included in the export

#### Scenario: Import assets
- **WHEN** a user uploads a CSV file containing asset data
- **THEN** the system validates the CSV header includes a `Hostname` column
- **AND** the system does NOT require a `Code` column
- **AND** the system validates and imports the assets
- **AND** validation errors are reported for invalid records
- **AND** successfully imported assets are created with hostname from the `Hostname` column
