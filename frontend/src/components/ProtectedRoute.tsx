// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Navigate } from 'react-router-dom'
import { Result, Button, Spin } from 'antd'
import { useAuth } from '@/hooks/useAuth'
import { usePermissionStore } from '@/stores/permission'
import { useState, useEffect } from 'react'
import { userApi } from '@/api/user'

interface ProtectedRouteProps {
  children: React.ReactNode
  requiredPermission?: string // Format: "resource:action"
}

const ProtectedRoute = ({ children, requiredPermission }: ProtectedRouteProps) => {
  const { isAuthenticated, user } = useAuth()
  const { hasPermission, permissions, isAdmin, setPermissions } = usePermissionStore()
  const [permissionsLoaded, setPermissionsLoaded] = useState(false)

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  // Load permissions if not loaded yet
  useEffect(() => {
    if (user && !permissionsLoaded) {
      userApi.getPermissions()
        .then((response) => {
          setPermissions(response.data.permissions || [])
          setPermissionsLoaded(true)
        })
        .catch((error) => {
          console.error('Failed to fetch user permissions:', error)
          // If permission fetch fails, set empty permissions and mark as loaded
          // Backend will still check permissions
          setPermissions([])
          setPermissionsLoaded(true)
        })
    } else if (!user) {
      setPermissionsLoaded(false)
    }
  }, [user, permissionsLoaded, setPermissions])

  // Show loading state only briefly while fetching permissions
  // After a short timeout, allow access (backend will verify permissions)
  const [showLoading, setShowLoading] = useState(true)
  useEffect(() => {
    if (requiredPermission && !permissionsLoaded) {
      // Show loading for max 2 seconds, then allow access
      const timer = setTimeout(() => {
        setShowLoading(false)
      }, 2000)
      return () => clearTimeout(timer)
    } else {
      setShowLoading(false)
    }
  }, [requiredPermission, permissionsLoaded])

  if (showLoading && requiredPermission && !permissionsLoaded) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' }}>
        <Spin size="large" />
      </div>
    )
  }

  // Check permission if requiredPermission is provided and permissions are loaded
  // If permissions failed to load or not loaded yet, allow access and let backend verify
  // This ensures admin users are not blocked while permissions are being fetched
  if (requiredPermission && permissionsLoaded && permissions.length > 0) {
    const [resource, action] = requiredPermission.split(':')
    if (!hasPermission(resource, action)) {
      return (
        <Result
          status="403"
          title="403"
          subTitle="抱歉，您没有权限访问此页面。"
          extra={
            <Button type="primary" onClick={() => window.history.back()}>
              返回
            </Button>
          }
        />
      )
    }
  }
  
  // Allow access if:
  // 1. No permission required
  // 2. Permissions not loaded yet (backend will verify)
  // 3. Permission check passed
  // This ensures backend always verifies permissions, especially for admin users

  // Allow access if no permission required or permission check passed
  // Backend will still verify permissions
  return <>{children}</>
}

export default ProtectedRoute
