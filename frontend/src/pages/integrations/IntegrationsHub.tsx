// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Button,
  Card,
  Col,
  Form,
  Input,
  Modal,
  Row,
  Select,
  Space,
  Switch,
  Table,
  Tag,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import {
  ApiOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExperimentOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons'
import {
  PageContainer,
  PageHeader,
  EmptyState,
} from '@/components/shell'
import {
  integrationsApi,
  IntegrationPublic,
  CreateIntegrationRequest,
} from '@/api/integration'
import { usePermissionStore } from '@/stores/permission'

const KIND_CATALOG = [
  { kind: 'prometheus', title: 'Prometheus', description: '指标与告警查询入口' },
  { kind: 'nightingale', title: 'Nightingale', description: '夜莺 / FlashDuty 监控' },
  { kind: 'loki', title: 'Loki', description: '日志聚合' },
  { kind: 'elasticsearch', title: 'Elasticsearch', description: '日志与检索' },
  { kind: 'grafana', title: 'Grafana', description: '可视化与仪表盘' },
  { kind: 'jenkins', title: 'Jenkins', description: 'CI 流水线' },
  { kind: 'gitlab', title: 'GitLab', description: '代码与 CI/CD' },
  { kind: 'harbor', title: 'Harbor', description: '镜像仓库' },
  { kind: 'argocd', title: 'Argo CD', description: 'GitOps 持续交付' },
] as const

const IntegrationsHub = () => {
  const { hasPermission } = usePermissionStore()
  const canManage = hasPermission('integrations', '*')
  const [rows, setRows] = useState<IntegrationPublic[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<IntegrationPublic | null>(null)
  const [testingId, setTestingId] = useState<number | null>(null)
  const [form] = Form.useForm<CreateIntegrationRequest & { base_url?: string; token?: string; project_id?: number }>()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await integrationsApi.list()
      setRows(res.data.data || [])
    } catch {
      message.error('加载集成列表失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const byKind = useMemo(() => {
    const m: Record<string, IntegrationPublic[]> = {}
    for (const r of rows) {
      m[r.kind] = m[r.kind] || []
      m[r.kind].push(r)
    }
    return m
  }, [rows])

  const openCreate = (kind?: string) => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({
      kind: kind || 'prometheus',
      enabled: true,
      config: {},
      base_url: '',
      token: '',
    })
    setModalOpen(true)
  }

  const openEdit = (record: IntegrationPublic) => {
    setEditing(record)
    form.setFieldsValue({
      name: record.name,
      kind: record.kind,
      enabled: record.enabled,
      description: record.description,
      base_url: '',
      token: '',
      project_id: undefined,
    })
    setModalOpen(true)
  }

  const submit = async () => {
    try {
      const v = await form.validateFields()
      const cfg: Record<string, unknown> = {
        base_url: v.base_url,
        token: v.token,
      }
      if (v.project_id != null) {
        cfg.project_id = v.project_id
      }
      const payload: CreateIntegrationRequest = {
        name: v.name,
        kind: v.kind,
        enabled: v.enabled,
        description: v.description,
        config: cfg,
      }
      if (editing) {
        await integrationsApi.update(editing.id, payload)
        message.success('已更新集成')
      } else {
        await integrationsApi.create(payload)
        message.success('已创建集成')
      }
      setModalOpen(false)
      load()
    } catch (e: unknown) {
      const err = e as { errorFields?: unknown }
      if (!err.errorFields) {
        message.error('保存失败')
      }
    }
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '删除集成',
      content: '确定删除该集成配置？',
      okText: '删除',
      okButtonProps: { danger: true },
      onOk: async () => {
        await integrationsApi.delete(id)
        message.success('已删除')
        load()
      },
    })
  }

  const handleTest = async (id: number) => {
    setTestingId(id)
    try {
      const res = await integrationsApi.test(id)
      if (res.data.ok) {
        message.success('连接成功')
      } else {
        message.error(res.data.error || '连接失败')
      }
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '测试失败')
    } finally {
      setTestingId(null)
    }
  }

  const columns: ColumnsType<IntegrationPublic> = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '类型',
      dataIndex: 'kind',
      key: 'kind',
      render: (k: string) => <Tag color="blue">{k}</Tag>,
    },
    {
      title: '凭据',
      dataIndex: 'has_config',
      key: 'has_config',
      render: (h: boolean) =>
        h ? <Tag icon={<CheckCircleOutlined />} color="success">已配置</Tag> : <Tag>未配置</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (en: boolean) => (en ? <Tag color="green">启用</Tag> : <Tag>停用</Tag>),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<ExperimentOutlined />}
            loading={testingId === record.id}
            onClick={() => handleTest(record.id)}
          >
            测试连接
          </Button>
          {canManage && (
            <>
              <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
                编辑
              </Button>
              <Button type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
                删除
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <ApiOutlined style={{ marginRight: 8 }} />
            集成中心
          </>
        }
        description="配置外部系统连接（凭据加密存储）。点击下方卡片查看该类型的已接入实例，或在表格中统一管理。"
        extra={
          canManage ? (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => openCreate()}>
              添加集成
            </Button>
          ) : undefined
        }
      />

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {KIND_CATALOG.map((k) => {
          const connected = byKind[k.kind]?.filter((r) => r.enabled && r.has_config).length || 0
          const anyConnected = connected > 0
          return (
            <Col xs={24} sm={12} md={8} lg={6} key={k.kind}>
              <Card size="small" title={k.title}>
                <div style={{ marginBottom: 8, color: 'var(--ant-color-text-secondary)' }}>{k.description}</div>
                <Space direction="vertical" style={{ width: '100%' }}>
                  <div>
                    {anyConnected ? (
                      <Tag icon={<CheckCircleOutlined />} color="success">
                        已连接 {connected}
                      </Tag>
                    ) : (
                      <Tag icon={<CloseCircleOutlined />} color="default">
                        未配置
                      </Tag>
                    )}
                  </div>
                  {canManage && (
                    <Button size="small" type="primary" ghost onClick={() => openCreate(k.kind)}>
                      添加
                    </Button>
                  )}
                </Space>
              </Card>
            </Col>
          )
        })}
      </Row>

      {rows.length === 0 && !loading ? (
        <EmptyState description="暂无集成实例。请先添加外部系统集成。" />
      ) : (
        <Card title="集成实例">
          <Table rowKey="id" loading={loading} columns={columns} dataSource={rows} pagination={{ pageSize: 10 }} />
        </Card>
      )}

      <Modal
        title={editing ? '编辑集成' : '添加集成'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={submit}
        okText="保存"
        width={560}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="例如：生产 Prometheus" />
          </Form.Item>
          <Form.Item name="kind" label="类型" rules={[{ required: true }]}>
            <Select
              disabled={!!editing}
              options={KIND_CATALOG.map((x) => ({ value: x.kind, label: x.title }))}
            />
          </Form.Item>
          <Form.Item name="base_url" label="Base URL" rules={[{ required: true, message: '请输入服务地址' }]}>
            <Input placeholder="https://prometheus.example.com" />
          </Form.Item>
          <Form.Item name="token" label="Token / API Key" rules={[{ required: true, message: '请输入令牌' }]}>
            <Input.Password placeholder="Bearer / Personal Access Token" />
          </Form.Item>
          <Form.Item name="project_id" label="GitLab 项目 ID（仅 GitLab CI 流水线功能需要）">
            <Input type="number" placeholder="可选，例如 12345" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  )
}

export default IntegrationsHub
