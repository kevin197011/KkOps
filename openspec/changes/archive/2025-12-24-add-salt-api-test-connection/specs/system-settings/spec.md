## MODIFIED Requirements

### Requirement: Salt API Configuration Management
The system SHALL allow administrators to configure Salt API settings through the UI, including API URL, authentication credentials, and connection parameters. The system SHALL provide a test connection feature to verify Salt API connectivity before saving configuration.

#### Scenario: View Salt API configuration
- **WHEN** an administrator views the Salt API configuration section
- **THEN** the system SHALL display current Salt API settings
- **AND** the system SHALL show fields: API URL, Username, Password (masked), EAuth type, Timeout, Verify SSL
- **AND** the system SHALL load settings from the database (with fallback to environment variables)
- **AND** the system SHALL display a "测试连接" (Test Connection) button

#### Scenario: Update Salt API configuration
- **WHEN** an administrator updates Salt API configuration and saves
- **THEN** the system SHALL validate the input (URL format, numeric ranges, required fields)
- **AND** the system SHALL save the configuration to the database
- **AND** the system SHALL encrypt sensitive values (passwords) before storage
- **AND** the system SHALL update the Salt client configuration
- **AND** the system SHALL display a success message
- **AND** the system SHALL log the configuration change for audit purposes

#### Scenario: Test Salt API connection with form values
- **WHEN** an administrator clicks the "测试连接" button while viewing Salt API configuration form
- **THEN** the system SHALL use the current form values (including unsaved changes) to test the connection
- **AND** the system SHALL attempt to authenticate with the Salt API using the provided credentials
- **AND** the system SHALL display a loading indicator during the test
- **AND** if the connection is successful, the system SHALL display a success message
- **AND** if the connection fails, the system SHALL display an error message with details about the failure (e.g., "认证失败", "无法连接到Salt API", "网络错误")

#### Scenario: Test Salt API connection with saved configuration
- **WHEN** an administrator clicks the "测试连接" button without modifying form values
- **THEN** the system SHALL use the saved configuration from the database (or environment variables) to test the connection
- **AND** the system SHALL attempt to authenticate with the Salt API
- **AND** the system SHALL display connection test results with appropriate success or error messages

#### Scenario: Test connection failure handling
- **WHEN** the test connection fails due to network error, authentication failure, or invalid configuration
- **THEN** the system SHALL display a clear error message indicating the type of failure
- **AND** the system SHALL not save any configuration changes
- **AND** the system SHALL allow the administrator to correct the configuration and test again

#### Scenario: Salt API configuration validation
- **WHEN** an administrator submits invalid Salt API configuration
- **THEN** the system SHALL validate URL format, required fields, and numeric ranges
- **AND** the system SHALL display validation error messages
- **AND** the system SHALL prevent saving invalid configuration

#### Scenario: Salt API configuration error handling
- **WHEN** an error occurs while saving Salt API configuration
- **THEN** the system SHALL display an error message to the administrator
- **AND** the system SHALL preserve the form data
- **AND** the system SHALL log the error for troubleshooting

