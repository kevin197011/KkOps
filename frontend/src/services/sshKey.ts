import api from '../config/api';

export interface SSHKey {
  id: number;
  user_id: number;
  name: string;
  username?: string; // SSH用户名（可选）
  key_type: string;
  public_key: string;
  fingerprint: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSSHKeyRequest {
  name: string;
  username?: string; // SSH用户名（可选）
  private_key_content: string;
}

export interface UpdateSSHKeyRequest {
  name: string;
  username?: string; // SSH用户名（可选）
}

export interface ListSSHKeysResponse {
  ssh_keys: SSHKey[];
  total: number;
  page: number;
  page_size: number;
}

export const sshKeyService = {
  create: async (data: CreateSSHKeyRequest): Promise<{ ssh_key: SSHKey }> => {
    const response = await api.post('/ssh-keys', data);
    return response.data;
  },

  list: async (page: number = 1, pageSize: number = 10): Promise<ListSSHKeysResponse> => {
    const response = await api.get('/ssh-keys', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ ssh_key: SSHKey }> => {
    const response = await api.get(`/ssh-keys/${id}`);
    return response.data;
  },

  update: async (id: number, data: UpdateSSHKeyRequest): Promise<{ ssh_key: SSHKey }> => {
    const response = await api.put(`/ssh-keys/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/ssh-keys/${id}`);
  },
};

