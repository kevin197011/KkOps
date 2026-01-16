// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface Category {
  id: number
  name: string
  parent_id?: number
  description: string
  created_at: string
  updated_at: string
}

export interface CreateCategoryRequest {
  name: string
  parent_id?: number
  description?: string
}

export interface UpdateCategoryRequest {
  name?: string
  description?: string
}

export const categoryApi = {
  list: () => apiClient.get<Category[]>('/asset-categories'),
  get: (id: number) => apiClient.get<{ data: Category }>(`/asset-categories/${id}`),
  create: (data: CreateCategoryRequest) => apiClient.post<{ data: Category }>('/asset-categories', data),
  update: (id: number, data: UpdateCategoryRequest) => apiClient.put<{ data: Category }>(`/asset-categories/${id}`, data),
  delete: (id: number) => apiClient.delete(`/asset-categories/${id}`),
}
