// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useCallback } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Tag, Checkbox, Transfer, Select, Divider, Tooltip } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, CrownOutlined, SafetyCertificateOutlined, AppstoreAddOutlined } from '@ant-design/icons'
import type { TransferItem } from 'antd/es/transfer'
import { roleApi, Role, CreateRoleRequest, UpdateRoleRequest, RoleAssetInfo } from '@/api/role'
import { assetApi, Asset } from '@/api/asset'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'

const { TextArea } = Input

const RoleList = () => {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [form] = Form.useForm()

  // 资产授权相关状态
  const [authzModalVisible, setAuthzModalVisible] = useState(false)
  const [authzRole, setAuthzRole] = useState<Role | null>(null)
  const [allAssets, setAllAssets] = useState<Asset[]>([])
  const [roleAssets, setRoleAssets] = useState<RoleAssetInfo[]>([])
  const [targetKeys, setTargetKeys] = useState<string[]>([])
  const [authzLoading, setAuthzLoading] = useState(false)

  // 批量授权相关状态
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>()
  const [selectedEnvId, setSelectedEnvId] = useState<number | undefined>()

  useEffect(() => {
    fetchRoles()
  }, [])

  const fetchRoles = async () => {
    setLoading(true)
    try {
      const response = await roleApi.list()
      setRoles(response.data || [])
    } catch (error: any) {
      message.error('获取角色列表失败')
    } finally {
      setLoading(false)
    }
  }

  // 获取所有资产
  const fetchAllAssets = useCallback(async () => {
    try {
      const response = await assetApi.list({ page: 1, page_size: 1000 })
      setAllAssets(response.data.data || [])
    } catch (error: any) {
      message.error('获取资产列表失败')
    }
  }, [])

  // 获取角色已授权的资产
  const fetchRoleAssets = useCallback(async (roleId: number) => {
    try {
      const response = await roleApi.getRoleAssets(roleId)
      const assets = response.data.data || []
      setRoleAssets(assets)
      setTargetKeys(assets.map((a) => String(a.id)))
    } catch (error: any) {
      message.error('获取角色授权资产失败')
    }
  }, [])

  // 获取项目列表
  const fetchProjects = useCallback(async () => {
    try {
      const response = await projectApi.list()
      setProjects(response.data || [])
    } catch (error: any) {
      // 忽略错误
    }
  }, [])

  // 获取环境列表
  const fetchEnvironments = useCallback(async () => {
    try {
      const response = await environmentApi.list()
      setEnvironments(response.data || [])
    } catch (error: any) {
      // 忽略错误
    }
  }, [])

  const handleCreate = () => {
    setEditingRole(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (role: Role) => {
    setEditingRole(role)
    form.setFieldsValue({
      name: role.name,
      description: role.description,
      is_admin: role.is_admin,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个角色吗？',
      onOk: async () => {
        try {
          await roleApi.delete(id)
          message.success('删除成功')
          fetchRoles()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateRoleRequest | UpdateRoleRequest) => {
    try {
      if (editingRole) {
        await roleApi.update(editingRole.id, values as UpdateRoleRequest)
        message.success('更新成功')
      } else {
        await roleApi.create(values as CreateRoleRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchRoles()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  // 打开资产授权弹窗
  const handleOpenAuthzModal = async (role: Role) => {
    setAuthzRole(role)
    setAuthzLoading(true)
    setAuthzModalVisible(true)
    setSelectedProjectId(undefined)
    setSelectedEnvId(undefined)
    await Promise.all([fetchAllAssets(), fetchRoleAssets(role.id), fetchProjects(), fetchEnvironments()])
    setAuthzLoading(false)
  }

  // 批量添加资产到授权列表
  const handleBatchAdd = () => {
    if (!selectedProjectId && !selectedEnvId) {
      message.warning('请先选择项目或环境')
      return
    }

    // 筛选符合条件的资产
    const filteredAssets = allAssets.filter((asset) => {
      if (selectedProjectId && asset.project_id !== selectedProjectId) return false
      if (selectedEnvId && asset.environment_id !== selectedEnvId) return false
      return true
    })

    if (filteredAssets.length === 0) {
      message.info('没有符合条件的资产')
      return
    }

    // 添加到 targetKeys（去重）
    const newKeys = new Set([...targetKeys, ...filteredAssets.map((a) => String(a.id))])
    const addedCount = newKeys.size - targetKeys.length
    setTargetKeys([...newKeys])

    if (addedCount > 0) {
      message.success(`已添加 ${addedCount} 个资产到授权列表`)
    } else {
      message.info('所选资产已在授权列表中')
    }
  }

  // 选择所有主机
  const handleSelectAll = () => {
    if (allAssets.length === 0) {
      message.info('没有可授权的资产')
      return
    }

    const allAssetKeys = allAssets.map((a) => String(a.id))
    const addedCount = allAssetKeys.length - targetKeys.length
    setTargetKeys(allAssetKeys)

    if (addedCount > 0) {
      message.success(`已添加所有 ${allAssets.length} 个资产到授权列表`)
    } else {
      message.info('所有资产已在授权列表中')
    }
  }

  // 保存资产授权
  const handleSaveAuthz = async () => {
    if (!authzRole) return

    setAuthzLoading(true)
    try {
      // 计算需要新增和撤销的资产
      const currentAssetIds = new Set(roleAssets.map((a) => String(a.id)))
      const newTargetIds = new Set(targetKeys)

      const toGrant = targetKeys.filter((id) => !currentAssetIds.has(id)).map(Number)
      const toRevoke = roleAssets
        .filter((a) => !newTargetIds.has(String(a.id)))
        .map((a) => a.id)

      // 执行授权和撤销
      if (toGrant.length > 0) {
        await roleApi.grantRoleAssets(authzRole.id, toGrant)
      }
      if (toRevoke.length > 0) {
        await roleApi.revokeRoleAssets(authzRole.id, toRevoke)
      }

      message.success('授权保存成功')
      setAuthzModalVisible(false)
      fetchRoles() // 刷新角色列表以更新资产数量
    } catch (error: any) {
      message.error(error.response?.data?.error || '授权保存失败')
    } finally {
      setAuthzLoading(false)
    }
  }

  // 获取项目名称
  const getProjectName = (projectId?: number) => {
    if (!projectId) return ''
    const project = projects.find((p) => p.id === projectId)
    return project?.name || ''
  }

  // 获取环境名称
  const getEnvName = (envId?: number) => {
    if (!envId) return ''
    const env = environments.find((e) => e.id === envId)
    return env?.name || ''
  }

  // Transfer 数据源
  const transferDataSource: TransferItem[] = allAssets.map((asset) => ({
    key: String(asset.id),
    title: asset.hostName,
    description: asset.ip,
    projectId: asset.project_id,
    environmentId: asset.environment_id,
  }))

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '管理员',
      dataIndex: 'is_admin',
      key: 'is_admin',
      width: 100,
      render: (isAdmin: boolean) =>
        isAdmin ? (
          <Tag color="red" icon={<CrownOutlined />}>
            管理员
          </Tag>
        ) : (
          '-'
        ),
    },
    {
      title: '授权资产',
      dataIndex: 'asset_count',
      key: 'asset_count',
      width: 120,
      render: (assetCount: number, record: Role) => {
        if (record.is_admin) {
          return <Tag color="green">全部</Tag>
        }
        return (
          <Tag color={assetCount > 0 ? 'blue' : 'default'}>
            {assetCount} 台
          </Tag>
        )
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: Role) => (
        <Space size="small">
          {!record.is_admin && (
            <Tooltip title="资产授权">
              <Button
                type="link"
                size="small"
                icon={<SafetyCertificateOutlined />}
                onClick={() => handleOpenAuthzModal(record)}
                aria-label={`资产授权 ${record.name}`}
              >
                授权
              </Button>
            </Tooltip>
          )}
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            aria-label={`编辑角色 ${record.name}`}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            aria-label={`删除角色 ${record.name}`}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

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
        <h2>角色管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增角色">
          新增角色
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={roles}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />

      {/* 创建/编辑角色弹窗 */}
      <Modal
        title={editingRole ? '编辑角色' : '新增角色'}
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
          <Form.Item name="name" label="角色名称" rules={[{ required: true, message: '请输入角色名称' }]}>
            <Input placeholder="角色名称" aria-label="角色名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={4} placeholder="角色描述" aria-label="角色描述" />
          </Form.Item>
          <Form.Item name="is_admin" valuePropName="checked">
            <Checkbox>
              <Space>
                <CrownOutlined style={{ color: '#f5222d' }} />
                设为管理员角色
              </Space>
            </Checkbox>
          </Form.Item>
          <div style={{ color: '#888', fontSize: 12, marginTop: -16, marginBottom: 16 }}>
            管理员角色可访问所有资产，无需单独授权
          </div>
        </Form>
      </Modal>

      {/* 资产授权弹窗 */}
      <Modal
        title={
          <Space>
            <SafetyCertificateOutlined />
            <span>资产授权 - {authzRole?.name}</span>
          </Space>
        }
        open={authzModalVisible}
        onCancel={() => setAuthzModalVisible(false)}
        onOk={handleSaveAuthz}
        confirmLoading={authzLoading}
        width={750}
        okText="保存"
        cancelText="取消"
      >
        {/* 快速授权区域 */}
        <div style={{ marginBottom: 16, padding: '12px 16px', background: '#f5f5f5', borderRadius: 6 }}>
          <div style={{ marginBottom: 8, fontWeight: 500 }}>
            <AppstoreAddOutlined style={{ marginRight: 8 }} />
            快速授权
          </div>
          <Space wrap>
            <Select
              placeholder="选择项目"
              style={{ width: 160 }}
              allowClear
              value={selectedProjectId}
              onChange={(value) => setSelectedProjectId(value)}
            >
              {projects.map((p) => (
                <Select.Option key={p.id} value={p.id}>
                  {p.name}
                </Select.Option>
              ))}
            </Select>
            <Select
              placeholder="选择环境"
              style={{ width: 160 }}
              allowClear
              value={selectedEnvId}
              onChange={(value) => setSelectedEnvId(value)}
            >
              {environments.map((e) => (
                <Select.Option key={e.id} value={e.id}>
                  {e.name}
                </Select.Option>
              ))}
            </Select>
            <Button type="primary" icon={<AppstoreAddOutlined />} onClick={handleBatchAdd}>
              添加到授权
            </Button>
            <Button onClick={handleSelectAll}>
              选择所有主机
            </Button>
          </Space>
          <div style={{ marginTop: 8, color: '#888', fontSize: 12 }}>
            选择项目和/或环境后点击"添加到授权"，或点击"选择所有主机"授权全部资产
          </div>
        </div>

        <Divider style={{ margin: '12px 0' }} />

        <div style={{ marginBottom: 12 }}>
          <Tag color="blue">拥有此角色的用户将可以访问右侧的资产</Tag>
        </div>
        <Transfer
          dataSource={transferDataSource}
          targetKeys={targetKeys}
          onChange={(newTargetKeys) => setTargetKeys(newTargetKeys as string[])}
          render={(item: any) => (
            <span>
              <span style={{ fontWeight: 500 }}>{item.title}</span>
              <Tag color="default" style={{ marginLeft: 4 }}>{item.description}</Tag>
              {item.projectId && (
                <Tag color="blue" style={{ marginLeft: 2 }}>{getProjectName(item.projectId)}</Tag>
              )}
              {item.environmentId && (
                <Tag color="green" style={{ marginLeft: 2 }}>{getEnvName(item.environmentId)}</Tag>
              )}
            </span>
          )}
          titles={['未授权资产', '已授权资产']}
          listStyle={{ width: 310, height: 350 }}
          showSearch
          filterOption={(inputValue, item: any) =>
            item.title!.toLowerCase().includes(inputValue.toLowerCase()) ||
            (item.description || '').toLowerCase().includes(inputValue.toLowerCase()) ||
            getProjectName(item.projectId).toLowerCase().includes(inputValue.toLowerCase()) ||
            getEnvName(item.environmentId).toLowerCase().includes(inputValue.toLowerCase())
          }
          locale={{
            searchPlaceholder: '搜索主机名、IP、项目或环境',
            itemUnit: '项',
            itemsUnit: '项',
            notFoundContent: '无数据',
          }}
        />
      </Modal>
    </div>
  )
}

export default RoleList
