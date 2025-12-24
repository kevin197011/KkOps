# Design: Operations Platform Architecture

## Context
构建一个基于 Salt 的运维中台管理系统，需要整合多个外部系统（Salt、ELK、Prometheus），提供统一的操作界面和管理能力。系统需要支持大规模主机管理（1000+），确保高性能、高可用性和安全性。

## Goals / Non-Goals

### Goals
- 提供统一的运维操作入口，整合各类运维工具
- 支持大规模主机管理（1000+ 主机）
- 实现完整的权限控制和审计追溯
- 提供友好的 Web 界面，提升运维效率
- 支持与 Salt、ELK、Prometheus 的深度集成

### Non-Goals
- 不替代 Salt、ELK、Prometheus 等底层系统
- 不实现完整的 CI/CD 流程（仅关注发布管理）
- 不实现完整的监控告警系统（仅对接 Prometheus）
- 不实现完整的日志分析系统（仅对接 ELK）

## Decisions

### Decision: 后端架构采用分层架构
**What**: 采用 Handler -> Service -> Repository 三层架构
**Why**: 
- 职责清晰，易于维护和测试
- 符合 Go 社区最佳实践
- 便于后续扩展和重构

**Alternatives considered**:
- 微服务架构：当前规模不需要，增加复杂度
- 领域驱动设计（DDD）：过度设计，当前需求相对简单

### Decision: 使用 PostgreSQL 作为主数据库
**What**: 使用 PostgreSQL 存储应用数据
**Why**:
- 支持复杂查询和事务
- 支持 JSON 字段，便于存储灵活配置
- 性能优秀，社区活跃

**Alternatives considered**:
- MySQL：功能类似，但 JSON 支持不如 PostgreSQL
- MongoDB：NoSQL 不适合关系型数据场景

### Decision: 使用 Redis 作为缓存和会话存储
**What**: 使用 Redis 存储缓存数据和用户会话
**Why**:
- 高性能，支持多种数据结构
- 支持过期时间，适合会话管理
- 支持分布式部署

### Decision: 使用 JWT 进行身份认证
**What**: 使用 JWT Token 进行用户认证
**Why**:
- 无状态，易于扩展
- 支持分布式部署
- 减少服务器端会话存储压力

**Alternatives considered**:
- Session 认证：需要服务器端存储，不利于扩展
- OAuth2：过度设计，当前不需要第三方登录

### Decision: 前端使用 Ant Design Pro 作为基础框架
**What**: 基于 Ant Design Pro 构建前端应用
**Why**:
- 提供完整的后台管理系统模板
- 组件丰富，开箱即用
- 文档完善，社区活跃

**Alternatives considered**:
- 从零开始：开发效率低，重复造轮子
- 其他 UI 框架：Ant Design 更适合后台管理系统

### Decision: Salt 集成使用 Salt API
**What**: 通过 Salt API (REST API) 与 Salt Master 交互
**Why**:
- 官方推荐方式，稳定可靠
- 支持异步任务和结果查询
- 易于集成和扩展

**Alternatives considered**:
- Salt CLI：需要进程调用，性能较差
- Salt Event System：复杂度高，当前不需要

### Decision: ELK 集成直接查询 Elasticsearch
**What**: 直接通过 Elasticsearch REST API 查询日志
**Why**:
- 性能好，支持复杂查询
- 避免 Logstash 和 Kibana 的额外依赖
- 灵活度高，可以自定义查询逻辑

**Alternatives considered**:
- 通过 Kibana API：功能受限，不如直接查询 ES
- 通过 Logstash：增加复杂度，不需要

### Decision: Prometheus 集成使用 PromQL 查询
**What**: 通过 Prometheus Query API 使用 PromQL 查询指标
**Why**:
- 官方标准方式
- 支持复杂的指标查询和聚合
- 性能优秀

## Database Schema Design

### 1. 用户和权限管理表

#### users (用户表)
存储系统用户信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 用户ID |
| username | VARCHAR(50) | UNIQUE, NOT NULL | 用户名 |
| email | VARCHAR(100) | UNIQUE, NOT NULL | 邮箱 |
| password_hash | VARCHAR(255) | NOT NULL | 密码哈希值 |
| display_name | VARCHAR(100) | | 显示名称 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'active' | 状态: active, disabled, deleted |
| last_login_at | TIMESTAMP | | 最后登录时间 |
| last_login_ip | VARCHAR(45) | | 最后登录IP |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_users_username` ON `users(username)`
- `idx_users_email` ON `users(email)`
- `idx_users_status` ON `users(status)`
- `idx_users_deleted_at` ON `users(deleted_at)`

**约束**:
- `CHECK (status IN ('active', 'disabled', 'deleted'))`

#### roles (角色表)
存储系统角色信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 角色ID |
| name | VARCHAR(50) | UNIQUE, NOT NULL | 角色名称 |
| description | TEXT | | 角色描述 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_roles_name` ON `roles(name)`

#### permissions (权限表)
存储系统权限定义。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 权限ID |
| code | VARCHAR(100) | UNIQUE, NOT NULL | 权限代码 (如: host:read, deployment:create) |
| name | VARCHAR(100) | NOT NULL | 权限名称 |
| resource_type | VARCHAR(50) | NOT NULL | 资源类型 (host, deployment, task, etc.) |
| action | VARCHAR(50) | NOT NULL | 操作类型 (read, create, update, delete, execute) |
| description | TEXT | | 权限描述 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_permissions_code` ON `permissions(code)`
- `idx_permissions_resource_type` ON `permissions(resource_type)`

#### user_roles (用户角色关联表)
用户和角色的多对多关系。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 关联ID |
| user_id | BIGINT | NOT NULL, FK -> users(id) | 用户ID |
| role_id | BIGINT | NOT NULL, FK -> roles(id) | 角色ID |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_user_roles_user_id` ON `user_roles(user_id)`
- `idx_user_roles_role_id` ON `user_roles(role_id)`
- `UNIQUE (user_id, role_id)` - 防止重复关联

#### role_permissions (角色权限关联表)
角色和权限的多对多关系。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 关联ID |
| role_id | BIGINT | NOT NULL, FK -> roles(id) | 角色ID |
| permission_id | BIGINT | NOT NULL, FK -> permissions(id) | 权限ID |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_role_permissions_role_id` ON `role_permissions(role_id)`
- `idx_role_permissions_permission_id` ON `role_permissions(permission_id)`
- `UNIQUE (role_id, permission_id)` - 防止重复关联

### 2. 主机管理表

#### hosts (主机表)
存储主机信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 主机ID |
| hostname | VARCHAR(255) | NOT NULL | 主机名 |
| ip_address | VARCHAR(45) | NOT NULL | IP地址 (支持IPv4/IPv6) |
| salt_minion_id | VARCHAR(255) | UNIQUE | Salt Minion ID |
| os_type | VARCHAR(50) | | 操作系统类型 |
| os_version | VARCHAR(100) | | 操作系统版本 |
| cpu_cores | INTEGER | | CPU核心数 |
| memory_gb | DECIMAL(10,2) | | 内存大小(GB) |
| disk_gb | DECIMAL(10,2) | | 磁盘大小(GB) |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'unknown' | 状态: online, offline, unknown |
| last_seen_at | TIMESTAMP | | 最后在线时间 |
| salt_version | VARCHAR(50) | | Salt版本 |
| metadata | JSONB | | 扩展元数据 (JSON格式) |
| description | TEXT | | 描述 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_hosts_hostname` ON `hosts(hostname)`
- `idx_hosts_ip_address` ON `hosts(ip_address)`
- `idx_hosts_salt_minion_id` ON `hosts(salt_minion_id)`
- `idx_hosts_status` ON `hosts(status)`
- `idx_hosts_deleted_at` ON `hosts(deleted_at)`
- `idx_hosts_metadata` ON `hosts USING GIN(metadata)` - GIN索引用于JSONB查询

**约束**:
- `CHECK (status IN ('online', 'offline', 'unknown'))`
- `UNIQUE (hostname, deleted_at)` - 软删除后允许同名主机

#### host_groups (主机组表)
存储主机分组信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 组ID |
| name | VARCHAR(100) | NOT NULL | 组名称 |
| description | TEXT | | 组描述 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_host_groups_name` ON `host_groups(name)`

#### host_group_members (主机组成员关联表)
主机和组的多对多关系。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 关联ID |
| host_id | BIGINT | NOT NULL, FK -> hosts(id) | 主机ID |
| group_id | BIGINT | NOT NULL, FK -> host_groups(id) | 组ID |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_host_group_members_host_id` ON `host_group_members(host_id)`
- `idx_host_group_members_group_id` ON `host_group_members(group_id)`
- `UNIQUE (host_id, group_id)` - 防止重复关联

#### host_tags (主机标签表)
存储主机标签定义。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 标签ID |
| name | VARCHAR(50) | UNIQUE, NOT NULL | 标签名称 |
| color | VARCHAR(7) | | 标签颜色 (hex格式) |
| description | TEXT | | 标签描述 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_host_tags_name` ON `host_tags(name)`

#### host_tag_assignments (主机标签关联表)
主机和标签的多对多关系。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 关联ID |
| host_id | BIGINT | NOT NULL, FK -> hosts(id) | 主机ID |
| tag_id | BIGINT | NOT NULL, FK -> host_tags(id) | 标签ID |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_host_tag_assignments_host_id` ON `host_tag_assignments(host_id)`
- `idx_host_tag_assignments_tag_id` ON `host_tag_assignments(tag_id)`
- `UNIQUE (host_id, tag_id)` - 防止重复关联

### 3. SSH 管理表

#### ssh_connections (SSH连接配置表)
存储SSH连接配置信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 连接ID |
| name | VARCHAR(100) | NOT NULL | 连接名称 |
| host_id | BIGINT | FK -> hosts(id) | 关联主机ID (可选) |
| hostname | VARCHAR(255) | NOT NULL | 目标主机名或IP |
| port | INTEGER | NOT NULL, DEFAULT 22 | SSH端口 |
| username | VARCHAR(100) | NOT NULL | 用户名 |
| auth_type | VARCHAR(20) | NOT NULL | 认证类型: password, key |
| password_encrypted | TEXT | | 加密后的密码 (auth_type=password时) |
| key_id | BIGINT | FK -> ssh_keys(id) | 关联密钥ID (auth_type=key时) |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'active' | 状态: active, disabled |
| last_connected_at | TIMESTAMP | | 最后连接时间 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_ssh_connections_host_id` ON `ssh_connections(host_id)`
- `idx_ssh_connections_key_id` ON `ssh_connections(key_id)`
- `idx_ssh_connections_status` ON `ssh_connections(status)`

**约束**:
- `CHECK (auth_type IN ('password', 'key'))`
- `CHECK (status IN ('active', 'disabled'))`
- `CHECK (port > 0 AND port <= 65535)`

#### ssh_keys (SSH密钥表)
存储SSH密钥信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 密钥ID |
| name | VARCHAR(100) | NOT NULL | 密钥名称 |
| key_type | VARCHAR(20) | NOT NULL | 密钥类型: rsa, ed25519, ecdsa |
| private_key_encrypted | TEXT | NOT NULL | 加密后的私钥 |
| public_key | TEXT | NOT NULL | 公钥内容 |
| fingerprint | VARCHAR(100) | UNIQUE | 密钥指纹 |
| passphrase_encrypted | TEXT | | 加密后的密钥密码 (如果有) |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_ssh_keys_fingerprint` ON `ssh_keys(fingerprint)`
- `idx_ssh_keys_key_type` ON `ssh_keys(key_type)`

**约束**:
- `CHECK (key_type IN ('rsa', 'ed25519', 'ecdsa'))`

#### ssh_sessions (SSH会话表)
存储SSH会话记录。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 会话ID |
| connection_id | BIGINT | NOT NULL, FK -> ssh_connections(id) | 连接ID |
| user_id | BIGINT | NOT NULL, FK -> users(id) | 用户ID |
| session_token | VARCHAR(255) | UNIQUE, NOT NULL | 会话令牌 |
| client_ip | VARCHAR(45) | | 客户端IP |
| started_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 开始时间 |
| ended_at | TIMESTAMP | | 结束时间 |
| duration_seconds | INTEGER | | 会话时长(秒) |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'active' | 状态: active, closed, timeout |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_ssh_sessions_connection_id` ON `ssh_sessions(connection_id)`
- `idx_ssh_sessions_user_id` ON `ssh_sessions(user_id)`
- `idx_ssh_sessions_session_token` ON `ssh_sessions(session_token)`
- `idx_ssh_sessions_started_at` ON `ssh_sessions(started_at)`

**约束**:
- `CHECK (status IN ('active', 'closed', 'timeout'))`

### 4. 发布管理表

#### deployment_configs (部署配置表)
存储部署配置信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 配置ID |
| name | VARCHAR(100) | NOT NULL | 配置名称 |
| application_name | VARCHAR(100) | NOT NULL | 应用名称 |
| description | TEXT | | 配置描述 |
| salt_state_files | TEXT[] | NOT NULL | Salt State文件列表 |
| target_groups | BIGINT[] | | 目标主机组ID列表 |
| target_hosts | BIGINT[] | | 目标主机ID列表 |
| environment | VARCHAR(50) | | 环境: dev, test, prod |
| config_data | JSONB | | 配置数据 (JSON格式) |
| created_by | BIGINT | NOT NULL, FK -> users(id) | 创建人 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_deployment_configs_name` ON `deployment_configs(name)`
- `idx_deployment_configs_application_name` ON `deployment_configs(application_name)`
- `idx_deployment_configs_created_by` ON `deployment_configs(created_by)`
- `idx_deployment_configs_config_data` ON `deployment_configs USING GIN(config_data)`

**约束**:
- `CHECK (environment IN ('dev', 'test', 'prod') OR environment IS NULL)`

#### deployments (部署记录表)
存储部署执行记录。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 部署ID |
| config_id | BIGINT | NOT NULL, FK -> deployment_configs(id) | 配置ID |
| version | VARCHAR(50) | NOT NULL | 版本号 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'pending' | 状态: pending, running, completed, failed, cancelled |
| started_by | BIGINT | NOT NULL, FK -> users(id) | 启动人 |
| started_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 开始时间 |
| completed_at | TIMESTAMP | | 完成时间 |
| duration_seconds | INTEGER | | 执行时长(秒) |
| target_hosts | BIGINT[] | NOT NULL | 目标主机ID列表 |
| salt_job_id | VARCHAR(255) | | Salt Job ID |
| results | JSONB | | 执行结果 (按主机) |
| error_message | TEXT | | 错误信息 |
| is_rollback | BOOLEAN | NOT NULL, DEFAULT FALSE | 是否为回滚操作 |
| rollback_from_deployment_id | BIGINT | FK -> deployments(id) | 回滚来源部署ID |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |

**索引**:
- `idx_deployments_config_id` ON `deployments(config_id)`
- `idx_deployments_started_by` ON `deployments(started_by)`
- `idx_deployments_status` ON `deployments(status)`
- `idx_deployments_started_at` ON `deployments(started_at)`
- `idx_deployments_salt_job_id` ON `deployments(salt_job_id)`

**约束**:
- `CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'))`

#### deployment_versions (部署版本表)
存储应用版本信息。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 版本ID |
| application_name | VARCHAR(100) | NOT NULL | 应用名称 |
| version | VARCHAR(50) | NOT NULL | 版本号 |
| release_notes | TEXT | | 发布说明 |
| created_by | BIGINT | NOT NULL, FK -> users(id) | 创建人 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_deployment_versions_application_name` ON `deployment_versions(application_name)`
- `idx_deployment_versions_version` ON `deployment_versions(version)`
- `UNIQUE (application_name, version)` - 应用版本唯一

#### salt_states (Salt State文件表)
存储Salt State文件。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 文件ID |
| name | VARCHAR(255) | NOT NULL | 文件名称 |
| file_path | VARCHAR(500) | NOT NULL | 文件路径 |
| content | TEXT | NOT NULL | 文件内容 |
| file_size | BIGINT | NOT NULL | 文件大小(字节) |
| checksum | VARCHAR(64) | NOT NULL | 文件校验和 (SHA256) |
| usage_count | INTEGER | NOT NULL, DEFAULT 0 | 使用次数 |
| created_by | BIGINT | NOT NULL, FK -> users(id) | 创建人 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_salt_states_name` ON `salt_states(name)`
- `idx_salt_states_file_path` ON `salt_states(file_path)`
- `idx_salt_states_checksum` ON `salt_states(checksum)`

### 5. 定时任务表

#### scheduled_tasks (定时任务表)
存储定时任务配置。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 任务ID |
| name | VARCHAR(100) | NOT NULL | 任务名称 |
| description | TEXT | | 任务描述 |
| schedule_type | VARCHAR(20) | NOT NULL | 调度类型: cron, once, interval |
| cron_expression | VARCHAR(100) | | Cron表达式 (schedule_type=cron时) |
| execute_at | TIMESTAMP | | 执行时间 (schedule_type=once时) |
| interval_seconds | INTEGER | | 间隔秒数 (schedule_type=interval时) |
| salt_command | TEXT | | Salt命令 |
| salt_state | VARCHAR(255) | | Salt State文件 |
| target_type | VARCHAR(20) | NOT NULL | 目标类型: all, groups, hosts |
| target_groups | BIGINT[] | | 目标主机组ID列表 |
| target_hosts | BIGINT[] | | 目标主机ID列表 |
| enabled | BOOLEAN | NOT NULL, DEFAULT TRUE | 是否启用 |
| timeout_seconds | INTEGER | DEFAULT 300 | 超时时间(秒) |
| created_by | BIGINT | NOT NULL, FK -> users(id) | 创建人 |
| last_executed_at | TIMESTAMP | | 最后执行时间 |
| next_execute_at | TIMESTAMP | | 下次执行时间 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_scheduled_tasks_name` ON `scheduled_tasks(name)`
- `idx_scheduled_tasks_enabled` ON `scheduled_tasks(enabled)`
- `idx_scheduled_tasks_next_execute_at` ON `scheduled_tasks(next_execute_at)`
- `idx_scheduled_tasks_created_by` ON `scheduled_tasks(created_by)`

**约束**:
- `CHECK (schedule_type IN ('cron', 'once', 'interval'))`
- `CHECK (target_type IN ('all', 'groups', 'hosts'))`
- `CHECK (timeout_seconds > 0)`

#### task_executions (任务执行记录表)
存储任务执行历史。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 执行ID |
| task_id | BIGINT | NOT NULL, FK -> scheduled_tasks(id) | 任务ID |
| execution_type | VARCHAR(20) | NOT NULL | 执行类型: scheduled, manual |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'running' | 状态: running, success, failed, timeout |
| started_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 开始时间 |
| completed_at | TIMESTAMP | | 完成时间 |
| duration_seconds | INTEGER | | 执行时长(秒) |
| salt_job_id | VARCHAR(255) | | Salt Job ID |
| results | JSONB | | 执行结果 (按主机) |
| error_message | TEXT | | 错误信息 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_task_executions_task_id` ON `task_executions(task_id)`
- `idx_task_executions_status` ON `task_executions(status)`
- `idx_task_executions_started_at` ON `task_executions(started_at)`
- `idx_task_executions_salt_job_id` ON `task_executions(salt_job_id)`

**约束**:
- `CHECK (execution_type IN ('scheduled', 'manual'))`
- `CHECK (status IN ('running', 'success', 'failed', 'timeout'))`

### 6. 监控管理表

#### alert_rules (告警规则表)
存储Prometheus告警规则配置。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 规则ID |
| name | VARCHAR(100) | NOT NULL | 规则名称 |
| description | TEXT | | 规则描述 |
| promql_expression | TEXT | NOT NULL | PromQL表达式 |
| threshold | DECIMAL(10,2) | | 阈值 |
| severity | VARCHAR(20) | NOT NULL, DEFAULT 'warning' | 严重程度: critical, warning, info |
| enabled | BOOLEAN | NOT NULL, DEFAULT TRUE | 是否启用 |
| notification_channels | TEXT[] | | 通知渠道列表 |
| last_triggered_at | TIMESTAMP | | 最后触发时间 |
| created_by | BIGINT | NOT NULL, FK -> users(id) | 创建人 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_alert_rules_name` ON `alert_rules(name)`
- `idx_alert_rules_enabled` ON `alert_rules(enabled)`
- `idx_alert_rules_severity` ON `alert_rules(severity)`

**约束**:
- `CHECK (severity IN ('critical', 'warning', 'info'))`

#### metric_dashboards (监控仪表盘表)
存储监控仪表盘配置。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 仪表盘ID |
| name | VARCHAR(100) | NOT NULL | 仪表盘名称 |
| description | TEXT | | 仪表盘描述 |
| layout_config | JSONB | NOT NULL | 布局配置 (JSON格式) |
| panels | JSONB | NOT NULL | 面板配置列表 (JSON格式) |
| is_public | BOOLEAN | NOT NULL, DEFAULT FALSE | 是否公开 |
| created_by | BIGINT | NOT NULL, FK -> users(id) | 创建人 |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | | 软删除时间 |

**索引**:
- `idx_metric_dashboards_name` ON `metric_dashboards(name)`
- `idx_metric_dashboards_created_by` ON `metric_dashboards(created_by)`
- `idx_metric_dashboards_layout_config` ON `metric_dashboards USING GIN(layout_config)`

### 7. 审计管理表

#### audit_logs (审计日志表)
存储系统操作审计日志。

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGSERIAL | PRIMARY KEY | 日志ID |
| user_id | BIGINT | FK -> users(id) | 用户ID (可为空，系统操作) |
| username | VARCHAR(50) | | 用户名 (冗余字段，便于查询) |
| action | VARCHAR(50) | NOT NULL | 操作类型 (login, create, update, delete, execute) |
| resource_type | VARCHAR(50) | | 资源类型 (host, deployment, task, user, etc.) |
| resource_id | BIGINT | | 资源ID |
| resource_name | VARCHAR(255) | | 资源名称 (冗余字段) |
| ip_address | VARCHAR(45) | | IP地址 |
| user_agent | TEXT | | 用户代理 |
| request_data | JSONB | | 请求数据 (JSON格式) |
| response_data | JSONB | | 响应数据 (JSON格式) |
| before_data | JSONB | | 变更前数据 (update操作) |
| after_data | JSONB | | 变更后数据 (update操作) |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'success' | 操作状态: success, failed |
| error_message | TEXT | | 错误信息 |
| duration_ms | INTEGER | | 操作耗时(毫秒) |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() | 创建时间 |

**索引**:
- `idx_audit_logs_user_id` ON `audit_logs(user_id)`
- `idx_audit_logs_username` ON `audit_logs(username)`
- `idx_audit_logs_action` ON `audit_logs(action)`
- `idx_audit_logs_resource_type` ON `audit_logs(resource_type)`
- `idx_audit_logs_resource_id` ON `audit_logs(resource_id)`
- `idx_audit_logs_created_at` ON `audit_logs(created_at)` - 用于时间范围查询和分区
- `idx_audit_logs_status` ON `audit_logs(status)`

**约束**:
- `CHECK (status IN ('success', 'failed'))`
- `CHECK (action IN ('login', 'logout', 'create', 'update', 'delete', 'execute', 'query', 'export'))`

**分区策略** (PostgreSQL):
- 建议按 `created_at` 进行范围分区，每月一个分区
- 便于数据归档和查询性能优化

### 8. 表关联关系图

```
users ──┬── user_roles ── roles ── role_permissions ── permissions
        │
        ├── deployments (started_by)
        ├── deployment_configs (created_by)
        ├── scheduled_tasks (created_by)
        ├── ssh_sessions (user_id)
        ├── audit_logs (user_id)
        └── alert_rules (created_by)

hosts ──┬── host_group_members ── host_groups
        ├── host_tag_assignments ── host_tags
        ├── ssh_connections (host_id)
        ├── deployments (target_hosts[])
        └── scheduled_tasks (target_hosts[])

deployment_configs ── deployments
deployments ── deployments (rollback_from_deployment_id) [自关联]
salt_states ── deployment_configs (salt_state_files[])

scheduled_tasks ── task_executions
ssh_connections ── ssh_sessions
ssh_keys ── ssh_connections (key_id)
```

### 9. 数据库约束总结

#### 外键约束
- 所有外键使用 `ON DELETE RESTRICT` 防止误删
- 软删除的表使用 `deleted_at` 字段，不真正删除数据

#### 唯一约束
- `users.username`, `users.email`
- `hosts.salt_minion_id`
- `roles.name`
- `permissions.code`
- `host_groups.name`
- `host_tags.name`
- `ssh_keys.fingerprint`
- `ssh_sessions.session_token`
- `deployment_versions(application_name, version)`
- `user_roles(user_id, role_id)`
- `role_permissions(role_id, permission_id)`
- `host_group_members(host_id, group_id)`
- `host_tag_assignments(host_id, tag_id)`

#### 检查约束
- 状态字段使用枚举值检查
- 数值字段范围检查 (如 port, timeout_seconds)

#### 索引策略
- 主键自动创建索引
- 外键字段创建索引
- 常用查询字段创建索引
- JSONB 字段使用 GIN 索引
- 时间字段创建索引用于范围查询
- 软删除字段创建索引

### 10. 数据安全考虑

#### 敏感数据加密
- `users.password_hash`: 使用 bcrypt 或 argon2 哈希
- `ssh_connections.password_encrypted`: 使用 AES-256 加密
- `ssh_keys.private_key_encrypted`: 使用 AES-256 加密
- `ssh_keys.passphrase_encrypted`: 使用 AES-256 加密

#### 数据访问控制
- 所有表操作通过应用层权限控制
- 审计日志表只读，不允许修改
- 敏感字段在查询时进行脱敏处理

## Risks / Trade-offs

### Risk: Salt API 性能瓶颈
**Mitigation**: 
- 实现请求缓存机制
- 使用异步任务处理长时间操作
- 实现连接池和请求限流

### Risk: 大规模主机查询性能问题
**Mitigation**:
- 使用分页查询
- 实现数据库索引优化
- 使用 Redis 缓存热点数据

### Risk: 权限系统复杂度
**Mitigation**:
- 采用标准的 RBAC 模型
- 提供清晰的权限配置界面
- 编写详细的权限管理文档

### Risk: 审计日志数据量大
**Mitigation**:
- 实现日志归档机制
- 使用时间分区表
- 定期清理过期日志

### Trade-off: 实时性 vs 性能
- 监控数据采用缓存策略，牺牲部分实时性换取性能
- 关键操作（如发布）保持实时性

## Migration Plan

### Phase 1: 基础功能（1-2 个月）
- 完成项目初始化和数据库设计
- 实现主机管理、SSH 管理基础功能
- 实现权限管理和审计管理

### Phase 2: 核心功能（2-3 个月）
- 实现发布管理和定时任务
- 完成 Salt 集成
- 完成前端核心页面开发

### Phase 3: 集成功能（1-2 个月）
- 完成 ELK 集成
- 完成 Prometheus 集成
- 完善前端功能

### Phase 4: 优化和测试（1 个月）
- 性能优化
- 完善测试
- 文档编写

### Phase 5: K8s 管理（后续）
- 根据需求设计 K8s 管理功能
- 实现 K8s 相关 API 和前端页面

## Open Questions
- Salt Master 的高可用部署方案？
- 是否需要支持多 Salt Master 环境？
- 日志保留策略和归档方案？
- 是否需要支持自定义 Salt State 模板？
- K8s 管理的具体功能范围？

