// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface PipelineRow {
  id: string
  name: string
  status: string
  integration_id: number
  web_url?: string
  ref?: string
}

export const cicdApi = {
  listPipelines: (integrationId: number, project?: string) =>
    apiClient.get<{ data: PipelineRow[] }>('/cicd/pipelines', {
      params: { integration_id: integrationId, project },
    }),
  runPipeline: (
    pipelineId: string,
    body: { integration_id: number; ref?: string; variables?: Record<string, string> }
  ) => apiClient.post(`/cicd/pipelines/${encodeURIComponent(pipelineId)}/run`, body),
  logs: (pipelineId: string, integrationId: number) =>
    apiClient.get<{ data: string }>(`/cicd/pipelines/${encodeURIComponent(pipelineId)}/logs`, {
      params: { integration_id: integrationId },
    }),
}
