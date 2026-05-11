// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

/** sso_link = 同 SSO 应用（与 KkOps 共用同一 IdP，打开即已登录）；jwt_token = Token 跳转（附带身份与权限） */
export type LaunchType = 'sso_link' | 'jwt_token'

export interface ExternalSystem {
  id: number
  name: string
  description?: string
  launch_type: LaunchType
  base_url: string
  login_path?: string
  role_mapping?: string
  icon?: string
  order_index: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface CreateExternalSystemRequest {
  name: string
  description?: string
  launch_type?: LaunchType
  base_url: string
  login_path?: string
  secret?: string
  role_mapping?: string
  icon?: string
  order_index?: number
  enabled?: boolean
}

export interface UpdateExternalSystemRequest {
  name?: string
  description?: string
  launch_type?: LaunchType
  base_url?: string
  login_path?: string
  secret?: string
  role_mapping?: string
  icon?: string
  order_index?: number
  enabled?: boolean
}

export interface LaunchResponse {
  redirect_url: string
}

export const externalSystemApi = {
  list: (enabledOnly = true) =>
    apiClient.get<{ data: ExternalSystem[] }>(
      `/external-systems${enabledOnly ? '' : '?enabled=false'}`
    ),
  get: (id: number) =>
    apiClient.get<{ data: ExternalSystem }>(`/external-systems/${id}`),
  create: (data: CreateExternalSystemRequest) =>
    apiClient.post<{ data: ExternalSystem }>('/external-systems', data),
  update: (id: number, data: UpdateExternalSystemRequest) =>
    apiClient.put<{ data: ExternalSystem }>(`/external-systems/${id}`, data),
  delete: (id: number) =>
    apiClient.delete<unknown>(`/external-systems/${id}`),
  launch: (id: number) =>
    apiClient.post<LaunchResponse>(`/external-systems/${id}/launch`),
}
