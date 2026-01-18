// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState, useCallback } from 'react'
import {
  Card,
  Row,
  Col,
  Tag,
  Skeleton,
  Empty,
  Tooltip,
  theme,
  Badge,
  Typography,
  Space,
  Divider,
} from 'antd'
import {
  DatabaseOutlined,
  UserOutlined,
  PlayCircleOutlined,
  FolderOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
  ClockCircleOutlined,
  StopOutlined,
  ThunderboltOutlined,
  CloudServerOutlined,
  SafetyCertificateOutlined,
  EnvironmentOutlined,
  RocketOutlined,
} from '@ant-design/icons'
import { dashboardApi, DashboardStats } from '@/api/dashboard'

const { Text, Title } = Typography

// 状态配置
const STATUS_CONFIG: Record<
  string,
  { color: string; bg: string; icon: React.ReactNode; label: string }
> = {
  success: {
    color: '#10B981',
    bg: 'rgba(16, 185, 129, 0.1)',
    icon: <CheckCircleOutlined />,
    label: '成功',
  },
  failed: {
    color: '#EF4444',
    bg: 'rgba(239, 68, 68, 0.1)',
    icon: <CloseCircleOutlined />,
    label: '失败',
  },
  running: {
    color: '#3B82F6',
    bg: 'rgba(59, 130, 246, 0.1)',
    icon: <SyncOutlined spin />,
    label: '运行中',
  },
  pending: {
    color: '#F59E0B',
    bg: 'rgba(245, 158, 11, 0.1)',
    icon: <ClockCircleOutlined />,
    label: '待执行',
  },
  cancelled: {
    color: '#6B7280',
    bg: 'rgba(107, 114, 128, 0.1)',
    icon: <StopOutlined />,
    label: '已取消',
  },
  online: {
    color: '#10B981',
    bg: 'rgba(16, 185, 129, 0.1)',
    icon: <CheckCircleOutlined />,
    label: '在线',
  },
  offline: {
    color: '#EF4444',
    bg: 'rgba(239, 68, 68, 0.1)',
    icon: <CloseCircleOutlined />,
    label: '离线',
  },
}

// 统计卡片配置
const STAT_CARDS = [
  {
    key: 'total_assets',
    title: '资产总数',
    icon: <CloudServerOutlined />,
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    iconBg: 'rgba(255,255,255,0.2)',
  },
  {
    key: 'total_users',
    title: '用户总数',
    icon: <UserOutlined />,
    gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    iconBg: 'rgba(255,255,255,0.2)',
  },
  {
    key: 'total_tasks',
    title: '任务总数',
    icon: <ThunderboltOutlined />,
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    iconBg: 'rgba(255,255,255,0.2)',
  },
  {
    key: 'total_projects',
    title: '项目总数',
    icon: <FolderOutlined />,
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    iconBg: 'rgba(255,255,255,0.2)',
  },
]

// 格式化相对时间
function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)

  if (diffSec < 60) return '刚刚'
  if (diffMin < 60) return `${diffMin}分钟前`
  if (diffHour < 24) return `${diffHour}小时前`
  if (diffDay < 7) return `${diffDay}天前`

  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 渐变统计卡片组件
interface GradientStatCardProps {
  title: string
  value: number
  icon: React.ReactNode
  gradient: string
  iconBg: string
  loading?: boolean
}

const GradientStatCard: React.FC<GradientStatCardProps> = ({
  title,
  value,
  icon,
  gradient,
  iconBg,
  loading = false,
}) => {
  return (
    <Card
      style={{
        background: gradient,
        borderRadius: 16,
        border: 'none',
        overflow: 'hidden',
        height: 140,
        position: 'relative',
      }}
      styles={{ body: { padding: '20px 24px', height: '100%' } }}
    >
      {/* 装饰性背景圆 */}
      <div
        style={{
          position: 'absolute',
          right: -20,
          top: -20,
          width: 120,
          height: 120,
          borderRadius: '50%',
          background: 'rgba(255,255,255,0.1)',
        }}
      />
      <div
        style={{
          position: 'absolute',
          right: 30,
          bottom: -30,
          width: 80,
          height: 80,
          borderRadius: '50%',
          background: 'rgba(255,255,255,0.08)',
        }}
      />

      {loading ? (
        <Skeleton active paragraph={{ rows: 1 }} />
      ) : (
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'space-between',
            height: '100%',
            position: 'relative',
            zIndex: 1,
          }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
            <Text style={{ color: 'rgba(255,255,255,0.85)', fontSize: 14, fontWeight: 500 }}>
              {title}
            </Text>
            <div
              style={{
                width: 40,
                height: 40,
                borderRadius: 12,
                background: iconBg,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: 20,
                color: '#fff',
              }}
            >
              {icon}
            </div>
          </div>
          <div>
            <span
              style={{
                fontSize: 36,
                fontWeight: 700,
                color: '#fff',
                fontFamily: '"JetBrains Mono", "SF Mono", monospace',
                letterSpacing: '-0.02em',
              }}
            >
              {value.toLocaleString()}
            </span>
          </div>
        </div>
      )}
    </Card>
  )
}

// 系统状态面板组件
interface SystemStatusPanelProps {
  stats: DashboardStats | null
  loading: boolean
}

const SystemStatusPanel: React.FC<SystemStatusPanelProps> = ({ stats, loading }) => {
  const { token } = theme.useToken()

  const systemHealth =
    stats && stats.task_execution_stats
      ? Math.round(
          ((stats.task_execution_stats.success + stats.task_execution_stats.cancelled) /
            Math.max(stats.task_execution_stats.total, 1)) *
            100
        )
      : 0

  const getHealthColor = (health: number) => {
    if (health >= 90) return '#10B981'
    if (health >= 70) return '#F59E0B'
    return '#EF4444'
  }

  const getHealthStatus = (health: number) => {
    if (health >= 90) return '健康'
    if (health >= 70) return '警告'
    return '异常'
  }

  if (loading) {
    return (
      <Card
        style={{
          borderRadius: 16,
          height: '100%',
          background: token.colorBgContainer,
        }}
      >
        <Skeleton active paragraph={{ rows: 4 }} />
      </Card>
    )
  }

  return (
    <Card
      style={{
        borderRadius: 16,
        height: '100%',
        background: token.colorBgContainer,
      }}
      styles={{ body: { padding: '24px' } }}
    >
      {/* 标题区域 */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
        <div
          style={{
            width: 40,
            height: 40,
            borderRadius: 10,
            background: `${getHealthColor(systemHealth)}15`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: getHealthColor(systemHealth),
            fontSize: 20,
          }}
        >
          <SafetyCertificateOutlined />
        </div>
        <div>
          <Title level={5} style={{ margin: 0 }}>
            系统状态
          </Title>
          <Text type="secondary" style={{ fontSize: 12 }}>
            实时监控概览
          </Text>
        </div>
      </div>

      {/* 健康度指示器 */}
      <div
        style={{
          background: token.colorBgLayout,
          borderRadius: 12,
          padding: 20,
          marginBottom: 20,
          textAlign: 'center',
        }}
      >
        <div style={{ position: 'relative', display: 'inline-block' }}>
          {/* 圆环进度 */}
          <svg width="120" height="120" viewBox="0 0 120 120">
            <circle
              cx="60"
              cy="60"
              r="52"
              fill="none"
              stroke={token.colorBorderSecondary}
              strokeWidth="8"
            />
            <circle
              cx="60"
              cy="60"
              r="52"
              fill="none"
              stroke={getHealthColor(systemHealth)}
              strokeWidth="8"
              strokeLinecap="round"
              strokeDasharray={`${(systemHealth / 100) * 327} 327`}
              transform="rotate(-90 60 60)"
              style={{ transition: 'stroke-dasharray 0.5s ease' }}
            />
          </svg>
          <div
            style={{
              position: 'absolute',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%, -50%)',
              textAlign: 'center',
            }}
          >
            <div
              style={{
                fontSize: 28,
                fontWeight: 700,
                fontFamily: '"JetBrains Mono", monospace',
                color: getHealthColor(systemHealth),
              }}
            >
              {systemHealth}%
            </div>
            <div style={{ fontSize: 12, color: token.colorTextSecondary }}>
              {getHealthStatus(systemHealth)}
            </div>
          </div>
        </div>
      </div>

      {/* 快速统计 */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
        <div
          style={{
            padding: 16,
            background: STATUS_CONFIG.success.bg,
            borderRadius: 10,
            textAlign: 'center',
          }}
        >
          <div
            style={{
              fontSize: 24,
              fontWeight: 700,
              color: STATUS_CONFIG.success.color,
              fontFamily: '"JetBrains Mono", monospace',
            }}
          >
            {stats?.task_execution_stats?.success ?? 0}
          </div>
          <div style={{ fontSize: 12, color: token.colorTextSecondary, marginTop: 4 }}>执行成功</div>
        </div>
        <div
          style={{
            padding: 16,
            background: STATUS_CONFIG.failed.bg,
            borderRadius: 10,
            textAlign: 'center',
          }}
        >
          <div
            style={{
              fontSize: 24,
              fontWeight: 700,
              color: STATUS_CONFIG.failed.color,
              fontFamily: '"JetBrains Mono", monospace',
            }}
          >
            {stats?.task_execution_stats?.failed ?? 0}
          </div>
          <div style={{ fontSize: 12, color: token.colorTextSecondary, marginTop: 4 }}>执行失败</div>
        </div>
      </div>
    </Card>
  )
}

// 任务执行概览组件
interface ExecutionOverviewProps {
  stats: DashboardStats['task_execution_stats'] | null
  loading: boolean
}

const ExecutionOverview: React.FC<ExecutionOverviewProps> = ({ stats, loading }) => {
  const { token } = theme.useToken()

  const items = [
    { key: 'success', ...STATUS_CONFIG.success },
    { key: 'failed', ...STATUS_CONFIG.failed },
    { key: 'running', ...STATUS_CONFIG.running },
    { key: 'pending', ...STATUS_CONFIG.pending },
    { key: 'cancelled', ...STATUS_CONFIG.cancelled },
  ]

  if (loading) {
    return (
      <Card style={{ borderRadius: 16, height: '100%' }}>
        <Skeleton active paragraph={{ rows: 5 }} />
      </Card>
    )
  }

  if (!stats || stats.total === 0) {
    return (
      <Card style={{ borderRadius: 16, height: '100%' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
          <div
            style={{
              width: 40,
              height: 40,
              borderRadius: 10,
              background: '#3B82F615',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#3B82F6',
              fontSize: 20,
            }}
          >
            <PlayCircleOutlined />
          </div>
          <div>
            <Title level={5} style={{ margin: 0 }}>
              任务执行
            </Title>
            <Text type="secondary" style={{ fontSize: 12 }}>
              执行状态分布
            </Text>
          </div>
        </div>
        <Empty description="暂无执行记录" />
      </Card>
    )
  }

  return (
    <Card style={{ borderRadius: 16, height: '100%' }} styles={{ body: { padding: 24 } }}>
      {/* 标题区域 */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          marginBottom: 20,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <div
            style={{
              width: 40,
              height: 40,
              borderRadius: 10,
              background: '#3B82F615',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#3B82F6',
              fontSize: 20,
            }}
          >
            <PlayCircleOutlined />
          </div>
          <div>
            <Title level={5} style={{ margin: 0 }}>
              任务执行
            </Title>
            <Text type="secondary" style={{ fontSize: 12 }}>
              执行状态分布
            </Text>
          </div>
        </div>
        <Tag
          style={{
            background: token.colorBgLayout,
            border: 'none',
            borderRadius: 8,
            padding: '4px 12px',
            fontFamily: '"JetBrains Mono", monospace',
          }}
        >
          共 {stats.total} 次
        </Tag>
      </div>

      {/* 进度条堆叠 */}
      <div
        style={{
          display: 'flex',
          height: 12,
          borderRadius: 6,
          overflow: 'hidden',
          marginBottom: 24,
          background: token.colorBorderSecondary,
        }}
      >
        {items.map((item) => {
          const value = stats[item.key as keyof typeof stats] as number
          const percent = stats.total > 0 ? (value / stats.total) * 100 : 0
          if (percent === 0) return null
          return (
            <Tooltip key={item.key} title={`${item.label}: ${value} (${percent.toFixed(1)}%)`}>
              <div
                style={{
                  width: `${percent}%`,
                  height: '100%',
                  background: item.color,
                  transition: 'width 0.3s ease',
                }}
              />
            </Tooltip>
          )
        })}
      </div>

      {/* 详细列表 */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        {items.map((item) => {
          const value = stats[item.key as keyof typeof stats] as number
          const percent = stats.total > 0 ? (value / stats.total) * 100 : 0
          return (
            <div
              key={item.key}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '12px 16px',
                background: item.bg,
                borderRadius: 10,
                transition: 'transform 0.2s ease',
                cursor: 'default',
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                <span style={{ color: item.color, fontSize: 16 }}>{item.icon}</span>
                <span style={{ fontWeight: 500 }}>{item.label}</span>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <span
                  style={{
                    fontSize: 18,
                    fontWeight: 600,
                    fontFamily: '"JetBrains Mono", monospace',
                    color: item.color,
                  }}
                >
                  {value}
                </span>
                <span
                  style={{
                    fontSize: 12,
                    color: token.colorTextSecondary,
                    minWidth: 48,
                    textAlign: 'right',
                  }}
                >
                  {percent.toFixed(1)}%
                </span>
              </div>
            </div>
          )
        })}
      </div>
    </Card>
  )
}

// 资产分布组件
interface AssetDistributionProps {
  title: string
  subtitle: string
  icon: React.ReactNode
  iconColor: string
  data: Array<{ name: string; count: number }> | null
  loading: boolean
}

const AssetDistribution: React.FC<AssetDistributionProps> = ({
  title,
  subtitle,
  icon,
  iconColor,
  data,
  loading,
}) => {
  const { token } = theme.useToken()

  const colors = [
    '#6366F1',
    '#8B5CF6',
    '#EC4899',
    '#F43F5E',
    '#F97316',
    '#EAB308',
    '#22C55E',
    '#14B8A6',
    '#06B6D4',
    '#3B82F6',
  ]

  if (loading) {
    return (
      <Card style={{ borderRadius: 16, height: '100%' }}>
        <Skeleton active paragraph={{ rows: 5 }} />
      </Card>
    )
  }

  if (!data || data.length === 0) {
    return (
      <Card style={{ borderRadius: 16, height: '100%' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
          <div
            style={{
              width: 40,
              height: 40,
              borderRadius: 10,
              background: `${iconColor}15`,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: iconColor,
              fontSize: 20,
            }}
          >
            {icon}
          </div>
          <div>
            <Title level={5} style={{ margin: 0 }}>
              {title}
            </Title>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {subtitle}
            </Text>
          </div>
        </div>
        <Empty description="暂无数据" />
      </Card>
    )
  }

  const total = data.reduce((sum, item) => sum + item.count, 0)
  const maxCount = Math.max(...data.map((item) => item.count))

  return (
    <Card style={{ borderRadius: 16, height: '100%' }} styles={{ body: { padding: 24 } }}>
      {/* 标题区域 */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          marginBottom: 20,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <div
            style={{
              width: 40,
              height: 40,
              borderRadius: 10,
              background: `${iconColor}15`,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: iconColor,
              fontSize: 20,
            }}
          >
            {icon}
          </div>
          <div>
            <Title level={5} style={{ margin: 0 }}>
              {title}
            </Title>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {subtitle}
            </Text>
          </div>
        </div>
        <Tag
          style={{
            background: token.colorBgLayout,
            border: 'none',
            borderRadius: 8,
            padding: '4px 12px',
            fontFamily: '"JetBrains Mono", monospace',
          }}
        >
          共 {total} 台
        </Tag>
      </div>

      {/* 分布列表 */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        {data.map((item, index) => {
          const color = colors[index % colors.length]
          const percent = maxCount > 0 ? (item.count / maxCount) * 100 : 0
          return (
            <div key={item.name || 'unassigned'}>
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  marginBottom: 8,
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <div
                    style={{
                      width: 10,
                      height: 10,
                      borderRadius: 3,
                      background: color,
                    }}
                  />
                  <span style={{ fontWeight: 500 }}>{item.name || '未分配'}</span>
                </div>
                <span
                  style={{
                    fontFamily: '"JetBrains Mono", monospace',
                    fontWeight: 600,
                    color: color,
                  }}
                >
                  {item.count}
                </span>
              </div>
              <div
                style={{
                  height: 8,
                  background: token.colorBorderSecondary,
                  borderRadius: 4,
                  overflow: 'hidden',
                }}
              >
                <div
                  style={{
                    width: `${percent}%`,
                    height: '100%',
                    background: `linear-gradient(90deg, ${color}, ${color}cc)`,
                    borderRadius: 4,
                    transition: 'width 0.4s ease',
                  }}
                />
              </div>
            </div>
          )
        })}
      </div>
    </Card>
  )
}

// 最近活动组件
interface RecentActivitiesProps {
  activities: DashboardStats['recent_activities'] | null
  loading: boolean
}

const RecentActivities: React.FC<RecentActivitiesProps> = ({ activities, loading }) => {
  const { token } = theme.useToken()

  if (loading) {
    return (
      <Card style={{ borderRadius: 16 }}>
        <Skeleton active paragraph={{ rows: 8 }} />
      </Card>
    )
  }

  if (!activities || activities.length === 0) {
    return (
      <Card style={{ borderRadius: 16 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
          <div
            style={{
              width: 40,
              height: 40,
              borderRadius: 10,
              background: '#8B5CF615',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#8B5CF6',
              fontSize: 20,
            }}
          >
            <RocketOutlined />
          </div>
          <div>
            <Title level={5} style={{ margin: 0 }}>
              最近活动
            </Title>
            <Text type="secondary" style={{ fontSize: 12 }}>
              系统操作记录
            </Text>
          </div>
        </div>
        <Empty description="暂无活动记录" />
      </Card>
    )
  }

  return (
    <Card style={{ borderRadius: 16 }} styles={{ body: { padding: 24 } }}>
      {/* 标题区域 */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
        <div
          style={{
            width: 40,
            height: 40,
            borderRadius: 10,
            background: '#8B5CF615',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#8B5CF6',
            fontSize: 20,
          }}
        >
          <RocketOutlined />
        </div>
        <div>
          <Title level={5} style={{ margin: 0 }}>
            最近活动
          </Title>
          <Text type="secondary" style={{ fontSize: 12 }}>
            系统操作记录
          </Text>
        </div>
      </div>

      {/* 活动列表 */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 0 }}>
        {activities.map((activity, index) => {
          const config = STATUS_CONFIG[activity.status] || STATUS_CONFIG.pending
          return (
            <div key={activity.id}>
              <div
                style={{
                  display: 'flex',
                  gap: 16,
                  padding: '16px 0',
                }}
              >
                {/* 时间线指示器 */}
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                  }}
                >
                  <div
                    style={{
                      width: 36,
                      height: 36,
                      borderRadius: '50%',
                      background: config.bg,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: config.color,
                      fontSize: 16,
                      flexShrink: 0,
                    }}
                  >
                    {config.icon}
                  </div>
                  {index < activities.length - 1 && (
                    <div
                      style={{
                        width: 2,
                        flex: 1,
                        minHeight: 20,
                        background: token.colorBorderSecondary,
                        marginTop: 8,
                      }}
                    />
                  )}
                </div>

                {/* 内容区域 */}
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      gap: 12,
                      flexWrap: 'wrap',
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span style={{ fontWeight: 500 }}>{activity.title}</span>
                      <Tag
                        style={{
                          background: config.bg,
                          color: config.color,
                          border: 'none',
                          borderRadius: 4,
                          fontSize: 11,
                          padding: '0 6px',
                          lineHeight: '18px',
                        }}
                      >
                        {config.label}
                      </Tag>
                    </div>
                    <Text
                      type="secondary"
                      style={{
                        fontSize: 12,
                        whiteSpace: 'nowrap',
                        fontFamily: '"JetBrains Mono", monospace',
                      }}
                    >
                      {formatRelativeTime(activity.created_at)}
                    </Text>
                  </div>
                  <Text
                    type="secondary"
                    style={{
                      fontSize: 13,
                      marginTop: 4,
                      display: 'block',
                    }}
                  >
                    {activity.description}
                  </Text>
                </div>
              </div>
              {index < activities.length - 1 && <Divider style={{ margin: 0 }} />}
            </div>
          )
        })}
      </div>
    </Card>
  )
}

// 主仪表板组件
const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)
  const { token } = theme.useToken()

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true)
      const response = await dashboardApi.getStats()
      setStats(response.data.data)
    } catch (error) {
      console.error('获取仪表板数据失败:', error)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchStats()
    const interval = setInterval(fetchStats, 30000)
    return () => clearInterval(interval)
  }, [fetchStats])

  const assetsByProject = stats?.assets_by_project?.map((item) => ({
    name: item.project_name,
    count: item.count,
  }))

  const assetsByEnvironment = stats?.assets_by_environment?.map((item) => ({
    name: item.environment_name,
    count: item.count,
  }))

  return (
    <div style={{ minHeight: '100%' }}>
      {/* 页面头部 */}
      <div style={{ marginBottom: 28 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginBottom: 8 }}>
          <Title
            level={3}
            style={{
              margin: 0,
              fontFamily: '"Inter", sans-serif',
              fontWeight: 700,
              letterSpacing: '-0.02em',
            }}
          >
            仪表板
          </Title>
          <Badge
            status="processing"
            text={
              <Text type="secondary" style={{ fontSize: 12 }}>
                实时更新
              </Text>
            }
          />
        </div>
        <Text type="secondary">系统运行概览和实时状态监控</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[20, 20]} style={{ marginBottom: 24 }}>
        {STAT_CARDS.map((card) => (
          <Col xs={24} sm={12} lg={6} key={card.key}>
            <GradientStatCard
              title={card.title}
              value={(stats?.[card.key as keyof DashboardStats] as number) ?? 0}
              icon={card.icon}
              gradient={card.gradient}
              iconBg={card.iconBg}
              loading={loading}
            />
          </Col>
        ))}
      </Row>

      {/* 第二行：系统状态 + 任务执行 */}
      <Row gutter={[20, 20]} style={{ marginBottom: 24 }}>
        <Col xs={24} lg={8}>
          <SystemStatusPanel stats={stats} loading={loading} />
        </Col>
        <Col xs={24} lg={16}>
          <ExecutionOverview stats={stats?.task_execution_stats ?? null} loading={loading} />
        </Col>
      </Row>

      {/* 第三行：资产分布 */}
      <Row gutter={[20, 20]} style={{ marginBottom: 24 }}>
        <Col xs={24} lg={12}>
          <AssetDistribution
            title="按项目分布"
            subtitle="资产项目归属统计"
            icon={<FolderOutlined />}
            iconColor="#6366F1"
            data={assetsByProject ?? null}
            loading={loading}
          />
        </Col>
        <Col xs={24} lg={12}>
          <AssetDistribution
            title="按环境分布"
            subtitle="资产环境归属统计"
            icon={<EnvironmentOutlined />}
            iconColor="#22C55E"
            data={assetsByEnvironment ?? null}
            loading={loading}
          />
        </Col>
      </Row>

      {/* 第四行：最近活动 */}
      <Row gutter={[20, 20]}>
        <Col xs={24}>
          <RecentActivities activities={stats?.recent_activities ?? null} loading={loading} />
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
