import React, { useState, useEffect, useRef } from 'react';
import {
  Card,
  Row,
  Col,
  Button,
  Input,
  Select,
  Form,
  Space,
  message,
  Table,
  Tag,
  Modal,
  Typography,
  Radio,
  InputNumber,
  Checkbox,
  Collapse,
} from 'antd';
import {
  PlayCircleOutlined,
  SaveOutlined,
  ReloadOutlined,
  HistoryOutlined,
  ThunderboltOutlined,
  StopOutlined,
  DeleteOutlined,
  DownOutlined,
  RightOutlined,
} from '@ant-design/icons';
import { batchOperationsService, BatchOperation, CreateBatchOperationRequest } from '../services/batchOperations';
import { commandTemplateService, CommandTemplate } from '../services/commandTemplate';
import HostSelector from '../components/HostSelector';
import ResultViewer from '../components/ResultViewer';

const { Option } = Select;
const { TextArea } = Input;
const { Title, Text } = Typography;

// 内置命令模板
const builtinTemplates: Array<{
  name: string;
  function: string;
  args?: any[];
  category: string;
  description: string;
}> = [
  { name: '查看系统负载', function: 'status.loadavg', category: 'system', description: '查看系统负载平均值' },
  { name: '查看内存使用', function: 'status.meminfo', category: 'system', description: '查看内存使用情况' },
  { name: '查看磁盘使用', function: 'disk.usage', category: 'disk', description: '查看磁盘使用情况' },
  { name: '查看CPU信息', function: 'status.cpuinfo', category: 'system', description: '查看CPU信息' },
  { name: '查看网络接口', function: 'network.interfaces', category: 'network', description: '查看网络接口信息' },
  { name: '查看活动TCP连接', function: 'network.active_tcp', category: 'network', description: '查看活动的TCP连接' },
  { name: '查看进程列表', function: 'cmd.run', args: ['ps aux | head -20'], category: 'process', description: '查看进程列表' },
  { name: '测试连通性', function: 'test.ping', category: 'system', description: '测试主机连通性' },
  // 包管理
  { name: '批量安装软件包', function: 'pkg.install', args: ['vim', 'htop', 'curl'], category: 'package', description: '批量安装多个软件包' },
  { name: '更新系统包', function: 'pkg.update', category: 'package', description: '更新系统软件包列表' },
  { name: '升级系统包', function: 'pkg.upgrade', category: 'package', description: '升级所有已安装的软件包' },
];

const BatchOperations: React.FC = () => {
  const [form] = Form.useForm();
  const [selectedHosts, setSelectedHosts] = useState<Array<{ id: number; hostname: string; ip_address?: string }>>([]);
  const [commandType, setCommandType] = useState<'custom' | 'template' | 'builtin'>('custom');
  const [selectedTemplate, setSelectedTemplate] = useState<number | undefined>(undefined);
  const [selectedBuiltin, setSelectedBuiltin] = useState<string | undefined>(undefined);
  const [currentOperation, setCurrentOperation] = useState<BatchOperation | undefined>(undefined);
  const [operations, setOperations] = useState<BatchOperation[]>([]);
  const [totalOperations, setTotalOperations] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20); // 每页显示20条记录
  const [templates, setTemplates] = useState<CommandTemplate[]>([]);
  const [loading, setLoading] = useState(false);
  const [operationsLoading, setOperationsLoading] = useState(false);
  const [saveTemplateModalVisible, setSaveTemplateModalVisible] = useState(false);
  const [templateForm] = Form.useForm();
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    loadTemplates();
    loadOperations();
  }, []);

  useEffect(() => {
    // 清理轮询
    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, []);

  // 轮询当前操作状态
  useEffect(() => {
    if (currentOperation && (currentOperation.status === 'pending' || currentOperation.status === 'running')) {
      pollingRef.current = setInterval(async () => {
        try {
          const { operation } = await batchOperationsService.getOperation(currentOperation.id);
          setCurrentOperation(operation);
          
          if (operation.status === 'completed' || operation.status === 'failed' || operation.status === 'cancelled') {
            if (pollingRef.current) {
              clearInterval(pollingRef.current);
              pollingRef.current = null;
            }
            loadOperations(); // 刷新操作列表
          }
        } catch (error: any) {
          console.error('Failed to poll operation status:', error);
        }
      }, 2000); // 每2秒轮询一次
    } else {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
        pollingRef.current = null;
      }
    }

    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
        pollingRef.current = null;
      }
    };
  }, [currentOperation]);

  const loadTemplates = async () => {
    try {
      const response = await commandTemplateService.listTemplates();
      setTemplates(response.templates || []);
    } catch (error: any) {
      console.error('Failed to load templates:', error);
    }
  };

  const loadOperations = async (page: number = 1) => {
    setOperationsLoading(true);
    try {
      const response = await batchOperationsService.listOperations(page, pageSize);
      setOperations(response.operations || []);
      setTotalOperations(response.total || 0);
      setCurrentPage(page);
    } catch (error: any) {
      message.error('加载操作历史失败');
    } finally {
      setOperationsLoading(false);
    }
  };

  const handleCleanupOldOperations = async () => {
    try {
      const result = await batchOperationsService.cleanupOldOperations();
      message.success(`成功清理了 ${result.count} 条1个月前的操作记录`);
      // 重新加载当前页数据
      loadOperations(currentPage);
    } catch (error: any) {
      message.error('清理数据失败: ' + (error.response?.data?.error || error.message));
    }
  };


  const handleExecute = async () => {
    if (selectedHosts.length === 0) {
      message.warning('请至少选择一个主机');
      return;
    }

    const values = form.getFieldsValue();
    
    let commandFunction = '';
    let commandArgs: any[] = [];

    if (commandType === 'custom') {
      if (!values.command_function) {
        message.warning('请输入命令函数');
        return;
      }
      commandFunction = values.command_function;
      if (values.command_args) {
        try {
          const argsValue = values.command_args;
          if (typeof argsValue === 'string') {
            // 尝试解析为JSON
            try {
              commandArgs = JSON.parse(argsValue);
            } catch (jsonError) {
              // 如果不是JSON，检查是否是多行命令
              const trimmedValue = argsValue.trim();
              if (trimmedValue.includes('\n')) {
                // 多行命令，合并为单个字符串
                const lines = trimmedValue.split('\n')
                  .map(line => line.trim())
                  .filter(line => line.length > 0);
                if (lines.length > 0) {
                  commandArgs = [lines.join(' && ')];
                }
              } else if (trimmedValue) {
                // 单行命令
                commandArgs = [trimmedValue];
              }
            }
          } else {
            commandArgs = argsValue;
          }
        } catch (e) {
          message.error('命令参数格式错误');
          return;
        }
      }
    } else if (commandType === 'template') {
      if (!selectedTemplate) {
        message.warning('请选择一个命令模板');
        return;
      }
      const template = templates.find(t => t.id === selectedTemplate);
      if (!template) {
        message.error('模板不存在');
        return;
      }
      commandFunction = template.command_function;
      if (template.command_args) {
        commandArgs = Array.isArray(template.command_args) 
          ? template.command_args 
          : JSON.parse(template.command_args as any);
      }
      // 增加模板使用次数
      // commandTemplateService 没有 incrementUsageCount 方法，暂时跳过
    } else if (commandType === 'builtin') {
      if (!selectedBuiltin) {
        message.warning('请选择一个内置命令');
        return;
      }
      const builtin = builtinTemplates.find(t => t.function === selectedBuiltin);
      if (!builtin) {
        message.error('内置命令不存在');
        return;
      }
      commandFunction = builtin.function;
      commandArgs = builtin.args || [];
    }

    const operationName = `${commandFunction} - ${new Date().toLocaleString()}`;

    const request: CreateBatchOperationRequest = {
      name: operationName,
      description: '',
      command_type: commandType,
      command_function: commandFunction,
      command_args: commandArgs.length > 0 ? commandArgs : undefined,
      target_hosts: selectedHosts,
    };

    setLoading(true);
    try {
      const { operation } = await batchOperationsService.createOperation(request);
      setCurrentOperation(operation);
      message.success('操作已创建，正在执行...');
      loadOperations();
    } catch (error: any) {
      message.error(error.response?.data?.error || '执行失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = async () => {
    if (!currentOperation) return;
    
    if (currentOperation.status !== 'running') {
      message.warning('操作未在运行中，无法取消');
      return;
    }

    try {
      await batchOperationsService.cancelOperation(currentOperation.id);
      message.success('操作已取消');
      loadOperations();
      // 重新加载当前操作
      const { operation } = await batchOperationsService.getOperation(currentOperation.id);
      setCurrentOperation(operation);
    } catch (error: any) {
      message.error('取消操作失败');
    }
  };

  const handleSaveAsTemplate = () => {
    const values = form.getFieldsValue();
    if (!values.command_function && commandType === 'custom') {
      message.warning('请先配置命令');
      return;
    }

    // 预填充模板表单，使用智能默认值
    let templateName = '';
    let templateDescription = '';

    if (commandType === 'custom') {
      templateName = `${values.command_function} 模板`;
      templateDescription = `执行 ${values.command_function} 命令`;
    } else if (commandType === 'template' && selectedTemplate) {
      // 如果是从现有模板创建，添加"副本"后缀
      const template = templates.find(t => t.id === selectedTemplate);
      templateName = template ? `${template.name} 副本` : '模板副本';
      templateDescription = template ? (template.description || '模板副本') : '模板副本';
    } else if (commandType === 'builtin' && selectedBuiltin) {
      templateName = `${selectedBuiltin} 模板`;
      templateDescription = `执行内置命令：${selectedBuiltin}`;
    }

    templateForm.setFieldsValue({
      name: templateName,
      description: templateDescription,
      category: 'custom', // 默认分类
    });

    setSaveTemplateModalVisible(true);
  };

  const handleSaveTemplate = async () => {
    const values = form.getFieldsValue();
    const templateValues = templateForm.getFieldsValue();

    try {
      let commandArgs: any[] = [];
      if (values.command_args) {
        try {
          const argsValue = values.command_args;
          if (typeof argsValue === 'string') {
            // 尝试解析为JSON
            try {
              commandArgs = JSON.parse(argsValue);
            } catch (jsonError) {
              // 如果不是JSON，检查是否是多行命令
              const trimmedValue = argsValue.trim();
              if (trimmedValue.includes('\n')) {
                // 多行命令，合并为单个字符串
                const lines = trimmedValue.split('\n')
                  .map(line => line.trim())
                  .filter(line => line.length > 0);
                if (lines.length > 0) {
                  commandArgs = [lines.join(' && ')];
                }
              } else if (trimmedValue) {
                // 单行命令
                commandArgs = [trimmedValue];
              }
            }
          } else {
            commandArgs = argsValue;
          }
        } catch (e) {
          message.error('命令参数格式错误');
          return;
        }
      }

      await commandTemplateService.createTemplate({
        name: templateValues.name,
        description: templateValues.description,
        category: templateValues.category,
        command_function: values.command_function,
        command_args: commandArgs.length > 0 ? commandArgs : undefined,
        icon: templateValues.icon,
        is_public: templateValues.is_public || false,
      });

      message.success('模板保存成功');
      setSaveTemplateModalVisible(false);
      templateForm.resetFields();
      loadTemplates();
    } catch (error: any) {
      message.error(error.response?.data?.error || '保存模板失败');
    }
  };

  const handleRetry = async (operation: BatchOperation) => {
    try {
      const { operation: newOperation } = await batchOperationsService.retryOperation(operation.id);
      setCurrentOperation(newOperation);
      message.success('操作已重试');
      loadOperations();
    } catch (error: any) {
      message.error('重试失败');
    }
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending: { color: 'default', text: '等待中' },
      running: { color: 'processing', text: '运行中' },
      completed: { color: 'success', text: '已完成' },
      failed: { color: 'error', text: '失败' },
      cancelled: { color: 'warning', text: '已取消' },
    };
    const statusInfo = statusMap[status] || { color: 'default', text: status };
    return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
  };

  return (
    <div>
      <Title level={2}>批量操作</Title>

      <Row gutter={16} style={{ marginTop: 16 }}>
        {/* 左侧：主机选择 */}
        <Col xs={24} lg={8}>
          <HostSelector value={selectedHosts} onChange={setSelectedHosts} />
        </Col>

        {/* 中间：命令配置 */}
        <Col xs={24} lg={8}>
          <Card
            title={
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                <span>命令配置</span>
                <Button size="small" icon={<SaveOutlined />} onClick={handleSaveAsTemplate}>
                  保存为模板
                </Button>
              </div>
            }
            style={{ height: '100%' }}
          >
            <Form form={form} layout="vertical">
              <div style={{ marginBottom: 16 }}>
                <div style={{ fontSize: '14px', fontWeight: 500, marginBottom: 8 }}>命令类型</div>
                <Radio.Group
                  value={commandType}
                  onChange={(e) => {
                    setCommandType(e.target.value);
                    setSelectedTemplate(undefined);
                    setSelectedBuiltin(undefined);
                  }}
                  style={{ marginBottom: 0 }}
                >
                  <Radio value="custom">自定义</Radio>
                  <Radio value="template">模板</Radio>
                  <Radio value="builtin">内置</Radio>
                </Radio.Group>
              </div>

              {commandType === 'custom' && (
                <>
                  <Form.Item label="命令函数" name="command_function" rules={[{ required: true, message: '请输入命令函数' }]}>
                    <Input
                      placeholder="例如: test.ping, status.loadavg"
                      onBlur={(e) => {
                        const value = e.target.value.trim();
                        if (value !== e.target.value) {
                          form.setFieldsValue({ command_function: value });
                        }
                      }}
                    />
                  </Form.Item>
                  <Form.Item label="命令参数" name="command_args">
                    <TextArea
                      rows={4}
                      placeholder={`多行Shell命令 (一行一个命令):
echo "开始执行..."
apt update && apt upgrade -y
echo "执行完成"

或者 JSON 格式:
["echo '开始'", "apt update", "echo '完成'"]

或者单行命令:
"uptime"`}
                      onBlur={(e) => {
                        const value = e.target.value.trim();
                        if (value !== e.target.value) {
                          form.setFieldsValue({ command_args: value });
                        }
                      }}
                    />
                  </Form.Item>
                </>
              )}

              {commandType === 'template' && (
                <>
                  <Form.Item label="选择模板">
                    <Select
                      placeholder="选择命令模板"
                      value={selectedTemplate}
                      onChange={setSelectedTemplate}
                      showSearch
                      optionFilterProp="children"
                    >
                      {templates.map(template => (
                        <Option key={template.id} value={template.id}>
                          {template.name} - {template.command_function}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                  {selectedTemplate && (() => {
                    const template = templates.find(t => t.id === selectedTemplate);
                    return template ? (
                      <div style={{ padding: '8px', backgroundColor: '#f5f5f5', borderRadius: 4 }}>
                        <Text type="secondary">{template.description || '无描述'}</Text>
                      </div>
                    ) : null;
                  })()}
                </>
              )}

              {commandType === 'builtin' && (
                <Form.Item label="选择内置命令">
                  <Select
                    placeholder="选择内置命令"
                    value={selectedBuiltin}
                    onChange={setSelectedBuiltin}
                  >
                    {builtinTemplates.map(template => (
                      <Option key={template.function} value={template.function}>
                        {template.name} - {template.description}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              )}


              <Form.Item>
                <Space>
                  <Button
                    type="primary"
                    icon={<PlayCircleOutlined />}
                    onClick={handleExecute}
                    loading={loading}
                    disabled={selectedHosts.length === 0}
                  >
                    执行
                  </Button>
                  {currentOperation && currentOperation.status === 'running' && (
                    <Button
                      danger
                      icon={<StopOutlined />}
                      onClick={handleCancel}
                    >
                      取消
                    </Button>
                  )}
                </Space>
              </Form.Item>
            </Form>
          </Card>
        </Col>

        {/* 右侧：结果查看 */}
        <Col xs={24} lg={8}>
          <ResultViewer
            operation={currentOperation}
            onRefresh={async () => {
              if (currentOperation) {
                const { operation } = await batchOperationsService.getOperation(currentOperation.id);
                setCurrentOperation(operation);
              }
            }}
          />
        </Col>
      </Row>

      {/* 操作历史 - 可折叠 */}
      <Collapse
        defaultActiveKey={[]}
        style={{ marginTop: 16 }}
        items={[
          {
            key: 'history',
            label: (
              <Space>
                <HistoryOutlined />
                <span>操作历史</span>
                <Tag>{totalOperations} 条记录</Tag>
              </Space>
            ),
            extra: (
              <Space onClick={(e) => e.stopPropagation()}>
                <Button size="small" icon={<ReloadOutlined />} onClick={() => loadOperations(currentPage)} loading={operationsLoading}>
                  刷新
                </Button>
                <Button
                  size="small"
                  danger
                  icon={<DeleteOutlined />}
                  onClick={handleCleanupOldOperations}
                  loading={operationsLoading}
                >
                  清理旧数据
                </Button>
              </Space>
            ),
            children: (
              <Table
                columns={[
                  {
                    title: 'ID',
                    dataIndex: 'id',
                    key: 'id',
                    width: 60,
                  },
                  {
                    title: '命令',
                    key: 'command',
                    render: (_, record) => (
                      <Text code style={{ fontSize: '12px' }}>{record.command_function}</Text>
                    ),
                  },
                  {
                    title: '主机',
                    dataIndex: 'target_count',
                    key: 'target_count',
                    width: 60,
                  },
                  {
                    title: '状态',
                    dataIndex: 'status',
                    key: 'status',
                    width: 80,
                    render: (status: string) => getStatusTag(status),
                  },
                  {
                    title: '结果',
                    key: 'results',
                    width: 100,
                    render: (_, record) => (
                      <Space size="small">
                        <Tag color="success" style={{ margin: 0 }}>{record.success_count}</Tag>
                        <Tag color="error" style={{ margin: 0 }}>{record.failed_count}</Tag>
                      </Space>
                    ),
                  },
                  {
                    title: '时间',
                    dataIndex: 'started_at',
                    key: 'started_at',
                    width: 150,
                    render: (text: string) => (
                      <span style={{ fontSize: '12px' }}>{new Date(text).toLocaleString('zh-CN')}</span>
                    ),
                  },
                  {
                    title: '操作',
                    key: 'action',
                    width: 100,
                    render: (_, record) => (
                      <Space size="small">
                        <Button
                          size="small"
                          type="link"
                          style={{ padding: 0 }}
                          onClick={() => {
                            setCurrentOperation(record);
                            window.scrollTo({ top: 0, behavior: 'smooth' });
                          }}
                        >
                          查看
                        </Button>
                        {(record.status === 'failed' || record.status === 'cancelled') && (
                          <Button
                            size="small"
                            type="link"
                            style={{ padding: 0 }}
                            onClick={() => handleRetry(record)}
                          >
                            重试
                          </Button>
                        )}
                      </Space>
                    ),
                  },
                ]}
                dataSource={operations}
                rowKey="id"
                loading={operationsLoading}
                pagination={{
                  current: currentPage,
                  pageSize: pageSize,
                  total: totalOperations,
                  showSizeChanger: false,
                  showQuickJumper: true,
                  size: 'small',
                  showTotal: (total, range) => `${range[0]}-${range[1]} / ${total}`,
                  onChange: (page) => loadOperations(page),
                }}
                size="small"
              />
            ),
          },
        ]}
      />

      {/* 保存模板对话框 */}
      <Modal
        title="保存为模板"
        open={saveTemplateModalVisible}
        onOk={handleSaveTemplate}
        onCancel={() => {
          setSaveTemplateModalVisible(false);
          templateForm.resetFields();
        }}
      >
        <Form form={templateForm} layout="vertical">
          <Form.Item label="模板名称" name="name" rules={[{ required: true, message: '请输入模板名称' }]}>
            <Input placeholder="输入模板名称" />
          </Form.Item>
          <Form.Item label="描述" name="description">
            <TextArea rows={2} placeholder="输入模板描述" />
          </Form.Item>
          <Form.Item label="分类" name="category">
            <Select placeholder="选择分类">
              <Option value="system">系统信息</Option>
              <Option value="network">网络</Option>
              <Option value="disk">磁盘</Option>
              <Option value="process">进程</Option>
              <Option value="custom">自定义</Option>
            </Select>
          </Form.Item>
          <Form.Item label="图标" name="icon">
            <Input placeholder="图标名称（可选）" />
          </Form.Item>
          <Form.Item name="is_public" valuePropName="checked">
            <Checkbox>公开（所有用户可见）</Checkbox>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BatchOperations;

