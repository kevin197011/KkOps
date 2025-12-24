import api from '../config/api';

export interface Role {
  id: number;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
  permissions?: Permission[];
}

export interface Permission {
  id: number;
  code: string;
  name: string;
  resource_type: string;
  action: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface ListRolesResponse {
  roles: Role[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreateRoleRequest {
  name: string;
  description?: string;
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
}

export const roleService = {
  list: async (page: number = 1, pageSize: number = 10): Promise<ListRolesResponse> => {
    const response = await api.get('/roles', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ role: Role }> => {
    const response = await api.get(`/roles/${id}`);
    return response.data;
  },

  create: async (data: CreateRoleRequest): Promise<{ role: Role }> => {
    const response = await api.post('/roles', data);
    return response.data;
  },

  update: async (id: number, data: UpdateRoleRequest): Promise<{ role: Role }> => {
    const response = await api.put(`/roles/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/roles/${id}`);
  },

  assignPermission: async (roleId: number, permissionId: number): Promise<void> => {
    await api.post(`/roles/${roleId}/permissions`, { permission_id: permissionId });
  },

  removePermission: async (roleId: number, permissionId: number): Promise<void> => {
    await api.delete(`/roles/${roleId}/permissions`, {
      data: { permission_id: permissionId },
    });
  },
};

