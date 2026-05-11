# Design: 整合用户角色与授权管理

## Context

当前 KkOps 系统的权限模型存在以下问题：

1. **用户-角色关联缺失**：后端已实现 `user_roles` 关联表和 `AssignRoleToUser` API，但前端用户管理界面没有角色分配功能
2. **授权模型不完整**：
   - 当前只有"用户→资产"的直接授权（`user_assets` 表）
   - 缺少"角色→资产"的授权关系
   - 管理员角色（`is_admin=true`）可以访问所有资产，但普通角色无法批量授权资产
3. **操作繁琐**：每个用户都需要单独授权资产，无法通过角色批量管理

## Goals / Non-Goals

### Goals
- 实现用户-角色的前端管理功能
- 实现角色-资产的授权关系
- 统一授权检查逻辑：用户可访问资产 = 管理员角色 OR 直接授权 OR 角色授权

### Non-Goals
- 不改变现有的 `user_assets` 直接授权机制（保持兼容）
- 不实现复杂的权限继承（如角色层级）
- 不实现细粒度的操作权限（如只读、读写）

## Decisions

### Decision 1: 新增 `role_assets` 表实现角色级别授权

**选择**：创建 `role_assets` 关联表存储角色与资产的授权关系

**原因**：
- 与现有 `user_assets` 表结构一致，易于理解和维护
- 支持角色级别的批量授权
- 不影响现有的直接授权机制

**替代方案**：
- 使用 JSON 字段存储资产 ID 列表 → 不利于查询和索引
- 使用权限表（permission）管理 → 过于复杂，不符合当前需求

### Decision 2: 授权检查顺序

**选择**：按以下顺序检查用户是否有权访问资产

```
1. 检查用户是否有管理员角色（is_admin=true）→ 有则允许访问所有资产
2. 检查 user_assets 表是否有直接授权 → 有则允许访问
3. 检查 role_assets 表是否有角色授权 → 有则允许访问
4. 拒绝访问
```

**原因**：
- 管理员优先，避免不必要的数据库查询
- 直接授权优先于角色授权，支持特殊情况的单独授权
- 角色授权作为批量授权的补充

### Decision 3: 前端 UI 设计

**选择**：
- 用户管理页面增加"角色"列和"角色分配"按钮
- 角色管理页面增加"授权资产"列和"资产授权"按钮
- 复用现有的 Transfer 组件进行资产选择

**原因**：
- 与现有的资产授权 UI 保持一致
- 用户熟悉的操作模式
- 减少开发工作量

## Data Model

### 新增表：role_assets

```sql
CREATE TABLE role_assets (
    role_id BIGINT NOT NULL,
    asset_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    PRIMARY KEY (role_id, asset_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

CREATE INDEX idx_role_assets_role_id ON role_assets(role_id);
CREATE INDEX idx_role_assets_asset_id ON role_assets(asset_id);
```

### Go Model

```go
// RoleAsset represents the many-to-many relationship between roles and assets for authorization
type RoleAsset struct {
    RoleID    uint      `gorm:"primaryKey" json:"role_id"`
    AssetID   uint      `gorm:"primaryKey" json:"asset_id"`
    CreatedAt time.Time `json:"created_at"`
    CreatedBy uint      `json:"created_by"` // The user who granted this authorization
}
```

## API Design

### 用户角色管理

```
GET    /api/v1/users/:id/roles          # 获取用户的角色列表
POST   /api/v1/users/:id/roles          # 为用户分配角色 { role_ids: [1, 2] }
DELETE /api/v1/users/:id/roles/:role_id # 移除用户的某个角色
```

### 角色资产授权

```
GET    /api/v1/roles/:id/assets         # 获取角色授权的资产列表
POST   /api/v1/roles/:id/assets         # 为角色授权资产 { asset_ids: [1, 2, 3] }
DELETE /api/v1/roles/:id/assets/:asset_id # 撤销角色对某个资产的授权
DELETE /api/v1/roles/:id/assets         # 批量撤销 { asset_ids: [1, 2] }
```

## Authorization Service Changes

```go
// GetUserAuthorizedAssetIDs returns all asset IDs a user can access
func (s *Service) GetUserAuthorizedAssetIDs(userID uint) ([]uint, error) {
    // 1. Check if user has admin role
    isAdmin, err := s.IsAdmin(userID)
    if err != nil {
        return nil, err
    }
    if isAdmin {
        // Return all asset IDs
        var allAssetIDs []uint
        s.db.Model(&model.Asset{}).Select("id").Find(&allAssetIDs)
        return allAssetIDs, nil
    }

    // 2. Get directly authorized assets (user_assets)
    var directAssetIDs []uint
    s.db.Model(&model.UserAsset{}).
        Where("user_id = ?", userID).
        Pluck("asset_id", &directAssetIDs)

    // 3. Get role-authorized assets (role_assets through user_roles)
    var roleAssetIDs []uint
    s.db.Model(&model.RoleAsset{}).
        Joins("JOIN user_roles ON user_roles.role_id = role_assets.role_id").
        Where("user_roles.user_id = ?", userID).
        Pluck("role_assets.asset_id", &roleAssetIDs)

    // 4. Merge and deduplicate
    assetIDSet := make(map[uint]bool)
    for _, id := range directAssetIDs {
        assetIDSet[id] = true
    }
    for _, id := range roleAssetIDs {
        assetIDSet[id] = true
    }

    result := make([]uint, 0, len(assetIDSet))
    for id := range assetIDSet {
        result = append(result, id)
    }
    return result, nil
}
```

## Risks / Trade-offs

### Risk 1: 性能影响
- **风险**：授权检查需要查询多张表
- **缓解**：
  - 管理员检查优先，减少后续查询
  - 添加适当的索引
  - 考虑缓存用户权限（未来优化）

### Risk 2: 授权冲突
- **风险**：用户直接授权和角色授权可能产生混淆
- **缓解**：
  - UI 清晰展示授权来源（直接/角色）
  - 文档说明授权优先级

## Migration Plan

1. 创建 `role_assets` 表（数据库迁移）
2. 更新授权服务，支持角色授权检查
3. 添加角色资产授权 API
4. 添加用户角色管理 API
5. 更新前端用户管理页面
6. 更新前端角色管理页面
7. 测试验证

## Open Questions

1. 是否需要在资产列表中显示授权来源（直接授权 vs 角色授权）？
   - 建议：暂不实现，保持简洁
2. 是否需要支持角色授权的批量操作（按项目/环境）？
   - 建议：复用现有的批量选择功能
