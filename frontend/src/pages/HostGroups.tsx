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
  Descriptions,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { hostGroupService, HostGroup, CreateHostGroupRequest } from '../services/hostGroup';
import { projectService, Project } from '../services/project';

const { Option } = Select;

const HostGroups: React.FC = () => {
  const [groups, setGroups] = useState<HostGroup[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingGroup, setEditingGroup] = useState<HostGroup | null>(null);
  const [detailGroup, setDetailGroup] = useState<HostGroup | null>(null);
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined);
  const [form] = Form.useForm();

  useEffect(() => {
    loadProjects();
  }, []);

  useEffect(() => {
    loadGroups();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, selectedProjectId]);

  const loadProjects = async () => {
    try {
      const response = await projectService.list(1, 100);
      setProjects(response.projects);
    } catch (error: any) {
      message.error('加载项目列表失败');
    }
  };

  const loadGroups = async () => {
    setLoading(true);
    try {
      const response = await hostGroupService.list(page, pageSize, {
        project_id: selectedProjectId,
      });
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
        project_id: response.group.project_id,
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
          project_id: values.project_id,
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
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '项目',
      key: 'project',
      width: 150,
      render: (_: any, record: HostGroup) => (
        <span>{record.project?.name || '-'}</span>
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
      title: '主机数量',
      key: 'host_count',
      width: 100,
      render: (_: any, record: HostGroup) => (
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
      width: 200,
      fixed: 'right' as const,
      render: (_: any, record: HostGroup) => (
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
          <h2 style={{ margin: 0 }}>主机组管理</h2>
          <Space>
            <Select
              placeholder="筛选项目"
              style={{ width: 200 }}
              value={selectedProjectId}
              onChange={(value) => setSelectedProjectId(value)}
              allowClear
            >
              {projects.map((project) => (
                <Option key={project.id} value={project.id}>
                  {project.name}
                </Option>
              ))}
            </Select>
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
          {!editingGroup && (
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
          )}
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入主机组名称' }]}
          >
            <Input placeholder="请输入主机组名称" />
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
        footer={null}
        width={600}
      >
        {detailGroup && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="名称">{detailGroup.name}</Descriptions.Item>
            <Descriptions.Item label="项目">{detailGroup.project?.name || '-'}</Descriptions.Item>
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

