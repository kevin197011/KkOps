# Change: Add System Settings Page

## Why
Currently, Salt API configuration is managed through environment variables (in `.env` file or `docker-compose.yml`), which requires manual configuration and server restarts to change settings. This creates friction for administrators who need to adjust Salt API settings. Additionally, the system needs a centralized location for configuration management that can be extended in the future for other system-level settings (e.g., email configuration, notification settings, etc.).

By providing a system settings page, administrators can:
- Configure Salt API settings through the UI without manual file editing
- View current configuration in one place
- Update settings dynamically (with appropriate validation and service restarts)
- Extend the page for future configuration needs

## What Changes
- **ADDED**: SystemSettings database model to store system configuration in the database
- **ADDED**: Settings repository, service, and handler for managing system settings
- **ADDED**: REST API endpoints for retrieving and updating system settings
- **ADDED**: System Settings page in the frontend (`/settings`)
- **ADDED**: Salt API configuration section in the settings page (API URL, Username, Password, EAuth type, Timeout, Verify SSL)
- **ADDED**: Settings menu item in the main navigation (admin-only)
- **MODIFIED**: Salt client initialization to support dynamic configuration updates (with service restart or reload)
- **MODIFIED**: Configuration loading to prioritize database settings over environment variables

## Impact
- **Affected specs**: New `system-settings` capability
- **Affected code**:
  - `backend/internal/models/` - Add SystemSettings model
  - `backend/internal/repository/` - Add settings repository
  - `backend/internal/service/` - Add settings service
  - `backend/internal/handler/` - Add settings handler
  - `backend/internal/config/` - Modify config loading to support database settings
  - `backend/cmd/api/main.go` - Register settings routes
  - `frontend/src/pages/` - Add Settings page component
  - `frontend/src/services/` - Add settings API service
  - `frontend/src/App.tsx` - Add `/settings` route
  - `frontend/src/components/MainLayout.tsx` - Add settings menu item
- **Database changes**: New `system_settings` table
- **User experience**: Administrators can configure Salt API settings through the UI, improving operational efficiency

