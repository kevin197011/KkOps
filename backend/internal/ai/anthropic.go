// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type anthropicProv struct {
	cfg *IntegrationConfig
	hc  *http.Client
}

func newAnthropic(cfg *IntegrationConfig) (Provider, error) {
	_ = strings.TrimSpace(cfg.BaseURL)
	return &anthropicProv{
		cfg: cfg,
		hc:  &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func (p *anthropicProv) Name() string {
	return "anthropic"
}

func (p *anthropicProv) Chat(ctx context.Context, messages []Message, opts ChatOptions) (ChatResponse, error) {
	model := opts.Model
	if model == "" {
		model = strings.TrimSpace(p.cfg.Model)
	}
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}
	var sys string
	var msgs []map[string]string
	for _, m := range messages {
		if m.Role == "system" {
			sys = m.Content
			continue
		}
		r := m.Role
		if r == "assistant" || r == "user" {
			msgs = append(msgs, map[string]string{"role": r, "content": m.Content})
		}
	}
	body := map[string]interface{}{
		"model":      model,
		"messages":   msgs,
		"max_tokens": 4096,
	}
	if sys != "" {
		body["system"] = sys
	}
	if opts.MaxTokens > 0 {
		body["max_tokens"] = opts.MaxTokens
	}
	if opts.Temperature > 0 {
		body["temperature"] = opts.Temperature
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return ChatResponse{}, err
	}
	base := strings.TrimSpace(p.cfg.BaseURL)
	if base == "" {
		base = "https://api.anthropic.com"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/messages", bytes.NewReader(raw))
	if err != nil {
		return ChatResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.cfg.APIKey)
	ver := strings.TrimSpace(p.cfg.AnthropicVersion)
	if ver == "" {
		ver = "2023-06-01"
	}
	req.Header.Set("anthropic-version", ver)
	resp, err := p.hc.Do(req)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ChatResponse{}, fmt.Errorf("%w: anthropic http %d: %s", ErrProviderUnavailable, resp.StatusCode, truncateForErr(b))
	}
	var wrap struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(b, &wrap); err != nil {
		return ChatResponse{}, err
	}
	var sb strings.Builder
	for _, c := range wrap.Content {
		sb.WriteString(c.Text)
	}
	return ChatResponse{Content: sb.String()}, nil
}

func (p *anthropicProv) ChatStream(ctx context.Context, messages []Message, opts ChatOptions, sink func(string) error) error {
	// Anthropic streaming uses SSE with event types; delegate to non-stream then chunk for simplicity.
	cr, err := p.Chat(ctx, messages, opts)
	if err != nil {
		return err
	}
	const chunkSize = 24
	s := cr.Content
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		if err := sink(s[i:end]); err != nil {
			return err
		}
	}
	return nil
}

func (p *anthropicProv) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, fmt.Errorf("%w: embeddings not configured for anthropic", ErrProviderUnavailable)
}
