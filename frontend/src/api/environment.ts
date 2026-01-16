// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface Environment {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateEnvironmentRequest {
  name: string
  description?: string
}

export interface UpdateEnvironmentRequest {
  name?: string
  description?: string
}

export const environmentApi = {
  list: () => apiClient.get<Environment[]>('/environments'),
  get: (id: number) => apiClient.get<Environment>(`/environments/${id}`),
  create: (data: CreateEnvironmentRequest) => apiClient.post<Environment>('/environments', data),
  update: (id: number, data: UpdateEnvironmentRequest) => apiClient.put<Environment>(`/environments/${id}`, data),
  delete: (id: number) => apiClient.delete(`/environments/${id}`),
}
