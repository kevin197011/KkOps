// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useRef, useState } from 'react'
import { Button, Input, Layout, List, Select, Space, Spin, Typography, message } from 'antd'
import { RobotOutlined, SendOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader } from '@/components/shell'
import { aiApi, streamAIChat, type AIMessage, type AIChatSessionRow } from '@/api/ai'

const { Sider, Content } = Layout

function renderAssistantText(text: string) {
  return (
    <Typography.Paragraph style={{ whiteSpace: 'pre-wrap', marginBottom: 0 }}>
      {text}
    </Typography.Paragraph>
  )
}

const AiChat = () => {
  const [sessions, setSessions] = useState<AIChatSessionRow[]>([])
  const [sessionId, setSessionId] = useState<number | undefined>()
  const [messages, setMessages] = useState<AIMessage[]>([])
  const [input, setInput] = useState('')
  const [loadingSessions, setLoadingSessions] = useState(false)
  const [sending, setSending] = useState(false)
  const [integrationId, setIntegrationId] = useState<number | undefined>()
  const [providers, setProviders] = useState<{ id: number; name: string }[]>([])
  const abortRef = useRef<AbortController | null>(null)

  const loadSessions = useCallback(async () => {
    setLoadingSessions(true)
    try {
      const res = await aiApi.listSessions()
      setSessions(res.data.data || [])
    } catch {
      message.error('Failed to load sessions')
    } finally {
      setLoadingSessions(false)
    }
  }, [])

  const loadProviders = useCallback(async () => {
    try {
      const res = await aiApi.listProviders()
      const rows = res.data.data || []
      setProviders(rows.map((r) => ({ id: r.id, name: `${r.name} (#${r.id})` })))
    } catch {
      /* ignore */
    }
  }, [])

  useEffect(() => {
    loadSessions()
    loadProviders()
  }, [loadSessions, loadProviders])

  const loadSession = async (id: number) => {
    try {
      const res = await aiApi.getSession(id)
      const d = res.data.data
      setSessionId(d.session.id)
      setMessages(
        (d.messages || []).map((m) => ({ role: m.role, content: m.content }))
      )
    } catch {
      message.error('Failed to load session')
    }
  }

  const send = async () => {
    const text = input.trim()
    if (!text || sending) return
    const userMsg: AIMessage = { role: 'user', content: text }
    setInput('')
    const historyPlusUser = [...messages, userMsg]
    setMessages([...historyPlusUser, { role: 'assistant', content: '' }])
    setSending(true)
    abortRef.current?.abort()
    abortRef.current = new AbortController()
    let acc = ''
    try {
      const { sessionId: returnedSid } = await streamAIChat(
        {
          session_id: sessionId,
          integration_id: integrationId,
          messages: historyPlusUser,
        },
        (tok) => {
          acc += tok
          setMessages((prev) => {
            const next = [...prev]
            const last = next[next.length - 1]
            if (last && last.role === 'assistant') {
              next[next.length - 1] = { role: 'assistant', content: acc }
            }
            return next
          })
        },
        abortRef.current.signal
      )
      if (returnedSid) setSessionId(returnedSid)
      await loadSessions()
    } catch (e: unknown) {
      message.error(e instanceof Error ? e.message : 'Chat failed')
      setMessages((prev) => prev.slice(0, -2))
    } finally {
      setSending(false)
    }
  }

  const newChat = () => {
    setSessionId(undefined)
    setMessages([])
  }

  return (
    <PageContainer>
      <PageHeader
        title={
          <Space>
            <RobotOutlined />
            <span>AI 运维助手</span>
          </Space>
        }
        description="Streaming assistant with read-only platform tools (alerts, incidents, metrics, logs, K8s, pipelines)."
        extra={
          <Space wrap>
            <Select
              allowClear
              placeholder="AI integration (optional)"
              style={{ minWidth: 220 }}
              options={providers.map((p) => ({ label: p.name, value: p.id }))}
              value={integrationId}
              onChange={(v) => setIntegrationId(v ?? undefined)}
            />
            <Button onClick={newChat}>New chat</Button>
          </Space>
        }
      />
      <Layout style={{ minHeight: 480, background: 'transparent' }}>
        <Sider width={260} theme="light" style={{ paddingRight: 16 }}>
          <Spin spinning={loadingSessions}>
            <List
              size="small"
              header={<Typography.Text strong>Sessions</Typography.Text>}
              dataSource={sessions}
              renderItem={(item) => (
                <List.Item
                  style={{
                    cursor: 'pointer',
                    background: sessionId === item.id ? 'rgba(0,0,0,0.06)' : undefined,
                  }}
                  onClick={() => loadSession(item.id)}
                >
                  <Typography.Text ellipsis title={item.title}>
                    {item.title || `Session #${item.id}`}
                  </Typography.Text>
                </List.Item>
              )}
            />
          </Spin>
        </Sider>
        <Content>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 12, minHeight: 400 }}>
            <div style={{ flex: 1, overflowY: 'auto', padding: 12, borderRadius: 8, background: 'var(--ant-color-fill-quaternary, #fafafa)' }}>
              {messages.map((m, i) => (
                <div key={i} style={{ marginBottom: 12 }}>
                  <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                    {m.role}
                  </Typography.Text>
                  {m.role === 'assistant' ? renderAssistantText(m.content) : (
                    <Typography.Paragraph style={{ whiteSpace: 'pre-wrap', marginBottom: 0 }}>
                      {m.content}
                    </Typography.Paragraph>
                  )}
                </div>
              ))}
              {sending && messages.length === 0 && <Spin />}
            </div>
            <Space.Compact style={{ width: '100%' }}>
              <Input.TextArea
                rows={2}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onPressEnter={(e) => {
                  if (!e.shiftKey) {
                    e.preventDefault()
                    void send()
                  }
                }}
                placeholder="Ask about alerts, incidents, metrics…"
              />
              <Button type="primary" icon={<SendOutlined />} loading={sending} onClick={() => void send()}>
                Send
              </Button>
            </Space.Compact>
          </div>
        </Content>
      </Layout>
    </PageContainer>
  )
}

export default AiChat
