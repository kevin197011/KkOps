import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  message,
  Tag,
  Card,
  Tabs,
  Modal,
  Form,
  Input,
  Popconfirm,
} from 'antd';
import {
  ConsoleSqlOutlined,
  ReloadOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  KeyOutlined,
} from '@ant-design/icons';
import { hostService, Host } from '../services/host';
import { projectService, Project } from '../services/project';
import { sshKeyService, SSHKey, CreateSSHKeyRequest } from '../services/sshKey';
import Terminal from '../components/Terminal';

const { TabPane } = Tabs;

const WebSSH: React.FC = () => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  
  // 终端标签页管理
  const [activeTab, setActiveTab] = useState('hosts');
  const [terminalTabs, setTerminalTabs] = useState<Array<{ id: string; hostId: number; hostName: string }>>([]);

  // SSH密钥管理
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([]);
  const [sshKeysLoading, setSshKeysLoading] = useState(false);
  const [sshKeysTotal, setSshKeysTotal] = useState(0);
  const [sshKeysPage, setSshKeysPage] = useState(1);
  const [sshKeysPageSize, setSshKeysPageSize] = useState(10);
  const [sshKeyModalVisible, setSshKeyModalVisible] = useState(false);
  const [editingSshKey, setEditingSshKey] = useState<SSHKey | null>(null);
  const [sshKeyForm] = Form.useForm();

  useEffect(() => {
    loadProjects();
    if (activeTab === 'ssh-keys') {
      loadSSHKeys();
    }
  }, [activeTab]);

  useEffect(() => {
    loadHosts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  useEffect(() => {
    if (activeTab === 'ssh-keys') {
      loadSSHKeys();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sshKeysPage, sshKeysPageSize]);

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
      const response = await hostService.list(page, pageSize, {});
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

  const loadSSHKeys = async () => {
    setSshKeysLoading(true);
    try {
      const response = await sshKeyService.list(sshKeysPage, sshKeysPageSize);
      setSshKeys(response.ssh_keys);
      setSshKeysTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载SSH密钥列表失败');
    } finally {
      setSshKeysLoading(false);
    }
  };

  const handleCreateSSHKey = () => {
    setEditingSshKey(null);
    sshKeyForm.resetFields();
    setSshKeyModalVisible(true);
  };

  const handleEditSSHKey = (key: SSHKey) => {
    setEditingSshKey(key);
    sshKeyForm.setFieldsValue({ name: key.name, username: key.username || '' });
    setSshKeyModalVisible(true);
  };

  const handleDeleteSSHKey = async (id: number) => {
    try {
      await sshKeyService.delete(id);
      message.success('删除成功');
      loadSSHKeys();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSSHKeySubmit = async () => {
    try {
      const values = await sshKeyForm.validateFields();
      if (editingSshKey) {
        await sshKeyService.update(editingSshKey.id, { 
          name: values.name,
          username: values.username || undefined,
        });
        message.success('更新成功');
      } else {
        await sshKeyService.create({
          name: values.name,
          username: values.username || undefined,
          private_key_content: values.private_key_content,
        });
        message.success('创建成功');
      }
      setSshKeyModalVisible(false);
      loadSSHKeys();
    } catch (error: any) {
      if (error.errorFields) {
        return; // 表单验证错误
      }
      message.error(error.response?.data?.error || '操作失败');
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
        <TabPane tab="SSH密钥" key="ssh-keys" closable={false}>
          <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h2 style={{ margin: 0 }}>SSH密钥管理</h2>
            <Space wrap>
              <Button icon={<ReloadOutlined />} onClick={loadSSHKeys}>
                刷新
              </Button>
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateSSHKey}>
                上传密钥
              </Button>
            </Space>
          </div>

          <Table
            columns={[
              {
                title: 'ID',
                dataIndex: 'id',
                key: 'id',
                width: 80,
              },
              {
                title: '名称',
                dataIndex: 'name',
                key: 'name',
              },
              {
                title: '用户名',
                dataIndex: 'username',
                key: 'username',
                width: 120,
                render: (username: string) => username || <span style={{ color: '#999' }}>-</span>,
              },
              {
                title: '类型',
                dataIndex: 'key_type',
                key: 'key_type',
                width: 100,
                render: (type: string) => <Tag>{type.toUpperCase()}</Tag>,
              },
              {
                title: '指纹',
                dataIndex: 'fingerprint',
                key: 'fingerprint',
                render: (fp: string) => (
                  <span style={{ fontFamily: 'monospace', fontSize: '12px' }}>{fp || '-'}</span>
                ),
              },
              {
                title: '创建时间',
                dataIndex: 'created_at',
                key: 'created_at',
                render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
              },
              {
                title: '操作',
                key: 'action',
                width: 150,
                render: (_: any, record: SSHKey) => (
                  <Space>
                    <Button
                      type="link"
                      size="small"
                      icon={<EditOutlined />}
                      onClick={() => handleEditSSHKey(record)}
                    >
                      编辑
                    </Button>
                    <Popconfirm
                      title="确定要删除这个SSH密钥吗？"
                      onConfirm={() => handleDeleteSSHKey(record.id)}
                      okText="确定"
                      cancelText="取消"
                    >
                      <Button
                        type="link"
                        danger
                        size="small"
                        icon={<DeleteOutlined />}
                      >
                        删除
                      </Button>
                    </Popconfirm>
                  </Space>
                ),
              },
            ]}
            dataSource={sshKeys}
            rowKey="id"
            loading={sshKeysLoading}
            pagination={{
              current: sshKeysPage,
              pageSize: sshKeysPageSize,
              total: sshKeysTotal,
              showSizeChanger: true,
              showTotal: (total) => `共 ${total} 条`,
              onChange: (page, pageSize) => {
                setSshKeysPage(page);
                setSshKeysPageSize(pageSize);
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

      {/* SSH密钥上传/编辑模态框 */}
      <Modal
        title={editingSshKey ? '编辑SSH密钥' : '上传SSH密钥'}
        open={sshKeyModalVisible}
        onOk={handleSSHKeySubmit}
        onCancel={() => setSshKeyModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={sshKeyForm} layout="vertical">
          <Form.Item
            name="name"
            label="密钥名称"
            rules={[{ required: true, message: '请输入密钥名称' }]}
          >
            <Input placeholder="例如: 我的SSH密钥" />
          </Form.Item>
          <Form.Item
            name="username"
            label="SSH用户名"
            tooltip="可选：设置默认SSH用户名，连接时将自动使用此用户名"
          >
            <Input placeholder="例如: root, ubuntu, admin" />
          </Form.Item>
          {!editingSshKey && (
            <Form.Item
              name="private_key_content"
              label="私钥内容"
              rules={[{ required: true, message: '请输入SSH私钥内容' }]}
            >
              <Input.TextArea
                rows={10}
                placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
                style={{ fontFamily: 'monospace', fontSize: '12px' }}
              />
            </Form.Item>
          )}
        </Form>
      </Modal>
    </Card>
  );
};

export default WebSSH;

