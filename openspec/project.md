# Project Context

## Purpose
KkOps 是一个基于 Salt 的运维中台管理系统，旨在提供统一的主机管理、SSH 管理、发布管理、定时任务、日志管理、监控管理、权限管理和审计管理等核心功能，后续将扩展 K8s 相关管理能力。

## Tech Stack
- **后端**: Go (推荐使用 Go 1.21+)
- **前端**: React + Ant Design
- **配置管理**: Salt (SaltStack)
- **日志系统**: ELK (Elasticsearch, Logstash, Kibana)
- **监控系统**: Prometheus
- **数据库**: PostgreSQL (推荐) 或 MySQL
- **缓存**: Redis

## Project Conventions

### Code Style
- **Go**: 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
  - 使用 `gofmt` 格式化代码
  - 使用 `golangci-lint` 进行代码检查
  - 包名使用小写，简短有意义
  - 导出函数和类型使用大写开头
- **React/TypeScript**: 
  - 使用 TypeScript 进行类型检查
  - 遵循 React Hooks 最佳实践
  - 组件使用函数式组件
  - 使用 ESLint + Prettier 进行代码格式化

### Architecture Patterns
- **后端架构**: 
  - 采用分层架构：Handler -> Service -> Repository
  - 使用依赖注入模式
  - API 遵循 RESTful 设计规范
  - 使用中间件处理认证、日志、错误处理
- **前端架构**:
  - 使用 React Hooks 和函数式组件
  - 状态管理使用 Context API 或 Redux（如需要）
  - 组件按功能模块组织
  - 使用 Ant Design 作为 UI 组件库

### Testing Strategy
- **单元测试**: 核心业务逻辑覆盖率 ≥ 80%
- **集成测试**: 关键路径必须 100% 覆盖
- **端到端测试**: 主要用户流程需要 E2E 测试
- Go 测试使用 `testing` 包，前端使用 Jest + React Testing Library

### Git Workflow
- 遵循 Conventional Commits 规范
- 主分支：`main`
- 功能分支：`feature/xxx`
- 修复分支：`fix/xxx`
- 通过 Pull Request 合并代码

## Domain Context
- **Salt**: SaltStack 是一个配置管理和远程执行系统，通过 Salt Master 和 Salt Minion 架构实现主机管理
- **运维中台**: 提供统一的运维操作入口，整合各类运维工具和系统
- **主机管理**: 管理服务器资产信息、分组、标签等
- **发布管理**: 基于 Salt 实现应用的自动化部署和发布流程
- **定时任务**: 支持 Salt 命令的定时执行和任务调度

## Important Constraints
- 必须与 Salt Master 集成，通过 Salt API 或 Salt CLI 执行操作
- 需要支持大规模主机管理（1000+ 主机）
- 所有操作需要记录审计日志
- 权限管理需要支持 RBAC（基于角色的访问控制）
- 前端需要支持响应式设计，适配不同屏幕尺寸

## External Dependencies
- **Salt Master**: SaltStack 主控节点，提供配置管理和远程执行能力
- **ELK Stack**: 
  - Elasticsearch: 日志存储和检索
  - Logstash: 日志收集和处理
  - Kibana: 日志可视化和分析
- **Prometheus**: 监控指标收集和存储
- **PostgreSQL/MySQL**: 应用数据存储
- **Redis**: 缓存和会话存储

## Project Structure

### Root Directory
```
KkOps/
├── backend/              # Go 后端应用
├── frontend/             # React 前端应用
├── openspec/             # OpenSpec 规范和变更文档
├── scripts/              # Ruby 辅助脚本（工具脚本统一使用 Ruby）
├── vms/                  # Vagrant 虚拟机配置（用于 Salt 测试环境）
├── docker-compose.yml    # Docker Compose 配置
├── Makefile              # Make 构建脚本
├── Rakefile              # Rake 任务管理（Ruby）
├── push.rb               # Git 推送脚本（Ruby）
└── README.md             # 项目说明文档
```

### Backend Structure (`backend/`)
```
backend/
├── cmd/
│   └── api/
│       └── main.go              # 应用入口点
├── internal/
│   ├── config/
│   │   └── config.go            # 配置加载和管理
│   ├── handler/                 # HTTP 请求处理器层（Controller）
│   │   ├── audit_handler.go     # 审计日志处理器
│   │   ├── auth_handler.go      # 认证处理器
│   │   ├── cloud_platform_handler.go  # 云平台处理器
│   │   ├── deployment_handler.go      # 发布管理处理器
│   │   ├── environment_handler.go     # 环境管理处理器
│   │   ├── host_group_handler.go      # 主机分组处理器
│   │   ├── host_handler.go            # 主机管理处理器
│   │   ├── host_tag_handler.go        # 主机标签处理器
│   │   ├── permission_handler.go      # 权限管理处理器
│   │   ├── project_handler.go         # 项目管理处理器
│   │   ├── role_handler.go            # 角色管理处理器
│   │   ├── salt_handler.go            # Salt 集成处理器
│   │   ├── settings_handler.go        # 系统设置处理器
│   │   ├── ssh_key_handler.go         # SSH 密钥处理器
│   │   ├── user_handler.go            # 用户管理处理器
│   │   └── webssh_handler.go          # WebSSH 终端处理器
│   ├── middleware/              # HTTP 中间件
│   │   ├── audit.go             # 审计日志中间件
│   │   ├── auth.go              # 认证中间件
│   │   └── rbac.go              # RBAC 权限中间件
│   ├── models/                  # 数据模型（GORM）
│   │   ├── audit.go             # 审计日志模型
│   │   ├── cloud_platform.go    # 云平台模型
│   │   ├── database.go          # 数据库连接和初始化
│   │   ├── deployment.go        # 发布管理模型
│   │   ├── environment.go       # 环境模型
│   │   ├── host.go              # 主机模型
│   │   ├── project.go           # 项目模型
│   │   ├── settings.go          # 系统设置模型
│   │   ├── ssh_key.go           # SSH 密钥模型
│   │   └── user.go              # 用户、角色、权限模型
│   ├── repository/              # 数据访问层（Repository Pattern）
│   │   ├── audit_repository.go
│   │   ├── cloud_platform_repository.go
│   │   ├── deployment_repository.go
│   │   ├── environment_repository.go
│   │   ├── host_group_repository.go
│   │   ├── host_repository.go
│   │   ├── host_tag_repository.go
│   │   ├── permission_repository.go
│   │   ├── project_repository.go
│   │   ├── role_repository.go
│   │   ├── settings_repository.go
│   │   ├── ssh_key_repository.go
│   │   └── user_repository.go
│   ├── service/                 # 业务逻辑层（Service Layer）
│   │   ├── audit_service.go
│   │   ├── auth_service.go
│   │   ├── cloud_platform_service.go
│   │   ├── deployment_service.go
│   │   ├── environment_service.go
│   │   ├── host_group_service.go
│   │   ├── host_service.go
│   │   ├── host_tag_service.go
│   │   ├── permission_service.go
│   │   ├── project_service.go
│   │   ├── role_service.go
│   │   ├── scheduler_service.go      # 定时任务服务
│   │   ├── settings_service.go
│   │   ├── ssh_key_service.go
│   │   └── user_service.go
│   ├── salt/                    # Salt 集成模块
│   │   ├── client.go            # Salt API 客户端
│   │   └── manager.go           # Salt 管理器（封装常用操作）
│   └── utils/                   # 工具函数
│       ├── encrypt.go           # 加密工具（AES 加密）
│       ├── jwt.go               # JWT 令牌生成和验证
│       └── password.go          # 密码哈希和验证
├── migrations/                  # 数据库迁移脚本（SQL）
│   ├── add_ssh_keys_table.sql
│   ├── add_ssh_keys_username.sql
│   ├── fix_host_tag_assignments.sql
│   ├── remove_ssh_tables.sql
│   └── README.md
├── Dockerfile                   # 后端 Docker 镜像构建
├── go.mod                       # Go 模块依赖
├── go.sum                       # Go 依赖校验和
└── README.md                    # 后端说明文档
```

**后端架构说明**:
- **Handler 层**: 处理 HTTP 请求，参数验证，调用 Service 层
- **Service 层**: 业务逻辑处理，事务管理，调用 Repository 层
- **Repository 层**: 数据访问封装，GORM 操作，数据库交互
- **Model 层**: 数据模型定义，GORM 标签，关联关系
- **Middleware**: 请求拦截，认证授权，审计日志
- **Utils**: 通用工具函数，JWT、加密、密码处理

### Frontend Structure (`frontend/`)
```
frontend/
├── public/                      # 静态资源
│   ├── index.html
│   ├── favicon.ico
│   └── manifest.json
├── src/
│   ├── components/              # 可复用组件
│   │   ├── AuditLogSettings.tsx  # 审计日志设置组件
│   │   ├── Footer.tsx            # 页脚组件
│   │   ├── MainLayout.tsx        # 主布局组件（左侧菜单）
│   │   ├── PrivateRoute.tsx      # 私有路由组件（权限控制）
│   │   └── Terminal.tsx          # WebSSH 终端组件
│   ├── pages/                   # 页面组件
│   │   ├── Audit.tsx             # 审计管理页面
│   │   ├── CloudPlatforms.tsx    # 云平台管理页面
│   │   ├── Dashboard.tsx         # 仪表盘页面
│   │   ├── Deployments.tsx       # 发布管理页面
│   │   ├── Environments.tsx      # 环境管理页面
│   │   ├── HostGroups.tsx        # 主机分组管理页面
│   │   ├── Hosts.tsx             # 主机管理页面
│   │   ├── HostTags.tsx          # 主机标签管理页面
│   │   ├── Login.tsx             # 登录页面
│   │   ├── Permissions.tsx       # 权限管理页面
│   │   ├── Projects.tsx          # 项目管理页面
│   │   ├── Roles.tsx             # 角色管理页面
│   │   ├── Settings.tsx          # 系统设置页面
│   │   ├── SSHKeys.tsx           # SSH 密钥管理页面
│   │   ├── Users.tsx             # 用户管理页面
│   │   └── WebSSH.tsx            # WebSSH 终端页面
│   ├── services/                # API 服务封装
│   │   ├── audit.ts              # 审计日志 API
│   │   ├── auth.ts               # 认证 API
│   │   ├── cloudPlatform.ts      # 云平台 API
│   │   ├── deployment.ts         # 发布管理 API
│   │   ├── environment.ts        # 环境管理 API
│   │   ├── host.ts               # 主机管理 API
│   │   ├── hostGroup.ts          # 主机分组 API
│   │   ├── hostTag.ts            # 主机标签 API
│   │   ├── permission.ts         # 权限管理 API
│   │   ├── project.ts            # 项目管理 API
│   │   ├── role.ts               # 角色管理 API
│   │   ├── settings.ts           # 系统设置 API
│   │   ├── sshKey.ts             # SSH 密钥 API
│   │   ├── user.ts               # 用户管理 API
│   │   └── webssh.ts             # WebSSH WebSocket API
│   ├── contexts/                # React Context（状态管理）
│   │   └── AuthContext.tsx       # 认证上下文（用户信息、权限）
│   ├── config/                  # 配置文件
│   │   └── api.ts                # API 基础配置（baseURL、拦截器）
│   ├── App.tsx                  # 根组件（路由配置）
│   ├── App.css                  # 全局样式
│   ├── index.tsx                # 应用入口
│   ├── index.css                # 全局样式
│   └── setupTests.ts            # 测试配置
├── build/                       # 构建输出目录（生产环境）
├── Dockerfile                   # 前端 Docker 镜像构建
├── nginx.conf                   # Nginx 配置（生产环境）
├── package.json                 # 前端依赖配置
├── package-lock.json            # 依赖锁定文件
├── tsconfig.json                # TypeScript 配置
└── README.md                    # 前端说明文档
```

**前端架构说明**:
- **Pages**: 页面级组件，对应路由路径，包含完整的业务逻辑
- **Components**: 可复用组件，跨页面共享，通用 UI 组件
- **Services**: API 调用封装，统一请求处理，错误处理
- **Contexts**: 全局状态管理，认证信息，用户权限
- **Config**: 配置文件，API 基础 URL，环境变量

### OpenSpec Structure (`openspec/`)
```
openspec/
├── project.md                   # 项目上下文和约定
├── AGENTS.md                    # AI 助手使用指南
├── specs/                       # 当前规范（已实现的功能）
│   └── [capability]/            # 功能模块规范
│       ├── spec.md              # 需求规范
│       └── design.md            # 设计文档（可选）
└── changes/                     # 变更提案（待实现的功能）
    ├── [change-id]/             # 变更目录
    │   ├── proposal.md          # 变更提案（Why, What, Impact）
    │   ├── tasks.md             # 实施任务清单
    │   ├── design.md            # 技术设计文档（可选）
    │   └── specs/               # 规范变更（Delta）
    │       └── [capability]/
    │           └── spec.md      # ADDED/MODIFIED/REMOVED 需求
    ├── PROGRESS.md              # 项目进度总结
    ├── RECENT_UPDATES.md        # 最近更新记录
    └── archive/                 # 已完成的变更
        └── YYYY-MM-DD-[name]/   # 归档的变更
```

### Scripts and Tools (`scripts/`, `vms/`)
```
scripts/                         # Ruby 脚本目录（工具脚本统一使用 Ruby）
└── setup-salt-hosts.rb          # Salt 主机配置脚本

vms/                             # Vagrant 虚拟机配置
├── Vagrantfile                  # Vagrant 配置文件
├── create_salt_vms.rb           # 创建 Salt VM 脚本
├── README.md                    # VM 使用说明
└── scripts/
    ├── provision-master.sh      # Salt Master 初始化脚本
    └── provision-minion.sh      # Salt Minion 初始化脚本
```

### Build and Deployment
```
├── docker-compose.yml           # Docker Compose 配置（开发/生产环境）
├── Makefile                     # Make 构建脚本
├── Rakefile                     # Rake 任务管理（Ruby）
├── push.rb                      # Git 推送脚本（Ruby）
└── deploy.sh                    # 部署脚本（如存在）
```

## Module Organization

### Backend Modules (按功能划分)

1. **认证授权模块** (`auth`, `user`, `role`, `permission`)
   - 用户认证、JWT 令牌、角色权限管理
   - Handler: `auth_handler.go`, `user_handler.go`, `role_handler.go`, `permission_handler.go`
   - Service: `auth_service.go`, `user_service.go`, `role_service.go`, `permission_service.go`
   - Repository: `user_repository.go`, `role_repository.go`, `permission_repository.go`
   - Model: `user.go`

2. **主机管理模块** (`host`, `host_group`, `host_tag`)
   - 主机 CRUD、分组管理、标签管理
   - Handler: `host_handler.go`, `host_group_handler.go`, `host_tag_handler.go`
   - Service: `host_service.go`, `host_group_service.go`, `host_tag_service.go`
   - Repository: `host_repository.go`, `host_group_repository.go`, `host_tag_repository.go`
   - Model: `host.go`

3. **SSH 管理模块** (`ssh_key`, `webssh`)
   - SSH 密钥管理、WebSSH 终端
   - Handler: `ssh_key_handler.go`, `webssh_handler.go`
   - Service: `ssh_key_service.go`
   - Repository: `ssh_key_repository.go`
   - Model: `ssh_key.go`

4. **项目管理模块** (`project`)
   - 项目 CRUD、项目成员管理
   - Handler: `project_handler.go`
   - Service: `project_service.go`
   - Repository: `project_repository.go`
   - Model: `project.go`

5. **发布管理模块** (`deployment`, `environment`)
   - 发布流程管理、环境管理
   - Handler: `deployment_handler.go`, `environment_handler.go`
   - Service: `deployment_service.go`, `environment_service.go`
   - Repository: `deployment_repository.go`, `environment_repository.go`
   - Model: `deployment.go`, `environment.go`

6. **Salt 集成模块** (`salt`)
   - Salt API 客户端、Salt 操作封装
   - Handler: `salt_handler.go`
   - Service: (通过 Salt Manager)
   - Salt Client: `salt/client.go`, `salt/manager.go`

7. **系统管理模块** (`settings`, `cloud_platform`, `audit`)
   - 系统设置、云平台管理、审计日志
   - Handler: `settings_handler.go`, `cloud_platform_handler.go`, `audit_handler.go`
   - Service: `settings_service.go`, `cloud_platform_service.go`, `audit_service.go`
   - Repository: `settings_repository.go`, `cloud_platform_repository.go`, `audit_repository.go`
   - Model: `settings.go`, `cloud_platform.go`, `audit.go`

### Frontend Modules (按页面划分)

1. **认证模块** (`Login.tsx`, `AuthContext.tsx`)
   - 登录页面、认证状态管理

2. **仪表盘模块** (`Dashboard.tsx`)
   - 统计信息展示

3. **用户权限模块** (`Users.tsx`, `Roles.tsx`, `Permissions.tsx`)
   - 用户管理、角色管理、权限管理

4. **主机管理模块** (`Hosts.tsx`, `HostGroups.tsx`, `HostTags.tsx`)
   - 主机管理、分组管理、标签管理

5. **SSH 管理模块** (`SSHKeys.tsx`, `WebSSH.tsx`)
   - SSH 密钥管理、WebSSH 终端

6. **项目管理模块** (`Projects.tsx`)
   - 项目管理

7. **发布管理模块** (`Deployments.tsx`, `Environments.tsx`)
   - 发布管理、环境管理

8. **系统管理模块** (`Settings.tsx`, `CloudPlatforms.tsx`, `Audit.tsx`)
   - 系统设置、云平台管理、审计日志

## Development Tools

- **Makefile**: 构建、测试、运行命令
- **Rakefile**: Ruby 任务管理（推荐使用）
- **push.rb**: Git 推送脚本（Ruby）
- **scripts/**: Ruby 辅助脚本目录
- **vms/**: Vagrant 虚拟机配置（Salt 测试环境）
