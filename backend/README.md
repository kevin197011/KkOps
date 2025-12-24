# KkOps Backend

基于 Go 和 Gin 的运维中台管理系统后端。

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- PostgreSQL
- JWT 认证

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 3. 创建数据库

```bash
createdb kkops
```

或使用 SQL 脚本：

```bash
psql -U postgres -d kkops -f openspec/changes/add-ops-platform/database-schema.sql
```

### 4. 运行服务

```bash
go run cmd/api/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 认证

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/change-password` - 修改密码（需要认证）

### 用户管理

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户
- `POST /api/v1/users/:id/roles` - 分配角色
- `DELETE /api/v1/users/:id/roles` - 移除角色

## 项目结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # 应用入口
├── internal/
│   ├── config/              # 配置管理
│   ├── models/              # 数据模型
│   ├── repository/          # 数据访问层
│   ├── service/             # 业务逻辑层
│   ├── handler/             # HTTP 处理器
│   ├── middleware/          # 中间件
│   └── utils/               # 工具函数
└── migrations/              # 数据库迁移脚本
```

