// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, Tag, ColorPicker, Card, Typography, App } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, TagsOutlined } from '@ant-design/icons'
import { tagApi, Tag as TagType, CreateTagRequest, UpdateTagRequest } from '@/api/tag'
import { usePermissionStore } from '@/stores/permission'
import { useThemeStore } from '@/stores/theme'

const { Title } = Typography

const TagList = () => {
  const { message } = App.useApp()
  const { hasPermission } = usePermissionStore()
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
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
    form.setFieldsValue({ color: '#1890ff' })
    setModalVisible(true)
  }

  const handleEdit = (tag: TagType) => {
    setEditingTag(tag)
    form.setFieldsValue({
      ...tag,
      color: tag.color || '#1890ff',
    })
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
      render: (_: any, record: TagType) => {
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
    <div style={{ 
      padding: 24,
      background: isDark ? '#0F172A' : '#F5F5F5',
      minHeight: '100vh',
    }}>
      <Card
        style={{
          background: isDark ? '#1E293B' : '#FFFFFF',
          border: isDark 
            ? '1px solid rgba(255, 255, 255, 0.08)' 
            : '1px solid rgba(0, 0, 0, 0.06)',
          boxShadow: isDark
            ? '0 4px 24px rgba(0, 0, 0, 0.5)'
            : '0 2px 12px rgba(0, 0, 0, 0.08)',
        }}
      >
        <div style={{ 
          marginBottom: 24, 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 16
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <TagsOutlined style={{ 
              fontSize: 24, 
              color: isDark ? '#60A5FA' : '#2563EB' 
            }} />
            <Title level={3} style={{ 
              margin: 0,
              color: isDark ? '#E2E8F0' : '#1E293B',
            }}>
              标签管理
            </Title>
          </div>
          {hasPermission('assets', 'create') && (
            <Button 
              type="primary" 
              icon={<PlusOutlined />} 
              onClick={handleCreate} 
              size="large"
              style={{
                borderRadius: 6,
                fontWeight: 500,
              }}
            >
              新增标签
            </Button>
          )}
        </div>
        <Table
          columns={columns}
          dataSource={tags}
          loading={loading}
          rowKey="id"
          scroll={{ x: 'max-content' }}
          style={{
            background: isDark ? '#1E293B' : '#FFFFFF',
          }}
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
            pageSizeOptions: ['10', '20', '50', '100'],
          }}
        />
      </Card>
      <Modal
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <TagsOutlined style={{ 
              fontSize: 18, 
              color: isDark ? '#60A5FA' : '#2563EB' 
            }} />
            <span>{editingTag ? '编辑标签' : '新增标签'}</span>
          </div>
        }
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width={600}
        styles={{
          body: {
            background: isDark ? '#1E293B' : '#FFFFFF',
          },
        }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="标签名称"
            rules={[{ required: true, message: '请输入标签名称' }]}
          >
            <Input placeholder="标签名称" />
          </Form.Item>
          <Form.Item name="color" label="颜色" initialValue="#1890ff">
            <ColorPicker
              showText
              format="hex"
              onChange={(color) => form.setFieldsValue({ color: color.toHexString() })}
            />
          </Form.Item>
          <Form.Item label="预览" dependencies={['name', 'color']}>
            {() => {
              const name = form.getFieldValue('name') || '标签预览'
              const color = form.getFieldValue('color') || '#1890ff'
              return <Tag color={color}>{name}</Tag>
            }}
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
