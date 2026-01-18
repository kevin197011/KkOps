// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect } from 'react'
import { Card, Select, Input, Space, Typography } from 'antd'
import { CodeOutlined } from '@ant-design/icons'

const { TextArea } = Input
const { Text } = Typography

// 不同类型的 shebang
const SHEBANGS: Record<string, string> = {
  shell: '#!/usr/bin/env bash\n\n',
  python: '#!/usr/bin/env python3\n\n',
}

interface ScriptEditorProps {
  mode: 'template' | 'custom'
  scriptContent: string
  scriptType: string
  onContentChange: (content: string) => void
  onTypeChange: (type: string) => void
  templateContent?: string
  readOnly?: boolean
}

const ScriptEditor: React.FC<ScriptEditorProps> = ({
  mode,
  scriptContent,
  scriptType,
  onContentChange,
  onTypeChange,
  templateContent,
  readOnly = false,
}) => {
  const isReadOnly = mode === 'template' || readOnly

  // 当选择模板时，自动填充内容
  useEffect(() => {
    if (mode === 'template' && templateContent) {
      onContentChange(templateContent)
    }
  }, [mode, templateContent, onContentChange])

  // 当脚本类型改变时，自动插入 shebang
  const handleTypeChange = (type: string) => {
    onTypeChange(type)
    if (mode === 'custom' && scriptContent) {
      // 移除旧的 shebang（如果存在）
      const withoutShebang = scriptContent.replace(/^#!.*\n\n?/, '')
      const newContent = SHEBANGS[type] + withoutShebang
      onContentChange(newContent)
    } else if (mode === 'custom' && !scriptContent) {
      // 如果内容为空，插入 shebang
      onContentChange(SHEBANGS[type])
    }
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
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Text strong>
            <CodeOutlined style={{ marginRight: 8 }} />
            脚本内容
          </Text>
          <Select
            value={scriptType}
            onChange={handleTypeChange}
            disabled={isReadOnly}
            style={{ width: 120 }}
            options={[
              { label: 'Shell', value: 'shell' },
              { label: 'Python', value: 'python' },
            ]}
          />
        </div>
        <TextArea
          value={scriptContent}
          onChange={(e) => onContentChange(e.target.value)}
          readOnly={isReadOnly}
          rows={12}
          placeholder={
            isReadOnly
              ? '请先选择模板'
              : scriptType === 'shell'
                ? '请输入 Shell 脚本内容...'
                : '请输入 Python 脚本内容...'
          }
          style={{
            fontFamily: '"JetBrains Mono", "SF Mono", "Consolas", monospace',
            fontSize: 13,
            lineHeight: 1.6,
          }}
        />
        {isReadOnly && mode === 'template' && (
          <Text type="secondary" style={{ fontSize: 12 }}>
            模板模式下脚本内容为只读，如需修改请切换到自定义模式
          </Text>
        )}
      </Space>
    </Card>
  )
}

export default ScriptEditor
