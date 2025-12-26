import api from '../config/api';

export interface CloudPlatform {
  id: number;
  name: string;
  display_name: string;
  icon: string;
  color: string;
  sort_order: number;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCloudPlatformRequest {
  name: string;
  display_name: string;
  icon?: string;
  color?: string;
  sort_order?: number;
  description?: string;
}

export interface UpdateCloudPlatformRequest {
  name?: string;
  display_name?: string;
  icon?: string;
  color?: string;
  sort_order?: number;
  description?: string;
}

export const cloudPlatformService = {
  list: async (page: number = 1, pageSize: number = 100) => {
    const response = await api.get<{
      cloud_platforms: CloudPlatform[];
      total: number;
      page: number;
      page_size: number;
    }>('/cloud-platforms', { params: { page, page_size: pageSize } });
    return response.data;
  },

  get: async (id: number) => {
    const response = await api.get<{ cloud_platform: CloudPlatform }>(`/cloud-platforms/${id}`);
    return response.data;
  },

  create: async (data: CreateCloudPlatformRequest) => {
    const response = await api.post<{ cloud_platform: CloudPlatform }>('/cloud-platforms', data);
    return response.data;
  },

  update: async (id: number, data: UpdateCloudPlatformRequest) => {
    const response = await api.put<{ cloud_platform: CloudPlatform }>(`/cloud-platforms/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await api.delete(`/cloud-platforms/${id}`);
    return response.data;
  },
};
