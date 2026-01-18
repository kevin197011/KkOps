// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Tag } from 'antd'
import {
  CheckCircleFilled,
  CloseCircleFilled,
  SyncOutlined,
  ClockCircleOutlined,
  StopOutlined,
} from '@ant-design/icons'
import { useThemeStore } from '@/stores/theme'
import type { WorkflowWithExecutions } from '../TaskManagementPage'

interface WorkflowCardProps {
  workflow: WorkflowWithExecutions
  isSelected: boolean
  onClick: () => void
}

const getStatusConfig = (status: string) => {
  switch (status) {
    case 'success':
      return {
        color: '#52c41a',
        icon: <CheckCircleFilled style={{ color: '#52c41a' }} />,
        text: '成功',
        bgColor: 'rgba(82, 196, 26, 0.1)',
      }
    case 'running':
      return {
        color: '#1890ff',
        icon: <SyncOutlined spin style={{ color: '#1890ff' }} />,
        text: '执行中',
        bgColor: 'rgba(24, 144, 255, 0.1)',
      }
    case 'failed':
      return {
        color: '#ff4d4f',
        icon: <CloseCircleFilled style={{ color: '#ff4d4f' }} />,
        text: '失败',
        bgColor: 'rgba(255, 77, 79, 0.1)',
      }
    case 'cancelled':
      return {
        color: '#faad14',
        icon: <StopOutlined style={{ color: '#faad14' }} />,
        text: '已取消',
        bgColor: 'rgba(250, 173, 20, 0.1)',
      }
    default:
      return {
        color: '#8c8c8c',
        icon: <ClockCircleOutlined style={{ color: '#8c8c8c' }} />,
        text: '待执行',
        bgColor: 'rgba(140, 140, 140, 0.1)',
      }
  }
}

const getRelativeTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  return `${diffDays}天前`
}

const WorkflowCard = ({ workflow, isSelected, onClick }: WorkflowCardProps) => {
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
  
  // Use last execution status or task status
  const displayStatus = workflow.lastExecution?.status || workflow.status
  const statusConfig = getStatusConfig(displayStatus)
  
  const lastRunTime = workflow.lastExecution?.created_at
  const runCount = workflow.executions.length

  return (
    <div
      onClick={onClick}
      style={{
        padding: '12px 16px',
        margin: '8px 12px',
        borderRadius: 8,
        cursor: 'pointer',
        background: isSelected
          ? isDark
            ? 'rgba(24, 144, 255, 0.15)'
            : 'rgba(24, 144, 255, 0.08)'
          : isDark
            ? '#1f1f1f'
            : '#fff',
        border: isSelected
          ? '1px solid #1890ff'
          : isDark
            ? '1px solid #303030'
            : '1px solid #f0f0f0',
        transition: 'all 0.2s ease',
      }}
      onMouseEnter={(e) => {
        if (!isSelected) {
          e.currentTarget.style.background = isDark ? '#252525' : '#fafafa'
        }
      }}
      onMouseLeave={(e) => {
        if (!isSelected) {
          e.currentTarget.style.background = isDark ? '#1f1f1f' : '#fff'
        }
      }}
    >
      {/* Header with status indicator and name */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
        {statusConfig.icon}
        <span
          style={{
            fontWeight: 500,
            fontSize: 14,
            color: isDark ? '#e6e6e6' : '#262626',
            flex: 1,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {workflow.name}
        </span>
      </div>

      {/* Type tag */}
      <div style={{ marginBottom: 8 }}>
        <Tag
          color={workflow.type === 'shell' ? 'blue' : workflow.type === 'python' ? 'green' : 'default'}
          style={{ fontSize: 11 }}
        >
          {workflow.type || 'shell'}
        </Tag>
      </div>

      {/* Footer with last run info */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          fontSize: 12,
          color: isDark ? '#8c8c8c' : '#8c8c8c',
        }}
      >
        <span>
          {lastRunTime ? `${getRelativeTime(lastRunTime)}` : '从未执行'}
        </span>
        <span>
          {runCount > 0 ? `#${runCount} ${statusConfig.text}` : ''}
        </span>
      </div>
    </div>
  )
}

export default WorkflowCard
