// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cicd

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/integration/cicd"
	"github.com/kkops/backend/internal/integration/provider"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler exposes CI/CD pipeline proxy APIs.
type Handler struct {
	svc *integrationsvc.Service
}

// NewHandler constructs handler.
func NewHandler(svc *integrationsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// ListPipelines GET /cicd/pipelines
func (h *Handler) ListPipelines(c *gin.Context) {
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
	cfg, err := h.svc.DecryptConfigForWorker(uint(iid64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	k := provider.NormalizeKind(pub.Kind)
	project := c.Query("project")
	switch k {
	case provider.KindJenkins:
		cli, err := cicd.NewJenkinsClient(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		items, err := cli.ListPipelines(c.Request.Context(), project)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		for i := range items {
			items[i].IntegrationID = uint(iid64)
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	case provider.KindGitLab:
		cli, err := cicd.NewGitLabCIClient(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		items, err := cli.ListPipelines(c.Request.Context(), project)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		for i := range items {
			items[i].IntegrationID = uint(iid64)
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be jenkins or gitlab"})
	}
}

type runBody struct {
	IntegrationID uint              `json:"integration_id" binding:"required"`
	Ref           string            `json:"ref"`
	Variables     map[string]string `json:"variables"`
}

// RunPipeline POST /cicd/pipelines/:id/run
func (h *Handler) RunPipeline(c *gin.Context) {
	pid := c.Param("id")
	var req runBody
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
	k := provider.NormalizeKind(pub.Kind)
	switch k {
	case provider.KindJenkins:
		cli, err := cicd.NewJenkinsClient(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := cli.TriggerPipeline(c.Request.Context(), pid, req.Ref, req.Variables); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	case provider.KindGitLab:
		cli, err := cicd.NewGitLabCIClient(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := cli.TriggerPipeline(c.Request.Context(), pid, req.Ref, req.Variables); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be jenkins or gitlab"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "triggered"})
}

// PipelineLogs GET /cicd/pipelines/:id/logs
func (h *Handler) PipelineLogs(c *gin.Context) {
	pid := c.Param("id")
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
	cfg, err := h.svc.DecryptConfigForWorker(uint(iid64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	k := provider.NormalizeKind(pub.Kind)
	switch k {
	case provider.KindJenkins:
		cli, err := cicd.NewJenkinsClient(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logs, err := cli.GetLogs(c.Request.Context(), pid)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": logs})
	case provider.KindGitLab:
		cli, err := cicd.NewGitLabCIClient(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logs, err := cli.GetLogs(c.Request.Context(), pid)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": logs})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be jenkins or gitlab"})
	}
}
