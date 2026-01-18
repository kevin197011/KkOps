# Tasks: 简化授权模型为纯角色授权

## 1. 后端授权服务简化

- [x] 1.1 修改 `HasAssetAccess` 移除直接用户授权检查
- [x] 1.2 修改 `HasMultipleAssetAccess` 移除直接用户授权检查
- [x] 1.3 修改 `GetUserAssetIDs` 只返回角色授权的资产
- [x] 1.4 移除 `GetUserAssets`、`GrantAssetAccess`、`RevokeAssetAccess` 等直接授权方法

## 2. 后端 API 路由清理

- [x] 2.1 从 `main.go` 移除用户资产授权路由 (`/users/:id/assets`)
- [x] 2.2 删除 `backend/internal/handler/user/assets.go` 文件
- [x] 2.3 从 `main.go` 移除资产用户授权路由 (`/assets/:id/users`)
- [x] 2.4 删除 `backend/internal/handler/asset/users.go` 文件

## 3. 前端用户管理页面简化

- [x] 3.1 从 `UserList.tsx` 移除资产授权相关状态
- [x] 3.2 从 `UserList.tsx` 移除资产授权相关方法
- [x] 3.3 从 `UserList.tsx` 移除"授权"按钮
- [x] 3.4 从 `UserList.tsx` 移除资产授权弹窗

## 4. 前端 API 清理

- [x] 4.1 删除 `frontend/src/api/userAsset.ts` 文件

## 5. 验证测试

- [x] 5.1 验证用户管理页面只有角色分配功能
- [x] 5.2 验证角色资产授权功能正常
- [x] 5.3 验证用户通过角色获得资产访问权限
- [x] 5.4 验证 WebSSH 权限检查正常
- [x] 5.5 验证任务执行权限检查正常
