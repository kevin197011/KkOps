// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { tagApi, Tag as TagType, CreateTagRequest, UpdateTagRequest } from '@/api/tag'

const TagList = () => {
  const [tags, setTags] = useState<TagType[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTag, setEditingTag] = useState<TagType | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchTags()
  }, [])

  const fetchTags = async () => {
    setLoading(true)
    try {
      const response = await tagApi.list()
      setTags(response.data)
    } catch (error: any) {
      message.error('获取标签列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingTag(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (tag: TagType) => {
    setEditingTag(tag)
    form.setFieldsValue(tag)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个标签吗？',
      onOk: async () => {
        try {
          await tagApi.delete(id)
          message.success('删除成功')
          fetchTags()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateTagRequest | UpdateTagRequest) => {
    try {
      if (editingTag) {
        await tagApi.update(editingTag.id, values)
        message.success('更新成功')
      } else {
        await tagApi.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchTags()
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
      title: '标签名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '颜色',
      dataIndex: 'color',
      key: 'color',
      render: (color: string) => (
        <Tag color={color || 'default'}>{color || '默认'}</Tag>
      ),
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
      render: (_: any, record: TagType) => (
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
        <h2>标签管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增标签">
          新增标签
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={tags}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingTag ? '编辑标签' : '新增标签'}
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
            label="标签名称"
            rules={[{ required: true, message: '请输入标签名称' }]}
          >
            <Input placeholder="标签名称" />
          </Form.Item>
          <Form.Item name="color" label="颜色">
            <Input placeholder="颜色代码（如：#1890ff）" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} placeholder="标签描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default TagList
