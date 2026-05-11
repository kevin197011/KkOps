// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import type { MenuProps } from 'antd'
import type { ReactNode } from 'react'
import { isValidElement } from 'react'

/** Route command for command palette (same entries as filtered sidebar). */
export type NavCommand = { key: string; label: string }

function labelToString(label: ReactNode): string {
  if (typeof label === 'string' || typeof label === 'number') return String(label)
  if (!isValidElement(label)) return ''
  const props = label.props as { children?: ReactNode }
  if (typeof props.children === 'string') return props.children
  if (typeof props.children === 'number') return String(props.children)
  return ''
}

/**
 * Flattens Ant Design menu items to route commands (keys starting with '/').
 */
export function flattenMenuToCommands(items: MenuProps['items']): NavCommand[] {
  const out: NavCommand[] = []

  const walk = (nodes: MenuProps['items']) => {
    if (!nodes) return
    for (const node of nodes) {
      if (!node || node.type === 'divider') continue
      if ('children' in node && node.children && Array.isArray(node.children)) {
        walk(node.children)
      }
      if ('key' in node && typeof node.key === 'string' && node.key.startsWith('/')) {
        const raw = 'label' in node ? node.label : undefined
        const label = labelToString(raw as ReactNode) || node.key
        out.push({ key: node.key, label })
      }
    }
  }

  walk(items)
  return out
}
