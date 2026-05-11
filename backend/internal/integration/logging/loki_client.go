// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kkops/backend/internal/integration/provider"
)

const lokiHTTPTimeout = 30 * time.Second

type lokiResp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// LokiClient searches Grafana Loki.
type LokiClient struct {
	BaseURL string
	Token   string
}

// Search runs a LogQL query over [start,end).
func (c *LokiClient) Search(ctx context.Context, query string, start, end time.Time, limit int) ([]LogLine, error) {
	if limit <= 0 {
		limit = 100
	}
	q := url.Values{}
	q.Set("query", query)
	q.Set("limit", strconv.Itoa(limit))
	q.Set("start", strconv.FormatInt(toNanos(start), 10))
	q.Set("end", strconv.FormatInt(toNanos(end), 10))
	u := strings.TrimSuffix(c.BaseURL, "/") + "/loki/api/v1/query_range?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	cli := &http.Client{Timeout: lokiHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("loki request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("loki status %d: %s", resp.StatusCode, truncate(string(body), 180))
	}
	var wrap lokiResp
	if err := json.Unmarshal(body, &wrap); err != nil {
		return nil, fmt.Errorf("decode loki: %w", err)
	}
	if wrap.Status != "success" {
		return nil, fmt.Errorf("loki query failed")
	}
	var lines []LogLine
	for _, stream := range wrap.Data.Result {
		for _, pair := range stream.Values {
			if len(pair) < 2 {
				continue
			}
			tsNano := pair[0]
			msg := pair[1]
			lines = append(lines, LogLine{
				Timestamp: formatLokiTS(tsNano),
				Message:   msg,
				Labels:    stream.Stream,
			})
		}
	}
	return lines, nil
}

func toNanos(t time.Time) int64 {
	if t.IsZero() {
		return time.Now().Add(-1 * time.Hour).UnixNano()
	}
	return t.UnixNano()
}

func formatLokiTS(ns string) string {
	// Loki returns nanoseconds as string
	n, err := strconv.ParseInt(ns, 10, 64)
	if err != nil {
		return ns
	}
	return time.Unix(0, n).UTC().Format(time.RFC3339Nano)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// NewLokiClientFromConfig parses BaseConfig.
func NewLokiClientFromConfig(cfgJSON []byte) (*LokiClient, error) {
	b, err := provider.ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &LokiClient{BaseURL: b.BaseURL, Token: b.Token}, nil
}
