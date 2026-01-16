// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect } from 'react'
import { Layout, Menu, theme, Dropdown, message, Button } from 'antd'
import type { MenuProps } from 'antd'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  DashboardOutlined,
  DatabaseOutlined,
  PlayCircleOutlined,
  ConsoleSqlOutlined,
  UserOutlined,
  TeamOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MoonOutlined,
  SunOutlined,
  LogoutOutlined,
  FolderOutlined,
  GlobalOutlined,
  CloudOutlined,
} from '@ant-design/icons'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'

const { Header, Sider, Content } = Layout

interface MainLayoutProps {
  children: React.ReactNode
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false)
  const [isMobile, setIsMobile] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { mode, toggleMode } = useThemeStore()
  const { user, clearAuth } = useAuthStore()
  const {
    token: { colorBgContainer },
  } = theme.useToken()

  // Responsive handling
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768)
      if (window.innerWidth < 768) {
        setCollapsed(true)
      }
    }
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: '/projects',
      icon: <FolderOutlined />,
      label: '项目管理',
    },
    {
      key: '/environments',
      icon: <GlobalOutlined />,
      label: '环境管理',
    },
    {
      key: '/cloud-platforms',
      icon: <CloudOutlined />,
      label: '云平台管理',
    },
    {
      key: '/assets',
      icon: <DatabaseOutlined />,
      label: '资产管理',
    },
    {
      key: '/tasks',
      icon: <PlayCircleOutlined />,
      label: '运维执行',
    },
    {
      key: '/task-templates',
      icon: <PlayCircleOutlined />,
      label: '任务模板',
    },
    {
      key: '/ssh/keys',
      icon: <ConsoleSqlOutlined />,
      label: 'SSH 密钥',
    },
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/roles',
      icon: <TeamOutlined />,
      label: '角色权限',
    },
  ]

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key)
  }

  const handleLogout = async () => {
    try {
      // Call logout API (optional, for API consistency)
      await authApi.logout()
    } catch (error) {
      // Ignore API errors, proceed with client-side logout
    }
    
    // Clear authentication state
    clearAuth()
    
    // Show success message
    message.success('已成功登出')
    
    // Navigate to login page
    navigate('/login')
  }

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '登出',
      onClick: handleLogout,
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={240}
        collapsedWidth={isMobile ? 0 : 80}
        breakpoint="lg"
        onBreakpoint={(broken) => {
          setIsMobile(broken)
          setCollapsed(broken)
        }}
        style={{
          background: mode === 'dark' ? '#1E293B' : '#FFFFFF',
          borderRight: `1px solid ${mode === 'dark' ? '#334155' : '#E2E8F0'}`,
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: `1px solid ${mode === 'dark' ? '#334155' : '#E2E8F0'}`,
            padding: collapsed ? '0 16px' : '0 20px',
          }}
        >
          {collapsed ? (
            <img
              src="/icon.svg"
              alt="KkOps"
              style={{
                width: 32,
                height: 32,
                color: mode === 'dark' ? '#3B82F6' : '#2563EB',
              }}
            />
          ) : (
            <img
              src={mode === 'dark' ? '/logo-dark.svg' : '/logo-light.svg'}
              alt="KkOps"
              style={{
                height: 28,
                width: 'auto',
              }}
            />
          )}
        </div>
        <Menu
          theme={mode}
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{
            background: mode === 'dark' ? '#1E293B' : '#FFFFFF',
            border: 'none',
          }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: '0 24px',
            background: colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: `1px solid ${mode === 'dark' ? '#334155' : '#E2E8F0'}`,
            position: 'sticky',
            top: 0,
            zIndex: 1000,
          }}
        >
        <div
          role="button"
          tabIndex={0}
          aria-label={collapsed ? '展开菜单' : '折叠菜单'}
          style={{
            fontSize: 18,
            cursor: 'pointer',
            color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
          }}
          onClick={() => setCollapsed(!collapsed)}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault()
              setCollapsed(!collapsed)
            }
          }}
        >
          {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
        </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Button
              type="text"
              icon={<ConsoleSqlOutlined />}
              onClick={() => window.open('/ssh/terminal', '_blank')}
              title="WebSSH 终端（在新标签页打开）"
              style={{
                color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
              }}
            >
              WebSSH
            </Button>
            <div
              role="button"
              tabIndex={0}
              aria-label={mode === 'dark' ? '切换到亮色模式' : '切换到暗色模式'}
              style={{
                fontSize: 20,
                cursor: 'pointer',
                color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
              }}
              onClick={toggleMode}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault()
                  toggleMode()
                }
              }}
            >
              {mode === 'dark' ? <SunOutlined /> : <MoonOutlined />}
            </div>
            {user && (
              <Dropdown
                menu={{ items: userMenuItems }}
                placement="bottomRight"
                trigger={['click']}
              >
                <div
                  role="button"
                  tabIndex={0}
                  aria-label="用户菜单"
                  aria-haspopup="menu"
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 8,
                    cursor: 'pointer',
                    padding: '4px 12px',
                    borderRadius: 4,
                    color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
                    fontSize: 14,
                    transition: 'background-color 0.2s',
                  }}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                      e.preventDefault()
                      // Dropdown will handle the click
                    }
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor =
                      mode === 'dark' ? '#334155' : '#F1F5F9'
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'transparent'
                  }}
                >
                  <UserOutlined />
                  <span>{user.real_name || user.username}</span>
                </div>
              </Dropdown>
            )}
          </div>
        </Header>
        <Content
          style={{
            margin: isMobile ? '16px 8px' : '24px',
            padding: isMobile ? 16 : 24,
            minHeight: 280,
            background: colorBgContainer,
            borderRadius: 8,
          }}
        >
          {children}
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
