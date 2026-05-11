// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback, useRef } from 'react'
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
  Slider,
} from 'antd'
import { SearchOutlined, ReloadOutlined, EyeOutlined, PlayCircleOutlined, PauseCircleOutlined } from '@ant-design/icons'
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

export interface TranscriptChunk {
  t: number
  d: string
}

// Parse transcript JSON into chunks with timestamps; returns null if not timed format
const parseTranscriptChunks = (transcript: string | undefined): TranscriptChunk[] | null => {
  if (!transcript?.trim()) return null
  try {
    const parsed = JSON.parse(transcript)
    if (!Array.isArray(parsed) || parsed.length === 0) return null
    const chunks: TranscriptChunk[] = parsed.map((c: { t?: number; d?: string }) => ({
      t: typeof c.t === 'number' ? c.t : 0,
      d: typeof c.d === 'string' ? c.d : '',
    }))
    return chunks
  } catch {
    return null
  }
}

// Fallback: flatten chunks to plain text (for display when not replaying)
const transcriptChunksToText = (chunks: TranscriptChunk[] | null): string => {
  if (!chunks?.length) return ''
  return chunks.map((c) => c.d).join('')
}

// Content up to positionMs (offset from first chunk time) for seeking
const contentForPosition = (chunks: TranscriptChunk[], positionMs: number): string => {
  if (!chunks.length) return ''
  const firstT = chunks[0].t
  return chunks
    .filter((c) => c.t - firstT <= positionMs)
    .map((c) => c.d)
    .join('')
}

const formatReplayTime = (ms: number): string => {
  const s = Math.floor(ms / 1000)
  const m = Math.floor(s / 60)
  return `${m}:${String(s % 60).padStart(2, '0')}`
}

// Format duration_seconds (from API) as "X分Y秒" or "Ys"
const formatDurationSeconds = (val: number | string | null | undefined): string => {
  if (val == null || val === '') return '-'
  const s = typeof val === 'number' ? Math.floor(val) : Math.floor(Number(val))
  if (Number.isNaN(s) || s < 0) return '-'
  if (s < 60) return `${s}秒`
  const m = Math.floor(s / 60)
  const sec = s % 60
  return sec > 0 ? `${m}分${sec}秒` : `${m}分`
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
  const [replayChunks, setReplayChunks] = useState<TranscriptChunk[] | null>(null)
  const [replayContent, setReplayContent] = useState('')
  const [replayPositionMs, setReplayPositionMs] = useState(0)
  const [replayPlaying, setReplayPlaying] = useState(false)
  const [replaySpeed, setReplaySpeed] = useState(1)
  const replayTimeoutsRef = useRef<ReturnType<typeof setTimeout>[]>([])
  const replayPreRef = useRef<HTMLPreElement>(null)

  const replayTotalDurationMs =
    replayChunks && replayChunks.length > 0
      ? replayChunks[replayChunks.length - 1].t - replayChunks[0].t
      : 0
  const replayDisplayContent =
    replayChunks && replayChunks.length > 0
      ? contentForPosition(replayChunks, replayPositionMs)
      : replayContent
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

  const clearReplayTimeouts = useCallback(() => {
    replayTimeoutsRef.current.forEach((id) => clearTimeout(id))
    replayTimeoutsRef.current = []
  }, [])

  const startReplay = useCallback(
    (chunks: TranscriptChunk[], speed: number, fromPositionMs: number = 0) => {
      clearReplayTimeouts()
      if (chunks.length === 0) return
      const firstT = chunks[0].t
      const lastT = chunks[chunks.length - 1].t
      const totalDurationMs = lastT - firstT
      const timeouts: ReturnType<typeof setTimeout>[] = []
      chunks.forEach((chunk) => {
        const positionMs = chunk.t - firstT
        if (positionMs < fromPositionMs) return
        const delayMs = (positionMs - fromPositionMs) / speed
        const t = setTimeout(() => {
          setReplayPositionMs(positionMs)
        }, Math.max(0, delayMs))
        timeouts.push(t)
      })
      replayTimeoutsRef.current = timeouts
      setReplayPlaying(true)
      const playDurationMs = (totalDurationMs - fromPositionMs) / speed
      const endT = setTimeout(() => {
        setReplayPositionMs(totalDurationMs)
        setReplayPlaying(false)
        replayTimeoutsRef.current = []
      }, playDurationMs + 50)
      replayTimeoutsRef.current.push(endT)
    },
    [clearReplayTimeouts]
  )

  const seekReplay = useCallback(
    (positionMs: number) => {
      clearReplayTimeouts()
      setReplayPlaying(false)
      setReplayPositionMs(Math.max(0, Math.min(positionMs, replayTotalDurationMs)))
    },
    [clearReplayTimeouts, replayTotalDurationMs]
  )

  const pauseReplay = useCallback(() => {
    clearReplayTimeouts()
    setReplayPlaying(false)
  }, [clearReplayTimeouts])

  useEffect(() => {
    if (!detailVisible) {
      clearReplayTimeouts()
      setReplayPlaying(false)
    }
    return () => clearReplayTimeouts()
  }, [detailVisible, clearReplayTimeouts])

  useEffect(() => {
    if (replayPreRef.current && replayDisplayContent) {
      replayPreRef.current.scrollTop = replayPreRef.current.scrollHeight
    }
  }, [replayDisplayContent])

  const applyTranscriptToState = useCallback(
    (transcript: string | undefined, _rec: SSHConnectionRecord, autoPlay: boolean) => {
      const chunks = parseTranscriptChunks(transcript)
      if (chunks && chunks.length > 0) {
        setReplayChunks(chunks)
        setReplayContent('')
        setReplayPositionMs(0)
        if (autoPlay) {
          requestAnimationFrame(() => startReplay(chunks, replaySpeed))
        } else {
          setReplayPositionMs(chunks[chunks.length - 1].t - chunks[0].t)
        }
      } else {
        setReplayChunks(null)
        setReplayContent(
          transcript && typeof transcript === 'string' ? transcript : transcriptChunksToText(chunks)
        )
      }
    },
    [replaySpeed, startReplay]
  )

  const handleViewReplay = async (record: SSHConnectionRecord) => {
    setSelectedRecord(record)
    setDetailVisible(true)
    setReplayContent('')
    setReplayChunks(null)
    setReplayPlaying(false)
    let transcript: string | undefined = record.transcript
    if (transcript === undefined) {
      setDetailLoading(true)
      try {
        const res = await connectionAuditApi.get(record.id)
        const rec = (res.data as SSHConnectionRecord) ?? record
        setSelectedRecord(rec)
        transcript = rec.transcript
        applyTranscriptToState(transcript, rec, true)
      } catch {
        message.error('获取录像内容失败')
      } finally {
        setDetailLoading(false)
      }
    } else {
      applyTranscriptToState(transcript, record, true)
    }
  }

  const isActiveRecord = selectedRecord?.duration_seconds === 0

  const refreshCurrentReplay = useCallback(async () => {
    if (!selectedRecord) return
    setDetailLoading(true)
    try {
      const res = await connectionAuditApi.get(selectedRecord.id)
      const rec = res.data as SSHConnectionRecord
      setSelectedRecord(rec)
      const isActive = rec.duration_seconds === 0
      applyTranscriptToState(rec.transcript, rec, !isActive)
      if (isActive) message.success('已刷新，显示最新录像')
    } catch {
      message.error('获取录像失败')
    } finally {
      setDetailLoading(false)
    }
  }, [selectedRecord, applyTranscriptToState])

  const columns: ColumnsType<SSHConnectionRecord> = [
    {
      title: '操作用户',
      dataIndex: 'login_username',
      key: 'login_username',
      width: 120,
      render: (v: string) => v || '-',
    },
    {
      title: '连线用户',
      dataIndex: 'username',
      key: 'username',
      width: 100,
    },
    {
      title: '目标资产',
      dataIndex: 'asset_hostname',
      key: 'asset_hostname',
      width: 160,
    },
    {
      title: '状态',
      key: 'status',
      width: 88,
      render: (_: unknown, record: SSHConnectionRecord) =>
        record.duration_seconds === 0 ? (
          <Tag color="green">进行中</Tag>
        ) : (
          <Tag color="default">已结束</Tag>
        ),
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
      title: '时长',
      dataIndex: 'duration_seconds',
      key: 'duration_seconds',
      width: 100,
      render: (val: number | string) => formatDurationSeconds(val),
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
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 8 }}>
            <Title level={5} style={{ margin: 0 }}>
              审计连线
            </Title>
            <Button icon={<ReloadOutlined />} onClick={() => fetchRecords()}>
              刷新
            </Button>
          </div>
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
          locale={{
            emptyText: (
              <div style={{ padding: '24px 0', color: token.colorTextSecondary }}>
                <div style={{ marginBottom: 8 }}>暂无连线记录</div>
                <div style={{ fontSize: 12 }}>
                  请先打开
                  <a href="/ssh/terminal" target="_blank" rel="noopener noreferrer" style={{ margin: '0 4px' }}>
                    WebSSH 终端
                  </a>
                  连接主机，连接建立后此处会显示记录；断开连接后会更新时长与录像。
                </div>
              </div>
            ),
          }}
        />
      </Card>

      <Modal
        title="查看/回放"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width="min(96vw, 1200px)"
        styles={{ body: { maxHeight: '85vh', overflow: 'auto' } }}
        destroyOnClose
      >
        {selectedRecord && (
          <>
            <Descriptions column={1} size="small" style={{ marginBottom: 16 }}>
              <Descriptions.Item label="操作用户">
                {selectedRecord.login_username ?? selectedRecord.username ?? '-'}
              </Descriptions.Item>
              <Descriptions.Item label="连线用户">{selectedRecord.username ?? '-'}</Descriptions.Item>
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
                {formatDurationSeconds(selectedRecord.duration_seconds)}
              </Descriptions.Item>
              {selectedRecord.transcript_truncated && (
                <Descriptions.Item label="说明">
                  <Text type="warning">录像已截断（超过 1MB）</Text>
                </Descriptions.Item>
              )}
            </Descriptions>
            <div style={{ marginBottom: 8, display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 8 }}>
              <Space>
                <Text type="secondary">终端输出</Text>
                {isActiveRecord && (
                  <Button
                    size="small"
                    icon={<ReloadOutlined />}
                    onClick={refreshCurrentReplay}
                    loading={detailLoading}
                  >
                    刷新
                  </Button>
                )}
              </Space>
              {replayChunks && replayChunks.length > 0 && !detailLoading && (
                <Space>
                  <Select
                    value={replaySpeed}
                    onChange={(v) => setReplaySpeed(v)}
                    options={[
                      { value: 0.5, label: '0.5x' },
                      { value: 1, label: '1x' },
                      { value: 2, label: '2x' },
                      { value: 4, label: '4x' },
                    ]}
                    style={{ width: 72 }}
                    disabled={replayPlaying}
                  />
                  {replayPlaying ? (
                    <Button icon={<PauseCircleOutlined />} onClick={pauseReplay}>
                      暂停
                    </Button>
                  ) : (
                    <Button
                      type="primary"
                      icon={<PlayCircleOutlined />}
                      onClick={() => startReplay(replayChunks, replaySpeed, replayPositionMs)}
                    >
                      播放
                    </Button>
                  )}
                </Space>
              )}
            </div>
            {replayChunks && replayChunks.length > 0 && replayTotalDurationMs > 0 && !detailLoading && (
              <div style={{ marginBottom: 12 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4, fontSize: 12, color: token.colorTextSecondary }}>
                  <span>{formatReplayTime(replayPositionMs)}</span>
                  <span>{formatReplayTime(replayTotalDurationMs)}</span>
                </div>
                <Slider
                  min={0}
                  max={replayTotalDurationMs}
                  value={replayPositionMs}
                  onChange={(v) => seekReplay(typeof v === 'number' ? v : v[0])}
                  tooltip={{ formatter: (v) => (v != null ? formatReplayTime(v) : '') }}
                />
              </div>
            )}
            {detailLoading ? (
              <div style={{ padding: 24, textAlign: 'center' }}>加载中...</div>
            ) : (
              <pre
                ref={replayPreRef}
                style={{
                  marginTop: 0,
                  padding: 12,
                  background: '#1e1e1e',
                  color: '#d4d4d4',
                  borderRadius: 4,
                  height: 480,
                  minHeight: 480,
                  overflow: 'auto',
                  fontSize: 12,
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-all',
                }}
              >
                {replayDisplayContent || '(无录像内容)'}
              </pre>
            )}
          </>
        )}
      </Modal>
    </div>
  )
}

export default ConnectionAuditList
