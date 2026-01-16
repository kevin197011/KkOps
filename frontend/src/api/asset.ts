// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface TagInfo {
  id: number
  name: string
  color: string
}

export interface CloudPlatformInfo {
  id: number
  name: string
  description: string
}

export interface Asset {
  id: number
  hostName: string
  project_id?: number
  cloud_platform_id?: number
  cloud_platform?: CloudPlatformInfo
  environment_id?: number
  ip: string
  ssh_port: number
  ssh_key_id?: number
  ssh_user: string
  cpu: string
  memory: string
  disk: string
  status: string
  description: string
  tags: TagInfo[]
  created_at: string
  updated_at: string
}

export interface CreateAssetRequest {
  hostName: string
  project_id?: number
  cloud_platform_id?: number
  environment_id?: number
  ip?: string
  ssh_port?: number
  ssh_key_id?: number
  ssh_user?: string
  cpu?: string
  memory?: string
  disk?: string
  status?: string
  description?: string
  tag_ids?: number[]
}

export interface UpdateAssetRequest {
  hostName?: string
  project_id?: number
  cloud_platform_id?: number
  environment_id?: number
  ip?: string
  ssh_port?: number
  ssh_key_id?: number
  ssh_user?: string
  cpu?: string
  memory?: string
  disk?: string
  status?: string
  description?: string
  tag_ids?: number[]
}

export interface ListAssetsResponse {
  data: Asset[]
  total: number
  page: number
  size: number
}

export interface ListAssetsParams {
  page?: number
  page_size?: number
  project_id?: number
  cloud_platform_id?: number
  environment_id?: number
  status?: string
  ip?: string
  ssh_key_id?: number
  tag_ids?: number[]
  search?: string
}

export const assetApi = {
  list: (params?: ListAssetsParams) => {
    const queryParams = new URLSearchParams()
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            queryParams.append(key, value.join(','))
          } else {
            queryParams.append(key, String(value))
          }
        }
      })
    }
    const query = queryParams.toString()
    return apiClient.get<ListAssetsResponse>(`/assets${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<Asset>(`/assets/${id}`),
  create: (data: CreateAssetRequest) => apiClient.post<Asset>('/assets', data),
  update: (id: number, data: UpdateAssetRequest) => apiClient.put<Asset>(`/assets/${id}`, data),
  delete: (id: number) => apiClient.delete(`/assets/${id}`),
  import: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return apiClient.post('/assets/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  export: () => apiClient.get<Blob>('/assets/export', { responseType: 'blob' }),
}
