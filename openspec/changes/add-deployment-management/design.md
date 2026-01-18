## Context

部署管理是运维平台的核心功能，需要支持：
- 项目下配置多个模块（如微服务）
- 每个模块配置版本数据源（HTTP JSON）
- 从数据源获取可用版本列表
- 选择版本执行部署脚本到目标主机

## Goals / Non-Goals

**Goals:**
- 支持配置项目下的部署模块
- 支持 HTTP JSON 版本数据源
- 支持版本选择执行部署
- 记录部署历史和状态

**Non-Goals:**
- 不实现 CI/CD 流水线
- 不实现代码构建功能
- 不实现容器编排（K8s 等）集成

## Data Model

### DeploymentModule（部署模块）

```go
type DeploymentModule struct {
    ID              uint   `gorm:"primaryKey"`
    ProjectID       uint   `gorm:"not null;index"`        // 所属项目
    Name            string `gorm:"not null;size:100"`     // 模块名称
    Description     string `gorm:"type:text"`             // 描述
    VersionSourceURL string `gorm:"size:500"`             // 版本数据源 URL
    DeployScript    string `gorm:"type:text"`             // 部署脚本内容
    ScriptType      string `gorm:"size:50;default:shell"` // shell/python
    Timeout         int    `gorm:"default:600"`           // 超时时间（秒）
    AssetIDs        string `gorm:"type:text"`             // 默认目标主机 ID（逗号分隔）
    CreatedBy       uint
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### Deployment（部署记录）

```go
type Deployment struct {
    ID              uint   `gorm:"primaryKey"`
    ModuleID        uint   `gorm:"not null;index"`       // 部署模块 ID
    Version         string `gorm:"size:100"`             // 部署版本
    Status          string `gorm:"default:pending"`      // pending/running/success/failed/cancelled
    AssetIDs        string `gorm:"type:text"`            // 执行的目标主机
    Output          string `gorm:"type:text"`            // 执行输出
    Error           string `gorm:"type:text"`            // 错误信息
    CreatedBy       uint
    StartedAt       *time.Time
    FinishedAt      *time.Time
    CreatedAt       time.Time
}
```

## Version Source Format

HTTP JSON 返回格式：
```json
{
  "versions": ["v1.0.0", "v1.0.1", "v1.1.0"],
  "latest": "v1.1.0"
}
```

支持的变量替换（在部署脚本中）：
- `${VERSION}` - 选择的版本号
- `${MODULE_NAME}` - 模块名称
- `${PROJECT_NAME}` - 项目名称

## API Design

### 模块管理
- `GET /api/v1/deployment-modules` - 列表（支持按项目筛选）
- `POST /api/v1/deployment-modules` - 创建
- `GET /api/v1/deployment-modules/:id` - 详情
- `PUT /api/v1/deployment-modules/:id` - 更新
- `DELETE /api/v1/deployment-modules/:id` - 删除

### 版本获取
- `GET /api/v1/deployment-modules/:id/versions` - 从数据源获取版本列表

### 部署执行
- `POST /api/v1/deployment-modules/:id/deploy` - 执行部署
- `GET /api/v1/deployments` - 部署历史（支持按模块筛选）
- `GET /api/v1/deployments/:id` - 部署详情
- `POST /api/v1/deployments/:id/cancel` - 取消部署

## UI Design

### 部署管理页面布局

```
+------------------------------------------+
| 部署管理                        [+ 新增模块] |
+------------------------------------------+
| 筛选: [项目选择 ▼]                          |
+------------------------------------------+
| 模块列表                                   |
| +--------------------------------------+ |
| | 模块名称 | 项目 | 版本源 | 操作         | |
| | module1 | proj | ✓配置 | 部署/编辑/删除 | |
| +--------------------------------------+ |
+------------------------------------------+
```

### 部署执行弹窗

```
+------------------------------------------+
| 部署模块: module1                          |
+------------------------------------------+
| 选择版本: [v1.1.0 ▼] (latest)             |
|          v1.1.0                          |
|          v1.0.1                          |
|          v1.0.0                          |
+------------------------------------------+
| 目标主机: [√] host1  [√] host2  [ ] host3 |
+------------------------------------------+
| [取消]                          [执行部署] |
+------------------------------------------+
```

## Risks / Trade-offs

- **风险**: 版本数据源不可用
  - **缓解**: 显示友好错误提示，支持手动输入版本

- **风险**: 部署脚本执行失败
  - **缓解**: 详细记录输出日志，支持查看执行详情

## Open Questions

- 是否需要支持部署回滚功能？（可后续迭代）
- 是否需要支持部署审批流程？（可后续迭代）
