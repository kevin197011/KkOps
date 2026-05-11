// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const alertmanagerHTTPTimeout = 25 * time.Second

// AlertmanagerV2Alert is a minimal subset of Alertmanager GET /api/v2/alerts elements.
type AlertmanagerV2Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	State        string            `json:"state"`
	ActiveAt     *string           `json:"activeAt,omitempty"`
	StartsAt     *string           `json:"startsAt,omitempty"`
	EndsAt       *string           `json:"endsAt,omitempty"`
	Fingerprint  string            `json:"fingerprint,omitempty"`
	GeneratorURL string            `json:"generatorURL,omitempty"`
}

// FetchAlertmanagerAlerts retrieves active alerts from Alertmanager HTTP API v2.
func FetchAlertmanagerAlerts(ctx context.Context, baseURL, token string) ([]AlertmanagerV2Alert, error) {
	base := strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	if base == "" {
		return nil, fmt.Errorf("alertmanager base URL is empty")
	}
	u := base + "/api/v2/alerts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	cli := &http.Client{Timeout: alertmanagerHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("alertmanager request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("alertmanager status %d: %s", resp.StatusCode, truncateStr(string(body), 256))
	}
	var out []AlertmanagerV2Alert
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode alertmanager alerts: %w", err)
	}
	return out, nil
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
