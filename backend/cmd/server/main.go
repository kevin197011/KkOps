// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// @title           KkOps API
// @version         1.0
// @description     Intelligent Operations Management Platform API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@kkops.local

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/kkops/backend/api"
	"github.com/kkops/backend/internal/config"
	assetHandler "github.com/kkops/backend/internal/handler/asset"
	authHandler "github.com/kkops/backend/internal/handler/auth"
	categoryHandler "github.com/kkops/backend/internal/handler/category"
	cloudplatformHandler "github.com/kkops/backend/internal/handler/cloudplatform"
	environmentHandler "github.com/kkops/backend/internal/handler/environment"
	projectHandler "github.com/kkops/backend/internal/handler/project"
	roleHandler "github.com/kkops/backend/internal/handler/role"
	sshkeyHandler "github.com/kkops/backend/internal/handler/sshkey"
	tagHandler "github.com/kkops/backend/internal/handler/tag"
	taskHandler "github.com/kkops/backend/internal/handler/task"
	userHandler "github.com/kkops/backend/internal/handler/user"
	websocketHandler "github.com/kkops/backend/internal/handler/websocket"
	"github.com/kkops/backend/internal/middleware"
	assetService "github.com/kkops/backend/internal/service/asset"
	authService "github.com/kkops/backend/internal/service/auth"
	categoryService "github.com/kkops/backend/internal/service/category"
	cloudplatformService "github.com/kkops/backend/internal/service/cloudplatform"
	environmentService "github.com/kkops/backend/internal/service/environment"
	projectService "github.com/kkops/backend/internal/service/project"
	roleService "github.com/kkops/backend/internal/service/role"
	sshkeyService "github.com/kkops/backend/internal/service/sshkey"
	tagService "github.com/kkops/backend/internal/service/tag"
	taskService "github.com/kkops/backend/internal/service/task"
	userService "github.com/kkops/backend/internal/service/user"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	var zapLogger *zap.Logger
	if cfg.Log.Format == "json" {
		zapLogger, _ = zap.NewProduction()
	} else {
		zapLogger, _ = zap.NewDevelopment()
	}
	defer zapLogger.Sync()

	// Initialize database
	db, err := config.InitDB(&cfg.Database, zapLogger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	r := gin.New()

	// Set trusted proxies (empty slice means trust no proxies)
	// In Docker/Kubernetes environments, set to nil or specific proxy IPs
	r.SetTrustedProxies(nil)

	// Apply middleware
	r.Use(middleware.Logger(zapLogger))
	r.Use(middleware.ErrorHandler(zapLogger))
	r.Use(middleware.CORS())

	// Initialize services
	authSvc := authService.NewService(db, cfg)
	userSvc := userService.NewService(db)
	roleSvc := roleService.NewService(db)
	projectSvc := projectService.NewService(db)
	environmentSvc := environmentService.NewService(db)
	cloudplatformSvc := cloudplatformService.NewService(db)
	categorySvc := categoryService.NewService(db)
	tagSvc := tagService.NewService(db)
	assetSvc := assetService.NewService(db)
	taskSvc := taskService.NewService(db)
	sshkeySvc := sshkeyService.NewService(db, cfg)
	taskExecutionSvc := taskService.NewExecutionService(db, cfg, sshkeySvc)

	// Initialize handlers
	authHdl := authHandler.NewHandler(authSvc)
	userHdl := userHandler.NewHandler(userSvc)
	roleHdl := roleHandler.NewHandler(roleSvc)
	projectHdl := projectHandler.NewHandler(projectSvc)
	environmentHdl := environmentHandler.NewHandler(environmentSvc)
	cloudplatformHdl := cloudplatformHandler.NewHandler(cloudplatformSvc)
	categoryHdl := categoryHandler.NewHandler(categorySvc)
	tagHdl := tagHandler.NewHandler(tagSvc)
	assetHdl := assetHandler.NewHandler(assetSvc)
	taskHdl := taskHandler.NewHandler(taskSvc, taskExecutionSvc)
	sshkeyHdl := sshkeyHandler.NewHandler(sshkeySvc)

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHdl.Login)
			authGroup.GET("/me", middleware.AuthMiddleware(cfg), authHdl.GetMe)
			authGroup.POST("/logout", middleware.AuthMiddleware(cfg), authHdl.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// User management
			usersGroup := protected.Group("/users")
			{
				usersGroup.GET("", userHdl.ListUsers)
				usersGroup.POST("", userHdl.CreateUser)
				usersGroup.GET("/:id", userHdl.GetUser)
				usersGroup.PUT("/:id", userHdl.UpdateUser)
				usersGroup.DELETE("/:id", userHdl.DeleteUser)
				usersGroup.POST("/:id/reset-password", userHdl.ResetPassword)
				usersGroup.GET("/:id/tokens", authHdl.ListAPITokens)
				usersGroup.POST("/:id/tokens", authHdl.CreateAPIToken)
			}

			// API Token management
			tokensGroup := protected.Group("/tokens")
			{
				tokensGroup.DELETE("/:id", authHdl.RevokeAPIToken)
			}

			// Role management
			rolesGroup := protected.Group("/roles")
			{
				rolesGroup.GET("", roleHdl.ListRoles)
				rolesGroup.POST("", roleHdl.CreateRole)
				rolesGroup.GET("/:id", roleHdl.GetRole)
				rolesGroup.PUT("/:id", roleHdl.UpdateRole)
				rolesGroup.DELETE("/:id", roleHdl.DeleteRole)
				rolesGroup.GET("/:id/permissions", roleHdl.GetRolePermissions)
				rolesGroup.POST("/:id/permissions", roleHdl.AssignPermissionToRole)
			}

			// Permission management
			permissionsGroup := protected.Group("/permissions")
			{
				permissionsGroup.GET("", roleHdl.ListPermissions)
			}

			// User role assignment
			usersGroup.POST("/:id/roles", roleHdl.AssignRoleToUser)

			// Project management
			projectsGroup := protected.Group("/projects")
			{
				projectsGroup.GET("", projectHdl.ListProjects)
				projectsGroup.POST("", projectHdl.CreateProject)
				projectsGroup.GET("/:id", projectHdl.GetProject)
				projectsGroup.PUT("/:id", projectHdl.UpdateProject)
				projectsGroup.DELETE("/:id", projectHdl.DeleteProject)
			}

			// Environment management
			environmentsGroup := protected.Group("/environments")
			{
				environmentsGroup.GET("", environmentHdl.ListEnvironments)
				environmentsGroup.POST("", environmentHdl.CreateEnvironment)
				environmentsGroup.GET("/:id", environmentHdl.GetEnvironment)
				environmentsGroup.PUT("/:id", environmentHdl.UpdateEnvironment)
				environmentsGroup.DELETE("/:id", environmentHdl.DeleteEnvironment)
			}

			// Cloud platform management
			cloudPlatformsGroup := protected.Group("/cloud-platforms")
			{
				cloudPlatformsGroup.GET("", cloudplatformHdl.ListCloudPlatforms)
				cloudPlatformsGroup.POST("", cloudplatformHdl.CreateCloudPlatform)
				cloudPlatformsGroup.GET("/:id", cloudplatformHdl.GetCloudPlatform)
				cloudPlatformsGroup.PUT("/:id", cloudplatformHdl.UpdateCloudPlatform)
				cloudPlatformsGroup.DELETE("/:id", cloudplatformHdl.DeleteCloudPlatform)
			}

			// Asset Category management
			categoriesGroup := protected.Group("/asset-categories")
			{
				categoriesGroup.GET("", categoryHdl.ListCategories)
				categoriesGroup.POST("", categoryHdl.CreateCategory)
				categoriesGroup.GET("/:id", categoryHdl.GetCategory)
				categoriesGroup.PUT("/:id", categoryHdl.UpdateCategory)
				categoriesGroup.DELETE("/:id", categoryHdl.DeleteCategory)
			}

			// Tag management
			tagsGroup := protected.Group("/tags")
			{
				tagsGroup.GET("", tagHdl.ListTags)
				tagsGroup.POST("", tagHdl.CreateTag)
				tagsGroup.GET("/:id", tagHdl.GetTag)
				tagsGroup.PUT("/:id", tagHdl.UpdateTag)
				tagsGroup.DELETE("/:id", tagHdl.DeleteTag)
			}

			// Asset management
			assetsGroup := protected.Group("/assets")
			{
				assetsGroup.GET("", assetHdl.ListAssets)
				assetsGroup.POST("", assetHdl.CreateAsset)
				assetsGroup.GET("/:id", assetHdl.GetAsset)
				assetsGroup.PUT("/:id", assetHdl.UpdateAsset)
				assetsGroup.DELETE("/:id", assetHdl.DeleteAsset)
				assetsGroup.POST("/import", assetHdl.ImportAssets)
				assetsGroup.GET("/export", assetHdl.ExportAssets)
			}

			// Task template management
			templatesGroup := protected.Group("/task-templates")
			{
				templatesGroup.GET("", taskHdl.ListTemplates)
				templatesGroup.POST("", taskHdl.CreateTemplate)
				templatesGroup.GET("/:id", taskHdl.GetTemplate)
				templatesGroup.PUT("/:id", taskHdl.UpdateTemplate)
				templatesGroup.DELETE("/:id", taskHdl.DeleteTemplate)
			}

			// Task management
			tasksGroup := protected.Group("/tasks")
			{
				tasksGroup.GET("", taskHdl.ListTasks)
				tasksGroup.POST("", taskHdl.CreateTask)
				tasksGroup.GET("/:id", taskHdl.GetTask)
				tasksGroup.PUT("/:id", taskHdl.UpdateTask)
				tasksGroup.DELETE("/:id", taskHdl.DeleteTask)
				tasksGroup.POST("/:id/execute", taskHdl.ExecuteTask)
				tasksGroup.POST("/:id/cancel", taskHdl.CancelTask)
				tasksGroup.GET("/:id/executions", taskHdl.GetTaskExecutions)
			}

			// Task execution management
			executionsGroup := protected.Group("/task-executions")
			{
				executionsGroup.GET("/:id", taskHdl.GetTaskExecution)
				executionsGroup.GET("/:id/logs", taskHdl.GetTaskExecutionLogs)
				executionsGroup.POST("/:id/cancel", taskHdl.CancelTaskExecution)
			}

			// SSH key management
			sshKeysGroup := protected.Group("/ssh/keys")
			{
				sshKeysGroup.GET("", sshkeyHdl.ListSSHKeys)
				sshKeysGroup.POST("", sshkeyHdl.CreateSSHKey)
				sshKeysGroup.GET("/:id", sshkeyHdl.GetSSHKey)
				sshKeysGroup.PUT("/:id", sshkeyHdl.UpdateSSHKey)
				sshKeysGroup.DELETE("/:id", sshkeyHdl.DeleteSSHKey)
				sshKeysGroup.POST("/:id/test", sshkeyHdl.TestSSHKey)
			}
		}
	}

	// WebSocket routes (protected)
	ws := r.Group("/ws")
	ws.Use(middleware.AuthMiddleware(cfg))
	{
		ws.GET("/task-executions/:id/logs", websocketHandler.StreamExecutionLogs(db))
		ws.GET("/ssh/connect", websocketHandler.SSHTerminalHandler(db, cfg, sshkeySvc))
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	zapLogger.Info("Starting server", zap.String("address", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
