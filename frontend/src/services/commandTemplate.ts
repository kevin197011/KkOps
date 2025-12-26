import api from '../config/api';

export interface CommandTemplate {
  id: number;
  name: string;
  description?: string;
  category?: 'system' | 'network' | 'disk' | 'process' | 'custom';
  command_function: string;
  command_args?: any[];
  icon?: string;
  created_by: number;
  is_public: boolean;
  usage_count: number;
  created_at: string;
  updated_at: string;
  created_by_user?: {
    id: number;
    username: string;
    display_name?: string;
  };
}

export interface CreateCommandTemplateRequest {
  name: string;
  description?: string;
  category?: 'system' | 'network' | 'disk' | 'process' | 'custom';
  command_function: string;
  command_args?: any[];
  icon?: string;
  is_public?: boolean;
}

export interface UpdateCommandTemplateRequest {
  name?: string;
  description?: string;
  category?: 'system' | 'network' | 'disk' | 'process' | 'custom';
  command_function?: string;
  command_args?: any[];
  icon?: string;
  is_public?: boolean;
}

export const commandTemplateService = {
  // 创建命令模板
  createTemplate: async (data: CreateCommandTemplateRequest): Promise<{ template: CommandTemplate }> => {
    const response = await api.post('/command-templates', data);
    return response.data;
  },

  // 列出命令模板
  listTemplates: async (filters?: { category?: string; is_public?: boolean }): Promise<{ templates: CommandTemplate[] }> => {
    const response = await api.get('/command-templates', {
      params: filters,
    });
    return response.data;
  },

  // 获取模板详情
  getTemplate: async (id: number): Promise<{ template: CommandTemplate }> => {
    const response = await api.get(`/command-templates/${id}`);
    return response.data;
  },

  // 更新模板
  updateTemplate: async (id: number, data: UpdateCommandTemplateRequest): Promise<void> => {
    await api.put(`/command-templates/${id}`, data);
  },

  // 删除模板
  deleteTemplate: async (id: number): Promise<void> => {
    await api.delete(`/command-templates/${id}`);
  },
};

