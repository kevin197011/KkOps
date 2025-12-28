# Guacamole RDP 网关快速开始

**创建日期**: 2025-01-28  
**目标**: 快速部署和测试 Apache Guacamole RDP 网关

---

## 🚀 5 分钟快速开始

### 步骤 1: 添加 Docker Compose 配置

将 `docker-compose.guacamole.yml` 的内容合并到主 `docker-compose.yml`，或使用：

```bash
# 方式 1: 合并到主文件
cat docker-compose.guacamole.yml >> docker-compose.yml

# 方式 2: 使用多个 Compose 文件
docker-compose -f docker-compose.yml -f docker-compose.guacamole.yml up -d
```

### 步骤 2: 启动服务

```bash
# 启动 Guacamole 服务
docker-compose up -d guacamole-db guacd guacamole

# 查看日志
docker-compose logs -f guacamole
```

### 步骤 3: 初始化数据库

```bash
# 下载 Guacamole schema SQL
mkdir -p guacamole-schema
cd guacamole-schema

# 下载并解压
wget https://downloads.apache.org/guacamole/1.5.0/binary/guacamole-auth-jdbc-1.5.0.tar.gz
tar -xzf guacamole-auth-jdbc-1.5.0.tar.gz

# 执行初始化 SQL
docker-compose exec -T guacamole-db mysql -u root -pguacamole_root guacamole_db < \
  guacamole-auth-jdbc-1.5.0/mysql/schema/001-create-schema.sql
docker-compose exec -T guacamole-db mysql -u root -pguacamole_root guacamole_db < \
  guacamole-auth-jdbc-1.5.0/mysql/schema/002-create-admin-user.sql
```

### 步骤 4: 访问 Web 界面

打开浏览器访问: `http://localhost:8080/guacamole`

**默认登录**:
- 用户名: `guacadmin`
- 密码: `guacadmin`

### 步骤 5: 测试 RDP 连接

1. 登录 Guacamole
2. 点击 "New Connection"
3. 配置连接：
   - **Name**: Windows Server
   - **Protocol**: RDP
   - **Network**: `192.168.1.100:3389` (你的 Windows 服务器)
   - **Username**: `administrator`
   - **Password**: `your_password`
4. 点击 "Save" 并测试连接

---

## 📝 环境变量配置

创建或更新 `.env` 文件：

```bash
# Guacamole 数据库配置
GUACAMOLE_DB_ROOT_PASSWORD=guacamole_root_password
GUACAMOLE_DB_PASSWORD=guacamole_password

# Guacamole API 配置（用于 KkOps 集成）
GUACAMOLE_BASE_URL=http://guacamole:8080/guacamole
GUACAMOLE_USERNAME=kkops_service
GUACAMOLE_PASSWORD=service_password
```

---

## 🔧 配置 Nginx 代理

更新 `frontend/nginx.conf`，添加 Guacamole 代理：

```nginx
location /guacamole {
    proxy_pass http://guacamole:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_buffering off;
    proxy_read_timeout 86400;
}
```

重建前端镜像：

```bash
docker-compose build frontend
docker-compose up -d frontend
```

---

## ✅ 验证安装

### 检查服务状态

```bash
# 检查容器状态
docker-compose ps | grep guacamole

# 应该看到：
# kkops-guacamole-db    Up (healthy)
# kkops-guacd          Up
# kkops-guacamole      Up (healthy)
```

### 测试 API

```bash
# 获取认证令牌
curl -X POST http://localhost:8080/guacamole/api/tokens \
  -d "username=guacadmin" \
  -d "password=guacadmin" \
  -d "dataSource=mysql"

# 应该返回 JSON，包含 authToken
```

---

## 🐛 常见问题

### Q: Guacamole 无法启动

**A**: 检查数据库是否就绪：
```bash
docker-compose logs guacamole-db
docker-compose exec guacamole ping guacamole-db
```

### Q: 无法访问 Web 界面

**A**: 检查端口是否被占用：
```bash
netstat -an | grep 8080
# 或
lsof -i :8080
```

### Q: RDP 连接失败

**A**: 检查 Windows 服务器：
- 是否启用远程桌面
- 防火墙是否允许 3389 端口
- 用户名密码是否正确

---

## 📚 下一步

- 查看 **GUACAMOLE_INTEGRATION.md** 了解详细集成方案
- 查看 **DEPLOYMENT_GUIDE.md** 了解完整部署指南
- 查看 **design.md** 了解架构设计

---

**快速开始完成！** 🎉

