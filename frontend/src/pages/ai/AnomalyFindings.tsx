// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { DatePicker, InputNumber, Select, Space, Table, Tag, Typography } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'
import { PageContainer, PageHeader } from '@/components/shell'
import { aiApi, type AIAnomalyFindingRow } from '@/api/ai'

const AnomalyFindings = () => {
  const [rows, setRows] = useState<AIAnomalyFindingRow[]>([])
  const [loading, setLoading] = useState(false)
  const [ruleId, setRuleId] = useState<number | undefined>()
  const [since, setSince] = useState<string | undefined>()
  const [limit, setLimit] = useState(100)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await aiApi.listFindings({
        rule_id: ruleId,
        since,
        limit,
      })
      setRows(res.data.data || [])
    } catch {
      setRows([])
    } finally {
      setLoading(false)
    }
  }, [ruleId, since, limit])

  useEffect(() => {
    load()
  }, [load])

  const columns: ColumnsType<AIAnomalyFindingRow> = [
    { title: 'Time', dataIndex: 'ts', width: 200, render: (t: string) => dayjs(t).format('YYYY-MM-DD HH:mm:ss') },
    { title: 'Rule', dataIndex: 'rule_id', width: 90 },
    {
      title: 'Severity',
      dataIndex: 'severity',
      width: 110,
      render: (s: string) => <Tag>{s}</Tag>,
    },
    { title: 'Summary', dataIndex: 'summary', ellipsis: true },
  ]

  return (
    <PageContainer>
      <PageHeader
        title="AI 异常检测结果"
        description="Findings produced by scheduled anomaly rules."
        extra={
          <Space wrap>
            <InputNumber
              min={1}
              placeholder="rule id"
              value={ruleId}
              onChange={(v) => setRuleId(v ?? undefined)}
            />
            <DatePicker
              showTime
              onChange={(d) => setSince(d ? d.toISOString() : undefined)}
            />
            <Select
              style={{ width: 100 }}
              value={limit}
              options={[50, 100, 200].map((n) => ({ label: `${n}`, value: n }))}
              onChange={setLimit}
            />
            <Typography.Link onClick={() => load()}>Refresh</Typography.Link>
          </Space>
        }
      />
      <Table<AIAnomalyFindingRow>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={rows}
        pagination={{ pageSize: 20 }}
      />
    </PageContainer>
  )
}

export default AnomalyFindings
