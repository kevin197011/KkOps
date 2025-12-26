import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  message,
  Tag,
  Card,
  Modal,
  Form,
  Input,
  Popconfirm,
  Descriptions,
} from 'antd';
import {
  ReloadOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
} from '@ant-design/icons';
import { sshKeyService, SSHKey } from '../services/sshKey';

const SSHKeys: React.FC = () => {
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [editingSshKey, setEditingSshKey] = useState<SSHKey | null>(null);
  const [detailSshKey, setDetailSshKey] = useState<SSHKey | null>(null);
  const [form] = Form.useForm();
  const [showPrivateKey, setShowPrivateKey] = useState(false); // 控制私钥显示/隐藏

  useEffect(() => {
    loadSSHKeys();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const loadSSHKeys = async () => {
    setLoading(true);
    try {
      const response = await sshKeyService.list(page, pageSize);
      setSshKeys(response.ssh_keys);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载SSH密钥列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingSshKey(null);
    setShowPrivateKey(true); // 创建模式下默认显示私钥输入框
    form.resetFields();
    setModalVisible(true);
  };

  const handleViewDetail = async (key: SSHKey) => {
    try {
      const response = await sshKeyService.get(key.id);
      setDetailSshKey(response.ssh_key);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取SSH密钥详情失败');
    }
  };

  // 掩码私钥内容（保留 BEGIN 和 END 行，中间用星号替换）
  const maskPrivateKey = (privateKey: string): string => {
    if (!privateKey || !privateKey.trim()) {
      return '';
    }
    
    const lines = privateKey.split('\n');
    const beginLine = lines.find(line => line.includes('BEGIN'));
    const endLine = lines.find(line => line.includes('END'));
    
    if (beginLine && endLine) {
      // 保留 BEGIN 和 END 行，中间用星号替换
      return `${beginLine}\n${'*'.repeat(64)}\n${'*'.repeat(64)}\n${endLine}`;
    }
    
    // 如果不是标准格式，返回掩码
    return '*'.repeat(64);
  };

  const handleEdit = async (key: SSHKey) => {
    try {
      // 从服务器重新获取完整的密钥信息
      const response = await sshKeyService.get(key.id);
      // 获取私钥内容（用于编辑）
      let privateKey = '';
      try {
        const privateKeyResponse = await sshKeyService.getPrivateKey(key.id);
        privateKey = privateKeyResponse.private_key;
      } catch (error: any) {
        // 如果获取私钥失败，忽略错误（可能权限问题，继续编辑其他字段）
        console.warn('获取私钥内容失败:', error);
      }
      
      setEditingSshKey(response.ssh_key);
      setShowPrivateKey(false); // 编辑模式下默认隐藏私钥
      form.setFieldsValue({ 
        name: response.ssh_key.name, 
        username: response.ssh_key.username || '',
        private_key_content: privateKey // 存储原始私钥，但显示时用掩码
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取SSH密钥详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await sshKeyService.delete(id);
      message.success('删除成功');
      loadSSHKeys();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingSshKey) {
        const updateData: any = {
          name: values.name,
          username: values.username || undefined,
        };
        // 只有在提供了私钥内容时才包含此字段
        if (values.private_key_content && values.private_key_content.trim()) {
          updateData.private_key_content = values.private_key_content.trim();
        }
        await sshKeyService.update(editingSshKey.id, updateData);
        message.success('更新成功');
      } else {
        await sshKeyService.create({
          name: values.name,
          username: values.username || undefined,
          private_key_content: values.private_key_content,
        });
        message.success('创建成功');
      }
      setModalVisible(false);
      setEditingSshKey(null);
      setShowPrivateKey(false);
      form.resetFields();
      loadSSHKeys();
    } catch (error: any) {
      if (error.errorFields) {
        return; // 表单验证错误
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const columns = [
    {
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 120,
      render: (username: string) => username || <span style={{ color: '#999' }}>-</span>,
    },
    {
      title: '类型',
      dataIndex: 'key_type',
      key: 'key_type',
      width: 100,
      render: (type: string) => <Tag>{type.toUpperCase()}</Tag>,
    },
    {
      title: '指纹',
      dataIndex: 'fingerprint',
      key: 'fingerprint',
      render: (fp: string) => (
        <span style={{ fontFamily: 'monospace', fontSize: '12px' }}>{fp || '-'}</span>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => (text ? new Date(text).toLocaleString() : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right' as const,
      render: (_: any, record: SSHKey) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个SSH密钥吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2 style={{ margin: 0 }}>SSH密钥管理</h2>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={loadSSHKeys}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              上传密钥
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={sshKeys}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            pageSizeOptions: ['10', '20', '50', '100'],
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: (page, pageSize) => {
              setPage(page);
              setPageSize(pageSize);
            },
          }}
        />
      </Card>

      {/* SSH密钥上传/编辑模态框 */}
      <Modal
        title={editingSshKey ? '编辑SSH密钥' : '上传SSH密钥'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => {
          setModalVisible(false);
          setEditingSshKey(null);
          setShowPrivateKey(false);
          form.resetFields();
        }}
        okText="确定"
        cancelText="取消"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="密钥名称"
            rules={[{ required: true, message: '请输入密钥名称' }]}
          >
            <Input placeholder="例如: 我的SSH密钥" />
          </Form.Item>
          <Form.Item
            name="username"
            label="SSH用户名"
            tooltip="可选：设置默认SSH用户名，连接时将自动使用此用户名"
          >
            <Input placeholder="例如: root, ubuntu, admin" />
          </Form.Item>
          <Form.Item
            name="private_key_content"
            label={
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>私钥内容</span>
                {editingSshKey && (
                  <Button
                    type="link"
                    size="small"
                    icon={showPrivateKey ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                    onClick={() => {
                      setShowPrivateKey(!showPrivateKey);
                    }}
                  >
                    {showPrivateKey ? '隐藏' : '显示'}
                  </Button>
                )}
              </div>
            }
            rules={[
              { required: !editingSshKey, message: '请输入SSH私钥内容' },
              {
                validator: (_, value) => {
                  if (editingSshKey && !value) {
                    // 编辑模式下，私钥内容为可选（留空则不更新私钥）
                    return Promise.resolve();
                  }
                  if (!value && !editingSshKey) {
                    return Promise.reject(new Error('请输入SSH私钥内容'));
                  }
                  return Promise.resolve();
                },
              },
            ]}
            tooltip={editingSshKey ? '可选：如需更新私钥，请输入新的私钥内容；留空则保持原私钥不变' : '请输入SSH私钥内容'}
          >
            <Form.Item noStyle shouldUpdate={(prevValues, currentValues) => prevValues.private_key_content !== currentValues.private_key_content}>
              {({ getFieldValue, setFieldValue }) => {
                const currentValue = getFieldValue('private_key_content') || '';
                // 编辑模式下，如果不显示私钥且有值，使用掩码；否则显示原始值
                const displayValue = editingSshKey && !showPrivateKey && currentValue
                  ? maskPrivateKey(currentValue)
                  : currentValue;

                return (
                  <Input.TextArea
                    rows={10}
                    value={displayValue}
                    readOnly={editingSshKey && !showPrivateKey && currentValue ? true : false}
                    onFocus={() => {
                      // 当用户聚焦到掩码的输入框时，切换到显示模式
                      if (editingSshKey && !showPrivateKey && currentValue) {
                        setShowPrivateKey(true);
                      }
                    }}
                    onChange={(e) => {
                      const newValue = e.target.value;
                      // 如果当前是掩码模式，用户修改内容时切换到显示模式
                      if (editingSshKey && !showPrivateKey && currentValue) {
                        const maskedValue = maskPrivateKey(currentValue);
                        // 如果用户修改的内容与掩码不同，切换到显示模式并使用原始值
                        if (newValue !== maskedValue) {
                          setShowPrivateKey(true);
                          // 保持原始值，让用户可以继续编辑
                          return;
                        }
                      }
                      setFieldValue('private_key_content', newValue);
                    }}
                    placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
                    style={{ 
                      fontFamily: 'monospace', 
                      fontSize: '12px',
                      cursor: editingSshKey && !showPrivateKey && currentValue ? 'pointer' : 'text'
                    }}
                  />
                );
              }}
            </Form.Item>
          </Form.Item>
        </Form>
      </Modal>

      {/* SSH密钥详情模态框 */}
      <Modal
        title="SSH密钥详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {detailSshKey && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="名称">{detailSshKey.name}</Descriptions.Item>
            <Descriptions.Item label="用户名">{detailSshKey.username || '-'}</Descriptions.Item>
            <Descriptions.Item label="类型">
              <Tag>{detailSshKey.key_type?.toUpperCase() || '-'}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="指纹">
              <span style={{ fontFamily: 'monospace', fontSize: '12px' }}>
                {detailSshKey.fingerprint || '-'}
              </span>
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailSshKey.created_at ? new Date(detailSshKey.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailSshKey.updated_at ? new Date(detailSshKey.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </>
  );
};

export default SSHKeys;

