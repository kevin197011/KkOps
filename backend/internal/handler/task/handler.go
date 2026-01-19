// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package task

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/task"
)

// Handler handles task and template management HTTP requests
type Handler struct {
	service          *task.Service
	executionService *task.ExecutionService
}

// NewHandler creates a new task handler
func NewHandler(service *task.Service, executionService *task.ExecutionService) *Handler {
	return &Handler{
		service:          service,
		executionService: executionService,
	}
}

// CreateTemplate handles template creation
func (h *Handler) CreateTemplate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req task.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateTemplate(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetTemplate handles template retrieval
func (h *Handler) GetTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	resp, err := h.service.GetTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListTemplates handles template list retrieval
func (h *Handler) ListTemplates(c *gin.Context) {
	templates, err := h.service.ListTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// UpdateTemplate handles template update
func (h *Handler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	var req task.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateTemplate(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTemplate handles template deletion
func (h *Handler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	if err := h.service.DeleteTemplate(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateTask handles task creation
func (h *Handler) CreateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req task.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateTask(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetTask handles task retrieval
func (h *Handler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	resp, err := h.service.GetTask(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListTasks handles task list retrieval
func (h *Handler) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tasks, total, err := h.service.ListTasks(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// UpdateTask handles task update
func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req task.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateTask(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTask handles task deletion
func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.service.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExecuteTask handles task execution
// POST /api/v1/tasks/:id/execute
func (h *Handler) ExecuteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req struct {
		ExecutionType string `json:"execution_type"` // sync or async
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.ExecutionType = "sync" // Default to sync
	}

	executionType := req.ExecutionType
	if executionType == "" {
		executionType = "sync"
	}
	if executionType != "sync" && executionType != "async" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "execution_type must be 'sync' or 'async'"})
		return
	}

	if err := h.executionService.ExecuteTask(uint(id), executionType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task execution started"})
}

// GetTaskExecutions handles getting all executions for a task
// GET /api/v1/tasks/:id/executions
func (h *Handler) GetTaskExecutions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	executions, err := h.executionService.GetTaskExecutions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": executions})
}

// GetTaskExecution handles getting a single execution
// GET /api/v1/task-executions/:id
func (h *Handler) GetTaskExecution(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	execution, err := h.executionService.GetTaskExecution(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": execution})
}

// CancelTaskExecution handles cancelling a task execution
// POST /api/v1/task-executions/:id/cancel
func (h *Handler) CancelTaskExecution(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	if err := h.executionService.CancelTaskExecution(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "execution cancelled"})
}

// CancelTask handles cancelling all executions for a task
// POST /api/v1/tasks/:id/cancel
func (h *Handler) CancelTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.executionService.CancelTask(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task cancelled"})
}

// GetTaskExecutionLogs handles getting logs for a task execution
// GET /api/v1/task-executions/:id/logs
func (h *Handler) GetTaskExecutionLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	execution, err := h.executionService.GetTaskExecution(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	// Return logs from the execution record
	logs := []string{}
	if execution.Output != "" {
		logs = append(logs, execution.Output)
	}
	if execution.Error != "" {
		logs = append(logs, "Error: "+execution.Error)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"logs": logs}})
}

// ExportTemplates handles template export
// @Summary Export task templates
// @Description Export all task templates to JSON
// @Tags Templates
// @Produce json
// @Success 200 {object} task.ExportTemplatesConfig
// @Router /templates/export [get]
func (h *Handler) ExportTemplates(c *gin.Context) {
	config, err := h.service.ExportTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set response headers for file download
	filename := fmt.Sprintf("task-templates-%s.json", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, config)
}

// ImportTemplates handles template import
// @Summary Import task templates
// @Description Import task templates from JSON
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body task.ImportTemplatesConfig true "Import config"
// @Success 200 {object} task.ImportResult
// @Failure 400 {object} map[string]string
// @Router /templates/import [post]
func (h *Handler) ImportTemplates(c *gin.Context) {
	var config task.ImportTemplatesConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的导入配置: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	result, err := h.service.ImportTemplates(&config, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
