// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { taskTemplateApi, TaskTemplate, CreateTemplateRequest, UpdateTemplateRequest } from '@/api/task'

const { TextArea } = Input

const TemplateList = () => {
  const [templates, setTemplates] = useState<TaskTemplate[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTemplate, setEditingTemplate] = useState<TaskTemplate | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchTemplates()
  }, [])

  const fetchTemplates = async () => {
    setLoading(true)
    try {
      const response = await taskTemplateApi.list()
      setTemplates(response.data)
    } catch (error: any) {
      message.error('获取模板列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingTemplate(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (template: TaskTemplate) => {
    setEditingTemplate(template)
    form.setFieldsValue(template)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个模板吗？',
      onOk: async () => {
        try {
          await taskTemplateApi.delete(id)
          message.success('删除成功')
          fetchTemplates()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateTemplateRequest | UpdateTemplateRequest) => {
    try {
      if (editingTemplate) {
        await taskTemplateApi.update(editingTemplate.id, values)
        message.success('更新成功')
      } else {
        await taskTemplateApi.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchTemplates()
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
      title: '模板名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={type === 'shell' ? 'blue' : type === 'python' ? 'green' : 'default'}>
          {type || 'shell'}
        </Tag>
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
      render: (_: any, record: TaskTemplate) => (
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
        <h2>任务模板管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增模板">
          新增模板
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={templates}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingTemplate ? '编辑模板' : '新增模板'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 800 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="模板名称"
            rules={[{ required: true, message: '请输入模板名称' }]}
          >
            <Input placeholder="模板名称" />
          </Form.Item>
          <Form.Item name="type" label="类型" initialValue="shell">
            <Select>
              <Select.Option value="shell">Shell</Select.Option>
              <Select.Option value="python">Python</Select.Option>
              <Select.Option value="bash">Bash</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="content"
            label="模板内容"
            rules={[{ required: true, message: '请输入模板内容' }]}
          >
            <TextArea rows={10} placeholder="输入脚本内容..." />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} placeholder="模板描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default TemplateList
