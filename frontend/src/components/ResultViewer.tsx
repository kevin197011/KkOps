import React, { useState } from 'react';
import { Card, Table, Tag, Button, Space, Typography, Progress, Tooltip, message } from 'antd';
import { DownloadOutlined, CheckCircleOutlined, CloseCircleOutlined, CopyOutlined, EyeOutlined, EyeInvisibleOutlined, SyncOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { BatchOperation } from '../services/batchOperations';
import type { ColumnsType } from 'antd/es/table';

const { Text } = Typography;

interface ResultViewerProps {
  operation?: BatchOperation;
  onRefresh?: () => void;
}

interface ResultRow {
  hostId: string;
  hostname: string;
  status: 'success' | 'failed';
  output: string;
  error?: string;
}

const ResultViewer: React.FC<ResultViewerProps> = ({ operation, onRefresh }) => {
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());

  // 复制到剪贴板
  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      message.success('已复制到剪贴板');
    } catch (err) {
      // 降级方案
      const textArea = document.createElement('textarea');
      textArea.value = text;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      message.success('已复制到剪贴板');
    }
  };

  // 格式化输出文本（用于详情显示）
  const formatOutput = (output: string) => {
    if (!output) return '-';

    // 如果是JSON格式，尝试格式化
    if (output.trim().startsWith('{') || output.trim().startsWith('[')) {
      try {
        const parsed = JSON.parse(output);
        return JSON.stringify(parsed, null, 2);
      } catch (e) {
        // 不是有效的JSON，继续正常处理
      }
    }

    // 清理转义字符
    return output.replace(/\\n/g, '\n').replace(/\\t/g, '\t').replace(/\\"/g, '"');
  };

  if (!operation) {
    return (
      <Card title="执行结果">
        <div style={{ textAlign: 'center', padding: '40px', color: '#999' }}>
          请先选择主机并执行操作
        </div>
      </Card>
    );
  }

  // 解析结果数据
  const parseResults = (): ResultRow[] => {
    if (!operation.results) {
      return [];
    }

    const results: ResultRow[] = [];
    
    // 从 target_hosts 构建主机映射
    const hostMap = new Map<number, string>();
    if (operation.target_hosts && Array.isArray(operation.target_hosts)) {
      operation.target_hosts.forEach(host => {
        hostMap.set(host.id, host.hostname || `主机 #${host.id}`);
      });
    }

    // 解析 results JSON
    let resultsData: Record<string, { success: boolean; output: string; error?: string }> = {};
    if (operation.results) {
      let parsedResults: any;
      if (typeof operation.results === 'string') {
        try {
          parsedResults = JSON.parse(operation.results);
        } catch (e) {
          console.error('Failed to parse results:', e);
          return [];
        }
      } else if (typeof operation.results === 'object') {
        parsedResults = operation.results;
      }

      // 检查是否是嵌套结构（新格式）
      if (parsedResults && typeof parsedResults === 'object' && parsedResults.results) {
        // 新格式：{ status, results, completed_at, ... }
        if (typeof parsedResults.results === 'string') {
          try {
            resultsData = JSON.parse(parsedResults.results);
          } catch (e) {
            console.error('Failed to parse nested results:', e);
            return [];
          }
        } else {
          resultsData = parsedResults.results;
        }
      } else {
        // 旧格式：直接的 { hostId: result } 对象
        resultsData = parsedResults;
      }
    }

    // 构建结果行
    Object.entries(resultsData).forEach(([hostId, result]) => {
      const hostname = hostMap.get(parseInt(hostId)) || `主机 #${hostId}`;
      results.push({
        hostId,
        hostname,
        status: result.success ? 'success' : 'failed',
        output: result.output || '',
        error: result.error,
      });
    });

    return results;
  };

  const results = parseResults();

  const getStatusTag = (status: string) => {
    if (status === 'success') {
      return <Tag color="success" icon={<CheckCircleOutlined />}>成功</Tag>;
    }
    return <Tag color="error" icon={<CloseCircleOutlined />}>失败</Tag>;
  };

  const handleExportCSV = () => {
    const csvRows = [
      ['主机ID', '主机名', '状态', '输出', '错误'].join(','),
      ...results.map(row => [
        row.hostId,
        row.hostname,
        row.status === 'success' ? '成功' : '失败',
        `"${(row.output || '').replace(/"/g, '""')}"`,
        `"${(row.error || '').replace(/"/g, '""')}"`,
      ].join(',')),
    ];

    const csv = csvRows.join('\n');
    const blob = new Blob(['\ufeff' + csv], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `batch-operation-${operation.id}-results.csv`;
    link.click();
  };

  const handleExportJSON = () => {
    const data = {
      operation_id: operation.id,
      operation_name: operation.name,
      status: operation.status,
      success_count: operation.success_count,
      failed_count: operation.failed_count,
      results: results,
    };
    const json = JSON.stringify(data, null, 2);
    const blob = new Blob([json], { type: 'application/json' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `batch-operation-${operation.id}-results.json`;
    link.click();
  };

  const columns: ColumnsType<ResultRow> = [
    {
      title: '主机名',
      dataIndex: 'hostname',
      key: 'hostname',
      width: 150,
      render: (text: string) => <Text strong>{text}</Text>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="复制输出">
            <Button
              size="small"
              type="text"
              icon={<CopyOutlined />}
              onClick={() => copyToClipboard(record.output || '')}
              style={{ padding: '2px 6px' }}
            />
          </Tooltip>
          <Tooltip title={expandedRows.has(record.hostId) ? '收起详情' : '查看详情'}>
            <Button
              size="small"
              type="text"
              icon={expandedRows.has(record.hostId) ? <EyeInvisibleOutlined /> : <EyeOutlined />}
              onClick={() => {
                const newExpanded = new Set(expandedRows);
                if (newExpanded.has(record.hostId)) {
                  newExpanded.delete(record.hostId);
                } else {
                  newExpanded.add(record.hostId);
                }
                setExpandedRows(newExpanded);
              }}
              style={{ padding: '2px 6px' }}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const getProgress = () => {
    if (operation.status === 'completed' || operation.status === 'failed') {
      return 100;
    }
    if (operation.status === 'running') {
      // 可以根据已完成的个数计算进度
      const total = operation.target_count;
      const completed = operation.success_count + operation.failed_count;
      return total > 0 ? Math.round((completed / total) * 100) : 0;
    }
    return 0;
  };

  return (
    <Card
        title="执行结果"
      extra={
        results.length > 0 && (
          <Space>
            <Button size="small" icon={<DownloadOutlined />} onClick={handleExportCSV}>
              导出 CSV
            </Button>
            <Button size="small" icon={<DownloadOutlined />} onClick={handleExportJSON}>
              导出 JSON
            </Button>
            {onRefresh && (
              <Button size="small" onClick={onRefresh}>
                刷新
              </Button>
            )}
          </Space>
        )
      }
      style={{ height: '100%' }}
    >
      {(operation.status === 'pending' || operation.status === 'running') && (
        <div style={{ marginBottom: 16 }}>
          {/* 进度条 */}
          <Progress percent={getProgress()} status="active" />
          <div style={{ marginTop: 8, textAlign: 'center', marginBottom: 12 }}>
            <Text type="secondary">
              {operation.status === 'pending' ? (
                <><ClockCircleOutlined style={{ marginRight: 4 }} /> 等待执行...</>
              ) : (
                <><SyncOutlined spin style={{ marginRight: 4 }} /> 执行中... ({operation.success_count + operation.failed_count}/{operation.target_count})</>
              )}
            </Text>
          </div>
          
          {/* 实时日志终端 */}
          <div style={{
            background: '#1e1e1e',
            color: '#d4d4d4',
            padding: '16px',
            borderRadius: '6px',
            fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
            fontSize: '13px',
            minHeight: '180px',
            maxHeight: '300px',
            overflow: 'auto',
          }}>
            <div style={{ color: '#569cd6' }}>
              <SyncOutlined spin style={{ marginRight: 6 }} />
              [INFO] 批量操作任务已提交
            </div>
            <div style={{ color: '#6a9955', marginTop: 4 }}>
              [INFO] 操作名称: {operation.name}
            </div>
            <div style={{ color: '#6a9955', marginTop: 4 }}>
              [INFO] 命令函数: {operation.command_function}
            </div>
            {operation.command_args && (
              <div style={{ color: '#6a9955', marginTop: 4 }}>
                [INFO] 命令参数: {JSON.stringify(operation.command_args)}
              </div>
            )}
            <div style={{ color: '#dcdcaa', marginTop: 4 }}>
              [INFO] 目标主机数: {operation.target_count}
            </div>
            {operation.target_hosts && operation.target_hosts.length > 0 && (
              <div style={{ color: '#ce9178', marginTop: 4 }}>
                [INFO] 目标主机: {operation.target_hosts.slice(0, 5).map(h => h.hostname).join(', ')}
                {operation.target_hosts.length > 5 && ` ... 等 ${operation.target_hosts.length} 台`}
              </div>
            )}
            <div style={{ color: '#569cd6', marginTop: 8 }}>
              {operation.status === 'pending' ? (
                <><ClockCircleOutlined style={{ marginRight: 6 }} /> 等待Salt Master响应...</>
              ) : (
                <>
                  <SyncOutlined spin style={{ marginRight: 6 }} />
                  正在执行命令，已完成 {operation.success_count + operation.failed_count}/{operation.target_count}...
                </>
              )}
            </div>
            {operation.success_count > 0 && (
              <div style={{ color: '#4ec9b0', marginTop: 4 }}>
                [SUCCESS] 成功: {operation.success_count} 台
              </div>
            )}
            {operation.failed_count > 0 && (
              <div style={{ color: '#f14c4c', marginTop: 4 }}>
                [FAILED] 失败: {operation.failed_count} 台
              </div>
            )}
          </div>
        </div>
      )}

      {(operation.status === 'completed' || operation.status === 'failed' || results.length > 0) && (
        <>
          {/* 统计信息 */}
          <div style={{ marginBottom: 16, padding: '12px', backgroundColor: '#f5f5f5', borderRadius: 4 }}>
            <Space>
              <Text>总计: {operation.target_count}</Text>
              <Tag color="success">成功: {operation.success_count}</Tag>
              <Tag color="error">失败: {operation.failed_count}</Tag>
              {operation.duration_seconds && (
                <Text type="secondary">耗时: {operation.duration_seconds} 秒</Text>
              )}
            </Space>
          </div>

          {/* 结果表格 */}
          <Table
            columns={columns}
            dataSource={results}
            rowKey="hostId"
            pagination={false}
            size="small"
            expandable={{
              expandedRowKeys: Array.from(expandedRows),
              onExpandedRowsChange: (expandedKeys) => {
                setExpandedRows(new Set(expandedKeys as string[]));
              },
              expandedRowRender: (record) => (
                <div style={{ padding: '16px', backgroundColor: '#fafafa', borderRadius: '8px', margin: '8px 0' }}>
                  <Space direction="vertical" style={{ width: '100%' }} size="middle">
                    {/* 输出区域 */}
                    <div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                        <Text strong style={{ color: '#1890ff' }}>📄 执行输出</Text>
                        <Space>
                          <Tooltip title="复制完整输出">
                            <Button
                              size="small"
                              icon={<CopyOutlined />}
                              onClick={() => copyToClipboard(record.output || '')}
                            >
                              复制
                            </Button>
                          </Tooltip>
                          <Tag color={record.status === 'success' ? 'success' : 'error'}>
                            {record.status === 'success' ? '执行成功' : '执行失败'}
                          </Tag>
                        </Space>
                      </div>
                      <div style={{
                        backgroundColor: '#fff',
                        border: '1px solid #d9d9d9',
                        borderRadius: '6px',
                        padding: '12px',
                        fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                        fontSize: '13px',
                        lineHeight: '1.5',
                        maxHeight: '300px',
                        overflow: 'auto',
                        whiteSpace: 'pre-wrap',
                        wordBreak: 'break-all'
                      }}>
                        {formatOutput(record.output)}
                      </div>
                      {record.output && record.output.length > 500 && (
                        <Text type="secondary" style={{ fontSize: '12px', marginTop: '4px' }}>
                          输出长度: {record.output.length} 字符
                        </Text>
                      )}
                    </div>

                    {/* 错误信息区域 */}
                    {record.error && (
                      <div>
                        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
                          <Text strong style={{ color: '#ff4d4f' }}>❌ 错误信息</Text>
                        </div>
                        <div style={{
                          backgroundColor: '#fff2f0',
                          border: '1px solid #ffccc7',
                          borderRadius: '6px',
                          padding: '12px',
                          fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                          fontSize: '13px',
                          color: '#ff4d4f',
                          maxHeight: '200px',
                          overflow: 'auto',
                          whiteSpace: 'pre-wrap'
                        }}>
                          {record.error}
                        </div>
                      </div>
                    )}

                    {/* 元信息 */}
                    <div style={{
                      padding: '8px 12px',
                      backgroundColor: '#f0f2f5',
                      borderRadius: '4px',
                      fontSize: '12px'
                    }}>
                      <Space>
                        <Text type="secondary">主机ID: {record.hostId}</Text>
                        <Text type="secondary">主机名: {record.hostname}</Text>
                        <Text type="secondary">
                          执行时间: {operation.started_at ? new Date(operation.started_at).toLocaleString() : '-'}
                        </Text>
                      </Space>
                    </div>
                  </Space>
                </div>
              ),
            }}
          />

          {results.length === 0 && operation.status === 'completed' && (
            <div style={{ textAlign: 'center', padding: '20px', color: '#999' }}>
              暂无结果数据
            </div>
          )}
        </>
      )}

      {operation.status === 'cancelled' && (
        <div style={{ textAlign: 'center', padding: '20px' }}>
          <Tag color="warning">操作已取消</Tag>
        </div>
      )}

      {operation.error_message && (
        <div style={{ marginTop: 16, padding: '12px', backgroundColor: '#fff2f0', borderRadius: 4 }}>
          <Text type="danger" strong>错误信息:</Text>
          <div style={{ marginTop: 4 }}>{operation.error_message}</div>
        </div>
      )}
    </Card>
  );
};

export default ResultViewer;

