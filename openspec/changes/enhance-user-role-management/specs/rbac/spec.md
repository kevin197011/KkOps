## MODIFIED Requirements

### Requirement: Role-Based Access Control
The system SHALL implement RBAC (Role-Based Access Control) where users are assigned roles, and roles are assigned permissions.

#### Scenario: Create role via UI
- **WHEN** an administrator clicks the "新增角色" (Create Role) button on the role management page
- **THEN** a modal form is displayed with fields for name, code, and description
- **AND** the administrator can fill in the form and submit it
- **AND** upon successful submission, a new role is created in the database
- **AND** the role list is refreshed to show the newly created role
- **AND** a success message is displayed

#### Scenario: Update role via UI
- **WHEN** an administrator clicks the "编辑" (Edit) button for a role in the role list
- **THEN** a modal form is displayed with the role's current information pre-filled
- **AND** the administrator can modify the role information (name, description)
- **AND** the role code field is disabled (as it should not be changed after creation)
- **AND** upon successful submission, the role record is updated in the database
- **AND** the role list is refreshed to show the updated information
- **AND** a success message is displayed

#### Scenario: Delete role via UI
- **WHEN** an administrator clicks the "删除" (Delete) button for a role in the role list
- **THEN** a confirmation dialog is displayed asking for confirmation
- **AND** upon confirmation, the role is removed from the database
- **AND** the role list is refreshed to reflect the deletion
- **AND** a success message is displayed
- **AND** if the deletion fails, an error message is displayed

#### Scenario: Form validation for role management
- **WHEN** an administrator attempts to create or update a role with invalid data
- **THEN** form validation errors are displayed for the invalid fields
- **AND** the form submission is prevented until all required fields are valid
- **AND** appropriate error messages are shown (e.g., "角色代码已存在", "角色名称不能为空")

#### Scenario: Error handling for role operations
- **WHEN** a role management operation (create, update, delete) fails due to a server error
- **THEN** an appropriate error message is displayed to the administrator
- **AND** the role list remains in its previous state
- **AND** the modal form remains open (for create/update) to allow correction and retry