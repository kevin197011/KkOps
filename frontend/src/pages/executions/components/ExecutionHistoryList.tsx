// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useRef, useCallback, useMemo } from 'react'
import { Empty, Tag, Button, Spin } from 'antd'
import { DesktopOutlined } from '@ant-design/icons'
import {
  CheckCircleFilled,
  CloseCircleFilled,
  SyncOutlined,
  ClockCircleOutlined,
  StopOutlined,
  DownOutlined,
  UpOutlined,
} from '@ant-design/icons'
import { useThemeStore } from '@/stores/theme'
import { ExecutionRecord, executionRecordApi } from '@/api/execution'

interface ExecutionHistoryListProps {
  executions: ExecutionRecord[]
  onRefresh: () => void
  autoExpandLatest?: boolean  // Auto-expand the latest running execution
}

const getStatusConfig = (status: string) => {
  switch (status) {
    case 'success':
      return {
        color: 'success',
        icon: <CheckCircleFilled style={{ color: '#52c41a' }} />,
        text: '成功',
      }
    case 'running':
      return {
        color: 'processing',
        icon: <SyncOutlined spin style={{ color: '#1890ff' }} />,
        text: '执行中',
      }
    case 'failed':
      return {
        color: 'error',
        icon: <CloseCircleFilled style={{ color: '#ff4d4f' }} />,
        text: '失败',
      }
    case 'cancelled':
      return {
        color: 'warning',
        icon: <StopOutlined style={{ color: '#faad14' }} />,
        text: '已取消',
      }
    default:
      return {
        color: 'default',
        icon: <ClockCircleOutlined style={{ color: '#8c8c8c' }} />,
        text: '待执行',
      }
  }
}

const formatTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

const getDuration = (start: string, end?: string) => {
  const startDate = new Date(start)
  const endDate = end ? new Date(end) : new Date()
  const diffMs = endDate.getTime() - startDate.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  
  if (diffSec < 60) return `${diffSec}秒`
  const mins = Math.floor(diffSec / 60)
  const secs = diffSec % 60
  return `${mins}分${secs}秒`
}

interface InlineLogViewerProps {
  executionId: number
  isRunning: boolean
  assetName?: string  // Asset hostname for display
  assetIp?: string    // Asset IP for display
}

const InlineLogViewer = ({ executionId, isRunning, assetName, assetIp }: InlineLogViewerProps) => {
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
  const [logs, setLogs] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [autoScroll, setAutoScroll] = useState(true)
  const logContainerRef = useRef<HTMLDivElement>(null)
  const wsRef = useRef<WebSocket | null>(null)

  const scrollToBottom = useCallback(() => {
    if (logContainerRef.current && autoScroll) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight
    }
  }, [autoScroll])

  useEffect(() => {
    // Fetch initial logs
    const fetchLogs = async () => {
      try {
        const response = await executionRecordApi.getLogs(executionId)
        const logData = response.data.data?.logs || []
        setLogs(logData.map((l: any) => l.output || l))
      } catch {
        setLogs(['无法加载日志'])
      } finally {
        setLoading(false)
      }
    }

    fetchLogs()

    // If running, connect WebSocket for real-time logs
    if (isRunning) {
      const token = localStorage.getItem('token')
      const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsHost = window.location.hostname
      const wsPort = import.meta.env.DEV ? '8080' : window.location.port
      const wsUrl = `${wsProtocol}//${wsHost}:${wsPort}/api/v1/ws/task-logs/${executionId}?token=${token}`

      const ws = new WebSocket(wsUrl)
      wsRef.current = ws

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (data.type === 'log' && data.data) {
            setLogs((prev) => [...prev, data.data.output || data.data])
            scrollToBottom()
          }
        } catch {
          // Handle plain text logs
          setLogs((prev) => [...prev, event.data])
          scrollToBottom()
        }
      }

      ws.onclose = () => {
        // Connection closed
      }
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
    }
  }, [executionId, isRunning, scrollToBottom])

  useEffect(() => {
    scrollToBottom()
  }, [logs, scrollToBottom])

  if (loading) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <Spin size="small" />
      </div>
    )
  }

  // Host label for log header
  const hostLabel = assetName ? `${assetName}${assetIp ? ` (${assetIp})` : ''}` : `Asset #${executionId}`

  return (
    <div>
      {/* Host Header */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          marginBottom: 8,
          padding: '6px 10px',
          background: isDark ? '#1a1a2e' : '#e3f2fd',
          borderRadius: 4,
          fontSize: 12,
          fontWeight: 500,
          color: isDark ? '#79c0ff' : '#1565c0',
        }}
      >
        <span style={{ fontSize: 14 }}>🖥️</span>
        <span>{hostLabel}</span>
        {isRunning && (
          <Tag color="processing" style={{ marginLeft: 'auto', marginRight: 0 }}>
            <SyncOutlined spin /> 执行中
          </Tag>
        )}
      </div>

      <div
        ref={logContainerRef}
        style={{
          background: isDark ? '#0d1117' : '#1e1e1e',
          borderRadius: 6,
          padding: 12,
          maxHeight: 300,
          overflow: 'auto',
          fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
          fontSize: 12,
          lineHeight: 1.6,
          color: '#e6edf3',
        }}
        onScroll={(e) => {
          const target = e.target as HTMLDivElement
          const isAtBottom = target.scrollHeight - target.scrollTop - target.clientHeight < 20
          setAutoScroll(isAtBottom)
        }}
      >
        {logs.length === 0 ? (
          <div style={{ color: '#8b949e' }}>暂无日志输出</div>
        ) : (
          logs.map((log, index) => (
            <div
              key={index}
              style={{
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-all',
              }}
              dangerouslySetInnerHTML={{
                __html: log
                  .replace(/\x1b\[31m/g, '<span style="color:#ff7b72">')
                  .replace(/\x1b\[32m/g, '<span style="color:#7ee787">')
                  .replace(/\x1b\[33m/g, '<span style="color:#ffa657">')
                  .replace(/\x1b\[34m/g, '<span style="color:#79c0ff">')
                  .replace(/\x1b\[35m/g, '<span style="color:#d2a8ff">')
                  .replace(/\x1b\[36m/g, '<span style="color:#56d4dd">')
                  .replace(/\x1b\[0m/g, '</span>')
                  .replace(/\x1b\[\d+m/g, ''),
              }}
            />
          ))
        )}
        {isRunning && (
          <div style={{ color: '#1890ff', marginTop: 8 }}>
            <SyncOutlined spin /> 实时日志流...
          </div>
        )}
      </div>
      <div
        style={{
          display: 'flex',
          justifyContent: 'flex-end',
          marginTop: 8,
          fontSize: 12,
          color: '#8c8c8c',
        }}
      >
        <span>
          {autoScroll ? '📜 自动滚动已开启' : '📜 自动滚动已暂停（滚动以恢复）'}
        </span>
      </div>
    </div>
  )
}

const ExecutionHistoryList = ({ executions, onRefresh, autoExpandLatest = true }: ExecutionHistoryListProps) => {
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set())
  const [showAll, setShowAll] = useState(false)
  const prevExecutionsRef = useRef<ExecutionRecord[]>([])

  // Sort executions by created_at descending (newest first)
  const sortedExecutions = useMemo(() => {
    return [...executions].sort((a, b) => {
      return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    })
  }, [executions])

  // Auto-expand only the latest execution when new ones appear, collapse all others
  useEffect(() => {
    if (!autoExpandLatest) return

    // Check if a new execution was added
    if (sortedExecutions.length > prevExecutionsRef.current.length && sortedExecutions.length > 0) {
      const latestExecution = sortedExecutions[0]
      // Only expand the newest one, collapse all others
      setExpandedIds(new Set([latestExecution.id]))
    }

    prevExecutionsRef.current = sortedExecutions
  }, [sortedExecutions, autoExpandLatest])

  const displayedExecutions = showAll ? sortedExecutions : sortedExecutions.slice(0, 5)

  // 计算每个执行记录的任务内运行编号（从最早到最新：#1, #2, #3...）
  const executionRunNumbers = useMemo(() => {
    const runNumbers = new Map<number, number>()
    // 按创建时间升序排列，然后分配编号
    const chronological = [...executions].sort((a, b) => 
      new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    )
    chronological.forEach((exec, index) => {
      runNumbers.set(exec.id, index + 1)
    })
    return runNumbers
  }, [executions])

  const toggleExpand = (id: number) => {
    setExpandedIds(prev => {
      const newSet = new Set(prev)
      if (newSet.has(id)) {
        newSet.delete(id)
      } else {
        newSet.add(id)
      }
      return newSet
    })
  }

  if (executions.length === 0) {
    return (
      <div style={{ padding: 40 }}>
        <Empty description="暂无执行记录" />
      </div>
    )
  }

  return (
    <div>
      {displayedExecutions.map((execution) => {
        const statusConfig = getStatusConfig(execution.status)
        const isExpanded = expandedIds.has(execution.id)
        const isRunning = execution.status === 'running'
        const assetName = execution.asset?.hostname || execution.asset?.ip
        const assetIp = execution.asset?.ip

        return (
          <div
            key={execution.id}
            style={{
              borderBottom: isDark ? '1px solid #303030' : '1px solid #f0f0f0',
            }}
          >
            {/* Execution Row */}
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                padding: '12px 16px',
                cursor: 'pointer',
                gap: 12,
                background: isExpanded
                  ? isDark
                    ? 'rgba(24, 144, 255, 0.08)'
                    : 'rgba(24, 144, 255, 0.04)'
                  : 'transparent',
              }}
              onClick={() => toggleExpand(execution.id)}
            >
              {/* Expand Icon */}
              <span style={{ color: '#8c8c8c' }}>
                {isExpanded ? <UpOutlined /> : <DownOutlined />}
              </span>

              {/* Status Icon */}
              {statusConfig.icon}

              {/* Run Number (任务内的运行编号) */}
              <span style={{ fontWeight: 500, minWidth: 50 }}>#{executionRunNumbers.get(execution.id) || 1}</span>

              {/* Asset Info */}
              {assetName && (
                <Tag icon={<DesktopOutlined />} color="default" style={{ margin: 0 }}>
                  {assetName}
                </Tag>
              )}

              {/* Status Tag */}
              <Tag color={statusConfig.color}>{statusConfig.text}</Tag>

              {/* Duration */}
              <span style={{ color: '#8c8c8c', fontSize: 12, flex: 1 }}>
                {getDuration(execution.created_at, execution.finished_at)}
              </span>

              {/* Time */}
              <span style={{ color: '#8c8c8c', fontSize: 12 }}>
                {formatTime(execution.created_at)}
              </span>

              {/* Exit Code (if failed) */}
              {execution.exit_code !== undefined && execution.exit_code !== 0 && (
                <Tag color="red">退出码: {execution.exit_code}</Tag>
              )}
            </div>

            {/* Inline Logs (Expanded) */}
            {isExpanded && (
              <div style={{ padding: '0 16px 16px' }}>
                <InlineLogViewer 
                  executionId={execution.id} 
                  isRunning={isRunning}
                  assetName={assetName}
                  assetIp={assetIp}
                />
              </div>
            )}
          </div>
        )
      })}

      {/* Show More Button */}
      {sortedExecutions.length > 5 && (
        <div style={{ padding: 16, textAlign: 'center' }}>
          <Button type="link" onClick={() => setShowAll(!showAll)}>
            {showAll ? '收起' : `显示更多 (${sortedExecutions.length - 5} 条)`}
          </Button>
        </div>
      )}
    </div>
  )
}

export default ExecutionHistoryList
