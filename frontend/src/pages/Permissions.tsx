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
  Modal,
  Descriptions,
} from 'antd';
import { ReloadOutlined, SearchOutlined, EyeOutlined } from '@ant-design/icons';
import { permissionService, Permission } from '../services/permission';

const { Option } = Select;

const Permissions: React.FC = () => {
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]); // 所有权限数据
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  const [resourceTypeFilter, setResourceTypeFilter] = useState<string>('');
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [detailPermission, setDetailPermission] = useState<Permission | null>(null);

  useEffect(() => {
    loadPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadPermissions = async () => {
    setLoading(true);
    try {
      // 获取所有权限数据，在前端做分页和过滤
      const response = await permissionService.listAll();
      setAllPermissions(response.permissions);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载权限列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取所有唯一的资源类型
  const getResourceTypes = () => {
    const types = new Set<string>();
    allPermissions.forEach((p) => types.add(p.resource_type));
    return Array.from(types);
  };

  // 过滤后的权限数据
  const filteredPermissions = React.useMemo(() => {
    let result = allPermissions;

    // 应用资源类型过滤
    if (resourceTypeFilter) {
      result = result.filter((p) => p.resource_type === resourceTypeFilter);
    }

    // 应用搜索过滤
    if (searchText) {
      result = result.filter(
        (p) =>
          p.name.toLowerCase().includes(searchText.toLowerCase()) ||
          p.code.toLowerCase().includes(searchText.toLowerCase()) ||
          p.description?.toLowerCase().includes(searchText.toLowerCase())
      );
    }

    return result;
  }, [allPermissions, resourceTypeFilter, searchText]);

  // 当前页的数据
  const currentPageData = React.useMemo(() => {
    const startIndex = (page - 1) * pageSize;
    return filteredPermissions.slice(startIndex, startIndex + pageSize);
  }, [filteredPermissions, page, pageSize]);

  // 当过滤条件变化时，重置到第一页
  useEffect(() => {
    setPage(1);
  }, [searchText, resourceTypeFilter]);

  const handleViewDetail = (permission: Permission) => {
    setDetailPermission(permission);
    setDetailModalVisible(true);
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
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
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
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right' as const,
      render: (_: any, record: Permission) => (
        <Button
          type="link"
          size="small"
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record)}
        >
          详情
        </Button>
      ),
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
        dataSource={currentPageData}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: filteredPermissions.length,
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
          onChange: (newPage, newPageSize) => {
            setPage(newPage);
            setPageSize(newPageSize);
          },
        }}
      />

      {/* 权限详情模态框 */}
      <Modal
        title="权限详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {detailPermission && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="权限代码">{detailPermission.code}</Descriptions.Item>
            <Descriptions.Item label="权限名称">{detailPermission.name}</Descriptions.Item>
            <Descriptions.Item label="资源类型">
              <Tag>{detailPermission.resource_type}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="操作">
              {getActionTag(detailPermission.action)}
            </Descriptions.Item>
            <Descriptions.Item label="描述">{detailPermission.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailPermission.created_at ? new Date(detailPermission.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailPermission.updated_at ? new Date(detailPermission.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </Card>
  );
};

export default Permissions;

