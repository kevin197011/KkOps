# Design: 资产授权权限系统

## Context

KkOps 是一个运维管理平台，需要对资产进行细粒度的访问控制：
- 管理员可以管理所有资产
- 普通用户只能操作被授权的资产
- 授权包括查看、SSH 连接、任务执行等操作

## Goals / Non-Goals

### Goals
- 实现基于用户-资产关联的细粒度权限控制
- 管理员角色自动拥有所有资产权限
- 支持批量授权和撤销
- 授权变更实时生效
- 保持现有 RBAC 系统兼容性

### Non-Goals
- 不实现基于项目/环境的自动授权（后续可扩展）
- 不实现时间限制的临时授权
- 不实现审批流程

## Decisions

### 1. 授权模型选择

**决定**: 使用独立的 `user_assets` 关联表实现用户-资产授权

**替代方案**:
- A) 在 Permission 中添加 asset_id 字段 → 复杂，与现有 RBAC 耦合过深
- B) 基于项目授权（授权项目下所有资产）→ 粒度不够细
- C) 独立关联表 ✓ → 简单、灵活、易于扩展

**理由**: 独立关联表方案实现简单，查询高效，且不影响现有权限系统。

### 2. 管理员识别方式

**决定**: 在 Role 表添加 `is_admin` 布尔字段

**替代方案**:
- A) 硬编码角色名 "admin" → 不灵活
- B) 在 User 表添加 is_admin → 与角色系统不一致
- C) Role 表添加 is_admin 字段 ✓ → 灵活，支持多个管理员角色

### 3. 权限检查位置

**决定**: 在 Service 层统一进行权限检查

**理由**: 
- Handler 层只做认证，Service 层做授权
- 统一入口便于维护
- 多个 Handler 可复用同一 Service 的权限逻辑

### 4. 授权数据结构

```go
// UserAsset 用户资产授权关联
type UserAsset struct {
    UserID    uint      `gorm:"primaryKey" json:"user_id"`
    AssetID   uint      `gorm:"primaryKey" json:"asset_id"`
    CreatedAt time.Time `json:"created_at"`
    CreatedBy uint      `json:"created_by"` // 授权操作人
}
```

## Risks / Trade-offs

### Risk 1: 大量授权记录影响查询性能
**缓解**: 
- 使用联合主键索引
- 授权查询使用缓存（Redis）
- 批量操作使用事务

### Risk 2: 管理员权限泄露
**缓解**:
- is_admin 字段只能由超级管理员修改
- 审计日志记录角色变更

### Risk 3: 授权管理工作量大
**缓解**:
- 提供批量授权功能
- 支持按项目/环境批量选择资产
- 考虑后续添加自动授权规则

## API Design

### 资产授权管理 API

```
# 获取用户已授权的资产列表
GET /api/v1/users/:id/assets

# 批量授权资产给用户
POST /api/v1/users/:id/assets
Body: { "asset_ids": [1, 2, 3] }

# 撤销用户的资产授权
DELETE /api/v1/users/:id/assets/:asset_id

# 批量撤销授权
DELETE /api/v1/users/:id/assets
Body: { "asset_ids": [1, 2, 3] }

# 获取资产已授权的用户列表
GET /api/v1/assets/:id/users

# 批量授权用户访问资产
POST /api/v1/assets/:id/users
Body: { "user_ids": [1, 2, 3] }
```

### 修改现有 API 行为

```
# 资产列表 - 根据用户权限过滤
GET /api/v1/assets
- 管理员: 返回所有资产
- 普通用户: 只返回已授权的资产

# 资产详情 - 权限检查
GET /api/v1/assets/:id
- 管理员: 允许
- 普通用户: 检查 user_assets 表

# WebSSH - 权限检查
WS /ws/ssh/connect
- 管理员: 允许连接任意资产
- 普通用户: 检查 user_assets 表
```

## Migration Plan

### Phase 1: 数据模型
1. 添加 `user_assets` 表
2. 添加 `roles.is_admin` 字段
3. 将现有 "admin" 角色设置为 is_admin=true

### Phase 2: 后端服务
1. 实现 `IsAdmin()` 检查方法
2. 实现 `HasAssetAccess()` 检查方法
3. 修改资产服务添加权限过滤
4. 修改 WebSSH 添加权限检查
5. 修改任务服务添加权限检查

### Phase 3: 前端
1. 用户管理页面添加资产授权 Tab
2. 资产详情页面显示已授权用户

### Rollback
- 删除 `user_assets` 表
- 删除 `roles.is_admin` 字段
- 回滚后端代码（权限检查）

## Open Questions

1. **是否需要支持"只读"授权？** 
   - 当前方案：授权即拥有完整操作权限
   - 后续可扩展：在 user_assets 添加 permission_level 字段

2. **是否需要支持项目级别自动授权？**
   - 当前方案：只支持单个资产授权
   - 后续可扩展：添加 user_projects 表实现项目级授权
