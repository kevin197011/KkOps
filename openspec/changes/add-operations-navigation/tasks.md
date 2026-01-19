## 1. 数据库设计
- [x] 1.1 创建数据库迁移文件，添加 `operation_tools` 表
  - 字段：id, name, description, category, icon, url, order, enabled, created_at, updated_at
- [x] 1.2 运行数据库迁移（在 `database.go` 中添加了迁移文件到列表和 AutoMigrate）

## 2. 后端实现
- [x] 2.1 创建 `OperationTool` model（`backend/internal/model/operation_tool.go`）
- [x] 2.2 创建 operation tool service（`backend/internal/service/operationtool/service.go`）
  - 实现 CRUD 操作
  - 实现按分类查询
- [x] 2.3 创建 operation tool handler（`backend/internal/handler/operationtool/handler.go`）
  - GET `/api/v1/operation-tools` - 获取工具列表（支持分类过滤）
  - POST `/api/v1/operation-tools` - 创建工具（管理员）
  - PUT `/api/v1/operation-tools/:id` - 更新工具（管理员）
  - DELETE `/api/v1/operation-tools/:id` - 删除工具（管理员）
- [x] 2.4 在 `main.go` 中注册路由和中间件
- [x] 2.5 添加权限控制（`operation-tools:read` 用于查看，`operation-tools:*` 用于管理）

## 3. 前端实现
- [x] 3.1 创建 API client（`frontend/src/api/operationTool.ts`）
- [x] 3.2 在 Dashboard 页面添加"运维导航"组件
  - 显示工具卡片网格
  - 支持分类显示（可选分类标签）
  - 卡片点击打开新标签页
  - 悬停显示描述信息
- [ ] 3.3 创建工具管理页面（管理员功能，可选）
  - 工具列表展示
  - 添加/编辑/删除工具
  - 支持拖拽排序（可选）

## 4. 测试
- [ ] 4.1 后端 API 单元测试
- [ ] 4.2 前端组件测试
- [ ] 4.3 集成测试

## 5. 文档
- [ ] 5.1 更新 API 文档
- [ ] 5.2 更新用户指南（可选）
