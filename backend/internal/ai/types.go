// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package ai defines LLM provider interfaces and HTTP-backed implementations.
package ai

import (
	"context"
	"errors"
)

// Message is one chat turn for the LLM API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatOptions tunes completion behavior.
type ChatOptions struct {
	Model       string
	Temperature float64
	MaxTokens   int
}

// ChatResponse is a non-streaming completion result.
type ChatResponse struct {
	Content string
}

// Provider is implemented by concrete LLM backends.
type Provider interface {
	Name() string
	Chat(ctx context.Context, messages []Message, opts ChatOptions) (ChatResponse, error)
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	// ChatStream streams completion deltas via sink; optional for webhook-only backends.
	ChatStream(ctx context.Context, messages []Message, opts ChatOptions, sink func(string) error) error
}

// ErrProviderUnavailable indicates configuration or transport failure.
var ErrProviderUnavailable = errors.New("ai provider unavailable")
