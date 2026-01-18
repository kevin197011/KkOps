// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Card, Radio, InputNumber, Checkbox, Space, Typography, Row, Col } from 'antd'
import { ThunderboltOutlined, SyncOutlined, ClockCircleOutlined, SaveOutlined } from '@ant-design/icons'

const { Text } = Typography

interface ExecutionOptionsProps {
  executionType: 'sync' | 'async'
  onExecutionTypeChange: (type: 'sync' | 'async') => void
  timeout: number
  onTimeoutChange: (timeout: number) => void
  saveAsTask: boolean
  onSaveAsTaskChange: (save: boolean) => void
}

const ExecutionOptions: React.FC<ExecutionOptionsProps> = ({
  executionType,
  onExecutionTypeChange,
  timeout,
  onTimeoutChange,
  saveAsTask,
  onSaveAsTaskChange,
}) => {
  return (
    <Card
      style={{
        marginBottom: 16,
        borderRadius: 12,
      }}
      styles={{
        body: { padding: '20px 24px' },
      }}
    >
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <Text strong>
          <ThunderboltOutlined style={{ marginRight: 8 }} />
          执行选项
        </Text>

        <Row gutter={[24, 16]}>
          <Col xs={24} sm={12}>
            <div>
              <Text style={{ marginBottom: 8, display: 'block' }}>
                <SyncOutlined style={{ marginRight: 8 }} />
                执行方式
              </Text>
              <Radio.Group value={executionType} onChange={(e) => onExecutionTypeChange(e.target.value)}>
                <Radio.Button value="sync">同步执行</Radio.Button>
                <Radio.Button value="async">异步执行</Radio.Button>
              </Radio.Group>
              <div style={{ marginTop: 8 }}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {executionType === 'sync'
                    ? '等待所有主机执行完成后返回结果'
                    : '立即返回，后台执行，可查看实时日志'}
                </Text>
              </div>
            </div>
          </Col>

          <Col xs={24} sm={12}>
            <div>
              <Text style={{ marginBottom: 8, display: 'block' }}>
                <ClockCircleOutlined style={{ marginRight: 8 }} />
                超时时间（秒）
              </Text>
              <InputNumber
                value={timeout}
                onChange={(value) => onTimeoutChange(value || 600)}
                min={60}
                max={3600}
                step={60}
                style={{ width: '100%' }}
                addonAfter="秒"
              />
              <div style={{ marginTop: 8 }}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  单个主机执行超时时间，默认 600 秒（10 分钟）
                </Text>
              </div>
            </div>
          </Col>
        </Row>

        <div>
          <Checkbox
            checked={saveAsTask}
            onChange={(e) => onSaveAsTaskChange(e.target.checked)}
          >
            <Space>
              <SaveOutlined />
              <Text>保存为任务</Text>
            </Space>
          </Checkbox>
          <div style={{ marginTop: 8 }}>
            <Text type="secondary" style={{ fontSize: 12 }}>
              勾选后，执行完成后将自动保存为任务，可在"任务执行"页面查看和管理
            </Text>
          </div>
        </div>
      </Space>
    </Card>
  )
}

export default ExecutionOptions
