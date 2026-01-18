# Tasks: 整合用户角色与授权管理

## 1. 数据库迁移

- [x] 1.1 创建 `role_assets` 表迁移脚本
- [x] 1.2 添加索引优化查询性能
- [x] 1.3 更新数据库初始化代码

## 2. 后端模型

- [x] 2.1 创建 `RoleAsset` 模型 (`backend/internal/model/role_asset.go`)
- [x] 2.2 更新 `Role` 模型添加 Assets 关联
- [x] 2.3 更新 GORM AutoMigrate 包含新模型

## 3. 后端授权服务

- [x] 3.1 修改 `GetUserAuthorizedAssetIDs` 支持角色授权检查
- [x] 3.2 修改 `HasAssetAccess` 支持角色授权检查
- [x] 3.3 添加 `GetRoleAssets` 方法
- [x] 3.4 添加 `GrantRoleAssets` 方法
- [x] 3.5 添加 `RevokeRoleAssets` 方法

## 4. 后端角色服务

- [x] 4.1 添加 `GetUserRoles` 方法
- [x] 4.2 添加 `SetUserRoles` 方法（批量设置用户角色）
- [x] 4.3 修改 `ListRoles` 返回授权资产数量

## 5. 后端 API Handler

- [x] 5.1 创建角色资产授权 handler (`backend/internal/handler/role/assets.go`)
  - GET /roles/:id/assets
  - POST /roles/:id/assets
  - DELETE /roles/:id/assets/:asset_id
  - DELETE /roles/:id/assets (批量撤销)
- [x] 5.2 创建用户角色管理 handler (`backend/internal/handler/user/roles.go`)
  - GET /users/:id/roles
  - POST /users/:id/roles
  - DELETE /users/:id/roles/:role_id
- [x] 5.3 注册新路由到 main.go

## 6. 前端 API

- [x] 6.1 更新 `frontend/src/api/role.ts` 添加角色资产授权 API
- [x] 6.2 更新 `frontend/src/api/user.ts` 添加用户角色管理 API

## 7. 前端用户管理

- [x] 7.1 用户列表增加"角色"列显示用户角色标签
- [x] 7.2 用户列表增加"角色"操作按钮
- [x] 7.3 实现角色分配弹窗（Checkbox 列表选择角色）
- [x] 7.4 实现角色分配保存逻辑

## 8. 前端角色管理

- [x] 8.1 角色列表增加"授权资产"列显示资产数量
- [x] 8.2 角色列表增加"授权"操作按钮（非管理员角色）
- [x] 8.3 实现资产授权弹窗（复用 Transfer 组件）
- [x] 8.4 实现快速授权功能（按项目/环境批量添加）
- [x] 8.5 实现资产授权保存逻辑

## 9. 验证测试

- [x] 9.1 测试用户角色分配功能
- [x] 9.2 测试角色资产授权功能
- [x] 9.3 测试组合授权检查（直接授权 + 角色授权）
- [x] 9.4 测试管理员角色访问所有资产
- [x] 9.5 测试 WebSSH 权限检查
- [x] 9.6 测试任务执行权限检查
