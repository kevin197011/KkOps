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
  Tag,
  Tooltip,
  InputNumber,
  Checkbox,
  Spin,
  Divider,
  Timeline,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  RocketOutlined,
  HistoryOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  LoadingOutlined,
  ClockCircleOutlined,
  StopOutlined,
  ExportOutlined,
  ImportOutlined,
  UploadOutlined,
} from '@ant-design/icons'
import {
  deploymentModuleApi,
  deploymentApi,
  DeploymentModule,
  Deployment,
  CreateModuleRequest,
  UpdateModuleRequest,
  VersionSourceResponse,
  ExportConfig,
  ImportResult,
} from '@/api/deployment'
import { executionTemplateApi, ExecutionTemplate } from '@/api/execution'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import { assetApi, Asset } from '@/api/asset'
import { useThemeStore } from '@/stores/theme'
import { usePermissionStore } from '@/stores/permission'

const { TextArea } = Input

// 不同类型的 shebang
const SHEBANGS: Record<string, string> = {
  shell: '#!/usr/bin/env bash\n\n',
  python: '#!/usr/bin/env python3\n\n',
}

const DeploymentModuleList = () => {
  const { mode } = useThemeStore()
  const { hasPermission } = usePermissionStore()
  const isDark = mode === 'dark'

  const [modules, setModules] = useState<DeploymentModule[]>([])
  const [loading, setLoading] = useState(false)
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [assets, setAssets] = useState<Asset[]>([])
  const [templates, setTemplates] = useState<ExecutionTemplate[]>([])
  const [projectFilter, setProjectFilter] = useState<number | undefined>()

  // 模块编辑弹窗
  const [moduleModalVisible, setModuleModalVisible] = useState(false)
  const [editingModule, setEditingModule] = useState<DeploymentModule | null>(null)
  const [moduleForm] = Form.useForm()
  // 表单中选择的项目和环境，用于筛选主机
  const [formProjectId, setFormProjectId] = useState<number | undefined>()
  const [formEnvironmentId, setFormEnvironmentId] = useState<number | undefined>()
  // 主机搜索过滤
  const [hostSearchText, setHostSearchText] = useState('')

  // 部署执行弹窗
  const [deployModalVisible, setDeployModalVisible] = useState(false)
  const [deployingModule, setDeployingModule] = useState<DeploymentModule | null>(null)
  const [versions, setVersions] = useState<VersionSourceResponse | null>(null)
  const [versionsLoading, setVersionsLoading] = useState(false)
  const [selectedVersion, setSelectedVersion] = useState<string>('')
  const [selectedAssetIds, setSelectedAssetIds] = useState<number[]>([])
  const [deploying, setDeploying] = useState(false)
  // deployHostSearchText 已移除，部署时目标主机为只读

  // 部署历史弹窗
  const [historyModalVisible, setHistoryModalVisible] = useState(false)
  const [historyModule, setHistoryModule] = useState<DeploymentModule | null>(null)
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [historyLoading, setHistoryLoading] = useState(false)

  // 部署详情弹窗
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [detailDeployment, setDetailDeployment] = useState<Deployment | null>(null)

  // 导入导出
  const [importModalVisible, setImportModalVisible] = useState(false)
  const [importPreviewResult, setImportPreviewResult] = useState<ImportResult | null>(null)
  const [importConfig, setImportConfig] = useState<ExportConfig | null>(null)
  const [importing, setImporting] = useState(false)

  const fetchModules = useCallback(async () => {
    setLoading(true)
    try {
      const response = await deploymentModuleApi.list(projectFilter)
      // 处理不同的响应格式
      const data = Array.isArray(response.data) ? response.data : (response.data as any)?.data || []
      setModules(data)
    } catch (error: any) {
      // 404 表示暂无数据，不显示错误
      if (error.response?.status === 404) {
        setModules([])
        return
      }
      console.error('获取部署模块列表失败:', error)
      message.error('获取部署模块列表失败: ' + (error.response?.data?.error || error.message))
    } finally {
      setLoading(false)
    }
  }, [projectFilter])

  const fetchProjects = async () => {
    try {
      const response = await projectApi.list()
      setProjects(response.data || [])
    } catch (error: any) {
      // Ignore
    }
  }

  const fetchEnvironments = async () => {
    try {
      const response = await environmentApi.list()
      setEnvironments(response.data || [])
    } catch (error: any) {
      // Ignore
    }
  }

  const fetchAssets = async () => {
    try {
      const response = await assetApi.list({ page: 1, page_size: 1000 })
      setAssets(response.data?.data || [])
    } catch (error: any) {
      // Ignore
    }
  }

  const fetchTemplates = async () => {
    try {
      const response = await executionTemplateApi.list()
      // 处理不同的响应格式
      const data = Array.isArray(response.data) ? response.data : (response.data as any)?.data || []
      setTemplates(data)
    } catch (error: any) {
      // Ignore
    }
  }

  useEffect(() => {
    fetchModules()
    fetchProjects()
    fetchEnvironments()
    fetchAssets()
    fetchTemplates()
  }, [fetchModules])

  // 根据选择的项目和环境筛选主机
  const filteredAssets = assets.filter((asset) => {
    if (formProjectId && asset.project_id !== formProjectId) return false
    if (formEnvironmentId && asset.environment_id !== formEnvironmentId) return false
    // 搜索过滤
    if (hostSearchText) {
      const search = hostSearchText.toLowerCase()
      const hostname = (asset.hostName || '').toLowerCase()
      const ip = (asset.ip || '').toLowerCase()
      if (!hostname.includes(search) && !ip.includes(search)) return false
    }
    return true
  })

  // 创建/编辑模块
  const handleCreateModule = () => {
    setEditingModule(null)
    moduleForm.resetFields()
    moduleForm.setFieldsValue({
      script_type: 'shell',
      timeout: 600,
      deploy_script: SHEBANGS.shell,
    })
    setFormProjectId(undefined)
    setFormEnvironmentId(undefined)
    setHostSearchText('')
    setModuleModalVisible(true)
  }

  const handleEditModule = (module: DeploymentModule) => {
    setEditingModule(module)
    moduleForm.setFieldsValue({
      project_id: module.project_id,
      environment_id: module.environment_id,
      template_id: module.template_id,
      name: module.name,
      description: module.description,
      version_source_url: module.version_source_url,
      deploy_script: module.deploy_script,
      script_type: module.script_type || 'shell',
      timeout: module.timeout || 600,
      asset_ids: module.asset_ids || [],
    })
    setFormProjectId(module.project_id)
    setFormEnvironmentId(module.environment_id || undefined)
    setHostSearchText('')
    setModuleModalVisible(true)
  }

  const handleDeleteModule = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个部署模块吗？相关的部署历史也将被保留。',
      onOk: async () => {
        try {
          await deploymentModuleApi.delete(id)
          message.success('删除成功')
          fetchModules()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  // 导出配置
  const handleExport = async () => {
    try {
      const response = await deploymentModuleApi.exportConfig()
      const config = response.data
      const blob = new Blob([JSON.stringify(config, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `deployment-modules-${new Date().toISOString().slice(0, 10)}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (error: any) {
      message.error(error.response?.data?.error || '导出失败')
    }
  }

  // 处理文件上传
  const handleImportFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const content = e.target?.result as string
        const config = JSON.parse(content) as ExportConfig
        
        if (!config.modules || !Array.isArray(config.modules)) {
          message.error('无效的配置文件格式')
          return
        }

        setImportConfig(config)
        
        // 预览导入
        const response = await deploymentModuleApi.previewImport(config)
        setImportPreviewResult(response.data)
        setImportModalVisible(true)
      } catch (error: any) {
        message.error('解析配置文件失败: ' + (error.message || '未知错误'))
      }
    }
    reader.readAsText(file)
    
    // 重置 input 以便可以重新选择同一文件
    event.target.value = ''
  }

  // 执行导入
  const handleImportConfirm = async () => {
    if (!importConfig) return

    setImporting(true)
    try {
      const response = await deploymentModuleApi.importConfig(importConfig)
      const result = response.data
      
      if (result.failed === 0) {
        message.success(`导入成功: ${result.success} 个模块`)
      } else {
        message.warning(`导入完成: ${result.success} 成功, ${result.failed} 失败`)
      }
      
      setImportModalVisible(false)
      setImportConfig(null)
      setImportPreviewResult(null)
      fetchModules()
    } catch (error: any) {
      message.error(error.response?.data?.error || '导入失败')
    } finally {
      setImporting(false)
    }
  }

  const handleModuleSubmit = async (values: CreateModuleRequest | UpdateModuleRequest) => {
    try {
      if (editingModule) {
        await deploymentModuleApi.update(editingModule.id, values as UpdateModuleRequest)
        message.success('更新成功')
      } else {
        await deploymentModuleApi.create(values as CreateModuleRequest)
        message.success('创建成功')
      }
      setModuleModalVisible(false)
      moduleForm.resetFields()
      fetchModules()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 模板选择变更，继承模板内容
  const handleTemplateChange = (templateId: number | undefined) => {
    if (!templateId) {
      // 清除模板时不清空脚本（保持现有内容）
      return
    }
    const template = templates.find((t) => t.id === templateId)
    if (template) {
      moduleForm.setFieldsValue({
        deploy_script: template.content,
        script_type: template.type || 'shell',
      })
    }
  }

  // 脚本类型变更自动插入 shebang
  const handleScriptTypeChange = (newType: string) => {
    const currentScript = moduleForm.getFieldValue('deploy_script') || ''
    const newShebang = SHEBANGS[newType] || ''

    const isEmptyOrShebangOnly =
      !currentScript.trim() ||
      Object.values(SHEBANGS).some((shebang) => currentScript.trim() === shebang.trim())

    if (isEmptyOrShebangOnly) {
      moduleForm.setFieldsValue({ deploy_script: newShebang })
    } else {
      let updatedScript = currentScript
      for (const shebang of Object.values(SHEBANGS)) {
        if (currentScript.startsWith(shebang.trim())) {
          updatedScript = currentScript.replace(shebang.trim(), newShebang.trim())
          break
        }
      }
      if (updatedScript === currentScript && !currentScript.startsWith('#!')) {
        updatedScript = newShebang + currentScript
      }
      moduleForm.setFieldsValue({ deploy_script: updatedScript })
    }
  }

  // 打开部署弹窗
  const handleOpenDeploy = async (module: DeploymentModule) => {
    setDeployingModule(module)
    setSelectedAssetIds(module.asset_ids || [])
    setSelectedVersion('')
    setVersions(null)
    setDeployModalVisible(true)

    // 获取版本列表
    if (module.version_source_url) {
      setVersionsLoading(true)
      try {
        const response = await deploymentModuleApi.getVersions(module.id)
        setVersions(response.data)
        if (response.data?.latest) {
          setSelectedVersion(response.data.latest)
        }
      } catch (error: any) {
        message.warning('获取版本列表失败，请手动输入版本')
      } finally {
        setVersionsLoading(false)
      }
    }
  }

  // 执行部署
  const handleDeploy = async () => {
    if (!deployingModule) return
    if (!selectedVersion) {
      message.warning('请选择或输入版本')
      return
    }
    if (selectedAssetIds.length === 0) {
      message.warning('请选择目标主机')
      return
    }

    setDeploying(true)
    try {
      const response = await deploymentModuleApi.deploy(deployingModule.id, {
        version: selectedVersion,
        asset_ids: selectedAssetIds,
      })
      message.success('部署已开始执行')
      setDeployModalVisible(false)
      // 自动打开执行详情弹窗显示日志
      if (response.data) {
        setDetailDeployment(response.data)
        setDetailModalVisible(true)
        // 轮询刷新部署状态直到完成
        pollDeploymentStatus(response.data.id)
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '部署失败')
    } finally {
      setDeploying(false)
    }
  }

  // 轮询部署状态
  const pollDeploymentStatus = async (deploymentId: number) => {
    const poll = async () => {
      try {
        const response = await deploymentApi.get(deploymentId)
        setDetailDeployment(response.data)
        // 如果还在执行中，继续轮询
        if (response.data.status === 'pending' || response.data.status === 'running') {
          setTimeout(poll, 2000) // 每2秒刷新一次
        }
      } catch (error) {
        // 忽略错误
      }
    }
    poll()
  }

  // 打开历史弹窗
  const handleOpenHistory = async (module: DeploymentModule) => {
    setHistoryModule(module)
    setHistoryModalVisible(true)
    setHistoryLoading(true)
    try {
      const response = await deploymentApi.list(module.id, 1, 20)
      setDeployments(response.data?.data || [])
    } catch (error: any) {
      message.error('获取部署历史失败')
    } finally {
      setHistoryLoading(false)
    }
  }

  // 查看部署详情
  const handleViewDetail = async (deployment: Deployment) => {
    setDetailDeployment(deployment)
    setDetailModalVisible(true)
  }

  // 获取状态标签
  const getStatusTag = (status: string) => {
    const config: Record<string, { color: string; icon: React.ReactNode; text: string }> = {
      pending: { color: 'default', icon: <ClockCircleOutlined />, text: '待执行' },
      running: { color: 'processing', icon: <LoadingOutlined />, text: '执行中' },
      success: { color: 'success', icon: <CheckCircleOutlined />, text: '成功' },
      failed: { color: 'error', icon: <CloseCircleOutlined />, text: '失败' },
      cancelled: { color: 'warning', icon: <StopOutlined />, text: '已取消' },
    }
    const c = config[status] || config.pending
    return (
      <Tag color={c.color} icon={c.icon}>
        {c.text}
      </Tag>
    )
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '模块名称', dataIndex: 'name', key: 'name' },
    {
      title: '所属项目',
      dataIndex: 'project_name',
      key: 'project_name',
      render: (name: string) => <Tag color="purple">{name}</Tag>,
    },
    {
      title: '所属环境',
      dataIndex: 'environment_name',
      key: 'environment_name',
      render: (name: string) => name ? <Tag color="blue">{name}</Tag> : <Tag color="default">全部环境</Tag>,
    },
    {
      title: '执行模板',
      dataIndex: 'template',
      key: 'template',
      render: (template: DeploymentModule['template']) =>
        template ? (
          <Tooltip title={`来自模板: ${template.name}`}>
            <Tag color="orange">{template.name}</Tag>
          </Tooltip>
        ) : (
          <Tag color="default">自定义</Tag>
        ),
    },
    {
      title: '脚本类型',
      dataIndex: 'script_type',
      key: 'script_type',
      render: (type: string) => (
        <Tag color={type === 'shell' ? 'blue' : 'green'}>{type || 'shell'}</Tag>
      ),
    },
    {
      title: '版本源',
      dataIndex: 'version_source_url',
      key: 'version_source_url',
      render: (url: string) =>
        url ? (
          <Tooltip title={url}>
            <Tag color="cyan">已配置</Tag>
          </Tooltip>
        ) : (
          <Tag color="default">未配置</Tag>
        ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 220,
      render: (_: any, record: DeploymentModule) => {
        const actions = []
        if (hasPermission('deployments', 'create')) {
          actions.push(
            <Button
              key="deploy"
              type="primary"
              size="small"
              icon={<RocketOutlined />}
              onClick={() => handleOpenDeploy(record)}
            >
              部署
            </Button>
          )
        }
        if (hasPermission('deployments', 'read')) {
          actions.push(
            <Button
              key="history"
              type="link"
              size="small"
              icon={<HistoryOutlined />}
              onClick={() => handleOpenHistory(record)}
            >
              历史
            </Button>
          )
        }
        if (hasPermission('deployments', 'update')) {
          actions.push(
            <Button
              key="edit"
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEditModule(record)}
            >
              编辑
            </Button>
          )
        }
        if (hasPermission('deployments', 'delete')) {
          actions.push(
            <Button
              key="delete"
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
              onClick={() => handleDeleteModule(record.id)}
            >
              删除
            </Button>
          )
        }
        return actions.length > 0 ? <Space size="small">{actions}</Space> : '-'
      },
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
          flexWrap: 'wrap',
          gap: 8,
        }}
      >
        <h2>部署管理</h2>
        <Space>
          <Select
            placeholder="按项目筛选"
            allowClear
            style={{ width: 200 }}
            value={projectFilter}
            onChange={setProjectFilter}
            options={projects.map((p) => ({ label: p.name, value: p.id }))}
          />
          <Button icon={<ReloadOutlined />} onClick={fetchModules}>
            刷新
          </Button>
          {hasPermission('deployments', 'export') && (
            <Button icon={<ExportOutlined />} onClick={handleExport}>
              导出配置
            </Button>
          )}
          {hasPermission('deployments', 'import') && (
            <>
              <Button
                icon={<ImportOutlined />}
                onClick={() => document.getElementById('import-file-input')?.click()}
              >
                导入配置
              </Button>
              <input
                id="import-file-input"
                type="file"
                accept=".json"
                style={{ display: 'none' }}
                onChange={handleImportFileSelect}
              />
            </>
          )}
          {hasPermission('deployments', 'create') && (
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateModule}>
              新增模块
            </Button>
          )}
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={modules}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />

      {/* 模块编辑弹窗 */}
      <Modal
        title={editingModule ? '编辑部署模块' : '新增部署模块'}
        open={moduleModalVisible}
        onCancel={() => {
          setModuleModalVisible(false)
          moduleForm.resetFields()
        }}
        onOk={() => moduleForm.submit()}
        width={700}
        destroyOnHidden
      >
        <Form form={moduleForm} layout="vertical" onFinish={handleModuleSubmit}>
          <Form.Item
            name="project_id"
            label="所属项目"
            rules={[{ required: true, message: '请选择项目' }]}
          >
            <Select
              placeholder="选择项目"
              options={projects.map((p) => ({ label: p.name, value: p.id }))}
              onChange={(value) => {
                setFormProjectId(value)
                // 清空已选主机（因为项目变了）
                moduleForm.setFieldsValue({ asset_ids: [] })
              }}
            />
          </Form.Item>
          <Form.Item name="environment_id" label="所属环境">
            <Select
              placeholder="选择环境（可选，不选则适用所有环境）"
              allowClear
              options={environments.map((e) => ({ label: e.name, value: e.id }))}
              onChange={(value) => {
                setFormEnvironmentId(value)
                // 清空已选主机（因为环境变了）
                moduleForm.setFieldsValue({ asset_ids: [] })
              }}
            />
          </Form.Item>
          <Form.Item name="template_id" label="从执行模板创建">
            <Select
              placeholder="选择执行模板（可选，选择后自动填充脚本）"
              allowClear
              showSearch
              optionFilterProp="label"
              options={templates.map((t) => ({
                label: `${t.name} (${t.type || 'shell'})`,
                value: t.id,
              }))}
              onChange={handleTemplateChange}
            />
          </Form.Item>
          <div style={{ fontSize: 12, color: '#888', marginTop: -16, marginBottom: 16 }}>
            选择模板后将自动填充脚本内容和类型，您可以在此基础上进行修改
          </div>
          <Form.Item
            name="name"
            label="模块名称"
            rules={[{ required: true, message: '请输入模块名称' }]}
          >
            <Input placeholder="模块名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input placeholder="模块描述" />
          </Form.Item>
          <Form.Item name="version_source_url" label="版本数据源 URL">
            <Input placeholder="HTTP JSON 数据源地址，如: https://api.example.com/versions" />
          </Form.Item>
          <div style={{ fontSize: 12, color: '#888', marginTop: -16, marginBottom: 16 }}>
            数据源需返回 JSON 格式：{`{"versions": ["v1.0.0", "v1.1.0"], "latest": "v1.1.0"}`}
          </div>
          <Form.Item name="script_type" label="脚本类型">
            <Select onChange={handleScriptTypeChange}>
              <Select.Option value="shell">Shell</Select.Option>
              <Select.Option value="python">Python</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="deploy_script"
            label="部署脚本"
            rules={[{ required: true, message: '请输入部署脚本' }]}
          >
            <TextArea
              rows={10}
              placeholder="部署脚本内容"
              style={{
                fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                fontSize: 13,
                background: isDark ? '#141414' : '#fafafa',
              }}
            />
          </Form.Item>
          <div style={{ fontSize: 12, color: '#888', marginTop: -16, marginBottom: 16 }}>
            支持变量替换：<code>${`{VERSION}`}</code> <code>${`{MODULE_NAME}`}</code>{' '}
            <code>${`{PROJECT_NAME}`}</code> <code>${`{ENVIRONMENT_NAME}`}</code>
          </div>
          <Form.Item name="timeout" label="超时时间（秒）">
            <InputNumber min={10} max={86400} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            label={
              <Space>
                <span>默认目标主机</span>
                <span style={{ fontSize: 12, color: '#888' }}>
                  ({filteredAssets.length} 台可选)
                </span>
              </Space>
            }
          >
            <Input
              placeholder="搜索主机名或IP..."
              value={hostSearchText}
              onChange={(e) => setHostSearchText(e.target.value)}
              allowClear
              style={{ marginBottom: 8 }}
            />
            <Form.Item name="asset_ids" noStyle>
              <Checkbox.Group style={{ width: '100%' }}>
                <div style={{ maxHeight: 200, overflowY: 'auto', border: '1px solid #d9d9d9', borderRadius: 6, padding: 8 }}>
                  {filteredAssets.length === 0 ? (
                    <div style={{ color: '#888', padding: 16, textAlign: 'center' }}>
                      {formProjectId ? (hostSearchText ? '无匹配主机' : '该项目/环境下暂无主机') : '请先选择项目'}
                    </div>
                  ) : (
                    filteredAssets.map((asset) => (
                      <div key={asset.id} style={{ marginBottom: 4 }}>
                        <Checkbox value={asset.id}>
                          {asset.hostName || asset.ip} ({asset.ip})
                        </Checkbox>
                      </div>
                    ))
                  )}
                </div>
              </Checkbox.Group>
            </Form.Item>
          </Form.Item>
        </Form>
      </Modal>

      {/* 部署执行弹窗 */}
      <Modal
        title={
          <Space>
            <RocketOutlined />
            <span>部署模块: {deployingModule?.name}</span>
          </Space>
        }
        open={deployModalVisible}
        onCancel={() => setDeployModalVisible(false)}
        onOk={handleDeploy}
        confirmLoading={deploying}
        okText="执行部署"
        width={600}
      >
        <Form layout="vertical">
          <Form.Item label="选择版本" required>
            {versionsLoading ? (
              <Spin />
            ) : versions && versions.versions.length > 0 ? (
              <Select
                placeholder="选择版本"
                value={selectedVersion}
                onChange={setSelectedVersion}
                showSearch
                allowClear
                options={versions.versions.map((v) => ({
                  label: v === versions.latest ? `${v} (latest)` : v,
                  value: v,
                }))}
              />
            ) : (
              <Input
                placeholder="手动输入版本号"
                value={selectedVersion}
                onChange={(e) => setSelectedVersion(e.target.value)}
              />
            )}
            {!deployingModule?.version_source_url && (
              <div style={{ fontSize: 12, color: '#888', marginTop: 4 }}>
                未配置版本数据源，请手动输入版本号
              </div>
            )}
          </Form.Item>
          <Form.Item label="目标主机（已预配置，不可修改）">
            {(() => {
              // 显示部署模块预配置的目标主机
              const targetAssets = assets.filter((asset) => 
                selectedAssetIds.includes(asset.id)
              )
              return (
                <div style={{ 
                  maxHeight: 200, 
                  overflowY: 'auto', 
                  border: '1px solid #d9d9d9', 
                  borderRadius: 6, 
                  padding: 12,
                  background: isDark ? '#1f1f1f' : '#fafafa'
                }}>
                  {targetAssets.length === 0 ? (
                    <div style={{ color: '#888', textAlign: 'center' }}>
                      该模块未配置目标主机
                    </div>
                  ) : (
                    targetAssets.map((asset) => (
                      <Tag 
                        key={asset.id} 
                        color="blue"
                        style={{ marginBottom: 4, marginRight: 4 }}
                      >
                        {asset.hostName || asset.ip} ({asset.ip})
                      </Tag>
                    ))
                  )}
                </div>
              )
            })()}
            <div style={{ fontSize: 12, color: '#888', marginTop: 4 }}>
              共 {selectedAssetIds.length} 台主机，如需修改请编辑部署模块配置
            </div>
          </Form.Item>
        </Form>
      </Modal>

      {/* 部署历史弹窗 */}
      <Modal
        title={
          <Space>
            <HistoryOutlined />
            <span>部署历史: {historyModule?.name}</span>
          </Space>
        }
        open={historyModalVisible}
        onCancel={() => setHistoryModalVisible(false)}
        footer={null}
        width={700}
      >
        {historyLoading ? (
          <div style={{ textAlign: 'center', padding: 50 }}>
            <Spin size="large" />
          </div>
        ) : deployments.length === 0 ? (
          <div style={{ textAlign: 'center', color: '#888', padding: 50 }}>暂无部署记录</div>
        ) : (
          <Timeline>
            {deployments.map((d) => (
              <Timeline.Item key={d.id} color={d.status === 'success' ? 'green' : d.status === 'failed' ? 'red' : 'gray'}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <div>
                    <strong>版本: {d.version}</strong>
                    <span style={{ marginLeft: 12 }}>{getStatusTag(d.status)}</span>
                    <div style={{ fontSize: 12, color: '#888', marginTop: 4 }}>
                      {d.creator_name} · {d.created_at}
                    </div>
                  </div>
                  <Button type="link" size="small" onClick={() => handleViewDetail(d)}>
                    查看详情
                  </Button>
                </div>
              </Timeline.Item>
            ))}
          </Timeline>
        )}
      </Modal>

      {/* 部署详情弹窗 */}
      <Modal
        title={
          <Space>
            <span>部署详情</span>
            {detailDeployment && getStatusTag(detailDeployment.status)}
          </Space>
        }
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={800}
        zIndex={1100}
      >
        {detailDeployment && (
          <div>
            <div style={{ marginBottom: 16 }}>
              <strong>版本:</strong> {detailDeployment.version}
              <span style={{ marginLeft: 24 }}>
                <strong>模块:</strong> {detailDeployment.module_name}
              </span>
              <span style={{ marginLeft: 24 }}>
                <strong>项目:</strong> {detailDeployment.project_name}
              </span>
            </div>
            <div style={{ marginBottom: 16 }}>
              <strong>执行者:</strong> {detailDeployment.creator_name}
              <span style={{ marginLeft: 24 }}>
                <strong>开始时间:</strong> {detailDeployment.started_at || '-'}
              </span>
              <span style={{ marginLeft: 24 }}>
                <strong>结束时间:</strong> {detailDeployment.finished_at || '-'}
              </span>
            </div>
            <Divider />
            {/* 执行中状态显示 */}
            {(detailDeployment.status === 'pending' || detailDeployment.status === 'running') && 
             !detailDeployment.output && !detailDeployment.error && (
              <div style={{ textAlign: 'center', padding: 40 }}>
                <Spin size="large" />
                <div style={{ marginTop: 16, color: '#888' }}>
                  {detailDeployment.status === 'pending' ? '等待执行...' : '正在执行中，请稍候...'}
                </div>
              </div>
            )}
            {/* 执行输出 */}
            {detailDeployment.output && (
              <>
                <div style={{ marginBottom: 8, fontWeight: 500 }}>
                  执行输出:
                  {(detailDeployment.status === 'running') && (
                    <Spin size="small" style={{ marginLeft: 8 }} />
                  )}
                </div>
                <pre
                  style={{
                    background: isDark ? '#141414' : '#f5f5f5',
                    padding: 12,
                    borderRadius: 6,
                    maxHeight: 400,
                    overflow: 'auto',
                    fontSize: 12,
                    fontFamily: 'Monaco, Menlo, monospace',
                  }}
                >
                  {detailDeployment.output}
                </pre>
              </>
            )}
            {/* 错误信息 */}
            {detailDeployment.error && (
              <>
                <div style={{ marginBottom: 8, marginTop: 16, fontWeight: 500, color: '#ff4d4f' }}>
                  错误信息:
                </div>
                <pre
                  style={{
                    background: isDark ? '#2a1515' : '#fff2f0',
                    padding: 12,
                    borderRadius: 6,
                    maxHeight: 200,
                    overflow: 'auto',
                    fontSize: 12,
                    fontFamily: 'Monaco, Menlo, monospace',
                    color: '#ff4d4f',
                  }}
                >
                  {detailDeployment.error}
                </pre>
              </>
            )}
          </div>
        )}
      </Modal>

      {/* 导入预览弹窗 */}
      <Modal
        title="导入配置预览"
        open={importModalVisible}
        onCancel={() => {
          setImportModalVisible(false)
          setImportConfig(null)
          setImportPreviewResult(null)
        }}
        onOk={handleImportConfirm}
        okText="确认导入"
        cancelText="取消"
        confirmLoading={importing}
        width={600}
      >
        {importConfig && importPreviewResult && (
          <div>
            <div style={{ marginBottom: 16 }}>
              <strong>配置文件版本:</strong> {importConfig.version}
              <span style={{ marginLeft: 24 }}>
                <strong>导出时间:</strong> {importConfig.export_at}
              </span>
            </div>
            
            <div style={{ marginBottom: 16 }}>
              <strong>模块总数:</strong> {importPreviewResult.total}
              <span style={{ marginLeft: 16 }}>
                <Tag color="green">可导入: {importPreviewResult.success}</Tag>
              </span>
              {importPreviewResult.failed > 0 && (
                <span style={{ marginLeft: 8 }}>
                  <Tag color="red">失败: {importPreviewResult.failed}</Tag>
                </span>
              )}
            </div>

            {/* 待导入的模块列表 */}
            <div style={{ marginBottom: 16 }}>
              <strong>模块列表:</strong>
              <div
                style={{
                  background: isDark ? '#141414' : '#fafafa',
                  padding: 12,
                  borderRadius: 6,
                  marginTop: 8,
                  maxHeight: 200,
                  overflow: 'auto',
                }}
              >
                {importConfig.modules.map((m, index) => (
                  <div key={index} style={{ marginBottom: 8 }}>
                    <Tag color="blue">{m.name}</Tag>
                    <span style={{ marginLeft: 8, color: isDark ? '#888' : '#666' }}>
                      {m.project_name}
                      {m.environment_name && ` / ${m.environment_name}`}
                    </span>
                    {m.target_hosts && m.target_hosts.length > 0 && (
                      <div style={{ marginTop: 4, marginLeft: 8 }}>
                        <span style={{ fontSize: 12, color: isDark ? '#666' : '#999' }}>
                          目标主机: {m.target_hosts.join(', ')}
                        </span>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {/* 警告信息 */}
            {importPreviewResult.skipped && importPreviewResult.skipped.length > 0 && (
              <div style={{ marginBottom: 16 }}>
                <strong style={{ color: '#faad14' }}>提示信息:</strong>
                <div
                  style={{
                    background: isDark ? '#2b2111' : '#fffbe6',
                    padding: 12,
                    borderRadius: 6,
                    marginTop: 8,
                    maxHeight: 150,
                    overflow: 'auto',
                  }}
                >
                  {importPreviewResult.skipped.map((msg, index) => (
                    <div key={index} style={{ marginBottom: 4, fontSize: 12 }}>
                      • {msg}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* 错误信息 */}
            {importPreviewResult.errors && importPreviewResult.errors.length > 0 && (
              <div>
                <strong style={{ color: '#ff4d4f' }}>错误信息:</strong>
                <div
                  style={{
                    background: isDark ? '#2a1515' : '#fff2f0',
                    padding: 12,
                    borderRadius: 6,
                    marginTop: 8,
                    maxHeight: 150,
                    overflow: 'auto',
                  }}
                >
                  {importPreviewResult.errors.map((msg, index) => (
                    <div key={index} style={{ marginBottom: 4, fontSize: 12, color: '#ff4d4f' }}>
                      • {msg}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default DeploymentModuleList
