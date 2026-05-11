// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface RepositorySummary {
  name: string
  project_id: number
  pull_count: number
  update_time: string
}

export interface TagSummary {
  name: string
  digest?: string
  size?: number
}

export const registryApi = {
  listRepositories: (integrationId: number, params?: { page?: number; page_size?: number; project_name?: string }) =>
    apiClient.get<{ data: RepositorySummary[] }>('/registry/repositories', {
      params: { integration_id: integrationId, ...params },
    }),
  listTags: (integrationId: number, repository: string, params?: { page?: number; page_size?: number }) =>
    apiClient.get<{ data: TagSummary[] }>('/registry/tags', {
      params: { integration_id: integrationId, repository, ...params },
    }),
  vulnerabilities: (integrationId: number, repository: string, reference: string) =>
    apiClient.get<{ data: { raw_json: string } }>('/registry/vulnerabilities', {
      params: { integration_id: integrationId, repository, reference },
    }),
}
