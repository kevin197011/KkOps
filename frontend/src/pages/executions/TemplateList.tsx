// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, Tag, Upload, Descriptions } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined, UploadOutlined } from '@ant-design/icons'
import { executionTemplateApi, ExecutionTemplate, CreateTemplateRequest, UpdateTemplateRequest, ExportTemplatesConfig, ImportTemplatesConfig, ImportResult } from '@/api/execution'
import type { UploadFile } from 'antd/es/upload/interface'
import { usePermissionStore } from '@/stores/permission'

const { TextArea } = Input

// 不同类型的 shebang
const SHEBANGS: Record<string, string> = {
  shell: '#!/usr/bin/env bash\n\n',
  python: '#!/usr/bin/env python3\n\n',
}

const TemplateList = () => {
  const { hasPermission } = usePermissionStore()
  const [templates, setTemplates] = useState<ExecutionTemplate[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTemplate, setEditingTemplate] = useState<ExecutionTemplate | null>(null)
  const [importResultVisible, setImportResultVisible] = useState(false)
  const [importResult, setImportResult] = useState<ImportResult | null>(null)
  const [form] = Form.useForm()

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


  const fetchTemplates = async () => {
    setLoading(true)
    try {
      const response = await executionTemplateApi.list(page, pageSize)
      // 处理不同的响应格式
      if (Array.isArray(response.data)) {
        setTemplates(response.data)
        setTotal(response.data.length)
      } else if (response.data?.data) {
        setTemplates(response.data.data)
        setTotal(response.data.total || response.data.data.length)
      } else {
        setTemplates([])
        setTotal(0)
      }
    } catch (error: any) {
      message.error('获取模板列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTemplates()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize])

  const handleCreate = () => {
    setEditingTemplate(null)
    form.resetFields()
    // 新建时默认插入 shell 的 shebang
    form.setFieldsValue({ content: SHEBANGS.shell })
    setModalVisible(true)
  }

  const handleEdit = (template: ExecutionTemplate) => {
    setEditingTemplate(template)
    form.setFieldsValue(template)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个模板吗？',
      onOk: async () => {
        try {
          await executionTemplateApi.delete(id)
          message.success('删除成功')
          fetchTemplates()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateTemplateRequest | UpdateTemplateRequest) => {
    try {
      if (editingTemplate) {
        await executionTemplateApi.update(editingTemplate.id, values)
        message.success('更新成功')
      } else {
        await executionTemplateApi.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      // 如果是第一页或数据量少，直接刷新；否则回到第一页查看新数据
      if (page === 1) {
        fetchTemplates()
      } else {
        setPage(1)
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 导出配置
  const handleExport = async () => {
    try {
      const config = await executionTemplateApi.exportConfig()
      const blob = new Blob([JSON.stringify(config.data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `task-templates-${new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (error: any) {
      message.error('导出失败: ' + (error.response?.data?.error || error.message))
    }
  }

  // 导入配置
  const handleImport = async (file: File) => {
    try {
      const text = await file.text()
      const config: ImportTemplatesConfig = JSON.parse(text)
      
      // 验证配置格式
      if (!config.templates || !Array.isArray(config.templates)) {
        message.error('无效的导入配置格式')
        return false
      }

      // 执行导入
      const result = await executionTemplateApi.importConfig(config)
      setImportResult(result.data)
      setImportResultVisible(true)
      
      // 如果导入成功，刷新列表
      if (result.data.success > 0) {
        fetchTemplates()
      }
      
      return false // 阻止默认上传行为
    } catch (error: any) {
      message.error('导入失败: ' + (error.response?.data?.error || error.message))
      return false
    }
  }

  const uploadProps = {
    beforeUpload: (file: File) => {
      handleImport(file)
      return false // 阻止默认上传行为
    },
    accept: '.json',
    showUploadList: false,
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '模板名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={type === 'shell' ? 'blue' : type === 'python' ? 'green' : 'default'}>
          {type || 'shell'}
        </Tag>
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
      width: 150,
      render: (_: any, record: ExecutionTemplate) => {
        const actions = []
        if (hasPermission('templates', 'update')) {
          actions.push(
            <Button key="edit" type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
              编辑
            </Button>
          )
        }
        if (hasPermission('templates', 'delete')) {
          actions.push(
            <Button key="delete" type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
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
      <div style={{ 
        marginBottom: 16, 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: 8
      }}>
        <h2>任务模板管理</h2>
        <Space>
          {hasPermission('templates', 'export') && (
            <Button icon={<DownloadOutlined />} onClick={handleExport}>
              导出配置
            </Button>
          )}
          {hasPermission('templates', 'import') && (
            <Upload {...uploadProps}>
              <Button icon={<UploadOutlined />}>
                导入配置
              </Button>
            </Upload>
          )}
          {hasPermission('templates', 'create') && (
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增模板">
              新增模板
            </Button>
          )}
        </Space>
      </div>
      <Table
        columns={columns}
        dataSource={templates}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
          responsive: true,
          onChange: (newPage, newPageSize) => {
            setPage(newPage)
            setPageSize(newPageSize || 20)
          },
          onShowSizeChange: (current, size) => {
            setPage(1)
            setPageSize(size)
          },
        }}
      />
      <Modal
        title={editingTemplate ? '编辑模板' : '新增模板'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 800 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="模板名称"
            rules={[{ required: true, message: '请输入模板名称' }]}
          >
            <Input placeholder="模板名称" />
          </Form.Item>
          <Form.Item name="type" label="类型" initialValue="shell">
            <Select onChange={handleTypeChange}>
              <Select.Option value="shell">Shell</Select.Option>
              <Select.Option value="python">Python</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="content"
            label="模板内容"
            rules={[{ required: true, message: '请输入模板内容' }]}
          >
            <TextArea rows={10} placeholder="输入脚本内容..." />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} placeholder="模板描述" />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="导入结果"
        open={importResultVisible}
        onCancel={() => setImportResultVisible(false)}
        footer={[
          <Button key="close" onClick={() => setImportResultVisible(false)}>
            关闭
          </Button>,
        ]}
      >
        {importResult && (
          <Descriptions bordered column={1}>
            <Descriptions.Item label="总数">{importResult.total}</Descriptions.Item>
            <Descriptions.Item label="成功">
              <Tag color="success">{importResult.success}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="失败">
              <Tag color="error">{importResult.failed}</Tag>
            </Descriptions.Item>
            {importResult.errors && importResult.errors.length > 0 && (
              <Descriptions.Item label="错误信息">
                <ul style={{ margin: 0, paddingLeft: 20 }}>
                  {importResult.errors.map((error, index) => (
                    <li key={index} style={{ color: 'red' }}>{error}</li>
                  ))}
                </ul>
              </Descriptions.Item>
            )}
            {importResult.skipped && importResult.skipped.length > 0 && (
              <Descriptions.Item label="跳过项">
                <ul style={{ margin: 0, paddingLeft: 20 }}>
                  {importResult.skipped.map((item, index) => (
                    <li key={index} style={{ color: 'orange' }}>{item}</li>
                  ))}
                </ul>
              </Descriptions.Item>
            )}
          </Descriptions>
        )}
      </Modal>
    </div>
  )
}

export default TemplateList
