# Change: 模板导入时已存在则覆盖

## Why

当前任务模板导入逻辑为：若导入的模板名称与已有模板同名，则**跳过**该条并计入 `skipped`。用户希望改为**已存在的覆盖**：同名模板用导入内容更新已有记录，便于通过重复导入同一文件同步/更新模板内容。

## What Changes

- **导入行为**：当导入的模板名称在系统中已存在时，使用导入的 `name`、`description`、`content`、`type` 更新该条记录（覆盖），不再跳过。
- **导入结果**：在 `ImportResult` 中增加「更新」统计（如 `updated` 数量及可选 `updated_items` 列表），与「新增」「失败」区分；原「跳过」逻辑移除后，`skipped` 可保留为空或废弃。
- **前端**：导入结果弹窗中展示「更新」数量及被更新模板名称（若有），与成功/失败并列。

## Impact

- **受影响规格**：task-template-import-export（Task Template Import 相关需求）
- **受影响代码**：`backend/internal/service/task/service.go` 中 `ImportTemplates`（同名时改为 Update 并计入 updated）；`ImportResult` 结构增加 `updated`、`updated_items`（可选）；`frontend/src/api/execution.ts` 中 `ImportResult` 类型；`frontend/src/pages/executions/TemplateList.tsx` 中导入结果展示
