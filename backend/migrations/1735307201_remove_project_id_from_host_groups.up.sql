-- 移除主机组与项目的关联，主机组改为纯业务分组
ALTER TABLE host_groups DROP COLUMN IF EXISTS project_id;

-- 添加唯一索引确保名称唯一
CREATE UNIQUE INDEX IF NOT EXISTS uk_host_groups_name ON host_groups(name) WHERE deleted_at IS NULL;
