// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { create } from 'zustand'

interface PermissionState {
  permissions: string[]
  isAdmin: boolean
  setPermissions: (permissions: string[]) => void
  hasPermission: (resource: string, action: string) => boolean
  hasAnyPermission: (permissionList: string[]) => boolean
  clearPermissions: () => void
}

export const usePermissionStore = create<PermissionState>((set, get) => ({
  permissions: [],
  isAdmin: false,
  setPermissions: (permissions: string[]) => {
    // Check if user is admin (has all menu permissions or is_admin role)
    // For now, we'll check if user has all permissions (admin users get all permissions from backend)
    const allMenuResources = [
      'dashboard',
      'projects',
      'environments',
      'cloud-platforms',
      'assets',
      'executions',
      'templates',
      'tasks',
      'deployments',
      'ssh-keys',
      'users',
      'roles',
      'audit-logs',
    ]
    
    // If user has all permissions, consider them admin
    const hasAllPermissions = allMenuResources.every((resource) =>
      permissions.some((p) => p.startsWith(`${resource}:`))
    )
    
    set({ permissions, isAdmin: hasAllPermissions })
  },
  hasPermission: (resource: string, action: string) => {
    const { permissions, isAdmin } = get()
    if (isAdmin) {
      return true
    }
    
    // Check for exact match or wildcard
    return permissions.some(
      (p) => p === `${resource}:${action}` || p === `${resource}:*`
    )
  },
  hasAnyPermission: (permissionList: string[]) => {
    const { permissions, isAdmin } = get()
    if (isAdmin) {
      return true
    }
    
    return permissionList.some((perm) => permissions.includes(perm))
  },
  clearPermissions: () => {
    set({ permissions: [], isAdmin: false })
  },
}))
