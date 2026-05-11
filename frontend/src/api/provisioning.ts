// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface IntegrationSummary {
  id: number
  name: string
  kind: string
  enabled: boolean
  description?: string
}

export interface ProvisioningTarget {
  id: number
  integration_id: number
  provider_kind: string
  status: string
  last_error?: string
  last_sync_at?: string
  enabled: boolean
  created_at: string
  updated_at: string
  integration?: IntegrationSummary
}

export interface ProvisioningRun {
  id: number
  target_id: number
  user_id?: number
  action: string
  status: string
  message?: string
  started_at: string
  ended_at: string
}

export interface CreateProvisioningTargetRequest {
  integration_id: number
  provider_kind: string
  enabled?: boolean
}

export const provisioningApi = {
  listTargets: () =>
    apiClient.get<{ data: ProvisioningTarget[] }>('/provisioning/targets'),
  createTarget: (data: CreateProvisioningTargetRequest) =>
    apiClient.post<{ data: ProvisioningTarget }>('/provisioning/targets', data),
  syncTarget: (id: number) =>
    apiClient.post<{ status: string }>(`/provisioning/targets/${id}/sync`),
  listRuns: (id: number, limit = 50) =>
    apiClient.get<{ data: ProvisioningRun[] }>(
      `/provisioning/targets/${id}/runs`,
      { params: { limit } }
    ),
}

export { integrationsApi } from './integration'
