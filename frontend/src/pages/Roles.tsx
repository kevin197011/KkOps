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
  Checkbox,
  Descriptions,
  Divider,
  Row,
  Col,
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
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [detailRole, setDetailRole] = useState<Role | null>(null);
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [selectedPermissions, setSelectedPermissions] = useState<number[]>([]);
  const [form] = Form.useForm();

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
    setSelectedPermissions([]);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = async (role: Role) => {
    try {
      const response = await roleService.get(role.id);
      const roleData = response.role;
      setEditingRole(roleData);
      setSelectedPermissions(roleData.permissions?.map((p: Permission) => p.id) || []);
      form.setFieldsValue({
        name: roleData.name,
        description: roleData.description,
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
        // 更新角色基本信息
        await roleService.update(editingRole.id, values);
        
        // 更新权限：先获取当前权限，计算差异
        const currentPermissions = editingRole.permissions?.map((p: Permission) => p.id) || [];
        const toAdd = selectedPermissions.filter(id => !currentPermissions.includes(id));
        const toRemove = currentPermissions.filter(id => !selectedPermissions.includes(id));
        
        // 添加新权限
        for (const permId of toAdd) {
          await roleService.assignPermission(editingRole.id, permId);
        }
        // 移除权限
        for (const permId of toRemove) {
          await roleService.removePermission(editingRole.id, permId);
        }
        
        message.success('更新成功');
      } else {
        // 创建角色
        const response = await roleService.create(values);
        const newRoleId = response.role.id;
        
        // 分配权限
        for (const permId of selectedPermissions) {
          await roleService.assignPermission(newRoleId, permId);
        }
        
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

  const handlePermissionToggle = (permissionId: number, checked: boolean) => {
    if (checked) {
      setSelectedPermissions([...selectedPermissions, permissionId]);
    } else {
      setSelectedPermissions(selectedPermissions.filter(id => id !== permissionId));
    }
  };

  const handleSelectAllInGroup = (permissions: Permission[], checked: boolean) => {
    const permIds = permissions.map(p => p.id);
    if (checked) {
      const newSelected = Array.from(new Set([...selectedPermissions, ...permIds]));
      setSelectedPermissions(newSelected);
    } else {
      setSelectedPermissions(selectedPermissions.filter(id => !permIds.includes(id)));
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

  // 资源类型中文映射
  const resourceTypeMap: Record<string, string> = {
    host: '主机管理',
    deployment: '部署管理',
    task: '任务管理',
    log: '日志管理',
    monitoring: '监控管理',
    audit: '审计管理',
    user: '用户管理',
    project: '项目管理',
    webssh: 'WebSSH',
    batch_operation: '批量操作',
  };

  const columns = [
    {
      title: '序号',
      key: 'index',
      width: 60,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
      width: 120,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '权限',
      key: 'permission_count',
      width: 70,
      render: (_: any, record: Role) => (
        <Tag color="blue">{record.permissions?.length || 0}</Tag>
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
      width: 150,
      render: (_: any, record: Role) => (
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

      {/* 创建/编辑角色模态框 */}
      <Modal
        title={editingRole ? '编辑角色' : '创建角色'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={700}
      >
        <Form form={form} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="name"
                label="角色名称"
                rules={[{ required: true, message: '请输入角色名称' }]}
              >
                <Input placeholder="请输入角色名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="description" label="描述">
                <Input placeholder="请输入角色描述" />
              </Form.Item>
            </Col>
          </Row>
        </Form>
        
        <Divider style={{ marginTop: 8 }}>
          <span>权限配置</span>
          <span style={{ fontSize: 12, color: '#999', marginLeft: 8 }}>
            (已选 {selectedPermissions.length} 项)
          </span>
        </Divider>
        
        <div style={{ maxHeight: 400, overflowY: 'auto', paddingRight: 8 }}>
          {Object.entries(groupedPermissions).map(([resourceType, permissions]) => {
            const allSelected = permissions.every(p => selectedPermissions.includes(p.id));
            const someSelected = permissions.some(p => selectedPermissions.includes(p.id));
            
            return (
              <Card 
                key={resourceType} 
                size="small" 
                style={{ marginBottom: 12 }}
                title={
                  <Checkbox
                    checked={allSelected}
                    indeterminate={someSelected && !allSelected}
                    onChange={(e) => handleSelectAllInGroup(permissions, e.target.checked)}
                  >
                    <span style={{ fontWeight: 500 }}>
                      {resourceTypeMap[resourceType] || resourceType}
                    </span>
                    <span style={{ color: '#999', marginLeft: 8, fontSize: 12 }}>
                      ({permissions.filter(p => selectedPermissions.includes(p.id)).length}/{permissions.length})
                    </span>
                  </Checkbox>
                }
              >
                <Row gutter={[8, 8]}>
                  {permissions.map((permission) => (
                    <Col span={12} key={permission.id}>
                      <Checkbox
                        checked={selectedPermissions.includes(permission.id)}
                        onChange={(e) => handlePermissionToggle(permission.id, e.target.checked)}
                      >
                        <span>{permission.name}</span>
                        <span style={{ color: '#999', fontSize: 12, marginLeft: 4 }}>
                          ({permission.action})
                        </span>
                      </Checkbox>
                    </Col>
                  ))}
                </Row>
              </Card>
            );
          })}
        </div>
      </Modal>

      {/* 角色详情模态框 */}
      <Modal
        title="角色详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={<Button onClick={() => setDetailModalVisible(false)}>关闭</Button>}
        width={600}
      >
        {detailRole && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="角色名称">{detailRole.name}</Descriptions.Item>
            <Descriptions.Item label="描述">{detailRole.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="权限数量">
              <Tag color="blue">{detailRole.permissions?.length || 0}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="权限列表">
              <div style={{ maxHeight: 200, overflowY: 'auto' }}>
                <Space wrap size={[4, 8]}>
                  {detailRole.permissions?.map((p) => (
                    <Tag key={p.id} color="blue">{p.name}</Tag>
                  )) || <span style={{ color: '#999' }}>无</span>}
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

