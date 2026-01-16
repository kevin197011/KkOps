// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Navigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'

interface ProtectedRouteProps {
  children: React.ReactNode
  requiredPermission?: string
}

const ProtectedRoute = ({ children, requiredPermission }: ProtectedRouteProps) => {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  // TODO: Check permission if requiredPermission is provided
  // if (requiredPermission && !hasPermission(requiredPermission)) {
  //   return <Navigate to="/unauthorized" replace />
  // }

  return <>{children}</>
}

export default ProtectedRoute
