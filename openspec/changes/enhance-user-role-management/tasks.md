# Implementation Tasks

## API Client Implementation

- [x] 1.1 Create user API client (`frontend/src/api/user.ts`)
  - Define TypeScript interfaces (User, CreateUserRequest, UpdateUserRequest)
  - Implement userApi with list, get, create, update, delete methods
  - Match the API structure used in other API clients (project, category, tag)

- [x] 1.2 Create role API client (`frontend/src/api/role.ts`)
  - Define TypeScript interfaces (Role, CreateRoleRequest, UpdateRoleRequest)
  - Implement roleApi with list, get, create, update, delete methods
  - Match the API structure used in other API clients

## User Management Page Enhancement

- [x] 2.1 Add Modal state and form handling to UserList
  - Add state for modal visibility and editing user
  - Add Ant Design Form instance
  - Import necessary components (Modal, Form, Input, Select, etc.)

- [x] 2.2 Implement create user functionality
  - Add handleCreate function to open modal for new user
  - Add user creation form with fields (username, email, password, real_name, status)
  - Implement form validation rules
  - Call userApi.create on form submission

- [x] 2.3 Implement edit user functionality
  - Add handleEdit function to open modal with existing user data
  - Pre-fill form with user data
  - Call userApi.update on form submission
  - Handle password field (optional for updates)

- [x] 2.4 Implement delete user functionality
  - Add handleDelete function with confirmation dialog
  - Use Modal.confirm for delete confirmation
  - Call userApi.delete API
  - Refresh user list after deletion

- [x] 2.5 Add Modal form component to UserList
  - Create Modal with Form inside
  - Add form fields: username, email, password (for create), real_name, status
  - Add form validation
  - Handle form submission (create/update)

- [x] 2.6 Update UserList buttons to use handlers
  - Add onClick handler to "新增用户" button (handleCreate)
  - Add onClick handlers to "编辑" and "删除" buttons in table
  - Ensure buttons are properly wired up

## Role Management Page Enhancement

- [x] 3.1 Add Modal state and form handling to RoleList
  - Add state for modal visibility and editing role
  - Add Ant Design Form instance
  - Import necessary components (Modal, Form, Input, etc.)

- [x] 3.2 Implement create role functionality
  - Add handleCreate function to open modal for new role
  - Add role creation form with fields (name, code, description)
  - Implement form validation rules
  - Call roleApi.create on form submission

- [x] 3.3 Implement edit role functionality
  - Add handleEdit function to open modal with existing role data
  - Pre-fill form with role data
  - Call roleApi.update on form submission

- [x] 3.4 Implement delete role functionality
  - Add handleDelete function with confirmation dialog
  - Use Modal.confirm for delete confirmation
  - Call roleApi.delete API
  - Refresh role list after deletion

- [x] 3.5 Add Modal form component to RoleList
  - Create Modal with Form inside
  - Add form fields: name, code, description
  - Add form validation
  - Handle form submission (create/update)

- [x] 3.6 Update RoleList buttons to use handlers
  - Add onClick handler to "新增角色" button (handleCreate)
  - Add onClick handlers to "编辑" and "删除" buttons in table
  - Ensure buttons are properly wired up

## Testing & Validation

- [ ] 4.1 Test user management functionality
  - Verify create user form works correctly
  - Verify edit user form pre-fills and updates correctly
  - Verify delete user with confirmation works
  - Verify form validation works
  - Verify error handling displays appropriate messages

- [ ] 4.2 Test role management functionality
  - Verify create role form works correctly
  - Verify edit role form pre-fills and updates correctly
  - Verify delete role with confirmation works
  - Verify form validation works
  - Verify error handling displays appropriate messages

- [ ] 4.3 Verify API integration
  - Verify API calls match backend endpoints
  - Verify request/response data structures match backend
  - Verify error responses are handled correctly

## Code Quality

- [x] 5.1 Ensure code follows project patterns
  - Match code style with other list pages (ProjectList, CategoryList, TagList)
  - Use consistent naming conventions
  - Follow TypeScript best practices
  - Note: Code structure matches ProjectList and CategoryList patterns

- [x] 5.2 Ensure accessibility
  - Add ARIA labels to buttons and form fields
  - Ensure keyboard navigation works
  - Verify screen reader compatibility
  - Note: ARIA labels are added to all buttons and form fields, keyboard navigation is supported by Ant Design components