import api from '../config/api';

export interface HostTag {
  id: number;
  name: string;
  color?: string;
  description?: string;
  created_at: string;
  hosts?: Array<{
    id: number;
    hostname: string;
  }>;
}

export interface ListHostTagsResponse {
  tags: HostTag[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreateHostTagRequest {
  name: string;
  color?: string;
  description?: string;
}

export interface UpdateHostTagRequest {
  name?: string;
  color?: string;
  description?: string;
}

export const hostTagService = {
  list: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: { name?: string }
  ): Promise<ListHostTagsResponse> => {
    const response = await api.get('/host-tags', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ tag: HostTag }> => {
    const response = await api.get(`/host-tags/${id}`);
    return response.data;
  },

  create: async (data: CreateHostTagRequest): Promise<{ tag: HostTag }> => {
    const response = await api.post('/host-tags', data);
    return response.data;
  },

  update: async (id: number, data: UpdateHostTagRequest): Promise<void> => {
    await api.put(`/host-tags/${id}`, data);
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/host-tags/${id}`);
  },
};

