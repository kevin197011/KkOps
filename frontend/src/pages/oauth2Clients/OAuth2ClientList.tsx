// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import {
  Card,
  Button,
  Space,
  message,
  Modal,
  Form,
  Input,
  Switch,
  Table,
  Tag,
  Typography,
  Empty,
  Select,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SafetyCertificateOutlined,
  CopyOutlined,
} from '@ant-design/icons'
import {
  oauth2ClientApi,
  OAuth2Client,
  CreateOAuth2ClientRequest,
  UpdateOAuth2ClientRequest,
} from '@/api/oauth2Client'

const { Title, Text } = Typography

const PROTOCOL_OPTIONS = [
  { value: '', label: '全部协议' },
  { value: 'oidc', label: 'OIDC' },
  { value: 'saml', label: 'SAML 2.0' },
  { value: 'ldap', label: 'LDAP' },
]

const OAuth2ClientList = () => {
  const [list, setList] = useState<OAuth2Client[]>([])
  const [loading, setLoading] = useState(false)
  const [protocolFilter, setProtocolFilter] = useState<string>('')
  const [modalVisible, setModalVisible] = useState(false)
  const [editing, setEditing] = useState<OAuth2Client | null>(null)
  const [secretModalVisible, setSecretModalVisible] = useState(false)
  const [secretValue, setSecretValue] = useState('')
  const [form] = Form.useForm()

  const fetchList = useCallback(async () => {
    setLoading(true)
    try {
      const res = await oauth2ClientApi.list(
        protocolFilter ? { protocol: protocolFilter } : undefined
      )
      setList(res.data.data || [])
    } catch {
      message.error('获取 IdP 应用列表失败')
    } finally {
      setLoading(false)
    }
  }, [protocolFilter])

  useEffect(() => {
    fetchList()
  }, [fetchList])

  const handleCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({
      protocol: 'oidc',
      redirect_uris: [''],
      scopes: 'openid profile email',
      enabled: true,
    })
    setModalVisible(true)
  }

  const handleEdit = (record: OAuth2Client) => {
    setEditing(record)
    form.setFieldsValue({
      name: record.name,
      protocol: record.protocol || 'oidc',
      redirect_uris: record.redirect_uris?.length
        ? record.redirect_uris
        : [''],
      scopes: record.scopes || 'openid profile email',
      enabled: record.enabled,
    })
    setModalVisible(true)
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除该 IdP 应用吗？删除后依赖该凭据的系统将无法使用 KkOps 登录。',
      onOk: async () => {
        try {
          await oauth2ClientApi.delete(id)
          message.success('删除成功')
          fetchList()
        } catch (e: any) {
          message.error(e.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleRegenerateSecret = (id: number) => {
    Modal.confirm({
      title: '重新生成密钥',
      content: '重新生成后，旧密钥将立即失效。新密钥仅显示一次，请妥善保存。',
      onOk: async () => {
        try {
          const res = await oauth2ClientApi.regenerateSecret(id)
          const secret = res.data.data?.client_secret
          if (secret) {
            setSecretValue(secret)
            setSecretModalVisible(true)
            message.success('已生成新密钥，请复制保存')
          }
        } catch (e: any) {
          message.error(e.response?.data?.error || '操作失败')
        }
      },
    })
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const redirect_uris = (values.redirect_uris as string[]).filter(
        (u: string) => u?.trim()
      )
      if (redirect_uris.length === 0) {
        message.error('请至少填写一个回调地址')
        return
      }
      if (editing) {
        const payload: UpdateOAuth2ClientRequest = {
          name: values.name,
          protocol: values.protocol,
          redirect_uris,
          scopes: values.scopes?.trim() || 'openid profile email',
          enabled: values.enabled ?? true,
        }
        await oauth2ClientApi.update(editing.id, payload)
        message.success('更新成功')
      } else {
        const payload: CreateOAuth2ClientRequest = {
          name: values.name,
          protocol: values.protocol || 'oidc',
          redirect_uris,
          scopes: values.scopes?.trim() || 'openid profile email',
        }
        const res = await oauth2ClientApi.create(payload)
        const created = res.data.data
        if (created?.client_secret) {
          setSecretValue(created.client_secret)
          setSecretModalVisible(true)
          message.success('创建成功，请复制并保存 client_secret（仅显示一次）')
        } else {
          message.success('创建成功')
        }
      }
      setModalVisible(false)
      form.resetFields()
      fetchList()
    } catch (e: any) {
      if (e.errorFields) return
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const copySecret = () => {
    if (secretValue) {
      navigator.clipboard.writeText(secretValue)
      message.success('已复制到剪贴板')
    }
  }

  return (
    <div style={{ padding: 24 }}>
      <Title level={4} style={{ marginBottom: 24 }}>
        <SafetyCertificateOutlined style={{ marginRight: 8 }} />
        IdP 应用管理
      </Title>
      <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
        登记依赖 KkOps 身份提供商的应用（OIDC 为主；SAML/LDAP 协议字段用于路由与未来扩展）。GitLab、Jenkins、Grafana 等可配置 KkOps 为认证来源，使用此处创建的 client_id / client_secret。
      </Text>

      <div style={{ marginBottom: 16, display: 'flex', gap: 12, flexWrap: 'wrap', alignItems: 'center' }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          添加 IdP 应用
        </Button>
        <span style={{ color: 'var(--ant-color-text-secondary)' }}>协议筛选</span>
        <Select
          style={{ width: 160 }}
          value={protocolFilter}
          onChange={(v) => setProtocolFilter(v)}
          options={PROTOCOL_OPTIONS}
        />
      </div>
      <Card title="应用列表" size="small">
        {list.length === 0 && !loading ? (
          <Empty description="暂无 IdP 应用，点击「添加 IdP 应用」创建" />
        ) : (
          <Table
            rowKey="id"
            loading={loading}
            dataSource={list}
            size="small"
            columns={[
              { title: 'ID', dataIndex: 'id', width: 70 },
              {
                title: 'Client ID',
                dataIndex: 'client_id',
                ellipsis: true,
                render: (t: string) => <code style={{ fontSize: 12 }}>{t}</code>,
              },
              {
                title: '名称',
                dataIndex: 'name',
                render: (t: string) => <strong>{t}</strong>,
              },
              {
                title: '协议',
                dataIndex: 'protocol',
                width: 100,
                render: (p: string) => (
                  <Tag>{p === 'saml' ? 'SAML' : p === 'ldap' ? 'LDAP' : 'OIDC'}</Tag>
                ),
              },
              {
                title: '回调地址',
                dataIndex: 'redirect_uris',
                ellipsis: true,
                render: (uris: string[]) =>
                  Array.isArray(uris) ? uris.join(', ') : '-',
              },
              {
                title: 'Scopes',
                dataIndex: 'scopes',
                width: 140,
                ellipsis: true,
              },
              {
                title: '状态',
                dataIndex: 'enabled',
                width: 80,
                render: (en: boolean) => (
                  <Tag color={en ? 'green' : 'default'}>
                    {en ? '启用' : '禁用'}
                  </Tag>
                ),
              },
              {
                title: '操作',
                key: 'action',
                width: 220,
                render: (_: unknown, r: OAuth2Client) => (
                  <Space>
                    <Button
                      type="link"
                      size="small"
                      icon={<EditOutlined />}
                      onClick={() => handleEdit(r)}
                    >
                      编辑
                    </Button>
                    <Button
                      type="link"
                      size="small"
                      icon={<SafetyCertificateOutlined />}
                      onClick={() => handleRegenerateSecret(r.id)}
                    >
                      重置密钥
                    </Button>
                    <Button
                      type="link"
                      danger
                      size="small"
                      icon={<DeleteOutlined />}
                      onClick={() => handleDelete(r.id)}
                    >
                      删除
                    </Button>
                  </Space>
                ),
              },
            ]}
            pagination={false}
          />
        )}
      </Card>

      <Modal
        title={editing ? '编辑 IdP 应用' : '添加 IdP 应用'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        width={560}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label="应用名称" rules={[{ required: true }]}>
            <Input placeholder="如：GitLab、Jenkins、Grafana" />
          </Form.Item>
          <Form.Item
            name="protocol"
            label="协议"
            rules={[{ required: true }]}
            extra="OIDC 用于标准 OAuth2/OIDC 集成；SAML/LDAP 为预留登记，具体端点见运维文档。"
          >
            <Select
              options={[
                { value: 'oidc', label: 'OIDC（OAuth2）' },
                { value: 'saml', label: 'SAML 2.0' },
                { value: 'ldap', label: 'LDAP' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="redirect_uris"
            label="回调地址 (redirect_uri)"
            rules={[
              {
                validator: (_, val) => {
                  const uris = (val || []).filter((u: string) => u?.trim())
                  if (uris.length === 0) return Promise.reject(new Error('至少填写一个'))
                  return Promise.resolve()
                },
              },
            ]}
            extra="OIDC 授权后跳转的地址，需与第三方系统配置完全一致"
          >
            <Form.List name="redirect_uris">
              {(fields, { add, remove }) => (
                <>
                  {fields.map(({ key, name, ...rest }) => (
                    <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                      <Form.Item
                        {...rest}
                        name={[name]}
                        rules={[{ required: true, message: '请填写' }]}
                        noStyle
                      >
                        <Input placeholder="https://app.example.com/callback" style={{ width: 360 }} />
                      </Form.Item>
                      <Button type="text" danger onClick={() => remove(name)}>
                        删除
                      </Button>
                    </Space>
                  ))}
                  <Button type="dashed" onClick={() => add('')} block>
                    + 添加回调地址
                  </Button>
                </>
              )}
            </Form.List>
          </Form.Item>
          <Form.Item
            name="scopes"
            label="Scopes"
            extra="空格分隔，如 openid profile email"
          >
            <Input placeholder="openid profile email" />
          </Form.Item>
          <Form.Item
            name="enabled"
            label="启用"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="请保存 Client Secret"
        open={secretModalVisible}
        onCancel={() => {
          setSecretModalVisible(false)
          setSecretValue('')
        }}
        footer={[
          <Button key="copy" type="primary" icon={<CopyOutlined />} onClick={copySecret}>
            复制
          </Button>,
          <Button
            key="close"
            onClick={() => {
              setSecretModalVisible(false)
              setSecretValue('')
            }}
          >
            已保存，关闭
          </Button>,
        ]}
        destroyOnClose
      >
        <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
          此密钥仅显示一次，请复制保存。关闭后将无法再次查看。
        </Text>
        <Input.Password
          readOnly
          value={secretValue}
          style={{ fontFamily: 'monospace' }}
        />
      </Modal>
    </div>
  )
}

export default OAuth2ClientList
