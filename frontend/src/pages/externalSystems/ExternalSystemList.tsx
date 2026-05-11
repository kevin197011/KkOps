// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  message,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  Table,
  Tag,
  Typography,
  Empty,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  LinkOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons'
import {
  externalSystemApi,
  ExternalSystem,
  CreateExternalSystemRequest,
  LaunchType,
} from '@/api/externalSystem'
import { usePermissionStore } from '@/stores/permission'

const { Title, Text } = Typography
const { TextArea } = Input

const ExternalSystemList = () => {
  const { hasPermission } = usePermissionStore()
  const canManage = hasPermission('external-systems', '*')
  const [list, setList] = useState<ExternalSystem[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editing, setEditing] = useState<ExternalSystem | null>(null)
  const [launchingId, setLaunchingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const launchType = Form.useWatch('launch_type', form) as LaunchType | undefined

  const fetchList = async () => {
    setLoading(true)
    try {
      const res = await externalSystemApi.list(false)
      setList(res.data.data || [])
    } catch {
      message.error('获取外部系统列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchList()
  }, [])

  const handleLaunch = async (id: number) => {
    setLaunchingId(id)
    try {
      const res = await externalSystemApi.launch(id)
      const url = res.data.redirect_url
      if (url) {
        window.open(url, '_blank', 'noopener,noreferrer')
        message.success('已打开新窗口，请在新窗口中完成登录')
      }
    } catch (e: any) {
      message.error(e.response?.data?.error || '跳转失败')
    } finally {
      setLaunchingId(null)
    }
  }

  const handleCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({
      enabled: true,
      order_index: 0,
      launch_type: 'sso_link',
      login_path: '/sso/consume',
    })
    setModalVisible(true)
  }

  const handleEdit = (record: ExternalSystem) => {
    setEditing(record)
    form.setFieldsValue({
      name: record.name,
      description: record.description,
      launch_type: record.launch_type || 'sso_link',
      base_url: record.base_url,
      login_path: record.login_path,
      role_mapping: record.role_mapping,
      icon: record.icon,
      order_index: record.order_index,
      enabled: record.enabled,
    })
    setModalVisible(true)
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除该外部系统配置吗？',
      onOk: async () => {
        try {
          await externalSystemApi.delete(id)
          message.success('删除成功')
          fetchList()
        } catch (e: any) {
          message.error(e.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const launchType = (values.launch_type as LaunchType) || 'sso_link'
      const payload = {
        name: values.name,
        description: values.description,
        launch_type: launchType,
        base_url: values.base_url,
        login_path: launchType === 'jwt_token' ? values.login_path : undefined,
        secret: launchType === 'jwt_token' ? values.secret : undefined,
        role_mapping: launchType === 'jwt_token' ? values.role_mapping : undefined,
        icon: values.icon,
        order_index: values.order_index ?? 0,
        enabled: values.enabled ?? true,
      }
      if (editing) {
        await externalSystemApi.update(editing.id, payload)
        message.success('更新成功')
      } else {
        await externalSystemApi.create({
          ...payload,
          secret: payload.secret || (launchType === 'jwt_token' ? '' : undefined),
        } as CreateExternalSystemRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchList()
    } catch (e: any) {
      if (e.errorFields) return
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const enabledList = list.filter((x) => x.enabled)

  return (
    <div style={{ padding: 24 }}>
      <Title level={4} style={{ marginBottom: 24 }}>
        <SafetyCertificateOutlined style={{ marginRight: 8 }} />
        SSO 应用门户
      </Title>
      <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
        用户登录一次（SSO）后，在此一键打开 GitLab、Jenkins、Grafana 等接入系统。若接入系统与 KkOps 共用同一 IdP（如 Keycloak），选择「同 SSO 应用」打开即已登录；也可配置「Token 跳转」向目标系统附带身份与权限。
      </Text>

      {enabledList.length > 0 && (
        <Card title="快捷打开" size="small" style={{ marginBottom: 24 }}>
          <Row gutter={[16, 16]}>
            {enabledList.map((sys) => (
              <Col key={sys.id} xs={24} sm={12} md={8} lg={6}>
                <Card
                  size="small"
                  actions={[
                    <Button
                      key="open"
                      type="primary"
                      icon={<LinkOutlined />}
                      loading={launchingId === sys.id}
                      onClick={() => handleLaunch(sys.id)}
                    >
                      打开
                    </Button>,
                  ]}
                >
                  <Card.Meta
                    title={sys.name}
                    description={sys.description || sys.base_url}
                  />
                </Card>
              </Col>
            ))}
          </Row>
        </Card>
      )}

      {canManage && (
        <>
          <div style={{ marginBottom: 16 }}>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              添加外部系统
            </Button>
          </div>
          <Card title="配置列表" size="small">
            {list.length === 0 && !loading ? (
              <Empty description="暂无外部系统，点击「添加外部系统」配置" />
            ) : (
              <Table
                rowKey="id"
                loading={loading}
                dataSource={list}
                size="small"
                columns={[
                  { title: 'ID', dataIndex: 'id', width: 70 },
                  { title: '名称', dataIndex: 'name', render: (t: string) => <strong>{t}</strong> },
                  {
                    title: '类型',
                    dataIndex: 'launch_type',
                    width: 100,
                    render: (t: LaunchType) => (
                      <Tag color={t === 'sso_link' ? 'blue' : 'orange'}>
                        {t === 'sso_link' ? '同 SSO 应用' : 'Token 跳转'}
                      </Tag>
                    ),
                  },
                  {
                    title: '地址',
                    key: 'url',
                    ellipsis: true,
                    render: (_: unknown, r: ExternalSystem) =>
                      r.launch_type === 'sso_link'
                        ? r.base_url
                        : `${r.base_url}${r.login_path || ''}`,
                  },
                  {
                    title: '状态',
                    dataIndex: 'enabled',
                    width: 80,
                    render: (en: boolean) => (
                      <Tag color={en ? 'green' : 'default'}>{en ? '启用' : '禁用'}</Tag>
                    ),
                  },
                  {
                    title: '操作',
                    key: 'action',
                    width: 180,
                    render: (_: unknown, r: ExternalSystem) => (
                      <Space>
                        <Button
                          type="link"
                          size="small"
                          icon={<LinkOutlined />}
                          loading={launchingId === r.id}
                          onClick={() => handleLaunch(r.id)}
                        >
                          打开
                        </Button>
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
        </>
      )}

      {!canManage && list.length === 0 && !loading && (
        <Empty description="暂无可用的外部系统" />
      )}

      <Modal
        title={editing ? '编辑 SSO 应用' : '添加 SSO 应用'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => { setModalVisible(false); form.resetFields() }}
        width={560}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input placeholder="如：GitLab、Jenkins、Grafana" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input placeholder="简要说明" />
          </Form.Item>
          <Form.Item
            name="launch_type"
            label="打开方式"
            rules={[{ required: true }]}
            extra="同 SSO 应用：与 KkOps 共用同一 IdP（如 Keycloak），打开即已登录；Token 跳转：目标系统接收 KkOps 签发的 JWT"
          >
            <Select
              options={[
                { value: 'sso_link', label: '同 SSO 应用（仅链接，打开即已登录）' },
                { value: 'jwt_token', label: 'Token 跳转（附带身份与权限）' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="base_url"
            label={launchType === 'sso_link' ? '应用 URL' : '基础 URL'}
            rules={[
              { required: true },
              {
                pattern: /^https?:\/\//,
                message: '请填写以 http:// 或 https:// 开头的 URL',
              },
            ]}
            extra={
              launchType === 'sso_link'
                ? '用户点击打开时直接访问的完整地址，如 https://gitlab.example.com'
                : '目标系统根地址，如 https://jumpserver.example.com'
            }
          >
            <Input
              placeholder={
                launchType === 'sso_link'
                  ? 'https://gitlab.example.com'
                  : 'https://jumpserver.example.com'
              }
            />
          </Form.Item>
          {launchType === 'jwt_token' && (
            <>
              <Form.Item
                name="login_path"
                label="登录路径"
                rules={[{ required: true }]}
                extra="接收 token 的路径，如 /api/sso/consume"
              >
                <Input placeholder="/sso/consume" />
              </Form.Item>
              <Form.Item
                name="secret"
                label="共享密钥"
                rules={editing ? [] : [{ required: true, message: '创建时必填' }]}
                extra="与目标系统约定的密钥，用于签名 JWT"
              >
                <Input.Password placeholder={editing ? '留空不修改' : '与目标系统一致'} />
              </Form.Item>
              <Form.Item
                name="role_mapping"
                label="角色映射 (JSON)"
                extra='如 {"admin":"administrator"}，将 KkOps 角色映射为目标系统角色'
              >
                <TextArea rows={2} placeholder='{"admin":"administrator"}' />
              </Form.Item>
            </>
          )}
          <Form.Item name="icon" label="图标">
            <Input placeholder="可选" />
          </Form.Item>
          <Form.Item name="order_index" label="排序" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked" initialValue={true}>
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ExternalSystemList
