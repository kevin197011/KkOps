// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, Tag, Tooltip, Checkbox, Switch } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, TeamOutlined, LockOutlined } from '@ant-design/icons'
import { userApi, User, CreateUserRequest, UpdateUserRequest, UserRoleInfo } from '@/api/user'
import { roleApi, Role } from '@/api/role'

const UserList = () => {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [form] = Form.useForm()
  const [resetPassword, setResetPassword] = useState(false)

  // 角色分配相关状态
  const [roleModalVisible, setRoleModalVisible] = useState(false)
  const [roleUser, setRoleUser] = useState<User | null>(null)
  const [allRoles, setAllRoles] = useState<Role[]>([])
  const [selectedRoleIds, setSelectedRoleIds] = useState<number[]>([])
  const [roleLoading, setRoleLoading] = useState(false)
  const [userRolesMap, setUserRolesMap] = useState<Record<number, UserRoleInfo[]>>({})

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const response = await userApi.list(page, pageSize)
      setUsers(response.data.data)
      setTotal(response.data.total)
    } catch (error: any) {
      message.error('获取用户列表失败')
    } finally {
      setLoading(false)
    }
  }

  // 获取所有角色
  const fetchAllRoles = useCallback(async () => {
    try {
      const response = await roleApi.list()
      setAllRoles(response.data || [])
    } catch (error: any) {
      // 忽略错误
    }
  }, [])

  // 获取用户的角色
  const fetchUserRoles = useCallback(async (userId: number) => {
    try {
      const response = await userApi.getUserRoles(userId)
      const roles = response.data.data || []
      setSelectedRoleIds(roles.map((r) => r.id))
    } catch (error: any) {
      message.error('获取用户角色失败')
    }
  }, [])

  // 批量获取用户角色（用于表格显示）
  const fetchAllUsersRoles = useCallback(async (userList: User[]) => {
    const rolesMap: Record<number, UserRoleInfo[]> = {}
    for (const user of userList) {
      try {
        const response = await userApi.getUserRoles(user.id)
        rolesMap[user.id] = response.data.data || []
      } catch {
        rolesMap[user.id] = []
      }
    }
    setUserRolesMap(rolesMap)
  }, [])

  // 初始化加载
  useEffect(() => {
    fetchUsers()
    fetchAllRoles()
  }, [page, pageSize, fetchAllRoles])

  // 当用户列表加载完成后，获取用户角色
  useEffect(() => {
    if (users.length > 0) {
      fetchAllUsersRoles(users)
    }
  }, [users, fetchAllUsersRoles])

  const handleCreate = () => {
    setEditingUser(null)
    setResetPassword(false)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (user: User) => {
    setEditingUser(user)
    setResetPassword(false)
    form.setFieldsValue({
      email: user.email,
      phone: user.phone,
      real_name: user.real_name,
      status: user.status,
      new_password: '',
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个用户吗？',
      onOk: async () => {
        try {
          await userApi.delete(id)
          message.success('删除成功')
          fetchUsers()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingUser) {
        // 更新用户信息
        const { new_password, ...updateValues } = values
        await userApi.update(editingUser.id, updateValues as UpdateUserRequest)
        
        // 如果启用了密码重置，单独调用重置密码API
        if (resetPassword && new_password) {
          await userApi.resetPassword(editingUser.id, new_password)
          message.success('用户信息和密码更新成功')
        } else {
          message.success('更新成功')
        }
      } else {
        await userApi.create(values as CreateUserRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      setResetPassword(false)
      form.resetFields()
      fetchUsers()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 打开角色分配弹窗
  const handleOpenRoleModal = async (user: User) => {
    setRoleUser(user)
    setRoleLoading(true)
    setRoleModalVisible(true)
    await Promise.all([fetchAllRoles(), fetchUserRoles(user.id)])
    setRoleLoading(false)
  }

  // 保存用户角色
  const handleSaveUserRoles = async () => {
    if (!roleUser) return

    setRoleLoading(true)
    try {
      await userApi.setUserRoles(roleUser.id, selectedRoleIds)
      message.success('角色分配成功')
      setRoleModalVisible(false)
      // 刷新用户角色显示
      fetchAllUsersRoles(users)
    } catch (error: any) {
      message.error(error.response?.data?.error || '角色分配失败')
    } finally {
      setRoleLoading(false)
    }
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '姓名',
      dataIndex: 'real_name',
      key: 'real_name',
    },
    {
      title: '角色',
      key: 'roles',
      width: 200,
      render: (_: any, record: User) => {
        const roles = userRolesMap[record.id] || []
        if (roles.length === 0) {
          return <span style={{ color: '#999' }}>暂无角色</span>
        }
        return (
          <Space size={[0, 4]} wrap>
            {roles.map((role) => (
              <Tag
                key={role.id}
                color={role.is_admin ? 'gold' : 'blue'}
              >
                {role.name}
              </Tag>
            ))}
          </Space>
        )
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status === 'active' ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: User) => (
        <Space size="small">
          <Tooltip title="角色分配">
            <Button
              type="link"
              size="small"
              icon={<TeamOutlined />}
              onClick={() => handleOpenRoleModal(record)}
              aria-label={`角色分配 ${record.username}`}
            >
              角色
            </Button>
          </Tooltip>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            aria-label={`编辑用户 ${record.username}`}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            aria-label={`删除用户 ${record.username}`}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div
        style={{
          marginBottom: 16,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 8,
        }}
      >
        <h2>用户管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增用户">
          新增用户
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={users}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          responsive: true,
          onChange: (page, pageSize) => {
            setPage(page)
            setPageSize(pageSize)
          },
        }}
      />
      <Modal
        title={editingUser ? '编辑用户' : '新增用户'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          setResetPassword(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 600 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          {!editingUser && (
            <Form.Item
              name="username"
              label="用户名"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input placeholder="用户名" aria-label="用户名" />
            </Form.Item>
          )}
          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input placeholder="邮箱" type="email" aria-label="邮箱" />
          </Form.Item>
          {!editingUser && (
            <Form.Item
              name="password"
              label="密码"
              rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '密码至少6位' }]}
            >
              <Input.Password placeholder="密码" aria-label="密码" />
            </Form.Item>
          )}
          {editingUser && (
            <>
              <Form.Item label="重置密码">
                <Space>
                  <Switch
                    checked={resetPassword}
                    onChange={setResetPassword}
                    checkedChildren={<LockOutlined />}
                    unCheckedChildren={<LockOutlined />}
                  />
                  <span style={{ color: '#888' }}>
                    {resetPassword ? '启用密码重置' : '不修改密码'}
                  </span>
                </Space>
              </Form.Item>
              {resetPassword && (
                <Form.Item
                  name="new_password"
                  label="新密码"
                  rules={[
                    { required: resetPassword, message: '请输入新密码' },
                    { min: 6, message: '密码至少6位' },
                  ]}
                >
                  <Input.Password placeholder="新密码" aria-label="新密码" />
                </Form.Item>
              )}
            </>
          )}
          <Form.Item name="real_name" label="姓名">
            <Input placeholder="姓名" aria-label="姓名" />
          </Form.Item>
          <Form.Item name="phone" label="电话">
            <Input placeholder="电话" aria-label="电话" />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="active">
            <Select aria-label="状态">
              <Select.Option value="active">启用</Select.Option>
              <Select.Option value="disabled">禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      {/* 角色分配弹窗 */}
      <Modal
        title={
          <Space>
            <TeamOutlined />
            <span>角色分配 - {roleUser?.username}</span>
          </Space>
        }
        open={roleModalVisible}
        onCancel={() => setRoleModalVisible(false)}
        onOk={handleSaveUserRoles}
        confirmLoading={roleLoading}
        width={500}
        okText="保存"
        cancelText="取消"
      >
        <div style={{ marginBottom: 16 }}>
          <Tag color="blue">勾选要分配给用户的角色，用户将获得角色所授权的资产访问权限</Tag>
        </div>
        <div style={{ maxHeight: 400, overflowY: 'auto' }}>
          <Checkbox.Group
            value={selectedRoleIds}
            onChange={(checkedValues) => setSelectedRoleIds(checkedValues as number[])}
            style={{ width: '100%' }}
          >
            <Space direction="vertical" style={{ width: '100%' }}>
              {allRoles.map((role) => (
                <div
                  key={role.id}
                  style={{
                    padding: '12px 16px',
                    border: '1px solid #f0f0f0',
                    borderRadius: 6,
                    marginBottom: 8,
                    background: selectedRoleIds.includes(role.id) ? '#f6ffed' : '#fff',
                  }}
                >
                  <Checkbox value={role.id}>
                    <Space>
                      <span style={{ fontWeight: 500 }}>{role.name}</span>
                      {role.is_admin && <Tag color="gold">管理员</Tag>}
                      {!role.is_admin && role.asset_count >= 0 && (
                        <Tag color="default">{role.asset_count} 台资产</Tag>
                      )}
                    </Space>
                  </Checkbox>
                  {role.description && (
                    <div style={{ marginLeft: 24, marginTop: 4, color: '#888', fontSize: 12 }}>
                      {role.description}
                    </div>
                  )}
                </div>
              ))}
            </Space>
          </Checkbox.Group>
        </div>
      </Modal>
    </div>
  )
}

export default UserList
