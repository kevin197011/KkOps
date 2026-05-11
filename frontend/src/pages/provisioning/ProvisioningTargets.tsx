// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import {
  Button,
  Card,
  Form,
  Modal,
  Select,
  Switch,
  Table,
  Tag,
  message,
} from 'antd'
import { PlusOutlined, SyncOutlined, CloudSyncOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import {
  PageContainer,
  PageHeader,
  EmptyState,
} from '@/components/shell'
import {
  provisioningApi,
  integrationsApi,
  ProvisioningTarget,
  ProvisioningRun,
  CreateProvisioningTargetRequest,
} from '@/api/provisioning'

const PROVIDER_KINDS = [
  { value: 'scim', label: 'SCIM 2.0 (generic)' },
  { value: 'gitlab', label: 'GitLab' },
  { value: 'jenkins', label: 'Jenkins' },
  { value: 'grafana', label: 'Grafana' },
  { value: 'harbor', label: 'Harbor' },
  { value: 'argocd', label: 'Argo CD' },
  { value: 'jumpserver', label: 'Jumpserver' },
  { value: 'nightingale', label: 'Nightingale (n9e)' },
]

const ProvisioningTargets = () => {
  const [targets, setTargets] = useState<ProvisioningTarget[]>([])
  const [integrations, setIntegrations] = useState<
    { id: number; name: string; kind: string }[]
  >([])
  const [loading, setLoading] = useState(false)
  const [runsLoading, setRunsLoading] = useState(false)
  const [runsByTarget, setRunsByTarget] = useState<Record<number, ProvisioningRun[]>>({})
  const [modalOpen, setModalOpen] = useState(false)
  const [syncingId, setSyncingId] = useState<number | null>(null)
  const [form] = Form.useForm<CreateProvisioningTargetRequest>()

  const loadTargets = useCallback(async () => {
    setLoading(true)
    try {
      const res = await provisioningApi.listTargets()
      setTargets(res.data.data || [])
    } catch {
      message.error('加载账号同步目标失败')
    } finally {
      setLoading(false)
    }
  }, [])

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      const rows = (res.data.data || []).map((i) => ({
        id: i.id,
        name: i.name,
        kind: i.kind,
      }))
      setIntegrations(rows)
    } catch {
      message.error('加载外部集成列表失败')
    }
  }, [])

  useEffect(() => {
    loadTargets()
    loadIntegrations()
  }, [loadTargets, loadIntegrations])

  const fetchRunsForTarget = async (targetId: number) => {
    setRunsLoading(true)
    try {
      const res = await provisioningApi.listRuns(targetId, 50)
      setRunsByTarget((prev) => ({
        ...prev,
        [targetId]: res.data.data || [],
      }))
    } catch {
      message.error('加载同步记录失败')
    } finally {
      setRunsLoading(false)
    }
  }

  const handleExpand = async (expanded: boolean, record: ProvisioningTarget) => {
    if (expanded && !runsByTarget[record.id]) {
      await fetchRunsForTarget(record.id)
    }
  }

  const handleSync = async (id: number) => {
    setSyncingId(id)
    try {
      await provisioningApi.syncTarget(id)
      message.success('已加入同步队列')
      await fetchRunsForTarget(id)
    } catch {
      message.error('触发同步失败')
    } finally {
      setSyncingId(null)
    }
  }

  const submitCreate = async () => {
    try {
      const values = await form.validateFields()
      await provisioningApi.createTarget(values)
      message.success('已创建目标')
      setModalOpen(false)
      form.resetFields()
      loadTargets()
    } catch (e: unknown) {
      if ((e as { errorFields?: unknown }).errorFields) return
      message.error('创建失败')
    }
  }

  const columns: ColumnsType<ProvisioningTarget> = [
    { title: 'ID', dataIndex: 'id', width: 64 },
    {
      title: '外部集成',
      key: 'integ',
      render: (_, r) =>
        r.integration ? (
          <span>
            <strong>{r.integration.name}</strong>{' '}
            <Tag>{r.integration.kind}</Tag>
          </span>
        ) : (
          <Tag>integration #{r.integration_id}</Tag>
        ),
    },
    {
      title: '下发类型',
      dataIndex: 'provider_kind',
      render: (k: string) => <code>{k}</code>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (s: string) => <Tag>{s}</Tag>,
    },
    {
      title: '上次同步',
      dataIndex: 'last_sync_at',
      width: 180,
      render: (t?: string) => t || '—',
    },
    {
      title: '启用',
      dataIndex: 'enabled',
      width: 72,
      render: (en: boolean) => <Tag color={en ? 'green' : 'default'}>{en ? '是' : '否'}</Tag>,
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      render: (_, r) => (
        <Button
          type="primary"
          size="small"
          icon={<SyncOutlined />}
          loading={syncingId === r.id}
          onClick={() => handleSync(r.id)}
        >
          立即同步
        </Button>
      ),
    },
  ]

  const runColumns: ColumnsType<ProvisioningRun> = [
    { title: '时间', dataIndex: 'started_at', width: 180 },
    { title: '动作', dataIndex: 'action', width: 88 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 96,
      render: (s: string) => <Tag>{s}</Tag>,
    },
    {
      title: '说明',
      dataIndex: 'message',
      ellipsis: true,
      render: (m?: string) => m || '—',
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <CloudSyncOutlined style={{ marginRight: 8 }} />
            账号同步
          </>
        }
        description="将用户账号同步到外部系统（依赖「外部集成」中的连接与凭据）。手动同步与后台自动同步共用同一队列。"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新建目标
          </Button>
        }
      />

      {!loading && targets.length === 0 ? (
        <EmptyState description="暂无账号同步目标。请先配置外部集成，再在此绑定下发类型。" />
      ) : (
        <Card size="small">
          <Table<ProvisioningTarget>
            rowKey="id"
            loading={loading}
            columns={columns}
            dataSource={targets}
            pagination={false}
            expandable={{
              expandedRowRender: (record) => (
                <Table<ProvisioningRun>
                  size="small"
                  rowKey="id"
                  loading={runsLoading}
                  columns={runColumns}
                  dataSource={runsByTarget[record.id] || []}
                  pagination={false}
                  locale={{ emptyText: '暂无同步记录' }}
                />
              ),
              onExpand: handleExpand,
            }}
          />
        </Card>
      )}

      <Modal
        title="新建账号同步目标"
        open={modalOpen}
        onOk={submitCreate}
        onCancel={() => {
          setModalOpen(false)
          form.resetFields()
        }}
        destroyOnClose
        width={520}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{ enabled: true }}
        >
          <Form.Item
            name="integration_id"
            label="外部集成"
            rules={[{ required: true, message: '请选择集成' }]}
          >
            <Select
              showSearch
              optionFilterProp="label"
              placeholder="选择已配置的集成"
              options={integrations.map((i) => ({
                value: i.id,
                label: `${i.name} (${i.kind})`,
              }))}
            />
          </Form.Item>
          <Form.Item
            name="provider_kind"
            label="下发类型"
            rules={[{ required: true, message: '请选择类型' }]}
          >
            <Select options={PROVIDER_KINDS} placeholder="选择提供商 / 协议" />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  )
}

export default ProvisioningTargets
