// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useRef } from 'react'
import { Form, Input, Button, Card, ConfigProvider, App } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { authApi, LoginRequest } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { userApi } from '@/api/user'
import { lightTheme } from '@/themes'

const Login = () => {
  const { message } = App.useApp()
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)
  const setPermissions = usePermissionStore((state) => state.setPermissions)
  const containerRef = useRef<HTMLDivElement>(null)
  const canvasRef = useRef<HTMLCanvasElement | null>(null)

  // Matrix 代码雨效果
  useEffect(() => {
    if (!containerRef.current) return

    const container = containerRef.current
    const canvas = document.createElement('canvas')
    canvas.style.position = 'absolute'
    canvas.style.top = '0'
    canvas.style.left = '0'
    canvas.style.width = '100%'
    canvas.style.height = '100%'
    canvas.style.pointerEvents = 'none'
    canvas.style.zIndex = '0'
    canvas.width = container.offsetWidth || window.innerWidth
    canvas.height = container.offsetHeight || window.innerHeight
    container.appendChild(canvas)
    canvasRef.current = canvas

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    // Matrix 字符集
    const chars = '01アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン'
    const fontSize = 14
    let columns = Math.floor(canvas.width / fontSize)
    const drops: number[] = []
    const speeds: number[] = []

    // 初始化每列的起始位置和速度
    const initColumns = () => {
      columns = Math.floor(canvas.width / fontSize)
      while (drops.length < columns) {
        drops.push(Math.random() * -100)
        speeds.push(0.5 + Math.random() * 2)
      }
      drops.splice(columns)
      speeds.splice(columns)
    }

    initColumns()

    const draw = () => {
      // 半透明黑色背景（产生拖尾效果）
      ctx.fillStyle = 'rgba(0, 0, 0, 0.05)'
      ctx.fillRect(0, 0, canvas.width, canvas.height)

      // Matrix 绿色
      ctx.fillStyle = '#00ff41'
      ctx.font = `${fontSize}px monospace`

      // 绘制每列字符
      for (let i = 0; i < drops.length; i++) {
        // 随机字符
        const char = chars[Math.floor(Math.random() * chars.length)]
        const x = i * fontSize
        const y = drops[i] * fontSize

        // 绘制字符（透明度渐变）
        const opacity = Math.min(1, (canvas.height - y) / (fontSize * 10))
        ctx.globalAlpha = opacity
        ctx.fillText(char, x, y)

        // 重置位置并随机速度
        if (y > canvas.height && Math.random() > 0.975) {
          drops[i] = 0
          speeds[i] = 0.5 + Math.random() * 2
        }

        // 移动位置
        drops[i] += speeds[i]
      }

      ctx.globalAlpha = 1
      requestAnimationFrame(draw)
    }

    // 窗口大小变化时更新画布尺寸
    const handleResize = () => {
      canvas.width = container.offsetWidth || window.innerWidth
      canvas.height = container.offsetHeight || window.innerHeight
      initColumns()
    }

    window.addEventListener('resize', handleResize)
    draw()

    return () => {
      window.removeEventListener('resize', handleResize)
      if (canvas.parentNode) {
        canvas.parentNode.removeChild(canvas)
      }
      canvasRef.current = null
    }
  }, [])

  const onFinish = async (values: LoginRequest) => {
    setLoading(true)
    try {
      const response = await authApi.login(values)
      const { token, user } = response.data

      // Store auth state
      setAuth(token, user)

      // Fetch and store user permissions
      try {
        const permResponse = await userApi.getPermissions()
        setPermissions(permResponse.data.permissions || [])
      } catch (permError) {
        // If permission fetch fails, continue but log error
        console.error('Failed to fetch user permissions:', permError)
        // Set empty permissions (user will see limited menu)
        setPermissions([])
      }

      message.success('登录成功')
      navigate('/dashboard')
    } catch (error: any) {
      message.error(error.response?.data?.error || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <ConfigProvider theme={lightTheme}>
      <App>
        <div 
        ref={containerRef}
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          padding: '16px',
          position: 'relative',
          background: '#000000',
          overflow: 'hidden',
        }}
      >
      {/* Matrix 网格背景 */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundImage: `
          linear-gradient(rgba(0, 255, 65, 0.05) 1px, transparent 1px),
          linear-gradient(90deg, rgba(0, 255, 65, 0.05) 1px, transparent 1px)
        `,
        backgroundSize: '40px 40px',
        opacity: 0.3,
        pointerEvents: 'none',
        zIndex: 0,
      }} />
      <Card
        style={{ 
          width: '100%',
          maxWidth: 400,
          position: 'relative',
          zIndex: 1,
          background: 'rgba(255, 255, 255, 0.95)',
          backdropFilter: 'blur(10px)',
          border: '1px solid rgba(226, 232, 240, 0.8)',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
          animation: 'login-fade-in-up 0.6s ease-out',
          transition: 'transform 0.3s ease, box-shadow 0.3s ease',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.transform = 'translateY(-2px)'
          e.currentTarget.style.boxShadow = '0 12px 40px rgba(0, 0, 0, 0.15)'
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.transform = 'translateY(0)'
          e.currentTarget.style.boxShadow = '0 8px 32px rgba(0, 0, 0, 0.1)'
        }}
      >
        <div style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          marginBottom: 24,
        }}>
          <img
            src="/logo-light.svg"
            alt="KkOps"
            style={{
              height: 40,
              width: 'auto',
              marginBottom: 8,
            }}
          />
          <div style={{
            fontSize: 18,
            fontWeight: 600,
            color: '#1E293B',
          }}>
            智能运维管理平台
          </div>
        </div>
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名"
              size="large"
              autoComplete="username"
              aria-label="用户名输入框"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
              size="large"
              autoComplete="current-password"
              aria-label="密码输入框"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
            >
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
      
      {/* Footer */}
      <div
        style={{
          position: 'absolute',
          bottom: 16,
          left: 0,
          right: 0,
          textAlign: 'center',
          color: 'rgba(148, 163, 184, 0.8)',
          fontSize: 12,
          zIndex: 1,
        }}
      >
        系统运行部驱动 © {new Date().getFullYear()} KkOps
      </div>
      </div>
      </App>
    </ConfigProvider>
  )
}

export default Login
