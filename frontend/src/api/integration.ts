// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface IntegrationPublic {
  id: number
  name: string
  kind: string
  enabled: boolean
  description?: string
  has_config: boolean
  created_at: string
  updated_at: string
}

export interface CreateIntegrationRequest {
  name: string
  kind: string
  enabled?: boolean
  description?: string
  config: Record<string, unknown>
}

export const integrationsApi = {
  list: () => apiClient.get<{ data: IntegrationPublic[] }>('/integrations'),
  get: (id: number) => apiClient.get<{ data: IntegrationPublic }>(`/integrations/${id}`),
  create: (data: CreateIntegrationRequest) =>
    apiClient.post<{ data: IntegrationPublic }>('/integrations', data),
  update: (id: number, data: CreateIntegrationRequest) =>
    apiClient.put<{ data: IntegrationPublic }>(`/integrations/${id}`, data),
  delete: (id: number) => apiClient.delete(`/integrations/${id}`),
  test: (id: number) =>
    apiClient.post<{ ok: boolean; error?: string; metadata?: Record<string, string> }>(
      `/integrations/${id}/test`
    ),
}
