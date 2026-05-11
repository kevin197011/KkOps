// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface ApplicationSummary {
  name: string
  health: string
  sync_status: string
}

export interface PipelineEvent {
  ts: string
  kind: string
  source: string
  ref: string
  status: string
  link: string
}

export const gitopsApi = {
  listApplications: (integrationId: number) =>
    apiClient.get<{ data: ApplicationSummary[] }>('/gitops/applications', {
      params: { integration_id: integrationId },
    }),
  sync: (name: string, integrationId: number) =>
    apiClient.post(`/gitops/applications/${encodeURIComponent(name)}/sync`, {
      integration_id: integrationId,
    }),
  pipelineView: (params: { app: string; argocd_integration_id?: number }) =>
    apiClient.get<{ data: PipelineEvent[] }>('/gitops/pipeline-view', { params }),
}
