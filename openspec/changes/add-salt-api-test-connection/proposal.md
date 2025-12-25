# Change: Add Salt API Test Connection Button

## Why
Currently, administrators can configure Salt API settings through the system settings page, but there is no way to verify if the configuration is correct before saving. Administrators may enter incorrect credentials or URLs, and only discover the issue later when the Salt API is actually used. This creates a poor user experience and wastes time troubleshooting configuration issues.

By adding a test connection button, administrators can:
- Verify Salt API connection immediately after entering configuration
- Catch configuration errors before saving
- Validate credentials and network connectivity
- Get immediate feedback on whether the Salt API is accessible

## What Changes
- **ADDED**: Test connection endpoint in settings handler (`POST /api/v1/settings/salt/test`)
- **ADDED**: Test connection method in Salt client/service
- **ADDED**: Test connection button in Settings page UI
- **ADDED**: Connection test feedback (success/error messages)
- **MODIFIED**: Settings page to include test connection functionality

## Impact
- **Affected specs**: `system-settings` capability (MODIFIED)
- **Affected code**:
  - `backend/internal/handler/settings_handler.go` - Add test connection handler
  - `backend/internal/salt/client.go` - Add test connection method (or use existing authenticate)
  - `backend/internal/service/settings_service.go` - Add test connection service method (optional)
  - `frontend/src/pages/Settings.tsx` - Add test connection button and handler
  - `frontend/src/services/settings.ts` - Add test connection API service method
  - `backend/cmd/api/main.go` - Register test connection route
- **User experience**: Administrators can test Salt API connection before saving configuration, improving configuration reliability and reducing errors

