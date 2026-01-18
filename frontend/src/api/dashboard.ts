// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

// 任务执行统计
export interface TaskExecutionStats {
  total: number
  success: number
  failed: number
  running: number
  pending: number
  cancelled: number
}

// 环境资产计数
export interface EnvironmentAssetCount {
  environment_id: number
  environment_name: string
  count: number
}

// 项目资产计数
export interface ProjectAssetCount {
  project_id: number
  project_name: string
  count: number
}

// 最近活动
export interface RecentActivity {
  id: number
  type: string
  title: string
  description: string
  status: string
  created_at: string
}

// 仪表板统计数据
export interface DashboardStats {
  total_assets: number
  total_users: number
  total_tasks: number
  total_projects: number
  assets_by_status: Record<string, number>
  task_execution_stats: TaskExecutionStats
  assets_by_environment: EnvironmentAssetCount[]
  assets_by_project: ProjectAssetCount[]
  recent_activities: RecentActivity[]
}

export const dashboardApi = {
  getStats: () => apiClient.get<{ data: DashboardStats }>('/dashboard/stats'),
}
