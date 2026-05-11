// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { Typography, Space } from 'antd'

const { Title, Text } = Typography

interface PageHeaderProps {
  title: React.ReactNode
  description?: React.ReactNode
  extra?: React.ReactNode
}

export const PageHeader: React.FC<PageHeaderProps> = ({ title, description, extra }) => (
  <div
    style={{
      display: 'flex',
      flexWrap: 'wrap',
      alignItems: 'flex-start',
      justifyContent: 'space-between',
      gap: 16,
      marginBottom: 24,
    }}
  >
    <div style={{ flex: '1 1 240px', minWidth: 0 }}>
      <Title level={4} style={{ marginBottom: description ? 8 : 0 }}>
        {title}
      </Title>
      {description ? (
        <Text type="secondary" style={{ display: 'block' }}>
          {description}
        </Text>
      ) : null}
    </div>
    {extra ? (
      <Space wrap>{extra}</Space>
    ) : null}
  </div>
)
