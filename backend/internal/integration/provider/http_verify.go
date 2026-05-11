// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const httpTimeout = 15 * time.Second

// VerifyHTTPGET probes baseURL+path with optional Bearer token; accepts 2xx as success.
func VerifyHTTPGET(ctx context.Context, baseURL, token, path string) error {
	u := baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/json, */*")

	cli := &http.Client{Timeout: httpTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d from %s", resp.StatusCode, trimURL(u))
	}
	return nil
}

func trimURL(u string) string {
	if len(u) > 80 {
		return u[:80] + "..."
	}
	return u
}

// Kind constants for integrations hub cards (lowercase, stable).
const (
	KindPrometheus    = "prometheus"
	KindNightingale   = "nightingale"
	KindLoki          = "loki"
	KindElasticsearch = "elasticsearch"
	KindGrafana       = "grafana"
	KindJenkins       = "jenkins"
	KindGitLab        = "gitlab"
	KindHarbor        = "harbor"
	KindArgoCD        = "argocd"
	KindKubernetes    = "kubernetes"
)

// NormalizeKind maps UI aliases to canonical kind strings.
func NormalizeKind(k string) string {
	switch strings.ToLower(strings.TrimSpace(k)) {
	case "prometheus", "prom":
		return KindPrometheus
	case "nightingale", "n9e":
		return KindNightingale
	case "loki":
		return KindLoki
	case "elasticsearch", "es", "elastic":
		return KindElasticsearch
	case "grafana":
		return KindGrafana
	case "jenkins":
		return KindJenkins
	case "gitlab":
		return KindGitLab
	case "harbor":
		return KindHarbor
	case "argocd", "argo_cd", "argo":
		return KindArgoCD
	case "kubernetes", "k8s":
		return KindKubernetes
	default:
		return strings.ToLower(strings.TrimSpace(k))
	}
}
