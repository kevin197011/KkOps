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
  Row,
  Select,
  Space,
  Table,
  Radio,
  DatePicker,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'
import { LineChartOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import { monitoringApi, MetricSeries, MetricPoint } from '@/api/monitoring'

const METRIC_KINDS = ['prometheus', 'nightingale']

function Sparkline({ points }: { points: MetricPoint[] }) {
  if (!points.length) return null
  const w = 240
  const h = 64
  const vals = points.map((p) => p.v)
  const min = Math.min(...vals)
  const max = Math.max(...vals)
  const span = max - min || 1
  const pad = 4
  const pts = points.map((p, i) => {
    const x = pad + (i / Math.max(points.length - 1, 1)) * (w - pad * 2)
    const y = h - pad - ((p.v - min) / span) * (h - pad * 2)
    return `${x},${y}`
  })
  return (
    <svg width={w} height={h} style={{ display: 'block' }}>
      <polyline
        fill="none"
        stroke="var(--ant-color-primary)"
        strokeWidth="2"
        points={pts.join(' ')}
      />
    </svg>
  )
}

const MonitoringQuery = () => {
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [series, setSeries] = useState<MetricSeries[]>([])
  const [loading, setLoading] = useState(false)
  const [mode, setMode] = useState<'instant' | 'range'>('instant')
  const [form] = Form.useForm()

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      const rows = (res.data.data || []).filter((i) =>
        METRIC_KINDS.includes(i.kind.toLowerCase())
      )
      setIntegrations(rows)
    } catch {
      message.error('加载集成失败')
    }
  }, [])

  useEffect(() => {
    loadIntegrations()
  }, [loadIntegrations])

  const runQuery = async () => {
    const v = await form.validateFields()
    setLoading(true)
    try {
      const body: Parameters<typeof monitoringApi.query>[0] = {
        integration_id: v.integration_id,
        query: v.query,
      }
      if (mode === 'instant' && v.at) {
        body.time = v.at.toISOString()
      }
      if (mode === 'range' && v.range) {
        body.range = {
          start: v.range[0].toISOString(),
          end: v.range[1].toISOString(),
          step: v.step || '30s',
        }
      }
      const res = await monitoringApi.query(body)
      setSeries(res.data.data.series || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '查询失败')
    } finally {
      setLoading(false)
    }
  }

  const columns: ColumnsType<MetricSeries> = useMemo(
    () => [
      {
        title: '标签',
        key: 'labels',
        render: (_, row) => (
          <code style={{ fontSize: 12 }}>
            {Object.entries(row.labels || {})
              .map(([k, v]) => `${k}="${v}"`)
              .join(', ')}
          </code>
        ),
      },
      {
        title: '样本数',
        key: 'n',
        width: 90,
        render: (_, row) => row.points?.length || 0,
      },
      {
        title: '趋势',
        key: 'spark',
        width: 260,
        render: (_, row) => <Sparkline points={row.points || []} />,
      },
    ],
    []
  )

  const firstSeries = series[0]

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <LineChartOutlined style={{ marginRight: 8 }} />
            监控查询
          </>
        }
        description="选择 Prometheus 或 Nightingale 集成，执行 PromQL 即时或区间查询。"
      />
      <Card style={{ marginBottom: 16 }}>
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            query: 'up',
            step: '30s',
            at: dayjs(),
            range: [dayjs().subtract(1, 'hour'), dayjs()],
          }}
        >
          <Row gutter={16}>
            <Col xs={24} md={8}>
              <Form.Item name="integration_id" label="集成" rules={[{ required: true }]}>
                <Select
                  placeholder="选择集成"
                  options={integrations.map((i) => ({
                    value: i.id,
                    label: `${i.name} (${i.kind})`,
                  }))}
                />
              </Form.Item>
            </Col>
            <Col xs={24} md={16}>
              <Form.Item name="query" label="PromQL" rules={[{ required: true }]}>
                <Input.TextArea rows={2} placeholder="例如 up or rate(http_requests_total[5m])" />
              </Form.Item>
            </Col>
          </Row>
          <Space wrap>
            <Radio.Group value={mode} onChange={(e) => setMode(e.target.value)}>
              <Radio.Button value="instant">即时查询</Radio.Button>
              <Radio.Button value="range">区间查询</Radio.Button>
            </Radio.Group>
            {mode === 'instant' ? (
              <Form.Item name="at" label="时间" style={{ marginBottom: 0 }}>
                <DatePicker showTime />
              </Form.Item>
            ) : (
              <>
                <Form.Item name="range" label="时间范围" style={{ marginBottom: 0 }}>
                  <DatePicker.RangePicker showTime />
                </Form.Item>
                <Form.Item name="step" label="步长" style={{ marginBottom: 0 }}>
                  <Input style={{ width: 100 }} placeholder="30s" />
                </Form.Item>
              </>
            )}
            <Button type="primary" loading={loading} onClick={runQuery}>
              查询
            </Button>
          </Space>
        </Form>
      </Card>

      {series.length === 0 ? (
        <EmptyState description="暂无查询结果。执行查询后将在此展示时间序列。" />
      ) : (
        <>
          {firstSeries && <Card title="首个序列预览" style={{ marginBottom: 16 }}>
            <Sparkline points={firstSeries.points || []} />
          </Card>}
          <Card title="序列列表">
            <Table rowKey={(_, i) => String(i)} columns={columns} dataSource={series} pagination={{ pageSize: 10 }} />
          </Card>
        </>
      )}
    </PageContainer>
  )
}

export default MonitoringQuery
