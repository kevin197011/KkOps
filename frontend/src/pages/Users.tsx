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
  Select,
  Descriptions,
  Divider,
  Row,
  Col,
  Checkbox,
  Switch,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { userService, CreateUserRequest, UpdateUserRequest } from '../services/user';
import { User } from '../services/auth';
import api from '../config/api';

const { Option } = Select;

interface Role {
  id: number;
  name: string;
  description?: string;
}

const Users: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [detailUser, setDetailUser] = useState<User | null>(null);
  const [selectedRoles, setSelectedRoles] = useState<number[]>([]);
  const [changePassword, setChangePassword] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    loadUsers();
    loadRoles();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const response = await userService.list(page, pageSize);
      setUsers(response.users);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadRoles = async () => {
    try {
      const response = await api.get('/roles');
      setRoles(response.data.roles || []);
    } catch (error: any) {
      message.error('加载角色列表失败');
    }
  };

  const handleViewDetail = async (user: User) => {
    try {
      const response = await api.get(`/users/${user.id}`);
      setDetailUser(response.data.user);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取用户详情失败');
    }
  };

  const handleCreate = () => {
    setEditingUser(null);
    setSelectedRoles([]);
    setChangePassword(false);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = async (user: User) => {
    try {
      const response = await api.get(`/users/${user.id}`);
      const userData = response.data.user;
      setEditingUser(userData);
      setSelectedRoles(userData.roles?.map((r: Role) => r.id) || []);
      setChangePassword(false);
      form.setFieldsValue({
        username: userData.username,
        email: userData.email,
        display_name: userData.display_name,
        status: userData.status || 'active',
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取用户详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await userService.delete(id);
      message.success('删除成功');
      loadUsers();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      if (editingUser) {
        // 更新用户基本信息
        const updateData: UpdateUserRequest = {
          email: values.email,
          display_name: values.display_name,
          status: values.status,
        };
        await userService.update(editingUser.id, updateData);
        
        // 更新角色：计算差异
        const currentRoles = editingUser.roles?.map((r: Role) => r.id) || [];
        const toAdd = selectedRoles.filter(id => !currentRoles.includes(id));
        const toRemove = currentRoles.filter(id => !selectedRoles.includes(id));
        
        // 添加新角色
        for (const roleId of toAdd) {
          await userService.assignRole(editingUser.id, roleId);
        }
        // 移除角色
        for (const roleId of toRemove) {
          await userService.removeRole(editingUser.id, roleId);
        }
        
        // 如果需要修改密码
        if (changePassword && values.new_password) {
          await userService.resetPassword(editingUser.id, values.new_password);
        }
        
        message.success('更新成功');
      } else {
        // 创建用户
        const createData: CreateUserRequest = {
          username: values.username,
          password: values.password,
          email: values.email,
          display_name: values.display_name,
        };
        const response = await userService.create(createData);
        const newUserId = response.user.id;
        
        // 分配角色
        for (const roleId of selectedRoles) {
          await userService.assignRole(newUserId, roleId);
        }
        
        message.success('创建成功');
      }
      setModalVisible(false);
      loadUsers();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleRoleToggle = (roleId: number, checked: boolean) => {
    if (checked) {
      setSelectedRoles([...selectedRoles, roleId]);
    } else {
      setSelectedRoles(selectedRoles.filter(id => id !== roleId));
    }
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      active: { color: 'success', text: '活跃' },
      inactive: { color: 'default', text: '未激活' },
      suspended: { color: 'error', text: '已暂停' },
    };
    const statusInfo = statusMap[status] || statusMap.active;
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
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '显示名称',
      dataIndex: 'display_name',
      key: 'display_name',
      render: (text: string) => text || '-',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '角色',
      key: 'roles',
      render: (_: any, record: User) => (
        <Space wrap size={[4, 4]}>
          {record.roles?.map((role: Role) => (
            <Tag key={role.id} color="blue">
              {role.name}
            </Tag>
          )) || <Tag>无</Tag>}
        </Space>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      key: 'last_login_at',
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: any, record: User) => (
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
            title="确定要删除这个用户吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger size="small" icon={<DeleteOutlined />}>
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
          <h2 style={{ margin: 0 }}>用户管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={loadUsers}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              新建用户
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={users}
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

      {/* 用户编辑模态框 */}
      <Modal
        title={editingUser ? '编辑用户' : '新建用户'}
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
                name="username"
                label="用户名"
                rules={[{ required: true, message: '请输入用户名' }]}
              >
                <Input placeholder="请输入用户名" disabled={!!editingUser} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="display_name"
                label="显示名称"
              >
                <Input placeholder="请输入显示名称" />
              </Form.Item>
            </Col>
          </Row>
          
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="email"
                label="邮箱"
                rules={[
                  { required: true, message: '请输入邮箱' },
                  { type: 'email', message: '请输入有效的邮箱地址' },
                ]}
              >
                <Input placeholder="请输入邮箱" />
              </Form.Item>
            </Col>
            <Col span={12}>
              {editingUser ? (
                <Form.Item
                  name="status"
                  label="状态"
                  rules={[{ required: true, message: '请选择状态' }]}
                >
                  <Select placeholder="请选择状态">
                    <Option value="active">活跃</Option>
                    <Option value="inactive">未激活</Option>
                    <Option value="suspended">已暂停</Option>
                  </Select>
                </Form.Item>
              ) : (
                <Form.Item
                  name="password"
                  label="密码"
                  rules={[
                    { required: true, message: '请输入密码' },
                    { min: 6, message: '密码长度至少6位' },
                  ]}
                >
                  <Input.Password placeholder="请输入密码" />
                </Form.Item>
              )}
            </Col>
          </Row>
        </Form>
        
        {/* 角色配置 */}
        <Divider style={{ marginTop: 8 }}>
          <span>角色配置</span>
          <span style={{ fontSize: 12, color: '#999', marginLeft: 8 }}>
            (已选 {selectedRoles.length} 项)
          </span>
        </Divider>
        
        <Card size="small" style={{ marginBottom: 16 }}>
          <Row gutter={[8, 8]}>
            {roles.map((role) => (
              <Col span={8} key={role.id}>
                <Checkbox
                  checked={selectedRoles.includes(role.id)}
                  onChange={(e) => handleRoleToggle(role.id, e.target.checked)}
                >
                  <span>{role.name}</span>
                  {role.description && (
                    <span style={{ color: '#999', fontSize: 12, marginLeft: 4 }}>
                      ({role.description})
                    </span>
                  )}
                </Checkbox>
              </Col>
            ))}
            {roles.length === 0 && (
              <Col span={24}>
                <span style={{ color: '#999' }}>暂无可用角色</span>
              </Col>
            )}
          </Row>
        </Card>
        
        {/* 修改密码（仅编辑时显示） */}
        {editingUser && (
          <>
            <Divider style={{ marginTop: 8 }}>
              <Space>
                <span>修改密码</span>
                <Switch
                  size="small"
                  checked={changePassword}
                  onChange={setChangePassword}
                />
              </Space>
            </Divider>
            
            {changePassword && (
              <Form form={form} layout="vertical">
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      name="new_password"
                      label="新密码"
                      rules={[
                        { required: changePassword, message: '请输入新密码' },
                        { min: 6, message: '密码长度至少6位' },
                      ]}
                    >
                      <Input.Password placeholder="请输入新密码" />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item
                      name="confirm_password"
                      label="确认密码"
                      dependencies={['new_password']}
                      rules={[
                        { required: changePassword, message: '请确认新密码' },
                        ({ getFieldValue }) => ({
                          validator(_, value) {
                            if (!value || getFieldValue('new_password') === value) {
                              return Promise.resolve();
                            }
                            return Promise.reject(new Error('两次输入的密码不一致'));
                          },
                        }),
                      ]}
                    >
                      <Input.Password placeholder="请再次输入新密码" />
                    </Form.Item>
                  </Col>
                </Row>
              </Form>
            )}
          </>
        )}
      </Modal>

      {/* 用户详情模态框 */}
      <Modal
        title="用户详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={<Button onClick={() => setDetailModalVisible(false)}>关闭</Button>}
        width={600}
      >
        {detailUser && (
          <Descriptions bordered column={2} size="small">
            <Descriptions.Item label="用户名">{detailUser.username}</Descriptions.Item>
            <Descriptions.Item label="显示名称">{detailUser.display_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{detailUser.email}</Descriptions.Item>
            <Descriptions.Item label="状态">
              {getStatusTag(detailUser.status || 'active')}
            </Descriptions.Item>
            <Descriptions.Item label="角色" span={2}>
              <Space wrap>
                {detailUser.roles?.map((role: Role) => (
                  <Tag key={role.id} color="blue">{role.name}</Tag>
                )) || <span>无</span>}
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="最后登录" span={2}>
              {detailUser.last_login_at ? new Date(detailUser.last_login_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailUser.created_at ? new Date(detailUser.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailUser.updated_at ? new Date(detailUser.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default Users;
