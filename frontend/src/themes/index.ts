// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { ThemeConfig, theme } from 'antd'

export const lightTheme: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  token: {
    // 主色调 - 信任蓝
    colorPrimary: '#2563EB',
    colorSuccess: '#10B981',
    colorWarning: '#F59E0B',
    colorError: '#EF4444',
    colorInfo: '#3B82F6',
    
    // 背景色
    colorBgBase: '#FFFFFF',
    colorBgContainer: '#F8FAFC',
    colorBgElevated: '#FFFFFF',
    
    // 文字颜色
    colorText: '#1E293B',
    colorTextSecondary: '#64748B',
    colorTextTertiary: '#94A3B8',
    
    // 边框颜色
    colorBorder: '#E2E8F0',
    colorBorderSecondary: '#F1F5F9',
    
    // 链接颜色
    colorLink: '#2563EB',
    colorLinkHover: '#1D4ED8',
    
    // 字体
    fontFamily: 'Open Sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    fontSizeHeading1: 32,
    fontSizeHeading2: 24,
    fontSizeHeading3: 20,
    fontSizeHeading4: 16,
    fontSize: 14,
    lineHeight: 1.6,
    
    // 圆角
    borderRadius: 6,
  },
  components: {
    Button: {
      primaryColor: '#FFFFFF',
      controlHeight: 36,
      borderRadius: 6,
    },
    Card: {
      borderRadius: 8,
      paddingLG: 24,
    },
    Table: {
      borderRadius: 8,
    },
  },
}

export const darkTheme: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  token: {
    // 主色调 - 稍亮的蓝色
    colorPrimary: '#3B82F6',
    colorSuccess: '#10B981',
    colorWarning: '#F59E0B',
    colorError: '#EF4444',
    colorInfo: '#60A5FA',
    
    // 背景色 - OLED友好的深色
    colorBgBase: '#0F172A',
    colorBgContainer: '#1E293B',
    colorBgElevated: '#334155',
    
    // 文字颜色
    colorText: '#F1F5F9',
    colorTextSecondary: '#CBD5E1',
    colorTextTertiary: '#94A3B8',
    
    // 边框颜色
    colorBorder: '#334155',
    colorBorderSecondary: '#475569',
    
    // 字体
    fontFamily: 'Open Sans, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    fontSizeHeading1: 32,
    fontSizeHeading2: 24,
    fontSizeHeading3: 20,
    fontSizeHeading4: 16,
    fontSize: 14,
    lineHeight: 1.6,
    
    // 圆角
    borderRadius: 6,
  },
  components: {
    Button: {
      primaryColor: '#FFFFFF',
      controlHeight: 36,
      borderRadius: 6,
    },
    Card: {
      borderRadius: 8,
      paddingLG: 24,
    },
    Table: {
      borderRadius: 8,
    },
  },
}
