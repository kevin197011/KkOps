package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/service"
	"gorm.io/datatypes"
)

type BatchOperationHandler struct {
	operationService service.BatchOperationService
}

func NewBatchOperationHandler(operationService service.BatchOperationService) *BatchOperationHandler {
	return &BatchOperationHandler{
		operationService: operationService,
	}
}

// CreateOperation 创建并执行批量操作
func (h *BatchOperationHandler) CreateOperation(c *gin.Context) {
	var req struct {
		Name            string                   `json:"name" binding:"required"`
		Description     string                   `json:"description"`
		CommandType     string                   `json:"command_type" binding:"required"`
		CommandFunction string                   `json:"command_function" binding:"required"`
		CommandArgs     []interface{}            `json:"command_args"`
		TargetHosts     []map[string]interface{} `json:"target_hosts" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 转换目标主机为 JSON
	targetHostsJSON, err := json.Marshal(req.TargetHosts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target_hosts format"})
		return
	}

	// 转换命令参数为 JSON
	var commandArgsJSON datatypes.JSON
	if len(req.CommandArgs) > 0 {
		argsBytes, err := json.Marshal(req.CommandArgs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command_args format"})
			return
		}
		commandArgsJSON = datatypes.JSON(argsBytes)
	}

	// 创建操作
	operation := &models.BatchOperation{
		Name:            req.Name,
		Description:     req.Description,
		CommandType:     req.CommandType,
		CommandFunction: req.CommandFunction,
		CommandArgs:     commandArgsJSON,
		TargetHosts:     datatypes.JSON(targetHostsJSON),
		TargetCount:     len(req.TargetHosts),
		StartedBy:       userID.(uint64),
		StartedAt:       time.Now(),
	}

	if err := h.operationService.CreateOperation(operation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 异步执行
	if err := h.operationService.ExecuteOperation(operation.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"operation": operation})
}

// ListOperations 列出批量操作
func (h *BatchOperationHandler) ListOperations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if startedBy := c.Query("started_by"); startedBy != "" {
		if id, err := strconv.ParseUint(startedBy, 10, 32); err == nil {
			filters["started_by"] = uint(id)
		}
	}

	// 获取当前用户ID（如果用户只查看自己的操作）
	if viewOwn := c.Query("view_own"); viewOwn == "true" {
		userID, exists := c.Get("user_id")
		if exists {
			filters["started_by"] = userID.(uint64)
		}
	}

	operations, total, err := h.operationService.ListOperations(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"operations": operations,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// GetOperation 获取操作详情
func (h *BatchOperationHandler) GetOperation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation id"})
		return
	}

	operation, err := h.operationService.GetOperation(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "operation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"operation": operation})
}

// GetOperationStatus 获取操作状态
func (h *BatchOperationHandler) GetOperationStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation id"})
		return
	}

	operation, err := h.operationService.GetOperationStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "operation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     operation.ID,
		"status": operation.Status,
	})
}

// GetOperationResults 获取操作结果
func (h *BatchOperationHandler) GetOperationResults(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation id"})
		return
	}

	operation, err := h.operationService.GetOperationResults(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "operation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            operation.ID,
		"status":        operation.Status,
		"results":       operation.Results,
		"success_count": operation.SuccessCount,
		"failed_count":  operation.FailedCount,
	})
}

// CancelOperation 取消操作
func (h *BatchOperationHandler) CancelOperation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation id"})
		return
	}

	if err := h.operationService.CancelOperation(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operation cancelled"})
}

// CleanupOldOperations 清理超过1个月的操作记录
func (h *BatchOperationHandler) CleanupOldOperations(c *gin.Context) {
	// 检查用户权限（需要管理员权限）
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 这里可以添加管理员权限检查，暂时允许所有认证用户

	// 计算1个月前的时间
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	// 执行清理
	count, err := h.operationService.CleanupOldOperations(oneMonthAgo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup old operations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cleaned up %d old operations", count),
		"count":   count,
	})
}

