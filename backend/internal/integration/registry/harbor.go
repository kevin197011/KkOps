// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package registry implements Harbor HTTP client helpers.
package registry

import (
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

const harborHTTPTimeout = 30 * time.Second

// HarborClient talks to Harbor registry API v2.0.
type HarborClient struct {
	BaseURL string
	Token   string
}

// RepositorySummary is a single repo reference inside a project.
type RepositorySummary struct {
	Name       string `json:"name"`
	ProjectID  int    `json:"project_id"`
	PullCount  int64  `json:"pull_count"`
	UpdateTime string `json:"update_time"`
}

// TagSummary describes an artifact tag.
type TagSummary struct {
	Name     string `json:"name"`
	Digest   string `json:"digest,omitempty"`
	Size     int64  `json:"size,omitempty"`
	PushTime string `json:"push_time,omitempty"`
}

// NewHarborClientFromConfig parses standard integration JSON.
func NewHarborClientFromConfig(cfgJSON []byte) (*HarborClient, error) {
	b, err := provider.ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &HarborClient{BaseURL: b.BaseURL, Token: b.Token}, nil
}

func (h *HarborClient) auth(req *http.Request) {
	if h.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.Token)
	}
	req.Header.Set("Accept", "application/json")
}

// ListRepositories lists repositories (optionally scoped by project name prefix).
func (h *HarborClient) ListRepositories(ctx context.Context, projectName string, page, pageSize int) ([]RepositorySummary, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	path := "/api/v2.0/repositories"
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("page_size", fmt.Sprintf("%d", pageSize))
	if projectName != "" {
		q.Set("project_name", projectName)
	}
	u := strings.TrimSuffix(h.BaseURL, "/") + path + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	h.auth(req)
	cli := &http.Client{Timeout: harborHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("harbor repos: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("harbor repos status %d: %s", resp.StatusCode, truncate(string(body), 160))
	}
	var rows []RepositorySummary
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, fmt.Errorf("decode harbor repos: %w", err)
	}
	return rows, nil
}

// ListTags lists tags for repository full name `project/repo`.
func (h *HarborClient) ListTags(ctx context.Context, repoFullName string, page, pageSize int) ([]TagSummary, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 30
	}
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("page_size", fmt.Sprintf("%d", pageSize))
	parts := strings.SplitN(repoFullName, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("repository must be project/name format")
	}
	project, repo := parts[0], parts[1]
	u := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts?%s",
		strings.TrimSuffix(h.BaseURL, "/"),
		url.PathEscape(project),
		url.PathEscape(repo),
		q.Encode(),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	h.auth(req)
	cli := &http.Client{Timeout: harborHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("harbor artifacts: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("harbor artifacts status %d", resp.StatusCode)
	}
	var artifacts []struct {
		Digest string `json:"digest"`
		Tags   []struct {
			Name string `json:"name"`
		} `json:"tags"`
		Size int64 `json:"size"`
	}
	if err := json.Unmarshal(body, &artifacts); err != nil {
		return nil, fmt.Errorf("decode artifacts: %w", err)
	}
	var tags []TagSummary
	for _, a := range artifacts {
		for _, t := range a.Tags {
			tags = append(tags, TagSummary{Name: t.Name, Digest: a.Digest, Size: a.Size})
		}
	}
	return tags, nil
}

// VulnReport is a minimal vulnerability summary blob for UI display.
type VulnReport struct {
	RawJSON string `json:"raw_json"`
}

// GetVulnerabilities returns vulnerability addition payload when scanning is enabled.
func (h *HarborClient) GetVulnerabilities(ctx context.Context, repoFullName, reference string) (*VulnReport, error) {
	parts := strings.SplitN(repoFullName, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("repository must be project/name format")
	}
	project, repo := parts[0], parts[1]
	u := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s/additions/vulnerabilities",
		strings.TrimSuffix(h.BaseURL, "/"),
		url.PathEscape(project),
		url.PathEscape(repo),
		url.PathEscape(reference),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	h.auth(req)
	cli := &http.Client{Timeout: harborHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("harbor vuln: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("harbor vuln status %d", resp.StatusCode)
	}
	return &VulnReport{RawJSON: string(body)}, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
