// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Tag, Descriptions } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons'
import { sshkeyApi, SSHKey, CreateSSHKeyRequest, UpdateSSHKeyRequest } from '@/api/sshkey'

const { TextArea } = Input

const SSHKeyList = () => {
  const [keys, setKeys] = useState<SSHKey[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [editingKey, setEditingKey] = useState<SSHKey | null>(null)
  const [detailKey, setDetailKey] = useState<SSHKey | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchKeys()
  }, [])

  const fetchKeys = async () => {
    setLoading(true)
    try {
      const response = await sshkeyApi.list()
      setKeys(response.data.data)
    } catch (error: any) {
      message.error('获取 SSH 密钥列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingKey(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (key: SSHKey) => {
    setEditingKey(key)
    form.setFieldsValue(key)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个 SSH 密钥吗？',
      onOk: async () => {
        try {
          await sshkeyApi.delete(id)
          message.success('删除成功')
          fetchKeys()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleViewDetail = (key: SSHKey) => {
    setDetailKey(key)
    setDetailModalVisible(true)
  }

  const handleSubmit = async (values: CreateSSHKeyRequest | UpdateSSHKeyRequest) => {
    try {
      if (editingKey) {
        await sshkeyApi.update(editingKey.id, values)
        message.success('更新成功')
      } else {
        await sshkeyApi.create(values as CreateSSHKeyRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchKeys()
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
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => <Tag>{type || 'unknown'}</Tag>,
    },
    {
      title: 'SSH 用户',
      dataIndex: 'ssh_user',
      key: 'ssh_user',
    },
    {
      title: '指纹',
      dataIndex: 'fingerprint',
      key: 'fingerprint',
      ellipsis: true,
    },
    {
      title: '公钥',
      dataIndex: 'public_key',
      key: 'public_key',
      ellipsis: true,
      render: (key: string) => key.substring(0, 50) + '...',
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
      width: 200,
      render: (_: any, record: SSHKey) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(record)}>
            详情
          </Button>
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
        <h2>SSH 密钥管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增SSH密钥">
          新增密钥
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={keys}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingKey ? '编辑 SSH 密钥' : '新增 SSH 密钥'}
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
            label="名称"
            rules={[{ required: true, message: '请输入密钥名称' }]}
          >
            <Input placeholder="密钥名称" />
          </Form.Item>
          <Form.Item name="ssh_user" label="SSH 用户">
            <Input placeholder="默认 SSH 用户名（如：root, ubuntu）" />
          </Form.Item>
          {!editingKey && (
            <>
              <Form.Item
                name="private_key"
                label="私钥"
                rules={[{ required: true, message: '请输入私钥' }]}
                help="公钥将自动从私钥中提取"
              >
                <TextArea rows={6} placeholder="SSH 私钥内容（将被加密存储，公钥将自动提取）" />
              </Form.Item>
              <Form.Item name="passphrase" label="私钥密码（如果有）">
                <Input.Password placeholder="如果私钥已加密，请输入密码" />
              </Form.Item>
            </>
          )}
          <Form.Item name="description" label="描述">
            <TextArea rows={3} placeholder="密钥描述" />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="SSH 密钥详情"
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false)
          setDetailKey(null)
        }}
        footer={[
          <Button key="close" onClick={() => {
            setDetailModalVisible(false)
            setDetailKey(null)
          }}>
            关闭
          </Button>,
        ]}
        width="90%"
        style={{ maxWidth: 800 }}
      >
        {detailKey ? (
          <Descriptions column={2} bordered>
            <Descriptions.Item label="ID">{detailKey.id}</Descriptions.Item>
            <Descriptions.Item label="名称">{detailKey.name}</Descriptions.Item>
            <Descriptions.Item label="类型">
              <Tag>{detailKey.type || 'unknown'}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="SSH 用户">{detailKey.ssh_user || '-'}</Descriptions.Item>
            <Descriptions.Item label="指纹" span={2}>
              {detailKey.fingerprint || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="公钥" span={2}>
              <div style={{ wordBreak: 'break-all', marginBottom: 8 }}>
                {detailKey.public_key || '-'}
              </div>
              {detailKey.public_key && (
                <Button
                  size="small"
                  onClick={async () => {
                    try {
                      await navigator.clipboard.writeText(detailKey.public_key)
                      message.success('公钥已复制到剪贴板')
                    } catch (error) {
                      message.error('复制失败，请手动复制')
                    }
                  }}
                >
                  复制
                </Button>
              )}
            </Descriptions.Item>
            <Descriptions.Item label="描述" span={2}>
              {detailKey.description || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="最后使用时间">
              {detailKey.last_used_at ? new Date(detailKey.last_used_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {new Date(detailKey.created_at).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {new Date(detailKey.updated_at).toLocaleString()}
            </Descriptions.Item>
          </Descriptions>
        ) : null}
      </Modal>
    </div>
  )
}

export default SSHKeyList
