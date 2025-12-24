import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin } from 'antd';
import { UserOutlined, CloudServerOutlined, ProjectOutlined } from '@ant-design/icons';
import { useAuth } from '../contexts/AuthContext';
import { userService } from '../services/user';
import { hostService } from '../services/host';
import { projectService } from '../services/project';

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    users: 0,
    hosts: 0,
    projects: 0,
  });

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    setLoading(true);
    try {
      const [usersRes, hostsRes, projectsRes] = await Promise.all([
        userService.list(1, 1).catch(() => ({ total: 0 })),
        hostService.list(1, 1).catch(() => ({ total: 0 })),
        projectService.list(1, 1).catch(() => ({ total: 0 })),
      ]);

      setStats({
        users: usersRes.total || 0,
        hosts: hostsRes.total || 0,
        projects: projectsRes.total || 0,
      });
    } catch (error) {
      console.error('加载统计数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <h1>欢迎, {user?.display_name || user?.username}</h1>
      <Spin spinning={loading}>
        <Row gutter={16} style={{ marginTop: '24px' }}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="用户总数"
                value={stats.users}
                prefix={<UserOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="主机总数"
                value={stats.hosts}
                prefix={<CloudServerOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="项目总数"
                value={stats.projects}
                prefix={<ProjectOutlined />}
              />
            </Card>
          </Col>
        </Row>
      </Spin>
    </>
  );
};

export default Dashboard;

