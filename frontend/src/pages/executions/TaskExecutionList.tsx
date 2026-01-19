// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Table, Button, Space, message, Tag, Spin } from 'antd'
import { ArrowLeftOutlined, PlayCircleOutlined, EyeOutlined, StopOutlined } from '@ant-design/icons'
import { executionApi, executionRecordApi, ExecutionRecord, Execution } from '@/api/execution'
import { assetApi, Asset } from '@/api/asset'
import { usePermissionStore } from '@/stores/permission'

const TaskExecutionList = () => {
  const { executionId } = useParams<{ executionId: string }>()
  const navigate = useNavigate()
  const { hasPermission } = usePermissionStore()
  const [executions, setExecutions] = useState<ExecutionRecord[]>([])
  const [execution, setExecution] = useState<Execution | null>(null)
  const [assets, setAssets] = useState<Map<number, Asset>>(new Map())
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)

  const fetchExecutions = useCallback(async () => {
    if (!executionId) return
    try {
      const historyResponse = await executionApi.getHistory(parseInt(executionId))
      // API 返回的是数组，需要前端分页
      const allExecutions = historyResponse.data.data || []
      const start = (page - 1) * pageSize
      const end = start + pageSize
      setExecutions(allExecutions.slice(start, end))
      setTotal(allExecutions.length)
    } catch (error: any) {
      message.error('获取执行历史失败')
    }
  }, [executionId, page, pageSize])

  const fetchData = useCallback(async () => {
    if (!executionId) {
      navigate('/executions')
      return
    }
    
    setLoading(true)
    try {
      const [executionResponse, assetsResponse] = await Promise.all([
        executionApi.get(parseInt(executionId)),
        assetApi.list({ page: 1, page_size: 1000 }),
      ])
      
      setExecution(executionResponse.data)
      
      // Build asset lookup map
      const assetMap = new Map<number, Asset>()
      assetsResponse.data.data.forEach((asset: Asset) => {
        assetMap.set(asset.id, asset)
      })
      setAssets(assetMap)
      
      // 获取执行历史（分页）
      await fetchExecutions()
    } catch (error: any) {
      message.error('获取执行记录失败')
    } finally {
      setLoading(false)
    }
  }, [executionId, fetchExecutions, navigate])

  useEffect(() => {
    if (executionId) {
      fetchData()
    }
  }, [executionId, fetchData])

  useEffect(() => {
    if (executionId) {
      fetchExecutions()
    }
  }, [executionId, fetchExecutions])

  const handleViewLogs = (recordId: number) => {
    navigate(`/execution-records/${recordId}/logs`)
  }

  const handleCancelExecution = async (recordId: number) => {
    try {
      await executionRecordApi.cancel(recordId)
      message.success('执行已取消')
      fetchData()
    } catch (error: any) {
      message.error(error.response?.data?.error || '取消失败')
    }
  }

  const handleExecuteTask = async (execType: 'sync' | 'async') => {
    if (!executionId) return
    try {
      await executionApi.execute(parseInt(executionId), execType)
      message.success('任务执行已启动')
      fetchData()
    } catch (error: any) {
      message.error(error.response?.data?.error || '执行失败')
    }
  }

  const getStatusTag = (status: string) => {
    const config: Record<string, { color: string; text: string }> = {
      pending: { color: 'default', text: '等待中' },
      running: { color: 'processing', text: '执行中' },
      success: { color: 'success', text: '成功' },
      failed: { color: 'error', text: '失败' },
      cancelled: { color: 'warning', text: '已取消' },
    }
    const { color, text } = config[status] || { color: 'default', text: status }
    return <Tag color={color}>{text}</Tag>
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '执行主机',
      dataIndex: 'asset_id',
      key: 'asset_id',
      render: (assetId: number) => {
        const asset = assets.get(assetId)
        if (asset) {
          return `${asset.hostName} (${asset.ip})`
        }
        return `资产 #${assetId}`
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: getStatusTag,
    },
    {
      title: '退出码',
      dataIndex: 'exit_code',
      key: 'exit_code',
      width: 80,
      render: (code: number | undefined) => {
        if (code === undefined || code === null) return '-'
        return code === 0 ? (
          <Tag color="success">{code}</Tag>
        ) : (
          <Tag color="error">{code}</Tag>
        )
      },
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 180,
      render: (time: string) => time || '-',
    },
    {
      title: '结束时间',
      dataIndex: 'finished_at',
      key: 'finished_at',
      width: 180,
      render: (time: string) => time || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: ExecutionRecord) => {
        const actions = []
        // 查看日志不需要特殊权限（已有该页面访问权限即可）
        actions.push(
          <Button
            key="logs"
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewLogs(record.id)}
          >
            日志
          </Button>
        )
        if (record.status === 'running' && hasPermission('executions', 'update')) {
          actions.push(
            <Button
              key="cancel"
              type="link"
              danger
              size="small"
              icon={<StopOutlined />}
              onClick={() => handleCancelExecution(record.id)}
            >
              取消
            </Button>
          )
        }
        return <Space size="small">{actions}</Space>
      },
    },
  ]

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div>
      <div
        style={{
          marginBottom: 16,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 8,
        }}
      >
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/executions')}>
            返回运维执行
          </Button>
          <h2 style={{ margin: 0 }}>
            {execution?.name || '执行任务'} - 执行历史
          </h2>
        </Space>
        {hasPermission('executions', 'create') && (
          <Space>
            <Button
              type="primary"
              icon={<PlayCircleOutlined />}
              onClick={() => handleExecuteTask('sync')}
            >
              同步执行
            </Button>
            <Button
              icon={<PlayCircleOutlined />}
              onClick={() => handleExecuteTask('async')}
            >
              异步执行
            </Button>
          </Space>
        )}
      </div>

      <Table
        columns={columns}
        dataSource={executions}
        rowKey="id"
        scroll={{ x: 'max-content' }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
          responsive: true,
          onChange: (newPage, newPageSize) => {
            setPage(newPage)
            setPageSize(newPageSize || 20)
          },
          onShowSizeChange: (current, size) => {
            setPage(1)
            setPageSize(size)
          },
        }}
        locale={{
          emptyText: '暂无执行记录',
        }}
      />
    </div>
  )
}

export default TaskExecutionList
