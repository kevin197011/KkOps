// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cicd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GitLabCIClient calls GitLab HTTP API for pipelines.
type GitLabCIClient struct {
	cfg Config
}

// NewGitLabCIClient builds a client from JSON config (requires project_id for listing).
func NewGitLabCIClient(cfgJSON []byte) (*GitLabCIClient, error) {
	c, err := parseConfig(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &GitLabCIClient{cfg: c}, nil
}

type glPipeline struct {
	ID        int    `json:"id"`
	IID       int    `json:"iid"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	WebURL    string `json:"web_url"`
	UpdatedAt string `json:"updated_at"`
}

// ListPipelines lists recent pipelines for the configured project (project param ignored if config has project_id).
func (g *GitLabCIClient) ListPipelines(ctx context.Context, project string) ([]Pipeline, error) {
	pid := g.cfg.ProjectID
	if project != "" {
		if n, err := strconv.Atoi(project); err == nil {
			pid = n
		}
	}
	if pid == 0 {
		return nil, fmt.Errorf("gitlab project_id is required in integration config")
	}
	u := fmt.Sprintf("%s/api/v4/projects/%d/pipelines?per_page=30", g.cfg.BaseURL, pid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if g.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.cfg.Token)
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab list pipelines: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gitlab status %d: %s", resp.StatusCode, truncateStr(string(body), 160))
	}
	var rows []glPipeline
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, fmt.Errorf("decode gitlab pipelines: %w", err)
	}
	out := make([]Pipeline, 0, len(rows))
	for _, r := range rows {
		var rec *time.Time
		if r.UpdatedAt != "" {
			if t, err := time.Parse(time.RFC3339Nano, r.UpdatedAt); err == nil {
				rec = &t
			} else if t2, err2 := time.Parse(time.RFC3339, r.UpdatedAt); err2 == nil {
				rec = &t2
			}
		}
		out = append(out, Pipeline{
			ID:         fmt.Sprintf("gitlab:%d", r.ID),
			Name:       fmt.Sprintf("pipeline #%d", r.IID),
			Status:     r.Status,
			WebURL:     r.WebURL,
			Ref:        r.Ref,
			RecordedAt: rec,
		})
	}
	return out, nil
}

// GetPipeline fetches one pipeline by numeric id (without gitlab: prefix).
func (g *GitLabCIClient) GetPipeline(ctx context.Context, pipelineID string) (*Pipeline, error) {
	if g.cfg.ProjectID == 0 {
		return nil, fmt.Errorf("gitlab project_id is required in integration config")
	}
	idStr := strings.TrimPrefix(pipelineID, "gitlab:")
	u := fmt.Sprintf("%s/api/v4/projects/%d/pipelines/%s", g.cfg.BaseURL, g.cfg.ProjectID, idStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if g.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.cfg.Token)
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab get pipeline: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gitlab pipeline status %d", resp.StatusCode)
	}
	var r glPipeline
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	p := &Pipeline{
		ID:     fmt.Sprintf("gitlab:%d", r.ID),
		Name:   fmt.Sprintf("pipeline #%d", r.IID),
		Status: r.Status,
		WebURL: r.WebURL,
		Ref:    r.Ref,
	}
	return p, nil
}

// TriggerPipeline creates a new pipeline on ref.
func (g *GitLabCIClient) TriggerPipeline(ctx context.Context, _ string, ref string, vars map[string]string) error {
	if g.cfg.ProjectID == 0 {
		return fmt.Errorf("gitlab project_id is required in integration config")
	}
	if ref == "" {
		ref = "main"
	}
	payload := map[string]interface{}{"ref": ref}
	if len(vars) > 0 {
		vlist := make([]map[string]string, 0, len(vars))
		for k, v := range vars {
			vlist = append(vlist, map[string]string{"key": k, "value": v})
		}
		payload["variables"] = vlist
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/api/v4/projects/%d/pipeline", g.cfg.BaseURL, g.cfg.ProjectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if g.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.cfg.Token)
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("gitlab trigger: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gitlab trigger status %d: %s", resp.StatusCode, truncateStr(string(b), 160))
	}
	return nil
}

type glJob struct {
	ID int `json:"id"`
}

// GetLogs returns concatenated traces for all jobs in the pipeline (best-effort).
func (g *GitLabCIClient) GetLogs(ctx context.Context, pipelineID string) (string, error) {
	if g.cfg.ProjectID == 0 {
		return "", fmt.Errorf("gitlab project_id is required in integration config")
	}
	idStr := strings.TrimPrefix(pipelineID, "gitlab:")
	jobsURL := fmt.Sprintf("%s/api/v4/projects/%d/pipelines/%s/jobs", g.cfg.BaseURL, g.cfg.ProjectID, idStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jobsURL, nil)
	if err != nil {
		return "", err
	}
	if g.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.cfg.Token)
	}
	cli := &http.Client{Timeout: cicdHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("gitlab jobs: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("gitlab jobs status %d", resp.StatusCode)
	}
	var jobs []glJob
	if err := json.Unmarshal(body, &jobs); err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, j := range jobs {
		traceURL := fmt.Sprintf("%s/api/v4/projects/%d/jobs/%d/trace", g.cfg.BaseURL, g.cfg.ProjectID, j.ID)
		treq, err := http.NewRequestWithContext(ctx, http.MethodGet, traceURL, nil)
		if err != nil {
			continue
		}
		if g.cfg.Token != "" {
			treq.Header.Set("Authorization", "Bearer "+g.cfg.Token)
		}
		tresp, err := cli.Do(treq)
		if err != nil {
			continue
		}
		tb, _ := io.ReadAll(tresp.Body)
		tresp.Body.Close()
		sb.WriteString(string(tb))
		sb.WriteString("\n---\n")
	}
	return sb.String(), nil
}
