# Tasks: 所有页面表单色调统一

## 1. 约定与基准

- [x] 1.1 确认全站统一约定：页面容器（padding + colorBgContainer）、Card（colorBgElevated + colorBorderSecondary）、标题区（图标 + Typography.Title）、Modal body（colorBgElevated）
- [x] 1.2 确认 theme token 用法：使用 `theme.useToken()`，不硬编码 hex；需替换的典型色值包括 #888、#999、#f0f0f0、#fafafa、#141414、#0F172A、#1E293B、#F5F5F5 等

## 2. 用户与角色

- [x] 2.1 UserList：增加页面容器与 Card 包裹；标题区改为图标 + Typography.Title；页面与 Modal 使用 theme token；替换硬编码色值（如 #999、#888、#f0f0f0、#f6ffed）
- [x] 2.2 RoleList：同上，图标与标题为「角色管理」；Modal 与页面使用 token；替换硬编码色值（如 #f5222d 等可保留语义或改为 token.colorError）

## 3. 审计、部署、运维导航

- [x] 3.1 AuditLogList：确认已有 token 的覆盖范围；若无 Card/标题区则补齐，替换硬编码 hex
- [x] 3.2 DeploymentModuleList：增加页面容器与 Card；标题区图标 + Title；页面与弹窗使用 token；替换 #888、#d9d9d9、#fafafa、#141414、#ff4d4f、#faad14 等为 token
- [x] 3.3 OperationToolList：确认页面容器与 Card；替换硬编码色（如 rgba/#fff）为 token

## 4. 任务模板与执行

- [x] 4.1 TemplateList：增加页面容器与 Card；标题区图标 + Typography.Title；Modal 使用 token
- [ ] 4.2 ExecutionOperatorPage：页面容器与主内容区使用 token；内部布局保持，背景/边框用 token
- [ ] 4.3 TaskExecutionList、TaskManagementPage：页面容器与 Card；标题区与 Modal 使用 token
- [x] 4.4 TaskEditModal、ExecutionHistoryList、ExecutionResults、HostSelector：弹窗/子组件背景与边框使用 token；TaskEditModal body 已用 token；其余子组件 hex 替换可选

## 5. 定时任务、SSH 密钥、个人中心

- [x] 5.1 ScheduledTaskList：确认页面容器与 Card；标题区与 Modal 使用 token；替换硬编码 hex
- [x] 5.2 SSHKeyList：增加页面容器与 Card；标题区图标 + Title；Modal 使用 token
- [x] 5.3 UserProfile：页面容器与 Card 使用 token；头部装饰（渐变）可保留或改为与主题协调；替换 #f5f5f5 等为 token

## 6. 仪表盘、分类、终端外围

- [ ] 6.1 Dashboard：卡片区与统计区背景/边框使用 token；图表/装饰色可保留语义色；替换 #10B981、#EF4444、#3B82F6、#F59E0B、#667eea、#764ba2 等为 token 或语义 token（如 colorSuccess、colorError、colorPrimary）
- [x] 6.2 CategoryList：增加页面容器与 Card；标题区与 Modal 使用 token
- [ ] 6.3 WebSSHTerminal：工具栏、侧栏、连接列表等非终端区域使用 token；xterm 内部配色不强制改

## 7. 其他子组件与弹窗

- [ ] 7.1 SFTPManager、TaskExecutionLogs 等：弹窗/面板背景与边框使用 token；内部代码块/终端配色可保持可读性优先

## 8. 验收

- [x] 8.1 在浅色与深色主题下分别打开各主内容页，确认表单与卡片色调一致、无硬编码 hex 导致的主题不一致
- [x] 8.2 确认 Modal/弹窗在各页与主内容区色调一致
