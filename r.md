# 智能运维管理平台规划文档

## 1. 项目概述

### 1.1 项目目标
构建一个智能运维管理平台，实现IT资产的统一管理和运维操作的自动化执行，提高运维效率，降低运维成本。

### 1.2 核心功能
- **用户管理**：用户账号、角色、组织架构管理
- **权限管理**：基于角色的访问控制（RBAC），细粒度权限控制
- **资产管理**：服务器、网络设备、应用系统等IT资产的统一管理
- **运维执行**：运维任务的创建、执行、监控和历史记录
- **WebSSH**：基于浏览器的SSH终端，支持远程服务器连接和命令执行

### 1.3 技术栈

#### 后端
- **语言**：Golang 1.21+
- **Web框架**：Gin / Echo / Fiber
- **ORM**：GORM
- **数据库**：PostgreSQL（主数据库）、Redis（缓存/会话）
- **认证授权**：JWT
- **API文档**：Swagger/OpenAPI
- **配置管理**：Viper
- **日志**：Zap
- **任务调度**：Cron
- **WebSocket**：gorilla/websocket（WebSSH实时通信）

#### 前端
- **框架**：React 18+
- **UI组件库**：Ant Design 5+
- **状态管理**：Redux Toolkit / Zustand
- **路由**：React Router v6
- **HTTP客户端**：Axios
- **WebSocket客户端**：原生WebSocket API 或 socket.io-client
- **终端组件**：xterm.js（WebSSH终端）
- **主题切换**：Ant Design Theme Provider（支持亮色/暗色主题）
- **构建工具**：Vite
- **TypeScript**：类型支持

## 2. 系统架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                     前端层 (React)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │用户管理  │  │权限管理  │  │资产管理  │  │运维执行  │  │ WebSSH   │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/REST API
┌────────────────────┴────────────────────────────────────┐
│                   后端层 (Golang)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │用户服务  │  │权限服务  │  │资产服务  │  │执行服务  │  │WebSSH服务│ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
│  ┌──────────────────────────────────────────────────┐   │
│  │              中间件层 (Middleware)                │   │
│  │  认证/授权/日志/限流/错误处理/CORS               │   │
│  └──────────────────────────────────────────────────┘   │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────────────┐
│                   数据层                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │PostgreSQL│  │  Redis   │  │ 文件存储  │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
```

### 2.2 前后端分离架构

- **前端**：独立的SPA应用，通过RESTful API与后端通信
- **后端**：提供统一的API网关，各模块以微服务或模块化方式组织
- **通信协议**：HTTPS + RESTful API（JSON数据格式），WebSocket（WebSSH实时通信）
- **认证方式**：
  - Web端：JWT Token（存储在HTTP-only Cookie或LocalStorage）
  - API调用：API Token（通过Authorization Header传递，格式：`Bearer <token>`）

## 3. 模块设计

### 3.1 用户管理模块

#### 3.1.1 功能特性
- 用户账号管理（增删改查）
- 用户信息管理（姓名、邮箱、手机、部门等）
- 用户状态管理（启用/禁用）
- 用户密码管理（重置、修改）
- 用户组织架构管理（部门、团队）
- **API Token管理**：用户可生成和管理自己的API Token，用于程序化API调用

#### 3.1.2 数据模型
```sql
-- 用户表
users (
  id, username, password_hash, email, phone, 
  real_name, department_id, status, avatar_url,
  created_at, updated_at, last_login_at
)

-- 部门表
departments (
  id, name, parent_id, code, description,
  created_at, updated_at
)

-- API Token表
api_tokens (
  id, user_id, name, token_hash, 
  expires_at, last_used_at, status,
  description, created_at, updated_at
)
```

**说明**：
- `api_tokens.token_hash`：Token的哈希值（用于验证，类似密码哈希）
- `api_tokens.name`：Token名称/描述（便于用户识别）
- `api_tokens.expires_at`：过期时间（NULL表示永不过期）
- `api_tokens.status`：状态（active-启用，disabled-禁用）
- `api_tokens.last_used_at`：最后使用时间（用于审计）
- Token生成时返回完整Token（仅显示一次），后续只显示Token前缀

#### 3.1.3 API设计
- `POST /api/v1/users` - 创建用户
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户
- `POST /api/v1/users/:id/reset-password` - 重置密码
- `PUT /api/v1/users/:id/password` - 修改密码
- `PUT /api/v1/users/:id/status` - 更新用户状态
- `GET /api/v1/users/:id/tokens` - 获取用户的API Token列表
- `POST /api/v1/users/:id/tokens` - 为用户创建API Token
- `PUT /api/v1/tokens/:id` - 更新API Token信息（名称、过期时间等）
- `DELETE /api/v1/tokens/:id` - 删除/撤销API Token
- `PUT /api/v1/tokens/:id/status` - 启用/禁用API Token

### 3.2 权限管理模块

#### 3.2.1 功能特性
- 角色管理（角色创建、编辑、删除）
- 权限管理（功能权限、数据权限）
- 用户角色分配
- 角色权限分配
- 权限验证中间件

#### 3.2.2 RBAC模型
```
用户 (User) ──N:N──> 角色 (Role) ──N:N──> 权限 (Permission)
```

#### 3.2.3 数据模型
```sql
-- 角色表
roles (
  id, name, code, description, 
  created_at, updated_at
)

-- 权限表
permissions (
  id, name, code, resource, action, description,
  created_at, updated_at
)

-- 用户角色关联表
user_roles (
  user_id, role_id, created_at
)

-- 角色权限关联表
role_permissions (
  role_id, permission_id, created_at
)
```

#### 3.2.4 API设计
- `GET /api/v1/roles` - 获取角色列表
- `POST /api/v1/roles` - 创建角色
- `PUT /api/v1/roles/:id` - 更新角色
- `DELETE /api/v1/roles/:id` - 删除角色
- `GET /api/v1/permissions` - 获取权限列表
- `POST /api/v1/roles/:id/permissions` - 分配权限给角色
- `POST /api/v1/users/:id/roles` - 分配角色给用户
- `GET /api/v1/users/:id/permissions` - 获取用户权限列表

### 3.3 资产管理模块

#### 3.3.1 功能特性
- 资产分类管理（服务器、网络设备、存储设备、应用系统等）
- 资产信息管理（基本信息、配置信息、关联关系）
- 项目关联管理（资产归属项目）
- 云平台标识（AWS、阿里云、腾讯云、Azure等）
- SSH连接信息管理（IP、端口、关联SSH密钥）
- 硬件资源信息（CPU、内存、磁盘）
- 资产状态管理（激活、禁用）
- 标签管理（用于区分应用类型，如：web服务、数据库、缓存等）
- 环境区分（dev开发、test测试、uat预发布、prod生产）
- 资产搜索和筛选（支持按项目、云平台、IP、标签、环境等筛选）
- 资产导入导出

#### 3.3.2 数据模型
```sql
-- 项目表
projects (
  id, name, code, description,
  created_by, created_at, updated_at
)

-- 资产分类表
asset_categories (
  id, name, code, parent_id, description,
  created_at, updated_at
)

-- 资产表
assets (
  id, category_id, name, code, brand, model,
  serial_number, status, location, 
  purchase_date, warranty_expiry, price,
  owner_id, department_id, project_id,
  cloud_platform, environment, ip, ssh_port, ssh_key_id,
  cpu, memory, disk, description,
  created_at, updated_at
)

-- 资产配置表（扩展属性，JSON格式存储）
asset_configs (
  id, asset_id, config_data (JSON),
  created_at, updated_at
)

-- 资产关联关系表
asset_relations (
  id, asset_id, related_asset_id, relation_type,
  created_at, updated_at
)

-- 标签表
tags (
  id, name, color, description,
  created_at, updated_at
)

-- 资产标签关联表（多对多）
asset_tags (
  asset_id, tag_id,
  created_at
)
```

**字段说明**：
- `status`：资产状态（枚举：`active`-激活，`disabled`-禁用），默认激活
- `project_id`：关联项目表，资产所属项目（可选）
- `cloud_platform`：云平台标识（字符串，如：aws、aliyun、tencent、azure、on-premise等）
- `environment`：环境类型（枚举：`dev`-开发、`test`-测试、`uat`-预发布、`prod`-生产），可选
- `ip`：IP地址（字符串，支持IPv4/IPv6，可多个IP用逗号分隔）
- `ssh_port`：SSH端口（整数，默认22）
- `ssh_key_id`：关联SSH密钥表，用于SSH连接认证（可选）
- `cpu`：CPU信息（字符串，如："4核" 或 JSON格式存储详细规格）
- `memory`：内存大小（整数，单位MB；或字符串，如："8GB"）
- `disk`：磁盘信息（字符串，如："500GB SSD" 或 JSON格式存储详细规格）

#### 3.3.3 API设计
- `GET /api/v1/projects` - 获取项目列表
- `POST /api/v1/projects` - 创建项目
- `PUT /api/v1/projects/:id` - 更新项目
- `DELETE /api/v1/projects/:id` - 删除项目
- `GET /api/v1/asset-categories` - 获取资产分类
- `GET /api/v1/tags` - 获取标签列表
- `POST /api/v1/tags` - 创建标签
- `PUT /api/v1/tags/:id` - 更新标签
- `DELETE /api/v1/tags/:id` - 删除标签
- `POST /api/v1/assets` - 创建资产
- `GET /api/v1/assets` - 获取资产列表（支持搜索、筛选、分页）
  - 筛选参数：`project_id`, `cloud_platform`, `environment`, `status`, `ip`, `ssh_key_id`, `tag_ids`等
- `GET /api/v1/assets/:id` - 获取资产详情
- `PUT /api/v1/assets/:id` - 更新资产信息
- `DELETE /api/v1/assets/:id` - 删除资产
- `PUT /api/v1/assets/:id/status` - 更新资产状态
- `POST /api/v1/assets/:id/tags` - 为资产添加标签
- `DELETE /api/v1/assets/:id/tags/:tagId` - 移除资产的标签
- `GET /api/v1/assets/:id/tags` - 获取资产的标签列表
- `POST /api/v1/assets/import` - 批量导入资产
- `GET /api/v1/assets/export` - 导出资产数据
- `GET /api/v1/assets/:id/relations` - 获取资产关联关系

### 3.4 运维执行模块

#### 3.4.1 功能特性
- 运维任务创建（脚本执行、命令执行、文件传输等）
- 任务模板管理
- 任务执行（同步/异步执行）
- 任务执行历史记录
- 任务执行结果查看（输出、日志）
- 任务调度（定时任务）
- 执行主机管理（目标服务器配置）

#### 3.4.2 数据模型
```sql
-- 执行主机表
execution_hosts (
  id, name, host, port, username, 
  auth_type, auth_config (JSON), status,
  created_at, updated_at
)

-- 任务模板表
task_templates (
  id, name, type, content, description,
  created_by, created_at, updated_at
)

-- 任务表
tasks (
  id, name, type, template_id, content,
  execution_type, scheduled_at, status,
  created_by, created_at, updated_at, started_at, completed_at
)

-- 任务执行记录表
task_executions (
  id, task_id, host_id, status, 
  output, error, exit_code,
  started_at, completed_at, duration
)

-- 任务执行日志表
task_execution_logs (
  id, execution_id, level, message, timestamp
)
```

#### 3.4.3 API设计
- `GET /api/v1/execution-hosts` - 获取执行主机列表
- `POST /api/v1/execution-hosts` - 添加执行主机
- `PUT /api/v1/execution-hosts/:id` - 更新执行主机
- `DELETE /api/v1/execution-hosts/:id` - 删除执行主机
- `GET /api/v1/task-templates` - 获取任务模板列表
- `POST /api/v1/task-templates` - 创建任务模板
- `GET /api/v1/tasks` - 获取任务列表
- `POST /api/v1/tasks` - 创建任务
- `POST /api/v1/tasks/:id/execute` - 执行任务
- `GET /api/v1/tasks/:id/executions` - 获取任务执行记录
- `GET /api/v1/task-executions/:id` - 获取执行详情
- `GET /api/v1/task-executions/:id/logs` - 获取执行日志
- `POST /api/v1/tasks/:id/cancel` - 取消任务执行

### 3.5 认证授权模块（基础模块）

#### 3.5.1 功能特性
- 用户登录（用户名/密码）
- JWT Token生成和验证（Web端会话）
- API Token认证（程序化API调用）
- 密码加密存储（bcrypt）
- 登录日志记录
- Token刷新机制

#### 3.5.2 API设计
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `POST /api/v1/auth/refresh` - 刷新Token
- `GET /api/v1/auth/me` - 获取当前用户信息

### 3.6 WebSSH模块

#### 3.6.1 功能特性
- 基于浏览器的SSH终端连接
- 支持密码和密钥认证
- **SSH密钥管理**：密钥的增删改查、导入、测试连接，密钥关联默认连接用户名
- 终端会话管理（创建、断开、重连）
- 多标签页支持（同时连接多台服务器）
- 终端窗口大小调整（动态调整PTY大小）
- 会话记录和回放
- 文件传输支持（SFTP，lrzsz支持）
- 会话权限控制（基于RBAC）

#### 3.6.2 技术实现
- **前端**：使用 xterm.js 实现终端界面
- **后端**：使用 golang.org/x/crypto/ssh 建立SSH连接
- **通信协议**：WebSocket（双向通信）
- **数据流**：前端 → WebSocket → 后端 → SSH → 目标服务器

#### 3.6.3 数据模型
```sql
-- SSH密钥表
ssh_keys (
  id, user_id, name, type, username,
  public_key, private_key_encrypted, passphrase_encrypted,
  fingerprint, description,
  created_at, updated_at, last_used_at
)

-- WebSSH会话表
ssh_sessions (
  id, user_id, host_id, asset_id, ssh_key_id,
  session_id, status, 
  connected_at, disconnected_at, duration,
  created_at, updated_at
)

-- WebSSH会话记录表（可选，用于审计）
ssh_session_logs (
  id, session_id, user_id, action_type,
  content, timestamp
)
```

**说明**：
- `ssh_keys.username`：默认连接用户名（使用该密钥连接时的默认用户名）
- `ssh_keys.private_key_encrypted`：私钥加密存储（AES加密）
- `ssh_keys.passphrase_encrypted`：密钥密码加密存储（如果密钥有密码）
- `ssh_keys.type`：密钥类型（rsa, ed25519, ecdsa等）
- `ssh_keys.fingerprint`：公钥指纹，用于快速识别和去重
- `ssh_sessions.ssh_key_id`：关联使用的SSH密钥（可选，密码认证时为空）
- 使用密钥连接时，默认使用密钥关联的用户名，连接时可覆盖
- `host_id` 关联执行主机表（execution_hosts）
- `asset_id` 关联资产表（assets），可选
- `session_id` 用于WebSocket连接标识
- 会话日志可用于审计和安全追踪

#### 3.6.4 SSH密钥管理API设计
- `GET /api/v1/ssh/keys` - 获取SSH密钥列表（当前用户）
- `POST /api/v1/ssh/keys` - 创建/导入SSH密钥（需指定username字段）
- `GET /api/v1/ssh/keys/:id` - 获取SSH密钥详情（不含私钥）
- `PUT /api/v1/ssh/keys/:id` - 更新SSH密钥信息（可更新username）
- `DELETE /api/v1/ssh/keys/:id` - 删除SSH密钥
- `POST /api/v1/ssh/keys/:id/test` - 测试SSH密钥连接（使用密钥关联的用户名或指定用户名）
- `POST /api/v1/ssh/keys/generate` - 生成新的SSH密钥对（可选，需指定username）

#### 3.6.5 WebSocket API设计
- `WS /ws/ssh/connect` - 建立WebSSH连接
  - 连接参数：`host_id`, `username`, `auth_type`, `auth_config`
  - `auth_type`: `password` 或 `key`
  - `auth_config`: 
    - 密码认证：`{"password": "xxx"}`
    - 密钥认证：`{"ssh_key_id": 1, "username": "xxx", "passphrase": "xxx"}` 
      - `username`：可选，如果未指定则使用密钥关联的默认用户名
      - `passphrase`：可选，如果密钥有密码保护
  - 消息类型：`connect`, `input`, `resize`, `disconnect`
  
- `POST /api/v1/ssh/sessions` - 创建SSH会话记录
- `GET /api/v1/ssh/sessions` - 获取SSH会话列表
- `GET /api/v1/ssh/sessions/:id` - 获取SSH会话详情
- `DELETE /api/v1/ssh/sessions/:id` - 关闭SSH会话
- `GET /api/v1/ssh/sessions/:id/logs` - 获取SSH会话日志（如果启用）

#### 3.6.6 WebSocket消息协议

**客户端 → 服务器**：
```json
{
  "type": "connect",
  "data": {
    "host_id": 1,
    "username": "root",
    "auth_type": "password",
    "auth_config": {
      "password": "xxx"
    }
  }
}

{
  "type": "connect",
  "data": {
    "host_id": 1,
    "username": "root",
    "auth_type": "key",
    "auth_config": {
      "ssh_key_id": 1,
      "passphrase": "xxx"
    }
  }
}

{
  "type": "input",
  "data": {
    "data": "ls -la\n"
  }
}

{
  "type": "resize",
  "data": {
    "rows": 24,
    "cols": 80
  }
}

{
  "type": "disconnect"
}
```

**服务器 → 客户端**：
```json
{
  "type": "output",
  "data": "root@server:~$ "
}

{
  "type": "error",
  "data": "Connection failed: ..."
}

{
  "type": "connected",
  "data": {
    "session_id": "xxx"
  }
}

{
  "type": "disconnected",
  "data": {
    "reason": "User disconnected"
  }
}
```

## 4. 前端设计

### 4.1 项目结构
```
frontend/
├── src/
│   ├── api/              # API接口定义
│   ├── components/       # 公共组件
│   ├── layouts/          # 布局组件
│   ├── pages/            # 页面组件
│   │   ├── auth/         # 认证相关
│   │   ├── users/        # 用户管理
│   │   ├── roles/        # 角色权限
│   │   ├── assets/       # 资产管理
│   │   ├── tasks/        # 运维执行
│   │   └── ssh/          # WebSSH终端
│   ├── stores/           # 状态管理
│   ├── hooks/            # 自定义Hooks
│   ├── utils/            # 工具函数
│   ├── themes/           # 主题配置
│   ├── App.tsx           # 根组件
│   └── main.tsx          # 入口文件
├── public/
└── package.json
```

### 4.2 主题切换实现

使用Ant Design的ConfigProvider和主题配置：

```typescript
// themes/index.ts
import { ThemeConfig } from 'antd';

export const lightTheme: ThemeConfig = {
  token: {
    colorPrimary: '#1890ff',
    // 亮色主题配置
  },
};

export const darkTheme: ThemeConfig = {
  token: {
    colorPrimary: '#1890ff',
    // 暗色主题配置
  },
  algorithm: theme.darkAlgorithm,
};
```

主题切换通过Context或状态管理实现，用户选择后持久化到LocalStorage。

### 4.3 路由设计
- `/login` - 登录页
- `/dashboard` - 仪表盘
- `/users` - 用户管理
- `/roles` - 角色权限管理
- `/assets` - 资产管理
- `/assets/:id` - 资产详情
- `/tasks` - 运维任务
- `/tasks/:id` - 任务详情
- `/tasks/:id/executions/:executionId` - 执行记录详情
- `/ssh` - WebSSH终端（支持多标签页，独立菜单）
- `/ssh/connect/:hostId` - 连接到指定主机
- `/ssh/keys` - SSH密钥管理
- `/ssh/keys/:id` - SSH密钥详情/编辑

## 5. 后端设计

### 5.1 项目结构
```
backend/
├── cmd/
│   └── server/
│       └── main.go        # 应用入口
├── internal/
│   ├── config/            # 配置管理
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   │   ├── user/
│   │   ├── auth/
│   │   ├── role/
│   │   ├── asset/
│   │   ├── task/
│   │   └── ssh/
│   ├── handler/           # HTTP处理器
│   │   ├── user/
│   │   ├── auth/
│   │   ├── role/
│   │   ├── asset/
│   │   ├── task/
│   │   └── ssh/
│   ├── middleware/        # 中间件
│   │   ├── auth.go
│   │   ├── permission.go
│   │   ├── logger.go
│   │   └── cors.go
│   └── utils/             # 工具函数
├── pkg/                   # 公共包
├── api/                   # API文档
├── migrations/            # 数据库迁移
├── configs/               # 配置文件
├── go.mod
└── go.sum
```

### 5.2 技术选型说明

- **Web框架**：推荐Gin（性能好、生态丰富）
- **ORM**：GORM（功能完善、易用）
- **数据库**：PostgreSQL（支持JSON、事务、性能好）
- **缓存**：Redis（会话存储、缓存）

## 6. 数据库设计

### 6.1 核心表关系图

```
users ──N:N── user_roles ──N:1── roles ──N:N── role_permissions ──N:1── permissions
users ──N:1── departments
users ──1:N── api_tokens
users ──1:N── projects (created_by)
users ──1:N── assets (owner)
assets ──N:1── asset_categories
assets ──N:1── projects
assets ──N:1── ssh_keys
assets ──N:N── asset_tags ──N:1── tags
tasks ──N:1── task_templates
tasks ──1:N── task_executions ──N:1── execution_hosts
task_executions ──1:N── task_execution_logs
users ──1:N── ssh_sessions ──N:1── execution_hosts
users ──1:N── ssh_keys
ssh_sessions ──N:1── ssh_keys
ssh_sessions ──1:N── ssh_session_logs
```

### 6.2 索引设计

- `users.username` - 唯一索引
- `users.email` - 唯一索引
- `api_tokens.token_hash` - 唯一索引
- `api_tokens.user_id`, `api_tokens.status` - 复合索引
- `projects.code` - 唯一索引
- `assets.code` - 唯一索引
- `assets.ip` - 索引（用于搜索）
- `assets.project_id` - 索引
- `assets.cloud_platform` - 索引
- `assets.environment` - 索引
- `assets.ssh_key_id` - 索引
- `tags.name` - 唯一索引
- `asset_tags.asset_id`, `asset_tags.tag_id` - 复合唯一索引
- `tasks.status`, `tasks.created_at` - 复合索引
- `task_executions.task_id`, `task_executions.status` - 复合索引
- `ssh_sessions.user_id`, `ssh_sessions.status` - 复合索引
- `ssh_sessions.session_id` - 唯一索引
- `ssh_keys.user_id`, `ssh_keys.name` - 复合索引
- `ssh_keys.fingerprint` - 索引（用于去重）

## 7. 安全设计

### 7.1 认证安全
- 密码使用bcrypt加密存储
- JWT Token设置合理的过期时间（Web端会话）
- 支持Token刷新机制
- 登录失败次数限制
- **API Token安全**：
  - Token使用SHA-256或bcrypt哈希存储（类似密码）
  - Token生成后仅显示一次完整Token，后续只显示前缀
  - 支持Token过期时间设置
  - Token可随时撤销（禁用或删除）
  - Token权限继承用户权限（RBAC）
  - 记录Token使用日志（last_used_at）
  - API调用时通过Header传递：`Authorization: Bearer <token>`

### 7.2 授权安全
- 基于RBAC的权限控制
- API接口级别的权限验证
- 前端路由级别的权限控制
- 数据权限隔离

### 7.3 数据安全
- SQL注入防护（使用ORM参数化查询）
- XSS防护（输入验证、输出转义）
- CSRF防护（Token验证）
- HTTPS传输加密

### 7.4 WebSSH安全
- WebSocket连接认证（JWT Token验证）
- SSH会话权限控制（基于RBAC，限制用户可访问的主机）
- 会话超时自动断开
- 敏感命令拦截和审计（可选）
- 会话日志记录（用于安全审计）
- **SSH密钥安全管理**：
  - 私钥使用AES-256加密存储（密钥由系统密钥派生）
  - 密钥密码（passphrase）加密存储
  - 私钥仅在连接时临时解密，不在前端暴露
  - 密钥传输使用HTTPS/WSS加密通道
  - 支持密钥访问权限控制（用户只能访问自己的密钥）
  - 密钥删除前验证是否被使用

## 8. 开发计划

### 8.1 第一阶段：基础框架搭建（1-2周）
- [ ] 项目初始化（前后端项目结构）
- [ ] 数据库设计和迁移脚本
- [ ] 基础认证模块（登录、JWT）
- [ ] 前端框架搭建（路由、布局、主题）
- [ ] API基础框架（中间件、错误处理）

### 8.2 第二阶段：用户和权限管理（2-3周）
- [ ] 用户管理模块（后端API + 前端页面）
- [ ] 角色权限管理模块
- [ ] API Token管理功能（生成、管理、认证）
- [ ] 权限中间件实现（JWT + API Token认证）
- [ ] 前端权限控制（路由守卫、按钮权限）

### 8.3 第三阶段：资产管理（2-3周）
- [ ] 项目管理功能（项目CRUD）
- [ ] 资产分类管理
- [ ] 标签管理功能（标签CRUD、标签与资产关联）
- [ ] 资产CRUD功能（包含项目、云平台、SSH信息、硬件资源字段、标签）
- [ ] 资产搜索和筛选（支持多维度筛选，包括标签筛选）
- [ ] 资产导入导出
- [ ] 资产详情页（展示完整信息，包括SSH连接信息、标签）
- [ ] SSH密钥关联（资产关联SSH密钥选择）

### 8.4 第四阶段：运维执行（3-4周）
- [ ] 执行主机管理
- [ ] 任务模板管理
- [ ] 任务创建和执行
- [ ] 任务执行记录和日志
- [ ] 实时日志输出（WebSocket支持）

### 8.5 第五阶段：WebSSH功能（2-3周）
- [ ] SSH密钥管理功能
  - [ ] 密钥数据模型和存储（加密存储）
  - [ ] 密钥管理API（增删改查、导入、测试）
  - [ ] 密钥管理前端页面
  - [ ] 密钥生成功能（可选）
- [ ] WebSocket服务实现
- [ ] SSH连接后端服务（SSH客户端封装，支持密钥和密码认证）
- [ ] 前端终端组件（xterm.js集成）
- [ ] SSH会话管理（创建、断开、重连）
- [ ] 多标签页支持
- [ ] 终端窗口大小调整
- [ ] 会话记录和权限控制

### 8.6 第六阶段：测试和优化（1-2周）
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能优化
- [ ] 安全审计
- [ ] 文档完善

## 9. 技术难点和解决方案

### 9.1 运维任务执行
- **问题**：如何安全、可靠地执行远程命令/脚本
- **方案**：使用SSH库（如golang.org/x/crypto/ssh），支持密钥认证，任务队列管理

### 9.2 实时日志输出
- **问题**：任务执行过程中的实时日志展示
- **方案**：WebSocket推送或Server-Sent Events (SSE)

### 9.3 大文件上传（资产导入）
- **问题**：批量导入资产时的大文件处理
- **方案**：分片上传、异步处理、进度反馈

### 9.4 权限性能
- **问题**：频繁的权限验证可能影响性能
- **方案**：权限信息缓存到Redis，合理设置缓存过期时间

### 9.5 WebSSH实现
- **问题1**：WebSocket长连接管理和并发控制
- **方案**：使用连接池管理SSH连接，限制单用户并发连接数，实现连接心跳检测

- **问题2**：SSH会话的PTY大小动态调整
- **方案**：监听浏览器窗口resize事件，通过WebSocket发送resize消息到后端，后端调用SSH Session的WindowChange方法

- **问题3**：SSH密钥的安全存储和传输
- **方案**：密钥加密存储（使用AES加密），传输时使用HTTPS/WSS，不在前端存储明文密钥

- **问题4**：大量并发SSH连接的性能问题
- **方案**：使用goroutine池，连接复用，合理设置超时和缓冲区大小

- **问题5**：SSH密钥的安全存储和访问
- **方案**：
  - 使用AES-256加密存储私钥（加密密钥从系统环境变量或密钥管理服务获取）
  - 密钥仅在连接时临时解密到内存，使用后立即清除
  - 实现密钥访问控制，用户只能访问自己的密钥
  - 密钥传输使用HTTPS/WSS，不在前端存储或显示私钥明文

## 10. 扩展规划

### 10.1 后续可扩展功能
- 监控告警（集成Prometheus、Grafana）
- 自动化运维（工作流引擎）
- 工单系统
- 审计日志
- 报表统计
- 移动端支持

### 10.2 技术升级方向
- 微服务架构改造
- 容器化部署（Docker、K8s）
- CI/CD集成
- 分布式任务调度
- 消息队列（任务异步处理）

---

**文档版本**：v1.1  
**创建日期**：2025-01-27  
**最后更新**：2025-01-27
