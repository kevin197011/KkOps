import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  message,
  Popconfirm,
  Tag,
  Card,
  Tabs,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import {
  sshService,
  SSHConnection,
  SSHKey,
  SSHSession,
  CreateSSHConnectionRequest,
  CreateSSHKeyRequest,
  CreateSSHSessionRequest,
} from '../services/ssh';
import { projectService, Project } from '../services/project';
import { hostService, Host } from '../services/host';
import { useAuth } from '../contexts/AuthContext';

const { Option } = Select;
const { TabPane } = Tabs;

const SSH: React.FC = () => {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState('connections');
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined);
  const [hosts, setHosts] = useState<Host[]>([]);

  // SSH连接相关状态
  const [connections, setConnections] = useState<SSHConnection[]>([]);
  const [connectionLoading, setConnectionLoading] = useState(false);
  const [connectionTotal, setConnectionTotal] = useState(0);
  const [connectionPage, setConnectionPage] = useState(1);
  const [connectionPageSize, setConnectionPageSize] = useState(10);
  const [connectionModalVisible, setConnectionModalVisible] = useState(false);
  const [editingConnection, setEditingConnection] = useState<SSHConnection | null>(null);
  const [connectionForm] = Form.useForm();

  // SSH密钥相关状态
  const [keys, setKeys] = useState<SSHKey[]>([]);
  const [keyLoading, setKeyLoading] = useState(false);
  const [keyTotal, setKeyTotal] = useState(0);
  const [keyPage, setKeyPage] = useState(1);
  const [keyPageSize, setKeyPageSize] = useState(10);
  const [keyModalVisible, setKeyModalVisible] = useState(false);
  const [editingKey, setEditingKey] = useState<SSHKey | null>(null);
  const [keyForm] = Form.useForm();

  // SSH会话相关状态
  const [sessions, setSessions] = useState<SSHSession[]>([]);
  const [sessionLoading, setSessionLoading] = useState(false);
  const [sessionTotal, setSessionTotal] = useState(0);
  const [sessionPage, setSessionPage] = useState(1);
  const [sessionPageSize, setSessionPageSize] = useState(10);
  const [sessionModalVisible, setSessionModalVisible] = useState(false);
  const [selectedConnectionId, setSelectedConnectionId] = useState<number | undefined>(undefined);
  const [selectedStatus, setSelectedStatus] = useState<string | undefined>(undefined);
  const [sessionForm] = Form.useForm();

  useEffect(() => {
    loadProjects();
    loadHosts();
  }, []);

  useEffect(() => {
    if (activeTab === 'connections') {
      loadConnections();
    } else if (activeTab === 'keys') {
      loadKeys();
    } else if (activeTab === 'sessions') {
      loadSessions();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, connectionPage, connectionPageSize, keyPage, keyPageSize, sessionPage, sessionPageSize, selectedProjectId, selectedConnectionId, selectedStatus]);

  const loadProjects = async () => {
    try {
      const response = await projectService.list(1, 100);
      setProjects(response.projects);
    } catch (error: any) {
      message.error('加载项目列表失败');
    }
  };

  const loadHosts = async () => {
    try {
      const response = await hostService.list(1, 1000);
      setHosts(response.hosts);
    } catch (error: any) {
      // 静默失败，不影响主流程
    }
  };

  const loadConnections = async () => {
    setConnectionLoading(true);
    try {
      const response = await sshService.listConnections(connectionPage, connectionPageSize, {
        project_id: selectedProjectId,
      });
      setConnections(response.connections || []);
      setConnectionTotal(response.total || 0);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载SSH连接列表失败');
    } finally {
      setConnectionLoading(false);
    }
  };

  const loadKeys = async () => {
    setKeyLoading(true);
    try {
      const response = await sshService.listKeys(keyPage, keyPageSize);
      setKeys(response.keys || []);
      setKeyTotal(response.total || 0);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载SSH密钥列表失败');
    } finally {
      setKeyLoading(false);
    }
  };

  const loadSessions = async () => {
    setSessionLoading(true);
    try {
      const filters: any = {};
      if (selectedConnectionId) filters.connection_id = selectedConnectionId;
      if (selectedStatus) filters.status = selectedStatus;
      if (user?.id) filters.user_id = user.id; // 默认只显示当前用户的会话
      const response = await sshService.listSessions(sessionPage, sessionPageSize, filters);
      setSessions(response.sessions || []);
      setSessionTotal(response.total || 0);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载SSH会话列表失败');
    } finally {
      setSessionLoading(false);
    }
  };

  const handleCreateSession = () => {
    sessionForm.resetFields();
    setSessionModalVisible(true);
  };

  const handleSessionSubmit = async () => {
    try {
      const values = await sessionForm.validateFields();
      await sshService.createSession({
        connection_id: values.connection_id,
        client_ip: values.client_ip,
      });
      message.success('会话创建成功');
      setSessionModalVisible(false);
      loadSessions();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '创建会话失败');
    }
  };

  const handleCloseSession = async (id: number) => {
    try {
      await sshService.closeSession(id);
      message.success('会话已关闭');
      loadSessions();
    } catch (error: any) {
      message.error(error.response?.data?.error || '关闭会话失败');
    }
  };

  const handleCreateConnection = () => {
    setEditingConnection(null);
    connectionForm.resetFields();
    connectionForm.setFieldsValue({ port: 22, auth_type: 'password' });
    setConnectionModalVisible(true);
  };

  const handleHostSelect = (hostId: number | undefined) => {
    if (hostId) {
      const selectedHost = hosts.find((h) => h.id === hostId);
      if (selectedHost) {
        connectionForm.setFieldsValue({
          host_id: hostId,
          hostname: selectedHost.ip_address || selectedHost.hostname,
          port: selectedHost.ssh_port || 22,
        });
      }
    } else {
      // 清空主机选择时，清空自动填充的字段
      connectionForm.setFieldsValue({
        host_id: undefined,
        hostname: undefined,
        port: 22,
      });
    }
  };

  const handleEditConnection = async (connection: SSHConnection) => {
    try {
      const response = await sshService.getConnection(connection.id);
      const fullConnection = response.connection;
      setEditingConnection(fullConnection);
      connectionForm.setFieldsValue({
        ...fullConnection,
        project_id: (fullConnection as any).project_id,
      });
      setConnectionModalVisible(true);
    } catch (error: any) {
      message.error('获取SSH连接详情失败');
    }
  };

  const handleDeleteConnection = async (id: number) => {
    try {
      await sshService.deleteConnection(id);
      message.success('删除成功');
      loadConnections();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleConnectionSubmit = async () => {
    try {
      const values = await connectionForm.validateFields();
      if (editingConnection) {
        await sshService.updateConnection(editingConnection.id, values);
        message.success('更新成功');
      } else {
        await sshService.createConnection(values as CreateSSHConnectionRequest);
        message.success('创建成功');
      }
      setConnectionModalVisible(false);
      loadConnections();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleCreateKey = () => {
    setEditingKey(null);
    keyForm.resetFields();
    keyForm.setFieldsValue({ key_type: 'rsa' });
    setKeyModalVisible(true);
  };

  const handleEditKey = async (key: SSHKey) => {
    try {
      const response = await sshService.getKey(key.id);
      const fullKey = response.key;
      setEditingKey(fullKey);
      keyForm.setFieldsValue(fullKey);
      setKeyModalVisible(true);
    } catch (error: any) {
      message.error('获取SSH密钥详情失败');
    }
  };

  const handleDeleteKey = async (id: number) => {
    try {
      await sshService.deleteKey(id);
      message.success('删除成功');
      loadKeys();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleKeySubmit = async () => {
    try {
      const values = await keyForm.validateFields();
      if (editingKey) {
        await sshService.updateKey(editingKey.id, values);
        message.success('更新成功');
      } else {
        await sshService.createKey(values as CreateSSHKeyRequest);
        message.success('创建成功');
      }
      setKeyModalVisible(false);
      loadKeys();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const connectionColumns = [
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
      title: '主机',
      dataIndex: 'hostname',
      key: 'hostname',
    },
    {
      title: '端口',
      dataIndex: 'port',
      key: 'port',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '认证方式',
      dataIndex: 'auth_type',
      key: 'auth_type',
      render: (type: string) => (
        <Tag color={type === 'key' ? 'blue' : 'green'}>
          {type === 'key' ? '密钥' : '密码'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'success' : 'default'}>
          {status === 'active' ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '最后连接',
      dataIndex: 'last_connected_at',
      key: 'last_connected_at',
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: SSHConnection) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEditConnection(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个SSH连接吗？"
            onConfirm={() => handleDeleteConnection(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const keyColumns = [
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
      title: '类型',
      dataIndex: 'key_type',
      key: 'key_type',
      render: (type: string) => <Tag>{type.toUpperCase()}</Tag>,
    },
    {
      title: '公钥',
      dataIndex: 'public_key',
      key: 'public_key',
      ellipsis: true,
    },
    {
      title: '指纹',
      dataIndex: 'fingerprint',
      key: 'fingerprint',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: SSHKey) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEditKey(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个SSH密钥吗？"
            onConfirm={() => handleDeleteKey(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const sessionColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '连接ID',
      dataIndex: 'connection_id',
      key: 'connection_id',
      width: 100,
    },
    {
      title: '用户ID',
      dataIndex: 'user_id',
      key: 'user_id',
      width: 100,
    },
    {
      title: '客户端IP',
      dataIndex: 'client_ip',
      key: 'client_ip',
      width: 150,
      render: (text: string) => text || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => {
        const statusMap: Record<string, { color: string; text: string }> = {
          active: { color: 'success', text: '活跃' },
          closed: { color: 'default', text: '已关闭' },
          timeout: { color: 'warning', text: '超时' },
        };
        const statusInfo = statusMap[status] || { color: 'default', text: status };
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
      },
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '结束时间',
      dataIndex: 'ended_at',
      key: 'ended_at',
      width: 180,
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '持续时间',
      key: 'duration',
      width: 120,
      render: (_: any, record: SSHSession) => {
        if (record.duration_seconds) {
          const hours = Math.floor(record.duration_seconds / 3600);
          const minutes = Math.floor((record.duration_seconds % 3600) / 60);
          const seconds = record.duration_seconds % 60;
          if (hours > 0) {
            return `${hours}时${minutes}分${seconds}秒`;
          } else if (minutes > 0) {
            return `${minutes}分${seconds}秒`;
          } else {
            return `${seconds}秒`;
          }
        }
        return '-';
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_: any, record: SSHSession) => (
        <Space>
          {record.status === 'active' && (
            <Popconfirm
              title="确定要关闭这个会话吗？"
              onConfirm={() => handleCloseSession(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" danger size="small">
                关闭
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2 style={{ margin: 0 }}>SSH管理</h2>
          <Space>
            {activeTab === 'connections' && (
              <Select
                placeholder="选择项目"
                allowClear
                style={{ width: 200 }}
                value={selectedProjectId}
                onChange={(value) => {
                  setSelectedProjectId(value);
                  setConnectionPage(1);
                }}
              >
                {projects.map((project) => (
                  <Option key={project.id} value={project.id}>
                    {project.name}
                  </Option>
                ))}
              </Select>
            )}
            <Button icon={<ReloadOutlined />} onClick={() => {
              if (activeTab === 'connections') loadConnections();
              else loadKeys();
            }}>
              刷新
            </Button>
            {activeTab === 'connections' && (
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateConnection}>
                新建连接
              </Button>
            )}
            {activeTab === 'keys' && (
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateKey}>
                新建密钥
              </Button>
            )}
            {activeTab === 'sessions' && (
              <Space>
                <Select
                  placeholder="筛选连接"
                  style={{ width: 200 }}
                  value={selectedConnectionId}
                  onChange={(value) => {
                    setSelectedConnectionId(value);
                    setSessionPage(1);
                  }}
                  allowClear
                >
                  {connections.map((conn) => (
                    <Option key={conn.id} value={conn.id}>
                      {conn.name} ({conn.hostname})
                    </Option>
                  ))}
                </Select>
                <Select
                  placeholder="筛选状态"
                  style={{ width: 150 }}
                  value={selectedStatus}
                  onChange={(value) => {
                    setSelectedStatus(value);
                    setSessionPage(1);
                  }}
                  allowClear
                >
                  <Option value="active">活跃</Option>
                  <Option value="closed">已关闭</Option>
                  <Option value="timeout">超时</Option>
                </Select>
                <Button icon={<ReloadOutlined />} onClick={loadSessions}>
                  刷新
                </Button>
                <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateSession}>
                  新建会话
                </Button>
              </Space>
            )}
          </Space>
        </div>

        <Tabs activeKey={activeTab} onChange={setActiveTab}>
          <TabPane tab="SSH连接" key="connections">
            <Table
              columns={connectionColumns}
              dataSource={connections}
              rowKey="id"
              loading={connectionLoading}
              pagination={{
                current: connectionPage,
                pageSize: connectionPageSize,
                total: connectionTotal,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条`,
                onChange: (page, pageSize) => {
                  setConnectionPage(page);
                  setConnectionPageSize(pageSize);
                },
              }}
            />
          </TabPane>
          <TabPane tab="SSH密钥" key="keys">
            <Table
              columns={keyColumns}
              dataSource={keys}
              rowKey="id"
              loading={keyLoading}
              pagination={{
                current: keyPage,
                pageSize: keyPageSize,
                total: keyTotal,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条`,
                onChange: (page, pageSize) => {
                  setKeyPage(page);
                  setKeyPageSize(pageSize);
                },
              }}
            />
          </TabPane>
          <TabPane tab="SSH会话" key="sessions">
            <Table
              columns={sessionColumns}
              dataSource={sessions}
              rowKey="id"
              loading={sessionLoading}
              pagination={{
                current: sessionPage,
                pageSize: sessionPageSize,
                total: sessionTotal,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条`,
                onChange: (page, pageSize) => {
                  setSessionPage(page);
                  setSessionPageSize(pageSize);
                },
              }}
            />
          </TabPane>
        </Tabs>
      </Card>

      {/* SSH连接模态框 */}
      <Modal
        title={editingConnection ? '编辑SSH连接' : '新建SSH连接'}
        open={connectionModalVisible}
        onOk={handleConnectionSubmit}
        onCancel={() => setConnectionModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={connectionForm} layout="vertical">
          <Form.Item
            name="project_id"
            label="项目"
            rules={[{ required: true, message: '请选择项目' }]}
          >
            <Select placeholder="请选择项目">
              {projects.map((project) => (
                <Option key={project.id} value={project.id}>
                  {project.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="name"
            label="连接名称"
            rules={[{ required: true, message: '请输入连接名称' }]}
          >
            <Input placeholder="请输入连接名称" />
          </Form.Item>
          <Form.Item
            name="host_id"
            label="选择主机"
            tooltip="从主机管理中选择主机，将自动填充主机地址和SSH端口"
          >
            <Select
              placeholder="请选择主机（可选）"
              allowClear
              showSearch
              filterOption={(input, option) => {
                const label = option?.label || option?.children;
                if (typeof label === 'string') {
                  return label.toLowerCase().includes(input.toLowerCase());
                }
                return false;
              }}
              onChange={handleHostSelect}
            >
              {hosts.map((host) => (
                <Option key={host.id} value={host.id}>
                  {host.hostname} ({host.ip_address})
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="hostname"
            label="主机地址"
            rules={[{ required: true, message: '请输入主机地址' }]}
          >
            <Input placeholder="例如: 192.168.1.100" />
          </Form.Item>
          <Form.Item
            name="port"
            label="端口"
            initialValue={22}
            rules={[{ required: true, message: '请输入端口' }]}
          >
            <InputNumber min={1} max={65535} placeholder="SSH端口" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="SSH用户名" />
          </Form.Item>
          <Form.Item
            name="auth_type"
            label="认证方式"
            initialValue="password"
            rules={[{ required: true, message: '请选择认证方式' }]}
          >
            <Select>
              <Option value="password">密码</Option>
              <Option value="key">密钥</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      {/* SSH密钥模态框 */}
      <Modal
        title={editingKey ? '编辑SSH密钥' : '新建SSH密钥'}
        open={keyModalVisible}
        onOk={handleKeySubmit}
        onCancel={() => setKeyModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={keyForm} layout="vertical">
          <Form.Item
            name="name"
            label="密钥名称"
            rules={[{ required: true, message: '请输入密钥名称' }]}
          >
            <Input placeholder="请输入密钥名称" />
          </Form.Item>
          <Form.Item
            name="key_type"
            label="密钥类型"
            initialValue="rsa"
            rules={[{ required: true, message: '请选择密钥类型' }]}
          >
            <Select>
              <Option value="rsa">RSA</Option>
              <Option value="ed25519">ED25519</Option>
              <Option value="ecdsa">ECDSA</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="public_key"
            label="公钥"
            rules={[{ required: true, message: '请输入公钥' }]}
          >
            <Input.TextArea rows={4} placeholder="SSH公钥内容" />
          </Form.Item>
          <Form.Item
            name="fingerprint"
            label="指纹"
            rules={[{ required: true, message: '请输入指纹' }]}
          >
            <Input placeholder="SSH密钥指纹" />
          </Form.Item>
        </Form>
      </Modal>

      {/* SSH会话创建模态框 */}
      <Modal
        title="新建SSH会话"
        open={sessionModalVisible}
        onOk={handleSessionSubmit}
        onCancel={() => setSessionModalVisible(false)}
        okText="确定"
        cancelText="取消"
      >
        <Form form={sessionForm} layout="vertical">
          <Form.Item
            name="connection_id"
            label="SSH连接"
            rules={[{ required: true, message: '请选择SSH连接' }]}
          >
            <Select placeholder="请选择SSH连接">
              {connections
                .filter((conn) => conn.status === 'active')
                .map((conn) => (
                  <Option key={conn.id} value={conn.id}>
                    {conn.name} ({conn.hostname}:{conn.port})
                  </Option>
                ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="client_ip"
            label="客户端IP"
          >
            <Input placeholder="客户端IP地址（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default SSH;

