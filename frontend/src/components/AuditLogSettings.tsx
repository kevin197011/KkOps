import React, { useState, useEffect } from 'react';
import { Form, InputNumber, Button, message, Space, Alert, Statistic, Row, Col } from 'antd';
import { SaveOutlined, DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import settingsService from '../services/settings';

const AuditLogSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [cleanupLoading, setCleanupLoading] = useState(false);
  const [stats, setStats] = useState<{
    total: number;
    oldRecords: number;
    retentionDays: number;
  } | null>(null);

  useEffect(() => {
    loadSettings();
    loadStats();
  }, []);

  const loadSettings = async () => {
    setLoading(true);
    try {
      const response = await settingsService.getAuditLogSettings();
      form.setFieldsValue({
        retention_days: response.retention_days,
      });
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载审计日志设置失败');
    } finally {
      setLoading(false);
    }
  };

  const loadStats = async () => {
    try {
      const response = await settingsService.getAuditLogStats();
      setStats(response);
    } catch (error: any) {
      console.error('Failed to load audit log stats:', error);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      await settingsService.updateAuditLogSettings({
        retention_days: values.retention_days,
      });
      message.success('审计日志设置保存成功');
      loadStats(); // 重新加载统计信息
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '保存审计日志设置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCleanupNow = async () => {
    try {
      setCleanupLoading(true);
      const response = await settingsService.cleanupAuditLogs();
      message.success(`清理完成，删除了 ${response.deleted_count} 条过期记录`);
      loadStats(); // 重新加载统计信息
    } catch (error: any) {
      message.error(error.response?.data?.error || '清理审计日志失败');
    } finally {
      setCleanupLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '800px', marginTop: '24px' }}>
      <Alert
        message="审计日志自动清理"
        description="系统会在每天凌晨3点自动清理超过保存天数的审计日志记录，以节省存储空间。"
        type="info"
        showIcon
        style={{ marginBottom: '24px' }}
      />

      {stats && (
        <Row gutter={16} style={{ marginBottom: '24px' }}>
          <Col span={8}>
            <Statistic
              title="总记录数"
              value={stats.total}
              suffix="条"
            />
          </Col>
          <Col span={8}>
            <Statistic
              title="过期记录数"
              value={stats.oldRecords}
              suffix="条"
              valueStyle={{ color: stats.oldRecords > 0 ? '#cf1322' : '#3f8600' }}
            />
          </Col>
          <Col span={8}>
            <Statistic
              title="当前保存天数"
              value={stats.retentionDays}
              suffix="天"
            />
          </Col>
        </Row>
      )}

      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
      >
        <Form.Item
          name="retention_days"
          label="审计日志保存天数"
          rules={[
            { required: true, message: '请输入保存天数' },
            { type: 'number', min: 1, max: 3650, message: '保存天数必须在1-3650天之间' },
          ]}
          tooltip="审计日志的保存天数，超过此天数的记录将被自动删除。建议设置为30-90天。"
        >
          <InputNumber
            min={1}
            max={3650}
            style={{ width: '200px' }}
            placeholder="30"
            addonAfter="天"
          />
        </Form.Item>

        <Form.Item>
          <Space>
            <Button
              type="primary"
              htmlType="submit"
              icon={<SaveOutlined />}
              loading={loading}
            >
              保存设置
            </Button>
            <Button
              icon={<DeleteOutlined />}
              onClick={handleCleanupNow}
              loading={cleanupLoading}
              danger
            >
              立即清理过期记录
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={() => {
                loadSettings();
                loadStats();
              }}
            >
              刷新
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </div>
  );
};

export default AuditLogSettings;