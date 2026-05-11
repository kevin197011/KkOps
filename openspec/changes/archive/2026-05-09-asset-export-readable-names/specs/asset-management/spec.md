## MODIFIED Requirements

### Requirement: Asset Export (Readable Names)
The system SHALL export asset data to CSV with human-readable names for Project, Cloud Platform, Environment, and SSH Key instead of numeric IDs.

#### Scenario: Export assets with readable names
- **WHEN** a user exports assets
- **THEN** the system generates a CSV file containing asset data
- **AND** the CSV includes columns "Project", "Cloud Platform", "Environment", "SSH Key" (not "Project ID", "Environment ID", "SSH Key ID")
- **AND** the values in those columns are the **names** of the associated project, cloud platform, environment, and SSH key respectively
- **AND** when an asset has no associated project, environment, or SSH key, the corresponding cell is empty
- **AND** the exported file is suitable for manual inspection and cross-reference without looking up IDs

#### Scenario: Import assets with ID or name columns (optional)
- **WHEN** a user uploads a CSV that includes "Project ID", "Environment ID", "SSH Key ID" columns
- **THEN** the system continues to accept those columns and resolve by ID (backward compatible)
- **WHEN** the CSV includes "Project", "Environment", "SSH Key" (name) columns instead of or in addition to ID columns
- **THEN** the system MAY resolve by name to the corresponding entity ID and assign the asset
- **AND** invalid or unknown names are reported as validation errors or the association is skipped (per design)
