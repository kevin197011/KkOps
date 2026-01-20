# Change: 完善 Swagger API 文档

## Why

当前 Swagger API 文档不完整，许多接口缺少 Swagger 注释，导致：
- API 文档无法完整反映系统的所有功能
- 开发人员无法通过 Swagger UI 快速了解和使用 API
- 第三方集成时缺少完整的接口规范
- API 的使用和维护效率低下

## What Changes

### 需要添加 Swagger 注释的 Handler

1. **标签管理** (`tag/handler.go`)
   - CreateTag, GetTag, ListTags, UpdateTag, DeleteTag

2. **分类管理** (`category/handler.go`)
   - CreateCategory, GetCategory, ListCategories, UpdateCategory, DeleteCategory

3. **资产管理** (`asset/handler.go`)
   - CreateAsset, GetAsset, ListAssets, UpdateAsset, DeleteAsset
   - ExportAssets, ImportAssets

4. **角色管理** (`role/handler.go`)
   - CreateRole, GetRole, ListRoles, UpdateRole, DeleteRole
   - AssignRoleToUser, AssignPermissionToRole, ListPermissions, GetRolePermissions, RemovePermissionFromRole

5. **项目管理** (`project/handler.go`)
   - CreateProject, GetProject, ListProjects, UpdateProject, DeleteProject

6. **环境管理** (`environment/handler.go`)
   - CreateEnvironment, GetEnvironment, ListEnvironments, UpdateEnvironment, DeleteEnvironment

7. **云平台管理** (`cloudplatform/handler.go`)
   - CreateCloudPlatform, GetCloudPlatform, ListCloudPlatforms, UpdateCloudPlatform, DeleteCloudPlatform

8. **SSH 密钥管理** (`sshkey/handler.go`)
   - ListSSHKeys, GetSSHKey, CreateSSHKey, UpdateSSHKey, DeleteSSHKey, TestSSHKey

9. **任务管理** (`task/handler.go`)
   - CreateTemplate, GetTemplate, ListTemplates, UpdateTemplate, DeleteTemplate
   - CreateTask, GetTask, ListTasks, UpdateTask, DeleteTask
   - ExecuteTask, GetTaskExecutions, GetTaskExecution, CancelTaskExecution, CancelTask, GetTaskExecutionLogs

10. **认证 API Token 管理** (`auth/handler.go`)
    - CreateAPIToken, ListAPITokens, RevokeAPIToken

### Swagger 注释要求

每个接口需要包含以下标准注释：
- `@Summary`: 接口简要描述
- `@Description`: 详细描述（可选）
- `@Tags`: 标签分类
- `@Accept`: 接受的请求格式（json）
- `@Produce`: 返回格式（json）
- `@Security`: 安全认证（如果需要）
- `@Param`: 参数说明（路径参数、查询参数、请求体）
- `@Success`: 成功响应
- `@Failure`: 失败响应
- `@Router`: 路由路径和方法

### 完善现有注释

部分接口已有注释但不完整，需要：
- 补充缺失的参数说明
- 完善响应示例
- 统一注释格式和风格
- 确保所有响应状态码都有文档

## Impact

- **受影响的文件**: 所有 handler 文件（约 15 个文件）
- **新增内容**: 约 80+ 个接口的完整 Swagger 注释
- **文档生成**: 需要重新运行 `swag init` 生成完整的 API 文档
- **维护**: 后续新增接口时需同步添加 Swagger 注释

## 质量标准

- 所有接口必须有完整的 Swagger 注释
- 注释格式符合 swaggo/swag 规范
- 参数类型和响应类型必须准确
- 错误响应必须完整列出所有可能的错误码
- 安全认证要求必须明确标注
