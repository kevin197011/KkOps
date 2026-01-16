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
)

// Handler handles asset management HTTP requests
type Handler struct {
	service *asset.Service
}

// NewHandler creates a new asset handler
func NewHandler(service *asset.Service) *Handler {
	return &Handler{service: service}
}

// CreateAsset handles asset creation
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

// GetAsset handles asset retrieval
func (h *Handler) GetAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	resp, err := h.service.GetAsset(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListAssets handles asset list retrieval with filtering
func (h *Handler) ListAssets(c *gin.Context) {
	filter := asset.ListAssetsFilter{
		Page:     1,
		PageSize: 20,
	}

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
func (h *Handler) ExportAssets(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=assets.csv")

	if err := h.service.ExportAssets(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// ImportAssets handles asset import
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
