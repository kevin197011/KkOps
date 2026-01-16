// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { environmentApi, Environment, CreateEnvironmentRequest, UpdateEnvironmentRequest } from '@/api/environment'

const { TextArea } = Input

const EnvironmentList = () => {
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingEnvironment, setEditingEnvironment] = useState<Environment | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchEnvironments()
  }, [])

  const fetchEnvironments = async () => {
    setLoading(true)
    try {
      const response = await environmentApi.list()
      setEnvironments(response.data)
    } catch (error: any) {
      message.error('获取环境列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingEnvironment(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (environment: Environment) => {
    setEditingEnvironment(environment)
    form.setFieldsValue({
      name: environment.name,
      description: environment.description,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个环境吗？',
      onOk: async () => {
        try {
          await environmentApi.delete(id)
          message.success('删除成功')
          fetchEnvironments()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateEnvironmentRequest | UpdateEnvironmentRequest) => {
    try {
      if (editingEnvironment) {
        await environmentApi.update(editingEnvironment.id, values as UpdateEnvironmentRequest)
        message.success('更新成功')
      } else {
        await environmentApi.create(values as CreateEnvironmentRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchEnvironments()
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
      title: '环境名称',
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
      render: (_: any, record: Environment) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            aria-label={`编辑环境 ${record.name}`}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            aria-label={`删除环境 ${record.name}`}
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
        <h2>环境管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增环境">
          新增环境
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={environments}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingEnvironment ? '编辑环境' : '新增环境'}
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
          <Form.Item name="name" label="环境名称" rules={[{ required: true, message: '请输入环境名称' }]}>
            <Input placeholder="环境名称" aria-label="环境名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={4} placeholder="环境描述" aria-label="环境描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default EnvironmentList
