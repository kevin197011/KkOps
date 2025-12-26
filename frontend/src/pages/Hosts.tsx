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
  ColorPicker,
  Descriptions,
  Tooltip,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  SyncOutlined,
  LinkOutlined,
} from '@ant-design/icons';
import { hostService, Host, CreateHostRequest } from '../services/host';
import { projectService, Project } from '../services/project';
import { environmentService, Environment } from '../services/environment';
import { cloudPlatformService, CloudPlatform } from '../services/cloudPlatform';
import { hostGroupService, HostGroup } from '../services/hostGroup';
import { hostTagService, HostTag } from '../services/hostTag';
import { sshKeyService, SSHKey } from '../services/sshKey';

const { Option } = Select;

const Hosts: React.FC = () => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [cloudPlatforms, setCloudPlatforms] = useState<CloudPlatform[]>([]);
  const [hostGroups, setHostGroups] = useState<HostGroup[]>([]);
  const [hostTags, setHostTags] = useState<HostTag[]>([]);
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([]);
  const [loading, setLoading] = useState(false);
  const [syncingHostId, setSyncingHostId] = useState<number | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingHost, setEditingHost] = useState<Host | null>(null);
  const [detailHost, setDetailHost] = useState<Host | null>(null);
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined);
  const [selectedEnvironment, setSelectedEnvironment] = useState<string | undefined>(undefined);
  const [selectedCloudPlatformId, setSelectedCloudPlatformId] = useState<number | undefined>(undefined);
  const [selectedGroupId, setSelectedGroupId] = useState<number | undefined>(undefined);
  const [selectedTagId, setSelectedTagId] = useState<number | undefined>(undefined);
  const [tagModalVisible, setTagModalVisible] = useState(false);
  const [tagForm] = Form.useForm();
  const [form] = Form.useForm();

  useEffect(() => {
    loadProjects();
    loadEnvironments();
    loadCloudPlatforms();
    loadHostGroups();
    loadHostTags();
    loadSSHKeys();
  }, []);

  useEffect(() => {
    loadHosts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, selectedProjectId, selectedEnvironment, selectedCloudPlatformId, selectedGroupId, selectedTagId]);

  const loadProjects = async () => {
    try {
      const response = await projectService.list(1, 100);
      setProjects(response.projects);
    } catch (error: any) {
      message.error('加载项目列表失败');
    }
  };

  const loadEnvironments = async () => {
    try {
      const response = await environmentService.list(1, 100);
      setEnvironments(response.environments);
    } catch (error: any) {
      // 静默失败，不影响主流程
    }
  };

  const loadCloudPlatforms = async () => {
    try {
      const response = await cloudPlatformService.list(1, 100);
      setCloudPlatforms(response.cloud_platforms || []);
    } catch (error: any) {
      // 静默失败，不影响主流程
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

  const handleCreateTag = async () => {
    try {
      const values = await tagForm.validateFields();
      const color = typeof values.color === 'string' ? values.color : values.color?.toHexString();
      const response = await hostTagService.create({
        name: values.name,
        color: color,
        description: values.description,
      });
      message.success('标签创建成功');
      setTagModalVisible(false);
      tagForm.resetFields();
      // 重新加载标签列表
      await loadHostTags();
      // 如果正在编辑主机，自动选中新创建的标签
      if (editingHost) {
        const currentTagIds = form.getFieldValue('tag_ids') || [];
        form.setFieldsValue({
          tag_ids: [...currentTagIds, response.tag.id],
        });
      }
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '创建标签失败');
    }
  };

  const loadSSHKeys = async () => {
    try {
      const response = await sshKeyService.list(1, 100);
      setSshKeys(response.ssh_keys);
    } catch (error: any) {
      // 静默失败，不影响主流程
    }
  };

  const loadHosts = async () => {
    setLoading(true);
    try {
      const filters: any = {};
      if (selectedProjectId) filters.project_id = selectedProjectId;
      if (selectedEnvironment) filters.environment = selectedEnvironment;
      if (selectedCloudPlatformId) filters.cloud_platform_id = selectedCloudPlatformId;
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

  const handleViewDetail = async (host: Host) => {
    try {
      const response = await hostService.get(host.id);
      setDetailHost(response.host);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取主机详情失败');
    }
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
        cloud_platform_id: fullHost.cloud_platform_id || undefined,
        group_ids: fullHost.groups?.map((g) => g.id) || [],
        tag_ids: fullHost.tags?.map((t) => t.id) || [],
        ssh_key_id: fullHost.ssh_key_id || undefined,
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

  // 同步单个主机状态（通过 Salt）
  const handleSyncStatus = async (host: Host) => {
    if (!host.salt_minion_id) {
      message.warning('该主机未配置 Salt Minion ID，无法同步状态');
      return;
    }
    setSyncingHostId(host.id);
    try {
      const response = await hostService.syncStatus(host.id);
      message.success(`状态同步成功：${response.host.status}`);
      loadHosts();
    } catch (error: any) {
      message.error(error.response?.data?.error || '同步状态失败');
    } finally {
      setSyncingHostId(null);
    }
  };

  // 同步单个主机（包含自动匹配和状态同步）
  const handleSyncHost = async (host: Host) => {
    setSyncingHostId(host.id);
    try {
      await hostService.syncStatus(host.id);
      message.success('主机同步完成');
      loadHosts();
    } catch (error: any) {
      message.error(error.response?.data?.error || '同步失败');
    } finally {
      setSyncingHostId(null);
    }
  };

  // 同步所有主机状态
  const handleSyncAllStatus = async () => {
    setLoading(true);
    try {
      await hostService.syncAllStatus();
      message.success('所有主机状态同步完成');
      loadHosts();
    } catch (error: any) {
      message.error(error.response?.data?.error || '同步失败');
    } finally {
      setLoading(false);
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
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '项目',
      key: 'project',
      width: 120,
      render: (_: any, record: Host) => {
        if (record.project) {
          return <span>{record.project.name}</span>;
        }
        // 如果没有关联数据，尝试从 projects 列表查找
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
      title: '环境',
      dataIndex: 'environment',
      key: 'environment',
      width: 90,
      render: (env: string) => {
        if (!env) {
          return <span style={{ color: '#999' }}>-</span>;
        }
        // 从 environments 状态中查找对应的环境配置
        const envConfig = environments.find((e) => e.name === env);
        if (envConfig) {
          return <Tag color={envConfig.color || 'default'}>{envConfig.display_name || env.toUpperCase()}</Tag>;
        }
        // 如果没找到，使用默认颜色
        const defaultColors: Record<string, string> = {
          dev: 'blue',
          uat: 'orange',
          staging: 'purple',
          prod: 'red',
          test: 'green',
        };
        return <Tag color={defaultColors[env] || 'default'}>{env.toUpperCase()}</Tag>;
      },
    },
    {
      title: '云平台',
      key: 'cloud_platform',
      width: 100,
      render: (_: any, record: Host) => {
        if (record.cloud_platform) {
          return <Tag color={record.cloud_platform.color || 'default'}>{record.cloud_platform.display_name}</Tag>;
        }
        if (record.cloud_platform_id) {
          const cp = cloudPlatforms.find((p) => p.id === record.cloud_platform_id);
          if (cp) {
            return <Tag color={cp.color || 'default'}>{cp.display_name}</Tag>;
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
      title: '主机组',
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
      width: 280,
      render: (_: any, record: Host) => (
        <Space wrap>
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
          {record.salt_minion_id ? (
            <Tooltip title="同步 Salt 状态">
              <Button
                type="link"
                size="small"
                icon={<SyncOutlined spin={syncingHostId === record.id} />}
                onClick={() => handleSyncStatus(record)}
                loading={syncingHostId === record.id}
              >
                同步
              </Button>
            </Tooltip>
          ) : (
            <Tooltip title="同步主机（自动匹配 Salt Minion 并更新状态信息）">
              <Button
                type="link"
                size="small"
                icon={<SyncOutlined />}
                onClick={() => handleSyncHost(record)}
                loading={syncingHostId === record.id}
              >
                同步
              </Button>
            </Tooltip>
          )}
          <Popconfirm
            title="确定要删除这台主机吗？"
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
              placeholder="选择环境"
              allowClear
              style={{ width: 150 }}
              value={selectedEnvironment}
              onChange={(value) => {
                setSelectedEnvironment(value);
                setPage(1);
              }}
            >
              {environments.map((env) => (
                <Option key={env.id} value={env.name}>
                  <Tag color={env.color || 'default'}>{env.name.toUpperCase()}</Tag>
                </Option>
              ))}
            </Select>
            <Select
              placeholder="选择云平台"
              allowClear
              style={{ width: 150 }}
              value={selectedCloudPlatformId}
              onChange={(value) => {
                setSelectedCloudPlatformId(value);
                setPage(1);
              }}
            >
              {cloudPlatforms.map((cp) => (
                <Option key={cp.id} value={cp.id}>
                  <Tag color={cp.color || 'default'}>{cp.display_name}</Tag>
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
            <Tooltip title="并发同步所有主机状态和信息（最多100个并发）">
              <Button icon={<SyncOutlined />} onClick={handleSyncAllStatus} loading={loading}>
                全部同步
              </Button>
            </Tooltip>
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
            name="environment"
            label="环境"
            tooltip="选择主机所属环境（可选）"
          >
            <Select placeholder="请选择环境（可选）" allowClear>
              {environments.map((env) => (
                <Option key={env.id} value={env.name}>
                  <Tag color={env.color || 'default'}>{env.name.toUpperCase()}</Tag>
                  <span style={{ marginLeft: 8 }}>{env.display_name}</span>
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="cloud_platform_id"
            label="云平台"
            tooltip="选择主机所属云平台（可选）"
          >
            <Select placeholder="请选择云平台（可选）" allowClear>
              {cloudPlatforms.map((cp) => (
                <Option key={cp.id} value={cp.id}>
                  <Tag color={cp.color || 'default'}>{cp.display_name}</Tag>
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
              {
                validator: async (rule, value) => {
                  if (!value) return;
                  try {
                    // 检查IP地址是否已被使用
                    const response = await hostService.list(1, 1, { ip_address: value });
                    if (response.hosts.length > 0) {
                      // 如果是编辑模式，允许使用当前主机的IP地址
                      if (editingHost && response.hosts[0].id === editingHost.id) {
                        return;
                      }
                      throw new Error(`IP地址 ${value} 已被主机 "${response.hosts[0].hostname}" (ID: ${response.hosts[0].id}) 使用`);
                    }
                  } catch (error: any) {
                    if (error.message.includes('已被')) {
                      throw new Error(error.message);
                    }
                    // 忽略其他错误（如网络错误），让后端处理
                  }
                },
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
            name="ssh_key_id"
            label="默认SSH密钥"
            tooltip="选择默认使用的SSH密钥（可选，打开终端时会自动选择）"
          >
            <Select placeholder="选择SSH密钥（可选）" allowClear>
              {sshKeys.map((key) => (
                <Option key={key.id} value={key.id}>
                  {key.name} ({key.key_type.toUpperCase()})
                </Option>
              ))}
            </Select>
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
            extra="可以在主机标签管理页面创建和管理标签"
          >
            <Select
              mode="multiple"
              placeholder="选择标签"
              allowClear
              dropdownRender={(menu) => (
                <>
                  {menu}
                  <div style={{ padding: '8px', borderTop: '1px solid #f0f0f0' }}>
                    <Button
                      type="link"
                      icon={<PlusOutlined />}
                      onClick={() => setTagModalVisible(true)}
                      style={{ width: '100%' }}
                    >
                      创建新标签
                    </Button>
                  </div>
                </>
              )}
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

      {/* 快速创建标签模态框 */}
      <Modal
        title="创建新标签"
        open={tagModalVisible}
        onOk={handleCreateTag}
        onCancel={() => {
          setTagModalVisible(false);
          tagForm.resetFields();
        }}
        okText="确定"
        cancelText="取消"
      >
        <Form form={tagForm} layout="vertical">
          <Form.Item
            name="name"
            label="标签名称"
            rules={[{ required: true, message: '请输入标签名称' }]}
          >
            <Input placeholder="请输入标签名称" />
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
              placeholder="请输入标签描述（可选）"
              rows={3}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 主机详情模态框 */}
      <Modal
        title="主机详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={700}
      >
        {detailHost && (
          <Descriptions bordered column={2} size="small">
            <Descriptions.Item label="主机名">{detailHost.hostname}</Descriptions.Item>
            <Descriptions.Item label="IP地址">{detailHost.ip_address}</Descriptions.Item>
            <Descriptions.Item label="SSH端口">{detailHost.ssh_port || 22}</Descriptions.Item>
            <Descriptions.Item label="状态">{getStatusTag(detailHost.status)}</Descriptions.Item>
            <Descriptions.Item label="项目">{detailHost.project?.name || '-'}</Descriptions.Item>
            <Descriptions.Item label="环境">
              {detailHost.environment ? (() => {
                const envConfig = environments.find((e) => e.name === detailHost.environment);
                if (envConfig) {
                  return <Tag color={envConfig.color || 'default'}>{envConfig.display_name || detailHost.environment.toUpperCase()}</Tag>;
                }
                const defaultColors: Record<string, string> = { dev: 'blue', test: 'green', uat: 'orange', staging: 'purple', prod: 'red' };
                return <Tag color={defaultColors[detailHost.environment] || 'default'}>{detailHost.environment.toUpperCase()}</Tag>;
              })() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="云平台">
              {detailHost.cloud_platform ? (
                <Tag color={detailHost.cloud_platform.color || 'default'}>{detailHost.cloud_platform.display_name}</Tag>
              ) : detailHost.cloud_platform_id ? (() => {
                const cp = cloudPlatforms.find((p) => p.id === detailHost.cloud_platform_id);
                return cp ? <Tag color={cp.color || 'default'}>{cp.display_name}</Tag> : '-';
              })() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="操作系统" span={2}>
              {detailHost.os_type ? `${detailHost.os_type}${detailHost.os_version ? ` ${detailHost.os_version}` : ''}` : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="CPU核心数">{detailHost.cpu_cores || '-'}</Descriptions.Item>
            <Descriptions.Item label="内存">{detailHost.memory_gb ? `${detailHost.memory_gb} GB` : '-'}</Descriptions.Item>
            <Descriptions.Item label="磁盘">{detailHost.disk_gb ? `${detailHost.disk_gb} GB` : '-'}</Descriptions.Item>
            <Descriptions.Item label="Salt Minion ID">{detailHost.salt_minion_id || '-'}</Descriptions.Item>
            <Descriptions.Item label="主机组" span={2}>
              <Space wrap>
                {detailHost.groups && detailHost.groups.length > 0 ? (
                  detailHost.groups.map((group) => <Tag key={group.id}>{group.name}</Tag>)
                ) : '-'}
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="标签" span={2}>
              <Space wrap>
                {detailHost.tags && detailHost.tags.length > 0 ? (
                  detailHost.tags.map((tag) => (
                    <Tag key={tag.id} color={tag.color || '#1890ff'}>{tag.name}</Tag>
                  ))
                ) : '-'}
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="描述" span={2}>{detailHost.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="最后在线">
              {detailHost.last_seen_at ? new Date(detailHost.last_seen_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailHost.created_at ? new Date(detailHost.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default Hosts;

