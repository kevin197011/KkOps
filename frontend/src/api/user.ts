// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface User {
  id: number
  username: string
  email: string
  phone?: string
  real_name: string
  department_id?: number
  status: string
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  email: string
  password: string
  phone?: string
  real_name?: string
  department_id?: number
}

export interface UpdateUserRequest {
  email?: string
  phone?: string
  real_name?: string
  department_id?: number
  status?: string
}

export interface ListUsersResponse {
  data: User[]
  total: number
  page: number
  size: number
}

// 用户角色信息
export interface UserRoleInfo {
  id: number
  name: string
  description: string
  is_admin: boolean
}

export const userApi = {
  list: (page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<ListUsersResponse>(`/users${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<User>(`/users/${id}`),
  create: (data: CreateUserRequest) => apiClient.post<User>('/users', data),
  update: (id: number, data: UpdateUserRequest) => apiClient.put<User>(`/users/${id}`, data),
  delete: (id: number) => apiClient.delete(`/users/${id}`),
  resetPassword: (id: number, password: string) =>
    apiClient.post(`/users/${id}/reset-password`, { password }),
  
  // 用户角色管理
  getUserRoles: (userId: number) =>
    apiClient.get<{ data: UserRoleInfo[] }>(`/users/${userId}/roles`),
  setUserRoles: (userId: number, roleIds: number[]) =>
    apiClient.post(`/users/${userId}/roles`, { role_ids: roleIds }),
  removeUserRole: (userId: number, roleId: number) =>
    apiClient.delete(`/users/${userId}/roles/${roleId}`),
  
  // 获取当前用户权限
  getPermissions: () =>
    apiClient.get<{ permissions: string[] }>('/user/permissions'),
}
