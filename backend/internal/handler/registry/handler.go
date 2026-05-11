// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package registry

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/integration/provider"
	"github.com/kkops/backend/internal/integration/registry"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler exposes Harbor registry proxy endpoints.
type Handler struct {
	svc *integrationsvc.Service
}

// NewHandler constructs handler.
func NewHandler(svc *integrationsvc.Service) *Handler {
	return &Handler{svc: svc}
}

func parseUintQuery(c *gin.Context, key string) (uint, bool) {
	q := c.Query(key)
	if q == "" {
		return 0, false
	}
	v, err := strconv.ParseUint(q, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(v), true
}

// ListRepositories GET /registry/repositories
func (h *Handler) ListRepositories(c *gin.Context) {
	iid, ok := parseUintQuery(c, "integration_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id is required"})
		return
	}
	pub, err := h.svc.Get(iid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindHarbor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be harbor"})
		return
	}
	cfg, err := h.svc.DecryptConfigForWorker(iid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	cli, err := registry.NewHarborClientFromConfig(cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	projectName := c.Query("project_name")
	repos, err := cli.ListRepositories(c.Request.Context(), projectName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": repos})
}

// ListTags GET /registry/tags?integration_id=&repository=project/repo
func (h *Handler) ListTags(c *gin.Context) {
	iid, ok := parseUintQuery(c, "integration_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id is required"})
		return
	}
	repo := c.Query("repository")
	if repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repository query required (project/name)"})
		return
	}
	pub, err := h.svc.Get(iid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindHarbor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be harbor"})
		return
	}
	cfg, err := h.svc.DecryptConfigForWorker(iid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	cli, err := registry.NewHarborClientFromConfig(cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "30"))
	tags, err := cli.ListTags(c.Request.Context(), repo, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// Vulnerabilities GET /registry/vulnerabilities?integration_id=&repository=&reference=
func (h *Handler) Vulnerabilities(c *gin.Context) {
	iid, ok := parseUintQuery(c, "integration_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id is required"})
		return
	}
	repo := c.Query("repository")
	ref := c.Query("reference")
	if repo == "" || ref == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repository and reference query required"})
		return
	}
	pub, err := h.svc.Get(iid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindHarbor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be harbor"})
		return
	}
	cfg, err := h.svc.DecryptConfigForWorker(iid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return
	}
	cli, err := registry.NewHarborClientFromConfig(cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := cli.GetVulnerabilities(c.Request.Context(), repo, ref)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": v})
}
