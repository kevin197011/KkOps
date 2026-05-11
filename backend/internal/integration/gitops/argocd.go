// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package gitops integrates Argo CD application APIs.
package gitops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kkops/backend/internal/integration/provider"
)

const argoHTTPTimeout = 30 * time.Second

// ArgoCDClient wraps Argo CD REST endpoints used by KkOps.
type ArgoCDClient struct {
	BaseURL string
	Token   string
}

// ApplicationSummary is a subset of Argo CD Application fields for the UI.
type ApplicationSummary struct {
	Name       string `json:"name"`
	Health     string `json:"health"`
	SyncStatus string `json:"sync_status"`
}

type argoApp struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		Health struct {
			Status string `json:"status"`
		} `json:"health"`
		Sync struct {
			Status string `json:"status"`
		} `json:"sync"`
	} `json:"status"`
}

// NewArgoCDClientFromConfig parses integration JSON.
func NewArgoCDClientFromConfig(cfgJSON []byte) (*ArgoCDClient, error) {
	b, err := provider.ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &ArgoCDClient{BaseURL: b.BaseURL, Token: b.Token}, nil
}

func (a *ArgoCDClient) auth(req *http.Request) {
	if a.Token != "" {
		req.Header.Set("Authorization", "Bearer "+a.Token)
	}
	req.Header.Set("Accept", "application/json")
}

// ListApplications returns application summaries (cluster scope default).
func (a *ArgoCDClient) ListApplications(ctx context.Context) ([]ApplicationSummary, error) {
	u := strings.TrimSuffix(a.BaseURL, "/") + "/api/v1/applications"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	a.auth(req)
	cli := &http.Client{Timeout: argoHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("argocd list apps: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("argocd list status %d: %s", resp.StatusCode, truncate(string(body), 160))
	}
	var wrap struct {
		Items []argoApp `json:"items"`
	}
	if err := json.Unmarshal(body, &wrap); err != nil {
		return nil, fmt.Errorf("decode argocd apps: %w", err)
	}
	out := make([]ApplicationSummary, 0, len(wrap.Items))
	for _, it := range wrap.Items {
		out = append(out, ApplicationSummary{
			Name:       it.Metadata.Name,
			Health:     it.Status.Health.Status,
			SyncStatus: it.Status.Sync.Status,
		})
	}
	return out, nil
}

// GetApplication loads one application by name (cluster-scoped default).
func (a *ArgoCDClient) GetApplication(ctx context.Context, name string) (*ApplicationSummary, error) {
	u := fmt.Sprintf("%s/api/v1/applications/%s", strings.TrimSuffix(a.BaseURL, "/"), url.PathEscape(name))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	a.auth(req)
	cli := &http.Client{Timeout: argoHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("argocd get app: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("argocd get status %d", resp.StatusCode)
	}
	var it argoApp
	if err := json.Unmarshal(body, &it); err != nil {
		return nil, err
	}
	return &ApplicationSummary{
		Name:       it.Metadata.Name,
		Health:     it.Status.Health.Status,
		SyncStatus: it.Status.Sync.Status,
	}, nil
}

// Sync triggers an application sync (async on Argo CD side).
func (a *ArgoCDClient) Sync(ctx context.Context, name string) error {
	u := fmt.Sprintf("%s/api/v1/applications/%s/sync", strings.TrimSuffix(a.BaseURL, "/"), url.PathEscape(name))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	a.auth(req)
	cli := &http.Client{Timeout: argoHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("argocd sync: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("argocd sync status %d: %s", resp.StatusCode, truncate(string(b), 160))
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// HistoryEntry is one record from Argo CD Application status.history.
type HistoryEntry struct {
	ID         int64  `json:"id"`
	Revision   string `json:"revision"`
	DeployedAt string `json:"deployed_at"`
}

type argoHistoryRaw struct {
	ID              int64  `json:"id"`
	Revision        string `json:"revision"`
	DeployedAt      string `json:"deployedAt"`
	DeployStartedAt string `json:"deployStartedAt,omitempty"`
}

type argoAppDetail struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		History []argoHistoryRaw `json:"history"`
		Sync    struct {
			Revision string `json:"revision"`
			Status   string `json:"status"`
		} `json:"sync"`
		Health struct {
			Status string `json:"status"`
		} `json:"health"`
	} `json:"status"`
}

// GetApplicationDetail loads application JSON including revision history.
func (a *ArgoCDClient) GetApplicationDetail(ctx context.Context, name string) (*argoAppDetail, error) {
	u := fmt.Sprintf("%s/api/v1/applications/%s", strings.TrimSuffix(a.BaseURL, "/"), url.PathEscape(name))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	a.auth(req)
	cli := &http.Client{Timeout: argoHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("argocd get application detail: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("argocd get application detail status %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var detail argoAppDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("decode argocd application detail: %w", err)
	}
	return &detail, nil
}

// ListHistory returns normalized history newest-first for timeline views.
func (a *ArgoCDClient) ListHistory(ctx context.Context, appName string) ([]HistoryEntry, error) {
	d, err := a.GetApplicationDetail(ctx, appName)
	if err != nil {
		return nil, err
	}
	out := make([]HistoryEntry, 0, len(d.Status.History))
	for _, h := range d.Status.History {
		ts := h.DeployedAt
		if ts == "" {
			ts = h.DeployStartedAt
		}
		out = append(out, HistoryEntry{
			ID:         h.ID,
			Revision:   h.Revision,
			DeployedAt: ts,
		})
	}
	// newest last in API; reverse for timeline descending
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}
