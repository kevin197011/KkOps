// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { Button, Card, Col, Row, Skeleton, Statistic, theme, Typography } from 'antd'
import {
  UserOutlined,
  ThunderboltOutlined,
  CloudServerOutlined,
  FolderOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
  EnvironmentOutlined,
  SafetyCertificateOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons'
import { dashboardApi, DashboardStats, DashboardSummary, RecentActivity } from '@/api/dashboard'
import { PageContainer, PageHeader } from '@/components/shell'
import { usePermissionStore } from '@/stores/permission'

const { Text } = Typography

const STAT_ITEMS = [
  {
    key: 'total_assets',
    title: '资产',
    icon: CloudServerOutlined,
    bg: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
  },
  {
    key: 'total_users',
    title: '用户',
    icon: UserOutlined,
    bg: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
  },
  {
    key: 'total_tasks',
    title: '任务',
    icon: ThunderboltOutlined,
    bg: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
  },
  {
    key: 'total_projects',
    title: '项目',
    icon: FolderOutlined,
    bg: 'linear-gradient(135deg, #0ea5e9 0%, #38bdf8 100%)',
  },
] as const

const STATUS_LABELS: Record<string, string> = {
  active: '在线',
  disabled: '停用',
  unknown: '未知',
}

const ACTIVITY_STATUS_MAP: Record<string, { color: string; text: string }> = {
  success: { color: '#10b981', text: '成功' },
  failed: { color: '#ef4444', text: '失败' },
  running: { color: '#3b82f6', text: '运行中' },
  pending: { color: '#94a3b8', text: '待执行' },
  cancelled: { color: '#64748b', text: '已取消' },
}

function RecentActivityList({
  items,
  loading,
  token,
}: {
  items: RecentActivity[]
  loading: boolean
  token: ReturnType<typeof theme.useToken>['token']
}) {
  const list = (items ?? []).slice(0, 5)
  if (loading) return <Skeleton active paragraph={{ rows: 2 }} />
  if (list.length === 0) {
    return (
      <Text type="secondary" style={{ fontSize: 12 }}>
        暂无最近执行
      </Text>
    )
  }
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      {list.map((a) => {
        const statusInfo = ACTIVITY_STATUS_MAP[a.status] ?? {
          color: token.colorTextSecondary,
          text: a.status,
        }
        return (
          <div
            key={a.id}
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              fontSize: 12,
              padding: '6px 0',
              borderBottom: `1px solid ${token.colorBorderSecondary}`,
            }}
          >
            <Text ellipsis style={{ flex: 1, marginRight: 8, fontSize: 12 }}>
              {a.title}
            </Text>
            <span
              style={{
                color: statusInfo.color,
                fontWeight: 500,
                flexShrink: 0,
                marginRight: 8,
              }}
            >
              {statusInfo.text}
            </span>
            <Text type="secondary" style={{ fontSize: 11, flexShrink: 0 }}>
              {a.created_at}
            </Text>
          </div>
        )
      })}
    </div>
  )
}

const Dashboard = () => {
  const { token } = theme.useToken()
  const { hasPermission } = usePermissionStore()
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [summary, setSummary] = useState<DashboardSummary | null>(null)
  const [loading, setLoading] = useState(true)

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true)
      const [statsRes, sumRes] = await Promise.all([
        dashboardApi.getStats(),
        dashboardApi.getSummary(),
      ])
      setStats(statsRes.data.data)
      setSummary(sumRes.data.data)
    } catch {
      // ignore
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchStats()
    const t = setInterval(fetchStats, 60000)
    return () => clearInterval(t)
  }, [fetchStats])

  const exec = stats?.task_execution_stats
  const totalExec = exec ? exec.total : 0
  const healthPercent =
    totalExec > 0 && exec
      ? Math.round(((exec.success + exec.cancelled) / totalExec) * 100)
      : 100

  const assetsByStatus = stats?.assets_by_status ?? {}

  return (
    <PageContainer>
      <div
        style={{
          height: 'calc(100vh - 64px)',
          display: 'flex',
          flexDirection: 'column',
          padding: 24,
          boxSizing: 'border-box',
          background: `linear-gradient(180deg, ${token.colorBgLayout} 0%, ${token.colorBgContainer} 100%)`,
        }}
      >
      <PageHeader
        title="仪表盘"
        description="资源规模、资产分布与任务执行概览"
      />

      {/* Operations hub summary */}
      <Row gutter={[16, 16]} style={{ flexShrink: 0, marginBottom: 16 }}>
        <Col xs={24}>
          <Card
            size="small"
            title="运营概览"
            style={{ borderRadius: 12, boxShadow: token.boxShadowSecondary }}
          >
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : summary ? (
              <>
                <Row gutter={[16, 16]}>
                  <Col xs={12} sm={8} md={6}>
                    <Statistic title="未恢复告警" value={summary.alerts_open} />
                  </Col>
                  <Col xs={12} sm={8} md={6}>
                    <Statistic title="处理中事件" value={summary.incidents_open} />
                  </Col>
                  <Col xs={12} sm={8} md={6}>
                    <Statistic
                      title="集成（启用 / 总数）"
                      value={summary.integrations_ok}
                      suffix={`/ ${summary.integrations_total}`}
                    />
                  </Col>
                  <Col xs={12} sm={8} md={6}>
                    <Statistic title="账号同步运行（24h）" value={summary.provisioning_runs_24h} />
                  </Col>
                  <Col xs={12} sm={8} md={6}>
                    <Statistic title="AI 异常命中（24h）" value={summary.ai_anomaly_findings_24h} />
                  </Col>
                  <Col xs={12} sm={8} md={6}>
                    <Statistic title="CMDB 配置项" value={summary.cmdb_assets_total} />
                  </Col>
                </Row>
                <div style={{ marginTop: 16, display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                  {hasPermission('alerts', '*') && (
                    <Link to="/alerts">
                      <Button size="small">告警中心</Button>
                    </Link>
                  )}
                  {hasPermission('incidents', '*') && (
                    <Link to="/incidents">
                      <Button size="small">事件管理</Button>
                    </Link>
                  )}
                  {hasPermission('integrations', '*') && (
                    <Link to="/integrations">
                      <Button size="small">集成中心</Button>
                    </Link>
                  )}
                  {hasPermission('ai', '*') && (
                    <Link to="/ai/chat">
                      <Button size="small" type="primary">
                        AI 助手
                      </Button>
                    </Link>
                  )}
                  {hasPermission('cmdb', '*') && (
                    <Link to="/topology">
                      <Button size="small">拓扑</Button>
                    </Link>
                  )}
                </div>
              </>
            ) : (
              <Text type="secondary">暂无摘要</Text>
            )}
          </Card>
        </Col>
      </Row>

      {/* 1. 资源规模 */}
      <Row gutter={[16, 16]} style={{ flexShrink: 0, marginBottom: 20 }}>
        {STAT_ITEMS.map(({ key, title, icon: Icon, bg }) => (
          <Col xs={12} sm={12} md={6} key={key}>
            <Card
              bordered={false}
              style={{
                borderRadius: 12,
                overflow: 'hidden',
                boxShadow: token.boxShadowSecondary,
              }}
              styles={{ body: { padding: 20 } }}
            >
              {loading ? (
                <Skeleton active paragraph={{ rows: 0 }} />
              ) : (
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                  }}
                >
                  <div>
                    <Text
                      style={{
                        color: token.colorTextSecondary,
                        fontSize: 13,
                      }}
                    >
                      {title}
                    </Text>
                    <div
                      style={{
                        fontSize: 28,
                        fontWeight: 700,
                        fontVariantNumeric: 'tabular-nums',
                        color: token.colorText,
                        marginTop: 4,
                        letterSpacing: '-0.02em',
                      }}
                    >
                      {(stats?.[key] as number) ?? 0}
                    </div>
                  </div>
                  <div
                    style={{
                      width: 48,
                      height: 48,
                      borderRadius: 12,
                      background: bg,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontSize: 22,
                      color: '#fff',
                      boxShadow: '0 4px 12px rgba(0,0,0,0.08)',
                    }}
                  >
                    <Icon />
                  </div>
                </div>
              )}
            </Card>
          </Col>
        ))}
      </Row>

      {/* 2. 任务执行 + 最近活动 | 资产概览（状态 + 项目 + 环境） */}
      <div
        style={{
          flex: 1,
          minHeight: 0,
          display: 'flex',
          gap: 16,
          flexDirection: 'column',
        }}
      >
        <Row gutter={[16, 16]} style={{ flex: 1, minHeight: 0 }}>
          <Col xs={24} lg={14}>
            <Card
              title="任务执行"
              size="small"
              style={{
                borderRadius: 12,
                height: '100%',
                boxShadow: token.boxShadowSecondary,
              }}
              styles={{
                body: {
                  padding: 20,
                  height: 'calc(100% - 57px)',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: 16,
                  overflow: 'auto',
                },
              }}
            >
              {loading ? (
                <Skeleton active paragraph={{ rows: 2 }} />
              ) : exec ? (
                <>
                  <div
                    style={{
                      display: 'flex',
                      flexWrap: 'wrap',
                      gap: 12,
                      alignItems: 'center',
                      justifyContent: 'space-between',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        flexWrap: 'wrap',
                        gap: 12,
                        alignItems: 'center',
                      }}
                    >
                      <span
                        style={{
                          padding: '8px 14px',
                          borderRadius: 8,
                          background: 'rgba(16, 185, 129, 0.12)',
                          color: '#10b981',
                          fontSize: 13,
                          fontVariantNumeric: 'tabular-nums',
                          fontWeight: 500,
                        }}
                      >
                        <CheckCircleOutlined style={{ marginRight: 6 }} />
                        成功 {exec.success}
                      </span>
                      <span
                        style={{
                          padding: '8px 14px',
                          borderRadius: 8,
                          background: 'rgba(239, 68, 68, 0.12)',
                          color: '#ef4444',
                          fontSize: 13,
                          fontVariantNumeric: 'tabular-nums',
                          fontWeight: 500,
                        }}
                      >
                        <CloseCircleOutlined style={{ marginRight: 6 }} />
                        失败 {exec.failed}
                      </span>
                      <span
                        style={{
                          padding: '8px 14px',
                          borderRadius: 8,
                          background: 'rgba(59, 130, 246, 0.12)',
                          color: '#3b82f6',
                          fontSize: 13,
                          fontVariantNumeric: 'tabular-nums',
                          fontWeight: 500,
                        }}
                      >
                        <SyncOutlined spin style={{ marginRight: 6 }} />
                        运行中 {exec.running}
                      </span>
                      <span
                        style={{
                          color: token.colorTextSecondary,
                          fontSize: 13,
                        }}
                      >
                        共 {totalExec} 次
                      </span>
                    </div>
                    <span
                      style={{
                        padding: '6px 12px',
                        borderRadius: 8,
                        background:
                          healthPercent >= 90
                            ? 'rgba(16, 185, 129, 0.15)'
                            : healthPercent >= 70
                              ? 'rgba(245, 158, 11, 0.15)'
                              : 'rgba(239, 68, 68, 0.15)',
                        color:
                          healthPercent >= 90
                            ? '#10b981'
                            : healthPercent >= 70
                              ? '#f59e0b'
                              : '#ef4444',
                        fontSize: 13,
                        fontWeight: 600,
                      }}
                    >
                      健康度 {healthPercent}%
                    </span>
                  </div>
                  <div>
                    <div
                      style={{
                        marginBottom: 8,
                        display: 'flex',
                        alignItems: 'center',
                        gap: 6,
                      }}
                    >
                      <ClockCircleOutlined style={{ color: token.colorTextSecondary }} />
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        最近执行
                      </Text>
                    </div>
                    <RecentActivityList
                      items={stats?.recent_activities ?? []}
                      loading={false}
                      token={token}
                    />
                  </div>
                </>
              ) : (
                <Text type="secondary" style={{ fontSize: 13 }}>
                  暂无执行记录
                </Text>
              )}
            </Card>
          </Col>
          <Col xs={24} lg={10}>
            <Card
              title="资产概览"
              size="small"
              style={{
                borderRadius: 12,
                height: '100%',
                boxShadow: token.boxShadowSecondary,
              }}
              styles={{
                body: {
                  padding: 20,
                  height: 'calc(100% - 57px)',
                  overflow: 'auto',
                },
              }}
            >
              {loading ? (
                <Skeleton active paragraph={{ rows: 4 }} />
              ) : (
                <>
                  {/* 资产状态 */}
                  {Object.keys(assetsByStatus).length > 0 && (
                    <div
                      style={{
                        marginBottom: 16,
                        paddingBottom: 12,
                        borderBottom: `1px solid ${token.colorBorderSecondary}`,
                      }}
                    >
                      <Text
                        type="secondary"
                        style={{ fontSize: 12, display: 'flex', alignItems: 'center', gap: 6 }}
                      >
                        <SafetyCertificateOutlined /> 按状态
                      </Text>
                      <div
                        style={{
                          display: 'flex',
                          flexWrap: 'wrap',
                          gap: 8,
                          marginTop: 8,
                        }}
                      >
                        {Object.entries(assetsByStatus).map(([status, count]) => (
                          <span
                            key={status}
                            style={{
                              padding: '6px 12px',
                              borderRadius: 8,
                              background: token.colorFillSecondary,
                              fontSize: 12,
                              fontVariantNumeric: 'tabular-nums',
                              fontWeight: 500,
                            }}
                          >
                            {STATUS_LABELS[status] ?? status} {count}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                  <Row gutter={[20, 16]}>
                    <Col span={12}>
                      <div
                        style={{
                          marginBottom: 12,
                          paddingBottom: 8,
                          borderBottom: `1px solid ${token.colorBorderSecondary}`,
                        }}
                      >
                        <Text
                          type="secondary"
                          style={{ fontSize: 12, display: 'flex', alignItems: 'center', gap: 6 }}
                        >
                          <FolderOutlined /> 按项目
                        </Text>
                      </div>
                      {(stats?.assets_by_project?.slice(0, 5) ?? []).map((p) => (
                        <div
                          key={p.project_id}
                          style={{
                            fontSize: 13,
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            marginBottom: 10,
                            padding: '6px 0',
                          }}
                        >
                          <Text ellipsis style={{ maxWidth: '72%' }}>
                            {p.project_name || '未分配'}
                          </Text>
                          <span
                            style={{
                              fontWeight: 600,
                              fontVariantNumeric: 'tabular-nums',
                              color: token.colorPrimary,
                              fontSize: 13,
                            }}
                          >
                            {p.count}
                          </span>
                        </div>
                      ))}
                      {(!stats?.assets_by_project?.length ||
                        stats.assets_by_project.length === 0) && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          暂无数据
                        </Text>
                      )}
                    </Col>
                    <Col span={12}>
                      <div
                        style={{
                          marginBottom: 12,
                          paddingBottom: 8,
                          borderBottom: `1px solid ${token.colorBorderSecondary}`,
                        }}
                      >
                        <Text
                          type="secondary"
                          style={{ fontSize: 12, display: 'flex', alignItems: 'center', gap: 6 }}
                        >
                          <EnvironmentOutlined /> 按环境
                        </Text>
                      </div>
                      {(stats?.assets_by_environment?.slice(0, 5) ?? []).map((e) => (
                        <div
                          key={e.environment_id}
                          style={{
                            fontSize: 13,
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            marginBottom: 10,
                            padding: '6px 0',
                          }}
                        >
                          <Text ellipsis style={{ maxWidth: '72%' }}>
                            {e.environment_name || '未分配'}
                          </Text>
                          <span
                            style={{
                              fontWeight: 600,
                              fontVariantNumeric: 'tabular-nums',
                              color: token.colorPrimary,
                              fontSize: 13,
                            }}
                          >
                            {e.count}
                          </span>
                        </div>
                      ))}
                      {(!stats?.assets_by_environment?.length ||
                        stats.assets_by_environment.length === 0) && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          暂无数据
                        </Text>
                      )}
                    </Col>
                  </Row>
                </>
              )}
            </Card>
          </Col>
        </Row>
      </div>
    </div>
    </PageContainer>
  )
}

export default Dashboard
