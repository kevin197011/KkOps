// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useMemo } from 'react'
import {
  Table,
  Button,
  Space,
  message,
  Modal,
  Form,
  Input,
  InputNumber,
  Switch,
  Tag,
  Tooltip,
  theme,
  Card,
  Row,
  Col,
  Typography,
  Empty,
  Divider,
  Collapse,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  LinkOutlined,
  DashboardOutlined,
  AppstoreOutlined,
  ToolOutlined,
  BugOutlined,
  CodeOutlined,
  FileSearchOutlined,
  DatabaseOutlined,
  RocketOutlined,
  CloudServerOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SettingOutlined,
} from '@ant-design/icons'
import {
  operationToolApi,
  OperationTool,
  CreateOperationToolRequest,
  UpdateOperationToolRequest,
} from '@/api/operationTool'
import { usePermissionStore } from '@/stores/permission'
import { useThemeStore } from '@/stores/theme'

const { TextArea } = Input
const { Title, Text } = Typography
const { Panel } = Collapse

// 图标映射
const iconMap: Record<string, React.ReactNode> = {
  dashboard: <DashboardOutlined />,
  appstore: <AppstoreOutlined />,
  tool: <ToolOutlined />,
  bug: <BugOutlined />,
  code: <CodeOutlined />,
  fileSearch: <FileSearchOutlined />,
  link: <LinkOutlined />,
  database: <DatabaseOutlined />,
  rocket: <RocketOutlined />,
  cloudServer: <CloudServerOutlined />,
}

const OperationToolList = () => {
  const { hasPermission } = usePermissionStore()
  const { mode } = useThemeStore()
  const { token } = theme.useToken()
  const [tools, setTools] = useState<OperationTool[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTool, setEditingTool] = useState<OperationTool | null>(null)
  const [form] = Form.useForm()
  const isDark = mode === 'dark'

  useEffect(() => {
    fetchTools()
  }, [])

  const fetchTools = async () => {
    setLoading(true)
    try {
      // 获取所有工具（包括禁用的）用于配置管理
      const response = await operationToolApi.list()
      setTools(response.data.data)
    } catch (error: any) {
      message.error('获取运维工具列表失败')
    } finally {
      setLoading(false)
    }
  }

  // 按分类分组工具（只显示启用的）
  const toolsByCategory = useMemo(() => {
    const enabledTools = tools.filter((tool) => tool.enabled)
    const grouped: Record<string, OperationTool[]> = {}
    enabledTools.forEach((tool) => {
      const category = tool.category || '未分类'
      if (!grouped[category]) {
        grouped[category] = []
      }
      grouped[category].push(tool)
    })
    
    // 对每个分类的工具按 order 排序
    Object.keys(grouped).forEach((category) => {
      grouped[category].sort((a, b) => a.order - b.order)
    })
    
    return grouped
  }, [tools])

  // 获取所有分类（排序）
  const categories = useMemo(() => {
    return Object.keys(toolsByCategory).sort((a, b) => {
      // 未分类放在最后
      if (a === '未分类') return 1
      if (b === '未分类') return -1
      return a.localeCompare(b)
    })
  }, [toolsByCategory])

  const handleCreate = () => {
    setEditingTool(null)
    form.resetFields()
    form.setFieldsValue({ enabled: true, order: 0 })
    setModalVisible(true)
  }

  const handleEdit = (tool: OperationTool) => {
    setEditingTool(tool)
    form.setFieldsValue(tool)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个运维工具吗？删除后将从导航页面和配置管理中移除。',
      onOk: async () => {
        try {
          await operationToolApi.delete(id)
          message.success('删除成功')
          fetchTools()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateOperationToolRequest | UpdateOperationToolRequest) => {
    try {
      if (editingTool) {
        await operationToolApi.update(editingTool.id, values)
        message.success('更新成功')
      } else {
        await operationToolApi.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchTools()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 渲染图标
  const renderIcon = (iconName?: string, iconUrl?: string, size: number = 24) => {
    if (iconUrl && iconUrl.startsWith('http')) {
      return <img src={iconUrl} alt="" style={{ width: size, height: size }} />
    }
    if (iconName && iconMap[iconName]) {
      return <span style={{ fontSize: size }}>{iconMap[iconName]}</span>
    }
    return <ToolOutlined style={{ fontSize: size }} />
  }

  // URL 验证
  const validateUrl = (_: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error('请输入 URL'))
    }
    if (!value.startsWith('http://') && !value.startsWith('https://')) {
      return Promise.reject(new Error('URL 必须以 http:// 或 https:// 开头'))
    }
    return Promise.resolve()
  }


  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: OperationTool) => (
        <Button
          type="link"
          onClick={() => window.open(record.url, '_blank', 'noopener,noreferrer')}
          style={{ padding: 0 }}
        >
          {text}
        </Button>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (category: string) => (category ? <Tag>{category}</Tag> : <Tag>未分类</Tag>),
    },
    {
      title: '图标',
      dataIndex: 'icon',
      key: 'icon',
      width: 80,
      render: (icon: string, record: OperationTool) => (
        <div style={{ fontSize: 20 }}>{renderIcon(icon, icon?.startsWith('http') ? icon : undefined)}</div>
      ),
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
      ellipsis: { showTitle: false },
      render: (url: string) => (
        <Tooltip title={url}>
          <Button
            type="link"
            icon={<LinkOutlined />}
            onClick={() => window.open(url, '_blank', 'noopener,noreferrer')}
            style={{ padding: 0 }}
          >
            打开链接
          </Button>
        </Tooltip>
      ),
    },
    {
      title: '排序',
      dataIndex: 'order',
      key: 'order',
      width: 80,
      sorter: (a: OperationTool, b: OperationTool) => a.order - b.order,
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 100,
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'success' : 'default'} icon={enabled ? <CheckCircleOutlined /> : <CloseCircleOutlined />}>
          {enabled ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: OperationTool) => {
        const canManage = hasPermission('operation-tools', '*')
        if (!canManage) {
          return '-'
        }
        return (
          <Space size="small">
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            >
              编辑
            </Button>
            <Button
              type="link"
              danger
              size="small"
              icon={<DeleteOutlined />}
              onClick={() => handleDelete(record.id)}
            >
              删除
            </Button>
          </Space>
        )
      },
    },
  ]

  const canManage = hasPermission('operation-tools', '*')

  return (
    <div>
      {/* 导航展示区域 */}
      <div style={{ marginBottom: 32 }}>
        <div
          style={{
            marginBottom: 24,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            flexWrap: 'wrap',
            gap: 16,
          }}
        >
          <Title level={2} style={{ margin: 0 }}>
            运维导航
          </Title>
        </div>

        {loading ? (
          <div style={{ padding: '40px 0', textAlign: 'center' }}>
            <Text type="secondary">加载中...</Text>
          </div>
        ) : categories.length === 0 ? (
          <Empty description="暂无运维工具" />
        ) : (
          categories.map((category) => (
            <div key={category} style={{ marginBottom: 32 }}>
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 12,
                  marginBottom: 16,
                }}
              >
                <Title level={4} style={{ margin: 0 }}>
                  {category}
                </Title>
                <Tag>{toolsByCategory[category].length}</Tag>
              </div>
              <Row gutter={[16, 16]}>
                {toolsByCategory[category].map((tool) => (
                  <Col key={tool.id} xs={24} sm={12} md={8} lg={6} xl={6}>
                    <Card
                      hoverable
                      onClick={() => window.open(tool.url, '_blank', 'noopener,noreferrer')}
                      style={{
                        height: '100%',
                        cursor: 'pointer',
                        borderRadius: 12,
                        border: isDark
                          ? '1px solid rgba(255, 255, 255, 0.1)'
                          : '1px solid rgba(0, 0, 0, 0.06)',
                        transition: 'all 0.2s ease',
                        background: isDark ? 'rgba(255, 255, 255, 0.03)' : '#fff',
                      }}
                      bodyStyle={{
                        padding: 20,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        textAlign: 'center',
                        gap: 12,
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.transform = 'translateY(-4px)'
                        e.currentTarget.style.boxShadow = isDark
                          ? '0 8px 24px rgba(0, 0, 0, 0.3)'
                          : '0 8px 24px rgba(0, 0, 0, 0.1)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.transform = 'translateY(0)'
                        e.currentTarget.style.boxShadow = 'none'
                      }}
                    >
                      <div
                        style={{
                          fontSize: 48,
                          color: token.colorPrimary,
                          marginBottom: 8,
                        }}
                      >
                        {renderIcon(tool.icon, tool.icon?.startsWith('http') ? tool.icon : undefined, 48)}
                      </div>
                      <Title level={5} style={{ margin: 0, fontSize: 16 }}>
                        {tool.name}
                      </Title>
                      {tool.description && (
                        <Text
                          type="secondary"
                          style={{
                            fontSize: 13,
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            lineHeight: '1.5',
                          }}
                        >
                          {tool.description}
                        </Text>
                      )}
                    </Card>
                  </Col>
                ))}
              </Row>
            </div>
          ))
        )}
      </div>

      <Divider />

      {/* 配置管理区域 */}
      <Collapse
        defaultActiveKey={canManage ? ['1'] : []}
        items={[
          {
            key: '1',
            label: (
              <Space>
                <SettingOutlined />
                <span>导航配置管理</span>
              </Space>
            ),
            children: (
              <div>
                <div
                  style={{
                    marginBottom: 16,
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    flexWrap: 'wrap',
                    gap: 8,
                  }}
                >
                  <Text type="secondary">管理和配置运维导航工具</Text>
                  {canManage && (
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
                      新增工具
                    </Button>
                  )}
                </div>
                <Table
                  columns={columns}
                  dataSource={tools}
                  loading={loading}
                  rowKey="id"
                  scroll={{ x: 'max-content' }}
                  pagination={{
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total) => `共 ${total} 条`,
                  }}
                />
              </div>
            ),
          },
        ]}
      />

      {/* 创建/编辑弹窗 */}
      <Modal
        title={editingTool ? '编辑运维工具' : '新增运维工具'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 600 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入工具名称' }]}
          >
            <Input placeholder="工具名称（如：Grafana、Prometheus）" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} placeholder="工具描述（可选）" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Input placeholder="工具分类（如：监控工具、日志工具）" />
          </Form.Item>
          <Form.Item
            name="icon"
            label="图标"
            help="输入 Ant Design 图标名称（如：dashboard、appstore）或图标 URL（以 http:// 或 https:// 开头）"
          >
            <Input placeholder="图标名称或 URL" />
          </Form.Item>
          <Form.Item
            name="url"
            label="URL"
            rules={[{ validator: validateUrl }]}
          >
            <Input placeholder="https://example.com" />
          </Form.Item>
          <Form.Item name="order" label="排序" initialValue={0}>
            <InputNumber min={0} placeholder="排序值，数字越小越靠前" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="enabled" label="启用状态" valuePropName="checked" initialValue={true}>
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default OperationToolList
