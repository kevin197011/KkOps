# Guacamole 部署指南

**创建日期**: 2025-01-28  
**目标**: 在 KkOps 项目中部署 Apache Guacamole 作为 RDP 网关

---

## 📋 快速开始

### 前置要求

- Docker 和 Docker Compose
- 至少 2GB 可用内存（Guacamole 需要）
- 网络访问（用于下载镜像）

---

## 1. Docker Compose 配置

### 1.1 更新 docker-compose.yml

在现有的 `docker-compose.yml` 中添加以下服务：

```yaml
services:
  # ... 现有服务 (postgres, redis, backend, frontend)

  # Guacamole 数据库
  guacamole-db:
    image: mysql:8.0
    container_name: kkops-guacamole-db
    environment:
      MYSQL_ROOT_PASSWORD: ${GUACAMOLE_DB_ROOT_PASSWORD:-guacamole_root}
      MYSQL_DATABASE: guacamole_db
      MYSQL_USER: guacamole
      MYSQL_PASSWORD: ${GUACAMOLE_DB_PASSWORD:-guacamole_password}
      TZ: Asia/Shanghai
    ports:
      - "3307:3306"  # 避免与现有 MySQL 冲突（如果有）
    volumes:
      - guacamole_db_data:/var/lib/mysql
      - ./guacamole-init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${MYSQL_ROOT_PASSWORD}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - kkops-network

  # Guacamole daemon (协议网关)
  guacd:
    image: guacamole/guacd:latest
    container_name: kkops-guacd
    volumes:
      - guacd_drive:/drive
      - guacd_record:/record
    networks:
      - kkops-network
    restart: unless-stopped

  # Guacamole webapp
  guacamole:
    image: guacamole/guacamole:latest
    container_name: kkops-guacamole
    environment:
      GUACD_HOSTNAME: guacd
      GUACD_PORT: 4822
      MYSQL_HOSTNAME: guacamole-db
      MYSQL_PORT: 3306
      MYSQL_DATABASE: guacamole_db
      MYSQL_USERNAME: guacamole
      MYSQL_PASSWORD: ${GUACAMOLE_DB_PASSWORD:-guacamole_password}
      MYSQL_AUTO_CREATE_ACCOUNTS: "true"
      TZ: Asia/Shanghai
    ports:
      - "8080:8080"
    depends_on:
      guacd:
        condition: service_started
      guacamole-db:
        condition: service_healthy
    networks:
      - kkops-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/guacamole/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s

volumes:
  # ... 现有 volumes
  guacamole_db_data:
    driver: local
  guacd_drive:
    driver: local
  guacd_record:
    driver: local

networks:
  kkops-network:
    external: true  # 使用现有网络，或改为 driver: bridge
```

### 1.2 创建数据库初始化脚本

**文件**: `guacamole-init.sql`

```sql
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS guacamole_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE guacamole_db;

-- 注意：Guacamole 的 schema SQL 文件需要从 Guacamole 镜像中获取
-- 或者手动下载：https://guacamole.apache.org/releases/
-- 这里只创建数据库，schema 会在容器启动时自动执行（如果配置了）
```

**或者手动下载 schema 文件**:

```bash
# 下载 Guacamole schema 文件
wget https://downloads.apache.org/guacamole/1.5.0/binary/guacamole-auth-jdbc-1.5.0.tar.gz
tar -xzf guacamole-auth-jdbc-1.5.0.tar.gz
cp guacamole-auth-jdbc-1.5.0/mysql/schema/*.sql ./guacamole-schema/
```

---

## 2. 启动 Guacamole

### 2.1 启动服务

```bash
# 启动所有服务（包括 Guacamole）
docker-compose up -d

# 或只启动 Guacamole 相关服务
docker-compose up -d guacamole-db guacd guacamole
```

### 2.2 检查服务状态

```bash
# 检查容器状态
docker-compose ps

# 查看 Guacamole 日志
docker-compose logs guacamole

# 查看数据库日志
docker-compose logs guacamole-db
```

### 2.3 访问 Guacamole Web 界面

打开浏览器访问: `http://localhost:8080/guacamole`

**默认登录**:
- 用户名: `guacadmin`
- 密码: `guacadmin`（首次登录后需要修改）

---

## 3. 配置 Guacamole

### 3.1 创建管理员用户（如果使用数据库认证）

**通过 Web 界面**:
1. 登录 Guacamole
2. 进入 "Settings" → "Users"
3. 创建新用户或修改默认用户

**通过 SQL**:
```sql
USE guacamole_db;

-- 创建用户（密码需要 MD5 哈希）
INSERT INTO guacamole_user (username, password_hash, password_salt)
VALUES ('kkops_service', MD5('password'), NULL);

-- 获取用户 ID
SET @user_id = LAST_INSERT_ID();

-- 授予系统权限
INSERT INTO guacamole_system_permission (user_id, permission)
VALUES (@user_id, 'ADMINISTER');
```

### 3.2 配置 REST API 访问

Guacamole REST API 使用 HTTP Basic Auth 或 Token 认证。

**获取 Token**:
```bash
curl -X POST http://localhost:8080/guacamole/api/tokens \
  -d "username=guacadmin" \
  -d "password=guacadmin" \
  -d "dataSource=mysql"
```

**响应**:
```json
{
  "authToken": "A1B2C3D4E5F6...",
  "username": "guacadmin",
  "dataSource": "mysql"
}
```

---

## 4. KkOps 后端配置

### 4.1 环境变量配置

**文件**: `.env`

```bash
# Guacamole 配置
GUACAMOLE_BASE_URL=http://guacamole:8080/guacamole
GUACAMOLE_USERNAME=kkops_service
GUACAMOLE_PASSWORD=your_password_here
GUACAMOLE_API_PATH=/api
```

### 4.2 后端配置更新

**文件**: `backend/internal/config/config.go`

```go
type Config struct {
    // ... 现有配置
    
    Guacamole struct {
        BaseURL  string `env:"GUACAMOLE_BASE_URL" envDefault:"http://guacamole:8080/guacamole"`
        Username string `env:"GUACAMOLE_USERNAME" envDefault:"guacamole"`
        Password string `env:"GUACAMOLE_PASSWORD" envDefault:"guacamole"`
        APIPath  string `env:"GUACAMOLE_API_PATH" envDefault:"/api"`
    }
}
```

---

## 5. 前端 Nginx 代理配置

### 5.1 更新 Nginx 配置

**文件**: `frontend/nginx.conf`

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # Docker 内部 DNS resolver
    resolver 127.0.0.11 valid=30s;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;

    # 前端路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理到后端
    location /api {
        set $backend "http://backend:8080";
        proxy_pass $backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Guacamole 代理（重要！）
    location /guacamole {
        set $guacamole "http://guacamole:8080";
        proxy_pass $guacamole;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_buffering off;
        proxy_read_timeout 86400;  # 24小时（RDP会话可能很长）
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### 5.2 重建前端镜像

```bash
docker-compose build frontend
docker-compose up -d frontend
```

---

## 6. 测试连接

### 6.1 通过 Guacamole Web 界面测试

1. 访问 `http://localhost:8080/guacamole`
2. 登录 Guacamole
3. 创建新的 RDP 连接：
   - Protocol: RDP
   - Network: 192.168.1.100:3389
   - Username: administrator
   - Password: your_password
4. 测试连接

### 6.2 通过 KkOps API 测试

```bash
# 1. 创建 RDP 连接
curl -X POST http://localhost:8080/api/v1/rdp/connections \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "host_id": 1,
    "connection_name": "Windows Server 2019",
    "hostname": "192.168.1.100",
    "port": 3389,
    "username": "administrator",
    "password": "password",
    "domain": "",
    "width": 1024,
    "height": 768,
    "color_depth": 24,
    "security_mode": "rdp"
  }'

# 2. 创建会话
curl -X POST http://localhost:8080/api/v1/rdp/sessions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "connection_id": 1
  }'
```

---

## 7. 故障排查

### 7.1 常见问题

**问题 1: Guacamole 无法连接数据库**

**解决**:
```bash
# 检查数据库容器状态
docker-compose ps guacamole-db

# 查看数据库日志
docker-compose logs guacamole-db

# 检查网络连接
docker-compose exec guacamole ping guacamole-db
```

**问题 2: RDP 连接失败**

**检查**:
- Windows 服务器是否启用远程桌面
- 防火墙是否允许 3389 端口
- 用户名密码是否正确
- 网络是否可达

**问题 3: 前端无法连接 Guacamole**

**检查**:
- Nginx 代理配置是否正确
- Guacamole 服务是否运行
- 网络是否连通

### 7.2 日志查看

```bash
# Guacamole 日志
docker-compose logs -f guacamole

# guacd 日志
docker-compose logs -f guacd

# 数据库日志
docker-compose logs -f guacamole-db
```

---

## 8. 生产环境配置

### 8.1 安全配置

1. **修改默认密码**
   ```bash
   # 修改 Guacamole 默认管理员密码
   # 通过 Web 界面或 SQL 更新
   ```

2. **使用强密码**
   ```bash
   # 在 .env 中使用强密码
   GUACAMOLE_DB_PASSWORD=strong_password_here
   GUACAMOLE_PASSWORD=strong_password_here
   ```

3. **限制网络访问**
   ```yaml
   # 在 docker-compose.yml 中限制端口暴露
   # 仅允许内部网络访问
   ports:
     - "127.0.0.1:8080:8080"  # 仅本地访问
   ```

### 8.2 性能优化

1. **数据库优化**
   ```sql
   -- 添加索引
   CREATE INDEX idx_connection_name ON guacamole_connection(connection_name);
   ```

2. **连接池配置**
   - 调整 Guacamole 连接池大小
   - 配置超时时间

3. **资源限制**
   ```yaml
   # 在 docker-compose.yml 中添加资源限制
   guacamole:
     deploy:
       resources:
         limits:
           memory: 2G
           cpus: '1.0'
   ```

---

## 9. 维护和升级

### 9.1 备份数据库

```bash
# 备份 Guacamole 数据库
docker-compose exec guacamole-db mysqldump -u root -p guacamole_db > guacamole_backup.sql

# 恢复数据库
docker-compose exec -T guacamole-db mysql -u root -p guacamole_db < guacamole_backup.sql
```

### 9.2 升级 Guacamole

```bash
# 1. 备份数据库
# 2. 停止服务
docker-compose stop guacamole guacd

# 3. 拉取新镜像
docker-compose pull guacamole guacd

# 4. 启动服务
docker-compose up -d guacamole guacd

# 5. 检查日志
docker-compose logs -f guacamole
```

---

## 10. 监控和告警

### 10.1 健康检查

```bash
# 检查 Guacamole 健康状态
curl http://localhost:8080/guacamole/api/tokens

# 检查数据库连接
docker-compose exec guacamole-db mysqladmin ping -h localhost -u root -p
```

### 10.2 监控指标

- Guacamole 服务状态
- 数据库连接数
- 活跃 RDP 会话数
- 资源使用情况（CPU、内存）

---

**文档完成日期**: 2025-01-28  
**维护者**: KkOps开发团队

