// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package websocket

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// LogMessage represents a log message sent via WebSocket
type LogMessage struct {
	Level     string    `json:"level"` // INFO, WARN, ERROR
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// StreamExecutionLogs creates a handler function for streaming execution logs via WebSocket
// WS /ws/task-executions/:id/logs
func StreamExecutionLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		executionIDStr := c.Param("id")
		executionID, err := strconv.ParseUint(executionIDStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid execution ID"})
			return
		}

		// Verify execution exists and user has access
		var execution model.TaskExecution
		if err := db.Preload("Task").First(&execution, executionID).Error; err != nil {
			c.JSON(404, gin.H{"error": "execution not found"})
			return
		}

		// Check if user has permission (basic check - can be enhanced with RBAC)
		if execution.Task.CreatedBy != userID.(uint) {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// Send existing output if execution is completed
		if execution.Status != "running" && execution.Output != "" {
			msg := LogMessage{
				Level:     "INFO",
				Content:   execution.Output,
				Timestamp: time.Now(),
			}
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}

		// Poll for updates if execution is running
		if execution.Status == "running" {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// Check for updates
					var updatedExecution model.TaskExecution
					if err := db.First(&updatedExecution, executionID).Error; err != nil {
						return
					}

					// Send updated output if available
					if updatedExecution.Output != execution.Output {
						newOutput := updatedExecution.Output[len(execution.Output):]
						msg := LogMessage{
							Level:     "INFO",
							Content:   newOutput,
							Timestamp: time.Now(),
						}
						if err := conn.WriteJSON(msg); err != nil {
							return
						}
						execution.Output = updatedExecution.Output
					}

					// Close connection if execution completed
					if updatedExecution.Status != "running" {
						finalMsg := LogMessage{
							Level:     "INFO",
							Content:   "\n[Execution completed]",
							Timestamp: time.Now(),
						}
						conn.WriteJSON(finalMsg)
						return
					}
				}
			}
		}

		// For completed executions, just close after sending output
	}
}
