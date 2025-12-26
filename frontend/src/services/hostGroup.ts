import api from '../config/api';

export interface HostGroup {
  id: number;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
  hosts?: Array<{
    id: number;
    hostname: string;
  }>;
}

export interface ListHostGroupsResponse {
  groups: HostGroup[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreateHostGroupRequest {
  name: string;
  description?: string;
}

export interface UpdateHostGroupRequest {
  name?: string;
  description?: string;
}

export const hostGroupService = {
  list: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: { name?: string }
  ): Promise<ListHostGroupsResponse> => {
    const response = await api.get('/host-groups', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ group: HostGroup }> => {
    const response = await api.get(`/host-groups/${id}`);
    return response.data;
  },

  create: async (data: CreateHostGroupRequest): Promise<{ group: HostGroup }> => {
    const response = await api.post('/host-groups', data);
    return response.data;
  },

  update: async (id: number, data: UpdateHostGroupRequest): Promise<void> => {
    await api.put(`/host-groups/${id}`, data);
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/host-groups/${id}`);
  },
};
