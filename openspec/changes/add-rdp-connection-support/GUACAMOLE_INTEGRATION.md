# Apache Guacamole 集成详细方案

**创建日期**: 2025-01-28  
**目标**: 使用 Apache Guacamole 作为 RDP 网关，实现浏览器中的 Windows 远程桌面访问

---

## 📋 概述

Apache Guacamole 是一个无客户端的远程桌面网关，支持 RDP、VNC、SSH 等多种协议。我们将使用 Guacamole 作为 RDP 协议网关，通过其 REST API 和 WebSocket 实现浏览器中的 Windows 远程桌面访问。

---

## 1. Guacamole 架构

### 1.1 Guacamole 组件

```
┌─────────────────────────────────────────┐
│      Guacamole Client (Browser)         │
│  - HTML5/JavaScript                     │
│  - WebSocket Client                      │
└──────────────┬──────────────────────────┘
               │ WebSocket (Guacamole Protocol)
               ↓
┌─────────────────────────────────────────┐
│      Guacamole Server (Java)             │
│  - guacd (daemon)                        │
│  - guacamole webapp (Tomcat)            │
│  - REST API                              │
└──────────────┬──────────────────────────┘
               │ RDP Protocol
               ↓
┌─────────────────────────────────────────┐
│      Windows Server (RDP Server)         │
└─────────────────────────────────────────┘
```

### 1.2 Guacamole 数据流

1. **连接建立**:
   - Browser → Guacamole WebSocket → guacd → RDP → Windows Server

2. **用户输入**:
   - Browser (鼠标/键盘) → Guacamole Protocol → guacd → RDP → Windows Server

3. **图形输出**:
   - Windows Server → RDP → guacd → Guacamole Protocol → Browser → Canvas

---

## 2. Guacamole 部署方案

### 2.1 部署方式选择

#### 方案 A: Docker 部署（推荐）✅

**优点**:
- ✅ 部署简单，一键启动
- ✅ 环境隔离，不影响现有系统
- ✅ 易于升级和维护
- ✅ 配置管理方便

**Docker Compose 配置**:
```yaml
services:
  guacamole:
    image: guacamole/guacamole:latest
    container_name: kkops-guacamole
    environment:
      GUACD_HOSTNAME: guacd
      GUACD_PORT: 4822
      MYSQL_HOSTNAME: guacamole-db
      MYSQL_DATABASE: guacamole_db
      MYSQL_USERNAME: guacamole
      MYSQL_PASSWORD: guacamole_password
    ports:
      - "8080:8080"
    depends_on:
      - guacd
      - guacamole-db
    networks:
      - kkops-network

  guacd:
    image: guacamole/guacd:latest
    container_name: kkops-guacd
    volumes:
      - guacd_drive:/drive
      - guacd_record:/record
    networks:
      - kkops-network

  guacamole-db:
    image: mysql:8.0
    container_name: kkops-guacamole-db
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: guacamole_db
      MYSQL_USER: guacamole
      MYSQL_PASSWORD: guacamole_password
    volumes:
      - guacamole_db_data:/var/lib/mysql
      - ./guacamole-init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - kkops-network

volumes:
  guacd_drive:
  guacd_record:
  guacamole_db_data:

networks:
  kkops-network:
    external: true
```

#### 方案 B: 传统部署

**步骤**:
1. 安装 Java 运行环境（JRE 11+）
2. 安装 Tomcat 9+
3. 安装 MySQL/PostgreSQL
4. 部署 guacd daemon
5. 部署 guacamole webapp
6. 配置数据库和认证

**工作量**: 比 Docker 部署多 2-3 天

### 2.2 数据库初始化

**MySQL 初始化脚本** (`guacamole-init.sql`):
```sql
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS guacamole_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE guacamole_db;

-- 执行 Guacamole 初始化 SQL
-- 从 Guacamole 安装包中获取: schema/*.sql
SOURCE /docker-entrypoint-initdb.d/schema/001-create-schema.sql;
SOURCE /docker-entrypoint-initdb.d/schema/002-create-admin-user.sql;
```

**PostgreSQL 初始化**:
```sql
CREATE DATABASE guacamole_db;
-- 然后执行 Guacamole PostgreSQL schema
```

---

## 3. Guacamole 配置

### 3.1 基础配置

**guacamole.properties**:
```properties
# MySQL 数据库配置
mysql-hostname: guacamole-db
mysql-port: 3306
mysql-database: guacamole_db
mysql-username: guacamole
mysql-password: guacamole_password

# guacd 配置
guacd-hostname: guacd
guacd-port: 4822

# 认证配置（使用数据库认证）
mysql-username: guacamole
mysql-password: guacamole_password
```

### 3.2 认证配置

**选项 1: 数据库认证**（推荐）
- 使用 Guacamole 内置的数据库认证
- 在数据库中创建用户
- 通过 REST API 管理用户

**选项 2: LDAP 认证**
- 集成 LDAP/Active Directory
- 统一用户管理

**选项 3: 自定义认证扩展**
- 开发自定义认证模块
- 与 KkOps 用户系统集成

**推荐**: 使用数据库认证，通过 REST API 同步 KkOps 用户到 Guacamole。

### 3.3 REST API 配置

**启用 REST API**:
- Guacamole 默认启用 REST API
- 需要配置认证（HTTP Basic Auth 或 Token）

**API 端点**:
- Base URL: `http://guacamole:8080/guacamole/api`
- 认证: HTTP Basic Auth 或 Token

---

## 4. KkOps 与 Guacamole 集成

### 4.1 集成架构

```
┌─────────────────────────────────────────┐
│         KkOps Frontend (React)           │
│  ┌──────────────────────────────────┐   │
│  │  RDP Client Component            │   │
│  │  (Guacamole JavaScript Client)   │   │
│  └──────────────┬───────────────────┘   │
└─────────────────┼────────────────────────┘
                  │ WebSocket
                  ↓
┌─────────────────────────────────────────┐
│      KkOps Backend (Go)                  │
│  ┌──────────────────────────────────┐   │
│  │  RDP Handler                     │   │
│  │  - Connection Management         │   │
│  │  - Session Management            │   │
│  └──────────────┬───────────────────┘   │
│                 │ REST API               │
└─────────────────┼────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│   Apache Guacamole Server (Java)        │
│  ┌──────────────────────────────────┐   │
│  │  REST API                        │   │
│  │  - Create Connection             │   │
│  │  - Get Session Token             │   │
│  └──────────────┬───────────────────┘   │
│                 │ RDP Protocol            │
└─────────────────┼────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│      Windows Server (RDP Server)         │
└─────────────────────────────────────────┘
```

### 4.2 数据流

#### 连接创建流程

1. **用户在前端创建 RDP 连接**
   ```
   Frontend → POST /api/v1/rdp/connections
   ```

2. **KkOps 后端处理**
   ```
   - 验证用户权限
   - 加密密码存储
   - 调用 Guacamole REST API 创建连接
   - 保存连接记录到数据库
   ```

3. **Guacamole 创建连接**
   ```
   POST /guacamole/api/session/data/mysql/connections
   {
     "name": "Windows Server 2019",
     "protocol": "rdp",
     "parameters": {
       "hostname": "192.168.1.100",
       "port": "3389",
       "username": "administrator",
       "password": "password",
       "domain": "DOMAIN",
       "width": "1024",
       "height": "768",
       "color-depth": "24",
       "security": "rdp"
     }
   }
   ```

#### 会话建立流程

1. **用户点击连接**
   ```
   Frontend → POST /api/v1/rdp/sessions
   ```

2. **KkOps 后端创建会话**
   ```
   - 创建会话记录
   - 调用 Guacamole API 获取会话令牌
   - 返回令牌和 WebSocket URL
   ```

3. **Guacamole 返回令牌**
   ```
   POST /guacamole/api/tokens
   {
     "username": "guacamole_user",
     "password": "guacamole_password",
     "dataSource": "mysql"
   }
   Response: {
     "authToken": "A1B2C3D4E5F6...",
     "username": "guacamole_user",
     "dataSource": "mysql"
   }
   ```

4. **前端建立 WebSocket 连接**
   ```javascript
   const token = response.authToken;
   const client = new Guacamole.Client(
     new Guacamole.HTTPTunnel(`/guacamole/tunnel?token=${token}`)
   );
   ```

---

## 5. 后端实现细节

### 5.1 Guacamole Service 实现

**文件**: `backend/internal/service/guacamole_service.go`

```go
package service

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type GuacamoleService interface {
    CreateConnection(conn *RDPConnection) (string, error)
    UpdateConnection(connID string, conn *RDPConnection) error
    DeleteConnection(connID string) error
    GetSessionToken(username, password string) (string, error)
}

type guacamoleService struct {
    baseURL    string
    username   string
    password   string
    httpClient *http.Client
}

func NewGuacamoleService(baseURL, username, password string) GuacamoleService {
    return &guacamoleService{
        baseURL:    baseURL,
        username:   username,
        password:   password,
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

// CreateConnection 在 Guacamole 中创建 RDP 连接
func (s *guacamoleService) CreateConnection(conn *RDPConnection) (string, error) {
    // 获取认证令牌
    token, err := s.getAuthToken()
    if err != nil {
        return "", err
    }

    // 构建连接参数
    params := map[string]string{
        "hostname":   conn.Hostname,
        "port":       fmt.Sprintf("%d", conn.Port),
        "username":   conn.Username,
        "password":   conn.Password, // 注意：Guacamole 需要明文密码
        "domain":     conn.Domain,
        "width":      fmt.Sprintf("%d", conn.Width),
        "height":     fmt.Sprintf("%d", conn.Height),
        "color-depth": fmt.Sprintf("%d", conn.ColorDepth),
        "security":    conn.SecurityMode,
    }

    if conn.IgnoreCert {
        params["ignore-cert"] = "true"
    }

    // 创建连接请求
    connectionReq := map[string]interface{}{
        "name":      conn.ConnectionName,
        "parentIdentifier": "ROOT",
        "protocol": "rdp",
        "parameters": params,
    }

    reqBody, _ := json.Marshal(connectionReq)
    req, _ := http.NewRequest("POST", 
        fmt.Sprintf("%s/api/session/data/mysql/connections?token=%s", s.baseURL, token),
        bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // 解析响应获取连接 ID
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    
    connID, _ := result["identifier"].(string)
    return connID, nil
}

// GetSessionToken 获取 Guacamole 会话令牌
func (s *guacamoleService) GetSessionToken(username, password string) (string, error) {
    // 这里可以使用 KkOps 用户映射到 Guacamole 用户
    // 或者使用统一的 Guacamole 服务账户
    
    data := url.Values{}
    data.Set("username", username)
    data.Set("password", password)
    data.Set("dataSource", "mysql")

    resp, err := http.PostForm(fmt.Sprintf("%s/api/tokens", s.baseURL), data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result map[string]string
    json.NewDecoder(resp.Body).Decode(&result)
    
    return result["authToken"], nil
}

// getAuthToken 获取 Guacamole API 认证令牌
func (s *guacamoleService) getAuthToken() (string, error) {
    data := url.Values{}
    data.Set("username", s.username)
    data.Set("password", s.password)
    data.Set("dataSource", "mysql")

    resp, err := http.PostForm(fmt.Sprintf("%s/api/tokens", s.baseURL), data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result map[string]string
    json.NewDecoder(resp.Body).Decode(&result)
    
    return result["authToken"], nil
}
```

### 5.2 RDP Connection Service

**文件**: `backend/internal/service/rdp_connection_service.go`

```go
package service

type RDPConnectionService interface {
    CreateConnection(userID uint64, conn *RDPConnection) (*RDPConnection, error)
    UpdateConnection(userID uint64, connID uint64, conn *RDPConnection) error
    DeleteConnection(userID uint64, connID uint64) error
    ListConnections(userID uint64, filters RDPConnectionFilters) ([]*RDPConnection, int64, error)
    GetConnection(userID uint64, connID uint64) (*RDPConnection, error)
}

type rdpConnectionService struct {
    connRepo    repository.RDPConnectionRepository
    guacamole   GuacamoleService
    encryptUtil utils.EncryptUtil
}

func NewRDPConnectionService(
    connRepo repository.RDPConnectionRepository,
    guacamole GuacamoleService,
    encryptUtil utils.EncryptUtil,
) RDPConnectionService {
    return &rdpConnectionService{
        connRepo:    connRepo,
        guacamole:   guacamole,
        encryptUtil: encryptUtil,
    }
}

func (s *rdpConnectionService) CreateConnection(userID uint64, conn *RDPConnection) (*RDPConnection, error) {
    // 1. 验证权限
    // 2. 加密密码
    encryptedPassword, err := s.encryptUtil.Encrypt(conn.Password)
    if err != nil {
        return nil, err
    }
    conn.PasswordEncrypted = encryptedPassword
    
    // 3. 调用 Guacamole API 创建连接
    guacamoleConnID, err := s.guacamole.CreateConnection(conn)
    if err != nil {
        return nil, fmt.Errorf("failed to create Guacamole connection: %w", err)
    }
    conn.GuacamoleConnectionID = guacamoleConnID
    
    // 4. 保存到数据库
    conn.UserID = userID
    if err := s.connRepo.Create(conn); err != nil {
        // 如果数据库保存失败，删除 Guacamole 连接
        s.guacamole.DeleteConnection(guacamoleConnID)
        return nil, err
    }
    
    return conn, nil
}
```

### 5.3 RDP Session Service

**文件**: `backend/internal/service/rdp_session_service.go`

```go
package service

type RDPSessionService interface {
    CreateSession(userID uint64, connID uint64) (*RDPSession, string, error)
    GetSessionToken(userID uint64, sessionID uint64) (string, error)
    EndSession(userID uint64, sessionID uint64) error
    ListSessions(userID uint64, filters RDPSessionFilters) ([]*RDPSession, int64, error)
}

type rdpSessionService struct {
    sessionRepo repository.RDPSessionRepository
    connRepo    repository.RDPConnectionRepository
    guacamole   GuacamoleService
}

func (s *rdpSessionService) CreateSession(userID uint64, connID uint64) (*RDPSession, error) {
    // 1. 验证连接存在且用户有权限
    conn, err := s.connRepo.GetByID(connID)
    if err != nil {
        return nil, err
    }
    if conn.UserID != userID {
        return nil, errors.New("unauthorized")
    }
    
    // 2. 创建会话记录
    session := &RDPSession{
        ConnectionID: connID,
        UserID:       userID,
        Status:       "active",
    }
    if err := s.sessionRepo.Create(session); err != nil {
        return nil, err
    }
    
    // 3. 获取 Guacamole 会话令牌
    // 注意：这里需要将 KkOps 用户映射到 Guacamole 用户
    // 或者使用统一的 Guacamole 服务账户
    token, err := s.guacamole.GetSessionToken("guacamole_user", "guacamole_password")
    if err != nil {
        return nil, err
    }
    session.GuacamoleSessionID = token
    
    // 4. 更新会话记录
    s.sessionRepo.Update(session)
    
    return session, nil
}
```

---

## 6. 前端实现细节

### 6.1 集成 Guacamole 客户端库

**安装**:
```bash
npm install guacamole-common-js
# 或使用 CDN
```

**CDN 方式**:
```html
<script src="https://cdn.jsdelivr.net/npm/guacamole-common-js@1.5.0/dist/guacamole.min.js"></script>
```

### 6.2 RDP Client 组件

**文件**: `frontend/src/components/RDPClient.tsx`

```typescript
import React, { useEffect, useRef, useState } from 'react';
import { Button, Space, message } from 'antd';
import { DisconnectOutlined } from '@ant-design/icons';
import { rdpService } from '../services/rdp';

// 声明 Guacamole 类型（如果没有类型定义）
declare const Guacamole: any;

interface RDPClientProps {
  connectionId: number;
  onClose?: () => void;
  onStatusChange?: (status: 'connecting' | 'connected' | 'disconnected') => void;
}

const RDPClient: React.FC<RDPClientProps> = ({ connectionId, onClose, onStatusChange }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const clientRef = useRef<any>(null);
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);

  useEffect(() => {
    connect();
    return () => {
      disconnect();
    };
  }, [connectionId]);

  const connect = async () => {
    if (connecting || connected) return;

    setConnecting(true);
    if (onStatusChange) onStatusChange('connecting');

    try {
      // 1. 创建会话并获取令牌
      const session = await rdpService.createSession(connectionId);
      const token = session.guacamoleSessionToken;

      // 2. 构建 Guacamole WebSocket URL
      const guacamoleURL = `/guacamole/tunnel?token=${token}`;
      // 注意：这里需要配置代理，将 /guacamole/* 代理到 Guacamole 服务器

      // 3. 创建 Guacamole 客户端
      const tunnel = new Guacamole.HTTPTunnel(guacamoleURL);
      const client = new Guacamole.Client(tunnel);

      // 4. 设置显示
      if (canvasRef.current) {
        const display = client.getDisplay();
        display.getDefaultLayer().getCanvas().style.width = '100%';
        display.getDefaultLayer().getCanvas().style.height = '100%';
        canvasRef.current.appendChild(display.getElement());
      }

      // 5. 处理连接事件
      client.onstatechange = (state: number) => {
        if (state === Guacamole.Client.CONNECTED) {
          setConnected(true);
          setConnecting(false);
          if (onStatusChange) onStatusChange('connected');
          message.success('RDP 连接已建立');
        } else if (state === Guacamole.Client.DISCONNECTED) {
          setConnected(false);
          setConnecting(false);
          if (onStatusChange) onStatusChange('disconnected');
          message.info('RDP 连接已断开');
        }
      };

      // 6. 处理错误
      client.onerror = (error: any) => {
        message.error(`RDP 连接错误: ${error.message}`);
        setConnecting(false);
        if (onStatusChange) onStatusChange('disconnected');
      };

      // 7. 连接
      client.connect();
      clientRef.current = client;

    } catch (error: any) {
      message.error(`连接失败: ${error.message}`);
      setConnecting(false);
      if (onStatusChange) onStatusChange('disconnected');
    }
  };

  const disconnect = () => {
    if (clientRef.current) {
      clientRef.current.disconnect();
      clientRef.current = null;
    }
    setConnected(false);
    if (onStatusChange) onStatusChange('disconnected');
  };

  return (
    <div style={{ width: '100%', height: '100%', position: 'relative' }}>
      <div style={{ padding: '8px', background: '#f0f0f0', borderBottom: '1px solid #d9d9d9' }}>
        <Space>
          <Button
            icon={<DisconnectOutlined />}
            onClick={disconnect}
            disabled={!connected}
          >
            断开连接
          </Button>
          {connecting && <span>正在连接...</span>}
          {connected && <span style={{ color: 'green' }}>已连接</span>}
        </Space>
      </div>
      <div style={{ width: '100%', height: 'calc(100% - 50px)', overflow: 'auto' }}>
        <canvas ref={canvasRef} style={{ width: '100%', height: '100%' }} />
      </div>
    </div>
  );
};

export default RDPClient;
```

### 6.3 Nginx 代理配置

**前端 Nginx 配置** (`frontend/nginx.conf`):
```nginx
server {
    listen 80;
    server_name localhost;

    # 前端静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    # API 代理到 KkOps 后端
    location /api {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Guacamole 代理
    location /guacamole {
        proxy_pass http://guacamole:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
        proxy_read_timeout 86400;
    }
}
```

---

## 7. 用户映射策略

### 7.1 方案 A: 统一服务账户（推荐）✅

**策略**: 使用统一的 Guacamole 服务账户

**优点**:
- ✅ 实现简单
- ✅ 无需同步用户
- ✅ 管理简单

**缺点**:
- ⚠️ 无法在 Guacamole 中区分用户（但可以在 KkOps 中区分）

**实现**:
```go
// 使用固定的 Guacamole 服务账户
const GUACAMOLE_SERVICE_USER = "kkops_service"
const GUACAMOLE_SERVICE_PASSWORD = "service_password"

func (s *guacamoleService) GetSessionToken() (string, error) {
    return s.getAuthToken(GUACAMOLE_SERVICE_USER, GUACAMOLE_SERVICE_PASSWORD)
}
```

### 7.2 方案 B: 用户同步

**策略**: 将 KkOps 用户同步到 Guacamole

**优点**:
- ✅ 可以在 Guacamole 中区分用户
- ✅ 支持 Guacamole 的用户管理功能

**缺点**:
- ⚠️ 需要维护用户同步
- ⚠️ 实现复杂

**实现**:
```go
// 在创建 RDP 连接时，同步用户到 Guacamole
func (s *guacamoleService) SyncUser(kkopsUserID uint64, kkopsUsername string) error {
    // 调用 Guacamole API 创建用户
    // 或使用 Guacamole 的用户管理 API
}
```

**推荐**: 使用方案 A（统一服务账户），简单且满足需求。

---

## 8. 配置管理

### 8.1 KkOps 配置

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

**环境变量** (`.env`):
```bash
# Guacamole 配置
GUACAMOLE_BASE_URL=http://guacamole:8080/guacamole
GUACAMOLE_USERNAME=guacamole
GUACAMOLE_PASSWORD=guacamole_password
```

### 8.2 Docker Compose 更新

**更新 `docker-compose.yml`**:
```yaml
services:
  # ... 现有服务

  # Guacamole 数据库
  guacamole-db:
    image: mysql:8.0
    container_name: kkops-guacamole-db
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: guacamole_db
      MYSQL_USER: guacamole
      MYSQL_PASSWORD: guacamole_password
    volumes:
      - guacamole_db_data:/var/lib/mysql
      - ./guacamole-init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - kkops-network

  # Guacamole daemon
  guacd:
    image: guacamole/guacd:latest
    container_name: kkops-guacd
    volumes:
      - guacd_drive:/drive
      - guacd_record:/record
    networks:
      - kkops-network

  # Guacamole webapp
  guacamole:
    image: guacamole/guacamole:latest
    container_name: kkops-guacamole
    environment:
      GUACD_HOSTNAME: guacd
      GUACD_PORT: 4822
      MYSQL_HOSTNAME: guacamole-db
      MYSQL_DATABASE: guacamole_db
      MYSQL_USERNAME: guacamole
      MYSQL_PASSWORD: guacamole_password
    ports:
      - "8080:8080"
    depends_on:
      guacd:
        condition: service_started
      guacamole-db:
        condition: service_healthy
    networks:
      - kkops-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/guacamole/"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  # ... 现有 volumes
  guacamole_db_data:
  guacd_drive:
  guacd_record:
```

---

## 9. 实施步骤

### 阶段 1: Guacamole 部署 (2-3 天)

1. **准备 Docker Compose 配置**
   - 添加 Guacamole 服务
   - 配置数据库
   - 配置网络

2. **初始化数据库**
   - 创建数据库
   - 执行 Guacamole schema SQL
   - 创建管理员用户

3. **启动和测试**
   - 启动 Guacamole 服务
   - 访问 Web 界面测试
   - 测试 RDP 连接（手动）

### 阶段 2: 后端集成 (4-5 天)

1. **实现 Guacamole Service**
   - REST API 客户端
   - 连接管理方法
   - 会话令牌获取

2. **实现 RDP Connection Service**
   - 连接 CRUD 操作
   - 权限检查
   - 密码加密

3. **实现 RDP Session Service**
   - 会话创建和管理
   - 令牌管理

4. **实现 Handler 和路由**
   - HTTP 处理器
   - 路由注册
   - 权限中间件

### 阶段 3: 前端集成 (3-4 天)

1. **集成 Guacamole 客户端库**
   - 安装或引入库
   - 配置代理

2. **实现 RDP Client 组件**
   - 连接建立
   - 图形显示
   - 事件处理

3. **实现连接管理界面**
   - 连接列表
   - 连接表单
   - 连接操作

### 阶段 4: 测试和优化 (2-3 天)

1. **功能测试**
   - 连接创建和删除
   - 会话建立和断开
   - 图形界面显示
   - 鼠标键盘输入

2. **性能测试**
   - 连接延迟
   - 图形传输性能
   - 并发连接测试

3. **安全测试**
   - 权限检查
   - 密码加密
   - 会话安全

---

## 10. 常见问题

### Q1: Guacamole 需要独立的 Java 服务吗？

**A**: 是的，Guacamole 需要 Java 运行环境。可以使用 Docker 部署，不影响现有 Go 服务。

### Q2: 如何与 KkOps 用户系统集成？

**A**: 推荐使用统一的 Guacamole 服务账户，在 KkOps 中管理用户权限。如果需要，也可以同步用户到 Guacamole。

### Q3: 性能如何？

**A**: Guacamole 性能良好，支持压缩和优化。对于局域网环境，延迟通常 < 100ms。

### Q4: 支持哪些 RDP 版本？

**A**: Guacamole 支持 RDP 5.0+，包括 Windows 10/11 和 Windows Server 2019+。

### Q5: 如何配置 Nginx 代理？

**A**: 需要将 `/guacamole/*` 路径代理到 Guacamole 服务器，支持 WebSocket 升级。

---

## 11. 参考资源

- **Apache Guacamole 官方文档**: https://guacamole.apache.org/doc/gug/
- **Guacamole REST API**: https://guacamole.apache.org/doc/gug/rest-api.html
- **Guacamole JavaScript API**: https://guacamole.apache.org/doc/gug/guacamole-common-js.html
- **Docker 部署指南**: https://guacamole.apache.org/doc/gug/guacamole-docker.html

---

**文档完成日期**: 2025-01-28  
**维护者**: KkOps开发团队

