// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ProviderKind selects implementation for an integration row with kind "ai".
type ProviderKind string

const (
	KindOpenAICompat ProviderKind = "openai_compat"
	KindAnthropic    ProviderKind = "anthropic"
	KindWebhook      ProviderKind = "webhook"
)

// IntegrationConfig is decrypted JSON stored on integrations.kind = ai.
type IntegrationConfig struct {
	Provider ProviderKind `json:"provider"`
	BaseURL  string       `json:"base_url"`
	APIKey   string       `json:"api_key"`
	Model    string       `json:"model"`

	AnthropicVersion string `json:"anthropic_version"`

	WebhookURL     string            `json:"webhook_url"`
	WebhookHeaders map[string]string `json:"webhook_headers"`

	Extra map[string]interface{} `json:"extra,omitempty"`
}

// ParseIntegrationConfig unmarshals encrypted integration JSON.
func ParseIntegrationConfig(raw []byte) (*IntegrationConfig, error) {
	var c IntegrationConfig
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("parse ai integration config: %w", err)
	}
	c.Provider = ProviderKind(strings.TrimSpace(string(c.Provider)))
	if c.Provider == "" {
		return nil, fmt.Errorf("ai integration config: provider is required")
	}
	return &c, nil
}

// NewProviderFromConfig builds a Provider from decrypted integration bytes.
func NewProviderFromConfig(raw []byte) (Provider, error) {
	cfg, err := ParseIntegrationConfig(raw)
	if err != nil {
		return nil, err
	}
	switch cfg.Provider {
	case KindOpenAICompat:
		return newOpenAICompat(cfg)
	case KindAnthropic:
		return newAnthropic(cfg)
	case KindWebhook:
		return newWebhook(cfg)
	default:
		return nil, fmt.Errorf("unsupported ai provider kind %q", cfg.Provider)
	}
}
