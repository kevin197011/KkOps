import React, { useState } from 'react';
import { Layout, Menu, Avatar, Dropdown, Space, Button } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  ProjectOutlined,
  LaptopOutlined,
  RocketOutlined,
  UserOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  TeamOutlined,
  FileTextOutlined,
  SafetyOutlined,
  LockOutlined,
  ConsoleSqlOutlined,
  KeyOutlined,
  SettingOutlined,
  CloudOutlined,
  CloudServerOutlined,
  TagsOutlined,
  AppstoreOutlined,
  ThunderboltOutlined,
  BranchesOutlined,
  GlobalOutlined,
  DeploymentUnitOutlined,
} from '@ant-design/icons';
import { useAuth } from '../contexts/AuthContext';
import Footer from './Footer';
import type { MenuProps } from 'antd';

const { Header, Sider, Content } = Layout;

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();

  // 检查用户是否为管理员
  const isAdmin = user?.roles?.some(role => role.name === 'admin') || false;
  
  // 检查用户是否有批量操作权限（admin 和 operator 都有权限）
  const hasBatchOperationPermission = user?.roles?.some(role => 
    role.name === 'admin' || role.name === 'operator'
  ) || false;

  // 根据当前路径获取展开的菜单
  const getOpenKeys = () => {
    const path = location.pathname;
    if (['/projects', '/environments', '/cloud-platforms'].includes(path)) {
      return ['base-management'];
    }
    if (['/hosts', '/host-groups', '/host-tags', '/ssh-keys'].includes(path)) {
      return ['asset-management'];
    }
    if (['/batch-operations'].includes(path)) {
      return [];
    }
    if (['/users', '/roles', '/permissions'].includes(path)) {
      return ['user-management'];
    }
    return [];
  };

  const [openKeys, setOpenKeys] = useState<string[]>(getOpenKeys());

  // 菜单项配置
  const menuItems: MenuProps['items'] = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: '仪表盘',
    },
    {
      key: 'base-management',
      icon: <SettingOutlined />,
      label: '基础管理',
      children: [
        {
          key: '/projects',
          icon: <ProjectOutlined />,
          label: '项目管理',
        },
        {
          key: '/environments',
          icon: <CloudOutlined />,
          label: '环境管理',
        },
        {
          key: '/cloud-platforms',
          icon: <CloudServerOutlined />,
          label: '云平台管理',
        },
      ],
    },
    {
      key: 'asset-management',
      icon: <LaptopOutlined />,
      label: '资产管理',
      children: [
        {
          key: '/hosts',
          icon: <LaptopOutlined />,
          label: '主机管理',
        },
        {
          key: '/host-groups',
          icon: <AppstoreOutlined />,
          label: '主机分组',
        },
        {
          key: '/host-tags',
          icon: <TagsOutlined />,
          label: '主机标签',
        },
        {
          key: '/ssh-keys',
          icon: <KeyOutlined />,
          label: 'SSH密钥',
        },
      ],
    },
    {
      key: 'cicd-management',
      icon: <BranchesOutlined />,
      label: 'CICD管理',
      children: [
        {
          key: '/deployments',
          icon: <RocketOutlined />,
          label: '部署管理',
        },
        // 批量操作（仅 admin 和 operator 可见）
        ...(hasBatchOperationPermission ? [{
          key: '/batch-operations',
          icon: <ThunderboltOutlined />,
          label: '批量操作',
        }] : []),
        {
          key: '/formulas',
          icon: <DeploymentUnitOutlined />,
          label: 'Formula部署',
        },
      ],
    },
    {
      key: 'user-management',
      icon: <TeamOutlined />,
      label: '用户权限',
      children: [
        {
          key: '/users',
          icon: <TeamOutlined />,
          label: '用户管理',
        },
        {
          key: '/roles',
          icon: <SafetyOutlined />,
          label: '角色管理',
        },
        {
          key: '/permissions',
          icon: <LockOutlined />,
          label: '权限管理',
        },
      ],
    },
    {
      key: 'system-management',
      icon: <GlobalOutlined />,
      label: '系统管理',
      children: [
        {
          key: '/audit',
          icon: <FileTextOutlined />,
          label: '审计管理',
        },
        // 系统设置（仅管理员可见）
        ...(isAdmin ? [{
          key: '/settings',
          icon: <SettingOutlined />,
          label: '系统设置',
        }] : []),
      ],
    },
  ];

  // 处理菜单点击
  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  // 用户下拉菜单
  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
    },
  ];

  const handleUserMenuClick = ({ key }: { key: string }) => {
    if (key === 'logout') {
      logout();
      navigate('/login');
    } else if (key === 'profile') {
      // TODO: 跳转到个人资料页面
      console.log('跳转到个人资料');
    }
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={200}
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
        }}
      >
        <div
          style={{
            height: 64,
            margin: 16,
            background: 'rgba(255, 255, 255, 0.2)',
            borderRadius: 6,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: collapsed ? 16 : 18,
            fontWeight: 'bold',
            cursor: 'pointer',
          }}
          onClick={() => navigate('/')}
        >
          {collapsed ? (
            <img 
              src="/logo-icon.svg" 
              alt="KkOps" 
              style={{ width: 24, height: 24 }}
            />
          ) : (
            <img 
              src="/logo.svg" 
              alt="KkOps" 
              style={{ height: 28, width: 'auto' }}
            />
          )}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          openKeys={collapsed ? [] : openKeys}
          onOpenChange={(keys) => setOpenKeys(keys as string[])}
          items={menuItems}
          onClick={handleMenuClick}
        />
      </Sider>
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'all 0.2s' }}>
        <Header
          style={{
            padding: '0 24px',
            background: '#fff',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            position: 'relative',
            zIndex: 100,
          }}
        >
          <div
            style={{
              fontSize: 18,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
            }}
            onClick={() => setCollapsed(!collapsed)}
          >
            {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          </div>
          <Space>
            <Button
              type="text"
              icon={<ConsoleSqlOutlined />}
              onClick={() => {
                window.open('/webssh', '_blank');
              }}
              title="WebSSH终端"
              style={{ 
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                zIndex: 10,
                position: 'relative'
              }}
            >
              WebSSH
            </Button>
            <span>{user?.display_name || user?.username}</span>
            <Dropdown
              menu={{
                items: userMenuItems,
                onClick: handleUserMenuClick,
              }}
              placement="bottomRight"
            >
              <Avatar
                style={{ backgroundColor: '#1890ff', cursor: 'pointer' }}
                icon={<UserOutlined />}
              />
            </Dropdown>
          </Space>
        </Header>
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            background: '#fff',
            minHeight: 280,
          }}
        >
          {children}
        </Content>
        <Footer />
      </Layout>
    </Layout>
  );
};

export default MainLayout;
