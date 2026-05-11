// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface OAuth2Client {
  id: number
  client_id: string
  client_secret?: string
  name: string
  protocol: string
  redirect_uris: string[]
  scopes: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface CreateOAuth2ClientRequest {
  name: string
  protocol?: string
  redirect_uris: string[]
  scopes?: string
}

export interface UpdateOAuth2ClientRequest {
  name?: string
  protocol?: string
  redirect_uris?: string[]
  scopes?: string
  enabled?: boolean
}

export const oauth2ClientApi = {
  list: (params?: { protocol?: string }) =>
    apiClient.get<{ data: OAuth2Client[] }>('/oauth2-clients', { params }),
  get: (id: number) =>
    apiClient.get<{ data: OAuth2Client }>(`/oauth2-clients/${id}`),
  create: (data: CreateOAuth2ClientRequest) =>
    apiClient.post<{ data: OAuth2Client }>('/oauth2-clients', data),
  update: (id: number, data: UpdateOAuth2ClientRequest) =>
    apiClient.put<{ data: OAuth2Client }>(`/oauth2-clients/${id}`, data),
  delete: (id: number) => apiClient.delete<unknown>(`/oauth2-clients/${id}`),
  regenerateSecret: (id: number) =>
    apiClient.post<{ data: { client_secret: string } }>(
      `/oauth2-clients/${id}/regenerate-secret`
    ),
}
