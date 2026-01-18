// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Tag, Input, Select, Modal, Form, Descriptions, Spin, Divider, ColorPicker } from 'antd'
import { PlusOutlined, SearchOutlined, DownloadOutlined, UploadOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { assetApi, Asset, CreateAssetRequest, UpdateAssetRequest } from '@/api/asset'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import { cloudPlatformApi, CloudPlatform } from '@/api/cloudPlatform'
import { sshkeyApi, SSHKey } from '@/api/sshkey'
import { tagApi, Tag as TagType, CreateTagRequest } from '@/api/tag'

const { Search } = Input

const AssetList = () => {
  const [assets, setAssets] = useState<Asset[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [searchText, setSearchText] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>()
  const [modalVisible, setModalVisible] = useState(false)
  const [editingAsset, setEditingAsset] = useState<Asset | null>(null)
  const [form] = Form.useForm()
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [cloudPlatforms, setCloudPlatforms] = useState<CloudPlatform[]>([])
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([])
  const [tags, setTags] = useState<TagType[]>([])
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [detailAsset, setDetailAsset] = useState<Asset | null>(null)
  const [detailLoading, setDetailLoading] = useState(false)
  // 新建标签相关状态
  const [newTagModalVisible, setNewTagModalVisible] = useState(false)
  const [newTagName, setNewTagName] = useState('')
  const [newTagColor, setNewTagColor] = useState('#1890ff')
  const [creatingTag, setCreatingTag] = useState(false)

  useEffect(() => {
    fetchAssets()
    fetchProjects()
    fetchEnvironments()
    fetchCloudPlatforms()
    fetchSshKeys()
    fetchTags()
  }, [page, pageSize, searchText, statusFilter])

  const fetchAssets = async () => {
    setLoading(true)
    try {
      const response = await assetApi.list({
        page,
        page_size: pageSize,
        search: searchText || undefined,
        status: statusFilter,
      })
      // Debug: log the first asset to check field names
      if (response.data.data && response.data.data.length > 0) {
        console.log('First asset data:', response.data.data[0])
      }
      setAssets(response.data.data)
      setTotal(response.data.total)
    } catch (error: any) {
      if (error.response?.status === 403) {
        message.error('无权限访问资产列表')
      } else {
        message.error('获取资产列表失败')
      }
    } finally {
      setLoading(false)
    }
  }

  const fetchProjects = async () => {
    try {
      const response = await projectApi.list()
      setProjects(response.data)
    } catch (error: any) {
      // Ignore project fetch errors
    }
  }

  const fetchEnvironments = async () => {
    try {
      const response = await environmentApi.list()
      setEnvironments(response.data)
    } catch (error: any) {
      // Ignore environment fetch errors
    }
  }

  const fetchCloudPlatforms = async () => {
    try {
      const response = await cloudPlatformApi.list()
      setCloudPlatforms(response.data)
    } catch (error: any) {
      // Ignore cloud platform fetch errors
    }
  }

  const fetchSshKeys = async () => {
    try {
      const response = await sshkeyApi.list()
      setSshKeys(response.data.data || [])
    } catch (error: any) {
      // Ignore SSH key fetch errors
    }
  }

  const fetchTags = async () => {
    try {
      const response = await tagApi.list()
      setTags(response.data || [])
    } catch (error: any) {
      // Ignore tag fetch errors
    }
  }

  // 创建新标签
  const handleCreateTag = async () => {
    if (!newTagName.trim()) {
      message.warning('请输入标签名称')
      return
    }
    setCreatingTag(true)
    try {
      const tagData: CreateTagRequest = {
        name: newTagName.trim(),
        color: newTagColor,
      }
      const response = await tagApi.create(tagData)
      message.success('标签创建成功')
      
      // 刷新标签列表
      await fetchTags()
      
      // 获取当前选中的标签并添加新标签
      const currentTags = form.getFieldValue('tag_ids') || []
      const newTagId = response.data?.data?.id || response.data?.id
      if (newTagId) {
        form.setFieldsValue({ tag_ids: [...currentTags, newTagId] })
      }
      
      // 重置并关闭弹窗
      setNewTagName('')
      setNewTagColor('#1890ff')
      setNewTagModalVisible(false)
    } catch (error: any) {
      message.error(error.response?.data?.error || '创建标签失败')
    } finally {
      setCreatingTag(false)
    }
  }

  const handleCreate = () => {
    setEditingAsset(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (asset: Asset) => {
    setEditingAsset(asset)
    const tagIds = asset.tags?.map((tag) => tag.id) || []
    form.setFieldsValue({
      hostName: asset.hostName,
      project_id: asset.project_id,
      cloud_platform_id: asset.cloud_platform_id,
      environment_id: asset.environment_id,
      ip: asset.ip,
      status: asset.status || 'active',
      ssh_port: asset.ssh_port || 22,
      ssh_key_id: asset.ssh_key_id,
      ssh_user: asset.ssh_user,
      cpu: asset.cpu,
      memory: asset.memory,
      disk: asset.disk,
      description: asset.description,
      tag_ids: tagIds,
    })
    setModalVisible(true)
  }

  const handleViewDetail = async (asset: Asset) => {
    setDetailLoading(true)
    setDetailModalVisible(true)
    try {
      const response = await assetApi.get(asset.id)
      setDetailAsset(response.data)
    } catch (error: any) {
      message.error('获取资产详情失败')
      setDetailModalVisible(false)
    } finally {
      setDetailLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个资产吗？',
      onOk: async () => {
        try {
          await assetApi.delete(id)
          message.success('删除成功')
          fetchAssets()
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateAssetRequest | UpdateAssetRequest) => {
    try {
      // Convert numeric fields from string to number if present
      const processedValues: any = {
        ...values,
        // Convert ssh_port to number (required int in backend)
        ssh_port: values.ssh_port !== undefined && values.ssh_port !== null && values.ssh_port !== ''
          ? Number(values.ssh_port)
          : undefined,
        // Optional numeric fields - convert empty strings to undefined
        project_id: values.project_id && values.project_id !== '' ? Number(values.project_id) : undefined,
        cloud_platform_id: values.cloud_platform_id && values.cloud_platform_id !== '' ? Number(values.cloud_platform_id) : undefined,
        environment_id: values.environment_id && values.environment_id !== '' ? Number(values.environment_id) : undefined,
        ssh_key_id: values.ssh_key_id && values.ssh_key_id !== '' ? Number(values.ssh_key_id) : undefined,
        // Ensure tag_ids is an array of numbers if present
        tag_ids: values.tag_ids && Array.isArray(values.tag_ids) 
          ? values.tag_ids.map(id => Number(id))
          : undefined,
      }

      // Remove undefined values to avoid sending them
      Object.keys(processedValues).forEach(key => {
        if (processedValues[key] === undefined || processedValues[key] === '') {
          delete processedValues[key]
        }
      })

      if (editingAsset) {
        await assetApi.update(editingAsset.id, processedValues as UpdateAssetRequest)
        message.success('更新成功')
      } else {
        await assetApi.create(processedValues as CreateAssetRequest)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchAssets()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
    }
  }

  const handleImport = () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.csv,.xlsx,.xls'
    input.onchange = async (e: any) => {
      const file = e.target.files[0]
      if (!file) return

      try {
        await assetApi.import(file)
        message.success('导入成功')
        fetchAssets()
      } catch (error: any) {
        message.error(error.response?.data?.error || '导入失败')
      }
    }
    input.click()
  }

  const handleExport = async () => {
    try {
      const response = await assetApi.export()
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `assets_${new Date().getTime()}.csv`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (error: any) {
      message.error('导出失败')
    }
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '主机名',
      dataIndex: 'hostName',
      key: 'hostName',
    },
    {
      title: '项目',
      dataIndex: 'project_id',
      key: 'project_id',
      render: (projectId: number | undefined) => {
        if (!projectId) return '-'
        const proj = projects.find((p) => p.id === projectId)
        return proj ? (
          <Tag color="purple">{proj.name}</Tag>
        ) : (
          <Tag color="default">{projectId}</Tag>
        )
      },
    },
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
    },
    {
      title: '环境',
      dataIndex: 'environment_id',
      key: 'environment_id',
      render: (environmentId: number | undefined) => {
        if (!environmentId) return '-'
        const env = environments.find((e) => e.id === environmentId)
        return env ? (
          <Tag color="blue">{env.name}</Tag>
        ) : (
          <Tag color="default">{environmentId}</Tag>
        )
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status === 'active' ? '激活' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '标签',
      dataIndex: 'tags',
      key: 'tags',
      render: (tags: Array<{ id: number; name: string; color: string }>) => (
        <Space size="small">
          {tags?.map((tag) => (
            <Tag key={tag.id} color={tag.color}>
              {tag.name}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: any, record: Asset) => (
        <Space size="small">
          <Button type="link" size="small" onClick={() => handleViewDetail(record)} aria-label={`查看资产详情 ${record.hostName}`}>
            详情
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)} aria-label={`编辑资产 ${record.hostName}`}>
            编辑
          </Button>
          <Button type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)} aria-label={`删除资产 ${record.hostName}`}>
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
        <h2>资产管理</h2>
        <Space wrap>
          <Button icon={<UploadOutlined />} onClick={handleImport} aria-label="导入资产">
            导入
          </Button>
          <Button icon={<DownloadOutlined />} onClick={handleExport} aria-label="导出资产">
            导出
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增资产">
            新增资产
          </Button>
        </Space>
      </div>
      <div style={{ marginBottom: 16, display: 'flex', gap: 8, flexWrap: 'wrap' }}>
        <Search
          placeholder="搜索资产主机名、代码、IP"
          allowClear
          style={{ width: 300, minWidth: 200 }}
          onSearch={(value) => {
            setSearchText(value)
            setPage(1)
          }}
          aria-label="搜索资产"
        />
        <Select
          placeholder="状态筛选"
          allowClear
          style={{ width: 120, minWidth: 100 }}
          onChange={(value) => {
            setStatusFilter(value)
            setPage(1)
          }}
          aria-label="状态筛选"
        >
          <Select.Option value="active">激活</Select.Option>
          <Select.Option value="disabled">禁用</Select.Option>
        </Select>
      </div>
      <Table
        columns={columns}
        dataSource={assets}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          responsive: true,
          onChange: (page, pageSize) => {
            setPage(page)
            setPageSize(pageSize || 20)
          },
        }}
      />
      <Modal
        title={editingAsset ? '编辑资产' : '新增资产'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 800 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          {/* Section 1: Basic Information */}
          <Form.Item
            name="hostName"
            label="主机名"
            rules={[{ required: true, message: '请输入主机名' }]}
          >
            <Input placeholder="资产主机名" aria-label="主机名" />
          </Form.Item>

          {/* Section 2: Organization */}
          <Form.Item name="project_id" label="项目">
            <Select placeholder="选择项目" allowClear aria-label="项目">
              {projects.map((proj) => (
                <Select.Option key={proj.id} value={proj.id}>
                  {proj.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="cloud_platform_id" label="云平台">
            <Select placeholder="选择云平台" allowClear aria-label="云平台">
              {cloudPlatforms.map((platform) => (
                <Select.Option key={platform.id} value={platform.id}>
                  {platform.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="environment_id" label="环境">
            <Select placeholder="选择环境" allowClear aria-label="环境">
              {environments.map((env) => (
                <Select.Option key={env.id} value={env.id}>
                  {env.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          {/* Section 3: Network Information */}
          <Form.Item name="ip" label="IP地址">
            <Input placeholder="IP地址" aria-label="IP地址" />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue="active">
            <Select aria-label="状态">
              <Select.Option value="active">激活</Select.Option>
              <Select.Option value="disabled">禁用</Select.Option>
            </Select>
          </Form.Item>

          {/* Section 4: SSH Connection */}
          <Form.Item name="ssh_port" label="SSH端口" initialValue={22}>
            <Input type="number" placeholder="SSH端口（默认：22）" aria-label="SSH端口" />
          </Form.Item>
          <Form.Item name="ssh_key_id" label="SSH密钥">
            <Select placeholder="选择SSH密钥" allowClear aria-label="SSH密钥">
              {sshKeys.map((key) => (
                <Select.Option key={key.id} value={key.id}>
                  {key.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="ssh_user" label="SSH用户">
            <Input placeholder="SSH用户名" aria-label="SSH用户" />
          </Form.Item>

          {/* Section 5: Hardware Specifications */}
          <Form.Item name="cpu" label="CPU">
            <Input placeholder="CPU信息（可自动收集）" aria-label="CPU" />
          </Form.Item>
          <Form.Item name="memory" label="内存">
            <Input placeholder="内存信息（可自动收集）" aria-label="内存" />
          </Form.Item>
          <Form.Item name="disk" label="磁盘">
            <Input placeholder="磁盘信息（可自动收集）" aria-label="磁盘" />
          </Form.Item>

          {/* Section 6: Additional Information */}
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} placeholder="资产描述" aria-label="描述" />
          </Form.Item>
          <Form.Item name="tag_ids" label="标签">
            <Select 
              mode="multiple" 
              placeholder="选择标签" 
              allowClear 
              aria-label="标签"
              dropdownRender={(menu) => (
                <>
                  {menu}
                  <Divider style={{ margin: '8px 0' }} />
                  <Button 
                    type="link" 
                    icon={<PlusOutlined />}
                    onClick={() => setNewTagModalVisible(true)}
                    style={{ width: '100%', textAlign: 'left' }}
                  >
                    新建标签
                  </Button>
                </>
              )}
            >
              {tags.map((tag) => (
                <Select.Option key={tag.id} value={tag.id}>
                  <Tag color={tag.color}>{tag.name}</Tag>
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="资产详情"
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false)
          setDetailAsset(null)
        }}
        footer={[
          <Button key="close" onClick={() => {
            setDetailModalVisible(false)
            setDetailAsset(null)
          }}>
            关闭
          </Button>,
        ]}
        width="90%"
        style={{ maxWidth: 800 }}
      >
        {detailLoading ? (
          <div style={{ textAlign: 'center', padding: '50px' }}>
            <Spin size="large" />
          </div>
        ) : detailAsset ? (
          <Descriptions column={2} bordered>
            <Descriptions.Item label="ID">{detailAsset.id}</Descriptions.Item>
            <Descriptions.Item label="主机名">{detailAsset.hostName}</Descriptions.Item>
            <Descriptions.Item label="状态">
              <Tag color={detailAsset.status === 'active' ? 'green' : 'default'}>
                {detailAsset.status === 'active' ? '激活' : '禁用'}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="项目">
              {detailAsset.project_id ? (
                (() => {
                  const proj = projects.find((p) => p.id === detailAsset.project_id)
                  return proj ? (
                    <Tag color="purple">{proj.name}</Tag>
                  ) : (
                    <Tag color="default">{detailAsset.project_id}</Tag>
                  )
                })()
              ) : (
                '-'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="云平台">
              {detailAsset.cloud_platform ? (
                <Tag color="cyan">{detailAsset.cloud_platform.name}</Tag>
              ) : detailAsset.cloud_platform_id ? (
                (() => {
                  const platform = cloudPlatforms.find((p) => p.id === detailAsset.cloud_platform_id)
                  return platform ? (
                    <Tag color="cyan">{platform.name}</Tag>
                  ) : (
                    <Tag color="default">{detailAsset.cloud_platform_id}</Tag>
                  )
                })()
              ) : (
                '-'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="环境">
              {detailAsset.environment_id ? (
                (() => {
                  const env = environments.find((e) => e.id === detailAsset.environment_id)
                  return env ? (
                    <Tag color="blue">{env.name}</Tag>
                  ) : (
                    <Tag color="default">{detailAsset.environment_id}</Tag>
                  )
                })()
              ) : (
                '-'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="IP地址">{detailAsset.ip || '-'}</Descriptions.Item>
            <Descriptions.Item label="SSH端口">{detailAsset.ssh_port || 22}</Descriptions.Item>
            <Descriptions.Item label="SSH密钥">
              {detailAsset.ssh_key_id ? (
                (() => {
                  const sshKey = sshKeys.find((k) => k.id === detailAsset.ssh_key_id)
                  return sshKey ? (
                    <Tag color="orange">{sshKey.name}</Tag>
                  ) : (
                    <Tag color="default">{detailAsset.ssh_key_id}</Tag>
                  )
                })()
              ) : (
                '-'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="SSH用户">{detailAsset.ssh_user || '-'}</Descriptions.Item>
            <Descriptions.Item label="CPU">{detailAsset.cpu || '-'}</Descriptions.Item>
            <Descriptions.Item label="内存">{detailAsset.memory || '-'}</Descriptions.Item>
            <Descriptions.Item label="磁盘">{detailAsset.disk || '-'}</Descriptions.Item>
            <Descriptions.Item label="标签" span={2}>
              <Space size="small">
                {detailAsset.tags?.map((tag) => (
                  <Tag key={tag.id} color={tag.color}>
                    {tag.name}
                  </Tag>
                ))}
                {(!detailAsset.tags || detailAsset.tags.length === 0) && '-'}
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="描述" span={2}>
              {detailAsset.description || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">{detailAsset.created_at}</Descriptions.Item>
            <Descriptions.Item label="更新时间">{detailAsset.updated_at}</Descriptions.Item>
          </Descriptions>
        ) : null}
      </Modal>

      {/* 新建标签弹窗 */}
      <Modal
        title="新建标签"
        open={newTagModalVisible}
        onCancel={() => {
          setNewTagModalVisible(false)
          setNewTagName('')
          setNewTagColor('#1890ff')
        }}
        onOk={handleCreateTag}
        confirmLoading={creatingTag}
        width={400}
        okText="创建"
        cancelText="取消"
      >
        <Form layout="vertical">
          <Form.Item label="标签名称" required>
            <Input
              placeholder="输入标签名称"
              value={newTagName}
              onChange={(e) => setNewTagName(e.target.value)}
              onPressEnter={handleCreateTag}
            />
          </Form.Item>
          <Form.Item label="标签颜色">
            <Space>
              <ColorPicker
                value={newTagColor}
                onChange={(color) => setNewTagColor(color.toHexString())}
              />
              <Tag color={newTagColor}>{newTagName || '预览'}</Tag>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AssetList
