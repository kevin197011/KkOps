// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, Button, Space, message, Spin } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { assetApi, Asset } from '@/api/asset'
import { projectApi, Project } from '@/api/project'
import { environmentApi, Environment } from '@/api/environment'
import { cloudPlatformApi, CloudPlatform } from '@/api/cloudPlatform'
import { sshkeyApi, SSHKey } from '@/api/sshkey'

const AssetDetail = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [asset, setAsset] = useState<Asset | null>(null)
  const [loading, setLoading] = useState(false)
  const [projects, setProjects] = useState<Project[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [cloudPlatforms, setCloudPlatforms] = useState<CloudPlatform[]>([])
  const [sshKeys, setSshKeys] = useState<SSHKey[]>([])

  useEffect(() => {
    if (!id) {
      navigate('/assets')
      return
    }
    const assetId = parseInt(id, 10)
    if (isNaN(assetId) || assetId <= 0) {
      message.error('无效的资产ID')
      navigate('/assets')
      return
    }
    fetchAsset(assetId)
    fetchProjects()
    fetchEnvironments()
    fetchCloudPlatforms()
    fetchSshKeys()
  }, [id, navigate])

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

  const fetchAsset = async (assetId: number) => {
    setLoading(true)
    try {
      const response = await assetApi.get(assetId)
      setAsset(response.data)
    } catch (error: any) {
      message.error('获取资产详情失败')
      navigate('/assets')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!asset) {
    return null
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/assets')}>
          返回列表
        </Button>
      </div>
      <Card title="资产详情">
        <Descriptions column={2} bordered>
          <Descriptions.Item label="ID">{asset.id}</Descriptions.Item>
          <Descriptions.Item label="主机名">{asset.hostName}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={asset.status === 'active' ? 'green' : 'default'}>
              {asset.status === 'active' ? '激活' : '禁用'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="项目">
            {asset.project_id ? (
              (() => {
                const proj = projects.find((p) => p.id === asset.project_id)
                return proj ? (
                  <Tag color="purple">{proj.name}</Tag>
                ) : (
                  <Tag color="default">{asset.project_id}</Tag>
                )
              })()
            ) : (
              '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label="云平台">
            {asset.cloud_platform ? (
              <Tag color="cyan">{asset.cloud_platform.name}</Tag>
            ) : asset.cloud_platform_id ? (
              (() => {
                const platform = cloudPlatforms.find((p) => p.id === asset.cloud_platform_id)
                return platform ? (
                  <Tag color="cyan">{platform.name}</Tag>
                ) : (
                  <Tag color="default">{asset.cloud_platform_id}</Tag>
                )
              })()
            ) : (
              '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label="环境">
            {asset.environment_id ? (
              (() => {
                const env = environments.find((e) => e.id === asset.environment_id)
                return env ? (
                  <Tag color="blue">{env.name}</Tag>
                ) : (
                  <Tag color="default">{asset.environment_id}</Tag>
                )
              })()
            ) : (
              '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label="IP地址">{asset.ip || '-'}</Descriptions.Item>
          <Descriptions.Item label="SSH端口">{asset.ssh_port || 22}</Descriptions.Item>
          <Descriptions.Item label="SSH密钥">
            {asset.ssh_key_id ? (
              (() => {
                const sshKey = sshKeys.find((k) => k.id === asset.ssh_key_id)
                return sshKey ? (
                  <Tag color="orange">{sshKey.name}</Tag>
                ) : (
                  <Tag color="default">{asset.ssh_key_id}</Tag>
                )
              })()
            ) : (
              '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label="SSH用户">{asset.ssh_user || '-'}</Descriptions.Item>
          <Descriptions.Item label="CPU">{asset.cpu || '-'}</Descriptions.Item>
          <Descriptions.Item label="内存">{asset.memory || '-'}</Descriptions.Item>
          <Descriptions.Item label="磁盘">{asset.disk || '-'}</Descriptions.Item>
          <Descriptions.Item label="标签" span={2}>
            <Space size="small">
              {asset.tags?.map((tag) => (
                <Tag key={tag.id} color={tag.color}>
                  {tag.name}
                </Tag>
              ))}
              {(!asset.tags || asset.tags.length === 0) && '-'}
            </Space>
          </Descriptions.Item>
          <Descriptions.Item label="描述" span={2}>
            {asset.description || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">{asset.created_at}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{asset.updated_at}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}

export default AssetDetail
