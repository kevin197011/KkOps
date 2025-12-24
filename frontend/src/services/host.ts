import api from '../config/api';

export interface Host {
  id: number;
  project_id?: number;
  hostname: string;
  ip_address: string;
  salt_minion_id?: string;
  os_type?: string;
  os_version?: string;
  cpu_cores?: number;
  memory_gb?: number;
  disk_gb?: number;
  status: 'online' | 'offline' | 'unknown';
  environment?: string;
  ssh_port?: number;
  ssh_username?: string;
  ssh_key_id?: number;
  last_seen_at?: string;
  salt_version?: string;
  description?: string;
  project?: {
    id: number;
    name: string;
  };
  groups?: Array<{ id: number; name: string }>;
  tags?: Array<{ id: number; name: string; color?: string }>;
}

export interface CreateHostRequest {
  project_id: number;
  hostname: string;
  ip_address: string;
  salt_minion_id?: string;
  os_type?: string;
  os_version?: string;
  cpu_cores?: number;
  memory_gb?: number;
  disk_gb?: number;
  status?: string;
  environment?: string;
  ssh_port?: number;
  ssh_username?: string;
  ssh_key_id?: number;
  description?: string;
}

export interface ListHostsResponse {
  hosts: Host[];
  total: number;
  page: number;
  page_size: number;
}

export interface HostFilters {
  project_id?: number;
  hostname?: string;
  ip_address?: string;
  status?: string;
  environment?: string;
  group_id?: number;
  tag_id?: number;
}

export const hostService = {
  list: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: HostFilters
  ): Promise<ListHostsResponse> => {
    const response = await api.get('/hosts', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ host: Host }> => {
    const response = await api.get(`/hosts/${id}`);
    return response.data;
  },

  create: async (data: CreateHostRequest): Promise<{ host: Host }> => {
    const response = await api.post('/hosts', data);
    return response.data;
  },

  update: async (id: number, data: Partial<CreateHostRequest>): Promise<void> => {
    await api.put(`/hosts/${id}`, data);
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/hosts/${id}`);
  },

  addToGroup: async (hostId: number, groupId: number): Promise<void> => {
    await api.post(`/hosts/${hostId}/groups`, { group_id: groupId });
  },

  removeFromGroup: async (hostId: number, groupId: number): Promise<void> => {
    await api.delete(`/hosts/${hostId}/groups`, {
      data: { group_id: groupId },
    });
  },

  addTag: async (hostId: number, tagId: number): Promise<void> => {
    await api.post(`/hosts/${hostId}/tags`, { tag_id: tagId });
  },

  removeTag: async (hostId: number, tagId: number): Promise<void> => {
    await api.delete(`/hosts/${hostId}/tags`, {
      data: { tag_id: tagId },
    });
  },
};

