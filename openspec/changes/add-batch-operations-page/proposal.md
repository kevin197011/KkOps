# Change: Add Batch Operations Page

## Why
当前系统虽然有 Salt 命令执行功能，但缺少一个统一的批量操作界面。运维人员需要频繁执行批量操作（如批量查看主机状态、批量执行命令、批量查看系统信息等），但现有的 Salt API 调用方式不够直观和便捷。需要提供一个自助作业页面，让运维人员能够：
- 通过可视化界面选择目标主机（支持按项目、分组、标签筛选）
- 选择常用的 Salt 命令或自定义命令
- 批量执行操作并实时查看结果
- 保存常用操作模板以便重复使用
- 查看操作历史和结果

## What Changes
- **新增批量操作页面**: 创建独立的批量操作管理页面 (`/batch-operations`)
- **主机选择器**: 支持多选主机，可按项目、分组、标签、状态筛选
- **命令模板管理**: 支持创建、保存、使用常用命令模板（如查看磁盘、查看内存、查看进程等）
- **批量命令执行**: 基于 Salt API 执行批量命令，支持实时查看执行进度和结果
- **结果展示**: 以表格形式展示每个主机的执行结果，支持结果导出
- **操作历史**: 记录批量操作历史，支持查看历史执行结果
- **权限控制**: 基于 RBAC 控制批量操作的权限

## Impact
- **Affected specs**: 
  - `batch-operations` (新增 - 批量操作能力)
- **Affected code**: 
  - 数据库: 创建批量操作历史表 (`batch_operations`)、命令模板表 (`command_templates`)
  - 后端: 创建 BatchOperation, CommandTemplate 模型和相关 Repository、Service、Handler
  - 后端: 扩展 Salt Handler 支持批量操作和异步任务跟踪
  - 前端: 创建批量操作页面 (`BatchOperations.tsx`)
  - 前端: 创建主机选择器组件 (`HostSelector.tsx`)
  - 前端: 创建命令模板管理组件 (`CommandTemplateManager.tsx`)
  - 前端: 创建批量操作服务 (`batchOperations.ts`)
- **User experience**: 运维人员可以通过可视化界面快速执行批量操作，提高运维效率

