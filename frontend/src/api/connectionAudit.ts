// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

// WebSSH 连线记录（列表项不含 transcript）
export interface SSHConnectionRecord {
  id: number
  user_id: number
  login_username: string // 操作用户：KkOps 登录用户
  username: string // 连线用户：SSH 登录名（如 root）
  asset_id: number
  asset_hostname: string
  started_at: string
  ended_at: string
  duration_seconds: number
  transcript?: string
  transcript_truncated?: boolean
  created_at: string
}

// 连线记录列表响应
export interface ConnectionAuditListResponse {
  total: number
  data: SSHConnectionRecord[]
}

// 查询参数
export interface ConnectionAuditQueryParams {
  page?: number
  page_size?: number
  user_id?: number
  asset_id?: number
  start_time?: string
  end_time?: string
}

export const connectionAuditApi = {
  list: (params: ConnectionAuditQueryParams = {}) =>
    apiClient.get<ConnectionAuditListResponse>('/connection-audit', { params }),

  get: (id: number) =>
    apiClient.get<SSHConnectionRecord>(`/connection-audit/${id}`),
}
