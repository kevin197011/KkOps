// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gitops

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/integration/gitops"
	"github.com/kkops/backend/internal/integration/provider"
	gitopsviewsvc "github.com/kkops/backend/internal/service/gitopsview"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler exposes Argo CD proxy endpoints.
type Handler struct {
	svc  *integrationsvc.Service
	view *gitopsviewsvc.Service
}

// NewHandler constructs handler.
func NewHandler(svc *integrationsvc.Service, view *gitopsviewsvc.Service) *Handler {
	return &Handler{svc: svc, view: view}
}

// PipelineView GET /gitops/pipeline-view
func (h *Handler) PipelineView(c *gin.Context) {
	app := c.Query("app")
	var argoID *uint
	if q := c.Query("argocd_integration_id"); q != "" {
		v64, err := strconv.ParseUint(q, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid argocd_integration_id"})
			return
		}
		v := uint(v64)
		argoID = &v
	}
	ev, err := h.view.PipelineView(c.Request.Context(), app, argoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ev})
}

// ListApplications GET /gitops/applications
func (h *Handler) ListApplications(c *gin.Context) {
	q := c.Query("integration_id")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id is required"})
		return
	}
	iid64, err := strconv.ParseUint(q, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id"})
		return
	}
	pub, err := h.svc.Get(uint(iid64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindArgoCD {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be argocd"})
		return
	}
	cfg, err := h.svc.DecryptConfigForWorker(uint(iid64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	cli, err := gitops.NewArgoCDClientFromConfig(cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	apps, err := cli.ListApplications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": apps})
}

// SyncApplication POST /gitops/applications/:name/sync
func (h *Handler) SyncApplication(c *gin.Context) {
	name := c.Param("name")
	var req struct {
		IntegrationID uint `json:"integration_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pub, err := h.svc.Get(req.IntegrationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindArgoCD {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be argocd"})
		return
	}
	cfg, err := h.svc.DecryptConfigForWorker(req.IntegrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	cli, err := gitops.NewArgoCDClientFromConfig(cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cli.Sync(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "sync initiated"})
}
