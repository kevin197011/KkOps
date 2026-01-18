# Change: 重组任务相关 URL 路由

## Why

当前 `/tasks` URL 被"运维执行"功能占用，用户需要新增"任务管理"功能用于定时任务管理。从语义上讲，`/tasks` 更适合用于任务管理，需要将现有的"运维执行"功能迁移到新的 URL。

## What Changes

### URL 路由调整

| 功能 | 当前 URL | 新 URL |
|------|----------|--------|
| 运维执行 | `/tasks` | `/executions` |
| 任务执行列表 | `/tasks/:taskId/executions` | `/executions/:executionId/history` |
| 任务执行日志 | `/task-executions/:id/logs` | `/executions/:id/logs` |
| 任务模板 | `/task-templates` | `/templates` |
| 任务管理（新增） | - | `/tasks` |

### 菜单名称调整

| 功能 | 当前名称 | 新名称 |
|------|----------|--------|
| 运维执行 | 运维执行 | 运维执行 |
| 任务模板 | 任务模板 | 执行模板 |
| 任务管理（新增） | - | 任务管理 |

### 后端 API 调整

| 功能 | 当前 API | 新 API |
|------|----------|--------|
| 任务 CRUD | `/api/v1/tasks/*` | `/api/v1/executions/*` |
| 任务模板 | `/api/v1/task-templates/*` | `/api/v1/templates/*` |
| 任务执行 | `/api/v1/task-executions/*` | `/api/v1/execution-records/*` |
| 定时任务管理（新增） | - | `/api/v1/tasks/*` |

### 新增功能：任务管理（定时任务）

- 支持创建定时任务（Cron 表达式）
- 任务启用/禁用
- 任务执行历史
- 任务执行状态监控
- 支持基于执行模板创建定时任务

## Impact

- **Affected specs**: task-management (new)
- **Affected code**: 
  - Frontend: `App.tsx`, `MainLayout.tsx`, `pages/tasks/*`
  - Backend: `cmd/server/main.go`, `handler/task/*`, `service/task/*`
- **Breaking changes**: 
  - 前端 URL 路由变更
  - 后端 API 路径变更（需要同步更新前端 API 客户端）
