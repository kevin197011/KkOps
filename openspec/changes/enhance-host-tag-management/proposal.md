# Change: Enhance Host Tag Management

## Why
当前主机标签管理功能已实现基本的 CRUD 操作，但缺少一些实用功能来提升用户体验和管理效率：
- 无法快速搜索标签
- 无法查看使用某个标签的所有主机
- 无法批量删除标签
- 缺少常用颜色预设，每次都需要手动选择颜色
- 标签名称没有唯一性校验提示

## What Changes
- **添加搜索功能**: 在标签列表页面添加实时搜索框
- **添加标签使用详情**: 点击使用数量可查看使用该标签的所有主机列表
- **添加批量删除**: 支持选择多个标签进行批量删除
- **添加颜色预设**: 提供常用颜色的快捷选择按钮
- **添加名称唯一性校验**: 创建/编辑时检查标签名称是否已存在

## Impact
- **Affected specs**: 
  - `host-tag-management` (新增 - 主机标签管理能力)
- **Affected code**: 
  - 前端: `frontend/src/pages/HostTags.tsx` - 增强标签管理页面
  - 前端: `frontend/src/services/hostTag.ts` - 可能需要添加新的 API 调用
  - 后端: 可能需要添加标签名称唯一性检查 API
