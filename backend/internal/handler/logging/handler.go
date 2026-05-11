// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logging

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/integration/logging"
	"github.com/kkops/backend/internal/integration/provider"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler exposes log search API.
type Handler struct {
	svc *integrationsvc.Service
}

// NewHandler constructs handler.
func NewHandler(svc *integrationsvc.Service) *Handler {
	return &Handler{svc: svc}
}

type searchBody struct {
	IntegrationID uint   `json:"integration_id" binding:"required"`
	Query         string `json:"query" binding:"required"`
	Start         string `json:"start"`
	End           string `json:"end"`
	Limit         int    `json:"limit"`
}

// Search POST /logging/search
func (h *Handler) Search(c *gin.Context) {
	var req searchBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pub, err := h.svc.Get(req.IntegrationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cfg, err := h.svc.DecryptConfigForWorker(req.IntegrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	if req.End != "" {
		if t, err := time.Parse(time.RFC3339, req.End); err == nil {
			end = t
		}
	}
	if req.Start != "" {
		if t, err := time.Parse(time.RFC3339, req.Start); err == nil {
			start = t
		}
	}
	k := provider.NormalizeKind(pub.Kind)
	switch k {
	case provider.KindLoki:
		cli, err := logging.NewLokiClientFromConfig(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		lines, err := cli.Search(c.Request.Context(), req.Query, start, end, req.Limit)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": lines})
	case provider.KindElasticsearch:
		cli, err := logging.NewElasticsearchClientFromConfig(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		lines, err := cli.Search(c.Request.Context(), req.Query, start, end, req.Limit)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": lines})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind is not a log backend"})
	}
}
