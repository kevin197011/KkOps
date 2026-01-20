// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import {
  Table,
  Card,
  Form,
  Input,
  Select,
  DatePicker,
  Button,
  Tag,
  Space,
  Modal,
  Descriptions,
  message,
  Row,
  Col,
  Tooltip,
  Typography,
  theme,
} from 'antd'
import {
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  EyeOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  FilterOutlined,
  ClearOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import dayjs from 'dayjs'
import {
  auditApi,
  AuditLog,
  AuditLogQueryParams,
  moduleLabels,
  actionLabels,
  statusLabels,
} from '@/api/audit'

const { RangePicker } = DatePicker
const { Text } = Typography

const AuditLogList = () => {
  const { token } = theme.useToken()
  const [loading, setLoading] = useState(false)
  const [logs, setLogs] = useState<AuditLog[]>([])
  const [total, setTotal] = useState(0)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20 })
  const [detailVisible, setDetailVisible] = useState(false)
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null)
  const [modules, setModules] = useState<string[]>([])
  const [actions, setActions] = useState<string[]>([])
  const [form] = Form.useForm()

  // 获取模块和操作类型列表
  useEffect(() => {
    const fetchOptions = async () => {
      try {
        const [modulesRes, actionsRes] = await Promise.all([
          auditApi.getModules(),
          auditApi.getActions(),
        ])
        setModules(modulesRes.data.data)
        setActions(actionsRes.data.data)
      } catch (error) {
        console.error('获取选项列表失败:', error)
      }
    }
    fetchOptions()
  }, [])

  // 获取审计日志列表
  const fetchLogs = useCallback(async () => {
    setLoading(true)
    try {
      const values = form.getFieldsValue()
      const params: AuditLogQueryParams = {
        page: pagination.current,
        page_size: pagination.pageSize,
        username: values.username,
        module: values.module,
        action: values.action,
        status: values.status,
        keyword: values.keyword,
      }

      // 处理时间范围
      if (values.timeRange && values.timeRange.length === 2) {
        params.start_time = values.timeRange[0].startOf('day').toISOString()
        params.end_time = values.timeRange[1].endOf('day').toISOString()
      }

      const response = await auditApi.list(params)
      setLogs(response.data.data.items || [])
      setTotal(response.data.data.total)
    } catch (error) {
      message.error('获取审计日志失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [form, pagination])

  useEffect(() => {
    fetchLogs()
  }, [fetchLogs])

  // 处理表格分页变化
  const handleTableChange = (newPagination: any) => {
    setPagination({
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    })
  }

  // 处理搜索
  const handleSearch = () => {
    setPagination({ ...pagination, current: 1 })
    fetchLogs()
  }

  // 重置筛选
  const handleReset = () => {
    form.resetFields()
    setPagination({ current: 1, pageSize: 20 })
  }

  // 导出
  const handleExport = async (format: 'csv' | 'json') => {
    try {
      setLoading(true)
      const values = form.getFieldsValue()
      const params: AuditLogQueryParams & { format: 'csv' | 'json' } = {
        format,
        username: values.username,
        module: values.module,
        action: values.action,
        status: values.status,
        keyword: values.keyword,
      }

      if (values.timeRange && values.timeRange.length === 2) {
        params.start_time = values.timeRange[0].startOf('day').toISOString()
        params.end_time = values.timeRange[1].endOf('day').toISOString()
      }

      await auditApi.export(params)
      message.success('导出成功')
    } catch (error: any) {
      console.error('导出失败:', error)
      message.error(error.response?.data?.error || '导出失败')
    } finally {
      setLoading(false)
    }
  }

  // 查看详情
  const handleViewDetail = (record: AuditLog) => {
    setSelectedLog(record)
    setDetailVisible(true)
  }

  // 格式化详情 JSON
  const formatDetail = (detail: string | undefined) => {
    if (!detail) return '-'
    try {
      const parsed = JSON.parse(detail)
      return JSON.stringify(parsed, null, 2)
    } catch {
      return detail
    }
  }

  // 表格列定义
  const columns: ColumnsType<AuditLog> = [
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (val) => dayjs(val).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 100,
    },
    {
      title: '模块',
      dataIndex: 'module',
      key: 'module',
      width: 100,
      render: (val) => (
        <Tag color="blue">{moduleLabels[val] || val}</Tag>
      ),
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 80,
      render: (val) => actionLabels[val] || val,
    },
    {
      title: '资源',
      dataIndex: 'resource_name',
      key: 'resource_name',
      width: 150,
      ellipsis: true,
      render: (val, record) => (
        <Tooltip title={val || (record.resource_id ? `ID: ${record.resource_id}` : '-')}>
          {val || (record.resource_id ? `ID: ${record.resource_id}` : '-')}
        </Tooltip>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (val) => (
        <Tag
          icon={val === 'success' ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
          color={val === 'success' ? 'success' : 'error'}
        >
          {statusLabels[val] || val}
        </Tag>
      ),
    },
    {
      title: 'IP 地址',
      dataIndex: 'ip_address',
      key: 'ip_address',
      width: 130,
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      fixed: 'right',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record)}
        >
          详情
        </Button>
      ),
    },
  ]

  return (
    <div>
      {/* 筛选表单 */}
      <Card style={{ marginBottom: 16 }}>
        <Form form={form} layout="inline" style={{ flexWrap: 'wrap', gap: 8 }}>
          <Row gutter={[16, 16]} style={{ width: '100%' }}>
            <Col xs={24} sm={12} md={6} lg={4}>
              <Form.Item name="username" style={{ marginBottom: 0, width: '100%' }}>
                <Input placeholder="用户名" allowClear />
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} md={6} lg={4}>
              <Form.Item name="module" style={{ marginBottom: 0, width: '100%' }}>
                <Select placeholder="模块" allowClear style={{ width: '100%' }}>
                  {modules.map((m) => (
                    <Select.Option key={m} value={m}>
                      {moduleLabels[m] || m}
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} md={6} lg={4}>
              <Form.Item name="action" style={{ marginBottom: 0, width: '100%' }}>
                <Select placeholder="操作类型" allowClear style={{ width: '100%' }}>
                  {actions.map((a) => (
                    <Select.Option key={a} value={a}>
                      {actionLabels[a] || a}
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} md={6} lg={4}>
              <Form.Item name="status" style={{ marginBottom: 0, width: '100%' }}>
                <Select placeholder="状态" allowClear style={{ width: '100%' }}>
                  <Select.Option value="success">成功</Select.Option>
                  <Select.Option value="failed">失败</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} md={8} lg={5}>
              <Form.Item name="timeRange" style={{ marginBottom: 0, width: '100%' }}>
                <RangePicker style={{ width: '100%' }} placeholder={['开始日期', '结束日期']} />
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} md={6} lg={3}>
              <Form.Item name="keyword" style={{ marginBottom: 0, width: '100%' }}>
                <Input placeholder="关键词" allowClear />
              </Form.Item>
            </Col>
          </Row>
          <Row style={{ width: '100%', marginTop: 16 }}>
            <Col>
              <Space wrap>
                <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                  搜索
                </Button>
                <Button icon={<ClearOutlined />} onClick={handleReset}>
                  重置
                </Button>
                <Button icon={<ReloadOutlined />} onClick={fetchLogs}>
                  刷新
                </Button>
                <Button icon={<DownloadOutlined />} onClick={() => handleExport('csv')}>
                  导出 CSV
                </Button>
                <Button icon={<DownloadOutlined />} onClick={() => handleExport('json')}>
                  导出 JSON
                </Button>
              </Space>
            </Col>
          </Row>
        </Form>
      </Card>

      {/* 数据表格 */}
      <Card>
        <Table
          columns={columns}
          dataSource={logs}
          rowKey="id"
          loading={loading}
          pagination={{
            ...pagination,
            total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
          }}
          onChange={handleTableChange}
          scroll={{ x: 1000 }}
        />
      </Card>

      {/* 详情弹窗 */}
      <Modal
        title="审计日志详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
        width={700}
      >
        {selectedLog && (
          <Descriptions bordered column={2} size="small">
            <Descriptions.Item label="ID">{selectedLog.id}</Descriptions.Item>
            <Descriptions.Item label="时间">
              {dayjs(selectedLog.created_at).format('YYYY-MM-DD HH:mm:ss')}
            </Descriptions.Item>
            <Descriptions.Item label="用户">{selectedLog.username}</Descriptions.Item>
            <Descriptions.Item label="用户 ID">{selectedLog.user_id}</Descriptions.Item>
            <Descriptions.Item label="模块">
              <Tag color="blue">{moduleLabels[selectedLog.module] || selectedLog.module}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="操作">
              {actionLabels[selectedLog.action] || selectedLog.action}
            </Descriptions.Item>
            <Descriptions.Item label="资源 ID">{selectedLog.resource_id || '-'}</Descriptions.Item>
            <Descriptions.Item label="资源名称">{selectedLog.resource_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="状态" span={2}>
              <Tag
                icon={selectedLog.status === 'success' ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
                color={selectedLog.status === 'success' ? 'success' : 'error'}
              >
                {statusLabels[selectedLog.status] || selectedLog.status}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="IP 地址">{selectedLog.ip_address}</Descriptions.Item>
            <Descriptions.Item label="User Agent" span={2}>
              <Text style={{ wordBreak: 'break-all', fontSize: 12 }}>
                {selectedLog.user_agent || '-'}
              </Text>
            </Descriptions.Item>
            {selectedLog.error_msg && (
              <Descriptions.Item label="错误信息" span={2}>
                <Text type="danger">{selectedLog.error_msg}</Text>
              </Descriptions.Item>
            )}
            <Descriptions.Item label="详情" span={2}>
              <pre
                style={{
                  maxHeight: 300,
                  overflow: 'auto',
                  margin: 0,
                  padding: 12,
                  background: token.colorFillTertiary,
                  borderRadius: 4,
                  fontSize: 12,
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-all',
                  border: `1px solid ${token.colorBorderSecondary}`,
                }}
              >
                {formatDetail(selectedLog.detail)}
              </pre>
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  )
}

export default AuditLogList
