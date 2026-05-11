// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Spin } from 'antd'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { userApi } from '@/api/user'

/**
 * Handles redirect from SSO callback: reads token from query, stores auth, fetches user and permissions, then redirects to dashboard.
 */
const AuthCallback = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)
  const setAuth = useAuthStore((state) => state.setAuth)
  const setPermissions = usePermissionStore((state) => state.setPermissions)

  useEffect(() => {
    const token = searchParams.get('token')
    if (!token) {
      setError('Missing token')
      return
    }

    const run = async () => {
      try {
        localStorage.setItem('token', token)
        const meRes = await authApi.getMe()
        const user = meRes.data
        setAuth(token, user)

        try {
          const permRes = await userApi.getPermissions()
          setPermissions(permRes.data.permissions || [])
        } catch {
          setPermissions([])
        }

        navigate('/dashboard', { replace: true })
      } catch (e: any) {
        setError(e.response?.data?.error || 'Failed to complete login')
      }
    }

    run()
  }, [searchParams, setAuth, setPermissions, navigate])

  if (error) {
    return (
      <div style={{ padding: 24, textAlign: 'center', color: '#ef4444' }}>
        {error}
        <br />
        <a href="/login">Return to login</a>
      </div>
    )
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
      <Spin size="large" tip="Completing login..." />
    </div>
  )
}

export default AuthCallback
