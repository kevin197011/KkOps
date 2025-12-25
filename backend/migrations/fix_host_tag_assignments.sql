-- 修复 host_tag_assignments 表结构
-- 1. 检查并删除不存在的 host_tag_id 列（如果存在）
-- 2. 确保唯一约束存在

DO $$
BEGIN
    -- 删除错误的 host_tag_id 列（如果存在）
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'host_tag_assignments' 
        AND column_name = 'host_tag_id'
    ) THEN
        ALTER TABLE host_tag_assignments DROP COLUMN host_tag_id;
        RAISE NOTICE 'Dropped host_tag_id column from host_tag_assignments';
    END IF;
    
    -- 添加唯一约束（如果不存在）
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_constraint 
        WHERE conrelid = 'host_tag_assignments'::regclass 
        AND conname = 'host_tag_assignments_host_id_tag_id_key'
    ) THEN
        ALTER TABLE host_tag_assignments 
        ADD CONSTRAINT host_tag_assignments_host_id_tag_id_key 
        UNIQUE (host_id, tag_id);
        RAISE NOTICE 'Added unique constraint on (host_id, tag_id)';
    END IF;
END $$;

-- 验证表结构
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'host_tag_assignments'
ORDER BY ordinal_position;
