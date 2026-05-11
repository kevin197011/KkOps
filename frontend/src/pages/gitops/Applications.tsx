// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { Button, Card, Select, Space, Table, Tag, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { CloudSyncOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import { gitopsApi, ApplicationSummary } from '@/api/gitops'

const Applications = () => {
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [integrationId, setIntegrationId] = useState<number | undefined>()
  const [apps, setApps] = useState<ApplicationSummary[]>([])
  const [loading, setLoading] = useState(false)
  const [syncing, setSyncing] = useState<string | null>(null)

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      const list = (res.data.data || []).filter((i) => i.kind.toLowerCase() === 'argocd')
      setIntegrations(list)
      setIntegrationId((prev) => prev ?? list[0]?.id)
    } catch {
      message.error('加载集成失败')
    }
  }, [])

  const loadApps = useCallback(async () => {
    if (!integrationId) {
      setApps([])
      return
    }
    setLoading(true)
    try {
      const res = await gitopsApi.listApplications(integrationId)
      setApps(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载应用失败')
    } finally {
      setLoading(false)
    }
  }, [integrationId])

  useEffect(() => {
    loadIntegrations()
  }, [loadIntegrations])

  useEffect(() => {
    loadApps()
  }, [loadApps])

  const sync = async (name: string) => {
    if (!integrationId) return
    setSyncing(name)
    try {
      await gitopsApi.sync(name, integrationId)
      message.success('已请求同步')
      loadApps()
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '同步失败')
    } finally {
      setSyncing(null)
    }
  }

  const columns: ColumnsType<ApplicationSummary> = [
    { title: '应用', dataIndex: 'name', key: 'name' },
    {
      title: '健康',
      dataIndex: 'health',
      key: 'health',
      render: (h: string) => <Tag color={h === 'Healthy' ? 'green' : 'orange'}>{h}</Tag>,
    },
    {
      title: '同步',
      dataIndex: 'sync_status',
      key: 'sync_status',
      render: (s: string) => <Tag>{s}</Tag>,
    },
    {
      title: '操作',
      key: 'act',
      render: (_, row) => (
        <Button
          size="small"
          type="primary"
          ghost
          icon={<CloudSyncOutlined />}
          loading={syncing === row.name}
          onClick={() => sync(row.name)}
        >
          同步
        </Button>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <CloudSyncOutlined style={{ marginRight: 8 }} />
            GitOps 应用
          </>
        }
        description="查看 Argo CD Application 状态并触发同步。"
      />
      <Card style={{ marginBottom: 16 }}>
        <Space wrap>
          <span>集成：</span>
          <Select
            style={{ minWidth: 260 }}
            value={integrationId}
            onChange={setIntegrationId}
            options={integrations.map((i) => ({
              value: i.id,
              label: `${i.name} (${i.kind})`,
            }))}
          />
          <Button type="primary" onClick={loadApps} loading={loading}>
            刷新
          </Button>
        </Space>
      </Card>

      {integrations.length === 0 ? (
        <EmptyState description="未配置 Argo CD。请在集成中心添加 argocd 类型集成。" />
      ) : (
        <Card title="Applications">
          <Table rowKey="name" loading={loading} columns={columns} dataSource={apps} pagination={{ pageSize: 20 }} />
        </Card>
      )}
    </PageContainer>
  )
}

export default Applications
