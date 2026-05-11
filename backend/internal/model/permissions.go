// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

// MenuPermission 菜单权限定义
type MenuPermission struct {
	Resource    string
	Action      string
	Name        string
	Description string
}

// AllMenuPermissions 所有菜单权限定义
var AllMenuPermissions = []MenuPermission{
	// 仪表板
	{
		Resource:    "dashboard",
		Action:      "read",
		Name:        "仪表板查看",
		Description: "查看仪表板信息",
	},
	{
		Resource:    "operation-tools",
		Action:      "read",
		Name:        "运维工具查看",
		Description: "查看运维导航工具列表",
	},
	{
		Resource:    "operation-tools",
		Action:      "*",
		Name:        "运维工具管理",
		Description: "运维工具管理所有操作（查看、创建、编辑、删除）",
	},
	// 基础设施
	{
		Resource:    "projects",
		Action:      "*",
		Name:        "项目管理",
		Description: "项目管理所有操作（查看、创建、编辑、删除）",
	},
	{
		Resource:    "environments",
		Action:      "*",
		Name:        "环境管理",
		Description: "环境管理所有操作（查看、创建、编辑、删除）",
	},
	{
		Resource:    "cloud-platforms",
		Action:      "*",
		Name:        "云平台管理",
		Description: "云平台管理所有操作（查看、创建、编辑、删除）",
	},
	{
		Resource:    "assets",
		Action:      "*",
		Name:        "资产管理",
		Description: "资产管理所有操作（查看、创建、编辑、删除）",
	},
	// 任务管理
	{
		Resource:    "executions",
		Action:      "*",
		Name:        "运维执行",
		Description: "运维执行所有操作（查看、创建、执行）",
	},
	{
		Resource:    "templates",
		Action:      "*",
		Name:        "任务模板",
		Description: "任务模板所有操作（查看、创建、编辑、删除）",
	},
	{
		Resource:    "tasks",
		Action:      "*",
		Name:        "任务执行",
		Description: "任务执行所有操作（查看、创建、编辑、删除、执行）",
	},
	{
		Resource:    "deployments",
		Action:      "*",
		Name:        "部署管理",
		Description: "部署管理所有操作（查看、创建、编辑、删除、部署）",
	},
	// 安全管理
	{
		Resource:    "ssh-keys",
		Action:      "*",
		Name:        "SSH 密钥管理",
		Description: "SSH 密钥管理所有操作（查看、创建、编辑、删除）",
	},
	// 系统管理
	{
		Resource:    "users",
		Action:      "*",
		Name:        "用户管理",
		Description: "用户管理所有操作（查看、创建、编辑、删除）",
	},
	{
		Resource:    "roles",
		Action:      "*",
		Name:        "角色权限管理",
		Description: "角色权限管理所有操作（查看、创建、编辑、删除）",
	},
	{
		Resource:    "audit-logs",
		Action:      "read",
		Name:        "审计日志查看",
		Description: "查看审计日志",
	},
	{
		Resource:    "connection-audit",
		Action:      "read",
		Name:        "审计连线查看",
		Description: "查看 WebSSH 连线录像记录与回放",
	},
	{
		Resource:    "external-systems",
		Action:      "*",
		Name:        "外部运维系统",
		Description: "配置与跳转外部运维系统（SSO 附带权限）",
	},
	{
		Resource:    "oauth2-clients",
		Action:      "*",
		Name:        "IdP 应用管理",
		Description: "IdP 应用（OAuth2/OIDC/SAML/LDAP 协议登记）的创建、编辑、删除与密钥管理",
	},
	{
		Resource:    "provisioning",
		Action:      "*",
		Name:        "账号同步",
		Description: "账号下发目标与同步记录管理",
	},
	{
		Resource:    "integrations",
		Action:      "*",
		Name:        "外部集成",
		Description: "第三方集成连接器的创建、编辑、删除与加密凭据管理",
	},
	{
		Resource:    "monitoring",
		Action:      "*",
		Name:        "监控集成",
		Description: "通过集成查询 Prometheus / Nightingale 指标",
	},
	{
		Resource:    "logging",
		Action:      "*",
		Name:        "日志集成",
		Description: "通过集成检索 Loki / Elasticsearch 日志",
	},
	{
		Resource:    "cicd",
		Action:      "*",
		Name:        "CI/CD 集成",
		Description: "Jenkins / GitLab CI 流水线与日志代理",
	},
	{
		Resource:    "registry",
		Action:      "*",
		Name:        "镜像仓库",
		Description: "Harbor 镜像仓库浏览",
	},
	{
		Resource:    "gitops",
		Action:      "*",
		Name:        "GitOps 集成",
		Description: "Argo CD 应用查看与同步",
	},
	{
		Resource:    "kubernetes",
		Action:      "*",
		Name:        "Kubernetes 集群",
		Description: "多集群 Kubernetes 资源浏览与日志",
	},
	{
		Resource:    "alerts",
		Action:      "*",
		Name:        "告警中心",
		Description: "聚合告警查看、确认与同步",
	},
	{
		Resource:    "incidents",
		Action:      "*",
		Name:        "事件管理",
		Description: "运维事件工单创建与跟踪",
	},
	{
		Resource:    "ai",
		Action:      "*",
		Name:        "AI 运维",
		Description: "AI 助手、异常检测与根因分析",
	},
	{
		Resource:    "cmdb",
		Action:      "*",
		Name:        "CMDB 与拓扑",
		Description: "配置项（CMDB）与依赖拓扑查看与管理",
	},
}

// RoutePermissionMap 路由与权限映射（用于后端权限检查）
var RoutePermissionMap = map[string]string{
	// 仪表板
	"/api/v1/dashboard":       "dashboard:read",
	"/api/v1/operation-tools": "operation-tools:read",
	// 基础设施
	"/api/v1/projects":         "projects:*",
	"/api/v1/environments":     "environments:*",
	"/api/v1/cloud-platforms":  "cloud-platforms:*",
	"/api/v1/assets":           "assets:*",
	"/api/v1/asset-categories": "assets:*",
	"/api/v1/tags":             "assets:*",
	// 任务管理
	"/api/v1/executions":         "executions:*",
	"/api/v1/execution-records":  "executions:*",
	"/api/v1/templates":          "templates:*",
	"/api/v1/tasks":              "tasks:*",
	"/api/v1/deployment-modules": "deployments:*",
	"/api/v1/deployments":        "deployments:*",
	// 安全管理
	"/api/v1/ssh/keys": "ssh-keys:*",
	// 系统管理（仅管理员）
	"/api/v1/users":            "users:*",
	"/api/v1/roles":            "roles:*",
	"/api/v1/permissions":      "roles:*", // 权限管理属于角色管理的一部分
	"/api/v1/audit-logs":       "audit-logs:read",
	"/api/v1/connection-audit": "connection-audit:read",
	"/api/v1/external-systems": "external-systems:*",
	"/api/v1/oauth2-clients":   "oauth2-clients:*",
	"/api/v1/provisioning":     "provisioning:*",
	"/api/v1/integrations":     "integrations:*",
	"/api/v1/monitoring":       "monitoring:*",
	"/api/v1/logging":          "logging:*",
	"/api/v1/cicd":             "cicd:*",
	"/api/v1/registry":         "registry:*",
	"/api/v1/gitops":           "gitops:*",
	"/api/v1/k8s":              "kubernetes:*",
	"/api/v1/alerts":           "alerts:*",
	"/api/v1/incidents":        "incidents:*",
	"/api/v1/ai":               "ai:*",
	"/api/v1/cmdb":             "cmdb:*",
	"/api/v1/topology":         "cmdb:*",
}
