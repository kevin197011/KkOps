# Change: Enhance User and Role Management Functionality

## Why

Currently, the user management and role management pages only display lists but lack interactive functionality. Buttons for creating, editing, and deleting users/roles have no onClick handlers, making them non-functional. This creates a poor user experience where administrators cannot actually manage users and roles through the UI. The backend APIs are complete, but the frontend implementation is incomplete.

This change will:
- Enable administrators to create, edit, and delete users through the UI
- Enable administrators to create, edit, and delete roles through the UI
- Provide a consistent user experience matching other management pages (projects, categories, tags)
- Complete the CRUD functionality for user and role management

## What Changes

- **ADDED**: User API client (`frontend/src/api/user.ts`) with create, update, delete, and list methods
- **ADDED**: Role API client (`frontend/src/api/role.ts`) with create, update, delete, and list methods
- **MODIFIED**: UserList page to include create/edit/delete functionality with Modal forms
- **MODIFIED**: RoleList page to include create/edit/delete functionality with Modal forms
- **ADDED**: Form validation for user and role creation/editing
- **ADDED**: Confirmation dialogs for delete operations
- **ADDED**: Success/error message feedback for all operations

## Impact

- Affected code: Frontend user and role management pages, API clients
- Affected files:
  - `frontend/src/api/user.ts` (new file - API client)
  - `frontend/src/api/role.ts` (new file - API client)
  - `frontend/src/pages/users/UserList.tsx` (add CRUD functionality)
  - `frontend/src/pages/roles/RoleList.tsx` (add CRUD functionality)
- User experience: Administrators can now fully manage users and roles through the UI
- Consistency: User and role management pages will match the functionality pattern of other management pages