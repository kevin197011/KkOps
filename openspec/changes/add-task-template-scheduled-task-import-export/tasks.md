# Tasks: 任务模板和任务执行导入导出功能

## 1. 后端任务模板导入导出

- [ ] 1.1 在 `taskTemplate` Service 中添加导出方法
  - [ ] 1.1.1 创建 `ExportTemplates` 方法，返回所有模板的 JSON 格式
  - [ ] 1.1.2 定义 `ExportConfig` 和 `ExportTemplate` 结构
  - [ ] 1.1.3 处理空列表情况（返回空数组）
- [ ] 1.2 在 `taskTemplate` Service 中添加导入方法
  - [ ] 1.2.1 创建 `ImportTemplates` 方法
  - [ ] 1.2.2 创建 `PreviewImport` 方法（预览，不实际创建）
  - [ ] 1.2.3 实现导入验证（JSON 格式、必填字段）
  - [ ] 1.2.4 实现重名检测（跳过已存在的模板）
  - [ ] 1.2.5 返回导入结果（成功/失败/跳过统计）
- [ ] 1.3 在 Handler 中添加导入导出端点
  - [ ] 1.3.1 添加 `GET /api/v1/templates/export` 端点
  - [ ] 1.3.2 添加 `POST /api/v1/templates/import` 端点
  - [ ] 1.3.3 添加 `POST /api/v1/templates/import/preview` 端点（可选）
  - [ ] 1.3.4 设置正确的响应头（Content-Disposition, Content-Type）

## 2. 后端任务执行导入导出

- [ ] 2.1 在 `scheduledtask` Service 中添加导出方法
  - [ ] 2.1.1 创建 `ExportScheduledTasks` 方法
  - [ ] 2.1.2 定义 `ExportConfig` 和 `ExportScheduledTask` 结构
  - [ ] 2.1.3 将 `AssetIDs` 转换为主机名/IP 列表
  - [ ] 2.1.4 将模板 ID 转换为模板名称（如果有关联）
  - [ ] 2.1.5 处理空列表情况
- [ ] 2.2 在 `scheduledtask` Service 中添加导入方法
  - [ ] 2.2.1 创建 `ImportScheduledTasks` 方法
  - [ ] 2.2.2 创建 `PreviewImport` 方法（预览，不实际创建）
  - [ ] 2.2.3 实现导入验证（JSON 格式、必填字段、Cron 表达式格式）
  - [ ] 2.2.4 实现主机匹配（根据主机名或 IP 查找）
  - [ ] 2.2.5 实现模板匹配（根据模板名称查找，可选）
  - [ ] 2.2.6 实现重名检测（跳过已存在的任务）
  - [ ] 2.2.7 返回导入结果（成功/失败/跳过统计）
- [ ] 2.3 添加辅助方法
  - [ ] 2.3.1 `FindAssetsByHostnamesOrIPs` - 根据主机名或 IP 查找资产
  - [ ] 2.3.2 `FindTemplateByName` - 根据模板名称查找模板（或复用部署管理的方法）
  - [ ] 2.3.3 `TaskExistsByName` - 检查同名任务是否存在
- [ ] 2.4 在 Handler 中添加导入导出端点
  - [ ] 2.4.1 添加 `GET /api/v1/tasks/export` 端点
  - [ ] 2.4.2 添加 `POST /api/v1/tasks/import` 端点
  - [ ] 2.4.3 添加 `POST /api/v1/tasks/import/preview` 端点（可选）
  - [ ] 2.4.4 设置正确的响应头

## 3. 前端任务模板导入导出

- [ ] 3.1 更新 `frontend/src/api/execution.ts`
  - [ ] 3.1.1 添加 `ExportConfig` 和 `ImportResult` 接口定义
  - [ ] 3.1.2 添加 `templateApi.exportConfig()` 方法
  - [ ] 3.1.3 添加 `templateApi.importConfig()` 方法
  - [ ] 3.1.4 添加 `templateApi.previewImport()` 方法（可选）
- [ ] 3.2 更新 `frontend/src/pages/executions/TemplateList.tsx`
  - [ ] 3.2.1 添加"导出配置"按钮
  - [ ] 3.2.2 实现 `handleExport` 函数（下载 JSON 文件）
  - [ ] 3.2.3 添加"导入配置"按钮
  - [ ] 3.2.4 实现 `handleImport` 函数（文件选择、上传、显示结果）
  - [ ] 3.2.5 实现 `handlePreviewImport` 函数（预览导入结果）
  - [ ] 3.2.6 添加导入结果 Modal（显示成功/失败/跳过统计）

## 4. 前端任务执行导入导出

- [ ] 4.1 更新 `frontend/src/api/task.ts`
  - [ ] 4.1.1 添加 `ExportConfig` 和 `ImportResult` 接口定义（如果还没有）
  - [ ] 4.1.2 添加 `scheduledTaskApi.exportConfig()` 方法
  - [ ] 4.1.3 添加 `scheduledTaskApi.importConfig()` 方法
  - [ ] 4.1.4 添加 `scheduledTaskApi.previewImport()` 方法（可选）
- [ ] 4.2 更新 `frontend/src/pages/tasks/ScheduledTaskList.tsx`
  - [ ] 4.2.1 添加"导出配置"按钮
  - [ ] 4.2.2 实现 `handleExport` 函数（下载 JSON 文件）
  - [ ] 4.2.3 添加"导入配置"按钮
  - [ ] 4.2.4 实现 `handleImport` 函数（文件选择、上传、显示结果）
  - [ ] 4.2.5 实现 `handlePreviewImport` 函数（预览导入结果）
  - [ ] 4.2.6 添加导入结果 Modal（显示成功/失败/跳过统计）
  - [ ] 4.2.7 导入成功后刷新任务列表

## 5. 错误处理和验证

- [ ] 5.1 后端验证
  - [ ] 5.1.1 JSON 格式验证
  - [ ] 5.1.2 必填字段验证（任务模板：name, content, type）
  - [ ] 5.1.3 必填字段验证（任务执行：name, cron_expression, content/type）
  - [ ] 5.1.4 Cron 表达式格式验证（6 字段格式）
- [ ] 5.2 前端验证
  - [ ] 5.2.1 文件类型验证（必须是 JSON）
  - [ ] 5.2.2 文件大小限制（如 10MB）
  - [ ] 5.2.3 导入前预览验证（可选）
- [ ] 5.3 错误提示
  - [ ] 5.3.1 后端返回详细的错误信息
  - [ ] 5.3.2 前端显示导入结果（成功/失败/跳过统计）
  - [ ] 5.3.3 显示失败项的详细错误信息

## 6. 测试和验证

- [ ] 6.1 任务模板导入导出测试
  - [ ] 6.1.1 导出所有模板测试
  - [ ] 6.1.2 导入新模板测试
  - [ ] 6.1.3 导入重名模板测试（应跳过）
  - [ ] 6.1.4 导入格式错误 JSON 测试
  - [ ] 6.1.5 导入缺失必填字段测试
- [ ] 6.2 任务执行导入导出测试
  - [ ] 6.2.1 导出所有任务测试
  - [ ] 6.2.2 导入新任务测试（包含主机和模板）
  - [ ] 6.2.3 导入重名任务测试（应跳过）
  - [ ] 6.2.4 导入时主机匹配测试（部分主机不存在）
  - [ ] 6.2.5 导入时模板匹配测试（模板不存在）
  - [ ] 6.2.6 导入格式错误 JSON 测试
  - [ ] 6.2.7 导入无效 Cron 表达式测试
- [ ] 6.3 跨环境测试
  - [ ] 6.3.1 从环境 A 导出，导入到环境 B
  - [ ] 6.3.2 验证主机和模板的自动匹配
