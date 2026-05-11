// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package aisvc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/ai"
)

var toolMarkerRe = regexp.MustCompile(`<<TOOL:\s*(\w+)\s+(\{[\s\S]*?\})\s*>>`)

const systemToolPrompt = `You are KkOps AI assistant. You may call read-only tools by emitting exactly one marker line when you need data:
<<TOOL: tool_name {"arg":"value"}>>
Tools: list_alerts (optional status), get_incident (id), query_metric (integration_id, query), search_logs (integration_id, query, optional range like 1h), list_pods (cluster_id, optional namespace), pipeline_status (app, optional argocd_integration_id), list_integrations ({}).
After results appear in the conversation, answer clearly. Do not invent metric values.`

// ChatTurn resolves tool markers with non-streaming Chat, then streams the final reply as SSE chunks from text (single completion, no second LLM call).
// Returns the assistant text after streaming completes (for persistence).
func (s *Service) ChatTurn(ctx context.Context, reg *ai.Registry, integrationID *uint, messages []ai.Message, c *gin.Context) (string, error) {
	p, _, err := s.pickProvider(reg, integrationID)
	if err != nil {
		return "", err
	}

	work := append([]ai.Message{{Role: "system", Content: systemToolPrompt}}, messages...)
	const maxRounds = 8
	var finalText string
	for round := 0; round < maxRounds; round++ {
		resp, err := p.Chat(ctx, work, ai.ChatOptions{Temperature: 0.2})
		if err != nil {
			return "", err
		}
		txt := resp.Content
		if !toolMarkerRe.MatchString(txt) {
			finalText = txt
			break
		}
		m := toolMarkerRe.FindStringSubmatch(txt)
		if len(m) < 3 {
			finalText = txt
			break
		}
		toolName := m[1]
		args := m[2]
		result, err := s.Tools.Execute(ctx, toolName, args)
		if err != nil {
			result = fmt.Sprintf(`{"error":%q}`, err.Error())
		}
		work = append(work,
			ai.Message{Role: "assistant", Content: txt},
			ai.Message{Role: "user", Content: fmt.Sprintf("TOOL %s RESULT:\n%s", toolName, result)},
		)
	}
	if finalText == "" {
		return "", fmt.Errorf("tool loop exceeded max rounds")
	}
	if err := streamSSEChunks(c, finalText); err != nil {
		return "", err
	}
	return finalText, nil
}

func (s *Service) pickProvider(reg *ai.Registry, integrationID *uint) (ai.Provider, uint, error) {
	if integrationID != nil && *integrationID > 0 {
		p, err := reg.ProviderForIntegration(*integrationID)
		if err != nil {
			return nil, 0, err
		}
		return p, *integrationID, nil
	}
	return reg.DefaultProvider()
}

func streamSSEChunks(c *gin.Context, full string) error {
	h := c.Writer.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")
	fl, ok := c.Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}
	const chunkSize = 32
	for i := 0; i < len(full); i += chunkSize {
		end := i + chunkSize
		if end > len(full) {
			end = len(full)
		}
		chunk := full[i:end]
		line, err := json.Marshal(map[string]string{"token": chunk})
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", line); err != nil {
			return err
		}
		fl.Flush()
	}
	_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", `{"done":true}`)
	fl.Flush()
	return nil
}

// StreamPlainSSE streams provider-native chunks (for ping/test).
func StreamPlainSSE(c *gin.Context, p ai.Provider, messages []ai.Message) error {
	h := c.Writer.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	fl, ok := c.Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}
	err := p.ChatStream(c.Request.Context(), messages, ai.ChatOptions{Temperature: 0}, func(chunk string) error {
		line, jerr := json.Marshal(map[string]string{"token": chunk})
		if jerr != nil {
			return jerr
		}
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", line); err != nil {
			return err
		}
		fl.Flush()
		return nil
	})
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", `{"done":true}`)
	fl.Flush()
	return nil
}

// DrainChatNonStream runs chat with tool resolution until plain text (for RCA).
func (s *Service) DrainChatNonStream(ctx context.Context, reg *ai.Registry, integrationID *uint, messages []ai.Message) (string, error) {
	p, _, err := s.pickProvider(reg, integrationID)
	if err != nil {
		return "", err
	}
	work := append([]ai.Message{{Role: "system", Content: systemToolPrompt}}, messages...)
	for round := 0; round < 8; round++ {
		resp, err := p.Chat(ctx, work, ai.ChatOptions{Temperature: 0.2, MaxTokens: 8192})
		if err != nil {
			return "", err
		}
		txt := resp.Content
		if !toolMarkerRe.MatchString(txt) {
			return txt, nil
		}
		m := toolMarkerRe.FindStringSubmatch(txt)
		if len(m) < 3 {
			return txt, nil
		}
		result, err := s.Tools.Execute(ctx, m[1], m[2])
		if err != nil {
			result = fmt.Sprintf(`{"error":%q}`, err.Error())
		}
		work = append(work,
			ai.Message{Role: "assistant", Content: txt},
			ai.Message{Role: "user", Content: fmt.Sprintf("TOOL %s RESULT:\n%s", m[1], result)},
		)
	}
	return "", fmt.Errorf("tool loop exceeded max rounds")
}

// DrainChatNoTools forces a single completion without tool parsing (prompt already constrained).
func (s *Service) DrainChatNoTools(ctx context.Context, reg *ai.Registry, integrationID *uint, messages []ai.Message) (string, error) {
	p, _, err := s.pickProvider(reg, integrationID)
	if err != nil {
		return "", err
	}
	resp, err := p.Chat(ctx, messages, ai.ChatOptions{Temperature: 0.2, MaxTokens: 8192})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Content), nil
}
