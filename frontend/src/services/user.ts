import api from '../config/api';
import { User } from './auth';

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  display_name?: string;
}

export interface UpdateUserRequest {
  display_name?: string;
  email?: string;
}

export interface ListUsersResponse {
  users: User[];
  total: number;
  page: number;
  page_size: number;
}

export const userService = {
  list: async (page: number = 1, pageSize: number = 10): Promise<ListUsersResponse> => {
    const response = await api.get('/users', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ user: User }> => {
    const response = await api.get(`/users/${id}`);
    return response.data;
  },

  create: async (data: CreateUserRequest): Promise<{ user: User }> => {
    const response = await api.post('/users', data);
    return response.data;
  },

  update: async (id: number, data: UpdateUserRequest): Promise<{ user: User }> => {
    const response = await api.put(`/users/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/users/${id}`);
  },

  assignRole: async (userId: number, roleId: number): Promise<void> => {
    await api.post(`/users/${userId}/roles`, { role_id: roleId });
  },

  removeRole: async (userId: number, roleId: number): Promise<void> => {
    await api.delete(`/users/${userId}/roles`, {
      data: { role_id: roleId },
    });
  },
};

