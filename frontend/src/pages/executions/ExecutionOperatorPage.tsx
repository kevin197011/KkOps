// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback, useRef } from 'react'
import { Layout, Button, message, Form, Modal, Input, Space, Typography } from 'antd'
import { ThunderboltOutlined, LinkOutlined } from '@ant-design/icons'
import { executionApi, ExecutionRecord, executionTemplateApi, ExecutionTemplate } from '@/api/execution'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import { assetApi, Asset } from '@/api/asset'
import ExecutionModeSelector, { ExecutionMode } from './components/operator/ExecutionModeSelector'
import ScriptEditor from './components/operator/ScriptEditor'
import HostSelector from './components/operator/HostSelector'
import ExecutionOptions from './components/operator/ExecutionOptions'
import ExecutionResults from './components/operator/ExecutionResults'
import { useThemeStore } from '@/stores/theme'
import { useNavigate } from 'react-router-dom'

const { Content } = Layout
const { Text } = Typography

const ExecutionOperatorPage = () => {
  const { mode } = useThemeStore()
  const navigate = useNavigate()
  const [form] = Form.useForm()

  // 执行模式
  const [executionMode, setExecutionMode] = useState<ExecutionMode>('template')
  const [selectedTemplateId, setSelectedTemplateId] = useState<number | null>(null)
  const [selectedTemplate, setSelectedTemplate] = useState<ExecutionTemplate | null>(null)

  // 脚本内容
  const [scriptContent, setScriptContent] = useState('')
  const [scriptType, setScriptType] = useState<string>('shell')

  // 主机选择
  const [selectedAssetIds, setSelectedAssetIds] = useState<number[]>([])
  const [assets, setAssets] = useState<Map<number, Asset>>(new Map())

  // 执行选项
  const [executionType, setExecutionType] = useState<'sync' | 'async'>('async')
  const [timeout, setTimeout] = useState<number>(600)
  const [saveAsTask, setSaveAsTask] = useState<boolean>(false)

  // 项目和环境
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])

  // 执行结果
  const [executionRecords, setExecutionRecords] = useState<ExecutionRecord[]>([])
  const [isExecuting, setIsExecuting] = useState(false)
  const [currentTaskId, setCurrentTaskId] = useState<number | null>(null)
  const [taskName, setTaskName] = useState<string>('')
  const [saveTaskModalVisible, setSaveTaskModalVisible] = useState(false)
  const [saveTaskForm] = Form.useForm()

  // 轮询定时器
  const pollingIntervalRef = useRef<NodeJS.Timeout | null>(null)

  // 获取项目和环境列表
  useEffect(() => {
    fetchProjectsAndEnvironments()
  }, [])

  const fetchProjectsAndEnvironments = async () => {
    try {
      const [projectsRes, environmentsRes] = await Promise.all([
        projectApi.list(),
        environmentApi.list(),
      ])
      setProjects(Array.isArray(projectsRes.data) ? projectsRes.data : [])
      setEnvironments(Array.isArray(environmentsRes.data) ? environmentsRes.data : [])
    } catch (error) {
      console.error('获取项目和环境列表失败:', error)
    }
  }

  // 获取资产列表
  useEffect(() => {
    fetchAssets()
  }, [])

  const fetchAssets = async () => {
    try {
      const response = await assetApi.list({ page: 1, page_size: 1000, status: 'active' })
      const assetsData = response.data?.data || []
      const assetMap = new Map<number, Asset>()
      assetsData.forEach((asset: Asset) => {
        assetMap.set(asset.id, asset)
      })
      setAssets(assetMap)
    } catch (error) {
      console.error('获取资产列表失败:', error)
    }
  }

  // 模板选择处理
  const handleTemplateSelect = useCallback((templateId: number | null, template: ExecutionTemplate | null) => {
    setSelectedTemplateId(templateId)
    setSelectedTemplate(template)
    if (template) {
      setScriptContent(template.content || '')
      setScriptType(template.type || 'shell')
    }
  }, [])

  // 执行模式改变处理
  const handleModeChange = useCallback((mode: ExecutionMode) => {
    setExecutionMode(mode)
    if (mode === 'custom') {
      setSelectedTemplateId(null)
      setSelectedTemplate(null)
      // 插入默认 shebang
      if (!scriptContent) {
        setScriptContent(scriptType === 'shell' ? '#!/bin/bash\n\n' : '#!/usr/bin/env python3\n\n')
      }
    }
  }, [scriptContent, scriptType])

  // 验证表单
  const validateForm = () => {
    if (!scriptContent.trim()) {
      message.error('请选择模板或输入脚本内容')
      return false
    }
    if (selectedAssetIds.length === 0) {
      message.error('请至少选择一个执行主机')
      return false
    }
    return true
  }

  // 执行任务
  const handleExecute = async () => {
    if (!validateForm()) {
      return
    }

    try {
      setIsExecuting(true)
      setExecutionRecords([])

      // 创建任务（如果需要保存）
      let taskId: number | null = null
      if (saveAsTask) {
        const taskName = selectedTemplate?.name || `执行任务 ${new Date().toLocaleString()}`
        const taskResponse = await executionApi.create({
          name: taskName,
          description: selectedTemplate?.description || '',
          content: scriptContent,
          type: scriptType,
          asset_ids: selectedAssetIds,
          template_id: selectedTemplateId || undefined,
        })
        taskId = taskResponse.data.id
        setCurrentTaskId(taskId)
      } else {
        // 临时执行：也创建任务，但不保存（后续可添加 is_temporary 标记）
        const tempTaskName = `临时执行 ${new Date().toLocaleString()}`
        const taskResponse = await executionApi.create({
          name: tempTaskName,
          description: '临时执行任务',
          content: scriptContent,
          type: scriptType,
          asset_ids: selectedAssetIds,
        })
        taskId = taskResponse.data.id
        setCurrentTaskId(taskId)
      }

      // 执行任务
      await executionApi.execute(taskId, executionType)
      message.success('任务执行已启动')

      // 获取执行记录
      if (taskId) {
        const { hasRunning } = await fetchExecutionRecords(taskId)

        // 如果保存为任务，显示保存任务名称的提示
        if (saveAsTask) {
          setTaskName(selectedTemplate?.name || `执行任务 ${new Date().toLocaleString()}`)
          message.info('任务已保存，可在"任务执行"页面查看')
        }

        // 异步执行时轮询状态
        if (executionType === 'async' && hasRunning) {
          startPolling(taskId)
        } else {
          setIsExecuting(false)
        }
      } else {
        setIsExecuting(false)
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '执行失败')
      setIsExecuting(false)
    }
  }

  // 获取执行记录
  const fetchExecutionRecords = useCallback(async (taskId: number) => {
    try {
      const response = await executionApi.getHistory(taskId)
      const records = response.data.data || []
      setExecutionRecords(records)

      // 检查是否还在执行
      const hasRunning = records.some((r: ExecutionRecord) => r.status === 'running' || r.status === 'pending')
      setIsExecuting(hasRunning)
      return { records, hasRunning }
    } catch (error) {
      console.error('获取执行记录失败:', error)
      return { records: [], hasRunning: false }
    }
  }, [])

  // 开始轮询
  const startPolling = useCallback((taskId: number) => {
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current)
    }

    let pollCount = 0
    const maxPolls = 150 // 5 分钟，每 2 秒一次

    pollingIntervalRef.current = setInterval(async () => {
      pollCount++
      try {
        const { hasRunning } = await fetchExecutionRecords(taskId)
        if (!hasRunning || pollCount >= maxPolls) {
          stopPolling()
        }
      } catch (error) {
        console.error('轮询执行记录失败:', error)
        stopPolling()
      }
    }, 2000)
  }, [fetchExecutionRecords])

  // 停止轮询
  const stopPolling = useCallback(() => {
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current)
      pollingIntervalRef.current = null
    }
    setIsExecuting(false)
  }, [])

  // 清理定时器
  useEffect(() => {
    return () => {
      stopPolling()
    }
  }, [stopPolling])

  // 刷新执行记录
  const handleRefresh = useCallback(() => {
    if (currentTaskId) {
      fetchExecutionRecords(currentTaskId)
    }
  }, [currentTaskId, fetchExecutionRecords])

  // 查看任务历史
  const handleViewTaskHistory = () => {
    if (currentTaskId) {
      navigate(`/tasks/${currentTaskId}/history`)
    }
  }

  const isDark = mode === 'dark'

  return (
    <Layout style={{ minHeight: 'calc(100vh - 120px)', background: 'transparent' }}>
      <Content
        style={{
          padding: 24,
          maxWidth: 1200,
          margin: '0 auto',
          width: '100%',
        }}
      >
        <Space direction="vertical" size={24} style={{ width: '100%' }}>
          {/* 页面标题 */}
          <div>
            <Text strong style={{ fontSize: 24, display: 'block', marginBottom: 8 }}>
              <ThunderboltOutlined style={{ marginRight: 8 }} />
              运维执行
            </Text>
            <Text type="secondary">
              选择模板或输入自定义脚本，快速执行到目标主机
            </Text>
          </div>

          {/* 执行模式选择 */}
          <ExecutionModeSelector
            mode={executionMode}
            onModeChange={handleModeChange}
            selectedTemplateId={selectedTemplateId}
            onTemplateSelect={handleTemplateSelect}
          />

          {/* 脚本编辑器 */}
          <ScriptEditor
            mode={executionMode}
            scriptContent={scriptContent}
            scriptType={scriptType}
            onContentChange={setScriptContent}
            onTypeChange={setScriptType}
            templateContent={selectedTemplate?.content}
            readOnly={executionMode === 'template'}
          />

          {/* 主机选择器 */}
          <HostSelector
            selectedAssetIds={selectedAssetIds}
            onSelectionChange={setSelectedAssetIds}
            projects={projects}
            environments={environments}
          />

          {/* 执行选项 */}
          <ExecutionOptions
            executionType={executionType}
            onExecutionTypeChange={setExecutionType}
            timeout={timeout}
            onTimeoutChange={setTimeout}
            saveAsTask={saveAsTask}
            onSaveAsTaskChange={setSaveAsTask}
          />

          {/* 执行按钮 */}
          <div style={{ textAlign: 'center', padding: '16px 0' }}>
            <Button
              type="primary"
              size="large"
              icon={<ThunderboltOutlined />}
              onClick={handleExecute}
              loading={isExecuting}
              disabled={isExecuting}
            >
              {isExecuting ? '执行中...' : '执行'}
            </Button>
            {currentTaskId && saveAsTask && (
              <Button
                type="link"
                icon={<LinkOutlined />}
                onClick={handleViewTaskHistory}
                style={{ marginLeft: 16 }}
              >
                查看任务历史
              </Button>
            )}
          </div>

          {/* 执行结果 */}
          {executionRecords.length > 0 && (
            <ExecutionResults
              executionRecords={executionRecords}
              assets={assets}
              isRunning={isExecuting}
              onRefresh={handleRefresh}
            />
          )}
        </Space>
      </Content>
    </Layout>
  )
}

export default ExecutionOperatorPage
