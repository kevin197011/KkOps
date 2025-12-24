import api from '../config/api';

export interface SSHConnection {
  id: number;
  name: string;
  host_id?: number;
  hostname: string;
  port: number;
  username: string;
  auth_type: 'password' | 'key';
  status: 'active' | 'disabled';
  last_connected_at?: string;
}

export interface SSHKey {
  id: number;
  name: string;
  key_type: 'rsa' | 'ed25519' | 'ecdsa';
  public_key: string;
  fingerprint: string;
}

export interface SSHSession {
  id: number;
  connection_id: number;
  user_id: number;
  session_token: string;
  client_ip?: string;
  started_at: string;
  ended_at?: string;
  duration_seconds?: number;
  status: 'active' | 'closed' | 'timeout';
}

export interface CreateSSHConnectionRequest {
  project_id: number;
  name: string;
  host_id?: number;
  hostname: string;
  port?: number;
  username: string;
  auth_type: 'password' | 'key';
  password_encrypted?: string;
  key_id?: number;
}

export interface CreateSSHKeyRequest {
  name: string;
  key_type: 'rsa' | 'ed25519' | 'ecdsa';
  private_key_encrypted: string;
  public_key: string;
  fingerprint: string;
  passphrase_encrypted?: string;
}

export interface CreateSSHSessionRequest {
  connection_id: number;
  client_ip?: string;
}

export const sshService = {
  // SSH连接
  listConnections: async (page: number = 1, pageSize: number = 10, filters?: any) => {
    const response = await api.get('/ssh/connections', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  getConnection: async (id: number) => {
    const response = await api.get(`/ssh/connections/${id}`);
    return response.data;
  },

  createConnection: async (data: CreateSSHConnectionRequest) => {
    const response = await api.post('/ssh/connections', data);
    return response.data;
  },

  updateConnection: async (id: number, data: Partial<CreateSSHConnectionRequest>) => {
    await api.put(`/ssh/connections/${id}`, data);
  },

  deleteConnection: async (id: number) => {
    await api.delete(`/ssh/connections/${id}`);
  },

  // SSH密钥
  listKeys: async (page: number = 1, pageSize: number = 10) => {
    const response = await api.get('/ssh/keys', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  },

  getKey: async (id: number) => {
    const response = await api.get(`/ssh/keys/${id}`);
    return response.data;
  },

  createKey: async (data: CreateSSHKeyRequest) => {
    const response = await api.post('/ssh/keys', data);
    return response.data;
  },

  updateKey: async (id: number, data: Partial<CreateSSHKeyRequest>) => {
    await api.put(`/ssh/keys/${id}`, data);
  },

  deleteKey: async (id: number) => {
    await api.delete(`/ssh/keys/${id}`);
  },

  // SSH会话
  listSessions: async (page: number = 1, pageSize: number = 10, filters?: any) => {
    const response = await api.get('/ssh/sessions', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  getSession: async (id: number) => {
    const response = await api.get(`/ssh/sessions/${id}`);
    return response.data;
  },

  createSession: async (data: CreateSSHSessionRequest) => {
    const response = await api.post('/ssh/sessions', data);
    return response.data;
  },

  closeSession: async (id: number) => {
    await api.post(`/ssh/sessions/${id}/close`);
  },
};

