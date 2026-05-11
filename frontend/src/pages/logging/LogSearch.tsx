// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { Button, Card, Col, DatePicker, Form, Input, Row, Select, message } from 'antd'
import { FileSearchOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import { loggingApi, LogLine } from '@/api/logging'
import dayjs from 'dayjs'

const LOG_KINDS = ['loki', 'elasticsearch']

const LogSearch = () => {
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [lines, setLines] = useState<LogLine[]>([])
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      setIntegrations((res.data.data || []).filter((i) => LOG_KINDS.includes(i.kind.toLowerCase())))
    } catch {
      message.error('加载集成失败')
    }
  }, [])

  useEffect(() => {
    loadIntegrations()
  }, [loadIntegrations])

  const onSearch = async () => {
    const v = await form.validateFields()
    setLoading(true)
    try {
      const res = await loggingApi.search({
        integration_id: v.integration_id,
        query: v.query,
        start: v.range?.[0]?.toISOString(),
        end: v.range?.[1]?.toISOString(),
        limit: v.limit || 200,
      })
      setLines(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '查询失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <FileSearchOutlined style={{ marginRight: 8 }} />
            日志检索
          </>
        }
        description="对接 Loki 或 Elasticsearch，输入查询语句与时间范围。"
      />
      <Card style={{ marginBottom: 16 }}>
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            range: [dayjs().subtract(1, 'hour'), dayjs()],
            limit: 200,
            query: '{job=~".+"}',
          }}
        >
          <Row gutter={16}>
            <Col xs={24} md={8}>
              <Form.Item name="integration_id" label="集成" rules={[{ required: true }]}>
                <Select
                  placeholder="选择日志集成"
                  options={integrations.map((i) => ({
                    value: i.id,
                    label: `${i.name} (${i.kind})`,
                  }))}
                />
              </Form.Item>
            </Col>
            <Col xs={24} md={16}>
              <Form.Item name="query" label="查询" rules={[{ required: true }]}>
                <Input.TextArea rows={2} placeholder="LogQL 或 Lucene 查询字符串" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col xs={24} md={12}>
              <Form.Item name="range" label="时间范围">
                <DatePicker.RangePicker showTime style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col xs={24} md={6}>
              <Form.Item name="limit" label="条数上限">
                <Input type="number" />
              </Form.Item>
            </Col>
          </Row>
          <Button type="primary" loading={loading} onClick={onSearch}>
            搜索
          </Button>
        </Form>
      </Card>

      <Card title={`结果 (${lines.length})`}>
        {lines.length === 0 ? (
          <EmptyState description="暂无日志。调整查询条件后重试。" />
        ) : (
          <div
            style={{
              fontFamily: 'monospace',
              fontSize: 12,
              maxHeight: 'calc(100vh - 420px)',
              overflow: 'auto',
              background: 'var(--ant-color-fill-quaternary)',
              padding: 12,
              borderRadius: 8,
            }}
          >
            {lines.map((ln, i) => (
              <div key={`${ln.timestamp}-${i}`} style={{ marginBottom: 6, whiteSpace: 'pre-wrap' }}>
                <span style={{ color: 'var(--ant-color-text-secondary)', marginRight: 8 }}>{ln.timestamp}</span>
                {ln.message}
              </div>
            ))}
          </div>
        )}
      </Card>
    </PageContainer>
  )
}

export default LogSearch
