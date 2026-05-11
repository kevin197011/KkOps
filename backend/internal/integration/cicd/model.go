// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package cicd integrates Jenkins and GitLab CI for pipeline visibility.
package cicd

import "time"

// Pipeline is a normalized CI/CD run descriptor for the UI.
type Pipeline struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	IntegrationID uint       `json:"integration_id"`
	WebURL        string     `json:"web_url,omitempty"`
	Ref           string     `json:"ref,omitempty"`
	RecordedAt    *time.Time `json:"recorded_at,omitempty"`
}
