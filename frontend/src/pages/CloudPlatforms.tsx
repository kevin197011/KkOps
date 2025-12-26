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
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CloudServerOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  CloudPlatform,
  cloudPlatformService,
} from '../services/cloudPlatform';

const CloudPlatforms: React.FC = () => {
  const [platforms, setPlatforms] = useState<CloudPlatform[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingPlatform, setEditingPlatform] = useState<CloudPlatform | null>(null);
  const [viewingPlatform, setViewingPlatform] = useState<CloudPlatform | null>(null);
  const [form] = Form.useForm();

  const fetchPlatforms = async () => {
    setLoading(true);
    try {
      const response = await cloudPlatformService.list();
      setPlatforms(response.cloud_platforms || []);
    } catch (error) {
      message.error('获取云平台列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlatforms();
  }, []);

  const handleCreate = () => {
    setEditingPlatform(null);
    form.resetFields();
    form.setFieldsValue({ color: '#1890ff', sort_order: 0 });
    setModalVisible(true);
  };

  const handleEdit = (record: CloudPlatform) => {
    setEditingPlatform(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleView = (record: CloudPlatform) => {
    setViewingPlatform(record);
    setDetailModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await cloudPlatformService.delete(id);
      message.success('删除成功');
      fetchPlatforms();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async (values: any) => {
    try {
      // 处理颜色值
      if (values.color && typeof values.color === 'object') {
        values.color = values.color.toHexString();
      }

      if (editingPlatform) {
        await cloudPlatformService.update(editingPlatform.id, values);
        message.success('更新成功');
      } else {
        await cloudPlatformService.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchPlatforms();
    } catch (error) {
      message.error(editingPlatform ? '更新失败' : '创建失败');
    }
  };

  const columns: ColumnsType<CloudPlatform> = [
    {
      title: '序号',
      key: 'index',
      width: 60,
      render: (_: any, __: any, index: number) => index + 1,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 100,
    },
    {
      title: '显示名称',
      dataIndex: 'display_name',
      key: 'display_name',
      width: 120,
      render: (text, record) => (
        record.color ? <Tag color={record.color}>{text}</Tag> : text
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="text"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          />
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          />
          <Popconfirm
            title="确定要删除这个云平台吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="text" size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title={
        <Space>
          <CloudServerOutlined />
          云平台管理
        </Space>
      }
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新建云平台
        </Button>
      }
    >
      <Table
        columns={columns}
        dataSource={platforms}
        rowKey="id"
        loading={loading}
        pagination={false}
        scroll={{ x: 'max-content' }}
      />

      <Modal
        title={editingPlatform ? '编辑云平台' : '新建云平台'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入名称' }]}
          >
            <Input placeholder="如: aws, aliyun, azure" />
          </Form.Item>
          <Form.Item
            name="display_name"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          >
            <Input placeholder="如: AWS, 阿里云, Azure" />
          </Form.Item>
          <Form.Item name="color" label="颜色">
            <ColorPicker showText />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="云平台描述（可选）" />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingPlatform ? '更新' : '创建'}
              </Button>
              <Button onClick={() => setModalVisible(false)}>取消</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="云平台详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
        ]}
      >
        {viewingPlatform && (
          <div>
            <p><strong>名称：</strong>{viewingPlatform.name}</p>
            <p><strong>显示名称：</strong>
              {viewingPlatform.color ? (
                <Tag color={viewingPlatform.color}>{viewingPlatform.display_name}</Tag>
              ) : (
                viewingPlatform.display_name
              )}
            </p>
            <p><strong>颜色：</strong>
              {viewingPlatform.color ? (
                <Tag color={viewingPlatform.color}>{viewingPlatform.color}</Tag>
              ) : '-'}
            </p>
            <p><strong>描述：</strong>{viewingPlatform.description || '-'}</p>
            <p><strong>创建时间：</strong>{new Date(viewingPlatform.created_at).toLocaleString()}</p>
            <p><strong>更新时间：</strong>{new Date(viewingPlatform.updated_at).toLocaleString()}</p>
          </div>
        )}
      </Modal>
    </Card>
  );
};

export default CloudPlatforms;
