// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface Role {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  description?: string
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
}

export const roleApi = {
  list: () => apiClient.get<Role[]>('/roles'),
  get: (id: number) => apiClient.get<Role>(`/roles/${id}`),
  create: (data: CreateRoleRequest) => apiClient.post<Role>('/roles', data),
  update: (id: number, data: UpdateRoleRequest) => apiClient.put<Role>(`/roles/${id}`, data),
  delete: (id: number) => apiClient.delete(`/roles/${id}`),
}
