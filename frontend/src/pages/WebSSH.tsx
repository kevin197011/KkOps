import React, { useState, useEffect, useMemo, useCallback } from 'react';
import {
  Layout,
  Tree,
  Input,
  Tabs,
  message,
  Badge,
  Empty,
  Spin,
  Dropdown,
  Space,
  Avatar,
} from 'antd';
import type { DataNode } from 'antd/es/tree';
import type { DragEndEvent } from '@dnd-kit/core';
import {
  DndContext,
  PointerSensor,
  useSensor,
  useSensors,
  closestCenter,
} from '@dnd-kit/core';
import {
  SortableContext,
  useSortable,
  horizontalListSortingStrategy,
  arrayMove,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import {
  FolderOutlined,
  DesktopOutlined,
  ReloadOutlined,
  CloseOutlined,
  CopyOutlined,
  UserOutlined,
  HolderOutlined,
} from '@ant-design/icons';
import { hostService, Host } from '../services/host';
import { projectService, Project } from '../services/project';
import Terminal from '../components/Terminal';
import Footer from '../components/Footer';
import { useAuth } from '../contexts/AuthContext';

const { Header, Sider, Content } = Layout;
const { Search } = Input;

// 终端标签页类型
interface TerminalTab {
  id: string;
  hostId: number;
  hostName: string;
  connectionStatus: 'connecting' | 'connected' | 'disconnected';
  // 认证信息（用于克隆时复用）
  authInfo?: {
    authMethod: 'key' | 'password';
    username: string;
    password?: string;
    keyId?: number;
  };
}

// 树节点数据类型
interface TreeNodeData extends DataNode {
  type?: 'project' | 'environment' | 'host';
  projectId?: number;
  environment?: string;
  hostId?: number;
  hostData?: Host;
  hostCount?: number;
  children?: TreeNodeData[];
}

// 自定义样式，确保 Tabs 内容区域正确填充
const tabsStyles = `
  .webssh-tabs {
    height: 100%;
  }
  .webssh-tabs .ant-tabs-content-holder {
    flex: 1;
    overflow: hidden;
    background: #1e1e1e;
    position: relative;
  }
  .webssh-tabs .ant-tabs-content {
    height: 100%;
    background: #1e1e1e;
    position: relative;
  }
  .webssh-tabs .ant-tabs-tabpane {
    height: 100%;
    background: #1e1e1e;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: none;
  }
  .webssh-tabs .ant-tabs-tabpane-active {
    display: flex !important;
    flex-direction: column;
    position: relative;
  }
  .webssh-tabs .ant-tabs-nav {
    margin-bottom: 0 !important;
  }
  .webssh-tabs .ant-tabs-tab {
    background: #2d2d2d !important;
    border-color: #3a3f4b !important;
    color: #8b949e !important;
    transition: all 0.2s ease;
  }
  .webssh-tabs .ant-tabs-tab:hover {
    color: #e0e0e0 !important;
    background: #3a3f4b !important;
  }
  .webssh-tabs .ant-tabs-tab-active {
    background: #1890ff !important;
    border-color: #1890ff !important;
    color: #ffffff !important;
  }
  .webssh-tabs .ant-tabs-tab-active .ant-tabs-tab-btn {
    color: #ffffff !important;
  }
  .webssh-tabs .ant-tabs-tab-active:hover {
    background: #40a9ff !important;
  }
  .webssh-tabs .ant-tabs-tab-remove {
    color: #8b949e !important;
  }
  .webssh-tabs .ant-tabs-tab-active .ant-tabs-tab-remove {
    color: #ffffff !important;
  }
  .webssh-tabs .ant-tabs-tab-active .ant-tabs-tab-remove:hover {
    color: #ffccc7 !important;
  }
  /* DevOps 深色主题样式 */
  .webssh-dark-tree .ant-tree {
    background: transparent;
    color: #e0e0e0;
  }
  .webssh-dark-tree .ant-tree-node-content-wrapper:hover {
    background: #3a3f4b;
  }
  .webssh-dark-tree .ant-tree-node-selected {
    background: #1890ff !important;
  }
  .webssh-dark-tree .ant-tree-switcher {
    color: #8c8c8c;
  }
  .webssh-dark-tree .ant-tree-indent-unit::before {
    border-color: #3a3f4b;
  }
  .webssh-dark-search .ant-input-affix-wrapper {
    background: #2d333b;
    border-color: #3a3f4b;
  }
  .webssh-dark-search .ant-input {
    background: #2d333b;
    color: #e0e0e0;
  }
  .webssh-dark-search .ant-input::placeholder {
    color: #6e7681;
  }
  .webssh-dark-search .ant-input-clear-icon {
    color: #6e7681;
  }
  .webssh-dark-search .ant-input-search-button {
    background: #3a3f4b;
    border-color: #3a3f4b;
  }
  .draggable-tab {
    cursor: grab;
  }
  .draggable-tab:active {
    cursor: grabbing;
  }
  .draggable-tab.dragging {
    opacity: 0.5;
  }
`;

// 可拖拽的标签页组件
interface DraggableTabProps {
  id: string;
  children: React.ReactNode;
  isActive: boolean;
}

const DraggableTab: React.FC<DraggableTabProps> = ({ id, children, isActive }) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    display: 'inline-flex',
    alignItems: 'center',
    padding: '8px 12px',
    marginRight: '2px',
    borderRadius: '4px 4px 0 0',
    cursor: 'grab',
    background: isActive ? '#1890ff' : '#2d2d2d',
    color: isActive ? '#ffffff' : '#8b949e',
    border: `1px solid ${isActive ? '#1890ff' : '#3a3f4b'}`,
    borderBottom: 'none',
    userSelect: 'none',
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      {children}
    </div>
  );
};

// 可拖拽的标签节点包装器（用于 Ant Design Tabs）
interface DraggableTabNodeProps {
  id: string;
  children: React.ReactNode;
  isActive: boolean;
}

const DraggableTabNode: React.FC<DraggableTabNodeProps> = ({ id, children, isActive }) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.6 : 1,
    cursor: isDragging ? 'grabbing' : 'grab',
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      {children}
    </div>
  );
};

const WebSSH: React.FC = () => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchValue, setSearchValue] = useState('');
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<React.Key[]>([]);
  const [terminalTabs, setTerminalTabs] = useState<TerminalTab[]>([]);
  const [activeTabKey, setActiveTabKey] = useState<string | undefined>(undefined);

  // 拖拽传感器配置
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8, // 拖动8px后才开始拖拽，避免误触
      },
    })
  );

  // 处理拖拽结束
  const handleDragEnd = useCallback((event: DragEndEvent) => {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      setTerminalTabs((tabs) => {
        const oldIndex = tabs.findIndex((tab) => tab.id === active.id);
        const newIndex = tabs.findIndex((tab) => tab.id === over.id);
        return arrayMove(tabs, oldIndex, newIndex);
      });
    }
  }, []);

  // 加载项目列表
  useEffect(() => {
    loadProjects();
  }, []);

  // 加载主机列表
  useEffect(() => {
    loadHosts();
  }, []);

  const loadProjects = async () => {
    try {
      const response = await projectService.list(1, 100);
      setProjects(response.projects);
    } catch (error: any) {
      message.error('加载项目列表失败');
    }
  };

  const loadHosts = async () => {
    setLoading(true);
    try {
      // 加载所有主机（使用大页面大小）
      const response = await hostService.list(1, 10000);
      setHosts(response.hosts);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载主机列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 将主机列表转换为树形结构
  const treeData = useMemo(() => {
    if (!hosts.length) return [];

    // 按项目分组
    const projectMap = new Map<number, Map<string, Host[]>>();

    hosts.forEach((host) => {
      const projectId = host.project_id || 0;
      const environment = host.environment || 'default';

      if (!projectMap.has(projectId)) {
        projectMap.set(projectId, new Map());
      }

      const envMap = projectMap.get(projectId)!;
      if (!envMap.has(environment)) {
        envMap.set(environment, []);
      }

      envMap.get(environment)!.push(host);
    });

    const tree: TreeNodeData[] = [];

    // 遍历项目
    projectMap.forEach((envMap, projectId) => {
      const project = projects.find((p) => p.id === projectId);
      const projectName = project?.name || `项目 #${projectId}` || '未分组';

      // 计算项目下的主机总数
      let projectHostCount = 0;
      envMap.forEach((hosts) => {
        projectHostCount += hosts.length;
      });

      const projectNode: TreeNodeData = {
        title: (
          <span>
            <FolderOutlined style={{ marginRight: 8 }} />
            {projectName}
            <Badge count={projectHostCount} style={{ marginLeft: 8 }} />
          </span>
        ),
        key: `project-${projectId}`,
        type: 'project',
        projectId,
        children: [],
      };

      // 遍历环境
      const environmentNodes: TreeNodeData[] = [];
      envMap.forEach((hosts, environment) => {
        const envNode: TreeNodeData = {
          title: (
            <span>
              <FolderOutlined style={{ marginRight: 8 }} />
              {environment.toUpperCase()}
              <Badge count={hosts.length} style={{ marginLeft: 8 }} />
            </span>
          ),
          key: `project-${projectId}-env-${environment}`,
          type: 'environment',
          projectId,
          environment,
          hostCount: hosts.length,
          children: [],
        };

        // 添加主机节点
        const hostNodes: TreeNodeData[] = hosts.map((host) => {
          const statusColor =
            host.status === 'online'
              ? '#52c41a'
              : host.status === 'offline'
              ? '#ff4d4f'
              : '#d9d9d9';

          return {
            title: (
              <span>
                <DesktopOutlined style={{ marginRight: 8, color: statusColor }} />
                {host.hostname} ({host.ip_address})
              </span>
            ),
            key: `host-${host.id}`,
            type: 'host',
            projectId,
            environment,
            hostId: host.id,
            hostData: host,
            isLeaf: true,
          };
        });

        envNode.children = hostNodes;
        environmentNodes.push(envNode);
      });

      projectNode.children = environmentNodes;
      tree.push(projectNode);
    });

    // 如果没有项目的主机，添加到"未分组"
    const ungroupedHosts = hosts.filter((h) => !h.project_id);
    if (ungroupedHosts.length > 0) {
      const ungroupedNode: TreeNodeData = {
        title: (
          <span>
            <FolderOutlined style={{ marginRight: 8 }} />
            未分组
            <Badge count={ungroupedHosts.length} style={{ marginLeft: 8 }} />
          </span>
        ),
        key: 'project-0',
        type: 'project',
        projectId: 0,
        children: [
          {
            title: (
              <span>
                <FolderOutlined style={{ marginRight: 8 }} />
                默认环境
                <Badge count={ungroupedHosts.length} style={{ marginLeft: 8 }} />
              </span>
            ),
            key: 'project-0-env-default',
            type: 'environment',
            projectId: 0,
            environment: 'default',
            hostCount: ungroupedHosts.length,
            children: ungroupedHosts.map((host) => {
              const statusColor =
                host.status === 'online'
                  ? '#52c41a'
                  : host.status === 'offline'
                  ? '#ff4d4f'
                  : '#d9d9d9';

              return {
                title: (
                  <span>
                    <DesktopOutlined style={{ marginRight: 8, color: statusColor }} />
                    {host.hostname} ({host.ip_address})
                  </span>
                ),
                key: `host-${host.id}`,
                type: 'host',
                projectId: 0,
                environment: 'default',
                hostId: host.id,
                hostData: host,
                isLeaf: true,
              };
            }),
          },
        ],
      };
      tree.push(ungroupedNode);
    }

    return tree;
  }, [hosts, projects]);

  // 过滤树数据（搜索功能）
  const filteredTreeData = useMemo(() => {
    if (!searchValue) return treeData;

    const filterTree = (nodes: TreeNodeData[]): TreeNodeData[] => {
      return nodes
        .map((node) => {
          const matchesSearch =
            node.title?.toString().toLowerCase().includes(searchValue.toLowerCase()) ||
            (node.type === 'host' &&
              node.hostData &&
              (node.hostData.hostname.toLowerCase().includes(searchValue.toLowerCase()) ||
                node.hostData.ip_address.toLowerCase().includes(searchValue.toLowerCase())));

          const filteredChildren = node.children
            ? filterTree(node.children as TreeNodeData[])
            : [];

          if (matchesSearch || filteredChildren.length > 0) {
            return {
              ...node,
              children: filteredChildren.length > 0 ? filteredChildren : node.children,
            } as TreeNodeData;
          }
          return null;
        })
        .filter((node): node is TreeNodeData => node !== null);
    };

    return filterTree(treeData);
  }, [treeData, searchValue]);

  // 处理树节点选择（仅用于展开/折叠，不创建标签页）
  const handleTreeSelect = (selectedKeys: React.Key[], info: any) => {
    // 选择事件不再创建标签页，由 onClick 处理
  };

  // 处理树节点点击（允许重复点击同一主机）
  const handleTreeClick = (e: React.MouseEvent, node: any) => {
    const nodeKey = node.key;
    if (!nodeKey || typeof nodeKey !== 'string') return;

    // 只处理主机节点的点击
    if (nodeKey.startsWith('host-')) {
      const hostId = parseInt(nodeKey.replace('host-', ''));
      const host = hosts.find((h) => h.id === hostId);
      if (!host) return;

      // 每次点击都创建新标签页，允许同一主机打开多个终端窗口
      const tabId = `terminal-${hostId}-${Date.now()}`;
      const newTab: TerminalTab = {
        id: tabId,
        hostId: host.id,
        hostName: host.hostname || `${host.ip_address}:${host.ssh_port || 22}`,
        connectionStatus: 'connecting',
      };
      setTerminalTabs((prev) => [...prev, newTab]);
      setActiveTabKey(tabId);
    }
  };

  // 处理树节点展开
  const handleTreeExpand = (expandedKeys: React.Key[]) => {
    setExpandedKeys(expandedKeys);
  };

  // 关闭终端标签页
  const handleCloseTab = (targetKey: string) => {
    const newTabs = terminalTabs.filter((tab) => tab.id !== targetKey);
    setTerminalTabs(newTabs);
    if (activeTabKey === targetKey) {
      setActiveTabKey(newTabs.length > 0 ? newTabs[newTabs.length - 1].id : undefined);
    }
  };

  // 克隆终端标签页
  const handleCloneTab = (tabId: string) => {
    const tab = terminalTabs.find((t) => t.id === tabId);
    if (!tab) return;

    const newTabId = `terminal-${tab.hostId}-${Date.now()}`;
    const newTab: TerminalTab = {
      id: newTabId,
      hostId: tab.hostId,
      hostName: tab.hostName,
      connectionStatus: 'connecting',
      // 复制认证信息，克隆时不需要再次输入
      authInfo: tab.authInfo,
    };
    setTerminalTabs((prev) => [...prev, newTab]);
    setActiveTabKey(newTabId);
  };

  // 更新标签页认证信息
  const handleAuthInfoUpdate = (tabId: string, authInfo: TerminalTab['authInfo']) => {
    setTerminalTabs((prev) =>
      prev.map((tab) => (tab.id === tabId ? { ...tab, authInfo } : tab))
    );
  };

  // 关闭所有标签页
  const handleCloseAll = () => {
    setTerminalTabs([]);
    setActiveTabKey(undefined);
  };

  // 关闭其他标签页
  const handleCloseOthers = (tabId: string) => {
    const tab = terminalTabs.find((t) => t.id === tabId);
    if (!tab) return;

    setTerminalTabs([tab]);
    setActiveTabKey(tabId);
  };

  // 右键菜单项
  const contextMenuItems = (tabId: string) => [
    {
      key: 'clone',
      icon: <CopyOutlined />,
      label: '克隆',
      onClick: () => handleCloneTab(tabId),
    },
    {
      key: 'close',
      icon: <CloseOutlined />,
      label: '关闭',
      onClick: () => handleCloseTab(tabId),
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'closeOthers',
      icon: <CloseOutlined />,
      label: '关闭其他',
      disabled: terminalTabs.length <= 1,
      onClick: () => handleCloseOthers(tabId),
    },
    {
      key: 'closeAll',
      icon: <CloseOutlined />,
      label: '关闭所有',
      disabled: terminalTabs.length === 0,
      onClick: handleCloseAll,
    },
  ];

  // 更新标签页连接状态
  const handleTerminalStatusChange = (tabId: string, status: TerminalTab['connectionStatus']) => {
    setTerminalTabs((prev) =>
      prev.map((tab) => (tab.id === tabId ? { ...tab, connectionStatus: status } : tab))
    );
  };

  const { user } = useAuth();

  return (
    <Layout style={{ height: '100vh', margin: 0, padding: 0, background: '#1a1d21', display: 'flex', flexDirection: 'column' }}>
      <style>{tabsStyles}</style>
      <Header
        style={{
          background: '#21262d',
          borderBottom: '1px solid #30363d',
          padding: '0 24px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
          flexShrink: 0,
        }}
      >
        <h2 style={{ margin: 0, color: '#e0e0e0', fontWeight: 500 }}>WebSSH 终端</h2>
        <Space>
          <span style={{ color: '#8b949e' }}>{user?.display_name || user?.username}</span>
          <Avatar
            style={{ backgroundColor: '#238636' }}
            icon={<UserOutlined />}
          />
        </Space>
      </Header>
      <Layout style={{ flex: 1, background: '#1a1d21', overflow: 'hidden' }}>
        <Sider
          width={300}
          style={{
            background: '#21262d',
            borderRight: '1px solid #30363d',
            overflow: 'auto',
            boxShadow: '2px 0 8px rgba(0,0,0,0.15)',
          }}
        >
        <div style={{ padding: '16px', borderBottom: '1px solid #30363d', background: '#161b22' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
            <h3 style={{ margin: 0, color: '#e0e0e0', fontWeight: 500, fontSize: '16px' }}>主机列表</h3>
            <ReloadOutlined
              onClick={loadHosts}
              style={{ cursor: 'pointer', fontSize: '16px', color: '#8b949e' }}
              title="刷新"
            />
          </div>
          <div className="webssh-dark-search">
            <Search
              placeholder="搜索主机..."
              allowClear
              value={searchValue}
              onChange={(e) => setSearchValue(e.target.value)}
              style={{ marginBottom: '8px' }}
            />
          </div>
        </div>
        <div className="webssh-dark-tree" style={{ padding: '8px', background: '#21262d' }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '40px 0' }}>
              <Spin />
            </div>
          ) : filteredTreeData.length === 0 ? (
            <Empty description={<span style={{ color: '#8b949e' }}>暂无主机</span>} style={{ marginTop: '40px' }} />
          ) : (
            <Tree
              treeData={filteredTreeData}
              expandedKeys={expandedKeys}
              selectedKeys={selectedKeys}
              onSelect={handleTreeSelect}
              onExpand={handleTreeExpand}
              onClick={handleTreeClick}
              showLine={{ showLeafIcon: false }}
              defaultExpandAll={false}
            />
          )}
        </div>
        </Sider>
        <Content style={{ background: '#1a1d21', overflow: 'hidden', padding: '8px', display: 'flex', flexDirection: 'column' }}>
        {terminalTabs.length === 0 ? (
          <div style={{ 
            background: '#21262d', 
            borderRadius: '6px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
            border: '1px solid #30363d',
            flex: 1,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <Empty
              description={<span style={{ color: '#8b949e' }}>从左侧树中选择主机以打开终端</span>}
              style={{ marginTop: 0 }}
            />
          </div>
        ) : (
          <div style={{ 
            background: '#1e1e1e', 
            borderRadius: '6px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden'
          }}>
            <DndContext
              sensors={sensors}
              collisionDetection={closestCenter}
              onDragEnd={handleDragEnd}
            >
              <Tabs
                className="webssh-tabs"
                activeKey={activeTabKey}
                onChange={(key) => {
                  setActiveTabKey(key);
                  // 切换标签后重新 fit 终端
                  setTimeout(() => {
                    window.dispatchEvent(new Event('resize'));
                  }, 100);
                }}
                type="editable-card"
                hideAdd
                style={{ 
                  flex: 1,
                  display: 'flex',
                  flexDirection: 'column',
                }}
                tabBarStyle={{
                  marginBottom: 0,
                  padding: '0 16px',
                  background: '#2d2d2d',
                  flexShrink: 0,
                }}
                onEdit={(targetKey, action) => {
                  if (action === 'remove' && typeof targetKey === 'string') {
                    handleCloseTab(targetKey);
                  }
                }}
                renderTabBar={(tabBarProps, DefaultTabBar) => (
                  <SortableContext
                    items={terminalTabs.map((tab) => tab.id)}
                    strategy={horizontalListSortingStrategy}
                  >
                    <DefaultTabBar {...tabBarProps}>
                      {(node) => (
                        <DraggableTabNode
                          key={node.key}
                          id={String(node.key)}
                          isActive={activeTabKey === node.key}
                        >
                          {node}
                        </DraggableTabNode>
                      )}
                    </DefaultTabBar>
                  </SortableContext>
                )}
                items={terminalTabs.map((tab) => {
                  const tabTitle = (
                    <Dropdown
                      menu={{ items: contextMenuItems(tab.id) }}
                      trigger={['contextMenu']}
                    >
                      <span
                        style={{ userSelect: 'none' }}
                      >
                        {tab.hostName}
                        {tab.connectionStatus === 'connecting' && (
                          <Spin size="small" style={{ marginLeft: 8 }} />
                        )}
                        {tab.connectionStatus === 'connected' && (
                          <span style={{ color: '#52c41a', marginLeft: 8 }}>●</span>
                        )}
                        {tab.connectionStatus === 'disconnected' && (
                          <span style={{ color: '#ff4d4f', marginLeft: 8 }}>●</span>
                        )}
                      </span>
                    </Dropdown>
                  );

                  return {
                    label: tabTitle,
                    key: tab.id,
                    children: (
                      <Terminal
                        hostId={tab.hostId}
                        hostName={tab.hostName}
                        authInfo={tab.authInfo}
                        onClose={() => handleCloseTab(tab.id)}
                        onStatusChange={(status) => handleTerminalStatusChange(tab.id, status)}
                        onAuthInfoUpdate={(authInfo) => handleAuthInfoUpdate(tab.id, authInfo)}
                      />
                    ),
                    closable: true,
                  };
                })}
              />
            </DndContext>
          </div>
        )}
        </Content>
      </Layout>
      <Footer dark />
    </Layout>
  );
};

export default WebSSH;
