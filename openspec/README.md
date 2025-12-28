# OpenSpec 文档索引

**最后更新**: 2025-01-28

本文档提供 OpenSpec 目录的导航和索引。

---

## 📚 核心文档

### 项目文档
- **[CURRENT_STATUS_REPORT.md](./CURRENT_STATUS_REPORT.md)** - **项目现状报告**（**深度分析，推荐阅读**）⭐
  - 全面的项目分析报告（804行，22KB）
  - 架构分析、代码质量、功能完成度
  - 安全性、性能、测试覆盖分析
  - 风险评估和改进建议

- **[PROJECT_STATUS.md](./PROJECT_STATUS.md)** - 项目进度总览
  - 完整的项目状态、功能模块、技术栈、数据库状态
  - 最近更新记录、下一步计划

- **[DATA_ALIGNMENT.md](./DATA_ALIGNMENT.md)** - 文档数据对齐说明
  - 文档数据对齐验证
  - 数据源优先级和对齐规则
  - 更新流程和维护说明

- **[project.md](./project.md)** - 项目上下文
  - 项目目的、技术栈、项目结构
  - 开发规范和约定

- **[AGENTS.md](./AGENTS.md)** - AI助手使用指南
  - OpenSpec 工作流程
  - 变更提案创建和实施指南

---

## 📋 变更管理

### 变更索引
- **[changes/CHANGES_INDEX.md](./changes/CHANGES_INDEX.md)** - 变更提案索引
  - 所有变更提案的状态追踪
  - 已完成、待实施、已归档的提案列表

- **[changes/PROGRESS_TRACKING.md](./changes/PROGRESS_TRACKING.md)** - 进度追踪
  - 项目里程碑
  - 模块完成度追踪
  - 代码统计

- **[REQUIREMENTS_AND_TASKS.md](./REQUIREMENTS_AND_TASKS.md)** - 需求与任务总览
  - 需求状态总览（已完成/待实现）
  - 任务状态统计
  - 归档状态和建议

### 变更提案目录结构

```
changes/
├── CHANGES_INDEX.md              # 变更提案索引
├── PROGRESS_TRACKING.md          # 进度追踪
├── {change-id}/                  # 变更提案目录
│   ├── proposal.md              # 变更提案
│   ├── design.md                # 设计文档（可选）
│   ├── tasks.md                 # 任务清单
│   └── specs/                   # 规范文档（可选）
│       └── {module}/
│           └── spec.md
└── archive/                      # 已归档的提案
```

---

## 📖 功能规范

### 已实现功能规范

- **[specs/user-management/spec.md](./specs/user-management/spec.md)** - 用户管理规范
- **[specs/host-management/spec.md](./specs/host-management/spec.md)** - 主机管理规范
- **[specs/ssh-management/spec.md](./specs/ssh-management/spec.md)** - SSH管理规范
- **[specs/webssh-management/spec.md](./specs/webssh-management/spec.md)** - WebSSH管理规范
- **[specs/batch-operations/spec.md](./specs/batch-operations/spec.md)** - 批量操作规范
- **[specs/formula-management/spec.md](./specs/formula-management/spec.md)** - Formula管理规范
- **[specs/system-settings/spec.md](./specs/system-settings/spec.md)** - 系统设置规范
- **[specs/audit-management/spec.md](./specs/audit-management/spec.md)** - 审计管理规范

---

## 🔍 快速查找

### 按状态查找变更提案

**已完成** ✅:
- `add-role-permission-management`
- `add-host-group-tag-management`
- `add-host-environment`
- `add-ssh-key-management`
- `add-ssh-host-sync`
- `add-system-settings-page`
- `add-salt-api-test-connection`
- `add-batch-operations-page`
- `add-deployment-management`
- `refactor-execution-ui`
- `simplify-ssh-to-use-hosts`
- `refactor-database-migrations` (2025-01-28)

**待实施** ⏸️:
- `open-webssh-in-modal`
- `refactor-webssh-to-tree-layout`
- `separate-ssh-key-management-page`
- `separate-webssh-management`
- `enhance-host-tag-management`
- `auto-generate-ssh-connections`
- `refactor-ssh-to-webssh`
- `remove-ssh-add-webssh`

### 按功能模块查找

**用户和权限**:
- `add-role-permission-management` ✅

**主机管理**:
- `add-host-group-tag-management` ✅
- `add-host-environment` ✅
- `enhance-host-tag-management` ⏸️

**SSH/WebSSH**:
- `add-ssh-key-management` ✅
- `add-ssh-host-sync` ✅
- `simplify-ssh-to-use-hosts` ✅
- `open-webssh-in-modal` ⏸️
- `refactor-webssh-to-tree-layout` ⏸️

**批量操作和部署**:
- `add-batch-operations-page` ✅
- `add-deployment-management` ✅
- `refactor-execution-ui` ✅

**系统管理**:
- `add-system-settings-page` ✅
- `add-salt-api-test-connection` ✅

**数据库**:
- `refactor-database-migrations` ✅ (2025-01-28)

---

## 📊 项目统计

### 代码统计
- **后端Handler**: 19个
- **后端Service**: 18个
- **后端Repository**: 16个
- **后端Model**: 13个
- **前端页面**: 19个
- **前端组件**: 10个
- **前端Service**: 18个
- **数据库表**: 25个
- **迁移文件**: 11个

### 功能完成度
- **核心功能模块**: 10个（100%完成）
- **API端点**: 100+个
- **前端路由**: 19个
- **变更提案**: 24个（12个已完成）

---

## 🎯 使用指南

### 查看项目进度
1. 阅读 [PROJECT_STATUS.md](./PROJECT_STATUS.md) 了解整体进度
2. 查看 [changes/PROGRESS_TRACKING.md](./changes/PROGRESS_TRACKING.md) 了解详细追踪

### 查找变更提案
1. 查看 [changes/CHANGES_INDEX.md](./changes/CHANGES_INDEX.md) 获取所有提案列表
2. 进入对应的 `changes/{change-id}/` 目录查看详情

### 创建新变更提案
1. 参考 [AGENTS.md](./AGENTS.md) 中的工作流程
2. 在 `changes/` 目录下创建新文件夹
3. 创建 `proposal.md`, `tasks.md`, `design.md`（可选）
4. 更新 [changes/CHANGES_INDEX.md](./changes/CHANGES_INDEX.md)

---

## 📝 文档维护

### 更新频率
- **PROJECT_STATUS.md**: 每次重大功能完成后更新
- **CHANGES_INDEX.md**: 每次变更提案状态变更时更新
- **PROGRESS_TRACKING.md**: 每周或每次里程碑完成后更新

### 维护者
- KkOps开发团队

---

**最后更新**: 2025-01-28

