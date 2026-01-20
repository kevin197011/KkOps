// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Card, Table, Button, Space, message, Modal, Form, Input, DatePicker, Tag, Popconfirm, Typography } from 'antd'
import { PlusOutlined, DeleteOutlined, CopyOutlined, KeyOutlined, UserOutlined, MailOutlined, SafetyOutlined } from '@ant-design/icons'
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
  const [newToken, setNewToken] = useState<string>('')
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
      
      setNewToken(response.data.token || '')
      setTokenModalVisible(true)
      setCreateModalVisible(false)
      form.resetFields()
      fetchTokens()
      message.success('Token 创建成功')
    } catch (error: any) {
      message.error('创建 Token 失败: ' + (error.response?.data?.error || error.message))
    }
  }

  const handleRevokeToken = async (tokenId: number) => {
    try {
      await userApi.revokeToken(tokenId)
      message.success('Token 已撤销')
      fetchTokens()
    } catch (error: any) {
      message.error('撤销 Token 失败: ' + (error.response?.data?.error || error.message))
    }
  }

  const handleCopyToken = (token: string) => {
    navigator.clipboard.writeText(token)
    message.success('Token 已复制到剪贴板')
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
          <Popconfirm
            title="确定要撤销此 Token 吗？"
            description="撤销后该 Token 将无法再使用"
            onConfirm={() => handleRevokeToken(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
            >
              撤销
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

      {/* 显示新创建的 Token 弹窗 */}
      <Modal
        title={
          <Space>
            <SafetyOutlined style={{ color: '#52c41a' }} />
            <span>Token 创建成功</span>
          </Space>
        }
        open={tokenModalVisible}
        onCancel={() => {
          setTokenModalVisible(false)
          setNewToken('')
        }}
        footer={[
          <Button
            key="copy"
            type="primary"
            icon={<CopyOutlined />}
            onClick={() => handleCopyToken(newToken)}
          >
            复制 Token
          </Button>,
          <Button
            key="close"
            onClick={() => {
              setTokenModalVisible(false)
              setNewToken('')
            }}
          >
            我已保存
          </Button>,
        ]}
        destroyOnClose
      >
        <div style={{ marginBottom: 16 }}>
          <Text type="warning" strong>
            ⚠️ 请立即保存此 Token，关闭窗口后将无法再次查看完整 Token
          </Text>
        </div>
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
            {newToken}
          </Text>
        </div>
        <Text type="secondary" style={{ fontSize: 12 }}>
          提示：Token 仅显示一次，请妥善保管。建议使用密码管理器保存。
        </Text>
      </Modal>
    </div>
  )
}

export default UserProfile
