# Change: 针对所有页面的表单色调统一

## Why

基础设施子页面已完成表单与卡片色调统一（`unify-infrastructure-pages-form-tone`），但其他主内容页（用户、角色、审计、部署、任务模板、执行、定时任务、SSH 密钥、个人中心、仪表盘、运维导航、分类等）仍存在裸 div、硬编码 hex 色值、无统一 Card/标题区或 Modal 背景不一致的情况，导致整站视觉不统一。用户希望**所有页面**的表单与卡片色调保持一致。

## What Changes

- **统一约定**：将基础设施子页面已采用的约定扩展到 MainLayout 下所有主内容页——页面容器（padding + 背景 theme token）、主内容区用 Card 包裹（或与布局一致）、标题区为图标 + Typography.Title、表格/筛选项置于 Card 内、Modal/表单弹窗 body 使用 theme token 背景。
- **色调来源**：所有涉及背景、边框、阴影的样式优先使用 Ant Design `theme.useToken()`（如 colorBgContainer、colorBgElevated、colorBorderSecondary），避免各页硬编码不同 hex；深色/浅色由全局主题驱动。
- **适用范围**：除登录页（Login）及终端内 xterm 配色外的、由 MainLayout 包裹的所有主内容页；特殊页面（如仪表盘、WebSSH 终端、执行日志全屏）在适用处使用 token，终端/日志查看器内部配色可保持可读性优先。
- **可选**：若重复较多，可抽公共「列表页布局」组件或样式常量，供各页复用。

## Impact

- **受影响规格**：新增或扩展 UI 一致性规格（如 global-pages-ui），与现有 infrastructure-pages-ui 衔接。
- **受影响代码**：`frontend/src/pages/` 下除 Login 外的各页面及子组件——UserList、RoleList、AuditLogList、DeploymentModuleList、TemplateList、ExecutionOperatorPage、TaskExecutionList、TaskManagementPage、ScheduledTaskList、SSHKeyList、UserProfile、Dashboard、OperationToolList、CategoryList，以及相关弹窗/子组件（如 TaskEditModal、ExecutionHistoryList、SFTPManager 等）。WebSSHTerminal 仅对工具栏/侧栏等非终端区域统一 token；终端与日志查看器内部配色可保持现状或单独约定。
