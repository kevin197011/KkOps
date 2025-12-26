import React, { useState, useEffect, useRef } from 'react';
import { Modal, Tag, Space, Typography, Button, Tooltip, Tabs, Progress, Collapse, message } from 'antd';
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  SyncOutlined,
  StopOutlined,
  CopyOutlined,
  DownloadOutlined,
  ExpandOutlined,
  CompressOutlined,
  ReloadOutlined,
} from '@ant-design/icons';

const { Text } = Typography;

// 执行状态类型
export type ExecutionStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';

// 执行结果项
export interface ExecutionResultItem {
  hostId: string;
  hostname: string;
  status: 'success' | 'failed' | 'pending' | 'running';
  output?: string;
  error?: string;
  duration?: number;
}

// 执行记录接口
export interface ExecutionRecord {
  id: number;
  name: string;
  type: 'batch' | 'formula';
  status: ExecutionStatus;
  command?: string;
  targetCount: number;
  successCount: number;
  failedCount: number;
  duration?: number;
  startedAt?: string;
  completedAt?: string;
  errorMessage?: string;
  results?: ExecutionResultItem[];
  rawResults?: any;
}

interface ExecutionLogViewerProps {
  visible: boolean;
  record: ExecutionRecord | null;
  onClose: () => void;
  onRefresh?: () => void;
  onCancel?: () => void;
  onRetry?: () => void;
}

// 状态配置
const statusConfig: Record<ExecutionStatus, { color: string; text: string; icon: React.ReactNode; bgColor: string }> = {
  pending: { color: '#faad14', text: '等待中', icon: <ClockCircleOutlined />, bgColor: '#fffbe6' },
  running: { color: '#1890ff', text: '运行中', icon: <SyncOutlined spin />, bgColor: '#e6f7ff' },
  completed: { color: '#52c41a', text: '成功', icon: <CheckCircleOutlined />, bgColor: '#f6ffed' },
  failed: { color: '#ff4d4f', text: '失败', icon: <CloseCircleOutlined />, bgColor: '#fff2f0' },
  cancelled: { color: '#8c8c8c', text: '已取消', icon: <StopOutlined />, bgColor: '#fafafa' },
};

// 格式化耗时
const formatDuration = (seconds?: number): string => {
  if (seconds === undefined || seconds === null) return '-';
  if (seconds < 0) return '-';
  if (seconds === 0) return '<1s';
  if (seconds < 60) return `${Math.round(seconds)}s`;
  if (seconds < 3600) {
    const mins = Math.floor(seconds / 60);
    const secs = Math.round(seconds % 60);
    return secs > 0 ? `${mins}m ${secs}s` : `${mins}m`;
  }
  const hours = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
};

const ExecutionLogViewer: React.FC<ExecutionLogViewerProps> = ({
  visible,
  record,
  onClose,
  onRefresh,
  onCancel,
  onRetry,
}) => {
  const [fullscreen, setFullscreen] = useState(false);
  const [activeTab, setActiveTab] = useState('summary');
  const logContainerRef = useRef<HTMLDivElement>(null);

  // 自动滚动到底部
  useEffect(() => {
    if (logContainerRef.current && record?.status === 'running') {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight;
    }
  }, [record?.results, record?.status]);

  // 重置tab
  useEffect(() => {
    if (visible) {
      setActiveTab('summary');
    }
  }, [visible]);

  // 复制到剪贴板
  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      message.success('已复制到剪贴板');
    } catch (err) {
      const textArea = document.createElement('textarea');
      textArea.value = text;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      message.success('已复制到剪贴板');
    }
  };

  // 导出日志
  const exportLogs = () => {
    if (!record) return;
    const logs = generateFullLog();
    const blob = new Blob([logs], { type: 'text/plain' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `execution-${record.id}-${record.type}.log`;
    link.click();
    message.success('日志已下载');
  };

  // 生成完整日志
  const generateFullLog = (): string => {
    if (!record) return '';
    
    const lines: string[] = [
      `========================================`,
      `执行记录 #${record.id}`,
      `========================================`,
      `类型: ${record.type === 'batch' ? '批量操作' : 'Formula部署'}`,
      `名称: ${record.name}`,
      `状态: ${statusConfig[record.status].text}`,
      `开始时间: ${record.startedAt ? new Date(record.startedAt).toLocaleString() : '-'}`,
      `完成时间: ${record.completedAt ? new Date(record.completedAt).toLocaleString() : '-'}`,
      `耗时: ${formatDuration(record.duration)}`,
      `目标主机: ${record.targetCount}`,
      `成功: ${record.successCount}`,
      `失败: ${record.failedCount}`,
      ``,
      `========================================`,
      `执行详情`,
      `========================================`,
    ];

    if (record.command) {
      lines.push(`命令: ${record.command}`);
      lines.push(``);
    }

    if (record.results && record.results.length > 0) {
      record.results.forEach((result) => {
        lines.push(`--- ${result.hostname} (${result.status}) ---`);
        if (result.output) {
          lines.push(result.output);
        }
        if (result.error) {
          lines.push(`[ERROR] ${result.error}`);
        }
        lines.push(``);
      });
    }

    if (record.errorMessage) {
      lines.push(`========================================`);
      lines.push(`错误信息`);
      lines.push(`========================================`);
      lines.push(record.errorMessage);
    }

    return lines.join('\n');
  };

  // 计算进度
  const getProgress = () => {
    if (!record) return 0;
    if (record.status === 'completed' || record.status === 'failed') return 100;
    if (record.status === 'pending') return 0;
    const completed = record.successCount + record.failedCount;
    return record.targetCount > 0 ? Math.round((completed / record.targetCount) * 100) : 0;
  };

  if (!record) return null;

  const status = statusConfig[record.status];
  const isRunning = record.status === 'pending' || record.status === 'running';

  return (
    <Modal
      title={null}
      open={visible}
      onCancel={onClose}
      footer={null}
      width={fullscreen ? '100%' : 900}
      style={fullscreen ? { top: 0, padding: 0, maxWidth: '100%' } : undefined}
      styles={{
        body: fullscreen 
          ? { height: 'calc(100vh - 55px)', overflow: 'hidden', padding: 0 } 
          : { maxHeight: '80vh', overflow: 'hidden', padding: 0 }
      }}
      className="execution-log-modal"
    >
      {/* Header - GitHub Actions 风格 */}
      <div style={{
        background: status.bgColor,
        borderBottom: `2px solid ${status.color}`,
        padding: '16px 24px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'flex-start',
      }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
            <span style={{ color: status.color, fontSize: 24 }}>{status.icon}</span>
            <Text strong style={{ fontSize: 18 }}>{record.name}</Text>
            <Tag color={status.color}>{status.text}</Tag>
          </div>
          <Space size="large" style={{ color: '#666', fontSize: 13 }}>
            <span>#{record.id}</span>
            <span>{record.type === 'batch' ? '批量操作' : 'Formula部署'}</span>
            {record.startedAt && <span>开始于 {new Date(record.startedAt).toLocaleString()}</span>}
            {record.duration !== undefined && record.duration !== null && <span>耗时 {formatDuration(record.duration)}</span>}
          </Space>
        </div>
        <Space>
          {isRunning && onCancel && (
            <Button danger icon={<StopOutlined />} onClick={onCancel}>
              取消
            </Button>
          )}
          {(record.status === 'failed' || record.status === 'cancelled') && onRetry && (
            <Button icon={<ReloadOutlined />} onClick={onRetry}>
              重试
            </Button>
          )}
          <Tooltip title={fullscreen ? '退出全屏' : '全屏'}>
            <Button
              icon={fullscreen ? <CompressOutlined /> : <ExpandOutlined />}
              onClick={() => setFullscreen(!fullscreen)}
            />
          </Tooltip>
        </Space>
      </div>

      <div style={{ padding: '0 24px 24px 24px' }}>
        {/* Progress Bar */}
        {isRunning && (
          <div style={{ padding: '16px 0 0 0' }}>
            <Progress
              percent={getProgress()}
              status="active"
              strokeColor={status.color}
              format={() => `${record.successCount + record.failedCount}/${record.targetCount}`}
            />
          </div>
        )}

        {/* Stats Bar */}
        <div style={{
          display: 'flex',
          gap: 24,
          padding: '16px 0',
          borderBottom: '1px solid #f0f0f0',
        }}>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 24, fontWeight: 600 }}>{record.targetCount}</div>
            <div style={{ color: '#666', fontSize: 12 }}>目标主机</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 24, fontWeight: 600, color: '#52c41a' }}>{record.successCount}</div>
            <div style={{ color: '#666', fontSize: 12 }}>成功</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 24, fontWeight: 600, color: '#ff4d4f' }}>{record.failedCount}</div>
            <div style={{ color: '#666', fontSize: 12 }}>失败</div>
          </div>
          {record.duration !== undefined && record.duration !== null && (
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 600 }}>{formatDuration(record.duration)}</div>
              <div style={{ color: '#666', fontSize: 12 }}>耗时</div>
            </div>
          )}
          <div style={{ flex: 1 }} />
          <Space>
            <Tooltip title="复制日志">
              <Button icon={<CopyOutlined />} onClick={() => copyToClipboard(generateFullLog())} />
            </Tooltip>
            <Tooltip title="下载日志">
              <Button icon={<DownloadOutlined />} onClick={exportLogs} />
            </Tooltip>
            {onRefresh && (
              <Tooltip title="刷新">
                <Button icon={<ReloadOutlined />} onClick={onRefresh} />
              </Tooltip>
            )}
          </Space>
        </div>

        {/* Tabs */}
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          style={{ marginTop: 8 }}
          items={[
            {
              key: 'summary',
              label: '执行日志',
              children: (
                <div style={{ height: fullscreen ? 'calc(100vh - 320px)' : 400, overflow: 'auto' }}>
                  {/* 实时日志终端 */}
                  <div
                    ref={logContainerRef}
                    style={{
                      background: '#1e1e1e',
                      color: '#d4d4d4',
                      padding: '16px',
                      borderRadius: '8px',
                      fontFamily: 'Monaco, Menlo, "Ubuntu Mono", Consolas, monospace',
                      fontSize: '13px',
                      lineHeight: '1.6',
                      minHeight: 200,
                      height: '100%',
                      overflow: 'auto',
                    }}
                  >
                    <LogLine type="info" time={record.startedAt}>任务开始执行</LogLine>
                    <LogLine type="info">名称: {record.name}</LogLine>
                    {record.command && <LogLine type="info">命令: {record.command}</LogLine>}
                    <LogLine type="info">目标主机数: {record.targetCount}</LogLine>
                    
                    {record.results && record.results.map((result) => (
                      <React.Fragment key={result.hostId}>
                        <LogLine type={result.status === 'success' ? 'success' : result.status === 'failed' ? 'error' : 'info'}>
                          [{result.hostname}] {result.status === 'success' ? '✓ 成功' : result.status === 'failed' ? '✗ 失败' : '执行中...'}
                        </LogLine>
                        {result.output && (
                          <div style={{ color: '#9cdcfe', marginLeft: 20, whiteSpace: 'pre-wrap', wordBreak: 'break-all' }}>
                            {result.output.length > 1000 ? result.output.substring(0, 1000) + '\n... (输出过长，请查看主机详情)' : result.output}
                          </div>
                        )}
                        {result.error && (
                          <div style={{ color: '#f14c4c', marginLeft: 20 }}>{result.error}</div>
                        )}
                      </React.Fragment>
                    ))}

                    {isRunning && (
                      <LogLine type="running">
                        <SyncOutlined spin style={{ marginRight: 8 }} />
                        执行中... ({record.successCount + record.failedCount}/{record.targetCount})
                      </LogLine>
                    )}

                    {record.status === 'completed' && (
                      <LogLine type="success" time={record.completedAt}>
                        ✓ 任务完成 - 成功: {record.successCount}, 失败: {record.failedCount}
                      </LogLine>
                    )}

                    {record.status === 'failed' && (
                      <>
                        {record.errorMessage && <LogLine type="error">{record.errorMessage}</LogLine>}
                        <LogLine type="error" time={record.completedAt}>✗ 任务失败</LogLine>
                      </>
                    )}

                    {record.status === 'cancelled' && (
                      <LogLine type="warning" time={record.completedAt}>⚠ 任务已取消</LogLine>
                    )}
                  </div>
                </div>
              ),
            },
            {
              key: 'hosts',
              label: `主机详情 (${record.results?.length || 0})`,
              children: (
                <div style={{ height: fullscreen ? 'calc(100vh - 320px)' : 400, overflow: 'auto' }}>
                  {record.results && record.results.length > 0 ? (
                    <Collapse
                      accordion
                      items={record.results.map((result) => ({
                        key: result.hostId,
                        label: (
                          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                            {result.status === 'success' ? (
                              <CheckCircleOutlined style={{ color: '#52c41a' }} />
                            ) : result.status === 'failed' ? (
                              <CloseCircleOutlined style={{ color: '#ff4d4f' }} />
                            ) : (
                              <SyncOutlined spin style={{ color: '#1890ff' }} />
                            )}
                            <Text strong>{result.hostname}</Text>
                            <Tag color={result.status === 'success' ? 'success' : result.status === 'failed' ? 'error' : 'processing'}>
                              {result.status === 'success' ? '成功' : result.status === 'failed' ? '失败' : '执行中'}
                            </Tag>
                            {result.duration !== undefined && <Text type="secondary">{result.duration}s</Text>}
                          </div>
                        ),
                        children: (
                          <div>
                            {result.output && (
                              <div style={{ marginBottom: 16 }}>
                                <div style={{ marginBottom: 8, display: 'flex', justifyContent: 'space-between' }}>
                                  <Text strong>输出</Text>
                                  <Button
                                    size="small"
                                    icon={<CopyOutlined />}
                                    onClick={() => copyToClipboard(result.output || '')}
                                  >
                                    复制
                                  </Button>
                                </div>
                                <pre style={{
                                  background: '#f5f5f5',
                                  padding: 12,
                                  borderRadius: 6,
                                  overflow: 'auto',
                                  maxHeight: 300,
                                  fontSize: 12,
                                  margin: 0,
                                  whiteSpace: 'pre-wrap',
                                  wordBreak: 'break-all',
                                }}>
                                  {result.output}
                                </pre>
                              </div>
                            )}
                            {result.error && (
                              <div>
                                <Text strong style={{ color: '#ff4d4f' }}>错误</Text>
                                <pre style={{
                                  background: '#fff2f0',
                                  padding: 12,
                                  borderRadius: 6,
                                  color: '#ff4d4f',
                                  overflow: 'auto',
                                  maxHeight: 200,
                                  fontSize: 12,
                                  margin: '8px 0 0 0',
                                  whiteSpace: 'pre-wrap',
                                  wordBreak: 'break-all',
                                }}>
                                  {result.error}
                                </pre>
                              </div>
                            )}
                            {!result.output && !result.error && (
                              <Text type="secondary">暂无输出</Text>
                            )}
                          </div>
                        ),
                      }))}
                    />
                  ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                      {isRunning ? '等待执行结果...' : '暂无详细结果'}
                    </div>
                  )}
                </div>
              ),
            },
          ]}
        />
      </div>
    </Modal>
  );
};

// 日志行组件
const LogLine: React.FC<{
  type: 'info' | 'success' | 'error' | 'warning' | 'running';
  time?: string;
  children: React.ReactNode;
}> = ({ type, time, children }) => {
  const colors = {
    info: '#569cd6',
    success: '#4ec9b0',
    error: '#f14c4c',
    warning: '#dcdcaa',
    running: '#569cd6',
  };

  return (
    <div style={{ color: colors[type], marginBottom: 4 }}>
      {time && (
        <span style={{ color: '#6a9955', marginRight: 8 }}>
          [{new Date(time).toLocaleTimeString()}]
        </span>
      )}
      {children}
    </div>
  );
};

export default ExecutionLogViewer;
