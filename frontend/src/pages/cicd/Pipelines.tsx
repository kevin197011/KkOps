// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import {
  Button,
  Card,
  Drawer,
  Form,
  Input,
  Modal,
  Select,
  Space,
  Table,
  Tag,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { RocketOutlined, SyncOutlined, FileTextOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import { cicdApi, PipelineRow } from '@/api/cicd'

const CI_KINDS = ['jenkins', 'gitlab']

const Pipelines = () => {
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [integrationId, setIntegrationId] = useState<number | undefined>()
  const [rows, setRows] = useState<PipelineRow[]>([])
  const [loading, setLoading] = useState(false)
  const [logsOpen, setLogsOpen] = useState(false)
  const [logText, setLogText] = useState('')
  const [logLoading, setLogLoading] = useState(false)
  const [activePipe, setActivePipe] = useState<PipelineRow | null>(null)
  const [runOpen, setRunOpen] = useState(false)
  const [runForm] = Form.useForm<{ integration_id?: number; ref?: string; vars?: string }>()

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      const list = (res.data.data || []).filter((i) => CI_KINDS.includes(i.kind.toLowerCase()))
      setIntegrations(list)
      setIntegrationId((prev) => prev ?? list[0]?.id)
    } catch {
      message.error('加载集成失败')
    }
  }, [])

  const loadPipelines = useCallback(async () => {
    if (!integrationId) {
      setRows([])
      return
    }
    setLoading(true)
    try {
      const res = await cicdApi.listPipelines(integrationId)
      setRows(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载流水线失败')
    } finally {
      setLoading(false)
    }
  }, [integrationId])

  useEffect(() => {
    loadIntegrations()
  }, [loadIntegrations])

  useEffect(() => {
    loadPipelines()
  }, [loadPipelines])

  const openLogs = async (p: PipelineRow) => {
    setActivePipe(p)
    setLogsOpen(true)
    setLogLoading(true)
    try {
      const res = await cicdApi.logs(p.id, p.integration_id)
      setLogText(res.data.data || '')
    } catch {
      message.error('加载日志失败')
      setLogText('')
    } finally {
      setLogLoading(false)
    }
  }

  const submitRun = async () => {
    if (!activePipe) return
    const v = await runForm.validateFields()
    try {
      let variables: Record<string, string> | undefined
      if (v.vars?.trim()) {
        variables = {}
        for (const line of v.vars.split('\n')) {
          const idx = line.indexOf('=')
          if (idx > 0) {
            variables[line.slice(0, idx).trim()] = line.slice(idx + 1).trim()
          }
        }
      }
      await cicdApi.runPipeline(activePipe.id, {
        integration_id: v.integration_id ?? activePipe.integration_id,
        ref: v.ref,
        variables,
      })
      message.success('已触发运行')
      setRunOpen(false)
      loadPipelines()
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '触发失败')
    }
  }

  const columns: ColumnsType<PipelineRow> = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (s: string) => <Tag>{s}</Tag>,
    },
    { title: '引用', dataIndex: 'ref', key: 'ref', ellipsis: true },
    {
      title: '操作',
      key: 'act',
      render: (_, p) => (
        <Space>
          <Button
            size="small"
            icon={<SyncOutlined />}
            onClick={() => {
              setActivePipe(p)
              runForm.setFieldsValue({
                integration_id: p.integration_id,
                ref: p.ref || 'main',
                vars: '',
              })
              setRunOpen(true)
            }}
          >
            运行
          </Button>
          <Button size="small" icon={<FileTextOutlined />} onClick={() => openLogs(p)}>
            日志
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <RocketOutlined style={{ marginRight: 8 }} />
            CI/CD
          </>
        }
        description="查看 Jenkins / GitLab CI 流水线并触发构建、查看日志。"
      />
      <Card style={{ marginBottom: 16 }}>
        <Space wrap>
          <span>集成：</span>
          <Select
            style={{ minWidth: 260 }}
            placeholder="选择 CI 集成"
            value={integrationId}
            onChange={(v) => setIntegrationId(v)}
            options={integrations.map((i) => ({
              value: i.id,
              label: `${i.name} (${i.kind})`,
            }))}
          />
          <Button type="primary" onClick={loadPipelines} loading={loading}>
            刷新
          </Button>
        </Space>
      </Card>

      {integrations.length === 0 ? (
        <EmptyState description="未配置 CI 集成。请先在集成中心添加 Jenkins 或 GitLab。" />
      ) : (
        <Card title="流水线">
          <Table rowKey="id" loading={loading} columns={columns} dataSource={rows} pagination={{ pageSize: 15 }} />
        </Card>
      )}

      <Drawer title="日志" width={720} open={logsOpen} onClose={() => setLogsOpen(false)} destroyOnClose>
        <pre style={{ whiteSpace: 'pre-wrap', fontSize: 12 }}>{logLoading ? '加载中…' : logText}</pre>
      </Drawer>

      <Modal title="触发流水线" open={runOpen} onCancel={() => setRunOpen(false)} onOk={submitRun} destroyOnClose>
        <Form form={runForm} layout="vertical">
          <Form.Item name="integration_id" hidden>
            <Input />
          </Form.Item>
          <Form.Item name="ref" label="分支 / Tag" rules={[{ required: true }]}>
            <Input placeholder="main" />
          </Form.Item>
          <Form.Item name="vars" label="变量（可选，每行 KEY=VALUE）">
            <Input.TextArea rows={4} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  )
}

export default Pipelines
