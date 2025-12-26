import api from '../config/api';

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

export interface ListPermissionsResponse {
  permissions: Permission[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreatePermissionRequest {
  code: string;
  name: string;
  resource_type: string;
  action: string;
  description?: string;
}

export interface UpdatePermissionRequest {
  name?: string;
  description?: string;
}

export const permissionService = {
  list: async (page: number = 1, pageSize: number = 10): Promise<ListPermissionsResponse> => {
    const response = await api.get('/permissions', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  },

  listAll: async (): Promise<{ permissions: Permission[] }> => {
    let allPermissions: Permission[] = [];
    let page = 1;
    const pageSize = 100;
    let hasMore = true;

    while (hasMore) {
      const response = await api.get('/permissions', {
        params: { page, page_size: pageSize },
      });
      
      const { permissions, total } = response.data;
      allPermissions = [...allPermissions, ...permissions];
      
      // Check if we have more pages
      hasMore = allPermissions.length < total;
      page++;
    }

    return { permissions: allPermissions };
  },

  listByResourceType: async (resourceType: string): Promise<{ permissions: Permission[] }> => {
    const response = await api.get(`/permissions/resource/${resourceType}`);
    return response.data;
  },

  get: async (id: number): Promise<{ permission: Permission }> => {
    const response = await api.get(`/permissions/${id}`);
    return response.data;
  },

  create: async (data: CreatePermissionRequest): Promise<{ permission: Permission }> => {
    const response = await api.post('/permissions', data);
    return response.data;
  },

  update: async (id: number, data: UpdatePermissionRequest): Promise<{ permission: Permission }> => {
    const response = await api.put(`/permissions/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/permissions/${id}`);
  },
};

