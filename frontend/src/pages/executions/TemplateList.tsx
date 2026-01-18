// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { executionTemplateApi, ExecutionTemplate, CreateTemplateRequest, UpdateTemplateRequest } from '@/api/execution'

const { TextArea } = Input

// 不同类型的 shebang
const SHEBANGS: Record<string, string> = {
  shell: '#!/usr/bin/env bash\n\n',
  python: '#!/usr/bin/env python3\n\n',
}

const TemplateList = () => {
  const [templates, setTemplates] = useState<ExecutionTemplate[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTemplate, setEditingTemplate] = useState<ExecutionTemplate | null>(null)
  const [form] = Form.useForm()

  // 处理类型变更，自动插入对应的 shebang
  const handleTypeChange = (newType: string) => {
    const currentContent = form.getFieldValue('content') || ''
    const newShebang = SHEBANGS[newType] || ''
    
    // 检查当前内容是否为空或只包含其他类型的 shebang
    const isEmptyOrShebangOnly = !currentContent.trim() || 
      Object.values(SHEBANGS).some(shebang => currentContent.trim() === shebang.trim())
    
    if (isEmptyOrShebangOnly) {
      form.setFieldsValue({ content: newShebang })
    } else {
      // 如果内容不为空，检查是否以其他 shebang 开头，替换它
      let updatedContent = currentContent
      for (const shebang of Object.values(SHEBANGS)) {
        if (currentContent.startsWith(shebang.trim())) {
          updatedContent = currentContent.replace(shebang.trim(), newShebang.trim())
          break
        }
      }
      // 如果没有任何 shebang，在开头插入
      if (updatedContent === currentContent && !currentContent.startsWith('#!')) {
        updatedContent = newShebang + currentContent
      }
      form.setFieldsValue({ content: updatedContent })
    }
  }

  useEffect(() => {
    fetchTemplates()
  }, [])

  const fetchTemplates = async () => {
    setLoading(true)
    try {
      const response = await executionTemplateApi.list()
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
    // 新建时默认插入 shell 的 shebang
    form.setFieldsValue({ content: SHEBANGS.shell })
    setModalVisible(true)
  }

  const handleEdit = (template: ExecutionTemplate) => {
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
          await executionTemplateApi.delete(id)
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
        await executionTemplateApi.update(editingTemplate.id, values)
        message.success('更新成功')
      } else {
        await executionTemplateApi.create(values)
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
      render: (_: any, record: ExecutionTemplate) => (
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
            <Select onChange={handleTypeChange}>
              <Select.Option value="shell">Shell</Select.Option>
              <Select.Option value="python">Python</Select.Option>
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
