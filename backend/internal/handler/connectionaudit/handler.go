// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package connectionaudit

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/connectionaudit"
)

// Handler 审计连线处理器
type Handler struct {
	svc *connectionaudit.Service
}

// NewHandler 创建审计连线处理器实例
func NewHandler(svc *connectionaudit.Service) *Handler {
	return &Handler{svc: svc}
}

// ListConnectionRecords 获取连线记录列表
// @Summary 获取 WebSSH 连线记录列表
// @Tags 审计连线
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param user_id query int false "用户ID"
// @Param asset_id query int false "资产ID"
// @Param start_time query string false "开始时间 RFC3339"
// @Param end_time query string false "结束时间 RFC3339"
// @Success 200 {object} connectionaudit.ListResponse
// @Router /connection-audit [get]
func (h *Handler) ListConnectionRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	req := &connectionaudit.ListRequest{
		Page:     page,
		PageSize: pageSize,
	}
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(id)
			req.UserID = &uid
		}
	}
	if assetIDStr := c.Query("asset_id"); assetIDStr != "" {
		if id, err := strconv.ParseUint(assetIDStr, 10, 32); err == nil {
			aid := uint(id)
			req.AssetID = &aid
		}
	}
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

	resp, err := h.svc.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetConnectionRecord 获取单条连线记录（含 Transcript 用于回放）
// @Summary 获取单条连线记录详情
// @Tags 审计连线
// @Accept json
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} model.SSHConnectionRecord
// @Router /connection-audit/{id} [get]
func (h *Handler) GetConnectionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rec, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, rec)
}
