# Change: Add SSH Host Sync and SSH Port to Host Management

## Why
当前 SSH 连接管理需要手动输入主机信息（hostname、IP等），无法直接从主机管理中选择已存在的主机。同时，主机管理中缺少 SSH 端口字段，导致在创建 SSH 连接时需要重复输入端口信息。需要实现：
1. SSH 连接创建时可以从主机管理中选择主机，自动同步主机信息
2. 主机管理中增加 SSH 端口字段，便于统一管理

## What Changes
- **主机管理增强**: 在 Host 模型中增加 `SSHPort` 字段（默认 22）
- **SSH 连接优化**: SSH 连接创建表单支持从主机管理中选择主机，自动填充 hostname、IP、SSH 端口等信息
- **数据同步**: 当选择主机时，自动同步主机的 hostname、IP 地址和 SSH 端口到 SSH 连接表单

## Impact
- **Affected specs**: 
  - `host-management` (MODIFIED - 增加 SSH 端口字段)
  - `ssh-management` (MODIFIED - 支持从主机管理选择主机)
- **Affected code**: 
  - 数据库: Host 模型增加 `ssh_port` 字段
  - 后端: Host 模型、Repository、Service、Handler 需要支持 SSH 端口字段
  - 前端: Host 管理页面增加 SSH 端口输入字段
  - 前端: SSH 连接创建表单增加主机选择下拉框，支持从主机管理同步信息

