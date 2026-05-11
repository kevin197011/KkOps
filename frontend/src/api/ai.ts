// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'
import type { IntegrationPublic } from './integration'

export interface AIMessage {
  role: string
  content: string
}

export interface AIChatSessionRow {
  id: number
  user_id: number
  title: string
  created_at: string
  updated_at: string
}

export interface AIChatMessageRow {
  id: number
  session_id: number
  role: string
  content: string
  created_at: string
}

export interface AIAnomalyRuleRow {
  id: number
  name: string
  integration_id: number
  query: string
  schedule_cron: string
  enabled: boolean
  prompt_template?: string
  created_at: string
  updated_at: string
}

export interface AIAnomalyFindingRow {
  id: number
  rule_id: number
  ts: string
  severity: string
  summary: string
  raw?: string
  created_at: string
}

export interface AIRcaReportRow {
  id: number
  incident_id: number
  generated_at: string
  report_md: string
  raw_tool_log?: string
}

export const aiApi = {
  listProviders: () => apiClient.get<{ data: IntegrationPublic[] }>('/ai/providers'),

  listSessions: () => apiClient.get<{ data: AIChatSessionRow[] }>('/ai/sessions'),

  getSession: (id: number) =>
    apiClient.get<{ data: { session: AIChatSessionRow; messages: AIChatMessageRow[] } }>(
      `/ai/sessions/${id}`
    ),

  deleteSession: (id: number) => apiClient.delete(`/ai/sessions/${id}`),

  listAnomalyRules: () => apiClient.get<{ data: AIAnomalyRuleRow[] }>('/ai/anomaly/rules'),

  createAnomalyRule: (body: Partial<AIAnomalyRuleRow>) =>
    apiClient.post<{ data: AIAnomalyRuleRow }>('/ai/anomaly/rules', body),

  updateAnomalyRule: (id: number, body: Partial<AIAnomalyRuleRow>) =>
    apiClient.put<{ data: AIAnomalyRuleRow }>(`/ai/anomaly/rules/${id}`, body),

  deleteAnomalyRule: (id: number) => apiClient.delete(`/ai/anomaly/rules/${id}`),

  listFindings: (params?: { rule_id?: number; since?: string; limit?: number }) =>
    apiClient.get<{ data: AIAnomalyFindingRow[] }>('/ai/anomaly/findings', { params }),

  postRca: (incidentId: number, integrationIdOverrides?: number) =>
    apiClient.post<{
      data: { report_md: string; cited_tool_calls: unknown[]; stored_id: number }
    }>('/ai/rca', {
      incident_id: incidentId,
      integration_id_overrides: integrationIdOverrides,
    }),

  listRcaReports: (params?: { incident_id?: number; limit?: number }) =>
    apiClient.get<{ data: AIRcaReportRow[] }>('/ai/rca/reports', { params }),
}

/** Streams POST /ai/chat; invokes onToken for each chunk and returns session id from header. */
export async function streamAIChat(
  body: {
    session_id?: number
    integration_id?: number
    messages: AIMessage[]
    context?: Record<string, unknown>
  },
  onToken: (t: string) => void,
  signal?: AbortSignal
): Promise<{ sessionId: number | null }> {
  const token = localStorage.getItem('token')
  const res = await fetch('/api/v1/ai/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(body),
    signal,
  })
  const sidHdr = res.headers.get('X-Session-ID')
  const sessionId = sidHdr ? parseInt(sidHdr, 10) : null
  if (!res.ok) {
    const errText = await res.text()
    let msg = errText || res.statusText
    try {
      const j = JSON.parse(errText) as { error?: string }
      if (typeof j?.error === 'string') {
        msg = j.error
      }
    } catch {
      /* leave msg as raw body */
    }
    throw new Error(msg)
  }
  const reader = res.body?.getReader()
  if (!reader) {
    return { sessionId }
  }
  const dec = new TextDecoder()
  let buf = ''
  for (;;) {
    const { done, value } = await reader.read()
    if (done) break
    buf += dec.decode(value, { stream: true })
    const parts = buf.split('\n')
    buf = parts.pop() ?? ''
    for (const line of parts) {
      if (!line.startsWith('data: ')) continue
      try {
        const data = JSON.parse(line.slice(6)) as { token?: string; done?: boolean }
        if (data.token) onToken(data.token)
      } catch {
        /* ignore malformed SSE lines */
      }
    }
  }
  return { sessionId }
}
