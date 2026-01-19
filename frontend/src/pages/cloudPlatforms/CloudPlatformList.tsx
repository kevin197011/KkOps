// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { cloudPlatformApi, CloudPlatform, CreateCloudPlatformRequest, UpdateCloudPlatformRequest } from '@/api/cloudPlatform'
import { usePermissionStore } from '@/stores/permission'

const { TextArea } = Input

const CloudPlatformList = () => {
  const { hasPermission } = usePermissionStore()
  const [cloudPlatforms, setCloudPlatforms] = useState<CloudPlatform[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingCloudPlatform, setEditingCloudPlatform] = useState<CloudPlatform | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchCloudPlatforms()
  }, [])

  const fetchCloudPlatforms = async () => {
    setLoading(true)
    try {
      const response = await cloudPlatformApi.list()
      setCloudPlatforms(response.data)
    } catch (error: any) {
      message.error('获取云平台列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingCloudPlatform(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (cloudPlatform: CloudPlatform) => {
    setEditingCloudPlatform(cloudPlatform)
    form.setFieldsValue({
      name: cloudPlatform.name,
      description: cloudPlatform.description,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个云平台吗？',
      onOk: async () => {
        try {
          await cloudPlatformApi.delete(id)
          message.success('删除成功')
          fetchCloudPlatforms()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateCloudPlatformRequest | UpdateCloudPlatformRequest) => {
    try {
      if (editingCloudPlatform) {
        await cloudPlatformApi.update(editingCloudPlatform.id, values as UpdateCloudPlatformRequest)
        message.success('更新成功')
      } else {
        await cloudPlatformApi.create(values as CreateCloudPlatformRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchCloudPlatforms()
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
      title: '云平台名称',
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
      render: (_: any, record: CloudPlatform) => {
        const actions = []
        if (hasPermission('cloud-platforms', 'update')) {
          actions.push(
            <Button
              key="edit"
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
              aria-label={`编辑云平台 ${record.name}`}
            >
              编辑
            </Button>
          )
        }
        if (hasPermission('cloud-platforms', 'delete')) {
          actions.push(
            <Button
              key="delete"
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
              onClick={() => handleDelete(record.id)}
              aria-label={`删除云平台 ${record.name}`}
            >
              删除
            </Button>
          )
        }
        return actions.length > 0 ? <Space size="small">{actions}</Space> : '-'
      },
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
        <h2>云平台管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增云平台">
          新增云平台
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={cloudPlatforms}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingCloudPlatform ? '编辑云平台' : '新增云平台'}
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
          <Form.Item name="name" label="云平台名称" rules={[{ required: true, message: '请输入云平台名称' }]}>
            <Input placeholder="云平台名称" aria-label="云平台名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={4} placeholder="云平台描述" aria-label="云平台描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CloudPlatformList
