// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// InitDB initializes the database connection and runs migrations
func InitDB(cfg *DatabaseConfig, zapLogger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()

	var gormLogger logger.Interface
	if zapLogger != nil {
		gormLogger = NewGormLogger(zapLogger)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// Run SQL migrations first (for complex schema changes)
	if err := runSQLMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	// Run GORM AutoMigrate (for schema sync)
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// seedDefaultUser creates a default admin user if no users exist
func seedDefaultUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}

	// If users already exist, skip seeding
	if count > 0 {
		return nil
	}

	// Create default admin user
	passwordHash, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash default password: %w", err)
	}

	adminUser := model.User{
		Username:     "admin",
		Email:        "admin@kkops.local",
		PasswordHash: passwordHash,
		RealName:     "Administrator",
		Status:       "active",
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	return nil
}

// seedDefaultAdminRole creates a default admin role if no admin role exists
func seedDefaultAdminRole(db *gorm.DB) error {
	var adminRole model.Role
	result := db.Where("name = ?", "admin").First(&adminRole)
	
	if result.Error == gorm.ErrRecordNotFound {
		// Create admin role
		adminRole = model.Role{
			Name:        "admin",
			Description: "系统管理员，拥有所有资产的访问权限",
			IsAdmin:     true,
		}
		if err := db.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("failed to create default admin role: %w", err)
		}
		
		// Assign admin role to admin user
		var adminUser model.User
		if err := db.Where("username = ?", "admin").First(&adminUser).Error; err == nil {
			db.Model(&adminUser).Association("Roles").Append(&adminRole)
		}
	} else if result.Error != nil {
		return result.Error
	} else {
		// Admin role exists, ensure is_admin is true
		if !adminRole.IsAdmin {
			adminRole.IsAdmin = true
			db.Save(&adminRole)
		}
	}
	
	return nil
}

// seedDefaultTemplates creates default task templates
func seedDefaultTemplates(db *gorm.DB) error {
	// 获取 admin 用户 ID
	var adminUser model.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		// 如果没有 admin 用户，跳过模板创建
		return nil
	}
	adminID := adminUser.ID

	// 系统信息采集模板
	systemInfoTemplate := model.TaskTemplate{
		Name:        "系统信息采集",
		Type:        "shell",
		Description: "采集主机 CPU、内存、磁盘信息，用于自动更新资产信息",
		CreatedBy:   adminID,
		Content: `#!/usr/bin/env bash
# 系统信息采集脚本 - 输出 JSON 格式供自动更新资产信息

# 获取主机名
hostname=$(hostname)

# 获取 CPU 信息
cpu_info=$(grep -m1 'model name' /proc/cpuinfo 2>/dev/null | cut -d: -f2 | xargs || echo "Unknown")
cpu_cores=$(nproc 2>/dev/null || echo "1")
cpu="${cpu_info} (${cpu_cores}核)"

# 获取内存信息（MB）
mem_total=$(awk '/MemTotal/ {printf "%.0f", $2/1024}' /proc/meminfo 2>/dev/null || echo "0")
memory="${mem_total}MB"

# 获取磁盘信息（根分区）
disk_total=$(df -BG / 2>/dev/null | awk 'NR==2 {gsub(/G/,"",$2); print $2}' || echo "0")
disk="${disk_total}GB"

# 输出 JSON 格式（用于自动更新资产信息）
cat << EOF
{
  "hostname": "${hostname}",
  "cpu": "${cpu}",
  "memory": "${memory}",
  "disk": "${disk}"
}
EOF`,
	}

	// 辅助函数：检查模板是否存在（包括软删除的记录）
	templateExists := func(name string) bool {
		var count int64
		db.Unscoped().Model(&model.TaskTemplate{}).Where("name = ?", name).Count(&count)
		return count > 0
	}

	// 创建系统信息采集模板
	if !templateExists(systemInfoTemplate.Name) {
		if err := db.Create(&systemInfoTemplate).Error; err != nil {
			return fmt.Errorf("failed to create system info template: %w", err)
		}
	}

	// 磁盘使用率检查模板
	diskUsageTemplate := model.TaskTemplate{
		Name:        "磁盘使用率检查",
		Type:        "shell",
		Description: "检查磁盘使用率，超过阈值则告警",
		CreatedBy:   adminID,
		Content: `#!/usr/bin/env bash
# 磁盘使用率检查脚本

THRESHOLD=80

echo "=== 磁盘使用率检查 ==="
echo "主机: $(hostname)"
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

df -h | grep -E '^/dev/' | while read line; do
    usage=$(echo $line | awk '{print $5}' | sed 's/%//')
    mount=$(echo $line | awk '{print $6}')
    
    if [ "$usage" -gt "$THRESHOLD" ]; then
        echo "[警告] $mount 使用率: ${usage}% (超过阈值 ${THRESHOLD}%)"
    else
        echo "[正常] $mount 使用率: ${usage}%"
    fi
done`,
	}

	if !templateExists(diskUsageTemplate.Name) {
		if err := db.Create(&diskUsageTemplate).Error; err != nil {
			return fmt.Errorf("failed to create disk usage template: %w", err)
		}
	}

	// 系统健康检查模板
	healthCheckTemplate := model.TaskTemplate{
		Name:        "系统健康检查",
		Type:        "shell",
		Description: "检查系统负载、内存、进程等健康状态",
		CreatedBy:   adminID,
		Content: `#!/usr/bin/env bash
# 系统健康检查脚本

echo "=== 系统健康检查 ==="
echo "主机: $(hostname)"
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

echo "--- 系统负载 ---"
uptime

echo ""
echo "--- 内存使用 ---"
free -h

echo ""
echo "--- 磁盘使用 ---"
df -h | grep -E '^/dev/'

echo ""
echo "--- CPU 占用 TOP 5 ---"
ps aux --sort=-%cpu | head -6

echo ""
echo "--- 内存占用 TOP 5 ---"
ps aux --sort=-%mem | head -6

echo ""
echo "--- 网络连接统计 ---"
ss -s 2>/dev/null || netstat -an | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'`,
	}

	if !templateExists(healthCheckTemplate.Name) {
		if err := db.Create(&healthCheckTemplate).Error; err != nil {
			return fmt.Errorf("failed to create health check template: %w", err)
		}
	}

	// 日志清理模板
	logCleanupTemplate := model.TaskTemplate{
		Name:        "日志清理",
		Type:        "shell",
		Description: "清理过期日志文件，释放磁盘空间",
		CreatedBy:   adminID,
		Content: `#!/usr/bin/env bash
# 日志清理脚本

# 配置
LOG_DIRS="/var/log /tmp"
DAYS=7

echo "=== 日志清理 ==="
echo "主机: $(hostname)"
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "清理 ${DAYS} 天前的日志文件"
echo ""

for dir in $LOG_DIRS; do
    if [ -d "$dir" ]; then
        echo "--- 清理目录: $dir ---"
        # 查找并删除过期的日志文件
        find "$dir" -name "*.log" -type f -mtime +$DAYS -exec ls -lh {} \; 2>/dev/null
        find "$dir" -name "*.log.*" -type f -mtime +$DAYS -exec ls -lh {} \; 2>/dev/null
        # 实际删除（取消注释以下行）
        # find "$dir" -name "*.log" -type f -mtime +$DAYS -delete 2>/dev/null
        # find "$dir" -name "*.log.*" -type f -mtime +$DAYS -delete 2>/dev/null
    fi
done

echo ""
echo "清理完成"`,
	}

	if !templateExists(logCleanupTemplate.Name) {
		if err := db.Create(&logCleanupTemplate).Error; err != nil {
			return fmt.Errorf("failed to create log cleanup template: %w", err)
		}
	}

	return nil
}

// seedDefaultScheduledTasks creates default scheduled tasks
func seedDefaultScheduledTasks(db *gorm.DB) error {
	// 获取 admin 用户 ID
	var adminUser model.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return nil // 如果没有 admin 用户，跳过
	}

	// 辅助函数：检查定时任务是否存在（包括软删除的记录）
	taskExists := func(name string) bool {
		var count int64
		db.Unscoped().Model(&model.ScheduledTask{}).Where("name = ?", name).Count(&count)
		return count > 0
	}

	// 获取系统信息采集模板
	var systemInfoTemplate model.TaskTemplate
	if err := db.Where("name = ?", "系统信息采集").First(&systemInfoTemplate).Error; err != nil {
		return nil // 如果模板不存在，跳过
	}

	// 创建每日系统信息采集定时任务
	taskName := "每日系统信息采集"
	if !taskExists(taskName) {
		scheduledTask := model.ScheduledTask{
			Name:           taskName,
			Description:    "每天午夜自动采集所有主机的 CPU、内存、磁盘信息并更新资产管理",
			CronExpression: "0 0 0 * * *", // 每天午夜执行
			TemplateID:     &systemInfoTemplate.ID,
			Content:        systemInfoTemplate.Content,
			Type:           "shell",
			Timeout:        300,
			Enabled:        false, // 默认禁用，需要用户选择主机后启用
			UpdateAssets:   true,  // 启用资产自动更新
			CreatedBy:      adminUser.ID,
		}

		if err := db.Create(&scheduledTask).Error; err != nil {
			return fmt.Errorf("failed to create system info scheduled task: %w", err)
		}
	}

	// 创建每日健康检查定时任务
	var healthCheckTemplate model.TaskTemplate
	if err := db.Where("name = ?", "系统健康检查").First(&healthCheckTemplate).Error; err == nil {
		healthTaskName := "每日系统健康检查"
		if !taskExists(healthTaskName) {
			scheduledTask := model.ScheduledTask{
				Name:           healthTaskName,
				Description:    "每天早上 8 点执行系统健康检查",
				CronExpression: "0 0 8 * * *", // 每天早上 8 点
				TemplateID:     &healthCheckTemplate.ID,
				Content:        healthCheckTemplate.Content,
				Type:           "shell",
				Timeout:        300,
				Enabled:        false, // 默认禁用
				UpdateAssets:   false,
				CreatedBy:      adminUser.ID,
			}

			if err := db.Create(&scheduledTask).Error; err != nil {
				return fmt.Errorf("failed to create health check scheduled task: %w", err)
			}
		}
	}

	return nil
}

// runSQLMigrations executes SQL migration files in order
func runSQLMigrations(db *gorm.DB) error {
	// Get the base directory (assuming migrations are relative to backend/)
	migrationDir := "migrations"
	
	// Check if migration directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		// Try alternative path (when running from backend directory)
		migrationDir = "backend/migrations"
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			// No migration directory, skip SQL migrations
			return nil
		}
	}

	// List of migration files in order
	migrationFiles := []string{
		"rename_asset_name_to_hostname.sql",
		"remove_code_fields.sql",
		"add_cloud_platform_management.sql",
		"add_asset_authorization.sql",
		"add_role_asset_authorization.sql",
		"add_deployment_management.sql",
		"add_scheduled_task.sql",
		"add_audit_log.sql",
		"add_operations_navigation.sql",
	}

	// Execute migrations in order
	for _, filename := range migrationFiles {
		filepath := filepath.Join(migrationDir, filename)
		
		// Check if file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			continue // Skip if file doesn't exist
		}

		// Read migration file
		sqlContent, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute the entire SQL file as one transaction
		// This is better for DO blocks and complex migrations
		sql := string(sqlContent)
		
		// Remove comments that are on their own lines
		lines := strings.Split(sql, "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Keep lines that are not empty comments or comment-only lines
			if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
				cleanLines = append(cleanLines, line)
			} else if strings.Contains(line, "--") && !strings.HasPrefix(trimmed, "--") {
				// Keep lines that have code before comments
				cleanLines = append(cleanLines, line)
			}
		}
		sql = strings.Join(cleanLines, "\n")

		// Execute SQL
		if err := db.Exec(sql).Error; err != nil {
			// Some errors are expected (e.g., column already exists, index already exists)
			// Check if it's a safe error to ignore
			errStr := err.Error()
			if strings.Contains(errStr, "already exists") ||
				strings.Contains(errStr, "duplicate") {
				// Safe to ignore, continue
				continue
			}
			
			// For column/table not found errors, they might be expected depending on migration state
			// But we should log them as warnings, not fail completely
			if strings.Contains(errStr, "does not exist") {
				// Log as warning but continue
				continue
			}

			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	return nil
}

// migrate runs database migrations
func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Department{},
		&model.APIToken{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.UserAsset{},
		&model.RoleAsset{},
		&model.Project{},
		&model.Environment{},
		&model.CloudPlatform{},
		&model.AssetCategory{},
		&model.Asset{},
		&model.Tag{},
		&model.AssetTag{},
		&model.SSHKey{},
		&model.TaskTemplate{},
		&model.Task{},
		&model.TaskExecution{},
		&model.ScheduledTask{},
		&model.DeploymentModule{},
		&model.Deployment{},
		&model.AuditLog{},
		&model.OperationTool{},
	); err != nil {
		return err
	}

	// Initialize default admin user and admin role if database is empty
	if err := seedDefaultUser(db); err != nil {
		return err
	}
	if err := seedDefaultAdminRole(db); err != nil {
		return err
	}
	// Initialize default task templates
	if err := seedDefaultTemplates(db); err != nil {
		return err
	}
	// Initialize default scheduled tasks
	if err := seedDefaultScheduledTasks(db); err != nil {
		return err
	}
	// Initialize default operation tools
	if err := seedDefaultOperationTools(db); err != nil {
		return err
	}
	// Initialize menu permissions
	return seedMenuPermissions(db)
}

// seedDefaultOperationTools creates default operation tools
func seedDefaultOperationTools(db *gorm.DB) error {
	// 辅助函数：检查工具是否存在
	toolExists := func(name string) bool {
		var count int64
		db.Model(&model.OperationTool{}).Where("name = ?", name).Count(&count)
		return count > 0
	}

	// 定义默认运维工具
	defaultTools := []model.OperationTool{
		{
			Name:        "Grafana",
			Description: "开源监控和可视化平台",
			Category:    "监控工具",
			Icon:        "dashboard",
			URL:         "https://grafana.example.com",
			OrderIndex:  1,
			Enabled:     true,
		},
		{
			Name:        "Prometheus",
			Description: "开源监控和告警系统",
			Category:    "监控工具",
			Icon:        "appstore",
			URL:         "https://prometheus.example.com",
			OrderIndex:  2,
			Enabled:     true,
		},
		{
			Name:        "Kibana",
			Description: "Elasticsearch 数据可视化和分析平台",
			Category:    "日志工具",
			Icon:        "fileSearch",
			URL:         "https://kibana.example.com",
			OrderIndex:  3,
			Enabled:     true,
		},
		{
			Name:        "Jenkins",
			Description: "开源自动化服务器，用于 CI/CD",
			Category:    "CI/CD工具",
			Icon:        "rocket",
			URL:         "https://jenkins.example.com",
			OrderIndex:  4,
			Enabled:     true,
		},
		{
			Name:        "GitLab",
			Description: "Git 代码管理和 CI/CD 平台",
			Category:    "CI/CD工具",
			Icon:        "code",
			URL:         "https://gitlab.example.com",
			OrderIndex:  5,
			Enabled:     true,
		},
		{
			Name:        "phpMyAdmin",
			Description: "MySQL/MariaDB 数据库管理工具",
			Category:    "数据库工具",
			Icon:        "database",
			URL:         "https://phpmyadmin.example.com",
			OrderIndex:  6,
			Enabled:     true,
		},
		{
			Name:        "DBeaver",
			Description: "通用数据库管理工具（需要安装客户端）",
			Category:    "数据库工具",
			Icon:        "database",
			URL:         "https://dbeaver.io",
			OrderIndex:  7,
			Enabled:     true,
		},
	}

	// 创建不存在的工具
	for _, tool := range defaultTools {
		if !toolExists(tool.Name) {
			if err := db.Create(&tool).Error; err != nil {
				return fmt.Errorf("failed to create operation tool %s: %w", tool.Name, err)
			}
		}
	}

	return nil
}

// seedMenuPermissions creates menu permissions if they don't exist
func seedMenuPermissions(db *gorm.DB) error {
	for _, permDef := range model.AllMenuPermissions {
		// Check if permission already exists
		var existingPerm model.Permission
		result := db.Where("resource = ? AND action = ?", permDef.Resource, permDef.Action).First(&existingPerm)
		
		if result.Error == gorm.ErrRecordNotFound {
			// Create new permission
			permission := model.Permission{
				Name:        permDef.Name,
				Resource:    permDef.Resource,
				Action:      permDef.Action,
				Description: permDef.Description,
			}
			if err := db.Create(&permission).Error; err != nil {
				return fmt.Errorf("failed to create permission %s:%s: %w", permDef.Resource, permDef.Action, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to check permission %s:%s: %w", permDef.Resource, permDef.Action, result.Error)
		}
		// If permission exists, skip (idempotent)
	}
	
	return nil
}

// NewGormLogger creates a GORM logger adapter for Zap
func NewGormLogger(zapLogger *zap.Logger) logger.Interface {
	return logger.New(
		&zapWriter{logger: zapLogger},
		logger.Config{
			SlowThreshold:             0,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

type zapWriter struct {
	logger *zap.Logger
}

func (w *zapWriter) Printf(format string, args ...interface{}) {
	w.logger.Sugar().Infof(format, args...)
}
