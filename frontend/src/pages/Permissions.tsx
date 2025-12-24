import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Input,
  Tag,
  message,
  Select,
} from 'antd';
import { ReloadOutlined, SearchOutlined } from '@ant-design/icons';
import { permissionService, Permission } from '../services/permission';

const { Option } = Select;

const Permissions: React.FC = () => {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  const [resourceTypeFilter, setResourceTypeFilter] = useState<string>('');

  useEffect(() => {
    loadPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, resourceTypeFilter]);

  const loadPermissions = async () => {
    setLoading(true);
    try {
      const response = await permissionService.list(page, pageSize);
      let filteredPermissions = response.permissions;

      // 应用资源类型过滤
      if (resourceTypeFilter) {
        filteredPermissions = filteredPermissions.filter(
          (p) => p.resource_type === resourceTypeFilter
        );
      }

      // 应用搜索过滤
      if (searchText) {
        filteredPermissions = filteredPermissions.filter(
          (p) =>
            p.name.toLowerCase().includes(searchText.toLowerCase()) ||
            p.code.toLowerCase().includes(searchText.toLowerCase()) ||
            p.description?.toLowerCase().includes(searchText.toLowerCase())
        );
      }

      setPermissions(filteredPermissions);
      setTotal(filteredPermissions.length);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载权限列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取所有唯一的资源类型
  const getResourceTypes = () => {
    const types = new Set<string>();
    permissions.forEach((p) => types.add(p.resource_type));
    return Array.from(types);
  };

  const getActionTag = (action: string) => {
    const actionMap: Record<string, { color: string; text: string }> = {
      create: { color: 'green', text: '创建' },
      read: { color: 'blue', text: '读取' },
      update: { color: 'orange', text: '更新' },
      delete: { color: 'red', text: '删除' },
      execute: { color: 'purple', text: '执行' },
      manage: { color: 'cyan', text: '管理' },
    };
    const actionInfo = actionMap[action] || { color: 'default', text: action };
    return <Tag color={actionInfo.color}>{actionInfo.text}</Tag>;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '权限代码',
      dataIndex: 'code',
      key: 'code',
      width: 200,
    },
    {
      title: '权限名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '资源类型',
      dataIndex: 'resource_type',
      key: 'resource_type',
      width: 150,
      render: (text: string) => <Tag>{text}</Tag>,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      render: (action: string) => getActionTag(action),
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
  ];

  return (
    <Card>
      <div style={{ marginBottom: '16px' }}>
        <h2 style={{ margin: 0, marginBottom: '16px' }}>权限管理</h2>

        <Space wrap style={{ marginBottom: '16px' }}>
          <Input
            placeholder="搜索权限名称、代码或描述"
            style={{ width: 300 }}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            allowClear
            onPressEnter={loadPermissions}
          />
          <Select
            placeholder="资源类型"
            style={{ width: 150 }}
            value={resourceTypeFilter}
            onChange={(value) => setResourceTypeFilter(value)}
            allowClear
          >
            {getResourceTypes().map((type) => (
              <Option key={type} value={type}>
                {type}
              </Option>
            ))}
          </Select>
          <Button
            type="primary"
            icon={<SearchOutlined />}
            onClick={loadPermissions}
          >
            搜索
          </Button>
          <Button icon={<ReloadOutlined />} onClick={loadPermissions}>
            刷新
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={permissions}
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
  );
};

export default Permissions;

