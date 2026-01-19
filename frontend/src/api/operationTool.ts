// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

// 运维工具接口
export interface OperationTool {
  id: number
  name: string
  description?: string
  category?: string
  icon?: string
  url: string
  order: number
  enabled: boolean
  created_at: string
  updated_at: string
}

// 创建运维工具请求
export interface CreateOperationToolRequest {
  name: string
  description?: string
  category?: string
  icon?: string
  url: string
  order?: number
  enabled?: boolean
}

// 更新运维工具请求
export interface UpdateOperationToolRequest {
  name?: string
  description?: string
  category?: string
  icon?: string
  url?: string
  order?: number
  enabled?: boolean
}

// 列表查询参数
export interface ListOperationToolsParams {
  category?: string
  enabled?: boolean
}

export const operationToolApi = {
  // 获取工具列表
  list: (params?: ListOperationToolsParams) => {
    const queryParams = new URLSearchParams()
    if (params?.category) {
      queryParams.append('category', params.category)
    }
    if (params?.enabled !== undefined) {
      queryParams.append('enabled', params.enabled.toString())
    }
    const query = queryParams.toString()
    return apiClient.get<{ data: OperationTool[] }>(
      `/operation-tools${query ? `?${query}` : ''}`
    )
  },

  // 获取单个工具
  get: (id: number) =>
    apiClient.get<{ data: OperationTool }>(`/operation-tools/${id}`),

  // 创建工具（管理员）
  create: (data: CreateOperationToolRequest) =>
    apiClient.post<{ data: OperationTool }>('/operation-tools', data),

  // 更新工具（管理员）
  update: (id: number, data: UpdateOperationToolRequest) =>
    apiClient.put<{ data: OperationTool }>(`/operation-tools/${id}`, data),

  // 删除工具（管理员）
  delete: (id: number) =>
    apiClient.delete<{ message: string }>(`/operation-tools/${id}`),
}
