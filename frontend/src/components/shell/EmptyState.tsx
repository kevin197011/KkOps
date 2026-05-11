// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { Empty } from 'antd'

interface EmptyStateProps {
  description: React.ReactNode
  image?: React.ComponentProps<typeof Empty>['image']
  children?: React.ReactNode
}

export const EmptyState: React.FC<EmptyStateProps> = ({ description, image, children }) => (
  <Empty image={image ?? Empty.PRESENTED_IMAGE_SIMPLE} description={description}>
    {children}
  </Empty>
)
