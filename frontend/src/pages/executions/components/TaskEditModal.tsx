// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import {
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Checkbox,
  message,
  Spin,
  Tree,
  Empty,
  Button,
  Space,
  theme,
} from 'antd'
import type { DataNode } from 'antd/es/tree'
import { executionApi, Execution, templateApi, ExecutionTemplate } from '@/api/execution'
import { assetApi, Asset } from '@/api/asset'
import { Project } from '@/api/project'
import { Environment } from '@/api/environment'
import { useThemeStore } from '@/stores/theme'

const { TextArea } = Input

// 不同类型的 shebang
const SHEBANGS: Record<string, string> = {
  shell: '#!/usr/bin/env bash\n\n',
  python: '#!/usr/bin/env python3\n\n',
}

interface TaskEditModalProps {
  visible: boolean
  task: Execution | null
  projects: Project[]
  environments: Environment[]
  onClose: () => void
  onSave: () => void
}

const TaskEditModal = ({
  visible,
  task,
  projects,
  environments,
  onClose,
  onSave,
}: TaskEditModalProps) => {
  const { token } = theme.useToken()
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [assets, setAssets] = useState<Asset[]>([])
  const [templates, setTemplates] = useState<ExecutionTemplate[]>([])
  const [selectedAssetIds, setSelectedAssetIds] = useState<number[]>([])
  const [projectFilter, setProjectFilter] = useState<number | null>(null)
  const [environmentFilter, setEnvironmentFilter] = useState<number | null>(null)
  const [useTemplate, setUseTemplate] = useState(false)
  const [selectedTemplateId, setSelectedTemplateId] = useState<number | null>(null)

  const isEditing = !!task

  // Fetch assets and templates
  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const [assetsRes, templatesRes] = await Promise.all([
        assetApi.list({ page: 1, pageSize: 1000 }),
        templateApi.list(),
      ])
      // Handle different response formats
      const assetsData = assetsRes.data?.data || assetsRes.data || []
      const templatesData = templatesRes.data?.data || templatesRes.data || []
      setAssets(Array.isArray(assetsData) ? assetsData : [])
      setTemplates(Array.isArray(templatesData) ? templatesData : [])
    } catch (error) {
      console.error('加载数据失败:', error)
      message.error('加载数据失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (visible) {
      fetchData()
      if (task) {
        form.setFieldsValue({
          name: task.name,
          type: task.type || 'shell',
          content: task.content,
          timeout: task.timeout,
        })
        // Set selected assets
        setSelectedAssetIds(task.asset_ids || [])
        setUseTemplate(false)
        setSelectedTemplateId(null)
      } else {
        form.resetFields()
        form.setFieldsValue({
          type: 'shell',
          timeout: 600,
          content: SHEBANGS.shell,
        })
        setSelectedAssetIds([])
        setUseTemplate(false)
        setSelectedTemplateId(null)
        setProjectFilter(null)
        setEnvironmentFilter(null)
      }
    }
  }, [visible, task, form, fetchData])

  // Filter assets based on project and environment
  const filteredAssets = assets.filter((asset) => {
    if (projectFilter && asset.project_id !== projectFilter) return false
    if (environmentFilter && asset.environment_id !== environmentFilter) return false
    return true
  })

  // Build tree data for asset selection
  const buildAssetTree = (): DataNode[] => {
    const projectMap = new Map<number, { name: string; children: Asset[] }>()

    filteredAssets.forEach((asset) => {
      const projectId = asset.project_id || 0
      const project = projects.find((p) => p.id === projectId)
      const projectName = project?.name || '未分配项目'

      if (!projectMap.has(projectId)) {
        projectMap.set(projectId, { name: projectName, children: [] })
      }
      projectMap.get(projectId)!.children.push(asset)
    })

    return Array.from(projectMap.entries()).map(([projectId, { name, children }]) => ({
      key: `project-${projectId}`,
      title: `${name} (${children.length})`,
      children: children.map((asset) => ({
        key: `asset-${asset.id}`,
        title: `${asset.hostName || asset.ip} (${asset.ip})`,
      })),
    }))
  }

  const handleTemplateSelect = (templateId: number) => {
    setSelectedTemplateId(templateId)
    const template = templates.find((t) => t.id === templateId)
    if (template) {
      form.setFieldsValue({
        type: template.type || 'shell',
        content: template.content || '',
      })
    }
  }

  const handleAssetCheck = (checkedKeys: React.Key[] | { checked: React.Key[] }) => {
    const keys = Array.isArray(checkedKeys) ? checkedKeys : checkedKeys.checked
    const assetIds = keys
      .filter((key) => typeof key === 'string' && key.startsWith('asset-'))
      .map((key) => parseInt((key as string).replace('asset-', ''), 10))
    setSelectedAssetIds(assetIds)
  }

  // 处理类型变更，自动插入对应的 shebang
  const handleTypeChange = (newType: string) => {
    const currentContent = form.getFieldValue('content') || ''
    const newShebang = SHEBANGS[newType] || ''
    
    // 检查当前内容是否为空或只包含其他类型的 shebang
    const isEmptyOrShebangOnly = !currentContent.trim() || 
      Object.values(SHEBANGS).some(shebang => currentContent.trim() === shebang.trim())
    
    if (isEmptyOrShebangOnly) {
      form.setFieldsValue({ content: newShebang })
    } else {
      // 如果内容不为空，检查是否以其他 shebang 开头，替换它
      let updatedContent = currentContent
      for (const shebang of Object.values(SHEBANGS)) {
        if (currentContent.startsWith(shebang.trim())) {
          updatedContent = currentContent.replace(shebang.trim(), newShebang.trim())
          break
        }
      }
      // 如果没有任何 shebang，在开头插入
      if (updatedContent === currentContent && !currentContent.startsWith('#!')) {
        updatedContent = newShebang + currentContent
      }
      form.setFieldsValue({ content: updatedContent })
    }
  }

  // 选择所有主机
  const handleSelectAll = () => {
    if (filteredAssets.length === 0) {
      message.info('没有可选的主机')
      return
    }
    const allAssetIds = filteredAssets.map((a) => a.id)
    setSelectedAssetIds(allAssetIds)
    message.success(`已选择 ${allAssetIds.length} 个主机`)
  }

  // 清空选择
  const handleClearSelection = () => {
    setSelectedAssetIds([])
  }

  const handleSave = async () => {
    try {
      const values = await form.validateFields()

      if (selectedAssetIds.length === 0) {
        message.warning('请至少选择一个执行主机')
        return
      }

      setSaving(true)

      const taskData = {
        ...values,
        asset_ids: selectedAssetIds,
      }

      if (isEditing && task) {
        await executionApi.update(task.id, taskData)
        message.success('任务更新成功')
      } else {
        await executionApi.create(taskData)
        message.success('任务创建成功')
      }

      onSave()
    } catch (error: any) {
      if (error.errorFields) {
        // Form validation error
        return
      }
      message.error(error.response?.data?.error || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const treeData = buildAssetTree()
  const checkedKeys = selectedAssetIds.map((id) => `asset-${id}`)

  return (
    <Modal
      title={isEditing ? '编辑任务' : '新建任务'}
      open={visible}
      onCancel={onClose}
      onOk={handleSave}
      confirmLoading={saving}
      width={800}
      okText="保存"
      cancelText="取消"
      styles={{
        body: { maxHeight: '70vh', overflow: 'auto' },
      }}
    >
      {loading ? (
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Spin />
        </div>
      ) : (
        <Form form={form} layout="vertical">
          {/* Task Name */}
          <Form.Item
            name="name"
            label="任务名称"
            rules={[{ required: true, message: '请输入任务名称' }]}
          >
            <Input placeholder="例如: 部署应用" />
          </Form.Item>

          {/* Template Selection */}
          <Form.Item label="使用模板">
            <Checkbox 
              checked={useTemplate} 
              onChange={(e) => {
                setUseTemplate(e.target.checked)
                if (!e.target.checked) {
                  setSelectedTemplateId(null)
                }
              }}
            >
              从模板创建
            </Checkbox>
            {useTemplate && (
              <Select
                style={{ width: '100%', marginTop: 8 }}
                placeholder="选择模板"
                value={selectedTemplateId}
                onChange={handleTemplateSelect}
                options={templates.map((t) => ({
                  label: `${t.name}${t.description ? ` - ${t.description}` : ''}`,
                  value: t.id,
                }))}
                notFoundContent={templates.length === 0 ? "暂无模板，请先创建任务模板" : "未找到匹配的模板"}
              />
            )}
          </Form.Item>

          {/* Task Type */}
          <Form.Item
            name="type"
            label="脚本类型"
            rules={[{ required: true, message: '请选择脚本类型' }]}
            initialValue="shell"
          >
            <Select
              onChange={handleTypeChange}
              options={[
                { label: 'Shell', value: 'shell' },
                { label: 'Python', value: 'python' },
              ]}
            />
          </Form.Item>

          {/* Script Content */}
          <Form.Item
            name="content"
            label="脚本内容"
            rules={[{ required: true, message: '请输入脚本内容' }]}
          >
            <TextArea
              rows={8}
              placeholder="#!/usr/bin/env bash&#10;echo 'Hello World'"
              style={{
                fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                fontSize: 13,
                background: isDark ? '#141414' : '#fafafa',
              }}
            />
          </Form.Item>

          {/* Timeout */}
          <Form.Item name="timeout" label="超时时间 (秒)" initialValue={600}>
            <InputNumber min={10} max={86400} style={{ width: '100%' }} />
          </Form.Item>

          {/* Host Selection */}
          <Form.Item label="执行主机" required>
            <div style={{ marginBottom: 12, display: 'flex', gap: 12, flexWrap: 'wrap' }}>
              <Select
                placeholder="按项目筛选"
                allowClear
                style={{ width: 200 }}
                value={projectFilter}
                onChange={setProjectFilter}
                options={projects.map((p) => ({ label: p.name, value: p.id }))}
              />
              <Select
                placeholder="按环境筛选"
                allowClear
                style={{ width: 200 }}
                value={environmentFilter}
                onChange={setEnvironmentFilter}
                options={environments.map((e) => ({ label: e.name, value: e.id }))}
              />
              <Space>
                <Button size="small" onClick={handleSelectAll}>
                  全选
                </Button>
                <Button size="small" onClick={handleClearSelection}>
                  清空
                </Button>
              </Space>
            </div>
            <div style={{ marginBottom: 8, fontSize: 12, color: token.colorTextSecondary }}>
              已选择 {selectedAssetIds.length} 个主机
            </div>

            <div
              style={{
                border: `1px solid ${token.colorBorder}`,
                borderRadius: 6,
                padding: 12,
                maxHeight: 250,
                overflow: 'auto',
                background: token.colorFillTertiary,
              }}
            >
              {treeData.length === 0 ? (
                <Empty description="暂无可选主机" />
              ) : (
                <Tree
                  checkable
                  defaultExpandAll
                  treeData={treeData}
                  checkedKeys={checkedKeys}
                  onCheck={handleAssetCheck}
                />
              )}
            </div>
            <div style={{ marginTop: 8, color: token.colorTextTertiary }}>
              已选择 {selectedAssetIds.length} 台主机
            </div>
          </Form.Item>
        </Form>
      )}
    </Modal>
  )
}

export default TaskEditModal
