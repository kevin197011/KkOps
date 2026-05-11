// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	provisioningapi "github.com/kkops/backend/internal/provisioning"
)

// Handler exposes provisioning REST API.
type Handler struct {
	svc *provisioningapi.APIService
}

// NewHandler creates the handler.
func NewHandler(svc *provisioningapi.APIService) *Handler {
	return &Handler{svc: svc}
}

// ListTargets GET /provisioning/targets
func (h *Handler) ListTargets(c *gin.Context) {
	items, err := h.svc.ListTargets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// CreateTarget POST /provisioning/targets
func (h *Handler) CreateTarget(c *gin.Context) {
	var req provisioningapi.CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.svc.CreateTarget(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": t})
}

// SyncTarget POST /provisioning/targets/:id/sync
func (h *Handler) SyncTarget(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.RequestManualSync(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "target not found"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}

// ListRuns GET /provisioning/targets/:id/runs
func (h *Handler) ListRuns(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	runs, err := h.svc.ListRuns(uint(id), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": runs})
}
