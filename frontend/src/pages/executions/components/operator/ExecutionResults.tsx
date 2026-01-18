// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useRef, useEffect, useCallback } from 'react'
import {
  Card,
  Tag,
  Typography,
  Space,
  Collapse,
  Spin,
  Empty,
  Statistic,
  Row,
  Col,
  Button,
} from 'antd'
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
  ClockCircleOutlined,
  StopOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
} from '@ant-design/icons'
import { ExecutionRecord, executionRecordApi } from '@/api/execution'
import { Asset } from '@/api/asset'
import { useThemeStore } from '@/stores/theme'
import dayjs from 'dayjs'

const { Panel } = Collapse
const { Text, Title } = Typography

// 状态配置
const STATUS_CONFIG: Record<
  string,
  { color: string; icon: React.ReactNode; label: string }
> = {
  success: {
    color: 'success',
    icon: <CheckCircleOutlined />,
    label: '成功',
  },
  failed: {
    color: 'error',
    icon: <CloseCircleOutlined />,
    label: '失败',
  },
  running: {
    color: 'processing',
    icon: <SyncOutlined spin />,
    label: '运行中',
  },
  pending: {
    color: 'default',
    icon: <ClockCircleOutlined />,
    label: '待执行',
  },
  cancelled: {
    color: 'warning',
    icon: <StopOutlined />,
    label: '已取消',
  },
}

interface ExecutionResultsProps {
  executionRecords: ExecutionRecord[]
  assets: Map<number, Asset>
  isRunning: boolean
  onRefresh?: () => void
}

const ExecutionResults: React.FC<ExecutionResultsProps> = ({
  executionRecords,
  assets,
  isRunning,
  onRefresh,
}) => {
  const { mode } = useThemeStore()
  const [expandedLogs, setExpandedLogs] = useState<Set<number>>(new Set())
  const [logs, setLogs] = useState<Map<number, string[]>>(new Map())
  const [loadingLogs, setLoadingLogs] = useState<Set<number>>(new Set())
  const wsRefs = useRef<Map<number, WebSocket>>(new Map())

  // 计算执行统计
  const stats = {
    total: executionRecords.length,
    success: executionRecords.filter((r) => r.status === 'success').length,
    failed: executionRecords.filter((r) => r.status === 'failed').length,
    running: executionRecords.filter((r) => r.status === 'running').length,
    pending: executionRecords.filter((r) => r.status === 'pending').length,
    cancelled: executionRecords.filter((r) => r.status === 'cancelled').length,
  }

  // 获取总体状态
  const getOverallStatus = () => {
    if (stats.running > 0 || stats.pending > 0) return 'running'
    if (stats.failed > 0) return stats.success > 0 ? 'partial' : 'failed'
    return 'success'
  }

  const overallStatus = getOverallStatus()

  // 切换日志展开/折叠
  const toggleLog = (recordId: number) => {
    setExpandedLogs((prev) => {
      const newSet = new Set(prev)
      if (newSet.has(recordId)) {
        newSet.delete(recordId)
      } else {
        newSet.add(recordId)
        // 加载日志
        loadLogs(recordId)
      }
      return newSet
    })
  }

  // 加载日志
  const loadLogs = async (recordId: number) => {
    if (logs.has(recordId)) return

    setLoadingLogs((prev) => new Set(prev).add(recordId))
    try {
      const response = await executionRecordApi.getLogs(recordId)
      const logData = response.data.data?.logs || []
      setLogs((prev) => {
        const newMap = new Map(prev)
        newMap.set(recordId, logData.map((l: any) => l.output || l))
        return newMap
      })
    } catch (error) {
      console.error('加载日志失败:', error)
      setLogs((prev) => {
        const newMap = new Map(prev)
        newMap.set(recordId, ['无法加载日志'])
        return newMap
      })
    } finally {
      setLoadingLogs((prev) => {
        const newSet = new Set(prev)
        newSet.delete(recordId)
        return newSet
      })
    }
  }

  // WebSocket 连接（用于运行中的执行）
  useEffect(() => {
    executionRecords.forEach((record) => {
      if (record.status === 'running' && !wsRefs.current.has(record.id)) {
        const token = localStorage.getItem('token')
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/ws/execution-records/${record.id}/logs?token=${encodeURIComponent(token || '')}`

        const ws = new WebSocket(wsUrl)
        wsRefs.current.set(record.id, ws)

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            if (data.type === 'log' && data.data) {
              setLogs((prev) => {
                const newMap = new Map(prev)
                const currentLogs = newMap.get(record.id) || []
                newMap.set(record.id, [...currentLogs, data.data.output || data.data])
                return newMap
              })
            }
          } catch {
            setLogs((prev) => {
              const newMap = new Map(prev)
              const currentLogs = newMap.get(record.id) || []
              newMap.set(record.id, [...currentLogs, event.data])
              return newMap
            })
          }
        }

        ws.onclose = () => {
          wsRefs.current.delete(record.id)
        }
      }
    })

    return () => {
      wsRefs.current.forEach((ws) => ws.close())
      wsRefs.current.clear()
    }
  }, [executionRecords])

  if (executionRecords.length === 0) {
    return null
  }

  return (
    <Card
      style={{
        marginTop: 24,
        borderRadius: 12,
      }}
      styles={{
        body: { padding: '20px 24px' },
      }}
      title={
        <Space>
          <Title level={5} style={{ margin: 0 }}>
            执行结果
          </Title>
          {isRunning && (
            <Tag icon={<SyncOutlined spin />} color="processing">
              执行中
            </Tag>
          )}
          {onRefresh && (
            <Button size="small" icon={<SyncOutlined />} onClick={onRefresh}>
              刷新
            </Button>
          )}
        </Space>
      }
      extra={
        <Space>
          <Statistic
            title="总计"
            value={stats.total}
            valueStyle={{ fontSize: 18, fontWeight: 600 }}
          />
          <Statistic
            title="成功"
            value={stats.success}
            valueStyle={{ fontSize: 18, color: '#10B981', fontWeight: 600 }}
            prefix={<CheckCircleOutlined />}
          />
          <Statistic
            title="失败"
            value={stats.failed}
            valueStyle={{ fontSize: 18, color: '#EF4444', fontWeight: 600 }}
            prefix={<CloseCircleOutlined />}
          />
        </Space>
      }
    >
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        {executionRecords.map((record) => {
          const asset = assets.get(record.asset_id)
          const statusConfig = STATUS_CONFIG[record.status] || STATUS_CONFIG.pending
          const isExpanded = expandedLogs.has(record.id)
          const recordLogs = logs.get(record.id) || []

          // 计算执行时间
          let executionTime = '-'
          if (record.started_at && record.finished_at) {
            const diff = dayjs(record.finished_at).diff(dayjs(record.started_at), 'second')
            const mins = Math.floor(diff / 60)
            const secs = diff % 60
            executionTime = `${mins}分${secs}秒`
          } else if (record.started_at) {
            const diff = dayjs().diff(dayjs(record.started_at), 'second')
            const mins = Math.floor(diff / 60)
            const secs = diff % 60
            executionTime = `${mins}分${secs}秒 (进行中)`
          }

          return (
            <Card
              key={record.id}
              size="small"
              style={{
                borderRadius: 8,
                border: `1px solid ${mode === 'dark' ? '#334155' : '#E2E8F0'}`,
              }}
            >
              <Space direction="vertical" size={8} style={{ width: '100%' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Space>
                    <Tag
                      icon={statusConfig.icon}
                      color={statusConfig.color}
                      style={{ margin: 0 }}
                    >
                      {statusConfig.label}
                    </Tag>
                    <Text strong>{asset?.hostName || `Host #${record.asset_id}`}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {asset?.ip || ''}
                    </Text>
                  </Space>
                  <Space>
                    {executionTime !== '-' && (
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        {executionTime}
                      </Text>
                    )}
                    {record.exit_code !== undefined && record.exit_code !== null && (
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        退出码: {record.exit_code}
                      </Text>
                    )}
                    <Button
                      type="link"
                      size="small"
                      icon={isExpanded ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                      onClick={() => toggleLog(record.id)}
                    >
                      {isExpanded ? '隐藏日志' : '查看日志'}
                    </Button>
                  </Space>
                </div>

                {record.error && (
                  <div
                    style={{
                      padding: 8,
                      background: mode === 'dark' ? '#7F1D1D' : '#FEE2E2',
                      borderRadius: 4,
                      border: `1px solid ${mode === 'dark' ? '#991B1B' : '#FECACA'}`,
                    }}
                  >
                    <Text type="danger" style={{ fontSize: 12 }}>
                      {record.error}
                    </Text>
                  </div>
                )}

                {isExpanded && (
                  <div
                    style={{
                      background: mode === 'dark' ? '#0d1117' : '#1e1e1e',
                      borderRadius: 6,
                      padding: 12,
                      maxHeight: 400,
                      overflow: 'auto',
                      fontFamily: '"JetBrains Mono", "SF Mono", monospace',
                      fontSize: 12,
                      lineHeight: 1.6,
                      color: '#e6edf3',
                    }}
                  >
                    {loadingLogs.has(record.id) ? (
                      <div style={{ textAlign: 'center', padding: 20 }}>
                        <Spin size="small" />
                      </div>
                    ) : recordLogs.length === 0 ? (
                      <div style={{ color: '#8b949e' }}>暂无日志输出</div>
                    ) : (
                      recordLogs.map((log, index) => (
                        <div
                          key={index}
                          style={{
                            whiteSpace: 'pre-wrap',
                            wordBreak: 'break-all',
                            marginBottom: index < recordLogs.length - 1 ? 4 : 0,
                          }}
                          dangerouslySetInnerHTML={{ __html: log }}
                        />
                      ))
                    )}
                  </div>
                )}
              </Space>
            </Card>
          )
        })}
      </Space>
    </Card>
  )
}

export default ExecutionResults
