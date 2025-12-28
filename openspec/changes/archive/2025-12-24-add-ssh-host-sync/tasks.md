# Implementation Tasks

## 1. 数据库模型更新
- [x] 1.1 在 Host 模型中增加 SSHPort 字段（类型：int，默认值：22）
- [x] 1.2 更新数据库迁移脚本（AutoMigrate 会自动处理）

## 2. 后端更新
- [x] 2.1 更新 Host 模型，增加 SSHPort 字段
- [x] 2.2 更新 HostService，支持 SSH 端口字段的创建和更新（GORM 自动处理）
- [x] 2.3 更新 HostHandler，确保 SSH 端口字段在 API 响应中返回（GORM 自动处理）

## 3. 前端 Host 管理页面
- [x] 3.1 在 Host 创建/编辑表单中增加 SSH 端口输入字段（InputNumber，默认值 22）
- [x] 3.2 更新 Host 接口类型定义，增加 ssh_port 字段
- [x] 3.3 在 Host 列表表格中显示 SSH 端口（可选）- 暂不显示，保持界面简洁

## 4. 前端 SSH 管理页面
- [x] 4.1 在 SSH 连接创建表单中增加"选择主机"下拉框（从主机管理获取主机列表）
- [x] 4.2 实现主机选择后自动填充 hostname、IP 地址和 SSH 端口
- [x] 4.3 更新 SSH 连接表单，支持手动输入或从主机选择两种方式
- [x] 4.4 更新 SSH 连接接口类型，确保 host_id 字段正确传递

