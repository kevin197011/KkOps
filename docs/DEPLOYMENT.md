# KkOps 部署文档

## 目录

- [部署概述](#部署概述)
- [环境要求](#环境要求)
- [Docker Compose 部署](#docker-compose-部署)
- [生产环境部署](#生产环境部署)
- [配置说明](#配置说明)
- [数据备份与恢复](#数据备份与恢复)
- [监控与日志](#监控与日志)
- [故障排查](#故障排查)

## 部署概述

KkOps 支持多种部署方式：

- **开发环境**: Docker Compose（推荐）
- **生产环境**: Docker Compose 或 Kubernetes
- **云平台**: 支持部署到 AWS、阿里云、腾讯云等

## 环境要求

### 最低要求

- **CPU**: 2 核
- **内存**: 4GB
- **磁盘**: 20GB
- **操作系统**: Linux（推荐 Ubuntu 20.04+ 或 CentOS 8+）

### 软件要求

- Docker 20.10+
- Docker Compose 2.0+
- Git

## Docker Compose 部署

### 快速部署

```bash
# 克隆项目
git clone <repository-url>
cd KkOps

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，修改必要的配置

# 启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 服务说明

Docker Compose 会启动以下服务：

- **postgres**: PostgreSQL 数据库
- **redis**: Redis 缓存
- **backend**: 后端 API 服务
- **frontend**: 前端 Web 应用

### 服务端口

- 前端: `3000`
- 后端 API: `8080`
- PostgreSQL: `5432`
- Redis: `6379`

### 数据持久化

数据存储在 Docker volumes 中：

- `postgres_data`: PostgreSQL 数据
- `redis_data`: Redis 数据

### 停止服务

```bash
# 停止服务（保留数据）
docker-compose down

# 停止服务并删除数据（注意：会删除所有数据）
docker-compose down -v
```

## 生产环境部署

### 准备工作

1. **服务器准备**
   - 准备 Linux 服务器
   - 安装 Docker 和 Docker Compose
   - 配置防火墙规则

2. **域名和 SSL 证书**
   - 配置域名 DNS 解析
   - 申请 SSL 证书（推荐 Let's Encrypt）

3. **环境变量配置**
   - 创建 `.env` 文件
   - 设置强密码和密钥
   - 配置数据库连接

### 环境变量配置

创建 `.env` 文件并配置以下变量：

```bash
# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USER=kkops
DB_PASSWORD=<强密码>
DB_NAME=kkops

# Redis 配置
REDIS_HOST=redis
REDIS_PORT=6379

# JWT 密钥（必须修改为随机字符串）
JWT_SECRET=<随机生成的密钥，至少32字符>

# 加密密钥（必须修改为随机字符串，32字节）
ENCRYPTION_KEY=<随机生成的32字节密钥>

# 后端配置
BACKEND_HOST=0.0.0.0
BACKEND_PORT=8080

# 前端配置
FRONTEND_PORT=3000
```

**重要提示**: 生产环境必须修改 `JWT_SECRET` 和 `ENCRYPTION_KEY`，使用强随机字符串！

### 使用 Nginx 反向代理

1. 安装 Nginx

```bash
sudo apt-get update
sudo apt-get install nginx
```

2. 配置 Nginx

创建配置文件 `/etc/nginx/sites-available/kkops`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 前端
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 后端 API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 支持
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

3. 启用配置

```bash
sudo ln -s /etc/nginx/sites-available/kkops /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 使用 Docker Compose 部署

```bash
# 克隆项目
git clone <repository-url>
cd KkOps

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 使用 Kubernetes 部署

Kubernetes 部署配置需要单独准备，包括：

- Deployment 配置
- Service 配置
- ConfigMap 和 Secret
- Ingress 配置
- PersistentVolume 配置

## 配置说明

### 后端配置

配置文件位于 `backend/configs/config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # release 或 debug

database:
  host: "postgres"
  port: 5432
  user: "kkops"
  password: "kkops_dev"
  name: "kkops"
  sslmode: "disable"

redis:
  host: "redis"
  port: 6379

jwt:
  secret: "your-secret-key"
  expiration: 7200  # 秒

encryption:
  key: "your-32-byte-encryption-key"
```

### 前端配置

前端配置通过环境变量设置：

- `VITE_API_BASE_URL`: API 基础 URL（开发环境使用代理，生产环境配置实际域名）
- `VITE_WS_BASE_URL`: WebSocket 基础 URL

## 数据备份与恢复

### 数据库备份

```bash
# 备份 PostgreSQL 数据库
docker-compose exec postgres pg_dump -U kkops kkops > backup_$(date +%Y%m%d_%H%M%S).sql

# 或使用 Docker volume 备份
docker run --rm -v kkops_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data
```

### 数据库恢复

```bash
# 恢复 PostgreSQL 数据库
docker-compose exec -T postgres psql -U kkops kkops < backup_20240101_120000.sql

# 或使用 Docker volume 恢复
docker run --rm -v kkops_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /
```

### 定期备份

建议配置定时任务（cron）定期备份数据库：

```bash
# 编辑 crontab
crontab -e

# 每天凌晨 2 点备份
0 2 * * * cd /path/to/KkOps && docker-compose exec -T postgres pg_dump -U kkops kkops > /backup/kkops_$(date +\%Y\%m\%d).sql
```

## 监控与日志

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### 日志位置

- 后端日志: 标准输出（通过 Docker 日志查看）
- 前端日志: 浏览器控制台
- 数据库日志: Docker 日志

### 性能监控

建议在生产环境配置监控系统：

- **Prometheus + Grafana**: 指标监控
- **ELK Stack**: 日志收集和分析
- **APM 工具**: 应用性能监控

## 故障排查

### 服务无法启动

1. 检查 Docker 服务状态
2. 查看服务日志
3. 检查端口占用
4. 检查环境变量配置

### 数据库连接失败

1. 检查数据库服务是否运行
2. 检查数据库连接配置
3. 检查网络连接
4. 查看数据库日志

### 前端无法访问后端 API

1. 检查后端服务是否运行
2. 检查 API 地址配置
3. 检查网络连接
4. 查看浏览器控制台错误

### 性能问题

1. 检查服务器资源使用情况
2. 查看应用日志
3. 检查数据库查询性能
4. 考虑增加服务器资源

### 获取帮助

如果遇到问题：

1. 查看日志文件
2. 查看本文档的故障排查部分
3. 查看项目 Issue
4. 联系技术支持

## 安全建议

1. **使用强密码**: 数据库、Redis 等服务的密码必须足够强
2. **使用 HTTPS**: 生产环境必须启用 SSL/TLS
3. **定期更新**: 定期更新 Docker 镜像和依赖
4. **限制访问**: 使用防火墙限制不必要的端口访问
5. **备份数据**: 定期备份数据库和配置文件
6. **监控日志**: 定期检查日志，发现异常情况
7. **最小权限**: 使用最小权限原则配置用户和权限
