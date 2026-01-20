// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/category"
)

// Handler handles asset category management HTTP requests
type Handler struct {
	service *category.Service
}

// NewHandler creates a new category handler
func NewHandler(service *category.Service) *Handler {
	return &Handler{service: service}
}

// CreateCategory handles category creation
// @Summary Create category
// @Description Create a new asset category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body category.CreateCategoryRequest true "Create category request"
// @Success 201 {object} category.CategoryResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/asset-categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var req category.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateCategory(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetCategory handles category retrieval
// @Summary Get category
// @Description Get asset category by ID
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} category.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/asset-categories/{id} [get]
func (h *Handler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	resp, err := h.service.GetCategory(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListCategories handles category list retrieval
// @Summary List categories
// @Description Get list of all asset categories
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {array} category.CategoryResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/asset-categories [get]
func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.service.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory handles category update
// @Summary Update category
// @Description Update asset category information
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body category.UpdateCategoryRequest true "Update category request"
// @Success 200 {object} category.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/asset-categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req category.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateCategory(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteCategory handles category deletion
// @Summary Delete category
// @Description Delete an asset category by ID
// @Tags categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/asset-categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
