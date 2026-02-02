# Tasks: 主机资产导出可读名称

## 1. 导出 CSV 表头与可读内容

- [x] 1.1 在 `ExportAssets` 中 Preload `Project`、`SSHKey`（保留现有 `Environment`、`CloudPlatform` Preload）
- [x] 1.2 将导出 CSV 表头改为：ID, Hostname, Project, Cloud Platform, Environment, IP, SSH Port, SSH Key, SSH User, CPU, Memory, Disk, Status, Description
- [x] 1.3 每行写入 Project / Cloud Platform / Environment / SSH Key 的**名称**（关联为 nil 时写空字符串）；Cloud Platform 已为名称，仅与表头统一

## 2. 导入兼容（可选）

- [x] 2.1 导入解析时继续支持 "Project ID"、"Environment ID"、"SSH Key ID"、"Cloud Platform ID" 列（保持向后兼容）
- [x] 2.2 若存在 "Project"、"Environment"、"SSH Key"、"Cloud Platform" 列，则按名称解析为 ID 后写入 CreateAssetRequest；名称不存在时记录错误或跳过该关联（在 spec 中约定）

## 3. 验收

- [x] 3.1 导出资产 CSV，确认 Project、Environment、SSH Key、Cloud Platform 列显示名称而非 ID
- [x] 3.2 可选：用导出的 CSV 再导入，确认名称列可正确解析为 ID 并创建/更新资产
