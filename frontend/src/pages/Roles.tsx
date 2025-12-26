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
  Tree,
  Checkbox,
  Descriptions,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { roleService, Role } from '../services/role';
import { permissionService, Permission } from '../services/permission';

const Roles: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [permissionModalVisible, setPermissionModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [detailRole, setDetailRole] = useState<Role | null>(null);
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [rolePermissions, setRolePermissions] = useState<number[]>([]);
  const [currentRoleId, setCurrentRoleId] = useState<number | null>(null);
  const [form] = Form.useForm();
  const [permissionForm] = Form.useForm();

  useEffect(() => {
    loadRoles();
    loadAllPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const loadRoles = async () => {
    setLoading(true);
    try {
      const response = await roleService.list(page, pageSize);
      setRoles(response.roles);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadAllPermissions = async () => {
    try {
      const response = await permissionService.listAll();
      setAllPermissions(response.permissions);
    } catch (error: any) {
      message.error('加载权限列表失败');
    }
  };

  const handleViewDetail = async (role: Role) => {
    try {
      const response = await roleService.get(role.id);
      setDetailRole(response.role);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取角色详情失败');
    }
  };

  const handleCreate = () => {
    setEditingRole(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = async (role: Role) => {
    try {
      const response = await roleService.get(role.id);
      setEditingRole(response.role);
      form.setFieldsValue({
        name: response.role.name,
        description: response.role.description,
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取角色详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await roleService.delete(id);
      message.success('删除成功');
      loadRoles();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingRole) {
        await roleService.update(editingRole.id, values);
        message.success('更新成功');
      } else {
        await roleService.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      loadRoles();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleManagePermissions = async (role: Role) => {
    try {
      const response = await roleService.get(role.id);
      const roleWithPermissions = response.role;
      setCurrentRoleId(role.id);
      setRolePermissions(
        roleWithPermissions.permissions?.map((p) => p.id) || []
      );
      // 确保权限列表已加载
      if (allPermissions.length === 0) {
        await loadAllPermissions();
      }
      setPermissionModalVisible(true);
    } catch (error: any) {
      message.error('获取角色权限失败');
    }
  };

  const handlePermissionChange = async (checked: boolean, permissionId: number) => {
    if (!currentRoleId) return;

    try {
      if (checked) {
        await roleService.assignPermission(currentRoleId, permissionId);
        setRolePermissions([...rolePermissions, permissionId]);
        message.success('权限分配成功');
      } else {
        await roleService.removePermission(currentRoleId, permissionId);
        setRolePermissions(rolePermissions.filter((id) => id !== permissionId));
        message.success('权限移除成功');
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  // 按资源类型分组权限
  const groupedPermissions = allPermissions.reduce((acc, permission) => {
    if (!acc[permission.resource_type]) {
      acc[permission.resource_type] = [];
    }
    acc[permission.resource_type].push(permission);
    return acc;
  }, {} as Record<string, Permission[]>);

  const columns = [
    {
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '权限数量',
      key: 'permission_count',
      width: 100,
      render: (_: any, record: Role) => (
        <span>{record.permissions?.length || 0}</span>
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
      width: 250,
      fixed: 'right' as const,
      render: (_: any, record: Role) => (
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
            onClick={() => handleManagePermissions(record)}
          >
            权限管理
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
            title="确定要删除这个角色吗？"
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
          <h2 style={{ margin: 0 }}>角色管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={loadRoles}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              创建角色
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={roles}
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

      {/* 创建/编辑角色模态框 */}
      <Modal
        title={editingRole ? '编辑角色' : '创建角色'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea
              placeholder="请输入角色描述"
              rows={4}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 权限管理模态框 */}
      <Modal
        title="权限管理"
        open={permissionModalVisible}
        onCancel={() => setPermissionModalVisible(false)}
        footer={null}
        width={600}
      >
        <div style={{ maxHeight: 500, overflowY: 'auto' }}>
          {Object.entries(groupedPermissions).map(([resourceType, permissions]) => (
            <div key={resourceType} style={{ marginBottom: 24 }}>
              <h4 style={{ marginBottom: 12, fontWeight: 'bold' }}>
                {resourceType}
              </h4>
              <Space direction="vertical" style={{ width: '100%' }}>
                {permissions.map((permission) => (
                  <div key={permission.id} style={{ display: 'flex', alignItems: 'center' }}>
                    <Checkbox
                      checked={rolePermissions.includes(permission.id)}
                      onChange={(e) =>
                        handlePermissionChange(e.target.checked, permission.id)
                      }
                    >
                      <span style={{ marginLeft: 8 }}>
                        {permission.name} ({permission.action})
                      </span>
                    </Checkbox>
                    {permission.description && (
                      <span style={{ marginLeft: 8, color: '#999', fontSize: 12 }}>
                        - {permission.description}
                      </span>
                    )}
                  </div>
                ))}
              </Space>
            </div>
          ))}
        </div>
      </Modal>

      {/* 角色详情模态框 */}
      <Modal
        title="角色详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {detailRole && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="角色名称">{detailRole.name}</Descriptions.Item>
            <Descriptions.Item label="描述">{detailRole.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="权限数量">{detailRole.permissions?.length || 0}</Descriptions.Item>
            <Descriptions.Item label="权限列表">
              <div style={{ maxHeight: 200, overflowY: 'auto' }}>
                <Space wrap>
                  {detailRole.permissions?.map((p) => (
                    <Tag key={p.id} color="blue">{p.name}</Tag>
                  )) || <span>无</span>}
                </Space>
              </div>
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailRole.created_at ? new Date(detailRole.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailRole.updated_at ? new Date(detailRole.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default Roles;

