// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { categoryApi, Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/api/category'

const CategoryList = () => {
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
        await categoryApi.create(values)
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
      render: (_: any, record: Category) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ 
        marginBottom: 16, 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: 8
      }}>
        <h2>资产分类管理</h2>
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
