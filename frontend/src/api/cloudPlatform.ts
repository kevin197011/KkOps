// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface CloudPlatform {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateCloudPlatformRequest {
  name: string
  description?: string
}

export interface UpdateCloudPlatformRequest {
  name?: string
  description?: string
}

export const cloudPlatformApi = {
  list: () => apiClient.get<CloudPlatform[]>('/cloud-platforms'),
  get: (id: number) => apiClient.get<CloudPlatform>(`/cloud-platforms/${id}`),
  create: (data: CreateCloudPlatformRequest) => apiClient.post<CloudPlatform>('/cloud-platforms', data),
  update: (id: number, data: UpdateCloudPlatformRequest) => apiClient.put<CloudPlatform>(`/cloud-platforms/${id}`, data),
  delete: (id: number) => apiClient.delete(`/cloud-platforms/${id}`),
}
