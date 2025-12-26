-- 回滚：恢复 project_id 列
DROP INDEX IF EXISTS uk_host_groups_name;
ALTER TABLE host_groups ADD COLUMN IF NOT EXISTS project_id BIGINT;
