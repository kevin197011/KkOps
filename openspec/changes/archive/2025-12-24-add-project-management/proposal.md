# Change: Add Project Management

## Why
在实际运维场景中，不同的资源（主机、SSH连接、部署等）通常归属于不同的项目。为了支持多项目管理，需要引入项目管理功能，使所有运维资源都能关联到特定项目，实现资源的隔离和管理。

## What Changes
- **新增项目管理能力**: 支持项目的创建、编辑、删除和查询
- **修改主机管理**: 所有主机必须关联到项目
- **修改SSH管理**: SSH连接和会话关联到项目
- **修改发布管理**: 部署配置和执行关联到项目
- **修改定时任务**: 定时任务关联到项目
- **修改权限管理**: 支持项目级别的权限控制

## Impact
- **Affected specs**: 
  - `project-management` (新增)
  - `host-management` (修改 - 添加项目关联)
  - `ssh-management` (修改 - 添加项目关联)
  - `deployment-management` (修改 - 添加项目关联)
  - `scheduled-tasks` (修改 - 添加项目关联)
  - `permission-management` (修改 - 添加项目级权限)
- **Affected code**: 
  - 数据库: 新增项目表，所有相关表添加project_id字段
  - 后端: 所有模型、Repository、Service、Handler需要支持项目关联
  - 前端: 添加项目管理页面，所有功能页面添加项目选择

