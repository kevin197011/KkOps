// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package alert

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	alertsvc "github.com/kkops/backend/internal/service/alert"
)

// Handler exposes alert APIs.
type Handler struct {
	svc *alertsvc.Service
}

// NewHandler constructs Handler.
func NewHandler(svc *alertsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// List GET /alerts
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	var integrationID uint
	if v := c.Query("integration_id"); v != "" {
		if u, err := strconv.ParseUint(v, 10, 32); err == nil {
			integrationID = uint(u)
		}
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	rows, total, err := h.svc.List(c.Request.Context(), alertsvc.ListParams{
		Status:        status,
		IntegrationID: integrationID,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total})
}

// Sync POST /alerts/sync
func (h *Handler) Sync(c *gin.Context) {
	n, errs := h.svc.Sync(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"synced": n, "warnings": errs})
}

// Acknowledge POST /alerts/:id/acknowledge
func (h *Handler) Acknowledge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Acknowledge(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Dismiss POST /alerts/:id/dismiss
func (h *Handler) Dismiss(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Dismiss(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Webhook POST /alerts/webhook (no JWT; secret middleware)
func (h *Handler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body"})
		return
	}
	n, err := h.svc.IngestWebhook(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ingested": n})
}
