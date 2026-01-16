// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { roleApi, Role, CreateRoleRequest, UpdateRoleRequest } from '@/api/role'

const { TextArea } = Input

const RoleList = () => {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchRoles()
  }, [])

  const fetchRoles = async () => {
    setLoading(true)
    try {
      const response = await roleApi.list()
      setRoles(response.data)
    } catch (error: any) {
      message.error('获取角色列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingRole(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (role: Role) => {
    setEditingRole(role)
    form.setFieldsValue({
      name: role.name,
      description: role.description,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个角色吗？',
      onOk: async () => {
        try {
          await roleApi.delete(id)
          message.success('删除成功')
          fetchRoles()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateRoleRequest | UpdateRoleRequest) => {
    try {
      if (editingRole) {
        await roleApi.update(editingRole.id, values as UpdateRoleRequest)
        message.success('更新成功')
      } else {
        await roleApi.create(values as CreateRoleRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchRoles()
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
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Role) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            aria-label={`编辑角色 ${record.name}`}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            aria-label={`删除角色 ${record.name}`}
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
        <h2>角色管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增角色">
          新增角色
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={roles}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingRole ? '编辑角色' : '新增角色'}
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
          <Form.Item name="name" label="角色名称" rules={[{ required: true, message: '请输入角色名称' }]}>
            <Input placeholder="角色名称" aria-label="角色名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={4} placeholder="角色描述" aria-label="角色描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default RoleList
