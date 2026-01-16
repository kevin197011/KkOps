// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

export interface UserInfo {
  id: number
  username: string
  email: string
  real_name: string
  roles: string[]
}

export const authApi = {
  login: (data: LoginRequest) => apiClient.post<LoginResponse>('/auth/login', data),
  getMe: () => apiClient.get<UserInfo>('/auth/me'),
  logout: () => apiClient.post<{ message: string }>('/auth/logout'),
}
