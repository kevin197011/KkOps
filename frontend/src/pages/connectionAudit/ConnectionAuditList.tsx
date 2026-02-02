// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import {
  Table,
  Card,
  Form,
  Select,
  DatePicker,
  Button,
  Space,
  Modal,
  message,
  Typography,
  theme,
  Tag,
  Descriptions,
} from 'antd'
import { SearchOutlined, ReloadOutlined, EyeOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'
import {
  connectionAuditApi,
  SSHConnectionRecord,
  ConnectionAuditQueryParams,
  ConnectionAuditListResponse,
} from '@/api/connectionAudit'
import { userApi } from '@/api/user'
import { assetApi } from '@/api/asset'

const { RangePicker } = DatePicker
const { Title, Text } = Typography

interface UserOption {
  id: number
  username: string
}

interface AssetOption {
  id: number
  host_name: string
}

// Transcript 块格式：{ t: number, d: string }
const parseTranscript = (transcript: string | undefined): string => {
  if (!transcript) return ''
  try {
    const chunks: { t?: number; d: string }[] = JSON.parse(transcript)
    if (!Array.isArray(chunks)) return ''
    return chunks.map((c) => c.d ?? '').join('')
  } catch {
    return transcript
  }
}

const ConnectionAuditList = () => {
  const { token } = theme.useToken()
  const [loading, setLoading] = useState(false)
  const [records, setRecords] = useState<SSHConnectionRecord[]>([])
  const [total, setTotal] = useState(0)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20 })
  const [detailVisible, setDetailVisible] = useState(false)
  const [selectedRecord, setSelectedRecord] = useState<SSHConnectionRecord | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)
  const [replayContent, setReplayContent] = useState('')
  const [form] = Form.useForm()
  const [userOptions, setUserOptions] = useState<UserOption[]>([])
  const [assetOptions, setAssetOptions] = useState<AssetOption[]>([])

  const fetchUsers = useCallback(async () => {
    try {
      const res = await userApi.list(1, 500)
      const items = (res.data as { data?: UserOption[] })?.data ?? []
      setUserOptions(Array.isArray(items) ? items : [])
    } catch {
      setUserOptions([])
    }
  }, [])

  const fetchAssets = useCallback(async () => {
    try {
      const res = await assetApi.list({ page: 1, page_size: 500 })
      const raw = res.data as { data?: { id: number; host_name?: string; hostName?: string }[] }
      const items = raw?.data ?? []
      setAssetOptions(
        (Array.isArray(items) ? items : []).map((a) => ({
          id: a.id,
          host_name: a.host_name ?? a.hostName ?? '',
        }))
      )
    } catch {
      setAssetOptions([])
    }
  }, [])

  useEffect(() => {
    fetchUsers()
    fetchAssets()
  }, [fetchUsers, fetchAssets])

  const fetchRecords = useCallback(async () => {
    setLoading(true)
    try {
      const values = form.getFieldsValue()
      const params: Record<string, number | string> = {
        page: pagination.current,
        page_size: pagination.pageSize,
      }
      if (values.user_id != null) params.user_id = values.user_id
      if (values.asset_id != null) params.asset_id = values.asset_id
      if (values.timeRange && values.timeRange.length === 2) {
        params.start_time = values.timeRange[0].startOf('day').toISOString()
        params.end_time = values.timeRange[1].endOf('day').toISOString()
      }
      const res = await connectionAuditApi.list(params as ConnectionAuditQueryParams)
      const body = res.data as ConnectionAuditListResponse
      const data = body?.data ?? []
      setRecords(Array.isArray(data) ? data : [])
      setTotal(body?.total ?? 0)
    } catch (err: any) {
      const msg = err.response?.data?.error ?? err.message ?? '获取连线记录失败'
      message.error(typeof msg === 'string' ? msg : '获取连线记录失败')
      setRecords([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }, [form, pagination])

  useEffect(() => {
    fetchRecords()
  }, [fetchRecords])

  const handleTableChange = (newPagination: { current?: number; pageSize?: number }) => {
    setPagination({
      current: newPagination.current ?? 1,
      pageSize: newPagination.pageSize ?? 20,
    })
  }

  const handleSearch = () => {
    setPagination((p) => ({ ...p, current: 1 }))
  }

  const handleReset = () => {
    form.resetFields()
    setPagination({ current: 1, pageSize: 20 })
  }

  const handleViewReplay = async (record: SSHConnectionRecord) => {
    setSelectedRecord(record)
    setDetailVisible(true)
    setReplayContent('')
    if (record.transcript !== undefined) {
      setReplayContent(parseTranscript(record.transcript))
      return
    }
    setDetailLoading(true)
    try {
      const res = await connectionAuditApi.get(record.id)
      const rec = (res.data as SSHConnectionRecord) ?? record
      setSelectedRecord(rec)
      setReplayContent(parseTranscript(rec.transcript))
    } catch {
      message.error('获取录像内容失败')
    } finally {
      setDetailLoading(false)
    }
  }

  const columns: ColumnsType<SSHConnectionRecord> = [
    {
      title: '操作用户',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '目标资产',
      dataIndex: 'asset_hostname',
      key: 'asset_hostname',
      width: 160,
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 180,
      render: (val: string) => (val ? dayjs(val).format('YYYY-MM-DD HH:mm:ss') : '-'),
    },
    {
      title: '结束时间',
      dataIndex: 'ended_at',
      key: 'ended_at',
      width: 180,
      render: (val: string) => (val ? dayjs(val).format('YYYY-MM-DD HH:mm:ss') : '-'),
    },
    {
      title: '时长(秒)',
      dataIndex: 'duration_seconds',
      key: 'duration_seconds',
      width: 100,
      render: (val: number) => (val != null ? `${val}s` : '-'),
    },
    {
      title: '截断',
      dataIndex: 'transcript_truncated',
      key: 'transcript_truncated',
      width: 80,
      render: (val: boolean) =>
        val ? <Tag color="orange">是</Tag> : <Tag color="default">否</Tag>,
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      fixed: 'right',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => handleViewReplay(record)}
        >
          查看/回放
        </Button>
      ),
    },
  ]

  return (
    <div style={{ padding: 24, background: token.colorBgContainer, minHeight: '100%' }}>
      <Card
        style={{ background: token.colorBgElevated, borderColor: token.colorBorderSecondary }}
      >
        <Space direction="vertical" size="middle" style={{ width: '100%' }}>
          <Title level={5} style={{ margin: 0 }}>
            审计连线
          </Title>
          <Form form={form} layout="inline" onFinish={handleSearch}>
            <Form.Item name="user_id" label="用户">
              <Select
                placeholder="全部"
                allowClear
                style={{ width: 140 }}
                options={userOptions.map((u) => ({ value: u.id, label: u.username }))}
              />
            </Form.Item>
            <Form.Item name="asset_id" label="资产">
              <Select
                placeholder="全部"
                allowClear
                style={{ width: 160 }}
                options={assetOptions.map((a) => ({ value: a.id, label: a.host_name || `#${a.id}` }))}
              />
            </Form.Item>
            <Form.Item name="timeRange" label="时间范围">
              <RangePicker showTime />
            </Form.Item>
            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                  查询
                </Button>
                <Button icon={<ReloadOutlined />} onClick={handleReset}>
                  重置
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </Space>
      </Card>

      <Card
        style={{
          marginTop: 16,
          background: token.colorBgElevated,
          borderColor: token.colorBorderSecondary,
        }}
      >
        <Table<SSHConnectionRecord>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={records}
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total,
            showSizeChanger: true,
            showTotal: (t) => `共 ${t} 条`,
          }}
          onChange={handleTableChange}
        />
      </Card>

      <Modal
        title="查看/回放"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width={800}
        destroyOnClose
      >
        {selectedRecord && (
          <>
            <Descriptions column={1} size="small" style={{ marginBottom: 16 }}>
              <Descriptions.Item label="操作用户">{selectedRecord.username}</Descriptions.Item>
              <Descriptions.Item label="目标资产">{selectedRecord.asset_hostname}</Descriptions.Item>
              <Descriptions.Item label="开始时间">
                {selectedRecord.started_at
                  ? dayjs(selectedRecord.started_at).format('YYYY-MM-DD HH:mm:ss')
                  : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="结束时间">
                {selectedRecord.ended_at
                  ? dayjs(selectedRecord.ended_at).format('YYYY-MM-DD HH:mm:ss')
                  : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="时长">
                {selectedRecord.duration_seconds != null
                  ? `${selectedRecord.duration_seconds} 秒`
                  : '-'}
              </Descriptions.Item>
              {selectedRecord.transcript_truncated && (
                <Descriptions.Item label="说明">
                  <Text type="warning">录像已截断（超过 1MB）</Text>
                </Descriptions.Item>
              )}
            </Descriptions>
            <Text type="secondary">终端输出：</Text>
            {detailLoading ? (
              <div style={{ padding: 24, textAlign: 'center' }}>加载中...</div>
            ) : (
              <pre
                style={{
                  marginTop: 8,
                  padding: 12,
                  background: '#1e1e1e',
                  color: '#d4d4d4',
                  borderRadius: 4,
                  maxHeight: 400,
                  overflow: 'auto',
                  fontSize: 12,
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-all',
                }}
              >
                {replayContent || '(无录像内容)'}
              </pre>
            )}
          </>
        )}
      </Modal>
    </div>
  )
}

export default ConnectionAuditList
