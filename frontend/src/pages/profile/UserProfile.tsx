// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Card, Table, Button, Space, message, Modal, Form, Input, DatePicker, Tag, Popconfirm, Typography } from 'antd'
import { PlusOutlined, DeleteOutlined, CopyOutlined, KeyOutlined, UserOutlined, MailOutlined, SafetyOutlined, EyeOutlined } from '@ant-design/icons'
import { userApi, APITokenResponse, CreateAPITokenRequest } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import moment from 'moment'

const { Title, Text } = Typography

const UserProfile = () => {
  const { user } = useAuthStore()
  const [tokens, setTokens] = useState<APITokenResponse[]>([])
  const [loading, setLoading] = useState(false)
  const [createModalVisible, setCreateModalVisible] = useState(false)
  const [tokenModalVisible, setTokenModalVisible] = useState(false)
  const [viewingToken, setViewingToken] = useState<string>('')
  const [viewingTokenName, setViewingTokenName] = useState<string>('')
  const [form] = Form.useForm()

  useEffect(() => {
    if (user) {
      fetchTokens()
    }
  }, [user])

  const fetchTokens = async () => {
    if (!user) return
    
    setLoading(true)
    try {
      const response = await userApi.listTokens(user.id)
      // 后端直接返回数组，但 axios 会包装在 data 中
      setTokens(Array.isArray(response.data) ? response.data : [])
    } catch (error: any) {
      message.error('获取 Token 列表失败: ' + (error.response?.data?.error || error.message))
    } finally {
      setLoading(false)
    }
  }

  const handleCreateToken = async (values: CreateAPITokenRequest) => {
    if (!user) return

    try {
      const response = await userApi.createToken(user.id, {
        ...values,
        expires_at: values.expires_at 
          ? (moment.isMoment(values.expires_at) 
              ? values.expires_at.toISOString() 
              : moment(values.expires_at).toISOString())
          : undefined,
      })
      
      setViewingToken(response.data.token || '')
      setViewingTokenName(response.data.name || '')
      setTokenModalVisible(true)
      setCreateModalVisible(false)
      form.resetFields()
      fetchTokens()
      message.success('Token 创建成功')
    } catch (error: any) {
      message.error('创建 Token 失败: ' + (error.response?.data?.error || error.message))
    }
  }

  const handleDeleteToken = async (tokenId: number) => {
    if (!user) return

    try {
      await userApi.deleteToken(user.id, tokenId)
      message.success('Token 已删除')
      fetchTokens()
    } catch (error: any) {
      message.error('删除 Token 失败: ' + (error.response?.data?.error || error.message))
    }
  }

  const handleCopyToken = (token: string) => {
    navigator.clipboard.writeText(token)
    message.success('Token 已复制到剪贴板')
  }

  const handleViewToken = async (tokenId: number, tokenName: string) => {
    if (!user) return

    try {
      const response = await userApi.getToken(user.id, tokenId)
      setViewingToken(response.data.token || '')
      setViewingTokenName(tokenName)
      setTokenModalVisible(true)
    } catch (error: any) {
      message.error('获取 Token 失败: ' + (error.response?.data?.error || error.message))
    }
  }

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (text: string) => text || '-',
    },
    {
      title: 'Token',
      dataIndex: 'prefix',
      key: 'prefix',
      render: (prefix: string) => (
        <Text code style={{ fontSize: 12 }}>
          {prefix}
        </Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status === 'active' ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '过期时间',
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (expiresAt: string) => {
        if (!expiresAt) return <Text type="secondary">永不过期</Text>
        const isExpired = moment(expiresAt).isBefore(moment())
        return (
          <Text type={isExpired ? 'danger' : undefined}>
            {moment(expiresAt).format('YYYY-MM-DD HH:mm:ss')}
          </Text>
        )
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (createdAt: string) => moment(createdAt).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: APITokenResponse) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewToken(record.id, record.name)}
          >
            查看
          </Button>
          <Popconfirm
            title="确定要删除此 Token 吗？"
            description="删除后该 Token 将无法再使用，此操作不可恢复"
            onConfirm={() => handleDeleteToken(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto' }}>
      {/* 用户信息卡片 */}
      <Card
        style={{
          marginBottom: 24,
          borderRadius: 8,
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)',
        }}
      >
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <div
              style={{
                width: 64,
                height: 64,
                borderRadius: '50%',
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: 28,
                color: '#fff',
              }}
            >
              <UserOutlined />
            </div>
            <div>
              <Title level={4} style={{ margin: 0 }}>
                {user?.real_name || user?.username}
              </Title>
              <Space style={{ marginTop: 8 }}>
                <Text type="secondary">
                  <MailOutlined /> {user?.email}
                </Text>
                <Text type="secondary">
                  <UserOutlined /> {user?.username}
                </Text>
              </Space>
            </div>
          </div>
        </Space>
      </Card>

      {/* Token 管理卡片 */}
      <Card
        title={
          <Space>
            <KeyOutlined />
            <span>API Token 管理</span>
          </Space>
        }
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setCreateModalVisible(true)}
          >
            创建 Token
          </Button>
        }
        style={{
          borderRadius: 8,
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)',
        }}
      >
        <Table
          columns={columns}
          dataSource={tokens}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      </Card>

      {/* 创建 Token 弹窗 */}
      <Modal
        title="创建 API Token"
        open={createModalVisible}
        onCancel={() => {
          setCreateModalVisible(false)
          form.resetFields()
        }}
        footer={null}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateToken}
        >
          <Form.Item
            name="name"
            label="Token 名称"
            rules={[{ required: true, message: '请输入 Token 名称' }]}
          >
            <Input placeholder="例如: 生产环境 API" />
          </Form.Item>
          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea
              rows={3}
              placeholder="可选：描述此 Token 的用途"
            />
          </Form.Item>
          <Form.Item
            name="expires_at"
            label="过期时间"
          >
            <DatePicker
              showTime
              format="YYYY-MM-DD HH:mm:ss"
              style={{ width: '100%' }}
              placeholder="选择过期时间（留空表示永不过期）"
            />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => {
                setCreateModalVisible(false)
                form.resetFields()
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                创建
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 显示 Token 弹窗 */}
      <Modal
        title={
          <Space>
            <SafetyOutlined style={{ color: '#52c41a' }} />
            <span>{viewingTokenName}</span>
          </Space>
        }
        open={tokenModalVisible}
        onCancel={() => {
          setTokenModalVisible(false)
          setViewingToken('')
          setViewingTokenName('')
        }}
        footer={[
          <Button
            key="copy"
            type="primary"
            icon={<CopyOutlined />}
            onClick={() => handleCopyToken(viewingToken)}
          >
            复制 Token
          </Button>,
          <Button
            key="close"
            onClick={() => {
              setTokenModalVisible(false)
              setViewingToken('')
              setViewingTokenName('')
            }}
          >
            关闭
          </Button>,
        ]}
        destroyOnClose
      >
        <div
          style={{
            padding: 16,
            background: '#f5f5f5',
            borderRadius: 4,
            marginBottom: 16,
            wordBreak: 'break-all',
          }}
        >
          <Text code style={{ fontSize: 12 }}>
            {viewingToken}
          </Text>
        </div>
        <Text type="secondary" style={{ fontSize: 12 }}>
          提示：请妥善保管此 Token，建议使用密码管理器保存。
        </Text>
      </Modal>
    </div>
  )
}

export default UserProfile
