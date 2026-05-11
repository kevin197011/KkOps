// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface CMDBAsset {
  id: number
  name: string
  kind: string
  env: string
  owner_user_id?: number
  labels?: Record<string, unknown>
  integration_id?: number
  external_ref?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface AssetRelation {
  id: number
  from_asset_id: number
  to_asset_id: number
  relation_type: string
  meta?: Record<string, unknown>
  created_at: string
  from_asset?: CMDBAsset
  to_asset?: CMDBAsset
}

export interface TopologyGraph {
  nodes: Array<{
    id: number
    kind: string
    name: string
    subtype?: string
    env?: string
    extra?: Record<string, unknown>
  }>
  edges: Array<{
    id?: number
    from: number
    to: number
    relation_type: string
  }>
  derived_counts: Record<string, number>
  placeholder_hints?: Record<string, unknown>
}

export const cmdbApi = {
  listAssets: (params?: { kind?: string; env?: string; q?: string; limit?: number; offset?: number }) =>
    apiClient.get<{ data: CMDBAsset[]; total: number }>('/cmdb/assets', { params }),

  createAsset: (body: Partial<CMDBAsset> & { name: string; kind: string }) =>
    apiClient.post<{ data: CMDBAsset }>('/cmdb/assets', body),

  updateAsset: (id: number, body: Partial<CMDBAsset>) =>
    apiClient.put<{ data: CMDBAsset }>(`/cmdb/assets/${id}`, body),

  deleteAsset: (id: number) => apiClient.delete(`/cmdb/assets/${id}`),

  listRelations: (params?: { from_asset_id?: number; to_asset_id?: number }) =>
    apiClient.get<{ data: AssetRelation[]; total: number }>('/cmdb/asset-relations', { params }),

  createRelation: (body: {
    from_asset_id: number
    to_asset_id: number
    relation_type: string
    meta?: Record<string, unknown>
  }) => apiClient.post<{ data: AssetRelation }>('/cmdb/asset-relations', body),

  deleteRelation: (id: number) => apiClient.delete(`/cmdb/asset-relations/${id}`),

  topologyGraph: () => apiClient.get<{ data: TopologyGraph }>('/topology/graph'),
}
