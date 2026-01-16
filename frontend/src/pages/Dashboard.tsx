// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Card, Row, Col, Statistic } from 'antd'
import { DatabaseOutlined, UserOutlined, PlayCircleOutlined, ConsoleSqlOutlined } from '@ant-design/icons'

const Dashboard = () => {
  return (
    <div>
      <h2>仪表板</h2>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card>
            <Statistic title="总资产数" value={0} prefix={<DatabaseOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card>
            <Statistic title="总用户数" value={0} prefix={<UserOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card>
            <Statistic title="执行任务" value={0} prefix={<PlayCircleOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={12} lg={6}>
          <Card>
            <Statistic title="SSH会话" value={0} prefix={<ConsoleSqlOutlined />} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
