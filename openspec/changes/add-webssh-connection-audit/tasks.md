# Tasks: WebSSH 连线录像（审计连线）

## 1. 后端数据与权限

- [x] 1.1 新增 model SSHConnectionRecord（或 WebSSHSessionRecord）：ID、UserID、Username、AssetID、AssetHostname、StartedAt、EndedAt、DurationSeconds、Transcript（LONGTEXT）、CreatedAt；表名 ssh_connection_records
- [x] 1.2 新增权限 connection-audit:read（审计连线查看）；在 AllMenuPermissions 与 RoutePermissionMap 中注册
- [x] 1.3 新增 repository/service：创建记录、分页列表（按 user_id/asset_id/时间筛选）、按 ID 获取单条（含 Transcript）

## 2. WebSSH 录制集成

- [x] 2.1 在 ssh.go 中会话建立成功后创建录制 buffer（线程安全）；在 stdout → WebSocket 的发送路径上同时写入 buffer
- [x] 2.2 可选：在 WebSocket → stdin 的路径上同时写入 buffer（区分输入/输出）（未实现，设计为可选）
- [x] 2.3 会话结束时（stdout EOF 或 WebSocket 关闭）生成 SSHConnectionRecord，序列化 Transcript（如 JSON 数组 [{ t, d }]），调用 service 持久化；Transcript 超长时截断并标记

## 3. API 与路由

- [x] 3.1 新增 GET /api/v1/connection-audit：分页列表，查询参数 user_id、asset_id、start_time、end_time；权限 connection-audit:read
- [x] 3.2 新增 GET /api/v1/connection-audit/:id：单条详情含 Transcript；权限 connection-audit:read
- [x] 3.3 在 main.go 中注册路由与权限中间件

## 4. 前端页面与菜单

- [x] 4.1 新增「审计连线」页面组件（如 ConnectionAuditList），路由 /connection-audit；列表表格（操作用户、目标资产、开始时间、结束时间、时长、操作）；筛选（用户、资产、时间范围）；调用列表 API
- [x] 4.2 实现「查看/回放」：弹窗或详情页，展示该条记录的 Transcript（只读终端或按块顺序重放）
- [x] 4.3 在 MainLayout 系统管理 children 中新增菜单项：key /connection-audit，label「审计连线」，icon 可选 LinkOutlined 或 AuditOutlined
- [x] 4.4 在 App.tsx 中注册路由 /connection-audit，权限 requiredPermission 为 connection-audit:read

## 5. 验收

- [ ] 5.1 完成一次 WebSSH 连接并执行若干命令后断开，在「审计连线」列表中能看到该条记录，且可查看/回放终端内容
- [ ] 5.2 无 connection-audit:read 权限的用户无法访问「审计连线」页面与 API
