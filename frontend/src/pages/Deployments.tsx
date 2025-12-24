import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Tag,
  Tabs,
  Descriptions,
  Progress,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  PlayCircleOutlined,
  RollbackOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import {
  deploymentService,
  DeploymentConfig,
  Deployment,
  DeploymentVersion,
  CreateDeploymentConfigRequest,
  CreateDeploymentRequest,
  CreateDeploymentVersionRequest,
} from '../services/deployment';
import { hostService, Host } from '../services/host';

const { Option } = Select;
const { TabPane } = Tabs;

const Deployments: React.FC = () => {
  const [activeTab, setActiveTab] = useState('configs');
  
  // 部署配置相关状态
  const [configs, setConfigs] = useState<DeploymentConfig[]>([]);
  const [configLoading, setConfigLoading] = useState(false);
  const [configTotal, setConfigTotal] = useState(0);
  const [configPage, setConfigPage] = useState(1);
  const [configPageSize, setConfigPageSize] = useState(10);
  const [configModalVisible, setConfigModalVisible] = useState(false);
  const [editingConfig, setEditingConfig] = useState<DeploymentConfig | null>(null);
  const [configForm] = Form.useForm();

  // 部署执行相关状态
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [deploymentLoading, setDeploymentLoading] = useState(false);
  const [deploymentTotal, setDeploymentTotal] = useState(0);
  const [deploymentPage, setDeploymentPage] = useState(1);
  const [deploymentPageSize, setDeploymentPageSize] = useState(10);
  const [deploymentModalVisible, setDeploymentModalVisible] = useState(false);
  const [deploymentDetailVisible, setDeploymentDetailVisible] = useState(false);
  const [selectedDeployment, setSelectedDeployment] = useState<Deployment | null>(null);
  const [hosts, setHosts] = useState<Host[]>([]);
  const [deploymentForm] = Form.useForm();

  // 部署版本相关状态
  const [versions, setVersions] = useState<DeploymentVersion[]>([]);
  const [versionLoading, setVersionLoading] = useState(false);
  const [versionTotal, setVersionTotal] = useState(0);
  const [versionPage, setVersionPage] = useState(1);
  const [versionPageSize, setVersionPageSize] = useState(10);
  const [versionModalVisible, setVersionModalVisible] = useState(false);
  const [versionForm] = Form.useForm();

  useEffect(() => {
    loadHosts();
  }, []);

  useEffect(() => {
    if (activeTab === 'configs') {
      loadConfigs();
    } else if (activeTab === 'deployments') {
      loadDeployments();
    } else if (activeTab === 'versions') {
      loadVersions();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, configPage, configPageSize, deploymentPage, deploymentPageSize, versionPage, versionPageSize]);

  const loadHosts = async () => {
    try {
      const response = await hostService.list(1, 1000);
      setHosts(response.hosts);
    } catch (error: any) {
      // 静默失败
    }
  };

  const loadConfigs = async () => {
    setConfigLoading(true);
    try {
      const response = await deploymentService.listConfigs(configPage, configPageSize);
      setConfigs(response.configs);
      setConfigTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载部署配置列表失败');
    } finally {
      setConfigLoading(false);
    }
  };

  const loadDeployments = async () => {
    setDeploymentLoading(true);
    try {
      const response = await deploymentService.listDeployments(deploymentPage, deploymentPageSize);
      setDeployments(response.deployments);
      setDeploymentTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载部署列表失败');
    } finally {
      setDeploymentLoading(false);
    }
  };

  const loadVersions = async () => {
    setVersionLoading(true);
    try {
      const response = await deploymentService.listVersions(versionPage, versionPageSize);
      setVersions(response.versions);
      setVersionTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载版本列表失败');
    } finally {
      setVersionLoading(false);
    }
  };

  const handleCreateConfig = () => {
    setEditingConfig(null);
    configForm.resetFields();
    setConfigModalVisible(true);
  };

  const handleEditConfig = async (config: DeploymentConfig) => {
    try {
      const response = await deploymentService.getConfig(config.id);
      setEditingConfig(response.config);
      configForm.setFieldsValue({
        name: response.config.name,
        application_name: response.config.application_name,
        description: response.config.description,
        salt_state_files: response.config.salt_state_files,
        environment: response.config.environment,
      });
      setConfigModalVisible(true);
    } catch (error: any) {
      message.error('获取部署配置详情失败');
    }
  };

  const handleDeleteConfig = async (id: number) => {
    try {
      await deploymentService.deleteConfig(id);
      message.success('删除成功');
      loadConfigs();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleConfigSubmit = async () => {
    try {
      const values = await configForm.validateFields();
      if (editingConfig) {
        await deploymentService.updateConfig(editingConfig.id, values);
        message.success('更新成功');
      } else {
        await deploymentService.createConfig(values);
        message.success('创建成功');
      }
      setConfigModalVisible(false);
      loadConfigs();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleCreateDeployment = () => {
    deploymentForm.resetFields();
    setDeploymentModalVisible(true);
  };

  const handleDeploymentSubmit = async () => {
    try {
      const values = await deploymentForm.validateFields();
      await deploymentService.createDeployment({
        config_id: values.config_id,
        version: values.version,
        target_hosts: values.target_hosts,
      });
      message.success('部署已启动');
      setDeploymentModalVisible(false);
      loadDeployments();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '启动部署失败');
    }
  };

  const handleViewDeployment = async (deployment: Deployment) => {
    try {
      const response = await deploymentService.getDeployment(deployment.id);
      setSelectedDeployment(response.deployment);
      setDeploymentDetailVisible(true);
    } catch (error: any) {
      message.error('获取部署详情失败');
    }
  };

  const handleRollback = async (id: number) => {
    try {
      await deploymentService.rollbackDeployment(id);
      message.success('回滚已启动');
      loadDeployments();
    } catch (error: any) {
      message.error(error.response?.data?.error || '回滚失败');
    }
  };

  const handleCreateVersion = () => {
    versionForm.resetFields();
    setVersionModalVisible(true);
  };

  const handleVersionSubmit = async () => {
    try {
      const values = await versionForm.validateFields();
      await deploymentService.createVersion(values);
      message.success('创建成功');
      setVersionModalVisible(false);
      loadVersions();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending: { color: 'default', text: '等待中' },
      running: { color: 'processing', text: '运行中' },
      completed: { color: 'success', text: '已完成' },
      failed: { color: 'error', text: '失败' },
      cancelled: { color: 'warning', text: '已取消' },
    };
    const statusInfo = statusMap[status] || { color: 'default', text: status };
    return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
  };

  const configColumns = [
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
      title: '应用名称',
      dataIndex: 'application_name',
      key: 'application_name',
    },
    {
      title: '环境',
      dataIndex: 'environment',
      key: 'environment',
      width: 100,
      render: (env: string) => env || '-',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: DeploymentConfig) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditConfig(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个部署配置吗？"
            onConfirm={() => handleDeleteConfig(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger size="small" icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const deploymentColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '配置',
      key: 'config',
      width: 150,
      render: (_: any, record: Deployment) => (
        <span>{record.config?.name || `配置 #${record.config_id}`}</span>
      ),
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '完成时间',
      dataIndex: 'completed_at',
      key: 'completed_at',
      width: 180,
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '持续时间',
      key: 'duration',
      width: 100,
      render: (_: any, record: Deployment) => {
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
      width: 200,
      render: (_: any, record: Deployment) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDeployment(record)}
          >
            详情
          </Button>
          {record.status === 'completed' && !record.is_rollback && (
            <Popconfirm
              title="确定要回滚这个部署吗？"
              onConfirm={() => handleRollback(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" size="small" icon={<RollbackOutlined />}>
                回滚
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  const versionColumns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '应用名称',
      dataIndex: 'application_name',
      key: 'application_name',
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: '发布说明',
      dataIndex: 'release_notes',
      key: 'release_notes',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
    },
  ];

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between' }}>
          <h2 style={{ margin: 0 }}>部署管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={() => {
              if (activeTab === 'configs') loadConfigs();
              else if (activeTab === 'deployments') loadDeployments();
              else if (activeTab === 'versions') loadVersions();
            }}>
              刷新
            </Button>
            {activeTab === 'configs' && (
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateConfig}>
                创建配置
              </Button>
            )}
            {activeTab === 'deployments' && (
              <Button type="primary" icon={<PlayCircleOutlined />} onClick={handleCreateDeployment}>
                执行部署
              </Button>
            )}
            {activeTab === 'versions' && (
              <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateVersion}>
                创建版本
              </Button>
            )}
          </Space>
        </div>

        <Tabs activeKey={activeTab} onChange={setActiveTab}>
          <TabPane tab="部署配置" key="configs">
            <Table
              columns={configColumns}
              dataSource={configs}
              rowKey="id"
              loading={configLoading}
              pagination={{
                current: configPage,
                pageSize: configPageSize,
                total: configTotal,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条`,
                onChange: (page, pageSize) => {
                  setConfigPage(page);
                  setConfigPageSize(pageSize);
                },
              }}
            />
          </TabPane>
          <TabPane tab="部署历史" key="deployments">
            <Table
              columns={deploymentColumns}
              dataSource={deployments}
              rowKey="id"
              loading={deploymentLoading}
              pagination={{
                current: deploymentPage,
                pageSize: deploymentPageSize,
                total: deploymentTotal,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条`,
                onChange: (page, pageSize) => {
                  setDeploymentPage(page);
                  setDeploymentPageSize(pageSize);
                },
              }}
            />
          </TabPane>
          <TabPane tab="版本管理" key="versions">
            <Table
              columns={versionColumns}
              dataSource={versions}
              rowKey="id"
              loading={versionLoading}
              pagination={{
                current: versionPage,
                pageSize: versionPageSize,
                total: versionTotal,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条`,
                onChange: (page, pageSize) => {
                  setVersionPage(page);
                  setVersionPageSize(pageSize);
                },
              }}
            />
          </TabPane>
        </Tabs>
      </Card>

      {/* 部署配置模态框 */}
      <Modal
        title={editingConfig ? '编辑部署配置' : '创建部署配置'}
        open={configModalVisible}
        onOk={handleConfigSubmit}
        onCancel={() => setConfigModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={configForm} layout="vertical">
          <Form.Item
            name="name"
            label="配置名称"
            rules={[{ required: true, message: '请输入配置名称' }]}
          >
            <Input placeholder="请输入配置名称" />
          </Form.Item>
          <Form.Item
            name="application_name"
            label="应用名称"
            rules={[{ required: true, message: '请输入应用名称' }]}
          >
            <Input placeholder="请输入应用名称" />
          </Form.Item>
          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea rows={3} placeholder="请输入描述" />
          </Form.Item>
          <Form.Item
            name="salt_state_files"
            label="Salt State 文件"
            rules={[{ required: true, message: '请输入 Salt State 文件' }]}
          >
            <Input placeholder="例如: apache,nginx" />
          </Form.Item>
          <Form.Item
            name="environment"
            label="环境"
          >
            <Select placeholder="请选择环境" allowClear>
              <Option value="dev">开发</Option>
              <Option value="test">测试</Option>
              <Option value="prod">生产</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      {/* 执行部署模态框 */}
      <Modal
        title="执行部署"
        open={deploymentModalVisible}
        onOk={handleDeploymentSubmit}
        onCancel={() => setDeploymentModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={deploymentForm} layout="vertical">
          <Form.Item
            name="config_id"
            label="部署配置"
            rules={[{ required: true, message: '请选择部署配置' }]}
          >
            <Select placeholder="请选择部署配置">
              {configs.map((config) => (
                <Option key={config.id} value={config.id}>
                  {config.name} ({config.application_name})
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="version"
            label="版本"
            rules={[{ required: true, message: '请输入版本号' }]}
          >
            <Input placeholder="例如: v1.0.0" />
          </Form.Item>
          <Form.Item
            name="target_hosts"
            label="目标主机"
            rules={[{ required: true, message: '请选择目标主机' }]}
          >
            <Select mode="multiple" placeholder="请选择目标主机">
              {hosts.map((host) => (
                <Option key={host.id} value={host.id}>
                  {host.hostname} ({host.ip_address})
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      {/* 部署详情模态框 */}
      <Modal
        title="部署详情"
        open={deploymentDetailVisible}
        onCancel={() => setDeploymentDetailVisible(false)}
        footer={null}
        width={800}
      >
        {selectedDeployment && (
          <Descriptions bordered column={2}>
            <Descriptions.Item label="ID" span={1}>
              {selectedDeployment.id}
            </Descriptions.Item>
            <Descriptions.Item label="状态" span={1}>
              {getStatusTag(selectedDeployment.status)}
            </Descriptions.Item>
            <Descriptions.Item label="配置" span={1}>
              {selectedDeployment.config?.name || `配置 #${selectedDeployment.config_id}`}
            </Descriptions.Item>
            <Descriptions.Item label="版本" span={1}>
              {selectedDeployment.version}
            </Descriptions.Item>
            <Descriptions.Item label="开始时间" span={1}>
              {new Date(selectedDeployment.started_at).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="完成时间" span={1}>
              {selectedDeployment.completed_at
                ? new Date(selectedDeployment.completed_at).toLocaleString()
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="持续时间" span={1}>
              {selectedDeployment.duration_seconds
                ? `${selectedDeployment.duration_seconds}秒`
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="Salt Job ID" span={1}>
              {selectedDeployment.salt_job_id || '-'}
            </Descriptions.Item>
            {selectedDeployment.error_message && (
              <Descriptions.Item label="错误信息" span={2}>
                <div style={{ color: 'red' }}>{selectedDeployment.error_message}</div>
              </Descriptions.Item>
            )}
            {selectedDeployment.results && (
              <Descriptions.Item label="执行结果" span={2}>
                <pre style={{ maxHeight: 300, overflow: 'auto', background: '#f5f5f5', padding: 8, borderRadius: 4 }}>
                  {JSON.stringify(JSON.parse(selectedDeployment.results), null, 2)}
                </pre>
              </Descriptions.Item>
            )}
          </Descriptions>
        )}
      </Modal>

      {/* 创建版本模态框 */}
      <Modal
        title="创建版本"
        open={versionModalVisible}
        onOk={handleVersionSubmit}
        onCancel={() => setVersionModalVisible(false)}
        okText="确定"
        cancelText="取消"
      >
        <Form form={versionForm} layout="vertical">
          <Form.Item
            name="application_name"
            label="应用名称"
            rules={[{ required: true, message: '请输入应用名称' }]}
          >
            <Input placeholder="请输入应用名称" />
          </Form.Item>
          <Form.Item
            name="version"
            label="版本号"
            rules={[{ required: true, message: '请输入版本号' }]}
          >
            <Input placeholder="例如: v1.0.0" />
          </Form.Item>
          <Form.Item
            name="release_notes"
            label="发布说明"
          >
            <Input.TextArea rows={4} placeholder="请输入发布说明" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default Deployments;

