package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/service"
)

type CloudPlatformHandler struct {
	service service.CloudPlatformService
}

func NewCloudPlatformHandler(service service.CloudPlatformService) *CloudPlatformHandler {
	return &CloudPlatformHandler{service: service}
}

type CreateCloudPlatformRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
	Description string `json:"description"`
}

type UpdateCloudPlatformRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
	Description string `json:"description"`
}

// CreateCloudPlatform 创建云平台
func (h *CloudPlatformHandler) CreateCloudPlatform(c *gin.Context) {
	var req CreateCloudPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cp := &models.CloudPlatform{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Icon:        req.Icon,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
		Description: req.Description,
	}

	if err := h.service.Create(cp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"cloud_platform": cp})
}

// GetCloudPlatform 获取云平台详情
func (h *CloudPlatformHandler) GetCloudPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cp, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud platform not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cloud_platform": cp})
}

// ListCloudPlatforms 获取云平台列表
func (h *CloudPlatformHandler) ListCloudPlatforms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}

	cps, total, err := h.service.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cloud_platforms": cps,
		"total":           total,
		"page":            page,
		"page_size":       pageSize,
	})
}

// UpdateCloudPlatform 更新云平台
func (h *CloudPlatformHandler) UpdateCloudPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateCloudPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cp := &models.CloudPlatform{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Icon:        req.Icon,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
		Description: req.Description,
	}

	if err := h.service.Update(id, cp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := h.service.GetByID(id)
	c.JSON(http.StatusOK, gin.H{"cloud_platform": updated})
}

// DeleteCloudPlatform 删除云平台
func (h *CloudPlatformHandler) DeleteCloudPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cloud platform deleted"})
}
