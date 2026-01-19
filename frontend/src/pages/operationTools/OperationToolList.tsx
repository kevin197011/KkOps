// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
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
} from '@ant-design/icons'
import {
  operationToolApi,
  OperationTool,
  CreateOperationToolRequest,
  UpdateOperationToolRequest,
} from '@/api/operationTool'
import { usePermissionStore } from '@/stores/permission'

const { TextArea } = Input

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
  const { token } = theme.useToken()
  const [tools, setTools] = useState<OperationTool[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingTool, setEditingTool] = useState<OperationTool | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchTools()
  }, [])

  const fetchTools = async () => {
    setLoading(true)
    try {
      const response = await operationToolApi.list()
      setTools(response.data.data)
    } catch (error: any) {
      message.error('获取运维工具列表失败')
    } finally {
      setLoading(false)
    }
  }

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
      content: '确定要删除这个运维工具吗？删除后将从仪表盘和此页面中移除。',
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
  const renderIcon = (iconName?: string, iconUrl?: string) => {
    if (iconUrl && iconUrl.startsWith('http')) {
      return <img src={iconUrl} alt="" style={{ width: 20, height: 20 }} />
    }
    if (iconName && iconMap[iconName]) {
      return iconMap[iconName]
    }
    return <ToolOutlined />
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
      render: (category: string) => (category ? <Tag>{category}</Tag> : '-'),
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
        <h2>运维导航</h2>
        {canManage && (
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增工具">
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
