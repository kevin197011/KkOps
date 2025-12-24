package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

type CreatePermissionRequest struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required"`
	Action       string `json:"action" binding:"required"`
	Description  string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.permissionService.CreatePermission(
		req.Code,
		req.Name,
		req.ResourceType,
		req.Action,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"permission": permission})
}

func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	permission, err := h.permissionService.GetPermission(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	permissions, total, err := h.permissionService.ListPermissions(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

func (h *PermissionHandler) ListPermissionsByResourceType(c *gin.Context) {
	resourceType := c.Param("resource_type")
	if resourceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource_type is required"})
		return
	}

	permissions, err := h.permissionService.ListPermissionsByResourceType(resourceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.permissionService.UpdatePermission(id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	err = h.permissionService.DeletePermission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission deleted successfully"})
}

