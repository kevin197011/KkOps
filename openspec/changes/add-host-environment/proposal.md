# Change: Add Environment Field to Host Management and Reorder Project Column

## Why
1. **Project Column Visibility**: Currently, the project information is not displayed in the host management table, making it difficult to quickly identify which project a host belongs to. Moving the project column to appear right after the ID column will improve visibility and usability.

2. **Environment Classification**: Hosts need to be classified by environment (e.g., `uat`, `prod`, `dev`, `staging`) to better organize and manage infrastructure across different deployment stages. This is a common requirement in operations management systems.

## What Changes
- **Frontend**: Add "项目" (Project) column after "ID" column in host management table
- **Backend**: Add `Environment` field to `Host` model (string, enum: dev, uat, staging, prod, etc.)
- **Frontend**: Add environment field to host creation/edit form with dropdown selection
- **Frontend**: Add "环境" (Environment) column to host management table
- **Backend**: Update database schema to include `environment` field
- **Backend**: Add validation for environment field values

## Impact
- **Affected specs**: `host-management`
- **Affected code**: 
  - `backend/internal/models/host.go` - Add Environment field
  - `backend/internal/repository/host_repository.go` - Filter by environment (if needed)
  - `backend/internal/service/host_service.go` - Environment validation
  - `backend/internal/handler/host_handler.go` - Environment parameter handling
  - `frontend/src/services/host.ts` - Add environment to Host interface
  - `frontend/src/pages/Hosts.tsx` - Add project column after ID, add environment field and column
- **Database**: Migration to add `environment` column to `hosts` table

