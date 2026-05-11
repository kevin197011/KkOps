// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useEffect, useMemo, useState } from 'react'
import { Modal, Input, List } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import type { NavCommand } from '@/navigation/flattenMenu'

interface CommandPaletteProps {
  open: boolean
  onClose: () => void
  commands: NavCommand[]
  onNavigate: (path: string) => void
}

export const CommandPalette: React.FC<CommandPaletteProps> = ({
  open,
  onClose,
  commands,
  onNavigate,
}) => {
  const [query, setQuery] = useState('')

  useEffect(() => {
    if (!open) setQuery('')
  }, [open])

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return commands
    return commands.filter(
      (c) =>
        c.label.toLowerCase().includes(q) ||
        c.key.toLowerCase().includes(q)
    )
  }, [commands, query])

  const handleSelect = (path: string) => {
    onNavigate(path)
    onClose()
  }

  return (
    <Modal
      title={null}
      footer={null}
      closable
      open={open}
      onCancel={onClose}
      destroyOnClose
      styles={{ body: { paddingTop: 8 } }}
      width={480}
    >
      <Input
        autoFocus
        allowClear
        prefix={<SearchOutlined />}
        placeholder="Jump to page…"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onPressEnter={() => {
          if (filtered.length === 1) handleSelect(filtered[0].key)
        }}
      />
      <List
        style={{ marginTop: 12, maxHeight: 360, overflow: 'auto' }}
        size="small"
        dataSource={filtered}
        locale={{ emptyText: 'No matching pages' }}
        renderItem={(item) => (
          <List.Item
            style={{ cursor: 'pointer', padding: '8px 12px' }}
            onClick={() => handleSelect(item.key)}
          >
            <div>
              <div style={{ fontWeight: 500 }}>{item.label}</div>
              <div style={{ fontSize: 12, opacity: 0.65 }}>{item.key}</div>
            </div>
          </List.Item>
        )}
      />
      <div style={{ fontSize: 12, opacity: 0.55, marginTop: 8 }}>
        Tip: Enter when one match remains
      </div>
    </Modal>
  )
}
