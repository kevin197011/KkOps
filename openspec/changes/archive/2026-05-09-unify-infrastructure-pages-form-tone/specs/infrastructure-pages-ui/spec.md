# Infrastructure Pages UI – Form and Card Tone Consistency

## ADDED Requirements

### Requirement: Infrastructure Sub-Pages Consistent Tone
All sub-pages under the 「基础设施」(Infrastructure) menu SHALL use a consistent form and card tone so that background, card, title, table, and modal form styling are aligned across pages and driven by the same theme tokens.

#### Scenario: Unified page structure
- **WHEN** a user opens any infrastructure sub-page (项目管理, 环境管理, 云平台管理, 资产管理, 标签管理)
- **THEN** the page SHALL use a common structure:
  - A page container with padding and background from theme token (e.g. colorBgContainer or inherit from layout)
  - Main content wrapped in a Card whose background, border, and shadow use theme tokens (e.g. colorBgElevated, colorBorderSecondary)
  - A title area with a leading icon (specific to the page) and Typography.Title level 3 for the page name
  - Table and filters inside the Card; primary actions (e.g. 新增) in the title area
- **AND** the same theme (light/dark) SHALL produce the same tone across all five sub-pages

#### Scenario: Theme token usage
- **WHEN** rendering backgrounds, borders, or shadows on infrastructure sub-pages
- **THEN** the implementation SHALL use Ant Design theme tokens (e.g. from theme.useToken()) such as colorBgContainer, colorBgElevated, colorBorder
- **AND** SHALL NOT hardcode different hex values per page (e.g. #0F172A, #F5F5F5) so that a single theme change applies consistently

#### Scenario: Modal form tone
- **WHEN** a user opens a create/edit modal on an infrastructure sub-page
- **THEN** the modal body background SHALL use the same theme token as the main content (e.g. colorBgElevated) so that in dark mode the modal and the page card have a consistent tone
- **AND** SHALL apply to all five sub-pages (project, environment, cloud-platform, asset, tag)

#### Scenario: Detail page consistency
- **WHEN** a user opens an infrastructure-related detail page (e.g. 资产详情)
- **THEN** the detail page or modal SHALL use the same theme tokens for background and borders as the list pages so that tone remains consistent within the infrastructure section
