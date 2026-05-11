// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import {
  Button,
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
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { cmdbApi, CMDBAsset } from '@/api/cmdb'
import { PageContainer, PageHeader } from '@/components/shell'

const KIND_OPTIONS = [
  { value: 'service', label: 'service' },
  { value: 'host', label: 'host' },
  { value: 'database', label: 'database' },
  { value: 'cluster', label: 'cluster' },
  { value: 'other', label: 'other' },
]

export default function CmdbAssets() {
  const [rows, setRows] = useState<CMDBAsset[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [kindFilter, setKindFilter] = useState<string | undefined>()
  const [envFilter, setEnvFilter] = useState('')
  const [q, setQ] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<CMDBAsset | null>(null)
  const [form] = Form.useForm()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await cmdbApi.listAssets({
        kind: kindFilter,
        env: envFilter || undefined,
        q: q || undefined,
        limit: 100,
      })
      setRows(res.data.data ?? [])
      setTotal(res.data.total ?? 0)
    } catch {
      message.error('加载 CMDB 失败')
    } finally {
      setLoading(false)
    }
  }, [kindFilter, envFilter, q])

  useEffect(() => {
    load()
  }, [load])

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ kind: 'service', labels: {} })
    setModalOpen(true)
  }

  const openEdit = (r: CMDBAsset) => {
    setEditing(r)
    let labels: Record<string, unknown> = {}
    try {
      if (typeof (r as unknown as { labels?: string }).labels === 'string') {
        labels = JSON.parse((r as unknown as { labels: string }).labels || '{}')
      } else if (r.labels && typeof r.labels === 'object') {
        labels = r.labels as Record<string, unknown>
      }
    } catch {
      labels = {}
    }
    form.setFieldsValue({
      name: r.name,
      kind: r.kind,
      env: r.env,
      notes: r.notes,
      external_ref: r.external_ref,
      labels,
    })
    setModalOpen(true)
  }

  const submit = async () => {
    try {
      const v = await form.validateFields()
      const payload = {
        name: v.name,
        kind: v.kind,
        env: v.env ?? '',
        notes: v.notes ?? '',
        external_ref: v.external_ref ?? '',
        labels: v.labels && Object.keys(v.labels).length ? v.labels : undefined,
      }
      if (editing) {
        await cmdbApi.updateAsset(editing.id, payload)
        message.success('已更新')
      } else {
        await cmdbApi.createAsset(payload as { name: string; kind: string })
        message.success('已创建')
      }
      setModalOpen(false)
      load()
    } catch (e) {
      if ((e as { errorFields?: unknown }).errorFields) return
      message.error('保存失败')
    }
  }

  const remove = async (r: CMDBAsset) => {
    Modal.confirm({
      title: '删除配置项？',
      content: r.name,
      onOk: async () => {
        await cmdbApi.deleteAsset(r.id)
        message.success('已删除')
        load()
      },
    })
  }

  const columns: ColumnsType<CMDBAsset> = [
    { title: '名称', dataIndex: 'name', key: 'name', ellipsis: true },
    {
      title: '类型',
      dataIndex: 'kind',
      key: 'kind',
      width: 110,
      render: (k: string) => <Tag color="blue">{k}</Tag>,
    },
    { title: '环境', dataIndex: 'env', key: 'env', width: 100 },
    {
      title: '操作',
      key: 'actions',
      width: 140,
      render: (_, r) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(r)}>
            编辑
          </Button>
          <Button type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => remove(r)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title="CMDB"
        description="逻辑配置项（与基础设施主机资产表分离）；用于拓扑依赖建模。"
      />
      <Space wrap style={{ marginBottom: 16 }}>
        <Select
          allowClear
          placeholder="类型"
          style={{ width: 140 }}
          options={KIND_OPTIONS}
          value={kindFilter}
          onChange={setKindFilter}
        />
        <Input
          placeholder="环境"
          style={{ width: 120 }}
          value={envFilter}
          onChange={(e) => setEnvFilter(e.target.value)}
          allowClear
        />
        <Input.Search
          placeholder="搜索名称或备注"
          style={{ width: 220 }}
          onSearch={setQ}
          allowClear
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建
        </Button>
      </Space>
      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={rows}
        pagination={{ total, pageSize: 100, showTotal: (t) => `共 ${t} 条` }}
      />
      <Modal
        title={editing ? '编辑配置项' : '新建配置项'}
        open={modalOpen}
        onOk={submit}
        onCancel={() => setModalOpen(false)}
        destroyOnClose
        width={560}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="kind" label="类型" rules={[{ required: true }]}>
            <Select options={KIND_OPTIONS} />
          </Form.Item>
          <Form.Item name="env" label="环境">
            <Input placeholder="如 prod / staging" />
          </Form.Item>
          <Form.Item name="external_ref" label="外部引用">
            <Input placeholder="可选外部系统 ID / URL" />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  )
}
