// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface Role {
  id: number
  name: string
  description: string
  is_admin: boolean // 是否为管理员角色
  asset_count: number // 授权资产数量（-1 表示全部）
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  description?: string
  is_admin?: boolean // 是否为管理员角色
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  is_admin?: boolean // 是否为管理员角色
}

// 角色授权的资产信息
export interface RoleAssetInfo {
  id: number
  hostname: string
  ip: string
}

// 权限信息
export interface Permission {
  id: number
  name: string
  resource: string
  action: string
  description: string
}

export const roleApi = {
  list: () => apiClient.get<Role[]>('/roles'),
  get: (id: number) => apiClient.get<Role>(`/roles/${id}`),
  create: (data: CreateRoleRequest) => apiClient.post<Role>('/roles', data),
  update: (id: number, data: UpdateRoleRequest) => apiClient.put<Role>(`/roles/${id}`, data),
  delete: (id: number) => apiClient.delete(`/roles/${id}`),
  
  // 权限管理
  listPermissions: () => apiClient.get<Permission[]>('/permissions'),
  getRolePermissions: (roleId: number) =>
    apiClient.get<Permission[]>(`/roles/${roleId}/permissions`),
  assignPermissionToRole: (roleId: number, permissionId: number) =>
    apiClient.post(`/roles/${roleId}/permissions`, { permission_id: permissionId }),
  removePermissionFromRole: (roleId: number, permissionId: number) =>
    apiClient.delete(`/roles/${roleId}/permissions/${permissionId}`),
  
  // 角色资产授权管理
  getRoleAssets: (roleId: number) => 
    apiClient.get<{ data: RoleAssetInfo[] }>(`/roles/${roleId}/assets`),
  grantRoleAssets: (roleId: number, assetIds: number[]) =>
    apiClient.post(`/roles/${roleId}/assets`, { asset_ids: assetIds }),
  revokeRoleAssets: (roleId: number, assetIds: number[]) =>
    apiClient.delete(`/roles/${roleId}/assets`, { data: { asset_ids: assetIds } }),
  revokeSingleRoleAsset: (roleId: number, assetId: number) =>
    apiClient.delete(`/roles/${roleId}/assets/${assetId}`),
}
