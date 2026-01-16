// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useAuthStore } from '@/stores/auth'
import { UserInfo } from '@/api/auth'

export const useAuth = () => {
  const { token, user, isAuthenticated, setAuth, clearAuth } = useAuthStore()

  return {
    token,
    user,
    isAuthenticated,
    setAuth,
    clearAuth,
  }
}

export const useUser = (): UserInfo | null => {
  return useAuthStore((state) => state.user)
}

export const useHasRole = (role: string): boolean => {
  const user = useAuthStore((state) => state.user)
  return user?.roles?.includes(role) ?? false
}

export const useHasPermission = (permission: string): boolean => {
  // TODO: Implement permission checking based on user roles
  // This will be implemented when permission system is fully integrated
  return true
}
