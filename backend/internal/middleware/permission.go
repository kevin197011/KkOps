// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/rbac"
)

// RequirePermission creates a middleware that requires a specific permission
func RequirePermission(rbacService *rbac.Service, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		hasPermission, err := rbacService.HasPermission(userID.(uint), resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissionByCode is deprecated. Use RequirePermission(resource, action) instead.
// This middleware is kept for backward compatibility but will always deny access.
// Permission codes are no longer used in the system.
func RequirePermissionByCode(rbacService *rbac.Service, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permission code is no longer used. Always deny access.
		// Use RequirePermission(resource, action) instead.
		c.JSON(http.StatusForbidden, gin.H{"error": "permission code-based checks are deprecated. Use resource/action-based checks instead"})
		c.Abort()
	}
}
