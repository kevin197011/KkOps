// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { Button, Drawer, InputNumber, Space, Table, Typography } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'
import { PageContainer, PageHeader } from '@/components/shell'
import { aiApi, type AIRcaReportRow } from '@/api/ai'

const RcaReports = () => {
  const [rows, setRows] = useState<AIRcaReportRow[]>([])
  const [loading, setLoading] = useState(false)
  const [incidentId, setIncidentId] = useState<number | undefined>()
  const [preview, setPreview] = useState<AIRcaReportRow | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await aiApi.listRcaReports({
        incident_id: incidentId,
        limit: 100,
      })
      setRows(res.data.data || [])
    } catch {
      setRows([])
    } finally {
      setLoading(false)
    }
  }, [incidentId])

  useEffect(() => {
    load()
  }, [load])

  const columns: ColumnsType<AIRcaReportRow> = [
    { title: 'Incident', dataIndex: 'incident_id', width: 100 },
    {
      title: 'Generated',
      dataIndex: 'generated_at',
      width: 200,
      render: (t: string) => dayjs(t).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: 'Preview',
      key: 'p',
      ellipsis: true,
      render: (_x, r) => (
        <Typography.Text ellipsis style={{ maxWidth: 400 }}>
          {r.report_md.slice(0, 120)}…
        </Typography.Text>
      ),
    },
    {
      title: '',
      key: 'open',
      width: 100,
      render: (_x, r) => (
        <Button type="link" onClick={() => setPreview(r)}>
          Open
        </Button>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title="AI 根因分析报告"
        description="Stored RCA outputs linked to incidents."
        extra={
          <Space>
            <InputNumber
              min={1}
              placeholder="filter by incident id"
              value={incidentId}
              onChange={(v) => setIncidentId(v ?? undefined)}
            />
            <Button onClick={() => load()}>Refresh</Button>
          </Space>
        }
      />
      <Table<AIRcaReportRow>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={rows}
        pagination={{ pageSize: 15 }}
      />
      <Drawer
        title={preview ? `RCA #${preview.id}` : 'RCA'}
        open={!!preview}
        onClose={() => setPreview(null)}
        width={720}
      >
        {preview && (
          <Typography.Paragraph style={{ whiteSpace: 'pre-wrap' }}>
            {preview.report_md}
          </Typography.Paragraph>
        )}
      </Drawer>
    </PageContainer>
  )
}

export default RcaReports
