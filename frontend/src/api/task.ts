// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface TaskTemplate {
  id: number
  name: string
  description: string
  content: string
  type: string
  created_by: number
  created_at: string
  updated_at: string
}

export interface CreateTemplateRequest {
  name: string
  description?: string
  content: string
  type?: string
}

export interface UpdateTemplateRequest {
  name?: string
  description?: string
  content?: string
  type?: string
}

export interface Task {
  id: number
  template_id?: number
  name: string
  description: string
  content: string
  type: string
  status: string
  asset_ids: number[]
  created_by: number
  started_at?: string
  finished_at?: string
  created_at: string
  updated_at: string
}

export interface TaskExecutionAsset {
  id: number
  hostname: string
  ip: string
}

export interface TaskExecution {
  id: number
  task_id: number
  asset_id: number
  asset?: TaskExecutionAsset
  status: string
  output?: string
  error?: string
  exit_code?: number
  started_at?: string
  finished_at?: string
  created_at: string
  updated_at: string
}

export interface CreateTaskRequest {
  template_id?: number
  name: string
  description?: string
  content: string
  type?: string
  asset_ids: number[]
}

export interface UpdateTaskRequest {
  name?: string
  description?: string
  content?: string
  type?: string
  status?: string
}

export interface ListTasksResponse {
  data: Task[]
  total: number
  page: number
  size: number
}

export interface ListTemplatesResponse {
  data: TaskTemplate[]
  total: number
  page: number
  size: number
}

export const taskTemplateApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<TaskTemplate[]>(`/task-templates${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<TaskTemplate>(`/task-templates/${id}`),
  create: (data: CreateTemplateRequest) => apiClient.post<TaskTemplate>('/task-templates', data),
  update: (id: number, data: UpdateTemplateRequest) => apiClient.put<TaskTemplate>(`/task-templates/${id}`, data),
  delete: (id: number) => apiClient.delete(`/task-templates/${id}`),
}

export const taskApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<ListTasksResponse>(`/tasks${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<Task>(`/tasks/${id}`),
  create: (data: CreateTaskRequest) => apiClient.post<Task>('/tasks', data),
  update: (id: number, data: UpdateTaskRequest) => apiClient.put<Task>(`/tasks/${id}`, data),
  delete: (id: number) => apiClient.delete(`/tasks/${id}`),
  execute: (id: number, executionType: 'sync' | 'async' = 'sync') =>
    apiClient.post(`/tasks/${id}/execute`, { execution_type: executionType }),
  cancel: (id: number) => apiClient.post(`/tasks/${id}/cancel`),
  getExecutions: (id: number) => apiClient.get<{ data: TaskExecution[] }>(`/tasks/${id}/executions`),
}

export const taskExecutionApi = {
  get: (id: number) => apiClient.get<{ data: TaskExecution }>(`/task-executions/${id}`),
  cancel: (id: number) => apiClient.post(`/task-executions/${id}/cancel`),
  getLogs: (id: number) => apiClient.get<{ data: { logs: string[] } }>(`/task-executions/${id}/logs`),
}

export const templateApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<{ data: TaskTemplate[] }>(`/task-templates${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<TaskTemplate>(`/task-templates/${id}`),
  create: (data: CreateTemplateRequest) => apiClient.post<TaskTemplate>('/task-templates', data),
  update: (id: number, data: UpdateTemplateRequest) => apiClient.put<TaskTemplate>(`/task-templates/${id}`, data),
  delete: (id: number) => apiClient.delete(`/task-templates/${id}`),
}
