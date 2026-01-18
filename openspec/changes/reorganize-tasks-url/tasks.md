# Tasks: 重组任务相关 URL 路由

## 1. 后端重构 - 重命名现有 API

- [x] 1.1 重命名 `task` handler/service/model 为 `execution`
  - `handler/task` → `handler/execution`
  - `service/task` → `service/execution`
  - 保持 model 名称不变（Task, TaskExecution）
- [x] 1.2 更新 API 路由
  - `/api/v1/tasks/*` → `/api/v1/executions/*`
  - `/api/v1/task-templates/*` → `/api/v1/templates/*`
  - `/api/v1/task-executions/*` → `/api/v1/execution-records/*`
- [x] 1.3 更新 WebSocket 路由
  - `/ws/task-executions/:id/logs` → `/ws/execution-records/:id/logs`

## 2. 前端重构 - 重命名现有页面

- [x] 2.1 重命名 `pages/tasks` 目录为 `pages/executions`
- [x] 2.2 更新 API 客户端
  - `api/task.ts` → `api/execution.ts`
  - 更新所有 API 路径
- [x] 2.3 更新路由配置 (`App.tsx`)
  - `/tasks` → `/executions`
  - `/tasks/:taskId/executions` → `/executions/:executionId/history`
  - `/task-templates` → `/templates`
  - `/task-executions/:id/logs` → `/execution-records/:id/logs`
- [x] 2.4 更新菜单配置 (`MainLayout.tsx`)
  - 运维执行: `/tasks` → `/executions`
  - 任务模板: `/task-templates` → `/templates`，名称改为"执行模板"

## 3. 后端新增 - 定时任务管理

- [x] 3.1 创建定时任务数据模型
  - `model/scheduled_task.go`
  - 字段: id, name, description, cron_expression, template_id, asset_ids, enabled, last_run_at, next_run_at, created_by, created_at, updated_at
- [x] 3.2 创建定时任务数据库迁移
  - `migrations/add_scheduled_task.sql`
- [x] 3.3 创建定时任务服务
  - `service/scheduledtask/service.go`
  - CRUD 操作
  - 启用/禁用任务
  - 获取任务执行历史
- [x] 3.4 创建定时任务处理器
  - `handler/scheduledtask/handler.go`
  - API 端点注册
- [x] 3.5 实现 Cron 调度器
  - `service/scheduledtask/scheduler.go`
  - 使用 `robfig/cron/v3` 库
  - 支持在应用启动时加载所有启用的定时任务
  - 支持动态添加/删除/更新任务

## 4. 前端新增 - 任务管理页面

- [x] 4.1 创建任务管理 API 客户端
  - `api/task.ts` (新的，用于定时任务)
- [x] 4.2 创建任务管理页面
  - `pages/tasks/ScheduledTaskList.tsx`
  - 任务列表（表格展示）
  - 创建/编辑任务弹窗
  - 启用/禁用操作
  - 执行历史查看
- [x] 4.3 添加 Cron 表达式输入组件
  - 支持常用预设（每小时、每天、每周等）
  - 支持自定义 Cron 表达式
  - 显示下次执行时间预览
- [x] 4.4 更新路由配置
  - 添加 `/tasks` 路由
- [x] 4.5 更新菜单配置
  - 添加"任务管理"菜单项

## 5. 验证与测试

- [ ] 5.1 验证所有 API 路由正常工作
- [ ] 5.2 验证所有前端页面正常访问
- [ ] 5.3 验证定时任务创建和执行
- [ ] 5.4 验证 WebSocket 日志流正常
