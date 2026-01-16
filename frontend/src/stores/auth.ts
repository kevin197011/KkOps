// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { UserInfo } from '@/api/auth'

interface AuthState {
  token: string | null
  user: UserInfo | null
  isAuthenticated: boolean
  setAuth: (token: string, user: UserInfo) => void
  clearAuth: () => void
}

export const useAuthStore = create<AuthState>((set) => {
  // Initialize from localStorage
  const token = localStorage.getItem('token')
  const userStr = localStorage.getItem('user')
  let user: UserInfo | null = null
  if (userStr) {
    try {
      user = JSON.parse(userStr)
    } catch {
      // Ignore parse errors
    }
  }

  return {
    token,
    user,
    isAuthenticated: !!token && !!user,
    setAuth: (token: string, user: UserInfo) => {
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
      set({ token, user, isAuthenticated: true })
    },
    clearAuth: () => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      set({ token: null, user: null, isAuthenticated: false })
    },
  }
})
