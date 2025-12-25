package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	dsn := config.AppConfig.Database.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return gorm.ErrInvalidDB
	}

	log.Println("Starting database migration...")

	// 定义所有需要迁移的模型
	models := []interface{}{
		// 用户和权限
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
		// 项目管理
		&Project{},
		&ProjectMember{},
		// 主机管理
		&Host{},
		&HostGroup{},
		&HostTag{},
		&HostGroupMember{},
		&HostTagAssignment{},
		// SSH密钥管理
		&SSHKey{},
		// 发布管理
		&DeploymentConfig{},
		&Deployment{},
		&DeploymentVersion{},
		// 审计管理
		&AuditLog{},
		// 系统设置
		&SystemSettings{},
	}

	// 逐个迁移模型，以便更好地跟踪进度
	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate model: %v, error: %v", model, err)
			return err
		}
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
				DB.Where("resource_type IN (?)", []string{"host", "deployment", "task", "log", "monitoring", "webssh"}).Find(&opsPerms)
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

	log.Println("Default data initialization completed")
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
