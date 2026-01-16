# KkOps Backend

智能运维管理平台后端服务

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM (PostgreSQL ORM)
- Redis
- JWT Authentication
- WebSocket (gorilla/websocket)

## 开发环境设置

```bash
# 安装依赖
go mod download

# 运行开发服务器
go run cmd/server/main.go
```

## API 文档

项目使用 Swagger/OpenAPI 生成 API 文档。启动服务器后，可以通过以下地址访问：

- **Swagger UI**: http://localhost:8080/swagger/index.html

### 生成 API 文档

如果需要重新生成 API 文档（在修改了 handler 的注释后）：

```bash
# 安装 swag 工具（如果尚未安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
cd backend
swag init -g cmd/server/main.go -o api
```

## 项目结构

```
backend/
├── cmd/server/      # 应用入口
├── internal/        # 内部包
│   ├── config/      # 配置管理
│   ├── model/       # 数据模型
│   ├── repository/  # 数据访问层
│   ├── service/     # 业务逻辑层
│   ├── handler/     # HTTP处理器
│   ├── middleware/  # 中间件
│   └── utils/       # 工具函数
├── pkg/             # 公共包
├── api/             # API文档
├── migrations/      # 数据库迁移
└── configs/         # 配置文件
```
