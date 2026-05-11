// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type openAICompat struct {
	cfg *IntegrationConfig
	hc  *http.Client
}

func newOpenAICompat(cfg *IntegrationConfig) (Provider, error) {
	base := strings.TrimSuffix(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	return &openAICompat{
		cfg: cfg,
		hc:  &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func (p *openAICompat) Name() string {
	return "openai_compat"
}

func (p *openAICompat) Chat(ctx context.Context, messages []Message, opts ChatOptions) (ChatResponse, error) {
	model := opts.Model
	if model == "" {
		model = strings.TrimSpace(p.cfg.Model)
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	body := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": opts.Temperature,
	}
	if opts.MaxTokens > 0 {
		body["max_tokens"] = opts.MaxTokens
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return ChatResponse{}, err
	}
	base := strings.TrimSuffix(strings.TrimSpace(p.cfg.BaseURL), "/")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return ChatResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
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
		return ChatResponse{}, fmt.Errorf("%w: openai_compat http %d: %s", ErrProviderUnavailable, resp.StatusCode, truncateForErr(b))
	}
	var wrap struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(b, &wrap); err != nil {
		return ChatResponse{}, err
	}
	if wrap.Error != nil && wrap.Error.Message != "" {
		return ChatResponse{}, fmt.Errorf("%w: %s", ErrProviderUnavailable, wrap.Error.Message)
	}
	if len(wrap.Choices) == 0 {
		return ChatResponse{}, fmt.Errorf("%w: empty choices", ErrProviderUnavailable)
	}
	return ChatResponse{Content: wrap.Choices[0].Message.Content}, nil
}

func (p *openAICompat) ChatStream(ctx context.Context, messages []Message, opts ChatOptions, sink func(string) error) error {
	model := opts.Model
	if model == "" {
		model = strings.TrimSpace(p.cfg.Model)
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	body := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"stream":      true,
		"temperature": opts.Temperature,
	}
	if opts.MaxTokens > 0 {
		body["max_tokens"] = opts.MaxTokens
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	base := strings.TrimSuffix(strings.TrimSpace(p.cfg.BaseURL), "/")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if p.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	}
	resp, err := p.hc.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: stream http %d: %s", ErrProviderUnavailable, resp.StatusCode, truncateForErr(b))
	}
	sc := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if strings.TrimSpace(payload) == "[DONE]" {
			break
		}
		var ev struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(payload), &ev); err != nil {
			continue
		}
		if len(ev.Choices) == 0 {
			continue
		}
		t := ev.Choices[0].Delta.Content
		if t == "" {
			continue
		}
		if err := sink(t); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (p *openAICompat) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	model := strings.TrimSpace(p.cfg.Model)
	if model == "" {
		model = "text-embedding-3-small"
	}
	body := map[string]interface{}{
		"model": model,
		"input": texts,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	base := strings.TrimSuffix(strings.TrimSpace(p.cfg.BaseURL), "/")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/embeddings", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	}
	resp, err := p.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: embeddings http %d: %s", ErrProviderUnavailable, resp.StatusCode, truncateForErr(b))
	}
	var wrap struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &wrap); err != nil {
		return nil, err
	}
	out := make([][]float32, 0, len(wrap.Data))
	for _, d := range wrap.Data {
		row := make([]float32, len(d.Embedding))
		for i, v := range d.Embedding {
			row[i] = float32(v)
		}
		out = append(out, row)
	}
	return out, nil
}

func truncateForErr(b []byte) string {
	s := string(b)
	if len(s) > 512 {
		return s[:512] + "..."
	}
	return s
}
