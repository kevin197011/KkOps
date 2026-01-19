# Design: 运维导航功能

## Context
运维团队需要快速访问各种运维工具，目前这些工具分散在不同的系统或书签中，访问不便。需要在仪表盘页面提供一个集中的工具导航入口。

## Goals / Non-Goals

### Goals
- 在仪表盘页面提供运维工具快速访问入口
- 支持外部链接，点击在新标签页打开
- 支持工具分类和描述信息
- 管理员可以配置工具列表

### Non-Goals
- 不实现工具内部的嵌入式访问（iframe）
- 不实现工具的单点登录（SSO）集成（初期版本）
- 不实现工具的自动检测和发现

## Decisions

### Decision 1: 数据库存储工具配置
**选择**: 使用数据库表存储工具配置，而非配置文件

**原因**: 
- 支持动态配置，无需重启服务
- 管理员可以在后台管理工具列表
- 支持权限控制和审计日志

**替代方案**: 
- 配置文件：需要重启服务，不够灵活

### Decision 2: 外部链接新标签页打开
**选择**: 工具链接在新标签页打开

**原因**: 
- 保持仪表盘页面不被替换
- 用户可以同时打开多个工具
- 符合常见的链接跳转习惯

**替代方案**: 
- 当前页跳转：会丢失仪表盘上下文

### Decision 3: 分类支持
**选择**: 支持工具分类（可选功能）

**原因**: 
- 工具较多时可以分类组织，便于查找
- 提升用户体验

**分类示例**: 
- 监控工具（Grafana、Prometheus）
- 日志工具（ELK、Loki）
- CI/CD 工具（Jenkins、GitLab CI）
- 数据库工具（phpMyAdmin、DBeaver）
- 其他

### Decision 4: 工具管理权限
**选择**: 只有管理员可以管理工具列表，普通用户只能查看

**原因**: 
- 避免工具链接被误删或修改
- 符合最小权限原则

## Database Schema

```sql
CREATE TABLE operation_tools (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    icon VARCHAR(255),  -- 图标名称或 URL
    url VARCHAR(512) NOT NULL,
    order_index INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_operation_tools_category ON operation_tools(category);
CREATE INDEX idx_operation_tools_enabled ON operation_tools(enabled);
```

## API Design

### GET /api/v1/operation-tools
获取运维工具列表

**查询参数**:
- `category` (可选): 按分类过滤
- `enabled` (可选): 是否只返回启用的工具，默认 true

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Grafana",
      "description": "监控和可视化平台",
      "category": "监控工具",
      "icon": "dashboard",
      "url": "https://grafana.example.com",
      "order": 1,
      "enabled": true
    }
  ]
}
```

### POST /api/v1/operation-tools (管理员)
创建运维工具

**请求体**:
```json
{
  "name": "Grafana",
  "description": "监控和可视化平台",
  "category": "监控工具",
  "icon": "dashboard",
  "url": "https://grafana.example.com",
  "order": 1
}
```

### PUT /api/v1/operation-tools/:id (管理员)
更新运维工具

### DELETE /api/v1/operation-tools/:id (管理员)
删除运维工具

## UI Design

### 仪表盘运维导航区域
- 位置：仪表盘页面底部，在"最近活动"之后
- 布局：网格布局，响应式设计
- 卡片样式：
  - 显示图标、名称、描述（悬停显示）
  - 点击卡片打开链接（新标签页）
  - 支持分类标签（可选）
- 搜索过滤（可选）：顶部搜索框，支持按名称/描述搜索

### 工具管理页面（可选，管理员功能）
- 列表展示所有工具
- 支持添加/编辑/删除
- 支持拖拽排序
- 支持启用/禁用

## Risks / Trade-offs

### Risk 1: 工具链接泄露
**风险**: 工具 URL 可能包含敏感信息

**缓解**: 
- URL 存储时不包含认证信息
- 文档说明需要在工具端配置认证
- 未来可以支持 SSO 集成

### Risk 2: 工具过多导致页面混乱
**风险**: 如果工具过多，页面会变得拥挤

**缓解**: 
- 支持分类组织
- 支持搜索过滤
- 支持分页或懒加载（如果工具数量很大）

### Trade-off: 工具配置方式
**选择**: 使用数据库存储而非配置文件

**优势**: 动态配置，无需重启
**劣势**: 需要额外的数据库表和管理界面

## Migration Plan

### 初始化数据
- 创建数据库迁移文件
- 可选：预置一些常用工具（Grafana、Prometheus 等）
- 迁移脚本自动运行

### 权限设置
- 添加 `operation-tools:read` 权限（所有用户）
- 添加 `operation-tools:*` 权限（管理员）
- 在权限初始化时自动创建

## Open Questions

1. **图标来源**: 使用 Ant Design 图标库还是自定义图标 URL？
   - **建议**: 优先使用 Ant Design 图标，支持自定义图标 URL

2. **工具访问统计**: 是否需要记录工具访问日志？
   - **建议**: 初期版本不实现，未来可以考虑

3. **工具管理界面**: 是否需要在管理后台提供工具管理页面？
   - **建议**: 初期版本可以通过 API 和数据库直接管理，未来可以添加管理界面
