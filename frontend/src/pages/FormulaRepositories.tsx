import React, { useState, useEffect } from 'react';
import {
  Card,
  Row,
  Col,
  Button,
  Table,
  Modal,
  Form,
  Input,
  Space,
  message,
  Tag,
  Tooltip,
  Popconfirm,
} from 'antd';
import {
  PlusOutlined,
  SyncOutlined,
  DeleteOutlined,
  EditOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import formulaService, { FormulaRepository } from '../services/formula';

const { TextArea } = Input;

const FormulaRepositories: React.FC = () => {
  const [repositories, setRepositories] = useState<FormulaRepository[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRepo, setEditingRepo] = useState<FormulaRepository | undefined>(undefined);
  const [form] = Form.useForm();

  // 加载仓库列表
  const loadRepositories = async () => {
    setLoading(true);
    try {
      const response = await formulaService.listRepositories(1, 100);
      setRepositories(response.repositories || []);
    } catch (error: any) {
      message.error('加载仓库列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadRepositories();
  }, []);

  // 同步仓库
  const handleSyncRepository = async (repoId: number) => {
    try {
      await formulaService.syncRepository(repoId);
      message.success('仓库同步成功');
      loadRepositories(); // 重新加载列表
    } catch (error: any) {
      message.error(error.message || '同步仓库失败');
    }
  };

  // 删除仓库
  const handleDeleteRepository = async (repoId: number) => {
    // TODO: 实现删除仓库API
    message.warning('删除仓库功能暂未实现');
  };

  // 打开添加/编辑模态框
  const openModal = (repo?: FormulaRepository) => {
    setEditingRepo(repo);
    if (repo) {
      form.setFieldsValue({
        name: repo.name,
        url: repo.url,
        branch: repo.branch || 'master',
        local_path: repo.local_path,
        is_active: repo.is_active,
      });
    } else {
      form.resetFields();
      form.setFieldsValue({
        branch: 'master',
      });
    }
    setModalVisible(true);
  };

  // 关闭模态框
  const closeModal = () => {
    setModalVisible(false);
    setEditingRepo(undefined);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      if (editingRepo) {
        // 编辑仓库
        await formulaService.updateRepository(editingRepo.id, {
          name: values.name,
          url: values.url,
          branch: values.branch,
          local_path: values.local_path,
          is_active: values.is_active !== undefined ? values.is_active : true,
        });
        message.success('仓库更新成功');
        loadRepositories();
      } else {
        // 创建新仓库
        await formulaService.createRepository({
          name: values.name,
          url: values.url,
          branch: values.branch,
          local_path: values.local_path,
        });
        message.success('仓库创建成功');
        loadRepositories();
      }

      closeModal();
    } catch (error: any) {
      if (error.message) {
        message.error(error.message);
      }
    }
  };

  // 表格列定义
  const columns = [
    {
      title: '仓库名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: FormulaRepository) => (
        <div>
          <div style={{ fontWeight: 'bold' }}>{text}</div>
          <div style={{ fontSize: '12px', color: '#666' }}>
            {record.url}
          </div>
        </div>
      ),
    },
    {
      title: '分支',
      dataIndex: 'branch',
      key: 'branch',
      render: (branch: string) => (
        <Tag color="blue">{branch}</Tag>
      ),
    },
    {
      title: '状态',
      key: 'status',
      render: (record: FormulaRepository) => (
        <div>
          <Tag color={record.is_active ? 'green' : 'red'}>
            {record.is_active ? '启用' : '禁用'}
          </Tag>
          {record.last_sync_at && (
            <div style={{ fontSize: '12px', color: '#666', marginTop: '4px' }}>
              最后同步: {new Date(record.last_sync_at).toLocaleString()}
            </div>
          )}
        </div>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: FormulaRepository) => (
        <Space size="small">
          <Tooltip title="同步仓库">
            <Button
              size="small"
              icon={<SyncOutlined />}
              onClick={() => handleSyncRepository(record.id)}
            />
          </Tooltip>
          <Tooltip title="编辑仓库">
            <Button
              size="small"
              icon={<EditOutlined />}
              onClick={() => openModal(record)}
            />
          </Tooltip>
          <Tooltip title="删除仓库">
            <Popconfirm
              title="确定要删除这个仓库吗？"
              description="删除后将无法恢复，请谨慎操作"
              onConfirm={() => handleDeleteRepository(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button
                size="small"
                danger
                icon={<DeleteOutlined />}
              />
            </Popconfirm>
          </Tooltip>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h2>Formula仓库管理</h2>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={loadRepositories}>
            刷新
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()}>
            添加仓库
          </Button>
        </Space>
      </div>

      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card>
            <Table
              columns={columns}
              dataSource={repositories}
              rowKey="id"
              loading={loading}
              pagination={false}
            />
          </Card>
        </Col>
      </Row>

      {/* 添加/编辑仓库模态框 */}
      <Modal
        title={editingRepo ? '编辑仓库' : '添加仓库'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={closeModal}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="仓库名称"
            name="name"
            rules={[{ required: true, message: '请输入仓库名称' }]}
          >
            <Input placeholder="例如: salt-formulas" />
          </Form.Item>

          <Form.Item
            label="Git仓库地址"
            name="url"
            rules={[
              { required: true, message: '请输入Git仓库地址' },
              { type: 'url', message: '请输入有效的URL' }
            ]}
          >
            <Input placeholder="https://github.com/your-org/salt-formulas.git" />
          </Form.Item>

          <Form.Item
            label="分支"
            name="branch"
            rules={[{ required: true, message: '请输入分支名称' }]}
          >
            <Input placeholder="master 或 main" />
          </Form.Item>

          <Form.Item
            label="本地路径 (可选)"
            name="local_path"
            tooltip="如果不填写，将自动生成路径"
          >
            <Input placeholder="/opt/salt-formulas" />
          </Form.Item>

          {editingRepo && (
            <Form.Item
              label="状态"
              name="is_active"
              valuePropName="checked"
            >
              <input type="checkbox" style={{ marginRight: '8px' }} />
              启用仓库
            </Form.Item>
          )}

          <div style={{ marginTop: '16px', padding: '12px', backgroundColor: '#f6f8fa', borderRadius: '6px' }}>
            <h4>📋 配置说明</h4>
            <ul style={{ margin: 0, paddingLeft: '20px' }}>
              <li>仓库地址支持 HTTPS 和 SSH 格式</li>
              <li>分支默认为 master，可根据需要修改</li>
              <li>本地路径不填时会自动生成</li>
              <li>创建后可以通过"同步"按钮拉取Formula</li>
              {editingRepo && <li>编辑模式下可以修改仓库状态</li>}
            </ul>
          </div>
        </Form>
      </Modal>
    </div>
  );
};

export default FormulaRepositories;
