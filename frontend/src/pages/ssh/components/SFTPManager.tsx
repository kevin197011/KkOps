// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import React, { useState, useEffect, useRef } from 'react'
import {
  Button,
  Space,
  Table,
  App,
  Modal,
  Input,
  Upload,
  Progress,
  Breadcrumb,
  Popconfirm,
  Tooltip,
  Typography,
  Spin,
  Empty,
} from 'antd'
import {
  UploadOutlined,
  DownloadOutlined,
  DeleteOutlined,
  FolderOutlined,
  FileOutlined,
  ReloadOutlined,
  HomeOutlined,
  FolderAddOutlined,
  ArrowLeftOutlined,
  FolderOpenOutlined,
} from '@ant-design/icons'
import { useThemeStore } from '@/stores/theme'

const { Text } = Typography

interface RemoteFile {
  name: string
  size: number
  mode: string
  modTime: number
  isDir: boolean
}

interface SFTPManagerProps {
  ws: WebSocket | null
  connectionId: string
  assetName?: string
  visible: boolean
  onClose: () => void
}

const SFTPManager: React.FC<SFTPManagerProps> = ({ ws, connectionId, assetName, visible, onClose }) => {
  const { message } = App.useApp()
  const { mode } = useThemeStore()
  const isDark = mode === 'dark'
  const [currentPath, setCurrentPath] = useState<string>('.')
  const [files, setFiles] = useState<RemoteFile[]>([])
  const [loading, setLoading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState<Record<string, { progress: number; bytesTransferred: number; totalBytes: number }>>({})
  const [downloadProgress, setDownloadProgress] = useState<Record<string, { progress: number; bytesTransferred: number; totalBytes: number }>>({})
  const uploadFileRef = useRef<HTMLInputElement>(null)
  const downloadBufferRef = useRef<Map<string, Blob[]>>(new Map())

  const listFiles = React.useCallback((path: string) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      message.error('WebSocket 连接未建立')
      return
    }

    setLoading(true)
    ws.send(JSON.stringify({
      type: 'sftp_list',
      data: { path },
    }))
  }, [ws])

  useEffect(() => {
    if (visible && ws && ws.readyState === WebSocket.OPEN) {
      listFiles(currentPath)
    }
  }, [visible, currentPath, ws, listFiles])

  const handleSFTPMessage = React.useCallback((msg: any) => {
    switch (msg.type) {
      case 'sftp_list':
        if (msg.data && msg.data.path === currentPath) {
          setFiles(msg.data.files || [])
          setLoading(false)
        }
        break
      case 'sftp_error':
        message.error(msg.data?.error || 'SFTP 操作失败')
        setLoading(false)
        break
      case 'sftp_upload_progress':
        if (msg.data?.file_name) {
          const bytesTransferred = msg.data.bytes_transferred || 0
          const totalBytes = msg.data.total_bytes || 0
          const progress = totalBytes > 0
            ? Math.round((bytesTransferred / totalBytes) * 100)
            : 0
          setUploadProgress((prev) => ({
            ...prev,
            [msg.data.file_name]: {
              progress,
              bytesTransferred,
              totalBytes,
            },
          }))
        }
        break
      case 'sftp_upload_complete':
        message.success(`文件 ${msg.data?.file_name} 上传成功`)
        setUploadProgress((prev) => {
          const newProgress = { ...prev }
          delete newProgress[msg.data?.file_name]
          return newProgress
        })
        listFiles(currentPath)
        break
      case 'sftp_download_start':
        if (msg.data?.file_name) {
          downloadBufferRef.current.set(msg.data.file_name, [])
          setDownloadProgress((prev) => ({
            ...prev,
            [msg.data.file_name]: {
              progress: 0,
              bytesTransferred: 0,
              totalBytes: msg.data.size || 0,
            },
          }))
        }
        break
      case 'sftp_download_progress':
        if (msg.data?.file_name) {
          const bytesTransferred = msg.data.bytes_transferred || 0
          const totalBytes = msg.data.total_bytes || 0
          const progress = totalBytes > 0
            ? Math.round((bytesTransferred / totalBytes) * 100)
            : 0
          setDownloadProgress((prev) => ({
            ...prev,
            [msg.data.file_name]: {
              progress,
              bytesTransferred,
              totalBytes,
            },
          }))
        }
        break
      case 'sftp_download_complete':
        if (msg.data?.file_name) {
          const chunks = downloadBufferRef.current.get(msg.data.file_name) || []
          const blob = new Blob(chunks)
          const url = window.URL.createObjectURL(blob)
          const link = document.createElement('a')
          link.href = url
          link.download = msg.data.file_name
          document.body.appendChild(link)
          link.click()
          document.body.removeChild(link)
          window.URL.revokeObjectURL(url)
          downloadBufferRef.current.delete(msg.data.file_name)
          setDownloadProgress((prev) => {
            const newProgress = { ...prev }
            delete newProgress[msg.data.file_name]
            return newProgress
          })
          message.success(`文件 ${msg.data.file_name} 下载成功`)
        }
        break
      case 'sftp_delete_complete':
        message.success('删除成功')
        listFiles(currentPath)
        break
      case 'sftp_mkdir_complete':
        message.success('目录创建成功')
        listFiles(currentPath)
        break
      case 'sftp_rename_complete':
        message.success('重命名成功')
        listFiles(currentPath)
        break
    }
  }, [currentPath, listFiles])

  const handleBinaryMessage = React.useCallback((data: Blob) => {
    // Only handle binary messages if there's an active SFTP download
    // Check if there's a file being downloaded
    const downloadingFiles = Object.keys(downloadProgress)
    if (downloadingFiles.length === 0) {
      return // No active SFTP download, ignore binary data (might be ZMODEM)
    }

    // Find the file being downloaded by checking progress
    const downloadingFile = downloadingFiles.find(
      (fileName) => downloadProgress[fileName] !== undefined && downloadProgress[fileName].progress < 100
    )
    
    if (downloadingFile) {
      data.arrayBuffer().then((buffer) => {
        const chunks = downloadBufferRef.current.get(downloadingFile) || []
        chunks.push(new Blob([buffer]))
        downloadBufferRef.current.set(downloadingFile, chunks)
        
        // Update progress based on accumulated chunks
        const totalSize = chunks.reduce((sum, chunk) => sum + chunk.size, 0)
        const progressInfo = downloadProgress[downloadingFile]
        if (progressInfo) {
          const newProgress = progressInfo.totalBytes > 0
            ? Math.round((totalSize / progressInfo.totalBytes) * 100)
            : 0
          setDownloadProgress((prev) => ({
            ...prev,
            [downloadingFile]: {
              ...progressInfo,
              progress: newProgress,
              bytesTransferred: totalSize,
            },
          }))
        }
      })
    }
  }, [downloadProgress])

  useEffect(() => {
    if (!ws || !visible) return

    const handleMessage = (event: MessageEvent) => {
      try {
        if (event.data instanceof Blob) {
          // Handle binary data (download)
          handleBinaryMessage(event.data)
          return
        }

        const msg = JSON.parse(event.data)
        // Only handle SFTP-related messages
        if (msg.type && msg.type.startsWith('sftp_')) {
          handleSFTPMessage(msg)
        }
      } catch (error) {
        console.error('Error handling SFTP message:', error)
      }
    }

    ws.addEventListener('message', handleMessage)
    return () => {
      ws.removeEventListener('message', handleMessage)
    }
  }, [ws, visible, handleSFTPMessage, handleBinaryMessage])

  const handleFileClick = (file: RemoteFile) => {
    if (file.isDir) {
      let newPath: string
      if (currentPath === '.' || currentPath === '/') {
        newPath = currentPath === '/' ? `/${file.name}` : file.name
      } else {
        newPath = `${currentPath}/${file.name}`
      }
      setCurrentPath(newPath)
    }
  }

  const handleUpload = (file: File) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      message.error('WebSocket 连接未建立')
      return false
    }

    // Initialize upload progress
    setUploadProgress((prev) => ({
      ...prev,
      [file.name]: {
        progress: 0,
        bytesTransferred: 0,
        totalBytes: file.size,
      },
    }))

    const reader = new FileReader()
    reader.onload = (e) => {
      const arrayBuffer = e.target?.result as ArrayBuffer
      
      // Send upload request
      const uploadPath = currentPath === '.' ? '.' : currentPath
      ws.send(JSON.stringify({
        type: 'sftp_upload',
        data: {
          remote_path: uploadPath,
          file_name: file.name,
          size: file.size,
        },
      }))

      // Send file data in chunks
      const chunkSize = 8192
      let offset = 0
      
      const sendChunk = () => {
        if (offset >= arrayBuffer.byteLength) {
          return
        }
        
        const chunk = arrayBuffer.slice(offset, offset + chunkSize)
        ws?.send(chunk)
        offset += chunkSize
        
        // Update progress locally
        const progress = file.size > 0
          ? Math.round((offset / file.size) * 100)
          : 0
        setUploadProgress((prev) => ({
          ...prev,
          [file.name]: {
            progress: Math.min(progress, 100),
            bytesTransferred: Math.min(offset, file.size),
            totalBytes: file.size,
          },
        }))
        
        if (offset < arrayBuffer.byteLength) {
          setTimeout(sendChunk, 0)
        }
      }
      
      sendChunk()
    }
    
    reader.readAsArrayBuffer(file)
    return false // Prevent default upload
  }

  const handleDownload = (file: RemoteFile) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      message.error('WebSocket 连接未建立')
      return
    }

    let remotePath: string
    if (currentPath === '.' || currentPath === '/') {
      remotePath = currentPath === '/' ? `/${file.name}` : file.name
    } else {
      remotePath = `${currentPath}/${file.name}`
    }

    ws.send(JSON.stringify({
      type: 'sftp_download',
      data: { remote_path: remotePath },
    }))
  }

  const handleDelete = (file: RemoteFile) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      message.error('WebSocket 连接未建立')
      return
    }

    let remotePath: string
    if (currentPath === '.' || currentPath === '/') {
      remotePath = currentPath === '/' ? `/${file.name}` : file.name
    } else {
      remotePath = `${currentPath}/${file.name}`
    }

    ws.send(JSON.stringify({
      type: 'sftp_delete',
      data: { path: remotePath },
    }))
  }

  const handleCreateDir = () => {
    Modal.confirm({
      title: '创建目录',
      content: (
        <Input
          placeholder="目录名称"
          onPressEnter={(e) => {
            const dirName = e.currentTarget.value.trim()
            if (dirName) {
              let remotePath: string
              if (currentPath === '.' || currentPath === '/') {
                remotePath = currentPath === '/' ? `/${dirName}` : dirName
              } else {
                remotePath = `${currentPath}/${dirName}`
              }
              ws?.send(JSON.stringify({
                type: 'sftp_mkdir',
                data: { path: remotePath },
              }))
              Modal.destroyAll()
            }
          }}
          autoFocus
        />
      ),
      onOk: () => {
        const input = document.querySelector('input[placeholder="目录名称"]') as HTMLInputElement
        const dirName = input?.value.trim()
        if (dirName && ws) {
          let remotePath: string
          if (currentPath === '.' || currentPath === '/') {
            remotePath = currentPath === '/' ? `/${dirName}` : dirName
          } else {
            remotePath = `${currentPath}/${dirName}`
          }
          ws.send(JSON.stringify({
            type: 'sftp_mkdir',
            data: { path: remotePath },
          }))
        }
      },
    })
  }

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDate = (timestamp: number): string => {
    return new Date(timestamp * 1000).toLocaleString('zh-CN')
  }

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: RemoteFile) => (
        <Space 
          style={{ 
            cursor: record.isDir ? 'pointer' : 'default',
          }}
          onClick={() => record.isDir && handleFileClick(record)}
        >
          {record.isDir ? (
            <FolderOutlined 
              style={{ 
                fontSize: 18, 
                color: isDark ? '#60A5FA' : '#3B82F6',
              }} 
            />
          ) : (
            <FileOutlined 
              style={{ 
                fontSize: 16, 
                color: isDark ? '#94A3B8' : '#64748B',
              }} 
            />
          )}
          <Text
            style={{
              cursor: record.isDir ? 'pointer' : 'default',
              color: record.isDir 
                ? (isDark ? '#60A5FA' : '#2563EB') 
                : (isDark ? '#E2E8F0' : '#1E293B'),
              fontWeight: record.isDir ? 500 : 400,
            }}
          >
            {text}
          </Text>
        </Space>
      ),
    },
    {
      title: '大小',
      dataIndex: 'size',
      key: 'size',
      width: 100,
      render: (size: number, record: RemoteFile) => 
        record.isDir ? '-' : formatBytes(size),
    },
    {
      title: '修改时间',
      dataIndex: 'modTime',
      key: 'modTime',
      width: 180,
      render: (modTime: number) => formatDate(modTime),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: RemoteFile) => (
        <Space size="small">
          {!record.isDir && (
            <Tooltip title="下载">
              <Button
                type="link"
                size="small"
                icon={<DownloadOutlined />}
                onClick={() => handleDownload(record)}
                loading={downloadProgress[record.name] !== undefined && downloadProgress[record.name].progress < 100}
              />
            </Tooltip>
          )}
          <Popconfirm
            title="确定要删除吗？"
            onConfirm={() => handleDelete(record)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="link"
                size="small"
                danger
                icon={<DeleteOutlined />}
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  const pathParts = currentPath === '.' || currentPath === '/' ? [] : currentPath.split('/').filter(Boolean)

  const canGoBack = currentPath !== '.' && currentPath !== '/'
  const handleGoBack = () => {
    if (currentPath === '/' || currentPath === '.') return
    const parts = currentPath.split('/').filter(Boolean)
    if (parts.length === 1) {
      setCurrentPath('/')
    } else {
      setCurrentPath('/' + parts.slice(0, -1).join('/'))
    }
  }

  return (
    <Modal
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <FolderOpenOutlined style={{ fontSize: 18, color: isDark ? '#60A5FA' : '#1890FF' }} />
          <span style={{ fontWeight: 500 }}>SFTP 文件管理</span>
          {assetName && (
            <span style={{ 
              fontSize: 13, 
              color: isDark ? '#94A3B8' : '#64748B',
              fontWeight: 400,
            }}>
              {assetName}
            </span>
          )}
        </div>
      }
      open={visible}
      onCancel={onClose}
      footer={null}
      width={920}
      style={{ top: 40 }}
      styles={{
        body: {
          padding: 0,
          height: 'calc(100vh - 160px)',
          display: 'flex',
          flexDirection: 'column',
        },
        header: {
          padding: '16px 20px',
          borderBottom: isDark
            ? '1px solid rgba(255, 255, 255, 0.08)'
            : '1px solid rgba(0, 0, 0, 0.06)',
        },
      }}
    >
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
          background: isDark ? '#0F172A' : '#FFFFFF',
        }}
      >
        {/* Navigation Bar */}
        <div
          style={{
            padding: '12px 20px',
            borderBottom: isDark
              ? '1px solid rgba(255, 255, 255, 0.08)'
              : '1px solid rgba(0, 0, 0, 0.06)',
            background: isDark ? '#1E293B' : '#F8FAFC',
            display: 'flex',
            alignItems: 'center',
            gap: 12,
          }}
        >
          <Space.Compact>
            <Tooltip title="返回上级目录">
              <Button
                size="small"
                icon={<ArrowLeftOutlined />}
                onClick={handleGoBack}
                disabled={!canGoBack}
                style={{
                  borderRight: 'none',
                }}
              />
            </Tooltip>
            <Tooltip title="根目录">
              <Button
                size="small"
                icon={<HomeOutlined />}
                onClick={() => setCurrentPath('/')}
              />
            </Tooltip>
          </Space.Compact>

          {/* Breadcrumb */}
          <div style={{ flex: 1, minWidth: 0 }}>
            <Breadcrumb
              items={[
                {
                  title: (
                    <Button
                      type="link"
                      size="small"
                      icon={<HomeOutlined />}
                      onClick={() => setCurrentPath('/')}
                      style={{ 
                        padding: 0,
                        height: 'auto',
                        color: isDark ? '#C9D1D9' : '#1E293B',
                      }}
                    >
                      /
                    </Button>
                  ),
                },
                ...pathParts.map((part, index) => {
                  const path = '/' + pathParts.slice(0, index + 1).join('/')
                  return {
                    title: (
                      <Button
                        type="link"
                        size="small"
                        onClick={() => setCurrentPath(path)}
                        style={{ 
                          padding: 0,
                          height: 'auto',
                          color: isDark ? '#C9D1D9' : '#1E293B',
                        }}
                      >
                        {part}
                      </Button>
                    ),
                  }
                }),
              ]}
              style={{ fontSize: 13 }}
            />
          </div>

          <Space size="small">
            <Tooltip title="刷新">
              <Button
                size="small"
                icon={<ReloadOutlined />}
                onClick={() => listFiles(currentPath)}
                loading={loading}
              />
            </Tooltip>
            <Tooltip title="新建目录">
              <Button
                size="small"
                icon={<FolderAddOutlined />}
                onClick={handleCreateDir}
              />
            </Tooltip>
            <Upload
              beforeUpload={handleUpload}
              showUploadList={false}
            >
              <Tooltip title="上传文件">
                <Button
                  size="small"
                  type="primary"
                  icon={<UploadOutlined />}
                  loading={Object.keys(uploadProgress).length > 0}
                >
                  上传
                </Button>
              </Tooltip>
            </Upload>
          </Space>
        </div>

        {/* Content */}
        <div
          style={{
            flex: 1,
            overflow: 'auto',
            padding: 20,
            background: isDark ? '#0F172A' : '#FFFFFF',
          }}
        >

      {/* Upload Progress */}
      {Object.keys(uploadProgress).length > 0 && (
        <div
          style={{
            marginBottom: 16,
            padding: 16,
            background: isDark 
              ? 'linear-gradient(135deg, rgba(96, 165, 250, 0.1) 0%, rgba(139, 92, 246, 0.05) 100%)' 
              : 'linear-gradient(135deg, rgba(59, 130, 246, 0.08) 0%, rgba(139, 92, 246, 0.05) 100%)',
            borderRadius: 8,
            border: isDark
              ? '1px solid rgba(96, 165, 250, 0.2)'
              : '1px solid rgba(59, 130, 246, 0.2)',
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 12 }}>
            <UploadOutlined style={{ fontSize: 16, color: isDark ? '#60A5FA' : '#3B82F6' }} />
            <Text strong style={{ fontSize: 14, color: isDark ? '#E2E8F0' : '#1E293B' }}>
              上传进度
            </Text>
          </div>
          {Object.entries(uploadProgress).map(([fileName, progressInfo]) => (
            <div key={fileName} style={{ marginBottom: 12 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                <Text 
                  ellipsis 
                  style={{ 
                    fontSize: 13, 
                    color: isDark ? '#C9D1D9' : '#475569',
                    maxWidth: '60%',
                  }}
                >
                  {fileName}
                </Text>
                <Text 
                  style={{ 
                    fontSize: 12,
                    color: isDark ? '#94A3B8' : '#64748B',
                    fontFamily: 'monospace',
                  }}
                >
                  {formatBytes(progressInfo.bytesTransferred)} / {formatBytes(progressInfo.totalBytes)} ({progressInfo.progress}%)
                </Text>
              </div>
              <Progress
                percent={progressInfo.progress}
                status="active"
                strokeColor={{
                  '0%': '#3B82F6',
                  '100%': '#10B981',
                }}
                showInfo={false}
              />
            </div>
          ))}
        </div>
      )}

      {/* Download Progress */}
      {Object.keys(downloadProgress).length > 0 && (
        <div
          style={{
            marginBottom: 16,
            padding: 16,
            background: isDark 
              ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(59, 130, 246, 0.05) 100%)' 
              : 'linear-gradient(135deg, rgba(16, 185, 129, 0.08) 0%, rgba(59, 130, 246, 0.05) 100%)',
            borderRadius: 8,
            border: isDark
              ? '1px solid rgba(16, 185, 129, 0.2)'
              : '1px solid rgba(16, 185, 129, 0.2)',
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 12 }}>
            <DownloadOutlined style={{ fontSize: 16, color: isDark ? '#10B981' : '#059669' }} />
            <Text strong style={{ fontSize: 14, color: isDark ? '#E2E8F0' : '#1E293B' }}>
              下载进度
            </Text>
          </div>
          {Object.entries(downloadProgress).map(([fileName, progressInfo]) => (
            <div key={fileName} style={{ marginBottom: 12 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                <Text 
                  ellipsis 
                  style={{ 
                    fontSize: 13, 
                    color: isDark ? '#C9D1D9' : '#475569',
                    maxWidth: '60%',
                  }}
                >
                  {fileName}
                </Text>
                <Text 
                  style={{ 
                    fontSize: 12,
                    color: isDark ? '#94A3B8' : '#64748B',
                    fontFamily: 'monospace',
                  }}
                >
                  {formatBytes(progressInfo.bytesTransferred)} / {formatBytes(progressInfo.totalBytes)} ({progressInfo.progress}%)
                </Text>
              </div>
              <Progress
                percent={progressInfo.progress}
                status="active"
                strokeColor={{
                  '0%': '#10B981',
                  '100%': '#3B82F6',
                }}
                showInfo={false}
              />
            </div>
          ))}
        </div>
      )}

      {/* File List */}
      {loading && files.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '40px 0' }}>
          <Spin />
        </div>
      ) : files.length === 0 ? (
        <Empty description="目录为空" />
      ) : (
          <Table
            columns={columns}
            dataSource={files}
            rowKey="name"
            pagination={false}
            size="small"
            loading={loading}
            style={{
              background: isDark ? '#0F172A' : '#FFFFFF',
            }}
            rowClassName={(record) => {
              return record.isDir 
                ? (isDark ? 'sftp-dir-row-dark' : 'sftp-dir-row-light')
                : ''
            }}
            onRow={(record) => {
              if (!record.isDir) return {}
              return {
                onClick: () => handleFileClick(record),
                style: {
                  cursor: 'pointer',
                },
                onMouseEnter: (e) => {
                  if (isDark) {
                    e.currentTarget.style.background = 'rgba(96, 165, 250, 0.08)'
                  } else {
                    e.currentTarget.style.background = 'rgba(59, 130, 246, 0.05)'
                  }
                },
                onMouseLeave: (e) => {
                  e.currentTarget.style.background = 'transparent'
                },
              }
            }}
          />
        )}
        </div>
      </div>
    </Modal>
  )
}

export default SFTPManager
