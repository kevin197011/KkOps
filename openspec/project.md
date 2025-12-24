# Project Context

## Purpose
Kronos 是一个基于 Salt 的运维中台管理系统，旨在提供统一的主机管理、SSH 管理、发布管理、定时任务、日志管理、监控管理、权限管理和审计管理等核心功能，后续将扩展 K8s 相关管理能力。

## Tech Stack
- **后端**: Go (推荐使用 Go 1.21+)
- **前端**: React + Ant Design
- **配置管理**: Salt (SaltStack)
- **日志系统**: ELK (Elasticsearch, Logstash, Kibana)
- **监控系统**: Prometheus
- **数据库**: PostgreSQL (推荐) 或 MySQL
- **缓存**: Redis

## Project Conventions

### Code Style
- **Go**: 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
  - 使用 `gofmt` 格式化代码
  - 使用 `golangci-lint` 进行代码检查
  - 包名使用小写，简短有意义
  - 导出函数和类型使用大写开头
- **React/TypeScript**: 
  - 使用 TypeScript 进行类型检查
  - 遵循 React Hooks 最佳实践
  - 组件使用函数式组件
  - 使用 ESLint + Prettier 进行代码格式化

### Architecture Patterns
- **后端架构**: 
  - 采用分层架构：Handler -> Service -> Repository
  - 使用依赖注入模式
  - API 遵循 RESTful 设计规范
  - 使用中间件处理认证、日志、错误处理
- **前端架构**:
  - 使用 React Hooks 和函数式组件
  - 状态管理使用 Context API 或 Redux（如需要）
  - 组件按功能模块组织
  - 使用 Ant Design 作为 UI 组件库

### Testing Strategy
- **单元测试**: 核心业务逻辑覆盖率 ≥ 80%
- **集成测试**: 关键路径必须 100% 覆盖
- **端到端测试**: 主要用户流程需要 E2E 测试
- Go 测试使用 `testing` 包，前端使用 Jest + React Testing Library

### Git Workflow
- 遵循 Conventional Commits 规范
- 主分支：`main`
- 功能分支：`feature/xxx`
- 修复分支：`fix/xxx`
- 通过 Pull Request 合并代码

## Domain Context
- **Salt**: SaltStack 是一个配置管理和远程执行系统，通过 Salt Master 和 Salt Minion 架构实现主机管理
- **运维中台**: 提供统一的运维操作入口，整合各类运维工具和系统
- **主机管理**: 管理服务器资产信息、分组、标签等
- **发布管理**: 基于 Salt 实现应用的自动化部署和发布流程
- **定时任务**: 支持 Salt 命令的定时执行和任务调度

## Important Constraints
- 必须与 Salt Master 集成，通过 Salt API 或 Salt CLI 执行操作
- 需要支持大规模主机管理（1000+ 主机）
- 所有操作需要记录审计日志
- 权限管理需要支持 RBAC（基于角色的访问控制）
- 前端需要支持响应式设计，适配不同屏幕尺寸

## External Dependencies
- **Salt Master**: SaltStack 主控节点，提供配置管理和远程执行能力
- **ELK Stack**: 
  - Elasticsearch: 日志存储和检索
  - Logstash: 日志收集和处理
  - Kibana: 日志可视化和分析
- **Prometheus**: 监控指标收集和存储
- **PostgreSQL/MySQL**: 应用数据存储
- **Redis**: 缓存和会话存储
