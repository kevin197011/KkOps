# Change: 任务模板和任务执行导入导出功能

## Why

当前系统缺少任务模板和任务执行的导入导出功能，用户在以下场景中需要这些功能：

1. **配置迁移**：从一个环境迁移配置到另一个环境
2. **配置备份**：备份任务模板和定时任务配置，防止误删除
3. **批量管理**：批量导入/导出任务模板和定时任务
4. **配置共享**：在不同项目或团队间共享任务配置
5. **版本控制**：将任务配置纳入版本控制系统

参考部署管理的导入导出功能，为任务模板和任务执行提供类似的能力。

## What Changes

### 任务模板（Task Templates）导入导出

- **导出功能**：
  - 导出所有任务模板为 JSON 格式
  - 包含模板的完整信息：名称、描述、内容、类型等
  - 生成可下载的 JSON 文件

- **导入功能**：
  - 支持从 JSON 文件导入任务模板
  - 自动处理重名情况（跳过或更新）
  - 支持预览导入结果（不实际创建，仅验证）

### 任务执行（Scheduled Tasks）导入导出

- **导出功能**：
  - 导出所有定时任务为 JSON 格式
  - 包含任务的完整信息：名称、描述、Cron 表达式、脚本内容、类型、超时时间、目标主机等
  - 主机信息以主机名/IP 形式导出（而非 ID）
  - 模板关联以模板名称导出（而非 ID）

- **导入功能**：
  - 支持从 JSON 文件导入定时任务
  - 自动匹配主机（根据主机名或 IP）
  - 自动匹配模板（根据模板名称）
  - 自动处理重名情况（跳过或更新）
  - 支持预览导入结果（不实际创建，仅验证）

### 前端 UI

- **任务模板页面**（`/templates`）：
  - 添加"导出配置"和"导入配置"按钮
  - 导入时支持文件选择和预览
  - 显示导入结果统计（成功、失败、跳过）

- **任务执行页面**（`/tasks`）：
  - 添加"导出配置"和"导入配置"按钮
  - 导入时支持文件选择和预览
  - 显示导入结果统计（成功、失败、跳过）

## Impact

- **后端 API**：
  - 任务模板：新增 `/api/v1/templates/export` 和 `/api/v1/templates/import` 端点
  - 任务执行：新增 `/api/v1/tasks/export` 和 `/api/v1/tasks/import` 端点
  - 参考部署管理的实现模式

- **后端 Service**：
  - `scheduledtask.Service` 新增 `ExportScheduledTasks` 和 `ImportScheduledTasks` 方法
  - `taskTemplate.Service` 或相关服务新增 `ExportTemplates` 和 `ImportTemplates` 方法
  - 需要添加辅助方法用于匹配主机和模板

- **前端**：
  - `frontend/src/api/task.ts` 新增导入导出 API 调用
  - `frontend/src/api/execution.ts` 新增模板导入导出 API 调用
  - `frontend/src/pages/tasks/ScheduledTaskList.tsx` 添加导入导出 UI
  - `frontend/src/pages/executions/TemplateList.tsx` 添加导入导出 UI

- **兼容性**：
  - 不影响现有功能
  - 纯新增功能，无破坏性变更
