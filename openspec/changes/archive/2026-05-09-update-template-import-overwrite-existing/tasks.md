# Tasks: 模板导入时已存在则覆盖

## 1. 后端导入逻辑

- [x] 1.1 在 `ImportResult` 中增加字段 `updated`（int）、`updated_items`（[]string）
- [x] 1.2 修改 `ImportTemplates`：当按 name 查到已存在模板时，执行更新（description、content、type）而非跳过，并累计 `result.Updated`、`result.UpdatedItems`
- [x] 1.3 移除「同名则跳过」逻辑及对 `result.Skipped` 的写入（或保留 Skipped 字段但不再添加同名跳过项）

## 2. 前端类型与展示

- [x] 2.1 在 `frontend/src/api/execution.ts` 的 `ImportResult` 中增加 `updated?: number`、`updated_items?: string[]`
- [x] 2.2 在 `TemplateList.tsx` 导入结果弹窗中展示「更新」数量及「已更新」列表（与成功、失败并列）

## 3. 验收

- [x] 3.1 导入包含与已有模板同名的 JSON 时，该模板内容被更新，且导入结果中显示为「更新」而非「跳过」
- [x] 3.2 导入结果中 success、updated、failed 数量与实际情况一致
