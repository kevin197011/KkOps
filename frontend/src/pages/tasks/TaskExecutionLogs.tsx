// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useRef, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Button, Spin, Tag } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import MainLayout from '@/layouts/MainLayout'

interface LogMessage {
  level: string
  content: string
  timestamp: string
}

const TaskExecutionLogs = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [logs, setLogs] = useState<string>('')
  const [connected, setConnected] = useState(false)
  const [loading, setLoading] = useState(true)
  const wsRef = useRef<WebSocket | null>(null)
  const logsEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!id) {
      navigate('/tasks')
      return
    }

    connectWebSocket()

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [id])

  useEffect(() => {
    if (logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [logs])

  const connectWebSocket = () => {
    setLoading(true)
    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/login')
      return
    }

    // Determine WebSocket protocol (ws or wss)
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/task-executions/${id}/logs?token=${encodeURIComponent(token)}`
    
    const ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      setConnected(true)
      setLoading(false)
    }

    ws.onmessage = (event) => {
      try {
        const message: LogMessage = JSON.parse(event.data)
        setLogs((prev) => prev + message.content)
      } catch (error) {
        // If message is not JSON, treat it as plain text
        setLogs((prev) => prev + event.data)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      setLoading(false)
    }

    ws.onclose = () => {
      setConnected(false)
      setLoading(false)
    }

    wsRef.current = ws
  }

  return (
    <MainLayout>
      <div>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/tasks')}>
            返回任务列表
          </Button>
          <Tag color={connected ? 'success' : 'default'}>
            {connected ? '已连接' : '未连接'}
          </Tag>
        </div>
        <Card title="执行日志" loading={loading}>
          <div
            style={{
              background: '#1e1e1e',
              color: '#d4d4d4',
              padding: '16px',
              borderRadius: '4px',
              fontFamily: 'monospace',
              fontSize: '14px',
              lineHeight: '1.5',
              minHeight: '400px',
              maxHeight: '600px',
              overflow: 'auto',
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
            }}
          >
            {logs || (loading ? <Spin /> : '等待日志输出...')}
            <div ref={logsEndRef} />
          </div>
        </Card>
      </div>
    </MainLayout>
  )
}

export default TaskExecutionLogs
