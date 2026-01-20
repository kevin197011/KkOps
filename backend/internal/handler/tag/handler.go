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
// @Summary Create tag
// @Description Create a new tag
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tag.CreateTagRequest true "Create tag request"
// @Success 201 {object} tag.TagResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/tags [post]
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
// @Summary Get tag
// @Description Get tag by ID
// @Tags tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} tag.TagResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/tags/{id} [get]
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
// @Summary List tags
// @Description Get list of all tags
// @Tags tags
// @Produce json
// @Security BearerAuth
// @Success 200 {array} tag.TagResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/tags [get]
func (h *Handler) ListTags(c *gin.Context) {
	tags, err := h.service.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// UpdateTag handles tag update
// @Summary Update tag
// @Description Update tag information
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Param request body tag.UpdateTagRequest true "Update tag request"
// @Success 200 {object} tag.TagResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/tags/{id} [put]
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
// @Summary Delete tag
// @Description Delete a tag by ID
// @Tags tags
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/tags/{id} [delete]
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
