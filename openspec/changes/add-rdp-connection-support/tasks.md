# Implementation Tasks: RDP Connection Support

**创建日期**: 2025-01-28  
**状态**: ⏸️ 待实施  
**优先级**: 🟡 中  
**复杂度**: 🔴 高  
**工作量估算**: 12-18 天

---

## 任务清单

### 1. Guacamole 部署和配置

- [ ] 1.1 安装 Guacamole 服务器
  - 下载 Guacamole 服务器包
  - 安装 Java 运行环境
  - 配置 Tomcat 或使用 Docker
  - 文档：https://guacamole.apache.org/doc/gug/installing-guacamole.html

- [ ] 1.2 配置 Guacamole 数据库
  - 创建 MySQL/PostgreSQL 数据库
  - 执行 Guacamole 初始化 SQL
  - 配置数据库连接
  - 测试数据库连接

- [ ] 1.3 配置 Guacamole 认证
  - 配置数据库认证（推荐）
  - 或配置 LDAP 认证（可选）
  - 创建 Guacamole 管理员用户
  - 测试认证功能

- [ ] 1.4 配置 Guacamole REST API
  - 启用 REST API
  - 配置 API 认证
  - 测试 API 访问
  - 文档：https://guacamole.apache.org/doc/gug/rest-api.html

### 2. 数据库设计和迁移

- [ ] 2.1 创建 `rdp_connections` 表迁移文件
  - 定义表结构
  - 添加外键约束
  - 添加索引
  - 创建迁移文件：`backend/migrations/YYYYMMDDHHMMSS_create_rdp_connections_table.up.sql`

- [ ] 2.2 创建 `rdp_sessions` 表迁移文件
  - 定义表结构
  - 添加外键约束
  - 添加索引
  - 创建迁移文件：`backend/migrations/YYYYMMDDHHMMSS_create_rdp_sessions_table.up.sql`

- [ ] 2.3 创建回滚迁移文件
  - 创建对应的 `.down.sql` 文件

### 3. 后端模型层

- [ ] 3.1 创建 `RDPConnection` 模型
  - 定义结构体字段
  - 添加 GORM 标签
  - 添加 JSON 标签
  - 实现 `TableName()` 方法
  - 文件：`backend/internal/models/rdp_connection.go`

- [ ] 3.2 创建 `RDPSession` 模型
  - 定义结构体字段
  - 添加 GORM 标签
  - 添加 JSON 标签
  - 实现 `TableName()` 方法
  - 文件：`backend/internal/models/rdp_session.go`

### 4. 后端仓库层

- [ ] 4.1 创建 `RDPConnectionRepository`
  - 实现接口定义
  - 实现 `Create()` 方法
  - 实现 `GetByID()` 方法
  - 实现 `List()` 方法（支持分页、过滤）
  - 实现 `Update()` 方法
  - 实现 `Delete()` 方法
  - 文件：`backend/internal/repository/rdp_connection_repository.go`

- [ ] 4.2 创建 `RDPSessionRepository`
  - 实现接口定义
  - 实现 `Create()` 方法
  - 实现 `GetByID()` 方法
  - 实现 `List()` 方法
  - 实现 `Update()` 方法
  - 实现 `Delete()` 方法
  - 文件：`backend/internal/repository/rdp_session_repository.go`

### 5. Guacamole 集成服务

- [ ] 5.1 创建 `GuacamoleService`
  - 实现 Guacamole REST API 客户端
  - 实现 `CreateConnection()` 方法
  - 实现 `UpdateConnection()` 方法
  - 实现 `DeleteConnection()` 方法
  - 实现 `GetSessionToken()` 方法
  - 处理 API 认证
  - 处理错误和重试
  - 文件：`backend/internal/service/guacamole_service.go`

- [ ] 5.2 配置 Guacamole 服务
  - 添加 Guacamole 配置（URL、用户名、密码）
  - 在 `config.go` 中添加配置项
  - 实现配置加载

### 6. 后端服务层

- [ ] 6.1 创建 `RDPConnectionService`
  - 实现接口定义
  - 实现 `CreateConnection()` 方法
    - 验证连接参数
    - 加密密码存储
    - 调用 Guacamole API 创建连接
    - 保存连接记录
  - 实现 `UpdateConnection()` 方法
  - 实现 `DeleteConnection()` 方法
  - 实现 `ListConnections()` 方法
  - 实现 `GetConnection()` 方法
  - 实现权限检查
  - 文件：`backend/internal/service/rdp_connection_service.go`

- [ ] 6.2 创建 `RDPSessionService`
  - 实现接口定义
  - 实现 `CreateSession()` 方法
    - 创建会话记录
    - 调用 Guacamole API 获取会话令牌
    - 返回会话信息
  - 实现 `GetSessionToken()` 方法
  - 实现 `EndSession()` 方法
  - 实现 `ListSessions()` 方法
  - 文件：`backend/internal/service/rdp_session_service.go`

### 7. 后端 Handler 层

- [ ] 7.1 创建 `RDPConnectionHandler`
  - 实现 `CreateConnection()` 方法
    - `POST /api/v1/rdp/connections`
  - 实现 `ListConnections()` 方法
    - `GET /api/v1/rdp/connections`
  - 实现 `GetConnection()` 方法
    - `GET /api/v1/rdp/connections/:id`
  - 实现 `UpdateConnection()` 方法
    - `PUT /api/v1/rdp/connections/:id`
  - 实现 `DeleteConnection()` 方法
    - `DELETE /api/v1/rdp/connections/:id`
  - 添加请求验证
  - 添加权限检查
  - 文件：`backend/internal/handler/rdp_connection_handler.go`

- [ ] 7.2 创建 `RDPSessionHandler`
  - 实现 `CreateSession()` 方法
    - `POST /api/v1/rdp/sessions`
  - 实现 `ListSessions()` 方法
    - `GET /api/v1/rdp/sessions`
  - 实现 `GetSessionToken()` 方法
    - `GET /api/v1/rdp/sessions/:id/token`
  - 实现 `EndSession()` 方法
    - `DELETE /api/v1/rdp/sessions/:id`
  - 添加权限检查
  - 文件：`backend/internal/handler/rdp_session_handler.go`

- [ ] 7.3 注册路由
  - 在 `main.go` 中注册 RDP 路由
  - 添加认证中间件
  - 添加权限检查中间件

### 8. 前端 Service 层

- [ ] 8.1 创建 `rdp.ts` Service
  - 定义 TypeScript 接口：
    - `RDPConnection`
    - `RDPSession`
    - `SessionToken`
  - 实现 `listConnections()` 方法
  - 实现 `createConnection()` 方法
  - 实现 `updateConnection()` 方法
  - 实现 `deleteConnection()` 方法
  - 实现 `createSession()` 方法
  - 实现 `getSessionToken()` 方法
  - 实现 `endSession()` 方法
  - 文件：`frontend/src/services/rdp.ts`

### 9. 前端组件

- [ ] 9.1 集成 Guacamole 客户端库
  - 安装 Guacamole 客户端库（npm）
  - 或使用 CDN 引入
  - 配置 Guacamole 客户端
  - 文件：`frontend/src/components/RDPClient.tsx`

- [ ] 9.2 创建 `RDPClient` 组件
  - 初始化 Guacamole 客户端
  - 建立 WebSocket 连接
  - 显示远程桌面（Canvas）
  - 处理鼠标事件
  - 处理键盘事件
  - 处理连接状态
  - 处理错误
  - 文件：`frontend/src/components/RDPClient.tsx`

- [ ] 9.3 创建 `RDPConnectionList` 组件
  - 显示连接列表（表格）
  - 实现搜索和过滤
  - 实现分页
  - 实现连接操作（创建、编辑、删除、连接）
  - 文件：`frontend/src/components/RDPConnectionList.tsx`

- [ ] 9.4 创建 `RDPConnectionForm` 组件
  - 连接配置表单
  - 字段验证
  - 密码输入和处理
  - 参数配置（分辨率、颜色深度等）
  - 文件：`frontend/src/components/RDPConnectionForm.tsx`

### 10. 前端页面

- [ ] 10.1 创建 `RDPConnections.tsx` 页面
  - 集成 `RDPConnectionList` 组件
  - 集成 `RDPConnectionForm` 组件（模态框）
  - 实现连接创建和编辑
  - 实现连接删除
  - 实现连接操作
  - 文件：`frontend/src/pages/RDPConnections.tsx`

- [ ] 10.2 创建或更新远程访问页面
  - 选项 A: 更新 `WebSSH.tsx` 添加 RDP 标签页
  - 选项 B: 创建 `RemoteAccess.tsx` 统一管理 SSH 和 RDP
  - 实现标签页切换
  - 统一连接管理界面

- [ ] 10.3 添加路由和菜单
  - 在 `App.tsx` 中添加 RDP 路由
  - 在 `MainLayout.tsx` 中添加菜单项
  - 配置权限检查

### 11. 主机管理集成

- [ ] 11.1 更新主机详情页面
  - 显示 RDP 连接选项
  - 支持快速创建 RDP 连接
  - 显示连接状态

- [ ] 11.2 更新主机模型（可选）
  - 添加 `rdp_enabled` 字段
  - 添加 `rdp_port` 字段（默认 3389）
  - 创建迁移文件

### 12. 安全和权限

- [ ] 12.1 实现密码加密
  - 使用 AES-256 加密存储密码
  - 实现加密/解密函数
  - 在创建/更新连接时加密密码

- [ ] 12.2 实现权限检查
  - 用户只能访问自己的连接
  - 管理员可以访问所有连接
  - 在 Service 层实现权限检查

- [ ] 12.3 集成审计日志
  - RDP 连接创建记录到审计日志
  - RDP 会话建立记录到审计日志
  - RDP 连接删除记录到审计日志

### 13. 测试

- [ ] 13.1 Guacamole 功能测试
  - 测试 Guacamole 部署
  - 测试 RDP 连接建立
  - 测试图形界面显示
  - 测试鼠标键盘输入

- [ ] 13.2 后端单元测试
  - 测试 RDPConnectionService 方法
  - 测试 RDPSessionService 方法
  - 测试 GuacamoleService 方法
  - 测试 Repository 方法

- [ ] 13.3 后端集成测试
  - 测试 RDP 连接创建流程
  - 测试会话创建流程
  - 测试 Guacamole API 集成

- [ ] 13.4 前端单元测试
  - 测试 RDPClient 组件
  - 测试 RDPConnectionList 组件
  - 测试 RDPConnectionForm 组件

- [ ] 13.5 端到端测试
  - 测试完整的 RDP 连接流程
  - 测试不同 Windows 版本
  - 测试不同分辨率
  - 测试性能（延迟、画质）

### 14. 文档和部署

- [ ] 14.1 更新部署文档
  - Guacamole 部署说明
  - 配置说明
  - 网络配置说明

- [ ] 14.2 更新用户文档
  - RDP 连接使用说明
  - 常见问题解答

- [ ] 14.3 更新 API 文档
  - 文档化 RDP API 端点
  - 添加请求/响应示例

---

## 完成情况

**总任务数**: 约 60+ 个任务项  
**已完成**: 0  
**待完成**: 60+  
**完成率**: 0%

---

## 实施优先级

### 高优先级 🔴
1. Guacamole 部署和配置（必须）
2. 数据库设计和迁移（必须）
3. 后端基础功能（连接管理、会话管理）（必须）
4. 前端基础功能（连接列表、客户端）（必须）

### 中优先级 🟡
5. Guacamole 集成服务（重要）
6. 安全和权限（重要）
7. 主机管理集成（重要）

### 低优先级 🟢
8. 高级功能（剪贴板、文件传输等）（可选）
9. 性能优化（可选）
10. 详细文档（可选）

---

## 技术注意事项

1. **Guacamole 部署**: 需要 Java 运行环境，需要额外的服务器资源
2. **性能考虑**: RDP 图形传输数据量大，需要优化网络和压缩
3. **安全考虑**: 密码加密存储，RDP 连接使用加密传输
4. **兼容性**: 需要测试不同 Windows 版本和 Guacamole 版本

---

## 风险评估

- **技术风险**: 🟡 中-高（需要集成第三方服务）
- **实施风险**: 🟡 中（Guacamole 配置复杂）
- **性能风险**: 🟡 中（图形传输需要优化）
- **维护风险**: 🟡 中（需要维护 Guacamole 服务）

---

## 替代方案

如果 Guacamole 集成复杂度太高，可以考虑：

1. **延迟实施**: 先实施 SSH Windows 支持，后续再考虑 RDP
2. **使用 noVNC + xrdp**: 相对简单，但需要 Linux 代理
3. **推荐用户使用原生 RDP 客户端**: 功能完整，但体验不统一

---

**最后更新**: 2025-01-28

