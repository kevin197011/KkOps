-- 修复 Formula 相关表的 created_by 字段
-- 迁移时间: 2025-01-28
-- 说明: 从 database.go 的 fixFormulaCreatedByFields() 函数中提取

-- 删除 name 为空字符串的 Formula 记录（脏数据）
DELETE FROM formulas WHERE name = '' OR name IS NULL;

-- 修复 formulas 表的 created_by 字段
UPDATE formulas SET created_by = 1 WHERE created_by = 0 OR created_by IS NULL;

-- 修复 formula_repositories 表的 created_by 字段
UPDATE formula_repositories SET created_by = 1 WHERE created_by = 0 OR created_by IS NULL;

-- 修复 formula_templates 表的 created_by 字段
UPDATE formula_templates SET created_by = 1 WHERE created_by = 0 OR created_by IS NULL;

-- 修复 formula_deployments 表的 started_by 字段
UPDATE formula_deployments SET started_by = 1 WHERE started_by = 0 OR started_by IS NULL;

