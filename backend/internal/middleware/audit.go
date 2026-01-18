// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/service/audit"
)

// AuditRoute 审计路由配置
type AuditRoute struct {
	PathPattern  string // 路径模式（支持正则）
	Method       string // HTTP 方法
	Module       string // 模块名称
	Action       string // 操作类型
	ResourceName string // 资源名称字段（从请求体或响应中提取）
}

// AuditMiddleware 审计中间件
func AuditMiddleware(auditSvc *audit.Service, routes []AuditRoute) gin.HandlerFunc {
	// 编译正则表达式
	compiledRoutes := make([]struct {
		Route   AuditRoute
		Pattern *regexp.Regexp
	}, 0, len(routes))

	for _, route := range routes {
		pattern, err := regexp.Compile(route.PathPattern)
		if err != nil {
			continue
		}
		compiledRoutes = append(compiledRoutes, struct {
			Route   AuditRoute
			Pattern *regexp.Regexp
		}{Route: route, Pattern: pattern})
	}

	return func(c *gin.Context) {
		// 查找匹配的路由配置
		var matchedRoute *AuditRoute
		for _, cr := range compiledRoutes {
			if cr.Route.Method == c.Request.Method && cr.Pattern.MatchString(c.Request.URL.Path) {
				matchedRoute = &cr.Route
				break
			}
		}

		// 如果没有匹配的路由，直接跳过审计
		if matchedRoute == nil {
			c.Next()
			return
		}

		// 读取请求体（需要复制一份，因为只能读取一次）
		var requestBody map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			json.Unmarshal(bodyBytes, &requestBody)
		}

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 执行请求
		c.Next()

		// 创建审计日志
		go func() {
			status := string(model.AuditStatusSuccess)
			errorMsg := ""

			if c.Writer.Status() >= 400 {
				status = string(model.AuditStatusFailed)
				// 尝试从响应中获取错误信息
				if errVal, exists := c.Get("error"); exists {
					if err, ok := errVal.(error); ok {
						errorMsg = err.Error()
					}
				}
			}

			// 提取资源名称
			resourceName := ""
			if matchedRoute.ResourceName != "" && requestBody != nil {
				if name, ok := requestBody[matchedRoute.ResourceName]; ok {
					resourceName = toString(name)
				}
			}

			// 提取资源 ID（从 URL 路径中）
			var resourceID *uint
			pathParts := strings.Split(c.Request.URL.Path, "/")
			for i, part := range pathParts {
				if part == "api" || part == "v1" || part == "" {
					continue
				}
				// 查找数字 ID
				if id := parseUint(pathParts, i+1); id > 0 {
					resourceID = &id
					break
				}
			}

			req := &audit.CreateLogRequest{
				UserID:       toUint(userID),
				Username:     toString(username),
				Action:       matchedRoute.Action,
				Module:       matchedRoute.Module,
				ResourceID:   resourceID,
				ResourceName: resourceName,
				Detail:       requestBody,
				IPAddress:    c.ClientIP(),
				UserAgent:    c.Request.UserAgent(),
				Status:       status,
				ErrorMsg:     errorMsg,
			}

			auditSvc.CreateLog(req)
		}()
	}
}

// GetDefaultAuditRoutes 获取默认审计路由配置
func GetDefaultAuditRoutes() []AuditRoute {
	return []AuditRoute{
		// 认证模块
		{PathPattern: `^/api/v1/auth/login$`, Method: "POST", Module: "auth", Action: "login", ResourceName: "username"},
		{PathPattern: `^/api/v1/auth/logout$`, Method: "POST", Module: "auth", Action: "logout"},
		{PathPattern: `^/api/v1/auth/change-password$`, Method: "POST", Module: "auth", Action: "update"},

		// 用户管理
		{PathPattern: `^/api/v1/users$`, Method: "POST", Module: "user", Action: "create", ResourceName: "username"},
		{PathPattern: `^/api/v1/users/\d+$`, Method: "PUT", Module: "user", Action: "update", ResourceName: "username"},
		{PathPattern: `^/api/v1/users/\d+$`, Method: "DELETE", Module: "user", Action: "delete"},
		{PathPattern: `^/api/v1/users/\d+/roles$`, Method: "PUT", Module: "user", Action: "update"},
		{PathPattern: `^/api/v1/users/\d+/reset-password$`, Method: "POST", Module: "user", Action: "update"},

		// 角色管理
		{PathPattern: `^/api/v1/roles$`, Method: "POST", Module: "role", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/roles/\d+$`, Method: "PUT", Module: "role", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/roles/\d+$`, Method: "DELETE", Module: "role", Action: "delete"},
		{PathPattern: `^/api/v1/roles/\d+/assets$`, Method: "PUT", Module: "role", Action: "update"},

		// 资产管理
		{PathPattern: `^/api/v1/assets$`, Method: "POST", Module: "asset", Action: "create", ResourceName: "host_name"},
		{PathPattern: `^/api/v1/assets/\d+$`, Method: "PUT", Module: "asset", Action: "update", ResourceName: "host_name"},
		{PathPattern: `^/api/v1/assets/\d+$`, Method: "DELETE", Module: "asset", Action: "delete"},

		// 项目管理
		{PathPattern: `^/api/v1/projects$`, Method: "POST", Module: "project", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/projects/\d+$`, Method: "PUT", Module: "project", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/projects/\d+$`, Method: "DELETE", Module: "project", Action: "delete"},

		// 环境管理
		{PathPattern: `^/api/v1/environments$`, Method: "POST", Module: "environment", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/environments/\d+$`, Method: "PUT", Module: "environment", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/environments/\d+$`, Method: "DELETE", Module: "environment", Action: "delete"},

		// 云平台管理
		{PathPattern: `^/api/v1/cloud-platforms$`, Method: "POST", Module: "cloud_platform", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/cloud-platforms/\d+$`, Method: "PUT", Module: "cloud_platform", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/cloud-platforms/\d+$`, Method: "DELETE", Module: "cloud_platform", Action: "delete"},

		// 任务执行
		{PathPattern: `^/api/v1/executions$`, Method: "POST", Module: "task", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/executions/\d+$`, Method: "PUT", Module: "task", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/executions/\d+$`, Method: "DELETE", Module: "task", Action: "delete"},
		{PathPattern: `^/api/v1/executions/\d+/execute$`, Method: "POST", Module: "task", Action: "execute"},
		{PathPattern: `^/api/v1/executions/\d+/cancel$`, Method: "POST", Module: "task", Action: "update"},

		// 任务模板
		{PathPattern: `^/api/v1/templates$`, Method: "POST", Module: "template", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/templates/\d+$`, Method: "PUT", Module: "template", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/templates/\d+$`, Method: "DELETE", Module: "template", Action: "delete"},

		// 定时任务
		{PathPattern: `^/api/v1/tasks$`, Method: "POST", Module: "scheduled_task", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/tasks/\d+$`, Method: "PUT", Module: "scheduled_task", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/tasks/\d+$`, Method: "DELETE", Module: "scheduled_task", Action: "delete"},
		{PathPattern: `^/api/v1/tasks/\d+/run-now$`, Method: "POST", Module: "scheduled_task", Action: "execute"},

		// 部署管理
		{PathPattern: `^/api/v1/deployment-modules$`, Method: "POST", Module: "deployment", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/deployment-modules/\d+$`, Method: "PUT", Module: "deployment", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/deployment-modules/\d+$`, Method: "DELETE", Module: "deployment", Action: "delete"},
		{PathPattern: `^/api/v1/deployment-modules/\d+/deploy$`, Method: "POST", Module: "deployment", Action: "execute"},

		// SSH 密钥
		{PathPattern: `^/api/v1/ssh-keys$`, Method: "POST", Module: "ssh_key", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/ssh-keys/\d+$`, Method: "PUT", Module: "ssh_key", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/ssh-keys/\d+$`, Method: "DELETE", Module: "ssh_key", Action: "delete"},

		// 标签管理
		{PathPattern: `^/api/v1/tags$`, Method: "POST", Module: "tag", Action: "create", ResourceName: "name"},
		{PathPattern: `^/api/v1/tags/\d+$`, Method: "PUT", Module: "tag", Action: "update", ResourceName: "name"},
		{PathPattern: `^/api/v1/tags/\d+$`, Method: "DELETE", Module: "tag", Action: "delete"},
	}
}

// 辅助函数
func toUint(v interface{}) uint {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case uint:
		return val
	case int:
		return uint(val)
	case int64:
		return uint(val)
	case float64:
		return uint(val)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func parseUint(parts []string, index int) uint {
	if index >= len(parts) {
		return 0
	}
	var id uint
	for _, ch := range parts[index] {
		if ch >= '0' && ch <= '9' {
			id = id*10 + uint(ch-'0')
		} else {
			return 0
		}
	}
	return id
}
