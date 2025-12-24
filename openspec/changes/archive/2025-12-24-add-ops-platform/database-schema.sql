-- KkOps Operations Platform Database Schema
-- PostgreSQL 12+
-- 
-- This file contains the complete database schema definition
-- for the operations platform. Use this as a reference for
-- creating migration scripts.

-- ============================================
-- 1. 用户和权限管理表
-- ============================================

-- 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_users_status CHECK (status IN ('active', 'disabled', 'deleted'))
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 角色表
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_roles_name ON roles(name);

-- 权限表
CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_resource_type ON permissions(resource_type);

-- 用户角色关联表
CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- 角色权限关联表
CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- ============================================
-- 2. 主机管理表
-- ============================================

-- 主机表
CREATE TABLE hosts (
    id BIGSERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    salt_minion_id VARCHAR(255) UNIQUE,
    os_type VARCHAR(50),
    os_version VARCHAR(100),
    cpu_cores INTEGER,
    memory_gb DECIMAL(10,2),
    disk_gb DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    last_seen_at TIMESTAMP,
    salt_version VARCHAR(50),
    metadata JSONB,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_hosts_status CHECK (status IN ('online', 'offline', 'unknown'))
);

CREATE INDEX idx_hosts_hostname ON hosts(hostname);
CREATE INDEX idx_hosts_ip_address ON hosts(ip_address);
CREATE INDEX idx_hosts_salt_minion_id ON hosts(salt_minion_id);
CREATE INDEX idx_hosts_status ON hosts(status);
CREATE INDEX idx_hosts_deleted_at ON hosts(deleted_at);
CREATE INDEX idx_hosts_metadata ON hosts USING GIN(metadata);

-- 主机组表
CREATE TABLE host_groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_host_groups_name ON host_groups(name);

-- 主机组成员关联表
CREATE TABLE host_group_members (
    id BIGSERIAL PRIMARY KEY,
    host_id BIGINT NOT NULL REFERENCES hosts(id) ON DELETE RESTRICT,
    group_id BIGINT NOT NULL REFERENCES host_groups(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (host_id, group_id)
);

CREATE INDEX idx_host_group_members_host_id ON host_group_members(host_id);
CREATE INDEX idx_host_group_members_group_id ON host_group_members(group_id);

-- 主机标签表
CREATE TABLE host_tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    color VARCHAR(7),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_host_tags_name ON host_tags(name);

-- 主机标签关联表
CREATE TABLE host_tag_assignments (
    id BIGSERIAL PRIMARY KEY,
    host_id BIGINT NOT NULL REFERENCES hosts(id) ON DELETE RESTRICT,
    tag_id BIGINT NOT NULL REFERENCES host_tags(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (host_id, tag_id)
);

CREATE INDEX idx_host_tag_assignments_host_id ON host_tag_assignments(host_id);
CREATE INDEX idx_host_tag_assignments_tag_id ON host_tag_assignments(tag_id);

-- ============================================
-- 3. SSH 管理表
-- ============================================

-- SSH连接配置表
CREATE TABLE ssh_connections (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    host_id BIGINT REFERENCES hosts(id) ON DELETE SET NULL,
    hostname VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL DEFAULT 22,
    username VARCHAR(100) NOT NULL,
    auth_type VARCHAR(20) NOT NULL,
    password_encrypted TEXT,
    key_id BIGINT REFERENCES ssh_keys(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_connected_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_ssh_connections_auth_type CHECK (auth_type IN ('password', 'key')),
    CONSTRAINT chk_ssh_connections_status CHECK (status IN ('active', 'disabled')),
    CONSTRAINT chk_ssh_connections_port CHECK (port > 0 AND port <= 65535)
);

CREATE INDEX idx_ssh_connections_host_id ON ssh_connections(host_id);
CREATE INDEX idx_ssh_connections_key_id ON ssh_connections(key_id);
CREATE INDEX idx_ssh_connections_status ON ssh_connections(status);

-- SSH密钥表
CREATE TABLE ssh_keys (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    key_type VARCHAR(20) NOT NULL,
    private_key_encrypted TEXT NOT NULL,
    public_key TEXT NOT NULL,
    fingerprint VARCHAR(100) UNIQUE,
    passphrase_encrypted TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_ssh_keys_key_type CHECK (key_type IN ('rsa', 'ed25519', 'ecdsa'))
);

CREATE INDEX idx_ssh_keys_fingerprint ON ssh_keys(fingerprint);
CREATE INDEX idx_ssh_keys_key_type ON ssh_keys(key_type);

-- SSH会话表
CREATE TABLE ssh_sessions (
    id BIGSERIAL PRIMARY KEY,
    connection_id BIGINT NOT NULL REFERENCES ssh_connections(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    client_ip VARCHAR(45),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP,
    duration_seconds INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_ssh_sessions_status CHECK (status IN ('active', 'closed', 'timeout'))
);

CREATE INDEX idx_ssh_sessions_connection_id ON ssh_sessions(connection_id);
CREATE INDEX idx_ssh_sessions_user_id ON ssh_sessions(user_id);
CREATE INDEX idx_ssh_sessions_session_token ON ssh_sessions(session_token);
CREATE INDEX idx_ssh_sessions_started_at ON ssh_sessions(started_at);

-- ============================================
-- 4. 发布管理表
-- ============================================

-- 部署配置表
CREATE TABLE deployment_configs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    application_name VARCHAR(100) NOT NULL,
    description TEXT,
    salt_state_files TEXT[] NOT NULL,
    target_groups BIGINT[],
    target_hosts BIGINT[],
    environment VARCHAR(50),
    config_data JSONB,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_deployment_configs_environment CHECK (environment IN ('dev', 'test', 'prod') OR environment IS NULL)
);

CREATE INDEX idx_deployment_configs_name ON deployment_configs(name);
CREATE INDEX idx_deployment_configs_application_name ON deployment_configs(application_name);
CREATE INDEX idx_deployment_configs_created_by ON deployment_configs(created_by);
CREATE INDEX idx_deployment_configs_config_data ON deployment_configs USING GIN(config_data);

-- 部署记录表
CREATE TABLE deployments (
    id BIGSERIAL PRIMARY KEY,
    config_id BIGINT NOT NULL REFERENCES deployment_configs(id) ON DELETE RESTRICT,
    version VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    started_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    duration_seconds INTEGER,
    target_hosts BIGINT[] NOT NULL,
    salt_job_id VARCHAR(255),
    results JSONB,
    error_message TEXT,
    is_rollback BOOLEAN NOT NULL DEFAULT FALSE,
    rollback_from_deployment_id BIGINT REFERENCES deployments(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_deployments_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX idx_deployments_config_id ON deployments(config_id);
CREATE INDEX idx_deployments_started_by ON deployments(started_by);
CREATE INDEX idx_deployments_status ON deployments(status);
CREATE INDEX idx_deployments_started_at ON deployments(started_at);
CREATE INDEX idx_deployments_salt_job_id ON deployments(salt_job_id);

-- 部署版本表
CREATE TABLE deployment_versions (
    id BIGSERIAL PRIMARY KEY,
    application_name VARCHAR(100) NOT NULL,
    version VARCHAR(50) NOT NULL,
    release_notes TEXT,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (application_name, version)
);

CREATE INDEX idx_deployment_versions_application_name ON deployment_versions(application_name);
CREATE INDEX idx_deployment_versions_version ON deployment_versions(version);

-- Salt State文件表
CREATE TABLE salt_states (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    usage_count INTEGER NOT NULL DEFAULT 0,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_salt_states_name ON salt_states(name);
CREATE INDEX idx_salt_states_file_path ON salt_states(file_path);
CREATE INDEX idx_salt_states_checksum ON salt_states(checksum);

-- ============================================
-- 5. 定时任务表
-- ============================================

-- 定时任务表
CREATE TABLE scheduled_tasks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    schedule_type VARCHAR(20) NOT NULL,
    cron_expression VARCHAR(100),
    execute_at TIMESTAMP,
    interval_seconds INTEGER,
    salt_command TEXT,
    salt_state VARCHAR(255),
    target_type VARCHAR(20) NOT NULL,
    target_groups BIGINT[],
    target_hosts BIGINT[],
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    timeout_seconds INTEGER DEFAULT 300,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    last_executed_at TIMESTAMP,
    next_execute_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_scheduled_tasks_schedule_type CHECK (schedule_type IN ('cron', 'once', 'interval')),
    CONSTRAINT chk_scheduled_tasks_target_type CHECK (target_type IN ('all', 'groups', 'hosts')),
    CONSTRAINT chk_scheduled_tasks_timeout CHECK (timeout_seconds > 0)
);

CREATE INDEX idx_scheduled_tasks_name ON scheduled_tasks(name);
CREATE INDEX idx_scheduled_tasks_enabled ON scheduled_tasks(enabled);
CREATE INDEX idx_scheduled_tasks_next_execute_at ON scheduled_tasks(next_execute_at);
CREATE INDEX idx_scheduled_tasks_created_by ON scheduled_tasks(created_by);

-- 任务执行记录表
CREATE TABLE task_executions (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE RESTRICT,
    execution_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'running',
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    duration_seconds INTEGER,
    salt_job_id VARCHAR(255),
    results JSONB,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_task_executions_execution_type CHECK (execution_type IN ('scheduled', 'manual')),
    CONSTRAINT chk_task_executions_status CHECK (status IN ('running', 'success', 'failed', 'timeout'))
);

CREATE INDEX idx_task_executions_task_id ON task_executions(task_id);
CREATE INDEX idx_task_executions_status ON task_executions(status);
CREATE INDEX idx_task_executions_started_at ON task_executions(started_at);
CREATE INDEX idx_task_executions_salt_job_id ON task_executions(salt_job_id);

-- ============================================
-- 6. 监控管理表
-- ============================================

-- 告警规则表
CREATE TABLE alert_rules (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    promql_expression TEXT NOT NULL,
    threshold DECIMAL(10,2),
    severity VARCHAR(20) NOT NULL DEFAULT 'warning',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notification_channels TEXT[],
    last_triggered_at TIMESTAMP,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT chk_alert_rules_severity CHECK (severity IN ('critical', 'warning', 'info'))
);

CREATE INDEX idx_alert_rules_name ON alert_rules(name);
CREATE INDEX idx_alert_rules_enabled ON alert_rules(enabled);
CREATE INDEX idx_alert_rules_severity ON alert_rules(severity);

-- 监控仪表盘表
CREATE TABLE metric_dashboards (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    layout_config JSONB NOT NULL,
    panels JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_metric_dashboards_name ON metric_dashboards(name);
CREATE INDEX idx_metric_dashboards_created_by ON metric_dashboards(created_by);
CREATE INDEX idx_metric_dashboards_layout_config ON metric_dashboards USING GIN(layout_config);

-- ============================================
-- 7. 审计管理表
-- ============================================

-- 审计日志表
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    username VARCHAR(50),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id BIGINT,
    resource_name VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_data JSONB,
    response_data JSONB,
    before_data JSONB,
    after_data JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,
    duration_ms INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_audit_logs_status CHECK (status IN ('success', 'failed')),
    CONSTRAINT chk_audit_logs_action CHECK (action IN ('login', 'logout', 'create', 'update', 'delete', 'execute', 'query', 'export'))
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_username ON audit_logs(username);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_status ON audit_logs(status);

-- ============================================
-- 触发器：自动更新 updated_at 字段
-- ============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为需要的表创建触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_hosts_updated_at BEFORE UPDATE ON hosts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_host_groups_updated_at BEFORE UPDATE ON host_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ssh_connections_updated_at BEFORE UPDATE ON ssh_connections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ssh_keys_updated_at BEFORE UPDATE ON ssh_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deployment_configs_updated_at BEFORE UPDATE ON deployment_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deployments_updated_at BEFORE UPDATE ON deployments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_salt_states_updated_at BEFORE UPDATE ON salt_states
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_scheduled_tasks_updated_at BEFORE UPDATE ON scheduled_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alert_rules_updated_at BEFORE UPDATE ON alert_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_metric_dashboards_updated_at BEFORE UPDATE ON metric_dashboards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 初始化数据
-- ============================================

-- 插入默认权限
INSERT INTO permissions (code, name, resource_type, action, description) VALUES
('host:read', '查看主机', 'host', 'read', '查看主机信息'),
('host:create', '创建主机', 'host', 'create', '创建新主机'),
('host:update', '更新主机', 'host', 'update', '更新主机信息'),
('host:delete', '删除主机', 'host', 'delete', '删除主机'),
('deployment:read', '查看部署', 'deployment', 'read', '查看部署信息'),
('deployment:create', '创建部署', 'deployment', 'create', '创建部署配置'),
('deployment:execute', '执行部署', 'deployment', 'execute', '执行部署操作'),
('task:read', '查看任务', 'task', 'read', '查看定时任务'),
('task:create', '创建任务', 'task', 'create', '创建定时任务'),
('task:execute', '执行任务', 'task', 'execute', '执行定时任务'),
('log:read', '查看日志', 'log', 'read', '查看系统日志'),
('monitoring:read', '查看监控', 'monitoring', 'read', '查看监控指标'),
('audit:read', '查看审计', 'audit', 'read', '查看审计日志'),
('user:manage', '用户管理', 'user', 'manage', '管理用户和权限');

-- 插入默认管理员角色
INSERT INTO roles (name, description) VALUES
('admin', '系统管理员，拥有所有权限'),
('operator', '运维人员，拥有运维相关权限'),
('viewer', '查看者，只能查看信息');

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 为运维人员角色分配运维相关权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE resource_type IN ('host', 'deployment', 'task', 'log', 'monitoring');

-- 为查看者角色分配只读权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE action = 'read';

