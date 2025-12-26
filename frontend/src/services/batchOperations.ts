import api from '../config/api';

export interface BatchOperation {
  id: number;
  name: string;
  description?: string;
  command_type: 'custom' | 'template' | 'builtin';
  command_function: string;
  command_args?: any[];
  target_hosts: Array<{ id: number; hostname: string; ip_address?: string }>;
  target_count: number;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  salt_job_id?: string;
  started_by: number;
  started_at: string;
  completed_at?: string;
  duration_seconds?: number;
  results?: Record<string, { success: boolean; output: string; error?: string }>;
  success_count: number;
  failed_count: number;
  error_message?: string;
  created_at: string;
  updated_at: string;
  started_by_user?: {
    id: number;
    username: string;
    display_name?: string;
  };
}

export interface CreateBatchOperationRequest {
  name: string;
  description?: string;
  command_type: 'custom' | 'template' | 'builtin';
  command_function: string;
  command_args?: any[];
  target_hosts: Array<{ id: number; hostname: string; ip_address?: string }>;
}

export interface ListBatchOperationsResponse {
  operations: BatchOperation[];
  total: number;
  page: number;
  page_size: number;
}

export interface BatchOperationStatus {
  id: number;
  status: string;
}

export interface BatchOperationResults {
  id: number;
  status: string;
  results?: Record<string, { success: boolean; output: string; error?: string }>;
  success_count: number;
  failed_count: number;
}

export const batchOperationsService = {
  // 创建并执行批量操作
  createOperation: async (data: CreateBatchOperationRequest): Promise<{ operation: BatchOperation }> => {
    const response = await api.post('/batch-operations', data);
    return response.data;
  },

  // 列出批量操作
  listOperations: async (
    page: number = 1,
    pageSize: number = 10,
    filters?: { status?: string; view_own?: boolean }
  ): Promise<ListBatchOperationsResponse> => {
    const response = await api.get('/batch-operations', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  // 获取操作详情
  getOperation: async (id: number): Promise<{ operation: BatchOperation }> => {
    const response = await api.get(`/batch-operations/${id}`);
    return response.data;
  },

  // 获取操作状态
  getOperationStatus: async (id: number): Promise<{ id: number; status: string }> => {
    const response = await api.get(`/batch-operations/${id}/status`);
    return response.data;
  },

  // 获取操作结果
  getOperationResults: async (id: number): Promise<BatchOperationResults> => {
    const response = await api.get(`/batch-operations/${id}/results`);
    return response.data;
  },

  // 取消操作
  cancelOperation: async (id: number): Promise<void> => {
    await api.post(`/batch-operations/${id}/cancel`);
  },

  // 重试操作（基于现有操作创建新操作）
  retryOperation: async (operationId: number): Promise<{ operation: BatchOperation }> => {
    const { operation } = await batchOperationsService.getOperation(operationId);

    // 重新创建操作
    return batchOperationsService.createOperation({
      name: `${operation.name} (重试)`,
      description: operation.description,
      command_type: operation.command_type,
      command_function: operation.command_function,
      command_args: operation.command_args,
      target_hosts: operation.target_hosts,
    });
  },

  // 清理1个月前的操作记录
  cleanupOldOperations: async (): Promise<{ message: string; count: number }> => {
    const response = await api.post('/batch-operations/cleanup');
    return response.data;
  },
};

