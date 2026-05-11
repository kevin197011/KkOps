// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { Typography } from 'antd'

const { Text } = Typography

interface SectionProps {
  title?: React.ReactNode
  extra?: React.ReactNode
  children: React.ReactNode
}

export const Section: React.FC<SectionProps> = ({ title, extra, children }) => (
  <section style={{ marginBottom: 24 }}>
    {(title || extra) && (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          marginBottom: 12,
          gap: 12,
        }}
      >
        {title ? (
          <Text strong style={{ fontSize: 14 }}>
            {title}
          </Text>
        ) : (
          <span />
        )}
        {extra}
      </div>
    )}
    {children}
  </section>
)
