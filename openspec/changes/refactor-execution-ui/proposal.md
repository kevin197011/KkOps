# Change: 重构批量操作和Formula部署界面为GitHub Actions风格

## Why
当前批量操作和Formula部署页面的界面风格不一致，执行记录的展示方式不够直观。用户需要一个类似GitHub Actions的工作流执行界面，能够清晰地查看执行状态、点击进入详情查看完整日志输出。

## What Changes
- 统一批量操作和Formula部署页面的界面风格
- 重构执行记录列表为GitHub Actions风格的工作流运行列表
- 添加执行详情弹窗，支持查看完整日志输出
- 创建共享的执行日志查看器组件
- 优化实时日志显示效果，采用终端风格

## Impact
- Affected specs: `batch-operations`, `formula-management`
- Affected code:
  - `frontend/src/pages/BatchOperations.tsx`
  - `frontend/src/pages/FormulaDeployment.tsx`
  - `frontend/src/components/ExecutionLogViewer.tsx` (新增)
  - `frontend/src/components/WorkflowRunList.tsx` (新增)
  - `frontend/src/components/ResultViewer.tsx` (重构)
