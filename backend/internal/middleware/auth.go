// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/utils"
)

// AuthMiddleware validates JWT tokens
// Supports token from:
// 1. Authorization header: "Bearer <token>"
// 2. Query parameter: ?token=<token> (for WebSocket connections)
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Try Authorization header first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// If no token in header, try query parameter (for WebSocket)
		if token == "" {
			token = c.Query("token")
		}

		// If still no token, return error
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateJWT(token, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("auth_type", "jwt")

		c.Next()
	}
}

// CombinedAuthMiddleware tries JWT first, then API token
func CombinedAuthMiddleware(cfg *config.Config, authService interface{}) gin.HandlerFunc {
	jwtMiddleware := AuthMiddleware(cfg)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Try JWT first (typically shorter, starts with "eyJ")
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			// JWT tokens typically start with "eyJ" (base64 encoded JSON header)
			if len(token) > 3 && token[:3] == "eyJ" {
				jwtMiddleware(c)
				return
			}
		}

		// If not JWT, try API token (handled by separate middleware)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
	}
}
