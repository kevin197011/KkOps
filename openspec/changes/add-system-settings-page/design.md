# Design: System Settings Page

## Context
The system currently manages Salt API configuration through environment variables, which requires manual file editing and service restarts. Administrators need a UI-based configuration management system that:
1. Allows configuring Salt API settings through the web interface
2. Stores settings persistently in the database
3. Supports dynamic configuration updates
4. Is extensible for future system settings

## Goals / Non-Goals

### Goals
- Provide a centralized system settings page for administrators
- Enable UI-based Salt API configuration management
- Store settings in database for persistence
- Support configuration validation and error handling
- Design extensible architecture for future settings categories

### Non-Goals
- Real-time configuration reload without service restart (initial version)
- Settings versioning/history (can be added later)
- Multi-environment settings (production vs development)
- Settings import/export functionality

## Decisions

### Decision: Key-Value Storage Pattern
- **Rationale**: Flexible and extensible for various setting types. Easy to add new settings without schema changes.
- **Implementation**: Use a single `system_settings` table with `key`, `value`, `category` columns
- **Trade-offs**:
  - ✅ Flexible and easy to extend
  - ✅ No schema changes needed for new settings
  - ❌ Less type safety compared to structured tables
  - ❌ Requires careful key naming conventions

### Decision: Database Settings Override Environment Variables
- **Rationale**: Database settings provide runtime configurability, while environment variables serve as defaults/fallback
- **Implementation**: Load database settings after environment variables, with database values taking precedence
- **Trade-offs**:
  - ✅ Allows runtime configuration changes
  - ✅ Environment variables still serve as sensible defaults
  - ❌ Need to handle configuration reload/reinitialization

### Decision: Salt Client Reinitialization on Settings Update
- **Rationale**: Salt client holds connection state and credentials that need to be refreshed
- **Implementation**: Reinitialize Salt client when Salt settings are updated via API
- **Trade-offs**:
  - ✅ Ensures new settings are immediately effective
  - ✅ Clear separation of concerns
  - ❌ May interrupt in-flight Salt API calls (acceptable for configuration changes)

### Decision: Password Encryption in Database
- **Rationale**: Salt API passwords are sensitive and should not be stored in plain text
- **Implementation**: Encrypt password values before storing in database, decrypt when loading
- **Trade-offs**:
  - ✅ Enhanced security for sensitive data
  - ✅ Follows security best practices
  - ❌ Adds complexity to encryption key management

### Decision: Admin-Only Access
- **Rationale**: System settings are sensitive and should only be modifiable by administrators
- **Implementation**: Apply admin role check middleware to all settings endpoints and UI
- **Trade-offs**:
  - ✅ Prevents unauthorized configuration changes
  - ✅ Clear access control
  - ❌ Requires proper role/permission setup

## Data Model

### SystemSettings Model
```go
type SystemSettings struct {
    ID          uint64    `gorm:"primaryKey"`
    Key         string    `gorm:"uniqueIndex;size:100;not null"`
    Value       string    `gorm:"type:text;not null"`
    Category    string    `gorm:"size:50;index;not null"` // e.g., "salt", "email", "notification"
    Description string    `gorm:"type:text"`
    UpdatedAt   time.Time
    UpdatedBy   uint64    `gorm:"index"` // User ID who last updated
}
```

### Setting Keys (Initial)
- `salt.api_url` - Salt API URL
- `salt.username` - Salt API username
- `salt.password` - Salt API password (encrypted)
- `salt.eauth` - Salt API eauth type (default: "pam")
- `salt.timeout` - Salt API timeout in seconds (default: 30)
- `salt.verify_ssl` - Verify SSL certificates (default: true)

## API Design

### Endpoints
- `GET /api/v1/settings` - Get all settings
- `GET /api/v1/settings/:category` - Get settings by category (e.g., "salt")
- `GET /api/v1/settings/key/:key` - Get specific setting by key
- `PUT /api/v1/settings/:key` - Update a setting
- `PUT /api/v1/settings/salt` - Update all Salt settings at once (convenience endpoint)

### Request/Response Examples

**Get Salt Settings:**
```http
GET /api/v1/settings/salt
```
```json
{
  "settings": {
    "salt.api_url": "http://192.168.56.10:8000",
    "salt.username": "salt",
    "salt.password": "***" (masked),
    "salt.eauth": "pam",
    "salt.timeout": "30",
    "salt.verify_ssl": "false"
  }
}
```

**Update Salt Settings:**
```http
PUT /api/v1/settings/salt
Content-Type: application/json
```
```json
{
  "api_url": "http://192.168.56.10:8000",
  "username": "salt",
  "password": "newpassword",
  "eauth": "pam",
  "timeout": 30,
  "verify_ssl": false
}
```

## Frontend Design

### Settings Page Layout
- Page title: "系统设置"
- Tabs or sections for different setting categories:
  - "Salt API 配置" (initial focus)
  - Future: "邮件配置", "通知配置", etc.

### Salt API Configuration Form
- Form fields with labels and descriptions:
  - API URL (Text input, required, URL validation)
  - Username (Text input, required)
  - Password (Password input, required, with show/hide toggle)
  - EAuth Type (Select dropdown: pam, ldap, etc.)
  - Timeout (Number input, seconds, min: 1, max: 300)
  - Verify SSL (Checkbox)
- Actions: "保存" (Save) and "取消" (Cancel) buttons
- Validation feedback and error messages
- Success notification after save

## Configuration Loading Flow

1. Application startup:
   - Load environment variables (defaults)
   - Initialize database connection
   - Load settings from database
   - Override config with database values
   - Initialize Salt client with final configuration

2. Settings update:
   - User updates settings via UI
   - API validates and saves to database
   - API triggers Salt client reinitialization (or signals reload needed)
   - Response confirms update success

## Security Considerations

- **Password Encryption**: Use existing encryption utilities to encrypt/decrypt password values
- **Access Control**: Admin-only access enforced at API and UI levels
- **Input Validation**: Validate URL formats, numeric ranges, boolean values
- **Audit Logging**: Log all settings changes for audit trail

## Migration Strategy

1. Create `system_settings` table via migration
2. Seed initial Salt settings from environment variables if table is empty
3. Update config loading logic to check database first, fallback to env vars
4. Deploy backend changes
5. Deploy frontend changes
6. Verify settings page functionality

## Future Extensibility

The key-value pattern allows easy addition of new setting categories:
- Email settings (SMTP server, credentials, etc.)
- Notification settings (webhook URLs, alert thresholds, etc.)
- Integration settings (third-party API keys, etc.)
- UI preferences (theme, language, etc.)

Each new category can be added by:
1. Defining new setting keys (e.g., `email.smtp_host`, `email.smtp_port`)
2. Adding UI section in Settings page
3. Implementing category-specific validation in service layer

## Risks / Trade-offs

- **Risk**: Configuration changes might break Salt connectivity
  - **Mitigation**: Validate settings before saving, test connection on save (optional)
- **Risk**: Database settings might get out of sync with environment variables
  - **Mitigation**: Environment variables serve as defaults, database takes precedence (clear hierarchy)
- **Risk**: Password encryption key management
  - **Mitigation**: Reuse existing encryption utilities and key management approach

