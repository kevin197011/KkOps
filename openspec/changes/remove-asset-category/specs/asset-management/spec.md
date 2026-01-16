## MODIFIED Requirements

### Requirement: Asset Management
The system SHALL provide comprehensive IT asset management including creation, retrieval, update, and deletion of assets.

#### Scenario: Create asset
- **WHEN** a user with appropriate permissions creates a new asset with required information (name, code, status)
- **THEN** a new asset record is created
- **AND** the asset can include optional information (project, cloud_platform, IP, SSH credentials, hardware specs, tags)
- **AND** the asset does NOT require a category to be specified

#### Scenario: Update asset
- **WHEN** a user with appropriate permissions updates asset information
- **THEN** the asset record is updated
- **AND** the updated_at timestamp is automatically updated
- **AND** category information is NOT part of the asset update process

## REMOVED Requirements

### Requirement: Asset Category Association (REMOVED)
Assets no longer require association with an asset category. The category field has been removed from the asset model and API.
