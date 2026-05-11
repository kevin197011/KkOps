// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'

interface PageContainerProps {
  children: React.ReactNode
  /** Dense layout for data-heavy tables */
  dense?: boolean
}

/**
 * Standard page wrapper: max-width, horizontal padding, optional density.
 */
export const PageContainer: React.FC<PageContainerProps> = ({ children, dense }) => (
  <div
    style={{
      maxWidth: 1400,
      margin: '0 auto',
      padding: dense ? '0 8px 16px' : '0 0 24px',
      width: '100%',
    }}
  >
    {children}
  </div>
)
