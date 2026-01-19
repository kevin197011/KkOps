// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface ExecutionTemplate {
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

export interface Execution {
  id: number
  template_id?: number
  name: string
  description: string
  content: string
  type: string
  status: string
  asset_ids: number[]
  timeout?: number
  created_by: number
  started_at?: string
  finished_at?: string
  created_at: string
  updated_at: string
}

export interface ExecutionRecordAsset {
  id: number
  hostname: string
  ip: string
}

export interface ExecutionRecord {
  id: number
  task_id: number
  asset_id: number
  asset?: ExecutionRecordAsset
  status: string
  output?: string
  error?: string
  exit_code?: number
  started_at?: string
  finished_at?: string
  created_at: string
  updated_at: string
}

export interface CreateExecutionRequest {
  template_id?: number
  name: string
  description?: string
  content: string
  type?: string
  asset_ids: number[]
}

export interface UpdateExecutionRequest {
  name?: string
  description?: string
  content?: string
  type?: string
  status?: string
}

export interface ListExecutionsResponse {
  data: Execution[]
  total: number
  page: number
  size: number
}

export interface ListTemplatesResponse {
  data: ExecutionTemplate[]
  total: number
  page: number
  size: number
}

// 导入导出相关接口
export interface ExportTemplateConfig {
  name: string
  description: string
  content: string
  type: string
}

export interface ExportTemplatesConfig {
  version: string
  export_at: string
  templates: ExportTemplateConfig[]
}

export interface ImportTemplateConfig {
  name: string
  description?: string
  content: string
  type?: string
}

export interface ImportTemplatesConfig {
  version: string
  templates: ImportTemplateConfig[]
}

export interface ImportResult {
  total: number
  success: number
  failed: number
  errors?: string[]
  skipped?: string[]
}

// 执行模板 API (原 task-templates)
export const executionTemplateApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<ExecutionTemplate[]>(`/templates${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<ExecutionTemplate>(`/templates/${id}`),
  create: (data: CreateTemplateRequest) => apiClient.post<ExecutionTemplate>('/templates', data),
  update: (id: number, data: UpdateTemplateRequest) => apiClient.put<ExecutionTemplate>(`/templates/${id}`, data),
  delete: (id: number) => apiClient.delete(`/templates/${id}`),
  // 导出配置
  exportConfig: () => apiClient.get<ExportTemplatesConfig>('/templates/export'),
  // 导入配置
  importConfig: (data: ImportTemplatesConfig) =>
    apiClient.post<ImportResult>('/templates/import', data),
}

// 运维执行 API (原 tasks)
export const executionApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<ListExecutionsResponse>(`/executions${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<Execution>(`/executions/${id}`),
  create: (data: CreateExecutionRequest) => apiClient.post<Execution>('/executions', data),
  update: (id: number, data: UpdateExecutionRequest) => apiClient.put<Execution>(`/executions/${id}`, data),
  delete: (id: number) => apiClient.delete(`/executions/${id}`),
  execute: (id: number, executionType: 'sync' | 'async' = 'sync') =>
    apiClient.post(`/executions/${id}/execute`, { execution_type: executionType }),
  cancel: (id: number) => apiClient.post(`/executions/${id}/cancel`),
  getHistory: (id: number) => apiClient.get<{ data: ExecutionRecord[] }>(`/executions/${id}/history`),
}

// 执行记录 API (原 task-executions)
export const executionRecordApi = {
  get: (id: number) => apiClient.get<{ data: ExecutionRecord }>(`/execution-records/${id}`),
  cancel: (id: number) => apiClient.post(`/execution-records/${id}/cancel`),
  getLogs: (id: number) => apiClient.get<{ data: { logs: string[] } }>(`/execution-records/${id}/logs`),
}

// 兼容旧接口名称（可在迁移完成后移除）
export const taskTemplateApi = executionTemplateApi
export const taskApi = executionApi
export const taskExecutionApi = executionRecordApi
export const templateApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<{ data: ExecutionTemplate[] }>(`/templates${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<ExecutionTemplate>(`/templates/${id}`),
  create: (data: CreateTemplateRequest) => apiClient.post<ExecutionTemplate>('/templates', data),
  update: (id: number, data: UpdateTemplateRequest) => apiClient.put<ExecutionTemplate>(`/templates/${id}`, data),
  delete: (id: number) => apiClient.delete(`/templates/${id}`),
}

// 类型兼容（用于迁移期）
export type TaskTemplate = ExecutionTemplate
export type Task = Execution
export type TaskExecution = ExecutionRecord
export type TaskExecutionAsset = ExecutionRecordAsset
export type CreateTaskRequest = CreateExecutionRequest
export type UpdateTaskRequest = UpdateExecutionRequest
