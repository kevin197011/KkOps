// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package monitoring implements Prometheus-compatible metric query helpers for integrations.
package monitoring

import "time"

// MetricPoint is a single timestamped sample.
type MetricPoint struct {
	T int64   `json:"t"`
	V float64 `json:"v"`
}

// MetricSeries is a normalized time series for UI charts.
type MetricSeries struct {
	Labels map[string]string `json:"labels"`
	Points []MetricPoint     `json:"points"`
}

// QueryResult wraps instant or range query output.
type QueryResult struct {
	ResultType string         `json:"result_type"`
	Series     []MetricSeries `json:"series"`
}

// RangeSpec defines a Prometheus range query window.
type RangeSpec struct {
	Start time.Time
	End   time.Time
	Step  time.Duration
}
