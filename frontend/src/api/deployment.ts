// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import apiClient from './client'

// 模板信息
export interface TemplateInfo {
  id: number
  name: string
  type: string
}

// 部署模块
export interface DeploymentModule {
  id: number
  project_id: number
  project_name: string
  environment_id: number | null
  environment_name: string
  template_id: number | null
  template: TemplateInfo | null  // 关联的执行模板
  name: string
  description: string
  version_source_url: string
  deploy_script: string
  script_type: string
  timeout: number
  asset_ids: number[]
  created_by: number
  created_at: string
  updated_at: string
}

// 创建模块请求
export interface CreateModuleRequest {
  project_id: number
  environment_id?: number | null
  template_id?: number | null  // 可选：关联执行模板
  name: string
  description?: string
  version_source_url?: string
  deploy_script?: string
  script_type?: string
  timeout?: number
  asset_ids?: number[]
}

// 更新模块请求
export interface UpdateModuleRequest {
  project_id?: number
  environment_id?: number | null
  template_id?: number | null  // 可选：关联执行模板
  name?: string
  description?: string
  version_source_url?: string
  deploy_script?: string
  script_type?: string
  timeout?: number
  asset_ids?: number[]
}

// 版本响应
export interface VersionSourceResponse {
  versions: string[]
  latest: string
}

// 部署请求
export interface DeployRequest {
  version: string
  asset_ids: number[]
}

// 部署记录
export interface Deployment {
  id: number
  module_id: number
  module_name: string
  project_name: string
  version: string
  status: string
  asset_ids: number[]
  output: string
  error: string
  created_by: number
  creator_name: string
  started_at: string | null
  finished_at: string | null
  created_at: string
}

// 部署历史列表响应
export interface DeploymentListResponse {
  data: Deployment[]
  total: number
  page: number
  size: number
}

// 导出配置类型
export interface ExportModuleConfig {
  name: string
  description: string
  project_name: string
  environment_name?: string
  template_name?: string
  version_source_url?: string
  deploy_script: string
  script_type: string
  timeout: number
  target_hosts?: string[] // 目标主机列表
}

export interface ExportConfig {
  version: string
  export_at: string
  modules: ExportModuleConfig[]
}

// 导入结果类型
export interface ImportResult {
  total: number
  success: number
  failed: number
  errors?: string[]
  skipped?: string[]
}

// 部署模块 API
export const deploymentModuleApi = {
  list: (projectId?: number) => {
    const params = projectId ? `?project_id=${projectId}` : ''
    return apiClient.get<DeploymentModule[]>(`/deployment-modules${params}`)
  },
  get: (id: number) => apiClient.get<DeploymentModule>(`/deployment-modules/${id}`),
  create: (data: CreateModuleRequest) =>
    apiClient.post<DeploymentModule>('/deployment-modules', data),
  update: (id: number, data: UpdateModuleRequest) =>
    apiClient.put<DeploymentModule>(`/deployment-modules/${id}`, data),
  delete: (id: number) => apiClient.delete(`/deployment-modules/${id}`),
  getVersions: (id: number) =>
    apiClient.get<VersionSourceResponse>(`/deployment-modules/${id}/versions`),
  deploy: (id: number, data: DeployRequest) =>
    apiClient.post<Deployment>(`/deployment-modules/${id}/deploy`, data),
  // 导出配置
  exportConfig: () => apiClient.get<ExportConfig>('/deployment-modules/export'),
  // 导入配置
  importConfig: (data: ExportConfig) =>
    apiClient.post<ImportResult>('/deployment-modules/import', data),
  // 预览导入
  previewImport: (data: ExportConfig) =>
    apiClient.post<ImportResult>('/deployment-modules/import/preview', data),
}

// 部署历史 API
export const deploymentApi = {
  list: (moduleId?: number, page?: number, pageSize?: number) => {
    const params = new URLSearchParams()
    if (moduleId) params.append('module_id', String(moduleId))
    if (page) params.append('page', String(page))
    if (pageSize) params.append('page_size', String(pageSize))
    const query = params.toString()
    return apiClient.get<DeploymentListResponse>(`/deployments${query ? `?${query}` : ''}`)
  },
  get: (id: number) => apiClient.get<Deployment>(`/deployments/${id}`),
  cancel: (id: number) => apiClient.post(`/deployments/${id}/cancel`),
}
