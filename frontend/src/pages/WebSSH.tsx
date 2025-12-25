import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  message,
  Tag,
  Card,
  Tabs,
  Select,
} from 'antd';
import {
  ConsoleSqlOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { hostService, Host } from '../services/host';
import { projectService, Project } from '../services/project';
import { hostGroupService, HostGroup } from '../services/hostGroup';
import { hostTagService, HostTag } from '../services/hostTag';
import Terminal from '../components/Terminal';

const { Option } = Select;

const { TabPane } = Tabs;

const WebSSH: React.FC = () => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [hostGroups, setHostGroups] = useState<HostGroup[]>([]);
  const [hostTags, setHostTags] = useState<HostTag[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined);
  const [selectedGroupId, setSelectedGroupId] = useState<number | undefined>(undefined);
  const [selectedTagId, setSelectedTagId] = useState<number | undefined>(undefined);
  
  // 终端标签页管理
  const [activeTab, setActiveTab] = useState('hosts');
  const [terminalTabs, setTerminalTabs] = useState<Array<{ id: string; hostId: number; hostName: string }>>([]);

  useEffect(() => {
    loadProjects();
    loadHostGroups();
    loadHostTags();
  }, []);

  useEffect(() => {
    loadHosts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, selectedProjectId, selectedGroupId, selectedTagId]);

  const loadProjects = async () => {
    try {
      const response = await projectService.list(1, 100);
      setProjects(response.projects);
    } catch (error: any) {
      message.error('加载项目列表失败');
    }
  };

  const loadHostGroups = async () => {
    try {
      const response = await hostGroupService.list(1, 1000);
      setHostGroups(response.groups);
    } catch (error: any) {
      // 静默失败，不影响主流程
    }
  };

  const loadHostTags = async () => {
    try {
      const response = await hostTagService.list(1, 1000);
      setHostTags(response.tags);
    } catch (error: any) {
      // 静默失败，不影响主流程
    }
  };

  const loadHosts = async () => {
    setLoading(true);
    try {
      const filters: any = {};
      if (selectedProjectId) {
        filters.project_id = selectedProjectId;
      }
      if (selectedGroupId) {
        filters.group_id = selectedGroupId;
      }
      if (selectedTagId) {
        filters.tag_id = selectedTagId;
      }
      const response = await hostService.list(page, pageSize, filters);
      setHosts(response.hosts);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载主机列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleOpenTerminal = (host: Host) => {
    const tabId = `terminal-${host.id}-${Date.now()}`;
    const newTab = {
      id: tabId,
      hostId: host.id,
      hostName: host.hostname || `${host.ip_address}:${host.ssh_port || 22}`,
    };
    setTerminalTabs((prevTabs) => [...prevTabs, newTab]);
    setActiveTab(tabId);
  };

  const handleCloseTerminalTab = (targetTabId: string) => {
    const newTabs = terminalTabs.filter((tab) => tab.id !== targetTabId);
    setTerminalTabs(newTabs);
    if (activeTab === targetTabId) {
      setActiveTab(newTabs.length > 0 ? newTabs[newTabs.length - 1].id : 'hosts');
    }
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      online: { color: 'success', text: '在线' },
      offline: { color: 'error', text: '离线' },
      unknown: { color: 'default', text: '未知' },
    };
    const statusInfo = statusMap[status] || statusMap.unknown;
    return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '项目',
      key: 'project',
      width: 150,
      render: (_: any, record: Host) => {
        if (record.project) {
          return <span>{record.project.name}</span>;
        }
        if (record.project_id) {
          const project = projects.find((p) => p.id === record.project_id);
          if (project) {
            return <span>{project.name}</span>;
          }
        }
        return <span style={{ color: '#999' }}>-</span>;
      },
    },
    {
      title: '主机名',
      dataIndex: 'hostname',
      key: 'hostname',
    },
    {
      title: 'IP地址',
      dataIndex: 'ip_address',
      key: 'ip_address',
    },
    {
      title: 'SSH端口',
      dataIndex: 'ssh_port',
      key: 'ssh_port',
      width: 100,
      render: (port: number) => port || 22,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '环境',
      dataIndex: 'environment',
      key: 'environment',
      width: 100,
      render: (env: string) => {
        if (!env) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        const envColors: Record<string, string> = {
          dev: 'blue',
          uat: 'orange',
          staging: 'purple',
          prod: 'red',
          test: 'green',
        };
        return <Tag color={envColors[env] || 'default'}>{env.toUpperCase()}</Tag>;
      },
    },
    {
      title: '组',
      key: 'groups',
      width: 150,
      render: (_: any, record: Host) => (
        <Space wrap>
          {record.groups && record.groups.length > 0 ? (
            record.groups.map((group) => (
              <Tag key={group.id}>{group.name}</Tag>
            ))
          ) : (
            <span style={{ color: '#999' }}>-</span>
          )}
        </Space>
      ),
    },
    {
      title: '标签',
      key: 'tags',
      width: 150,
      render: (_: any, record: Host) => (
        <Space wrap>
          {record.tags && record.tags.length > 0 ? (
            record.tags.map((tag) => (
              <Tag key={tag.id} color={tag.color || '#1890ff'}>
                {tag.name}
              </Tag>
            ))
          ) : (
            <span style={{ color: '#999' }}>-</span>
          )}
        </Space>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Host) => (
        <Button
          type="primary"
          size="small"
          icon={<ConsoleSqlOutlined />}
          onClick={() => handleOpenTerminal(record)}
        >
          打开终端
        </Button>
      ),
    },
  ];

  return (
    <Card>
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        type="editable-card"
        hideAdd
        onEdit={(targetKey, action) => {
          if (action === 'remove' && typeof targetKey === 'string') {
            handleCloseTerminalTab(targetKey);
          }
        }}
      >
        <TabPane tab="主机列表" key="hosts" closable={false}>
          <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h2 style={{ margin: 0 }}>WebSSH管理</h2>
            <Space wrap>
              <Select
                placeholder="选择项目"
                allowClear
                style={{ width: 150 }}
                value={selectedProjectId}
                onChange={(value) => {
                  setSelectedProjectId(value);
                  setPage(1);
                }}
              >
                {projects.map((project) => (
                  <Option key={project.id} value={project.id}>
                    {project.name}
                  </Option>
                ))}
              </Select>
              <Select
                placeholder="选择主机组"
                allowClear
                style={{ width: 150 }}
                value={selectedGroupId}
                onChange={(value) => {
                  setSelectedGroupId(value);
                  setPage(1);
                }}
              >
                {hostGroups.map((group) => (
                  <Option key={group.id} value={group.id}>
                    {group.name}
                  </Option>
                ))}
              </Select>
              <Select
                placeholder="选择标签"
                allowClear
                style={{ width: 150 }}
                value={selectedTagId}
                onChange={(value) => {
                  setSelectedTagId(value);
                  setPage(1);
                }}
              >
                {hostTags.map((tag) => (
                  <Option key={tag.id} value={tag.id}>
                    {tag.name}
                  </Option>
                ))}
              </Select>
              <Button icon={<ReloadOutlined />} onClick={loadHosts}>
                刷新
              </Button>
            </Space>
          </div>

          <Table
            columns={columns}
            dataSource={hosts}
            rowKey="id"
            loading={loading}
            pagination={{
              current: page,
              pageSize: pageSize,
              total: total,
              showSizeChanger: true,
              showTotal: (total) => `共 ${total} 条`,
              onChange: (page, pageSize) => {
                setPage(page);
                setPageSize(pageSize);
              },
            }}
          />
        </TabPane>
        {terminalTabs.map((tab) => (
          <TabPane
            tab={tab.hostName}
            key={tab.id}
            closable
          >
            <div style={{ height: 'calc(100vh - 300px)', minHeight: '500px' }}>
              <Terminal
                hostId={tab.hostId}
                hostName={tab.hostName}
                onClose={() => handleCloseTerminalTab(tab.id)}
              />
            </div>
          </TabPane>
        ))}
      </Tabs>
    </Card>
  );
};

export default WebSSH;

