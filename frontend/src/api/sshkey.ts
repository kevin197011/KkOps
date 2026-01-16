// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface SSHKey {
  id: number
  user_id: number
  name: string
  type: string
  public_key: string
  fingerprint: string
  ssh_user: string
  description: string
  last_used_at?: string
  created_at: string
  updated_at: string
}

export interface CreateSSHKeyRequest {
  name: string
  type?: string
  public_key?: string // Optional - auto-extracted from private key if not provided
  private_key: string
  ssh_user?: string
  passphrase?: string
  description?: string
}

export interface UpdateSSHKeyRequest {
  name?: string
  ssh_user?: string
  description?: string
}

export interface TestSSHKeyRequest {
  host: string
  port?: number
  username?: string
  passphrase?: string
}

export const sshkeyApi = {
  list: () => apiClient.get<{ data: SSHKey[] }>('/ssh/keys'),
  get: (id: number) => apiClient.get<{ data: SSHKey }>(`/ssh/keys/${id}`),
  create: (data: CreateSSHKeyRequest) => apiClient.post<{ data: SSHKey }>('/ssh/keys', data),
  update: (id: number, data: UpdateSSHKeyRequest) => apiClient.put<{ data: SSHKey }>(`/ssh/keys/${id}`, data),
  delete: (id: number) => apiClient.delete(`/ssh/keys/${id}`),
  test: (id: number, data: TestSSHKeyRequest) => apiClient.post<{ message: string }>(`/ssh/keys/${id}/test`, data),
}
