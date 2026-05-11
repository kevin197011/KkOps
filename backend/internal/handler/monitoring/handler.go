// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package monitoring

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/integration/monitoring"
	"github.com/kkops/backend/internal/integration/provider"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler serves Prometheus / Nightingale metric queries through configured integrations.
type Handler struct {
	svc *integrationsvc.Service
}

// NewHandler constructs the monitoring handler.
func NewHandler(svc *integrationsvc.Service) *Handler {
	return &Handler{svc: svc}
}

type queryBody struct {
	IntegrationID uint    `json:"integration_id" binding:"required"`
	Query         string  `json:"query" binding:"required"`
	Time          *string `json:"time,omitempty"`
	Range         *struct {
		Start string `json:"start"`
		End   string `json:"end"`
		Step  string `json:"step"`
	} `json:"range,omitempty"`
}

// Query POST /monitoring/query
func (h *Handler) Query(c *gin.Context) {
	var req queryBody
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
	case provider.KindPrometheus:
		cli, err := monitoring.NewPrometheusClientFromConfig(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Range != nil {
			start, e1 := timeParse(req.Range.Start)
			end, e2 := timeParse(req.Range.End)
			if e1 != nil || e2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range start/end"})
				return
			}
			step, _ := time.ParseDuration(req.Range.Step)
			if step == 0 {
				step = 30 * time.Second
			}
			res, err := cli.QueryRange(c.Request.Context(), req.Query, start, end, step)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}
		at := time.Now()
		if req.Time != nil {
			if t, err := timeParse(*req.Time); err == nil {
				at = t
			}
		}
		res, err := cli.Query(c.Request.Context(), req.Query, at)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})

	case provider.KindNightingale:
		cli, err := monitoring.NewNightingaleClientFromConfig(cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Range != nil {
			start, e1 := timeParse(req.Range.Start)
			end, e2 := timeParse(req.Range.End)
			if e1 != nil || e2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid range start/end"})
				return
			}
			step, _ := time.ParseDuration(req.Range.Step)
			if step == 0 {
				step = 30 * time.Second
			}
			res, err := cli.QueryRange(c.Request.Context(), req.Query, start, end, step)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}
		at := time.Now()
		if req.Time != nil {
			if t, err := timeParse(*req.Time); err == nil {
				at = t
			}
		}
		res, err := cli.Query(c.Request.Context(), req.Query, at)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind is not a metrics backend"})
	}
}

func timeParse(s string) (time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339Nano, s)
}
