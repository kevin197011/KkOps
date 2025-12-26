import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Tag,
  Card,
  ColorPicker,
  Descriptions,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { environmentService, Environment, CreateEnvironmentRequest } from '../services/environment';

const Environments: React.FC = () => {
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingEnv, setEditingEnv] = useState<Environment | null>(null);
  const [detailEnv, setDetailEnv] = useState<Environment | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadEnvironments();
  }, [page, pageSize]);

  const loadEnvironments = async () => {
    setLoading(true);
    try {
      const response = await environmentService.list(page, pageSize);
      setEnvironments(response.environments);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载环境列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingEnv(null);
    form.resetFields();
    form.setFieldsValue({ color: '#1890ff', sort_order: 0 });
    setModalVisible(true);
  };

  const handleViewDetail = (env: Environment) => {
    setDetailEnv(env);
    setDetailModalVisible(true);
  };

  const handleEdit = (env: Environment) => {
    setEditingEnv(env);
    form.setFieldsValue({
      ...env,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await environmentService.delete(id);
      message.success('删除成功');
      loadEnvironments();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const color = typeof values.color === 'string' ? values.color : values.color?.toHexString();
      const data: CreateEnvironmentRequest = {
        ...values,
        color,
      };

      if (editingEnv) {
        await environmentService.update(editingEnv.id, data);
        message.success('更新成功');
      } else {
        await environmentService.create(data);
        message.success('创建成功');
      }
      setModalVisible(false);
      loadEnvironments();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const columns = [
    {
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '环境名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Environment) => (
        <Tag color={record.color || 'default'}>{name.toUpperCase()}</Tag>
      ),
    },
    {
      title: '颜色',
      dataIndex: 'color',
      key: 'color',
      width: 100,
      render: (color: string) => (
        <div
          style={{
            width: 24,
            height: 24,
            backgroundColor: color || '#1890ff',
            borderRadius: 4,
            border: '1px solid #d9d9d9',
          }}
        />
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: Environment) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个环境吗？"
            description="删除后，关联此环境的主机将失去环境标识"
            onConfirm={() => handleDelete(record.id)}
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

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2 style={{ margin: 0 }}>环境管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={loadEnvironments}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              新建环境
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={environments}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            pageSizeOptions: ['10', '20', '50'],
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: (page, pageSize) => {
              setPage(page);
              setPageSize(pageSize);
            },
          }}
        />
      </Card>

      <Modal
        title={editingEnv ? '编辑环境' : '新建环境'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="环境名称"
            rules={[
              { required: true, message: '请输入环境名称' },
              { pattern: /^[a-z][a-z0-9_-]*$/, message: '只能包含小写字母、数字、下划线和连字符，且以字母开头' },
            ]}
            tooltip="用于系统内部标识，如 dev, test, prod"
          >
            <Input placeholder="例如: dev, test, staging, prod" />
          </Form.Item>
          <Form.Item name="color" label="颜色" initialValue="#1890ff">
            <ColorPicker showText />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="环境描述信息" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 详情模态框 */}
      <Modal
        title="环境详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={500}
      >
        {detailEnv && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="环境名称">
              <Tag color={detailEnv.color || 'default'}>{detailEnv.name.toUpperCase()}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="颜色">
              <Space>
                <div
                  style={{
                    width: 20,
                    height: 20,
                    backgroundColor: detailEnv.color || '#1890ff',
                    borderRadius: 4,
                    border: '1px solid #d9d9d9',
                  }}
                />
                <span>{detailEnv.color || '#1890ff'}</span>
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="描述">{detailEnv.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailEnv.created_at ? new Date(detailEnv.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailEnv.updated_at ? new Date(detailEnv.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default Environments;
