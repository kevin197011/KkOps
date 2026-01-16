# KkOps - 智能运维管理平台

智能运维管理平台，提供资产管理、运维执行、WebSSH 终端等功能。

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM (PostgreSQL ORM)
- Redis
- JWT Authentication
- WebSocket (gorilla/websocket)

### 前端
- React 18+
- TypeScript
- Ant Design 5+
- Vite
- React Router v6
- Axios
- xterm.js (WebSSH)

## 快速开始

### 前置要求

- Docker & Docker Compose
- Git

### 启动服务

```bash
# 启动所有服务（数据库、Redis、后端、前端）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷（注意：会删除数据库数据）
docker-compose down -v
```

### 服务端口

- **前端**: http://localhost:3000
- **后端 API**: http://localhost:8080
- **API 文档 (Swagger)**: http://localhost:8080/swagger/index.html
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### 开发环境配置

#### 后端开发

```bash
cd backend

# 安装依赖
go mod download

# 运行开发服务器（需要先启动 postgres 和 redis）
docker-compose up -d postgres redis
go run cmd/server/main.go
```

#### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

### 环境变量

后端支持以下环境变量（可通过 docker-compose.yml 或 .env 文件配置）：

- `DB_HOST`: 数据库主机（默认: postgres）
- `DB_PORT`: 数据库端口（默认: 5432）
- `DB_USER`: 数据库用户名（默认: kkops）
- `DB_PASSWORD`: 数据库密码（默认: kkops_dev）
- `DB_NAME`: 数据库名称（默认: kkops）
- `REDIS_HOST`: Redis 主机（默认: redis）
- `REDIS_PORT`: Redis 端口（默认: 6379）
- `JWT_SECRET`: JWT 密钥（生产环境必须修改）
- `ENCRYPTION_KEY`: 加密密钥（生产环境必须修改）

### 配置文件

后端配置文件位于 `backend/configs/config.example.yaml`，复制并修改为 `config.yaml`：

```bash
cp backend/configs/config.example.yaml backend/configs/config.yaml
```

## 项目结构

```
KkOps/
├── backend/              # 后端服务
│   ├── cmd/server/       # 应用入口
│   ├── internal/         # 内部包
│   │   ├── config/       # 配置管理
│   │   ├── model/        # 数据模型
│   │   ├── service/      # 业务逻辑层
│   │   ├── handler/      # HTTP处理器
│   │   ├── middleware/   # 中间件
│   │   └── utils/        # 工具函数
│   └── configs/          # 配置文件
├── frontend/             # 前端应用
│   ├── src/
│   │   ├── api/          # API接口定义
│   │   ├── components/   # 公共组件
│   │   ├── layouts/      # 布局组件
│   │   ├── pages/        # 页面组件
│   │   ├── stores/       # 状态管理
│   │   ├── hooks/        # 自定义Hooks
│   │   └── themes/       # 主题配置
│   └── public/           # 静态资源
└── docker-compose.yml    # Docker Compose 配置
```

## 功能特性

- ✅ 用户认证和授权（JWT + RBAC）
- ✅ API Token 管理
- ✅ 资产管理（项目、分类、标签、资产 CRUD）
- ✅ 运维执行（任务模板、任务管理、实时日志）
- ✅ WebSSH（SSH 密钥管理、终端界面）
- ✅ 主题切换（亮色/暗色）
- ✅ 响应式布局

## 默认账号

系统首次启动时会自动创建默认管理员账号：

- **用户名**: `admin`
- **密码**: `admin123`
- **邮箱**: `admin@kkops.local`

**⚠️ 重要提示**: 首次登录后请立即修改默认密码！

## 文档

- [用户使用指南](docs/USER_GUIDE.md) - 用户操作手册
- [开发者指南](docs/DEVELOPER_GUIDE.md) - 开发环境设置和开发规范
- [部署文档](docs/DEPLOYMENT.md) - 生产环境部署指南

## 开发规范

请参考项目根目录下的规范文档。

## 许可证

MIT License
