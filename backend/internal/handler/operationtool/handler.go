// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package operationtool

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/operationtool"
)

// Handler 运维工具处理器
type Handler struct {
	service *operationtool.Service
}

// NewHandler 创建运维工具处理器实例
func NewHandler(service *operationtool.Service) *Handler {
	return &Handler{service: service}
}

// List 获取工具列表
// @Summary 获取运维工具列表
// @Description 获取运维工具列表，支持按分类过滤
// @Tags Operation Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "分类过滤"
// @Param enabled query bool false "是否只返回启用的工具（默认 true）"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /operation-tools [get]
func (h *Handler) List(c *gin.Context) {
	opts := &operationtool.ListOptions{}

	// 获取分类过滤参数
	if category := c.Query("category"); category != "" {
		opts.Category = &category
	}

	// 获取启用状态过滤参数
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			opts.Enabled = &enabled
		}
	}

	tools, err := h.service.List(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tools})
}

// Get 获取单个工具
// @Summary 获取运维工具详情
// @Description 根据 ID 获取运维工具详情
// @Tags Operation Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工具 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operation-tools/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工具 ID"})
		return
	}

	tool, err := h.service.Get(uint(id))
	if err != nil {
		if err.Error() == "工具不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tool})
}

// Create 创建工具
// @Summary 创建运维工具
// @Description 创建新的运维工具（管理员）
// @Tags Operation Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tool body operationtool.CreateToolRequest true "工具信息"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operation-tools [post]
func (h *Handler) Create(c *gin.Context) {
	var req operationtool.CreateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tool, err := h.service.Create(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "工具名称不能为空" || err.Error() == "工具 URL 不能为空" ||
			err.Error() == "URL must start with http:// or https://" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tool})
}

// Update 更新工具
// @Summary 更新运维工具
// @Description 更新运维工具信息（管理员）
// @Tags Operation Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工具 ID"
// @Param tool body operationtool.UpdateToolRequest true "工具信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operation-tools/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工具 ID"})
		return
	}

	var req operationtool.UpdateToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tool, err := h.service.Update(uint(id), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "工具不存在" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "URL must start with http:// or https://" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tool})
}

// Delete 删除工具
// @Summary 删除运维工具
// @Description 删除运维工具（管理员）
// @Tags Operation Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工具 ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operation-tools/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工具 ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "工具不存在" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
