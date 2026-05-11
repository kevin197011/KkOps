// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface MetricPoint {
  t: number
  v: number
}

export interface MetricSeries {
  labels: Record<string, string>
  points: MetricPoint[]
}

export interface QueryResult {
  result_type: string
  series: MetricSeries[]
}

export interface MonitoringQueryRequest {
  integration_id: number
  query: string
  time?: string
  range?: {
    start: string
    end: string
    step?: string
  }
}

export const monitoringApi = {
  query: (body: MonitoringQueryRequest) =>
    apiClient.post<{ data: QueryResult }>('/monitoring/query', body),
}
