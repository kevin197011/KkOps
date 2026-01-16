// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface Project {
  id: number
  name: string
  description: string
  created_by: number
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  description?: string
}

export interface UpdateProjectRequest {
  name?: string
  description?: string
}

export const projectApi = {
  list: () => apiClient.get<Project[]>('/projects'),
  get: (id: number) => apiClient.get<{ data: Project }>(`/projects/${id}`),
  create: (data: CreateProjectRequest) => apiClient.post<{ data: Project }>('/projects', data),
  update: (id: number, data: UpdateProjectRequest) => apiClient.put<{ data: Project }>(`/projects/${id}`, data),
  delete: (id: number) => apiClient.delete(`/projects/${id}`),
}
