-- 创建所有初始索引
-- 迁移时间: 2025-01-28
-- 说明: 从 database.go 中提取的所有索引创建语句

-- 用户和权限管理索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions(code);
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_type_action ON permissions(resource_type, action);
CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON permissions(deleted_at);

-- 环境管理索引
CREATE INDEX IF NOT EXISTS idx_environments_name ON environments(name);
CREATE INDEX IF NOT EXISTS idx_environments_deleted_at ON environments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_environments_sort_order ON environments(sort_order);
CREATE INDEX IF NOT EXISTS idx_cloud_platforms_name ON cloud_platforms(name);
CREATE INDEX IF NOT EXISTS idx_cloud_platforms_deleted_at ON cloud_platforms(deleted_at);
CREATE INDEX IF NOT EXISTS idx_cloud_platforms_sort_order ON cloud_platforms(sort_order);
CREATE INDEX IF NOT EXISTS idx_host_tags_name ON host_tags(name);
CREATE INDEX IF NOT EXISTS idx_host_tags_deleted_at ON host_tags(deleted_at);

-- 项目管理索引
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id);

-- SSH管理索引
CREATE INDEX IF NOT EXISTS idx_ssh_keys_user_id ON ssh_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_ssh_keys_deleted_at ON ssh_keys(deleted_at);

-- 主机管理索引
CREATE INDEX IF NOT EXISTS idx_host_groups_name ON host_groups(name);
CREATE INDEX IF NOT EXISTS idx_host_groups_deleted_at ON host_groups(deleted_at);
CREATE INDEX IF NOT EXISTS idx_hosts_hostname ON hosts(hostname);
CREATE INDEX IF NOT EXISTS idx_hosts_ip_address ON hosts(ip_address);
CREATE INDEX IF NOT EXISTS idx_hosts_environment ON hosts(environment);
CREATE INDEX IF NOT EXISTS idx_hosts_cloud_platform_id ON hosts(cloud_platform_id);
CREATE INDEX IF NOT EXISTS idx_hosts_project_id ON hosts(project_id);
CREATE INDEX IF NOT EXISTS idx_hosts_status ON hosts(status);
CREATE INDEX IF NOT EXISTS idx_hosts_salt_minion_id ON hosts(salt_minion_id);
CREATE INDEX IF NOT EXISTS idx_hosts_ssh_key_id ON hosts(ssh_key_id);
CREATE INDEX IF NOT EXISTS idx_hosts_deleted_at ON hosts(deleted_at);

-- 部署管理索引
CREATE INDEX IF NOT EXISTS idx_deployment_configs_project_id ON deployment_configs(project_id);
CREATE INDEX IF NOT EXISTS idx_deployment_configs_created_by ON deployment_configs(created_by);
CREATE INDEX IF NOT EXISTS idx_deployment_configs_deleted_at ON deployment_configs(deleted_at);
CREATE INDEX IF NOT EXISTS idx_deployments_project_id ON deployments(project_id);
CREATE INDEX IF NOT EXISTS idx_deployments_config_id ON deployments(config_id);
CREATE INDEX IF NOT EXISTS idx_deployments_deployed_by ON deployments(deployed_by);
CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status);
CREATE INDEX IF NOT EXISTS idx_deployments_deleted_at ON deployments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_deployment_versions_deployment_id ON deployment_versions(deployment_id);

-- 审计日志索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_status ON audit_logs(status);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_deleted_at ON audit_logs(deleted_at);

-- 系统设置索引
CREATE INDEX IF NOT EXISTS idx_system_settings_key ON system_settings(key);
CREATE INDEX IF NOT EXISTS idx_system_settings_category ON system_settings(category);
CREATE INDEX IF NOT EXISTS idx_system_settings_updated_by ON system_settings(updated_by);

-- 批量操作索引
CREATE INDEX IF NOT EXISTS idx_batch_operations_started_by ON batch_operations(started_by);
CREATE INDEX IF NOT EXISTS idx_batch_operations_status ON batch_operations(status);
CREATE INDEX IF NOT EXISTS idx_batch_operations_started_at ON batch_operations(started_at);
CREATE INDEX IF NOT EXISTS idx_batch_operations_created_at ON batch_operations(created_at);
CREATE INDEX IF NOT EXISTS idx_batch_operations_deleted_at ON batch_operations(deleted_at);

-- 命令模板索引
CREATE INDEX IF NOT EXISTS idx_command_templates_created_by ON command_templates(created_by);
CREATE INDEX IF NOT EXISTS idx_command_templates_category ON command_templates(category);
CREATE INDEX IF NOT EXISTS idx_command_templates_is_public ON command_templates(is_public);
CREATE INDEX IF NOT EXISTS idx_command_templates_deleted_at ON command_templates(deleted_at);

-- Formula管理索引
CREATE INDEX IF NOT EXISTS idx_formula_repositories_name ON formula_repositories(name);
CREATE INDEX IF NOT EXISTS idx_formula_repositories_is_active ON formula_repositories(is_active);
CREATE INDEX IF NOT EXISTS idx_formula_repositories_deleted_at ON formula_repositories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_formulas_name ON formulas(name);
CREATE INDEX IF NOT EXISTS idx_formulas_category ON formulas(category);
CREATE INDEX IF NOT EXISTS idx_formulas_repository ON formulas(repository);
CREATE INDEX IF NOT EXISTS idx_formulas_is_active ON formulas(is_active);
CREATE INDEX IF NOT EXISTS idx_formulas_deleted_at ON formulas(deleted_at);
CREATE INDEX IF NOT EXISTS idx_formula_parameters_formula_id ON formula_parameters(formula_id);
CREATE INDEX IF NOT EXISTS idx_formula_templates_formula_id ON formula_templates(formula_id);
CREATE INDEX IF NOT EXISTS idx_formula_templates_created_by ON formula_templates(created_by);
CREATE INDEX IF NOT EXISTS idx_formula_templates_is_public ON formula_templates(is_public);
CREATE INDEX IF NOT EXISTS idx_formula_templates_deleted_at ON formula_templates(deleted_at);
CREATE INDEX IF NOT EXISTS idx_formula_deployments_formula_id ON formula_deployments(formula_id);
CREATE INDEX IF NOT EXISTS idx_formula_deployments_status ON formula_deployments(status);
CREATE INDEX IF NOT EXISTS idx_formula_deployments_started_by ON formula_deployments(started_by);
CREATE INDEX IF NOT EXISTS idx_formula_deployments_started_at ON formula_deployments(started_at);
CREATE INDEX IF NOT EXISTS idx_formula_deployments_deleted_at ON formula_deployments(deleted_at);

