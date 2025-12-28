## ADDED Requirements

### Requirement: System Settings Management
The system SHALL provide a centralized system settings page where administrators can view and update system configuration.

#### Scenario: Access system settings page
- **WHEN** an administrator navigates to the system settings page
- **THEN** the system SHALL display the settings page
- **AND** the system SHALL show different configuration categories
- **AND** the system SHALL allow administrators to view and update settings

#### Scenario: Non-admin access restriction
- **WHEN** a non-administrator user attempts to access the system settings page
- **THEN** the system SHALL deny access
- **AND** the system SHALL redirect to an appropriate page or show an error message

### Requirement: Salt API Configuration Management
The system SHALL allow administrators to configure Salt API settings through the UI, including API URL, authentication credentials, and connection parameters.

#### Scenario: View Salt API configuration
- **WHEN** an administrator views the Salt API configuration section
- **THEN** the system SHALL display current Salt API settings
- **AND** the system SHALL show fields: API URL, Username, Password (masked), EAuth type, Timeout, Verify SSL
- **AND** the system SHALL load settings from the database (with fallback to environment variables)

#### Scenario: Update Salt API configuration
- **WHEN** an administrator updates Salt API configuration and saves
- **THEN** the system SHALL validate the input (URL format, numeric ranges, required fields)
- **AND** the system SHALL save the configuration to the database
- **AND** the system SHALL encrypt sensitive values (passwords) before storage
- **AND** the system SHALL update the Salt client configuration
- **AND** the system SHALL display a success message
- **AND** the system SHALL log the configuration change for audit purposes

#### Scenario: Salt API configuration validation
- **WHEN** an administrator submits invalid Salt API configuration
- **THEN** the system SHALL display validation errors
- **AND** the system SHALL NOT save invalid configuration
- **AND** the system SHALL highlight fields with validation errors

#### Scenario: Salt API configuration error handling
- **WHEN** an error occurs while saving Salt API configuration
- **THEN** the system SHALL display an error message to the administrator
- **AND** the system SHALL preserve the form data
- **AND** the system SHALL log the error for troubleshooting

### Requirement: Settings Persistence
The system SHALL store system settings in the database and load them at application startup.

#### Scenario: Settings loading priority
- **WHEN** the application starts
- **THEN** the system SHALL load settings from the database
- **AND** the system SHALL use database settings if available
- **AND** the system SHALL fallback to environment variables if database settings are not available
- **AND** the system SHALL initialize services (e.g., Salt client) with the loaded configuration

#### Scenario: Settings initialization
- **WHEN** the system settings table is empty on first startup
- **THEN** the system SHALL seed default settings from environment variables
- **AND** the system SHALL populate the database with current environment variable values

### Requirement: Settings API Endpoints
The system SHALL provide REST API endpoints for retrieving and updating system settings.

#### Scenario: Get all settings
- **WHEN** an administrator requests all settings via API
- **THEN** the system SHALL return all system settings
- **AND** the system SHALL mask sensitive values (passwords) in the response
- **AND** the system SHALL require administrator authentication

#### Scenario: Get settings by category
- **WHEN** an administrator requests settings for a specific category (e.g., "salt")
- **THEN** the system SHALL return only settings in that category
- **AND** the system SHALL mask sensitive values in the response

#### Scenario: Update a setting
- **WHEN** an administrator updates a setting via API
- **THEN** the system SHALL validate the setting value
- **AND** the system SHALL encrypt sensitive values before saving
- **AND** the system SHALL save the setting to the database
- **AND** the system SHALL record who updated the setting and when
- **AND** the system SHALL return the updated setting value

#### Scenario: Update Salt configuration
- **WHEN** an administrator updates Salt API configuration via the convenience endpoint
- **THEN** the system SHALL update all Salt-related settings atomically
- **AND** the system SHALL validate all Salt settings together
- **AND** the system SHALL reinitialize the Salt client with new configuration

