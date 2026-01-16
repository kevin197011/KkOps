// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface Tag {
  id: number
  name: string
  color: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateTagRequest {
  name: string
  color?: string
  description?: string
}

export interface UpdateTagRequest {
  name?: string
  color?: string
  description?: string
}

export const tagApi = {
  list: () => apiClient.get<Tag[]>('/tags'),
  get: (id: number) => apiClient.get<{ data: Tag }>(`/tags/${id}`),
  create: (data: CreateTagRequest) => apiClient.post<{ data: Tag }>('/tags', data),
  update: (id: number, data: UpdateTagRequest) => apiClient.put<{ data: Tag }>(`/tags/${id}`, data),
  delete: (id: number) => apiClient.delete(`/tags/${id}`),
}
