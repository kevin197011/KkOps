# Tasks: 资产授权权限重构

## 1. 数据模型变更

- [x] 1.1 创建 `user_assets` 模型 (`backend/internal/model/user_asset.go`)
- [x] 1.2 修改 `Role` 模型添加 `is_admin` 字段 (`backend/internal/model/rbac.go`)
- [x] 1.3 创建数据库迁移文件 (`backend/migrations/add_asset_authorization.sql`)
- [x] 1.4 更新 GORM AutoMigrate 包含新模型

## 2. 授权服务层

- [x] 2.1 创建授权服务 `authorization` (`backend/internal/service/authorization/service.go`)
  - `IsAdmin(userID uint) bool` - 检查用户是否为管理员
  - `HasAssetAccess(userID, assetID uint) bool` - 检查用户对资产的访问权限
  - `GetUserAssets(userID uint) ([]Asset, error)` - 获取用户已授权的资产
  - `GetAssetUsers(assetID uint) ([]User, error)` - 获取资产已授权的用户
  - `GrantAssetAccess(userID uint, assetIDs []uint) error` - 授权资产
  - `RevokeAssetAccess(userID uint, assetIDs []uint) error` - 撤销授权

## 3. 修改资产服务

- [x] 3.1 修改 `ListAssets` 添加用户权限过滤参数
- [x] 3.2 修改 `GetAsset` 添加权限检查
- [x] 3.3 确保管理员可以访问所有资产

## 4. 修改资产 Handler

- [x] 4.1 修改 `ListAssets` handler 传递用户上下文
- [x] 4.2 修改 `GetAsset` handler 进行权限校验
- [x] 4.3 添加用户资产授权管理 handler (`backend/internal/handler/user/assets.go`)
  - `GetUserAssets` - 获取用户已授权资产
  - `GrantUserAssets` - 授权资产给用户
  - `RevokeUserAsset` - 撤销单个授权
  - `RevokeUserAssets` - 批量撤销授权
- [x] 4.4 添加资产用户授权管理 handler (`backend/internal/handler/asset/users.go`)
  - `GetAssetUsers` - 获取资产已授权用户
  - `GrantAssetUsers` - 授权用户访问资产

## 5. 修改 WebSSH 服务

- [x] 5.1 修改 `SSHTerminalHandler` 添加资产访问权限检查
- [x] 5.2 管理员跳过权限检查
- [x] 5.3 普通用户检查 user_assets 表

## 6. 修改任务服务

- [x] 6.1 修改 `CreateTask` 检查用户对目标资产的权限
- [x] 6.2 修改 `ExecuteTask` 确认用户权限
- [x] 6.3 修改资产选择列表只返回用户有权限的资产

## 7. 注册路由

- [x] 7.1 注册用户资产授权 API 路由
- [x] 7.2 注册资产用户授权 API 路由
- [x] 7.3 确保授权管理 API 需要管理员权限

## 8. 前端 - API 层

- [x] 8.1 添加用户资产授权 API (`frontend/src/api/userAsset.ts`)
  - `getUserAssets(userId)` - 获取用户已授权资产
  - `grantUserAssets(userId, assetIds)` - 授权资产
  - `revokeUserAsset(userId, assetId)` - 撤销授权
  - `revokeUserAssets(userId, assetIds)` - 批量撤销

## 9. 前端 - 用户管理页面

- [x] 9.1 用户编辑页面添加"资产授权"Tab
- [x] 9.2 显示用户已授权的资产列表
- [x] 9.3 添加资产选择器（支持按项目/环境筛选）
- [x] 9.4 支持批量授权和撤销
- [x] 9.5 授权变更后刷新列表

## 10. 前端 - 资产列表适配

- [x] 10.1 确保资产列表 API 调用正常（后端已过滤）
- [x] 10.2 WebSSH 连接按钮只对有权限的资产显示
- [x] 10.3 处理 403 错误提示用户无权限

## 11. 初始化数据

- [x] 11.1 数据库迁移：创建 `user_assets` 表
- [x] 11.2 数据库迁移：添加 `roles.is_admin` 字段
- [x] 11.3 设置默认 admin 角色的 `is_admin = true`
- [x] 11.4 为管理员角色的用户自动授权所有现有资产（可选）

## 12. 测试验证

- [ ] 12.1 测试管理员可以查看所有资产
- [ ] 12.2 测试普通用户只能看到授权的资产
- [ ] 12.3 测试普通用户无法访问未授权资产详情
- [ ] 12.4 测试普通用户无法 SSH 连接未授权资产
- [ ] 12.5 测试授权管理功能（授权、撤销、批量操作）
- [ ] 12.6 测试任务执行的权限检查

## Dependencies

- 任务 1.x 完成后才能进行 2.x
- 任务 2.x 完成后才能进行 3.x - 7.x（可并行）
- 任务 8.x 可与后端开发并行
- 任务 9.x - 10.x 依赖 8.x 完成
- 任务 11.x 在部署前执行
- 任务 12.x 在所有开发完成后执行
