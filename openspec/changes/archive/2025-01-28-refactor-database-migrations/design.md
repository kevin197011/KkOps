# 数据库迁移系统重构设计

**创建日期**: 2025-01-28  
**状态**: ✅ 已完成

## 设计目标

1. **统一管理**: 所有 SQL 变更语句统一到 `backend/migrations/` 目录
2. **版本控制**: 使用时间戳命名，支持迁移顺序管理
3. **自动执行**: 应用启动时自动执行所有迁移文件
4. **幂等性**: 所有迁移文件可重复执行，不会出错

## 架构设计

### 迁移文件结构

```
backend/migrations/
├── {timestamp}_{description}.up.sql    # 升级迁移文件
├── {timestamp}_{description}.down.sql  # 回滚迁移文件（可选）
└── {legacy_name}.sql                  # 传统格式迁移文件（向后兼容）
```

### 迁移文件命名规范

- **时间戳格式**: `{timestamp}_{description}.up.sql`
  - `{timestamp}`: Unix 时间戳（如 `1766942491`）
  - `{description}`: 迁移描述（如 `create_initial_tables`）
  - `.up.sql`: 升级迁移文件
  - `.down.sql`: 回滚迁移文件（可选）

- **传统格式**: `{name}.sql`（向后兼容）
  - 如 `fix_formula_cascade_delete.sql`

### 执行流程

```
应用启动
  ↓
AutoMigrate()
  ↓
executeMigrations()
  ↓
读取 backend/migrations/*.up.sql
  ↓
按时间戳排序
  ↓
逐文件执行 SQL
  ↓
执行传统格式迁移文件
  ↓
initDefaultData() (初始化默认数据)
```

## 实现细节

### 1. 迁移文件读取

```go
// 读取时间戳格式的迁移文件
files, err := filepath.Glob(filepath.Join(migrationsDir, "*_*.up.sql"))
sort.Strings(files) // 按时间戳排序

// 读取传统格式的迁移文件
legacyFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
```

### 2. SQL 执行

```go
// 按分号分割 SQL 语句
statements := strings.Split(string(sqlContent), ";")
for _, stmt := range statements {
    stmt = strings.TrimSpace(stmt)
    // 跳过空语句和注释
    if stmt == "" || strings.HasPrefix(stmt, "--") {
        continue
    }
    // 执行 SQL
    DB.Exec(stmt)
}
```

### 3. 错误处理

- **已存在错误**: 只记录警告，不中断迁移（`IF NOT EXISTS` 保护）
- **其他错误**: 返回错误，中断迁移流程

## 迁移文件示例

### 表创建迁移

```sql
-- 创建所有初始表结构
-- 迁移时间: 2025-01-28

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    ...
);
```

### 索引创建迁移

```sql
-- 创建所有初始索引
-- 迁移时间: 2025-01-28

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

### 数据修复迁移

```sql
-- 修复 Formula 相关表的 created_by 字段
-- 迁移时间: 2025-01-28

UPDATE formulas SET created_by = 1 WHERE created_by = 0 OR created_by IS NULL;
```

## 向后兼容性

1. **支持传统格式**: 继续支持 `fix_*.sql` 格式的迁移文件
2. **自动执行**: 传统格式文件在时间戳文件之后执行
3. **错误容忍**: 对已存在的表/索引只记录警告

## 优势

1. ✅ **代码简化**: 删除约 650 行硬编码 SQL
2. ✅ **维护性**: 数据库结构变更通过文件管理
3. ✅ **版本控制**: 每个迁移文件都有明确的时间戳和描述
4. ✅ **可追溯性**: 迁移历史清晰可见
5. ✅ **回滚支持**: 支持 `.down.sql` 回滚文件

## 注意事项

1. **幂等性**: 所有迁移文件必须使用 `IF NOT EXISTS` 确保可重复执行
2. **顺序**: 迁移文件按时间戳排序执行，确保依赖关系正确
3. **测试**: 在生产环境部署前，必须在测试环境验证迁移结果

