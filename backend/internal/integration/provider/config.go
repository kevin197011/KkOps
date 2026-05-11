// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provider

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BaseConfig is the common JSON shape for Phase 2 stubs (base URL + bearer token).
type BaseConfig struct {
	BaseURL             string `json:"base_url"`
	Token               string `json:"token"`
	AlertmanagerBaseURL string `json:"alertmanager_base_url,omitempty"` // optional; Alertmanager UI/API base for alert sync
}

// ParseBase unmarshals config and normalizes base URL.
func ParseBase(cfgJSON []byte) (BaseConfig, error) {
	var c BaseConfig
	if len(cfgJSON) == 0 {
		return c, fmt.Errorf("empty integration config")
	}
	if err := json.Unmarshal(cfgJSON, &c); err != nil {
		return c, err
	}
	c.BaseURL = strings.TrimSuffix(strings.TrimSpace(c.BaseURL), "/")
	if c.BaseURL == "" {
		return c, fmt.Errorf("base_url is required")
	}
	c.AlertmanagerBaseURL = strings.TrimSuffix(strings.TrimSpace(c.AlertmanagerBaseURL), "/")
	return c, nil
}
