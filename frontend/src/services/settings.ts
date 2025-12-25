import api from '../config/api';

export interface SaltConfig {
  api_url: string;
  username: string;
  password: string;
  eauth: string;
  timeout: number;
  verify_ssl: boolean;
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
};

export default settingsService;

