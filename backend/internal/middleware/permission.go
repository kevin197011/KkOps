// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/model"
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

// RequireMenuPermission automatically checks permission based on route path
// This middleware checks if the route path matches a permission in RoutePermissionMap
func RequireMenuPermission(rbacService *rbac.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the actual request path (without query string)
		path := c.Request.URL.Path

		// Find matching permission in route permission map
		// Check exact match first, then prefix match
		var requiredPermission string
		var found bool

		// Try exact match first
		if perm, exists := model.RoutePermissionMap[path]; exists {
			requiredPermission = perm
			found = true
		}

		// Try prefix match if exact match not found
		if !found {
			// Find the longest matching route (most specific match)
			var longestMatch string
			var longestMatchPerm string
			longestLen := 0
			
			for route, perm := range model.RoutePermissionMap {
				// Check if path starts with the route pattern
				// For example: /api/v1/users/123 should match /api/v1/users
				// /api/v1/dashboard/stats should match /api/v1/dashboard
				if strings.HasPrefix(path, route+"/") || path == route {
					if len(route) > longestLen {
						longestLen = len(route)
						longestMatch = route
						longestMatchPerm = perm
					}
				}
			}
			
			if longestMatch != "" {
				requiredPermission = longestMatchPerm
				found = true
			}
		}

		// Check if user is admin FIRST - admin bypasses all permission checks
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Admin bypasses all permission checks
		isAdmin, err := rbacService.IsAdmin(userID.(uint))
		if err == nil && isAdmin {
			c.Next()
			return
		}

		// If no permission required for this route, allow access
		if !found {
			c.Next()
			return
		}

		// Parse permission (format: "resource:action")
		parts := strings.Split(requiredPermission, ":")
		if len(parts) != 2 {
			c.Next()
			return
		}

		hasPermission, err := rbacService.HasPermission(userID.(uint), parts[0], parts[1])
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
