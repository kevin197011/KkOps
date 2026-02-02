# Global Pages UI – Form and Card Tone Consistency

## ADDED Requirements

### Requirement: All Main Content Pages Consistent Tone
All main content pages under MainLayout SHALL use a consistent form and card tone so that background, card, title, table, and modal form styling are aligned across the application and driven by the same theme tokens (extending the convention already applied to infrastructure sub-pages).

#### Scenario: Unified page structure for list/CRUD pages
- **WHEN** a user opens any list or CRUD page (e.g. 用户管理, 角色管理, 审计日志, 部署管理, 任务模板, 任务执行, 定时任务, SSH 密钥, 个人中心, 仪表盘, 运维导航, 分类管理)
- **THEN** the page SHALL use a common structure where applicable:
  - A page container with padding and background from theme token (e.g. colorBgContainer or inherit from layout)
  - Main content wrapped in a Card whose background, border, and shadow use theme tokens (e.g. colorBgElevated, colorBorderSecondary)
  - A title area with a leading icon and Typography.Title level 3 for the page name
  - Table, filters, and primary actions inside or adjacent to the Card; modals with body background from theme token (e.g. colorBgElevated)
- **AND** the same theme (light/dark) SHALL produce the same tone across all such pages

#### Scenario: Theme token usage across all pages
- **WHEN** rendering backgrounds, borders, or shadows on any main content page (except where explicitly excluded, e.g. terminal or log viewer inner content)
- **THEN** the implementation SHALL use Ant Design theme tokens (e.g. from theme.useToken()) such as colorBgContainer, colorBgElevated, colorBorderSecondary, colorTextHeading
- **AND** SHALL NOT hardcode different hex values per page (e.g. #0F172A, #F5F5F5, #888, #fafafa) so that a single theme change applies consistently

#### Scenario: Modal and form tone
- **WHEN** a user opens a create/edit modal or form panel on any main content page
- **THEN** the modal or panel body background SHALL use the same theme token as the main content (e.g. colorBgElevated) so that in dark mode the modal and the page card have a consistent tone
- **AND** SHALL apply to all list/CRUD pages and their sub-components (e.g. TaskEditModal, role assignment modal)

#### Scenario: Special pages and exclusions
- **WHEN** a page has special layout (e.g. Dashboard with widgets, WebSSH terminal with xterm.js, full-screen log viewer)
- **THEN** the page SHALL use theme tokens for the page container, toolbar, sidebar, and any Card/panel wrappers
- **AND** terminal or log viewer **inner** content (xterm palette, code block background) MAY keep fixed colors for readability and SHALL NOT be forced to use theme tokens
- **AND** the login page (Login) is excluded from this requirement and MAY keep its own visual style
