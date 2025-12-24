## 1. Database Model and Migration
- [ ] 1.1 Create SystemSettings model
  - Create `backend/internal/models/settings.go`
  - Define SystemSettings struct with fields: ID, Key, Value, Category, Description, UpdatedAt, UpdatedBy
  - Use key-value storage pattern for flexibility
  - Support JSON values for complex settings
- [ ] 1.2 Create database migration
  - Create migration file `backend/migrations/add_system_settings_table.sql`
  - Define `system_settings` table schema
  - Add indexes for key and category lookups
- [ ] 1.3 Register model in AutoMigrate
  - Add SystemSettings to AutoMigrate list in `backend/internal/models/database.go`

## 2. Backend Repository Layer
- [ ] 2.1 Create settings repository
  - Create `backend/internal/repository/settings_repository.go`
  - Implement GetByKey, GetByCategory, GetAll, CreateOrUpdate, Delete methods
  - Handle key-value lookups and bulk operations

## 3. Backend Service Layer
- [ ] 3.1 Create settings service
  - Create `backend/internal/service/settings_service.go`
  - Implement GetSetting, GetSettingsByCategory, GetAllSettings, UpdateSetting methods
  - Add SaltConfig-specific methods: GetSaltConfig, UpdateSaltConfig
  - Handle configuration validation (URL format, timeout ranges, etc.)
  - Support encryption for sensitive values (passwords)

## 4. Backend Handler Layer
- [ ] 4.1 Create settings handler
  - Create `backend/internal/handler/settings_handler.go`
  - Implement GET `/api/v1/settings` - Get all settings
  - Implement GET `/api/v1/settings/:category` - Get settings by category
  - Implement GET `/api/v1/settings/key/:key` - Get setting by key
  - Implement PUT `/api/v1/settings/:key` - Update setting
  - Implement PUT `/api/v1/settings/salt` - Update Salt API configuration (convenience endpoint)
  - Add admin-only authorization middleware

## 5. Configuration Integration
- [ ] 5.1 Modify config loading to support database settings
  - Update `backend/internal/config/config.go` to load settings from database after initialization
  - Add method to reload Salt configuration from database
  - Support fallback to environment variables if database settings not available
  - Add configuration change detection and Salt client reinitialization

## 6. API Route Registration
- [ ] 6.1 Register settings routes
  - Add settings routes to `backend/cmd/api/main.go`
  - Initialize settings repository, service, and handler
  - Apply admin authorization middleware to settings routes

## 7. Frontend API Service
- [ ] 7.1 Create settings service
  - Create `frontend/src/services/settings.ts`
  - Define TypeScript interfaces for settings data
  - Implement getSettings, getSettingsByCategory, getSetting, updateSetting, updateSaltConfig methods
  - Handle API errors and responses

## 8. Frontend Settings Page
- [ ] 8.1 Create Settings page component
  - Create `frontend/src/pages/Settings.tsx`
  - Add page layout with Card component
  - Implement Salt API configuration form section
  - Add form fields: API URL, Username, Password, EAuth type, Timeout, Verify SSL checkbox
  - Add Save/Cancel buttons with validation
  - Display success/error messages
  - Load current settings on page mount

## 9. Frontend Routing and Menu
- [ ] 9.1 Add Settings route
  - Add `/settings` route to `frontend/src/App.tsx`
  - Import and use Settings component in route configuration
  - Ensure route is protected with PrivateRoute and admin role check
- [ ] 9.2 Add Settings menu item
  - Add Settings menu item to `frontend/src/components/MainLayout.tsx`
  - Use appropriate icon (e.g., `SettingOutlined`)
  - Place menu item at the bottom or in a "System" section
  - Set menu key to `/settings`
  - Show menu item only for admin users

## 10. Initial Data Seeding
- [ ] 10.1 Create default Salt settings
  - Add logic to seed default Salt API settings from environment variables during initial setup
  - Populate database with current environment variable values if database is empty
  - Ensure settings are initialized on first application startup

## 11. Testing
- [ ] 11.1 Backend tests
  - Write unit tests for settings repository
  - Write unit tests for settings service (including validation)
  - Write integration tests for settings handler endpoints
  - Test configuration loading with database settings
- [ ] 11.2 Frontend tests
  - Write tests for Settings page component
  - Test form validation
  - Test settings loading and saving

## 12. Validation
- [ ] 12.1 Code validation
  - Run `go build` to verify backend compiles without errors
  - Run `npm run build` to verify frontend compiles without errors
  - Fix any linting errors
- [ ] 12.2 Functional validation
  - Verify settings page is accessible only to admin users
  - Verify Salt API settings can be viewed and updated
  - Verify configuration changes are persisted to database
  - Verify Salt client reinitialization works after settings update
  - Verify environment variable fallback works when database settings unavailable
- [ ] 12.3 OpenSpec validation
  - Run `openspec validate add-system-settings-page --strict`
  - Fix any validation errors

