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
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { hostService, Host, CreateHostRequest } from '../services/host';
import { projectService, Project } from '../services/project';
import { hostGroupService, HostGroup } from '../services/hostGroup';
import { hostTagService, HostTag } from '../services/hostTag';

const { Option } = Select;

const Hosts: React.FC = () => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [hostGroups, setHostGroups] = useState<HostGroup[]>([]);
  const [hostTags, setHostTags] = useState<HostTag[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingHost, setEditingHost] = useState<Host | null>(null);
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined);
  const [selectedGroupId, setSelectedGroupId] = useState<number | undefined>(undefined);
  const [selectedTagId, setSelectedTagId] = useState<number | undefined>(undefined);
  const [form] = Form.useForm();

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
      if (selectedProjectId) filters.project_id = selectedProjectId;
      if (selectedGroupId) filters.group_id = selectedGroupId;
      if (selectedTagId) filters.tag_id = selectedTagId;
      const response = await hostService.list(page, pageSize, filters);
      setHosts(response.hosts);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载主机列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingHost(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = async (host: Host) => {
    try {
      // 获取完整的主机信息，包括project_id、groups、tags
      const response = await hostService.get(host.id);
      const fullHost = response.host;
      setEditingHost(fullHost);
      form.setFieldsValue({
        ...fullHost,
        project_id: (fullHost as any).project_id,
        group_ids: fullHost.groups?.map((g) => g.id) || [],
        tag_ids: fullHost.tags?.map((t) => t.id) || [],
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取主机详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await hostService.delete(id);
      message.success('删除成功');
      loadHosts();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const { group_ids, tag_ids, ...hostData } = values;
      
      if (editingHost) {
        await hostService.update(editingHost.id, hostData);
        
        // 更新组和标签关联
        const currentGroupIds = editingHost.groups?.map((g) => g.id) || [];
        const currentTagIds = editingHost.tags?.map((t) => t.id) || [];
        
        // 添加新的组
        const groupsToAdd = (group_ids || []).filter((id: number) => !currentGroupIds.includes(id));
        for (const groupId of groupsToAdd) {
          await hostService.addToGroup(editingHost.id, groupId);
        }
        // 移除旧的组
        const groupsToRemove = currentGroupIds.filter((id: number) => !(group_ids || []).includes(id));
        for (const groupId of groupsToRemove) {
          await hostService.removeFromGroup(editingHost.id, groupId);
        }
        
        // 添加新的标签
        const tagsToAdd = (tag_ids || []).filter((id: number) => !currentTagIds.includes(id));
        for (const tagId of tagsToAdd) {
          await hostService.addTag(editingHost.id, tagId);
        }
        // 移除旧的标签
        const tagsToRemove = currentTagIds.filter((id: number) => !(tag_ids || []).includes(id));
        for (const tagId of tagsToRemove) {
          await hostService.removeTag(editingHost.id, tagId);
        }
        
        message.success('更新成功');
      } else {
        const response = await hostService.create(hostData as CreateHostRequest);
        const newHostId = response.host.id;
        
        // 添加组和标签关联
        if (group_ids && group_ids.length > 0) {
          for (const groupId of group_ids) {
            await hostService.addToGroup(newHostId, groupId);
          }
        }
        if (tag_ids && tag_ids.length > 0) {
          for (const tagId of tag_ids) {
            await hostService.addTag(newHostId, tagId);
          }
        }
        
        message.success('创建成功');
      }
      setModalVisible(false);
      loadHosts();
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
      title: 'Salt Minion ID',
      dataIndex: 'salt_minion_id',
      key: 'salt_minion_id',
      render: (text: string) => text || '-',
    },
    {
      title: '操作系统',
      key: 'os',
      render: (_: any, record: Host) => {
        if (record.os_type) {
          return `${record.os_type}${record.os_version ? ` ${record.os_version}` : ''}`;
        }
        return '-';
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
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
      title: '最后在线',
      dataIndex: 'last_seen_at',
      key: 'last_seen_at',
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Host) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这台主机吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
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
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2 style={{ margin: 0 }}>主机管理</h2>
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
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              新建主机
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
      </Card>

      <Modal
        title={editingHost ? '编辑主机' : '新建主机'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={form} layout="vertical">
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
            name="hostname"
            label="主机名"
            rules={[{ required: true, message: '请输入主机名' }]}
          >
            <Input placeholder="请输入主机名" />
          </Form.Item>
          <Form.Item
            name="ip_address"
            label="IP地址"
            rules={[
              { required: true, message: '请输入IP地址' },
              {
                pattern: /^(\d{1,3}\.){3}\d{1,3}$/,
                message: '请输入有效的IP地址',
              },
            ]}
          >
            <Input placeholder="例如: 192.168.1.100" />
          </Form.Item>
          <Form.Item
            name="ssh_port"
            label="SSH端口"
            initialValue={22}
            rules={[
              { required: true, message: '请输入SSH端口' },
              { type: 'number', min: 1, max: 65535, message: '端口范围: 1-65535' },
            ]}
          >
            <InputNumber min={1} max={65535} placeholder="SSH端口（默认: 22）" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="salt_minion_id"
            label="Salt Minion ID"
          >
            <Input placeholder="Salt Minion ID（可选）" />
          </Form.Item>
          <Form.Item
            name="os_type"
            label="操作系统类型"
          >
            <Select placeholder="请选择操作系统类型" allowClear>
              <Option value="Linux">Linux</Option>
              <Option value="Windows">Windows</Option>
              <Option value="macOS">macOS</Option>
              <Option value="Unix">Unix</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="os_version"
            label="操作系统版本"
          >
            <Input placeholder="例如: Ubuntu 20.04" />
          </Form.Item>
          <Form.Item
            name="cpu_cores"
            label="CPU核心数"
          >
            <InputNumber min={1} placeholder="CPU核心数" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="memory_gb"
            label="内存(GB)"
          >
            <InputNumber min={0} placeholder="内存大小" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="disk_gb"
            label="磁盘(GB)"
          >
            <InputNumber min={0} placeholder="磁盘大小" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="status"
            label="状态"
            initialValue="unknown"
          >
            <Select>
              <Option value="online">在线</Option>
              <Option value="offline">离线</Option>
              <Option value="unknown">未知</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea rows={3} placeholder="主机描述信息" />
          </Form.Item>
          <Form.Item
            name="group_ids"
            label="主机组"
          >
            <Select
              mode="multiple"
              placeholder="选择主机组"
              allowClear
            >
              {hostGroups
                .filter((group) => !editingHost || group.project_id === (editingHost as any).project_id)
                .map((group) => (
                  <Option key={group.id} value={group.id}>
                    {group.name}
                  </Option>
                ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="tag_ids"
            label="标签"
          >
            <Select
              mode="multiple"
              placeholder="选择标签"
              allowClear
            >
              {hostTags.map((tag) => (
                <Option key={tag.id} value={tag.id}>
                  <Tag color={tag.color || '#1890ff'}>{tag.name}</Tag>
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default Hosts;

