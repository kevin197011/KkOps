-- 修复 formula 相关表的外键约束，添加级联删除
-- 删除 formula 时会自动删除关联的部署记录、参数和模板
-- 执行时间: 2025-12-26

-- 1. 修复 formula_deployments 表
ALTER TABLE formula_deployments DROP CONSTRAINT IF EXISTS formula_deployments_formula_id_fkey;
ALTER TABLE formula_deployments 
ADD CONSTRAINT formula_deployments_formula_id_fkey 
FOREIGN KEY (formula_id) REFERENCES formulas(id) ON DELETE CASCADE;

-- 2. 修复 formula_templates 表
ALTER TABLE formula_templates DROP CONSTRAINT IF EXISTS formula_templates_formula_id_fkey;
ALTER TABLE formula_templates 
ADD CONSTRAINT formula_templates_formula_id_fkey 
FOREIGN KEY (formula_id) REFERENCES formulas(id) ON DELETE CASCADE;

-- 3. 确认 formula_parameters 表（应该已经有 CASCADE）
ALTER TABLE formula_parameters DROP CONSTRAINT IF EXISTS formula_parameters_formula_id_fkey;
ALTER TABLE formula_parameters 
ADD CONSTRAINT formula_parameters_formula_id_fkey 
FOREIGN KEY (formula_id) REFERENCES formulas(id) ON DELETE CASCADE;

-- 验证外键约束（所有 confdeltype 应该都是 'c'）
-- SELECT conname, confdeltype FROM pg_constraint WHERE conname LIKE 'formula_%_formula_id_fkey';
