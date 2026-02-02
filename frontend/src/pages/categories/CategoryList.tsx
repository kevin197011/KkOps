// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, Card, Typography, theme } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, AppstoreOutlined } from '@ant-design/icons'
import { categoryApi, Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/api/category'
import { usePermissionStore } from '@/stores/permission'

const { Title } = Typography

const CategoryList = () => {
  const { token } = theme.useToken()
  const { hasPermission } = usePermissionStore()
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchCategories()
  }, [])

  const fetchCategories = async () => {
    setLoading(true)
    try {
      const response = await categoryApi.list()
      setCategories(response.data)
    } catch (error: any) {
      message.error('获取分类列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingCategory(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (category: Category) => {
    setEditingCategory(category)
    form.setFieldsValue(category)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个分类吗？',
      onOk: async () => {
        try {
          await categoryApi.delete(id)
          message.success('删除成功')
          fetchCategories()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateCategoryRequest | UpdateCategoryRequest) => {
    try {
      if (editingCategory) {
        await categoryApi.update(editingCategory.id, values)
        message.success('更新成功')
      } else {
        await categoryApi.create({ name: values.name!, description: values.description })
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchCategories()
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
      title: '分类名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '分类代码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '父分类',
      dataIndex: 'parent_id',
      key: 'parent_id',
      render: (parentId: number | undefined) => {
        if (!parentId) return '-'
        const parent = categories.find((c) => c.id === parentId)
        return parent ? parent.name : parentId
      },
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
      render: (_: any, record: Category) => {
        const actions = []
        if (hasPermission('assets', 'update')) {
          actions.push(
            <Button key="edit" type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
              编辑
            </Button>
          )
        }
        if (hasPermission('assets', 'delete')) {
          actions.push(
            <Button key="delete" type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
              删除
            </Button>
          )
        }
        return actions.length > 0 ? <Space size="small">{actions}</Space> : '-'
      },
    },
  ]

  return (
    <div style={{ padding: 24, background: token.colorBgContainer, minHeight: '100%' }}>
      <Card
        styles={{ body: { padding: 24 } }}
        style={{ background: token.colorBgElevated, borderColor: token.colorBorderSecondary }}
      >
        <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 16 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <AppstoreOutlined style={{ fontSize: 24, color: token.colorPrimary }} />
            <Title level={3} style={{ margin: 0, color: token.colorTextHeading }}>
              资产分类管理
            </Title>
          </div>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增分类">
            新增分类
          </Button>
        </div>
        <Table
          columns={columns}
          dataSource={categories}
          loading={loading}
          rowKey="id"
          scroll={{ x: 'max-content' }}
        />
      </Card>
      <Modal
        title={editingCategory ? '编辑分类' : '新增分类'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 600 }}
        styles={{ body: { background: token.colorBgElevated } }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="分类名称"
            rules={[{ required: true, message: '请输入分类名称' }]}
          >
            <Input placeholder="分类名称" />
          </Form.Item>
          <Form.Item name="parent_id" label="父分类">
            <Select placeholder="选择父分类" allowClear>
              {categories
                .filter((c) => !c.parent_id && c.id !== editingCategory?.id)
                .map((c) => (
                  <Select.Option key={c.id} value={c.id}>
                    {c.name}
                  </Select.Option>
                ))}
            </Select>
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} placeholder="分类描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CategoryList
