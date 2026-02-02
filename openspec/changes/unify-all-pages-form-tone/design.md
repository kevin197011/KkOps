# Design: 所有页面表单色调统一

## Context

基础设施下 5 个子页面已统一为：页面容器 + Card + 图标+标题 + theme token、Modal body token。其余主内容页（用户、角色、审计、部署、任务模板、执行、定时任务、SSH 密钥、个人中心、仪表盘、运维导航、分类等）仍存在裸 div、h2、硬编码 #xxx、无 Card 或 Modal 背景不一致的问题，需将同一套约定扩展到全站。

## Goals / Non-Goals

### Goals
- MainLayout 下所有主内容页视觉上统一：同一主题下背景、卡片、标题、表格、弹窗表单的色调与风格一致。
- 使用 theme token 驱动颜色，避免各页硬编码不同 hex；深色/浅色由全局主题决定。
- 列表/CRUD 类页面采用统一结构：外层 padding + 背景 token → Card（token）→ 标题区（图标 + Typography.Title）→ 表格/筛选/操作。

### Non-Goals
- 不改变各页面业务逻辑、表格列、筛选项功能。
- 不强制改动登录页（Login）的独立视觉风格。
- 不强制改动 WebSSH 终端内 xterm.js 的配色（终端可读性优先）；仅对终端页的工具栏/侧栏等非终端区域使用 token。
- 不强制改动全屏日志查看器（如 TaskExecutionLogs）内部代码块配色，可在适用处使用 token。

## Decisions

### 1. 参考基准

- **以基础设施子页面为基准**：已完成的 ProjectList、EnvironmentList、CloudPlatformList、AssetList、TagList、AssetDetail 所用结构（页面容器 token、Card token、图标+Title、Modal body token）作为全站统一标准。

### 2. 统一结构（列表/CRUD 页）

- **页面容器**：根节点 `padding: 24`、`background: token.colorBgContainer`（或继承 Layout Content），`minHeight: '100%'` 按需。
- **主内容区**：列表+筛选+操作置于一张 **Card** 内；Card 使用 `style={{ background: token.colorBgElevated, borderColor: token.colorBorderSecondary }}` 及 `styles.body` 的 padding。
- **标题区**：左侧图标 + Typography.Title level 3，右侧主操作按钮；图标与标题文案按各页语义选择。
- **Modal**：`styles.body` 使用 `background: token.colorBgElevated`，与主内容区一致。

### 3. 特殊页面

| 页面 | 处理方式 |
|------|----------|
| Dashboard | 卡片区、统计区使用 token；图表/装饰色可保留语义色（如成功绿、警告黄），但背景/边框用 token。 |
| WebSSHTerminal | 工具栏、侧栏、连接列表等非终端区域使用 token；xterm 终端内部配色不强制改为 token。 |
| ExecutionOperatorPage | 表单区、结果区使用 token；内部组件（如 HostSelector、ExecutionResults）中与主题相关的背景/边框改为 token。 |
| TaskExecutionLogs / 全屏日志 | 外层容器用 token；日志代码块背景/文字色可保留可读性优先，或使用 token。 |
| UserProfile | 与列表页一致：容器 + Card + token；头部装饰可保留渐变但避免与主题冲突。 |
| Login | 排除；保持独立登录页风格。 |

### 4. 实施顺序

- 按模块分批：用户与角色 → 审计 → 部署 → 任务模板与执行相关 → 定时任务 → SSH 密钥 → 个人中心 → 仪表盘 → 运维导航 → 分类；弹窗/子组件随主页面一并调整。
- 每批内：先页面容器与主 Card，再标题区，再 Modal，最后替换硬编码 hex 为 token。

### 5. 可选：公共组件

- 若多页重复明显，可抽「列表页布局」组件（接收 title、icon、children、extra），或集中定义 token 常量；否则各页内联相同结构，优先保证统一再考虑抽取。

## Risks / Trade-offs

| 风险 | 缓解 |
|------|------|
| 与 MainLayout Content 背景重复 | 使用与 Content 相同的 token，或仅 Card 内使用 colorBgElevated。 |
| 部分页面结构差异大（如 Dashboard） | 仅对适用区域（卡片、统计块）应用 token，不强行改布局。 |
| 终端/日志可读性 | 终端与日志查看器内部配色不强制 token，仅外围 UI 统一。 |
