// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useMemo, useState } from 'react'
import { Button, Card, Space, Table, Tag, Typography, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { BellOutlined, ReloadOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { alertsApi, AlertRow } from '@/api/alerts'

function parseLabels(raw?: string): Record<string, string> {
  if (!raw?.trim()) return {}
  try {
    return JSON.parse(raw) as Record<string, string>
  } catch {
    return {}
  }
}

const AlertCenter = () => {
  const [rows, setRows] = useState<AlertRow[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await alertsApi.list({ limit: 100 })
      setRows(res.data.data || [])
      setTotal(res.data.total ?? 0)
    } catch {
      message.error('Failed to load alerts')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const onSync = async () => {
    setLoading(true)
    try {
      const res = await alertsApi.sync()
      message.success(`Synced ${res.data.synced} alert(s)`)
      if (res.data.warnings?.length) {
        console.warn(res.data.warnings)
      }
      await load()
    } catch {
      message.error('Sync failed')
    } finally {
      setLoading(false)
    }
  }

  const columns: ColumnsType<AlertRow> = useMemo(
    () => [
      {
        title: 'Severity',
        dataIndex: 'severity',
        width: 110,
        render: (s: string) => (
          <Tag color={s === 'critical' ? 'red' : s === 'warning' ? 'orange' : 'blue'}>{s}</Tag>
        ),
      },
      {
        title: 'Title',
        dataIndex: 'title',
        ellipsis: true,
      },
      {
        title: 'Source',
        dataIndex: 'source',
        width: 160,
      },
      {
        title: 'Integration',
        dataIndex: 'integration_name',
        width: 160,
        render: (_v, r) => r.integration_name || (r.integration_id ? `#${r.integration_id}` : '—'),
      },
      {
        title: 'Status',
        dataIndex: 'status',
        width: 120,
      },
      {
        title: 'Labels',
        key: 'labels',
        ellipsis: true,
        render: (_v, r) => {
          const lb = parseLabels(r.labels_json)
          const txt = Object.entries(lb)
            .map(([k, v]) => `${k}=${v}`)
            .join(', ')
          return txt || '—'
        },
      },
      {
        title: 'Actions',
        key: 'actions',
        width: 200,
        render: (_v, r) => (
          <Space size="small">
            <Button
              size="small"
              disabled={r.status === 'dismissed' || r.status === 'acknowledged'}
              onClick={async () => {
                try {
                  await alertsApi.acknowledge(r.id)
                  message.success('Acknowledged')
                  load()
                } catch {
                  message.error('Acknowledge failed')
                }
              }}
            >
              Ack
            </Button>
            <Button
              size="small"
              danger
              disabled={r.status === 'dismissed'}
              onClick={async () => {
                try {
                  await alertsApi.dismiss(r.id)
                  message.success('Dismissed')
                  load()
                } catch {
                  message.error('Dismiss failed')
                }
              }}
            >
              Dismiss
            </Button>
          </Space>
        ),
      },
    ],
    [load]
  )

  return (
    <PageContainer>
      <PageHeader
        title={
          <Space>
            <BellOutlined />
            <span>Alerts</span>
          </Space>
        }
        description="Aggregated alerts from integrations and webhooks. Silence rules UI is planned."
        extra={
          <Space>
            <Button icon={<ReloadOutlined />} onClick={() => load()} loading={loading}>
              Refresh
            </Button>
            <Button type="primary" onClick={onSync} loading={loading}>
              Sync from integrations
            </Button>
          </Space>
        }
      />
      <Card styles={{ body: { paddingTop: 16 } }}>
        <Typography.Paragraph type="secondary">
          Configure <Typography.Text code>alertmanager_base_url</Typography.Text> on a Prometheus
          integration to pull Alertmanager <Typography.Text code>/api/v2/alerts</Typography.Text>.
          Webhook: <Typography.Text code>POST /api/v1/alerts/webhook</Typography.Text> with header{' '}
          <Typography.Text code>X-KkOps-Alert-Secret</Typography.Text>.
        </Typography.Paragraph>
        {rows.length === 0 && !loading ? (
          <EmptyState description="No alerts yet. Run sync or send a webhook payload." />
        ) : (
          <Table<AlertRow>
            rowKey="id"
            loading={loading}
            columns={columns}
            dataSource={rows}
            pagination={{ total, pageSize: 100, showTotal: (t) => `${t} alerts` }}
          />
        )}
        <Card size="small" title="Silence rules (stub)" style={{ marginTop: 16 }}>
          <Typography.Text type="secondary">
            Silence creation against Alertmanager will be added in a follow-up change.
          </Typography.Text>
        </Card>
      </Card>
    </PageContainer>
  )
}

export default AlertCenter
