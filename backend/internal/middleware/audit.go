package middleware

import (
	"bytes"
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
			if uid, ok := userID.(uint64); ok {
				userIDPtr = &uid
			}
		}
		if username != nil {
			usernameStr = username.(string)
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
		resourceName := "" // 资源名称需要从响应中提取，这里先留空

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
		parts := []rune(path[10:])
		var resourceType string
		for _, char := range parts {
			if char == '/' || char == '?' {
				break
			}
			resourceType += string(char)
		}
		// 移除复数形式
		if len(resourceType) > 0 && resourceType[len(resourceType)-1] == 's' {
			resourceType = resourceType[:len(resourceType)-1]
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

