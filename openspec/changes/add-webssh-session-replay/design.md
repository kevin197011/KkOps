# Design: WebSSH Session Replay

## Context

WebSSH 终端已经实现了基本的实时交互功能，但缺少会话记录和回放能力。需要添加完整的会话记录和回放系统，类似于 `script` 命令或 `asciinema` 工具的功能。

## Goals / Non-Goals

### Goals
- 记录所有终端输入输出（包括 ANSI 转义序列）
- 支持时间戳精确回放
- 提供友好的回放界面
- 支持会话搜索和过滤
- 高效的存储和检索机制

### Non-Goals
- 实时会话共享（多用户同时观看）
- 会话编辑功能（修改记录内容）
- 视频格式导出（仅支持文本回放）
- 会话加密存储（使用数据库加密即可）

## Decisions

### Decision: 使用时间戳序列记录
**What**: 记录每个输入/输出事件的时间戳（毫秒精度）
**Why**: 
- 支持精确的时间回放
- 可以计算回放速度
- 支持跳转到指定时间点

**Alternatives considered**:
- 仅记录相对时间：不够精确，无法支持跳转
- 记录绝对时间：需要处理时区问题

### Decision: 分离存储会话元数据和记录数据
**What**: 
- `webssh_sessions` 表存储会话元数据（用户、主机、时间等）
- `webssh_session_records` 表存储详细的输入输出记录

**Why**:
- 元数据查询频繁，需要索引优化
- 记录数据量大，可以分表或归档
- 便于清理旧记录

**Alternatives considered**:
- 单表存储：查询性能差，难以管理
- JSONB 存储：PostgreSQL JSONB 支持，但查询复杂

### Decision: 使用二进制格式存储记录数据
**What**: 将输入输出数据压缩后存储为二进制格式
**Why**:
- 减少存储空间（终端输出可能很大）
- 提高写入性能
- 保持数据完整性

**Alternatives considered**:
- 文本格式：存储空间大，但易于调试
- JSON 格式：结构化但体积大

### Decision: 前端回放使用 xterm.js
**What**: 复用现有的 xterm.js 终端组件进行回放
**Why**:
- 保持与实时终端一致的显示效果
- 支持 ANSI 转义序列
- 减少代码重复

**Alternatives considered**:
- 自定义回放组件：开发成本高
- 使用 iframe：隔离性好但集成复杂

## Architecture

### Data Flow

```
WebSSH Session
    ↓
Record Input/Output Events (with timestamps)
    ↓
Store to Database (compressed)
    ↓
Session List View
    ↓
Select Session → Load Records
    ↓
Replay Interface (xterm.js)
    ↓
Playback with Time Control
```

### Database Schema

#### webssh_sessions
```sql
CREATE TABLE webssh_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    host_id BIGINT NOT NULL REFERENCES hosts(id),
    hostname VARCHAR(255) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_ms INTEGER,
    record_enabled BOOLEAN DEFAULT true,
    record_data_size BIGINT, -- 记录数据大小（字节）
    status VARCHAR(20) DEFAULT 'active', -- active, completed, error
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webssh_sessions_user_id ON webssh_sessions(user_id);
CREATE INDEX idx_webssh_sessions_host_id ON webssh_sessions(host_id);
CREATE INDEX idx_webssh_sessions_started_at ON webssh_sessions(started_at);
CREATE INDEX idx_webssh_sessions_status ON webssh_sessions(status);
```

#### webssh_session_records
```sql
CREATE TABLE webssh_session_records (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES webssh_sessions(id) ON DELETE CASCADE,
    event_type VARCHAR(20) NOT NULL, -- input, output, resize
    timestamp_ms BIGINT NOT NULL, -- 相对于会话开始时间的毫秒数
    data BYTEA, -- 压缩后的数据
    data_size INTEGER, -- 原始数据大小
    rows INTEGER, -- 终端行数（resize 事件）
    columns INTEGER, -- 终端列数（resize 事件）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webssh_session_records_session_id ON webssh_session_records(session_id);
CREATE INDEX idx_webssh_session_records_timestamp ON webssh_session_records(session_id, timestamp_ms);
```

### Backend Implementation

#### Models
- `WebSSHSession` - 会话元数据模型
- `WebSSHSessionRecord` - 会话记录模型

#### Repository
- `WebSSHSessionRepository` - 会话数据访问
- `WebSSHSessionRecordRepository` - 记录数据访问

#### Service
- `WebSSHSessionService` - 会话业务逻辑
  - `CreateSession()` - 创建会话
  - `EndSession()` - 结束会话
  - `ListSessions()` - 列出会话
  - `GetSession()` - 获取会话详情
  - `GetSessionRecords()` - 获取会话记录
  - `DeleteSession()` - 删除会话

#### Handler
- `WebSSHSessionHandler` - HTTP 处理器
  - `GET /api/v1/webssh/sessions` - 会话列表
  - `GET /api/v1/webssh/sessions/:id` - 会话详情
  - `GET /api/v1/webssh/sessions/:id/records` - 会话记录
  - `DELETE /api/v1/webssh/sessions/:id` - 删除会话

#### Recording Integration
在 `WebSSHHandler.HandleTerminal()` 中：
- 会话开始时创建 `WebSSHSession` 记录
- 记录所有输入事件（input, resize）
- 记录所有输出事件（stdout, stderr）
- 会话结束时更新会话状态和持续时间

### Frontend Implementation

#### Components
- `SessionReplay` - 回放组件
  - 使用 xterm.js 显示回放
  - 时间轴控制（播放、暂停、速度、跳转）
  - 进度显示

- `SessionList` - 会话列表组件
  - 会话列表展示
  - 搜索和过滤
  - 会话操作（查看、删除）

#### Pages
- `WebSSH.tsx` - 更新 WebSSH 页面
  - 添加"会话记录"标签页
  - 集成会话列表和回放功能

#### Services
- `webssh.ts` - 添加回放相关 API
  - `listSessions()` - 获取会话列表
  - `getSession()` - 获取会话详情
  - `getSessionRecords()` - 获取会话记录
  - `deleteSession()` - 删除会话

### Recording Format

#### Event Types
- `input` - 用户输入（命令、按键）
- `output` - 终端输出（stdout）
- `error` - 错误输出（stderr）
- `resize` - 终端尺寸变化

#### Record Structure
```json
{
  "session_id": 123,
  "events": [
    {
      "type": "input",
      "timestamp_ms": 0,
      "data": "ls -la\n"
    },
    {
      "type": "output",
      "timestamp_ms": 50,
      "data": "total 1234\n..."
    },
    {
      "type": "resize",
      "timestamp_ms": 1000,
      "rows": 24,
      "columns": 80
    }
  ]
}
```

### Replay Algorithm

1. **Load Session**: 加载会话元数据和所有记录
2. **Initialize Terminal**: 创建 xterm.js 实例，设置初始尺寸
3. **Playback Loop**:
   - 按时间戳顺序处理事件
   - 计算事件之间的延迟
   - 使用 `setTimeout` 控制回放速度
   - 根据事件类型更新终端显示
4. **Control**:
   - 播放/暂停：控制回放循环
   - 速度控制：调整延迟倍数
   - 跳转：直接跳到指定时间点

### Performance Considerations

#### Storage Optimization
- 使用压缩算法（gzip）压缩记录数据
- 定期清理旧记录（可配置保留期限）
- 考虑归档到对象存储（S3/OSS）

#### Query Optimization
- 会话列表使用分页
- 记录数据按需加载（懒加载）
- 使用索引优化查询性能

#### Replay Performance
- 前端缓存已加载的记录
- 支持分段加载（长会话）
- 使用 Web Worker 处理大量数据（可选）

## Security Considerations

### Access Control
- 用户只能查看自己的会话记录
- 管理员可以查看所有会话记录
- 删除操作需要权限验证

### Data Privacy
- 敏感信息（密码）在记录时可能需要脱敏
- 会话记录加密存储（可选）
- 支持会话记录导出时的数据脱敏

### Audit Integration
- 会话查看操作记录到审计日志
- 会话删除操作记录到审计日志

## Migration Strategy

1. **Phase 1**: 数据库迁移
   - 创建新表
   - 添加索引

2. **Phase 2**: 后端实现
   - 实现模型和仓库
   - 实现服务层
   - 集成记录功能到 WebSSH Handler

3. **Phase 3**: 前端实现
   - 实现回放组件
   - 集成到 WebSSH 页面
   - 添加会话列表

4. **Phase 4**: 测试和优化
   - 功能测试
   - 性能测试
   - 存储优化

