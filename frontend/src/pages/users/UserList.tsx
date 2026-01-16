// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { userApi, User, CreateUserRequest, UpdateUserRequest } from '@/api/user'

const UserList = () => {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchUsers()
  }, [page, pageSize])

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

  const handleCreate = () => {
    setEditingUser(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (user: User) => {
    setEditingUser(user)
    form.setFieldsValue({
      email: user.email,
      phone: user.phone,
      real_name: user.real_name,
      status: user.status,
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

  const handleSubmit = async (values: CreateUserRequest | UpdateUserRequest) => {
    try {
      if (editingUser) {
        await userApi.update(editingUser.id, values as UpdateUserRequest)
        message.success('更新成功')
      } else {
        await userApi.create(values as CreateUserRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchUsers()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
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
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: User) => (
        <Space size="small">
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
    </div>
  )
}

export default UserList
