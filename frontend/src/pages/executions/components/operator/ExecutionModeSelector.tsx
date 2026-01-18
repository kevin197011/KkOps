// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { Card, Radio, Select, Typography, Space, Alert } from 'antd'
import { FileTextOutlined, EditOutlined } from '@ant-design/icons'
import { executionTemplateApi, ExecutionTemplate } from '@/api/execution'
import { useThemeStore } from '@/stores/theme'

const { Text } = Typography

export type ExecutionMode = 'template' | 'custom'

interface ExecutionModeSelectorProps {
  mode: ExecutionMode
  onModeChange: (mode: ExecutionMode) => void
  selectedTemplateId: number | null
  onTemplateSelect: (templateId: number | null, template: ExecutionTemplate | null) => void
}

const ExecutionModeSelector: React.FC<ExecutionModeSelectorProps> = ({
  mode,
  onModeChange,
  selectedTemplateId,
  onTemplateSelect,
}) => {
  const { mode: themeMode } = useThemeStore()
  const [templates, setTemplates] = useState<ExecutionTemplate[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (mode === 'template') {
      fetchTemplates()
    }
  }, [mode])

  const fetchTemplates = async () => {
    setLoading(true)
    try {
      const response = await executionTemplateApi.list()
      const templatesData = response.data?.data || response.data || []
      setTemplates(Array.isArray(templatesData) ? templatesData : [])
    } catch (error) {
      console.error('获取模板列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleModeChange = (e: any) => {
    const newMode = e.target.value as ExecutionMode
    onModeChange(newMode)
    if (newMode === 'custom') {
      onTemplateSelect(null, null)
    }
  }

  const handleTemplateChange = (value: number) => {
    const template = templates.find((t) => t.id === value)
    onTemplateSelect(value, template || null)
  }

  return (
    <Card
      style={{
        marginBottom: 16,
        borderRadius: 12,
      }}
      styles={{
        body: { padding: '20px 24px' },
      }}
    >
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <div>
          <Text strong style={{ marginBottom: 12, display: 'block' }}>
            执行模式
          </Text>
          <Radio.Group value={mode} onChange={handleModeChange} size="large">
            <Radio.Button value="template">
              <Space>
                <FileTextOutlined />
                模板模式
              </Space>
            </Radio.Button>
            <Radio.Button value="custom">
              <Space>
                <EditOutlined />
                自定义模式
              </Space>
            </Radio.Button>
          </Radio.Group>
        </div>

        {mode === 'template' && (
          <div>
            <Text strong style={{ marginBottom: 8, display: 'block' }}>
              选择模板
            </Text>
            <Select
              style={{ width: '100%' }}
              placeholder="请选择任务模板"
              value={selectedTemplateId}
              onChange={handleTemplateChange}
              loading={loading}
              showSearch
              optionFilterProp="children"
              filterOption={(input, option) =>
                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
              }
              options={templates.map((t) => ({
                value: t.id,
                label: t.name,
                description: t.description,
              }))}
            />
            {selectedTemplateId && (
              <Alert
                message={
                  templates.find((t) => t.id === selectedTemplateId)?.description ||
                  '已选择模板'
                }
                type="info"
                showIcon
                style={{ marginTop: 12 }}
              />
            )}
          </div>
        )}
      </Space>
    </Card>
  )
}

export default ExecutionModeSelector
