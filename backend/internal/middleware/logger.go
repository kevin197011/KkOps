// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin middleware for request logging
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("ip", ip),
			zap.Duration("latency", latency),
		)
	}
}
