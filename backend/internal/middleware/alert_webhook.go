// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/config"
)

const alertWebhookSecretHeader = "X-KkOps-Alert-Secret"

// AlertWebhookSecret validates the configured secret header for public alert ingestion.
func AlertWebhookSecret(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := cfg.Alerts.WebhookSecret
		if secret == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "alert webhook is not configured"})
			c.Abort()
			return
		}
		if c.GetHeader(alertWebhookSecretHeader) != secret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid webhook secret"})
			c.Abort()
			return
		}
		c.Next()
	}
}
