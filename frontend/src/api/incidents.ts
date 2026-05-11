// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface IncidentRow {
  id: number
  title: string
  severity: string
  status: string
  linked_alert_ids?: number[]
  assignee_user_id?: number
  created_by: number
  created_at: string
  updated_at: string
}

export interface IncidentCreate {
  title: string
  severity: string
  linked_alert_ids?: number[]
  assignee_user_id?: number
}

export interface IncidentPatch {
  title?: string
  severity?: string
  status?: string
  linked_alert_ids?: number[]
  assignee_user_id?: number
}

export const incidentsApi = {
  list: (params?: { status?: string; limit?: number; offset?: number }) =>
    apiClient.get<{ data: IncidentRow[]; total: number }>('/incidents', { params }),
  create: (body: IncidentCreate) => apiClient.post<{ data: IncidentRow }>('/incidents', body),
  get: (id: number) => apiClient.get<{ data: IncidentRow }>(`/incidents/${id}`),
  patch: (id: number, body: IncidentPatch) =>
    apiClient.patch<{ data: IncidentRow }>(`/incidents/${id}`, body),
}
