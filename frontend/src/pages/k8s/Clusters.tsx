// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { Card, Space, Table, Tag, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useNavigate } from 'react-router-dom'
import { ClusterOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'

const Clusters = () => {
  const navigate = useNavigate()
  const [rows, setRows] = useState<IntegrationPublic[]>([])
  const [loading, setLoading] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await integrationsApi.list()
      const list = (res.data.data || []).filter((i) => i.kind.toLowerCase() === 'kubernetes')
      setRows(list)
    } catch {
      message.error('加载集群列表失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const columns: ColumnsType<IntegrationPublic> = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (v: boolean) => (v ? <Tag color="green">启用</Tag> : <Tag>停用</Tag>),
    },
    {
      title: '凭据',
      dataIndex: 'has_config',
      key: 'has_config',
      render: (v: boolean) => (v ? <Tag color="blue">已配置</Tag> : <Tag>未配置</Tag>),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title={
          <Space>
            <ClusterOutlined />
            Kubernetes 集群
          </Space>
        }
        description="从集成连接器中选择集群查看节点、工作负载与 Pod。"
      />
      <Card bordered={false}>
        {rows.length === 0 && !loading ? (
          <EmptyState description='暂无 Kubernetes 集群。请在「集成中心 → 连接器」中添加 kind 为 kubernetes 的集成并上传 kubeconfig。' />
        ) : (
          <Table
            rowKey="id"
            loading={loading}
            columns={columns}
            dataSource={rows}
            pagination={false}
            onRow={(record) => ({
              onClick: () => navigate(`/k8s/clusters/${record.id}`),
              style: { cursor: 'pointer' },
            })}
          />
        )}
      </Card>
    </PageContainer>
  )
}

export default Clusters
