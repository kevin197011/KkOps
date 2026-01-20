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
// @Summary Create environment
// @Description Create a new environment
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body environment.CreateEnvironmentRequest true "Create environment request"
// @Success 201 {object} environment.EnvironmentResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/environments [post]
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
// @Summary Get environment
// @Description Get environment by ID
// @Tags environments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Environment ID"
// @Success 200 {object} environment.EnvironmentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/environments/{id} [get]
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
// @Summary List environments
// @Description Get list of all environments
// @Tags environments
// @Produce json
// @Security BearerAuth
// @Success 200 {array} environment.EnvironmentResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/environments [get]
func (h *Handler) ListEnvironments(c *gin.Context) {
	environments, err := h.service.ListEnvironments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, environments)
}

// UpdateEnvironment handles environment update
// @Summary Update environment
// @Description Update environment information
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Environment ID"
// @Param request body environment.UpdateEnvironmentRequest true "Update environment request"
// @Success 200 {object} environment.EnvironmentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/environments/{id} [put]
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
// @Summary Delete environment
// @Description Delete an environment by ID
// @Tags environments
// @Security BearerAuth
// @Param id path int true "Environment ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/environments/{id} [delete]
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
