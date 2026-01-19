// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package scheduledtask

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/scheduledtask"
)

// Handler 定时任务处理器
type Handler struct {
	service *scheduledtask.Service
}

// NewHandler 创建定时任务处理器
func NewHandler(service *scheduledtask.Service) *Handler {
	return &Handler{service: service}
}

// CreateScheduledTask godoc
// @Summary 创建定时任务
// @Description 创建一个新的定时任务
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body scheduledtask.CreateScheduledTaskRequest true "定时任务信息"
// @Success 200 {object} scheduledtask.ScheduledTaskResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *Handler) CreateScheduledTask(c *gin.Context) {
	var req scheduledtask.CreateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	task, err := h.service.CreateScheduledTask(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetScheduledTask godoc
// @Summary 获取定时任务详情
// @Description 根据 ID 获取定时任务详情
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务 ID"
// @Success 200 {object} scheduledtask.ScheduledTaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *Handler) GetScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return
	}

	task, err := h.service.GetScheduledTask(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ListScheduledTasks godoc
// @Summary 获取定时任务列表
// @Description 获取定时任务列表，支持分页和筛选
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param enabled query bool false "是否启用"
// @Success 200 {object} scheduledtask.ListScheduledTasksResponse
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *Handler) ListScheduledTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var enabled *bool
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		e := enabledStr == "true"
		enabled = &e
	}

	result, err := h.service.ListScheduledTasks(page, pageSize, enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateScheduledTask godoc
// @Summary 更新定时任务
// @Description 更新定时任务信息
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务 ID"
// @Param request body scheduledtask.UpdateScheduledTaskRequest true "更新信息"
// @Success 200 {object} scheduledtask.ScheduledTaskResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func (h *Handler) UpdateScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return
	}

	var req scheduledtask.UpdateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.UpdateScheduledTask(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteScheduledTask godoc
// @Summary 删除定时任务
// @Description 删除定时任务
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务 ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return
	}

	if err := h.service.DeleteScheduledTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// EnableScheduledTask godoc
// @Summary 启用定时任务
// @Description 启用定时任务
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务 ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/enable [post]
func (h *Handler) EnableScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return
	}

	if err := h.service.EnableScheduledTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已启用"})
}

// DisableScheduledTask godoc
// @Summary 禁用定时任务
// @Description 禁用定时任务
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务 ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/disable [post]
func (h *Handler) DisableScheduledTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return
	}

	if err := h.service.DisableScheduledTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已禁用"})
}

// GetScheduledTaskExecutions godoc
// @Summary 获取定时任务执行历史
// @Description 获取定时任务的执行历史记录
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务 ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/executions [get]
func (h *Handler) GetScheduledTaskExecutions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	executions, total, err := h.service.GetScheduledTaskExecutions(uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  executions,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// ValidateCron godoc
// @Summary 验证 Cron 表达式
// @Description 验证 Cron 表达式是否有效，并返回下次执行时间
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cron_expression query string true "Cron 表达式"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /tasks/validate-cron [get]
func (h *Handler) ValidateCron(c *gin.Context) {
	cronExpr := c.Query("cron_expression")
	if cronExpr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 cron_expression 参数"})
		return
	}

	if err := scheduledtask.ValidateCronExpression(cronExpr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "valid": false})
		return
	}

	nextRunAt, _ := scheduledtask.GetNextRunTime(cronExpr)
	c.JSON(http.StatusOK, gin.H{
		"valid":       true,
		"next_run_at": nextRunAt,
	})
}

// ExportScheduledTasks handles scheduled tasks export
// @Summary Export scheduled tasks
// @Description Export all scheduled tasks to JSON
// @Tags Scheduled Tasks
// @Produce json
// @Success 200 {object} scheduledtask.ExportScheduledTasksConfig
// @Router /tasks/export [get]
func (h *Handler) ExportScheduledTasks(c *gin.Context) {
	config, err := h.service.ExportScheduledTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set response headers for file download
	filename := fmt.Sprintf("scheduled-tasks-%s.json", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, config)
}

// ImportScheduledTasks handles scheduled tasks import
// @Summary Import scheduled tasks
// @Description Import scheduled tasks from JSON
// @Tags Scheduled Tasks
// @Accept json
// @Produce json
// @Param request body scheduledtask.ImportScheduledTasksConfig true "Import config"
// @Success 200 {object} scheduledtask.ImportScheduledTasksResult
// @Failure 400 {object} map[string]string
// @Router /tasks/import [post]
func (h *Handler) ImportScheduledTasks(c *gin.Context) {
	var config scheduledtask.ImportScheduledTasksConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的导入配置: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	result, err := h.service.ImportScheduledTasks(&config, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
