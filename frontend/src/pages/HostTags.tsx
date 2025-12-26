import React, { useState, useEffect, useCallback } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Tag,
  ColorPicker,
  Descriptions,
  List,
  Tooltip,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  SearchOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import { hostTagService, HostTag } from '../services/hostTag';

// 预设颜色列表
const PRESET_COLORS = [
  { color: '#1890ff', name: '蓝色' },
  { color: '#52c41a', name: '绿色' },
  { color: '#faad14', name: '橙色' },
  { color: '#f5222d', name: '红色' },
  { color: '#722ed1', name: '紫色' },
  { color: '#13c2c2', name: '青色' },
  { color: '#eb2f96', name: '粉色' },
  { color: '#fa8c16', name: '金色' },
  { color: '#a0d911', name: '青柠' },
  { color: '#2f54eb', name: '极客蓝' },
];

const HostTags: React.FC = () => {
  const [tags, setTags] = useState<HostTag[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [hostsModalVisible, setHostsModalVisible] = useState(false);
  const [editingTag, setEditingTag] = useState<HostTag | null>(null);
  const [detailTag, setDetailTag] = useState<HostTag | null>(null);
  const [selectedTag, setSelectedTag] = useState<HostTag | null>(null);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [batchDeleting, setBatchDeleting] = useState(false);
  const [nameValidating, setNameValidating] = useState(false);
  const [nameExists, setNameExists] = useState(false);
  const [form] = Form.useForm();
  const [searchTimer, setSearchTimer] = useState<NodeJS.Timeout | null>(null);


  useEffect(() => {
    loadTags();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, searchKeyword]);

  const loadTags = async () => {
    setLoading(true);
    try {
      const filters = searchKeyword ? { name: searchKeyword } : undefined;
      const response = await hostTagService.list(page, pageSize, filters);
      setTags(response.tags);
      setTotal(response.total);
    } catch (error: any) {
      message.error(error.response?.data?.error || '加载主机标签列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 搜索处理（带防抖）
  const handleSearch = useCallback((value: string) => {
    if (searchTimer) {
      clearTimeout(searchTimer);
    }
    const timer = setTimeout(() => {
      setSearchKeyword(value);
      setPage(1);
    }, 300);
    setSearchTimer(timer);
  }, [searchTimer]);

  // 名称唯一性校验（带防抖）
  const validateTagName = useCallback(async (name: string) => {
    if (!name || name.trim() === '') {
      setNameExists(false);
      return;
    }
    if (editingTag && editingTag.name === name) {
      setNameExists(false);
      return;
    }
    setNameValidating(true);
    try {
      const response = await hostTagService.list(1, 100, { name: name.trim() });
      const exists = response.tags.some(
        (tag) => tag.name.toLowerCase() === name.trim().toLowerCase() && 
                 (!editingTag || tag.id !== editingTag.id)
      );
      setNameExists(exists);
    } catch (error) {
      // 忽略校验错误
    } finally {
      setNameValidating(false);
    }
  }, [editingTag]);

  const handleCreate = () => {
    setEditingTag(null);
    setNameExists(false);
    form.resetFields();
    setModalVisible(true);
  };

  const handleViewDetail = async (tag: HostTag) => {
    try {
      const response = await hostTagService.get(tag.id);
      setDetailTag(response.tag);
      setDetailModalVisible(true);
    } catch (error: any) {
      message.error('获取主机标签详情失败');
    }
  };

  const handleViewHosts = async (tag: HostTag) => {
    try {
      const response = await hostTagService.get(tag.id);
      setSelectedTag(response.tag);
      setHostsModalVisible(true);
    } catch (error: any) {
      message.error('获取标签详情失败');
    }
  };

  const handleEdit = async (tag: HostTag) => {
    try {
      const response = await hostTagService.get(tag.id);
      setEditingTag(response.tag);
      setNameExists(false);
      form.setFieldsValue({
        name: response.tag.name,
        color: response.tag.color || '#1890ff',
        description: response.tag.description,
      });
      setModalVisible(true);
    } catch (error: any) {
      message.error('获取主机标签详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await hostTagService.delete(id);
      message.success('删除成功');
      loadTags();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };


  // 批量删除
  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) return;
    const selectedTags = tags.filter((tag) => selectedRowKeys.includes(tag.id));
    const usedTags = selectedTags.filter((tag) => (tag.hosts?.length || 0) > 0);

    Modal.confirm({
      title: `确定要删除选中的 ${selectedRowKeys.length} 个标签吗？`,
      icon: <ExclamationCircleOutlined />,
      content: usedTags.length > 0 ? (
        <div>
          <p style={{ color: '#faad14' }}>
            警告：以下 {usedTags.length} 个标签正在被主机使用，删除后将自动解除关联：
          </p>
          <ul>
            {usedTags.map((tag) => (
              <li key={tag.id}>
                <Tag color={tag.color}>{tag.name}</Tag> - 被 {tag.hosts?.length} 台主机使用
              </li>
            ))}
          </ul>
        </div>
      ) : null,
      okText: '确定删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        setBatchDeleting(true);
        let successCount = 0;
        let failCount = 0;
        for (const id of selectedRowKeys) {
          try {
            await hostTagService.delete(id as number);
            successCount++;
          } catch (error) {
            failCount++;
          }
        }
        setBatchDeleting(false);
        setSelectedRowKeys([]);
        if (failCount === 0) {
          message.success(`成功删除 ${successCount} 个标签`);
        } else {
          message.warning(`成功删除 ${successCount} 个，失败 ${failCount} 个`);
        }
        loadTags();
      },
    });
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (nameExists) {
        message.error('标签名称已存在，请使用其他名称');
        return;
      }
      const color = typeof values.color === 'string' ? values.color : values.color?.toHexString();
      if (editingTag) {
        await hostTagService.update(editingTag.id, {
          name: values.name,
          color: color,
          description: values.description,
        });
        message.success('更新成功');
      } else {
        await hostTagService.create({
          name: values.name,
          color: color,
          description: values.description,
        });
        message.success('创建成功');
      }
      setModalVisible(false);
      loadTags();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handlePresetColorClick = (color: string) => {
    form.setFieldsValue({ color });
  };

  const rowSelection = {
    selectedRowKeys,
    onChange: (newSelectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
  };


  const columns = [
    {
      title: '序号',
      key: 'index',
      width: 80,
      render: (_: any, __: any, index: number) => (page - 1) * pageSize + index + 1,
    },
    {
      title: '标签',
      key: 'tag',
      width: 200,
      render: (_: any, record: HostTag) => (
        <Tag color={record.color || '#1890ff'}>{record.name}</Tag>
      ),
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '颜色',
      dataIndex: 'color',
      key: 'color',
      width: 120,
      render: (color: string) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <div
            style={{
              width: 20,
              height: 20,
              backgroundColor: color || '#1890ff',
              borderRadius: 4,
              border: '1px solid #d9d9d9',
            }}
          />
          <span>{color || '#1890ff'}</span>
        </div>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '使用数量',
      key: 'host_count',
      width: 100,
      render: (_: any, record: HostTag) => {
        const count = record.hosts?.length || 0;
        return count > 0 ? (
          <Button type="link" size="small" onClick={() => handleViewHosts(record)} style={{ padding: 0 }}>
            {count} 台主机
          </Button>
        ) : (
          <span style={{ color: '#999' }}>0</span>
        );
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right' as const,
      render: (_: any, record: HostTag) => (
        <Space>
          <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => handleViewDetail(record)}>
            详情
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个主机标签吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger size="small" icon={<DeleteOutlined />}>
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
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
          <h2 style={{ margin: 0 }}>主机标签管理</h2>
          <Space wrap>
            <Input.Search
              placeholder="搜索标签名称"
              allowClear
              style={{ width: 200 }}
              prefix={<SearchOutlined />}
              onChange={(e) => handleSearch(e.target.value)}
              onSearch={(value) => {
                setSearchKeyword(value);
                setPage(1);
              }}
            />
            {selectedRowKeys.length > 0 && (
              <Button danger icon={<DeleteOutlined />} onClick={handleBatchDelete} loading={batchDeleting}>
                批量删除 ({selectedRowKeys.length})
              </Button>
            )}
            <Button icon={<ReloadOutlined />} onClick={loadTags}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              创建主机标签
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={tags}
          rowKey="id"
          loading={loading}
          rowSelection={rowSelection}
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

      {/* 创建/编辑主机标签模态框 */}
      <Modal
        title={editingTag ? '编辑主机标签' : '创建主机标签'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入主机标签名称' }]}
            validateStatus={nameExists ? 'error' : nameValidating ? 'validating' : undefined}
            help={nameExists ? '标签名称已存在' : undefined}
          >
            <Input
              placeholder="请输入主机标签名称"
              onChange={(e) => {
                if (searchTimer) clearTimeout(searchTimer);
                const timer = setTimeout(() => {
                  validateTagName(e.target.value);
                }, 500);
                setSearchTimer(timer);
              }}
            />
          </Form.Item>
          <Form.Item name="color" label="颜色" initialValue="#1890ff">
            <ColorPicker showText />
          </Form.Item>
          <Form.Item label="快捷颜色">
            <Space wrap>
              {PRESET_COLORS.map((preset) => (
                <Tooltip key={preset.color} title={preset.name}>
                  <div
                    onClick={() => handlePresetColorClick(preset.color)}
                    style={{
                      width: 28,
                      height: 28,
                      backgroundColor: preset.color,
                      borderRadius: 4,
                      cursor: 'pointer',
                      border: '2px solid transparent',
                      transition: 'all 0.2s',
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.border = '2px solid #1890ff';
                      e.currentTarget.style.transform = 'scale(1.1)';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.border = '2px solid transparent';
                      e.currentTarget.style.transform = 'scale(1)';
                    }}
                  />
                </Tooltip>
              ))}
            </Space>
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="请输入主机标签描述" rows={4} />
          </Form.Item>
        </Form>
      </Modal>


      {/* 主机标签详情模态框 */}
      <Modal
        title="主机标签详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {detailTag && (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="标签">
              <Tag color={detailTag.color || '#1890ff'}>{detailTag.name}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="名称">{detailTag.name}</Descriptions.Item>
            <Descriptions.Item label="颜色">
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <div
                  style={{
                    width: 20,
                    height: 20,
                    backgroundColor: detailTag.color || '#1890ff',
                    borderRadius: 4,
                    border: '1px solid #d9d9d9',
                  }}
                />
                <span>{detailTag.color || '#1890ff'}</span>
              </div>
            </Descriptions.Item>
            <Descriptions.Item label="描述">{detailTag.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="使用数量">
              {(detailTag.hosts?.length || 0) > 0 ? (
                <Button
                  type="link"
                  size="small"
                  onClick={() => {
                    setDetailModalVisible(false);
                    setSelectedTag(detailTag);
                    setHostsModalVisible(true);
                  }}
                  style={{ padding: 0 }}
                >
                  {detailTag.hosts?.length} 台主机
                </Button>
              ) : (
                '0'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {detailTag.created_at ? new Date(detailTag.created_at).toLocaleString() : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {detailTag.updated_at ? new Date(detailTag.updated_at).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>

      {/* 使用该标签的主机列表模态框 */}
      <Modal
        title={
          <Space>
            <span>使用标签</span>
            {selectedTag && <Tag color={selectedTag.color}>{selectedTag.name}</Tag>}
            <span>的主机</span>
          </Space>
        }
        open={hostsModalVisible}
        onCancel={() => setHostsModalVisible(false)}
        footer={null}
        width={600}
      >
        {selectedTag && (
          <List
            dataSource={selectedTag.hosts || []}
            locale={{ emptyText: '暂无主机使用此标签' }}
            renderItem={(host) => (
              <List.Item>
                <List.Item.Meta
                  title={host.hostname}
                  description={`IP: ${(host as any).ip_address || '-'}`}
                />
              </List.Item>
            )}
          />
        )}
      </Modal>
    </>
  );
};

export default HostTags;
