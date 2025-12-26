# Formula Management Specification

## Overview
Formula管理系统提供基于SaltStack Formula的自动化部署和配置管理，支持Formula仓库管理、参数配置和部署执行。

## Requirements

### Requirement: Formula Repository Management
The system SHALL manage SaltStack Formula repositories with Git integration.

#### Scenario: Repository Configuration
- **WHEN** administrator configures Formula repository
- **THEN** system validates Git URL and credentials
- **AND** clones repository to local storage
- **AND** scans for available Formulas

#### Scenario: Repository Synchronization
- **WHEN** system syncs repository changes
- **THEN** performs Git pull operations
- **AND** discovers new or modified Formulas
- **AND** updates Formula metadata
- **AND** maintains version history

### Requirement: Formula Discovery and Parsing
The system SHALL automatically discover and parse SaltStack Formulas.

#### Scenario: Formula Detection
- **WHEN** scanning repository
- **THEN** identifies directories containing `init.sls`
- **AND** parses Formula structure and metadata
- **AND** extracts parameter definitions
- **AND** categorizes Formulas by purpose

#### Scenario: Metadata Extraction
- **WHEN** parsing Formula files
- **THEN** extracts pillar schema definitions
- **AND** identifies required and optional parameters
- **AND** determines Formula dependencies
- **AND** generates Formula documentation

### Requirement: Formula Parameter Management
The system SHALL provide comprehensive parameter configuration for Formulas.

#### Scenario: Parameter Definition
- **WHEN** parsing Formula metadata
- **THEN** extracts parameter schemas
- **AND** defines parameter types and constraints
- **AND** provides default values
- **AND** validates parameter requirements

#### Scenario: Parameter Validation
- **WHEN** user configures Formula parameters
- **THEN** validates parameter types and formats
- **AND** enforces required parameter constraints
- **AND** provides parameter documentation
- **AND** suggests valid values

### Requirement: Formula Deployment Execution
The system SHALL execute Formula deployments through SaltStack.

#### Scenario: Deployment Creation
- **WHEN** user initiates Formula deployment
- **THEN** validates Formula selection and parameters
- **AND** selects target hosts
- **AND** creates deployment record
- **AND** submits Salt state.apply job

#### Scenario: Deployment Monitoring
- **WHEN** deployment executes
- **THEN** monitors Salt job progress
- **AND** collects per-host results
- **AND** provides real-time status updates
- **AND** handles deployment failures

#### Scenario: Deployment Results
- **WHEN** deployment completes
- **THEN** aggregates execution results
- **AND** provides detailed success/failure analysis
- **AND** stores results for historical reference
- **AND** generates deployment reports

### Requirement: Formula Template Management
The system SHALL support reusable Formula configuration templates.

#### Scenario: Template Creation
- **WHEN** user saves successful deployment as template
- **THEN** stores Formula selection and parameters
- **AND** associates with user and project
- **AND** provides template versioning
- **AND** supports template sharing

#### Scenario: Template Application
- **WHEN** user applies Formula template
- **THEN** loads saved configuration
- **AND** allows parameter customization
- **AND** validates template compatibility
- **AND** applies to selected hosts

## Formula Structure

### Directory Layout
```
formula-repository/
├── base/
│   ├── timezone/
│   │   ├── init.sls
│   │   ├── map.jinja
│   │   └── files/
│   └── users/
│       ├── init.sls
│       └── files/
├── middleware/
│   ├── nginx/
│   │   ├── init.sls
│   │   ├── config.sls
│   │   ├── install.sls
│   │   ├── service.sls
│   │   └── files/
│   └── redis/
│       ├── init.sls
│       ├── config.sls
│       └── files/
└── app/
    └── web-app/
        ├── init.sls
        ├── config.sls
        ├── deploy.sls
        └── files/
```

### Formula Categories
- **base**: System-level configurations (timezone, users, sysctl)
- **middleware**: Application servers (nginx, redis, mysql)
- **runtime**: Runtime environments (java, nodejs, python)
- **app**: Application-specific deployments

## Parameter Schema

### Parameter Types
- **string**: Text values with optional validation
- **number**: Numeric values with range constraints
- **boolean**: True/false values
- **array**: List of values
- **object**: Complex nested structures

### Parameter Definition
```yaml
# Pillar schema example
web_app:
  domain: example.com
  port: 8080
  ssl_enabled: true
  database:
    host: db.example.com
    port: 5432
    name: webapp
```

### Validation Rules
- **Required**: Mandatory parameters
- **Type**: Data type validation
- **Pattern**: Regular expression validation
- **Range**: Numeric range constraints
- **Enum**: Allowed value enumeration
- **Dependencies**: Conditional parameter requirements

## Deployment Lifecycle

### Pre-deployment Validation
1. **Formula Validation**: Verify Formula exists and is accessible
2. **Parameter Validation**: Validate all required parameters provided
3. **Host Validation**: Check target hosts are accessible
4. **Dependency Check**: Verify Formula dependencies are satisfied
5. **Permission Check**: Validate user has deployment permissions

### Deployment Execution
1. **Pillar Generation**: Create Salt pillar data from parameters
2. **Target Selection**: Generate Salt target expressions
3. **Job Submission**: Submit state.apply job to Salt Master
4. **Progress Monitoring**: Track job execution progress
5. **Result Collection**: Gather execution results from all hosts

### Post-deployment Processing
1. **Result Analysis**: Categorize success/failure by host
2. **Error Reporting**: Provide detailed error information
3. **Audit Logging**: Record deployment details and results
4. **Notification**: Alert stakeholders of deployment status

## Git Integration

### Repository Operations
- **Clone**: Initial repository download
- **Pull**: Update to latest changes
- **Branch**: Support for different branches
- **Tag**: Version-specific deployments
- **Submodule**: Support for nested repositories

### Change Detection
- **File Monitoring**: Track Formula file changes
- **Metadata Updates**: Automatic Formula metadata refresh
- **Version Tracking**: Formula version history
- **Conflict Resolution**: Handle merge conflicts

## Integration Points

### Salt Master Integration
- Salt API communication for job execution
- Pillar data transmission
- Result parsing and aggregation
- Error code translation

### Host Management Integration
- Host selection and filtering
- Environment-based deployment
- Group and tag targeting
- Host capability validation

### User Management Integration
- Permission-based access control
- User-specific templates
- Deployment attribution
- Audit trail creation

### Audit Integration
- Deployment operation logging
- Configuration change tracking
- Security event monitoring
- Compliance reporting

## Performance Considerations

### Deployment Optimization
- **Parallel Execution**: Concurrent host deployment
- **Batch Processing**: Group similar deployments
- **Caching**: Formula metadata and file caching
- **Incremental Updates**: Deploy only changed components

### Resource Management
- **Memory Limits**: Large deployment result handling
- **Timeout Controls**: Deployment execution time limits
- **Concurrency Limits**: Maximum simultaneous deployments
- **Queue Management**: Deployment job queuing

## Security Considerations

### Access Control
- **Formula Permissions**: User Formula access validation
- **Host Restrictions**: Deployment target authorization
- **Parameter Security**: Sensitive parameter encryption
- **Audit Compliance**: Complete deployment audit trail

### Data Protection
- **Pillar Security**: Sensitive pillar data encryption
- **Result Sanitization**: Remove sensitive information from logs
- **Access Logging**: Comprehensive security event tracking
- **Compliance**: Regulatory requirement adherence

## Operational Requirements

### Monitoring and Alerting
- **Deployment Metrics**: Success rates and performance
- **Formula Health**: Repository and Formula status monitoring
- **Resource Usage**: Deployment resource consumption tracking
- **Error Patterns**: Common deployment failure analysis

### Backup and Recovery
- **Formula Archives**: Formula repository backups
- **Deployment History**: Critical deployment data preservation
- **Configuration Backup**: Pillar and parameter backups
- **Disaster Recovery**: Service restoration procedures

### Maintenance
- **Formula Updates**: Automated Formula update detection
- **Dependency Management**: Formula dependency validation
- **Cleanup Procedures**: Old deployment data removal
- **Health Checks**: System and Formula health monitoring
