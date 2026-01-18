# Design: 重组任务相关 URL 路由

## Context

当前系统中存在命名混淆：
- "Task" 在当前实现中指的是一次性运维执行任务
- 用户需要的"任务管理"是定时任务（Scheduled Tasks）功能
- `/tasks` URL 语义上更适合定时任务管理

## Goals / Non-Goals

### Goals
- 将 `/tasks` URL 释放给定时任务管理功能
- 保持现有"运维执行"功能完整性
- 实现定时任务的 CRUD 和调度执行
- 保持 API 设计一致性

### Non-Goals
- 不改变现有任务执行的核心逻辑
- 不引入分布式任务调度（单机调度足够）
- 不实现任务依赖关系

## Decisions

### 1. URL 命名方案

**决定**: 使用 `/executions` 作为现有"运维执行"功能的新 URL

**理由**:
- "Execution" 准确描述了该功能的本质：执行脚本
- 避免与新的"任务管理"功能混淆
- 符合 RESTful API 命名惯例

**替代方案**:
- `/operations` - 太宽泛
- `/runs` - 不够正式
- `/jobs` - 可能与定时任务混淆

### 2. 定时任务数据模型

**决定**: 创建独立的 `ScheduledTask` 模型，关联到 `TaskTemplate`

```go
type ScheduledTask struct {
    ID             uint
    Name           string
    Description    string
    CronExpression string           // Cron 表达式
    TemplateID     *uint            // 关联的执行模板（可选）
    Content        string           // 自定义脚本内容
    Type           string           // shell/python
    AssetIDs       string           // 目标主机 ID（逗号分隔）
    Timeout        int              // 超时时间
    Enabled        bool             // 是否启用
    LastRunAt      *time.Time       // 上次执行时间
    NextRunAt      *time.Time       // 下次执行时间
    LastStatus     string           // 上次执行状态
    CreatedBy      uint
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

**理由**:
- 支持基于模板创建或自定义脚本
- 独立存储调度状态便于监控
- 与执行记录分离，职责清晰

### 3. Cron 调度实现

**决定**: 使用 `robfig/cron/v3` 库实现调度

**理由**:
- Go 生态最成熟的 Cron 库
- 支持秒级精度
- 支持动态添加/删除任务
- 内置时区处理

**实现方案**:
```go
type Scheduler struct {
    cron    *cron.Cron
    db      *gorm.DB
    entries map[uint]cron.EntryID  // taskID -> cronEntryID
}

func (s *Scheduler) Start() {
    // 加载所有启用的任务
    // 添加到 cron 调度器
    s.cron.Start()
}

func (s *Scheduler) AddTask(task *ScheduledTask) {
    entryID, _ := s.cron.AddFunc(task.CronExpression, func() {
        s.executeTask(task.ID)
    })
    s.entries[task.ID] = entryID
}
```

### 4. 执行记录复用

**决定**: 定时任务执行复用现有的 `TaskExecution` 模型

**理由**:
- 执行逻辑相同（SSH 执行脚本）
- 避免重复代码
- 统一执行记录查看

**扩展字段**:
```go
type TaskExecution struct {
    // 现有字段...
    ScheduledTaskID *uint  // 关联的定时任务（新增）
    TriggerType     string // manual/scheduled（新增）
}
```

## Risks / Trade-offs

### 1. API 路径变更的兼容性
- **风险**: 前端需要同步更新所有 API 调用
- **缓解**: 一次性完成所有变更，统一发布

### 2. Cron 调度器单点问题
- **风险**: 服务重启时可能错过调度
- **缓解**: 
  - 启动时检查错过的任务
  - 记录 `next_run_at` 用于恢复
  - 未来可扩展为分布式调度

### 3. 命名一致性
- **风险**: 代码中 Task 和 Execution 混用可能造成混淆
- **缓解**: 
  - 明确命名规范：Task = 定时任务，Execution = 执行任务
  - 更新代码注释和文档

## Migration Plan

1. **Phase 1: 重命名**（不中断服务）
   - 后端添加新 API 路由（保留旧路由）
   - 前端更新为新路由
   - 移除旧路由

2. **Phase 2: 新增定时任务**
   - 创建数据模型和迁移
   - 实现调度器
   - 实现 API 和前端页面

3. **Rollback**:
   - 如果出现问题，恢复旧的 API 路由
   - 定时任务功能独立，不影响核心功能

## Open Questions

1. ~~是否需要支持任务依赖关系？~~ → 不需要，保持简单
2. ~~是否需要分布式调度？~~ → 当前不需要，单机足够
3. 是否需要支持任务执行失败重试？→ 可以作为后续优化
