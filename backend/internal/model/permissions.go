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
}

// RoutePermissionMap 路由与权限映射（用于后端权限检查）
var RoutePermissionMap = map[string]string{
	// 仪表板
	"/api/v1/dashboard": "dashboard:read",
	"/api/v1/operation-tools": "operation-tools:read",
	// 基础设施
	"/api/v1/projects":         "projects:*",
	"/api/v1/environments":     "environments:*",
	"/api/v1/cloud-platforms":  "cloud-platforms:*",
	"/api/v1/assets":           "assets:*",
	// 任务管理
	"/api/v1/executions":           "executions:*",
	"/api/v1/execution-records":    "executions:*",
	"/api/v1/templates":            "templates:*",
	"/api/v1/tasks":                "tasks:*",
	"/api/v1/deployment-modules":   "deployments:*",
	"/api/v1/deployments":          "deployments:*",
	// 安全管理
	"/api/v1/ssh/keys": "ssh-keys:*",
	// 系统管理（仅管理员）
	"/api/v1/users":      "users:*",
	"/api/v1/roles":      "roles:*",
	"/api/v1/permissions": "roles:*", // 权限管理属于角色管理的一部分
	"/api/v1/audit-logs": "audit-logs:read",
}
