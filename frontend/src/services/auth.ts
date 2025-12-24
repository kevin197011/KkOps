import api from '../config/api';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  display_name?: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  display_name?: string;
  roles?: Array<{ id: number; name: string }>;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<{ user: User }> => {
    const response = await api.post('/auth/register', data);
    return response.data;
  },

  changePassword: async (oldPassword: string, newPassword: string): Promise<void> => {
    await api.post('/auth/change-password', {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },
};

