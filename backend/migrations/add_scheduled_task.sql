-- 创建定时任务表
CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cron_expression VARCHAR(100) NOT NULL,
    template_id BIGINT REFERENCES task_templates(id) ON DELETE SET NULL,
    content TEXT,
    type VARCHAR(50) DEFAULT 'shell',
    asset_ids TEXT,
    timeout INTEGER DEFAULT 300,
    enabled BOOLEAN DEFAULT false,
    last_run_at TIMESTAMP,
    next_run_at TIMESTAMP,
    last_status VARCHAR(50),
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_enabled ON scheduled_tasks (enabled);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_template_id ON scheduled_tasks (template_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_created_by ON scheduled_tasks (created_by);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_deleted_at ON scheduled_tasks (deleted_at);

-- 给 task_executions 表添加定时任务关联字段
ALTER TABLE task_executions ADD COLUMN IF NOT EXISTS scheduled_task_id BIGINT REFERENCES scheduled_tasks(id) ON DELETE SET NULL;
ALTER TABLE task_executions ADD COLUMN IF NOT EXISTS trigger_type VARCHAR(20) DEFAULT 'manual';

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_task_executions_scheduled_task_id ON task_executions (scheduled_task_id);
CREATE INDEX IF NOT EXISTS idx_task_executions_trigger_type ON task_executions (trigger_type);

-- 修复 task_id 外键约束，允许 NULL 值（定时任务执行时没有关联的普通任务）
ALTER TABLE task_executions ALTER COLUMN task_id DROP NOT NULL;

-- 删除旧的外键约束（如果存在）
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_tasks_executions' AND table_name = 'task_executions'
    ) THEN
        ALTER TABLE task_executions DROP CONSTRAINT fk_tasks_executions;
    END IF;
END $$;

-- 重新添加外键约束，支持 NULL 和 ON DELETE SET NULL
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_task_executions_task' AND table_name = 'task_executions'
    ) THEN
        ALTER TABLE task_executions 
            ADD CONSTRAINT fk_task_executions_task 
            FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL;
    END IF;
END $$;
