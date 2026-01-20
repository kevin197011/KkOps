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
// @Summary Create template
// @Description Create a new task template
// @Tags templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body task.CreateTemplateRequest true "Create template request"
// @Success 201 {object} task.TemplateResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/templates [post]
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
// @Summary Get template
// @Description Get task template by ID
// @Tags templates
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Success 200 {object} task.TemplateResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/templates/{id} [get]
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
// @Summary List templates
// @Description Get list of all task templates
// @Tags templates
// @Produce json
// @Security BearerAuth
// @Success 200 {array} task.TemplateResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/templates [get]
func (h *Handler) ListTemplates(c *gin.Context) {
	templates, err := h.service.ListTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// UpdateTemplate handles template update
// @Summary Update template
// @Description Update task template information
// @Tags templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Param request body task.UpdateTemplateRequest true "Update template request"
// @Success 200 {object} task.TemplateResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/templates/{id} [put]
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
// @Summary Delete template
// @Description Delete a task template by ID
// @Tags templates
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/templates/{id} [delete]
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
// @Summary Create task
// @Description Create a new execution task
// @Tags executions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body task.CreateTaskRequest true "Create task request"
// @Success 201 {object} task.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/executions [post]
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
// @Summary Get task
// @Description Get execution task by ID
// @Tags executions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} task.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/executions/{id} [get]
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
// @Summary List tasks
// @Description Get paginated list of execution tasks
// @Tags executions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{} "Response with data, total, page, and size"
// @Failure 500 {object} map[string]string
// @Router /api/v1/executions [get]
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
// @Summary Update task
// @Description Update execution task information
// @Tags executions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param request body task.UpdateTaskRequest true "Update task request"
// @Success 200 {object} task.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/executions/{id} [put]
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
// @Summary Delete task
// @Description Delete an execution task by ID
// @Tags executions
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/executions/{id} [delete]
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
// @Summary Execute task
// @Description Execute a task (sync or async)
// @Tags executions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param request body object true "Execution request" SchemaExample({"execution_type": "sync"})
// @Success 200 {object} map[string]string "Response with message"
// @Failure 400 {object} map[string]string
// @Router /api/v1/executions/{id}/execute [post]
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
// @Summary Get task execution history
// @Description Get execution history for a task
// @Tags executions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{} "Response with data array"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/executions/{id}/history [get]
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
// @Summary Get execution record
// @Description Get execution record by ID
// @Tags execution-records
// @Produce json
// @Security BearerAuth
// @Param id path int true "Execution Record ID"
// @Success 200 {object} map[string]interface{} "Response with data object"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/execution-records/{id} [get]
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
// @Summary Cancel execution record
// @Description Cancel a running execution record
// @Tags execution-records
// @Security BearerAuth
// @Param id path int true "Execution Record ID"
// @Success 200 {object} map[string]string "Response with message"
// @Failure 400 {object} map[string]string
// @Router /api/v1/execution-records/{id}/cancel [post]
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
// @Summary Cancel task
// @Description Cancel all executions for a task
// @Tags executions
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string "Response with message"
// @Failure 400 {object} map[string]string
// @Router /api/v1/executions/{id}/cancel [post]
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
// @Summary Get execution logs
// @Description Get logs for an execution record
// @Tags execution-records
// @Produce json
// @Security BearerAuth
// @Param id path int true "Execution Record ID"
// @Success 200 {object} map[string]interface{} "Response with data.logs array"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/execution-records/{id}/logs [get]
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
// @Tags templates
// @Produce json
// @Security BearerAuth
// @Success 200 {object} task.ExportTemplatesConfig
// @Failure 500 {object} map[string]string
// @Router /api/v1/templates/export [get]
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
// @Tags templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body task.ImportTemplatesConfig true "Import config"
// @Success 200 {object} task.ImportResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/templates/import [post]
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
