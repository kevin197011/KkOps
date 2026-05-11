// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useMemo, useState } from 'react'
import { Alert, Card, Col, Row, Table, Tag, Tree, Typography, theme } from 'antd'
import type { DataNode } from 'antd/es/tree'
import type { ColumnsType } from 'antd/es/table'
import { cmdbApi, TopologyGraph } from '@/api/cmdb'
import { PageContainer, PageHeader } from '@/components/shell'

const { Text } = Typography

export default function TopologyView() {
  const { token } = theme.useToken()
  const [graph, setGraph] = useState<TopologyGraph | null>(null)
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await cmdbApi.topologyGraph()
      setGraph(res.data?.data ?? null)
    } catch {
      setGraph(null)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const treeData = useMemo<DataNode[]>(() => {
    if (!graph?.nodes?.length) return []
    const cmdbNodes = graph.nodes.filter((n) => n.kind === 'cmdb_asset')
    const bySubtype = new Map<string, typeof cmdbNodes>()
    for (const n of cmdbNodes) {
      const sub = n.subtype || 'other'
      if (!bySubtype.has(sub)) bySubtype.set(sub, [])
      bySubtype.get(sub)!.push(n)
    }
    const keys: DataNode[] = []
    for (const [kind, items] of bySubtype.entries()) {
      keys.push({
        title: `${kind} (${items.length})`,
        key: `kind-${kind}`,
        children: items.map((n) => ({
          title: (
            <span>
              <Text strong>{n.name}</Text>
              {n.env ? (
                <Tag style={{ marginLeft: 8 }} color="default">
                  {n.env}
                </Tag>
              ) : null}
            </span>
          ),
          key: `asset-${n.id}`,
          isLeaf: true,
        })),
      })
    }
    const derived = graph.nodes.filter((n) => n.kind === 'derived')
    if (derived.length) {
      keys.push({
        title: `集成占位 (${derived.length})`,
        key: 'derived-root',
        children: derived.map((n) => ({
          title: n.name,
          key: `derived-${n.id}`,
          isLeaf: true,
        })),
      })
    }
    return keys
  }, [graph])

  const edgeColumns: ColumnsType<NonNullable<TopologyGraph['edges']>[0]> = [
    {
      title: 'From',
      dataIndex: 'from',
      key: 'from',
      width: 90,
    },
    {
      title: 'To',
      dataIndex: 'to',
      key: 'to',
      width: 90,
    },
    {
      title: '关系',
      dataIndex: 'relation_type',
      key: 'relation_type',
      render: (t: string) => <Tag color="purple">{t}</Tag>,
    },
  ]

  const dc = graph?.derived_counts ?? {}

  return (
    <PageContainer>
      <PageHeader
        title="拓扑"
        description="基于 CMDB 依赖边（asset_relations）；可视化采用 Tree + 边表，避免引入重型图形库。"
      />
      <Alert
        type="info"
        showIcon
        style={{ marginBottom: 16 }}
        message="轻量拓扑视图"
        description="节点来自 CMDB 配置项；边仅包含显式保存的关系。集成节点为占位计数，不自动连线。"
      />
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={10}>
          <Card title="按类型分组" loading={loading} styles={{ body: { maxHeight: 480, overflow: 'auto' } }}>
            <Tree showLine defaultExpandAll treeData={treeData} />
          </Card>
        </Col>
        <Col xs={24} lg={14}>
          <Card title="依赖边" loading={loading}>
            <Table
              size="small"
              rowKey={(r, i) => `${r.from}-${r.to}-${i}`}
              columns={edgeColumns}
              dataSource={graph?.edges ?? []}
              pagination={{ pageSize: 12 }}
            />
          </Card>
        </Col>
      </Row>
      <Card title="派生统计" style={{ marginTop: 16 }} loading={loading}>
        <Row gutter={16}>
          <Col span={6}>
            <Text type="secondary">Kubernetes 连接器（启用）</Text>
            <div style={{ fontSize: 22, fontWeight: 700, color: token.colorPrimary }}>
              {dc.kubernetes_integrations ?? 0}
            </div>
          </Col>
          <Col span={6}>
            <Text type="secondary">Harbor 连接器（启用）</Text>
            <div style={{ fontSize: 22, fontWeight: 700, color: token.colorPrimary }}>
              {dc.harbor_integrations ?? 0}
            </div>
          </Col>
          <Col span={6}>
            <Text type="secondary">集成总数</Text>
            <div style={{ fontSize: 22, fontWeight: 700 }}>{dc.integrations_total ?? 0}</div>
          </Col>
        </Row>
      </Card>
    </PageContainer>
  )
}
