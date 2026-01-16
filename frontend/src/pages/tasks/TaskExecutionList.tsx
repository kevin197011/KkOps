// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Table, Button, Space, message, Tag, Spin } from 'antd'
import { ArrowLeftOutlined, PlayCircleOutlined, EyeOutlined, StopOutlined } from '@ant-design/icons'
import { taskApi, taskExecutionApi, TaskExecution, Task } from '@/api/task'
import { assetApi, Asset } from '@/api/asset'

const TaskExecutionList = () => {
  const { taskId } = useParams<{ taskId: string }>()
  const navigate = useNavigate()
  const [executions, setExecutions] = useState<TaskExecution[]>([])
  const [task, setTask] = useState<Task | null>(null)
  const [assets, setAssets] = useState<Map<number, Asset>>(new Map())
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!taskId) {
      navigate('/tasks')
      return
    }
    fetchData()
  }, [taskId])

  const fetchData = async () => {
    if (!taskId) return
    
    setLoading(true)
    try {
      const [taskResponse, executionsResponse, assetsResponse] = await Promise.all([
        taskApi.get(parseInt(taskId)),
        taskApi.getExecutions(parseInt(taskId)),
        assetApi.list({ page: 1, page_size: 1000 }),
      ])
      
      setTask(taskResponse.data)
      setExecutions(executionsResponse.data.data || [])
      
      // Build asset lookup map
      const assetMap = new Map<number, Asset>()
      assetsResponse.data.data.forEach((asset: Asset) => {
        assetMap.set(asset.id, asset)
      })
      setAssets(assetMap)
    } catch (error: any) {
      message.error('获取执行记录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleViewLogs = (executionId: number) => {
    navigate(`/task-executions/${executionId}/logs`)
  }

  const handleCancelExecution = async (executionId: number) => {
    try {
      await taskExecutionApi.cancel(executionId)
      message.success('执行已取消')
      fetchData()
    } catch (error: any) {
      message.error(error.response?.data?.error || '取消失败')
    }
  }

  const handleExecuteTask = async (executionType: 'sync' | 'async') => {
    if (!taskId) return
    try {
      await taskApi.execute(parseInt(taskId), executionType)
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
      render: (_: any, record: TaskExecution) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewLogs(record.id)}
          >
            日志
          </Button>
          {record.status === 'running' && (
            <Button
              type="link"
              danger
              size="small"
              icon={<StopOutlined />}
              onClick={() => handleCancelExecution(record.id)}
            >
              取消
            </Button>
          )}
        </Space>
      ),
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
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/tasks')}>
            返回任务列表
          </Button>
          <h2 style={{ margin: 0 }}>
            {task?.name || '任务'} - 执行记录
          </h2>
        </Space>
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
      </div>

      <Table
        columns={columns}
        dataSource={executions}
        rowKey="id"
        scroll={{ x: 'max-content' }}
        pagination={{
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          responsive: true,
        }}
        locale={{
          emptyText: '暂无执行记录',
        }}
      />
    </div>
  )
}

export default TaskExecutionList
