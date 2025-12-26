import React, { useState, useEffect, useRef } from 'react';
import {
  Card,
  Row,
  Col,
  Button,
  Select,
  Form,
  Input,
  Space,
  message,
  Table,
  Tag,
  Modal,
  Typography,
  Tabs,
  Descriptions,
  Spin,
  Alert,
  Collapse,
} from 'antd';
import {
  PlayCircleOutlined,
  ReloadOutlined,
  HistoryOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  StopOutlined,
  BookOutlined,
  SyncOutlined,
  DeleteOutlined,
} from '@ant-design/icons';
import HostSelector from '../components/HostSelector';
import formulaService from '../services/formula';
// Formula 部署结果查看器组件
const FormulaResultViewer: React.FC<{ deployment: FormulaDeployment; isLive?: boolean }> = ({ deployment, isLive = false }) => {
  const parseResults = (results: any) => {
    if (!results) return [];

    try {
      // 如果results是字符串，尝试解析为JSON
      if (typeof results === 'string') {
        results = JSON.parse(results);
      }

      const rows: Array<{
        hostId: string;
        hostname: string;
        status: 'success' | 'failed';
        output?: string;
        error?: string;
      }> = [];

      // 解析结果 - 直接是 {minion_id: result} 格式
      if (results && typeof results === 'object') {
        // 检查是否有 return 字段（Salt API 原始格式）
        let minions = results;
        if (results.return && Array.isArray(results.return) && results.return.length > 0) {
          minions = results.return[0];
        }

        Object.entries(minions).forEach(([minionID, minionResult]: [string, any]) => {
          let status: 'success' | 'failed' = 'failed';
          let output = '';

          // 处理字符串错误信息
          if (typeof minionResult === 'string') {
            output = minionResult;
            status = 'failed';
          }
          // 处理数组格式的错误信息
          else if (Array.isArray(minionResult)) {
            output = minionResult.join('\n');
            status = 'failed';
          }
          // 处理对象格式的 state 结果
          else if (minionResult && typeof minionResult === 'object') {
            let allSuccess = true;
            const outputParts: string[] = [];

            // Salt state.apply 返回格式: {state_id: {result: bool, comment: string, changes: {}}}
            Object.entries(minionResult).forEach(([stateId, stateResult]: [string, any]) => {
              if (stateId === 'retcode') return; // 跳过 retcode 字段

              if (stateResult && typeof stateResult === 'object') {
                const stateSuccess = stateResult.result === true;
                if (!stateSuccess) {
                  allSuccess = false;
                }

                // 收集输出信息
                if (stateResult.comment) {
                  outputParts.push(`[${stateId}] ${stateResult.comment}`);
                }
                if (stateResult.changes && Object.keys(stateResult.changes).length > 0) {
                  outputParts.push(`变更: ${JSON.stringify(stateResult.changes)}`);
                }
              }
            });

            status = allSuccess ? 'success' : 'failed';
            output = outputParts.join('\n') || (allSuccess ? '执行成功' : '执行失败');
          }

          rows.push({
            hostId: minionID,
            hostname: minionID,
            status,
            output: output || undefined,
          });
        });
      }

      return rows;
    } catch (error) {
      console.error('解析部署结果失败:', error);
      return [];
    }
  };

  const results = parseResults(deployment.results);

  // 获取状态显示配置
  const getStatusDisplay = (status: string) => {
    const statusConfig: Record<string, { color: string; text: string; icon: React.ReactNode }> = {
      pending: { color: 'orange', text: '等待中', icon: <ClockCircleOutlined /> },
      running: { color: 'blue', text: '执行中', icon: <SyncOutlined spin /> },
      completed: { color: 'green', text: '成功', icon: <CheckCircleOutlined /> },
      failed: { color: 'red', text: '失败', icon: <CloseCircleOutlined /> },
      cancelled: { color: 'gray', text: '已取消', icon: <StopOutlined /> },
    };
    return statusConfig[status] || statusConfig.pending;
  };

  const statusDisplay = getStatusDisplay(deployment.status);

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Space size="large">
          <span>
            状态: <Tag color={statusDisplay.color} icon={statusDisplay.icon}>{statusDisplay.text}</Tag>
          </span>
          <span>成功: <strong style={{ color: '#52c41a' }}>{deployment.success_count}</strong></span>
          <span>失败: <strong style={{ color: '#ff4d4f' }}>{deployment.failed_count}</strong></span>
          {deployment.duration_seconds !== undefined && deployment.duration_seconds > 0 && (
            <span>耗时: <strong>{Math.round(deployment.duration_seconds)}s</strong></span>
          )}
        </Space>
      </div>

      {/* 实时日志区域 */}
      {isLive && (deployment.status === 'pending' || deployment.status === 'running') && (
        <div style={{ marginBottom: 16 }}>
          <div style={{
            background: '#1e1e1e',
            color: '#d4d4d4',
            padding: '16px',
            borderRadius: '4px',
            fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
            fontSize: '13px',
            minHeight: '150px',
            maxHeight: '300px',
            overflow: 'auto',
          }}>
            <div style={{ color: '#569cd6' }}>
              [INFO] 部署任务已提交，正在执行中...
            </div>
            <div style={{ color: '#6a9955' }}>
              [INFO] Formula: {deployment.name}
            </div>
            <div style={{ color: '#6a9955' }}>
              [INFO] 目标主机: {Array.isArray(deployment.target_hosts) ? deployment.target_hosts.join(', ') : deployment.target_hosts}
            </div>
            <div style={{ color: '#dcdcaa' }}>
              [INFO] Salt Job ID: {deployment.salt_job_id || '获取中...'}
            </div>
            <div style={{ color: '#569cd6', marginTop: '8px' }}>
              <SyncOutlined spin /> 等待执行结果...
            </div>
          </div>
        </div>
      )}

      {deployment.error_message && (
        <Alert
          message="执行错误"
          description={deployment.error_message}
          type="error"
          showIcon
          style={{ marginBottom: 16 }}
        />
      )}

      {results.length > 0 ? (
        <Table
          size="small"
          columns={[
            {
              title: '主机',
              dataIndex: 'hostname',
              key: 'hostname',
              width: 150,
            },
            {
              title: '状态',
              dataIndex: 'status',
              key: 'status',
              width: 80,
              render: (status: string) => (
                status === 'success' ?
                  <Tag color="green" icon={<CheckCircleOutlined />}>成功</Tag> :
                  <Tag color="red" icon={<CloseCircleOutlined />}>失败</Tag>
              ),
            },
            {
              title: '执行结果',
              dataIndex: 'output',
              key: 'output',
              render: (output: string) => (
                <pre style={{
                  maxWidth: '500px',
                  maxHeight: '200px',
                  overflow: 'auto',
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-all',
                  margin: 0,
                  padding: '8px',
                  background: '#f5f5f5',
                  borderRadius: '4px',
                  fontSize: '12px',
                }}>
                  {output || '-'}
                </pre>
              ),
            },
          ]}
          dataSource={results}
          rowKey="hostId"
          pagination={false}
        />
      ) : (
        deployment.status !== 'pending' && deployment.status !== 'running' && (
          <Alert message="暂无详细执行结果" type="info" />
        )
      )}
    </div>
  );
};

const { Option } = Select;
const { TextArea } = Input;
const { Title, Text } = Typography;
const { TabPane } = Tabs;

// Formula接口定义
interface Formula {
  id: number;
  name: string;
  description: string;
  category: string;
  version: string;
  path: string;
  repository: string;
  icon?: string;
  tags?: string[];
  metadata?: any;
  is_active: boolean;
}

// Formula参数接口
interface FormulaParameter {
  id: number;
  name: string;
  type: string;
  default?: any;
  required: boolean;
  label: string;
  description: string;
  validation?: any;
  order: number;
}

// Formula部署接口
interface FormulaDeployment {
  id: number;
  formula_id: number;
  name: string;
  description: string;
  target_hosts: string[];
  pillar_data?: any;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  salt_job_id?: string;
  results?: any;
  success_count: number;
  failed_count: number;
  error_message?: string;
  started_by: number;
  started_at?: string;
  completed_at?: string;
  duration_seconds?: number;
  created_at: string;
  updated_at: string;
}

const FormulaDeployment: React.FC = () => {
  const [form] = Form.useForm();
  const [selectedHosts, setSelectedHosts] = useState<Array<{ id: number; hostname: string; ip_address?: string }>>([]);
  const [formulas, setFormulas] = useState<Formula[]>([]);
  const [selectedFormula, setSelectedFormula] = useState<Formula | undefined>(undefined);
  const [formulaParameters, setFormulaParameters] = useState<FormulaParameter[]>([]);
  const [currentDeployment, setCurrentDeployment] = useState<FormulaDeployment | undefined>(undefined);
  const [deployments, setDeployments] = useState<FormulaDeployment[]>([]);
  const [totalDeployments, setTotalDeployments] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  const [loading, setLoading] = useState(false);
  const [deploymentsLoading, setDeploymentsLoading] = useState(false);
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  // 加载Formula列表
  const loadFormulas = async () => {
    try {
      const response = await formulaService.listFormulas(1, 100); // 获取所有Formula
      setFormulas(response.formulas || []);
    } catch (error: any) {
      console.error('Failed to load formulas:', error);
      message.error('加载Formula列表失败');
    }
  };

  // 加载Formula详情
  const loadFormulaDetails = async (formulaId: number) => {
    try {
      const response = await formulaService.getFormula(formulaId);
      const formula = response.formula;
      const parameters = response.parameters || [];

      setSelectedFormula(formula);
      setFormulaParameters(parameters);
    } catch (error: any) {
      console.error('Failed to load formula details:', error);
      message.error('加载Formula详情失败');
    }
  };

  // 加载部署历史
  const loadDeployments = async (page: number = currentPage) => {
    setDeploymentsLoading(true);
    try {
      const response = await formulaService.listDeployments(page, pageSize);
      setDeployments(response.deployments || []);
      setTotalDeployments(response.total || 0);
      setCurrentPage(page);
    } catch (error: any) {
      console.error('Failed to load deployments:', error);
      message.error('加载部署历史失败');
    } finally {
      setDeploymentsLoading(false);
    }
  };

  useEffect(() => {
    loadFormulas();
    loadDeployments();
  }, []);

  useEffect(() => {
    // 清理轮询
    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, []);

  // 轮询当前部署状态
  useEffect(() => {
    if (currentDeployment && (currentDeployment.status === 'pending' || currentDeployment.status === 'running')) {
      pollingRef.current = setInterval(async () => {
        try {
          const deployment = await formulaService.getDeployment(currentDeployment.id);
          setCurrentDeployment(deployment);

          if (deployment.status === 'completed' || deployment.status === 'failed' || deployment.status === 'cancelled') {
            if (pollingRef.current) {
              clearInterval(pollingRef.current);
            }
            loadDeployments(); // 刷新列表
          }
        } catch (error) {
          console.error('获取部署状态失败:', error);
        }
      }, 2000);
    }

    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, [currentDeployment]);

  // Formula选择变化处理
  const handleFormulaChange = (formulaId: number) => {
    const formula = formulas.find(f => f.id === formulaId);
    if (formula) {
      loadFormulaDetails(formulaId);
      // 清空之前的Pillar数据
      form.setFieldsValue({
        pillar_data: {},
      });
    }
  };

  // 执行部署
  const handleDeploy = async () => {
    if (!selectedFormula) {
      message.warning('请先选择Formula');
      return;
    }

    if (selectedHosts.length === 0) {
      message.warning('请至少选择一台主机');
      return;
    }

    try {
      const values = form.getFieldsValue();
      const deploymentName = `${selectedFormula.name} - ${new Date().toLocaleString()}`;

      const requestData = {
        formula_id: selectedFormula.id,
        name: deploymentName,
        description: '',
        target_hosts: selectedHosts.map(h => h.hostname), // 使用hostname作为Salt target
        pillar_data: values.pillar_data || {},
      };

      const deployment = await formulaService.createDeployment(requestData);

      // 立即执行部署
      await executeDeployment(deployment.id);

      message.success('部署任务已创建并开始执行');
      loadDeployments(); // 刷新列表

    } catch (error: any) {
      message.error(error.message || '部署失败');
    }
  };

  // 执行部署
  const executeDeployment = async (deploymentId: number) => {
    try {
      await formulaService.executeDeployment(deploymentId);

      // 加载部署详情
      const deployment = await formulaService.getDeployment(deploymentId);
      setCurrentDeployment(deployment);

    } catch (error: any) {
      message.error(error.message || '执行部署失败');
    }
  };

  // 取消部署
  const handleCancelDeployment = async (deploymentId: number) => {
    try {
      await formulaService.cancelDeployment(deploymentId);

      message.success('部署已取消');
      loadDeployments();

      if (currentDeployment && currentDeployment.id === deploymentId) {
        setCurrentDeployment(undefined);
      }

    } catch (error: any) {
      message.error(error.message || '取消部署失败');
    }
  };

  // 渲染Formula参数表单
  const renderParameterFields = () => {
    if (!formulaParameters || formulaParameters.length === 0) {
      return (
        <Alert
          message="此Formula无需配置参数"
          description="可以直接使用默认配置进行部署"
          type="info"
          showIcon
        />
      );
    }

    return formulaParameters.map(param => {
      const fieldName = `pillar_data.${param.name}`;

      switch (param.type) {
        case 'string':
          return (
            <Form.Item
              key={param.name}
              label={param.label || param.name}
              name={fieldName}
              rules={param.required ? [{ required: true, message: `${param.label || param.name}是必填项` }] : []}
              tooltip={param.description}
              initialValue={param.default}
            >
              <Input placeholder={`请输入${param.label || param.name}`} />
            </Form.Item>
          );

        case 'number':
          return (
            <Form.Item
              key={param.name}
              label={param.label || param.name}
              name={fieldName}
              rules={param.required ? [{ required: true, message: `${param.label || param.name}是必填项` }] : []}
              tooltip={param.description}
              initialValue={param.default}
            >
              <Input type="number" placeholder={`请输入${param.label || param.name}`} />
            </Form.Item>
          );

        case 'boolean':
          return (
            <Form.Item
              key={param.name}
              label={param.label || param.name}
              name={fieldName}
              valuePropName="checked"
              tooltip={param.description}
              initialValue={param.default || false}
            >
              <input type="checkbox" />
            </Form.Item>
          );

        case 'array':
        case 'object':
          return (
            <Form.Item
              key={param.name}
              label={param.label || param.name}
              name={fieldName}
              rules={param.required ? [{ required: true, message: `${param.label || param.name}是必填项` }] : []}
              tooltip={param.description}
              initialValue={param.default ? JSON.stringify(param.default, null, 2) : ''}
            >
              <TextArea
                rows={4}
                placeholder={`请输入JSON格式的${param.label || param.name}`}
              />
            </Form.Item>
          );

        default:
          return (
            <Form.Item
              key={param.name}
              label={param.label || param.name}
              name={fieldName}
              rules={param.required ? [{ required: true, message: `${param.label || param.name}是必填项` }] : []}
              tooltip={param.description}
              initialValue={param.default}
            >
              <Input placeholder={`请输入${param.label || param.name}`} />
            </Form.Item>
          );
      }
    });
  };

  return (
    <div style={{ padding: '20px' }}>
      <Title level={2}>
        <BookOutlined style={{ marginRight: '8px' }} />
        Formula部署管理
      </Title>

      {/* Formula仓库同步区域 */}
      <Card
        size="small"
        style={{ marginBottom: '20px', background: '#f6ffed', border: '1px solid #b7eb8f' }}
        bodyStyle={{ padding: '12px' }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <span style={{ fontWeight: 'bold', color: '#52c41a' }}>Formula仓库同步</span>
            <span style={{ marginLeft: '12px', fontSize: '14px', color: '#666' }}>
              同步最新的Formula定义和配置
            </span>
          </div>
          <Space>
            <Button
              type="primary"
              icon={<ReloadOutlined />}
              onClick={() => {
                loadFormulas();
                message.success('Formula列表已刷新');
              }}
            >
              刷新Formula列表
            </Button>
            <Button
              icon={<SyncOutlined />}
              onClick={async () => {
                try {
                  // 获取所有仓库并逐个同步
                  const response = await formulaService.listRepositories(1, 100);
                  const repos = response.repositories || [];

                  if (repos.length === 0) {
                    message.warning('没有找到Formula仓库，请先在系统设置中配置仓库');
                    return;
                  }

                  let syncedCount = 0;
                  for (const repo of repos) {
                    if (repo.is_active) {
                      try {
                        await formulaService.syncRepository(repo.id);
                        syncedCount++;
                      } catch (error) {
                        console.error(`同步仓库 ${repo.name} 失败:`, error);
                      }
                    }
                  }

                  if (syncedCount > 0) {
                    message.success(`已同步 ${syncedCount} 个仓库`);
                    loadFormulas(); // 同步完成后刷新Formula列表
                  } else {
                    message.info('没有活跃的仓库需要同步');
                  }
                } catch (error: any) {
                  message.error('同步仓库失败: ' + (error.message || '未知错误'));
                }
              }}
            >
              同步所有仓库
            </Button>
          </Space>
        </div>
      </Card>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="部署配置" size="small">
            <Form form={form} layout="vertical">
              <Form.Item
                label="Formula选择"
                name="formula_id"
                rules={[{ required: true, message: '请选择Formula' }]}
                style={{ marginBottom: 12 }}
              >
                <Select
                  placeholder="选择要部署的Formula"
                  onChange={handleFormulaChange}
                  showSearch
                  optionFilterProp="children"
                >
                  {formulas.map(formula => (
                    <Option key={formula.id} value={formula.id}>
                      <div>
                        <div style={{ fontWeight: 'bold' }}>{formula.name}</div>
                        <div style={{ fontSize: '12px', color: '#666' }}>
                          {formula.category} • {formula.description}
                        </div>
                      </div>
                    </Option>
                  ))}
                </Select>
              </Form.Item>

              {selectedFormula && (
                <div style={{ marginBottom: 12 }}>
                  <Descriptions size="small" column={2} bordered>
                    <Descriptions.Item label="名称">{selectedFormula.name}</Descriptions.Item>
                    <Descriptions.Item label="分类">{selectedFormula.category}</Descriptions.Item>
                    <Descriptions.Item label="版本">{selectedFormula.version || '最新'}</Descriptions.Item>
                    <Descriptions.Item label="路径">{selectedFormula.path}</Descriptions.Item>
                    <Descriptions.Item label="描述" span={2}>{selectedFormula.description}</Descriptions.Item>
                  </Descriptions>
                </div>
              )}

              {selectedFormula && (
                <div>
                  <Title level={5} style={{ marginBottom: 12 }}>配置参数</Title>
                  {renderParameterFields()}
                </div>
              )}
            </Form>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card 
            title="目标主机" 
            size="small"
            style={{ marginBottom: 16 }}
          >
            <HostSelector
              value={selectedHosts}
              onChange={setSelectedHosts}
            />
          </Card>
          <Button
            type="primary"
            icon={<PlayCircleOutlined />}
            onClick={handleDeploy}
            loading={loading}
            size="large"
            block
          >
            部署
          </Button>
        </Col>

        {currentDeployment && (
          <Col span={24}>
            <Card
              title={
                <Space>
                  {(currentDeployment.status === 'pending' || currentDeployment.status === 'running') && (
                    <SyncOutlined spin style={{ color: '#1890ff' }} />
                  )}
                  <span>当前部署: {currentDeployment.name}</span>
                </Space>
              }
              extra={
                <Space>
                  {currentDeployment.status === 'running' && (
                    <Button
                      danger
                      onClick={() => handleCancelDeployment(currentDeployment.id)}
                    >
                      取消部署
                    </Button>
                  )}
                  {(currentDeployment.status === 'completed' || currentDeployment.status === 'failed') && (
                    <Button onClick={() => setCurrentDeployment(undefined)}>
                      关闭
                    </Button>
                  )}
                </Space>
              }
            >
              <FormulaResultViewer deployment={currentDeployment} isLive={true} />
            </Card>
          </Col>
        )}

        <Col span={24}>
          <Collapse
            defaultActiveKey={[]}
            style={{ marginTop: 0 }}
            items={[
              {
                key: 'history',
                label: (
                  <Space>
                    <HistoryOutlined />
                    <span>部署历史</span>
                    <Tag>{totalDeployments} 条记录</Tag>
                  </Space>
                ),
                extra: (
                  <Space onClick={(e) => e.stopPropagation()}>
                    <Button
                      size="small"
                      icon={<ReloadOutlined />}
                      onClick={() => loadDeployments(currentPage)}
                      loading={deploymentsLoading}
                    >
                      刷新
                    </Button>
                    <Button
                      size="small"
                      danger
                      icon={<DeleteOutlined />}
                      onClick={async () => {
                        try {
                          const result = await formulaService.cleanupOldDeployments();
                          message.success(`成功清理了 ${result.count} 条1个月前的部署记录`);
                          loadDeployments(1);
                        } catch (error: any) {
                          message.error('清理数据失败: ' + (error.response?.data?.error || error.message));
                        }
                      }}
                      loading={deploymentsLoading}
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
                        title: 'Formula',
                        key: 'formula',
                        render: (_, record: FormulaDeployment) => (
                          <Typography.Text code style={{ fontSize: '12px' }}>
                            {formulas.find(f => f.id === record.formula_id)?.name || `ID:${record.formula_id}`}
                          </Typography.Text>
                        ),
                      },
                      {
                        title: '主机',
                        key: 'hosts',
                        width: 60,
                        render: (_, record: FormulaDeployment) => (
                          Array.isArray(record.target_hosts) ? record.target_hosts.length : 0
                        ),
                      },
                      {
                        title: '状态',
                        dataIndex: 'status',
                        key: 'status',
                        width: 80,
                        render: (status: string) => {
                          const statusConfig: Record<string, { color: string; text: string }> = {
                            pending: { color: 'default', text: '等待中' },
                            running: { color: 'processing', text: '运行中' },
                            completed: { color: 'success', text: '已完成' },
                            failed: { color: 'error', text: '失败' },
                            cancelled: { color: 'warning', text: '已取消' },
                          };
                          const config = statusConfig[status] || statusConfig.pending;
                          return <Tag color={config.color}>{config.text}</Tag>;
                        },
                      },
                      {
                        title: '结果',
                        key: 'results',
                        width: 100,
                        render: (_, record: FormulaDeployment) => (
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
                          <span style={{ fontSize: '12px' }}>
                            {text ? new Date(text).toLocaleString('zh-CN') : '-'}
                          </span>
                        ),
                      },
                      {
                        title: '操作',
                        key: 'action',
                        width: 100,
                        render: (_, record: FormulaDeployment) => (
                          <Space size="small">
                            <Button
                              size="small"
                              type="link"
                              style={{ padding: 0 }}
                              onClick={() => {
                                setCurrentDeployment(record);
                                window.scrollTo({ top: 0, behavior: 'smooth' });
                              }}
                            >
                              查看
                            </Button>
                            {record.status === 'running' && (
                              <Button
                                size="small"
                                type="link"
                                danger
                                style={{ padding: 0 }}
                                onClick={() => handleCancelDeployment(record.id)}
                              >
                                取消
                              </Button>
                            )}
                          </Space>
                        ),
                      },
                    ]}
                    dataSource={deployments}
                    rowKey="id"
                    loading={deploymentsLoading}
                    pagination={{
                      current: currentPage,
                      pageSize: pageSize,
                      total: totalDeployments,
                      showSizeChanger: false,
                      showQuickJumper: true,
                      size: 'small',
                      showTotal: (total, range) => `${range[0]}-${range[1]} / ${total}`,
                      onChange: (page) => loadDeployments(page),
                    }}
                    size="small"
                  />
                ),
              },
            ]}
          />
        </Col>
      </Row>
    </div>
  );
};

export default FormulaDeployment;
