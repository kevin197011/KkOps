// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cicd

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Config extends base integration JSON with CI-specific fields.
type Config struct {
	BaseURL   string `json:"base_url"`
	Token     string `json:"token"`
	ProjectID int    `json:"project_id"` // GitLab numeric project id
	JobPath   string `json:"job_path"`   // Jenkins job path segments (e.g. folder/job/my-job)
}

func parseConfig(cfgJSON []byte) (Config, error) {
	var c Config
	if len(cfgJSON) == 0 {
		return c, fmt.Errorf("empty config")
	}
	if err := json.Unmarshal(cfgJSON, &c); err != nil {
		return c, err
	}
	c.BaseURL = strings.TrimSuffix(strings.TrimSpace(c.BaseURL), "/")
	if c.BaseURL == "" {
		return c, fmt.Errorf("base_url is required")
	}
	return c, nil
}
