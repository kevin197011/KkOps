// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect, useMemo } from 'react'
import { Card, Select, Input, Checkbox, Space, Typography, Button, Row, Col, Tag } from 'antd'
import {
  DatabaseOutlined,
  CheckSquareOutlined,
  BorderOutlined,
  ClearOutlined,
  FolderOutlined,
  GlobalOutlined,
} from '@ant-design/icons'
import { assetApi, Asset } from '@/api/asset'
import { Project } from '@/api/project'
import { Environment } from '@/api/environment'
import { useThemeStore } from '@/stores/theme'

const { Text } = Typography
const { Search } = Input

interface HostSelectorProps {
  selectedAssetIds: number[]
  onSelectionChange: (assetIds: number[]) => void
  projects: Project[]
  environments: Environment[]
}

const HostSelector: React.FC<HostSelectorProps> = ({
  selectedAssetIds,
  onSelectionChange,
  projects,
  environments,
}) => {
  const { mode } = useThemeStore()
  const [assets, setAssets] = useState<Asset[]>([])
  const [loading, setLoading] = useState(false)
  const [projectFilter, setProjectFilter] = useState<number | null>(null)
  const [environmentFilter, setEnvironmentFilter] = useState<number | null>(null)
  const [searchText, setSearchText] = useState('')

  // 获取资产列表
  useEffect(() => {
    fetchAssets()
  }, [])

  const fetchAssets = async () => {
    setLoading(true)
    try {
      const response = await assetApi.list({ page: 1, page_size: 1000, status: 'active' })
      const assetsData = response.data?.data || []
      setAssets(Array.isArray(assetsData) ? assetsData : [])
    } catch (error) {
      console.error('获取资产列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  // 过滤后的资产列表
  const filteredAssets = useMemo(() => {
    let result = assets

    // 按项目筛选
    if (projectFilter) {
      result = result.filter((asset) => asset.project_id === projectFilter)
    }

    // 按环境筛选
    if (environmentFilter) {
      result = result.filter((asset) => asset.environment_id === environmentFilter)
    }

    // 按搜索关键词筛选
    if (searchText.trim()) {
      const keyword = searchText.toLowerCase()
      result = result.filter(
        (asset) =>
          asset.hostName?.toLowerCase().includes(keyword) ||
          asset.ip?.toLowerCase().includes(keyword)
      )
    }

    return result
  }, [assets, projectFilter, environmentFilter, searchText])

  // 获取项目名称
  const getProjectName = (projectId: number | null | undefined) => {
    if (!projectId) return null
    const project = projects.find((p) => p.id === projectId)
    return project?.name || null
  }

  // 获取环境名称
  const getEnvironmentName = (environmentId: number | null | undefined) => {
    if (!environmentId) return null
    const env = environments.find((e) => e.id === environmentId)
    return env?.name || null
  }

  // 批量选择处理
  const handleSelectAll = () => {
    onSelectionChange(filteredAssets.map((asset) => asset.id))
  }

  const handleSelectByProject = () => {
    if (!projectFilter) return
    const projectAssets = assets.filter((asset) => asset.project_id === projectFilter)
    const newSelection = Array.from(
      new Set([...selectedAssetIds, ...projectAssets.map((a) => a.id)])
    )
    onSelectionChange(newSelection)
  }

  const handleSelectByEnvironment = () => {
    if (!environmentFilter) return
    const envAssets = assets.filter((asset) => asset.environment_id === environmentFilter)
    const newSelection = Array.from(
      new Set([...selectedAssetIds, ...envAssets.map((a) => a.id)])
    )
    onSelectionChange(newSelection)
  }

  const handleClear = () => {
    onSelectionChange([])
  }

  // 单选处理
  const handleAssetToggle = (assetId: number, checked: boolean) => {
    if (checked) {
      onSelectionChange([...selectedAssetIds, assetId])
    } else {
      onSelectionChange(selectedAssetIds.filter((id) => id !== assetId))
    }
  }

  return (
    <Card
      style={{
        marginBottom: 16,
        borderRadius: 12,
      }}
      styles={{
        body: { padding: '20px 24px' },
      }}
    >
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Text strong>
            <DatabaseOutlined style={{ marginRight: 8 }} />
            选择执行主机
          </Text>
          {selectedAssetIds.length > 0 && (
            <Tag color="blue">已选择 {selectedAssetIds.length} 台主机</Tag>
          )}
        </div>

        {/* 筛选器 */}
        <Row gutter={[12, 12]}>
          <Col xs={24} sm={12} md={8}>
            <Select
              style={{ width: '100%' }}
              placeholder="按项目筛选"
              allowClear
              value={projectFilter}
              onChange={setProjectFilter}
              options={projects.map((p) => ({
                value: p.id,
                label: p.name,
              }))}
            />
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Select
              style={{ width: '100%' }}
              placeholder="按环境筛选"
              allowClear
              value={environmentFilter}
              onChange={setEnvironmentFilter}
              options={environments.map((e) => ({
                value: e.id,
                label: e.name,
              }))}
            />
          </Col>
          <Col xs={24} sm={24} md={8}>
            <Search
              placeholder="搜索主机名或 IP"
              allowClear
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
            />
          </Col>
        </Row>

        {/* 批量选择按钮 */}
        <Space wrap>
          <Button
            size="small"
            icon={<CheckSquareOutlined />}
            onClick={handleSelectAll}
            disabled={filteredAssets.length === 0}
          >
            全选 ({filteredAssets.length})
          </Button>
          <Button
            size="small"
            icon={<FolderOutlined />}
            onClick={handleSelectByProject}
            disabled={!projectFilter}
          >
            按项目选
          </Button>
          <Button
            size="small"
            icon={<GlobalOutlined />}
            onClick={handleSelectByEnvironment}
            disabled={!environmentFilter}
          >
            按环境选
          </Button>
          <Button size="small" icon={<ClearOutlined />} onClick={handleClear} danger>
            清空
          </Button>
        </Space>

        {/* 主机列表 */}
        <div
          style={{
            maxHeight: 400,
            overflowY: 'auto',
            border: `1px solid ${mode === 'dark' ? '#334155' : '#E2E8F0'}`,
            borderRadius: 8,
            padding: 12,
          }}
        >
          {loading ? (
            <div style={{ textAlign: 'center', padding: 40 }}>
              <Text type="secondary">加载中...</Text>
            </div>
          ) : filteredAssets.length === 0 ? (
            <div style={{ textAlign: 'center', padding: 40 }}>
              <Text type="secondary">
                {searchText || projectFilter || environmentFilter
                  ? '没有匹配的主机'
                  : '暂无主机'}
              </Text>
            </div>
          ) : (
            <Checkbox.Group
              value={selectedAssetIds}
              style={{ width: '100%' }}
              onChange={(values) => {
                onSelectionChange(values as number[])
              }}
            >
              <Space direction="vertical" size={8} style={{ width: '100%' }}>
                {filteredAssets.map((asset) => {
                  const isSelected = selectedAssetIds.includes(asset.id)
                  const projectName = getProjectName(asset.project_id)
                  const envName = getEnvironmentName(asset.environment_id)
                  return (
                    <div
                      key={asset.id}
                      style={{
                        padding: '8px 12px',
                        borderRadius: 6,
                        background: isSelected
                          ? mode === 'dark'
                            ? 'rgba(59, 130, 246, 0.15)'
                            : 'rgba(59, 130, 246, 0.1)'
                          : 'transparent',
                        border: `1px solid ${
                          isSelected
                            ? '#3B82F6'
                            : mode === 'dark'
                              ? '#334155'
                              : '#E2E8F0'
                        }`,
                        transition: 'all 0.2s',
                      }}
                    >
                      <Checkbox
                        value={asset.id}
                        style={{
                          width: '100%',
                        }}
                      >
                        <Space style={{ width: '100%', justifyContent: 'space-between' }}>
                          <div>
                            <Text strong>{asset.hostName}</Text>
                            <Text
                              type="secondary"
                              style={{
                                marginLeft: 8,
                                fontFamily: '"JetBrains Mono", monospace',
                                fontSize: 12,
                              }}
                            >
                              {asset.ip}
                            </Text>
                            {(projectName || envName) && (
                              <Space size={4} style={{ marginLeft: 8 }}>
                                {projectName && (
                                  <Tag
                                    icon={<FolderOutlined />}
                                    style={{
                                      margin: 0,
                                      fontSize: 11,
                                    }}
                                  >
                                    {projectName}
                                  </Tag>
                                )}
                                {envName && (
                                  <Tag
                                    icon={<GlobalOutlined />}
                                    style={{
                                      margin: 0,
                                      fontSize: 11,
                                    }}
                                  >
                                    {envName}
                                  </Tag>
                                )}
                              </Space>
                            )}
                          </div>
                        </Space>
                      </Checkbox>
                    </div>
                  )
                })}
              </Space>
            </Checkbox.Group>
          )}
        </div>
      </Space>
    </Card>
  )
}

export default HostSelector
