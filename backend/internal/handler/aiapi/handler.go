// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package aiapi exposes REST endpoints for AI Ops features.
package aiapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kkops/backend/internal/ai"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/service/aisvc"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler serves AI routes.
type Handler struct {
	Svc      *aisvc.Service
	Integ    *integrationsvc.Service
	Registry *ai.Registry
}

// NewHandler constructs Handler.
func NewHandler(s *aisvc.Service, integ *integrationsvc.Service, reg *ai.Registry) *Handler {
	return &Handler{Svc: s, Integ: integ, Registry: reg}
}

// ensureChatProvider verifies an AI provider exists before chat persists session state.
func (h *Handler) ensureChatProvider(integrationID *uint) error {
	if integrationID != nil && *integrationID > 0 {
		_, err := h.Registry.ProviderForIntegration(*integrationID)
		return err
	}
	_, _, err := h.Registry.DefaultProvider()
	return err
}

// Register attaches routes under /ai (caller prefixes with /api/v1).
func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/providers", h.Providers)
	g.POST("/test", h.Test)
	g.POST("/chat", h.Chat)

	g.GET("/sessions", h.ListSessions)
	g.GET("/sessions/:id", h.GetSession)
	g.DELETE("/sessions/:id", h.DeleteSession)

	anom := g.Group("/anomaly")
	{
		anom.GET("/rules", h.ListRules)
		anom.POST("/rules", h.CreateRule)
		anom.PUT("/rules/:id", h.UpdateRule)
		anom.DELETE("/rules/:id", h.DeleteRule)
		anom.GET("/findings", h.ListFindings)
	}

	g.POST("/rca", h.PostRCA)
	g.GET("/rca/reports", h.ListRCAReports)
}

// Providers GET /ai/providers — integrations with kind ai.
func (h *Handler) Providers(c *gin.Context) {
	list, err := h.Integ.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]model.IntegrationPublic, 0)
	for _, in := range list {
		if in.Kind == "ai" {
			out = append(out, in)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

type testBody struct {
	IntegrationID *uint `json:"integration_id"`
}

// Test POST /ai/test — streams probe response as SSE.
func (h *Handler) Test(c *gin.Context) {
	var req testBody
	_ = c.ShouldBindJSON(&req)
	var p ai.Provider
	var err error
	if req.IntegrationID != nil && *req.IntegrationID > 0 {
		p, err = h.Registry.ProviderForIntegration(*req.IntegrationID)
	} else {
		p, _, err = h.Registry.DefaultProvider()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msgs := []ai.Message{{Role: "user", Content: "Reply with exactly: pong"}}
	if err := aisvc.StreamPlainSSE(c, p, msgs); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
}

type chatBody struct {
	SessionID     *uint                  `json:"session_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Messages      []ai.Message           `json:"messages"`
	Context       map[string]interface{} `json:"context,omitempty"`
}

// Chat POST /ai/chat — SSE stream of assistant reply.
func (h *Handler) Chat(c *gin.Context) {
	var req chatBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messages required"})
		return
	}

	if err := h.ensureChatProvider(req.IntegrationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := userID(c)

	var sid uint
	if req.SessionID != nil && *req.SessionID > 0 {
		sid = *req.SessionID
		var sess model.AIChatSession
		if err := h.Svc.DB.Where("id = ? AND user_id = ?", sid, uid).First(&sess).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}
	} else {
		title := "Chat"
		last := req.Messages[len(req.Messages)-1]
		title = trimTitle(last.Content)
		s, err := h.Svc.CreateSession(uid, title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sid = s.ID
	}

	last := req.Messages[len(req.Messages)-1]
	if last.Role == "user" && last.Content != "" {
		if _, err := h.Svc.AppendMessage(sid, "user", last.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.Svc.TouchSessionTitle(sid, last.Content)
	}

	var hist []model.AIChatMessage
	if err := h.Svc.DB.Where("session_id = ?", sid).Order("id asc").Find(&hist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	conv := make([]ai.Message, 0, len(hist))
	for _, m := range hist {
		conv = append(conv, ai.Message{Role: m.Role, Content: m.Content})
	}

	c.Header("X-Session-ID", strconv.FormatUint(uint64(sid), 10))

	reply, err := h.Svc.ChatTurn(c.Request.Context(), h.Registry, req.IntegrationID, conv, c)
	if err != nil {
		if !c.Writer.Written() {
			status := http.StatusBadGateway
			if errors.Is(err, ai.ErrNoAIIntegration) {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": err.Error()})
		}
		return
	}
	if _, err := h.Svc.AppendMessage(sid, "assistant", reply); err != nil && h.Svc.Log != nil {
		h.Svc.Log.Debug("failed to persist assistant message", zap.Error(err))
	}
}

func trimTitle(s string) string {
	s = trimOneLine(s)
	if len([]rune(s)) > 80 {
		rs := []rune(s)
		return string(rs[:80]) + "…"
	}
	return s
}

func trimOneLine(s string) string {
	for i, r := range s {
		if r == '\n' || r == '\r' {
			return s[:i]
		}
	}
	return s
}

func userID(c *gin.Context) uint {
	v, ok := c.Get("user_id")
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case uint:
		return t
	case float64:
		return uint(t)
	default:
		return 0
	}
}

// ListSessions GET /ai/sessions
func (h *Handler) ListSessions(c *gin.Context) {
	rows, err := h.Svc.ListSessions(userID(c), 80)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// GetSession GET /ai/sessions/:id
func (h *Handler) GetSession(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	sess, msgs, err := h.Svc.GetSession(userID(c), uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"session": sess, "messages": msgs}})
}

// DeleteSession DELETE /ai/sessions/:id
func (h *Handler) DeleteSession(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	if err := h.Svc.DeleteSession(userID(c), uint(id64)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ListRules GET /ai/anomaly/rules
func (h *Handler) ListRules(c *gin.Context) {
	rows, err := h.Svc.ListRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// CreateRule POST /ai/anomaly/rules
func (h *Handler) CreateRule(c *gin.Context) {
	var body model.AIAnomalyRule
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.ID = 0
	row, err := h.Svc.SaveRule(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": row})
}

// UpdateRule PUT /ai/anomaly/rules/:id
func (h *Handler) UpdateRule(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	var body model.AIAnomalyRule
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body.ID = uint(id64)
	row, err := h.Svc.SaveRule(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

// DeleteRule DELETE /ai/anomaly/rules/:id
func (h *Handler) DeleteRule(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	if err := h.Svc.DeleteRule(uint(id64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ListFindings GET /ai/anomaly/findings
func (h *Handler) ListFindings(c *gin.Context) {
	var ruleID uint
	if v := c.Query("rule_id"); v != "" {
		if u, err := strconv.ParseUint(v, 10, 32); err == nil {
			ruleID = uint(u)
		}
	}
	var since *time.Time
	if v := c.Query("since"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			since = &t
		}
	}
	limit := 100
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	rows, err := h.Svc.ListFindings(ruleID, since, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

type rcaBody struct {
	IncidentID             uint  `json:"incident_id"`
	IntegrationIDOverrides *uint `json:"integration_id_overrides"`
}

// PostRCA POST /ai/rca
func (h *Handler) PostRCA(c *gin.Context) {
	var req rcaBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.IncidentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incident_id required"})
		return
	}
	res, err := h.Svc.GenerateRCA(c.Request.Context(), h.Registry, req.IncidentID, req.IntegrationIDOverrides)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	row, err := h.Svc.PersistRCA(req.IncidentID, res.ReportMD, res.RawToolLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"report_md":        res.ReportMD,
			"cited_tool_calls": res.CitedToolCalls,
			"stored_id":        row.ID,
		},
	})
}

// ListRCAReports GET /ai/rca/reports
func (h *Handler) ListRCAReports(c *gin.Context) {
	var incID uint
	if v := c.Query("incident_id"); v != "" {
		if u, err := strconv.ParseUint(v, 10, 32); err == nil {
			incID = uint(u)
		}
	}
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	rows, err := h.Svc.ListRcaReports(incID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
