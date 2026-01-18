# Change: 资产授权权限重构

## Why

当前系统所有认证用户都能看到和操作所有资产，缺乏细粒度的访问控制。运维平台需要：
- 管理员能查看和操作所有资产
- 普通用户只能查看和操作被授权的资产
- WebSSH 连接需要进行资产级别的权限校验

## What Changes

### 数据模型变更
- **ADDED**: `UserAsset` 关联表 - 存储用户与资产的授权关系
- **MODIFIED**: `Role` 模型 - 添加 `is_admin` 字段标识管理员角色

### 后端服务变更
- **MODIFIED**: 资产列表 API - 根据用户角色过滤可见资产
- **MODIFIED**: 资产详情 API - 校验用户是否有权限访问
- **MODIFIED**: WebSSH 连接 - 校验用户是否有权限连接该资产
- **MODIFIED**: 任务执行 - 校验用户是否有权限在目标资产执行
- **ADDED**: 资产授权管理 API - 管理员为用户分配资产权限

### 前端变更
- **ADDED**: 用户资产授权管理页面
- **MODIFIED**: 资产列表 - 显示用户有权限的资产
- **MODIFIED**: WebSSH - 只能连接授权的资产

## Impact

- **Affected specs**: asset-authorization (new)
- **Affected code**:
  - `backend/internal/model/` - 新增 UserAsset 模型
  - `backend/internal/service/asset/` - 资产服务增加权限过滤
  - `backend/internal/service/rbac/` - 扩展权限检查逻辑
  - `backend/internal/handler/asset/` - 资产 handler 增加权限校验
  - `backend/internal/handler/websocket/ssh.go` - SSH 连接增加权限校验
  - `backend/internal/service/task/` - 任务执行增加权限校验
  - `frontend/src/pages/users/` - 用户管理页面增加资产授权
  - `frontend/src/api/` - 新增资产授权 API

## Database Changes

```sql
-- 新增用户资产授权表
CREATE TABLE user_assets (
    user_id BIGINT NOT NULL,
    asset_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    PRIMARY KEY (user_id, asset_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

-- 角色表添加管理员标识
ALTER TABLE roles ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;

-- 更新默认管理员角色
UPDATE roles SET is_admin = TRUE WHERE name = 'admin';
```

## Authorization Flow

```
用户请求资产
    ↓
检查用户角色 → is_admin = true? → 允许访问所有资产
    ↓ (否)
检查 user_assets 表 → 有授权记录? → 允许访问
    ↓ (否)
拒绝访问 (403 Forbidden)
```
