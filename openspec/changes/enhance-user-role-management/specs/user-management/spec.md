## MODIFIED Requirements

### Requirement: User Account Management
The system SHALL provide user account management functionality including creation, retrieval, update, and deletion of user accounts.

#### Scenario: Create user via UI
- **WHEN** an administrator clicks the "新增用户" (Create User) button on the user management page
- **THEN** a modal form is displayed with fields for username, email, password, real_name, and status
- **AND** the administrator can fill in the form and submit it
- **AND** upon successful submission, a new user account is created and stored in the database with encrypted password
- **AND** the user list is refreshed to show the newly created user
- **AND** a success message is displayed

#### Scenario: Update user via UI
- **WHEN** an administrator clicks the "编辑" (Edit) button for a user in the user list
- **THEN** a modal form is displayed with the user's current information pre-filled
- **AND** the administrator can modify the user information (email, real_name, status, etc.)
- **AND** upon successful submission, the user record is updated in the database
- **AND** the user list is refreshed to show the updated information
- **AND** a success message is displayed

#### Scenario: Delete user via UI
- **WHEN** an administrator clicks the "删除" (Delete) button for a user in the user list
- **THEN** a confirmation dialog is displayed asking for confirmation
- **AND** upon confirmation, the user account is removed from the database
- **AND** the user list is refreshed to reflect the deletion
- **AND** a success message is displayed
- **AND** if the deletion fails, an error message is displayed

#### Scenario: Form validation for user management
- **WHEN** an administrator attempts to create or update a user with invalid data
- **THEN** form validation errors are displayed for the invalid fields
- **AND** the form submission is prevented until all required fields are valid
- **AND** appropriate error messages are shown (e.g., "用户名已存在", "邮箱格式不正确")

#### Scenario: Error handling for user operations
- **WHEN** a user management operation (create, update, delete) fails due to a server error
- **THEN** an appropriate error message is displayed to the administrator
- **AND** the user list remains in its previous state
- **AND** the modal form remains open (for create/update) to allow correction and retry