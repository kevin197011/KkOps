# Host Management Specification

## Overview
主机管理系统提供完整的主机资产管理、分组、标签和环境管理功能，支持大规模主机集群的管理。

## Requirements

### Requirement: Host Information Management
The system SHALL maintain comprehensive host information including connection details, system specs, and metadata.

#### Scenario: Host Creation
- **WHEN** administrator creates new host
- **THEN** system validates required fields (hostname, IP address)
- **AND** creates host record with unique ID
- **AND** initializes default values

#### Scenario: Host Update
- **WHEN** user updates host information
- **THEN** system validates input data
- **AND** updates host record
- **AND** records change in audit log

#### Scenario: Host Deletion
- **WHEN** administrator deletes host
- **THEN** system removes host from all groups and tags
- **AND** marks host as deleted (soft delete)
- **AND** logs deletion in audit

### Requirement: Host Grouping
The system SHALL support hierarchical host grouping for organizational management.

#### Scenario: Group Creation
- **WHEN** user creates host group
- **THEN** system validates group name uniqueness
- **AND** creates group with optional description
- **AND** allows nested group structure

#### Scenario: Host Assignment
- **WHEN** user assigns host to group
- **THEN** system creates group membership record
- **AND** validates no circular dependencies
- **AND** updates host's group associations

### Requirement: Host Tagging
The system SHALL provide flexible tagging system for host categorization.

#### Scenario: Tag Creation
- **WHEN** user creates host tag
- **THEN** system validates tag name and color
- **AND** creates tag with description
- **AND** ensures color uniqueness for visual distinction

#### Scenario: Tag Assignment
- **WHEN** user assigns tag to host
- **THEN** system creates tag association
- **AND** allows multiple tags per host
- **AND** provides visual tag indicators

### Requirement: Environment Management
The system SHALL support multiple deployment environments (dev, staging, prod).

#### Scenario: Environment Creation
- **WHEN** administrator creates environment
- **THEN** system validates environment name
- **AND** sets default configuration
- **AND** associates with host groups

#### Scenario: Environment Assignment
- **WHEN** user assigns environment to host
- **THEN** system updates host environment
- **AND** applies environment-specific settings

### Requirement: Salt Minion Integration
The system SHALL integrate with Salt Minion for automated host discovery and status monitoring.

#### Scenario: Minion Auto-Discovery
- **WHEN** Salt Master reports new minions
- **THEN** system automatically creates host records
- **AND** attempts IP address mapping
- **AND** sets initial status as "discovered"

#### Scenario: Minion Status Sync
- **WHEN** system performs status synchronization
- **THEN** queries Salt Master for minion status
- **AND** updates host online/offline status
- **AND** collects system information (CPU, memory, OS)

#### Scenario: SSH Port Detection
- **WHEN** syncing host information
- **THEN** system executes Salt command to detect SSH port
- **AND** updates host SSH port configuration
- **AND** handles detection failures gracefully

### Requirement: Host Search and Filtering
The system SHALL provide advanced search and filtering capabilities.

#### Scenario: Text Search
- **WHEN** user searches by hostname or IP
- **THEN** system returns matching hosts
- **AND** supports partial matching
- **AND** highlights search terms

#### Scenario: Advanced Filtering
- **WHEN** user filters by environment, group, or tags
- **THEN** system applies multiple filter criteria
- **AND** supports combination filters
- **AND** provides filter persistence

### Requirement: Bulk Host Operations
The system SHALL support bulk operations on multiple hosts.

#### Scenario: Bulk Group Assignment
- **WHEN** user selects multiple hosts
- **THEN** allows bulk group assignment
- **AND** updates all selected hosts
- **AND** provides operation progress

#### Scenario: Bulk Tagging
- **WHEN** user applies tag to multiple hosts
- **THEN** creates tag associations for all hosts
- **AND** handles conflicts gracefully

## Host Information Fields

### Required Fields
- **ID**: Unique identifier (UUID)
- **Hostname**: Server hostname
- **IP Address**: Primary IP address (unique constraint)
- **Environment**: Deployment environment

### Optional Fields
- **Description**: Host description
- **OS**: Operating system information
- **CPU Cores**: Number of CPU cores
- **Memory (GB)**: RAM size in GB
- **Disk (GB)**: Storage size in GB
- **Salt Minion ID**: Salt minion identifier
- **SSH Port**: SSH service port (default 22)

### System Fields
- **Created At**: Record creation timestamp
- **Updated At**: Last modification timestamp
- **Created By**: User who created the record
- **Status**: Host status (online/offline/discovered)

## Data Integrity Constraints

### Unique Constraints
- Hostname must be unique within environment
- IP address must be unique across all hosts
- Salt Minion ID should be unique when present

### Foreign Key Constraints
- Environment ID references environments table
- Group memberships reference host_groups table
- Tag associations reference host_tags table
- Created By references users table

### Validation Rules
- IP address format validation
- Hostname format validation (RFC standards)
- SSH port range validation (1-65535)
- Memory/disk values must be positive

## Performance Considerations

### Indexing Strategy
- Primary key on ID
- Unique index on IP address
- Index on hostname
- Index on environment_id
- Composite index on (environment_id, status)
- Full-text index on hostname and description

### Query Optimization
- Pagination for large result sets (default 20 items)
- Eager loading for related entities (groups, tags)
- Cached environment and group lists
- Asynchronous status updates

## Integration Points

### Salt Master Integration
- Real-time minion status monitoring
- Automated host discovery
- Grain data collection
- Remote command execution

### SSH Management
- Port detection and validation
- Key-based authentication support
- Connection testing and validation

### Audit Integration
- All host changes logged
- Bulk operations tracked
- User actions recorded with context
