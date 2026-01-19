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
      // 半透明白色背景（产生拖尾效果）- 适合浅色背景
      ctx.fillStyle = 'rgba(241, 245, 249, 0.05)'
      ctx.fillRect(0, 0, canvas.width, canvas.height)

      // Matrix 颜色 - 使用更深的灰色以增强可见度
      ctx.fillStyle = '#94A3B8'
      ctx.font = `${fontSize}px monospace`

      // 绘制每列字符
      for (let i = 0; i < drops.length; i++) {
        // 随机字符
        const char = chars[Math.floor(Math.random() * chars.length)]
        const x = i * fontSize
        const y = drops[i] * fontSize

        // 绘制字符（透明度渐变）- 增强可见度
        const opacity = Math.min(0.5, (canvas.height - y) / (fontSize * 10))
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
        <style>{`
          .login-input.ant-input,
          .login-input.ant-input-password {
            background: rgba(255, 255, 255, 0.9) !important;
            border-color: rgba(226, 232, 240, 0.8) !important;
            color: #1E293B !important;
            transition: all 0.3s ease !important;
          }
          
          .login-input.ant-input:hover,
          .login-input.ant-input-password:hover {
            border-color: rgba(203, 213, 225, 1) !important;
            box-shadow: 0 0 0 2px rgba(226, 232, 240, 0.5) !important;
          }
          
          .login-input.ant-input:focus,
          .login-input.ant-input-password:focus,
          .login-input.ant-input-focused,
          .login-input.ant-input-password-focused {
            border-color: rgba(148, 163, 184, 0.8) !important;
            box-shadow: 0 0 0 2px rgba(226, 232, 240, 0.6) !important;
            background: rgba(255, 255, 255, 1) !important;
          }
          
          .login-input.ant-input::placeholder,
          .login-input.ant-input-password input::placeholder {
            color: rgba(148, 163, 184, 0.6) !important;
          }
          
          .ant-form-item-label > label {
            color: #64748B !important;
          }
          
          .ant-form-item-explain-error {
            color: #EF4444 !important;
          }
          
          .ant-btn-loading {
            background: #64748B !important;
            border-color: #64748B !important;
            opacity: 0.8;
          }
        `}</style>
        <div 
        ref={containerRef}
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          padding: '16px',
          position: 'relative',
          background: '#F1F5F9',
          overflow: 'hidden',
        }}
      >
      {/* 网格背景 */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundImage: `
          linear-gradient(rgba(148, 163, 184, 0.3) 1px, transparent 1px),
          linear-gradient(90deg, rgba(148, 163, 184, 0.3) 1px, transparent 1px)
        `,
        backgroundSize: '40px 40px',
        opacity: 0.6,
        pointerEvents: 'none',
        zIndex: 0,
      }} />
      <Card
        style={{ 
          width: '100%',
          maxWidth: 400,
          position: 'relative',
          zIndex: 1,
          background: 'rgba(255, 255, 255, 0.7)',
          backdropFilter: 'blur(12px)',
          border: '1px solid rgba(203, 213, 225, 0.6)',
          boxShadow: '0 4px 24px rgba(0, 0, 0, 0.08)',
          animation: 'login-fade-in-up 0.6s ease-out',
          transition: 'transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.transform = 'translateY(-2px)'
          e.currentTarget.style.boxShadow = '0 8px 32px rgba(0, 0, 0, 0.12)'
          e.currentTarget.style.borderColor = 'rgba(203, 213, 225, 0.8)'
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.transform = 'translateY(0)'
          e.currentTarget.style.boxShadow = '0 4px 24px rgba(0, 0, 0, 0.08)'
          e.currentTarget.style.borderColor = 'rgba(203, 213, 225, 0.6)'
        }}
        bodyStyle={{
          padding: '32px',
        }}
      >
        <div style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          marginBottom: 24,
        }}>
          <div style={{
            padding: '12px',
            borderRadius: '8px',
            background: 'rgba(248, 250, 252, 0.6)',
            backdropFilter: 'blur(4px)',
            border: '1px solid rgba(226, 232, 240, 0.6)',
            marginBottom: 16,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}>
            <img
              src="/logo-light.svg"
              alt="KkOps"
              style={{
                height: 48,
                width: 'auto',
                opacity: 0.9,
              }}
            />
          </div>
          <div style={{
            fontSize: 18,
            fontWeight: 600,
            color: '#475569',
            letterSpacing: '0.5px',
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
            label={<span style={{ color: '#64748B' }}>用户名</span>}
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined style={{ color: '#94A3B8' }} />}
              placeholder="用户名"
              size="large"
              autoComplete="username"
              aria-label="用户名输入框"
              style={{
                background: 'rgba(255, 255, 255, 0.9)',
                borderColor: 'rgba(226, 232, 240, 0.8)',
                color: '#1E293B',
              }}
              className="login-input"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label={<span style={{ color: '#64748B' }}>密码</span>}
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#94A3B8' }} />}
              placeholder="密码"
              size="large"
              autoComplete="current-password"
              aria-label="密码输入框"
              style={{
                background: 'rgba(255, 255, 255, 0.9)',
                borderColor: 'rgba(226, 232, 240, 0.8)',
                color: '#1E293B',
              }}
              className="login-input"
            />
          </Form.Item>

          <Form.Item>
            <Button
              htmlType="submit"
              loading={loading}
              block
              size="large"
              style={{
                background: '#64748B',
                borderColor: '#64748B',
                color: '#FFFFFF',
                fontWeight: 500,
                boxShadow: '0 2px 8px rgba(100, 116, 139, 0.2)',
                transition: 'all 0.3s ease',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.background = '#475569'
                e.currentTarget.style.boxShadow = '0 4px 12px rgba(100, 116, 139, 0.3)'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = '#64748B'
                e.currentTarget.style.boxShadow = '0 2px 8px rgba(100, 116, 139, 0.2)'
              }}
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
