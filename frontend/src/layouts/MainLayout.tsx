// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect } from 'react'
import { Layout, Menu, theme, Dropdown, message, Button, Modal, Form, Input } from 'antd'
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
  MoonFilled,
  SunFilled,
  LogoutOutlined,
  FolderOutlined,
  GlobalOutlined,
  CloudOutlined,
  RocketOutlined,
  LockOutlined,
  ScheduleOutlined,
  FileTextOutlined,
  AuditOutlined,
  ThunderboltOutlined,
  SafetyOutlined,
  SettingOutlined,
  ApartmentOutlined,
  AppstoreOutlined,
  BookOutlined,
  TagsOutlined,
  LinkOutlined,
  SafetyCertificateOutlined,
  SearchOutlined,
  CloudSyncOutlined,
  ApiOutlined,
  LineChartOutlined,
  ContainerOutlined,
  BellOutlined,
  AlertOutlined,
  ClusterOutlined,
  RobotOutlined,
  DeploymentUnitOutlined,
  ShareAltOutlined,
} from '@ant-design/icons'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { authApi } from '@/api/auth'
import { userApi } from '@/api/user'
import { CommandPalette } from '@/components/shell'
import { flattenMenuToCommands } from '@/navigation/flattenMenu'

const { Header, Sider, Content, Footer } = Layout

interface MainLayoutProps {
  children: React.ReactNode
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false)
  const [isMobile, setIsMobile] = useState(false)
  const [passwordModalVisible, setPasswordModalVisible] = useState(false)
  const [paletteOpen, setPaletteOpen] = useState(false)
  const [passwordLoading, setPasswordLoading] = useState(false)
  const [passwordForm] = Form.useForm()
  const [openKeys, setOpenKeys] = useState<string[]>([
    'ai-ops',
    'integration-suite',
    'infrastructure',
    'operations',
    'security',
    'system',
  ])
  const navigate = useNavigate()
  const location = useLocation()
  const { mode, toggleMode } = useThemeStore()
  const { user, clearAuth } = useAuthStore()
  const { hasPermission, isAdmin, permissions } = usePermissionStore()

  // 获取用户权限（如果还没有加载）
  useEffect(() => {
    if (user && permissions.length === 0) {
      userApi.getPermissions()
        .then((response) => {
          const { setPermissions } = usePermissionStore.getState()
          setPermissions(response.data.permissions || [])
        })
        .catch((error) => {
          console.error('Failed to fetch user permissions:', error)
        })
    }
  }, [user, permissions.length])
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

  // 菜单权限映射
  const menuPermissionMap: Record<string, string> = {
    '/dashboard': 'dashboard:read',
    '/operation-tools': 'operation-tools:read',
    '/projects': 'projects:*',
    '/environments': 'environments:*',
    '/cloud-platforms': 'cloud-platforms:*',
    '/assets': 'assets:*',
    '/executions': 'executions:*',
    '/templates': 'templates:*',
    '/tasks': 'tasks:*',
    '/deployments': 'deployments:*',
    '/ssh/keys': 'ssh-keys:*',
    '/users': 'users:*',
    '/roles': 'roles:*',
    '/audit-logs': 'audit-logs:read',
    '/connection-audit': 'connection-audit:read',
    '/external-systems': 'external-systems:*',
    '/oauth2-clients': 'oauth2-clients:*',
    '/provisioning': 'provisioning:*',
    '/integrations': 'integrations:*',
    '/monitoring': 'monitoring:*',
    '/logging': 'logging:*',
    '/cicd': 'cicd:*',
    '/registry': 'registry:*',
    '/gitops': 'gitops:*',
    '/gitops/pipeline': 'gitops:*',
    '/k8s/clusters': 'kubernetes:*',
    '/alerts': 'alerts:*',
    '/incidents': 'incidents:*',
    '/ai/chat': 'ai:*',
    '/ai/anomaly/rules': 'ai:*',
    '/ai/anomaly/findings': 'ai:*',
    '/ai/rca/reports': 'ai:*',
    '/cmdb': 'cmdb:*',
    '/topology': 'cmdb:*',
  }

  // 根据权限过滤菜单项
  const filterMenuItems = (items: MenuProps['items']): MenuProps['items'] => {
    if (!items) return []
    
    return items
      .map((item) => {
        if (!item) return null
        
        // 处理分隔符
        if (item.type === 'divider') {
          return item
        }
        
        // 如果是子菜单，过滤子项
        if ('children' in item && item.children && Array.isArray(item.children)) {
          const filteredChildren = filterMenuItems(item.children) ?? []
          // 如果子菜单没有可用项，不显示父菜单
          if (filteredChildren.length === 0) {
            return null
          }
          return {
            ...item,
            children: filteredChildren,
          }
        }
        
        // 检查权限
        if (item.key && typeof item.key === 'string' && item.key.startsWith('/')) {
          const requiredPermission = menuPermissionMap[item.key]
          if (requiredPermission) {
            const [resource, action] = requiredPermission.split(':')
            const hasRequired = hasPermission(resource, action)
            if (!isAdmin && !hasRequired) {
              return null // 无权限，不显示
            }
          }
        }
        
        return item
      })
      .filter((item) => item !== null) as MenuProps['items']
  }

  const allMenuItems: MenuProps['items'] = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      type: 'divider',
    },
    {
      key: '/operation-tools',
      icon: <AppstoreOutlined />,
      label: '运维导航',
    },
    {
      key: 'ai-ops',
      icon: <RobotOutlined />,
      label: 'AI 运维',
      children: [
        {
          key: '/ai/chat',
          icon: <RobotOutlined />,
          label: '助手',
        },
        {
          key: '/ai/anomaly/rules',
          icon: <AlertOutlined />,
          label: '异常规则',
        },
        {
          key: '/ai/anomaly/findings',
          icon: <BellOutlined />,
          label: '异常结果',
        },
        {
          key: '/ai/rca/reports',
          icon: <FileTextOutlined />,
          label: '根因报告',
        },
      ],
    },
    {
      key: '/external-systems',
      icon: <LinkOutlined />,
      label: 'SSO 应用门户',
    },
    {
      key: '/oauth2-clients',
      icon: <SafetyCertificateOutlined />,
      label: 'IdP 应用',
    },
    {
      key: 'integration-suite',
      icon: <ApiOutlined />,
      label: '集成中心',
      children: [
        {
          key: '/integrations',
          icon: <ApiOutlined />,
          label: '连接器',
        },
        {
          key: '/monitoring',
          icon: <LineChartOutlined />,
          label: '监控',
        },
        {
          key: '/logging',
          icon: <FileTextOutlined />,
          label: '日志',
        },
        {
          key: '/cicd',
          icon: <RocketOutlined />,
          label: 'CI/CD',
        },
        {
          key: '/registry',
          icon: <ContainerOutlined />,
          label: '镜像仓库',
        },
        {
          key: '/gitops',
          icon: <CloudSyncOutlined />,
          label: 'GitOps',
        },
        {
          key: '/gitops/pipeline',
          icon: <ScheduleOutlined />,
          label: '流水线',
        },
        {
          key: '/k8s/clusters',
          icon: <ClusterOutlined />,
          label: 'Kubernetes',
        },
        {
          key: '/alerts',
          icon: <BellOutlined />,
          label: '告警',
        },
        {
          key: '/incidents',
          icon: <AlertOutlined />,
          label: '事件',
        },
      ],
    },
    {
      key: 'infrastructure',
      icon: <ApartmentOutlined />,
      label: '基础设施',
      children: [
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
          key: '/cmdb',
          icon: <DeploymentUnitOutlined />,
          label: 'CMDB',
        },
        {
          key: '/topology',
          icon: <ShareAltOutlined />,
          label: '拓扑',
        },
        {
          key: '/tags',
          icon: <TagsOutlined />,
          label: '标签管理',
        },
      ],
    },
    {
      key: 'operations',
      icon: <ThunderboltOutlined />,
      label: '任务管理',
      children: [
        {
          key: '/executions',
          icon: <PlayCircleOutlined />,
          label: '运维执行',
        },
        {
          key: '/templates',
          icon: <FileTextOutlined />,
          label: '任务模板',
        },
        {
          key: '/tasks',
          icon: <ScheduleOutlined />,
          label: '任务执行',
        },
        {
          key: '/deployments',
          icon: <RocketOutlined />,
          label: '部署管理',
        },
      ],
    },
    {
      key: 'security',
      icon: <SafetyOutlined />,
      label: '安全管理',
      children: [
        {
          key: '/ssh/keys',
          icon: <ConsoleSqlOutlined />,
          label: 'SSH 密钥',
        },
      ],
    },
    {
      key: 'system',
      icon: <SettingOutlined />,
      label: '系统管理',
      children: [
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
        {
          key: '/audit-logs',
          icon: <AuditOutlined />,
          label: '审计日志',
        },
        {
          key: '/connection-audit',
          icon: <LinkOutlined />,
          label: '审计连线',
        },
        {
          key: '/provisioning',
          icon: <CloudSyncOutlined />,
          label: '账号同步',
        },
      ],
    },
  ]

  // 根据权限过滤菜单
  const menuItems = isAdmin || permissions.length === 0 
    ? allMenuItems 
    : filterMenuItems(allMenuItems)

  const paletteCommands = flattenMenuToCommands(menuItems)

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault()
        setPaletteOpen((open) => !open)
      }
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'l') {
        e.preventDefault()
        navigate('/ai/chat')
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [navigate])

  const handleMenuClick = ({ key }: { key: string }) => {
    // 只导航到具体的路由，忽略分类键（如 'infrastructure', 'operations' 等）
    if (key.startsWith('/')) {
      navigate(key)
    }
  }

  const handleLogout = async () => {
    try {
      // Call logout API (optional, for API consistency)
      await authApi.logout()
    } catch (error) {
      // Ignore API errors, proceed with client-side logout
    }
    
    // Clear authentication and permission state
    clearAuth()
    const { clearPermissions } = usePermissionStore.getState()
    clearPermissions()
    
    // Show success message
    message.success('已成功登出')
    
    // Navigate to login page
    navigate('/login')
  }

  const handleChangePassword = async (values: { old_password: string; new_password: string }) => {
    setPasswordLoading(true)
    try {
      await authApi.changePassword(values)
      message.success('密码修改成功')
      setPasswordModalVisible(false)
      passwordForm.resetFields()
    } catch (error: any) {
      message.error(error.response?.data?.error || '密码修改失败')
    } finally {
      setPasswordLoading(false)
    }
  }

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人设置',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'change-password',
      icon: <LockOutlined />,
      label: '修改密码',
      onClick: () => setPasswordModalVisible(true),
    },
    {
      type: 'divider',
    },
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
          role="button"
          tabIndex={0}
          aria-label="返回仪表盘"
          onClick={() => navigate('/dashboard')}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault()
              navigate('/dashboard')
            }
          }}
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: `1px solid ${mode === 'dark' ? '#334155' : '#E2E8F0'}`,
            padding: collapsed ? '0 16px' : '0 20px',
            cursor: 'pointer',
            transition: 'background-color 0.2s',
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = mode === 'dark' ? '#334155' : '#F1F5F9'
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = 'transparent'
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
          openKeys={openKeys}
          onOpenChange={setOpenKeys}
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
              icon={<SearchOutlined />}
              onClick={() => setPaletteOpen(true)}
              title="Jump to page (⌘K / Ctrl+K)"
              style={{
                color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
              }}
            >
              Jump
            </Button>
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
            <Button
              type="text"
              icon={<BookOutlined />}
              onClick={() => window.open('/swagger/index.html', '_blank')}
              title="API 文档（在新标签页打开）"
              style={{
                color: mode === 'dark' ? '#F1F5F9' : '#1E293B',
              }}
            >
              Docs
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
              {mode === 'dark' ? <SunFilled /> : <MoonFilled />}
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
        <Footer
          style={{
            textAlign: 'center',
            padding: '12px 24px',
            background: 'transparent',
            color: mode === 'dark' ? '#64748B' : '#94A3B8',
            fontSize: 12,
          }}
        >
          系统运行部驱动 © {new Date().getFullYear()} KkOps
        </Footer>
      </Layout>

      {/* 修改密码弹窗 */}
      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false)
          passwordForm.resetFields()
        }}
        footer={null}
        destroyOnHidden
      >
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleChangePassword}
        >
          <Form.Item
            name="old_password"
            label="原密码"
            rules={[{ required: true, message: '请输入原密码' }]}
          >
            <Input.Password placeholder="请输入原密码" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码长度至少6位' },
            ]}
          >
            <Input.Password placeholder="请输入新密码（至少6位）" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认新密码"
            dependencies={['new_password']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Button
              style={{ marginRight: 8 }}
              onClick={() => {
                setPasswordModalVisible(false)
                passwordForm.resetFields()
              }}
            >
              取消
            </Button>
            <Button type="primary" htmlType="submit" loading={passwordLoading}>
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <CommandPalette
        open={paletteOpen}
        onClose={() => setPaletteOpen(false)}
        commands={paletteCommands}
        onNavigate={(path) => navigate(path)}
      />
    </Layout>
  )
}

export default MainLayout
