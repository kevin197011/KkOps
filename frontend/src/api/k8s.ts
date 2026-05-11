// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

export interface K8sNamespace {
  name: string
  status: string
}

export interface K8sWorkload {
  kind: string
  name: string
  replicas: number
  available: number
  image: string
}

export interface K8sPod {
  name: string
  status: string
  node: string
  restarts: number
  age: string
}

export interface K8sEvent {
  type: string
  reason: string
  message: string
  object: string
  namespace: string
  last: string
}

export interface K8sNode {
  name: string
  status: string
  kubelet_version: string
  addresses: string[]
}

export const k8sApi = {
  namespaces: (clusterId: number) =>
    apiClient.get<{ data: K8sNamespace[] }>(`/k8s/clusters/${clusterId}/namespaces`),
  workloads: (clusterId: number, namespace?: string) =>
    apiClient.get<{ data: K8sWorkload[] }>(`/k8s/clusters/${clusterId}/workloads`, {
      params: namespace ? { namespace } : {},
    }),
  pods: (clusterId: number, namespace?: string) =>
    apiClient.get<{ data: K8sPod[] }>(`/k8s/clusters/${clusterId}/pods`, {
      params: namespace ? { namespace } : {},
    }),
  podLogs: (clusterId: number, pod: string, params?: { namespace?: string; container?: string; tail?: number }) =>
    apiClient.get<string>(`/k8s/clusters/${clusterId}/pods/${encodeURIComponent(pod)}/logs`, {
      params,
      responseType: 'text',
    }),
  events: (clusterId: number, namespace?: string) =>
    apiClient.get<{ data: K8sEvent[] }>(`/k8s/clusters/${clusterId}/events`, {
      params: namespace ? { namespace } : {},
    }),
  nodes: (clusterId: number) =>
    apiClient.get<{ data: K8sNode[] }>(`/k8s/clusters/${clusterId}/nodes`),
}
