# Implementation Tasks

## 1. 数据库设计
- [x] 1.1 创建部署配置表（deployment_configs）
- [x] 1.2 创建部署执行表（deployments）
- [x] 1.3 创建部署版本表（deployment_versions）
- [x] 1.4 创建部署主机关联表（deployment_hosts）- 使用数组字段存储
- [x] 1.5 更新数据库迁移脚本

## 2. 后端模型
- [x] 2.1 创建 DeploymentConfig 模型
- [x] 2.2 创建 Deployment 模型
- [x] 2.3 创建 DeploymentVersion 模型
- [x] 2.4 创建 DeploymentHost 模型（关联表）- 使用数组字段存储
- [x] 2.5 在 database.go 中添加模型到 AutoMigrate

## 3. 后端 Repository
- [x] 3.1 创建 DeploymentConfigRepository
- [x] 3.2 创建 DeploymentRepository
- [x] 3.3 创建 DeploymentVersionRepository

## 4. 后端 Service
- [x] 4.1 创建 DeploymentConfigService
- [x] 4.2 创建 DeploymentService（集成 Salt API）
- [x] 4.3 创建 DeploymentVersionService
- [x] 4.4 实现部署执行逻辑（调用 Salt state.apply）
- [x] 4.5 实现部署状态查询（查询 Salt job 状态）

## 5. 后端 Handler
- [x] 5.1 创建 DeploymentConfigHandler
- [x] 5.2 创建 DeploymentHandler
- [x] 5.3 创建 DeploymentVersionHandler
- [x] 5.4 在 main.go 中注册路由

## 6. 前端 Service
- [x] 6.1 创建部署管理 Service（deployment.ts）

## 7. 前端页面
- [x] 7.1 创建部署管理页面（Deployments.tsx）
- [x] 7.2 实现部署配置管理（列表、创建、编辑、删除）
- [x] 7.3 实现部署执行界面（选择配置、选择主机、执行部署）
- [x] 7.4 实现部署历史查询（列表、详情、过滤）
- [x] 7.5 实现部署回滚功能

## 8. 菜单和路由
- [x] 8.1 在 MainLayout 中添加部署管理菜单项
- [x] 8.2 在 App.tsx 中添加路由配置
