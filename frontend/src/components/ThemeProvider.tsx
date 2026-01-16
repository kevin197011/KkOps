// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React from 'react'
import { ConfigProvider, App } from 'antd'
import { lightTheme, darkTheme } from '@/themes'
import { useThemeStore } from '@/stores/theme'

interface ThemeProviderProps {
  children: React.ReactNode
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const mode = useThemeStore((state) => state.mode)
  const theme = mode === 'dark' ? darkTheme : lightTheme

  return (
    <ConfigProvider theme={theme}>
      <App>
        {children}
      </App>
    </ConfigProvider>
  )
}
