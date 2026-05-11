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

type webhookProv struct {
	cfg *IntegrationConfig
	hc  *http.Client
}

func newWebhook(cfg *IntegrationConfig) (Provider, error) {
	if strings.TrimSpace(cfg.WebhookURL) == "" {
		return nil, fmt.Errorf("webhook_url is required")
	}
	return &webhookProv{
		cfg: cfg,
		hc:  &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func (p *webhookProv) Name() string {
	return "webhook"
}

func (p *webhookProv) Chat(ctx context.Context, messages []Message, opts ChatOptions) (ChatResponse, error) {
	body := map[string]interface{}{
		"messages": messages,
		"options":  opts,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return ChatResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.WebhookURL, bytes.NewReader(raw))
	if err != nil {
		return ChatResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range p.cfg.WebhookHeaders {
		req.Header.Set(k, v)
	}
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
		return ChatResponse{}, fmt.Errorf("%w: webhook http %d: %s", ErrProviderUnavailable, resp.StatusCode, truncateForErr(b))
	}
	text := strings.TrimSpace(string(b))
	var wrap struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(b, &wrap); err == nil && wrap.Text != "" {
		text = wrap.Text
	}
	return ChatResponse{Content: text}, nil
}

func (p *webhookProv) ChatStream(ctx context.Context, messages []Message, opts ChatOptions, sink func(string) error) error {
	cr, err := p.Chat(ctx, messages, opts)
	if err != nil {
		return err
	}
	const chunkSize = 16
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

func (p *webhookProv) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, fmt.Errorf("%w: embeddings not supported for webhook", ErrProviderUnavailable)
}
