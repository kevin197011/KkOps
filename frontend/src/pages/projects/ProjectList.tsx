// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { projectApi, Project, CreateProjectRequest, UpdateProjectRequest } from '@/api/project'

const ProjectList = () => {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchProjects()
  }, [])

  const fetchProjects = async () => {
    setLoading(true)
    try {
      const response = await projectApi.list()
      setProjects(response.data)
    } catch (error: any) {
      message.error('获取项目列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingProject(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (project: Project) => {
    setEditingProject(project)
    form.setFieldsValue(project)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个项目吗？',
      onOk: async () => {
        try {
          await projectApi.delete(id)
          message.success('删除成功')
          fetchProjects()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async (values: CreateProjectRequest | UpdateProjectRequest) => {
    try {
      if (editingProject) {
        await projectApi.update(editingProject.id, values)
        message.success('更新成功')
      } else {
        await projectApi.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      form.resetFields()
      fetchProjects()
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败')
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
      title: '项目名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Project) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ 
        marginBottom: 16, 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: 8
      }}>
        <h2>项目管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} aria-label="新增项目">
          新增项目
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={projects}
        loading={loading}
        rowKey="id"
        scroll={{ x: 'max-content' }}
      />
      <Modal
        title={editingProject ? '编辑项目' : '新增项目'}
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
            label="项目名称"
            rules={[{ required: true, message: '请输入项目名称' }]}
          >
            <Input placeholder="项目名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} placeholder="项目描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ProjectList
