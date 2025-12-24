# Change: Add Operations Platform

## Why
当前运维工作缺乏统一的管理平台，各类运维操作分散在不同工具中，效率低下且难以统一管理。需要构建一个基于 Salt 的运维中台管理系统，整合主机管理、SSH 管理、发布管理、定时任务、日志管理、监控管理、权限管理和审计管理等核心功能，提供统一的运维操作入口，提升运维效率和可追溯性。

## What Changes
- **新增主机管理能力**: 支持主机信息管理、分组、标签、状态查询等功能
- **新增 SSH 管理能力**: 提供 SSH 连接管理、密钥管理、会话管理等功能
- **新增发布管理能力**: 基于 Salt 实现应用的自动化部署和发布流程
- **新增定时任务能力**: 支持 Salt 命令的定时执行和任务调度
- **新增日志管理能力**: 对接 ELK，提供日志查询、分析和可视化
- **新增监控管理能力**: 对接 Prometheus，提供监控指标查询和告警管理
- **新增权限管理能力**: 实现 RBAC 权限控制，支持用户、角色、权限管理
- **新增审计管理能力**: 记录所有操作日志，支持审计查询和追溯
- **预留 K8s 管理能力**: 为后续 K8s 相关功能预留接口和规范

## Impact
- **Affected specs**: 
  - `host-management` (新增)
  - `ssh-management` (新增)
  - `deployment-management` (新增)
  - `scheduled-tasks` (新增)
  - `log-management` (新增)
  - `monitoring-management` (新增)
  - `permission-management` (新增)
  - `audit-management` (新增)
  - `k8s-management` (新增，后续完善)
- **Affected code**: 
  - 后端: 需要创建完整的 Go 项目结构，包括 API 层、Service 层、Repository 层
  - 前端: 需要创建 React 项目结构，包括路由、组件、状态管理
  - 数据库: 需要设计并创建相关数据表
  - 集成: 需要实现与 Salt、ELK、Prometheus 的集成接口

