# Change: 添加运维导航功能

## Why
运维团队在日常工作中经常需要使用各种运维工具（如监控系统 Grafana、日志系统 ELK、CI/CD 工具等）。目前这些工具分散在不同的地方，访问不便。在仪表盘页面增加"运维导航"功能，可以集中管理和快速访问这些常用运维工具，提高工作效率。

## What Changes
- 在仪表盘页面下方添加"运维导航"区域
- 支持添加运维工具链接（外部链接，新标签页打开）
- 支持工具分类（如：监控工具、日志工具、CI/CD 工具、数据库工具等）
- 每个工具支持配置名称、描述、图标、URL
- 工具列表支持管理员在后台配置和管理
- 工具卡片支持点击打开链接，支持悬停显示描述信息
- 支持工具搜索和过滤功能（可选）

## Impact
- Affected specs: 新增 `operations-navigation` capability
- Affected code: 
  - Frontend: `frontend/src/pages/Dashboard.tsx`（添加运维导航组件）
  - Backend: 新增工具管理 API（`backend/internal/handler/operation-tool/handler.go`）
  - Database: 新增 `operation_tools` 表存储工具配置
