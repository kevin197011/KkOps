import api from '../config/api';

export interface SaltConfig {
  api_url: string;
  username: string;
  password: string;
  eauth: string;
  timeout: number;
  verify_ssl: boolean;
}

export interface AuditLogSettings {
  retention_days: number;
}

export interface AuditLogStats {
  total: number;
  oldRecords: number;
  retentionDays: number;
}

export interface SettingsResponse {
  settings: Record<string, string>;
}

export interface SaltConfigResponse {
  config: SaltConfig;
}

const settingsService = {
  // 获取所有设置
  getAllSettings: async (): Promise<SettingsResponse> => {
    const response = await api.get<SettingsResponse>('/settings');
    return response.data;
  },

  // 按分类获取设置
  getSettingsByCategory: async (category: string): Promise<SettingsResponse> => {
    const response = await api.get<SettingsResponse>(`/settings/${category}`);
    return response.data;
  },

  // 获取单个设置
  getSetting: async (key: string): Promise<{ setting: { key: string; value: string; category: string; description: string } }> => {
    const response = await api.get(`/settings/key/${encodeURIComponent(key)}`);
    return response.data;
  },

  // 更新单个设置
  updateSetting: async (key: string, value: string, category?: string, description?: string): Promise<void> => {
    await api.put(`/settings/key/${encodeURIComponent(key)}`, {
      value,
      category,
      description,
    });
  },

  // 获取Salt配置
  getSaltConfig: async (): Promise<SaltConfigResponse> => {
    const response = await api.get<SaltConfigResponse>('/settings/salt');
    return response.data;
  },

  // 更新Salt配置
  updateSaltConfig: async (config: {
    api_url: string;
    username: string;
    password: string;
    eauth: string;
    timeout: number;
    verify_ssl: boolean;
  }): Promise<void> => {
    await api.put('/settings/salt', config);
  },

  // 测试Salt API连接
  testSaltConnection: async (config?: {
    api_url: string;
    username: string;
    password: string;
    eauth: string;
    timeout: number;
    verify_ssl: boolean;
  }): Promise<{ success: boolean; message: string; error?: string }> => {
    const response = await api.post('/settings/salt/test', config || {});
    return response.data;
  },

  // 获取审计日志设置
  getAuditLogSettings: async (): Promise<AuditLogSettings> => {
    const response = await api.get('/settings/audit');
    return response.data;
  },

  // 更新审计日志设置
  updateAuditLogSettings: async (settings: AuditLogSettings): Promise<void> => {
    await api.put('/settings/audit', settings);
  },

  // 获取审计日志统计信息
  getAuditLogStats: async (): Promise<AuditLogStats> => {
    const response = await api.get('/settings/audit/stats');
    return response.data;
  },

  // 立即清理过期审计日志
  cleanupAuditLogs: async (): Promise<{ deleted_count: number }> => {
    const response = await api.post('/settings/audit/cleanup');
    return response.data;
  },
};

export default settingsService;

