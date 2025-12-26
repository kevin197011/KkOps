package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/handler"
	"github.com/kkops/backend/internal/middleware"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/salt"
	"github.com/kkops/backend/internal/service"
	"github.com/kkops/backend/internal/utils"
)

func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库连接
	log.Println("Connecting to database...")
	if err := models.InitDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v\nPlease check your database configuration in .env file", err)
	}
	log.Println("Database connection established")

	// 自动迁移数据库表（启动时自动执行）
	log.Println("Running database migration...")
	if err := models.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v\nMigration is required for the application to work properly", err)
	}
	log.Println("Database migration completed")

	// 从数据库加载配置（覆盖环境变量）
	log.Println("Loading configuration from database...")
	if err := loadConfigFromDatabase(); err != nil {
		log.Printf("Warning: Failed to load configuration from database: %v. Using environment variables.", err)
	} else {
		log.Println("Configuration loaded from database successfully")
	}

	// 初始化Salt客户端管理器并创建客户端（可选，如果配置了Salt API）
	saltManager := salt.GetManager()
	if config.AppConfig.Salt.APIURL != "" && config.AppConfig.Salt.Username != "" {
		saltManager.InitializeClient(config.AppConfig.Salt)
	} else {
		log.Println("Salt API not configured, skipping Salt client initialization")
	}

	// 初始化Repository
	userRepo := repository.NewUserRepository(models.DB)
	roleRepo := repository.NewRoleRepository(models.DB)
	permissionRepo := repository.NewPermissionRepository(models.DB)
	projectRepo := repository.NewProjectRepository(models.DB)
	environmentRepo := repository.NewEnvironmentRepository(models.DB)
	cloudPlatformRepo := repository.NewCloudPlatformRepository(models.DB)
	hostRepo := repository.NewHostRepository(models.DB)
	hostGroupRepo := repository.NewHostGroupRepository(models.DB)
	hostTagRepo := repository.NewHostTagRepository(models.DB)
	deploymentConfigRepo := repository.NewDeploymentConfigRepository(models.DB)
	deploymentRepo := repository.NewDeploymentRepository(models.DB)
	deploymentVersionRepo := repository.NewDeploymentVersionRepository(models.DB)
	auditRepo := repository.NewAuditRepository(models.DB)
	batchOperationRepo := repository.NewBatchOperationRepository(models.DB)
	commandTemplateRepo := repository.NewCommandTemplateRepository(models.DB)
	formulaRepo := repository.NewFormulaRepository(models.DB)
	sshKeyRepo := repository.NewSSHKeyRepository(models.DB)
	settingsRepo := repository.NewSettingsRepository(models.DB)

	// 初始化Service
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	projectService := service.NewProjectService(projectRepo)
	environmentService := service.NewEnvironmentService(environmentRepo)
	cloudPlatformService := service.NewCloudPlatformService(cloudPlatformRepo)
	// 使用Manager.GetClient()以支持热重载，但保持Service接口兼容
	hostService := service.NewHostService(hostRepo, environmentRepo, saltManager.GetClient())
	hostGroupService := service.NewHostGroupService(hostGroupRepo)
	hostTagService := service.NewHostTagService(hostTagRepo)
	deploymentConfigService := service.NewDeploymentConfigService(deploymentConfigRepo)
	deploymentService := service.NewDeploymentService(deploymentRepo, deploymentConfigRepo, hostRepo, saltManager.GetClient())
	deploymentVersionService := service.NewDeploymentVersionService(deploymentVersionRepo)
	auditService := service.NewAuditService(auditRepo)
	batchOperationService := service.NewBatchOperationService(batchOperationRepo, hostRepo, saltManager.GetClient())
	commandTemplateService := service.NewCommandTemplateService(commandTemplateRepo)
	formulaService := service.NewFormulaService(formulaRepo, saltManager.GetClient())
	sshKeyService := service.NewSSHKeyService(sshKeyRepo)
	settingsService := service.NewSettingsService(settingsRepo, auditRepo)

	// 初始化调度器服务
	schedulerService := service.NewSchedulerService(settingsService, hostService, batchOperationService)
	// 启动调度器
	schedulerService.Start()
	defer schedulerService.Stop()

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService, auditService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	projectHandler := handler.NewProjectHandler(projectService)
	environmentHandler := handler.NewEnvironmentHandler(environmentService)
	cloudPlatformHandler := handler.NewCloudPlatformHandler(cloudPlatformService)
	hostHandler := handler.NewHostHandler(hostService)
	hostGroupHandler := handler.NewHostGroupHandler(hostGroupService)
	hostTagHandler := handler.NewHostTagHandler(hostTagService)
	deploymentConfigHandler := handler.NewDeploymentConfigHandler(deploymentConfigService)
	deploymentHandler := handler.NewDeploymentHandler(deploymentService)
	deploymentVersionHandler := handler.NewDeploymentVersionHandler(deploymentVersionService)
	websshHandler := handler.NewWebSSHHandler(hostRepo, sshKeyService)
	sshKeyHandler := handler.NewSSHKeyHandler(sshKeyService)
	auditHandler := handler.NewAuditHandler(auditService)
	settingsHandler := handler.NewSettingsHandler(settingsService, saltManager) // 传入Manager以支持热重载

	// 批量操作相关Repository和Service
	batchOperationHandler := handler.NewBatchOperationHandler(batchOperationService)
	commandTemplateHandler := handler.NewCommandTemplateHandler(commandTemplateService)
	formulaHandler := handler.NewFormulaHandler(formulaService)

	// 初始化Gin路由
	r := gin.Default()

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查（支持 GET 和 HEAD）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// 公开路由（也需要审计）
	api := r.Group("/api/v1")
	api.Use(middleware.AuditMiddleware(auditService))
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
	}

	// 需要认证的路由
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	auth.Use(middleware.AuditMiddleware(auditService))
	{
		// 用户管理
		users := auth.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.POST("/:id/roles", userHandler.AssignRole)
			users.DELETE("/:id/roles", userHandler.RemoveRole)
			users.POST("/:id/reset-password", userHandler.ResetPassword)
		}

		// 密码管理
		auth.POST("/auth/change-password", userHandler.ChangePassword)

		// 角色管理
		roles := auth.Group("/roles")
		{
			roles.POST("", roleHandler.CreateRole)
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.POST("/:id/permissions", roleHandler.AssignPermission)
			roles.DELETE("/:id/permissions", roleHandler.RemovePermission)
		}

		// 权限管理
		permissions := auth.Group("/permissions")
		{
			permissions.POST("", permissionHandler.CreatePermission)
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)
			permissions.PUT("/:id", permissionHandler.UpdatePermission)
			permissions.DELETE("/:id", permissionHandler.DeletePermission)
			permissions.GET("/resource/:resource_type", permissionHandler.ListPermissionsByResourceType)
		}

		// 项目管理
		projects := auth.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.POST("/:id/members", projectHandler.AddMember)
			projects.DELETE("/:id/members", projectHandler.RemoveMember)
			projects.GET("/:id/members", projectHandler.GetMembers)
		}

		// 环境管理
		environments := auth.Group("/environments")
		{
			environments.POST("", environmentHandler.CreateEnvironment)
			environments.GET("", environmentHandler.ListEnvironments)
			environments.GET("/:id", environmentHandler.GetEnvironment)
			environments.PUT("/:id", environmentHandler.UpdateEnvironment)
			environments.DELETE("/:id", environmentHandler.DeleteEnvironment)
		}

		// 云平台管理
		cloudPlatforms := auth.Group("/cloud-platforms")
		{
			cloudPlatforms.POST("", cloudPlatformHandler.CreateCloudPlatform)
			cloudPlatforms.GET("", cloudPlatformHandler.ListCloudPlatforms)
			cloudPlatforms.GET("/:id", cloudPlatformHandler.GetCloudPlatform)
			cloudPlatforms.PUT("/:id", cloudPlatformHandler.UpdateCloudPlatform)
			cloudPlatforms.DELETE("/:id", cloudPlatformHandler.DeleteCloudPlatform)
		}

		// 主机管理
		hosts := auth.Group("/hosts")
		{
			hosts.POST("", hostHandler.CreateHost)
			hosts.GET("", hostHandler.ListHosts)
			hosts.GET("/:id", hostHandler.GetHost)
			hosts.PUT("/:id", hostHandler.UpdateHost)
			hosts.DELETE("/:id", hostHandler.DeleteHost)
			hosts.POST("/:id/groups", hostHandler.AddToGroup)
			hosts.DELETE("/:id/groups", hostHandler.RemoveFromGroup)
			hosts.POST("/:id/tags", hostHandler.AddTag)
			hosts.DELETE("/:id/tags", hostHandler.RemoveTag)
			hosts.POST("/:id/sync-status", hostHandler.SyncHostStatus)
			hosts.POST("/:id/auto-assign-minion", hostHandler.AutoAssignMinion)
			hosts.POST("/sync-all-status", hostHandler.SyncAllHostsStatus)
			hosts.GET("/discover-minions", hostHandler.DiscoverMinions)
		}

		// 主机组管理
		hostGroups := auth.Group("/host-groups")
		{
			hostGroups.POST("", hostGroupHandler.CreateGroup)
			hostGroups.GET("", hostGroupHandler.ListGroups)
			hostGroups.GET("/:id", hostGroupHandler.GetGroup)
			hostGroups.PUT("/:id", hostGroupHandler.UpdateGroup)
			hostGroups.DELETE("/:id", hostGroupHandler.DeleteGroup)
		}

		// 主机标签管理
		hostTags := auth.Group("/host-tags")
		{
			hostTags.POST("", hostTagHandler.CreateTag)
			hostTags.GET("", hostTagHandler.ListTags)
			hostTags.GET("/:id", hostTagHandler.GetTag)
			hostTags.PUT("/:id", hostTagHandler.UpdateTag)
			hostTags.DELETE("/:id", hostTagHandler.DeleteTag)
		}

		// WebSSH管理
		webssh := auth.Group("/webssh")
		{
			// WebSSH终端 WebSocket
			webssh.GET("/terminal/:host_id", websshHandler.HandleTerminal)
		}

		// SSH密钥管理
		sshKeys := auth.Group("/ssh-keys")
		{
			sshKeys.POST("", sshKeyHandler.CreateSSHKey)
			sshKeys.GET("", sshKeyHandler.ListSSHKeys)
			sshKeys.GET("/:id", sshKeyHandler.GetSSHKey)
			sshKeys.GET("/:id/private-key", sshKeyHandler.GetSSHKeyPrivateKey)
			sshKeys.PUT("/:id", sshKeyHandler.UpdateSSHKey)
			sshKeys.DELETE("/:id", sshKeyHandler.DeleteSSHKey)
		}

		// 部署管理
		deploymentConfigs := auth.Group("/deployment-configs")
		{
			deploymentConfigs.POST("", deploymentConfigHandler.CreateConfig)
			deploymentConfigs.GET("", deploymentConfigHandler.ListConfigs)
			deploymentConfigs.GET("/:id", deploymentConfigHandler.GetConfig)
			deploymentConfigs.PUT("/:id", deploymentConfigHandler.UpdateConfig)
			deploymentConfigs.DELETE("/:id", deploymentConfigHandler.DeleteConfig)
		}

		deployments := auth.Group("/deployments")
		{
			deployments.POST("", deploymentHandler.CreateDeployment)
			deployments.GET("", deploymentHandler.ListDeployments)
			deployments.GET("/:id", deploymentHandler.GetDeployment)
			deployments.GET("/:id/status", deploymentHandler.GetDeploymentStatus)
			deployments.POST("/:id/rollback", deploymentHandler.RollbackDeployment)
		}

		deploymentVersions := auth.Group("/deployment-versions")
		{
			deploymentVersions.POST("", deploymentVersionHandler.CreateVersion)
			deploymentVersions.GET("", deploymentVersionHandler.ListVersions)
			deploymentVersions.GET("/:id", deploymentVersionHandler.GetVersion)
			deploymentVersions.GET("/application/:application_name", deploymentVersionHandler.GetVersionsByApplication)
		}

		// Salt管理（如果配置了Salt API）
		// 使用Manager获取客户端，支持动态配置更新
		saltHandler := handler.NewSaltHandler(saltManager.GetClient())
		salt := auth.Group("/salt")
		{
			salt.POST("/execute", saltHandler.ExecuteCommand)
			salt.POST("/ping", saltHandler.TestPing)
			salt.POST("/grains", saltHandler.GetGrains)
		}

		// 审计管理
		audit := auth.Group("/audit")
		{
			audit.GET("/logs", auditHandler.ListAuditLogs)
			audit.GET("/logs/:id", auditHandler.GetAuditLog)
			audit.GET("/logs/resource", auditHandler.GetLogsByResource)
			audit.GET("/logs/user/:user_id", auditHandler.GetLogsByUser)
		}

		// 批量操作
		batchOperations := auth.Group("/batch-operations")
		{
			batchOperations.POST("", middleware.RequirePermission("batch_operation:create"), batchOperationHandler.CreateOperation)
			batchOperations.GET("", middleware.RequirePermission("batch_operation:read"), batchOperationHandler.ListOperations)
			batchOperations.GET("/:id", middleware.RequirePermission("batch_operation:read"), batchOperationHandler.GetOperation)
			batchOperations.GET("/:id/status", middleware.RequirePermission("batch_operation:read"), batchOperationHandler.GetOperationStatus)
			batchOperations.GET("/:id/results", middleware.RequirePermission("batch_operation:read"), batchOperationHandler.GetOperationResults)
			batchOperations.POST("/:id/cancel", middleware.RequirePermission("batch_operation:cancel"), batchOperationHandler.CancelOperation)
			batchOperations.POST("/cleanup", middleware.RequirePermission("batch_operation:create"), batchOperationHandler.CleanupOldOperations)
		}

		// 命令模板
		commandTemplates := auth.Group("/command-templates")
		{
			commandTemplates.POST("", commandTemplateHandler.CreateTemplate)
			commandTemplates.GET("", commandTemplateHandler.ListTemplates)
			commandTemplates.GET("/:id", commandTemplateHandler.GetTemplate)
			commandTemplates.PUT("/:id", commandTemplateHandler.UpdateTemplate)
			commandTemplates.DELETE("/:id", commandTemplateHandler.DeleteTemplate)
		}

		// Formula管理
		formulas := auth.Group("/formulas")
		{
			formulas.POST("", middleware.RequirePermission("batch_operation:create"), formulaHandler.CreateFormula)
			formulas.GET("", middleware.RequirePermission("batch_operation:read"), formulaHandler.ListFormulas)
			formulas.GET("/:id", middleware.RequirePermission("batch_operation:read"), formulaHandler.GetFormula)
			formulas.PUT("/:id/parameters", middleware.RequirePermission("batch_operation:create"), formulaHandler.UpdateFormulaParameters)

			// Formula模板
			formulas.POST("/:id/templates", middleware.RequirePermission("batch_operation:create"), formulaHandler.CreateFormulaTemplate)
			formulas.GET("/:id/templates", middleware.RequirePermission("batch_operation:read"), formulaHandler.GetFormulaTemplates)

			// Formula部署
			formulas.POST("/deployments", middleware.RequirePermission("batch_operation:create"), formulaHandler.CreateDeployment)
			formulas.GET("/deployments", middleware.RequirePermission("batch_operation:read"), formulaHandler.ListDeployments)
			formulas.GET("/deployments/:id", middleware.RequirePermission("batch_operation:read"), formulaHandler.GetDeployment)
			formulas.POST("/deployments/:id/execute", middleware.RequirePermission("batch_operation:create"), formulaHandler.ExecuteDeployment)
			formulas.POST("/deployments/:id/cancel", middleware.RequirePermission("batch_operation:create"), formulaHandler.CancelDeployment)
			formulas.POST("/deployments/cleanup", middleware.RequirePermission("batch_operation:create"), formulaHandler.CleanupOldDeployments)

			// Formula仓库
			formulas.POST("/repositories", middleware.RequirePermission("batch_operation:create"), formulaHandler.CreateRepository)
			formulas.GET("/repositories", middleware.RequirePermission("batch_operation:read"), formulaHandler.ListRepositories)
			formulas.GET("/repositories/:id", middleware.RequirePermission("batch_operation:read"), formulaHandler.GetRepository)
			formulas.PUT("/repositories/:id", middleware.RequirePermission("batch_operation:create"), formulaHandler.UpdateRepository)
			formulas.POST("/repositories/:id/sync", middleware.RequirePermission("batch_operation:create"), formulaHandler.SyncRepository)
		}

		// 系统设置（仅管理员）
		settings := auth.Group("/settings")
		settings.Use(middleware.RequireRole("admin"))
		{
			settings.GET("", settingsHandler.GetAllSettings)
			settings.GET("/:category", settingsHandler.GetSettingsByCategory)
			settings.GET("/key/:key", settingsHandler.GetSetting)
			settings.PUT("/key/:key", settingsHandler.UpdateSetting)
			settings.GET("/salt", settingsHandler.GetSaltConfig)
			settings.PUT("/salt", settingsHandler.UpdateSaltConfig)
			settings.POST("/salt/test", settingsHandler.TestSaltConnection)
			// Audit log settings endpoints
			settings.GET("/audit", settingsHandler.GetAuditLogSettings)
			settings.PUT("/audit", settingsHandler.UpdateAuditLogSettings)
			settings.GET("/audit/stats", settingsHandler.GetAuditLogStats)
			settings.POST("/audit/cleanup", settingsHandler.CleanupAuditLogs)
		}
	}

	// 启动服务器
	addr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// loadConfigFromDatabase 从数据库加载配置并覆盖环境变量配置
func loadConfigFromDatabase() error {
	if models.DB == nil {
		// 数据库未初始化，跳过数据库配置加载
		return nil
	}

	settingsRepo := repository.NewSettingsRepository(models.DB)

	// 加载Salt配置
	saltSettings, err := settingsRepo.GetByCategory("salt")
	if err != nil {
		// 如果出错（比如表不存在），使用环境变量配置，不返回错误
		return nil
	}

	// 如果有数据库设置，则覆盖环境变量配置
	if len(saltSettings) > 0 {
		for _, setting := range saltSettings {
			switch setting.Key {
			case "salt.api_url":
				config.AppConfig.Salt.APIURL = setting.Value
			case "salt.username":
				config.AppConfig.Salt.Username = setting.Value
			case "salt.password":
				// 解密密码
				decrypted, err := utils.Decrypt(setting.Value)
				if err == nil {
					config.AppConfig.Salt.Password = decrypted
				}
			case "salt.eauth":
				config.AppConfig.Salt.EAuth = setting.Value
			case "salt.timeout":
				if timeout, err := strconv.Atoi(setting.Value); err == nil {
					config.AppConfig.Salt.Timeout = timeout
				}
			case "salt.verify_ssl":
				config.AppConfig.Salt.VerifySSL = setting.Value == "true"
			}
		}
	}

	return nil
}
