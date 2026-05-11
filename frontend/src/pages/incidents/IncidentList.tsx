// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Button,
  Drawer,
  Form,
  Input,
  Modal,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { PlusOutlined, RobotOutlined, WarningOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { incidentsApi, IncidentRow } from '@/api/incidents'
import { aiApi } from '@/api/ai'

const STATUS_OPTS = ['open', 'acknowledged', 'resolved']
const SEVERITY_OPTS = ['critical', 'high', 'medium', 'low']

const IncidentList = () => {
  const navigate = useNavigate()
  const [rows, setRows] = useState<IncidentRow[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [selected, setSelected] = useState<IncidentRow | null>(null)
  const [createForm] = Form.useForm()
  const [patchForm] = Form.useForm()
  const [rcaOpen, setRcaOpen] = useState(false)
  const [rcaMd, setRcaMd] = useState('')
  const [rcaLoading, setRcaLoading] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await incidentsApi.list({ limit: 100 })
      setRows(res.data.data || [])
      setTotal(res.data.total ?? 0)
    } catch {
      message.error('Failed to load incidents')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const openDetail = useCallback(
    async (row: IncidentRow) => {
      setSelected(row)
      setDrawerOpen(true)
      try {
        const res = await incidentsApi.get(row.id)
        const d = res.data.data
        setSelected(d)
        patchForm.setFieldsValue({
          ...d,
          linked_alert_ids: d.linked_alert_ids?.map(String),
        })
      } catch {
        message.error('Failed to load incident')
      }
    },
    [patchForm]
  )

  const columns: ColumnsType<IncidentRow> = useMemo(
    () => [
      { title: 'Title', dataIndex: 'title', ellipsis: true },
      {
        title: 'Severity',
        dataIndex: 'severity',
        width: 110,
        render: (s: string) => <Tag>{s}</Tag>,
      },
      {
        title: 'Status',
        dataIndex: 'status',
        width: 130,
      },
      {
        title: 'Alerts',
        key: 'alerts',
        width: 90,
        render: (_v, r) => (r.linked_alert_ids?.length ? r.linked_alert_ids.length : '—'),
      },
      {
        title: 'Actions',
        key: 'act',
        width: 120,
        render: (_v, r) => (
          <Button type="link" onClick={() => openDetail(r)}>
            Detail
          </Button>
        ),
      },
    ],
    [openDetail]
  )

  const submitCreate = async () => {
    const v = await createForm.validateFields()
    const rawIds = typeof v.linked_alert_ids === 'string' ? v.linked_alert_ids : ''
    const linked_alert_ids = rawIds
      .split(/[\s,]+/)
      .map((x: string) => parseInt(x.trim(), 10))
      .filter((n: number) => !Number.isNaN(n))
    const assignee =
      v.assignee_user_id !== undefined && v.assignee_user_id !== ''
        ? Number(v.assignee_user_id)
        : undefined
    try {
      await incidentsApi.create({
        title: v.title,
        severity: v.severity,
        linked_alert_ids: linked_alert_ids.length ? linked_alert_ids : undefined,
        assignee_user_id: assignee && !Number.isNaN(assignee) ? assignee : undefined,
      })
      message.success('Created')
      setModalOpen(false)
      createForm.resetFields()
      load()
    } catch (e: unknown) {
      message.error('Create failed')
    }
  }

  const submitPatch = async () => {
    if (!selected) return
    const v = await patchForm.validateFields()
    const rawTags = (v.linked_alert_ids as string[] | undefined) || []
    const linked_alert_ids = rawTags
      .map((x) => parseInt(String(x), 10))
      .filter((n) => !Number.isNaN(n))
    const assignee =
      v.assignee_user_id !== undefined && v.assignee_user_id !== ''
        ? Number(v.assignee_user_id)
        : undefined
    try {
      await incidentsApi.patch(selected.id, {
        title: v.title,
        severity: v.severity,
        status: v.status,
        linked_alert_ids,
        assignee_user_id: assignee && !Number.isNaN(assignee) ? assignee : undefined,
      })
      message.success('Updated')
      setDrawerOpen(false)
      load()
    } catch {
      message.error('Update failed')
    }
  }

  const runRca = async () => {
    if (!selected) return
    setRcaLoading(true)
    try {
      const res = await aiApi.postRca(selected.id)
      setRcaMd(res.data.data.report_md)
      setRcaOpen(true)
    } catch {
      message.error('AI RCA failed')
    } finally {
      setRcaLoading(false)
    }
  }

  return (
    <PageContainer>
      <PageHeader
        title={
          <Space>
            <WarningOutlined />
            <span>Incidents</span>
          </Space>
        }
        description="Minimal incident tracking; link alert IDs manually."
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            New incident
          </Button>
        }
      />
      {rows.length === 0 && !loading ? (
        <EmptyState description="No incidents yet. Create one to track an outage." />
      ) : (
        <Table<IncidentRow>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={rows}
          pagination={{ total, pageSize: 100, showTotal: (t) => `${t} incidents` }}
        />
      )}

      <Modal
        title="New incident"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={submitCreate}
        destroyOnClose
      >
        <Form form={createForm} layout="vertical">
          <Form.Item name="title" label="Title" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="severity" label="Severity" rules={[{ required: true }]}>
            <Select options={SEVERITY_OPTS.map((s) => ({ label: s, value: s }))} />
          </Form.Item>
          <Form.Item name="linked_alert_ids" label="Linked alert IDs (comma-separated)">
            <Input placeholder="e.g. 1, 2, 3" />
          </Form.Item>
          <Form.Item name="assignee_user_id" label="Assignee user ID">
            <Input type="number" />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer
        title={selected ? `Incident #${selected.id}` : 'Incident'}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={480}
        extra={
          <Space>
            <Button icon={<RobotOutlined />} loading={rcaLoading} onClick={() => void runRca()}>
              AI 根因分析
            </Button>
            <Button type="link" onClick={() => navigate('/ai/rca/reports')}>
              RCA 列表
            </Button>
            <Button onClick={() => setDrawerOpen(false)}>Cancel</Button>
            <Button type="primary" onClick={submitPatch}>
              Save
            </Button>
          </Space>
        }
      >
        {selected && (
          <Form form={patchForm} layout="vertical">
            <Form.Item name="title" label="Title" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="severity" label="Severity" rules={[{ required: true }]}>
              <Select options={SEVERITY_OPTS.map((s) => ({ label: s, value: s }))} />
            </Form.Item>
            <Form.Item name="status" label="Status" rules={[{ required: true }]}>
              <Select options={STATUS_OPTS.map((s) => ({ label: s, value: s }))} />
            </Form.Item>
            <Form.Item name="linked_alert_ids" label="Linked alert IDs">
              <Select
                mode="tags"
                placeholder="Numeric alert IDs"
                tokenSeparators={[',']}
              />
            </Form.Item>
            <Form.Item name="assignee_user_id" label="Assignee user ID">
              <Input type="number" />
            </Form.Item>
            <Typography.Text type="secondary">
              Created by user #{selected.created_by}
            </Typography.Text>
          </Form>
        )}
      </Drawer>

      <Modal
        title="AI 根因分析"
        open={rcaOpen}
        onCancel={() => setRcaOpen(false)}
        footer={null}
        width={840}
      >
        <Typography.Paragraph style={{ whiteSpace: 'pre-wrap' }}>{rcaMd}</Typography.Paragraph>
      </Modal>
    </PageContainer>
  )
}

export default IncidentList
