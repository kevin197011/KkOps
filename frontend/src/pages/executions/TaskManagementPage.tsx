// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import { Layout, Button, Space, message, Input, Select, Spin } from 'antd'
import { PlusOutlined, ReloadOutlined, SearchOutlined } from '@ant-design/icons'
import { executionApi, Execution, ExecutionRecord } from '@/api/execution'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import WorkflowList from './components/WorkflowList'
import ExecutionDetailsPanel from './components/ExecutionDetailsPanel'
import TaskEditModal from './components/TaskEditModal'
import { useThemeStore } from '@/stores/theme'

const { Sider, Content } = Layout

export interface WorkflowWithExecutions extends Execution {
  executions: ExecutionRecord[]
  lastExecution?: ExecutionRecord
}

const TaskManagementPage = () => {
  const { mode } = useThemeStore()
  const [workflows, setWorkflows] = useState<WorkflowWithExecutions[]>([])
  const [selectedWorkflowId, setSelectedWorkflowId] = useState<number | null>(null)
  const [loading, setLoading] = useState(true)
  const [searchText, setSearchText] = useState('')
  const [statusFilter, setStatusFilter] = useState<string | null>(null)
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTask, setEditingTask] = useState<Execution | null>(null)

  // Fetch workflows and their executions
  const fetchWorkflows = useCallback(async () => {
    setLoading(true)
    try {
      const [tasksResponse] = await Promise.all([
        executionApi.list(1, 100),
      ])
      
      const tasks = tasksResponse.data.data || []
      
      // Fetch executions for each task
      const workflowsWithExecutions = await Promise.all(
        tasks.map(async (task) => {
          try {
            const execResponse = await executionApi.getHistory(task.id)
            const executions = execResponse.data.data || []
            // Sort by created_at descending to get newest first
            const sortedExecutions = [...executions].sort((a, b) => 
              new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
            )
            return {
              ...task,
              executions: sortedExecutions,
              lastExecution: sortedExecutions[0],
            }
          } catch {
            return { ...task, executions: [], lastExecution: undefined }
          }
        })
      )
      
      setWorkflows(workflowsWithExecutions)
      
      // Auto-select first workflow if none selected
      if (!selectedWorkflowId && workflowsWithExecutions.length > 0) {
        setSelectedWorkflowId(workflowsWithExecutions[0].id)
      }
    } catch (error) {
      message.error('获取工作流列表失败')
    } finally {
      setLoading(false)
    }
  }, [selectedWorkflowId])

  // Fetch projects and environments for filters
  const fetchFilters = useCallback(async () => {
    try {
      const [projectsRes, environmentsRes] = await Promise.all([
        projectApi.list(),
        environmentApi.list(),
      ])
      // API returns data directly, not wrapped in data.data
      setProjects(Array.isArray(projectsRes.data) ? projectsRes.data : [])
      setEnvironments(Array.isArray(environmentsRes.data) ? environmentsRes.data : [])
    } catch {
      // Silently fail for filters
    }
  }, [])

  useEffect(() => {
    fetchWorkflows()
    fetchFilters()
  }, [fetchWorkflows, fetchFilters])

  // Filter workflows
  const filteredWorkflows = workflows.filter((w) => {
    if (searchText && !w.name.toLowerCase().includes(searchText.toLowerCase())) {
      return false
    }
    if (statusFilter && w.status !== statusFilter) {
      return false
    }
    return true
  })

  const selectedWorkflow = workflows.find((w) => w.id === selectedWorkflowId)

  const handleCreateTask = () => {
    setEditingTask(null)
    setModalVisible(true)
  }

  const handleEditTask = (task: Execution) => {
    setEditingTask(task)
    setModalVisible(true)
  }

  const handleModalClose = () => {
    setModalVisible(false)
    setEditingTask(null)
  }

  const handleTaskSaved = () => {
    setModalVisible(false)
    setEditingTask(null)
    fetchWorkflows()
  }

  const handleExecute = async (taskId: number, executionType: 'sync' | 'async') => {
    try {
      await executionApi.execute(taskId, executionType)
      message.success('任务执行已启动')
      // Immediately refresh to show the new execution
      await fetchWorkflows()
      // For async execution, start polling to update status
      if (executionType === 'async') {
        let pollCount = 0
        const maxPolls = 150 // 5 minutes at 2 second intervals
        const pollInterval = setInterval(async () => {
          pollCount++
          try {
            const response = await executionApi.getHistory(taskId)
            const executions = response.data.data || []
            const hasRunning = executions.some((e) => e.status === 'running' || e.status === 'pending')
            await fetchWorkflows()
            // Stop if no running executions or max polls reached
            if (!hasRunning || pollCount >= maxPolls) {
              clearInterval(pollInterval)
            }
          } catch {
            clearInterval(pollInterval)
          }
        }, 2000)
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '执行失败')
    }
  }

  const handleCancel = async (taskId: number) => {
    try {
      await executionApi.cancel(taskId)
      message.success('任务已取消')
      fetchWorkflows()
    } catch (error: any) {
      message.error(error.response?.data?.error || '取消失败')
    }
  }

  const handleDelete = async (taskId: number) => {
    try {
      await executionApi.delete(taskId)
      message.success('任务已删除')
      // Clear selection if deleted task was selected
      if (selectedWorkflowId === taskId) {
        setSelectedWorkflowId(null)
      }
      fetchWorkflows()
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败')
    }
  }

  const isDark = mode === 'dark'

  return (
    <div style={{ height: 'calc(100vh - 120px)', display: 'flex', flexDirection: 'column' }}>
      {/* Header */}
      <div
        style={{
          padding: '16px 24px',
          borderBottom: isDark ? '1px solid #303030' : '1px solid #f0f0f0',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 12,
        }}
      >
        <h2 style={{ margin: 0, fontSize: 20 }}>任务管理</h2>
        <Space wrap>
          <Input
            placeholder="搜索工作流..."
            prefix={<SearchOutlined />}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 200 }}
            allowClear
          />
          <Select
            placeholder="状态筛选"
            style={{ width: 120 }}
            allowClear
            value={statusFilter}
            onChange={setStatusFilter}
            options={[
              { label: '待执行', value: 'pending' },
              { label: '执行中', value: 'running' },
              { label: '成功', value: 'success' },
              { label: '失败', value: 'failed' },
              { label: '已取消', value: 'cancelled' },
            ]}
          />
          <Button icon={<ReloadOutlined />} onClick={fetchWorkflows}>
            刷新
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateTask}>
            新建任务
          </Button>
        </Space>
      </div>

      {/* Main Content - Dual Column Layout */}
      <Layout style={{ flex: 1, background: 'transparent' }}>
        {/* Left Panel - Workflow List */}
        <Sider
          width={320}
          style={{
            background: isDark ? '#141414' : '#fff',
            borderRight: isDark ? '1px solid #303030' : '1px solid #f0f0f0',
            overflow: 'auto',
          }}
        >
          {loading ? (
            <div style={{ display: 'flex', justifyContent: 'center', padding: 40 }}>
              <Spin />
            </div>
          ) : (
            <WorkflowList
              workflows={filteredWorkflows}
              selectedWorkflowId={selectedWorkflowId}
              onSelect={setSelectedWorkflowId}
            />
          )}
        </Sider>

        {/* Right Panel - Execution Details */}
        <Content
          style={{
            background: isDark ? '#1f1f1f' : '#fafafa',
            overflow: 'auto',
            padding: 24,
          }}
        >
          {selectedWorkflow ? (
            <ExecutionDetailsPanel
              workflow={selectedWorkflow}
              onEdit={() => handleEditTask(selectedWorkflow)}
              onExecute={(type) => handleExecute(selectedWorkflow.id, type)}
              onCancel={() => handleCancel(selectedWorkflow.id)}
              onDelete={() => handleDelete(selectedWorkflow.id)}
              onRefresh={fetchWorkflows}
            />
          ) : (
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '100%',
                color: isDark ? '#666' : '#999',
              }}
            >
              {workflows.length === 0 ? (
                <div style={{ textAlign: 'center' }}>
                  <p>暂无工作流</p>
                  <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateTask}>
                    创建第一个任务
                  </Button>
                </div>
              ) : (
                '请从左侧选择一个工作流'
              )}
            </div>
          )}
        </Content>
      </Layout>

      {/* Task Edit Modal */}
      <TaskEditModal
        visible={modalVisible}
        task={editingTask}
        projects={projects}
        environments={environments}
        onClose={handleModalClose}
        onSave={handleTaskSaved}
      />
    </div>
  )
}

export default TaskManagementPage
