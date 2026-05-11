# Design: 审计日志功能

## Context

KkOps 作为运维管理平台，需要记录用户的所有关键操作以满足安全审计和合规要求。当前系统仅在仪表板显示任务执行的最近活动，缺乏完整的审计能力。

### 约束条件
- 审计日志必须不可篡改（只读）
- 需要支持高效的时间范围查询
- 敏感信息不能记录（如密码明文）
- 需要考虑日志量增长对存储的影响

## Goals / Non-Goals

### Goals
- 记录所有关键用户操作
- 提供灵活的查询和筛选能力
- 支持审计日志导出
- 保证审计数据的完整性和不可篡改性

### Non-Goals
- 不实现实时告警功能（可后续扩展）
- 不实现日志归档到外部存储（可后续扩展）
- 不实现审计日志的统计分析报表（可后续扩展）

## Decisions

### 1. 数据模型设计

```go
type AuditLog struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      uint      `gorm:"index"`           // 操作用户 ID
    Username    string    `gorm:"size:100;index"`  // 操作用户名（冗余存储，防止用户删除后无法追溯）
    Action      string    `gorm:"size:50;index"`   // 操作类型：create, update, delete, execute, login, logout
    Module      string    `gorm:"size:50;index"`   // 模块：user, role, asset, task, deployment, auth
    ResourceID  *uint     `gorm:"index"`           // 资源 ID（可选）
    ResourceName string   `gorm:"size:200"`        // 资源名称（冗余存储）
    Detail      string    `gorm:"type:text"`       // 操作详情（JSON 格式）
    IPAddress   string    `gorm:"size:50"`         // 客户端 IP
    UserAgent   string    `gorm:"size:500"`        // 用户代理
    Status      string    `gorm:"size:20"`         // 操作结果：success, failed
    ErrorMsg    string    `gorm:"type:text"`       // 错误信息（失败时）
    CreatedAt   time.Time `gorm:"index"`           // 操作时间
}
```

**决策理由**:
- 冗余存储 `Username` 和 `ResourceName`，确保即使关联数据被删除也能追溯
- `Detail` 使用 JSON 格式存储，灵活记录不同操作的详细信息
- 多个索引字段支持高效查询

### 2. 审计记录方式

**方案 A**: 中间件自动记录（选用）
- 在 API 层通过中间件统一记录
- 优点：统一入口，不遗漏
- 缺点：需要配置哪些路由需要审计

**方案 B**: Service 层手动调用
- 在每个 Service 方法中手动调用审计服务
- 优点：精确控制
- 缺点：容易遗漏，代码侵入性强

**决策**: 采用方案 A，通过配置定义需要审计的路由和操作类型。

### 3. 审计配置

```yaml
audit:
  enabled: true
  retention_days: 90  # 日志保留天数
  routes:
    - path: "/api/v1/auth/login"
      method: "POST"
      module: "auth"
      action: "login"
    - path: "/api/v1/users"
      method: "POST"
      module: "user"
      action: "create"
    # ... 更多路由配置
```

### 4. 敏感信息处理

- 密码字段自动过滤，不记录到 `Detail`
- SSH 私钥内容不记录
- 使用白名单机制，只记录允许的字段

### 5. 查询 API 设计

```
GET /api/v1/audit-logs
  ?page=1
  &page_size=20
  &user_id=1
  &username=admin
  &module=user
  &action=create
  &status=success
  &start_time=2026-01-01T00:00:00Z
  &end_time=2026-01-31T23:59:59Z
  &keyword=admin

GET /api/v1/audit-logs/export
  ?format=csv|json
  &... (同上筛选参数)
```

## Risks / Trade-offs

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 日志量过大导致存储压力 | 中 | 配置保留天数，定期清理过期日志 |
| 审计中间件影响 API 性能 | 低 | 异步写入审计日志 |
| 敏感信息泄露 | 高 | 严格的字段过滤白名单 |

## Migration Plan

1. 创建 `audit_logs` 数据库表
2. 部署审计中间件（初始不启用）
3. 配置需要审计的路由
4. 启用审计功能
5. 验证审计日志记录正确性

## Open Questions

1. 是否需要支持审计日志的归档功能？（建议后续版本实现）
2. 是否需要支持审计日志的实时推送？（建议后续版本实现）
3. WebSSH 操作的审计粒度如何定义？（建议只记录连接/断开事件）
