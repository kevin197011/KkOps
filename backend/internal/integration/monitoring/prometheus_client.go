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
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kkops/backend/internal/integration/provider"
)

const promHTTPTimeout = 30 * time.Second

type promAPIResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
	Error string `json:"error"`
}

// PrometheusClient queries a Prometheus HTTP API.
type PrometheusClient struct {
	BaseURL string
	Token   string
}

// Query runs an instant query.
func (c *PrometheusClient) Query(ctx context.Context, query string, at time.Time) (*QueryResult, error) {
	q := url.Values{}
	q.Set("query", query)
	if !at.IsZero() {
		q.Set("time", strconv.FormatInt(at.Unix(), 10))
	}
	return c.do(ctx, "/api/v1/query?"+q.Encode())
}

// QueryRange runs a range query.
func (c *PrometheusClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*QueryResult, error) {
	q := url.Values{}
	q.Set("query", query)
	q.Set("start", formatRFC3339OrUnix(start))
	q.Set("end", formatRFC3339OrUnix(end))
	if step > 0 {
		q.Set("step", strconv.FormatFloat(step.Seconds(), 'f', -1, 64))
	}
	return c.do(ctx, "/api/v1/query_range?"+q.Encode())
}

func formatRFC3339OrUnix(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.Unix(), 10)
}

func (c *PrometheusClient) do(ctx context.Context, pathAndQuery string) (*QueryResult, error) {
	u := strings.TrimSuffix(c.BaseURL, "/") + pathAndQuery
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	cli := &http.Client{Timeout: promHTTPTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prometheus request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("prometheus status %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var wrap promAPIResponse
	if err := json.Unmarshal(body, &wrap); err != nil {
		return nil, fmt.Errorf("decode prometheus response: %w", err)
	}
	if wrap.Status != "success" {
		return nil, fmt.Errorf("prometheus error: %s", wrap.Error)
	}
	return normalizeProm(wrap.Data.ResultType, wrap.Data.Result), nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func normalizeProm(resultType string, raw []struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
	Values [][]interface{}   `json:"values"`
}) *QueryResult {
	out := &QueryResult{ResultType: resultType}
	for _, r := range raw {
		labels := r.Metric
		if labels == nil {
			labels = map[string]string{}
		}
		ms := MetricSeries{Labels: labels}
		if len(r.Values) > 0 {
			for _, pair := range r.Values {
				appendPoint(&ms, pair)
			}
		} else if len(r.Value) >= 2 {
			appendPoint(&ms, r.Value)
		}
		out.Series = append(out.Series, ms)
	}
	return out
}

func appendPoint(ms *MetricSeries, pair []interface{}) {
	if len(pair) < 2 {
		return
	}
	tsf, ok := numFloat(pair[0])
	if !ok {
		return
	}
	vf, ok := numFloat(pair[1])
	if !ok {
		return
	}
	ms.Points = append(ms.Points, MetricPoint{T: int64(tsf), V: vf})
}

func numFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

// NewPrometheusClientFromConfig parses provider.BaseConfig JSON.
func NewPrometheusClientFromConfig(cfgJSON []byte) (*PrometheusClient, error) {
	c, err := provider.ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &PrometheusClient{BaseURL: c.BaseURL, Token: c.Token}, nil
}
