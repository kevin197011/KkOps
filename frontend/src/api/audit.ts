// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

// 审计日志
export interface AuditLog {
  id: number
  user_id: number
  username: string
  action: string
  module: string
  resource_id?: number
  resource_name?: string
  detail?: string
  ip_address: string
  user_agent?: string
  status: string
  error_msg?: string
  created_at: string
}

// 审计日志列表响应
export interface AuditLogListResponse {
  total: number
  items: AuditLog[]
  page: number
  page_size: number
}

// 审计日志查询参数
export interface AuditLogQueryParams {
  page?: number
  page_size?: number
  user_id?: number
  username?: string
  module?: string
  action?: string
  status?: string
  start_time?: string
  end_time?: string
  keyword?: string
}

// 审计日志 API
export const auditApi = {
  // 获取审计日志列表
  list: (params: AuditLogQueryParams = {}) =>
    apiClient.get<{ data: AuditLogListResponse }>('/audit-logs', { params }),

  // 获取单条审计日志
  get: (id: number) => apiClient.get<{ data: AuditLog }>(`/audit-logs/${id}`),

  // 导出审计日志
  export: async (params: AuditLogQueryParams & { format: 'csv' | 'json' }) => {
    const queryParams: Record<string, string> = {}
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        queryParams[key] = String(value)
      }
    })
    
    const response = await apiClient.get('/audit-logs/export', {
      params: queryParams,
      responseType: 'blob',
    })
    
    // 从 Content-Disposition header 获取文件名，或生成默认文件名
    const contentDisposition = response.headers['content-disposition'] || response.headers['Content-Disposition']
    let filename = `audit_logs_${new Date().toISOString().slice(0, 10)}.${params.format}`
    if (contentDisposition) {
      // 尝试匹配 filename*=UTF-8''格式
      const utf8Match = contentDisposition.match(/filename\*=UTF-8''(.+)/)
      if (utf8Match && utf8Match[1]) {
        filename = decodeURIComponent(utf8Match[1])
      } else {
        // 尝试匹配标准 filename= 格式
        const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
        if (filenameMatch && filenameMatch[1]) {
          filename = filenameMatch[1].replace(/['"]/g, '')
        }
      }
    }
    
    // 创建 blob URL 并触发下载
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  },

  // 获取模块列表
  getModules: () => apiClient.get<{ data: string[] }>('/audit-logs/modules'),

  // 获取操作类型列表
  getActions: () => apiClient.get<{ data: string[] }>('/audit-logs/actions'),
}

// 模块名称映射
export const moduleLabels: Record<string, string> = {
  auth: '认证',
  user: '用户管理',
  role: '角色管理',
  asset: '资产管理',
  task: '任务执行',
  template: '任务模板',
  scheduled_task: '定时任务',
  deployment: '部署管理',
  ssh: 'SSH 连接',
  ssh_key: 'SSH 密钥',
  project: '项目管理',
  environment: '环境管理',
  cloud_platform: '云平台',
  tag: '标签管理',
}

// 操作类型映射
export const actionLabels: Record<string, string> = {
  create: '创建',
  update: '更新',
  delete: '删除',
  execute: '执行',
  login: '登录',
  logout: '登出',
  enable: '启用',
  disable: '禁用',
  export: '导出',
  connect: '连接',
}

// 状态映射
export const statusLabels: Record<string, string> = {
  success: '成功',
  failed: '失败',
}
