// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package audit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/audit"
)

// Handler 审计日志处理器
type Handler struct {
	svc *audit.Service
}

// NewHandler 创建审计处理器实例
func NewHandler(svc *audit.Service) *Handler {
	return &Handler{svc: svc}
}

// ListAuditLogs 获取审计日志列表
// @Summary 获取审计日志列表
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param user_id query int false "用户ID"
// @Param username query string false "用户名"
// @Param module query string false "模块"
// @Param action query string false "操作类型"
// @Param status query string false "状态"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param keyword query string false "关键词"
// @Success 200 {object} map[string]interface{}
// @Router /audit-logs [get]
func (h *Handler) ListAuditLogs(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建查询请求
	req := &audit.ListRequest{
		Page:     page,
		PageSize: pageSize,
		Username: c.Query("username"),
		Module:   c.Query("module"),
		Action:   c.Query("action"),
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
	}

	// 解析用户 ID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(userID)
			req.UserID = &uid
		}
	}

	// 解析时间范围
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			req.StartTime = &t
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			req.EndTime = &t
		}
	}

	// 查询
	result, err := h.svc.ListLogs(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询审计日志失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total": result.Total,
			"items": result.Items,
			"page":  page,
			"page_size": pageSize,
		},
	})
}

// GetAuditLog 获取单条审计日志
// @Summary 获取审计日志详情
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} map[string]interface{}
// @Router /audit-logs/{id} [get]
func (h *Handler) GetAuditLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"})
		return
	}

	log, err := h.svc.GetLog(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "审计日志不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": log})
}

// ExportAuditLogs 导出审计日志
// @Summary 导出审计日志
// @Tags 审计日志
// @Accept json
// @Produce octet-stream
// @Param format query string false "导出格式" Enums(csv, json) default(csv)
// @Param user_id query int false "用户ID"
// @Param username query string false "用户名"
// @Param module query string false "模块"
// @Param action query string false "操作类型"
// @Param status query string false "状态"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param keyword query string false "关键词"
// @Success 200 {file} binary
// @Router /audit-logs/export [get]
func (h *Handler) ExportAuditLogs(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	if format != "csv" && format != "json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的导出格式，请使用 csv 或 json"})
		return
	}

	// 构建查询请求
	req := &audit.ListRequest{
		Username: c.Query("username"),
		Module:   c.Query("module"),
		Action:   c.Query("action"),
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
	}

	// 解析用户 ID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(userID)
			req.UserID = &uid
		}
	}

	// 解析时间范围
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			req.StartTime = &t
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			req.EndTime = &t
		}
	}

	// 设置响应头
	filename := fmt.Sprintf("audit_logs_%s.%s", time.Now().Format("20060102_150405"), format)
	contentType := "text/csv"
	if format == "json" {
		contentType = "application/json"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 导出
	if err := h.svc.ExportLogs(req, format, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导出失败: " + err.Error()})
		return
	}
}

// GetModules 获取所有模块列表
// @Summary 获取模块列表
// @Tags 审计日志
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /audit-logs/modules [get]
func (h *Handler) GetModules(c *gin.Context) {
	modules := h.svc.GetModules()
	c.JSON(http.StatusOK, gin.H{"data": modules})
}

// GetActions 获取所有操作类型列表
// @Summary 获取操作类型列表
// @Tags 审计日志
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /audit-logs/actions [get]
func (h *Handler) GetActions(c *gin.Context) {
	actions := h.svc.GetActions()
	c.JSON(http.StatusOK, gin.H{"data": actions})
}
