// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState } from 'react'
import { Button, Card, Space, Tag, Dropdown, Descriptions, Popconfirm } from 'antd'
import type { MenuProps } from 'antd'
import {
  EditOutlined,
  PlayCircleOutlined,
  StopOutlined,
  DownOutlined,
  CodeOutlined,
  DeleteOutlined,
} from '@ant-design/icons'
import { useThemeStore } from '@/stores/theme'
import type { WorkflowWithExecutions } from '../TaskManagementPage'
import ExecutionHistoryList from './ExecutionHistoryList'

interface ExecutionDetailsPanelProps {
  workflow: WorkflowWithExecutions
  onEdit: () => void
  onExecute: (type: 'sync' | 'async') => void
  onCancel: () => void
  onDelete: () => void
  onRefresh: () => void
}

const ExecutionDetailsPanel = ({
  workflow,
  onEdit,
  onExecute,
  onCancel,
  onDelete,
  onRefresh,
}: ExecutionDetailsPanelProps) => {
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
  const [showScript, setShowScript] = useState(false)

  const isRunning = workflow.status === 'running' || workflow.lastExecution?.status === 'running'

  const executeMenuItems: MenuProps['items'] = [
    {
      key: 'sync',
      label: '同步执行',
      onClick: () => onExecute('sync'),
    },
    {
      key: 'async',
      label: '异步执行',
      onClick: () => onExecute('async'),
    },
  ]

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
      {/* Workflow Header */}
      <Card
        size="small"
        style={{
          background: isDark ? '#1f1f1f' : '#fff',
          borderColor: isDark ? '#303030' : '#f0f0f0',
        }}
      >
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            flexWrap: 'wrap',
            gap: 16,
          }}
        >
          <div style={{ flex: 1, minWidth: 200 }}>
            <h3 style={{ margin: 0, marginBottom: 8, fontSize: 18 }}>{workflow.name}</h3>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="类型">
                <Tag color={workflow.type === 'shell' ? 'blue' : 'green'}>
                  {workflow.type || 'shell'}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="超时时间">
                {workflow.timeout ? `${workflow.timeout}秒` : '默认'}
              </Descriptions.Item>
              <Descriptions.Item label="执行次数">
                {workflow.executions.length}
              </Descriptions.Item>
            </Descriptions>
          </div>

          {/* Action Buttons */}
          <Space>
            <Button icon={<EditOutlined />} onClick={onEdit}>
              编辑
            </Button>
            {isRunning ? (
              <Button icon={<StopOutlined />} danger onClick={onCancel}>
                取消
              </Button>
            ) : (
              <Dropdown menu={{ items: executeMenuItems }} trigger={['click']}>
                <Button type="primary" icon={<PlayCircleOutlined />}>
                  执行 <DownOutlined />
                </Button>
              </Dropdown>
            )}
            <Popconfirm
              title="删除任务"
              description="确定要删除这个任务吗？此操作不可撤销。"
              onConfirm={onDelete}
              okText="删除"
              cancelText="取消"
              okButtonProps={{ danger: true }}
            >
              <Button icon={<DeleteOutlined />} danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        </div>

        {/* Script Preview */}
        <div style={{ marginTop: 16 }}>
          <Button
            type="link"
            icon={<CodeOutlined />}
            onClick={() => setShowScript(!showScript)}
            style={{ padding: 0 }}
          >
            {showScript ? '隐藏脚本' : '查看脚本'}
          </Button>
          {showScript && (
            <pre
              style={{
                marginTop: 8,
                padding: 12,
                background: isDark ? '#141414' : '#f5f5f5',
                borderRadius: 6,
                overflow: 'auto',
                maxHeight: 200,
                fontSize: 12,
                lineHeight: 1.6,
              }}
            >
              {workflow.content || '(无脚本内容)'}
            </pre>
          )}
        </div>
      </Card>

      {/* Execution History */}
      <Card
        title="执行记录"
        size="small"
        style={{
          background: isDark ? '#1f1f1f' : '#fff',
          borderColor: isDark ? '#303030' : '#f0f0f0',
          flex: 1,
        }}
        styles={{
          body: { padding: 0 },
        }}
      >
        <ExecutionHistoryList
          executions={workflow.executions}
          onRefresh={onRefresh}
        />
      </Card>
    </div>
  )
}

export default ExecutionDetailsPanel
