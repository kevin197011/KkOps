package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
)

// RequirePermission 权限检查中间件
func RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// 获取用户权限
		userRepo := repository.NewUserRepository(models.DB)
		permissions, err := userRepo.GetUserPermissions(userID.(uint64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user permissions"})
			c.Abort()
			return
		}

		// 检查是否有指定权限
		hasPermission := false
		for _, perm := range permissions {
			if perm.Code == permissionCode {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 检查是否有任一权限
func RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// 获取用户权限
		userRepo := repository.NewUserRepository(models.DB)
		permissions, err := userRepo.GetUserPermissions(userID.(uint64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user permissions"})
			c.Abort()
			return
		}

		// 检查是否有任一权限
		hasPermission := false
		userPermCodes := make(map[string]bool)
		for _, perm := range permissions {
			userPermCodes[perm.Code] = true
		}

		for _, code := range permissionCodes {
			if userPermCodes[code] {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 角色检查中间件
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		roleList := roles.([]string)
		hasRole := false
		for _, role := range roleList {
			if role == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			c.Abort()
			return
		}

		c.Next()
	}
}

