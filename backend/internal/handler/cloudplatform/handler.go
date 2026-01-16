// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cloudplatform

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/cloudplatform"
)

// Handler handles cloud platform management HTTP requests
type Handler struct {
	service *cloudplatform.Service
}

// NewHandler creates a new cloud platform handler
func NewHandler(service *cloudplatform.Service) *Handler {
	return &Handler{service: service}
}

// CreateCloudPlatform handles cloud platform creation
func (h *Handler) CreateCloudPlatform(c *gin.Context) {
	var req cloudplatform.CreateCloudPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateCloudPlatform(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetCloudPlatform handles cloud platform retrieval
func (h *Handler) GetCloudPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud platform ID"})
		return
	}

	resp, err := h.service.GetCloudPlatform(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud platform not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListCloudPlatforms handles cloud platform list retrieval
func (h *Handler) ListCloudPlatforms(c *gin.Context) {
	cloudPlatforms, err := h.service.ListCloudPlatforms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cloudPlatforms)
}

// UpdateCloudPlatform handles cloud platform update
func (h *Handler) UpdateCloudPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud platform ID"})
		return
	}

	var req cloudplatform.UpdateCloudPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateCloudPlatform(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteCloudPlatform handles cloud platform deletion
func (h *Handler) DeleteCloudPlatform(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud platform ID"})
		return
	}

	if err := h.service.DeleteCloudPlatform(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud platform not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
