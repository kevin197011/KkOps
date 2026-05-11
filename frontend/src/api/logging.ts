// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface LogLine {
  timestamp: string
  message: string
  labels?: Record<string, string>
}

export interface LogSearchRequest {
  integration_id: number
  query: string
  start?: string
  end?: string
  limit?: number
}

export const loggingApi = {
  search: (body: LogSearchRequest) =>
    apiClient.post<{ data: LogLine[] }>('/logging/search', body),
}
