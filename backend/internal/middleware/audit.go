package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service"
)

// AuditMiddleware 审计日志中间件
func AuditMiddleware(auditService service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过健康检查等不需要审计的路径
		if c.Request.URL.Path == "/health" || 
		   c.Request.URL.Path == "/api/v1/auth/login" ||
		   c.Request.URL.Path == "/api/v1/auth/register" {
			c.Next()
			return
		}

		// 如果 auditService 未初始化，跳过审计
		if auditService == nil {
			c.Next()
			return
		}

		startTime := time.Now()

		// 获取用户信息（从认证中间件设置）
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		var userIDPtr *uint64
		var usernameStr string
		if userID != nil {
			// 处理不同的类型可能（uint64, uint, float64等）
			switch uid := userID.(type) {
			case uint64:
				userIDPtr = &uid
			case uint:
				uid64 := uint64(uid)
				userIDPtr = &uid64
			case float64:
				uid64 := uint64(uid)
				userIDPtr = &uid64
			case int64:
				uid64 := uint64(uid)
				userIDPtr = &uid64
			case int:
				uid64 := uint64(uid)
				userIDPtr = &uid64
			}
		}
		if username != nil {
			if uname, ok := username.(string); ok {
				usernameStr = uname
			}
		}

		// 读取请求体（用于记录请求数据）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应写入器包装器，用于捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 执行请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)
		durationMs := int(duration.Milliseconds())

		// 确定操作类型和资源信息
		action := getActionFromMethod(c.Request.Method)
		resourceType := getResourceTypeFromPath(c.Request.URL.Path)
		resourceID := extractResourceIDFromPath(c.Request.URL.Path)
		resourceName := extractResourceNameFromResponse(blw.body.String(), resourceType, resourceID)

		// 确定状态
		status := "success"
		if c.Writer.Status() >= 400 {
			status = "failed"
		}

		// 获取错误信息
		var errorMsg string
		if status == "failed" {
			errorMsg = blw.body.String()
			if len(errorMsg) > 500 {
				errorMsg = errorMsg[:500] + "..."
			}
		}

		// 限制响应体大小（避免记录过大的响应）
		responseBody := blw.body.String()
		if len(responseBody) > 10000 {
			responseBody = responseBody[:10000] + "... (truncated)"
		}

		// 限制请求体大小
		requestBodyStr := string(requestBody)
		if len(requestBodyStr) > 10000 {
			requestBodyStr = requestBodyStr[:10000] + "... (truncated)"
		}

		// 异步记录审计日志（不阻塞请求）
		go func() {
			_ = auditService.LogOperation(
				userIDPtr,
				usernameStr,
				action,
				resourceType,
				resourceID,
				resourceName,
				c.ClientIP(),
				c.Request.UserAgent(),
				requestBodyStr,
				responseBody,
				nil, // beforeData
				nil, // afterData
				status,
				errorMsg,
				durationMs,
			)
		}()
	}
}

// bodyLogWriter 响应写入器包装器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// getActionFromMethod 根据HTTP方法确定操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "query"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "execute"
	}
}

// getResourceTypeFromPath 从路径提取资源类型
func getResourceTypeFromPath(path string) string {
	// 从路径中提取资源类型，例如 /api/v1/hosts -> host
	if len(path) > 10 && path[:10] == "/api/v1/" {
		pathSuffix := path[10:]
		
		// 特殊处理一些路径
		specialPaths := map[string]string{
			"settings":        "settings",
			"settings/salt":   "salt_config",
			"settings/salt/":  "salt_config",
			"auth/login":      "auth",
			"auth/logout":     "auth",
			"auth/register":   "auth",
			"deployment-configs": "deployment_config",
			"deployment-versions": "deployment_version",
			"host-groups":     "host_group",
			"host-tags":       "host_tag",
			"ssh-keys":        "ssh_key",
		}
		
		// 检查是否是特殊路径
		for specialPath, resourceType := range specialPaths {
			if len(pathSuffix) >= len(specialPath) && pathSuffix[:len(specialPath)] == specialPath {
				if len(pathSuffix) == len(specialPath) || pathSuffix[len(specialPath)] == '/' {
					return resourceType
				}
			}
		}
		
		// 普通路径处理：提取第一个路径段
		var resourceType string
		for i := 0; i < len(pathSuffix); i++ {
			if pathSuffix[i] == '/' || pathSuffix[i] == '?' {
				break
			}
			resourceType += string(pathSuffix[i])
		}
		
		// 移除复数形式（但保留特殊情况）
		if len(resourceType) > 0 && resourceType[len(resourceType)-1] == 's' {
			// 对于以 's' 结尾但不是特殊路径的，移除 's'
			if _, isSpecial := specialPaths[resourceType]; !isSpecial {
				resourceType = resourceType[:len(resourceType)-1]
			}
		}
		
		return resourceType
	}
	return ""
}

// extractResourceIDFromPath 从路径提取资源ID
func extractResourceIDFromPath(path string) *uint64 {
	// 从路径中提取资源ID，例如 /api/v1/hosts/123 -> 123
	parts := []rune(path)
	var idStr string
	var foundSlash bool
	
	// 从后往前查找数字
	for i := len(parts) - 1; i >= 0; i-- {
		char := parts[i]
		if char == '?' {
			// 遇到查询参数，停止
			break
		}
		if char >= '0' && char <= '9' {
			idStr = string(char) + idStr
			foundSlash = false
		} else if char == '/' && idStr != "" {
			foundSlash = true
			break
		} else if char != '/' && idStr != "" {
			// 遇到非数字非斜杠字符，重置
			idStr = ""
		}
	}
	
	if idStr != "" && foundSlash {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			return &id
		}
	}
	
	return nil
}

// extractResourceNameFromResponse 从响应中提取资源名称
func extractResourceNameFromResponse(responseBody string, resourceType string, resourceID *uint64) string {
	if responseBody == "" || resourceType == "" {
		return ""
	}

	// 尝试从响应JSON中提取资源名称
	// 根据不同的资源类型，名称字段可能不同
	var responseMap map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &responseMap); err != nil {
		return ""
	}

	// 尝试常见的名称字段
	if name, ok := responseMap["name"].(string); ok && name != "" {
		return name
	}
	if username, ok := responseMap["username"].(string); ok && username != "" && resourceType == "user" {
		return username
	}
	if hostname, ok := responseMap["hostname"].(string); ok && hostname != "" && resourceType == "host" {
		return hostname
	}

	// 如果响应中有对应的资源对象，尝试从中提取名称
	if resource, ok := responseMap[resourceType].(map[string]interface{}); ok {
		if name, ok := resource["name"].(string); ok && name != "" {
			return name
		}
		if username, ok := resource["username"].(string); ok && username != "" {
			return username
		}
		if hostname, ok := resource["hostname"].(string); ok && hostname != "" {
			return hostname
		}
	}

	// 如果响应是一个数组（列表查询），返回空字符串
	if _, ok := responseMap[resourceType+"s"]; ok {
		return ""
	}

	return ""
}

