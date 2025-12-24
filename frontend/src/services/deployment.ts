import api from '../config/api';

export interface DeploymentConfig {
  id: number;
  name: string;
  application_name: string;
  description?: string;
  salt_state_files?: string;
  target_groups?: string;
  target_hosts?: string;
  environment?: string;
  config_data?: string;
  created_by: number;
  created_at: string;
  updated_at: string;
}

export interface Deployment {
  id: number;
  config_id: number;
  version: string;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  started_by: number;
  started_at: string;
  completed_at?: string;
  duration_seconds?: number;
  target_hosts: string; // JSON array string
  salt_job_id?: string;
  results?: string; // JSON string
  error_message?: string;
  is_rollback: boolean;
  rollback_from_deployment_id?: number;
  created_at: string;
  updated_at: string;
  config?: DeploymentConfig;
}

export interface DeploymentVersion {
  id: number;
  application_name: string;
  version: string;
  release_notes?: string;
  created_by: number;
  created_at: string;
}

export interface ListDeploymentConfigsResponse {
  configs: DeploymentConfig[];
  total: number;
  page: number;
  page_size: number;
}

export interface ListDeploymentsResponse {
  deployments: Deployment[];
  total: number;
  page: number;
  page_size: number;
}

export interface ListDeploymentVersionsResponse {
  versions: DeploymentVersion[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreateDeploymentConfigRequest {
  name: string;
  application_name: string;
  description?: string;
  salt_state_files?: string;
  target_groups?: string;
  target_hosts?: string;
  environment?: string;
  config_data?: string;
}

export interface CreateDeploymentRequest {
  config_id: number;
  version: string;
  target_hosts: number[];
}

export interface CreateDeploymentVersionRequest {
  application_name: string;
  version: string;
  release_notes?: string;
}

export const deploymentService = {
  // 部署配置
  listConfigs: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: { application_name?: string; name?: string; environment?: string }
  ): Promise<ListDeploymentConfigsResponse> => {
    const response = await api.get('/deployment-configs', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  getConfig: async (id: number): Promise<{ config: DeploymentConfig }> => {
    const response = await api.get(`/deployment-configs/${id}`);
    return response.data;
  },

  createConfig: async (data: CreateDeploymentConfigRequest): Promise<{ config: DeploymentConfig }> => {
    const response = await api.post('/deployment-configs', data);
    return response.data;
  },

  updateConfig: async (id: number, data: Partial<CreateDeploymentConfigRequest>): Promise<void> => {
    await api.put(`/deployment-configs/${id}`, data);
  },

  deleteConfig: async (id: number): Promise<void> => {
    await api.delete(`/deployment-configs/${id}`);
  },

  // 部署执行
  listDeployments: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: { config_id?: number; status?: string; version?: string }
  ): Promise<ListDeploymentsResponse> => {
    const response = await api.get('/deployments', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  getDeployment: async (id: number): Promise<{ deployment: Deployment }> => {
    const response = await api.get(`/deployments/${id}`);
    return response.data;
  },

  getDeploymentStatus: async (id: number): Promise<{ deployment: Deployment }> => {
    const response = await api.get(`/deployments/${id}/status`);
    return response.data;
  },

  createDeployment: async (data: CreateDeploymentRequest): Promise<{ deployment: Deployment }> => {
    const response = await api.post('/deployments', data);
    return response.data;
  },

  rollbackDeployment: async (id: number): Promise<{ deployment: Deployment }> => {
    const response = await api.post(`/deployments/${id}/rollback`);
    return response.data;
  },

  // 部署版本
  listVersions: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: { application_name?: string }
  ): Promise<ListDeploymentVersionsResponse> => {
    const response = await api.get('/deployment-versions', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  getVersion: async (id: number): Promise<{ version: DeploymentVersion }> => {
    const response = await api.get(`/deployment-versions/${id}`);
    return response.data;
  },

  createVersion: async (data: CreateDeploymentVersionRequest): Promise<{ version: DeploymentVersion }> => {
    const response = await api.post('/deployment-versions', data);
    return response.data;
  },

  getVersionsByApplication: async (applicationName: string): Promise<{ versions: DeploymentVersion[] }> => {
    const response = await api.get(`/deployment-versions/application/${applicationName}`);
    return response.data;
  },
};

