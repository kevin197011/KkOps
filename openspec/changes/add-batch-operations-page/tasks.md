# Implementation Tasks

## 1. Database Design and Migration
- [x] 1.1 Create `batch_operations` table migration ✅
  - Define table schema with all required fields
  - Add indexes for performance (started_by, status, started_at)
  - Add foreign key constraint to users table
- [x] 1.2 Create `command_templates` table migration ✅
  - Define table schema with all required fields
  - Add indexes for performance (category, created_by)
  - Add foreign key constraint to users table
- [x] 1.3 Create default command templates (built-in templates) ✅
  - Insert system information templates (loadavg, meminfo, disk.usage, cpuinfo)
  - Insert network information templates (interfaces, active_tcp)
  - Insert process management templates (ps.procs, service.status)
  - Insert test connectivity template (test.ping)
  - Note: File operation templates can be added later as needed

## 2. Backend Models
- [x] 2.1 Create BatchOperation model ✅
  - Define struct with all fields
  - Add GORM tags and relationships
  - Add JSON tags for API responses
- [x] 2.2 Create CommandTemplate model ✅
  - Define struct with all fields
  - Add GORM tags and relationships
  - Add JSON tags for API responses
- [x] 2.3 Register models in AutoMigrate ✅
  - Add BatchOperation to AutoMigrate list
  - Add CommandTemplate to AutoMigrate list

## 3. Backend Repository Layer
- [x] 3.1 Create BatchOperationsRepository ✅
  - Implement Create, Get, List, Update methods
  - Implement filtering by user, status, date range
  - Implement pagination support
- [x] 3.2 Create CommandTemplateRepository ✅
  - Implement Create, Get, List, Update, Delete methods
  - Implement filtering by category, user, public/private
  - Implement increment usage count method

## 4. Backend Service Layer
- [x] 4.1 Create BatchOperationsService ✅
  - Implement CreateOperation method
  - Implement ExecuteOperation method (async with goroutine)
  - Implement GetOperationStatus method
  - Implement GetOperationResults method
  - Implement CancelOperation method
  - Implement ListOperations method with filtering
  - Implement GetOperation method
  - Integrate with Salt Client for command execution
  - Implement result collection and status updates
  - Implement timeout handling (default 5 minutes)
- [x] 4.2 Create CommandTemplateService ✅
  - Implement CreateTemplate method
  - Implement UpdateTemplate method
  - Implement DeleteTemplate method
  - Implement ListTemplates method with filtering
  - Implement GetTemplate method
  - Implement IncrementUsageCount method

## 5. Backend Handler Layer
- [x] 5.1 Create BatchOperationsHandler ✅
  - Implement POST `/api/v1/batch-operations` - Create and execute operation
  - Implement GET `/api/v1/batch-operations` - List operations with pagination
  - Implement GET `/api/v1/batch-operations/:id` - Get operation details
  - Implement GET `/api/v1/batch-operations/:id/status` - Get operation status
  - Implement GET `/api/v1/batch-operations/:id/results` - Get operation results
  - Implement POST `/api/v1/batch-operations/:id/cancel` - Cancel operation
  - Add request validation
  - Add error handling
  - Add RBAC permission checks
- [x] 5.2 Create CommandTemplateHandler ✅
  - Implement POST `/api/v1/command-templates` - Create template
  - Implement GET `/api/v1/command-templates` - List templates
  - Implement GET `/api/v1/command-templates/:id` - Get template details
  - Implement PUT `/api/v1/command-templates/:id` - Update template
  - Implement DELETE `/api/v1/command-templates/:id` - Delete template
  - Add request validation
  - Add error handling
  - Add RBAC permission checks

## 6. Backend Route Registration
- [x] 6.1 Register batch operations routes ✅
  - Add routes to `backend/cmd/api/main.go`
  - Apply RBAC middleware for permission control
  - Initialize handler with dependencies
- [x] 6.2 Register command template routes ✅
  - Add routes to `backend/cmd/api/main.go`
  - Apply RBAC middleware for permission control
  - Initialize handler with dependencies

## 7. Frontend API Service
- [x] 7.1 Create batch operations service (`frontend/src/services/batchOperations.ts`) ✅
  - Define TypeScript interfaces for BatchOperation, CommandTemplate
  - Implement createOperation method
  - Implement listOperations method
  - Implement getOperation method
  - Implement getOperationStatus method
  - Implement getOperationResults method
  - Implement cancelOperation method
  - Implement retryOperation method
  - Add error handling
- [x] 7.2 Create command template service (`frontend/src/services/commandTemplate.ts`) ✅
  - Define TypeScript interfaces for CommandTemplate
  - Implement createTemplate method
  - Implement listTemplates method
  - Implement getTemplate method
  - Implement updateTemplate method
  - Implement deleteTemplate method
  - Add error handling

## 8. Frontend Components
- [x] 8.1 Create HostSelector component (`frontend/src/components/HostSelector.tsx`) ✅
  - Implement project filter
  - Implement group filter
  - Implement tag filter
  - Implement status filter
  - Implement host list with multi-select
  - Implement select all / clear selection
  - Display selected host count
  - Add search functionality
- [x] 8.2 Create CommandTemplateManager component ✅
  - Integrated into BatchOperations page (template selection and management)
  - Template creation form (via modal)
  - Template selection in command configuration
  - Template list display
- [x] 8.3 Create ResultViewer component (`frontend/src/components/ResultViewer.tsx`) ✅
  - Implement result table display
  - Implement status indicators (success/failed)
  - Implement expandable result details
  - Implement result filtering (by status)
  - Implement result export (CSV/JSON)
  - Implement real-time status updates
  - Display success/failure statistics

## 9. Frontend Page
- [x] 9.1 Create BatchOperations page (`frontend/src/pages/BatchOperations.tsx`) ✅
  - Implement page layout with sections:
    - Top toolbar (operation name, quick actions, save template)
    - Left panel (host selector)
    - Middle panel (command configuration)
    - Right panel (result viewer)
    - Bottom panel (operation history)
  - Implement command type selection (custom/template/built-in)
  - Implement command input/selection
  - Implement parameter configuration
  - Implement execute button and operation creation
  - Implement operation status polling
  - Implement operation history display
  - Implement retry functionality
  - Add loading states
  - Add error handling
  - Add success/error notifications

## 10. Frontend Routing and Menu
- [x] 10.1 Add BatchOperations route ✅
  - Add `/batch-operations` route to `frontend/src/App.tsx`
  - Import and use BatchOperations component
  - Ensure route is protected with PrivateRoute and permission check
- [x] 10.2 Add BatchOperations menu item ✅
  - Add menu item to `frontend/src/components/MainLayout.tsx`
  - Use appropriate icon (`ThunderboltOutlined`)
  - Place menu item in appropriate section (after deployments)
  - Set menu key to `/batch-operations`
  - Menu item visible to all authenticated users (permission check can be added later)

## 11. Built-in Command Templates
- [x] 11.1 Create system information templates ✅
  - View system load (`status.loadavg`)
  - View memory usage (`status.meminfo`)
  - View disk usage (`disk.usage`)
  - View CPU info (`status.cpuinfo`)
  - Test connectivity (`test.ping`)
- [x] 11.2 Create network information templates ✅
  - View network interfaces (`network.interfaces`)
  - View active TCP connections (`network.active_tcp`)
- [x] 11.3 Create process management templates ✅
  - View process list (`ps.procs`)
  - View service status (`service.status`)
- [ ] 11.4 Create file operation templates (optional, can be added later)
  - Read file (`file.read`)
  - Check file exists (`file.file_exists`)

## 12. Testing
- [ ] 12.1 Backend unit tests
  - Test BatchOperationsService methods
  - Test CommandTemplateService methods
  - Test repository methods
  - Test handler endpoints
  - Test error handling
- [ ] 12.2 Backend integration tests
  - Test batch operation execution flow
  - Test Salt API integration
  - Test result collection
  - Test timeout handling
- [ ] 12.3 Frontend unit tests
  - Test BatchOperations page component
  - Test HostSelector component
  - Test CommandTemplateManager component
  - Test ResultViewer component
  - Test API service methods
- [ ] 12.4 Manual testing
  - Test host selection and filtering
  - Test command execution on multiple hosts
  - Test result display and export
  - Test operation history
  - Test template management
  - Test error scenarios
  - Test permission control

## 13. Validation
- [x] 13.1 Code validation ✅
  - Run `go build` to verify backend compiles without errors ✅
  - Run `npm run build` to verify frontend compiles without errors (Note: `@dnd-kit/sortable` dependency issue unrelated to this change)
  - Fix any linting errors ✅
- [ ] 13.2 Functional validation
  - Verify batch operations page is accessible with proper permissions
  - Verify host selection works correctly
  - Verify command execution works on multiple hosts
  - Verify results are displayed correctly
  - Verify operation history is maintained
  - Verify template management works
  - Verify permission control works
  - Verify audit logging works
- [ ] 13.3 OpenSpec validation
  - Run `openspec validate add-batch-operations-page --strict`
  - Fix any validation errors

