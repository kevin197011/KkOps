import api from '../config/api';

export interface AuditLog {
  id: number;
  user_id?: number;
  username: string;
  action: string;
  resource_type?: string;
  resource_id?: number;
  resource_name?: string;
  ip_address?: string;
  user_agent?: string;
  request_data?: string;
  response_data?: string;
  before_data?: string;
  after_data?: string;
  status: 'success' | 'failed';
  error_message?: string;
  duration_ms?: number;
  created_at: string;
}

export interface ListAuditLogsResponse {
  logs: AuditLog[];
  total: number;
  page: number;
  page_size: number;
}

export interface AuditLogFilters {
  user_id?: number;
  username?: string;
  action?: string;
  resource_type?: string;
  resource_id?: number;
  status?: string;
  start_time?: string;
  end_time?: string;
}

export const auditService = {
  list: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: AuditLogFilters
  ): Promise<ListAuditLogsResponse> => {
    const response = await api.get('/audit/logs', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ log: AuditLog }> => {
    const response = await api.get(`/audit/logs/${id}`);
    return response.data;
  },

  getByResource: async (
    resourceType: string,
    resourceId: number
  ): Promise<{ logs: AuditLog[] }> => {
    const response = await api.get('/audit/logs/resource', {
      params: { resource_type: resourceType, resource_id: resourceId },
    });
    return response.data;
  },

  getByUser: async (
    userId: number,
    limit: number = 100
  ): Promise<{ logs: AuditLog[] }> => {
    const response = await api.get(`/audit/logs/user/${userId}`, {
      params: { limit },
    });
    return response.data;
  },
};

