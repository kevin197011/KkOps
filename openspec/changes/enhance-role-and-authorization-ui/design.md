# Design: 增强角色管理和资产授权界面

## Context

KkOps 平台已实现基础的 RBAC 和资产授权功能，但 UI 交互需要优化：
- 角色的 is_admin 字段已存在于后端，但前端未展示和编辑
- 资产授权只能单个选择，缺乏批量操作的便利性

## Goals / Non-Goals

### Goals
- 角色列表直观展示管理员标识
- 创建/编辑角色时可设置 is_admin
- 资产授权支持按项目/环境批量选择
- 保持单个资产选择能力
- 界面简洁易用

### Non-Goals
- 不改变后端权限检查逻辑（已实现）
- 不添加复杂的授权规则引擎
- 不实现角色权限的细粒度配置（本次只关注 is_admin）

## Decisions

### 1. 角色管理员标识展示

**决定**: 在角色列表增加"管理员"列，使用 Tag 组件显示

**理由**: 
- 直观展示角色类型
- Tag 组件符合 Ant Design 设计规范
- 颜色区分便于快速识别

### 2. 资产授权批量选择

**决定**: 在授权弹窗顶部添加项目/环境下拉框 + "添加到授权"按钮

**替代方案**:
- A) 使用树形选择器 → 复杂，且需要大幅改造 Transfer 组件
- B) 分组展示资产 → 展示更清晰但实现复杂
- C) 下拉筛选 + 批量添加按钮 ✓ → 简单有效，增量改进

**理由**: 方案 C 实现简单，用户体验提升明显，可以快速完成开发

### 3. 授权界面交互流程

```
1. 用户选择项目（可选）
   ↓
2. 用户选择环境（可选，可依赖项目筛选）
   ↓
3. 点击"添加到授权"
   ↓
4. 符合条件的资产自动添加到已授权列表（targetKeys）
   ↓
5. 用户可以继续微调（单独移除/添加）
   ↓
6. 点击保存提交
```

## Component Changes

### 角色列表改动

```tsx
// 新增列
{
  title: '管理员',
  dataIndex: 'is_admin',
  key: 'is_admin',
  width: 100,
  render: (isAdmin: boolean) => 
    isAdmin ? <Tag color="red">管理员</Tag> : '-'
}

// 表单新增字段
<Form.Item name="is_admin" valuePropName="checked">
  <Checkbox>设为管理员角色</Checkbox>
</Form.Item>
<div style={{ color: '#888', fontSize: 12 }}>
  管理员角色可访问所有资产，无需单独授权
</div>
```

### 资产授权弹窗改动

```tsx
// 新增状态
const [selectedProjectId, setSelectedProjectId] = useState<number>()
const [selectedEnvId, setSelectedEnvId] = useState<number>()
const [projects, setProjects] = useState<Project[]>([])
const [environments, setEnvironments] = useState<Environment[]>([])

// 新增批量添加函数
const handleBatchAdd = () => {
  // 筛选符合条件的资产
  const filteredAssets = allAssets.filter(asset => {
    if (selectedProjectId && asset.project_id !== selectedProjectId) return false
    if (selectedEnvId && asset.environment_id !== selectedEnvId) return false
    return true
  })
  
  // 添加到 targetKeys
  const newKeys = new Set([...targetKeys, ...filteredAssets.map(a => String(a.id))])
  setTargetKeys([...newKeys])
}
```

## API Changes

后端 Role API 已支持 is_admin 字段，只需前端适配：

```typescript
// frontend/src/api/role.ts
export interface Role {
  id: number
  name: string
  description: string
  is_admin: boolean  // 新增
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  description?: string
  is_admin?: boolean  // 新增
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  is_admin?: boolean  // 新增
}
```

## Migration / Rollback

无数据迁移需求，纯前端 UI 增强。回滚只需还原前端代码。
