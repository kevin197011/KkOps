// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

// 定时任务相关接口
export interface ScheduledTask {
  id: number
  name: string
  description: string
  cron_expression: string
  template_id?: number
  template_name?: string
  content: string
  type: string
  asset_ids: number[]
  timeout: number
  enabled: boolean
  update_assets: boolean // 是否更新资产信息
  last_run_at?: string
  next_run_at?: string
  last_status?: string
  created_by: number
  creator_name?: string
  created_at: string
  updated_at: string
}

export interface CreateScheduledTaskRequest {
  name: string
  description?: string
  cron_expression: string
  template_id?: number
  content: string
  type?: string
  asset_ids: number[]
  timeout?: number
  enabled?: boolean
  update_assets?: boolean // 是否更新资产信息
}

export interface UpdateScheduledTaskRequest {
  name?: string
  description?: string
  cron_expression?: string
  template_id?: number
  content?: string
  type?: string
  asset_ids?: number[]
  timeout?: number
  enabled?: boolean
  update_assets?: boolean // 是否更新资产信息
}

export interface ListScheduledTasksResponse {
  data: ScheduledTask[]
  total: number
  page: number
  size: number
}

export interface ScheduledTaskExecution {
  id: number
  scheduled_task_id: number
  asset_id: number
  asset?: {
    id: number
    hostname: string
    ip: string
  }
  trigger_type: string
  status: string
  output?: string
  error?: string
  exit_code?: number
  started_at?: string
  finished_at?: string
  created_at: string
}

export interface ValidateCronResponse {
  valid: boolean
  next_run_at?: string
  error?: string
}

// 定时任务 API
export const scheduledTaskApi = {
  // 获取列表
  list: (page?: number, pageSize?: number, enabled?: boolean) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    if (enabled !== undefined) params.append('enabled', String(enabled))
    const query = params.toString()
    return apiClient.get<ListScheduledTasksResponse>(`/tasks${query ? `?${query}` : ''}`)
  },

  // 获取详情
  get: (id: number) => apiClient.get<ScheduledTask>(`/tasks/${id}`),

  // 创建
  create: (data: CreateScheduledTaskRequest) =>
    apiClient.post<ScheduledTask>('/tasks', data),

  // 更新
  update: (id: number, data: UpdateScheduledTaskRequest) =>
    apiClient.put<ScheduledTask>(`/tasks/${id}`, data),

  // 删除
  delete: (id: number) => apiClient.delete(`/tasks/${id}`),

  // 启用
  enable: (id: number) => apiClient.post(`/tasks/${id}/enable`),

  // 禁用
  disable: (id: number) => apiClient.post(`/tasks/${id}/disable`),

  // 获取执行历史
  getExecutions: (id: number, page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<{ data: ScheduledTaskExecution[]; total: number }>(
      `/tasks/${id}/executions${query ? `?${query}` : ''}`
    )
  },

  // 验证 Cron 表达式
  validateCron: (cronExpression: string) =>
    apiClient.get<ValidateCronResponse>(
      `/tasks/validate-cron?cron_expression=${encodeURIComponent(cronExpression)}`
    ),

  // 导出配置
  exportConfig: () => apiClient.get<ExportScheduledTasksConfig>('/tasks/export'),

  // 导入配置
  importConfig: (data: ImportScheduledTasksConfig) =>
    apiClient.post<ImportScheduledTasksResult>('/tasks/import', data),
}

// 导入导出相关接口
export interface ExportScheduledTaskConfig {
  name: string
  description: string
  cron_expression: string
  template_name?: string
  content: string
  type: string
  timeout: number
  enabled: boolean
  update_assets: boolean
  target_hosts: string[]
}

export interface ExportScheduledTasksConfig {
  version: string
  export_at: string
  tasks: ExportScheduledTaskConfig[]
}

export interface ImportScheduledTaskConfig {
  name: string
  description?: string
  cron_expression: string
  template_name?: string
  content?: string
  type?: string
  timeout?: number
  enabled?: boolean
  update_assets?: boolean
  target_hosts?: string[]
}

export interface ImportScheduledTasksConfig {
  version: string
  tasks: ImportScheduledTaskConfig[]
}

export interface ImportScheduledTasksResult {
  total: number
  success: number
  failed: number
  errors?: string[]
  skipped?: string[]
}

// 导出默认 API
export default scheduledTaskApi
