# Design: 菜单功能权限授权系统

## Context

当前系统已有：
1. **RBAC 权限模型**：`Permission`（resource/action）+ `RolePermission` 关联表
2. **权限检查中间件**：`RequirePermission` 中间件，但未在所有路由上使用
3. **资产授权**：基于角色的资产访问控制（`role_assets` 表）
4. **管理员检查**：`IsAdmin` 字段标识管理员角色

缺失：
- 菜单功能的权限定义和初始化
- 路由上的权限检查
- 前端菜单权限控制
- 用户权限获取 API

## Goals / Non-Goals

### Goals
- 为所有菜单功能定义权限（resource/action）
- 在数据库初始化时创建所有权限
- 在所有路由上添加权限检查
- 前端根据权限显示/隐藏菜单
- 管理员自动拥有所有权限
- 角色管理界面支持权限分配

### Non-Goals
- 不实现细粒度的操作权限（read/create/update/delete），当前阶段统一使用 `*`
- 不实现权限继承机制
- 不实现临时权限或时间限制的权限

## Decisions

### 1. 权限命名规范

**格式**：`<resource>:<action>`

**示例**：
- `dashboard:read` - 仪表板查看
- `projects:*` - 项目管理所有操作
- `assets:*` - 资产管理所有操作（菜单权限，不影响资产访问控制）

**决策**：使用 `*` 表示所有操作，简化权限管理

### 2. 权限初始化策略

**选项 A**：数据库迁移脚本
- 优点：版本可控，易于追踪
- 缺点：需要手动维护

**选项 B**：Go 代码初始化
- 优点：集中管理，易于扩展
- 缺点：每次启动检查，性能略差

**决策**：采用选项 B（Go 代码初始化），在数据库初始化时检查并创建缺失的权限

### 3. 管理员权限处理

**选项 A**：为管理员角色分配所有权限
- 优点：权限检查逻辑统一
- 缺点：需要维护权限列表

**选项 B**：权限检查时优先检查 `IsAdmin`
- 优点：简化权限管理，不需要为管理员分配权限
- 缺点：权限检查逻辑需要两套

**决策**：采用选项 B，权限检查时优先检查 `IsAdmin`，管理员自动通过所有权限检查

### 4. 前端权限获取时机

**选项 A**：登录后立即获取
- 优点：权限信息立即可用
- 缺点：增加登录响应时间

**选项 B**：首次访问需要权限的页面时获取
- 优点：按需加载
- 缺点：可能导致页面闪烁

**决策**：采用选项 A，在登录后立即获取用户权限，存储在 Zustand Store 中

### 5. 菜单权限映射

**映射关系**：
```typescript
const menuPermissionMap: Record<string, string> = {
  '/dashboard': 'dashboard:read',
  '/projects': 'projects:*',
  '/environments': 'environments:*',
  '/cloud-platforms': 'cloud-platforms:*',
  '/assets': 'assets:*',
  '/executions': 'executions:*',
  '/templates': 'templates:*',
  '/tasks': 'tasks:*',
  '/deployments': 'deployments:*',
  '/ssh/keys': 'ssh-keys:*',
  '/users': 'users:*',
  '/roles': 'roles:*',
  '/audit-logs': 'audit-logs:read',
}
```

## Architecture

### 权限检查流程

```
用户请求 → 认证中间件 → 权限检查中间件 → Handler
                 ↓
           检查 IsAdmin
                 ↓ (否)
           检查用户权限
                 ↓ (无)
           返回 403
```

### 前端权限控制流程

```
用户登录 → 获取用户权限 → 存储到 Store → 渲染菜单
                                      ↓
                              根据权限过滤菜单项
                                      ↓
                        无权限时显示提示或隐藏
```

## Risks / Trade-offs

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 权限检查影响性能 | 中 | 使用缓存，优化查询 |
| 权限定义遗漏 | 高 | 完善的权限初始化脚本和测试 |
| 前端权限被绕过 | 高 | 后端必须进行权限验证 |
| 权限变更影响现有用户 | 中 | 保留管理员自动权限机制 |

## Implementation Plan

### Phase 1: 后端权限定义和初始化
1. 定义所有菜单权限
2. 实现权限初始化逻辑
3. 在数据库初始化时执行

### Phase 2: 路由权限检查
1. 在所有路由上添加权限检查中间件
2. 实现管理员自动通过逻辑
3. 测试所有 API 端点

### Phase 3: 前端权限获取和控制
1. 实现用户权限获取 API
2. 创建权限 Store
3. 实现菜单权限过滤
4. 实现路由守卫

### Phase 4: 角色管理界面增强
1. 在角色管理界面添加权限分配
2. 实现权限选择组件
3. 支持批量权限分配
