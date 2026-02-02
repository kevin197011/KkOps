# Design: 模板导入时已存在则覆盖

## Context

用户需求：`/templates` 页面导入模板时，**已存在的覆盖下**——即同名模板不跳过，而是用导入内容更新已有记录。

## Goals / Non-Goals

### Goals
- 导入时按**模板名称**判定是否已存在；若存在则**更新**该条（description、content、type 等），并计入「更新」统计。
- 导入结果中区分：新增数、更新数、失败数；前端可展示「更新」项。

### Non-Goals
- 不改变导出格式或导出逻辑。
- 不引入「按 ID 覆盖」或「可选策略（跳过/覆盖）」；统一为「同名即覆盖」。

## Decisions

### 1. 判定「已存在」

- **依据**：按 `name` 唯一匹配（与当前实现一致，仅行为从「跳过」改为「更新」）。
- 不按导出文件中的 `id` 匹配（导出 JSON 中无 id，且跨环境 id 不可靠）。

### 2. 更新内容

- 使用导入条目的 `name`、`description`、`content`、`type` 更新已有记录。
- `created_by`、`created_at` 保留不变；`updated_at` 由 ORM 自动更新。

### 3. ImportResult 结构

- **新增字段**：`updated int`（更新条数）、`updated_items []string`（可选，被更新模板名称列表，便于前端展示）。
- **skipped**：不再用于「同名跳过」；可保留字段但导入时始终为空数组，或后续废弃，以兼容旧前端。

### 4. 前端展示

- 导入结果弹窗中增加「更新」统计及「已更新」列表（若有），与「成功」「失败」并列展示。

## Risks / Trade-offs

| 风险 | 缓解 |
|------|------|
| 用户误导入覆盖重要模板 | 导入前为预览/确认流程；必要时可后续增加「仅新增不覆盖」选项 |
| 旧客户端未传 updated 字段 | 后端返回固定包含 updated；前端做可选展示，无 updated 时仅显示 success/failed |

## Migration Plan

1. 后端 `ImportTemplates`：同名时改为 `db.Model(&existing).Updates(...)`，并累计 `result.Updated`、`result.UpdatedItems`。
2. 后端 `ImportResult` 增加 `Updated`、`UpdatedItems`。
3. 前端 `ImportResult` 类型与导入结果弹窗展示「更新」信息。
4. 兼容：若前端暂不展示 updated，仅展示 success/failed 也可接受。
