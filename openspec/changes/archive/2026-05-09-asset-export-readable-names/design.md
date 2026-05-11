# Design: 资产导出可读名称

## Context

当前 `ExportAssets` 在 CSV 中写入 Project ID、Environment ID、SSH Key ID 的数值；Cloud Platform 已写入名称但表头仍为「Cloud Platform ID」。用户需要导出文件中直接看到项目名、环境名、SSH 密钥名，便于手工查看与对照。

## Goals / Non-Goals

### Goals
- 导出 CSV 表头与单元格内容统一为可读名称：Project、Cloud Platform、Environment、SSH Key。
- 关联为空时导出为空字符串，不写 ID。
- 导入仍支持现有「Project ID / Environment ID / SSH Key ID」列；可选支持「Project / Environment / SSH Key」名称列并解析为 ID，实现导出文件的 round-trip。

### Non-Goals
- 不改变资产列表 API 或前端列表页的列展示逻辑。
- 不改变导入的必填列（仍以 Hostname 等现有规则为准）。

## Decisions

### 1. 导出实现

- **Preload**：在 `ExportAssets` 中除现有 `Preload("Environment")`、`Preload("CloudPlatform")` 外，增加 `Preload("Project")`、`Preload("SSHKey")`，以便在循环中取 `asset.Project.Name`、`asset.Environment.Name`、`asset.SSHKey.Name`。
- **表头**：改为 `"ID", "Hostname", "Project", "Cloud Platform", "Environment", "IP", "SSH Port", "SSH Key", "SSH User", "CPU", "Memory", "Disk", "Status", "Description"`。
- **行数据**：Project / Cloud Platform / Environment / SSH Key 四列写入关联实体的 Name；指针为 nil 时写空字符串。Cloud Platform 当前已写名称，仅需与上述表头统一。

### 2. 导入兼容

- **保留现有列**：解析时继续支持 "Project ID"、"Environment ID"、"SSH Key ID"、"Cloud Platform ID" 等数字 ID 列，保证旧版导出或手工编辑的 CSV 仍可导入。
- **可选名称列**：若 CSV 中存在 "Project"、"Environment"、"SSH Key"、"Cloud Platform" 列，则按名称查库解析为 ID 后写入请求；名称不存在时可按业务选择报错或忽略该关联。实现时需在 service 层按 name 查询 Project / Environment / SSHKey / CloudPlatform 并得到 ID。

### 3. 风险与缓解

| 风险 | 缓解 |
|------|------|
| 项目/环境/密钥重名导致解析歧义 | 按 name 查询时使用唯一约束字段（如 name 唯一）；若存在多条则取第一条或报错，在 spec 中约定行为。 |
| 导出数据量增大（名称比 ID 长） | 仅增加可读性，CSV 体积影响可接受。 |

## Migration Plan

1. 修改 `ExportAssets`：Preload Project、SSHKey；表头与行数据改为上述可读名称列。
2. 可选：在 `parseImportRecord` 中支持 "Project"、"Environment"、"SSH Key"、"Cloud Platform" 名称列，解析为 ID 后写入 CreateAssetRequest。
