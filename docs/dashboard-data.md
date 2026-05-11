# 仪表盘数据梳理

仪表盘展示的数据来自接口 `GET /api/v1/dashboard/stats`，用于一屏内呈现系统核心概览。

## 数据维度

| 维度 | 字段 | 说明 |
|------|------|------|
| **资源规模** | total_assets, total_users, total_tasks, total_projects | 资产、用户、任务、项目总数 |
| **资产状态** | assets_by_status | 按状态分组数量（如 active / disabled） |
| **任务执行** | task_execution_stats | 总次数及成功/失败/运行中/待执行/已取消 |
| **资产分布** | assets_by_project, assets_by_environment | 按项目、按环境的资产数量 |
| **最近活动** | recent_activities | 最近任务执行记录（标题、主机、状态、时间） |

## 展示规划

1. **顶部**：标题「仪表盘」+ 简短说明。
2. **资源规模**：4 个指标卡，横向排列。
3. **资产概览**：资产状态（可选）+ 按项目 / 按环境 Top N，两列或单卡。
4. **任务执行**：执行统计（成功/失败/运行中、健康度）+ 最近执行精简列表。
5. **一屏内**：使用 flex 布局控制高度，避免整页滚动。

## 后端数据来源

- 基础统计：Asset / User / Task / Project 表 count。
- 资产状态：Asset 按 status 分组统计。
- 任务执行：TaskExecution 按 status 统计，最近活动来自 TaskExecution 最近 10 条。
- 资产分布：Asset 关联 Project、Environment 分组统计。
