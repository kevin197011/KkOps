import React from 'react';
import { Tag, Typography, Tooltip } from 'antd';
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  SyncOutlined,
  StopOutlined,
} from '@ant-design/icons';

const { Text } = Typography;

export type WorkflowStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';

export interface WorkflowRunItemProps {
  id: number;
  name: string;
  type: 'batch' | 'formula';
  status: WorkflowStatus;
  command?: string;
  successCount: number;
  failedCount: number;
  targetCount: number;
  duration?: number;
  startedAt?: string;
  onClick?: () => void;
}

// 状态配置
const statusConfig: Record<WorkflowStatus, { 
  color: string; 
  text: string; 
  icon: React.ReactNode;
  bgColor: string;
  borderColor: string;
}> = {
  pending: { 
    color: '#faad14', 
    text: '等待中', 
    icon: <ClockCircleOutlined />,
    bgColor: '#fffbe6',
    borderColor: '#ffe58f',
  },
  running: { 
    color: '#1890ff', 
    text: '运行中', 
    icon: <SyncOutlined spin />,
    bgColor: '#e6f7ff',
    borderColor: '#91d5ff',
  },
  completed: { 
    color: '#52c41a', 
    text: '成功', 
    icon: <CheckCircleOutlined />,
    bgColor: '#f6ffed',
    borderColor: '#b7eb8f',
  },
  failed: { 
    color: '#ff4d4f', 
    text: '失败', 
    icon: <CloseCircleOutlined />,
    bgColor: '#fff2f0',
    borderColor: '#ffccc7',
  },
  cancelled: { 
    color: '#8c8c8c', 
    text: '已取消', 
    icon: <StopOutlined />,
    bgColor: '#fafafa',
    borderColor: '#d9d9d9',
  },
};

// 格式化相对时间
const formatRelativeTime = (dateStr?: string): string => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffSec < 60) return '刚刚';
  if (diffMin < 60) return `${diffMin}分钟前`;
  if (diffHour < 24) return `${diffHour}小时前`;
  if (diffDay < 7) return `${diffDay}天前`;
  return date.toLocaleDateString('zh-CN');
};

// 格式化耗时
const formatDuration = (seconds?: number): string => {
  if (seconds === undefined || seconds === null) return '-';
  if (seconds < 0) return '-';
  if (seconds === 0) return '<1s';
  if (seconds < 60) return `${Math.round(seconds)}s`;
  if (seconds < 3600) {
    const mins = Math.floor(seconds / 60);
    const secs = Math.round(seconds % 60);
    return secs > 0 ? `${mins}m ${secs}s` : `${mins}m`;
  }
  const hours = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
};

const WorkflowRunItem: React.FC<WorkflowRunItemProps> = ({
  id,
  name,
  type,
  status,
  command,
  successCount,
  failedCount,
  targetCount,
  duration,
  startedAt,
  onClick,
}) => {
  const config = statusConfig[status];

  return (
    <div
      onClick={onClick}
      style={{
        display: 'flex',
        alignItems: 'center',
        padding: '12px 16px',
        borderBottom: '1px solid #f0f0f0',
        cursor: 'pointer',
        transition: 'background-color 0.2s',
        backgroundColor: '#fff',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.backgroundColor = '#fafafa';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.backgroundColor = '#fff';
      }}
    >
      {/* 状态图标 */}
      <div style={{ 
        width: 32, 
        height: 32, 
        borderRadius: '50%',
        backgroundColor: config.bgColor,
        border: `2px solid ${config.borderColor}`,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        marginRight: 12,
        flexShrink: 0,
      }}>
        <span style={{ color: config.color, fontSize: 14 }}>{config.icon}</span>
      </div>

      {/* 主要信息 */}
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
          <Text strong style={{ 
            fontSize: 14, 
            overflow: 'hidden', 
            textOverflow: 'ellipsis', 
            whiteSpace: 'nowrap',
            maxWidth: 300,
          }}>
            {name}
          </Text>
          <Tag color={config.color} style={{ margin: 0, fontSize: 11 }}>
            {config.text}
          </Tag>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#8c8c8c', fontSize: 12 }}>
          <span>#{id}</span>
          <span>·</span>
          <span>{type === 'batch' ? '批量操作' : 'Formula部署'}</span>
          {command && (
            <>
              <span>·</span>
              <Tooltip title={command}>
                <Text code style={{ 
                  fontSize: 11, 
                  maxWidth: 150, 
                  overflow: 'hidden', 
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                  display: 'inline-block',
                  verticalAlign: 'middle',
                }}>
                  {command}
                </Text>
              </Tooltip>
            </>
          )}
          <span>·</span>
          <span>{formatRelativeTime(startedAt)}</span>
        </div>
      </div>

      {/* 统计信息 */}
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        gap: 16,
        marginLeft: 16,
        flexShrink: 0,
      }}>
        <Tooltip title={`成功: ${successCount} / 失败: ${failedCount} / 总计: ${targetCount}`}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <span style={{ color: '#52c41a', fontWeight: 500 }}>{successCount}</span>
            <span style={{ color: '#d9d9d9' }}>/</span>
            <span style={{ color: '#ff4d4f', fontWeight: 500 }}>{failedCount}</span>
            <span style={{ color: '#d9d9d9' }}>/</span>
            <span style={{ color: '#8c8c8c' }}>{targetCount}</span>
          </div>
        </Tooltip>
        {(duration !== undefined && duration !== null) && (
          <Text type="secondary" style={{ fontSize: 12, minWidth: 50, textAlign: 'right' }}>
            {formatDuration(duration)}
          </Text>
        )}
      </div>
    </div>
  );
};

export default WorkflowRunItem;
