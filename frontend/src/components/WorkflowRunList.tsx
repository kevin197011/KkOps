import React from 'react';
import { Card, Empty, Spin, Pagination, Space, Button, Typography } from 'antd';
import { ReloadOutlined, DeleteOutlined } from '@ant-design/icons';
import WorkflowRunItem, { WorkflowRunItemProps } from './WorkflowRunItem';

const { Text } = Typography;

export interface WorkflowRunListProps {
  title?: string;
  items: WorkflowRunItemProps[];
  loading?: boolean;
  total: number;
  current: number;
  pageSize: number;
  onPageChange: (page: number, pageSize?: number) => void;
  onItemClick: (item: WorkflowRunItemProps) => void;
  onRefresh?: () => void;
  onCleanup?: () => void;
  cleanupLoading?: boolean;
}

const WorkflowRunList: React.FC<WorkflowRunListProps> = ({
  title = '执行记录',
  items,
  loading = false,
  total,
  current,
  pageSize,
  onPageChange,
  onItemClick,
  onRefresh,
  onCleanup,
  cleanupLoading = false,
}) => {
  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Space>
            <span>{title}</span>
            <Text type="secondary" style={{ fontSize: 12, fontWeight: 'normal' }}>
              共 {total} 条记录
            </Text>
          </Space>
          <Space>
            {onRefresh && (
              <Button 
                size="small" 
                icon={<ReloadOutlined />} 
                onClick={onRefresh}
                loading={loading}
              >
                刷新
              </Button>
            )}
            {onCleanup && (
              <Button 
                size="small" 
                danger
                icon={<DeleteOutlined />} 
                onClick={onCleanup}
                loading={cleanupLoading}
              >
                清理旧数据
              </Button>
            )}
          </Space>
        </div>
      }
      bodyStyle={{ padding: 0 }}
      style={{ marginTop: 16 }}
    >
      <Spin spinning={loading}>
        {items.length > 0 ? (
          <>
            <div style={{ 
              borderBottom: '1px solid #f0f0f0',
              padding: '8px 16px',
              backgroundColor: '#fafafa',
              display: 'flex',
              alignItems: 'center',
              fontSize: 12,
              color: '#8c8c8c',
            }}>
              <div style={{ width: 44 }}></div>
              <div style={{ flex: 1 }}>执行信息</div>
              <div style={{ width: 120, textAlign: 'center' }}>成功/失败/总计</div>
              <div style={{ width: 50, textAlign: 'right' }}>耗时</div>
            </div>
            {items.map((item) => (
              <WorkflowRunItem
                key={item.id}
                {...item}
                onClick={() => onItemClick(item)}
              />
            ))}
            <div style={{ 
              padding: '12px 16px', 
              display: 'flex', 
              justifyContent: 'flex-end',
              borderTop: '1px solid #f0f0f0',
              backgroundColor: '#fafafa',
            }}>
              <Pagination
                size="small"
                current={current}
                pageSize={pageSize}
                total={total}
                onChange={onPageChange}
                showSizeChanger
                showQuickJumper
                pageSizeOptions={['10', '20', '50', '100']}
                showTotal={(total, range) => (
                  <span style={{ marginRight: 8 }}>
                    第 {range[0]}-{range[1]} 条，共 {total} 条
                  </span>
                )}
              />
            </div>
          </>
        ) : (
          <Empty 
            description="暂无记录" 
            style={{ padding: '40px 0' }}
          />
        )}
      </Spin>
    </Card>
  );
};

export default WorkflowRunList;
