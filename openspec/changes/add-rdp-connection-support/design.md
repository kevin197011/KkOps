# Design: RDP Connection Support

## Context

当前系统仅支持 SSH 协议连接，需要添加 RDP（Remote Desktop Protocol）支持以提供 Windows 图形界面远程访问能力。RDP 是 Windows 的远程桌面协议，比 SSH 复杂，需要专门的客户端实现。

## Goals / Non-Goals

### Goals
- 支持 Windows RDP 连接
- 在浏览器中显示 Windows 桌面
- 支持鼠标和键盘操作
- 支持连接管理和会话跟踪
- 与现有 SSH 功能统一管理

### Non-Goals
- 不支持其他远程桌面协议（VNC、SPICE 等）- 可后续扩展
- 不支持音频重定向（可后续添加）
- 不支持打印机重定向（可后续添加）
- 不支持文件传输（可后续添加）

## Decisions

### Decision: 使用 Apache Guacamole 作为 RDP 网关
**What**: 集成 Apache Guacamole 服务器处理 RDP 协议
**Why**: 
- 成熟稳定，生产环境广泛使用
- 完整的 RDP 协议支持
- 提供 HTML5 Web 客户端
- 支持多种协议（未来可扩展）

**Alternatives considered**:
- noVNC + xrdp: 需要 Linux 代理服务器，配置复杂
- FreeRDP WebAssembly: 不成熟，浏览器兼容性问题
- 第三方 SaaS: 依赖外部服务，安全和隐私问题

### Decision: Guacamole 独立部署
**What**: Guacamole 作为独立服务部署，通过 API 集成
**Why**:
- 保持系统解耦
- 便于独立升级和维护
- 可以复用 Guacamole 的认证系统（可选）

**Alternatives considered**:
- 嵌入式集成: 增加系统复杂度，耦合度高

### Decision: 统一连接管理界面
**What**: SSH 和 RDP 连接在同一个界面管理
**Why**:
- 用户体验统一
- 便于管理
- 减少界面复杂度

## Architecture

### System Architecture

```
┌─────────────────────────────────────────┐
│           Browser (React)               │
│  ┌──────────────┐  ┌──────────────┐   │
│  │  SSH Client  │  │  RDP Client   │   │
│  │  (xterm.js)  │  │ (Guacamole)  │   │
│  └──────┬───────┘  └──────┬───────┘   │
└─────────┼──────────────────┼────────────┘
          │                  │
          │ HTTP/WebSocket   │ HTTP/WebSocket
          ↓                  ↓
┌─────────────────────────────────────────┐
│        KkOps Backend (Go)                │
│  ┌──────────────┐  ┌──────────────┐     │
│  │ SSH Handler  │  │ RDP Handler  │     │
│  │              │  │              │     │
│  │ (WebSocket)  │  │ (Guacamole   │     │
│  │              │  │   API)       │     │
│  └──────────────┘  └──────┬───────┘     │
└────────────────────────────┼────────────┘
                              │ REST API
                              ↓
┌─────────────────────────────────────────┐
│     Apache Guacamole Server (Java)      │
│  ┌──────────────┐  ┌──────────────┐     │
│  │  Connection  │  │   Session    │     │
│  │   Manager    │  │   Manager    │     │
│  └──────┬───────┘  └──────┬───────┘     │
└─────────┼──────────────────┼────────────┘
          │                  │
          │ RDP Protocol     │ RDP Protocol
          ↓                  ↓
┌─────────────────────────────────────────┐
│      Windows Server (RDP Server)        │
└─────────────────────────────────────────┘
```

### Data Flow

#### Connection Establishment
1. User selects Windows host in KkOps UI
2. User clicks "Connect via RDP"
3. Frontend requests RDP connection from KkOps backend
4. KkOps backend creates RDP connection record
5. KkOps backend calls Guacamole API to create connection
6. Guacamole establishes RDP connection to Windows server
7. Frontend receives Guacamole WebSocket URL
8. Frontend connects to Guacamole WebSocket
9. Guacamole streams RDP data to browser

#### Data Transmission
- **User Input** (Mouse/Keyboard):
  Browser → Guacamole WebSocket → Guacamole Server → RDP Protocol → Windows Server

- **Graphics Output**:
  Windows Server → RDP Protocol → Guacamole Server → Guacamole WebSocket → Browser → Canvas

### Database Schema

#### rdp_connections
```sql
CREATE TABLE rdp_connections (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    host_id BIGINT NOT NULL REFERENCES hosts(id),
    connection_name VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    port INTEGER DEFAULT 3389,
    username VARCHAR(255) NOT NULL,
    password_encrypted TEXT, -- AES-256 加密
    domain VARCHAR(255), -- Windows 域
    width INTEGER DEFAULT 1024,
    height INTEGER DEFAULT 768,
    color_depth INTEGER DEFAULT 24, -- 8, 15, 16, 24, 32
    security_mode VARCHAR(20) DEFAULT 'rdp', -- rdp, tls, nla
    ignore_cert BOOLEAN DEFAULT false,
    guacamole_connection_id VARCHAR(255), -- Guacamole 连接 ID
    status VARCHAR(20) DEFAULT 'inactive', -- active, inactive, error
    last_connected_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rdp_connections_user_id ON rdp_connections(user_id);
CREATE INDEX idx_rdp_connections_host_id ON rdp_connections(host_id);
CREATE INDEX idx_rdp_connections_status ON rdp_connections(status);
```

#### rdp_sessions
```sql
CREATE TABLE rdp_sessions (
    id BIGSERIAL PRIMARY KEY,
    connection_id BIGINT NOT NULL REFERENCES rdp_connections(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    guacamole_session_id VARCHAR(255), -- Guacamole 会话 ID
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    duration_ms INTEGER,
    status VARCHAR(20) DEFAULT 'active', -- active, completed, error
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rdp_sessions_connection_id ON rdp_sessions(connection_id);
CREATE INDEX idx_rdp_sessions_user_id ON rdp_sessions(user_id);
CREATE INDEX idx_rdp_sessions_status ON rdp_sessions(status);
CREATE INDEX idx_rdp_sessions_started_at ON rdp_sessions(started_at);
```

### Backend Implementation

#### Models
- `RDPConnection` - RDP 连接配置模型
- `RDPSession` - RDP 会话模型

#### Repository
- `RDPConnectionRepository` - 连接数据访问
- `RDPSessionRepository` - 会话数据访问

#### Service
- `RDPConnectionService` - 连接业务逻辑
  - `CreateConnection()` - 创建连接配置
  - `UpdateConnection()` - 更新连接配置
  - `DeleteConnection()` - 删除连接
  - `ListConnections()` - 列出连接
  - `GetConnection()` - 获取连接详情

- `RDPSessionService` - 会话业务逻辑
  - `CreateSession()` - 创建会话（调用 Guacamole API）
  - `GetSessionToken()` - 获取 Guacamole 会话令牌
  - `EndSession()` - 结束会话
  - `ListSessions()` - 列出会话

- `GuacamoleService` - Guacamole 集成服务
  - `CreateGuacamoleConnection()` - 在 Guacamole 中创建连接
  - `GetSessionToken()` - 获取会话令牌
  - `DeleteGuacamoleConnection()` - 删除 Guacamole 连接

#### Handler
- `RDPConnectionHandler` - HTTP 处理器
  - `POST /api/v1/rdp/connections` - 创建连接
  - `GET /api/v1/rdp/connections` - 列出连接
  - `GET /api/v1/rdp/connections/:id` - 获取连接
  - `PUT /api/v1/rdp/connections/:id` - 更新连接
  - `DELETE /api/v1/rdp/connections/:id` - 删除连接

- `RDPSessionHandler` - HTTP 处理器
  - `POST /api/v1/rdp/sessions` - 创建会话
  - `GET /api/v1/rdp/sessions` - 列出会话
  - `GET /api/v1/rdp/sessions/:id/token` - 获取会话令牌
  - `DELETE /api/v1/rdp/sessions/:id` - 结束会话

### Frontend Implementation

#### Components
- `RDPClient` - RDP 客户端组件
  - 集成 Guacamole 客户端库
  - 显示远程桌面
  - 处理鼠标键盘输入
  - 处理连接状态

- `RDPConnectionList` - 连接列表组件
  - 显示 RDP 连接列表
  - 搜索和过滤
  - 连接操作（创建、编辑、删除、连接）

- `RDPConnectionForm` - 连接表单组件
  - 连接配置表单
  - 参数验证
  - 密码加密处理

#### Pages
- `RDPConnections.tsx` - RDP 连接管理页面
  - 连接列表
  - 连接创建和编辑
  - 连接操作

- 更新 `WebSSH.tsx` 或创建 `RemoteAccess.tsx`
  - 统一管理 SSH 和 RDP 连接
  - 标签页切换

#### Services
- `rdp.ts` - RDP API 服务
  - `listConnections()` - 获取连接列表
  - `createConnection()` - 创建连接
  - `updateConnection()` - 更新连接
  - `deleteConnection()` - 删除连接
  - `createSession()` - 创建会话
  - `getSessionToken()` - 获取会话令牌

### Guacamole Integration

#### Guacamole API Usage

**创建连接**:
```java
// Guacamole API (通过 REST API 调用)
POST /api/session/data/mysql/connections
{
  "name": "Windows Server 2019",
  "protocol": "rdp",
  "parameters": {
    "hostname": "192.168.1.100",
    "port": "3389",
    "username": "administrator",
    "password": "encrypted_password",
    "domain": "DOMAIN",
    "width": "1024",
    "height": "768",
    "color-depth": "24",
    "security": "rdp"
  }
}
```

**获取会话令牌**:
```java
POST /api/tokens
{
  "username": "guacamole_user",
  "password": "guacamole_password",
  "dataSource": "mysql"
}
```

**建立连接**:
```javascript
// Frontend
const token = await getSessionToken(connectionId);
const guacamole = new Guacamole.Client(
  new Guacamole.HTTPTunnel(`/guacamole/tunnel?token=${token}`)
);
```

### Security Considerations

#### Authentication
- Guacamole 用户与 KkOps 用户映射
- 使用 KkOps JWT token 验证
- 密码加密存储（AES-256）

#### Authorization
- 用户只能访问自己的连接
- 管理员可以访问所有连接
- 连接权限检查

#### Data Security
- RDP 连接使用加密传输（TLS/NLA）
- 密码加密存储
- 会话令牌有效期限制

### Performance Considerations

#### Network Optimization
- 图形数据压缩（Guacamole 自动处理）
- 自适应画质（根据网络状况）
- 连接池管理

#### Resource Management
- 会话超时自动断开
- 限制并发连接数
- 资源清理机制

## Migration Strategy

### Phase 1: Guacamole Deployment (3-4 days)
- 部署 Guacamole 服务器
- 配置数据库（MySQL/PostgreSQL）
- 配置认证系统
- 测试基础功能

### Phase 2: Backend Integration (4-5 days)
- 实现 RDP 连接管理 API
- 实现 Guacamole 服务集成
- 实现会话管理
- 集成认证和授权

### Phase 3: Frontend Integration (3-4 days)
- 集成 Guacamole 客户端
- 实现连接管理界面
- 实现会话界面
- 优化用户体验

### Phase 4: Testing and Optimization (2-3 days)
- 功能测试
- 性能测试
- 安全测试
- 文档更新

## Alternatives

如果 Guacamole 集成复杂度太高，可以考虑：

1. **使用 noVNC + xrdp**
   - 需要 Linux 代理服务器
   - 配置相对简单
   - 但性能可能不如 Guacamole

2. **使用第三方 RDP Web 服务**
   - 快速实现
   - 但依赖外部服务
   - 安全和隐私问题

3. **延迟实施**
   - 先实施 SSH Windows 支持
   - 根据实际需求再考虑 RDP

## Reference Documents

- **GUACAMOLE_INTEGRATION.md** - 详细的 Guacamole 集成方案
- **DEPLOYMENT_GUIDE.md** - Guacamole 部署指南

