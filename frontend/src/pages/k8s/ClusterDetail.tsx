// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import {
  Button,
  Card,
  Drawer,
  Select,
  Space,
  Table,
  Tabs,
  Tag,
  Typography,
  message,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeftOutlined, ClusterOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import {
  k8sApi,
  K8sEvent,
  K8sNamespace,
  K8sNode,
  K8sPod,
  K8sWorkload,
} from '@/api/k8s'

const { Text } = Typography

const ClusterDetail = () => {
  const { clusterId } = useParams<{ clusterId: string }>()
  const navigate = useNavigate()
  const id = Number(clusterId)
  const [meta, setMeta] = useState<IntegrationPublic | null>(null)
  const [ns, setNs] = useState<string>('default')
  const [namespaces, setNamespaces] = useState<K8sNamespace[]>([])
  const [nodes, setNodes] = useState<K8sNode[]>([])
  const [workloads, setWorkloads] = useState<K8sWorkload[]>([])
  const [pods, setPods] = useState<K8sPod[]>([])
  const [events, setEvents] = useState<K8sEvent[]>([])
  const [loading, setLoading] = useState(false)
  const [tab, setTab] = useState('overview')
  const [logOpen, setLogOpen] = useState(false)
  const [logPod, setLogPod] = useState('')
  const [logText, setLogText] = useState('')
  const [logLoading, setLogLoading] = useState(false)

  const validId = Number.isFinite(id) && id > 0

  const loadMeta = useCallback(async () => {
    if (!validId) return
    try {
      const res = await integrationsApi.get(id)
      setMeta(res.data.data)
    } catch {
      message.error('加载集成信息失败')
    }
  }, [id, validId])

  const loadNamespaces = useCallback(async () => {
    if (!validId) return
    try {
      const res = await k8sApi.namespaces(id)
      setNamespaces(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载 Namespace 失败')
    }
  }, [id, validId])

  const loadNodes = useCallback(async () => {
    if (!validId) return
    setLoading(true)
    try {
      const res = await k8sApi.nodes(id)
      setNodes(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载节点失败')
    } finally {
      setLoading(false)
    }
  }, [id, validId])

  const loadWorkloads = useCallback(async () => {
    if (!validId) return
    setLoading(true)
    try {
      const res = await k8sApi.workloads(id, ns)
      setWorkloads(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载工作负载失败')
    } finally {
      setLoading(false)
    }
  }, [id, ns, validId])

  const loadPods = useCallback(async () => {
    if (!validId) return
    setLoading(true)
    try {
      const res = await k8sApi.pods(id, ns)
      setPods(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载 Pod 失败')
    } finally {
      setLoading(false)
    }
  }, [id, ns, validId])

  const loadEvents = useCallback(async () => {
    if (!validId) return
    setLoading(true)
    try {
      const res = await k8sApi.events(id, ns)
      setEvents(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载事件失败')
    } finally {
      setLoading(false)
    }
  }, [id, ns, validId])

  useEffect(() => {
    loadMeta()
    loadNamespaces()
  }, [loadMeta, loadNamespaces])

  useEffect(() => {
    if (tab === 'nodes') loadNodes()
    if (tab === 'workloads') loadWorkloads()
    if (tab === 'pods') loadPods()
    if (tab === 'events') loadEvents()
  }, [tab, loadNodes, loadWorkloads, loadPods, loadEvents])

  const openLogs = async (podName: string) => {
    if (!validId) return
    setLogPod(podName)
    setLogOpen(true)
    setLogLoading(true)
    setLogText('')
    try {
      const res = await k8sApi.podLogs(id, podName, { namespace: ns, tail: 200 })
      setLogText(typeof res.data === 'string' ? res.data : String(res.data))
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '获取日志失败')
    } finally {
      setLogLoading(false)
    }
  }

  const nodeCols: ColumnsType<K8sNode> = [
    { title: '节点', dataIndex: 'name', key: 'name' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: 'Kubelet', dataIndex: 'kubelet_version', key: 'kubelet_version' },
    {
      title: '地址',
      dataIndex: 'addresses',
      key: 'addresses',
      render: (a: string[]) => a.join(', '),
    },
  ]

  const wlCols: ColumnsType<K8sWorkload> = [
    { title: '类型', dataIndex: 'kind', key: 'kind' },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '副本', key: 'rep', render: (_, r) => `${r.available}/${r.replicas}` },
    { title: '镜像', dataIndex: 'image', key: 'image', ellipsis: true },
  ]

  const podCols: ColumnsType<K8sPod> = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '节点', dataIndex: 'node', key: 'node' },
    { title: '重启', dataIndex: 'restarts', key: 'restarts' },
    { title: '存活时间', dataIndex: 'age', key: 'age' },
    {
      title: '操作',
      key: 'act',
      render: (_, r) => (
        <Button type="link" size="small" onClick={() => openLogs(r.name)}>
          日志
        </Button>
      ),
    },
  ]

  const evCols: ColumnsType<K8sEvent> = [
    { title: '类型', dataIndex: 'type', key: 'type', width: 90 },
    { title: '原因', dataIndex: 'reason', key: 'reason', width: 120 },
    { title: '对象', dataIndex: 'object', key: 'object' },
    { title: '消息', dataIndex: 'message', key: 'message', ellipsis: true },
    { title: '最近', dataIndex: 'last', key: 'last', width: 100 },
  ]

  if (!validId) {
    return (
      <PageContainer>
        <EmptyState description="无效的集群 ID" />
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <Button type="link" icon={<ArrowLeftOutlined />} onClick={() => navigate('/k8s/clusters')}>
            返回列表
          </Button>
        </div>
        <PageHeader
          title={
            <Space>
              <ClusterOutlined />
              {meta?.name || `集群 #${id}`}
            </Space>
          }
          description="查看节点、工作负载、Pod 与事件；未指定命名空间时使用 default。"
        />
        <Card size="small">
          <Space wrap align="center">
            <Text type="secondary">命名空间</Text>
            <Select
              style={{ minWidth: 200 }}
              value={ns}
              onChange={setNs}
              options={namespaces.map((n) => ({ label: n.name, value: n.name }))}
              showSearch
              optionFilterProp="label"
            />
          </Space>
        </Card>
        <Tabs
          activeKey={tab}
          onChange={setTab}
          items={[
            {
              key: 'overview',
              label: '概览',
              children: (
                <Card bordered={false}>
                  <Space direction="vertical">
                    <div>
                      Namespace 数量：<Tag>{namespaces.length}</Tag>
                    </div>
                    <div>
                      集成 ID：<Tag>{id}</Tag>
                    </div>
                  </Space>
                </Card>
              ),
            },
            {
              key: 'nodes',
              label: '节点',
              children: (
                <Table rowKey="name" loading={loading} columns={nodeCols} dataSource={nodes} pagination={false} />
              ),
            },
            {
              key: 'workloads',
              label: '工作负载',
              children: (
                <Table rowKey={(r) => `${r.kind}/${r.name}`} loading={loading} columns={wlCols} dataSource={workloads} />
              ),
            },
            {
              key: 'pods',
              label: 'Pods',
              children: (
                <Table rowKey="name" loading={loading} columns={podCols} dataSource={pods} />
              ),
            },
            {
              key: 'events',
              label: '事件',
              children: (
                <Table
                  rowKey={(r, i) => `${r.object}-${r.reason}-${i}`}
                  loading={loading}
                  columns={evCols}
                  dataSource={events}
                />
              ),
            },
          ]}
        />
      </Space>
      <Drawer title={`Pod 日志 — ${logPod}`} width={720} open={logOpen} onClose={() => setLogOpen(false)}>
        <pre style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word', fontSize: 12 }}>
          {logLoading ? 'Loading…' : logText || '(empty)'}
        </pre>
      </Drawer>
    </PageContainer>
  )
}

export default ClusterDetail
