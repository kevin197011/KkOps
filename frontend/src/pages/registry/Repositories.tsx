// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useCallback, useEffect, useState } from 'react'
import { Button, Card, Drawer, Form, Input, Select, Space, Table, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { ContainerOutlined } from '@ant-design/icons'
import { PageContainer, PageHeader, EmptyState } from '@/components/shell'
import { integrationsApi, IntegrationPublic } from '@/api/integration'
import { registryApi, RepositorySummary, TagSummary } from '@/api/registry'

const Repositories = () => {
  const [integrations, setIntegrations] = useState<IntegrationPublic[]>([])
  const [integrationId, setIntegrationId] = useState<number | undefined>()
  const [repos, setRepos] = useState<RepositorySummary[]>([])
  const [loading, setLoading] = useState(false)
  const [tagsOpen, setTagsOpen] = useState(false)
  const [tags, setTags] = useState<TagSummary[]>([])
  const [tagRepo, setTagRepo] = useState('')
  const [tagLoading, setTagLoading] = useState(false)
  const [vulnOpen, setVulnOpen] = useState(false)
  const [vulnText, setVulnText] = useState('')
  const [form] = Form.useForm()

  const loadIntegrations = useCallback(async () => {
    try {
      const res = await integrationsApi.list()
      const list = (res.data.data || []).filter((i) => i.kind.toLowerCase() === 'harbor')
      setIntegrations(list)
      setIntegrationId((prev) => prev ?? list[0]?.id)
    } catch {
      message.error('加载集成失败')
    }
  }, [])

  const loadRepos = useCallback(async () => {
    if (!integrationId) {
      setRepos([])
      return
    }
    setLoading(true)
    try {
      const res = await registryApi.listRepositories(integrationId)
      setRepos(res.data.data || [])
    } catch (e: unknown) {
      const ax = e as { response?: { data?: { error?: string } } }
      message.error(ax.response?.data?.error || '加载仓库失败')
    } finally {
      setLoading(false)
    }
  }, [integrationId])

  useEffect(() => {
    loadIntegrations()
  }, [loadIntegrations])

  useEffect(() => {
    loadRepos()
  }, [loadRepos])

  const openTags = async (repoFullName: string) => {
    if (!integrationId) return
    setTagRepo(repoFullName)
    setTagsOpen(true)
    setTagLoading(true)
    try {
      const res = await registryApi.listTags(integrationId, repoFullName)
      setTags(res.data.data || [])
    } catch {
      message.error('加载标签失败')
      setTags([])
    } finally {
      setTagLoading(false)
    }
  }

  const openVuln = async (repoFullName: string, tagName: string) => {
    if (!integrationId) return
    try {
      const res = await registryApi.vulnerabilities(integrationId, repoFullName, tagName)
      setVulnText(res.data.data?.raw_json || '')
      setVulnOpen(true)
    } catch {
      message.error('加载漏洞信息失败')
    }
  }

  const columns: ColumnsType<RepositorySummary> = [
    { title: '仓库', dataIndex: 'name', key: 'name' },
    { title: '拉取次数', dataIndex: 'pull_count', key: 'pull_count' },
    { title: '更新时间', dataIndex: 'update_time', key: 'update_time' },
    {
      title: '操作',
      key: 'act',
      render: (_, r) => (
        <Button type="link" size="small" onClick={() => openTags(r.name)}>
          标签
        </Button>
      ),
    },
  ]

  const tagCols: ColumnsType<TagSummary> = [
    { title: '标签', dataIndex: 'name', key: 'name' },
    { title: 'Digest', dataIndex: 'digest', key: 'digest', ellipsis: true },
    {
      title: '漏洞',
      key: 'v',
      render: (_, t) => (
        <Button type="link" size="small" onClick={() => openVuln(tagRepo, t.name)}>
          查看
        </Button>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader
        title={
          <>
            <ContainerOutlined style={{ marginRight: 8 }} />
            镜像仓库
          </>
        }
        description="浏览 Harbor 项目中的镜像与标签（需集成类型为 harbor）。"
      />
      <Card style={{ marginBottom: 16 }}>
        <Space wrap>
          <span>集成：</span>
          <Select
            style={{ minWidth: 260 }}
            value={integrationId}
            onChange={setIntegrationId}
            options={integrations.map((i) => ({
              value: i.id,
              label: `${i.name} (${i.kind})`,
            }))}
          />
          <Button type="primary" onClick={loadRepos} loading={loading}>
            刷新
          </Button>
        </Space>
        <Form form={form} layout="inline" style={{ marginTop: 16 }}>
          <Form.Item name="project_name" label="项目过滤">
            <Input placeholder="可选 Harbor 项目名" style={{ width: 200 }} />
          </Form.Item>
          <Button
            onClick={async () => {
              const v = form.getFieldsValue()
              if (!integrationId) return
              setLoading(true)
              try {
                const res = await registryApi.listRepositories(integrationId, {
                  project_name: v.project_name,
                })
                setRepos(res.data.data || [])
              } catch {
                message.error('查询失败')
              } finally {
                setLoading(false)
              }
            }}
          >
            应用过滤
          </Button>
        </Form>
      </Card>

      {integrations.length === 0 ? (
        <EmptyState description="未配置 Harbor。请在集成中心添加 Harbor 集成。" />
      ) : (
        <Card title="镜像仓库">
          <Table rowKey="name" loading={loading} columns={columns} dataSource={repos} pagination={{ pageSize: 20 }} />
        </Card>
      )}

      <Drawer title={`标签 — ${tagRepo}`} width={720} open={tagsOpen} onClose={() => setTagsOpen(false)}>
        <Table rowKey="name" loading={tagLoading} columns={tagCols} dataSource={tags} pagination={{ pageSize: 15 }} />
      </Drawer>

      <Drawer title="漏洞扫描（原始 JSON）" width={640} open={vulnOpen} onClose={() => setVulnOpen(false)}>
        <pre style={{ fontSize: 11, whiteSpace: 'pre-wrap' }}>{vulnText}</pre>
      </Drawer>
    </PageContainer>
  )
}

export default Repositories
