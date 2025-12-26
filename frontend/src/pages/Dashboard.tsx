import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Row, 
  Col, 
  Statistic, 
  Spin, 
  Tag, 
  Typography,
  Divider,
  Empty
} from 'antd';
import {
  UserOutlined,
  CloudServerOutlined,
  ProjectOutlined,
  KeyOutlined,
  AppstoreOutlined,
  RocketOutlined,
  FileTextOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  QuestionCircleOutlined,
  ClockCircleOutlined,
  PlayCircleOutlined,
  StopOutlined,
} from '@ant-design/icons';
import { Progress } from 'antd';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { userService } from '../services/user';
import { hostService, Host } from '../services/host';
import { projectService } from '../services/project';
import { sshKeyService } from '../services/sshKey';
import { hostGroupService } from '../services/hostGroup';
import { deploymentService, Deployment } from '../services/deployment';
import { auditService, AuditLog } from '../services/audit';

const { Title } = Typography;

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    users: 0,
    hosts: 0,
    hostsOnline: 0,
    hostsOffline: 0,
    hostsUnknown: 0,
    projects: 0,
    sshKeys: 0,
    hostGroups: 0,
    deploymentConfigs: 0,
    deployments: 0,
    deploymentsCompleted: 0,
    deploymentsFailed: 0,
    deploymentsRunning: 0,
    deploymentsPending: 0,
  });
  const [recentLogs, setRecentLogs] = useState<AuditLog[]>([]);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    setLoading(true);
    try {
      const [
        usersRes,
        hostsRes,
        projectsRes,
        sshKeysRes,
        hostGroupsRes,
        deploymentConfigsRes,
        deploymentsRes,
        auditLogsRes,
      ] = await Promise.all([
        userService.list(1, 1).catch(() => ({ total: 0 })),
        hostService.list(1, 100).catch(() => ({ hosts: [], total: 0 })),
        projectService.list(1, 1).catch(() => ({ total: 0 })),
        sshKeyService.list(1, 1).catch(() => ({ total: 0 })),
        hostGroupService.list(1, 1).catch(() => ({ total: 0 })),
        deploymentService.listConfigs(1, 1).catch(() => ({ total: 0 })),
        deploymentService.listDeployments(1, 100).catch(() => ({ deployments: [], total: 0 })),
        auditService.list(1, 10).catch(() => ({ logs: [], total: 0 })),
      ]);

      // 统计主机状态
      const hosts = (hostsRes as any).hosts || [];
      const hostsOnline = hosts.filter((h: Host) => h.status === 'online').length;
      const hostsOffline = hosts.filter((h: Host) => h.status === 'offline').length;
      const hostsUnknown = hosts.filter((h: Host) => h.status === 'unknown').length;

      // 统计部署状态
      const deployments = (deploymentsRes as any).deployments || [];
      const deploymentsCompleted = deployments.filter((d: Deployment) => d.status === 'completed').length;
      const deploymentsFailed = deployments.filter((d: Deployment) => d.status === 'failed').length;
      const deploymentsRunning = deployments.filter((d: Deployment) => d.status === 'running').length;
      const deploymentsPending = deployments.filter((d: Deployment) => d.status === 'pending').length;

      setStats({
        users: (usersRes as any).total || 0,
        hosts: hostsRes.total || 0,
        hostsOnline,
        hostsOffline,
        hostsUnknown,
        projects: (projectsRes as any).total || 0,
        sshKeys: (sshKeysRes as any).total || 0,
        hostGroups: (hostGroupsRes as any).total || 0,
        deploymentConfigs: (deploymentConfigsRes as any).total || 0,
        deployments: deploymentsRes.total || 0,
        deploymentsCompleted,
        deploymentsFailed,
        deploymentsRunning,
        deploymentsPending,
      });

      // 最近的活动日志
      const logs = (auditLogsRes as any).logs || [];
      setRecentLogs(logs);
    } catch (error) {
      console.error('加载统计数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 主机状态分布数据
  const hostStatusData = [
    { type: '在线', value: stats.hostsOnline, color: '#52c41a' },
    { type: '离线', value: stats.hostsOffline, color: '#ff4d4f' },
    { type: '未知', value: stats.hostsUnknown, color: '#faad14' },
  ].filter(item => item.value > 0);

  // 部署状态分布数据
  const deploymentStatusData = [
    { type: '已完成', value: stats.deploymentsCompleted, color: '#52c41a' },
    { type: '失败', value: stats.deploymentsFailed, color: '#ff4d4f' },
    { type: '运行中', value: stats.deploymentsRunning, color: '#1890ff' },
    { type: '等待中', value: stats.deploymentsPending, color: '#faad14' },
  ].filter(item => item.value > 0);

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      success: '#52c41a',
      failed: '#ff4d4f',
      running: '#1890ff',
      pending: '#faad14',
    };
    return colors[status] || '#d9d9d9';
  };

  const getStatusIcon = (status: string) => {
    const icons: Record<string, React.ReactNode> = {
      success: <CheckCircleOutlined style={{ color: '#52c41a' }} />,
      failed: <CloseCircleOutlined style={{ color: '#ff4d4f' }} />,
      running: <PlayCircleOutlined style={{ color: '#1890ff' }} />,
      pending: <ClockCircleOutlined style={{ color: '#faad14' }} />,
      cancelled: <StopOutlined style={{ color: '#d9d9d9' }} />,
    };
    return icons[status] || <QuestionCircleOutlined />;
  };

  return (
    <div>
      <Title level={2}>欢迎回来, {user?.display_name || user?.username}</Title>
      <Spin spinning={loading}>
        {/* 基础统计卡片 */}
        <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/users')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="用户总数"
                value={stats.users}
                prefix={<UserOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/hosts')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="主机总数"
                value={stats.hosts}
                prefix={<CloudServerOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
              <div style={{ marginTop: 8, fontSize: '12px' }}>
                <Tag color="success">在线 {stats.hostsOnline}</Tag>
                <Tag color="error">离线 {stats.hostsOffline}</Tag>
                <Tag color="warning">未知 {stats.hostsUnknown}</Tag>
              </div>
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/projects')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="项目总数"
                value={stats.projects}
                prefix={<ProjectOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/ssh-keys')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="SSH密钥"
                value={stats.sshKeys}
                prefix={<KeyOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={[16, 16]} style={{ marginTop: '16px' }}>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/host-groups')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="主机组"
                value={stats.hostGroups}
                prefix={<AppstoreOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/deployments')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="部署配置"
                value={stats.deploymentConfigs}
                prefix={<RocketOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card 
              hoverable
              onClick={() => navigate('/deployments')}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title="部署记录"
                value={stats.deployments}
                prefix={<FileTextOutlined style={{ color: '#1890ff' }} />}
                valueStyle={{ color: '#1890ff' }}
              />
              <div style={{ marginTop: 8, fontSize: '12px' }}>
                <Tag color="success">成功 {stats.deploymentsCompleted}</Tag>
                <Tag color="error">失败 {stats.deploymentsFailed}</Tag>
                {stats.deploymentsRunning > 0 && (
                  <Tag color="processing">运行中 {stats.deploymentsRunning}</Tag>
                )}
              </div>
            </Card>
          </Col>
        </Row>

        {/* 图表区域 */}
        <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
          <Col xs={24} lg={12}>
            <Card title="主机状态分布">
              {hostStatusData.length > 0 ? (
                <div>
                  {hostStatusData.map((item, index) => {
                    const percentage = stats.hosts > 0 
                      ? Math.round((item.value / stats.hosts) * 100) 
                      : 0;
                    return (
                      <div key={index} style={{ marginBottom: 16 }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                          <span>
                            <span 
                              style={{ 
                                display: 'inline-block', 
                                width: 12, 
                                height: 12, 
                                borderRadius: 2,
                                backgroundColor: item.color,
                                marginRight: 8 
                              }} 
                            />
                            {item.type}
                          </span>
                          <span>
                            <strong>{item.value}</strong>
                            <span style={{ marginLeft: 4, color: '#999' }}>({percentage}%)</span>
                          </span>
                        </div>
                        <Progress 
                          percent={percentage} 
                          strokeColor={item.color}
                          showInfo={false}
                        />
                      </div>
                    );
                  })}
                </div>
              ) : (
                <Empty description="暂无主机数据" style={{ padding: '40px 0' }} />
              )}
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card title="部署状态分布">
              {deploymentStatusData.length > 0 ? (
                <div>
                  {deploymentStatusData.map((item, index) => {
                    const percentage = stats.deployments > 0 
                      ? Math.round((item.value / stats.deployments) * 100) 
                      : 0;
                    return (
                      <div key={index} style={{ marginBottom: 16 }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                          <span>
                            <span 
                              style={{ 
                                display: 'inline-block', 
                                width: 12, 
                                height: 12, 
                                borderRadius: 2,
                                backgroundColor: item.color,
                                marginRight: 8 
                              }} 
                            />
                            {item.type}
                          </span>
                          <span>
                            <strong>{item.value}</strong>
                            <span style={{ marginLeft: 4, color: '#999' }}>({percentage}%)</span>
                          </span>
                        </div>
                        <Progress 
                          percent={percentage} 
                          strokeColor={item.color}
                          showInfo={false}
                        />
                      </div>
                    );
                  })}
                </div>
              ) : (
                <Empty description="暂无部署数据" style={{ padding: '40px 0' }} />
              )}
            </Card>
          </Col>
        </Row>

        {/* 最近活动 */}
        <Card 
          title="最近活动" 
          extra={
            <a onClick={() => navigate('/audit')} style={{ cursor: 'pointer' }}>
              查看全部
            </a>
          }
          style={{ marginTop: '24px' }}
        >
          {recentLogs.length > 0 ? (
            <div>
              {recentLogs.map((log) => (
                <div key={log.id} style={{ padding: '12px 0', borderBottom: '1px solid #f0f0f0' }}>
                  <Row justify="space-between" align="middle">
                    <Col span={16}>
                      <div>
                        <strong>{log.username || '系统'}</strong>
                        <span style={{ margin: '0 8px', color: '#999' }}>执行了</span>
                        <Tag color="blue">{log.action}</Tag>
                        <span style={{ margin: '0 8px', color: '#999' }}>操作</span>
                        {log.resource_type && (
                          <>
                            <span style={{ margin: '0 8px', color: '#999' }}>于</span>
                            <Tag>{log.resource_type}</Tag>
                          </>
                        )}
                        {log.resource_name && (
                          <span style={{ marginLeft: 8 }}>"{log.resource_name}"</span>
                        )}
                      </div>
                      <div style={{ marginTop: 4, fontSize: '12px', color: '#999' }}>
                        {new Date(log.created_at).toLocaleString('zh-CN')}
                      </div>
                    </Col>
                    <Col span={4} style={{ textAlign: 'right' }}>
                      {getStatusIcon(log.status)}
                    </Col>
                  </Row>
                </div>
              ))}
            </div>
          ) : (
            <Empty description="暂无活动记录" style={{ padding: '40px 0' }} />
          )}
        </Card>
      </Spin>
    </div>
  );
};

export default Dashboard;

