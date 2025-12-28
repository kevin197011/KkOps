# Implementation Tasks: WebSSH Session Replay

**创建日期**: 2025-01-28  
**状态**: ⏸️ 待实施  
**优先级**: 🟡 中

---

## 任务清单

### 1. 数据库设计和迁移

- [ ] 1.1 创建 `webssh_sessions` 表迁移文件
  - 定义表结构（id, user_id, host_id, hostname, started_at, ended_at, duration_ms, record_enabled, record_data_size, status）
  - 添加外键约束（users, hosts）
  - 添加索引（user_id, host_id, started_at, status）
  - 创建迁移文件：`backend/migrations/YYYYMMDDHHMMSS_create_webssh_sessions_table.up.sql`

- [ ] 1.2 创建 `webssh_session_records` 表迁移文件
  - 定义表结构（id, session_id, event_type, timestamp_ms, data, data_size, rows, columns）
  - 添加外键约束（webssh_sessions，级联删除）
  - 添加索引（session_id, session_id+timestamp_ms）
  - 创建迁移文件：`backend/migrations/YYYYMMDDHHMMSS_create_webssh_session_records_table.up.sql`

- [ ] 1.3 创建回滚迁移文件
  - 创建对应的 `.down.sql` 文件

### 2. 后端模型层

- [ ] 2.1 创建 `WebSSHSession` 模型
  - 定义结构体字段
  - 添加 GORM 标签
  - 添加 JSON 标签
  - 实现 `TableName()` 方法
  - 文件：`backend/internal/models/webssh_session.go`

- [ ] 2.2 创建 `WebSSHSessionRecord` 模型
  - 定义结构体字段
  - 添加 GORM 标签
  - 添加 JSON 标签
  - 实现 `TableName()` 方法
  - 文件：`backend/internal/models/webssh_session_record.go`

- [ ] 2.3 注册模型到 AutoMigrate
  - 在 `database.go` 中添加模型注册（如果需要）

### 3. 后端仓库层

- [ ] 3.1 创建 `WebSSHSessionRepository`
  - 实现接口定义
  - 实现 `Create()` 方法
  - 实现 `GetByID()` 方法
  - 实现 `List()` 方法（支持分页、过滤、排序）
  - 实现 `Update()` 方法
  - 实现 `Delete()` 方法
  - 文件：`backend/internal/repository/webssh_session_repository.go`

- [ ] 3.2 创建 `WebSSHSessionRecordRepository`
  - 实现接口定义
  - 实现 `Create()` 方法（批量创建）
  - 实现 `GetBySessionID()` 方法（按时间戳排序）
  - 实现 `DeleteBySessionID()` 方法（级联删除）
  - 实现 `GetEventCount()` 方法（统计事件数量）
  - 文件：`backend/internal/repository/webssh_session_record_repository.go`

### 4. 后端服务层

- [ ] 4.1 创建 `WebSSHSessionService`
  - 实现接口定义
  - 实现 `CreateSession()` 方法
  - 实现 `EndSession()` 方法
  - 实现 `ListSessions()` 方法（支持过滤、分页）
  - 实现 `GetSession()` 方法
  - 实现 `GetSessionRecords()` 方法
  - 实现 `DeleteSession()` 方法
  - 实现权限检查（用户只能访问自己的会话，管理员可以访问所有）
  - 文件：`backend/internal/service/webssh_session_service.go`

- [ ] 4.2 实现数据压缩功能
  - 实现 `compressData()` 方法（使用 gzip）
  - 实现 `decompressData()` 方法
  - 在存储记录时压缩数据
  - 在读取记录时解压数据
  - 文件：`backend/internal/service/webssh_session_service.go` 或新建工具文件

### 5. WebSSH Handler 集成

- [ ] 5.1 修改 `WebSSHHandler` 添加记录功能
  - 在 `HandleTerminal()` 开始时创建会话记录
  - 记录所有输入事件（input, resize）
  - 记录所有输出事件（stdout, stderr）
  - 在会话结束时更新会话状态
  - 使用缓冲区批量写入记录（提高性能）
  - 文件：`backend/internal/handler/webssh_handler.go`

- [ ] 5.2 实现事件记录逻辑
  - 创建 `recordEvent()` 辅助方法
  - 处理时间戳计算（相对于会话开始时间）
  - 处理数据压缩
  - 实现批量写入机制（避免频繁数据库操作）

### 6. 后端 API Handler

- [ ] 6.1 创建 `WebSSHSessionHandler`
  - 实现 `ListSessions()` 方法
    - `GET /api/v1/webssh/sessions` - 会话列表
    - 支持查询参数：page, page_size, user_id, host_id, start_date, end_date, status
  - 实现 `GetSession()` 方法
    - `GET /api/v1/webssh/sessions/:id` - 会话详情
  - 实现 `GetSessionRecords()` 方法
    - `GET /api/v1/webssh/sessions/:id/records` - 会话记录
    - 支持查询参数：limit, offset（用于分段加载）
  - 实现 `DeleteSession()` 方法
    - `DELETE /api/v1/webssh/sessions/:id` - 删除会话
  - 添加权限检查中间件
  - 文件：`backend/internal/handler/webssh_session_handler.go`

- [ ] 6.2 注册路由
  - 在 `main.go` 中注册 WebSSH Session 路由
  - 添加认证中间件
  - 添加权限检查（webssh:read, webssh:delete）

### 7. 前端 Service 层

- [ ] 7.1 更新 `webssh.ts` Service
  - 添加 `listSessions()` 方法
  - 添加 `getSession()` 方法
  - 添加 `getSessionRecords()` 方法
  - 添加 `deleteSession()` 方法
  - 定义 TypeScript 接口：
    - `WebSSHSession`
    - `WebSSHSessionRecord`
    - `SessionListParams`
  - 文件：`frontend/src/services/webssh.ts`

### 8. 前端组件

- [ ] 8.1 创建 `SessionReplay` 组件
  - 使用 xterm.js 显示回放
  - 实现播放控制（播放、暂停、速度、跳转）
  - 实现进度条和时间显示
  - 实现会话信息显示
  - 处理事件回放逻辑
  - 支持终端尺寸变化
  - 文件：`frontend/src/components/SessionReplay.tsx`

- [ ] 8.2 创建 `SessionList` 组件
  - 显示会话列表（表格）
  - 实现搜索和过滤功能
  - 实现分页
  - 实现会话操作（查看、删除）
  - 文件：`frontend/src/components/SessionList.tsx`

- [ ] 8.3 更新 `Terminal` 组件
  - 添加会话记录状态跟踪
  - 在会话开始时通知后端创建记录
  - 在会话结束时通知后端结束记录
  - （可选）显示记录状态指示器

### 9. 前端页面

- [ ] 9.1 更新 `WebSSH.tsx` 页面
  - 添加"会话记录"标签页
  - 集成 `SessionList` 组件
  - 集成 `SessionReplay` 组件
  - 实现会话选择和回放流程
  - 文件：`frontend/src/pages/WebSSH.tsx`

### 10. 回放算法实现

- [ ] 10.1 实现回放核心逻辑
  - 事件排序（按 timestamp_ms）
  - 时间延迟计算
  - 播放循环实现
  - 暂停/恢复机制
  - 速度控制（调整延迟倍数）
  - 跳转功能（重建终端状态到指定时间点）

- [ ] 10.2 实现终端状态重建
  - 处理跳转时的状态重建
  - 从开始到目标时间点快速回放所有事件
  - 保持终端显示一致性

### 11. 性能优化

- [ ] 11.1 实现批量记录写入
  - 使用缓冲区收集事件
  - 定期批量写入数据库（如每 100 个事件或每 5 秒）
  - 会话结束时写入剩余事件

- [ ] 11.2 实现记录数据压缩
  - 使用 gzip 压缩记录数据
  - 在存储前压缩，读取后解压
  - 记录压缩前后大小（用于统计）

- [ ] 11.3 实现分段加载
  - 长会话记录分段加载（避免一次性加载大量数据）
  - 实现懒加载机制
  - 前端缓存已加载的记录

### 12. 权限和安全

- [ ] 12.1 实现权限检查
  - 用户只能查看自己的会话
  - 管理员可以查看所有会话
  - 删除操作需要权限验证
  - 在 Service 层实现权限检查

- [ ] 12.2 集成审计日志
  - 会话查看记录到审计日志
  - 会话删除记录到审计日志
  - 使用现有的审计中间件

### 13. 测试

- [ ] 13.1 后端单元测试
  - 测试 WebSSHSessionService 方法
  - 测试 Repository 方法
  - 测试数据压缩/解压
  - 测试权限检查

- [ ] 13.2 后端集成测试
  - 测试会话记录流程
  - 测试会话查询和过滤
  - 测试会话删除

- [ ] 13.3 前端单元测试
  - 测试 SessionReplay 组件
  - 测试 SessionList 组件
  - 测试回放控制逻辑

- [ ] 13.4 手动测试
  - 测试完整的记录和回放流程
  - 测试各种终端操作（命令、特殊键、resize）
  - 测试回放控制（播放、暂停、速度、跳转）
  - 测试权限控制
  - 测试性能（大量事件的处理）

### 14. 文档和清理

- [ ] 14.1 更新 API 文档
  - 文档化新的 API 端点
  - 添加请求/响应示例

- [ ] 14.2 更新用户文档
  - 添加会话回放功能说明
  - 添加使用指南

- [ ] 14.3 代码清理
  - 移除调试代码
  - 添加代码注释
  - 运行代码格式化

---

## 完成情况

**总任务数**: 约 60+ 个任务项  
**已完成**: 0  
**待完成**: 60+  
**完成率**: 0%

---

## 注意事项

1. **存储空间**: 会话记录可能占用大量存储空间，需要考虑：
   - 数据压缩
   - 定期清理旧记录
   - 可配置的保留期限

2. **性能**: 大量事件的回放可能影响性能，需要：
   - 优化回放算法
   - 使用虚拟滚动（如果列表很长）
   - 考虑使用 Web Worker 处理大量数据

3. **安全性**: 
   - 敏感信息（如密码）在记录时可能需要脱敏
   - 会话记录访问需要权限控制
   - 考虑会话记录加密存储

4. **兼容性**: 
   - 确保与现有 WebSSH 功能兼容
   - 不影响现有终端交互性能

