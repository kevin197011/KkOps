import api from '../config/api';

export interface Environment {
  id: number;
  name: string;
  display_name: string;
  color: string;
  sort_order: number;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEnvironmentRequest {
  name: string;
  display_name: string;
  color?: string;
  sort_order?: number;
  description?: string;
}

export interface UpdateEnvironmentRequest {
  name?: string;
  display_name?: string;
  color?: string;
  sort_order?: number;
  description?: string;
}

export const environmentService = {
  list: async (page: number = 1, pageSize: number = 100) => {
    const response = await api.get<{
      environments: Environment[];
      total: number;
      page: number;
      page_size: number;
    }>('/environments', { params: { page, page_size: pageSize } });
    return response.data;
  },

  get: async (id: number) => {
    const response = await api.get<{ environment: Environment }>(`/environments/${id}`);
    return response.data;
  },

  create: async (data: CreateEnvironmentRequest) => {
    const response = await api.post<{ environment: Environment }>('/environments', data);
    return response.data;
  },

  update: async (id: number, data: UpdateEnvironmentRequest) => {
    const response = await api.put<{ environment: Environment }>(`/environments/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await api.delete(`/environments/${id}`);
    return response.data;
  },
};
