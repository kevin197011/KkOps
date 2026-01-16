// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package environment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/environment"
)

// Handler handles environment management HTTP requests
type Handler struct {
	service *environment.Service
}

// NewHandler creates a new environment handler
func NewHandler(service *environment.Service) *Handler {
	return &Handler{service: service}
}

// CreateEnvironment handles environment creation
func (h *Handler) CreateEnvironment(c *gin.Context) {
	var req environment.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateEnvironment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetEnvironment handles environment retrieval
func (h *Handler) GetEnvironment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
		return
	}

	resp, err := h.service.GetEnvironment(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListEnvironments handles environment list retrieval
func (h *Handler) ListEnvironments(c *gin.Context) {
	environments, err := h.service.ListEnvironments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, environments)
}

// UpdateEnvironment handles environment update
func (h *Handler) UpdateEnvironment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
		return
	}

	var req environment.UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateEnvironment(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteEnvironment handles environment deletion
func (h *Handler) DeleteEnvironment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
		return
	}

	if err := h.service.DeleteEnvironment(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
