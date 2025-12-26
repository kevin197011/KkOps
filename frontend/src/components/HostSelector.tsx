import React, { useState, useEffect } from 'react';
import { Card, Select, Checkbox, Space, Input, Tag, Button, Typography } from 'antd';
import { ReloadOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { hostService, Host } from '../services/host';
import { projectService, Project } from '../services/project';
import { hostGroupService, HostGroup } from '../services/hostGroup';
import { hostTagService, HostTag } from '../services/hostTag';

const { Option } = Select;
const { Text } = Typography;
const { Search } = Input;

interface HostSelectorProps {
  value?: Array<{ id: number; hostname: string; ip_address?: string }>;
  onChange?: (selected: Array<{ id: number; hostname: string; ip_address?: string }>) => void;
}

const HostSelector: React.FC<HostSelectorProps> = ({ value = [], onChange }) => {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [hostGroups, setHostGroups] = useState<HostGroup[]>([]);
  const [hostTags, setHostTags] = useState<HostTag[]>([]);
  const [loading, setLoading] = useState(false);
  
  // 筛选条件
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined);
  const [selectedGroupId, setSelectedGroupId] = useState<number | undefined>(undefined);
  const [selectedTagId, setSelectedTagId] = useState<number | undefined>(undefined);
  const [selectedStatus, setSelectedStatus] = useState<string | undefined>(undefined);
  const [searchText, setSearchText] = useState('');

  // 已选择的主机ID集合
  const selectedHostIds = new Set(value.map(h => h.id));

  useEffect(() => {
    loadProjects();
    loadHostGroups();
    loadHostTags();
  }, []);

  useEffect(() => {
    loadHosts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedProjectId, selectedGroupId, selectedTagId, selectedStatus, searchText]);

  const loadProjects = async () => {
    try {
      const response = await projectService.list(1, 100);
      setProjects(response.projects);
    } catch (error: any) {
      console.error('Failed to load projects:', error);
    }
  };

  const loadHostGroups = async () => {
    try {
      const response = await hostGroupService.list(1, 100);
      setHostGroups(response.groups || []);
    } catch (error: any) {
      console.error('Failed to load host groups:', error);
    }
  };

  const loadHostTags = async () => {
    try {
      const response = await hostTagService.list(1, 100);
      setHostTags(response.tags || []);
    } catch (error: any) {
      console.error('Failed to load host tags:', error);
    }
  };

  const loadHosts = async () => {
    setLoading(true);
    try {
      const filters: any = {};
      if (selectedProjectId) filters.project_id = selectedProjectId;
      if (selectedGroupId) filters.group_id = selectedGroupId;
      if (selectedTagId) filters.tag_id = selectedTagId;
      if (selectedStatus) filters.status = selectedStatus;
      if (searchText) filters.hostname = searchText;

      const response = await hostService.list(1, 1000, filters);
      setHosts(response.hosts);
    } catch (error: any) {
      console.error('Failed to load hosts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleHostToggle = (host: Host) => {
    const newSelected = [...value];
    const index = newSelected.findIndex(h => h.id === host.id);
    
    if (index >= 0) {
      newSelected.splice(index, 1);
    } else {
      newSelected.push({
        id: host.id,
        hostname: host.hostname,
        ip_address: host.ip_address,
      });
    }
    
    onChange?.(newSelected);
  };

  const handleSelectAll = () => {
    const allSelected = hosts.map(host => ({
      id: host.id,
      hostname: host.hostname,
      ip_address: host.ip_address,
    }));
    onChange?.(allSelected);
  };

  const handleClearSelection = () => {
    onChange?.([]);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online':
        return 'success';
      case 'offline':
        return 'error';
      default:
        return 'default';
    }
  };

  const filteredHosts = hosts.filter(host => {
    if (searchText) {
      const search = searchText.toLowerCase();
      return (
        host.hostname.toLowerCase().includes(search) ||
        host.ip_address.toLowerCase().includes(search)
      );
    }
    return true;
  });

  return (
    <Card
      title="选择主机"
      extra={
        <Space>
          <Text type="secondary">已选择: {value.length} 台</Text>
          <Button size="small" onClick={loadHosts} icon={<ReloadOutlined />}>
            刷新
          </Button>
        </Space>
      }
      style={{ height: '100%' }}
    >
      <Space direction="vertical" style={{ width: '100%' }} size="middle">
        {/* 筛选器 */}
        <Space wrap>
          <Select
            placeholder="按项目筛选"
            allowClear
            style={{ width: 150 }}
            value={selectedProjectId}
            onChange={setSelectedProjectId}
          >
            {projects.map(project => (
              <Option key={project.id} value={project.id}>
                {project.name}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="按分组筛选"
            allowClear
            style={{ width: 150 }}
            value={selectedGroupId}
            onChange={setSelectedGroupId}
          >
            {hostGroups.map(group => (
              <Option key={group.id} value={group.id}>
                {group.name}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="按标签筛选"
            allowClear
            style={{ width: 150 }}
            value={selectedTagId}
            onChange={setSelectedTagId}
          >
            {hostTags.map(tag => (
              <Option key={tag.id} value={tag.id}>
                {tag.name}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="按状态筛选"
            allowClear
            style={{ width: 150 }}
            value={selectedStatus}
            onChange={setSelectedStatus}
          >
            <Option value="online">在线</Option>
            <Option value="offline">离线</Option>
            <Option value="unknown">未知</Option>
          </Select>
          <Search
            placeholder="搜索主机名或IP"
            allowClear
            style={{ width: 200 }}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
        </Space>

        {/* 操作按钮 */}
        <Space>
          <Button size="small" onClick={handleSelectAll}>
            全选当前
          </Button>
          <Button size="small" onClick={handleClearSelection}>
            清空选择
          </Button>
        </Space>

        {/* 主机列表 */}
        <div style={{ maxHeight: '400px', overflowY: 'auto', border: '1px solid #f0f0f0', borderRadius: 4, padding: 8 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '20px' }}>加载中...</div>
          ) : filteredHosts.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '20px', color: '#999' }}>暂无主机</div>
          ) : (
            filteredHosts.map(host => (
              <div
                key={host.id}
                style={{
                  padding: '8px',
                  borderBottom: '1px solid #f0f0f0',
                  cursor: 'pointer',
                  backgroundColor: selectedHostIds.has(host.id) ? '#e6f7ff' : 'transparent',
                }}
                onClick={() => handleHostToggle(host)}
              >
                <Space>
                  <Checkbox checked={selectedHostIds.has(host.id)} onChange={() => handleHostToggle(host)} />
                  <Text strong>{host.hostname}</Text>
                  <Text type="secondary">{host.ip_address}</Text>
                  <Tag color={getStatusColor(host.status)}>{host.status === 'online' ? '在线' : host.status === 'offline' ? '离线' : '未知'}</Tag>
                  {host.project && <Tag>{host.project.name}</Tag>}
                </Space>
              </div>
            ))
          )}
        </div>
      </Space>
    </Card>
  );
};

export default HostSelector;

