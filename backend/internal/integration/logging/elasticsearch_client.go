// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kkops/backend/internal/integration/provider"
)

const esHTTPTimeout = 30 * time.Second

// ElasticsearchClient runs _search against Elasticsearch / OpenSearch.
type ElasticsearchClient struct {
	BaseURL string
	Token   string
}

type esSearchResp struct {
	Hits struct {
		Hits []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Search executes a Lucene query string against index pattern "*" unless overridden later.
func (c *ElasticsearchClient) Search(ctx context.Context, query string, start, end time.Time, limit int) ([]LogLine, error) {
	if limit <= 0 {
		limit = 100
	}
	payload := map[string]interface{}{
		"size": limit,
		"sort": []interface{}{
			map[string]string{"@timestamp": "desc"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"query_string": map[string]interface{}{
							"query": query,
						},
					},
				},
				"filter": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte": start.UTC().Format(time.RFC3339),
								"lte": end.UTC().Format(time.RFC3339),
							},
						},
					},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	u := strings.TrimSuffix(c.BaseURL, "/") + "/_search"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	cli := &http.Client{Timeout: esHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch request: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("elasticsearch status %d: %s", resp.StatusCode, truncate(string(raw), 200))
	}
	var wrap esSearchResp
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, fmt.Errorf("decode es: %w", err)
	}
	out := make([]LogLine, 0, len(wrap.Hits.Hits))
	for _, h := range wrap.Hits.Hits {
		msg := ""
		ts := ""
		if m, ok := h.Source["message"].(string); ok {
			msg = m
		}
		if m, ok := h.Source["log"].(string); ok && msg == "" {
			msg = m
		}
		if t, ok := h.Source["@timestamp"].(string); ok {
			ts = t
		}
		labs := map[string]string{}
		for k, v := range h.Source {
			if k == "message" || k == "@timestamp" || k == "log" {
				continue
			}
			labs[k] = fmt.Sprint(v)
		}
		out = append(out, LogLine{Timestamp: ts, Message: msg, Labels: labs})
	}
	return out, nil
}

// NewElasticsearchClientFromConfig parses BaseConfig.
func NewElasticsearchClientFromConfig(cfgJSON []byte) (*ElasticsearchClient, error) {
	b, err := provider.ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &ElasticsearchClient{BaseURL: b.BaseURL, Token: b.Token}, nil
}
