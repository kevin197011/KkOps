# Change: 主机资产导出显示可读名称（Project / Environment / SSH Key）

## Why

主机资产导出的 CSV 中，Project ID、Environment ID、SSH Key ID 列目前输出的是数字 ID（如 1、2），手工查看或与外部表格对照时不直观，不方便人工核对与操作。用户希望导出数据中直接显示项目名、环境名、SSH 密钥名等可读信息。

## What Changes

- **导出 CSV 列名与内容**：将导出 CSV 中与关联实体相关的列改为显示可读名称，而非 ID：
  - **Project**：列名为「Project」，值为项目名称（对应原 Project ID）；无项目时为空。
  - **Cloud Platform**：列名为「Cloud Platform」，值为云平台名称（当前实现已为名称，仅统一列名为「Cloud Platform」）。
  - **Environment**：列名为「Environment」，值为环境名称（对应原 Environment ID）；无环境时为空。
  - **SSH Key**：列名为「SSH Key」，值为 SSH 密钥名称（对应原 SSH Key ID）；无密钥时为空。
- **导入兼容**：导入时继续支持按 ID 的列（如「Project ID」「Environment ID」「SSH Key ID」）；可选增强为支持按名称列（「Project」「Environment」「SSH Key」）解析为 ID 后再写入，以便导出的 CSV 可被重新导入。

## Impact

- **受影响规格**：asset-management（资产导入/导出相关需求）
- **受影响代码**：`backend/internal/service/asset/import_export.go`（ExportAssets 的 Preload、表头与行数据；Import 解析可选增强）、前端仅使用现有导出接口与文件名，无需改列名展示逻辑（导出为下载 CSV，由后端决定列名与内容）。
