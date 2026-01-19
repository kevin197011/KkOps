# Change: 添加运维导航管理页面

## Why
目前运维导航工具只能在仪表盘页面查看，管理员需要通过 API 直接管理工具。为了提供更好的用户体验，需要创建一个独立的管理页面，允许管理员在界面中方便地添加、编辑、更新和删除运维工具。

## What Changes
- 在菜单中添加"运维导航"菜单项（位于仪表盘下方）
- 创建运维导航管理页面（`/operation-tools`）
- 实现工具列表展示（表格形式）
- 实现添加工具功能（弹窗表单）
- 实现编辑/更新工具功能（弹窗表单）
- 实现删除工具功能（确认对话框）
- 保留仪表盘页面的运维导航展示功能（仅查看）
- 根据权限控制菜单显示和功能按钮

## Impact
- Affected specs: 修改 `operations-navigation` capability
- Affected code:
  - Frontend: 
    - `frontend/src/layouts/MainLayout.tsx`（添加菜单项）
    - `frontend/src/App.tsx`（添加路由）
    - 新建 `frontend/src/pages/operationTools/OperationToolList.tsx`（管理页面）
  - Backend: 无需修改（API 已存在）
