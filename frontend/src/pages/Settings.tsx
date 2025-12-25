import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, InputNumber, Switch, message, Space, Tabs } from 'antd';
import { SaveOutlined, ReloadOutlined, ApiOutlined } from '@ant-design/icons';
import settingsService, { SaltConfig } from '../services/settings';

const { TabPane } = Tabs;

const Settings: React.FC = () => {
  const [saltConfig, setSaltConfig] = useState<SaltConfig | null>(null);
  const [loading, setLoading] = useState(false);
  const [testing, setTesting] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    loadSaltConfig();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadSaltConfig = async () => {
    setLoading(true);
    try {
      const response = await settingsService.getSaltConfig();
      setSaltConfig(response.config);
      // 设置表单初始值，确保已配置的字段正常显示（非灰色占位符）
      // 注意：密码字段始终设置为空字符串，但我们会通过 saltConfig.password 来判断是否已配置
      const config = response.config;
      // 确保所有字段都能正确显示配置的值（非灰色占位符）
      form.setFieldsValue({
        api_url: config.api_url ?? '',
        username: config.username ?? '',
        password: '', // 密码字段始终为空（不显示实际密码），占位符会根据 config.password 是否等于 '***' 来显示不同提示
        eauth: config.eauth ?? '',
        timeout: config.timeout ?? 30, // 使用 ?? 而不是 ||，避免 0 被当作 falsy
        verify_ssl: config.verify_ssl !== undefined ? config.verify_ssl : false,
      });
    } catch (error: any) {
      console.error('Failed to load salt config:', error); // 调试日志
      message.error(error.response?.data?.error || '加载Salt配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      // 如果密码字段为空，传递空字符串（后端会保留现有密码）
      await settingsService.updateSaltConfig({
        api_url: values.api_url,
        username: values.username,
        password: values.password || '', // 空字符串表示保留现有密码
        eauth: values.eauth,
        timeout: values.timeout,
        verify_ssl: values.verify_ssl,
      });
      message.success('Salt配置保存成功');
      loadSaltConfig(); // 重新加载配置
    } catch (error: any) {
      if (error.errorFields) {
        return; // 表单验证错误
      }
      message.error(error.response?.data?.error || '保存Salt配置失败');
    }
  };

  const handleTestConnection = async () => {
    try {
      // 先验证表单字段（密码字段不验证，因为可以为空）
      const values = await form.validateFields(['api_url', 'username', 'eauth', 'timeout', 'verify_ssl']);
      setTesting(true);
      
      // 使用表单中的值进行测试
      // 如果密码字段为空，传递空字符串，后端会从数据库读取现有密码
      const result = await settingsService.testSaltConnection({
        api_url: values.api_url,
        username: values.username,
        password: values.password || '', // 如果为空，后端会从数据库读取现有密码
        eauth: values.eauth,
        timeout: values.timeout,
        verify_ssl: values.verify_ssl,
      });

      if (result.success) {
        message.success(result.message || '连接测试成功');
      } else {
        message.error(result.error || result.message || '连接测试失败');
      }
    } catch (error: any) {
      if (error.errorFields) {
        message.warning('请先填写完整的配置信息');
        return; // 表单验证错误
      }
      message.error(error.response?.data?.error || '连接测试失败');
    } finally {
      setTesting(false);
    }
  };

  return (
    <Card>
      <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>系统设置</h2>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={loadSaltConfig} loading={loading}>
            刷新
          </Button>
        </Space>
      </div>

      <Tabs defaultActiveKey="salt">
        <TabPane tab="Salt API 配置" key="salt">
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            style={{ maxWidth: '800px', marginTop: '24px' }}
          >
            <Form.Item
              name="api_url"
              label="API URL"
              rules={[
                { required: true, message: '请输入Salt API URL' },
                { type: 'url', message: '请输入有效的URL格式' },
              ]}
              tooltip="Salt API的URL地址，例如: http://192.168.56.10:8000"
            >
              <Input placeholder="http://192.168.56.10:8000" />
            </Form.Item>

            <Form.Item
              name="username"
              label="用户名"
              rules={[{ required: true, message: '请输入Salt API用户名' }]}
            >
              <Input placeholder="salt" />
            </Form.Item>

            <Form.Item
              name="password"
              label="密码"
              rules={[]}
              tooltip="如果不想修改密码，请保持此字段为空。首次配置时请输入密码。"
            >
              <Input.Password 
                placeholder={saltConfig?.password === '***' ? '********（已配置，如需修改请输入新密码）' : '请输入密码'}
              />
            </Form.Item>

            <Form.Item
              name="eauth"
              label="EAuth 类型"
              rules={[{ required: true, message: '请输入EAuth类型' }]}
              tooltip="Salt API认证类型，常见值: pam, ldap等"
            >
              <Input placeholder="pam" />
            </Form.Item>

            <Form.Item
              name="timeout"
              label="超时时间（秒）"
              rules={[
                { required: true, message: '请输入超时时间' },
                { type: 'number', min: 1, max: 300, message: '超时时间必须在1-300秒之间' },
              ]}
              tooltip="Salt API请求的超时时间，单位：秒"
            >
              <InputNumber min={1} max={300} style={{ width: '100%' }} placeholder="30" />
            </Form.Item>

            <Form.Item
              name="verify_ssl"
              label="验证SSL证书"
              valuePropName="checked"
              tooltip="是否验证Salt API的SSL证书"
            >
              <Switch />
            </Form.Item>

            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={loading}>
                  保存配置
                </Button>
                <Button 
                  icon={<ApiOutlined />} 
                  onClick={handleTestConnection} 
                  loading={testing}
                >
                  测试连接
                </Button>
                <Button onClick={() => form.resetFields()}>
                  重置
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </TabPane>
      </Tabs>
    </Card>
  );
};

export default Settings;

