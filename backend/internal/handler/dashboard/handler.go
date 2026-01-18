// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/dashboard"
)

// Handler 仪表板处理器
type Handler struct {
	service *dashboard.Service
}

// NewHandler 创建仪表板处理器实例
func NewHandler(service *dashboard.Service) *Handler {
	return &Handler{service: service}
}

// GetStats 获取仪表板统计数据
// @Summary 获取仪表板统计数据
// @Description 获取系统概览统计数据，包括资产、用户、任务等信息
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dashboard.StatsResponse
// @Failure 500 {object} map[string]string
// @Router /dashboard/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
