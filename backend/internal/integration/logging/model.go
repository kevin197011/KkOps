// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package logging implements log search against Loki and Elasticsearch backends.
package logging

// LogLine is a normalized log record for the UI.
type LogLine struct {
	Timestamp string            `json:"timestamp"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels,omitempty"`
}
