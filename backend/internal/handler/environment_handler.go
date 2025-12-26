package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/service"
)

type EnvironmentHandler struct {
	service service.EnvironmentService
}

func NewEnvironmentHandler(service service.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{service: service}
}

type CreateEnvironmentRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
	Description string `json:"description"`
}

type UpdateEnvironmentRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
	Description string `json:"description"`
}

// CreateEnvironment 创建环境
func (h *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	var req CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	env := &models.Environment{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
		Description: req.Description,
	}

	if err := h.service.Create(env); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"environment": env})
}

// GetEnvironment 获取环境详情
func (h *EnvironmentHandler) GetEnvironment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	env, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// ListEnvironments 获取环境列表
func (h *EnvironmentHandler) ListEnvironments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}

	envs, total, err := h.service.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"environments": envs,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

// UpdateEnvironment 更新环境
func (h *EnvironmentHandler) UpdateEnvironment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	env := &models.Environment{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
		Description: req.Description,
	}

	if err := h.service.Update(id, env); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := h.service.GetByID(id)
	c.JSON(http.StatusOK, gin.H{"environment": updated})
}

// DeleteEnvironment 删除环境
func (h *EnvironmentHandler) DeleteEnvironment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "environment deleted"})
}
