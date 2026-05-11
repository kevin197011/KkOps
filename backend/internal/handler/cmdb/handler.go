// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmdb

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/cmdb"
)

// Handler exposes CMDB and topology HTTP APIs.
type Handler struct {
	svc *cmdb.Service
}

// NewHandler constructs Handler.
func NewHandler(svc *cmdb.Service) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on g (caller prefixes with /cmdb and /topology).
func (h *Handler) Register(cmdbGroup *gin.RouterGroup, topologyGroup *gin.RouterGroup) {
	cmdbGroup.GET("/assets", h.ListAssets)
	cmdbGroup.POST("/assets", h.CreateAsset)
	cmdbGroup.GET("/assets/:id", h.GetAsset)
	cmdbGroup.PUT("/assets/:id", h.UpdateAsset)
	cmdbGroup.DELETE("/assets/:id", h.DeleteAsset)

	cmdbGroup.GET("/asset-relations", h.ListRelations)
	cmdbGroup.POST("/asset-relations", h.CreateRelation)
	cmdbGroup.DELETE("/asset-relations/:id", h.DeleteRelation)

	topologyGroup.GET("/graph", h.GetTopologyGraph)
}

// ListAssets GET /cmdb/assets
func (h *Handler) ListAssets(c *gin.Context) {
	kind := c.Query("kind")
	env := c.Query("env")
	q := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	rows, total, err := h.svc.ListAssets(c.Request.Context(), kind, env, q, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total})
}

// CreateAsset POST /cmdb/assets
func (h *Handler) CreateAsset(c *gin.Context) {
	var body struct {
		cmdb.CreateAssetInput
		Labels map[string]any `json:"labels"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	labelsJSON := "{}"
	if body.Labels != nil {
		b, err := json.Marshal(body.Labels)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labels"})
			return
		}
		labelsJSON = string(b)
	}
	row, err := h.svc.CreateAsset(c.Request.Context(), body.CreateAssetInput, labelsJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": row})
}

// GetAsset GET /cmdb/assets/:id
func (h *Handler) GetAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	row, err := h.svc.GetAsset(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if row == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

// UpdateAsset PUT /cmdb/assets/:id
func (h *Handler) UpdateAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Name          *string        `json:"name"`
		Kind          *string        `json:"kind"`
		Env           *string        `json:"env"`
		OwnerUserID   *uint          `json:"owner_user_id"`
		Labels        map[string]any `json:"labels"`
		IntegrationID *uint          `json:"integration_id"`
		ExternalRef   *string        `json:"external_ref"`
		Notes         *string        `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in := cmdb.UpdateAssetInput{
		Name:          body.Name,
		Kind:          body.Kind,
		Env:           body.Env,
		OwnerUserID:   body.OwnerUserID,
		IntegrationID: body.IntegrationID,
		ExternalRef:   body.ExternalRef,
		Notes:         body.Notes,
	}
	if body.Labels != nil {
		b, err := json.Marshal(body.Labels)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid labels"})
			return
		}
		s := string(b)
		in.LabelsJSON = &s
	}
	row, err := h.svc.UpdateAsset(c.Request.Context(), uint(id), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if row == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

// DeleteAsset DELETE /cmdb/assets/:id
func (h *Handler) DeleteAsset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteAsset(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ListRelations GET /cmdb/asset-relations
func (h *Handler) ListRelations(c *gin.Context) {
	var fromID, toID uint
	if v := c.Query("from_asset_id"); v != "" {
		u, _ := strconv.ParseUint(v, 10, 32)
		fromID = uint(u)
	}
	if v := c.Query("to_asset_id"); v != "" {
		u, _ := strconv.ParseUint(v, 10, 32)
		toID = uint(u)
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "200"))
	rows, total, err := h.svc.ListRelations(c.Request.Context(), fromID, toID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total})
}

// CreateRelation POST /cmdb/asset-relations
func (h *Handler) CreateRelation(c *gin.Context) {
	var body struct {
		cmdb.CreateRelationInput
		Meta map[string]any `json:"meta"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	metaJSON := "{}"
	if body.Meta != nil {
		b, err := json.Marshal(body.Meta)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meta"})
			return
		}
		metaJSON = string(b)
	}
	row, err := h.svc.CreateRelation(c.Request.Context(), body.CreateRelationInput, metaJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": row})
}

// DeleteRelation DELETE /cmdb/asset-relations/:id
func (h *Handler) DeleteRelation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteRelation(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetTopologyGraph GET /topology/graph
func (h *Handler) GetTopologyGraph(c *gin.Context) {
	g, err := h.svc.BuildTopologyGraph(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": g})
}
