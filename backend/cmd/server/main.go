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

// @host      localhost:3000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/kkops/backend/api"
	"github.com/kkops/backend/internal/ai"
	"github.com/kkops/backend/internal/config"
	aiapi "github.com/kkops/backend/internal/handler/aiapi"
	alerthandler "github.com/kkops/backend/internal/handler/alert"
	assetHandler "github.com/kkops/backend/internal/handler/asset"
	auditHandler "github.com/kkops/backend/internal/handler/audit"
	authHandler "github.com/kkops/backend/internal/handler/auth"
	categoryHandler "github.com/kkops/backend/internal/handler/category"
	cicdHandler "github.com/kkops/backend/internal/handler/cicd"
	cloudplatformHandler "github.com/kkops/backend/internal/handler/cloudplatform"
	cmdbHandler "github.com/kkops/backend/internal/handler/cmdb"
	connectionauditHandler "github.com/kkops/backend/internal/handler/connectionaudit"
	dashboardHandler "github.com/kkops/backend/internal/handler/dashboard"
	deploymentHandler "github.com/kkops/backend/internal/handler/deployment"
	environmentHandler "github.com/kkops/backend/internal/handler/environment"
	externalsystemHandler "github.com/kkops/backend/internal/handler/externalsystem"
	gitopsHandler "github.com/kkops/backend/internal/handler/gitops"
	incidenthandler "github.com/kkops/backend/internal/handler/incident"
	integrationHandler "github.com/kkops/backend/internal/handler/integration"
	k8sHandler "github.com/kkops/backend/internal/handler/k8s"
	loggingHandler "github.com/kkops/backend/internal/handler/logging"
	monitoringHandler "github.com/kkops/backend/internal/handler/monitoring"
	oauth2clientHandler "github.com/kkops/backend/internal/handler/oauth2client"
	operationtoolHandler "github.com/kkops/backend/internal/handler/operationtool"
	projectHandler "github.com/kkops/backend/internal/handler/project"
	provhandler "github.com/kkops/backend/internal/handler/provisioning"
	registryHandler "github.com/kkops/backend/internal/handler/registry"
	roleHandler "github.com/kkops/backend/internal/handler/role"
	scheduledtaskHandler "github.com/kkops/backend/internal/handler/scheduledtask"
	sshkeyHandler "github.com/kkops/backend/internal/handler/sshkey"
	tagHandler "github.com/kkops/backend/internal/handler/tag"
	taskHandler "github.com/kkops/backend/internal/handler/task"
	userHandler "github.com/kkops/backend/internal/handler/user"
	websocketHandler "github.com/kkops/backend/internal/handler/websocket"
	"github.com/kkops/backend/internal/idp"
	integrationk8s "github.com/kkops/backend/internal/integration/k8s"
	intprovider "github.com/kkops/backend/internal/integration/provider"
	"github.com/kkops/backend/internal/middleware"
	"github.com/kkops/backend/internal/provisioning"
	aisvc "github.com/kkops/backend/internal/service/aisvc"
	alertService "github.com/kkops/backend/internal/service/alert"
	assetService "github.com/kkops/backend/internal/service/asset"
	auditService "github.com/kkops/backend/internal/service/audit"
	authService "github.com/kkops/backend/internal/service/auth"
	authorizationService "github.com/kkops/backend/internal/service/authorization"
	categoryService "github.com/kkops/backend/internal/service/category"
	cloudplatformService "github.com/kkops/backend/internal/service/cloudplatform"
	cmdbService "github.com/kkops/backend/internal/service/cmdb"
	connectionauditService "github.com/kkops/backend/internal/service/connectionaudit"
	dashboardService "github.com/kkops/backend/internal/service/dashboard"
	deploymentService "github.com/kkops/backend/internal/service/deployment"
	environmentService "github.com/kkops/backend/internal/service/environment"
	externalsystemService "github.com/kkops/backend/internal/service/externalsystem"
	gitopsviewService "github.com/kkops/backend/internal/service/gitopsview"
	incidentService "github.com/kkops/backend/internal/service/incident"
	integrationService "github.com/kkops/backend/internal/service/integration"
	operationtoolService "github.com/kkops/backend/internal/service/operationtool"
	projectService "github.com/kkops/backend/internal/service/project"
	rbacService "github.com/kkops/backend/internal/service/rbac"
	roleService "github.com/kkops/backend/internal/service/role"
	scheduledtaskService "github.com/kkops/backend/internal/service/scheduledtask"
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
	sshkeySvc := sshkeyService.NewService(db, cfg)
	authzSvc := authorizationService.NewService(db) // 授权服务
	rbacSvc := rbacService.NewService(db)           // RBAC 服务
	taskSvc := taskService.NewService(db, authzSvc)
	taskExecutionSvc := taskService.NewExecutionService(db, cfg, sshkeySvc)
	dashboardSvc := dashboardService.NewService(db)
	cmdbSvc := cmdbService.NewService(db)
	deploymentSvc := deploymentService.NewService(db, cfg)
	scheduledTaskSvc := scheduledtaskService.NewService(db)
	auditSvc := auditService.NewService(db)
	connectionauditSvc := connectionauditService.NewService(db)
	operationtoolSvc := operationtoolService.NewService(db)
	externalsystemSvc := externalsystemService.NewService(db)
	integrationSvc := integrationService.NewService(db, cfg)
	gitopsViewSvc := gitopsviewService.NewService(integrationSvc)
	k8sRegistry := integrationk8s.NewClusterRegistry()
	alertSvc := alertService.NewService(db, integrationSvc)
	incidentSvc := incidentService.NewService(db)

	aiToolBridge := &aisvc.ToolBridge{
		Integration: integrationSvc,
		Alerts:      alertSvc,
		Incidents:   incidentSvc,
		GitOpsView:  gitopsViewSvc,
		K8sRegistry: k8sRegistry,
	}
	aiOrchestrationSvc := aisvc.NewService(db, zapLogger, integrationSvc, alertSvc, incidentSvc, aiToolBridge)
	aiRegistry := ai.NewRegistry(db, integrationSvc)
	aiHdl := aiapi.NewHandler(aiOrchestrationSvc, integrationSvc, aiRegistry)

	aiCron := cron.New()
	defer aiCron.Stop()
	if err := aiOrchestrationSvc.ScheduleAnomalyCron(aiCron, aiRegistry); err != nil {
		zapLogger.Warn("AI anomaly cron registration failed", zap.Error(err))
	}
	aiCron.Start()

	provReg := provisioning.NewRegistry()
	provCoord := provisioning.NewCoordinator(db, integrationSvc, provReg, zapLogger, 4, 256)
	provAPISvc := provisioning.NewAPIService(db, provReg, provCoord)
	provHdl := provhandler.NewHandler(provAPISvc)

	// Initialize scheduler for scheduled tasks
	scheduler := scheduledtaskService.NewScheduler(db, cfg, zapLogger)
	// 将调度器关联到服务，使新建的任务能被添加到调度器
	scheduledTaskSvc.SetScheduler(scheduler)
	if err := scheduler.Start(); err != nil {
		zapLogger.Error("启动定时任务调度器失败", zap.Error(err))
	} else {
		defer scheduler.Stop()
	}

	// Initialize handlers
	authHdl := authHandler.NewHandler(authSvc)
	userHdl := userHandler.NewHandler(userSvc, authzSvc, rbacSvc, provCoord)
	roleHdl := roleHandler.NewHandler(roleSvc)
	projectHdl := projectHandler.NewHandler(projectSvc)
	environmentHdl := environmentHandler.NewHandler(environmentSvc)
	cloudplatformHdl := cloudplatformHandler.NewHandler(cloudplatformSvc)
	categoryHdl := categoryHandler.NewHandler(categorySvc)
	tagHdl := tagHandler.NewHandler(tagSvc)
	assetHdl := assetHandler.NewHandler(assetSvc, authzSvc)
	taskHdl := taskHandler.NewHandler(taskSvc, taskExecutionSvc)
	sshkeyHdl := sshkeyHandler.NewHandler(sshkeySvc)
	dashboardHdl := dashboardHandler.NewHandler(dashboardSvc)
	cmdbHdl := cmdbHandler.NewHandler(cmdbSvc)
	deploymentHdl := deploymentHandler.NewHandler(deploymentSvc)
	scheduledTaskHdl := scheduledtaskHandler.NewHandler(scheduledTaskSvc)
	operationtoolHdl := operationtoolHandler.NewHandler(operationtoolSvc)
	externalsystemHdl := externalsystemHandler.NewHandler(externalsystemSvc, authSvc, rbacSvc)
	roleAssetHdl := roleHandler.NewAssetHandler(authzSvc)
	userRoleHdl := userHandler.NewRoleHandler(authzSvc)
	auditHdl := auditHandler.NewHandler(auditSvc)
	connectionauditHdl := connectionauditHandler.NewHandler(connectionauditSvc)
	idpHdl := idp.NewHandler(cfg, db, authSvc)
	oauth2clientHdl := oauth2clientHandler.NewHandler(db)
	integRegistry := intprovider.NewRegistry()
	integRegistry.RegisterDefaults()
	integrationHdl := integrationHandler.NewHandler(integrationSvc, integRegistry)
	monitoringHdl := monitoringHandler.NewHandler(integrationSvc)
	loggingHdl := loggingHandler.NewHandler(integrationSvc)
	cicdHdl := cicdHandler.NewHandler(integrationSvc)
	registryHdl := registryHandler.NewHandler(integrationSvc)
	gitopsHdl := gitopsHandler.NewHandler(integrationSvc, gitopsViewSvc)
	k8sHdl := k8sHandler.NewHandler(integrationSvc, k8sRegistry)
	alertHdl := alerthandler.NewHandler(alertSvc)
	incidentHdl := incidenthandler.NewHandler(incidentSvc)

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check endpoint (public, for container health checks)
		healthHandler := func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		}
		api.GET("/health", healthHandler)
		api.HEAD("/health", healthHandler)

		// Public alert webhook (secret header; no JWT)
		api.POST("/alerts/webhook", middleware.AlertWebhookSecret(cfg), alertHdl.Webhook)

		// Auth routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHdl.Login)
			authGroup.POST("/refresh", authHdl.Refresh)
			authGroup.GET("/me", middleware.AuthMiddleware(cfg), authHdl.GetMe)
			authGroup.POST("/logout", middleware.AuthMiddleware(cfg), authHdl.Logout)
			authGroup.POST("/change-password", middleware.AuthMiddleware(cfg), authHdl.ChangePassword)
			authGroup.GET("/sso/config", authHdl.GetSSOConfig)
			authGroup.GET("/sso/login", authHdl.SSOLogin)
			authGroup.GET("/sso/callback", authHdl.SSOCallback)
		}

		// User permissions (protected, but available before full menu load)
		userPermGroup := api.Group("/user")
		userPermGroup.Use(middleware.AuthMiddleware(cfg))
		{
			userPermGroup.GET("/permissions", userHdl.GetUserPermissions)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		// 应用审计中间件
		protected.Use(middleware.AuditMiddleware(auditSvc, middleware.GetDefaultAuditRoutes()))
		// 应用权限检查中间件（自动根据路由路径检查权限）
		protected.Use(middleware.RequireMenuPermission(rbacSvc))
		{
			// Audit logs (管理员查看)
			auditLogsGroup := protected.Group("/audit-logs")
			{
				auditLogsGroup.GET("", auditHdl.ListAuditLogs)
				auditLogsGroup.GET("/export", auditHdl.ExportAuditLogs)
				auditLogsGroup.GET("/modules", auditHdl.GetModules)
				auditLogsGroup.GET("/actions", auditHdl.GetActions)
				auditLogsGroup.GET("/:id", auditHdl.GetAuditLog)
			}

			// Connection audit (审计连线)
			connectionAuditGroup := protected.Group("/connection-audit")
			{
				connectionAuditGroup.GET("", connectionauditHdl.ListConnectionRecords)
				connectionAuditGroup.GET("/:id", connectionauditHdl.GetConnectionRecord)
			}

			// Dashboard
			dashboardGroup := protected.Group("/dashboard")
			{
				dashboardGroup.GET("/stats", dashboardHdl.GetStats)
				dashboardGroup.GET("/summary", dashboardHdl.GetSummary)
			}

			cmdbGrp := protected.Group("/cmdb")
			topoGrp := protected.Group("/topology")
			cmdbHdl.Register(cmdbGrp, topoGrp)

			// Operation tools (运维导航)
			operationToolsGroup := protected.Group("/operation-tools")
			{
				operationToolsGroup.GET("", operationtoolHdl.List)
				operationToolsGroup.GET("/:id", operationtoolHdl.Get)
				operationToolsGroup.POST("", operationtoolHdl.Create)
				operationToolsGroup.PUT("/:id", operationtoolHdl.Update)
				operationToolsGroup.DELETE("/:id", operationtoolHdl.Delete)
			}

			// External ops systems (SSO launch with JWT + permissions)
			externalSystemsGroup := protected.Group("/external-systems")
			{
				externalSystemsGroup.GET("", externalsystemHdl.List)
				externalSystemsGroup.GET("/:id", externalsystemHdl.Get)
				externalSystemsGroup.POST("", externalsystemHdl.Create)
				externalSystemsGroup.PUT("/:id", externalsystemHdl.Update)
				externalSystemsGroup.DELETE("/:id", externalsystemHdl.Delete)
				externalSystemsGroup.POST("/:id/launch", externalsystemHdl.Launch)
			}

			// OAuth2 clients (IdP applications)
			oauth2ClientsGroup := protected.Group("/oauth2-clients")
			{
				oauth2ClientsGroup.GET("", oauth2clientHdl.List)
				oauth2ClientsGroup.POST("", oauth2clientHdl.Create)
				oauth2ClientsGroup.GET("/:id", oauth2clientHdl.Get)
				oauth2ClientsGroup.PUT("/:id", oauth2clientHdl.Update)
				oauth2ClientsGroup.DELETE("/:id", oauth2clientHdl.Delete)
				oauth2ClientsGroup.POST("/:id/regenerate-secret", oauth2clientHdl.RegenerateSecret)
			}

			integrationsGroup := protected.Group("/integrations")
			{
				integrationsGroup.GET("", integrationHdl.List)
				integrationsGroup.POST("", integrationHdl.Create)
				integrationsGroup.GET("/:id", integrationHdl.Get)
				integrationsGroup.PUT("/:id", integrationHdl.Update)
				integrationsGroup.DELETE("/:id", integrationHdl.Delete)
				integrationsGroup.POST("/:id/test", integrationHdl.Test)
			}

			monitoringGroup := protected.Group("/monitoring")
			{
				monitoringGroup.POST("/query", monitoringHdl.Query)
			}

			loggingGroup := protected.Group("/logging")
			{
				loggingGroup.POST("/search", loggingHdl.Search)
			}

			cicdGroup := protected.Group("/cicd")
			{
				cicdGroup.GET("/pipelines", cicdHdl.ListPipelines)
				cicdGroup.POST("/pipelines/:id/run", cicdHdl.RunPipeline)
				cicdGroup.GET("/pipelines/:id/logs", cicdHdl.PipelineLogs)
			}

			registryGroup := protected.Group("/registry")
			{
				registryGroup.GET("/repositories", registryHdl.ListRepositories)
				registryGroup.GET("/tags", registryHdl.ListTags)
				registryGroup.GET("/vulnerabilities", registryHdl.Vulnerabilities)
			}

			gitopsGroup := protected.Group("/gitops")
			{
				gitopsGroup.GET("/applications", gitopsHdl.ListApplications)
				gitopsGroup.GET("/pipeline-view", gitopsHdl.PipelineView)
				gitopsGroup.POST("/applications/:name/sync", gitopsHdl.SyncApplication)
			}

			k8sGroup := protected.Group("/k8s/clusters")
			{
				k8sGroup.GET("/:id/namespaces", k8sHdl.ListNamespaces)
				k8sGroup.GET("/:id/workloads", k8sHdl.ListWorkloads)
				k8sGroup.GET("/:id/pods", k8sHdl.ListPods)
				k8sGroup.GET("/:id/pods/:name/logs", k8sHdl.PodLogs)
				k8sGroup.GET("/:id/events", k8sHdl.ListEvents)
				k8sGroup.GET("/:id/nodes", k8sHdl.ListNodes)
			}

			alertsGroup := protected.Group("/alerts")
			{
				alertsGroup.GET("", alertHdl.List)
				alertsGroup.POST("/sync", alertHdl.Sync)
				alertsGroup.POST("/:id/acknowledge", alertHdl.Acknowledge)
				alertsGroup.POST("/:id/dismiss", alertHdl.Dismiss)
			}

			incidentsGroup := protected.Group("/incidents")
			{
				incidentsGroup.GET("", incidentHdl.List)
				incidentsGroup.POST("", incidentHdl.Create)
				incidentsGroup.GET("/:id", incidentHdl.Get)
				incidentsGroup.PATCH("/:id", incidentHdl.Patch)
			}

			aiGroup := protected.Group("/ai")
			aiHdl.Register(aiGroup)

			provisioningGroup := protected.Group("/provisioning")
			{
				provisioningGroup.GET("/targets", provHdl.ListTargets)
				provisioningGroup.POST("/targets", provHdl.CreateTarget)
				provisioningGroup.POST("/targets/:id/sync", provHdl.SyncTarget)
				provisioningGroup.GET("/targets/:id/runs", provHdl.ListRuns)
			}

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
				usersGroup.GET("/:id/tokens/:token_id", authHdl.GetAPIToken)
				usersGroup.DELETE("/:id/tokens/:token_id", authHdl.DeleteAPIToken)
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
				rolesGroup.DELETE("/:id/permissions/:permission_id", roleHdl.RemovePermissionFromRole)
				// 角色资产授权管理
				rolesGroup.GET("/:id/assets", roleAssetHdl.GetRoleAssets)
				rolesGroup.POST("/:id/assets", roleAssetHdl.GrantRoleAssets)
				rolesGroup.DELETE("/:id/assets/:asset_id", roleAssetHdl.RevokeSingleRoleAsset)
				rolesGroup.DELETE("/:id/assets", roleAssetHdl.RevokeRoleAssets)
			}

			// Permission management
			permissionsGroup := protected.Group("/permissions")
			{
				permissionsGroup.GET("", roleHdl.ListPermissions)
			}

			// User role management
			usersGroup.GET("/:id/roles", userRoleHdl.GetUserRoles)
			usersGroup.POST("/:id/roles", userRoleHdl.SetUserRoles)
			usersGroup.DELETE("/:id/roles/:role_id", userRoleHdl.RemoveUserRole)

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

			// Execution template management (原 task-templates)
			templatesGroup := protected.Group("/templates")
			{
				templatesGroup.GET("", taskHdl.ListTemplates)
				templatesGroup.POST("", taskHdl.CreateTemplate)
				templatesGroup.GET("/export", taskHdl.ExportTemplates)
				templatesGroup.POST("/import", taskHdl.ImportTemplates)
				templatesGroup.GET("/:id", taskHdl.GetTemplate)
				templatesGroup.PUT("/:id", taskHdl.UpdateTemplate)
				templatesGroup.DELETE("/:id", taskHdl.DeleteTemplate)
			}

			// Execution management (原 tasks，运维执行)
			executionsGroup := protected.Group("/executions")
			{
				executionsGroup.GET("", taskHdl.ListTasks)
				executionsGroup.POST("", taskHdl.CreateTask)
				executionsGroup.GET("/:id", taskHdl.GetTask)
				executionsGroup.PUT("/:id", taskHdl.UpdateTask)
				executionsGroup.DELETE("/:id", taskHdl.DeleteTask)
				executionsGroup.POST("/:id/execute", taskHdl.ExecuteTask)
				executionsGroup.POST("/:id/cancel", taskHdl.CancelTask)
				executionsGroup.GET("/:id/history", taskHdl.GetTaskExecutions)
			}

			// Execution record management (原 task-executions)
			executionRecordsGroup := protected.Group("/execution-records")
			{
				executionRecordsGroup.GET("/:id", taskHdl.GetTaskExecution)
				executionRecordsGroup.GET("/:id/logs", taskHdl.GetTaskExecutionLogs)
				executionRecordsGroup.POST("/:id/cancel", taskHdl.CancelTaskExecution)
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

			// Deployment module management
			deploymentModulesGroup := protected.Group("/deployment-modules")
			{
				deploymentModulesGroup.GET("", deploymentHdl.ListModules)
				deploymentModulesGroup.POST("", deploymentHdl.CreateModule)
				deploymentModulesGroup.GET("/export", deploymentHdl.ExportModules)
				deploymentModulesGroup.POST("/import", deploymentHdl.ImportModules)
				deploymentModulesGroup.POST("/import/preview", deploymentHdl.PreviewImport)
				deploymentModulesGroup.GET("/:id", deploymentHdl.GetModule)
				deploymentModulesGroup.PUT("/:id", deploymentHdl.UpdateModule)
				deploymentModulesGroup.DELETE("/:id", deploymentHdl.DeleteModule)
				deploymentModulesGroup.GET("/:id/versions", deploymentHdl.GetVersions)
				deploymentModulesGroup.POST("/:id/deploy", deploymentHdl.Deploy)
			}

			// Deployment history management
			deploymentsGroup := protected.Group("/deployments")
			{
				deploymentsGroup.GET("", deploymentHdl.ListDeployments)
				deploymentsGroup.GET("/:id", deploymentHdl.GetDeployment)
				deploymentsGroup.POST("/:id/cancel", deploymentHdl.CancelDeployment)
			}

			// Scheduled task management (定时任务)
			tasksGroup := protected.Group("/tasks")
			{
				tasksGroup.GET("", scheduledTaskHdl.ListScheduledTasks)
				tasksGroup.POST("", scheduledTaskHdl.CreateScheduledTask)
				tasksGroup.GET("/validate-cron", scheduledTaskHdl.ValidateCron)
				tasksGroup.GET("/export", scheduledTaskHdl.ExportScheduledTasks)
				tasksGroup.POST("/import", scheduledTaskHdl.ImportScheduledTasks)
				tasksGroup.GET("/:id", scheduledTaskHdl.GetScheduledTask)
				tasksGroup.PUT("/:id", scheduledTaskHdl.UpdateScheduledTask)
				tasksGroup.DELETE("/:id", scheduledTaskHdl.DeleteScheduledTask)
				tasksGroup.POST("/:id/enable", scheduledTaskHdl.EnableScheduledTask)
				tasksGroup.POST("/:id/disable", scheduledTaskHdl.DisableScheduledTask)
				tasksGroup.GET("/:id/executions", scheduledTaskHdl.GetScheduledTaskExecutions)
			}
		}
	}

	// OIDC IdP routes (public; session used for /authorize and /login)
	if cfg.IdP.Enabled {
		store := cookie.NewStore([]byte(cfg.IdP.SessionSecret))
		oidc := r.Group("/oidc")
		oidc.Use(sessions.Sessions("idp", store))
		{
			oidc.GET("/.well-known/openid-configuration", idpHdl.Discovery)
			oidc.GET("/jwks", idpHdl.JWKS)
			oidc.GET("/authorize", idpHdl.Authorize)
			oidc.POST("/token", idpHdl.Token)
			oidc.GET("/userinfo", idpHdl.UserInfo)
			oidc.GET("/login", idpHdl.LoginGet)
			oidc.POST("/login", idpHdl.LoginPost)
		}
	}

	// WebSocket routes (protected)
	ws := r.Group("/ws")
	ws.Use(middleware.AuthMiddleware(cfg))
	{
		ws.GET("/execution-records/:id/logs", websocketHandler.StreamExecutionLogs(db))
		ws.GET("/ssh/connect", websocketHandler.SSHTerminalHandler(db, cfg, sshkeySvc, authzSvc, connectionauditSvc))
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	idp.RegisterSAMLRoutes(r, cfg)
	idp.StartLDAPListener(cfg, zapLogger)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	zapLogger.Info("Starting server", zap.String("address", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
