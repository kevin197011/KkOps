package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service"
)

type AuditHandler struct {
	auditService service.AuditService
}

func NewAuditHandler(auditService service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// ListAuditLogs 查询审计日志列表
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})

	// 用户ID过滤
	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 64); err == nil {
			filters["user_id"] = id
		}
	}

	// 用户名过滤
	if username := c.Query("username"); username != "" {
		filters["username"] = username
	}

	// 操作类型过滤
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}

	// 资源类型过滤
	if resourceType := c.Query("resource_type"); resourceType != "" {
		filters["resource_type"] = resourceType
	}

	// 资源ID过滤
	if resourceID := c.Query("resource_id"); resourceID != "" {
		if id, err := strconv.ParseUint(resourceID, 10, 64); err == nil {
			filters["resource_id"] = id
		}
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// 时间范围过滤
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", startTimeStr); err == nil {
			filters["start_time"] = startTime
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse("2006-01-02T15:04:05Z07:00", endTimeStr); err == nil {
			filters["end_time"] = endTime
		}
	}

	logs, total, err := h.auditService.ListLogs(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAuditLog 获取审计日志详情
func (h *AuditHandler) GetAuditLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log id"})
		return
	}

	log, err := h.auditService.GetLog(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audit log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": log})
}

// GetLogsByResource 获取资源的审计日志
func (h *AuditHandler) GetLogsByResource(c *gin.Context) {
	resourceType := c.Query("resource_type")
	if resourceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource_type is required"})
		return
	}

	resourceID, err := strconv.ParseUint(c.Query("resource_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource_id"})
		return
	}

	logs, err := h.auditService.GetLogsByResource(resourceType, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetLogsByUser 获取用户的审计日志
func (h *AuditHandler) GetLogsByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	logs, err := h.auditService.GetLogsByUser(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

