## MODIFIED Requirements

### Requirement: Operations Tool Management
The system SHALL provide a dedicated management page for administrators to manage operations tools through a user-friendly interface, in addition to the API endpoints.

#### Scenario: Access management page
- **WHEN** a user with `operation-tools:read` permission accesses the application
- **THEN** an "运维导航" (Operations Navigation) menu item is displayed in the sidebar menu below the "仪表板" (Dashboard) menu item
- **AND** clicking the menu item navigates to the operations tools management page

#### Scenario: View tools list
- **WHEN** a user accesses the operations tools management page
- **THEN** a table displaying all operations tools is shown
- **AND** the table includes columns for ID, name, category, icon, URL, order, status, and creation time
- **AND** the table supports pagination and can display both enabled and disabled tools

#### Scenario: Administrator creates a tool via UI
- **WHEN** an administrator with `operation-tools:*` permission clicks the "新增工具" (Add Tool) button
- **THEN** a modal form is displayed with fields for name, description, category, icon, URL, order, and enabled status
- **AND** name and URL fields are marked as required
- **WHEN** the administrator fills in the required fields and submits the form
- **THEN** the tool is created via API
- **AND** the tool list is refreshed
- **AND** a success message is displayed
- **AND** the tool becomes visible in both the management page and the dashboard

#### Scenario: Administrator edits a tool via UI
- **WHEN** an administrator with `operation-tools:*` permission clicks the "编辑" (Edit) button for a tool
- **THEN** a modal form is displayed with the tool's current values pre-filled
- **WHEN** the administrator modifies the fields and submits the form
- **THEN** the tool is updated via API
- **AND** the tool list is refreshed
- **AND** a success message is displayed
- **AND** changes are immediately reflected in both the management page and the dashboard

#### Scenario: Administrator deletes a tool via UI
- **WHEN** an administrator with `operation-tools:*` permission clicks the "删除" (Delete) button for a tool
- **THEN** a confirmation dialog is displayed
- **WHEN** the administrator confirms the deletion
- **THEN** the tool is deleted via API
- **AND** the tool is removed from the list
- **AND** a success message is displayed
- **AND** the tool no longer appears in the dashboard

#### Scenario: Permission-based UI visibility
- **WHEN** a user without `operation-tools:*` permission accesses the management page
- **THEN** the "新增工具" button is hidden
- **AND** the "编辑" and "删除" action buttons in the table are hidden
- **AND** the user can only view the tools list

#### Scenario: URL validation in UI
- **WHEN** creating or editing a tool
- **AND** the URL field does not start with "http://" or "https://"
- **THEN** a validation error message is displayed
- **AND** the form submission is prevented

## ADDED Requirements

### Requirement: Operations Navigation Menu Item
The system SHALL display an "运维导航" (Operations Navigation) menu item in the sidebar menu.

#### Scenario: Menu item display
- **WHEN** a user with `operation-tools:read` permission views the sidebar menu
- **THEN** an "运维导航" menu item is displayed below the "仪表板" menu item
- **AND** the menu item is separated from the dashboard by a divider

#### Scenario: Menu item navigation
- **WHEN** a user clicks the "运维导航" menu item
- **THEN** the user is navigated to the operations tools management page at `/operation-tools`

#### Scenario: Menu item icon
- **WHEN** the sidebar menu is displayed
- **THEN** the "运维导航" menu item shows an appropriate icon (e.g., `AppstoreOutlined`)

### Requirement: Operations Tools Management Page
The system SHALL provide a dedicated page for managing operations tools.

#### Scenario: Page layout
- **WHEN** a user accesses the operations tools management page
- **THEN** the page displays a title "运维导航"
- **AND** a table showing all operations tools
- **AND** an "新增工具" button (visible only to administrators)

#### Scenario: Table columns
- **WHEN** viewing the operations tools table
- **THEN** the table displays the following columns: ID, 名称 (Name), 分类 (Category), 图标 (Icon), URL, 排序 (Order), 状态 (Status), 创建时间 (Created At), 操作 (Actions)
- **AND** the URL column displays clickable links that open in a new tab
- **AND** the status column displays "启用" (Enabled) or "禁用" (Disabled) as tags

#### Scenario: Tool name click behavior
- **WHEN** a user clicks on a tool name in the table
- **THEN** the tool URL opens in a new browser tab
- **AND** the current page remains open
