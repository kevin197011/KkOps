## 1. Backend - 数据模型

- [x] 1.1 创建 `internal/model/deployment.go`，定义 `DeploymentModule` 和 `Deployment` 模型
- [x] 1.2 创建数据库迁移脚本 `migrations/add_deployment_management.sql`
- [x] 1.3 更新 `internal/config/database.go`，添加模型自动迁移

## 2. Backend - 服务层

- [x] 2.1 创建 `internal/service/deployment/service.go`，实现模块 CRUD
- [x] 2.2 实现版本获取功能：从 HTTP JSON 数据源获取版本列表
- [x] 2.3 实现部署执行功能：版本变量替换、SSH 执行脚本
- [x] 2.4 实现部署历史查询和状态更新

## 3. Backend - Handler 层

- [x] 3.1 创建 `internal/handler/deployment/handler.go`，实现模块管理 API
- [x] 3.2 实现版本获取 API `/deployment-modules/:id/versions`
- [x] 3.3 实现部署执行 API `/deployment-modules/:id/deploy`
- [x] 3.4 实现部署历史 API `/deployments`

## 4. Backend - 路由注册

- [x] 4.1 更新 `cmd/server/main.go`，注册部署管理路由

## 5. Frontend - API 层

- [x] 5.1 创建 `api/deployment.ts`，定义接口类型和 API 方法

## 6. Frontend - 页面组件

- [x] 6.1 创建 `pages/deployments/DeploymentModuleList.tsx` - 模块管理页面
- [x] 6.2 创建模块编辑弹窗组件
- [x] 6.3 创建部署执行弹窗组件（版本选择、主机选择）
- [x] 6.4 创建部署历史列表组件
- [x] 6.5 创建部署详情/日志查看组件

## 7. Frontend - 路由和菜单

- [x] 7.1 更新 `App.tsx`，添加部署管理路由
- [x] 7.2 更新 `layouts/MainLayout.tsx`，添加部署管理菜单项

## 8. 验证测试

- [x] 8.1 测试模块 CRUD 功能
- [x] 8.2 测试版本数据源获取
- [x] 8.3 测试部署执行和日志记录
- [x] 8.4 测试部署历史查询
