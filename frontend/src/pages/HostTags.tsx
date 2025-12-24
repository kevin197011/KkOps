import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Tag,
  ColorPicker,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { hostTagService, HostTag, CreateHostTagRequest } from '../services/hostTag';

const HostTags: React.FC = () => {
  const [tags, setTags] = useState<HostTag[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingTag, setEditingTag] = useState<HostTag | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadTags();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const loadTags = async () => {
    setLoading(true);
    try {
      const response = await hostTagService.list(page, pageSize);
      setTags(response.tags);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载主机标签列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingTag(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = async (tag: HostTag) => {
    try {
      const response = await hostTagService.get(tag.id);
      setEditingTag(response.tag);
      form.setFieldsValue({
        name: response.tag.name,
        color: response.tag.color || '#1890ff',
        description: response.tag.description,
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取主机标签详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await hostTagService.delete(id);
      message.success('删除成功');
      loadTags();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const color = typeof values.color === 'string' ? values.color : values.color?.toHexString();
      if (editingTag) {
        await hostTagService.update(editingTag.id, {
          name: values.name,
          color: color,
          description: values.description,
        });
        message.success('更新成功');
      } else {
        await hostTagService.create({
          name: values.name,
          color: color,
          description: values.description,
        });
        message.success('创建成功');
      }
      setModalVisible(false);
      loadTags();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '标签',
      key: 'tag',
      width: 200,
      render: (_: any, record: HostTag) => (
        <Tag color={record.color || '#1890ff'}>
          {record.name}
        </Tag>
      ),
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '颜色',
      dataIndex: 'color',
      key: 'color',
      width: 100,
      render: (color: string) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <div
            style={{
              width: 20,
              height: 20,
              backgroundColor: color || '#1890ff',
              borderRadius: 4,
              border: '1px solid #d9d9d9',
            }}
          />
          <span>{color || '#1890ff'}</span>
        </div>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '使用数量',
      key: 'host_count',
      width: 100,
      render: (_: any, record: HostTag) => (
        <span>{record.hosts?.length || 0}</span>
      ),
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
      width: 150,
      fixed: 'right' as const,
      render: (_: any, record: HostTag) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个主机标签吗？"
            onConfirm={() => handleDelete(record.id)}
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
  ];

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between' }}>
          <h2 style={{ margin: 0 }}>主机标签管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={loadTags}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              创建主机标签
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={tags}
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
      </Card>

      {/* 创建/编辑主机标签模态框 */}
      <Modal
        title={editingTag ? '编辑主机标签' : '创建主机标签'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入主机标签名称' }]}
          >
            <Input placeholder="请输入主机标签名称" />
          </Form.Item>
          <Form.Item
            name="color"
            label="颜色"
            initialValue="#1890ff"
          >
            <ColorPicker showText />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea
              placeholder="请输入主机标签描述"
              rows={4}
            />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default HostTags;

