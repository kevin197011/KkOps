// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import {
  Table,
  Button,
  Space,
  message,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  Tag,
  Tooltip,
  Tree,
  Empty,
  Tabs,
  Descriptions,
  Timeline,
  Checkbox,
} from 'antd'
import type { DataNode } from 'antd/es/tree'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  PauseOutlined,
  HistoryOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SyncOutlined,
  QuestionCircleOutlined,
  DownOutlined,
  UpOutlined,
} from '@ant-design/icons'
import {
  scheduledTaskApi,
  ScheduledTask,
  CreateScheduledTaskRequest,
  UpdateScheduledTaskRequest,
  ScheduledTaskExecution,
  ValidateCronResponse,
} from '@/api/task'
import { executionTemplateApi, ExecutionTemplate } from '@/api/execution'
import { assetApi, Asset } from '@/api/asset'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import moment from 'moment'

const { TextArea } = Input

// Cron 表达式预设（6字段格式：秒 分 时 日 月 周）
const CRON_PRESETS = [
  { label: '每分钟', value: '0 * * * * *' },
  { label: '每 5 分钟', value: '0 */5 * * * *' },
  { label: '每 15 分钟', value: '0 */15 * * * *' },
  { label: '每 30 分钟', value: '0 */30 * * * *' },
  { label: '每小时', value: '0 0 * * * *' },
  { label: '每 2 小时', value: '0 0 */2 * * *' },
  { label: '每天 0 点', value: '0 0 0 * * *' },
  { label: '每天 8 点', value: '0 0 8 * * *' },
  { label: '每周一 0 点', value: '0 0 0 * * 1' },
  { label: '每月 1 日 0 点', value: '0 0 0 1 * *' },
]

// Shebang 配置
const SHEBANGS: Record<string, string> = {
  shell: '#!/usr/bin/env bash\n\n',
  python: '#!/usr/bin/env python3\n\n',
}

const ScheduledTaskList = () => {
  const [tasks, setTasks] = useState<ScheduledTask[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [historyModalVisible, setHistoryModalVisible] = useState(false)
  const [editingTask, setEditingTask] = useState<ScheduledTask | null>(null)
  const [selectedTask, setSelectedTask] = useState<ScheduledTask | null>(null)
  const [executions, setExecutions] = useState<ScheduledTaskExecution[]>([])
  const [executionsLoading, setExecutionsLoading] = useState(false)
  const [expandedLogs, setExpandedLogs] = useState<Set<number>>(new Set())
  const [form] = Form.useForm()

  // 资产相关状态
  const [templates, setTemplates] = useState<ExecutionTemplate[]>([])
  const [assets, setAssets] = useState<Asset[]>([])
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [checkedAssetKeys, setCheckedAssetKeys] = useState<number[]>([])
  const [cronValidation, setCronValidation] = useState<ValidateCronResponse | null>(null)

  // 获取任务列表
  const fetchTasks = useCallback(async () => {
    setLoading(true)
    try {
      const response = await scheduledTaskApi.list(1, 100)
      setTasks(response.data.data || [])
    } catch (error: any) {
      message.error('获取任务列表失败')
    } finally {
      setLoading(false)
    }
  }, [])

  // 获取模板列表
  const fetchTemplates = useCallback(async () => {
    try {
      const response = await executionTemplateApi.list()
      setTemplates(response.data || [])
    } catch (error) {
      console.error('获取模板失败', error)
    }
  }, [])

  // 获取资产列表
  const fetchAssets = useCallback(async () => {
    try {
      const response = await assetApi.list({ page: 1, page_size: 1000 })
      setAssets(response.data.data || [])
    } catch (error) {
      console.error('获取资产失败', error)
    }
  }, [])

  // 获取项目和环境
  const fetchProjectsAndEnvironments = useCallback(async () => {
    try {
      const [projectsRes, environmentsRes] = await Promise.all([
        projectApi.list(),
        environmentApi.list(),
      ])
      setProjects(projectsRes.data || [])
      setEnvironments(environmentsRes.data || [])
    } catch (error) {
      console.error('获取项目/环境失败', error)
    }
  }, [])

  useEffect(() => {
    fetchTasks()
    fetchTemplates()
    fetchAssets()
    fetchProjectsAndEnvironments()
  }, [fetchTasks, fetchTemplates, fetchAssets, fetchProjectsAndEnvironments])

  // 构建资产树
  const buildAssetTree = (): DataNode[] => {
    const tree: DataNode[] = []

    projects.forEach((project) => {
      const projectAssets = assets.filter((a) => a.project_id === project.id)
      if (projectAssets.length === 0) return

      const projectNode: DataNode = {
        key: `project-${project.id}`,
        title: project.name,
        selectable: false,
        children: [],
      }

      environments.forEach((env) => {
        const envAssets = projectAssets.filter((a) => a.environment_id === env.id)
        if (envAssets.length === 0) return

        const envNode: DataNode = {
          key: `env-${project.id}-${env.id}`,
          title: env.name,
          selectable: false,
          children: envAssets.map((asset) => ({
            key: asset.id,
            title: `${asset.hostName} (${asset.ip})`,
          })),
        }

        projectNode.children?.push(envNode)
      })

      // 未分配环境的资产
      const unassignedAssets = projectAssets.filter((a) => !a.environment_id)
      if (unassignedAssets.length > 0) {
        const unassignedNode: DataNode = {
          key: `env-${project.id}-unassigned`,
          title: '未分配环境',
          selectable: false,
          children: unassignedAssets.map((asset) => ({
            key: asset.id,
            title: `${asset.hostName} (${asset.ip})`,
          })),
        }
        projectNode.children?.push(unassignedNode)
      }

      tree.push(projectNode)
    })

    // 未分配项目的资产
    const unassignedProjectAssets = assets.filter((a) => !a.project_id)
    if (unassignedProjectAssets.length > 0) {
      tree.push({
        key: 'project-unassigned',
        title: '未分配项目',
        selectable: false,
        children: unassignedProjectAssets.map((asset) => ({
          key: asset.id,
          title: `${asset.hostName} (${asset.ip})`,
        })),
      })
    }

    return tree
  }

  // 验证 Cron 表达式
  const validateCron = async (cronExpr: string) => {
    if (!cronExpr) {
      setCronValidation(null)
      return
    }
    try {
      const response = await scheduledTaskApi.validateCron(cronExpr)
      setCronValidation(response.data)
    } catch (error) {
      setCronValidation({ valid: false, error: '验证失败' })
    }
  }

  // 处理模板选择
  const handleTemplateSelect = (templateId: number | undefined) => {
    if (!templateId) return
    const template = templates.find((t) => t.id === templateId)
    if (template) {
      form.setFieldsValue({
        content: template.content,
        type: template.type,
      })
    }
  }

  // 处理类型变更
  const handleTypeChange = (type: string) => {
    const currentContent = form.getFieldValue('content') || ''
    const newShebang = SHEBANGS[type] || ''

    // 移除旧的 shebang
    let updatedContent = currentContent.replace(/^#!.*\n\n?/, '')

    // 插入新的 shebang
    if (newShebang && !updatedContent.startsWith(newShebang)) {
      updatedContent = newShebang + updatedContent
    }
    form.setFieldsValue({ content: updatedContent })
  }

  // 打开创建弹窗
  const handleCreate = () => {
    setEditingTask(null)
    setCheckedAssetKeys([])
    setCronValidation(null)
    form.resetFields()
    form.setFieldsValue({ content: SHEBANGS.shell, type: 'shell', timeout: 300 })
    setModalVisible(true)
  }

  // 打开编辑弹窗
  const handleEdit = (task: ScheduledTask) => {
    setEditingTask(task)
    setCheckedAssetKeys(task.asset_ids || [])
    setCronValidation(null)
    form.setFieldsValue({
      ...task,
      template_id: task.template_id || undefined,
    })
    validateCron(task.cron_expression)
    setModalVisible(true)
  }

  // 删除任务
  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个定时任务吗？',
      onOk: async () => {
        try {
          await scheduledTaskApi.delete(id)
          message.success('删除成功')
          fetchTasks()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  // 启用/禁用任务
  const handleToggleEnabled = async (task: ScheduledTask) => {
    try {
      if (task.enabled) {
        await scheduledTaskApi.disable(task.id)
        message.success('任务已禁用')
      } else {
        // 检查是否配置了目标主机
        if (!task.asset_ids || task.asset_ids.length === 0) {
          message.warning('请先配置目标主机再启用任务')
          return
        }
        await scheduledTaskApi.enable(task.id)
        message.success('任务已启用')
      }
      fetchTasks()
    } catch (error: any) {
      message.error('操作失败')
    }
  }

  // 查看执行历史
  const handleViewHistory = async (task: ScheduledTask) => {
    setSelectedTask(task)
    setHistoryModalVisible(true)
    setExecutionsLoading(true)
    setExpandedLogs(new Set()) // 重置展开状态
    try {
      const response = await scheduledTaskApi.getExecutions(task.id, 1, 50)
      setExecutions(response.data.data || [])
    } catch (error) {
      message.error('获取执行历史失败')
    } finally {
      setExecutionsLoading(false)
    }
  }

  // 提交表单
  const handleSubmit = async (
    values: CreateScheduledTaskRequest | UpdateScheduledTaskRequest
  ) => {
    try {
      const data = {
        ...values,
        asset_ids: checkedAssetKeys,
      }

      if (editingTask) {
        await scheduledTaskApi.update(editingTask.id, data as UpdateScheduledTaskRequest)
        message.success('更新成功')
      } else {
        await scheduledTaskApi.create(data as CreateScheduledTaskRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchTasks()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 获取状态标签
  const getStatusTag = (status: string) => {
    const config: Record<string, { color: string; text: string; icon: React.ReactNode }> = {
      success: { color: 'success', text: '成功', icon: <CheckCircleOutlined /> },
      failed: { color: 'error', text: '失败', icon: <CloseCircleOutlined /> },
      partial: { color: 'warning', text: '部分成功', icon: <QuestionCircleOutlined /> },
      running: { color: 'processing', text: '执行中', icon: <SyncOutlined spin /> },
      skipped: { color: 'default', text: '跳过（无主机）', icon: <QuestionCircleOutlined /> },
    }
    const { color, text, icon } = config[status] || { color: 'default', text: status, icon: null }
    return (
      <Tag color={color} icon={icon}>
        {text}
      </Tag>
    )
  }

  // 表格列定义
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: 'Cron 表达式',
      dataIndex: 'cron_expression',
      key: 'cron_expression',
      width: 150,
      render: (cron: string) => <code>{cron}</code>,
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 100,
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'success' : 'default'}>
          {enabled ? '已启用' : '已禁用'}
        </Tag>
      ),
    },
    {
      title: '上次执行',
      key: 'last_run',
      width: 180,
      render: (_: any, record: ScheduledTask) => (
        <Space direction="vertical" size="small">
          {record.last_run_at ? (
            <>
              <span>{moment(record.last_run_at).format('YYYY-MM-DD HH:mm:ss')}</span>
              {record.last_status && getStatusTag(record.last_status)}
            </>
          ) : (
            <span style={{ color: '#999' }}>从未执行</span>
          )}
        </Space>
      ),
    },
    {
      title: '下次执行',
      dataIndex: 'next_run_at',
      key: 'next_run_at',
      width: 180,
      render: (time: string) =>
        time ? moment(time).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 80,
      render: (type: string) => (
        <Tag color={type === 'shell' ? 'blue' : 'green'}>{type}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 280,
      render: (_: any, record: ScheduledTask) => (
        <Space size="small">
          <Tooltip title={record.enabled ? '禁用' : '启用'}>
            <Button
              type="text"
              size="small"
              icon={record.enabled ? <PauseOutlined /> : <PlayCircleOutlined />}
              onClick={() => handleToggleEnabled(record)}
            />
          </Tooltip>
          <Tooltip title="执行历史">
            <Button
              type="text"
              size="small"
              icon={<HistoryOutlined />}
              onClick={() => handleViewHistory(record)}
            />
          </Tooltip>
          <Tooltip title="编辑">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Tooltip title="删除">
            <Button
              type="text"
              danger
              size="small"
              icon={<DeleteOutlined />}
              onClick={() => handleDelete(record.id)}
            />
          </Tooltip>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div
        style={{
          marginBottom: 16,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <h2 style={{ margin: 0 }}>任务管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新建任务
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={tasks}
        rowKey="id"
        loading={loading}
        scroll={{ x: 'max-content' }}
        pagination={{
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
      />

      {/* 创建/编辑弹窗 */}
      <Modal
        title={editingTask ? '编辑任务' : '新建任务'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={800}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Tabs
            defaultActiveKey="basic"
            items={[
              {
                key: 'basic',
                label: '基本信息',
                children: (
                  <>
                    <Form.Item
                      name="name"
                      label="任务名称"
                      rules={[{ required: true, message: '请输入任务名称' }]}
                    >
                      <Input placeholder="输入任务名称" />
                    </Form.Item>

                    <Form.Item name="description" label="描述">
                      <Input.TextArea rows={2} placeholder="任务描述（可选）" />
                    </Form.Item>

                    <Form.Item
                      label={
                        <Space>
                          <span>Cron 表达式</span>
                          <Tooltip title="使用 6 字段格式：秒 分 时 日 月 周，例如：0 * * * * * 表示每分钟执行">
                            <QuestionCircleOutlined />
                          </Tooltip>
                        </Space>
                      }
                      required
                      extra={
                        cronValidation && (
                          <div style={{ marginTop: 8 }}>
                            {cronValidation.valid ? (
                              <span style={{ color: '#52c41a' }}>
                                <ClockCircleOutlined /> 下次执行时间：
                                {moment(cronValidation.next_run_at).format(
                                  'YYYY-MM-DD HH:mm:ss'
                                )}
                              </span>
                            ) : (
                              <span style={{ color: '#ff4d4f' }}>
                                {cronValidation.error || '无效的 Cron 表达式'}
                              </span>
                            )}
                          </div>
                        )
                      }
                    >
                      <Space.Compact style={{ width: '100%' }}>
                        <Form.Item
                          name="cron_expression"
                          noStyle
                          rules={[{ required: true, message: '请输入 Cron 表达式' }]}
                        >
                          <Input
                            placeholder="例如：0 0 * * * *（每小时执行）"
                            onChange={(e) => validateCron(e.target.value)}
                            style={{ flex: 1 }}
                          />
                        </Form.Item>
                        <Select
                          placeholder="快速选择"
                          style={{ width: 150 }}
                          onChange={(value) => {
                            form.setFieldsValue({ cron_expression: value })
                            validateCron(value)
                          }}
                          options={CRON_PRESETS}
                        />
                      </Space.Compact>
                    </Form.Item>

                    <Form.Item name="enabled" label="启用状态" valuePropName="checked">
                      <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                  </>
                ),
              },
              {
                key: 'script',
                label: '执行脚本',
                children: (
                  <>
                    <Form.Item name="template_id" label="从模板创建">
                      <Select
                        placeholder="选择模板（可选）"
                        allowClear
                        onChange={handleTemplateSelect}
                        options={templates.map((t) => ({
                          label: t.name,
                          value: t.id,
                        }))}
                      />
                    </Form.Item>

                    <Form.Item name="type" label="脚本类型" initialValue="shell">
                      <Select onChange={handleTypeChange}>
                        <Select.Option value="shell">Shell</Select.Option>
                        <Select.Option value="python">Python</Select.Option>
                      </Select>
                    </Form.Item>

                    <Form.Item
                      name="content"
                      label="脚本内容"
                      rules={[{ required: true, message: '请输入脚本内容' }]}
                    >
                      <TextArea
                        rows={10}
                        placeholder="输入要执行的脚本内容"
                        style={{ fontFamily: 'monospace' }}
                      />
                    </Form.Item>

                    <Form.Item name="timeout" label="超时时间（秒）" initialValue={300}>
                      <InputNumber min={1} max={86400} style={{ width: '100%' }} />
                    </Form.Item>

                    <Form.Item
                      name="update_assets"
                      valuePropName="checked"
                      tooltip="启用后，脚本输出 JSON 格式的资产信息将自动更新到资产管理"
                    >
                      <Checkbox>自动更新资产信息（CPU/内存/磁盘）</Checkbox>
                    </Form.Item>
                  </>
                ),
              },
              {
                key: 'hosts',
                label: '目标主机',
                children: (
                  <>
                    <div style={{ marginBottom: 8 }}>
                      已选择 {checkedAssetKeys.length} 台主机
                    </div>
                    {buildAssetTree().length > 0 ? (
                      <Tree
                        checkable
                        defaultExpandAll
                        treeData={buildAssetTree()}
                        checkedKeys={checkedAssetKeys}
                        onCheck={(keys) => {
                          const numericKeys = (keys as (string | number)[]).filter(
                            (k) => typeof k === 'number'
                          ) as number[]
                          setCheckedAssetKeys(numericKeys)
                        }}
                        height={300}
                      />
                    ) : (
                      <Empty description="暂无可用主机" />
                    )}
                  </>
                ),
              },
            ]}
          />

          <Form.Item style={{ marginTop: 16, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => setModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit">
                {editingTask ? '更新' : '创建'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 执行历史弹窗 */}
      <Modal
        title={`执行历史 - ${selectedTask?.name}`}
        open={historyModalVisible}
        onCancel={() => setHistoryModalVisible(false)}
        footer={null}
        width={800}
      >
        {selectedTask && (
          <>
            <Descriptions bordered size="small" style={{ marginBottom: 16 }}>
              <Descriptions.Item label="Cron 表达式">
                <code>{selectedTask.cron_expression}</code>
              </Descriptions.Item>
              <Descriptions.Item label="上次执行">
                {selectedTask.last_run_at
                  ? moment(selectedTask.last_run_at).format('YYYY-MM-DD HH:mm:ss')
                  : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="下次执行">
                {selectedTask.next_run_at
                  ? moment(selectedTask.next_run_at).format('YYYY-MM-DD HH:mm:ss')
                  : '-'}
              </Descriptions.Item>
            </Descriptions>

            <Timeline
              pending={executionsLoading ? '加载中...' : false}
              items={executions.map((exec) => {
                const isExpanded = expandedLogs.has(exec.id)
                const toggleExpand = () => {
                  setExpandedLogs((prev) => {
                    const newSet = new Set(prev)
                    if (newSet.has(exec.id)) {
                      newSet.delete(exec.id)
                    } else {
                      newSet.add(exec.id)
                    }
                    return newSet
                  })
                }
                return {
                  color:
                    exec.status === 'success'
                      ? 'green'
                      : exec.status === 'failed'
                      ? 'red'
                      : exec.status === 'running'
                      ? 'blue'
                      : 'gray',
                  children: (
                    <div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
                        <strong>
                          {exec.asset?.hostname || exec.asset?.ip || `资产 #${exec.asset_id}`}
                        </strong>
                        {getStatusTag(exec.status)}
                        <Tag>{exec.trigger_type === 'scheduled' ? '定时' : '手动'}</Tag>
                        {(exec.output || exec.error) && (
                          <Button
                            type="link"
                            size="small"
                            onClick={toggleExpand}
                            style={{ padding: 0 }}
                          >
                            {isExpanded ? <UpOutlined /> : <DownOutlined />}
                            {isExpanded ? ' 收起日志' : ' 查看日志'}
                          </Button>
                        )}
                      </div>
                      <div style={{ color: '#999', fontSize: '12px', marginTop: 4 }}>
                        {moment(exec.created_at).format('YYYY-MM-DD HH:mm:ss')}
                        {exec.finished_at && (
                          <span>
                            {' '}
                            - 耗时{' '}
                            {moment(exec.finished_at).diff(
                              moment(exec.started_at || exec.created_at),
                              'seconds'
                            )}
                            s
                          </span>
                        )}
                        {exec.exit_code !== undefined && exec.exit_code !== null && (
                          <span> - 退出码: {exec.exit_code}</span>
                        )}
                      </div>
                      {exec.error && (
                        <div style={{ color: '#ff4d4f', marginTop: 4 }}>
                          错误: {exec.error}
                        </div>
                      )}
                      {isExpanded && exec.output && (
                        <div
                          style={{
                            marginTop: 8,
                            padding: 12,
                            background: '#1e1e1e',
                            borderRadius: 6,
                            maxHeight: 300,
                            overflow: 'auto',
                          }}
                        >
                          <pre
                            style={{
                              margin: 0,
                              fontSize: 12,
                              fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                              color: '#d4d4d4',
                              whiteSpace: 'pre-wrap',
                              wordBreak: 'break-all',
                            }}
                          >
                            {exec.output}
                          </pre>
                        </div>
                      )}
                    </div>
                  ),
                }
              })}
            />
            {executions.length === 0 && !executionsLoading && (
              <Empty description="暂无执行记录" />
            )}
          </>
        )}
      </Modal>
    </div>
  )
}

export default ScheduledTaskList
