// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface AlertRow {
  id: number
  integration_id: number
  source: string
  fingerprint: string
  severity: string
  title: string
  status: string
  labels_json?: string
  starts_at?: string
  ends_at?: string
  created_at: string
  updated_at: string
  integration_name?: string
}

export const alertsApi = {
  list: (params?: { status?: string; integration_id?: number; limit?: number; offset?: number }) =>
    apiClient.get<{ data: AlertRow[]; total: number }>('/alerts', { params }),
  sync: () => apiClient.post<{ synced: number; warnings: string[] }>('/alerts/sync'),
  acknowledge: (id: number) => apiClient.post(`/alerts/${id}/acknowledge`),
  dismiss: (id: number) => apiClient.post(`/alerts/${id}/dismiss`),
}
