# Design: Batch Operations Page

## Architecture Overview

### 系统架构
```
Frontend (React)
  ├── BatchOperations Page
  │   ├── HostSelector Component (主机选择器)
  │   ├── CommandTemplateManager Component (命令模板管理)
  │   ├── CommandExecutor Component (命令执行器)
  │   └── ResultViewer Component (结果查看器)
  │
Backend (Go)
  ├── BatchOperationsHandler (批量操作处理器)
  ├── BatchOperationsService (批量操作服务)
  ├── BatchOperationsRepository (批量操作数据访问)
  ├── CommandTemplateService (命令模板服务)
  └── Salt Client (Salt API 客户端)
```

## Database Design

### batch_operations 表
存储批量操作历史记录。

```sql
CREATE TABLE batch_operations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,                    -- 操作名称
    description TEXT,                               -- 操作描述
    command_type VARCHAR(50) NOT NULL,             -- 命令类型: 'custom', 'template', 'builtin'
    command_function VARCHAR(255) NOT NULL,        -- Salt 函数名
    command_args JSONB,                             -- 命令参数
    target_hosts JSONB NOT NULL,                     -- 目标主机列表 [{"id": 1, "hostname": "..."}, ...]
    target_count INTEGER NOT NULL,                  -- 目标主机数量
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, running, completed, failed, cancelled
    salt_job_id VARCHAR(255),                       -- Salt Job ID
    started_by INTEGER NOT NULL,                    -- 执行用户ID
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    duration_seconds INTEGER,                       -- 执行耗时（秒）
    results JSONB,                                  -- 执行结果 {host_id: {success: bool, output: string, error: string}}
    success_count INTEGER DEFAULT 0,                -- 成功数量
    failed_count INTEGER DEFAULT 0,                 -- 失败数量
    error_message TEXT,                             -- 错误信息
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (started_by) REFERENCES users(id)
);

CREATE INDEX idx_batch_operations_started_by ON batch_operations(started_by);
CREATE INDEX idx_batch_operations_status ON batch_operations(status);
CREATE INDEX idx_batch_operations_started_at ON batch_operations(started_at);
```

### command_templates 表
存储常用命令模板。

```sql
CREATE TABLE command_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,                    -- 模板名称
    description TEXT,                               -- 模板描述
    category VARCHAR(50),                           -- 分类: 'system', 'network', 'disk', 'process', 'custom'
    command_function VARCHAR(255) NOT NULL,        -- Salt 函数名
    command_args JSONB,                             -- 默认参数
    icon VARCHAR(50),                               -- 图标名称
    created_by INTEGER NOT NULL,                    -- 创建用户ID
    is_public BOOLEAN DEFAULT FALSE,                -- 是否公开（所有用户可见）
    usage_count INTEGER DEFAULT 0,                  -- 使用次数
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_command_templates_category ON command_templates(category);
CREATE INDEX idx_command_templates_created_by ON command_templates(created_by);
```

## Backend Design

### Models

#### BatchOperation Model
```go
type BatchOperation struct {
    ID              uint           `gorm:"primaryKey"`
    Name            string         `gorm:"not null"`
    Description     string
    CommandType     string         `gorm:"not null"` // custom, template, builtin
    CommandFunction string         `gorm:"not null"`
    CommandArgs     datatypes.JSON `gorm:"type:jsonb"`
    TargetHosts     datatypes.JSON `gorm:"type:jsonb;not null"`
    TargetCount     int            `gorm:"not null"`
    Status          string         `gorm:"default:'pending'"`
    SaltJobID       string
    StartedBy       uint           `gorm:"not null"`
    StartedAt       time.Time      `gorm:"default:now()"`
    CompletedAt     *time.Time
    DurationSeconds *int
    Results         datatypes.JSON `gorm:"type:jsonb"`
    SuccessCount    int            `gorm:"default:0"`
    FailedCount     int            `gorm:"default:0"`
    ErrorMessage    string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    StartedByUser   User           `gorm:"foreignKey:StartedBy"`
}
```

#### CommandTemplate Model
```go
type CommandTemplate struct {
    ID              uint           `gorm:"primaryKey"`
    Name            string         `gorm:"not null"`
    Description     string
    Category        string
    CommandFunction string         `gorm:"not null"`
    CommandArgs     datatypes.JSON `gorm:"type:jsonb"`
    Icon            string
    CreatedBy       uint           `gorm:"not null"`
    IsPublic        bool           `gorm:"default:false"`
    UsageCount      int            `gorm:"default:0"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    CreatedByUser   User           `gorm:"foreignKey:CreatedBy"`
}
```

### Service Layer

#### BatchOperationsService
- `CreateOperation()`: 创建批量操作
- `ExecuteOperation()`: 执行批量操作（异步）
- `GetOperationStatus()`: 获取操作状态
- `GetOperationResults()`: 获取操作结果
- `CancelOperation()`: 取消正在执行的操作
- `ListOperations()`: 列出操作历史
- `GetOperation()`: 获取单个操作详情

#### CommandTemplateService
- `CreateTemplate()`: 创建命令模板
- `UpdateTemplate()`: 更新命令模板
- `DeleteTemplate()`: 删除命令模板
- `ListTemplates()`: 列出模板（支持按分类、用户筛选）
- `GetTemplate()`: 获取模板详情
- `IncrementUsageCount()`: 增加使用次数

### Handler Layer

#### BatchOperationsHandler
- `POST /api/v1/batch-operations`: 创建并执行批量操作
- `GET /api/v1/batch-operations`: 列出批量操作历史
- `GET /api/v1/batch-operations/:id`: 获取操作详情
- `GET /api/v1/batch-operations/:id/status`: 获取操作状态
- `GET /api/v1/batch-operations/:id/results`: 获取操作结果
- `POST /api/v1/batch-operations/:id/cancel`: 取消操作

#### CommandTemplateHandler
- `POST /api/v1/command-templates`: 创建命令模板
- `GET /api/v1/command-templates`: 列出命令模板
- `GET /api/v1/command-templates/:id`: 获取模板详情
- `PUT /api/v1/command-templates/:id`: 更新模板
- `DELETE /api/v1/command-templates/:id`: 删除模板

## Frontend Design

### Page Structure

#### BatchOperations.tsx
主页面，包含以下区域：
1. **顶部工具栏**
   - 操作名称输入
   - 快速操作按钮（常用命令）
   - 保存为模板按钮

2. **左侧面板 - 主机选择器**
   - 项目筛选
   - 分组筛选
   - 标签筛选
   - 状态筛选
   - 主机列表（多选）
   - 已选主机数量显示

3. **中间面板 - 命令配置**
   - 命令类型选择（自定义/模板/内置）
   - 命令输入/选择
   - 参数配置
   - 执行按钮

4. **右侧面板 - 结果展示**
   - 执行状态（进度条）
   - 结果表格（主机、状态、输出）
   - 结果统计（成功/失败数量）
   - 结果导出按钮

5. **底部面板 - 操作历史**
   - 历史操作列表
   - 快速重试按钮

### Components

#### HostSelector.tsx
主机选择器组件，支持：
- 多级筛选（项目、分组、标签、状态）
- 主机多选
- 已选主机列表展示
- 主机搜索

#### CommandTemplateManager.tsx
命令模板管理组件，支持：
- 模板列表展示（按分类）
- 模板创建/编辑
- 模板删除
- 模板使用

#### ResultViewer.tsx
结果查看器组件，支持：
- 实时更新执行状态
- 结果表格展示
- 结果详情查看（展开/收起）
- 结果导出（CSV/JSON）

## Built-in Command Templates

系统内置常用命令模板：

1. **系统信息**
   - 查看系统负载 (`status.loadavg`)
   - 查看内存使用 (`status.meminfo`)
   - 查看磁盘使用 (`disk.usage`)
   - 查看 CPU 信息 (`status.cpuinfo`)

2. **网络信息**
   - 查看网络接口 (`network.interfaces`)
   - 查看网络连接 (`network.active_tcp`)

3. **进程管理**
   - 查看进程列表 (`ps.procs`)
   - 查看服务状态 (`service.status`)

4. **文件操作**
   - 查看文件内容 (`file.read`)
   - 检查文件是否存在 (`file.file_exists`)

## Execution Flow

1. **用户选择主机和命令**
   - 用户在主机选择器中选择目标主机
   - 用户选择或输入命令

2. **创建批量操作**
   - 前端调用 `POST /api/v1/batch-operations`
   - 后端创建 BatchOperation 记录，状态为 `pending`

3. **异步执行**
   - 后端启动 goroutine 异步执行
   - 调用 Salt API 执行命令
   - 更新操作状态为 `running`
   - 保存 Salt Job ID

4. **结果收集**
   - 定期查询 Salt Job 状态
   - 收集执行结果
   - 更新操作状态和结果

5. **前端轮询**
   - 前端定期轮询操作状态
   - 实时更新结果展示
   - 操作完成后停止轮询

## Error Handling

- **Salt API 连接失败**: 返回错误，操作状态设为 `failed`
- **部分主机执行失败**: 记录失败主机，成功数量减少
- **超时处理**: 设置操作超时时间（默认 5 分钟），超时后取消操作
- **网络错误**: 重试机制（最多 3 次）

## Security Considerations

- **权限控制**: 基于 RBAC，只有有权限的用户才能执行批量操作
- **审计日志**: 所有批量操作记录到审计日志
- **命令验证**: 验证命令和参数，防止危险命令执行
- **主机访问控制**: 用户只能选择有权限访问的主机

## Performance Considerations

- **异步执行**: 批量操作异步执行，不阻塞请求
- **结果分页**: 操作历史支持分页
- **缓存**: 常用命令模板可以缓存
- **批量查询优化**: 主机列表查询使用索引优化

