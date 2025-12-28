package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 注意：所有建表和建索引的 SQL 语句已迁移到 backend/migrations/ 目录
// 请使用迁移文件来管理数据库结构变更，而不是在代码中硬编码 SQL

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

	// 从 migrations 目录读取并执行 SQL 文件
	log.Println("Executing migrations from migrations directory...")
	if err := executeMigrations(); err != nil {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	log.Println("Database migration completed successfully")

	// 初始化默认数据
	if err := initDefaultData(); err != nil {
		log.Printf("Warning: Failed to initialize default data: %v", err)
		// 不返回错误，因为迁移已经成功
	}

	return nil
}

// executeMigrations 从 migrations 目录读取并执行所有 .up.sql 文件
func executeMigrations() error {
	// 获取 migrations 目录路径（相对于项目根目录）
	migrationsDir := "backend/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// 如果 backend/migrations 不存在，尝试 migrations（可能在 backend 目录内运行）
		migrationsDir = "migrations"
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found. Please ensure backend/migrations/ directory exists with migration files")
		}
	}

	// 读取所有 .up.sql 文件
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*_*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// 按文件名排序（时间戳顺序）
	sort.Strings(files)

	log.Printf("Found %d migration files", len(files))

	// 执行每个迁移文件
	for _, file := range files {
		log.Printf("Executing migration: %s", filepath.Base(file))

		// 读取 SQL 文件内容
		sqlContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// 执行 SQL（按分号分割，逐条执行）
		statements := strings.Split(string(sqlContent), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			// 跳过空语句和注释
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}

			// 执行 SQL 语句
			if err := DB.Exec(stmt).Error; err != nil {
				// 对于某些错误（如表已存在），只记录警告而不中断
				if strings.Contains(err.Error(), "already exists") ||
					strings.Contains(err.Error(), "duplicate") {
					log.Printf("Warning: %s (this is usually safe to ignore)", err)
				} else {
					return fmt.Errorf("failed to execute migration %s: %w\nSQL: %s", file, err, stmt)
				}
			}
		}

		log.Printf("Migration %s executed successfully", filepath.Base(file))
	}

	// 执行非时间戳格式的迁移文件（如 fix_*.sql）
	legacyFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err == nil {
		for _, file := range legacyFiles {
			// 跳过已经处理过的 .up.sql 文件
			if strings.HasSuffix(file, ".up.sql") {
				continue
			}
			// 跳过 .down.sql 文件
			if strings.HasSuffix(file, ".down.sql") {
				continue
			}

			log.Printf("Executing legacy migration: %s", filepath.Base(file))

			sqlContent, err := os.ReadFile(file)
			if err != nil {
				log.Printf("Warning: Failed to read legacy migration file %s: %v", file, err)
				continue
			}

			statements := strings.Split(string(sqlContent), ";")
			for _, stmt := range statements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" || strings.HasPrefix(stmt, "--") {
					continue
				}

				if err := DB.Exec(stmt).Error; err != nil {
					if strings.Contains(err.Error(), "already exists") ||
						strings.Contains(err.Error(), "duplicate") {
						log.Printf("Warning: %s (this is usually safe to ignore)", err)
					} else {
						log.Printf("Warning: Failed to execute legacy migration %s: %v\nSQL: %s", file, err, stmt)
					}
				}
			}
		}
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
