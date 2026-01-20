// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package asset

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/asset"
	"github.com/kkops/backend/internal/service/authorization"
)

// Handler handles asset management HTTP requests
type Handler struct {
	service  *asset.Service
	authzSvc *authorization.Service
}

// NewHandler creates a new asset handler
func NewHandler(service *asset.Service, authzSvc *authorization.Service) *Handler {
	return &Handler{
		service:  service,
		authzSvc: authzSvc,
	}
}

// CreateAsset handles asset creation
// @Summary Create asset
// @Description Create a new asset
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body asset.CreateAssetRequest true "Create asset request"
// @Success 201 {object} asset.AssetResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/assets [post]
func (h *Handler) CreateAsset(c *gin.Context) {
	var req asset.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateAsset(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetAsset handles asset retrieval with permission check
// @Summary Get asset
// @Description Get asset by ID with permission check
// @Tags assets
// @Produce json
// @Security BearerAuth
// @Param id path int true "Asset ID"
// @Success 200 {object} asset.AssetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/assets/{id} [get]
func (h *Handler) GetAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 检查资产访问权限
	hasAccess, err := h.authzSvc.HasAssetAccess(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	resp, err := h.service.GetAsset(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListAssets handles asset list retrieval with filtering and permission check
// @Summary List assets
// @Description Get paginated list of assets with filtering and permission check
// @Tags assets
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param project_id query int false "Filter by project ID"
// @Param cloud_platform_id query int false "Filter by cloud platform ID"
// @Param environment_id query int false "Filter by environment ID"
// @Param status query string false "Filter by status"
// @Param ip query string false "Filter by IP address"
// @Param ssh_key_id query int false "Filter by SSH key ID"
// @Param tag_ids query string false "Filter by tag IDs (comma-separated)"
// @Param search query string false "Search by hostname or IP"
// @Success 200 {object} map[string]interface{} "Response with data, total, page, and size"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/assets [get]
func (h *Handler) ListAssets(c *gin.Context) {
	filter := asset.ListAssetsFilter{
		Page:     1,
		PageSize: 20,
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 获取用户权限范围
	allowedAssetIDs, isAdmin, err := h.authzSvc.GetUserAssetIDs(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
		return
	}
	filter.IsAdmin = isAdmin
	filter.AllowedAssetIDs = allowedAssetIDs

	// Parse pagination
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil {
		filter.PageSize = pageSize
	}

	// Parse filters
	if projectID, err := strconv.ParseUint(c.Query("project_id"), 10, 32); err == nil {
		id := uint(projectID)
		filter.ProjectID = &id
	}
	if cloudPlatformID, err := strconv.ParseUint(c.Query("cloud_platform_id"), 10, 32); err == nil {
		id := uint(cloudPlatformID)
		filter.CloudPlatformID = &id
	}
	if envID, err := strconv.ParseUint(c.Query("environment_id"), 10, 32); err == nil {
		id := uint(envID)
		filter.EnvironmentID = &id
	}
	filter.Status = c.Query("status")
	filter.IP = c.Query("ip")
	if sshKeyID, err := strconv.ParseUint(c.Query("ssh_key_id"), 10, 32); err == nil {
		id := uint(sshKeyID)
		filter.SSHKeyID = &id
	}
	filter.Search = c.Query("search")

	// Parse tag IDs
	if tagIDsStr := c.Query("tag_ids"); tagIDsStr != "" {
		tagIDs := parseUintSlice(tagIDsStr)
		filter.TagIDs = tagIDs
	}

	assets, total, err := h.service.ListAssets(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  assets,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// UpdateAsset handles asset update
// @Summary Update asset
// @Description Update asset information
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Asset ID"
// @Param request body asset.UpdateAssetRequest true "Update asset request"
// @Success 200 {object} asset.AssetResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/assets/{id} [put]
func (h *Handler) UpdateAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	var req asset.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateAsset(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteAsset handles asset deletion
// @Summary Delete asset
// @Description Delete an asset by ID
// @Tags assets
// @Security BearerAuth
// @Param id path int true "Asset ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/assets/{id} [delete]
func (h *Handler) DeleteAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	if err := h.service.DeleteAsset(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportAssets handles asset export
// @Summary Export assets
// @Description Export all assets to CSV format
// @Tags assets
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {file} file "CSV file"
// @Failure 500 {object} map[string]string
// @Router /api/v1/assets/export [get]
func (h *Handler) ExportAssets(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=assets.csv")

	if err := h.service.ExportAssets(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// ImportAssets handles asset import
// @Summary Import assets
// @Description Import assets from CSV file
// @Tags assets
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CSV file to import"
// @Success 200 {object} map[string]interface{} "Import result"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/assets/import [post]
func (h *Handler) ImportAssets(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	result, err := h.service.ImportAssets(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// parseUintSlice parses a comma-separated string of integers
func parseUintSlice(s string) []uint {
	if s == "" {
		return []uint{}
	}
	parts := strings.Split(s, ",")
	result := make([]uint, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			if id, err := strconv.ParseUint(part, 10, 32); err == nil {
				result = append(result, uint(id))
			}
		}
	}
	return result
}
