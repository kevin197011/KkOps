# Change: Add Deployment Management

## Why
发布管理是运维平台的核心功能之一。当前系统缺少发布管理能力，无法实现应用的自动化部署。需要实现完整的发布管理功能，包括部署配置管理、部署执行、版本管理、回滚等，基于 Salt 实现应用的自动化部署流程。

## What Changes
- **新增部署配置管理**: 支持创建、编辑、删除部署配置，定义应用部署规则
- **新增部署执行功能**: 基于 Salt 执行部署，支持多主机批量部署
- **新增部署版本管理**: 跟踪应用版本和部署历史
- **新增部署回滚功能**: 支持回滚到之前的版本
- **新增部署状态监控**: 实时监控部署进度和状态

## Impact
- **Affected specs**: 
  - `deployment-management` (新增 - 从归档中恢复并实现)
- **Affected code**: 
  - 数据库: 创建部署配置、部署执行、部署版本相关表
  - 后端: 创建 DeploymentConfig, Deployment, DeploymentVersion 模型和相关 Repository、Service、Handler
  - 后端: 集成 Salt API 执行部署
  - 前端: 创建部署管理页面（配置管理、部署执行、历史查询）

