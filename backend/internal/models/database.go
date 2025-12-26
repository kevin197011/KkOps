package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// manualMigrateUser 手动迁移 User 表，避免 GORM AutoMigrate 的 pgx/v5 兼容性问题
func manualMigrateUser() error {
	log.Println("Manually migrating User table...")

	// 检查表是否已存在
	var exists bool
	err := DB.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = CURRENT_SCHEMA()
			AND table_name = 'users'
		)
	`).Scan(&exists).Error

	if err != nil {
		return fmt.Errorf("failed to check if users table exists: %w", err)
	}

	if exists {
		log.Println("Users table already exists, skipping manual migration")
		return nil
	}

	// 手动创建 users 表
	err = DB.Exec(`
		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			display_name VARCHAR(100),
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			last_login_at TIMESTAMP,
			last_login_ip VARCHAR(45),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		)
	`).Error

	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// 创建索引
	err = DB.Exec(`
		CREATE UNIQUE INDEX idx_users_username ON users(username);
		CREATE UNIQUE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_deleted_at ON users(deleted_at);
	`).Error

	if err != nil {
		return fmt.Errorf("failed to create indexes for users table: %w", err)
	}

	log.Println("User table manually migrated successfully")
	return nil
}

// manualMigrateAllTables 手动迁移所有表，使用纯SQL避免pgx/v5兼容性问题
func manualMigrateAllTables() error {
	log.Println("Starting manual migration of all tables...")

	tables := []struct {
		name string
		sql  string
	}{
		{
			name: "users",
			sql: `
				CREATE TABLE IF NOT EXISTS users (
					id BIGSERIAL PRIMARY KEY,
					username VARCHAR(50) NOT NULL UNIQUE,
					email VARCHAR(100) NOT NULL UNIQUE,
					password_hash VARCHAR(255) NOT NULL,
					display_name VARCHAR(100),
					status VARCHAR(20) NOT NULL DEFAULT 'active',
					last_login_at TIMESTAMP,
					last_login_ip VARCHAR(45),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "roles",
			sql: `
				CREATE TABLE IF NOT EXISTS roles (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(50) NOT NULL UNIQUE,
					display_name VARCHAR(100),
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "permissions",
			sql: `
				CREATE TABLE IF NOT EXISTS permissions (
					id BIGSERIAL PRIMARY KEY,
					code VARCHAR(100) NOT NULL UNIQUE,
					name VARCHAR(100) NOT NULL,
					resource_type VARCHAR(50) NOT NULL,
					action VARCHAR(50) NOT NULL,
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "user_roles",
			sql: `
				CREATE TABLE IF NOT EXISTS user_roles (
					user_id BIGINT NOT NULL,
					role_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (user_id, role_id),
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
					FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "role_permissions",
			sql: `
				CREATE TABLE IF NOT EXISTS role_permissions (
					role_id BIGINT NOT NULL,
					permission_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (role_id, permission_id),
					FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
					FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "environments",
			sql: `
				CREATE TABLE IF NOT EXISTS environments (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(50) NOT NULL UNIQUE,
					display_name VARCHAR(100),
					description TEXT,
					color VARCHAR(20),
					sort_order INTEGER DEFAULT 0,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "cloud_platforms",
			sql: `
				CREATE TABLE IF NOT EXISTS cloud_platforms (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(50) NOT NULL UNIQUE,
					display_name VARCHAR(100),
					description TEXT,
					icon VARCHAR(100),
					color VARCHAR(20),
					sort_order INTEGER DEFAULT 0,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "host_tags",
			sql: `
				CREATE TABLE IF NOT EXISTS host_tags (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(50) NOT NULL UNIQUE,
					color VARCHAR(20),
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "projects",
			sql: `
				CREATE TABLE IF NOT EXISTS projects (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(100) NOT NULL UNIQUE,
					description TEXT,
					status VARCHAR(20) NOT NULL DEFAULT 'active',
					created_by BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (created_by) REFERENCES users(id)
				)`,
		},
		{
			name: "ssh_keys",
			sql: `
				CREATE TABLE IF NOT EXISTS ssh_keys (
					id BIGSERIAL PRIMARY KEY,
					user_id BIGINT NOT NULL,
					name VARCHAR(100) NOT NULL,
					username VARCHAR(100),
					key_type VARCHAR(20) NOT NULL,
					public_key TEXT,
					private_key TEXT NOT NULL,
					fingerprint VARCHAR(100),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id)
				)`,
		},
		{
			name: "host_groups",
			sql: `
				CREATE TABLE IF NOT EXISTS host_groups (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(100) NOT NULL UNIQUE,
					display_name VARCHAR(200),
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)`,
		},
		{
			name: "hosts",
			sql: `
				CREATE TABLE IF NOT EXISTS hosts (
					id BIGSERIAL PRIMARY KEY,
					project_id BIGINT NOT NULL,
					hostname VARCHAR(255) NOT NULL,
					ip_address VARCHAR(45) NOT NULL UNIQUE,
					salt_minion_id VARCHAR(255),
					os_type VARCHAR(50),
					os_version VARCHAR(100),
					cpu_cores INTEGER,
					memory_gb DECIMAL(10,2),
					disk_gb DECIMAL(10,2),
					status VARCHAR(20) NOT NULL DEFAULT 'unknown',
					environment VARCHAR(20),
					cloud_platform_id BIGINT,
					ssh_port INTEGER NOT NULL DEFAULT 22,
					ssh_key_id BIGINT,
					last_seen_at TIMESTAMP,
					salt_version VARCHAR(50),
					metadata JSONB,
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id),
					FOREIGN KEY (cloud_platform_id) REFERENCES cloud_platforms(id),
					FOREIGN KEY (ssh_key_id) REFERENCES ssh_keys(id)
				)`,
		},
		{
			name: "host_group_members",
			sql: `
				CREATE TABLE IF NOT EXISTS host_group_members (
					host_id BIGINT NOT NULL,
					group_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (host_id, group_id),
					FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE,
					FOREIGN KEY (group_id) REFERENCES host_groups(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "host_tag_assignments",
			sql: `
				CREATE TABLE IF NOT EXISTS host_tag_assignments (
					host_id BIGINT NOT NULL,
					tag_id BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (host_id, tag_id),
					FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE,
					FOREIGN KEY (tag_id) REFERENCES host_tags(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "project_members",
			sql: `
				CREATE TABLE IF NOT EXISTS project_members (
					id BIGSERIAL PRIMARY KEY,
					project_id BIGINT NOT NULL,
					user_id BIGINT NOT NULL,
					role VARCHAR(50) DEFAULT 'member',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "deployment_configs",
			sql: `
				CREATE TABLE IF NOT EXISTS deployment_configs (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(100) NOT NULL,
					project_id BIGINT NOT NULL,
					description TEXT,
					config_type VARCHAR(50) NOT NULL,
					config_data JSONB,
					created_by BIGINT NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id),
					FOREIGN KEY (created_by) REFERENCES users(id)
				)`,
		},
		{
			name: "deployments",
			sql: `
				CREATE TABLE IF NOT EXISTS deployments (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(100) NOT NULL,
					project_id BIGINT NOT NULL,
					config_id BIGINT NOT NULL,
					status VARCHAR(50) NOT NULL DEFAULT 'pending',
					version VARCHAR(100),
					deployed_by BIGINT NOT NULL,
					deployed_at TIMESTAMP,
					completed_at TIMESTAMP,
					error_message TEXT,
					logs TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id),
					FOREIGN KEY (config_id) REFERENCES deployment_configs(id),
					FOREIGN KEY (deployed_by) REFERENCES users(id)
				)`,
		},
		{
			name: "deployment_versions",
			sql: `
				CREATE TABLE IF NOT EXISTS deployment_versions (
					id BIGSERIAL PRIMARY KEY,
					deployment_id BIGINT NOT NULL,
					version VARCHAR(100) NOT NULL,
					config_snapshot JSONB,
					status VARCHAR(50) NOT NULL,
					deployed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "audit_logs",
			sql: `
				CREATE TABLE IF NOT EXISTS audit_logs (
					id BIGSERIAL PRIMARY KEY,
					user_id BIGINT,
					username VARCHAR(100),
					action VARCHAR(100) NOT NULL,
					resource_type VARCHAR(100) NOT NULL,
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
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id)
				)`,
		},
		{
			name: "system_settings",
			sql: `
				CREATE TABLE IF NOT EXISTS system_settings (
					id BIGSERIAL PRIMARY KEY,
					key VARCHAR(100) NOT NULL UNIQUE,
					value TEXT NOT NULL,
					category VARCHAR(50) NOT NULL,
					description TEXT,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_by BIGINT,
					FOREIGN KEY (updated_by) REFERENCES users(id)
				)`,
		},
		{
			name: "batch_operations",
			sql: `
				CREATE TABLE IF NOT EXISTS batch_operations (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					command_type VARCHAR(50) NOT NULL,
					command_function VARCHAR(255) NOT NULL,
					command_args JSONB,
					target_hosts JSONB NOT NULL,
					target_count INTEGER NOT NULL,
					status VARCHAR(20) NOT NULL DEFAULT 'pending',
					salt_job_id VARCHAR(255),
					started_by BIGINT NOT NULL,
					started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					completed_at TIMESTAMP,
					duration_seconds INTEGER,
					results JSONB,
					success_count INTEGER DEFAULT 0,
					failed_count INTEGER DEFAULT 0,
					error_message TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (started_by) REFERENCES users(id)
				)`,
		},
		{
			name: "command_templates",
			sql: `
				CREATE TABLE IF NOT EXISTS command_templates (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					category VARCHAR(100),
					command_function VARCHAR(100) NOT NULL,
					command_args JSONB,
					icon VARCHAR(100),
					created_by BIGINT NOT NULL,
					is_public BOOLEAN DEFAULT false,
					usage_count INTEGER DEFAULT 0,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (created_by) REFERENCES users(id)
				)`,
		},
		{
			name: "formula_repositories",
			sql: `
				CREATE TABLE IF NOT EXISTS formula_repositories (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					url VARCHAR(500) NOT NULL,
					branch VARCHAR(100) DEFAULT 'master',
					local_path VARCHAR(500),
					is_active BOOLEAN DEFAULT true,
					last_sync_at TIMESTAMP,
					created_by BIGINT NOT NULL DEFAULT 1,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (created_by) REFERENCES users(id)
				)`,
		},
		{
			name: "formulas",
			sql: `
				CREATE TABLE IF NOT EXISTS formulas (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					description TEXT,
					category VARCHAR(100),
					version VARCHAR(50),
					path VARCHAR(500) NOT NULL,
					repository VARCHAR(500) NOT NULL,
					icon VARCHAR(100),
					tags JSONB,
					metadata JSONB,
					is_active BOOLEAN DEFAULT true,
					created_by BIGINT NOT NULL DEFAULT 1,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (created_by) REFERENCES users(id)
				)`,
		},
		{
			name: "formula_parameters",
			sql: `
				CREATE TABLE IF NOT EXISTS formula_parameters (
					id BIGSERIAL PRIMARY KEY,
					formula_id BIGINT NOT NULL,
					name VARCHAR(255) NOT NULL,
					type VARCHAR(50) NOT NULL,
					default_value JSONB,
					required BOOLEAN DEFAULT false,
					label VARCHAR(255),
					description TEXT,
					validation JSONB,
					order_index INTEGER DEFAULT 0,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (formula_id) REFERENCES formulas(id) ON DELETE CASCADE
				)`,
		},
		{
			name: "formula_templates",
			sql: `
				CREATE TABLE IF NOT EXISTS formula_templates (
					id BIGSERIAL PRIMARY KEY,
					formula_id BIGINT NOT NULL,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					pillar_data JSONB,
					is_public BOOLEAN DEFAULT false,
					created_by BIGINT NOT NULL DEFAULT 1,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (formula_id) REFERENCES formulas(id) ON DELETE CASCADE,
					FOREIGN KEY (created_by) REFERENCES users(id)
				)`,
		},
		{
			name: "formula_deployments",
			sql: `
				CREATE TABLE IF NOT EXISTS formula_deployments (
					id BIGSERIAL PRIMARY KEY,
					formula_id BIGINT NOT NULL,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					target_hosts JSONB,
					pillar_data JSONB,
					status VARCHAR(20) DEFAULT 'pending',
					salt_job_id VARCHAR(100),
					results JSONB,
					success_count INTEGER DEFAULT 0,
					failed_count INTEGER DEFAULT 0,
					error_message TEXT,
					started_by BIGINT NOT NULL DEFAULT 1,
					started_at TIMESTAMP,
					completed_at TIMESTAMP,
					duration_seconds INTEGER,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (formula_id) REFERENCES formulas(id),
					FOREIGN KEY (started_by) REFERENCES users(id)
				)`,
		},
	}

	// 执行所有表的创建
	for _, table := range tables {
		log.Printf("Creating table: %s", table.name)
		if err := DB.Exec(table.sql).Error; err != nil {
			return fmt.Errorf("failed to create table %s: %w", table.name, err)
		}
		log.Printf("Table %s created successfully", table.name)
	}

	// 创建索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name)",
		"CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions(code)",
		"CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name)",
		"CREATE INDEX IF NOT EXISTS idx_permissions_resource_type_action ON permissions(resource_type, action)",
		"CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON permissions(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_environments_name ON environments(name)",
		"CREATE INDEX IF NOT EXISTS idx_environments_deleted_at ON environments(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_environments_sort_order ON environments(sort_order)",
		"CREATE INDEX IF NOT EXISTS idx_cloud_platforms_name ON cloud_platforms(name)",
		"CREATE INDEX IF NOT EXISTS idx_cloud_platforms_deleted_at ON cloud_platforms(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_cloud_platforms_sort_order ON cloud_platforms(sort_order)",
		"CREATE INDEX IF NOT EXISTS idx_host_tags_name ON host_tags(name)",
		"CREATE INDEX IF NOT EXISTS idx_host_tags_deleted_at ON host_tags(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name)",
		"CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)",
		"CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by)",
		"CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_ssh_keys_user_id ON ssh_keys(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_ssh_keys_deleted_at ON ssh_keys(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_host_groups_name ON host_groups(name)",
		"CREATE INDEX IF NOT EXISTS idx_host_groups_deleted_at ON host_groups(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_hostname ON hosts(hostname)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_ip_address ON hosts(ip_address)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_environment ON hosts(environment)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_cloud_platform_id ON hosts(cloud_platform_id)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_project_id ON hosts(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_status ON hosts(status)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_salt_minion_id ON hosts(salt_minion_id)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_ssh_key_id ON hosts(ssh_key_id)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_deleted_at ON hosts(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_hosts_ip_address ON hosts(ip_address)",
		"CREATE INDEX IF NOT EXISTS idx_deployment_configs_project_id ON deployment_configs(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_deployment_configs_created_by ON deployment_configs(created_by)",
		"CREATE INDEX IF NOT EXISTS idx_deployment_configs_deleted_at ON deployment_configs(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_deployments_project_id ON deployments(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_deployments_config_id ON deployments(config_id)",
		"CREATE INDEX IF NOT EXISTS idx_deployments_deployed_by ON deployments(deployed_by)",
		"CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status)",
		"CREATE INDEX IF NOT EXISTS idx_deployments_deleted_at ON deployments(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_deployment_versions_deployment_id ON deployment_versions(deployment_id)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_type ON audit_logs(resource_type)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs(resource_id)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_status ON audit_logs(status)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_deleted_at ON audit_logs(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_system_settings_key ON system_settings(key)",
		"CREATE INDEX IF NOT EXISTS idx_system_settings_category ON system_settings(category)",
		"CREATE INDEX IF NOT EXISTS idx_system_settings_updated_by ON system_settings(updated_by)",
		"CREATE INDEX IF NOT EXISTS idx_batch_operations_started_by ON batch_operations(started_by)",
		"CREATE INDEX IF NOT EXISTS idx_batch_operations_status ON batch_operations(status)",
		"CREATE INDEX IF NOT EXISTS idx_batch_operations_started_at ON batch_operations(started_at)",
		"CREATE INDEX IF NOT EXISTS idx_batch_operations_created_at ON batch_operations(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_batch_operations_deleted_at ON batch_operations(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_command_templates_created_by ON command_templates(created_by)",
		"CREATE INDEX IF NOT EXISTS idx_command_templates_category ON command_templates(category)",
		"CREATE INDEX IF NOT EXISTS idx_command_templates_is_public ON command_templates(is_public)",
		"CREATE INDEX IF NOT EXISTS idx_command_templates_deleted_at ON command_templates(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_formula_repositories_name ON formula_repositories(name)",
		"CREATE INDEX IF NOT EXISTS idx_formula_repositories_is_active ON formula_repositories(is_active)",
		"CREATE INDEX IF NOT EXISTS idx_formula_repositories_deleted_at ON formula_repositories(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_formulas_name ON formulas(name)",
		"CREATE INDEX IF NOT EXISTS idx_formulas_category ON formulas(category)",
		"CREATE INDEX IF NOT EXISTS idx_formulas_repository ON formulas(repository)",
		"CREATE INDEX IF NOT EXISTS idx_formulas_is_active ON formulas(is_active)",
		"CREATE INDEX IF NOT EXISTS idx_formulas_deleted_at ON formulas(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_formula_parameters_formula_id ON formula_parameters(formula_id)",
		"CREATE INDEX IF NOT EXISTS idx_formula_templates_formula_id ON formula_templates(formula_id)",
		"CREATE INDEX IF NOT EXISTS idx_formula_templates_created_by ON formula_templates(created_by)",
		"CREATE INDEX IF NOT EXISTS idx_formula_templates_is_public ON formula_templates(is_public)",
		"CREATE INDEX IF NOT EXISTS idx_formula_templates_deleted_at ON formula_templates(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_formula_deployments_formula_id ON formula_deployments(formula_id)",
		"CREATE INDEX IF NOT EXISTS idx_formula_deployments_status ON formula_deployments(status)",
		"CREATE INDEX IF NOT EXISTS idx_formula_deployments_started_by ON formula_deployments(started_by)",
		"CREATE INDEX IF NOT EXISTS idx_formula_deployments_started_at ON formula_deployments(started_at)",
		"CREATE INDEX IF NOT EXISTS idx_formula_deployments_deleted_at ON formula_deployments(deleted_at)",
	}

	log.Println("Creating indexes...")
	for _, indexSQL := range indexes {
		if err := DB.Exec(indexSQL).Error; err != nil {
			log.Printf("Warning: Failed to create index: %s, error: %v", indexSQL, err)
			// 不返回错误，因为索引创建失败不应该阻止整个迁移
		}
	}

	log.Println("All tables and indexes created successfully")
	return nil
}

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	dsn := config.AppConfig.Database.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info), // 使用 Info 模式，但减少某些日志可能避免问题
		DisableForeignKeyConstraintWhenMigrating: false,                               // 启用迁移时的外键约束检查
	})
	if err != nil {
		return err
	}

	// 获取底层 SQL 数据库连接，设置连接参数
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return gorm.ErrInvalidDB
	}

	log.Println("Starting database migration...")

	// 使用纯SQL手动迁移所有表，避免 GORM AutoMigrate 的 pgx/v5 兼容性问题
	log.Println("Using manual SQL migration to avoid pgx/v5 compatibility issues")
	if err := manualMigrateAllTables(); err != nil {
		return fmt.Errorf("failed to manually migrate all tables: %w", err)
	}

	log.Println("Database migration completed successfully")

	// 初始化默认数据
	if err := initDefaultData(); err != nil {
		log.Printf("Warning: Failed to initialize default data: %v", err)
		// 不返回错误，因为迁移已经成功
	}

	return nil
}

// initDefaultData 初始化默认数据（权限、角色等）
func initDefaultData() error {
	log.Println("Initializing default data...")

	// 检查是否已有权限数据，避免重复插入
	var count int64
	DB.Model(&Permission{}).Count(&count)
	if count > 0 {
		log.Println("Default permissions already exist, skipping...")
	} else {
		// 插入默认权限
		defaultPermissions := []Permission{
			{Code: "host:read", Name: "查看主机", ResourceType: "host", Action: "read", Description: "查看主机信息"},
			{Code: "host:create", Name: "创建主机", ResourceType: "host", Action: "create", Description: "创建新主机"},
			{Code: "host:update", Name: "更新主机", ResourceType: "host", Action: "update", Description: "更新主机信息"},
			{Code: "host:delete", Name: "删除主机", ResourceType: "host", Action: "delete", Description: "删除主机"},
			{Code: "deployment:read", Name: "查看部署", ResourceType: "deployment", Action: "read", Description: "查看部署信息"},
			{Code: "deployment:create", Name: "创建部署", ResourceType: "deployment", Action: "create", Description: "创建部署配置"},
			{Code: "deployment:execute", Name: "执行部署", ResourceType: "deployment", Action: "execute", Description: "执行部署操作"},
			{Code: "task:read", Name: "查看任务", ResourceType: "task", Action: "read", Description: "查看定时任务"},
			{Code: "task:create", Name: "创建任务", ResourceType: "task", Action: "create", Description: "创建定时任务"},
			{Code: "task:execute", Name: "执行任务", ResourceType: "task", Action: "execute", Description: "执行定时任务"},
			{Code: "log:read", Name: "查看日志", ResourceType: "log", Action: "read", Description: "查看系统日志"},
			{Code: "monitoring:read", Name: "查看监控", ResourceType: "monitoring", Action: "read", Description: "查看监控指标"},
			{Code: "audit:read", Name: "查看审计", ResourceType: "audit", Action: "read", Description: "查看审计日志"},
			{Code: "user:manage", Name: "用户管理", ResourceType: "user", Action: "manage", Description: "管理用户和权限"},
			{Code: "project:read", Name: "查看项目", ResourceType: "project", Action: "read", Description: "查看项目信息"},
			{Code: "project:create", Name: "创建项目", ResourceType: "project", Action: "create", Description: "创建新项目"},
			{Code: "project:update", Name: "更新项目", ResourceType: "project", Action: "update", Description: "更新项目信息"},
			{Code: "project:delete", Name: "删除项目", ResourceType: "project", Action: "delete", Description: "删除项目"},
			{Code: "webssh:execute", Name: "执行WebSSH", ResourceType: "webssh", Action: "execute", Description: "执行WebSSH终端连接"},
			{Code: "batch_operation:read", Name: "查看批量操作", ResourceType: "batch_operation", Action: "read", Description: "查看批量操作信息"},
			{Code: "batch_operation:create", Name: "创建批量操作", ResourceType: "batch_operation", Action: "create", Description: "创建并执行批量操作"},
			{Code: "batch_operation:execute", Name: "执行批量操作", ResourceType: "batch_operation", Action: "execute", Description: "执行批量操作"},
			{Code: "batch_operation:cancel", Name: "取消批量操作", ResourceType: "batch_operation", Action: "cancel", Description: "取消正在执行的批量操作"},
		}

		for _, perm := range defaultPermissions {
			if err := DB.Create(&perm).Error; err != nil {
				log.Printf("Failed to create permission %s: %v", perm.Code, err)
				// 继续创建其他权限，不中断
			}
		}
		log.Printf("Created %d default permissions", len(defaultPermissions))
	}

	// 检查是否已有角色数据
	DB.Model(&Role{}).Count(&count)
	if count > 0 {
		log.Println("Default roles already exist, skipping...")
	} else {
		// 插入默认角色
		defaultRoles := []Role{
			{Name: "admin", Description: "系统管理员，拥有所有权限"},
			{Name: "operator", Description: "运维人员，拥有运维相关权限"},
			{Name: "viewer", Description: "查看者，只能查看信息"},
		}

		for _, role := range defaultRoles {
			if err := DB.Create(&role).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
				continue
			}

			// 为角色分配权限
			if role.Name == "admin" {
				// 管理员拥有所有权限
				var allPerms []Permission
				DB.Find(&allPerms)
				for _, perm := range allPerms {
					DB.Create(&RolePermission{
						RoleID:       role.ID,
						PermissionID: perm.ID,
					})
				}
			} else if role.Name == "operator" {
				// 运维人员拥有运维相关权限
				var opsPerms []Permission
				DB.Where("resource_type IN (?)", []string{"host", "deployment", "task", "log", "monitoring", "webssh", "batch_operation"}).Find(&opsPerms)
				for _, perm := range opsPerms {
					DB.Create(&RolePermission{
						RoleID:       role.ID,
						PermissionID: perm.ID,
					})
				}
			} else if role.Name == "viewer" {
				// 查看者只有只读权限
				var readPerms []Permission
				DB.Where("action = ?", "read").Find(&readPerms)
				for _, perm := range readPerms {
					DB.Create(&RolePermission{
						RoleID:       role.ID,
						PermissionID: perm.ID,
					})
				}
			}
		}
		log.Printf("Created %d default roles", len(defaultRoles))
	}

	// 创建默认管理员账号（如果不存在）
	var adminUser User
	err := DB.Where("username = ?", "admin").First(&adminUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 管理员账号不存在，创建默认管理员
		hashedPassword, err := hashPassword("admin123")
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
		} else {
			adminUser = User{
				Username:     "admin",
				Email:        "admin@kkops.local",
				PasswordHash: hashedPassword,
				DisplayName:  "系统管理员",
				Status:       "active",
			}
			if err := DB.Create(&adminUser).Error; err != nil {
				log.Printf("Failed to create admin user: %v", err)
			} else {
				// 查找 admin 角色并分配给管理员
				var adminRole Role
				if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err == nil {
					DB.Create(&UserRole{
						UserID: adminUser.ID,
						RoleID: adminRole.ID,
					})
					log.Println("Default admin user created: username=admin, password=admin123")
					log.Println("⚠️  WARNING: Please change the default admin password after first login!")
				} else if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.Printf("Failed to find admin role: %v", err)
				}
			}
		}
	} else if err != nil {
		// 其他数据库错误
		log.Printf("Error checking admin user: %v", err)
	} else {
		// 管理员账号已存在
		log.Println("Admin user already exists, skipping creation")
	}

	// 初始化系统设置（Salt配置）
	if err := initSystemSettings(); err != nil {
		log.Printf("Warning: Failed to initialize system settings: %v", err)
		// 不返回错误，因为迁移已经成功
	}

	// 初始化默认环境
	if err := initDefaultEnvironments(); err != nil {
		log.Printf("Warning: Failed to initialize default environments: %v", err)
	}

	// 初始化默认命令模板
	if err := initDefaultCommandTemplates(); err != nil {
		log.Printf("Warning: Failed to initialize default command templates: %v", err)
	}

	// 修复现有Formula相关记录的created_by字段
	if err := fixFormulaCreatedByFields(); err != nil {
		log.Printf("Warning: Failed to fix formula created_by fields: %v", err)
	}

	log.Println("Default data initialization completed")
	return nil
}

// initDefaultCommandTemplates 初始化默认命令模板
func initDefaultCommandTemplates() error {
	log.Println("Initializing default command templates...")

	var count int64
	DB.Model(&CommandTemplate{}).Count(&count)
	if count > 0 {
		log.Println("Command templates already exist, skipping initialization...")
		return nil
	}

	// 获取管理员用户ID（假设ID为1，或者从users表查询第一个管理员）
	var adminUser User
	var adminUserID uint64 = 1
	if err := DB.Where("username = ?", "admin").First(&adminUser).Error; err == nil {
		adminUserID = adminUser.ID
	} else {
		log.Printf("Warning: Admin user not found, using user ID 1 for templates")
	}

	defaultTemplates := []CommandTemplate{
		// 系统信息
		{Name: "查看系统负载", Description: "查看系统负载平均值", Category: "system", CommandFunction: "status.loadavg", CreatedBy: adminUserID, IsPublic: true},
		{Name: "查看内存使用", Description: "查看内存使用情况", Category: "system", CommandFunction: "status.meminfo", CreatedBy: adminUserID, IsPublic: true},
		{Name: "查看磁盘使用", Description: "查看磁盘使用情况", Category: "disk", CommandFunction: "disk.usage", CreatedBy: adminUserID, IsPublic: true},
		{Name: "查看CPU信息", Description: "查看CPU信息", Category: "system", CommandFunction: "status.cpuinfo", CreatedBy: adminUserID, IsPublic: true},
		// 网络信息
		{Name: "查看网络接口", Description: "查看网络接口信息", Category: "network", CommandFunction: "network.interfaces", CreatedBy: adminUserID, IsPublic: true},
		{Name: "查看活动TCP连接", Description: "查看活动的TCP连接", Category: "network", CommandFunction: "network.active_tcp", CreatedBy: adminUserID, IsPublic: true},
		// 进程管理
		{Name: "查看进程列表", Description: "查看进程列表", Category: "process", CommandFunction: "cmd.run", CommandArgs: datatypes.JSON([]byte(`["ps aux | head -20"]`)), CreatedBy: adminUserID, IsPublic: true},
		{Name: "查看服务状态", Description: "查看服务状态", Category: "process", CommandFunction: "service.status", CommandArgs: datatypes.JSON([]byte(`["*"]`)), CreatedBy: adminUserID, IsPublic: true},
		// 包管理
		{Name: "批量安装软件包", Description: "批量安装多个软件包", Category: "package", CommandFunction: "pkg.install", CommandArgs: datatypes.JSON([]byte(`["vim","htop","curl"]`)), CreatedBy: adminUserID, IsPublic: true},
		{Name: "更新系统包", Description: "更新系统软件包列表", Category: "package", CommandFunction: "pkg.update", CreatedBy: adminUserID, IsPublic: true},
		{Name: "升级系统包", Description: "升级所有已安装的软件包", Category: "package", CommandFunction: "pkg.upgrade", CreatedBy: adminUserID, IsPublic: true},
		// 测试
		{Name: "测试连通性", Description: "测试主机连通性", Category: "system", CommandFunction: "test.ping", CreatedBy: adminUserID, IsPublic: true},
	}

	for _, template := range defaultTemplates {
		if err := DB.Create(&template).Error; err != nil {
			log.Printf("Failed to create template %s: %v", template.Name, err)
		}
	}

	log.Printf("Created %d default command templates", len(defaultTemplates))
	return nil
}

// initDefaultEnvironments 初始化默认环境
func initDefaultEnvironments() error {
	log.Println("Initializing default environments...")

	var count int64
	DB.Model(&Environment{}).Count(&count)
	if count > 0 {
		log.Println("Environments already exist, skipping initialization...")
		return nil
	}

	defaultEnvs := []Environment{
		{Name: "dev", DisplayName: "开发环境", Color: "blue", SortOrder: 1, Description: "开发测试环境"},
		{Name: "test", DisplayName: "测试环境", Color: "green", SortOrder: 2, Description: "功能测试环境"},
		{Name: "uat", DisplayName: "UAT环境", Color: "orange", SortOrder: 3, Description: "用户验收测试环境"},
		{Name: "staging", DisplayName: "预发布环境", Color: "purple", SortOrder: 4, Description: "预发布/灰度环境"},
		{Name: "prod", DisplayName: "生产环境", Color: "red", SortOrder: 5, Description: "生产环境"},
	}

	for _, env := range defaultEnvs {
		if err := DB.Create(&env).Error; err != nil {
			log.Printf("Failed to create environment %s: %v", env.Name, err)
		}
	}

	log.Printf("Created %d default environments", len(defaultEnvs))
	return nil
}

// initSystemSettings 初始化系统设置
func initSystemSettings() error {
	log.Println("Initializing system settings...")

	// 检查是否已有系统设置
	var count int64
	DB.Model(&SystemSettings{}).Count(&count)
	if count > 0 {
		log.Println("System settings already exist, skipping initialization...")
		return nil
	}

	// 从环境变量中读取Salt配置并保存到数据库
	saltConfig := config.AppConfig.Salt

	settings := []SystemSettings{
		{
			Key:         "salt.api_url",
			Value:       saltConfig.APIURL,
			Category:    "salt",
			Description: "Salt API URL",
			UpdatedBy:   1, // 系统用户
		},
		{
			Key:         "salt.username",
			Value:       saltConfig.Username,
			Category:    "salt",
			Description: "Salt API username",
			UpdatedBy:   1,
		},
		{
			Key:         "salt.password",
			Value:       saltConfig.Password, // 注意：这里存储的是明文，实际应该加密
			Category:    "salt",
			Description: "Salt API password (encrypted)",
			UpdatedBy:   1,
		},
		{
			Key:         "salt.eauth",
			Value:       saltConfig.EAuth,
			Category:    "salt",
			Description: "Salt API eauth type",
			UpdatedBy:   1,
		},
		{
			Key:         "salt.timeout",
			Value:       fmt.Sprintf("%d", saltConfig.Timeout),
			Category:    "salt",
			Description: "Salt API timeout in seconds",
			UpdatedBy:   1,
		},
		{
			Key:         "salt.verify_ssl",
			Value:       fmt.Sprintf("%v", saltConfig.VerifySSL),
			Category:    "salt",
			Description: "Verify SSL certificates",
			UpdatedBy:   1,
		},
	}

	// 加密密码字段
	for i := range settings {
		if settings[i].Key == "salt.password" && settings[i].Value != "" {
			// 使用utils.Encrypt加密密码
			encrypted, err := utils.Encrypt(settings[i].Value)
			if err != nil {
				log.Printf("Warning: Failed to encrypt password: %v", err)
			} else {
				settings[i].Value = encrypted
			}
		}
	}

	// 批量创建设置
	for _, setting := range settings {
		if err := DB.Create(&setting).Error; err != nil {
			log.Printf("Failed to create setting %s: %v", setting.Key, err)
			// 继续创建其他设置，不中断
		}
	}

	log.Printf("Initialized %d system settings", len(settings))
	return nil
}

// hashPassword 加密密码（临时函数，避免循环依赖）
func hashPassword(password string) (string, error) {
	// 这里需要导入 utils 包，但为了避免循环依赖，我们直接使用 bcrypt
	// 或者可以把这个函数移到 utils 包
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// fixFormulaCreatedByFields 修复现有Formula相关记录中created_by字段为0的问题
func fixFormulaCreatedByFields() error {
	log.Println("Fixing formula created_by fields...")

	// 删除name为空字符串的Formula记录（脏数据）
	if err := DB.Where("name = '' OR name IS NULL").Delete(&Formula{}).Error; err != nil {
		log.Printf("Failed to delete empty name formulas: %v", err)
		return err
	}

	// 再次删除，以防有多个空name记录
	if err := DB.Where("name = ''").Delete(&Formula{}).Error; err != nil {
		log.Printf("Failed to delete remaining empty name formulas: %v", err)
		return err
	}

	// 修复formulas表
	if err := DB.Model(&Formula{}).Where("created_by = 0 OR created_by IS NULL").Update("created_by", 1).Error; err != nil {
		log.Printf("Failed to fix formulas created_by: %v", err)
		return err
	}

	// 修复formula_repositories表
	if err := DB.Model(&FormulaRepository{}).Where("created_by = 0 OR created_by IS NULL").Update("created_by", 1).Error; err != nil {
		log.Printf("Failed to fix formula_repositories created_by: %v", err)
		return err
	}

	// 修复formula_templates表
	if err := DB.Model(&FormulaTemplate{}).Where("created_by = 0 OR created_by IS NULL").Update("created_by", 1).Error; err != nil {
		log.Printf("Failed to fix formula_templates created_by: %v", err)
		return err
	}

	// 修复formula_deployments表
	if err := DB.Model(&FormulaDeployment{}).Where("started_by = 0 OR started_by IS NULL").Update("started_by", 1).Error; err != nil {
		log.Printf("Failed to fix formula_deployments started_by: %v", err)
		return err
	}

	log.Println("Formula created_by fields fixed successfully")
	return nil
}
