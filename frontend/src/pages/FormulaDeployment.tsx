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
} from 'antd';
import {
  PlayCircleOutlined,
  SaveOutlined,
  ReloadOutlined,
  HistoryOutlined,
  ThunderboltOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  StopOutlined,
  BookOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import HostSelector from '../components/HostSelector';
import formulaService from '../services/formula';
// Formula 部署结果查看器组件
const FormulaResultViewer: React.FC<{ deployment: FormulaDeployment }> = ({ deployment }) => {
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

      // 解析Salt API的返回结果
      if (results && typeof results === 'object') {
        const returnData = results.return;
        if (Array.isArray(returnData) && returnData.length > 0) {
          const minions = returnData[0];
          if (typeof minions === 'object') {
            Object.entries(minions).forEach(([minionID, minionResult]: [string, any]) => {
              let status: 'success' | 'failed' = 'failed';
              let output = '';
              let error = '';

              if (minionResult && typeof minionResult === 'object') {
                const resultList = minionResult.result;
                if (Array.isArray(resultList) && resultList.length > 0) {
                  const firstResult = resultList[0];
                  if (firstResult && typeof firstResult === 'object') {
                    const resultStatus = firstResult.result;
                    if (resultStatus === true) {
                      status = 'success';
                    }

                    // 获取输出信息
                    if (firstResult.comment) {
                      output = firstResult.comment;
                    }
                    if (firstResult.changes && typeof firstResult.changes === 'object') {
                      const changesStr = JSON.stringify(firstResult.changes, null, 2);
                      if (output) {
                        output += '\n\n变更详情:\n' + changesStr;
                      } else {
                        output = '变更详情:\n' + changesStr;
                      }
                    }
                  }
                }
              }

              rows.push({
                hostId: minionID,
                hostname: minionID,
                status,
                output: output || undefined,
                error: error || undefined,
              });
            });
          }
        }
      }

      return rows;
    } catch (error) {
      console.error('解析部署结果失败:', error);
      return [];
    }
  };

  const results = parseResults(deployment.results);

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Space>
          <span>状态: <strong>{deployment.status}</strong></span>
          <span>成功: <strong style={{ color: '#52c41a' }}>{deployment.success_count}</strong></span>
          <span>失败: <strong style={{ color: '#ff4d4f' }}>{deployment.failed_count}</strong></span>
          {deployment.duration_seconds && (
            <span>耗时: <strong>{Math.round(deployment.duration_seconds)}s</strong></span>
          )}
        </Space>
      </div>

      {deployment.error_message && (
        <div style={{ marginBottom: 16, color: '#ff4d4f' }}>
          错误信息: {deployment.error_message}
        </div>
      )}

      <Table
        size="small"
        columns={[
          {
            title: '主机',
            dataIndex: 'hostname',
            key: 'hostname',
          },
          {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => (
              status === 'success' ?
                <span style={{ color: '#52c41a' }}>成功</span> :
                <span style={{ color: '#ff4d4f' }}>失败</span>
            ),
          },
          {
            title: '结果',
            dataIndex: 'output',
            key: 'output',
            render: (output: string) => (
              <div style={{
                maxWidth: '400px',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap'
              }}>
                {output || '-'}
              </div>
            ),
          },
        ]}
        dataSource={results}
        rowKey="hostId"
        pagination={false}
      />
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

// Formula模板接口
interface FormulaTemplate {
  id: number;
  formula_id: number;
  name: string;
  description: string;
  pillar_data?: any;
  is_public: boolean;
  created_by: number;
  created_at: string;
  updated_at: string;
}

const FormulaDeployment: React.FC = () => {
  const [form] = Form.useForm();
  const [selectedHosts, setSelectedHosts] = useState<Array<{ id: number; hostname: string; ip_address?: string }>>([]);
  const [formulas, setFormulas] = useState<Formula[]>([]);
  const [selectedFormula, setSelectedFormula] = useState<Formula | undefined>(undefined);
  const [formulaParameters, setFormulaParameters] = useState<FormulaParameter[]>([]);
  const [formulaTemplates, setFormulaTemplates] = useState<FormulaTemplate[]>([]);
  const [currentDeployment, setCurrentDeployment] = useState<FormulaDeployment | undefined>(undefined);
  const [deployments, setDeployments] = useState<FormulaDeployment[]>([]);
  const [totalDeployments, setTotalDeployments] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  const [loading, setLoading] = useState(false);
  const [deploymentsLoading, setDeploymentsLoading] = useState(false);
  const [saveTemplateModalVisible, setSaveTemplateModalVisible] = useState(false);
  const [templateForm] = Form.useForm();
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

      // 加载Formula模板
      loadFormulaTemplates(formulaId);
    } catch (error: any) {
      console.error('Failed to load formula details:', error);
      message.error('加载Formula详情失败');
    }
  };

  // 加载Formula模板
  const loadFormulaTemplates = async (formulaId: number) => {
    try {
      const response = await formulaService.getFormulaTemplates(formulaId);
      setFormulaTemplates(response.templates || []);
    } catch (error: any) {
      console.error('Failed to load formula templates:', error);
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

  // 模板选择处理
  const handleTemplateSelect = (template: FormulaTemplate) => {
    form.setFieldsValue({
      pillar_data: template.pillar_data || {},
    });
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
      const deploymentName = values.name || `${selectedFormula.name} - ${new Date().toLocaleString()}`;

      const requestData = {
        formula_id: selectedFormula.id,
        name: deploymentName,
        description: values.description || '',
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

  // 保存为模板
  const handleSaveAsTemplate = () => {
    if (!selectedFormula) {
      message.warning('请先选择Formula');
      return;
    }

    const values = form.getFieldsValue();
    templateForm.setFieldsValue({
      name: `${selectedFormula.name} 模板`,
      description: `从 ${selectedFormula.name} 创建的模板`,
      pillar_data: values.pillar_data || {},
    });

    setSaveTemplateModalVisible(true);
  };

  // 提交保存模板
  const handleSaveTemplate = async () => {
    try {
      const values = await templateForm.validateFields();
      if (!selectedFormula) return;

      const requestData = {
        formula_id: selectedFormula.id,
        name: values.name,
        description: values.description || '',
        pillar_data: values.pillar_data || {},
        is_public: values.is_public || false,
      };

      await formulaService.createTemplate(selectedFormula.id, requestData);

      message.success('模板保存成功');
      setSaveTemplateModalVisible(false);
      templateForm.resetFields();

      // 重新加载模板
      loadFormulaTemplates(selectedFormula.id);

    } catch (error: any) {
      message.error(error.message || '保存模板失败');
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

  // 部署历史表格列定义
  const deploymentColumns = [
    {
      title: '部署名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: FormulaDeployment) => (
        <div>
          <div style={{ fontWeight: 'bold' }}>{text}</div>
          <div style={{ fontSize: '12px', color: '#666' }}>
            Formula: {formulas.find(f => f.id === record.formula_id)?.name || record.formula_id}
          </div>
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusConfig = {
          pending: { color: 'orange', text: '等待中', icon: <ClockCircleOutlined /> },
          running: { color: 'blue', text: '执行中', icon: <ThunderboltOutlined /> },
          completed: { color: 'green', text: '成功', icon: <CheckCircleOutlined /> },
          failed: { color: 'red', text: '失败', icon: <CloseCircleOutlined /> },
          cancelled: { color: 'gray', text: '已取消', icon: <StopOutlined /> },
        };

        const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.pending;

        return (
          <Tag color={config.color} icon={config.icon}>
            {config.text}
          </Tag>
        );
      },
    },
    {
      title: '目标主机',
      dataIndex: 'target_hosts',
      key: 'target_hosts',
      render: (hosts: string[]) => (
        <div>
          {Array.isArray(hosts) ? hosts.join(', ') : hosts}
          <div style={{ fontSize: '12px', color: '#666' }}>
            共 {Array.isArray(hosts) ? hosts.length : 0} 台
          </div>
        </div>
      ),
    },
    {
      title: '执行结果',
      key: 'results',
      render: (record: FormulaDeployment) => (
        <div>
          <div>成功: <Tag color="green">{record.success_count}</Tag></div>
          <div>失败: <Tag color="red">{record.failed_count}</Tag></div>
          {record.duration_seconds && (
            <div style={{ fontSize: '12px', color: '#666' }}>
              耗时: {Math.round(record.duration_seconds)}s
            </div>
          )}
        </div>
      ),
    },
    {
      title: '执行时间',
      key: 'time',
      render: (record: FormulaDeployment) => (
        <div style={{ fontSize: '12px' }}>
          <div>开始: {record.started_at ? new Date(record.started_at).toLocaleString() : '-'}</div>
          <div>结束: {record.completed_at ? new Date(record.completed_at).toLocaleString() : '-'}</div>
        </div>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: FormulaDeployment) => (
        <Space size="small">
          {record.status === 'running' && (
            <Button
              size="small"
              danger
              onClick={() => handleCancelDeployment(record.id)}
            >
              取消
            </Button>
          )}
          {record.status === 'completed' && record.results && (
            <Button
              size="small"
              onClick={() => setCurrentDeployment(record)}
            >
              查看结果
            </Button>
          )}
        </Space>
      ),
    },
  ];

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
        <Col span={24}>
          <Card title="部署配置" extra={
            <Space>
              <Button icon={<SaveOutlined />} onClick={handleSaveAsTemplate}>
                保存为模板
              </Button>
              <Button
                type="primary"
                icon={<PlayCircleOutlined />}
                onClick={handleDeploy}
                loading={loading}
              >
                执行部署
              </Button>
            </Space>
          }>
            <Form form={form} layout="vertical">
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    label="Formula选择"
                    name="formula_id"
                    rules={[{ required: true, message: '请选择Formula' }]}
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
                </Col>
                <Col span={12}>
                  <Form.Item
                    label="部署名称"
                    name="name"
                    rules={[{ required: true, message: '请输入部署名称' }]}
                  >
                    <Input placeholder="为本次部署命名" />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item label="部署描述" name="description">
                <TextArea rows={2} placeholder="可选：描述本次部署的目的和内容" />
              </Form.Item>

              {selectedFormula && formulaTemplates.length > 0 && (
                <Form.Item label="使用模板">
                  <Space wrap>
                    {formulaTemplates.map(template => (
                      <Button
                        key={template.id}
                        size="small"
                        onClick={() => handleTemplateSelect(template)}
                      >
                        {template.name}
                      </Button>
                    ))}
                  </Space>
                </Form.Item>
              )}

              {selectedFormula && (
                <div style={{ marginTop: '16px' }}>
                  <Title level={4}>Formula信息</Title>
                  <Descriptions size="small" column={2}>
                    <Descriptions.Item label="名称">{selectedFormula.name}</Descriptions.Item>
                    <Descriptions.Item label="分类">{selectedFormula.category}</Descriptions.Item>
                    <Descriptions.Item label="版本">{selectedFormula.version || '最新'}</Descriptions.Item>
                    <Descriptions.Item label="路径">{selectedFormula.path}</Descriptions.Item>
                  </Descriptions>
                  <div style={{ marginTop: '16px' }}>
                    <Text>{selectedFormula.description}</Text>
                  </div>
                </div>
              )}

              {selectedFormula && (
                <div style={{ marginTop: '24px' }}>
                  <Title level={4}>配置参数</Title>
                  {renderParameterFields()}
                </div>
              )}
            </Form>
          </Card>
        </Col>

        <Col span={24}>
          <Card title="目标主机选择">
            <HostSelector
              value={selectedHosts}
              onChange={setSelectedHosts}
            />
          </Card>
        </Col>

        {currentDeployment && (
          <Col span={24}>
            <Card
              title={`当前部署: ${currentDeployment.name}`}
              extra={
                currentDeployment.status === 'running' && (
                  <Button
                    danger
                    onClick={() => handleCancelDeployment(currentDeployment.id)}
                  >
                    取消部署
                  </Button>
                )
              }
            >
              <FormulaResultViewer deployment={currentDeployment} />
            </Card>
          </Col>
        )}

        <Col span={24}>
          <Card
            title="部署历史"
            extra={
              <Button icon={<ReloadOutlined />} onClick={() => loadDeployments()}>
                刷新
              </Button>
            }
          >
            <Table
              columns={deploymentColumns}
              dataSource={deployments}
              rowKey="id"
              loading={deploymentsLoading}
              pagination={{
                current: currentPage,
                pageSize: pageSize,
                total: totalDeployments,
                showSizeChanger: false,
                showQuickJumper: true,
                showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
                onChange: (page) => loadDeployments(page),
              }}
            />
          </Card>
        </Col>
      </Row>

      {/* 保存模板模态框 */}
      <Modal
        title="保存为模板"
        open={saveTemplateModalVisible}
        onOk={handleSaveTemplate}
        onCancel={() => {
          setSaveTemplateModalVisible(false);
          templateForm.resetFields();
        }}
        okText="保存"
        cancelText="取消"
      >
        <Form form={templateForm} layout="vertical">
          <Form.Item
            label="模板名称"
            name="name"
            rules={[{ required: true, message: '请输入模板名称' }]}
          >
            <Input placeholder="为模板命名" />
          </Form.Item>
          <Form.Item label="模板描述" name="description">
            <TextArea rows={2} placeholder="描述模板的使用场景" />
          </Form.Item>
          <Form.Item label="设为公开模板" name="is_public" valuePropName="checked">
            <input type="checkbox" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default FormulaDeployment;
