# Change: 新增部署管理功能

## Why

目前任务管理只支持执行固定脚本，缺乏版本化部署能力。需要新增部署管理模块，支持配置多个项目的多个模块部署任务，并能从 HTTP JSON 数据源获取版本列表，选择指定版本执行部署。

## What Changes

- 新增部署管理菜单和页面
- 新增 DeploymentModule（部署模块）数据模型，关联项目
- 新增 VersionSource（版本数据源）配置，支持 HTTP JSON 格式
- 新增 Deployment（部署记录）模型，记录执行历史
- 新增版本获取 API，从配置的 HTTP 地址拉取版本列表
- 新增部署执行 API，支持选择版本和目标主机执行部署脚本
- 前端部署管理页面：模块配置、版本选择、部署执行、历史记录

## Impact

- Affected specs: deployment-management (新增)
- Affected code:
  - Backend: `internal/model/deployment.go`, `internal/service/deployment/`, `internal/handler/deployment/`
  - Frontend: `pages/deployments/`, `api/deployment.ts`
  - Database: 新增 `deployment_modules`, `deployments` 表
