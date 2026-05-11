// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package gitopsview aggregates CI/CD and Argo CD signals into a unified timeline.
package gitopsview

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/kkops/backend/internal/integration/cicd"
	"github.com/kkops/backend/internal/integration/gitops"
	"github.com/kkops/backend/internal/integration/provider"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// PipelineEvent is a normalized timeline row for the pipeline view UI.
type PipelineEvent struct {
	Ts     time.Time `json:"ts"`
	Kind   string    `json:"kind"`
	Source string    `json:"source"`
	Ref    string    `json:"ref"`
	Status string    `json:"status"`
	Link   string    `json:"link"`
}

// Service builds merged timelines from configured integrations.
type Service struct {
	integration *integrationsvc.Service
}

// NewService constructs the aggregator.
func NewService(integration *integrationsvc.Service) *Service {
	return &Service{integration: integration}
}

func parseRFC3339Flexible(s string) (time.Time, bool) {
	if strings.TrimSpace(s) == "" {
		return time.Time{}, false
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05Z07:00"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// PipelineView returns merged events sorted newest-first.
func (s *Service) PipelineView(ctx context.Context, app string, argoIntegrationID *uint) ([]PipelineEvent, error) {
	if strings.TrimSpace(app) == "" {
		return nil, fmt.Errorf("app is required")
	}
	list, err := s.integration.List()
	if err != nil {
		return nil, err
	}

	out := make([]PipelineEvent, 0)

	// Argo CD revision history for the application (first Argo integration that resolves the app).
	for _, in := range list {
		if !in.Enabled || provider.NormalizeKind(in.Kind) != provider.KindArgoCD {
			continue
		}
		if argoIntegrationID != nil && in.ID != *argoIntegrationID {
			continue
		}
		cfg, err := s.integration.DecryptConfigForWorker(in.ID)
		if err != nil {
			continue
		}
		cli, err := gitops.NewArgoCDClientFromConfig(cfg)
		if err != nil {
			continue
		}
		hist, err := cli.ListHistory(ctx, app)
		if err != nil {
			continue
		}
		b, _ := provider.ParseBase(cfg)
		base := strings.TrimSuffix(b.BaseURL, "/")
		for _, h := range hist {
			ts, ok := parseRFC3339Flexible(h.DeployedAt)
			if !ok {
				ts = time.Unix(0, 0).UTC()
			}
			link := fmt.Sprintf("%s/applications/%s", base, url.PathEscape(app))
			out = append(out, PipelineEvent{
				Ts:     ts,
				Kind:   "gitops",
				Source: in.Name,
				Ref:    h.Revision,
				Status: "synced",
				Link:   link,
			})
		}
		break
	}

	// GitLab CI recent pipelines (timestamped via GitLab API).
	for _, in := range list {
		if !in.Enabled || provider.NormalizeKind(in.Kind) != provider.KindGitLab {
			continue
		}
		cfg, err := s.integration.DecryptConfigForWorker(in.ID)
		if err != nil {
			continue
		}
		cli, err := cicd.NewGitLabCIClient(cfg)
		if err != nil {
			continue
		}
		pipes, err := cli.ListPipelines(ctx, "")
		if err != nil {
			continue
		}
		for _, p := range pipes {
			ts := time.Now().UTC()
			if p.RecordedAt != nil {
				ts = p.RecordedAt.UTC()
			}
			out = append(out, PipelineEvent{
				Ts:     ts,
				Kind:   "cicd",
				Source: in.Name,
				Ref:    p.Ref,
				Status: p.Status,
				Link:   p.WebURL,
			})
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Ts.After(out[j].Ts)
	})
	return out, nil
}
