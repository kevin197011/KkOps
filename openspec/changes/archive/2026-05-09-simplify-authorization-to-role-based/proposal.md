# Change: 简化授权模型为纯角色授权

## Why

当前系统存在两套资产授权机制：
1. **直接用户授权** (`user_assets` 表) - 在用户管理页面为用户直接分配资产
2. **角色授权** (`role_assets` 表) - 在角色管理页面为角色分配资产

这种双重授权模型存在以下问题：
- **管理复杂**：需要在两个地方管理授权，容易混淆
- **维护困难**：用户离职或角色变更时，需要同时处理两处授权
- **权限不清晰**：难以追溯用户的权限来源

用户希望简化为**纯角色授权模型**：
- 用户管理只负责**绑定角色**
- 角色管理负责**绑定资产**
- 用户通过角色获得资产访问权限

## What Changes

### 移除的功能
- **REMOVED**: 用户管理页面的"授权"按钮和资产授权弹窗
- **REMOVED**: 用户直接资产授权相关的后端 API
- **REMOVED**: `user_assets` 表（可选，保留数据迁移）

### 保留的功能
- ✅ 用户管理：用户 CRUD + 角色分配
- ✅ 角色管理：角色 CRUD + 资产授权
- ✅ 权限检查：通过角色关联检查资产访问权限

### 后端变更
- **MODIFIED**: 授权服务 - 移除直接用户授权检查，只保留角色授权检查
- **REMOVED**: 用户资产授权 API (`/users/:id/assets`)
- **REMOVED**: `user_assets` 模型和相关服务

### 前端变更
- **REMOVED**: `UserList.tsx` 中的资产授权功能
- **REMOVED**: `frontend/src/api/userAsset.ts` API 文件
- **MODIFIED**: `UserList.tsx` - 只保留角色分配功能

## Impact

- **Affected specs**: role-based-authorization (modified)
- **Affected code**:
  - `backend/internal/service/authorization/service.go` - 简化授权检查
  - `backend/internal/handler/user/assets.go` - 移除
  - `backend/cmd/server/main.go` - 移除用户资产路由
  - `frontend/src/pages/users/UserList.tsx` - 移除授权功能
  - `frontend/src/api/userAsset.ts` - 移除

## Authorization Flow (Simplified)

```
用户请求资产
    ↓
检查用户角色 → is_admin = true? → 允许访问所有资产
    ↓ (否)
检查 role_assets 表（通过 user_roles 关联）→ 有角色授权? → 允许访问
    ↓ (否)
拒绝访问 (403 Forbidden)
```

## UI Changes

### 用户管理页面（简化后）
```
| ID | 用户名 | 邮箱           | 角色              | 状态 | 操作              |
|----|--------|----------------|-------------------|------|-------------------|
| 1  | admin  | admin@test.com | [管理员]          | 启用 | 角色 编辑 删除    |
| 2  | ops01  | ops@test.com   | [运维]            | 启用 | 角色 编辑 删除    |
```

### 角色管理页面（保持不变）
```
| ID | 角色名称 | 管理员    | 描述     | 授权资产 | 操作              |
|----|----------|-----------|----------|----------|-------------------|
| 1  | admin    | ✅ 管理员 | 系统管理员 | 全部   | 编辑 删除         |
| 2  | ops      | -         | 运维人员   | 15 台  | 授权 编辑 删除    |
```

## Migration Note

如果数据库中已有 `user_assets` 数据，建议：
1. 将现有的用户直接授权迁移为角色授权
2. 或者为每个有直接授权的用户创建对应的角色
3. 迁移完成后可删除 `user_assets` 表
