import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Select,
  DatePicker,
  Tag,
  Modal,
  Descriptions,
  message,
  Tooltip,
  Typography,
  Divider,
  Row,
  Col,
} from 'antd';
import {
  ReloadOutlined,
  EyeOutlined,
  SearchOutlined,
  UserOutlined,
  ClockCircleOutlined,
  GlobalOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  FieldTimeOutlined,
  LaptopOutlined,
  FileTextOutlined,
  SwapOutlined,
} from '@ant-design/icons';
import { auditService, AuditLog, AuditLogFilters } from '../services/audit';
import { userService } from '../services/user';
import { User } from '../services/auth';

const { Option } = Select;
const { RangePicker } = DatePicker;
const { Text, Title } = Typography;

// 安全解析 JSON 并格式化显示
const formatJsonData = (data: string | null | undefined): string => {
  if (!data || data === '{}' || data === '[]') {
    return '';
  }
  try {
    const parsed = JSON.parse(data);
    // 如果解析后是空对象或空数组，返回空
    if (typeof parsed === 'object' && parsed !== null) {
      if (Array.isArray(parsed) && parsed.length === 0) return '';
      if (Object.keys(parsed).length === 0) return '';
    }
    return JSON.stringify(parsed, null, 2);
  } catch {
    // 如果解析失败，返回原始数据
    return data;
  }
};

// 检查数据是否有有效内容
const hasValidData = (data: string | null | undefined): boolean => {
  if (!data || data === '{}' || data === '[]') {
    return false;
  }
  try {
    const parsed = JSON.parse(data);
    if (typeof parsed === 'object' && parsed !== null) {
      if (Array.isArray(parsed) && parsed.length === 0) return false;
      if (Object.keys(parsed).length === 0) return false;
    }
    return true;
  } catch {
    return data.length > 0;
  }
};

const Audit: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState<AuditLogFilters>({});
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    loadUsers();
  }, []);

  useEffect(() => {
    loadLogs();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, filters]);

  const loadUsers = async () => {
    try {
      const response = await userService.list(1, 1000);
      setUsers(response.users || []);
    } catch (error) {
      console.error('加载用户列表失败:', error);
    }
  };

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
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
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
      render: (text: string) => text || 'system',
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
            <Select
              placeholder="选择用户"
              style={{ width: 150 }}
              value={filters.username}
              onChange={(value) => handleFilterChange('username', value)}
              allowClear
              showSearch
              optionFilterProp="children"
            >
              <Option value="system">system</Option>
              {users.map((user) => (
                <Option key={user.id} value={user.username}>
                  {user.username}
                </Option>
              ))}
            </Select>
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
            showQuickJumper: true,
            pageSizeOptions: ['10', '20', '50', '100'],
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: (page, pageSize) => {
              setPage(page);
              setPageSize(pageSize);
            },
          }}
        />
      </Card>

      {/* 详情模态框 */}
      <Modal
        title={
          <Space>
            <FileTextOutlined />
            <span>审计日志详情</span>
          </Space>
        }
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={
          <Button onClick={() => setDetailModalVisible(false)}>关闭</Button>
        }
        width={720}
      >
        {selectedLog && (
          <div>
            {/* 基本信息卡片 */}
            <Card size="small" style={{ marginBottom: 16 }}>
              <Row gutter={[16, 12]}>
                <Col span={12}>
                  <Space>
                    <ClockCircleOutlined style={{ color: '#1890ff' }} />
                    <Text type="secondary">时间：</Text>
                    <Text strong>{new Date(selectedLog.created_at).toLocaleString()}</Text>
                  </Space>
                </Col>
                <Col span={12}>
                  <Space>
                    <UserOutlined style={{ color: '#52c41a' }} />
                    <Text type="secondary">用户：</Text>
                    <Text strong>{selectedLog.username || 'system'}</Text>
                  </Space>
                </Col>
                <Col span={12}>
                  <Space>
                    <SwapOutlined style={{ color: '#722ed1' }} />
                    <Text type="secondary">操作：</Text>
                    {getActionTag(selectedLog.action)}
                  </Space>
                </Col>
                <Col span={12}>
                  <Space>
                    {selectedLog.status === 'success' ? (
                      <CheckCircleOutlined style={{ color: '#52c41a' }} />
                    ) : (
                      <CloseCircleOutlined style={{ color: '#ff4d4f' }} />
                    )}
                    <Text type="secondary">状态：</Text>
                    {getStatusTag(selectedLog.status)}
                  </Space>
                </Col>
                <Col span={12}>
                  <Space>
                    <FileTextOutlined style={{ color: '#fa8c16' }} />
                    <Text type="secondary">资源类型：</Text>
                    <Text>{selectedLog.resource_type || '-'}</Text>
                  </Space>
                </Col>
                <Col span={12}>
                  <Space>
                    <FieldTimeOutlined style={{ color: '#13c2c2' }} />
                    <Text type="secondary">耗时：</Text>
                    <Text>{selectedLog.duration_ms ? `${selectedLog.duration_ms}ms` : '-'}</Text>
                  </Space>
                </Col>
                <Col span={24}>
                  <Space>
                    <GlobalOutlined style={{ color: '#eb2f96' }} />
                    <Text type="secondary">IP地址：</Text>
                    <Text code>{selectedLog.ip_address || '-'}</Text>
                  </Space>
                </Col>
                {selectedLog.user_agent && (
                  <Col span={24}>
                    <Space align="start">
                      <LaptopOutlined style={{ color: '#595959' }} />
                      <Text type="secondary">客户端：</Text>
                      <Tooltip title={selectedLog.user_agent}>
                        <Text 
                          style={{ 
                            maxWidth: 500, 
                            display: 'inline-block',
                            overflow: 'hidden', 
                            textOverflow: 'ellipsis', 
                            whiteSpace: 'nowrap',
                            verticalAlign: 'middle'
                          }}
                          type="secondary"
                        >
                          {selectedLog.user_agent}
                        </Text>
                      </Tooltip>
                    </Space>
                  </Col>
                )}
              </Row>
            </Card>

            {/* 错误信息 */}
            {selectedLog.error_message && (
              <Card 
                size="small" 
                title={<Text type="danger">错误信息</Text>}
                style={{ marginBottom: 16, borderColor: '#ffccc7' }}
                styles={{ body: { background: '#fff2f0' } }}
              >
                <Text type="danger">{selectedLog.error_message}</Text>
              </Card>
            )}

            {/* 请求数据 */}
            {hasValidData(selectedLog.request_data) && (
              <Card 
                size="small" 
                title="请求数据"
                style={{ marginBottom: 16 }}
              >
                <pre style={{ 
                  maxHeight: 180, 
                  overflow: 'auto', 
                  background: '#fafafa', 
                  padding: 12, 
                  borderRadius: 6,
                  margin: 0,
                  fontSize: 12,
                  border: '1px solid #f0f0f0'
                }}>
                  {formatJsonData(selectedLog.request_data)}
                </pre>
              </Card>
            )}

            {/* 响应数据 */}
            {hasValidData(selectedLog.response_data) && (
              <Card 
                size="small" 
                title="响应数据"
                style={{ marginBottom: 16 }}
              >
                <pre style={{ 
                  maxHeight: 180, 
                  overflow: 'auto', 
                  background: '#fafafa', 
                  padding: 12, 
                  borderRadius: 6,
                  margin: 0,
                  fontSize: 12,
                  border: '1px solid #f0f0f0'
                }}>
                  {formatJsonData(selectedLog.response_data)}
                </pre>
              </Card>
            )}

            {/* 数据变更对比 */}
            {(hasValidData(selectedLog.before_data) || hasValidData(selectedLog.after_data)) && (
              <Card 
                size="small" 
                title={
                  <Space>
                    <SwapOutlined />
                    <span>数据变更</span>
                  </Space>
                }
              >
                <Row gutter={16}>
                  {hasValidData(selectedLog.before_data) && (
                    <Col span={hasValidData(selectedLog.after_data) ? 12 : 24}>
                      <div style={{ marginBottom: 8 }}>
                        <Tag color="orange">变更前</Tag>
                      </div>
                      <pre style={{ 
                        maxHeight: 200, 
                        overflow: 'auto', 
                        background: '#fffbe6', 
                        padding: 12, 
                        borderRadius: 6,
                        margin: 0,
                        fontSize: 11,
                        border: '1px solid #ffe58f'
                      }}>
                        {formatJsonData(selectedLog.before_data)}
                      </pre>
                    </Col>
                  )}
                  {hasValidData(selectedLog.after_data) && (
                    <Col span={hasValidData(selectedLog.before_data) ? 12 : 24}>
                      <div style={{ marginBottom: 8 }}>
                        <Tag color="green">变更后</Tag>
                      </div>
                      <pre style={{ 
                        maxHeight: 200, 
                        overflow: 'auto', 
                        background: '#f6ffed', 
                        padding: 12, 
                        borderRadius: 6,
                        margin: 0,
                        fontSize: 11,
                        border: '1px solid #b7eb8f'
                      }}>
                        {formatJsonData(selectedLog.after_data)}
                      </pre>
                    </Col>
                  )}
                </Row>
              </Card>
            )}
          </div>
        )}
      </Modal>
    </>
  );
};

export default Audit;

