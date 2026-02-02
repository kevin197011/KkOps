# Design: WebSSH 连线录像（审计连线）

## Context

当前 WebSSH 在 `ssh.go` 中建立 SSH 会话后，将 stdout/stderr 转发到 WebSocket，将客户端输入写入 stdin，会话结束后无持久化。需要在不影响现有实时终端体验的前提下，在服务端对会话进行录制，并在「系统管理」下提供「审计连线」列表与回放。

## Goals / Non-Goals

### Goals
- 每条 WebSSH 会话产生一条「连线记录」，包含元数据（用户、资产、开始/结束时间、时长）与终端输出录像。
- 管理员（或具备审计连线权限的用户）可在「审计连线」页面按条件筛选记录并查看/回放某次会话的终端内容。
- 录制与存储由服务端完成，不依赖前端是否开启录像。

### Non-Goals
- 不要求实现实时监控（多人同时观看同一会话）；仅支持事后查看/回放。
- 不要求支持 ZMODEM/文件传输内容的语义化回放；可将传输期间的二进制/协议数据从录像中排除或按二进制块记录，回放时仅展示可展示部分。
- 录像保留策略（如保留天数、最大条数、存储上限）可在后续迭代中再定；本变更先实现「有记录、可回放」。

## Decisions

### 1. 数据模型

- **SSHConnectionRecord（或 WebSSHSessionRecord）**：  
  - 字段：ID、UserID、Username（冗余）、AssetID、AssetHostname（冗余）、StartedAt、EndedAt、DurationSeconds、Transcript（见下）、CreatedAt。  
- **Transcript 存储**：  
  - 方案 A：单表，Transcript 存 LONGTEXT（或 BLOB），内容为「输出块序列」的 JSON 或二进制格式（如 `[{ "t": 0, "d": "base64..." }, ...]`），便于带时间戳回放。  
  - 方案 B：主表仅存元数据，Transcript 存文件或对象存储，主表存 TranscriptPath。  
  - 建议首版采用方案 A，单表 + LONGTEXT，对单条记录设合理上限（如 1MB），超限可截断并标记 truncated；若后续单条录像过大再迁至方案 B。

### 2. 录制时机与内容

- **录制起点**：SSH 会话建立成功（PTY 已开、Shell 已启）后开始。  
- **录制内容**：  
  - 必须：发往 WebSocket 的终端输出（即当前发给前端的 stdout/stderr 内容），在写入 WebSocket 的同时追加到录制 buffer。  
  - 可选：用户输入（从 WebSocket 收到并写入 stdin 的内容）一并记录，便于回放时区分「输入/输出」。  
- **结束与持久化**：会话结束（stdout 读 EOF、或 WebSocket 关闭、或超时）时，将 buffer 序列化为 Transcript，与元数据一起写入 DB。

### 3. 在 ssh.go 中的集成方式

- 在现有「stdout → WebSocket」的 goroutine 中，对每次要发送的 data 同时 append 到 recordBuffer（线程安全）；若记录「用户输入」，在从 WebSocket 读入并写入 stdin 的分支中同样 append 到 recordBuffer。  
- 会话结束时（defer 或 done channel）生成 SSHConnectionRecord（UserID、AssetID、StartedAt、EndedAt、DurationSeconds、Transcript），调用 service 写入 DB。  
- 注意：ZMODEM 激活期间发往客户端的 BinaryMessage 可选择性不录入或录入为占位，避免录像体积过大且难以回放。

### 4. API

- **GET /api/v1/connection-audit**（或 `/api/v1/ssh-connection-records`）：分页列表，查询参数 user_id、asset_id、start_time、end_time；返回记录列表（不含 Transcript 大字段，或仅返回摘要）。  
- **GET /api/v1/connection-audit/:id**：单条详情 + Transcript，用于回放；权限校验 connection-audit:read。

### 5. 前端

- **路由**：`/connection-audit`（或 `/ssh-connection-audit`），对应「审计连线」页。  
- **菜单位置**：系统管理 → 审计连线（与「用户管理」「角色权限」「审计日志」并列）。  
- **列表页**：表格列（操作用户、目标资产、开始时间、结束时间、时长、操作「查看/回放」）；筛选（用户、资产、时间范围）。  
- **回放**：点击「查看/回放」进入详情/弹窗，使用只读终端组件或按时间顺序重放 Transcript 中的输出块（若格式带时间戳可做简单时间轴回放）。

### 6. 权限与菜单

- 新增权限：`connection-audit` / `read`，名称「审计连线查看」，描述「查看 WebSSH 连线录像记录与回放」。  
- MainLayout 系统管理 children 中新增一项：key `/connection-audit`，label「审计连线」，icon 可选 LinkOutlined 或 AuditOutlined。  
- 后端路由与 API 校验该权限。

## Risks / Trade-offs

| 风险 | 缓解 |
|------|------|
| 录像数据含敏感信息（命令、输出） | 仅授权用户可访问；后续可加保留期限与脱敏策略。 |
| 长会话导致单条 Transcript 过大 | 单条上限 + 截断标记；必要时改为文件/对象存储。 |
| 录制增加内存与 IO | 使用 buffer 追加，结束时一次性写入；可配置关闭录制或采样。 |

## Migration Plan

1. 新增 model、migration（若使用 migration 工具）、repository/service、HTTP handler 与路由。  
2. 在 ssh.go 中挂载录制 buffer 与结束时的持久化逻辑。  
3. 前端新增「审计连线」页面与菜单、权限配置。  
4. 联调列表与回放，确认权限与数据隔离。
