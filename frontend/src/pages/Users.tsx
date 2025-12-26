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
  Transfer,
  Descriptions,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  UserOutlined,
  EyeOutlined,
  LockOutlined,
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
  const [roleModalVisible, setRoleModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [detailUser, setDetailUser] = useState<User | null>(null);
  const [userRoles, setUserRoles] = useState<number[]>([]);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [passwordUser, setPasswordUser] = useState<User | null>(null);
  const [form] = Form.useForm();
  const [passwordForm] = Form.useForm();

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

  const loadUserRoles = async (userId: number) => {
    try {
      const response = await api.get(`/users/${userId}`);
      const user = response.data.user;
      setUserRoles(user.roles?.map((r: Role) => r.id) || []);
    } catch (error: any) {
      message.error('加载用户角色失败');
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
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    form.setFieldsValue({
      username: user.username,
      email: user.email,
      display_name: user.display_name,
    });
    setModalVisible(true);
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

  const handleAssignRole = async (user: User) => {
    setSelectedUser(user);
    await loadUserRoles(user.id);
    setRoleModalVisible(true);
  };

  const handleRoleSubmit = async () => {
    if (!selectedUser) return;

    try {
      // 获取当前用户角色
      const currentResponse = await api.get(`/users/${selectedUser.id}`);
      const currentRoles = currentResponse.data.user.roles?.map((r: Role) => r.id) || [];

      // 找出需要添加和删除的角色
      const toAdd = userRoles.filter((id: number) => !currentRoles.includes(id));
      const toRemove = currentRoles.filter((id: number) => !userRoles.includes(id));

      // 添加新角色
      for (const roleId of toAdd) {
        await userService.assignRole(selectedUser.id, roleId);
      }

      // 删除旧角色
      for (const roleId of toRemove) {
        await userService.removeRole(selectedUser.id, roleId);
      }

      message.success('角色分配成功');
      setRoleModalVisible(false);
      loadUsers();
    } catch (error: any) {
      message.error(error.response?.data?.error || '角色分配失败');
    }
  };

  const handleResetPassword = (user: User) => {
    setPasswordUser(user);
    passwordForm.resetFields();
    setPasswordModalVisible(true);
  };

  const handlePasswordSubmit = async () => {
    if (!passwordUser) return;

    try {
      const values = await passwordForm.validateFields();
      await userService.resetPassword(passwordUser.id, values.new_password);
      message.success('密码重置成功');
      setPasswordModalVisible(false);
      passwordForm.resetFields();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '密码重置失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingUser) {
        await userService.update(editingUser.id, values as UpdateUserRequest);
        message.success('更新成功');
      } else {
        await userService.create(values as CreateUserRequest);
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
        <Space>
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
      width: 320,
      render: (_: any, record: User) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<LockOutlined />}
            onClick={() => handleResetPassword(record)}
          >
            密码
          </Button>
          <Button
            type="link"
            icon={<UserOutlined />}
            onClick={() => handleAssignRole(record)}
          >
            角色
          </Button>
          <Popconfirm
            title="确定要删除这个用户吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
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
        width={600}
      >
        <Form form={form} layout="vertical">
          {!editingUser && (
            <>
              <Form.Item
                name="username"
                label="用户名"
                rules={[{ required: true, message: '请输入用户名' }]}
              >
                <Input placeholder="请输入用户名" />
              </Form.Item>
              <Form.Item
                name="password"
                label="密码"
                rules={[{ required: true, message: '请输入密码' }]}
              >
                <Input.Password placeholder="请输入密码" />
              </Form.Item>
            </>
          )}
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
          <Form.Item
            name="display_name"
            label="显示名称"
          >
            <Input placeholder="请输入显示名称" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 角色分配模态框 */}
      <Modal
        title={`分配角色 - ${selectedUser?.username || ''}`}
        open={roleModalVisible}
        onOk={handleRoleSubmit}
        onCancel={() => setRoleModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={500}
      >
        <Select
          mode="multiple"
          style={{ width: '100%' }}
          placeholder="选择角色"
          value={userRoles}
          onChange={setUserRoles}
        >
          {roles.map((role) => (
            <Option key={role.id} value={role.id}>
              {role.name} - {role.description}
            </Option>
          ))}
        </Select>
      </Modal>

      {/* 用户详情模态框 */}
      <Modal
        title="用户详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {detailUser && (
          <Descriptions bordered column={2} size="small">
            <Descriptions.Item label="用户名">{detailUser.username}</Descriptions.Item>
            <Descriptions.Item label="显示名称">{detailUser.display_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{detailUser.email}</Descriptions.Item>
            <Descriptions.Item label="状态" span={2}>
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

      {/* 重置密码模态框 */}
      <Modal
        title={`重置密码 - ${passwordUser?.username || ''}`}
        open={passwordModalVisible}
        onOk={handlePasswordSubmit}
        onCancel={() => {
          setPasswordModalVisible(false);
          passwordForm.resetFields();
        }}
        okText="确定"
        cancelText="取消"
        width={400}
      >
        <Form form={passwordForm} layout="vertical">
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码长度至少6位' },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认密码"
            dependencies={['new_password']}
            rules={[
              { required: true, message: '请确认新密码' },
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
        </Form>
      </Modal>
    </>
  );
};

export default Users;

