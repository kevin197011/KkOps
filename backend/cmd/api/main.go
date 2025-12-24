package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kronos/backend/internal/config"
	"github.com/kronos/backend/internal/handler"
	"github.com/kronos/backend/internal/middleware"
	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
	"github.com/kronos/backend/internal/salt"
	"github.com/kronos/backend/internal/service"
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

	// 初始化Salt客户端（可选，如果配置了Salt API）
	var saltClient *salt.Client
	if config.AppConfig.Salt.APIURL != "" && config.AppConfig.Salt.Username != "" {
		saltClient = salt.NewClient(config.AppConfig.Salt)
		log.Println("Salt API client initialized")
	} else {
		log.Println("Salt API not configured, skipping Salt client initialization")
	}

	// 初始化Repository
	userRepo := repository.NewUserRepository(models.DB)
	roleRepo := repository.NewRoleRepository(models.DB)
	permissionRepo := repository.NewPermissionRepository(models.DB)
	projectRepo := repository.NewProjectRepository(models.DB)
	hostRepo := repository.NewHostRepository(models.DB)
	hostGroupRepo := repository.NewHostGroupRepository(models.DB)
	hostTagRepo := repository.NewHostTagRepository(models.DB)
	deploymentConfigRepo := repository.NewDeploymentConfigRepository(models.DB)
	deploymentRepo := repository.NewDeploymentRepository(models.DB)
	deploymentVersionRepo := repository.NewDeploymentVersionRepository(models.DB)
	auditRepo := repository.NewAuditRepository(models.DB)
	sshKeyRepo := repository.NewSSHKeyRepository(models.DB)

	// 初始化Service
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	projectService := service.NewProjectService(projectRepo)
	hostService := service.NewHostService(hostRepo, saltClient)
	hostGroupService := service.NewHostGroupService(hostGroupRepo)
	hostTagService := service.NewHostTagService(hostTagRepo)
	deploymentConfigService := service.NewDeploymentConfigService(deploymentConfigRepo)
	deploymentService := service.NewDeploymentService(deploymentRepo, deploymentConfigRepo, hostRepo, saltClient)
	deploymentVersionService := service.NewDeploymentVersionService(deploymentVersionRepo)
	auditService := service.NewAuditService(auditRepo)
	sshKeyService := service.NewSSHKeyService(sshKeyRepo)

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService, auditService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	projectHandler := handler.NewProjectHandler(projectService)
	hostHandler := handler.NewHostHandler(hostService)
	hostGroupHandler := handler.NewHostGroupHandler(hostGroupService)
	hostTagHandler := handler.NewHostTagHandler(hostTagService)
	deploymentConfigHandler := handler.NewDeploymentConfigHandler(deploymentConfigService)
	deploymentHandler := handler.NewDeploymentHandler(deploymentService)
	deploymentVersionHandler := handler.NewDeploymentVersionHandler(deploymentVersionService)
	websshHandler := handler.NewWebSSHHandler(hostRepo, sshKeyService)
	sshKeyHandler := handler.NewSSHKeyHandler(sshKeyService)
	saltHandler := handler.NewSaltHandler(saltClient)
	auditHandler := handler.NewAuditHandler(auditService)

	// 初始化Gin路由
	r := gin.Default()

	// 审计日志中间件（放在最前面，记录所有请求）
	r.Use(middleware.AuditMiddleware(auditService))

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

	// 公开路由
	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
	}

	// 需要认证的路由
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
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
			hosts.POST("/sync-all-status", hostHandler.SyncAllHostsStatus)
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
		if saltClient != nil {
			salt := auth.Group("/salt")
			{
				salt.POST("/execute", saltHandler.ExecuteCommand)
				salt.POST("/ping", saltHandler.TestPing)
				salt.POST("/grains", saltHandler.GetGrains)
			}
		}

		// 审计管理
		audit := auth.Group("/audit")
		{
			audit.GET("/logs", auditHandler.ListAuditLogs)
			audit.GET("/logs/:id", auditHandler.GetAuditLog)
			audit.GET("/logs/resource", auditHandler.GetLogsByResource)
			audit.GET("/logs/user/:user_id", auditHandler.GetLogsByUser)
		}
	}

	// 启动服务器
	addr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
