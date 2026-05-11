// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { Button, Form, Input, Modal, Select, Space, Switch, Table, Tag, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { PageContainer, PageHeader } from '@/components/shell'
import { aiApi, type AIAnomalyRuleRow } from '@/api/ai'
import { integrationsApi, type IntegrationPublic } from '@/api/integration'

const AnomalyRules = () => {
  const [rows, setRows] = useState<AIAnomalyRuleRow[]>([])
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<AIAnomalyRuleRow | null>(null)
  const [form] = Form.useForm()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const [rulesRes, intRes] = await Promise.all([
        aiApi.listAnomalyRules(),
        integrationsApi.list(),
      ])
      setRows(rulesRes.data.data || [])
      const prom = (intRes.data.data || []).filter((i) => i.kind === 'prometheus')
      setIntegrations(prom.length ? prom : intRes.data.data || [])
    } catch {
      message.error('Failed to load rules')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({
      enabled: true,
      schedule_cron: '*/15 * * * *',
    })
    setModalOpen(true)
  }

  const openEdit = (r: AIAnomalyRuleRow) => {
    setEditing(r)
    form.setFieldsValue(r)
    setModalOpen(true)
  }

  const submit = async () => {
    const v = await form.validateFields()
    try {
      if (editing) {
        await aiApi.updateAnomalyRule(editing.id, v)
        message.success('Updated')
      } else {
        await aiApi.createAnomalyRule(v)
        message.success('Created')
      }
      setModalOpen(false)
      load()
    } catch {
      message.error('Save failed')
    }
  }

  const remove = async (id: number) => {
    try {
      await aiApi.deleteAnomalyRule(id)
      message.success('Deleted')
      load()
    } catch {
      message.error('Delete failed')
    }
  }

  const columns: ColumnsType<AIAnomalyRuleRow> = [
    { title: 'Name', dataIndex: 'name', ellipsis: true },
    { title: 'Integration', dataIndex: 'integration_id', width: 100 },
    {
      title: 'Cron',
      dataIndex: 'schedule_cron',
      width: 140,
      ellipsis: true,
    },
    {
      title: 'Enabled',
      dataIndex: 'enabled',
      width: 90,
      render: (e: boolean) => (e ? <Tag color="green">yes</Tag> : <Tag>no</Tag>),
    },
    {
      title: 'Actions',
      key: 'a',
      width: 160,
      render: (_x, r) => (
        <Space>
          <Button type="link" size="small" onClick={() => openEdit(r)}>
            Edit
          </Button>
          <Button type="link" danger size="small" onClick={() => remove(r.id)}>
            Delete
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title="AI 异常检测规则"
        description="Prometheus query_range + LLM classification. Uses monitoring integrations (prometheus kind)."
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            New rule
          </Button>
        }
      />
      <Table<AIAnomalyRuleRow>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={rows}
        pagination={{ pageSize: 20 }}
      />
      <Modal
        title={editing ? 'Edit rule' : 'New rule'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={submit}
        destroyOnClose
        width={640}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="integration_id" label="Prometheus integration" rules={[{ required: true }]}>
            <Select
              options={integrations.map((i) => ({
                label: `${i.name} (${i.kind})`,
                value: i.id,
              }))}
            />
          </Form.Item>
          <Form.Item name="query" label="PromQL query" rules={[{ required: true }]}>
            <Input.TextArea rows={3} placeholder="e.g. rate(http_requests_total[5m])" />
          </Form.Item>
          <Form.Item name="schedule_cron" label="Cron schedule" rules={[{ required: true }]}>
            <Input placeholder="*/15 * * * *" />
          </Form.Item>
          <Form.Item name="enabled" label="Enabled" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="prompt_template" label="Prompt template (optional)">
            <Input.TextArea rows={4} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  )
}

export default AnomalyRules
