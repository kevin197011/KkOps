// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useRef } from 'react'
import { Form, Input, Button, Card, message, ConfigProvider } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { authApi, LoginRequest } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { lightTheme } from '@/themes'

const Login = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)
  const containerRef = useRef<HTMLDivElement>(null)
  const particlesRef = useRef<HTMLDivElement[]>([])

  // 创建粒子节点
  useEffect(() => {
    if (!containerRef.current) return

    const container = containerRef.current
    const particleCount = 15
    particlesRef.current = []

    for (let i = 0; i < particleCount; i++) {
      const particle = document.createElement('div')
      const x = Math.random() * 100
      const y = Math.random() * 100
      const dx = (Math.random() - 0.5) * 40
      const dy = (Math.random() - 0.5) * 40
      const delay = Math.random() * 5
      const duration = 8 + Math.random() * 4

      particle.style.cssText = `
        position: absolute;
        left: ${x}%;
        top: ${y}%;
        width: 4px;
        height: 4px;
        background: rgba(59, 130, 246, 0.8);
        border-radius: 50%;
        box-shadow: 0 0 10px rgba(59, 130, 246, 0.6);
        pointer-events: none;
        --dx: ${dx}px;
        --dy: ${dy}px;
        animation: login-particle-float ${duration}s ease-in-out infinite;
        animation-delay: ${delay}s;
      `
      container.appendChild(particle)
      particlesRef.current.push(particle)
    }

    return () => {
      particlesRef.current.forEach(particle => {
        if (particle.parentNode) {
          particle.parentNode.removeChild(particle)
        }
      })
      particlesRef.current = []
    }
  }, [])

  const onFinish = async (values: LoginRequest) => {
    setLoading(true)
    try {
      const response = await authApi.login(values)
      const { token, user } = response.data

      // Store auth state
      setAuth(token, user)

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
      <div 
        ref={containerRef}
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          padding: '16px',
          position: 'relative',
          background: 'linear-gradient(135deg, #0F172A 0%, #1E293B 50%, #0F172A 100%)',
          overflow: 'hidden',
        }}
      >
      {/* 扫描线效果 */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        height: '2px',
        background: 'linear-gradient(90deg, transparent, rgba(59, 130, 246, 0.5), transparent)',
        animation: 'login-scanline 3s linear infinite',
        pointerEvents: 'none',
        zIndex: 0,
      }} />

      {/* 数据流动效果 */}
      {[...Array(3)].map((_, i) => (
        <div
          key={`data-flow-${i}`}
          style={{
            position: 'absolute',
            top: `${20 + i * 30}%`,
            left: 0,
            width: '200px',
            height: '1px',
            background: 'linear-gradient(90deg, transparent, rgba(37, 99, 235, 0.6), transparent)',
            boxShadow: '0 0 10px rgba(37, 99, 235, 0.4)',
            animation: `login-data-flow ${8 + i * 2}s linear infinite`,
            animationDelay: `${i * 2}s`,
            pointerEvents: 'none',
            zIndex: 0,
          }}
        />
      ))}
      {/* 网格背景 */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundImage: `
          linear-gradient(rgba(37, 99, 235, 0.1) 1px, transparent 1px),
          linear-gradient(90deg, rgba(37, 99, 235, 0.1) 1px, transparent 1px)
        `,
        backgroundSize: '40px 40px',
        opacity: 0.5,
        pointerEvents: 'none',
        animation: 'login-fade-in 1s ease-out, login-grid-move 20s linear infinite',
      }} />
      
      {/* 渐变光晕效果 */}
      <div style={{
        position: 'absolute',
        top: '-50%',
        left: '-50%',
        width: '200%',
        height: '200%',
        background: 'radial-gradient(circle, rgba(37, 99, 235, 0.15) 0%, transparent 70%)',
        animation: 'login-pulse 8s ease-in-out infinite',
        pointerEvents: 'none',
      }} />
      
      {/* 右上角光晕 */}
      <div style={{
        position: 'absolute',
        top: '-30%',
        right: '-30%',
        width: '60%',
        height: '60%',
        background: 'radial-gradient(circle, rgba(59, 130, 246, 0.2) 0%, transparent 70%)',
        borderRadius: '50%',
        filter: 'blur(60px)',
        pointerEvents: 'none',
        animation: 'login-float 15s ease-in-out infinite, login-glow-pulse 4s ease-in-out infinite',
        animationDelay: '0s, 1s',
      }} />
      
      {/* 左下角光晕 */}
      <div style={{
        position: 'absolute',
        bottom: '-30%',
        left: '-30%',
        width: '60%',
        height: '60%',
        background: 'radial-gradient(circle, rgba(37, 99, 235, 0.15) 0%, transparent 70%)',
        borderRadius: '50%',
        filter: 'blur(60px)',
        pointerEvents: 'none',
        animation: 'login-float 18s ease-in-out infinite reverse, login-glow-pulse 5s ease-in-out infinite',
        animationDelay: '0s, 2s',
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
      </div>
    </ConfigProvider>
  )
}

export default Login
