# KkOps 开发者指南

## 目录

- [项目结构](#项目结构)
- [开发环境设置](#开发环境设置)
- [后端开发](#后端开发)
- [前端开发](#前端开发)
- [代码规范](#代码规范)
- [API 设计](#api-设计)
- [数据库设计](#数据库设计)
- [部署流程](#部署流程)

## 项目结构

```
KkOps/
├── backend/              # 后端服务
│   ├── cmd/server/       # 应用入口
│   ├── internal/         # 内部包（不对外暴露）
│   │   ├── config/       # 配置管理
│   │   ├── model/        # 数据模型（GORM）
│   │   ├── service/      # 业务逻辑层
│   │   ├── handler/      # HTTP 处理器（Gin）
│   │   ├── middleware/   # 中间件
│   │   └── utils/        # 工具函数
│   ├── pkg/              # 公共包（可对外暴露）
│   ├── api/              # API 文档（未来）
│   ├── migrations/       # 数据库迁移脚本
│   └── configs/          # 配置文件
├── frontend/             # 前端应用
│   ├── src/
│   │   ├── api/          # API 接口定义
│   │   ├── components/   # 公共组件
│   │   ├── layouts/      # 布局组件
│   │   ├── pages/        # 页面组件
│   │   ├── stores/       # 状态管理（Zustand）
│   │   ├── hooks/        # 自定义 Hooks
│   │   ├── themes/       # 主题配置
│   │   └── utils/        # 工具函数
│   └── public/           # 静态资源
├── docs/                 # 文档目录
└── docker-compose.yml    # Docker Compose 配置
```

## 开发环境设置

### 前置要求

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose
- Git

### 初始化项目

```bash
# 克隆项目
git clone <repository-url>
cd KkOps

# 启动基础服务（PostgreSQL、Redis）
docker-compose up -d postgres redis

# 后端开发环境
cd backend
go mod download

# 前端开发环境
cd frontend
npm install
```

## 后端开发

### 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT
- **日志**: Zap
- **配置**: Viper

### 代码结构

#### Handler 层

负责处理 HTTP 请求和响应：

```go
// internal/handler/user/handler.go
type Handler struct {
    service *userService.Service
}

func (h *Handler) ListUsers(c *gin.Context) {
    // 处理请求
    // 调用 service 层
    // 返回响应
}
```

#### Service 层

包含业务逻辑：

```go
// internal/service/user/service.go
type Service struct {
    db *gorm.DB
}

func (s *Service) ListUsers(page, pageSize int) ([]User, int64, error) {
    // 业务逻辑
    // 数据库操作
}
```

#### Model 层

数据模型定义（使用 GORM tags）：

```go
// internal/model/user.go
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"uniqueIndex;not null" json:"username"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### 添加新功能

1. **定义数据模型**（`internal/model/`）
2. **创建 Service**（`internal/service/`）
3. **创建 Handler**（`internal/handler/`）
4. **注册路由**（`cmd/server/main.go`）
5. **添加数据库迁移**（`migrations/`）

### 数据库迁移

目前使用 GORM 的 AutoMigrate 功能自动创建表结构。未来可以添加迁移脚本到 `migrations/` 目录。

### 中间件

- **认证中间件**（`middleware/auth.go`）: JWT 和 API Token 验证
- **权限中间件**（`middleware/permission.go`）: 权限检查
- **日志中间件**（`middleware/logger.go`）: 请求日志
- **错误处理中间件**（`middleware/error.go`）: 统一错误处理
- **CORS 中间件**（`middleware/cors.go`）: 跨域支持

### 运行后端

```bash
cd backend
go run cmd/server/main.go
```

后端服务运行在 `http://localhost:8080`

## 前端开发

### 技术栈

- **框架**: React 18+
- **语言**: TypeScript
- **UI 库**: Ant Design 5+
- **状态管理**: Zustand
- **路由**: React Router v6
- **构建工具**: Vite
- **HTTP 客户端**: Axios

### 代码结构

#### API 层

API 接口定义（`src/api/`）：

```typescript
// src/api/user.ts
export const userApi = {
  list: () => apiClient.get<User[]>('/users'),
  get: (id: number) => apiClient.get<User>(`/users/${id}`),
  create: (data: CreateUserRequest) => apiClient.post<User>('/users', data),
  // ...
}
```

#### 页面组件

页面组件（`src/pages/`）：

```typescript
// src/pages/users/UserList.tsx
const UserList = () => {
  const [users, setUsers] = useState<User[]>([])
  // ...
  return (
    <div>
      {/* UI */}
    </div>
  )
}
```

#### 状态管理

使用 Zustand 管理全局状态（`src/stores/`）：

```typescript
// src/stores/auth.ts
interface AuthState {
  token: string | null
  user: UserInfo | null
  setAuth: (token: string, user: UserInfo) => void
  logout: () => void
}
```

### 添加新页面

1. **创建页面组件**（`src/pages/`）
2. **定义 API 接口**（`src/api/`）
3. **添加路由**（`src/App.tsx`）
4. **更新导航菜单**（`src/layouts/MainLayout.tsx`）

### 运行前端

```bash
cd frontend
npm run dev
```

前端服务运行在 `http://localhost:3000`

## 代码规范

### 后端代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码
- 函数和变量命名使用驼峰命名法
- 常量使用大写字母和下划线
- 包名使用小写字母

### 前端代码规范

- 使用 TypeScript 严格模式
- 使用 ESLint 和 Prettier 格式化代码
- 组件使用函数式组件和 Hooks
- 使用 TypeScript 接口定义类型
- 组件文件使用 PascalCase 命名
- 工具函数使用 camelCase 命名

### 提交规范

遵循 Conventional Commits 规范：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档变更
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构
- `test`: 增加测试
- `chore`: 构建过程或辅助工具的变动

## API 设计

### RESTful API 规范

- 使用 HTTP 动词（GET、POST、PUT、DELETE）
- 使用复数名词作为资源路径
- 使用状态码表示请求结果
- 统一响应格式

### 响应格式

成功响应：

```json
{
  "data": { ... },
  "message": "success"
}
```

错误响应：

```json
{
  "error": "错误信息",
  "code": "ERROR_CODE"
}
```

### 认证

- Web 应用使用 JWT Token
- API 调用使用 API Token
- Token 通过 `Authorization: Bearer <token>` 头传递

## 数据库设计

### 核心表结构

- **users**: 用户表
- **roles**: 角色表
- **permissions**: 权限表
- **user_roles**: 用户角色关联表
- **role_permissions**: 角色权限关联表
- **api_tokens**: API Token 表
- **projects**: 项目表
- **assets**: 资产表
- **asset_categories**: 资产分类表
- **tags**: 标签表
- **ssh_keys**: SSH 密钥表
- **task_templates**: 任务模板表
- **tasks**: 任务表
- **task_executions**: 任务执行表

详细的数据库设计请参考 `r.md` 文档。

## 部署流程

请参考 [部署文档](DEPLOYMENT.md) 了解详细的部署流程。

## 测试

### 运行测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm test
```

### 测试覆盖

- 单元测试：测试单个函数/组件
- 集成测试：测试模块间交互
- E2E 测试：测试完整用户流程

## 调试

### 后端调试

- 使用日志（Zap）输出调试信息
- 使用 Go 调试器（Delve）
- 查看数据库日志

### 前端调试

- 使用浏览器开发者工具
- 使用 React DevTools
- 查看网络请求和响应

## 获取帮助

- 查看项目文档
- 查看代码注释
- 提交 Issue
- 联系项目维护者
