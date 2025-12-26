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
  Descriptions,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { hostGroupService, HostGroup } from '../services/hostGroup';

const HostGroups: React.FC = () => {
  const [groups, setGroups] = useState<HostGroup[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingGroup, setEditingGroup] = useState<HostGroup | null>(null);
  const [detailGroup, setDetailGroup] = useState<HostGroup | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadGroups();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const loadGroups = async () => {
    setLoading(true);
    try {
      const response = await hostGroupService.list(page, pageSize);
      setGroups(response.groups);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载主机组列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingGroup(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleViewDetail = async (group: HostGroup) => {
    try {
      const response = await hostGroupService.get(group.id);
      setDetailGroup(response.group);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取主机组详情失败');
    }
  };

  const handleEdit = async (group: HostGroup) => {
    try {
      const response = await hostGroupService.get(group.id);
      setEditingGroup(response.group);
      form.setFieldsValue({
        name: response.group.name,
        description: response.group.description,
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取主机组详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await hostGroupService.delete(id);
      message.success('删除成功');
      loadGroups();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingGroup) {
        await hostGroupService.update(editingGroup.id, {
          name: values.name,
          description: values.description,
        });
        message.success('更新成功');
      } else {
        await hostGroupService.create({
          name: values.name,
          description: values.description,
        });
        message.success('创建成功');
      }
      setModalVisible(false);
      loadGroups();
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
      width: 60,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 150,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '主机数',
      key: 'host_count',
      width: 80,
      render: (_: any, record: HostGroup) => (
        <span>{record.hosts?.length || 0}</span>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_: any, record: HostGroup) => (
        <Space size={0}>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          />
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          />
          <Popconfirm
            title="确定要删除这个主机组吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between' }}>
          <h2 style={{ margin: 0 }}>主机组管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={loadGroups}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              创建主机组
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={groups}
          rowKey="id"
          loading={loading}
          scroll={{ x: 'max-content' }}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            pageSizeOptions: ['10', '20', '50', '100'],
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: (page, pageSize) => {
              setPage(page);
              setPageSize(pageSize);
            },
          }}
        />
      </Card>

      {/* 创建/编辑主机组模态框 */}
      <Modal
        title={editingGroup ? '编辑主机组' : '创建主机组'}
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
            rules={[{ required: true, message: '请输入主机组名称' }]}
          >
            <Input placeholder="如：运维组、数据库组、Web服务组" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea
              placeholder="请输入主机组描述"
              rows={4}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 主机组详情模态框 */}
      <Modal
        title="主机组详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={<Button onClick={() => setDetailModalVisible(false)}>关闭</Button>}
        width={500}
      >
        {detailGroup && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="名称">{detailGroup.name}</Descriptions.Item>
            <Descriptions.Item label="描述">{detailGroup.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="主机数量">{detailGroup.hosts?.length || 0}</Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailGroup.created_at ? new Date(detailGroup.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailGroup.updated_at ? new Date(detailGroup.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default HostGroups;
