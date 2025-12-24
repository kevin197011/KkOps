import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Select,
  DatePicker,
  Input,
  Tag,
  Modal,
  Descriptions,
  message,
  Tooltip,
} from 'antd';
import {
  ReloadOutlined,
  EyeOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import { auditService, AuditLog, AuditLogFilters } from '../services/audit';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';

const { Option } = Select;
const { RangePicker } = DatePicker;

const Audit: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState<AuditLogFilters>({});
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);

  useEffect(() => {
    loadLogs();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, filters]);

  const loadLogs = async () => {
    setLoading(true);
    try {
      const response = await auditService.list(page, pageSize, filters);
      setLogs(response.logs);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载审计日志失败');
    } finally {
      setLoading(false);
    }
  };

  const handleViewDetail = async (log: AuditLog) => {
    try {
      const response = await auditService.get(log.id);
      setSelectedLog(response.log);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取日志详情失败');
    }
  };

  const handleFilterChange = (key: string, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
    }));
    setPage(1);
  };

  const handleResetFilters = () => {
    setFilters({});
    setPage(1);
  };

  const getStatusTag = (status: string) => {
    return status === 'success' ? (
      <Tag color="success">成功</Tag>
    ) : (
      <Tag color="error">失败</Tag>
    );
  };

  const getActionTag = (action: string) => {
    const actionMap: Record<string, { color: string; text: string }> = {
      login: { color: 'blue', text: '登录' },
      logout: { color: 'default', text: '登出' },
      create: { color: 'green', text: '创建' },
      update: { color: 'orange', text: '更新' },
      delete: { color: 'red', text: '删除' },
      query: { color: 'cyan', text: '查询' },
      execute: { color: 'purple', text: '执行' },
    };
    const actionInfo = actionMap[action] || { color: 'default', text: action };
    return <Tag color={actionInfo.color}>{actionInfo.text}</Tag>;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
      sorter: true,
    },
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      render: (action: string) => getActionTag(action),
    },
    {
      title: '资源类型',
      dataIndex: 'resource_type',
      key: 'resource_type',
      width: 120,
      render: (text: string) => text || '-',
    },
    {
      title: '资源ID',
      dataIndex: 'resource_id',
      key: 'resource_id',
      width: 100,
      render: (id: number) => id || '-',
    },
    {
      title: '资源名称',
      dataIndex: 'resource_name',
      key: 'resource_name',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: 'IP地址',
      dataIndex: 'ip_address',
      key: 'ip_address',
      width: 130,
      render: (text: string) => text || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '耗时',
      dataIndex: 'duration_ms',
      key: 'duration_ms',
      width: 80,
      render: (ms: number) => (ms ? `${ms}ms` : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right' as const,
      render: (_: any, record: AuditLog) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record)}
        >
          详情
        </Button>
      ),
    },
  ];

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px' }}>
          <h2 style={{ margin: 0, marginBottom: '16px' }}>审计日志</h2>
          
          <Space wrap style={{ marginBottom: '16px' }}>
            <Input
              placeholder="用户名"
              style={{ width: 150 }}
              value={filters.username}
              onChange={(e) => handleFilterChange('username', e.target.value || undefined)}
              allowClear
            />
            <Select
              placeholder="操作类型"
              style={{ width: 150 }}
              value={filters.action}
              onChange={(value) => handleFilterChange('action', value)}
              allowClear
            >
              <Option value="login">登录</Option>
              <Option value="logout">登出</Option>
              <Option value="create">创建</Option>
              <Option value="update">更新</Option>
              <Option value="delete">删除</Option>
              <Option value="query">查询</Option>
              <Option value="execute">执行</Option>
            </Select>
            <Select
              placeholder="资源类型"
              style={{ width: 150 }}
              value={filters.resource_type}
              onChange={(value) => handleFilterChange('resource_type', value)}
              allowClear
            >
              <Option value="host">主机</Option>
              <Option value="ssh">SSH</Option>
              <Option value="user">用户</Option>
              <Option value="project">项目</Option>
              <Option value="deployment">部署</Option>
            </Select>
            <Select
              placeholder="状态"
              style={{ width: 120 }}
              value={filters.status}
              onChange={(value) => handleFilterChange('status', value)}
              allowClear
            >
              <Option value="success">成功</Option>
              <Option value="failed">失败</Option>
            </Select>
            <RangePicker
              showTime
              onChange={(dates) => {
                if (dates && dates[0] && dates[1]) {
                  handleFilterChange('start_time', dates[0].toISOString());
                  handleFilterChange('end_time', dates[1].toISOString());
                } else {
                  handleFilterChange('start_time', undefined);
                  handleFilterChange('end_time', undefined);
                }
              }}
            />
            <Button
              type="primary"
              icon={<SearchOutlined />}
              onClick={loadLogs}
            >
              搜索
            </Button>
            <Button onClick={handleResetFilters}>重置</Button>
            <Button icon={<ReloadOutlined />} onClick={loadLogs}>
              刷新
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={logs}
          rowKey="id"
          loading={loading}
          scroll={{ x: 1400 }}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
            onChange: (page, pageSize) => {
              setPage(page);
              setPageSize(pageSize);
            },
          }}
        />
      </Card>

      {/* 详情模态框 */}
      <Modal
        title="审计日志详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={800}
      >
        {selectedLog && (
          <Descriptions bordered column={2}>
            <Descriptions.Item label="ID" span={1}>
              {selectedLog.id}
            </Descriptions.Item>
            <Descriptions.Item label="时间" span={1}>
              {new Date(selectedLog.created_at).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="用户" span={1}>
              {selectedLog.username || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="用户ID" span={1}>
              {selectedLog.user_id || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="操作" span={1}>
              {getActionTag(selectedLog.action)}
            </Descriptions.Item>
            <Descriptions.Item label="状态" span={1}>
              {getStatusTag(selectedLog.status)}
            </Descriptions.Item>
            <Descriptions.Item label="资源类型" span={1}>
              {selectedLog.resource_type || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="资源ID" span={1}>
              {selectedLog.resource_id || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="资源名称" span={2}>
              {selectedLog.resource_name || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="IP地址" span={1}>
              {selectedLog.ip_address || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="耗时" span={1}>
              {selectedLog.duration_ms ? `${selectedLog.duration_ms}ms` : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="User Agent" span={2}>
              <Tooltip title={selectedLog.user_agent}>
                <div style={{ maxWidth: 500, overflow: 'hidden', textOverflow: 'ellipsis' }}>
                  {selectedLog.user_agent || '-'}
                </div>
              </Tooltip>
            </Descriptions.Item>
            {selectedLog.error_message && (
              <Descriptions.Item label="错误信息" span={2}>
                <div style={{ color: 'red' }}>{selectedLog.error_message}</div>
              </Descriptions.Item>
            )}
            {selectedLog.request_data && (
              <Descriptions.Item label="请求数据" span={2}>
                <pre style={{ maxHeight: 200, overflow: 'auto', background: '#f5f5f5', padding: 8, borderRadius: 4 }}>
                  {JSON.stringify(JSON.parse(selectedLog.request_data), null, 2)}
                </pre>
              </Descriptions.Item>
            )}
            {selectedLog.response_data && (
              <Descriptions.Item label="响应数据" span={2}>
                <pre style={{ maxHeight: 200, overflow: 'auto', background: '#f5f5f5', padding: 8, borderRadius: 4 }}>
                  {JSON.stringify(JSON.parse(selectedLog.response_data), null, 2)}
                </pre>
              </Descriptions.Item>
            )}
            {selectedLog.before_data && (
              <Descriptions.Item label="变更前数据" span={2}>
                <pre style={{ maxHeight: 200, overflow: 'auto', background: '#f5f5f5', padding: 8, borderRadius: 4 }}>
                  {JSON.stringify(JSON.parse(selectedLog.before_data), null, 2)}
                </pre>
              </Descriptions.Item>
            )}
            {selectedLog.after_data && (
              <Descriptions.Item label="变更后数据" span={2}>
                <pre style={{ maxHeight: 200, overflow: 'auto', background: '#f5f5f5', padding: 8, borderRadius: 4 }}>
                  {JSON.stringify(JSON.parse(selectedLog.after_data), null, 2)}
                </pre>
              </Descriptions.Item>
            )}
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default Audit;

