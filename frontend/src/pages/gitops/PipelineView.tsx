// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useMemo, useState } from 'react'
import { Card, Input, Select, Space, Timeline, Typography, message } from 'antd'
import { BranchesOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import { gitopsApi, PipelineEvent } from '@/api/gitops'

const PipelineView = () => {
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [argoId, setArgoId] = useState<number | undefined>()
  const [app, setApp] = useState('')
  const [events, setEvents] = useState<PipelineEvent[]>([])
  const [loading, setLoading] = useState(false)

  const argoList = useMemo(
    () => integrations.filter((i) => i.kind.toLowerCase() === 'argocd'),
    [integrations]
  )

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      const list = res.data.data || []
      setIntegrations(list)
      const firstArgo = list.find((i) => i.kind.toLowerCase() === 'argocd')
      setArgoId((prev) => prev ?? firstArgo?.id)
    } catch {
      message.error('加载集成失败')
    }
  }, [])

  useEffect(() => {
    loadIntegrations()
  }, [loadIntegrations])

  const loadTimeline = async () => {
    const name = app.trim()
    if (!name) {
      message.warning('请输入 Argo CD 应用名称')
      return
    }
    setLoading(true)
    try {
      const res = await gitopsApi.pipelineView({
        app: name,
        argocd_integration_id: argoId,
      })
      setEvents(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载流水线视图失败')
    } finally {
      setLoading(false)
    }
  }

  const colorFor = (kind: string, status: string) => {
    if (kind === 'gitops') return 'green'
    const s = status.toLowerCase()
    if (s.includes('fail') || s.includes('error')) return 'red'
    if (s.includes('success') || s.includes('passed')) return 'green'
    return 'blue'
  }

  return (
    <PageContainer>
      <PageHeader
        title={
          <Space>
            <BranchesOutlined />
            流水线视图
          </Space>
        }
        description="聚合 GitLab CI 最近流水线与 Argo CD 部署历史（按应用名称）。"
      />
      <Card bordered={false}>
        <Space wrap style={{ marginBottom: 16 }}>
          <Select
            placeholder="Argo CD 集成（可选）"
            allowClear
            style={{ minWidth: 220 }}
            value={argoId}
            onChange={(v) => setArgoId(v)}
            options={argoList.map((i) => ({ label: i.name, value: i.id }))}
          />
          <Typography.Text>应用名称（Argo CD Application）</Typography.Text>
          <Input
            style={{ minWidth: 260 }}
            placeholder="例如 guestbook"
            value={app}
            onChange={(e) => setApp(e.target.value)}
            onPressEnter={loadTimeline}
            allowClear
          />
        </Space>
        <Space style={{ marginBottom: 24 }}>
          <Typography.Link onClick={loadTimeline}>加载时间线</Typography.Link>
        </Space>
        {!loading && events.length === 0 ? (
          <EmptyState description="暂无事件。输入应用名称并加载；需配置 Argo CD 与（可选）GitLab CI 集成。" />
        ) : (
          <Timeline
            pending={loading ? '加载中…' : false}
            items={events.map((ev) => ({
              color: colorFor(ev.kind, ev.status),
              children: (
                <div>
                  <div>
                    <Typography.Text strong>{new Date(ev.ts).toLocaleString()}</Typography.Text>{' '}
                    <Typography.Text type="secondary">
                      [{ev.kind}] {ev.source}
                    </Typography.Text>
                  </div>
                  <div>
                    ref: {ev.ref || '—'} · status: {ev.status || '—'}
                  </div>
                  {ev.link ? (
                    <Typography.Link href={ev.link} target="_blank" rel="noreferrer">
                      打开链接
                    </Typography.Link>
                  ) : null}
                </div>
              ),
            }))}
          />
        )}
      </Card>
    </PageContainer>
  )
}

export default PipelineView
