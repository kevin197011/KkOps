import api from '../config/api';

export interface Formula {
  id: number;
  name: string;
  description: string;
  category: string;
  version: string;
  path: string;
  repository: string;
  icon?: string;
  tags?: string[];
  metadata?: any;
  is_active: boolean;
}

export interface FormulaParameter {
  id: number;
  name: string;
  type: string;
  default?: any;
  required: boolean;
  label: string;
  description: string;
  validation?: any;
  order: number;
}

export interface FormulaTemplate {
  id: number;
  formula_id: number;
  name: string;
  description: string;
  pillar_data?: any;
  is_public: boolean;
  created_by: number;
  created_at: string;
  updated_at: string;
}

export interface FormulaDeployment {
  id: number;
  formula_id: number;
  name: string;
  description: string;
  target_hosts: string[];
  pillar_data?: any;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  salt_job_id?: string;
  results?: any;
  success_count: number;
  failed_count: number;
  error_message?: string;
  started_by: number;
  started_at?: string;
  completed_at?: string;
  duration_seconds?: number;
  created_at: string;
  updated_at: string;
}

export interface ListFormulasResponse {
  formulas: Formula[];
  total: number;
  page: number;
  page_size: number;
}

export interface ListDeploymentsResponse {
  deployments: FormulaDeployment[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreateTemplateRequest {
  name: string;
  description: string;
  pillar_data?: any;
  is_public?: boolean;
}

export interface CreateDeploymentRequest {
  formula_id: number;
  name: string;
  description?: string;
  target_hosts: string[];
  pillar_data?: any;
}

export interface CreateFormulaTemplateRequest {
  formula_id: number;
  name: string;
  description?: string;
  pillar_data?: any;
  is_public?: boolean;
}

export interface FormulaRepository {
  id: number;
  name: string;
  url: string;
  branch: string;
  local_path?: string;
  is_active: boolean;
  last_sync_at?: string;
  created_by: number;
  created_at: string;
  updated_at: string;
}

const formulaService = {
  // Formula管理
  listFormulas: async (page: number = 1, pageSize: number = 20, filters?: any): Promise<ListFormulasResponse> => {
    const response = await api.get('/formulas', { params: { page, page_size: pageSize, ...filters } });
    return response.data;
  },

  getFormula: async (id: number): Promise<{ formula: Formula; parameters: FormulaParameter[] }> => {
    const response = await api.get(`/formulas/${id}`);
    return response.data;
  },

  createFormula: async (formula: Partial<Formula>): Promise<Formula> => {
    const response = await api.post('/formulas', formula);
    return response.data.formula;
  },

  // Formula参数管理
  updateFormulaParameters: async (formulaId: number, parameters: FormulaParameter[]): Promise<void> => {
    await api.put(`/formulas/${formulaId}/parameters`, { parameters });
  },

  // Formula模板管理
  getFormulaTemplates: async (formulaId: number): Promise<{ templates: FormulaTemplate[] }> => {
    const response = await api.get(`/formulas/${formulaId}/templates`);
    return response.data;
  },

  createTemplate: async (formulaId: number, template: CreateTemplateRequest): Promise<FormulaTemplate> => {
    const response = await api.post(`/formulas/${formulaId}/templates`, template);
    return response.data.template;
  },

  createFormulaTemplate: async (template: CreateFormulaTemplateRequest): Promise<FormulaTemplate> => {
    const response = await api.post(`/formulas/${template.formula_id}/templates`, template);
    return response.data.template;
  },

  // Formula部署
  createDeployment: async (deployment: CreateDeploymentRequest): Promise<FormulaDeployment> => {
    const response = await api.post('/formulas/deployments', deployment);
    return response.data.deployment;
  },

  executeDeployment: async (deploymentId: number): Promise<void> => {
    await api.post(`/formulas/deployments/${deploymentId}/execute`);
  },

  listDeployments: async (page: number = 1, pageSize: number = 20, filters?: any): Promise<ListDeploymentsResponse> => {
    const response = await api.get('/formulas/deployments', { params: { page, page_size: pageSize, ...filters } });
    return response.data;
  },

  getDeployment: async (id: number): Promise<FormulaDeployment> => {
    const response = await api.get(`/formulas/deployments/${id}`);
    return response.data.deployment;
  },

  cancelDeployment: async (deploymentId: number): Promise<void> => {
    await api.post(`/formulas/deployments/${deploymentId}/cancel`);
  },

  cleanupOldDeployments: async (): Promise<{ message: string; count: number }> => {
    const response = await api.post('/formulas/deployments/cleanup');
    return response.data;
  },

  // Formula仓库管理
  createRepository: async (repo: any): Promise<any> => {
    const response = await api.post('/formulas/repositories', repo);
    return response.data.repository;
  },

  getRepository: async (id: number): Promise<FormulaRepository> => {
    const response = await api.get(`/formulas/repositories/${id}`);
    return response.data.repository;
  },

  updateRepository: async (id: number, repo: Partial<FormulaRepository>): Promise<FormulaRepository> => {
    const response = await api.put(`/formulas/repositories/${id}`, repo);
    return response.data.repository;
  },

  listRepositories: async (page: number = 1, pageSize: number = 20, filters?: any): Promise<any> => {
    const response = await api.get('/formulas/repositories', { params: { page, page_size: pageSize, ...filters } });
    return response.data;
  },

  syncRepository: async (repoId: number): Promise<void> => {
    await api.post(`/formulas/repositories/${repoId}/sync`);
  },
};

export default formulaService;
