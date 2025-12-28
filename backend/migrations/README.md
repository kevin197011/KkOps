# 数据库迁移说明

## 自动迁移

KkOps 使用基于 SQL 文件的迁移系统，在应用启动时自动执行数据库迁移。

### 迁移时机

- **应用启动时自动执行**：每次应用启动时，会自动从 `backend/migrations` 目录读取并执行所有 `.up.sql` 迁移文件
- **按时间戳顺序执行**：迁移文件按文件名（时间戳）排序执行，确保迁移顺序正确
- **无需手动执行**：不需要单独运行迁移脚本

### 迁移文件格式

迁移文件命名格式：`{timestamp}_{description}.up.sql`

- `{timestamp}`: Unix 时间戳（如 `1766942491`）
- `{description}`: 迁移描述（如 `create_initial_tables`）
- `.up.sql`: 升级迁移文件
- `.down.sql`: 回滚迁移文件（可选）

示例：
- `1766942491_create_initial_tables.up.sql` - 创建初始表结构
- `1766942492_create_initial_indexes.up.sql` - 创建初始索引
- `1766942493_fix_formula_created_by_fields.up.sql` - 修复数据

### 迁移内容

迁移文件包含所有数据库表结构的创建和修改：

#### 用户和权限管理
- `users` - 用户表
- `roles` - 角色表
- `permissions` - 权限表
- `user_roles` - 用户角色关联表
- `role_permissions` - 角色权限关联表

#### 项目管理
- `projects` - 项目表
- `project_members` - 项目成员关联表

#### 主机管理
- `hosts` - 主机表
- `host_groups` - 主机组表
- `host_tags` - 主机标签表
- `host_group_members` - 主机组成员关联表
- `host_tag_assignments` - 主机标签关联表

#### SSH管理
- `ssh_connections` - SSH连接配置表
- `ssh_keys` - SSH密钥表
- `ssh_sessions` - SSH会话表

### 默认数据初始化

迁移完成后，会自动初始化以下默认数据：

1. **默认权限**：系统预定义的权限（如 host:read, host:create 等）
2. **默认角色**：
   - `admin` - 系统管理员（拥有所有权限）
   - `operator` - 运维人员（拥有运维相关权限）
   - `viewer` - 查看者（只有只读权限）

### 注意事项

1. **首次启动**：首次启动应用时，会自动创建所有表和默认数据
2. **后续启动**：后续启动时，只会检查表结构差异并更新，不会重复插入默认数据
3. **数据安全**：AutoMigrate 只会添加新字段，不会删除现有字段或数据
4. **生产环境**：建议在生产环境首次部署前，先在测试环境验证迁移结果

### 迁移日志

启动时会输出详细的迁移日志：

```
Connecting to database...
Database connection established
Running database migration...
Starting database migration...
Database migration completed successfully
Initializing default data...
Created 20 default permissions
Created 3 default roles
Default data initialization completed
Database migration completed
```

### 故障排查

如果迁移失败，请检查：

1. **数据库连接**：确认 `.env` 文件中的数据库配置正确
2. **数据库权限**：确认数据库用户有 CREATE TABLE 和 ALTER TABLE 权限
3. **数据库存在**：确认数据库已创建
4. **端口和网络**：确认可以连接到数据库服务器

### 迁移文件位置

所有迁移文件位于 `backend/migrations/` 目录：

```
backend/migrations/
├── 1766942491_create_initial_tables.up.sql
├── 1766942492_create_initial_indexes.up.sql
├── 1766942493_fix_formula_created_by_fields.up.sql
├── 1766942493_fix_formula_created_by_fields.down.sql
├── fix_formula_cascade_delete.sql
├── fix_host_tag_assignments.sql
└── ...
```

### 添加新迁移

1. 创建新的迁移文件，文件名格式：`{timestamp}_{description}.up.sql`
2. 在文件中编写 SQL 语句（使用 `IF NOT EXISTS` 确保幂等性）
3. 可选：创建对应的 `.down.sql` 回滚文件
4. 重启应用，迁移会自动执行

### 手动迁移（可选）

如果需要手动执行 SQL 迁移脚本，可以使用：

```bash
# 执行单个迁移文件
psql -U postgres -d kkops -f backend/migrations/1766942491_create_initial_tables.up.sql

# 或使用 psql 交互式执行
psql -U postgres -d kkops
\i backend/migrations/1766942491_create_initial_tables.up.sql
```

但通常不需要，因为应用启动时会自动执行所有迁移。

