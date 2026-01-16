// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package tag

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/tag"
)

// Handler handles tag management HTTP requests
type Handler struct {
	service *tag.Service
}

// NewHandler creates a new tag handler
func NewHandler(service *tag.Service) *Handler {
	return &Handler{service: service}
}

// CreateTag handles tag creation
func (h *Handler) CreateTag(c *gin.Context) {
	var req tag.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateTag(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetTag handles tag retrieval
func (h *Handler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	resp, err := h.service.GetTag(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListTags handles tag list retrieval
func (h *Handler) ListTags(c *gin.Context) {
	tags, err := h.service.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// UpdateTag handles tag update
func (h *Handler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	var req tag.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateTag(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTag handles tag deletion
func (h *Handler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	if err := h.service.DeleteTag(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
